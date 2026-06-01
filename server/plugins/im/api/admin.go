package api

import (
	"io"
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

	// AI 知识库（RAG）管理
	g.GET("/im/knowledge", middleware.RequirePermission("im:knowledge"), adminListKnowledge)
	g.POST("/im/knowledge", middleware.RequirePermission("im:knowledge"), adminCreateKnowledge)
	g.PUT("/im/knowledge/:id", middleware.RequirePermission("im:knowledge"), adminUpdateKnowledge)
	g.DELETE("/im/knowledge/:id", middleware.RequirePermission("im:knowledge"), adminDeleteKnowledge)
	g.POST("/im/knowledge/reindex", middleware.RequirePermission("im:knowledge"), adminReindexKnowledge)
	g.POST("/im/knowledge/import", middleware.RequirePermission("im:knowledge"), adminImportKnowledge)
	// 本地大模型连通性测试
	g.POST("/im/ai/test", middleware.RequirePermission("im:knowledge"), adminTestAI)
}

func adminListKnowledge(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := imsvc.ListKnowledge(c.Request.Context(), keyword, page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminCreateKnowledge(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Tags    string `json:"tags"`
		Sort    int    `json:"sort"`
		Status  int8   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if req.Status == 0 {
		req.Status = 1
	}
	k := &immodel.ImKnowledge{
		Title: req.Title, Content: req.Content, Tags: req.Tags, Sort: req.Sort, Status: req.Status,
	}
	if err := imsvc.CreateKnowledge(c.Request.Context(), k); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, k)
}

func adminUpdateKnowledge(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
		Tags    *string `json:"tags"`
		Sort    *int    `json:"sort"`
		Status  *int8   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	fields := map[string]any{}
	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.Content != nil {
		fields["content"] = *req.Content
	}
	if req.Tags != nil {
		fields["tags"] = *req.Tags
	}
	if req.Sort != nil {
		fields["sort"] = *req.Sort
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if err := imsvc.UpdateKnowledge(c.Request.Context(), id, fields); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminDeleteKnowledge(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := imsvc.DeleteKnowledge(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminReindexKnowledge(c *gin.Context) {
	done, err := imsvc.ReindexKnowledge(c.Request.Context())
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"indexed": done})
}

// adminImportKnowledge accepts a multipart "file" upload (txt/md/csv/html/docx/
// pdf/xlsx/…), extracts and slices its text, and creates one knowledge entry
// per chunk. Optional form fields: title, tags, chunk_size, overlap.
func adminImportKnowledge(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	const maxSize = 20 << 20 // 20MB
	if fh.Size > maxSize {
		response.Fail(c, 400, "文件过大，最大支持 20MB")
		return
	}
	f, err := fh.Open()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	chunkSize, _ := strconv.Atoi(c.PostForm("chunk_size"))
	overlap, _ := strconv.Atoi(c.PostForm("overlap"))
	res, err := imsvc.ImportDocument(c.Request.Context(), fh.Filename, data,
		c.PostForm("title"), c.PostForm("tags"), chunkSize, overlap)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, res)
}

func adminTestAI(c *gin.Context) {
	reply, err := imsvc.TestAIConnection(c.Request.Context())
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"reply": reply})
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
