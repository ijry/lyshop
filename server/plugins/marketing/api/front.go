package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	mktsvc "github.com/ijry/lyshop/plugins/marketing/service"
)

func RegisterFrontRoutes(g *gin.RouterGroup) {
	g.GET("/marketing/seckills", listSeckills)

	auth := g.Group("")
	auth.Use(middleware.RequireAuth)
	auth.GET("/coupons", listClaimableCoupons)
	auth.POST("/coupons/:id/claim", claimCoupon)
	auth.GET("/user/coupons", listUserCoupons)
	auth.GET("/user/points/logs", listPointsLogs)
}

func listSeckills(c *gin.Context) {
	list, err := mktsvc.GetActiveSeckills(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func claimCoupon(c *gin.Context) {
	couponID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID, _ := c.Get("user_id")
	if err := mktsvc.ClaimCoupon(c.Request.Context(), couponID, userID.(uint64)); err != nil {
		response.Fail(c, 30001, err.Error())
		return
	}
	response.OK(c, nil)
}

func listClaimableCoupons(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, err := mktsvc.ListClaimableCoupons(c.Request.Context(), userID.(uint64), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func listUserCoupons(c *gin.Context) {
	userID, _ := c.Get("user_id")
	list, err := mktsvc.ListUserCoupons(c.Request.Context(), userID.(uint64))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func listPointsLogs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := mktsvc.ListPointsLogs(c.Request.Context(), userID.(uint64), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}
