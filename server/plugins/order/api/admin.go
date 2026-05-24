package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	ordersvc "github.com/ijry/lyshop/plugins/order/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/orders", middleware.RequirePermission("order:view"), adminListOrders)
	g.GET("/orders/:id", middleware.RequirePermission("order:view"), adminGetOrderDetail)
	g.PUT("/orders/:id/ship", middleware.RequirePermission("order:ship"), adminShipOrder)
	g.GET("/reviews", middleware.RequirePermission("order:view"), adminListReviews)
	g.GET("/reviews/:id", middleware.RequirePermission("order:view"), adminGetReview)
	g.POST("/reviews/:id/reply", middleware.RequirePermission("order:review-reply"), adminReplyReview)
}

func adminListOrders(c *gin.Context) {
	status, _ := strconv.ParseInt(c.Query("status"), 10, 8)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := ordersvc.AdminListOrders(c.Request.Context(), int8(status), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminGetOrderDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	detail, err := ordersvc.AdminGetOrderDetail(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.OK(c, detail)
}

func adminShipOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		TrackingNo string `json:"tracking_no"`
	}
	c.ShouldBindJSON(&req)
	if err := ordersvc.ShipOrder(c.Request.Context(), id, req.TrackingNo); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminListReviews(c *gin.Context) {
	productID, _ := strconv.ParseUint(c.Query("product_id"), 10, 64)
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	data, err := ordersvc.AdminListReviews(c.Request.Context(), productID, keyword, page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, data)
}

func adminGetReview(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	detail, err := ordersvc.AdminGetReview(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.OK(c, detail)
}

func adminReplyReview(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	adminIDRaw, _ := c.Get("user_id")
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := ordersvc.AdminUpsertReply(c.Request.Context(), id, ordersvc.AdminReviewReplyReq{
		AdminID: adminIDRaw.(uint64),
		Content: req.Content,
	}); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}
