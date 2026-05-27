package api

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
	mktsvc "github.com/ijry/lyshop/plugins/marketing/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/marketing/coupons", middleware.RequirePermission("marketing:view"), adminListCoupons)
	g.POST("/marketing/coupons", middleware.RequirePermission("marketing:edit"), adminCreateCoupon)

	g.GET("/marketing/seckill/activities", middleware.RequirePermission("marketing:view"), adminListSeckillActivities)
	g.POST("/marketing/seckill/activities", middleware.RequirePermission("marketing:edit"), adminCreateSeckillActivity)
	g.PUT("/marketing/seckill/activities/:id", middleware.RequirePermission("marketing:edit"), adminUpdateSeckillActivity)
	g.GET("/marketing/seckill/products", middleware.RequirePermission("marketing:view"), adminListSeckillProducts)
	g.PUT("/marketing/seckill/activities/:id/products", middleware.RequirePermission("marketing:edit"), adminUpsertSeckillProducts)

	g.GET("/marketing/group-buy/activities", middleware.RequirePermission("marketing:view"), adminListGroupBuyActivities)
	g.POST("/marketing/group-buy/activities", middleware.RequirePermission("marketing:edit"), adminCreateGroupBuyActivity)
	g.PUT("/marketing/group-buy/activities/:id", middleware.RequirePermission("marketing:edit"), adminUpdateGroupBuyActivity)
	g.GET("/marketing/group-buy/products", middleware.RequirePermission("marketing:view"), adminListGroupBuyProducts)
	g.PUT("/marketing/group-buy/activities/:id/products", middleware.RequirePermission("marketing:edit"), adminUpsertGroupBuyProducts)

	g.GET("/marketing/bargain/activities", middleware.RequirePermission("marketing:view"), adminListBargainActivities)
	g.POST("/marketing/bargain/activities", middleware.RequirePermission("marketing:edit"), adminCreateBargainActivity)
	g.PUT("/marketing/bargain/activities/:id", middleware.RequirePermission("marketing:edit"), adminUpdateBargainActivity)
	g.GET("/marketing/bargain/products", middleware.RequirePermission("marketing:view"), adminListBargainProducts)
	g.PUT("/marketing/bargain/activities/:id/products", middleware.RequirePermission("marketing:edit"), adminUpsertBargainProducts)
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

type activityReq struct {
	Name       string     `json:"name"`
	StartAt    *time.Time `json:"start_at"`
	EndAt      *time.Time `json:"end_at"`
	Status     int8       `json:"status"`
	PriceRule  string     `json:"price_rule"`
	PriceValue float64    `json:"price_value"`
	Exclusive  bool       `json:"exclusive"`
}

type listActivitiesQuery struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

type listManagedProductsQuery struct {
	ActivityID uint64 `form:"activity_id"`
	ProductID  uint64 `form:"product_id"`
	CategoryID uint64 `form:"category_id"`
	Keyword    string `form:"keyword"`
	Page       int    `form:"page"`
	Size       int    `form:"size"`
}

func adminListSeckillActivities(c *gin.Context) { adminListActivities(c, mktmodel.ActivityTypeSeckill) }
func adminListGroupBuyActivities(c *gin.Context) {
	adminListActivities(c, mktmodel.ActivityTypeGroupBuy)
}
func adminListBargainActivities(c *gin.Context) { adminListActivities(c, mktmodel.ActivityTypeBargain) }

func adminCreateSeckillActivity(c *gin.Context) { adminCreateActivity(c, mktmodel.ActivityTypeSeckill) }
func adminCreateGroupBuyActivity(c *gin.Context) {
	adminCreateActivity(c, mktmodel.ActivityTypeGroupBuy)
}
func adminCreateBargainActivity(c *gin.Context) { adminCreateActivity(c, mktmodel.ActivityTypeBargain) }

func adminUpdateSeckillActivity(c *gin.Context) { adminUpdateActivity(c, mktmodel.ActivityTypeSeckill) }
func adminUpdateGroupBuyActivity(c *gin.Context) {
	adminUpdateActivity(c, mktmodel.ActivityTypeGroupBuy)
}
func adminUpdateBargainActivity(c *gin.Context) { adminUpdateActivity(c, mktmodel.ActivityTypeBargain) }

func adminListSeckillProducts(c *gin.Context) {
	adminListManagedProducts(c, mktmodel.ActivityTypeSeckill)
}
func adminListGroupBuyProducts(c *gin.Context) {
	adminListManagedProducts(c, mktmodel.ActivityTypeGroupBuy)
}
func adminListBargainProducts(c *gin.Context) {
	adminListManagedProducts(c, mktmodel.ActivityTypeBargain)
}

func adminUpsertSeckillProducts(c *gin.Context) {
	adminUpsertManagedProducts(c, mktmodel.ActivityTypeSeckill)
}
func adminUpsertGroupBuyProducts(c *gin.Context) {
	adminUpsertManagedProducts(c, mktmodel.ActivityTypeGroupBuy)
}
func adminUpsertBargainProducts(c *gin.Context) {
	adminUpsertManagedProducts(c, mktmodel.ActivityTypeBargain)
}

func adminListActivities(c *gin.Context, actType string) {
	var q listActivitiesQuery
	_ = c.ShouldBindQuery(&q)
	list, total, err := mktsvc.ListActivities(c.Request.Context(), actType, mktsvc.ActivityListQuery{
		Page: q.Page,
		Size: q.Size,
	})
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: q.Page, Size: q.Size})
}

func adminCreateActivity(c *gin.Context, actType string) {
	var req activityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	row := &mktmodel.Activity{
		Type:       actType,
		Name:       req.Name,
		StartAt:    req.StartAt,
		EndAt:      req.EndAt,
		Status:     req.Status,
		PriceRule:  req.PriceRule,
		PriceValue: req.PriceValue,
		Exclusive:  req.Exclusive,
	}
	if err := mktsvc.CreateActivity(c.Request.Context(), actType, row); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, row)
}

func adminUpdateActivity(c *gin.Context, actType string) {
	activityID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req activityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	updates := map[string]any{
		"name":        req.Name,
		"status":      req.Status,
		"price_rule":  req.PriceRule,
		"price_value": req.PriceValue,
		"exclusive":   req.Exclusive,
		"start_at":    req.StartAt,
		"end_at":      req.EndAt,
	}
	if err := mktsvc.UpdateActivity(c.Request.Context(), activityID, actType, updates); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminListManagedProducts(c *gin.Context, actType string) {
	var q listManagedProductsQuery
	_ = c.ShouldBindQuery(&q)
	list, total, err := mktsvc.ListManagedActivityProducts(c.Request.Context(), actType, mktsvc.ActivityProductListQuery{
		ActivityID: q.ActivityID,
		ProductID:  q.ProductID,
		CategoryID: q.CategoryID,
		Keyword:    q.Keyword,
		Page:       q.Page,
		Size:       q.Size,
	})
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: q.Page, Size: q.Size})
}

func adminUpsertManagedProducts(c *gin.Context, actType string) {
	activityID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req []mktsvc.ActivityProductUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := mktsvc.UpsertActivityProducts(c.Request.Context(), activityID, actType, req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}
