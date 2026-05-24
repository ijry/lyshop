package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	apiadmin "github.com/ijry/lyshop/api/admin"
	apiauth "github.com/ijry/lyshop/api/auth"
	"github.com/ijry/lyshop/config"
	"github.com/ijry/lyshop/core/cache"
	"github.com/ijry/lyshop/core/db"
	driverStorage "github.com/ijry/lyshop/core/driver/storage"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/core/response"
	"github.com/ijry/lyshop/model"
	imapi "github.com/ijry/lyshop/plugins/im/api"
	adminsvc "github.com/ijry/lyshop/service/admin"
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

	// Auto-migrate core tables
	db.DB.AutoMigrate(&model.Admin{}, &model.Role{}, &model.ConfigKV{}, &model.User{})

	// Seed super admin on first run
	if err := adminsvc.EnsureSuperAdmin(); err != nil {
		slog.Warn("failed to seed super admin", "error", err)
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

	// Dynamic menus endpoint — filtered by admin's permissions
	adminAuth.GET("/menus", func(c *gin.Context) {
		perms, _ := c.Get("perms")
		permList, _ := perms.([]string)
		menus := plugin.EnabledMenus(config.Global.Plugins.Enabled, permList)
		c.JSON(200, menus)
	})
	adminAuth.GET("/dashboard", func(c *gin.Context) {
		data, err := adminsvc.GetDashboard(c.Request.Context())
		if err != nil {
			response.Fail(c, 500, err.Error())
			return
		}
		response.OK(c, data)
	})

	// Admin/Role management routes
	apiadmin.RegisterRoutes(adminAuth)

	// === Config Center API ===
	// GET /admin/api/config/schemas — list all plugins that have configItems
	adminAuth.GET("/config/schemas", middleware.RequirePermission("system:config"), func(c *gin.Context) {
		type pluginSchema struct {
			Plugin string               `json:"plugin"`
			Title  string               `json:"title"`
			Fields []plugin.ConfigField `json:"fields"`
		}
		var schemas []pluginSchema
		for _, name := range config.Global.Plugins.Enabled {
			p := plugin.Find(name)
			if p == nil {
				continue
			}
			meta := p.Meta()
			if len(meta.ConfigItems) > 0 {
				schemas = append(schemas, pluginSchema{Plugin: meta.Name, Title: meta.Title, Fields: meta.ConfigItems})
			}
		}
		response.OK(c, schemas)
	})
	// GET /admin/api/config/:plugin — get all config values for a plugin
	adminAuth.GET("/config/:plugin", middleware.RequirePermission("system:config"), func(c *gin.Context) {
		pluginName := c.Param("plugin")
		var kvs []model.ConfigKV
		db.DB.Where("plugin = ?", pluginName).Find(&kvs)
		result := map[string]string{}
		for _, kv := range kvs {
			result[kv.Key] = kv.Value
		}
		response.OK(c, result)
	})
	// PUT /admin/api/config/:plugin — save config values for a plugin
	adminAuth.PUT("/config/:plugin", middleware.RequirePermission("system:config"), func(c *gin.Context) {
		pluginName := c.Param("plugin")
		var values map[string]string
		if err := c.ShouldBindJSON(&values); err != nil {
			response.Fail(c, 400, err.Error())
			return
		}
		for k, v := range values {
			db.DB.Where(model.ConfigKV{Plugin: pluginName, Key: k}).
				Assign(model.ConfigKV{Value: v}).
				FirstOrCreate(&model.ConfigKV{})
		}
		response.OK(c, nil)
	})

	// Universal file upload endpoint
	adminAuth.POST("/upload", func(c *gin.Context) {
		fh, err := c.FormFile("file")
		if err != nil {
			response.Fail(c, 400, "请选择文件")
			return
		}
		driverName := strings.TrimSpace(c.Query("driver"))
		if driverName == "" {
			driverName = strings.TrimSpace(c.PostForm("driver"))
		}
		driver, err := driverStorage.GetByName(driverName)
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
