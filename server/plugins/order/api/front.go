package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	ordermodel "github.com/ijry/lyshop/plugins/order/model"
	ordersvc "github.com/ijry/lyshop/plugins/order/service"
)

func RegisterFrontRoutes(g *gin.RouterGroup) {
	auth := g.Group("")
	auth.Use(middleware.RequireAuth)

	auth.GET("/cart", getCart)
	auth.POST("/cart/add", addToCart)
	auth.PUT("/cart/qty", updateCartQty)
	auth.DELETE("/cart/:sku_id", removeFromCart)

	auth.GET("/addresses", listAddresses)
	auth.POST("/addresses", createAddress)

	auth.POST("/orders", createOrder)
	auth.GET("/orders", myOrders)
}

func getCart(c *gin.Context) {
	userID, _ := c.Get("user_id")
	items, err := ordersvc.GetCart(c.Request.Context(), userID.(uint64))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, items)
}

func addToCart(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		SkuID uint64 `json:"sku_id" binding:"required"`
		Qty   int    `json:"qty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if req.Qty <= 0 {
		req.Qty = 1
	}
	if err := ordersvc.AddToCart(c.Request.Context(), userID.(uint64), req.SkuID, req.Qty); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func updateCartQty(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		SkuID uint64 `json:"sku_id"`
		Qty   int    `json:"qty"`
	}
	c.ShouldBindJSON(&req)
	ordersvc.UpdateCartQty(c.Request.Context(), userID.(uint64), req.SkuID, req.Qty)
	response.OK(c, nil)
}

func removeFromCart(c *gin.Context) {
	userID, _ := c.Get("user_id")
	skuID, _ := strconv.ParseUint(c.Param("sku_id"), 10, 64)
	ordersvc.RemoveFromCart(c.Request.Context(), userID.(uint64), skuID)
	response.OK(c, nil)
}

func listAddresses(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var list []ordermodel.Address
	response.OK(c, list) // placeholder — full impl in later iteration
	_ = userID
}

func createAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var addr ordermodel.Address
	if err := c.ShouldBindJSON(&addr); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	addr.UserID = userID.(uint64)
	if err := ordersvc.CreateAddress(c.Request.Context(), &addr); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, addr)
}

func createOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req ordersvc.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	req.UserID = userID.(uint64)
	order, err := ordersvc.CreateOrder(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, 20001, err.Error())
		return
	}
	response.OK(c, order)
}

func myOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	status, _ := strconv.ParseInt(c.Query("status"), 10, 8)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := ordersvc.ListOrders(c.Request.Context(), userID.(uint64), int8(status), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}
