package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	seckillmodel "github.com/ijry/lyshop/plugins/seckill/model"
	seckillservice "github.com/ijry/lyshop/plugins/seckill/service"
)

// RegisterAdminRoutes 注册管理端路由
func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.Use(middleware.RequirePermission("seckill:view"))

	// 活动管理
	g.GET("/seckill/activities", adminListActivities)
	g.POST("/seckill/activities", middleware.RequirePermission("seckill:edit"), adminCreateActivity)
	g.PUT("/seckill/activities/:id", middleware.RequirePermission("seckill:edit"), adminUpdateActivity)
	g.DELETE("/seckill/activities/:id", middleware.RequirePermission("seckill:edit"), adminDeleteActivity)

	// 商品管理
	g.GET("/seckill/products", adminListProducts)
	g.PUT("/seckill/activities/:id/products", middleware.RequirePermission("seckill:edit"), adminUpsertProducts)
}

// adminListActivities 活动列表
func adminListActivities(c *gin.Context) {
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	activities, total, err := seckillservice.ListActivities(c.Request.Context(), int8(status), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"list":  activities,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// adminCreateActivity 创建活动
func adminCreateActivity(c *gin.Context) {
	var activity seckillmodel.SeckillActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := seckillservice.CreateActivity(c.Request.Context(), &activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": activity})
}

// adminUpdateActivity 更新活动
func adminUpdateActivity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := seckillservice.UpdateActivity(c.Request.Context(), id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功"})
}

// adminDeleteActivity 删除活动
func adminDeleteActivity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := seckillservice.DeleteActivity(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
}

// adminListProducts 商品列表
func adminListProducts(c *gin.Context) {
	activityID, _ := strconv.ParseUint(c.Query("activity_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	products, total, err := seckillservice.ListProducts(c.Request.Context(), activityID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"list":  products,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// adminUpsertProducts 批量更新商品
func adminUpsertProducts(c *gin.Context) {
	activityID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var products []seckillmodel.SeckillProduct
	if err := c.ShouldBindJSON(&products); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := seckillservice.UpsertProducts(c.Request.Context(), activityID, products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "保存成功"})
}
