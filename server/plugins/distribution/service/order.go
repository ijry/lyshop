package service

import (
	"context"
	"time"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/plugins/distribution/model"
)

// CreateDistributionOrder 创建分销订单
func CreateDistributionOrder(ctx context.Context, order *model.DistributionOrder) error {
	return db.DB.WithContext(ctx).Create(order).Error
}

// SettleDistributionOrder 结算分销订单
func SettleDistributionOrder(ctx context.Context, orderID uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *db.Tx) error {
		// 获取所有待结算的分销订单
		var orders []model.DistributionOrder
		if err := tx.Where("order_id = ? AND status = ?", orderID, "pending").Find(&orders).Error; err != nil {
			return err
		}

		now := time.Now()
		for _, order := range orders {
			// 更新订单状态
			order.Status = "settled"
			order.SettledAt = &now
			if err := tx.Save(&order).Error; err != nil {
				return err
			}

			// 增加分销商收益
			if err := AddEarnings(ctx, order.DistributorID, order.Commission); err != nil {
				return err
			}

			// 更新分销商订单统计
			if err := tx.Model(&model.Distributor{}).
				Where("id = ?", order.DistributorID).
				UpdateColumn("total_orders", db.Expr("total_orders + 1")).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// CancelDistributionOrder 取消分销订单
func CancelDistributionOrder(ctx context.Context, orderID uint64) error {
	return db.DB.WithContext(ctx).Model(&model.DistributionOrder{}).
		Where("order_id = ? AND status = ?", orderID, "pending").
		Update("status", "cancelled").Error
}

// ListDistributionOrders 获取分销订单列表
func ListDistributionOrders(ctx context.Context, page, size int, status string, distributorID uint64) ([]model.DistributionOrder, int64, error) {
	var orders []model.DistributionOrder
	var total int64

	query := db.DB.WithContext(ctx).Model(&model.DistributionOrder{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if distributorID > 0 {
		query = query.Where("distributor_id = ?", distributorID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("id DESC").Offset(offset).Limit(size).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetDistributionOrder 获取分销订单详情
func GetDistributionOrder(ctx context.Context, id uint64) (*model.DistributionOrder, error) {
	var order model.DistributionOrder
	err := db.DB.WithContext(ctx).First(&order, id).Error
	return &order, err
}

// CalculateCommission 计算佣金
func CalculateCommission(ctx context.Context, orderAmount float64, buyerUserID uint64) ([]model.DistributionOrder, error) {
	// 获取配置
	config, err := GetConfig(ctx)
	if err != nil || !config.Enabled {
		return nil, nil
	}

	// 获取买家的上级分销商链
	chain, err := GetDistributorChain(ctx, buyerUserID, config.Level)
	if err != nil {
		return nil, err
	}

	var orders []model.DistributionOrder
	rates := []float64{config.Level1Rate, config.Level2Rate, config.Level3Rate}

	for i, distributor := range chain {
		if i >= config.Level {
			break
		}

		if distributor.Status != "active" {
			continue
		}

		rate := rates[i]
		commission := orderAmount * rate / 100

		orders = append(orders, model.DistributionOrder{
			DistributorID:  distributor.ID,
			Level:          i + 1,
			OrderAmount:    orderAmount,
			CommissionRate: rate,
			Commission:     commission,
			Status:         "pending",
		})
	}

	return orders, nil
}
