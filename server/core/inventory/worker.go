package inventory

import (
	"time"

	"gorm.io/gorm"
)

type AsyncTaskProcessor interface {
	ProcessTask(tx *gorm.DB, task *InventoryIntegrationTask, now time.Time) error
}

func ClaimDueTask(db *gorm.DB, worker string, now time.Time) (*InventoryIntegrationTask, error) {
	var task InventoryIntegrationTask
	err := db.Transaction(func(tx *gorm.DB) error {
		query := tx.
			Where("status = ?", TaskStatusPending).
			Where("next_retry_at IS NULL OR next_retry_at <= ?", now).
			Where("lock_expires_at IS NULL OR lock_expires_at <= ?", now).
			Order("id ASC").
			First(&task).Error
		if err := query; err != nil {
			return err
		}

		lockUntil := now.Add(2 * time.Minute)
		return tx.Model(&InventoryIntegrationTask{}).
			Where("id = ?", task.ID).
			Updates(map[string]any{
				"status":          TaskStatusProcessing,
				"lock_owner":      worker,
				"lock_expires_at": lockUntil,
			}).Error
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	task.Status = TaskStatusProcessing
	task.LockOwner = worker
	expiresAt := now.Add(2 * time.Minute)
	task.LockExpiresAt = &expiresAt
	return &task, nil
}

func ProcessTask(db *gorm.DB, processor AsyncTaskProcessor, task *InventoryIntegrationTask, worker string, now time.Time) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := processor.ProcessTask(tx, task, now); err != nil {
			return MarkTaskRetry(tx, task, err, now)
		}
		return MarkTaskSuccess(tx, task, now)
	})
	if err != nil {
		return err
	}

	task.LockOwner = worker
	return nil
}
