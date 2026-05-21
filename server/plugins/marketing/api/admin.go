package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/response"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
	mktsvc "github.com/ijry/lyshop/plugins/marketing/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/marketing/coupons", adminListCoupons)
	g.POST("/marketing/coupons", adminCreateCoupon)
	g.GET("/marketing/activities", adminListActivities)
	g.POST("/marketing/activities", adminCreateActivity)
}

func adminListCoupons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := mktsvc.ListCoupons(c.Request.Context(), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminCreateCoupon(c *gin.Context) {
	var coupon mktmodel.Coupon
	if err := c.ShouldBindJSON(&coupon); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := mktsvc.CreateCoupon(c.Request.Context(), &coupon); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, coupon)
}

func adminListActivities(c *gin.Context) {
	actType, _ := strconv.ParseInt(c.Query("type"), 10, 8)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := mktsvc.ListActivities(c.Request.Context(), int8(actType), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminCreateActivity(c *gin.Context) {
	var req struct {
		Activity mktmodel.Activity          `json:"activity"`
		Products []mktmodel.ActivityProduct `json:"products"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := mktsvc.CreateActivity(c.Request.Context(), &req.Activity, req.Products); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, req.Activity)
}
