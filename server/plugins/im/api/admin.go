package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/response"
	immodel "github.com/ijry/lyshop/plugins/im/model"
	imsvc "github.com/ijry/lyshop/plugins/im/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/im/sessions", adminListSessions)
	g.GET("/im/sessions/:id/messages", adminListMessages)
	g.POST("/im/sessions/:id/reply", adminReply)
	g.GET("/im/auto-replies", adminListAutoReplies)
	g.POST("/im/auto-replies", adminCreateAutoReply)
}

func adminListSessions(c *gin.Context) {
	staffID, _ := strconv.ParseUint(c.Query("staff_id"), 10, 64)
	status, _ := strconv.ParseInt(c.Query("status"), 10, 8)
	list, err := imsvc.ListSessions(c.Request.Context(), staffID, int8(status))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func adminListMessages(c *gin.Context) {
	sessionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))
	list, total, err := imsvc.ListMessages(c.Request.Context(), sessionID, page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminReply(c *gin.Context) {
	sessionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	staffID, _ := c.Get("user_id")
	var req struct {
		Content string `json:"content" binding:"required"`
		Type    string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if req.Type == "" { req.Type = immodel.MsgTypeText }
	msg := &immodel.ImMessage{
		SessionID:  sessionID,
		SenderType: immodel.SenderStaff,
		SenderID:   staffID.(uint64),
		Type:       req.Type,
		Content:    req.Content,
	}
	if err := imsvc.SaveMessage(c.Request.Context(), msg); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	// Push to user via Hub
	imsvc.PushToUser(sessionID, msg)
	response.OK(c, msg)
}

func adminListAutoReplies(c *gin.Context) {
	var list []immodel.ImAutoReply
	response.OK(c, list)
}

func adminCreateAutoReply(c *gin.Context) {
	var r immodel.ImAutoReply
	if err := c.ShouldBindJSON(&r); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, r)
}
