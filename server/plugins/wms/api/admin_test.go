package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/core/db"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	wmssvc "github.com/ijry/lyshop/plugins/wms/service"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type apiResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type pageResp struct {
	List  []map[string]any `json:"list"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

func TestRegisterAdminRoutesContainsRequiredEndpoints(t *testing.T) {
	router, _ := setupAdminTestRouter(t)
	routeSet := map[string]bool{}
	for _, route := range router.Routes() {
		routeSet[route.Method+" "+route.Path] = true
	}

	required := []string{
		"GET /admin/wms/warehouses",
		"POST /admin/wms/warehouses",
		"PUT /admin/wms/warehouses/:id",
		"PUT /admin/wms/warehouses/:id/status",
		"GET /admin/wms/stocks",
		"PUT /admin/wms/stocks/:id/safety",
		"GET /admin/wms/docs",
		"POST /admin/wms/docs",
		"GET /admin/wms/docs/:id",
		"PUT /admin/wms/docs/:id",
		"POST /admin/wms/docs/:id/complete",
		"POST /admin/wms/docs/:id/cancel",
		"GET /admin/wms/movements",
	}
	for _, key := range required {
		require.Truef(t, routeSet[key], "route missing: %s", key)
	}
}

func TestAdminWarehouseStatusDocCompleteAndMovements(t *testing.T) {
	router, testDB := setupAdminTestRouter(t)

	warehouseForStatus := wmsmodel.Warehouse{Code: "WH-S", Name: "状态仓", Status: wmsmodel.WarehouseStatusEnabled}
	require.NoError(t, testDB.Create(&warehouseForStatus).Error)

	statusResp := doJSONRequest(t, router, http.MethodPut, fmt.Sprintf("/admin/wms/warehouses/%d/status", warehouseForStatus.ID), `{"status":0}`)
	require.Equal(t, 0, statusResp.Code)

	var latestStatusWarehouse wmsmodel.Warehouse
	require.NoError(t, testDB.Where("id = ?", warehouseForStatus.ID).First(&latestStatusWarehouse).Error)
	require.Equal(t, wmsmodel.WarehouseStatusDisabled, latestStatusWarehouse.Status)

	warehouse := wmsmodel.Warehouse{Code: "WH-C", Name: "完成仓", Status: wmsmodel.WarehouseStatusEnabled}
	require.NoError(t, testDB.Create(&warehouse).Error)
	require.NoError(t, testDB.Create(&wmsmodel.InventoryStock{
		WarehouseID: warehouse.ID,
		SkuID:       301,
		Qty:         10,
		SafeQty:     1,
	}).Error)

	doc, err := wmssvc.CreateDraftDoc(context.Background(), wmssvc.CreateDocInput{
		WarehouseID: warehouse.ID,
		DocType:     wmsmodel.DocTypeOutbound,
		Remark:      "API 完成测试",
		Items: []wmssvc.DocItemInput{
			{SkuID: 301, Qty: 4},
		},
	})
	require.NoError(t, err)

	completeResp := doJSONRequest(t, router, http.MethodPost, fmt.Sprintf("/admin/wms/docs/%d/complete", doc.ID), "")
	require.Equal(t, 0, completeResp.Code)

	var latestDoc wmsmodel.InventoryDoc
	require.NoError(t, testDB.Where("id = ?", doc.ID).First(&latestDoc).Error)
	require.Equal(t, wmsmodel.DocStatusCompleted, latestDoc.Status)

	movementResp := doJSONRequest(t, router, http.MethodGet, "/admin/wms/movements?doc_no="+doc.DocNo, "")
	require.Equal(t, 0, movementResp.Code)

	var movementPage pageResp
	require.NoError(t, json.Unmarshal(movementResp.Data, &movementPage))
	require.GreaterOrEqual(t, movementPage.Total, int64(1))
	require.NotEmpty(t, movementPage.List)
	require.Equal(t, doc.DocNo, movementPage.List[0]["doc_no"])
}

func doJSONRequest(t *testing.T, router *gin.Engine, method, path, body string) apiResp {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var r apiResp
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &r))
	return r
}

func setupAdminTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	dsn := fmt.Sprintf("file:wms_api_%d?mode=memory&cache=shared", time.Now().UnixNano())
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })

	require.NoError(t, gdb.AutoMigrate(
		&wmsmodel.Warehouse{},
		&wmsmodel.InventoryStock{},
		&wmsmodel.InventoryMovement{},
		&wmsmodel.InventoryDoc{},
		&wmsmodel.InventoryDocItem{},
	))

	r := gin.New()
	admin := r.Group("/admin", func(c *gin.Context) {
		c.Set("perms", []string{"*"})
		c.Next()
	})
	RegisterAdminRoutes(admin)
	return r, gdb
}
