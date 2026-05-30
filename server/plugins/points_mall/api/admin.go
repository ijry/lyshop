package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/server/core/middleware"
	pmmodel "github.com/ijry/lyshop/server/plugins/points_mall/model"
	pmservice "github.com/ijry/lyshop/server/plugins/points_mall/service"
)

// RegisterAdminRoutes 注册管理端路由
func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.Use(middleware.RequirePermission("points_mall:view"))

	// 积分商品管理
	g.GET("/points/products", adminListProducts)
	g.POST("/points/products", middleware.RequirePermission("points_mall:edit"), adminCreateProduct)
	g.PUT("/points/products/:id", middleware.RequirePermission("points_mall:edit"), adminUpdateProduct)
	g.DELETE("/points/products/:id", middleware.RequirePermission("points_mall:edit"), adminDeleteProduct)
	g.PUT("/points/products/:id/status", middleware.RequirePermission("points_mall:edit"), adminUpdateProductStatus)

	// 兑换记录管理
	g.GET("/points/exchanges", adminListExchanges)
	g.GET("/points/exchanges/:id", adminGetExchange)
	g.PUT("/points/exchanges/:id/ship", middleware.RequirePermission("points_mall:edit"), adminShipExchange)
	g.PUT("/points/exchanges/:id/complete", middleware.RequirePermission("points_mall:edit"), adminCompleteExchange)
	g.PUT("/points/exchanges/:id/cancel", middleware.RequirePermission("points_mall:edit"), adminCancelExchange)

	// 积分管理
	g.GET("/points/logs", adminListPointsLogs)
	g.POST("/points/adjust", middleware.RequirePermission("points_mall:edit"), adminAdjustPoints)
	g.GET("/points/stats", adminGetPointsStats)

	// 配置管理
	g.GET("/points/config", adminGetConfig)
	g.PUT("/points/config", middleware.RequirePermission("points_mall:edit"), adminUpdateConfig)
}

// adminListProducts 积分商品列表
func adminListProducts(c *gin.Context) {
	productType := c.Query("type")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	products, total, err := pmservice.ListProducts(c.Request.Context(), productType, int8(status), page, size)
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

// adminCreateProduct 创建积分商品
func adminCreateProduct(c *gin.Context) {
	var product pmmodel.PointsProduct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := pmservice.CreateProduct(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "创建成功", "data": product})
}

// adminUpdateProduct 更新积分商品
func adminUpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := pmservice.UpdateProduct(c.Request.Context(), id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功"})
}

// adminDeleteProduct 删除积分商品
func adminDeleteProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := pmservice.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
}

// adminUpdateProductStatus 更新商品状态
func adminUpdateProductStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := pmservice.UpdateProductStatus(c.Request.Context(), id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功"})
}

// adminListExchanges 兑换记录列表
func adminListExchanges(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	exchanges, total, err := pmservice.ListExchanges(c.Request.Context(), userID, status, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"list":  exchanges,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// adminGetExchange 兑换详情
func adminGetExchange(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	exchange, err := pmservice.GetExchange(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "兑换记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": exchange})
}

// adminShipExchange 发货
func adminShipExchange(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		TrackingNo string `json:"tracking_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := pmservice.ShipExchange(c.Request.Context(), id, req.TrackingNo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "发货成功"})
}

// adminCompleteExchange 完成兑换
func adminCompleteExchange(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := pmservice.CompleteExchange(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "操作成功"})
}

// adminCancelExchange 取消兑换
func adminCancelExchange(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := pmservice.CancelExchange(c.Request.Context(), id, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "取消成功"})
}

// adminListPointsLogs 积分日志列表
func adminListPointsLogs(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	logType, _ := strconv.Atoi(c.Query("type"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	logs, total, err := pmservice.AdminListPointsLogs(c.Request.Context(), userID, int8(logType), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"list":  logs,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// adminAdjustPoints 调整用户积分
func adminAdjustPoints(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id"`
		Points int    `json:"points"`
		Remark string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if req.UserID == 0 || req.Points == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := pmservice.AddPoints(c.Request.Context(), req.UserID, req.Points, 5, req.Remark); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "调整成功"})
}

// adminGetPointsStats 积分统计
func adminGetPointsStats(c *gin.Context) {
	stats, err := pmservice.GetPointsStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": stats})
}

// adminGetConfig 获取积分配置
func adminGetConfig(c *gin.Context) {
	// TODO: 实现配置获取逻辑
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"points_to_yuan":      100,
			"order_points_rate":   0.01,
			"points_expire_days":  0,
			"enable_order_points": false,
		},
	})
}

// adminUpdateConfig 更新积分配置
func adminUpdateConfig(c *gin.Context) {
	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// TODO: 实现配置保存逻辑

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "保存成功"})
}
