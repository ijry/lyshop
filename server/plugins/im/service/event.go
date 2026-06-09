package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/ijry/lyshop/core/db"
	immodel "github.com/ijry/lyshop/plugins/im/model"
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
