package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/server/core/middleware"
	pmservice "github.com/ijry/lyshop/server/plugins/points_mall/service"
)

// RegisterFrontRoutes 注册前端路由
func RegisterFrontRoutes(g *gin.RouterGroup) {
	// 积分商城（无需认证）
	g.GET("/points/products", listPointsProducts)
	g.GET("/points/products/:id", getPointsProductDetail)

	// 需要认证的路由
	auth := g.Group("")
	auth.Use(middleware.RequireAuth)

	// 兑换操作
	auth.POST("/points/products/:id/exchange", exchangeProduct)
	auth.GET("/points/exchanges", listMyExchanges)
	auth.GET("/points/exchanges/:id", getExchangeDetail)
	auth.POST("/points/exchanges/:id/confirm", confirmReceive)

	// 积分日志
	auth.GET("/points/logs", listPointsLogs)
	auth.GET("/points/balance", getPointsBalance)
}

// listPointsProducts 积分商品列表
func listPointsProducts(c *gin.Context) {
	productType := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	products, total, err := pmservice.ListProducts(c.Request.Context(), productType, 1, page, size)
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

// getPointsProductDetail 积分商品详情
func getPointsProductDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	product, err := pmservice.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "商品不存在"})
		return
	}

	if product.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "商品已下架"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": product})
}

// exchangeProduct 兑换商品
func exchangeProduct(c *gin.Context) {
	userID := c.GetUint64("user_id")
	productID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Qty             int             `json:"qty"`
		AddressSnapshot json.RawMessage `json:"address_snapshot"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if req.Qty <= 0 {
		req.Qty = 1
	}

	exchange, err := pmservice.ExchangeProduct(c.Request.Context(), userID, productID, req.Qty, req.AddressSnapshot)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "兑换成功",
		"data": exchange,
	})
}

// listMyExchanges 我的兑换记录
func listMyExchanges(c *gin.Context) {
	userID := c.GetUint64("user_id")
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

// getExchangeDetail 兑换详情
func getExchangeDetail(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	exchange, err := pmservice.GetExchange(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "兑换记录不存在"})
		return
	}

	// 验证权限
	if exchange.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权访问"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": exchange})
}

// confirmReceive 确认收货
func confirmReceive(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	exchange, err := pmservice.GetExchange(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "兑换记录不存在"})
		return
	}

	if exchange.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权操作"})
		return
	}

	if err := pmservice.CompleteExchange(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "确认收货成功"})
}

// listPointsLogs 积分日志
func listPointsLogs(c *gin.Context) {
	userID := c.GetUint64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	logs, total, err := pmservice.ListPointsLogs(c.Request.Context(), userID, page, size)
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

// getPointsBalance 获取积分余额
func getPointsBalance(c *gin.Context) {
	userID := c.GetUint64("user_id")

	points, err := pmservice.GetUserPoints(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{"points": points},
	})
}
