package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
	mktsvc "github.com/ijry/lyshop/plugins/marketing/service"
)

func RegisterFrontRoutes(g *gin.RouterGroup) {
	g.GET("/marketing/seckill/products", listSeckillProducts)
	g.GET("/marketing/group-buy/products", listGroupBuyProducts)
	g.GET("/marketing/bargain/products", listBargainProducts)
	g.GET("/marketing/activity-products/:id", getFrontActivityProductDetail)

	auth := g.Group("")
	auth.Use(middleware.RequireAuth)
	auth.GET("/coupons", listClaimableCoupons)
	auth.POST("/coupons/:id/claim", claimCoupon)
	auth.GET("/user/coupons", listUserCoupons)
	auth.GET("/user/points/logs", listPointsLogs)
}

func getFrontActivityProductDetail(c *gin.Context) {
	activityProductID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || activityProductID == 0 {
		response.Fail(c, 400, "invalid activity product id")
		return
	}
	detail, err := mktsvc.GetFrontActivityProductDetail(c.Request.Context(), activityProductID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, detail)
}

type activityProductsQuery struct {
	ActivityID uint64  `form:"activity_id"`
	CategoryID uint64  `form:"category_id"`
	Keyword    string  `form:"keyword"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	SortBy     string  `form:"sort_by"`
	SortOrder  string  `form:"sort_order"`
	Page       int     `form:"page"`
	Size       int     `form:"size"`
}

func listSeckillProducts(c *gin.Context)  { listActivityProducts(c, mktmodel.ActivityTypeSeckill) }
func listGroupBuyProducts(c *gin.Context) { listActivityProducts(c, mktmodel.ActivityTypeGroupBuy) }
func listBargainProducts(c *gin.Context)  { listActivityProducts(c, mktmodel.ActivityTypeBargain) }

func listActivityProducts(c *gin.Context, actType string) {
	var q activityProductsQuery
	_ = c.ShouldBindQuery(&q)
	page := q.Page
	size := q.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	list, total, err := mktsvc.ListFrontActivityProducts(c.Request.Context(), actType, mktsvc.ActivityProductListQuery{
		ActivityID: q.ActivityID,
		CategoryID: q.CategoryID,
		Keyword:    q.Keyword,
		MinPrice:   q.MinPrice,
		MaxPrice:   q.MaxPrice,
		SortBy:     q.SortBy,
		SortOrder:  q.SortOrder,
		Page:       page,
		Size:       size,
	})
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
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
