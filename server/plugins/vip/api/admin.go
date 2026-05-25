package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	vipmodel "github.com/ijry/lyshop/plugins/vip/model"
	vipsvc "github.com/ijry/lyshop/plugins/vip/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/vip/plans", middleware.RequirePermission("vip:view"), adminListPlans)
	g.POST("/vip/plans", middleware.RequirePermission("vip:edit"), adminCreatePlan)

	g.GET("/vip/levels", middleware.RequirePermission("vip:view"), adminListLevels)
	g.POST("/vip/levels", middleware.RequirePermission("vip:edit"), adminCreateLevel)

	g.GET("/vip/coupon-rules", middleware.RequirePermission("vip:view"), adminListCouponRules)
	g.POST("/vip/coupon-rules", middleware.RequirePermission("vip:edit"), adminCreateCouponRule)

	g.GET("/vip/sku-prices", middleware.RequirePermission("vip:view"), adminListSkuPrices)
	g.POST("/vip/sku-prices", middleware.RequirePermission("vip:edit"), adminCreateSkuPrice)
}

func adminListPlans(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := vipsvc.ListPlans(c.Request.Context(), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminCreatePlan(c *gin.Context) {
	var row vipmodel.Plan
	if err := c.ShouldBindJSON(&row); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := vipsvc.CreatePlan(c.Request.Context(), &row); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, row)
}

func adminListLevels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := vipsvc.ListLevels(c.Request.Context(), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminCreateLevel(c *gin.Context) {
	var row vipmodel.Level
	if err := c.ShouldBindJSON(&row); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := vipsvc.CreateLevel(c.Request.Context(), &row); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, row)
}

func adminListCouponRules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := vipsvc.ListCouponRules(c.Request.Context(), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminCreateCouponRule(c *gin.Context) {
	var row vipmodel.CouponRule
	if err := c.ShouldBindJSON(&row); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := vipsvc.CreateCouponRule(c.Request.Context(), &row); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, row)
}

func adminListSkuPrices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := vipsvc.ListSkuPrices(c.Request.Context(), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminCreateSkuPrice(c *gin.Context) {
	var row vipmodel.SkuPrice
	if err := c.ShouldBindJSON(&row); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := vipsvc.CreateSkuPrice(c.Request.Context(), &row); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, row)
}
