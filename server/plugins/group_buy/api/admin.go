package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/api"
	"github.com/ijry/lyshop/server/plugins/group_buy/model"
	"github.com/ijry/lyshop/server/plugins/group_buy/service"
)

// RegisterAdminRoutes 注册管理端路由
func RegisterAdminRoutes(r *gin.RouterGroup) {
	g := r.Group("/admin/api/group-buy")
	{
		// 活动管理
		g.GET("/activities", api.RequirePermission("group_buy:view"), listActivities)
		g.POST("/activities", api.RequirePermission("group_buy:edit"), createActivity)
		g.PUT("/activities/:id", api.RequirePermission("group_buy:edit"), updateActivity)
		g.DELETE("/activities/:id", api.RequirePermission("group_buy:delete"), deleteActivity)

		// 商品管理
		g.GET("/products", api.RequirePermission("group_buy:view"), listProducts)
		g.PUT("/activities/:id/products", api.RequirePermission("group_buy:edit"), upsertProducts)

		// 拼团订单管理
		g.GET("/orders", api.RequirePermission("group_buy:view"), listGroupOrders)
		g.GET("/orders/:id", api.RequirePermission("group_buy:view"), getGroupOrder)
		g.GET("/orders/:id/members", api.RequirePermission("group_buy:view"), listGroupMembers)
	}
}

func listActivities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var status *int8
	if s := c.Query("status"); s != "" {
		val, _ := strconv.Atoi(s)
		st := int8(val)
		status = &st
	}

	activities, total, err := service.ListActivities(c.Request.Context(), page, size, status)
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

func createActivity(c *gin.Context) {
	var activity model.GroupBuyActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	if err := service.CreateActivity(c.Request.Context(), &activity); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, activity)
}

func updateActivity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	if err := service.UpdateActivity(c.Request.Context(), id, updates); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, nil)
}

func deleteActivity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := service.DeleteActivity(c.Request.Context(), id); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, nil)
}

func listProducts(c *gin.Context) {
	activityID, _ := strconv.ParseUint(c.Query("activity_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	products, total, err := service.ListProducts(c.Request.Context(), activityID, page, size)
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

func upsertProducts(c *gin.Context) {
	activityID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var products []model.GroupBuyProduct
	if err := c.ShouldBindJSON(&products); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	if err := service.UpsertProducts(c.Request.Context(), activityID, products); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, nil)
}

func listGroupOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := c.Query("status")

	orders, total, err := service.ListGroupOrders(c.Request.Context(), page, size, status)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

func getGroupOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	order, err := service.GetGroupOrder(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, order)
}

func listGroupMembers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	members, err := service.ListGroupMembers(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, members)
}
