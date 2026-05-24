package api

import (
	"strconv"
	"strings"

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
	auth.PUT("/addresses/:id", updateAddress)
	auth.DELETE("/addresses/:id", deleteAddress)

	auth.POST("/orders", createOrder)
	auth.GET("/orders", myOrders)
	auth.GET("/orders/:id", myOrderDetail)
	auth.GET("/orders/:id/review", myOrderReviewMeta)
	auth.POST("/orders/:id/pay", payOrder)
	auth.POST("/orders/:id/review", reviewOrder)
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
	list, err := ordersvc.ListAddresses(c.Request.Context(), userID.(uint64))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
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

func updateAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var addr ordermodel.Address
	if err := c.ShouldBindJSON(&addr); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	saved, err := ordersvc.UpdateAddress(c.Request.Context(), userID.(uint64), id, addr)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, saved)
}

func deleteAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := ordersvc.DeleteAddress(c.Request.Context(), userID.(uint64), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
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

func myOrderDetail(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	detail, err := ordersvc.GetOrderDetail(c.Request.Context(), userID.(uint64), id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.OK(c, detail)
}

func payOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := ordersvc.PayOrder(c.Request.Context(), userID.(uint64), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func myOrderReviewMeta(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	meta, err := ordersvc.GetOrderReviewMeta(c.Request.Context(), userID.(uint64), id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.OK(c, meta)
}

func reviewOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Mode           string `json:"mode"`
		LogisticsScore int8   `json:"logistics_score"`
		Items          []struct {
			OrderItemID  uint64 `json:"order_item_id"`
			ProductScore int8   `json:"product_score"`
			Content      string `json:"content"`
		} `json:"items"`
		AppendContent string `json:"append_content"`
		Content       string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	var items []ordersvc.ReviewItemInput
	for _, item := range req.Items {
		items = append(items, ordersvc.ReviewItemInput{
			OrderItemID:  item.OrderItemID,
			ProductScore: item.ProductScore,
			Content:      item.Content,
		})
	}
	if len(items) == 0 {
		detail, err := ordersvc.GetOrderDetail(c.Request.Context(), userID.(uint64), id)
		if err != nil {
			response.Fail(c, 404, err.Error())
			return
		}
		seedContent := strings.TrimSpace(req.Content)
		for _, item := range detail.Items {
			if ordersvc.ReviewMode(req.Mode) == ordersvc.ReviewModeAppend && item.Review == nil {
				continue
			}
			items = append(items, ordersvc.ReviewItemInput{
				OrderItemID:  item.ID,
				ProductScore: 5,
				Content:      seedContent,
			})
		}
		if ordersvc.ReviewMode(req.Mode) == ordersvc.ReviewModeAppend && strings.TrimSpace(req.AppendContent) == "" {
			req.AppendContent = seedContent
		}
	}

	if err := ordersvc.SubmitOrderReview(c.Request.Context(), ordersvc.SubmitOrderReviewReq{
		OrderID:        id,
		UserID:         userID.(uint64),
		Mode:           ordersvc.ReviewMode(req.Mode),
		LogisticsScore: req.LogisticsScore,
		Items:          items,
		AppendContent:  req.AppendContent,
	}); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}
