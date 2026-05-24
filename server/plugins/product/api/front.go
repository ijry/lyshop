package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/response"
	productsvc "github.com/ijry/lyshop/plugins/product/service"
)

func RegisterFrontRoutes(g *gin.RouterGroup) {
	g.GET("/categories", listCategories)
	g.GET("/products", listProducts)
	g.GET("/products/recommend", listRecommendProducts)
	g.GET("/products/:id", getProduct)
	g.GET("/products/:id/reviews", listProductReviews)
}

func listCategories(c *gin.Context) {
	list, err := productsvc.ListCategories(c.Request.Context(), false)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func listProducts(c *gin.Context) {
	var q productsvc.ProductListQuery
	c.ShouldBindQuery(&q)
	// Front-end only sees on-shelf products
	list, total, err := productsvc.ListProducts(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: q.Page, Size: q.Size})
}

func getProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	detail, err := productsvc.GetProduct(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, 10002, err.Error())
		return
	}
	response.OK(c, detail)
}

func listRecommendProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "8"))
	if limit <= 0 {
		limit = 8
	}
	if limit > 50 {
		limit = 50
	}
	list, err := productsvc.ListRecommendProducts(c.Request.Context(), limit)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func listProductReviews(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	data, err := productsvc.ListProductReviews(c.Request.Context(), id, page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, data)
}
