package model

import (
	"encoding/json"

	"github.com/ijry/lyshop/model"
)

const (
	SessionStatusWaiting int8 = 1
	SessionStatusOngoing int8 = 2
	SessionStatusClosed  int8 = 3

	// SessionMode marks who is currently serving the session.
	// New sessions start in "ai" mode (answered by the local LLM); when the
	// user requests a human agent we switch to "human" (queue/assign flow).
	SessionModeAI    = "ai"
	SessionModeHuman = "human"

	MsgTypeText        = "text"
	MsgTypeImage       = "image"
	MsgTypeFile        = "file"
	MsgTypeProductCard = "product_card"
	MsgTypeOrderCard   = "order_card"
	MsgTypeSystem      = "system" // 系统通知（接入、排队等）

	SenderUser   int8 = 1
	SenderStaff  int8 = 2
	SenderSystem int8 = 0
	SenderAI     int8 = 3 // 本地大模型客服
)

const (
	ImEventSessionCreated  = "session_created"
	ImEventMessageSent     = "message_sent"
	ImEventAIReply         = "ai_reply"
	ImEventAIFailed        = "ai_failed"
	ImEventRAGHit          = "rag_hit"
	ImEventToHuman         = "to_human"
	ImEventStaffAccept     = "staff_accept"
	ImEventSessionClose    = "session_close"
	ImEventSessionTransfer = "session_transfer"
	ImEventFileUploaded    = "file_uploaded"
	ImEventWebSearch       = "web_search"

	ImEventSourceUser   = "user"
	ImEventSourceStaff  = "staff"
	ImEventSourceAI     = "ai"
	ImEventSourceSystem = "system"
)

// ImStaff tracks online status and load for each staff member.
type ImStaff struct {
	model.Base
	AdminID     uint64 `gorm:"not null;uniqueIndex" json:"admin_id"`
	IsOnline    int8   `gorm:"not null;default:0"   json:"is_online"`
	MaxLoad     int    `gorm:"not null;default:5"   json:"max_load"`
	CurrentLoad int    `gorm:"not null;default:0"   json:"current_load"`
}

type ImSession struct {
	model.Base
	UserID          uint64 `gorm:"not null;index"            json:"user_id"`
	StaffID         uint64 `gorm:"not null;default:0;index"  json:"staff_id"`
	Mode            string `gorm:"size:16;not null;default:'ai'" json:"mode"` // ai|human
	Status          int8   `gorm:"not null"                  json:"status"`
	QueuePosition   int    `gorm:"not null;default:0"        json:"queue_position"`
	LastMsg         string `gorm:"size:255"                  json:"last_msg"`
	UnreadCount     int    `gorm:"not null;default:0"        json:"unread_count"`
	VisitorID       string `gorm:"size:128;index"            json:"visitor_id"`
	VisitorIP       string `gorm:"size:64"                   json:"visitor_ip"`
	VisitorLocation string `gorm:"size:128"                json:"visitor_location"`
	VisitorBrowser  string `gorm:"size:64"                  json:"visitor_browser"`
	VisitorOS       string `gorm:"size:64"                  json:"visitor_os"`
	VisitorLanguage string `gorm:"size:32"                  json:"visitor_language"`
	VisitorReferrer string `gorm:"size:512"                 json:"visitor_referrer"`
	VisitorURL      string `gorm:"size:512"                 json:"visitor_url"`
	VisitorDevice   string `gorm:"size:64"                  json:"visitor_device"`
	VisitorExtra    string `gorm:"type:json"                json:"visitor_extra"`
}

// ImTransferLog records every session transfer for audit and history.
type ImTransferLog struct {
	model.Base
	SessionID   uint64 `gorm:"not null;index" json:"session_id"`
	FromStaffID uint64 `gorm:"not null"       json:"from_staff_id"`
	ToStaffID   uint64 `gorm:"not null"       json:"to_staff_id"`
	Remark      string `gorm:"size:255"       json:"remark"`
}

type ImMessage struct {
	model.Base
	SessionID  uint64 `gorm:"not null;index"    json:"session_id"`
	SenderType int8   `gorm:"not null"          json:"sender_type"`
	SenderID   uint64 `gorm:"not null"          json:"sender_id"`
	Type       string `gorm:"size:32;not null"  json:"type"`
	Content    string `gorm:"type:text;not null" json:"content"`
	Extra      string `gorm:"type:json"         json:"extra"`
	IsRead     int8   `gorm:"not null;default:0" json:"is_read"`
}

type ImEventLog struct {
	model.Base
	Event     string `gorm:"size:64;not null;index" json:"event"`
	Level     string `gorm:"size:16;not null;default:'info';index" json:"level"`
	Category  string `gorm:"size:32;not null;default:'im';index" json:"category"`
	TraceID   string `gorm:"size:64;index" json:"trace_id"`
	SessionID uint64 `gorm:"not null;default:0;index" json:"session_id"`
	UserID    uint64 `gorm:"not null;default:0;index" json:"user_id"`
	StaffID   uint64 `gorm:"not null;default:0;index" json:"staff_id"`
	MessageID uint64 `gorm:"not null;default:0;index" json:"message_id"`
	Source    string `gorm:"size:16;not null;default:'system';index" json:"source"`
	Success   int8   `gorm:"not null;default:1;index" json:"success"`
	LatencyMS int64  `gorm:"not null;default:0" json:"latency_ms"`
	Message   string `gorm:"size:512" json:"message"`
	Meta      string `gorm:"type:json" json:"meta"`
	Extra     string `gorm:"type:json" json:"extra"`
}

type ImAutoReply struct {
	model.Base
	Keyword   string `gorm:"size:128;not null" json:"keyword"`
	MatchType int8   `gorm:"not null"          json:"match_type"` // 1精确 2包含 3正则
	Reply     string `gorm:"type:text;not null" json:"reply"`
	Sort      int    `gorm:"not null;default:0" json:"sort"`
	Status    int8   `gorm:"not null;default:1" json:"status"`
}

// ImFeedback records user or auto-eval ratings for an AI answer.
//
// source = "user"  → submitted by the end user (thumbs up/down + optional comment).
// source = "auto"  → LLM-as-judge scores stored by AutoScore.
type ImFeedback struct {
	model.Base
	SessionID    uint64  `gorm:"not null;index"   json:"session_id"`
	Source       string  `gorm:"size:16;not null" json:"source"`   // user|auto
	Rating       int8    `gorm:"not null;default:0" json:"rating"` // user: 1=👍 -1=👎 0=unset
	Comment      string  `gorm:"size:512"          json:"comment"`
	Faithfulness float64 `gorm:"not null;default:0" json:"faithfulness"` // auto 0-5
	Relevance    float64 `gorm:"not null;default:0" json:"relevance"`    // auto 0-5
	Query        string  `gorm:"type:text"          json:"query"`
	Answer       string  `gorm:"type:text"          json:"answer"`
}

const (
	FeedbackSourceUser = "user"
	FeedbackSourceAuto = "auto"
)

// ImKnowledge is one entry in the AI customer-service RAG knowledge base.
// Content is embedded into a vector (Embedding) when the local LLM exposes an
// embeddings endpoint; otherwise retrieval falls back to keyword matching.
type ImKnowledge struct {
	model.Base
	Title     string          `gorm:"size:255;not null"   json:"title"`
	Content   string          `gorm:"type:text;not null"  json:"content"`
	Tags      string          `gorm:"size:255"            json:"tags"`    // 逗号分隔，便于关键词召回
	Embedding json.RawMessage `gorm:"type:json"           json:"-"`       // []float64，前端无需返回
	Indexed   int8            `gorm:"not null;default:0"  json:"indexed"` // 1=已向量化
	Sort      int             `gorm:"not null;default:0"  json:"sort"`
	Status    int8            `gorm:"not null;default:1"  json:"status"` // 1启用 0停用
}
