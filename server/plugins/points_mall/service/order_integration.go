package service

import (
	"context"
	"fmt"

	"github.com/ijry/lyshop/server/core/db"
	pmservice "github.com/ijry/lyshop/server/plugins/points_mall/service"
)

// GrantOrderPoints 订单完成后赠送积分
// 应该在订单状态变更为已完成时调用
func GrantOrderPoints(ctx context.Context, orderID uint64) error {
	// 查询订单信息
	var order struct {
		ID          uint64
		UserID      uint64
		TotalAmount float64
		Status      string
	}

	if err := db.DB.WithContext(ctx).Table("orders").
		Select("id, user_id, total_amount, status").
		Where("id = ?", orderID).
		First(&order).Error; err != nil {
		return err
	}

	// 只有已完成的订单才赠送积分
	if order.Status != "completed" {
		return nil
	}

	// 检查是否已经赠送过积分（避免重复赠送）
	var existingLog struct {
		ID uint64
	}
	err := db.DB.WithContext(ctx).Table("points_logs").
		Select("id").
		Where("user_id = ? AND type = 4 AND related_id = ?", order.UserID, orderID).
		First(&existingLog).Error

	if err == nil {
		// 已经赠送过积分，跳过
		return nil
	}

	// 获取积分配置
	var config struct {
		EnableOrderPoints bool
		OrderPointsRate   float64
	}
	// TODO: 从配置表读取，这里使用默认值
	config.EnableOrderPoints = true
	config.OrderPointsRate = 0.01 // 1% 即消费100元送100积分

	if !config.EnableOrderPoints {
		return nil
	}

	// 计算赠送积分
	points := int(order.TotalAmount * config.OrderPointsRate * 100)
	if points <= 0 {
		return nil
	}

	// 赠送积分
	remark := fmt.Sprintf("订单完成奖励（订单号：%d）", orderID)
	return pmservice.AddPointsWithRelated(ctx, order.UserID, points, 4, orderID, remark)
}
