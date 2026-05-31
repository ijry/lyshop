package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/api"
	"github.com/ijry/lyshop/plugins/group_buy/service"
)

// RegisterFrontRoutes 注册前端路由
func RegisterFrontRoutes(r *gin.RouterGroup) {
	g := r.Group("/api/v1/group-buy")
	{
		// 商品列表
		g.GET("/products", frontListProducts)
		g.GET("/products/:id", frontGetProduct)

		// 活动列表
		g.GET("/activities", frontListActivities)

		// 拼团订单
		g.POST("/orders", createGroupOrder)
		g.POST("/orders/:id/join", joinGroupOrder)
		g.GET("/orders/:id", frontGetGroupOrder)
		g.GET("/orders/:id/members", listMembers)
	}
}

func frontListProducts(c *gin.Context) {
	activityID, _ := strconv.ParseUint(c.Query("activity_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	products, total, err := service.ListFrontProducts(c.Request.Context(), activityID, page, size)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, gin.H{
		"list":  products,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

func frontGetProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	product, err := service.GetProduct(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, product)
}

func frontListActivities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := int8(1)

	activities, total, err := service.ListActivities(c.Request.Context(), page, size, &status)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, gin.H{
		"list":  activities,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

func createGroupOrder(c *gin.Context) {
	var req struct {
		ActivityID uint64 `json:"activity_id" binding:"required"`
		ProductID  uint64 `json:"product_id" binding:"required"`
		SkuID      uint64 `json:"sku_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	// TODO: 从上下文获取用户ID
	userID := uint64(1)

	order, err := service.CreateGroupBuyOrder(c.Request.Context(), req.ActivityID, req.ProductID, req.SkuID, userID)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, order)
}

func joinGroupOrder(c *gin.Context) {
	groupOrderID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		OrderID uint64 `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	// TODO: 从上下文获取用户ID
	userID := uint64(1)

	if err := service.JoinGroupBuy(c.Request.Context(), groupOrderID, userID, req.OrderID); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, nil)
}

func frontGetGroupOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	order, err := service.GetGroupOrder(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, order)
}

func listMembers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	members, err := service.ListGroupMembers(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, members)
}
