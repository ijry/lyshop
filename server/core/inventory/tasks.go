package inventory

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/config"
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
	maxAttempts := config.Global.ExternalWMS.Retry.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 8
	}
	backoffSeconds := config.Global.ExternalWMS.Retry.BackoffSeconds
	if backoffSeconds <= 0 {
		backoffSeconds = 30
	}
	return &InventoryIntegrationTask{
		Provider:       provider,
		BizType:        bizType,
		BizNo:          bizNo,
		Action:         action,
		Payload:        string(raw),
		Status:         TaskStatusPending,
		AttemptCount:   0,
		MaxAttempts:    maxAttempts,
		BackoffSeconds: backoffSeconds,
		NextRetryAt:    &now,
	}
}

func EnqueueInventoryTask(tx *gorm.DB, provider, bizType, bizNo, action string, payload TaskPayload) error {
	if tx == nil {
		tx = db.DB
	}
	return tx.Create(NewIntegrationTask(provider, bizType, bizNo, action, payload)).Error
}

func MarkTaskSuccess(tx *gorm.DB, task *InventoryIntegrationTask, now time.Time) error {
	return tx.Model(&InventoryIntegrationTask{}).
		Where("id = ?", task.ID).
		Updates(map[string]any{
			"status":          TaskStatusSuccess,
			"attempt_count":   gorm.Expr("attempt_count + 1"),
			"last_error":      "",
			"lock_owner":      "",
			"lock_expires_at": nil,
			"completed_at":    now,
		}).Error
}

func MarkTaskRetry(tx *gorm.DB, task *InventoryIntegrationTask, cause error, now time.Time) error {
	attemptCount := task.AttemptCount + 1
	updates := map[string]any{
		"attempt_count":   attemptCount,
		"last_error":      cause.Error(),
		"lock_owner":      "",
		"lock_expires_at": nil,
	}
	if attemptCount >= task.MaxAttempts {
		updates["status"] = TaskStatusFailed
		updates["completed_at"] = now
		updates["next_retry_at"] = nil
	} else {
		next := now.Add(time.Duration(task.BackoffSeconds) * time.Second)
		updates["status"] = TaskStatusPending
		updates["next_retry_at"] = next
	}
	return tx.Model(&InventoryIntegrationTask{}).
		Where("id = ?", task.ID).
		Updates(updates).Error
}

func CompleteTaskByCallback(tx *gorm.DB, requestID, callbackID, status, message string, now time.Time) error {
	var task InventoryIntegrationTask
	if err := tx.Where("request_id = ?", requestID).First(&task).Error; err != nil {
		return err
	}
	if task.LastCallbackID == callbackID {
		return nil
	}
	if task.Status == TaskStatusSuccess || task.Status == TaskStatusFailed {
		return nil
	}

	updates := map[string]any{
		"last_callback_id": callbackID,
		"last_error":       message,
		"lock_owner":       "",
		"lock_expires_at":  nil,
	}

	inventoryStatus := InventoryStatusPending
	switch status {
	case "success":
		updates["status"] = TaskStatusSuccess
		updates["completed_at"] = now
		inventoryStatus = callbackInventoryStatus(task.Action, true)
	case "failed":
		updates["status"] = TaskStatusFailed
		updates["completed_at"] = now
		inventoryStatus = InventoryStatusFailed
	default:
		updates["status"] = TaskStatusProcessing
	}
	if err := tx.Model(&InventoryIntegrationTask{}).Where("id = ?", task.ID).Updates(updates).Error; err != nil {
		return err
	}
	return updateOrderInventoryStatusByTask(tx, task.BizType, task.BizNo, inventoryStatus, message)
}

func callbackInventoryStatus(action string, success bool) string {
	if !success {
		return InventoryStatusFailed
	}
	switch action {
	case "reserve":
		return InventoryStatusReserved
	case "confirm", "deduct":
		return InventoryStatusConfirmed
	case "release", "restore":
		return InventoryStatusReleased
	default:
		return InventoryStatusConfirmed
	}
}

func updateOrderInventoryStatusByTask(tx *gorm.DB, bizType, bizNo, inventoryStatus, message string) error {
	if bizType != "order" {
		return nil
	}
	return tx.Table("orders").Where("order_no = ?", bizNo).Update("inventory_status", inventoryStatus).Error
}
