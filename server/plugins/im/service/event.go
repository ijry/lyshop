package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/ijry/lyshop/core/db"
	immodel "github.com/ijry/lyshop/plugins/im/model"
	"gorm.io/gorm"
)

type EventInput struct {
	Event     string
	SessionID uint64
	UserID    uint64
	StaffID   uint64
	MessageID uint64
	Source    string
	Success   bool
	LatencyMS int64
	Extra     map[string]any
}

type EventLogQuery struct {
	Event     string
	SessionID uint64
	UserID    uint64
	StaffID   uint64
	Source    string
	Success   *int8
	Page      int
	Size      int
}

func RecordEvent(ctx context.Context, input EventInput) error {
	if strings.TrimSpace(input.Event) == "" {
		return nil
	}
	source := strings.TrimSpace(input.Source)
	if source == "" {
		source = immodel.ImEventSourceSystem
	}
	success := int8(0)
	if input.Success {
		success = 1
	}
	extra := ""
	if len(input.Extra) > 0 {
		raw, _ := json.Marshal(input.Extra)
		extra = string(raw)
	}
	row := &immodel.ImEventLog{
		Event:     input.Event,
		SessionID: input.SessionID,
		UserID:    input.UserID,
		StaffID:   input.StaffID,
		MessageID: input.MessageID,
		Source:    source,
		Success:   success,
		LatencyMS: input.LatencyMS,
		Extra:     extra,
	}
	return db.DB.WithContext(ctx).Create(row).Error
}

func ListEventLogs(ctx context.Context, q EventLogQuery) ([]immodel.ImEventLog, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 || q.Size > 100 {
		q.Size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&immodel.ImEventLog{})
	if q.Event != "" {
		tx = tx.Where("event = ?", q.Event)
	}
	if q.SessionID > 0 {
		tx = tx.Where("session_id = ?", q.SessionID)
	}
	if q.UserID > 0 {
		tx = tx.Where("user_id = ?", q.UserID)
	}
	if q.StaffID > 0 {
		tx = tx.Where("staff_id = ?", q.StaffID)
	}
	if q.Source != "" {
		tx = tx.Where("source = ?", q.Source)
	}
	if q.Success != nil {
		tx = tx.Where("success = ?", *q.Success)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []immodel.ImEventLog
	err := tx.Order("id desc").Offset((q.Page - 1) * q.Size).Limit(q.Size).Find(&list).Error
	return list, total, err
}

type AnalyticsQuery struct {
	From    time.Time
	To      time.Time
	StaffID uint64
}

type AnalyticsResult struct {
	Summary map[string]int64 `json:"summary"`
	Trend   []map[string]any `json:"trend"`
}

func emptyAnalyticsCounters() map[string]int64 {
	return map[string]int64{
		"sessions":   0,
		"messages":   0,
		"ai_replies": 0,
		"ai_failed":  0,
		"rag_hits":   0,
		"to_human":   0,
		"accepts":    0,
		"closes":     0,
		"transfers":  0,
		"files":      0,
	}
}

func incrementAnalyticsCounter(counters map[string]int64, event string, count int64) {
	switch event {
	case immodel.ImEventSessionCreated:
		counters["sessions"] += count
	case immodel.ImEventMessageSent:
		counters["messages"] += count
	case immodel.ImEventAIReply:
		counters["ai_replies"] += count
	case immodel.ImEventAIFailed:
		counters["ai_failed"] += count
	case immodel.ImEventRAGHit:
		counters["rag_hits"] += count
	case immodel.ImEventToHuman:
		counters["to_human"] += count
	case immodel.ImEventStaffAccept:
		counters["accepts"] += count
	case immodel.ImEventSessionClose:
		counters["closes"] += count
	case immodel.ImEventSessionTransfer:
		counters["transfers"] += count
	case immodel.ImEventFileUploaded:
		counters["files"] += count
	}
}

func analyticsBaseQuery(ctx context.Context, q AnalyticsQuery) *gorm.DB {
	tx := db.DB.WithContext(ctx).Model(&immodel.ImEventLog{})
	if !q.From.IsZero() {
		tx = tx.Where("created_at >= ?", q.From)
	}
	if !q.To.IsZero() {
		tx = tx.Where("created_at < ?", q.To)
	}
	if q.StaffID > 0 {
		tx = tx.Where("staff_id = ?", q.StaffID)
	}
	return tx
}

func Analytics(ctx context.Context, q AnalyticsQuery) (*AnalyticsResult, error) {
	var rows []struct {
		Event string
		Count int64
	}
	if err := analyticsBaseQuery(ctx, q).
		Select("event, COUNT(*) AS count").
		Group("event").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	summary := emptyAnalyticsCounters()
	for _, r := range rows {
		incrementAnalyticsCounter(summary, r.Event, r.Count)
	}

	var events []immodel.ImEventLog
	if err := analyticsBaseQuery(ctx, q).Order("created_at asc").Find(&events).Error; err != nil {
		return nil, err
	}
	byDay := map[string]map[string]int64{}
	var days []string
	for _, event := range events {
		day := event.CreatedAt.Format("2006-01-02")
		if _, ok := byDay[day]; !ok {
			byDay[day] = emptyAnalyticsCounters()
			days = append(days, day)
		}
		incrementAnalyticsCounter(byDay[day], event.Event, 1)
	}
	trend := make([]map[string]any, 0, len(days))
	for _, day := range days {
		row := map[string]any{"date": day}
		for key, value := range byDay[day] {
			row[key] = value
		}
		trend = append(trend, row)
	}
	return &AnalyticsResult{Summary: summary, Trend: trend}, nil
}
