package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/response"
	decorsvc "github.com/ijry/lyshop/plugins/decor/service"
)

func RegisterFrontRoutes(g *gin.RouterGroup) {
	g.GET("/index/decor", getPublishedPage)
	g.GET("/decor/:page_key", getPage)
}

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/decor/:page_key/variants", adminListVariants)
	g.POST("/decor/:page_key/copies", adminCreateVariantCopy)
	g.PUT("/decor/:page_key/variants/:variant_key", adminRenameVariant)
	g.DELETE("/decor/:page_key/variants/:variant_key", adminDeleteVariant)
	g.GET("/decor/:page_key", adminGetPage)
	g.PUT("/decor/:page_key", adminSavePage)
	g.POST("/decor/:page_key/publish", adminPublishPage)
}

func getPublishedPage(c *gin.Context) {
	page, err := decorsvc.GetPublishedPage(c.Request.Context(), 0, "index")
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, page)
}

func getPage(c *gin.Context) {
	pageKey := c.Param("page_key")
	variant := c.Query("variant")
	page, err := decorsvc.GetPage(c.Request.Context(), 0, pageKey, variant)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, page)
}

func adminGetPage(c *gin.Context) {
	pageKey := c.Param("page_key")
	variant := c.Query("variant")
	page, err := decorsvc.GetPage(c.Request.Context(), 0, pageKey, variant)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, page)
}

func adminSavePage(c *gin.Context) {
	pageKey := c.Param("page_key")
	variant := c.Query("variant")
	var req struct {
		Components json.RawMessage `json:"components"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	page, err := decorsvc.SavePage(c.Request.Context(), 0, pageKey, req.Components, variant)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, page)
}

func adminPublishPage(c *gin.Context) {
	pageKey := c.Param("page_key")
	variant := c.Query("variant")
	if err := decorsvc.PublishPage(c.Request.Context(), 0, pageKey, variant); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminListVariants(c *gin.Context) {
	pageKey := c.Param("page_key")
	rows, err := decorsvc.ListVariants(c.Request.Context(), 0, pageKey)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, rows)
}

func adminCreateVariantCopy(c *gin.Context) {
	pageKey := c.Param("page_key")
	var req struct {
		FromVariantKey string `json:"from_variant_key"`
		NewVariantKey  string `json:"new_variant_key"`
		NewVariantName string `json:"new_variant_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	row, err := decorsvc.CreateVariantCopy(
		c.Request.Context(),
		0,
		pageKey,
		req.FromVariantKey,
		req.NewVariantKey,
		req.NewVariantName,
	)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, row)
}

func adminRenameVariant(c *gin.Context) {
	pageKey := c.Param("page_key")
	variantKey := c.Param("variant_key")
	var req struct {
		VariantName string `json:"variant_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := decorsvc.RenameVariant(c.Request.Context(), 0, pageKey, variantKey, req.VariantName); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminDeleteVariant(c *gin.Context) {
	pageKey := c.Param("page_key")
	variantKey := c.Param("variant_key")
	if err := decorsvc.DeleteVariant(c.Request.Context(), 0, pageKey, variantKey); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}
