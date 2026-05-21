package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	apiauth "github.com/ijry/lyshop/api/auth"
	"github.com/ijry/lyshop/config"
	"github.com/ijry/lyshop/core/cache"
	"github.com/ijry/lyshop/core/db"
	driverStorage "github.com/ijry/lyshop/core/driver/storage"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/core/response"
	imapi "github.com/ijry/lyshop/plugins/im/api"
)

// Init loads config then initializes DB and Redis.
func Init(cfgPath string) error {
	if err := config.Load(cfgPath); err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	gin.SetMode(config.Global.Server.Mode)

	if err := db.Init(); err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	if err := cache.Init(); err != nil {
		return fmt.Errorf("init cache: %w", err)
	}
	return nil
}

// Run builds the Gin engine, loads plugins, and starts the HTTP server.
func Run() error {
	r := gin.New()
	r.Use(middleware.Logger(), middleware.CORS(), gin.Recovery())

	// Public route groups (no auth)
	front := r.Group("/api/v1")
	adminPublic := r.Group("/admin/api")

	// Auth-protected admin group
	adminAuth := r.Group("/admin/api")
	adminAuth.Use(middleware.RequireAdmin)

	// Core auth routes
	apiauth.RegisterFrontRoutes(front)
	apiauth.RegisterAdminRoutes(adminPublic)

	// Dynamic menus endpoint — returns menus for all enabled plugins
	adminAuth.GET("/menus", func(c *gin.Context) {
		menus := plugin.EnabledMenus(config.Global.Plugins.Enabled)
		c.JSON(200, menus)
	})

	// Universal file upload endpoint
	adminAuth.POST("/upload", func(c *gin.Context) {
		fh, err := c.FormFile("file")
		if err != nil {
			response.Fail(c, 400, "请选择文件")
			return
		}
		driver, err := driverStorage.Get()
		if err != nil {
			response.Fail(c, 500, err.Error())
			return
		}
		result, err := driver.Upload(c.Request.Context(), fh)
		if err != nil {
			response.Fail(c, 500, err.Error())
			return
		}
		response.OK(c, result)
	})
	// Ping
	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"pong": true}) })

	// Load enabled plugins
	if err := plugin.Load(config.Global.Plugins.Enabled, db.DB, front, adminAuth); err != nil {
		return fmt.Errorf("load plugins: %w", err)
	}

	// WebSocket IM endpoint (registered after plugins so Hub is running)
	imapi.RegisterWSRoute(r)

	addr := fmt.Sprintf(":%d", config.Global.Server.Port)
	return r.Run(addr)
}
