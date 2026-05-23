package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	msgsvc "github.com/ijry/lyshop/plugins/message/service"
)

func RegisterFrontRoutes(g *gin.RouterGroup) {
	auth := g.Group("").Use(middleware.RequireAuth)
	auth.GET("/messages", listMessages)
	auth.GET("/messages/unread", unreadCounts)
	auth.POST("/messages/read", markRead)
}

func listMessages(c *gin.Context) {
	userID, _ := c.Get("user_id")
	group := c.Query("group")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := msgsvc.ListByGroup(c.Request.Context(), userID.(uint64), group, page, size)
	if err != nil { response.Fail(c, 500, err.Error()); return }
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func unreadCounts(c *gin.Context) {
	userID, _ := c.Get("user_id")
	counts, err := msgsvc.UnreadCounts(c.Request.Context(), userID.(uint64))
	if err != nil { response.Fail(c, 500, err.Error()); return }
	response.OK(c, counts)
}

func markRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct { IDs []uint64 `json:"ids"` }
	c.ShouldBindJSON(&req)
	msgsvc.MarkRead(c.Request.Context(), userID.(uint64), req.IDs)
	response.OK(c, nil)
}
