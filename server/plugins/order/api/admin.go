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
	g.POST("/orders/:id/shipments/:shipment_id/sync", middleware.RequirePermission("order:ship"), adminSyncShipment)
	g.GET("/orders/:id/shipments/:shipment_id/tracks", middleware.RequirePermission("order:view"), adminGetShipmentTracks)
	g.GET("/after-sales", middleware.RequirePermission("order:view"), adminListAfterSales)
	g.GET("/after-sales/:id", middleware.RequirePermission("order:view"), adminGetAfterSaleDetail)
	g.POST("/after-sales/:id/audit", middleware.RequirePermission("order:ship"), adminAuditAfterSale)
	g.POST("/after-sales/:id/receive", middleware.RequirePermission("order:ship"), adminReceiveAfterSale)
	g.POST("/after-sales/:id/refund", middleware.RequirePermission("order:refund"), adminRefundAfterSale)
	g.POST("/after-sales/:id/complete", middleware.RequirePermission("order:ship"), adminCompleteAfterSale)
	g.POST("/after-sales/:id/close", middleware.RequirePermission("order:ship"), adminCloseAfterSale)
	g.GET("/reviews", middleware.RequirePermission("order:view"), adminListReviews)
	g.GET("/reviews/:id", middleware.RequirePermission("order:view"), adminGetReview)
	g.POST("/reviews/:id/reply", middleware.RequirePermission("order:review-reply"), adminReplyReview)
}

func adminListOrders(c *gin.Context) {
	status, _ := strconv.ParseInt(c.Query("status"), 10, 8)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := ordersvc.AdminListOrders(c, int8(status), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminGetOrderDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	detail, err := ordersvc.AdminGetOrderDetail(c, id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.OK(c, detail)
}

func adminShipOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		DeliveryType    string `json:"delivery_type"`
		TrackingNo      string `json:"tracking_no"`
		ShipType        string `json:"ship_type"`
		AfterSaleCaseID uint64 `json:"after_sale_case_id"`
		Company         string `json:"company"`
		RiderName       string `json:"rider_name"`
		RiderPhone      string `json:"rider_phone"`
		Remark          string `json:"remark"`
	}
	c.ShouldBindJSON(&req)
	if err := ordersvc.ShipOrder(c.Request.Context(), id, ordersvc.ShipOrderReq{
		DeliveryType:    req.DeliveryType,
		ShipType:        req.ShipType,
		AfterSaleCaseID: req.AfterSaleCaseID,
		Company:         req.Company,
		TrackingNo:      req.TrackingNo,
		RiderName:       req.RiderName,
		RiderPhone:      req.RiderPhone,
		Remark:          req.Remark,
	}); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminSyncShipment(c *gin.Context) {
	shipmentID, _ := strconv.ParseUint(c.Param("shipment_id"), 10, 64)
	if err := ordersvc.SyncShipmentTracks(c.Request.Context(), shipmentID, ordersvc.SyncShipmentReq{Manual: true}); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminGetShipmentTracks(c *gin.Context) {
	orderID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	shipmentID, _ := strconv.ParseUint(c.Param("shipment_id"), 10, 64)
	rows, err := ordersvc.ListShipmentTracks(c.Request.Context(), orderID, shipmentID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, rows)
}

func adminListAfterSales(c *gin.Context) {
	status := c.Query("status")
	caseType := c.Query("case_type")
	orderID, _ := strconv.ParseUint(c.Query("order_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := ordersvc.ListAfterSales(c, status, caseType, orderID, page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminGetAfterSaleDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	detail, err := ordersvc.GetAfterSale(c, id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.OK(c, detail)
}

func adminAuditAfterSale(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	adminIDRaw, _ := c.Get("user_id")
	adminID, _ := adminIDRaw.(uint64)
	var req struct {
		Approve     bool   `json:"approve"`
		AuditRemark string `json:"audit_remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := ordersvc.AuditAfterSale(c.Request.Context(), id, ordersvc.AuditAfterSaleReq{
		Approve:     req.Approve,
		AuditRemark: req.AuditRemark,
		AdminID:     adminID,
	}); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminReceiveAfterSale(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	adminIDRaw, _ := c.Get("user_id")
	adminID, _ := adminIDRaw.(uint64)
	if err := ordersvc.ReceiveAfterSale(c.Request.Context(), id, adminID); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminRefundAfterSale(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	adminIDRaw, _ := c.Get("user_id")
	adminID, _ := adminIDRaw.(uint64)
	var req struct {
		Amount   float64 `json:"amount"`
		Reason   string  `json:"reason"`
		RefundNo string  `json:"refund_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := ordersvc.MarkRefund(c.Request.Context(), id, ordersvc.MarkRefundReq{
		Amount:   req.Amount,
		Reason:   req.Reason,
		RefundNo: req.RefundNo,
		AdminID:  adminID,
	}); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminCompleteAfterSale(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	adminIDRaw, _ := c.Get("user_id")
	adminID, _ := adminIDRaw.(uint64)
	if err := ordersvc.CompleteAfterSale(c.Request.Context(), id, adminID); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminCloseAfterSale(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	adminIDRaw, _ := c.Get("user_id")
	adminID, _ := adminIDRaw.(uint64)
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := ordersvc.CloseAfterSale(c.Request.Context(), id, ordersvc.CloseAfterSaleReq{
		Reason:  req.Reason,
		AdminID: adminID,
	}); err != nil {
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
