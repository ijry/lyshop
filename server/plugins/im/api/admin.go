package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	immodel "github.com/ijry/lyshop/plugins/im/model"
	imsvc "github.com/ijry/lyshop/plugins/im/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/im/sessions", middleware.RequirePermission("im:view"), adminListSessions)
	g.GET("/im/sessions/:id/messages", middleware.RequirePermission("im:view"), adminListMessages)
	g.POST("/im/sessions/:id/reply", middleware.RequirePermission("im:reply"), adminReply)
	g.POST("/im/sessions/:id/accept", middleware.RequirePermission("im:reply"), adminAcceptSession)
	g.POST("/im/sessions/:id/close", middleware.RequirePermission("im:reply"), adminCloseSession)
	g.POST("/im/sessions/:id/transfer", middleware.RequirePermission("im:reply"), adminTransferSession)
	g.GET("/im/staff/status", middleware.RequirePermission("im:view"), adminGetStaffStatus)
	g.POST("/im/staff/online", middleware.RequirePermission("im:reply"), adminSetOnline)
	g.GET("/im/staff", middleware.RequirePermission("im:staff:manage"), adminListStaff)
	g.POST("/im/staff", middleware.RequirePermission("im:staff:manage"), adminCreateStaff)
	g.PUT("/im/staff/:id", middleware.RequirePermission("im:staff:manage"), adminUpdateStaff)
	g.DELETE("/im/staff/:id", middleware.RequirePermission("im:staff:manage"), adminDeleteStaff)
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
	if req.Type == "" {
		req.Type = immodel.MsgTypeText
	}
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
	imsvc.PushToUser(sessionID, msg)
	response.OK(c, msg)
}

// adminAcceptSession lets a staff manually accept a waiting session.
func adminAcceptSession(c *gin.Context) {
	sessionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	staffID, _ := c.Get("user_id")
	if err := imsvc.AcceptSession(c.Request.Context(), sessionID, staffID.(uint64)); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

// adminCloseSession closes a session and frees staff capacity.
func adminCloseSession(c *gin.Context) {
	sessionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := imsvc.CloseSession(c.Request.Context(), sessionID); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

// adminTransferSession reassigns a session to another staff member.
func adminTransferSession(c *gin.Context) {
	sessionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	fromStaffID, _ := c.Get("user_id")
	var req struct {
		ToStaffID uint64 `json:"to_staff_id" binding:"required"`
		Remark    string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := imsvc.TransferSession(c.Request.Context(), sessionID, fromStaffID.(uint64), req.ToStaffID, req.Remark); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

// adminGetStaffStatus returns the current staff's online status and load.
func adminGetStaffStatus(c *gin.Context) {
	staffID, _ := c.Get("user_id")
	status, err := imsvc.GetStaffStatus(c.Request.Context(), staffID.(uint64))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, status)
}

// adminSetOnline toggles the current staff's online status.
func adminSetOnline(c *gin.Context) {
	staffID, _ := c.Get("user_id")
	var req struct {
		Online bool `json:"online"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	imsvc.SetStaffOnline(c.Request.Context(), staffID.(uint64), req.Online)
	response.OK(c, nil)
}

// adminListStaff returns all customer service staff.
func adminListStaff(c *gin.Context) {
	list, err := imsvc.ListStaff(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

// adminCreateStaff creates a new staff record.
func adminCreateStaff(c *gin.Context) {
	var req struct {
		AdminID uint64 `json:"admin_id" binding:"required"`
		MaxLoad int    `json:"max_load"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if req.MaxLoad <= 0 {
		req.MaxLoad = 5
	}
	staff := &immodel.ImStaff{
		AdminID: req.AdminID,
		MaxLoad: req.MaxLoad,
	}
	if err := imsvc.CreateStaff(c.Request.Context(), staff); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, staff)
}

// adminUpdateStaff updates staff settings.
func adminUpdateStaff(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		MaxLoad int `json:"max_load"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := imsvc.UpdateStaff(c.Request.Context(), id, req.MaxLoad); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

// adminDeleteStaff removes a staff record.
func adminDeleteStaff(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := imsvc.DeleteStaff(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
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
