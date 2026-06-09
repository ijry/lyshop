package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	immodel "github.com/ijry/lyshop/plugins/im/model"
	imsvc "github.com/ijry/lyshop/plugins/im/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(_ *http.Request) bool { return true },
}

func RegisterFrontRoutes(g *gin.RouterGroup) {
	auth := g.Group("")
	auth.Use(middleware.RequireAuth)
	auth.GET("/im/session", getOrCreateSession)
	auth.GET("/im/messages", getMessages)
	auth.POST("/im/upload", uploadAttachment)
	auth.POST("/im/feedback", submitFeedback)
}

// RegisterWSRoute registers the WebSocket endpoint on the root engine.
func RegisterWSRoute(engine interface {
	GET(string, ...gin.HandlerFunc) gin.IRoutes
}) {
	engine.GET("/ws/im", wsHandler)
}

func getOrCreateSession(c *gin.Context) {
	userID, _ := c.Get("user_id")
	input := imsvc.SessionContextInput{
		VisitorID:       c.Query("visitor_id"),
		VisitorIP:       clientIP(c),
		VisitorLocation: c.Query("visitor_location"),
		VisitorBrowser:  c.Query("visitor_browser"),
		VisitorOS:       c.Query("visitor_os"),
		VisitorLanguage: firstNonEmpty(c.Query("visitor_language"), c.GetHeader("Accept-Language")),
		VisitorReferrer: firstNonEmpty(c.Query("visitor_referrer"), c.GetHeader("Referer")),
		VisitorURL:      c.Query("visitor_url"),
		VisitorDevice:   c.Query("visitor_device"),
	}
	if raw := strings.TrimSpace(c.Query("visitor_extra")); raw != "" {
		var extra map[string]any
		if json.Unmarshal([]byte(raw), &extra) == nil {
			input.VisitorExtra = extra
		}
	}
	session, err := imsvc.GetOrCreateSessionWithContext(c.Request.Context(), userID.(uint64), input)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, session)
}

func clientIP(c *gin.Context) string {
	if forwarded := strings.TrimSpace(c.GetHeader("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	return c.ClientIP()
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func getMessages(c *gin.Context) {
	sessionID, _ := strconv.ParseUint(c.Query("session_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))
	list, total, err := imsvc.ListMessages(c.Request.Context(), sessionID, page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func submitFeedback(c *gin.Context) {
	var req struct {
		SessionID uint64 `json:"session_id" binding:"required"`
		Rating    int8   `json:"rating"` // 1=👍 -1=👎
		Comment   string `json:"comment"`
		Query     string `json:"query"`
		Answer    string `json:"answer"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	fb := &immodel.ImFeedback{
		SessionID: req.SessionID,
		Source:    immodel.FeedbackSourceUser,
		Rating:    req.Rating,
		Comment:   req.Comment,
		Query:     req.Query,
		Answer:    req.Answer,
	}
	if err := imsvc.SaveFeedback(c.Request.Context(), fb); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func uploadAttachment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessionID, _ := strconv.ParseUint(c.PostForm("session_id"), 10, 64)
	if sessionID == 0 {
		response.Fail(c, 400, "session_id required")
		return
	}
	session, err := imsvc.GetSession(c.Request.Context(), sessionID)
	if err != nil || session.UserID != userID.(uint64) {
		response.Fail(c, 403, "无权上传到该会话")
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	info, err := imsvc.UploadAttachment(c.Request.Context(), fh)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	_ = imsvc.RecordEvent(c.Request.Context(), imsvc.EventInput{
		Event:     immodel.ImEventFileUploaded,
		SessionID: sessionID,
		UserID:    userID.(uint64),
		Source:    immodel.ImEventSourceUser,
		Success:   true,
		Extra: map[string]any{
			"name": info.Name,
			"size": info.Size,
			"type": info.MessageType,
		},
	})
	response.OK(c, info)
}

func wsHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Err(401, "missing token"))
		return
	}
	claims, err := middleware.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Err(401, "invalid token"))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	if claims.Role == "admin" {
		clientID := fmt.Sprintf("staff_%d", claims.UserID)
		imsvc.HandleWSStaff(conn, clientID, claims.UserID)
		return
	}

	clientID := fmt.Sprintf("user_%d", claims.UserID)
	s, err := imsvc.GetOrCreateSession(c.Request.Context(), claims.UserID)
	if err != nil {
		conn.Close()
		return
	}
	imsvc.HandleWS(conn, clientID, s)
}
