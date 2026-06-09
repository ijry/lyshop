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

// ListStaff returns all customer service staff.
func ListStaff(ctx context.Context) ([]immodel.ImStaff, error) {
	var list []immodel.ImStaff
	err := db.DB.WithContext(ctx).Order("id asc").Find(&list).Error
	return list, err
}

func GetSession(ctx context.Context, sessionID uint64) (*immodel.ImSession, error) {
	var session immodel.ImSession
	err := db.DB.WithContext(ctx).First(&session, sessionID).Error
	return &session, err
}

type SessionContextInput struct {
	VisitorID       string         `json:"visitor_id"`
	VisitorIP       string         `json:"visitor_ip"`
	VisitorLocation string         `json:"visitor_location"`
	VisitorBrowser  string         `json:"visitor_browser"`
	VisitorOS       string         `json:"visitor_os"`
	VisitorLanguage string         `json:"visitor_language"`
	VisitorReferrer string         `json:"visitor_referrer"`
	VisitorURL      string         `json:"visitor_url"`
	VisitorDevice   string         `json:"visitor_device"`
	VisitorExtra    map[string]any `json:"visitor_extra"`
}

func (input SessionContextInput) hasValues() bool {
	return strings.TrimSpace(input.VisitorID) != "" ||
		strings.TrimSpace(input.VisitorIP) != "" ||
		strings.TrimSpace(input.VisitorLocation) != "" ||
		strings.TrimSpace(input.VisitorBrowser) != "" ||
		strings.TrimSpace(input.VisitorOS) != "" ||
		strings.TrimSpace(input.VisitorLanguage) != "" ||
		strings.TrimSpace(input.VisitorReferrer) != "" ||
		strings.TrimSpace(input.VisitorURL) != "" ||
		strings.TrimSpace(input.VisitorDevice) != "" ||
		len(input.VisitorExtra) > 0
}

func trimField(value string, max int) string {
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) > max {
		return string(runes[:max])
	}
	return value
}

func sessionContextUpdates(input SessionContextInput) map[string]any {
	updates := map[string]any{}
	if v := trimField(input.VisitorID, 128); v != "" {
		updates["visitor_id"] = v
	}
	if v := trimField(input.VisitorIP, 64); v != "" {
		updates["visitor_ip"] = v
	}
	if v := trimField(input.VisitorLocation, 128); v != "" {
		updates["visitor_location"] = v
	}
	if v := trimField(input.VisitorBrowser, 64); v != "" {
		updates["visitor_browser"] = v
	}
	if v := trimField(input.VisitorOS, 64); v != "" {
		updates["visitor_os"] = v
	}
	if v := trimField(input.VisitorLanguage, 32); v != "" {
		updates["visitor_language"] = v
	}
	if v := trimField(input.VisitorReferrer, 512); v != "" {
		updates["visitor_referrer"] = v
	}
	if v := trimField(input.VisitorURL, 512); v != "" {
		updates["visitor_url"] = v
	}
	if v := trimField(input.VisitorDevice, 64); v != "" {
		updates["visitor_device"] = v
	}
	if len(input.VisitorExtra) > 0 {
		raw, _ := json.Marshal(input.VisitorExtra)
		updates["visitor_extra"] = string(raw)
	}
	return updates
}

func applySessionContext(ctx context.Context, session *immodel.ImSession, input SessionContextInput) {
	if !input.hasValues() {
		return
	}
	updates := sessionContextUpdates(input)
	if len(updates) == 0 {
		return
	}
	db.DB.WithContext(ctx).Model(session).Updates(updates)
	if v, ok := updates["visitor_id"].(string); ok {
		session.VisitorID = v
	}
	if v, ok := updates["visitor_ip"].(string); ok {
		session.VisitorIP = v
	}
	if v, ok := updates["visitor_location"].(string); ok {
		session.VisitorLocation = v
	}
	if v, ok := updates["visitor_browser"].(string); ok {
		session.VisitorBrowser = v
	}
	if v, ok := updates["visitor_os"].(string); ok {
		session.VisitorOS = v
	}
	if v, ok := updates["visitor_language"].(string); ok {
		session.VisitorLanguage = v
	}
	if v, ok := updates["visitor_referrer"].(string); ok {
		session.VisitorReferrer = v
	}
	if v, ok := updates["visitor_url"].(string); ok {
		session.VisitorURL = v
	}
	if v, ok := updates["visitor_device"].(string); ok {
		session.VisitorDevice = v
	}
	if v, ok := updates["visitor_extra"].(string); ok {
		session.VisitorExtra = v
	}
}

func visitorPayload(session *immodel.ImSession) map[string]any {
	if session == nil {
		return nil
	}
	payload := map[string]any{}
	add := func(key, value string) {
		if strings.TrimSpace(value) != "" {
			payload[key] = value
		}
	}
	add("visitor_id", session.VisitorID)
	add("ip", session.VisitorIP)
	add("location", session.VisitorLocation)
	add("browser", session.VisitorBrowser)
	add("os", session.VisitorOS)
	add("language", session.VisitorLanguage)
	add("referrer", session.VisitorReferrer)
	add("url", session.VisitorURL)
	add("device", session.VisitorDevice)
	if extra := messageExtraPayload(session.VisitorExtra); extra != nil {
		payload["extra"] = extra
	}
	if len(payload) == 0 {
		return nil
	}
	return payload
}

func recordEventBestEffort(ctx context.Context, input EventInput) {
	_ = RecordEvent(ctx, input)
}

func eventSourceFromSender(senderType int8) string {
	switch senderType {
	case immodel.SenderUser:
		return immodel.ImEventSourceUser
	case immodel.SenderStaff:
		return immodel.ImEventSourceStaff
	case immodel.SenderAI:
		return immodel.ImEventSourceAI
	default:
		return immodel.ImEventSourceSystem
	}
}

func messageExtraPayload(extra string) any {
	if strings.TrimSpace(extra) == "" {
		return nil
	}
	var parsed any
	if json.Unmarshal([]byte(extra), &parsed) == nil {
		return parsed
	}
	return extra
}

func normalizeMessageExtra(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(v)
	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(raw)
	}
}

func truncateDraft(value string) string {
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) > 500 {
		return string(runes[:500])
	}
	return value
}

func forwardTypingFrame(session immodel.ImSession, senderType int8, senderID uint64, draft string, stop bool) {
	targetID := ""
	if senderType == immodel.SenderUser {
		if session.StaffID == 0 {
			return
		}
		targetID = fmt.Sprintf("staff_%d", session.StaffID)
	} else if senderType == immodel.SenderStaff {
		targetID = fmt.Sprintf("user_%d", session.UserID)
	} else {
		return
	}
	frameType := "typing_draft"
	if stop {
		frameType = "typing_stop"
	}
	frame := Frame{Type: frameType, SessionID: session.ID, Payload: map[string]any{
		"sender_type": senderType,
		"sender_id":   senderID,
		"updated_at":  time.Now().UnixMilli(),
	}}
	if !stop {
		frame.Payload["draft"] = truncateDraft(draft)
	}
	GlobalHub.Send(targetID, mustMarshal(frame))
}

// CreateStaff creates a new staff record.
func CreateStaff(ctx context.Context, staff *immodel.ImStaff) error {
	return db.DB.WithContext(ctx).Create(staff).Error
}

// UpdateStaff updates staff max_load.
func UpdateStaff(ctx context.Context, id uint64, maxLoad int) error {
	return db.DB.WithContext(ctx).Model(&immodel.ImStaff{}).
		Where("id = ?", id).Update("max_load", maxLoad).Error
}

// DeleteStaff removes a staff record.
func DeleteStaff(ctx context.Context, id uint64) error {
	return db.DB.WithContext(ctx).Delete(&immodel.ImStaff{}, id).Error
}

// GetOrCreateSession finds an open session for userID, or creates one.
//
// When AI first-line service is enabled, a new session starts in "ai" mode and
// is served by the local LLM immediately (no queue). The user can switch to a
// human agent later by typing a human-request keyword (see SwitchToHuman).
//
// When AI is disabled, a new session falls back to the classic human flow:
// assign an available staff immediately, otherwise queue.
func GetOrCreateSession(ctx context.Context, userID uint64) (*immodel.ImSession, error) {
	return GetOrCreateSessionWithContext(ctx, userID, SessionContextInput{})
}

func GetOrCreateSessionWithContext(ctx context.Context, userID uint64, input SessionContextInput) (*immodel.ImSession, error) {
	var session immodel.ImSession
	err := db.DB.WithContext(ctx).
		Where("user_id = ? AND status != ?", userID, immodel.SessionStatusClosed).
		First(&session).Error
	if err == nil {
		applySessionContext(ctx, &session, input)
		return &session, nil
	}

	if AIEnabled() {
		session = immodel.ImSession{
			UserID: userID,
			Mode:   immodel.SessionModeAI,
			Status: immodel.SessionStatusOngoing, // served by AI
		}
		if err = db.DB.WithContext(ctx).Create(&session).Error; err != nil {
			return &session, err
		}
		applySessionContext(ctx, &session, input)
		recordEventBestEffort(ctx, EventInput{
			Event:     immodel.ImEventSessionCreated,
			SessionID: session.ID,
			UserID:    userID,
			Source:    immodel.ImEventSourceUser,
			Success:   true,
			Extra:     map[string]any{"mode": session.Mode},
		})
		return &session, nil
	}

	// Classic human flow.
	session = immodel.ImSession{
		UserID: userID,
		Mode:   immodel.SessionModeHuman,
		Status: immodel.SessionStatusWaiting,
	}
	if err = db.DB.WithContext(ctx).Create(&session).Error; err != nil {
		return &session, err
	}
	applySessionContext(ctx, &session, input)
	recordEventBestEffort(ctx, EventInput{
		Event:     immodel.ImEventSessionCreated,
		SessionID: session.ID,
		UserID:    userID,
		Source:    immodel.ImEventSourceUser,
		Success:   true,
		Extra:     map[string]any{"mode": session.Mode},
	})
	if staffID := pickAvailableStaff(ctx); staffID > 0 {
		assignSession(ctx, &session, staffID)
	} else {
		pos := countWaiting(ctx)
		db.DB.WithContext(ctx).Model(&session).Update("queue_position", pos)
		session.QueuePosition = pos
	}
	return &session, nil
}

// SwitchToHuman transitions an AI session to the human queue/assign flow.
// It records a system message, assigns an available staff or enqueues the user,
// and notifies the user of the result. Idempotent: a session already in human
// mode is left unchanged.
func SwitchToHuman(ctx context.Context, sessionID uint64) error {
	var session immodel.ImSession
	if err := db.DB.WithContext(ctx).First(&session, sessionID).Error; err != nil {
		return err
	}
	if session.Mode == immodel.SessionModeHuman {
		return nil
	}

	// System notice into the transcript.
	sysMsg := &immodel.ImMessage{
		SessionID:  sessionID,
		SenderType: immodel.SenderSystem,
		Type:       immodel.MsgTypeSystem,
		Content:    "正在为您转接人工客服…",
	}
	SaveMessage(ctx, sysMsg)
	sysData, _ := json.Marshal(Frame{Type: "msg", SessionID: sessionID, Payload: map[string]any{
		"msg_type": immodel.MsgTypeSystem, "content": sysMsg.Content, "sender_type": immodel.SenderSystem,
	}})
	GlobalHub.Send(fmt.Sprintf("user_%d", session.UserID), sysData)

	db.DB.WithContext(ctx).Model(&session).Update("mode", immodel.SessionModeHuman)
	session.Mode = immodel.SessionModeHuman
	recordEventBestEffort(ctx, EventInput{
		Event:     immodel.ImEventToHuman,
		SessionID: sessionID,
		UserID:    session.UserID,
		StaffID:   session.StaffID,
		Source:    immodel.ImEventSourceUser,
		Success:   true,
	})

	if staffID := pickAvailableStaff(ctx); staffID > 0 {
		assignSession(ctx, &session, staffID)
		// Tell the user a human accepted.
		acc, _ := json.Marshal(Frame{Type: "assign", SessionID: sessionID, Payload: map[string]any{"action": "accepted"}})
		GlobalHub.Send(fmt.Sprintf("user_%d", session.UserID), acc)
	} else {
		db.DB.WithContext(ctx).Model(&session).Update("status", immodel.SessionStatusWaiting)
		session.Status = immodel.SessionStatusWaiting
		pos := countWaiting(ctx)
		db.DB.WithContext(ctx).Model(&session).Update("queue_position", pos)
		session.QueuePosition = pos
		qf, _ := json.Marshal(Frame{Type: "queue", SessionID: sessionID, Payload: map[string]any{"position": pos}})
		GlobalHub.Send(fmt.Sprintf("user_%d", session.UserID), qf)
	}
	return nil
}

// pickAvailableStaff returns the staff_id with lowest load that is online and not full.
func pickAvailableStaff(ctx context.Context) uint64 {
	var staff immodel.ImStaff
	err := db.DB.WithContext(ctx).
		Where("is_online = 1 AND current_load < max_load").
		Order("current_load asc").
		First(&staff).Error
	if err != nil {
		return 0
	}
	return staff.AdminID
}

// assignSession sets staff on a session and increments staff load.
func assignSession(ctx context.Context, session *immodel.ImSession, staffID uint64) {
	db.DB.WithContext(ctx).Model(session).Updates(map[string]any{
		"staff_id":       staffID,
		"status":         immodel.SessionStatusOngoing,
		"queue_position": 0,
	})
	session.StaffID = staffID
	session.Status = immodel.SessionStatusOngoing
	session.QueuePosition = 0
	db.DB.WithContext(ctx).Model(&immodel.ImStaff{}).
		Where("admin_id = ?", staffID).
		UpdateColumn("current_load", db.DB.Raw("current_load + 1"))

	// Push assign frame to staff
	frame := Frame{Type: "assign", SessionID: session.ID, Payload: map[string]any{
		"action": "new", "user_id": session.UserID,
	}}
	if visitor := visitorPayload(session); visitor != nil {
		frame.Payload["visitor"] = visitor
	}
	data, _ := json.Marshal(frame)
	GlobalHub.Send(fmt.Sprintf("staff_%d", staffID), data)
}

// countWaiting returns the number of sessions currently waiting (for queue position).
func countWaiting(ctx context.Context) int {
	var count int64
	db.DB.WithContext(ctx).Model(&immodel.ImSession{}).
		Where("status = ?", immodel.SessionStatusWaiting).Count(&count)
	return int(count)
}

// DrainQueue tries to assign waiting sessions to newly available staff.
// Call this when a staff comes online or closes a session.
func DrainQueue(ctx context.Context) {
	var waiting []immodel.ImSession
	db.DB.WithContext(ctx).
		Where("status = ?", immodel.SessionStatusWaiting).
		Order("queue_position asc, created_at asc").
		Find(&waiting)

	for _, s := range waiting {
		staffID := pickAvailableStaff(ctx)
		if staffID == 0 {
			break // no more available staff
		}
		sess := s
		assignSession(ctx, &sess, staffID)

		// Notify user of their position change (now assigned)
		queueFrame := Frame{Type: "assign", SessionID: sess.ID, Payload: map[string]any{
			"action": "accepted",
		}}
		data, _ := json.Marshal(queueFrame)
		GlobalHub.Send(fmt.Sprintf("user_%d", sess.UserID), data)
	}

	// Recompute queue positions for remaining waiters
	var stillWaiting []immodel.ImSession
	db.DB.WithContext(ctx).
		Where("status = ?", immodel.SessionStatusWaiting).
		Order("created_at asc").Find(&stillWaiting)
	for i, s := range stillWaiting {
		pos := i + 1
		db.DB.WithContext(ctx).Model(&immodel.ImSession{}).
			Where("id = ?", s.ID).Update("queue_position", pos)
		// Notify user of updated position
		posFrame := Frame{Type: "queue", SessionID: s.ID, Payload: map[string]any{
			"position": pos,
		}}
		data, _ := json.Marshal(posFrame)
		GlobalHub.Send(fmt.Sprintf("user_%d", s.UserID), data)
	}
}

// GetStaffStatus returns the online status and load for a staff member.
func GetStaffStatus(ctx context.Context, adminID uint64) (*immodel.ImStaff, error) {
	var staff immodel.ImStaff
	err := db.DB.WithContext(ctx).Where("admin_id = ?", adminID).First(&staff).Error
	if err != nil {
		// Return default offline record if not found
		return &immodel.ImStaff{AdminID: adminID, IsOnline: 0, MaxLoad: 5, CurrentLoad: 0}, nil
	}
	return &staff, nil
}

// SetStaffOnline marks a staff as online/offline and drains queue if coming online.
func SetStaffOnline(ctx context.Context, adminID uint64, online bool) {
	onlineVal := int8(0)
	if online {
		onlineVal = 1
	}
	// Upsert staff record
	var staff immodel.ImStaff
	err := db.DB.WithContext(ctx).Where("admin_id = ?", adminID).First(&staff).Error
	if err != nil {
		staff = immodel.ImStaff{AdminID: adminID, MaxLoad: 5}
		db.DB.WithContext(ctx).Create(&staff)
	}
	db.DB.WithContext(ctx).Model(&staff).Update("is_online", onlineVal)

	if online {
		DrainQueue(ctx)
	}
}

// AcceptSession lets a staff manually accept a waiting session.
func AcceptSession(ctx context.Context, sessionID, staffID uint64) error {
	var session immodel.ImSession
	if err := db.DB.WithContext(ctx).First(&session, sessionID).Error; err != nil {
		return err
	}
	if session.Status != immodel.SessionStatusWaiting {
		return fmt.Errorf("session is not waiting")
	}
	assignSession(ctx, &session, staffID)
	recordEventBestEffort(ctx, EventInput{
		Event:     immodel.ImEventStaffAccept,
		SessionID: sessionID,
		UserID:    session.UserID,
		StaffID:   staffID,
		Source:    immodel.ImEventSourceStaff,
		Success:   true,
	})
	DrainQueue(ctx) // recompute positions
	return nil
}

// TransferSession reassigns a session from one staff to another.
// The original staff loses the load slot; the new staff gains one.
// A system message is inserted into the conversation and both staff are notified.
func TransferSession(ctx context.Context, sessionID, fromStaffID, toStaffID uint64, remark string) error {
	var session immodel.ImSession
	if err := db.DB.WithContext(ctx).First(&session, sessionID).Error; err != nil {
		return err
	}
	if session.Status != immodel.SessionStatusOngoing {
		return fmt.Errorf("session is not ongoing")
	}
	if session.StaffID != fromStaffID {
		return fmt.Errorf("session is not assigned to you")
	}

	// Write transfer log
	log := &immodel.ImTransferLog{
		SessionID:   sessionID,
		FromStaffID: fromStaffID,
		ToStaffID:   toStaffID,
		Remark:      remark,
	}
	db.DB.WithContext(ctx).Create(log)

	// Update session staff
	db.DB.WithContext(ctx).Model(&session).Update("staff_id", toStaffID)
	session.StaffID = toStaffID

	// Adjust load counters
	db.DB.WithContext(ctx).Model(&immodel.ImStaff{}).
		Where("admin_id = ? AND current_load > 0", fromStaffID).
		UpdateColumn("current_load", db.DB.Raw("current_load - 1"))
	db.DB.WithContext(ctx).Model(&immodel.ImStaff{}).
		Where("admin_id = ?", toStaffID).
		UpdateColumn("current_load", db.DB.Raw("current_load + 1"))

	// Insert system message visible to all parties
	sysMsg := &immodel.ImMessage{
		SessionID:  sessionID,
		SenderType: immodel.SenderSystem,
		SenderID:   0,
		Type:       immodel.MsgTypeSystem,
		Content:    fmt.Sprintf("会话已转接"),
	}
	if remark != "" {
		sysMsg.Content = fmt.Sprintf("会话已转接：%s", remark)
	}
	SaveMessage(ctx, sysMsg)

	sysMsgData, _ := json.Marshal(Frame{Type: "msg", SessionID: sessionID, Payload: map[string]any{
		"msg_type":    immodel.MsgTypeSystem,
		"content":     sysMsg.Content,
		"sender_type": immodel.SenderSystem,
	}})

	// Notify original staff: remove session from their list
	GlobalHub.Send(fmt.Sprintf("staff_%d", fromStaffID), mustMarshal(Frame{
		Type: "assign", SessionID: sessionID,
		Payload: map[string]any{"action": "transfer_out", "to_staff_id": toStaffID},
	}))

	// Notify new staff: add session to their list
	GlobalHub.Send(fmt.Sprintf("staff_%d", toStaffID), mustMarshal(Frame{
		Type: "assign", SessionID: sessionID,
		Payload: map[string]any{"action": "transfer_in", "from_staff_id": fromStaffID, "user_id": session.UserID},
	}))

	// Push system message to user
	GlobalHub.Send(fmt.Sprintf("user_%d", session.UserID), sysMsgData)

	// Notify user about transfer
	GlobalHub.Send(fmt.Sprintf("user_%d", session.UserID), mustMarshal(Frame{
		Type: "assign", SessionID: sessionID,
		Payload: map[string]any{"action": "transfer"},
	}))
	recordEventBestEffort(ctx, EventInput{
		Event:     immodel.ImEventSessionTransfer,
		SessionID: sessionID,
		UserID:    session.UserID,
		StaffID:   toStaffID,
		Source:    immodel.ImEventSourceStaff,
		Success:   true,
		Extra:     map[string]any{"from_staff_id": fromStaffID, "remark": remark},
	})

	return nil
}

func mustMarshal(v any) []byte {
	data, _ := json.Marshal(v)
	return data
}

// CloseSession closes a session and frees staff load, then drains queue.
func CloseSession(ctx context.Context, sessionID uint64) error {
	var session immodel.ImSession
	if err := db.DB.WithContext(ctx).First(&session, sessionID).Error; err != nil {
		return err
	}
	db.DB.WithContext(ctx).Model(&session).Update("status", immodel.SessionStatusClosed)
	if session.StaffID > 0 {
		db.DB.WithContext(ctx).Model(&immodel.ImStaff{}).
			Where("admin_id = ? AND current_load > 0", session.StaffID).
			UpdateColumn("current_load", db.DB.Raw("current_load - 1"))
	}
	// Notify user session closed
	closeFrame := Frame{Type: "close", SessionID: sessionID, Payload: map[string]any{}}
	data, _ := json.Marshal(closeFrame)
	GlobalHub.Send(fmt.Sprintf("user_%d", session.UserID), data)
	recordEventBestEffort(ctx, EventInput{
		Event:     immodel.ImEventSessionClose,
		SessionID: sessionID,
		UserID:    session.UserID,
		StaffID:   session.StaffID,
		Source:    immodel.ImEventSourceSystem,
		Success:   true,
	})

	DrainQueue(ctx)
	return nil
}

// SaveMessage persists a message and updates session last_msg.
func SaveMessage(ctx context.Context, msg *immodel.ImMessage) error {
	if err := db.DB.WithContext(ctx).Create(msg).Error; err != nil {
		return err
	}
	db.DB.WithContext(ctx).Model(&immodel.ImSession{}).Where("id = ?", msg.SessionID).
		Updates(map[string]any{"last_msg": msg.Content, "unread_count": db.DB.Raw("unread_count + 1")})
	recordEventBestEffort(ctx, EventInput{
		Event:     immodel.ImEventMessageSent,
		SessionID: msg.SessionID,
		UserID: func() uint64 {
			if msg.SenderType == immodel.SenderUser {
				return msg.SenderID
			}
			return 0
		}(),
		StaffID: func() uint64 {
			if msg.SenderType == immodel.SenderStaff {
				return msg.SenderID
			}
			return 0
		}(),
		MessageID: msg.ID,
		Source:    eventSourceFromSender(msg.SenderType),
		Success:   true,
		Extra:     map[string]any{"type": msg.Type},
	})
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
	if staffID > 0 {
		tx = tx.Where("staff_id = ?", staffID)
	}
	if status > 0 {
		tx = tx.Where("status = ?", status)
	}
	var list []immodel.ImSession
	err := tx.Order("updated_at desc").Limit(50).Find(&list).Error
	return list, err
}

// ListMessages returns message history for a session.
func ListMessages(ctx context.Context, sessionID uint64, page, size int) ([]immodel.ImMessage, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 50
	}
	var total int64
	db.DB.WithContext(ctx).Model(&immodel.ImMessage{}).Where("session_id = ?", sessionID).Count(&total)
	var list []immodel.ImMessage
	err := db.DB.WithContext(ctx).Where("session_id = ?", sessionID).
		Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

// CheckAutoReply looks for a matching auto-reply rule and returns reply text.
func CheckAutoReply(ctx context.Context, content string) string {
	var rules []immodel.ImAutoReply
	db.DB.WithContext(ctx).Where("status = 1").Order("sort asc").Find(&rules)
	for _, r := range rules {
		switch r.MatchType {
		case 1: // exact
			if content == r.Keyword {
				return r.Reply
			}
		case 2: // contains
			if strings.Contains(content, r.Keyword) {
				return r.Reply
			}
		}
	}
	return ""
}

// HandleWSStaff manages a staff WebSocket connection (receive-only from hub).
func HandleWSStaff(conn *websocket.Conn, clientID string, adminID uint64) {
	ctx := context.Background()
	SetStaffOnline(ctx, adminID, true)

	client := &Client{ID: clientID, Conn: conn, Send: make(chan []byte, 64)}
	GlobalHub.Register(client)
	defer func() {
		GlobalHub.Unregister(client)
		SetStaffOnline(ctx, adminID, false)
	}()

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
		if _, raw, err := conn.ReadMessage(); err != nil {
			break
		} else {
			var frame Frame
			if json.Unmarshal(raw, &frame) == nil {
				switch frame.Type {
				case "typing_draft", "typing_stop":
					var session immodel.ImSession
					if frame.SessionID > 0 && db.DB.WithContext(ctx).First(&session, frame.SessionID).Error == nil && session.StaffID == adminID {
						draft, _ := frame.Payload["draft"].(string)
						forwardTypingFrame(session, immodel.SenderStaff, adminID, draft, frame.Type == "typing_stop")
					}
				}
			}
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
	if extra := messageExtraPayload(msg.Extra); extra != nil {
		frame.Payload["extra"] = extra
	}
	data, _ := json.Marshal(frame)
	GlobalHub.Send(userKey, data)
}

// answerWithAI generates a RAG reply for an AI-mode session and pushes it to
// the user. Runs in its own goroutine; failures degrade to a polite fallback
// that nudges the user toward a human agent.
func answerWithAI(session immodel.ImSession, userText string) {
	ctx := context.Background()
	start := time.Now()
	reply, err := AIAnswer(ctx, &session, userText)
	aiOK := err == nil && strings.TrimSpace(reply) != ""
	if !aiOK {
		reply = "抱歉，我暂时无法回答这个问题。您可以换个说法，或输入“人工”转接人工客服。"
		recordEventBestEffort(ctx, EventInput{
			Event:     immodel.ImEventAIFailed,
			SessionID: session.ID,
			UserID:    session.UserID,
			Source:    immodel.ImEventSourceAI,
			Success:   false,
			LatencyMS: time.Since(start).Milliseconds(),
			Extra:     map[string]any{"error": fmt.Sprint(err)},
		})
	}
	aiMsg := &immodel.ImMessage{
		SessionID:  session.ID,
		SenderType: immodel.SenderAI,
		Type:       immodel.MsgTypeText,
		Content:    reply,
	}
	SaveMessage(ctx, aiMsg)
	if aiOK {
		recordEventBestEffort(ctx, EventInput{
			Event:     immodel.ImEventAIReply,
			SessionID: session.ID,
			UserID:    session.UserID,
			MessageID: aiMsg.ID,
			Source:    immodel.ImEventSourceAI,
			Success:   true,
			LatencyMS: time.Since(start).Milliseconds(),
		})
	}
	data, _ := json.Marshal(Frame{Type: "msg", SessionID: session.ID, Payload: map[string]any{
		"msg_id":      fmt.Sprintf("%d", aiMsg.ID),
		"msg_type":    immodel.MsgTypeText,
		"content":     reply,
		"sender_type": immodel.SenderAI,
	}})
	GlobalHub.Send(fmt.Sprintf("user_%d", session.UserID), data)
}

// HandleWS manages one WebSocket client lifecycle (readPump + writePump).
func HandleWS(conn *websocket.Conn, clientID string, session *immodel.ImSession) {
	client := &Client{ID: clientID, Conn: conn, Send: make(chan []byte, 64)}
	GlobalHub.Register(client)
	defer GlobalHub.Unregister(client)

	ctx := context.Background()

	// Send queue position if still waiting
	if session.Status == immodel.SessionStatusWaiting && session.QueuePosition > 0 {
		queueFrame := Frame{Type: "queue", SessionID: session.ID, Payload: map[string]any{
			"position": session.QueuePosition,
		}}
		data, _ := json.Marshal(queueFrame)
		client.Send <- data
	}

	// Replay unread messages on connect
	unread, _ := GetUnread(ctx, session.ID)
	for _, msg := range unread {
		frame := Frame{Type: "msg", SessionID: session.ID, Payload: map[string]any{
			"msg_id":      fmt.Sprintf("%d", msg.ID),
			"msg_type":    msg.Type,
			"content":     msg.Content,
			"sender_type": msg.SenderType,
		}}
		if extra := messageExtraPayload(msg.Extra); extra != nil {
			frame.Payload["extra"] = extra
		}
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
		if err != nil {
			break
		}
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		var frame Frame
		if err := json.Unmarshal(raw, &frame); err != nil {
			continue
		}

		switch frame.Type {
		case "msg":
			content, _ := frame.Payload["content"].(string)
			msgType, _ := frame.Payload["msg_type"].(string)
			if msgType == "" {
				msgType = immodel.MsgTypeText
			}
			extra := normalizeMessageExtra(frame.Payload["extra"])

			msg := &immodel.ImMessage{
				SessionID:  session.ID,
				SenderType: immodel.SenderUser,
				SenderID:   session.UserID,
				Type:       msgType,
				Content:    content,
				Extra:      extra,
			}
			SaveMessage(ctx, msg)

			// Reload current session state: mode/staff can change mid-connection
			// (e.g. user switched to human, or a staff accepted/transferred).
			var cur immodel.ImSession
			if db.DB.WithContext(ctx).First(&cur, session.ID).Error == nil {
				session = &cur
			}

			// Human mode: forward to the assigned staff; otherwise the user is
			// still queued (queue notice already delivered). Keep the legacy
			// keyword auto-reply for unassigned human sessions (AI disabled).
			if session.Mode == immodel.SessionModeHuman {
				if session.StaffID > 0 {
					GlobalHub.Send(fmt.Sprintf("staff_%d", session.StaffID), raw)
				} else if reply := CheckAutoReply(ctx, content); reply != "" {
					autoMsg := &immodel.ImMessage{
						SessionID: session.ID, SenderType: immodel.SenderStaff,
						Type: immodel.MsgTypeText, Content: reply,
					}
					SaveMessage(ctx, autoMsg)
					data, _ := json.Marshal(Frame{Type: "msg", SessionID: session.ID, Payload: map[string]any{
						"msg_type": immodel.MsgTypeText, "content": reply,
					}})
					client.Send <- data
				}
				break
			}

			// AI mode (default). A human-request keyword hands off to a person.
			cfg := LoadAIConfig()
			if msgType == immodel.MsgTypeText && IsHumanRequest(cfg, content) {
				SwitchToHuman(ctx, session.ID)
				break
			}
			// Generate the AI reply off the read loop so further messages and
			// pings stay responsive while the local LLM is thinking.
			typing, _ := json.Marshal(Frame{Type: "typing", SessionID: session.ID, Payload: map[string]any{"sender_type": immodel.SenderAI}})
			client.Send <- typing
			sess := *session
			go answerWithAI(sess, content)

		case "to_human":
			// Explicit "转人工" button: hand off regardless of locale/keywords.
			SwitchToHuman(ctx, session.ID)

		case "typing_draft", "typing_stop":
			var cur immodel.ImSession
			if db.DB.WithContext(ctx).First(&cur, session.ID).Error == nil {
				session = &cur
			}
			draft, _ := frame.Payload["draft"].(string)
			forwardTypingFrame(*session, immodel.SenderUser, session.UserID, draft, frame.Type == "typing_stop")

		case "ping":
			pong := Frame{Type: "pong"}
			data, _ := json.Marshal(pong)
			client.Send <- data
		}
	}
}
