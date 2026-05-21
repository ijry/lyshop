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
	g.GET("/decor/:page_key", adminGetPage)
	g.PUT("/decor/:page_key", adminSavePage)
	g.POST("/decor/:page_key/publish", adminPublishPage)
}

func getPublishedPage(c *gin.Context) {
	page, err := decorsvc.GetPage(c.Request.Context(), 0, "index")
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, page)
}

func getPage(c *gin.Context) {
	pageKey := c.Param("page_key")
	page, err := decorsvc.GetPage(c.Request.Context(), 0, pageKey)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, page)
}

func adminGetPage(c *gin.Context) {
	pageKey := c.Param("page_key")
	page, err := decorsvc.GetPage(c.Request.Context(), 0, pageKey)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, page)
}

func adminSavePage(c *gin.Context) {
	pageKey := c.Param("page_key")
	var req struct {
		Components json.RawMessage `json:"components"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	page, err := decorsvc.SavePage(c.Request.Context(), 0, pageKey, req.Components)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, page)
}

func adminPublishPage(c *gin.Context) {
	pageKey := c.Param("page_key")
	if err := decorsvc.PublishPage(c.Request.Context(), 0, pageKey); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}
