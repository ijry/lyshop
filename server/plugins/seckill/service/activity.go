package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ijry/lyshop/server/core/db"
	seckillmodel "github.com/ijry/lyshop/server/plugins/seckill/model"
)

// ListActivities 获取秒杀活动列表
func ListActivities(ctx context.Context, status int8, page, size int) ([]seckillmodel.SeckillActivity, int64, error) {
	var activities []seckillmodel.SeckillActivity
	var total int64

	query := db.DB.WithContext(ctx).Model(&seckillmodel.SeckillActivity{})

	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := query.Order("sort DESC, id DESC").Offset(offset).Limit(size).Find(&activities).Error

	return activities, total, err
}

// GetActivity 获取活动详情
func GetActivity(ctx context.Context, id uint64) (*seckillmodel.SeckillActivity, error) {
	var activity seckillmodel.SeckillActivity
	err := db.DB.WithContext(ctx).First(&activity, id).Error
	return &activity, err
}

// CreateActivity 创建秒杀活动
func CreateActivity(ctx context.Context, activity *seckillmodel.SeckillActivity) error {
	// 检查时间冲突
	if activity.StartAt != nil && activity.EndAt != nil {
		var count int64
		db.DB.WithContext(ctx).Model(&seckillmodel.SeckillActivity{}).
			Where("status = 1").
			Where("NOT (end_at <= ? OR start_at >= ?)", activity.StartAt, activity.EndAt).
			Count(&count)

		if count > 0 {
			return fmt.Errorf("活动时间与现有活动冲突")
		}
	}

	return db.DB.WithContext(ctx).Create(activity).Error
}

// UpdateActivity 更新秒杀活动
func UpdateActivity(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&seckillmodel.SeckillActivity{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// DeleteActivity 删除秒杀活动
func DeleteActivity(ctx context.Context, id uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *db.DB) error {
		// 删除活动商品
		if err := tx.Where("activity_id = ?", id).Delete(&seckillmodel.SeckillProduct{}).Error; err != nil {
			return err
		}
		// 删除活动
		return tx.Delete(&seckillmodel.SeckillActivity{}, id).Error
	})
}

// GetActiveActivities 获取当前有效的秒杀活动
func GetActiveActivities(ctx context.Context) ([]seckillmodel.SeckillActivity, error) {
	var activities []seckillmodel.SeckillActivity
	now := time.Now()

	err := db.DB.WithContext(ctx).
		Where("status = 1").
		Where("start_at <= ? AND end_at >= ?", now, now).
		Order("sort DESC, id DESC").
		Find(&activities).Error

	return activities, err
}
