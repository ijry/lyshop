package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	productsvc "github.com/ijry/lyshop/plugins/product/service"
)

func RegisterFrontRoutes(g *gin.RouterGroup) {
	g.GET("/categories", listCategories)
	g.GET("/products", listProducts)
	g.GET("/products/recommend", listRecommendProducts)
	g.GET("/products/:id", getProduct)
	g.GET("/products/:id/reviews", listProductReviews)

	auth := g.Group("")
	auth.Use(middleware.RequireAuth)
	auth.POST("/products/:id/favorite", favoriteProduct)
	auth.DELETE("/products/:id/favorite", unfavoriteProduct)
	auth.GET("/user/favorites", listUserFavorites)
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
	list, total, err := productsvc.ListProducts(c.Request.Context(), q, currentUserID(c))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: q.Page, Size: q.Size})
}

func getProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	detail, err := productsvc.GetProduct(c.Request.Context(), id, currentUserID(c))
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

func favoriteProduct(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := productsvc.FavoriteProduct(c.Request.Context(), userID.(uint64), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func unfavoriteProduct(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := productsvc.UnfavoriteProduct(c.Request.Context(), userID.(uint64), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func listUserFavorites(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := productsvc.ListUserFavorites(c.Request.Context(), userID.(uint64), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func currentUserID(c *gin.Context) uint64 {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return 0
	}
	claims, err := middleware.ParseToken(strings.TrimPrefix(auth, "Bearer "))
	if err != nil {
		return 0
	}
	return claims.UserID
}
