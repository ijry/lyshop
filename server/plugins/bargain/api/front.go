package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/api"
	"github.com/ijry/lyshop/server/plugins/bargain/service"
)

// RegisterFrontRoutes 注册前端路由
func RegisterFrontRoutes(r *gin.RouterGroup) {
	g := r.Group("/api/v1/bargain")
	{
		// 商品列表
		g.GET("/products", listProducts)
		g.GET("/products/:id", getProduct)

		// 活动列表
		g.GET("/activities", listActivities)

		// 砍价订单
		g.POST("/orders", createBargainOrder)
		g.POST("/orders/:id/help", helpBargain)
		g.GET("/orders/:id", getBargainOrder)
		g.GET("/orders/:id/helpers", listHelpers)
	}
}

func listProducts(c *gin.Context) {
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

func getProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	product, err := service.GetProduct(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, product)
}

func listActivities(c *gin.Context) {
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

func createBargainOrder(c *gin.Context) {
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

	order, err := service.CreateBargainOrder(c.Request.Context(), req.ActivityID, req.ProductID, req.SkuID, userID)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, order)
}

func helpBargain(c *gin.Context) {
	bargainOrderID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// TODO: 从上下文获取用户ID
	helperUserID := uint64(2)

	cutAmount, err := service.HelpBargain(c.Request.Context(), bargainOrderID, helperUserID)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, gin.H{
		"cut_amount": cutAmount,
	})
}

func getBargainOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	order, err := service.GetBargainOrder(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, order)
}

func listHelpers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	helpers, err := service.ListBargainHelpers(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, helpers)
}
