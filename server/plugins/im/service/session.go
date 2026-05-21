package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ijry/lyshop/core/db"
	immodel "github.com/ijry/lyshop/plugins/im/model"
)

// GetOrCreateSession finds an open session for userID, or creates one.
func GetOrCreateSession(ctx context.Context, userID uint64) (*immodel.ImSession, error) {
	var session immodel.ImSession
	err := db.DB.WithContext(ctx).
		Where("user_id = ? AND status != ?", userID, immodel.SessionStatusClosed).
		First(&session).Error
	if err != nil {
		session = immodel.ImSession{
			UserID: userID, Status: immodel.SessionStatusWaiting,
		}
		err = db.DB.WithContext(ctx).Create(&session).Error
	}
	return &session, err
}

// SaveMessage persists a message and updates session last_msg.
func SaveMessage(ctx context.Context, msg *immodel.ImMessage) error {
	if err := db.DB.WithContext(ctx).Create(msg).Error; err != nil {
		return err
	}
	db.DB.WithContext(ctx).Model(&immodel.ImSession{}).Where("id = ?", msg.SessionID).
		Updates(map[string]any{"last_msg": msg.Content, "unread_count": db.DB.Raw("unread_count + 1")})
	return nil
}

// GetUnread returns all unread messages for a session (for offline replay).
func GetUnread(ctx context.Context, sessionID uint64) ([]immodel.ImMessage, error) {
	var list []immodel.ImMessage
	err := db.DB.WithContext(ctx).
		Where("session_id = ? AND is_read = 0", sessionID).
		Order("id asc").Find(&list).Error
	return list, err
}

// MarkRead marks all messages in a session as read.
func MarkRead(ctx context.Context, sessionID uint64) error {
	return db.DB.WithContext(ctx).Model(&immodel.ImMessage{}).
		Where("session_id = ? AND is_read = 0", sessionID).
		Update("is_read", 1).Error
}

// ListSessions returns sessions for admin seat view.
func ListSessions(ctx context.Context, staffID uint64, status int8) ([]immodel.ImSession, error) {
	tx := db.DB.WithContext(ctx)
	if staffID > 0 { tx = tx.Where("staff_id = ?", staffID) }
	if status > 0 { tx = tx.Where("status = ?", status) }
	var list []immodel.ImSession
	err := tx.Order("updated_at desc").Limit(50).Find(&list).Error
	return list, err
}

// ListMessages returns message history for a session.
func ListMessages(ctx context.Context, sessionID uint64, page, size int) ([]immodel.ImMessage, int64, error) {
	if page <= 0 { page = 1 }
	if size <= 0 || size > 100 { size = 50 }
	var total int64
	db.DB.WithContext(ctx).Model(&immodel.ImMessage{}).Where("session_id = ?", sessionID).Count(&total)
	var list []immodel.ImMessage
	err := db.DB.WithContext(ctx).Where("session_id = ?", sessionID).
		Order("id desc").Offset((page-1)*size).Limit(size).Find(&list).Error
	return list, total, err
}

// CheckAutoReply looks for a matching auto-reply rule and returns reply text.
func CheckAutoReply(ctx context.Context, content string) string {
	var rules []immodel.ImAutoReply
	db.DB.WithContext(ctx).Where("status = 1").Order("sort asc").Find(&rules)
	for _, r := range rules {
		switch r.MatchType {
		case 1: // exact
			if content == r.Keyword { return r.Reply }
		case 2: // contains
			if strings.Contains(content, r.Keyword) { return r.Reply }
		}
	}
	return ""
}

// HandleWSStaff manages a staff WebSocket connection (receive-only from hub).
func HandleWSStaff(conn *websocket.Conn, clientID string, _ uint64) {
	client := &Client{ID: clientID, Conn: conn, Send: make(chan []byte, 64)}
	GlobalHub.Register(client)
	defer GlobalHub.Unregister(client)

	go func() {
		for data := range client.Send {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				break
			}
		}
	}()

	conn.SetReadLimit(4096)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(_ string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	}
}

// PushToUser sends a staff message to the user of the given session via Hub.
func PushToUser(sessionID uint64, msg *immodel.ImMessage) {
	var session immodel.ImSession
	if err := db.DB.First(&session, sessionID).Error; err != nil {
		return
	}
	userKey := fmt.Sprintf("user_%d", session.UserID)
	frame := Frame{Type: "msg", SessionID: sessionID, Payload: map[string]any{
		"msg_id":      fmt.Sprintf("%d", msg.ID),
		"msg_type":    msg.Type,
		"content":     msg.Content,
		"sender_type": msg.SenderType,
	}}
	data, _ := json.Marshal(frame)
	GlobalHub.Send(userKey, data)
}

// HandleWS manages one WebSocket client lifecycle (readPump + writePump).
func HandleWS(conn *websocket.Conn, clientID string, session *immodel.ImSession) {
	client := &Client{ID: clientID, Conn: conn, Send: make(chan []byte, 64)}
	GlobalHub.Register(client)
	defer GlobalHub.Unregister(client)

	// Replay unread messages on connect
	ctx := context.Background()
	unread, _ := GetUnread(ctx, session.ID)
	for _, msg := range unread {
		frame := Frame{Type: "msg", SessionID: session.ID, Payload: map[string]any{
			"msg_id":      fmt.Sprintf("%d", msg.ID),
			"msg_type":    msg.Type,
			"content":     msg.Content,
			"sender_type": msg.SenderType,
		}}
		data, _ := json.Marshal(frame)
		client.Send <- data
	}
	MarkRead(ctx, session.ID)

	// writePump: drain Send channel to connection
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case data, ok := <-client.Send:
				if !ok {
					conn.WriteMessage(websocket.CloseMessage, nil)
					return
				}
				conn.WriteMessage(websocket.TextMessage, data)
			case <-ticker.C:
				conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()

	// readPump: receive messages from client
	conn.SetReadLimit(4096)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(_ string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil { break }
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		var frame Frame
		if err := json.Unmarshal(raw, &frame); err != nil { continue }

		switch frame.Type {
		case "msg":
			content, _ := frame.Payload["content"].(string)
			msgType, _ := frame.Payload["msg_type"].(string)
			if msgType == "" { msgType = immodel.MsgTypeText }

			msg := &immodel.ImMessage{
				SessionID:  session.ID,
				SenderType: immodel.SenderUser,
				SenderID:   session.UserID,
				Type:       msgType,
				Content:    content,
			}
			SaveMessage(ctx, msg)

			// Forward to staff if online
			if session.StaffID > 0 {
				staffKey := fmt.Sprintf("staff_%d", session.StaffID)
				GlobalHub.Send(staffKey, raw)
			}

			// Auto-reply
			if reply := CheckAutoReply(ctx, content); reply != "" {
				autoMsg := &immodel.ImMessage{
					SessionID:  session.ID,
					SenderType: immodel.SenderStaff,
					SenderID:   0,
					Type:       immodel.MsgTypeText,
					Content:    reply,
				}
				SaveMessage(ctx, autoMsg)
				replyFrame := Frame{Type: "msg", SessionID: session.ID, Payload: map[string]any{
					"msg_type": immodel.MsgTypeText, "content": reply,
				}}
				data, _ := json.Marshal(replyFrame)
				client.Send <- data
			}

		case "ping":
			pong := Frame{Type: "pong"}
			data, _ := json.Marshal(pong)
			client.Send <- data
		}
	}
}
