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
	"github.com/ijry/lyshop/core/db"
	inventorycore "github.com/ijry/lyshop/core/inventory"
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

func setupExternalAdminRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	gdb, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:external_wms_admin_%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	require.NoError(t, err)
	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })
	require.NoError(t, gdb.AutoMigrate(&inventorycore.InventoryIntegrationTask{}))

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
