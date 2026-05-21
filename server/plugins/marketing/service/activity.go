package service

import (
	"context"

	"github.com/ijry/lyshop/core/db"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
)

func ListActivities(ctx context.Context, actType int8, page, size int) ([]mktmodel.Activity, int64, error) {
	if page <= 0 { page = 1 }
	if size <= 0 || size > 100 { size = 20 }
	tx := db.DB.WithContext(ctx).Model(&mktmodel.Activity{})
	if actType > 0 {
		tx = tx.Where("type = ?", actType)
	}
	var total int64
	tx.Count(&total)
	var list []mktmodel.Activity
	err := tx.Order("id desc").Offset((page-1)*size).Limit(size).Find(&list).Error
	return list, total, err
}

func CreateActivity(ctx context.Context, a *mktmodel.Activity, products []mktmodel.ActivityProduct) error {
	if err := db.DB.WithContext(ctx).Create(a).Error; err != nil {
		return err
	}
	for i := range products {
		products[i].ActivityID = a.ID
	}
	if len(products) > 0 {
		return db.DB.WithContext(ctx).Create(&products).Error
	}
	return nil
}

// GetActiveSeckills returns currently active seckill activities with their products.
func GetActiveSeckills(ctx context.Context) ([]mktmodel.Activity, error) {
	var list []mktmodel.Activity
	err := db.DB.WithContext(ctx).
		Where("type = ? AND status = ? AND start_at <= NOW() AND end_at >= NOW()",
			mktmodel.ActivityTypeSeckill, 1).
		Find(&list).Error
	return list, err
}
