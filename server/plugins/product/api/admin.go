package api

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	productsvc "github.com/ijry/lyshop/plugins/product/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/categories", middleware.RequirePermission("product:view"), adminListCategories)
	g.POST("/categories", middleware.RequirePermission("product:create"), adminCreateCategory)
	g.PUT("/categories/:id", middleware.RequirePermission("product:edit"), adminUpdateCategory)
	g.DELETE("/categories/:id", middleware.RequirePermission("product:delete"), adminDeleteCategory)

	g.GET("/products", middleware.RequirePermission("product:view"), adminListProducts)
	g.GET("/products/:id", middleware.RequirePermission("product:view"), adminGetProduct)
	g.POST("/products", middleware.RequirePermission("product:create"), adminCreateProduct)
	g.PUT("/products/:id", middleware.RequirePermission("product:edit"), adminUpdateProduct)
	g.DELETE("/products/:id", middleware.RequirePermission("product:delete"), adminDeleteProduct)
}

func adminListCategories(c *gin.Context) {
	list, err := productsvc.ListCategories(c.Request.Context(), true)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func adminCreateCategory(c *gin.Context) {
	var cat productmodel.ProductCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := productsvc.CreateCategory(c.Request.Context(), &cat); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, cat)
}

func adminUpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates map[string]any
	c.ShouldBindJSON(&updates)
	if err := productsvc.UpdateCategory(c.Request.Context(), id, updates); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminDeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := productsvc.DeleteCategory(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminListProducts(c *gin.Context) {
	var q productsvc.ProductListQuery
	c.ShouldBindQuery(&q)
	list, total, err := productsvc.ListProducts(c.Request.Context(), q, 0)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: q.Page, Size: q.Size})
}

func adminGetProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	detail, err := productsvc.GetProduct(c.Request.Context(), id, 0)
	if err != nil {
		response.Fail(c, 10002, err.Error())
		return
	}
	response.OK(c, detail)
}

func adminCreateProduct(c *gin.Context) {
	var req struct {
		Product productmodel.Product        `json:"product"`
		SKUs    []productmodel.ProductSku   `json:"skus"`
		Images  []productmodel.ProductImage `json:"images"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if len(req.Product.Detail) == 0 {
		req.Product.Detail = json.RawMessage(`{"version":1,"blocks":[]}`)
	}
	if err := productsvc.CreateProduct(c.Request.Context(), &req.Product, req.SKUs, req.Images); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, req.Product)
}

func adminUpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Product map[string]any              `json:"product"`
		SKUs    []productmodel.ProductSku   `json:"skus"`
		Images  []productmodel.ProductImage `json:"images"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	updates := req.Product
	if updates == nil {
		updates = map[string]any{}
	}
	if err := productsvc.UpdateProduct(c.Request.Context(), id, updates); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	if req.SKUs != nil {
		if err := productsvc.ReplaceProductSkus(c.Request.Context(), id, req.SKUs); err != nil {
			response.Fail(c, 500, err.Error())
			return
		}
	}
	if req.Images != nil {
		if err := productsvc.ReplaceProductImages(c.Request.Context(), id, req.Images); err != nil {
			response.Fail(c, 500, err.Error())
			return
		}
	}
	response.OK(c, nil)
}

func adminDeleteProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := productsvc.DeleteProduct(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}
