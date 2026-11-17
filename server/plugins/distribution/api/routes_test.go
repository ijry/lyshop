package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/plugins/distribution"
)

func TestDistributionRoutesUseHostPrefixes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	front := router.Group("/api/v1")
	admin := router.Group("/admin/api", middleware.RequireAdmin)
	(&distribution.Plugin{}).RegisterRoutes(front, admin)

	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/admin/api/distribution/config"},
		{http.MethodPut, "/admin/api/distribution/config"},
		{http.MethodGet, "/admin/api/distribution/distributors"},
		{http.MethodGet, "/admin/api/distribution/distributors/:id"},
		{http.MethodPut, "/admin/api/distribution/distributors/:id"},
		{http.MethodGet, "/admin/api/distribution/orders"},
		{http.MethodGet, "/admin/api/distribution/orders/:id"},
		{http.MethodGet, "/admin/api/distribution/withdrawals"},
		{http.MethodGet, "/admin/api/distribution/withdrawals/:id"},
		{http.MethodPost, "/admin/api/distribution/withdrawals/:id/approve"},
		{http.MethodPost, "/admin/api/distribution/withdrawals/:id/reject"},
		{http.MethodPost, "/admin/api/distribution/withdrawals/:id/complete"},
		{http.MethodGet, "/api/v1/distribution/info"},
		{http.MethodPost, "/api/v1/distribution/apply"},
		{http.MethodGet, "/api/v1/distribution/team"},
		{http.MethodGet, "/api/v1/distribution/orders"},
		{http.MethodPost, "/api/v1/distribution/withdrawals"},
		{http.MethodGet, "/api/v1/distribution/withdrawals"},
		{http.MethodGet, "/api/v1/distribution/withdrawals/:id"},
	}
	registered := make(map[string]bool)
	for _, route := range router.Routes() {
		registered[route.Method+" "+route.Path] = true
	}
	if len(registered) != len(routes) {
		t.Fatalf("路由数量 = %d，期望 %d", len(registered), len(routes))
	}

	for _, route := range routes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			if !registered[route.method+" "+route.path] {
				t.Fatalf("缺少路由 %s %s", route.method, route.path)
			}
			path := strings.ReplaceAll(route.path, ":id", "1")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(route.method, path, nil))
			var response struct {
				Code int `json:"code"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("解析认证响应失败: %v，响应: %s", err, recorder.Body.String())
			}
			if recorder.Code != http.StatusOK || response.Code != 401 {
				t.Fatalf("未登录请求应被认证中间件拒绝，HTTP = %d，响应: %s", recorder.Code, recorder.Body.String())
			}

			prefix := "/api/v1"
			if strings.HasPrefix(path, "/admin/api/") {
				prefix = "/admin/api"
			}
			recorder = httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(route.method, prefix+path, nil))
			if recorder.Code != http.StatusNotFound {
				t.Fatalf("重复前缀路径 %s 不应注册，HTTP = %d", prefix+path, recorder.Code)
			}
		})
	}
}

func TestDistributionAdminRoutesRequirePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	admin := router.Group("/admin/api", func(c *gin.Context) {
		c.Set("perms", []string{"product:view"})
		c.Next()
	})
	(&distribution.Plugin{}).RegisterRoutes(router.Group("/api/v1"), admin)

	for _, route := range router.Routes() {
		if !strings.HasPrefix(route.Path, "/admin/api/") {
			continue
		}
		t.Run(route.Method+" "+route.Path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			path := strings.ReplaceAll(route.Path, ":id", "1")
			router.ServeHTTP(recorder, httptest.NewRequest(route.Method, path, nil))
			var response struct {
				Code int `json:"code"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("解析权限响应失败: %v，响应: %s", err, recorder.Body.String())
			}
			if recorder.Code != http.StatusOK || response.Code != 403 {
				t.Fatalf("缺少分销权限时应拒绝请求，HTTP = %d，响应: %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}
