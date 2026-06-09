package inventory

import (
	"time"

	"github.com/ijry/lyshop/model"
)

const (
	InventoryStatusNone      = "none"
	InventoryStatusPending   = "pending"
	InventoryStatusReserved  = "reserved"
	InventoryStatusConfirmed = "confirmed"
	InventoryStatusReleased  = "released"
	InventoryStatusFailed    = "failed"

	TaskStatusPending    = "pending"
	TaskStatusProcessing = "processing"
	TaskStatusSuccess    = "success"
	TaskStatusFailed     = "failed"
)

type InventoryReservation struct {
	model.Base
	BizType   string     `gorm:"size:32;not null;index:idx_inventory_reservation_biz,priority:1" json:"biz_type"`
	BizNo     string     `gorm:"size:64;not null;index:idx_inventory_reservation_biz,priority:2" json:"biz_no"`
	SkuID     uint64     `gorm:"not null;index:idx_inventory_reservation_sku" json:"sku_id"`
	Qty       int        `gorm:"not null" json:"qty"`
	Status    string     `gorm:"size:16;not null;default:'reserved';index" json:"status"`
	Reason    string     `gorm:"size:128" json:"reason"`
	ExpiredAt *time.Time `json:"expired_at"`
}

type OrderInventoryState struct {
	model.Base
	OrderNo   string `gorm:"size:64;uniqueIndex;not null" json:"order_no"`
	BizType   string `gorm:"size:32;not null;default:'order'" json:"biz_type"`
	Status    string `gorm:"size:16;not null;default:'none';index" json:"status"`
	Provider  string `gorm:"size:32;not null;default:'local'" json:"provider"`
	LastError string `gorm:"size:255" json:"last_error"`
}

type InventoryIntegrationTask struct {
	model.Base
	Provider       string     `gorm:"size:32;not null;index" json:"provider"`
	BizType        string     `gorm:"size:32;not null;index" json:"biz_type"`
	BizNo          string     `gorm:"size:64;not null;index" json:"biz_no"`
	Action         string     `gorm:"size:16;not null;index" json:"action"`
	RequestID      string     `gorm:"size:64;index" json:"request_id"`
	Payload        string     `gorm:"type:json" json:"payload"`
	Status         string     `gorm:"size:16;not null;default:'pending';index" json:"status"`
	AttemptCount   int        `gorm:"not null;default:0" json:"attempt_count"`
	MaxAttempts    int        `gorm:"not null;default:0" json:"max_attempts"`
	BackoffSeconds int        `gorm:"not null;default:0" json:"backoff_seconds"`
	NextRetryAt    *time.Time `gorm:"index" json:"next_retry_at"`
	LastError      string     `gorm:"size:255" json:"last_error"`
	LockOwner      string     `gorm:"size:64;index" json:"lock_owner"`
	LockExpiresAt  *time.Time `gorm:"index" json:"lock_expires_at"`
	CompletedAt    *time.Time `json:"completed_at"`
	LastCallbackID string     `gorm:"size:64;index" json:"last_callback_id"`
}
