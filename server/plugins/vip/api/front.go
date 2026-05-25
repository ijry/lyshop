package api

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	vipsvc "github.com/ijry/lyshop/plugins/vip/service"
)

func RegisterFrontRoutes(g *gin.RouterGroup) {
	auth := g.Group("")
	auth.Use(middleware.RequireAuth)
	auth.GET("/vip/profile", getProfile)
	auth.POST("/vip/open", openMembership)
	auth.GET("/vip/coupons/monthly", listMonthlyCoupons)
	auth.POST("/vip/coupons/monthly/:rule_id/claim", claimMonthlyCoupon)
	auth.GET("/vip/growth/logs", listGrowthLogs)
}

func getProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	data, err := vipsvc.GetProfile(c.Request.Context(), userID.(uint64), time.Now())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, data)
}

func openMembership(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		PlanID uint64 `json:"plan_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.PlanID == 0 {
		response.Fail(c, 400, "plan_id 必填")
		return
	}
	asset, err := vipsvc.OpenMembership(c.Request.Context(), userID.(uint64), req.PlanID, time.Now())
	if err != nil {
		response.Fail(c, 30001, err.Error())
		return
	}
	response.OK(c, asset)
}

func listMonthlyCoupons(c *gin.Context) {
	userID, _ := c.Get("user_id")
	list, err := vipsvc.ListMonthlyCoupons(c.Request.Context(), userID.(uint64), time.Now())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func claimMonthlyCoupon(c *gin.Context) {
	userID, _ := c.Get("user_id")
	ruleID, _ := strconv.ParseUint(c.Param("rule_id"), 10, 64)
	if ruleID == 0 {
		response.Fail(c, 400, "rule_id 无效")
		return
	}
	if err := vipsvc.ClaimMonthlyCoupon(c.Request.Context(), userID.(uint64), ruleID, time.Now()); err != nil {
		response.Fail(c, 30001, err.Error())
		return
	}
	response.OK(c, nil)
}

func listGrowthLogs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := vipsvc.ListGrowthLogs(c.Request.Context(), userID.(uint64), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}
