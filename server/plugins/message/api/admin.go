package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	msgsvc "github.com/ijry/lyshop/plugins/message/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/messages", middleware.RequirePermission("message:view"), adminListMessages)
	g.POST("/messages/send", middleware.RequirePermission("message:send"), adminSendMessage)
}

func adminListMessages(c *gin.Context) {
	group := c.Query("group")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := msgsvc.AdminList(c.Request.Context(), group, page, size)
	if err != nil { response.Fail(c, 500, err.Error()); return }
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminSendMessage(c *gin.Context) {
	var req struct {
		UserID  uint64 `json:"user_id"`
		Group   string `json:"group"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, err.Error()); return }
	if err := msgsvc.Send(c.Request.Context(), req.UserID, req.Group, req.Title, req.Content); err != nil {
		response.Fail(c, 500, err.Error()); return
	}
	response.OK(c, nil)
}
