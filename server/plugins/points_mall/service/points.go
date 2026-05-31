package service

import (
	"context"
	"fmt"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/model"
	pmmodel "github.com/ijry/lyshop/plugins/points_mall/model"
)

// AddPoints 添加或扣减积分
// points 为正数表示增加，负数表示减少
// logType: 1=签到 2=订单抵扣 3=兑换消耗 4=订单完成 5=管理员调整 6=过期扣除 7=活动奖励
func AddPoints(ctx context.Context, userID uint64, points int, logType int8, remark string) error {
	return AddPointsWithRelated(ctx, userID, points, logType, 0, remark)
}

// AddPointsWithRelated 添加或扣减积分（带关联ID）
func AddPointsWithRelated(ctx context.Context, userID uint64, points int, logType int8, relatedID uint64, remark string) error {
	if points == 0 {
		return nil
	}

	return db.DB.WithContext(ctx).Transaction(func(tx *db.Tx) error {
		// 更新用户积分
		result := tx.Model(&model.User{}).
			Where("id = ?", userID).
			Update("points", db.Expr("points + ?", points))

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("user not found: %d", userID)
		}

		// 记录积分日志
		log := &pmmodel.PointsLog{
			UserID:    userID,
			Type:      logType,
			Points:    points,
			RelatedID: relatedID,
			Remark:    remark,
		}

		return tx.Create(log).Error
	})
}

// GetUserPoints 获取用户当前积分
func GetUserPoints(ctx context.Context, userID uint64) (int, error) {
	var user model.User
	err := db.DB.WithContext(ctx).Select("points").First(&user, userID).Error
	if err != nil {
		return 0, err
	}
	return user.Points, nil
}

// ListPointsLogs 获取用户积分日志
func ListPointsLogs(ctx context.Context, userID uint64, page, size int) ([]pmmodel.PointsLog, int64, error) {
	var logs []pmmodel.PointsLog
	var total int64

	query := db.DB.WithContext(ctx).Model(&pmmodel.PointsLog{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := query.Order("id DESC").Offset(offset).Limit(size).Find(&logs).Error

	return logs, total, err
}

// AdminListPointsLogs 管理员获取积分日志（支持筛选）
func AdminListPointsLogs(ctx context.Context, userID uint64, logType int8, page, size int) ([]pmmodel.PointsLog, int64, error) {
	var logs []pmmodel.PointsLog
	var total int64

	query := db.DB.WithContext(ctx).Model(&pmmodel.PointsLog{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if logType > 0 {
		query = query.Where("type = ?", logType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := query.Order("id DESC").Offset(offset).Limit(size).Find(&logs).Error

	return logs, total, err
}

// GetPointsStats 获取积分统计
func GetPointsStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 累计发放积分
	var totalIssued int64
	db.DB.WithContext(ctx).Model(&pmmodel.PointsLog{}).
		Where("points > 0").
		Select("COALESCE(SUM(points), 0)").
		Scan(&totalIssued)
	stats["total_issued"] = totalIssued

	// 累计消耗积分
	var totalConsumed int64
	db.DB.WithContext(ctx).Model(&pmmodel.PointsLog{}).
		Where("points < 0").
		Select("COALESCE(SUM(ABS(points)), 0)").
		Scan(&totalConsumed)
	stats["total_consumed"] = totalConsumed

	// 当前余额
	var totalBalance int64
	db.DB.WithContext(ctx).Model(&model.User{}).
		Select("COALESCE(SUM(points), 0)").
		Scan(&totalBalance)
	stats["total_balance"] = totalBalance

	// 商品数量
	var productCount int64
	db.DB.WithContext(ctx).Model(&pmmodel.PointsProduct{}).Count(&productCount)
	stats["product_count"] = productCount

	// 兑换次数
	var exchangeCount int64
	db.DB.WithContext(ctx).Model(&pmmodel.PointsExchange{}).Count(&exchangeCount)
	stats["exchange_count"] = exchangeCount

	// 今日发放
	var todayIssued int64
	db.DB.WithContext(ctx).Model(&pmmodel.PointsLog{}).
		Where("points > 0 AND DATE(created_at) = CURDATE()").
		Select("COALESCE(SUM(points), 0)").
		Scan(&todayIssued)
	stats["today_issued"] = todayIssued

	// 今日消耗
	var todayConsumed int64
	db.DB.WithContext(ctx).Model(&pmmodel.PointsLog{}).
		Where("points < 0 AND DATE(created_at) = CURDATE()").
		Select("COALESCE(SUM(ABS(points)), 0)").
		Scan(&todayConsumed)
	stats["today_consumed"] = todayConsumed

	return stats, nil
}
