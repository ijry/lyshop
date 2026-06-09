package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/config"
	"github.com/ijry/lyshop/core/db"
	inventorycore "github.com/ijry/lyshop/core/inventory"
	ordermodel "github.com/ijry/lyshop/plugins/order/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type apiResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func TestExternalAdminRoutes(t *testing.T) {
	router, gdb := setupExternalAdminRouter(t)
	task := inventorycore.NewIntegrationTask("external_wms", "order", "ORD-11", "reserve", inventorycore.TaskPayload{})
	require.NoError(t, gdb.Create(task).Error)

	resp := doExternalAdminReq(t, router, http.MethodGet, "/admin/inventory/tasks", "", "*")
	require.Equal(t, 0, resp.Code)
	require.Contains(t, string(resp.Data), "ORD-11")

	retryResp := doExternalAdminReq(t, router, http.MethodPost, fmt.Sprintf("/admin/inventory/tasks/%d/retry", task.ID), "", "*")
	require.Equal(t, 0, retryResp.Code)
}

func TestExternalWMSCallbackIgnoresDuplicateCallbackID(t *testing.T) {
	router, gdb := setupExternalAdminRouter(t)
	now := time.Now()
	task := inventorycore.InventoryIntegrationTask{
		Provider:       "external_wms",
		Action:         "reserve",
		BizType:        "order",
		BizNo:          "O1001",
		Status:         inventorycore.TaskStatusProcessing,
		RequestID:      "REQ-1",
		LastCallbackID: "CALLBACK-1",
		NextRetryAt:    &now,
	}
	require.NoError(t, gdb.Create(&task).Error)
	require.NoError(t, gdb.Create(&ordermodel.Order{
		OrderNo:         "O1001",
		UserID:          1,
		Status:          ordermodel.OrderStatusPending,
		InventoryStatus: inventorycore.InventoryStatusPending,
	}).Error)

	resp := doExternalAdminReq(t, router, http.MethodPost, "/admin/external-wms/callback", `{"request_id":"REQ-1","callback_id":"CALLBACK-1","status":"success"}`, "*")
	require.Equal(t, 0, resp.Code)

	var latest inventorycore.InventoryIntegrationTask
	require.NoError(t, gdb.First(&latest, task.ID).Error)
	require.Equal(t, inventorycore.TaskStatusProcessing, latest.Status)
	requireOrderInventoryStatus(t, gdb, "O1001", inventorycore.InventoryStatusPending)
}

func TestExternalWMSCallbackMarksOrderInventoryFailed(t *testing.T) {
	router, gdb := setupExternalAdminRouter(t)
	now := time.Now()
	task := inventorycore.InventoryIntegrationTask{
		Provider:    "external_wms",
		Action:      "deduct",
		BizType:     "order",
		BizNo:       "O1002",
		Status:      inventorycore.TaskStatusProcessing,
		RequestID:   "REQ-2",
		NextRetryAt: &now,
	}
	require.NoError(t, gdb.Create(&task).Error)
	require.NoError(t, gdb.Create(&ordermodel.Order{
		OrderNo:         "O1002",
		UserID:          1,
		Status:          ordermodel.OrderStatusPending,
		InventoryStatus: inventorycore.InventoryStatusPending,
	}).Error)

	resp := doExternalAdminReq(t, router, http.MethodPost, "/admin/external-wms/callback", `{"request_id":"REQ-2","callback_id":"CALLBACK-2","status":"failed","message":"stock not enough"}`, "*")
	require.Equal(t, 0, resp.Code)

	var latest inventorycore.InventoryIntegrationTask
	require.NoError(t, gdb.First(&latest, task.ID).Error)
	require.Equal(t, inventorycore.TaskStatusFailed, latest.Status)
	require.Equal(t, "CALLBACK-2", latest.LastCallbackID)
	require.Equal(t, "stock not enough", latest.LastError)
	requireOrderInventoryStatus(t, gdb, "O1002", inventorycore.InventoryStatusFailed)
}

func TestExternalWMSCallbackMarksOrderInventoryConfirmedOnSuccess(t *testing.T) {
	router, gdb := setupExternalAdminRouter(t)
	now := time.Now()
	task := inventorycore.InventoryIntegrationTask{
		Provider:    "external_wms",
		Action:      "deduct",
		BizType:     "order",
		BizNo:       "O1003",
		Status:      inventorycore.TaskStatusProcessing,
		RequestID:   "REQ-3",
		NextRetryAt: &now,
	}
	require.NoError(t, gdb.Create(&task).Error)
	require.NoError(t, gdb.Create(&ordermodel.Order{
		OrderNo:         "O1003",
		UserID:          1,
		Status:          ordermodel.OrderStatusPaid,
		InventoryStatus: inventorycore.InventoryStatusPending,
	}).Error)

	resp := doExternalAdminReq(t, router, http.MethodPost, "/admin/external-wms/callback", `{"request_id":"REQ-3","callback_id":"CALLBACK-3","status":"success"}`, "*")
	require.Equal(t, 0, resp.Code)

	var latest inventorycore.InventoryIntegrationTask
	require.NoError(t, gdb.First(&latest, task.ID).Error)
	require.Equal(t, inventorycore.TaskStatusSuccess, latest.Status)
	requireOrderInventoryStatus(t, gdb, "O1003", inventorycore.InventoryStatusConfirmed)
}

func TestExternalWMSCallbackRejectsInvalidSignature(t *testing.T) {
	original := config.Global
	t.Cleanup(func() { config.Global = original })
	config.Global.ExternalWMS.AppSecret = "demo-secret"
	config.Global.ExternalWMS.SignatureTTL = 300

	router, _ := setupExternalAdminRouter(t)
	resp := doExternalAdminReq(t, router, http.MethodPost, "/admin/external-wms/callback", `{"app_key":"demo-key","timestamp":"1717910400","nonce":"nonce-1","sign":"bad-sign","body":"{\"request_id\":\"REQ-1\"}","request_id":"REQ-1","callback_id":"CALLBACK-3","status":"success"}`, "*")
	require.NotEqual(t, 0, resp.Code)
}

func setupExternalAdminRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	gdb, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:external_wms_admin_%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	require.NoError(t, err)
	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })
	require.NoError(t, gdb.AutoMigrate(&inventorycore.InventoryIntegrationTask{}, &ordermodel.Order{}))

	r := gin.New()
	admin := r.Group("/admin", func(c *gin.Context) {
		raw := strings.TrimSpace(c.GetHeader("X-Test-Perms"))
		if raw == "" {
			c.Set("perms", []string{"*"})
		} else {
			c.Set("perms", strings.Split(raw, ","))
		}
		c.Set("role", "admin")
		c.Next()
	})
	RegisterAdminRoutes(admin)
	return r, gdb
}

func doExternalAdminReq(t *testing.T, router *gin.Engine, method, path, body, perms string) apiResp {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-Perms", perms)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var out apiResp
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &out))
	return out
}

func requireOrderInventoryStatus(t *testing.T, gdb *gorm.DB, orderNo string, want string) {
	t.Helper()
	var order ordermodel.Order
	require.NoError(t, gdb.Where("order_no = ?", orderNo).First(&order).Error)
	require.Equal(t, want, order.InventoryStatus)
}
