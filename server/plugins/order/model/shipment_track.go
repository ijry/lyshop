package model

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/model"
)

type OrderShipmentTrack struct {
	model.Base
	ShipmentID uint64          `gorm:"not null;index:idx_shipment_event_time,priority:1;uniqueIndex:uk_shipment_track,priority:1" json:"shipment_id"`
	Provider   string          `gorm:"size:32;not null"                                                                          json:"provider"`
	TrackHash  string          `gorm:"size:64;not null;uniqueIndex:uk_shipment_track,priority:2"                                json:"track_hash"`
	StatusCode string          `gorm:"size:32;not null"                                                                          json:"status_code"`
	StatusText string          `gorm:"size:255;not null"                                                                         json:"status_text"`
	EventTime  time.Time       `gorm:"not null;index:idx_shipment_event_time,priority:2"                                        json:"event_time"`
	Location   string          `gorm:"size:255"                                                                                  json:"location"`
	RawPayload json.RawMessage `gorm:"type:json"                                                                                 json:"raw_payload"`
}

type OrderShipmentSyncLog struct {
	model.Base
	ShipmentID     uint64 `gorm:"not null;index"       json:"shipment_id"`
	Provider       string `gorm:"size:32;not null"     json:"provider"`
	Success        bool   `gorm:"not null"             json:"success"`
	ErrorCode      string `gorm:"size:64"              json:"error_code"`
	ErrorMessage   string `gorm:"size:500"             json:"error_message"`
	CostMS         int64  `gorm:"not null;default:0"   json:"cost_ms"`
	ResponseDigest string `gorm:"size:500"             json:"response_digest"`
}
