package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	imsvc "github.com/ijry/lyshop/plugins/im/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:    func(_ *http.Request) bool { return true },
}

func RegisterFrontRoutes(g *gin.RouterGroup) {
	auth := g.Group("")
	auth.Use(middleware.RequireAuth)
	auth.GET("/im/session", getOrCreateSession)
	auth.GET("/im/messages", getMessages)
}

// RegisterWSRoute registers the WebSocket endpoint on the root engine.
// Call this separately with the raw *gin.Engine.
func RegisterWSRoute(engine interface {
	GET(string, ...gin.HandlerFunc) gin.IRoutes
}) {
	engine.GET("/ws/im", wsHandler)
}

func getOrCreateSession(c *gin.Context) {
	userID, _ := c.Get("user_id")
	session, err := imsvc.GetOrCreateSession(c.Request.Context(), userID.(uint64))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, session)
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

func wsHandler(c *gin.Context) {
	// Extract JWT from query param
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

	var clientID string
	var session interface {
		GetID() uint64
	}

	if claims.Role == "admin" {
		clientID = fmt.Sprintf("staff_%d", claims.UserID)
		// Staff: get session_id from query
		sessionID, _ := strconv.ParseUint(c.Query("session_id"), 10, 64)
		s, _ := imsvc.GetOrCreateSession(c.Request.Context(), sessionID)
		_ = strings.TrimSpace(clientID) // suppress unused warning
		imsvc.HandleWSStaff(conn, clientID, s.ID)
		return
	}

	clientID = fmt.Sprintf("user_%d", claims.UserID)
	s, err := imsvc.GetOrCreateSession(c.Request.Context(), claims.UserID)
	if err != nil {
		conn.Close()
		return
	}
	_ = session
	imsvc.HandleWS(conn, clientID, s)
}
