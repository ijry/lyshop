package service

import (
	"context"
	"time"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/plugins/group_buy/model"
)

// ListActivities 获取活动列表
func ListActivities(ctx context.Context, page, size int, status *int8) ([]model.GroupBuyActivity, int64, error) {
	var activities []model.GroupBuyActivity
	var total int64

	query := db.DB.WithContext(ctx).Model(&model.GroupBuyActivity{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("sort DESC, id DESC").Offset(offset).Limit(size).Find(&activities).Error; err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}

// GetActivity 获取活动详情
func GetActivity(ctx context.Context, id uint64) (*model.GroupBuyActivity, error) {
	var activity model.GroupBuyActivity
	if err := db.DB.WithContext(ctx).First(&activity, id).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

// CreateActivity 创建活动
func CreateActivity(ctx context.Context, activity *model.GroupBuyActivity) error {
	// 检查时间冲突
	if activity.StartAt != nil && activity.EndAt != nil {
		var count int64
		err := db.DB.WithContext(ctx).Model(&model.GroupBuyActivity{}).
			Where("status = 1").
			Where("((start_at <= ? AND end_at >= ?) OR (start_at <= ? AND end_at >= ?))",
				activity.StartAt, activity.StartAt, activity.EndAt, activity.EndAt).
			Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return db.ErrConflict
		}
	}

	return db.DB.WithContext(ctx).Create(activity).Error
}

// UpdateActivity 更新活动
func UpdateActivity(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&model.GroupBuyActivity{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteActivity 删除活动
func DeleteActivity(ctx context.Context, id uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *db.Tx) error {
		// 删除活动商品
		if err := tx.Where("activity_id = ?", id).Delete(&model.GroupBuyProduct{}).Error; err != nil {
			return err
		}
		// 删除活动
		return tx.Delete(&model.GroupBuyActivity{}, id).Error
	})
}

// GetActiveActivities 获取当前有效的活动
func GetActiveActivities(ctx context.Context) ([]model.GroupBuyActivity, error) {
	var activities []model.GroupBuyActivity
	now := time.Now()
	err := db.DB.WithContext(ctx).
		Where("status = 1").
		Where("start_at <= ? AND end_at >= ?", now, now).
		Order("sort DESC, id DESC").
		Find(&activities).Error
	return activities, err
}
