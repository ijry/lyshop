package inventory

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/core/db"
	"gorm.io/gorm"
)

type TaskPayload struct {
	SkuID  uint64        `json:"sku_id,omitempty"`
	Stock  int           `json:"stock,omitempty"`
	Items  []ReserveItem `json:"items,omitempty"`
	Reason string        `json:"reason,omitempty"`
}

func NewIntegrationTask(provider, bizType, bizNo, action string, payload TaskPayload) *InventoryIntegrationTask {
	raw, _ := json.Marshal(payload)
	now := time.Now()
	return &InventoryIntegrationTask{
		Provider:     provider,
		BizType:      bizType,
		BizNo:        bizNo,
		Action:       action,
		Payload:      string(raw),
		Status:       "pending",
		AttemptCount: 0,
		NextRetryAt:  &now,
	}
}

func EnqueueInventoryTask(tx *gorm.DB, provider, bizType, bizNo, action string, payload TaskPayload) error {
	if tx == nil {
		tx = db.DB
	}
	return tx.Create(NewIntegrationTask(provider, bizType, bizNo, action, payload)).Error
}
