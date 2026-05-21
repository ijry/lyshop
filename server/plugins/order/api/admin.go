package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/response"
	ordersvc "github.com/ijry/lyshop/plugins/order/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/orders", adminListOrders)
	g.PUT("/orders/:id/ship", adminShipOrder)
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
