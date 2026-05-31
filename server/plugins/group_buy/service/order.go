package service

import (
	"context"
	"errors"
	"time"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/server/plugins/group_buy/model"
)

var (
	ErrGroupFull    = errors.New("拼团已满")
	ErrGroupExpired = errors.New("拼团已过期")
	ErrAlreadyJoined = errors.New("已参加该拼团")
)

// CreateGroupBuyOrder 创建拼团订单（开团）
func CreateGroupBuyOrder(ctx context.Context, activityID, productID, skuID, leaderID uint64) (*model.GroupBuyOrder, error) {
	activity, err := GetActivity(ctx, activityID)
	if err != nil {
		return nil, err
	}

	expireAt := time.Now().Add(time.Duration(activity.ExpireHours) * time.Hour)

	groupOrder := &model.GroupBuyOrder{
		ActivityID:  activityID,
		ProductID:   productID,
		SkuID:       skuID,
		LeaderID:    leaderID,
		GroupSize:   activity.GroupSize,
		JoinedCount: 1,
		Status:      "pending",
		ExpireAt:    &expireAt,
	}

	if err := db.DB.WithContext(ctx).Create(groupOrder).Error; err != nil {
		return nil, err
	}

	return groupOrder, nil
}

// JoinGroupBuy 参加拼团
func JoinGroupBuy(ctx context.Context, groupOrderID, userID, orderID uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *db.DB) error {
		var groupOrder model.GroupBuyOrder
		if err := tx.Where("id = ?", groupOrderID).First(&groupOrder).Error; err != nil {
			return err
		}

		// 检查状态
		if groupOrder.Status != "pending" {
			return errors.New("拼团已结束")
		}

		// 检查是否过期
		if groupOrder.ExpireAt != nil && time.Now().After(*groupOrder.ExpireAt) {
			return ErrGroupExpired
		}

		// 检查是否已满
		if groupOrder.JoinedCount >= groupOrder.GroupSize {
			return ErrGroupFull
		}

		// 检查是否已参加
		var count int64
		if err := tx.Model(&model.GroupBuyMember{}).
			Where("group_order_id = ? AND user_id = ?", groupOrderID, userID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return ErrAlreadyJoined
		}

		// 添加成员
		member := &model.GroupBuyMember{
			GroupOrderID: groupOrderID,
			UserID:       userID,
			OrderID:      orderID,
			IsLeader:     false,
		}
		if err := tx.Create(member).Error; err != nil {
			return err
		}

		// 更新参团人数
		groupOrder.JoinedCount++
		if err := tx.Model(&groupOrder).Update("joined_count", groupOrder.JoinedCount).Error; err != nil {
			return err
		}

		// 检查是否成团
		if groupOrder.JoinedCount >= groupOrder.GroupSize {
			return CompleteGroupBuy(tx, groupOrderID)
		}

		return nil
	})
}

// CompleteGroupBuy 完成拼团
func CompleteGroupBuy(tx *db.DB, groupOrderID uint64) error {
	return tx.Model(&model.GroupBuyOrder{}).
		Where("id = ?", groupOrderID).
		Update("status", "success").Error
}

// ExpireGroupBuy 拼团过期处理
func ExpireGroupBuy(ctx context.Context, groupOrderID uint64) error {
	return db.DB.WithContext(ctx).Model(&model.GroupBuyOrder{}).
		Where("id = ?", groupOrderID).
		Where("status = ?", "pending").
		Update("status", "failed").Error
}

// ListGroupOrders 获取拼团订单列表
func ListGroupOrders(ctx context.Context, page, size int, status string) ([]model.GroupBuyOrder, int64, error) {
	var orders []model.GroupBuyOrder
	var total int64

	query := db.DB.WithContext(ctx).Model(&model.GroupBuyOrder{})
	if status != "" {
		query = query.Where("status = ?", status)
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

// GetGroupOrder 获取拼团订单详情
func GetGroupOrder(ctx context.Context, id uint64) (*model.GroupBuyOrder, error) {
	var order model.GroupBuyOrder
	if err := db.DB.WithContext(ctx).First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// ListGroupMembers 获取拼团成员列表
func ListGroupMembers(ctx context.Context, groupOrderID uint64) ([]model.GroupBuyMember, error) {
	var members []model.GroupBuyMember
	err := db.DB.WithContext(ctx).
		Where("group_order_id = ?", groupOrderID).
		Order("is_leader DESC, id ASC").
		Find(&members).Error
	return members, err
}

// CheckExpiredGroups 检查并处理过期拼团（定时任务）
func CheckExpiredGroups(ctx context.Context) error {
	now := time.Now()
	return db.DB.WithContext(ctx).Model(&model.GroupBuyOrder{}).
		Where("status = ?", "pending").
		Where("expire_at < ?", now).
		Update("status", "failed").Error
}
