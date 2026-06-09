package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	inventorycore "github.com/ijry/lyshop/core/inventory"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/inventory/tasks", middleware.RequirePermission("system:config"), listTasks)
	g.POST("/inventory/tasks/:id/retry", middleware.RequirePermission("system:config"), retryTask)
	g.POST("/external-wms/callback", callback)
}

func listTasks(c *gin.Context) {
	var rows []inventorycore.InventoryIntegrationTask
	if err := db.DB.WithContext(c.Request.Context()).Order("id desc").Limit(50).Find(&rows).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, rows)
}

func retryTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := db.DB.WithContext(c.Request.Context()).
		Model(&inventorycore.InventoryIntegrationTask{}).
		Where("id = ?", id).
		Updates(map[string]any{"status": "pending", "last_error": ""}).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, gin.H{"retried": id})
}

func callback(c *gin.Context) {
	response.OK(c, gin.H{"received": true})
}
