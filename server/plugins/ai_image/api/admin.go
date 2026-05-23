package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	aimodel "github.com/ijry/lyshop/plugins/ai_image/model"
	aisvc "github.com/ijry/lyshop/plugins/ai_image/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/ai/models", middleware.RequirePermission("ai:view"), listModels)
	g.POST("/ai/models", middleware.RequirePermission("ai:config"), createModel)
	g.POST("/ai/generate", middleware.RequirePermission("ai:generate"), generate)
	g.GET("/ai/tasks/:id", middleware.RequirePermission("ai:view"), getTask)
	g.GET("/ai/tasks", middleware.RequirePermission("ai:view"), listTasks)
}

func listModels(c *gin.Context) {
	list, err := aisvc.ListModels(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func createModel(c *gin.Context) {
	var m aimodel.AiModel
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := aisvc.CreateModel(c.Request.Context(), &m); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, m)
}

func generate(c *gin.Context) {
	var req struct {
		ModelID   uint64         `json:"model_id"`
		Scene     string         `json:"scene"`   // carousel|detail
		BizType   string         `json:"biz_type"` // cover|gallery|detail|intro
		TargetProductID uint64   `json:"target_product_id"`
		RefImageURL string       `json:"ref_image_url"`
		Prompt    string         `json:"prompt"   binding:"required"`
		NegPrompt string         `json:"neg_prompt"`
		Params    map[string]any `json:"params"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if req.BizType == "" {
		req.BizType = "detail"
	}
	task, err := aisvc.Generate(c.Request.Context(), req.ModelID, req.Scene, req.BizType, req.Prompt, req.NegPrompt, req.RefImageURL, req.TargetProductID, req.Params)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, task)
}

func getTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	task, err := aisvc.GetTask(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, 404, "任务不存在")
		return
	}
	response.OK(c, task)
}

func listTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := aisvc.ListTasks(c.Request.Context(), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}
