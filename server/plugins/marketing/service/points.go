package service

import (
	"context"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/model"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
	"gorm.io/gorm"
)

// AddPoints adds (or deducts if negative) points for a user.
func AddPoints(ctx context.Context, userID uint64, points int, logType int8, remark string) error {
	// Update user points balance
	db.DB.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		UpdateColumn("points", gorm.Expr("points + ?", points))
	return db.DB.WithContext(ctx).Create(&mktmodel.PointsLog{
		UserID: userID, Type: logType, Points: points, Remark: remark,
	}).Error
}

// ListPointsLogs returns point change history for a user.
func ListPointsLogs(ctx context.Context, userID uint64, page, size int) ([]mktmodel.PointsLog, int64, error) {
	if page <= 0 { page = 1 }
	if size <= 0 || size > 100 { size = 20 }
	var total int64
	db.DB.WithContext(ctx).Model(&mktmodel.PointsLog{}).Where("user_id = ?", userID).Count(&total)
	var list []mktmodel.PointsLog
	err := db.DB.WithContext(ctx).Where("user_id = ?", userID).
		Order("id desc").Offset((page-1)*size).Limit(size).Find(&list).Error
	return list, total, err
}
