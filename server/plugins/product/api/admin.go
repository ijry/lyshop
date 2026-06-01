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

	g.GET("/spec-templates", middleware.RequirePermission("product:view"), adminListSpecTemplates)
	g.GET("/spec-templates/:id", middleware.RequirePermission("product:view"), adminGetSpecTemplate)
	g.POST("/spec-templates", middleware.RequirePermission("product:edit"), adminCreateSpecTemplate)
	g.PUT("/spec-templates/:id", middleware.RequirePermission("product:edit"), adminUpdateSpecTemplate)
	g.DELETE("/spec-templates/:id", middleware.RequirePermission("product:delete"), adminDeleteSpecTemplate)
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
		Product           productmodel.Product         `json:"product"`
		SKUs              []productmodel.ProductSku    `json:"skus"`
		Images            []productmodel.ProductImage  `json:"images"`
		SpecSchema        []productsvc.SpecSchemaGroup `json:"spec_schema"`
		SkuOverrides      []productsvc.SkuOverride     `json:"sku_overrides"`
		SkuGenerationMode string                       `json:"sku_generation_mode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if req.SkuGenerationMode == "auto" {
		generated, err := productsvc.BuildSkusFromSpecSchema(req.SpecSchema, req.Product.Price, req.SkuOverrides)
		if err != nil {
			response.Fail(c, 400, err.Error())
			return
		}
		req.SKUs = generated
	}
	if len(req.Product.Detail) == 0 {
		req.Product.Detail = json.RawMessage(`{"version":1,"blocks":[]}`)
	}
	if err := productsvc.CreateProduct(c.Request.Context(), &req.Product, req.SKUs, req.Images); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, gin.H{
		"product": req.Product,
		"sku_diff": productsvc.SkuDiffSummary{
			Added: len(req.SKUs),
		},
	})
}

func adminUpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Product           map[string]any               `json:"product"`
		SKUs              []productmodel.ProductSku    `json:"skus"`
		Images            []productmodel.ProductImage  `json:"images"`
		SpecSchema        []productsvc.SpecSchemaGroup `json:"spec_schema"`
		SkuOverrides      []productsvc.SkuOverride     `json:"sku_overrides"`
		SkuGenerationMode string                       `json:"sku_generation_mode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	updates := req.Product
	if updates == nil {
		updates = map[string]any{}
	}
	mode := req.SkuGenerationMode
	if mode == "" {
		mode = parseString(updates["sku_generation_mode"])
	}
	if mode == "auto" {
		basePrice := 0.0
		if p, ok := parseFloat64(updates["price"]); ok {
			basePrice = p
		} else {
			detail, err := productsvc.GetProduct(c.Request.Context(), id, 0)
			if err == nil && detail != nil {
				basePrice = detail.Price
			}
		}
		generated, err := productsvc.BuildSkusFromSpecSchema(req.SpecSchema, basePrice, req.SkuOverrides)
		if err != nil {
			response.Fail(c, 400, err.Error())
			return
		}
		req.SKUs = generated
	}
	if err := productsvc.UpdateProduct(c.Request.Context(), id, updates); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	var diff *productsvc.SkuDiffSummary
	if req.SKUs != nil {
		var err error
		diff, err = productsvc.ReplaceProductSkus(c.Request.Context(), id, req.SKUs)
		if err != nil {
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
	response.OK(c, gin.H{"sku_diff": diff})
}

func adminDeleteProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := productsvc.DeleteProduct(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminListSpecTemplates(c *gin.Context) {
	var q productsvc.SpecTemplateListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 || q.Size > 200 {
		q.Size = 20
	}
	list, total, err := productsvc.ListSpecTemplates(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{
		List:  list,
		Total: total,
		Page:  q.Page,
		Size:  q.Size,
	})
}

func adminGetSpecTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := productsvc.GetSpecTemplate(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.OK(c, row)
}

func adminCreateSpecTemplate(c *gin.Context) {
	var req struct {
		Name        string                       `json:"name"`
		CategoryIDs []uint64                     `json:"category_ids"`
		Attrs       []productsvc.SpecSchemaGroup `json:"attrs"`
		Status      int8                         `json:"status"`
		Sort        int                          `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	name := parseString(req.Name)
	if name == "" {
		response.Fail(c, 400, "模板名称不能为空")
		return
	}
	categoryIDs, err := productsvc.BuildSpecTemplateCategoryIDs(req.CategoryIDs)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	attrs, err := productsvc.BuildSpecTemplateAttrs(req.Attrs)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	status := int8(1)
	if req.Status == 0 {
		status = 0
	}
	row := &productmodel.SpecTemplate{
		Name:        name,
		CategoryIDs: categoryIDs,
		Attrs:       attrs,
		Status:      status,
		Sort:        req.Sort,
	}
	if err := productsvc.CreateSpecTemplate(c.Request.Context(), row); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, row)
}

func adminUpdateSpecTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Name        *string                       `json:"name"`
		CategoryIDs *[]uint64                     `json:"category_ids"`
		Attrs       *[]productsvc.SpecSchemaGroup `json:"attrs"`
		Status      *int8                         `json:"status"`
		Sort        *int                          `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = parseString(*req.Name)
	}
	if req.CategoryIDs != nil {
		updates["category_ids"] = *req.CategoryIDs
	}
	if req.Attrs != nil {
		updates["attrs"] = *req.Attrs
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if err := productsvc.UpdateSpecTemplate(c.Request.Context(), id, updates); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminDeleteSpecTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := productsvc.DeleteSpecTemplate(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func parseString(v any) string {
	switch value := v.(type) {
	case string:
		return value
	default:
		return ""
	}
}

func parseFloat64(v any) (float64, bool) {
	switch value := v.(type) {
	case float64:
		return value, true
	case float32:
		return float64(value), true
	case int:
		return float64(value), true
	case int64:
		return float64(value), true
	case json.Number:
		n, err := value.Float64()
		return n, err == nil
	default:
		return 0, false
	}
}
