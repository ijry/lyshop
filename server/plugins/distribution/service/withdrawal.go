package service

import (
	"context"
	"errors"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/server/plugins/distribution/model"
)

var (
	ErrInsufficientBalance = errors.New("余额不足")
	ErrBelowMinWithdraw    = errors.New("低于最低提现金额")
)

// CreateWithdrawal 创建提现申请
func CreateWithdrawal(ctx context.Context, withdrawal *model.DistributionWithdrawal) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *db.DB) error {
		// 获取分销商信息
		var distributor model.Distributor
		if err := tx.First(&distributor, withdrawal.DistributorID).Error; err != nil {
			return err
		}

		// 检查余额
		if distributor.AvailableAmount < withdrawal.Amount {
			return ErrInsufficientBalance
		}

		// 获取配置
		var config model.DistributionConfig
		if err := tx.First(&config).Error; err != nil {
			return err
		}

		// 检查最低提现金额
		if withdrawal.Amount < config.MinWithdraw {
			return ErrBelowMinWithdraw
		}

		// 计算手续费
		withdrawal.Fee = withdrawal.Amount * config.WithdrawFeeRate / 100
		withdrawal.ActualAmount = withdrawal.Amount - withdrawal.Fee
		withdrawal.Status = "pending"

		// 创建提现记录
		if err := tx.Create(withdrawal).Error; err != nil {
			return err
		}

		// 冻结金额
		return FreezeAmount(ctx, withdrawal.DistributorID, withdrawal.Amount)
	})
}

// ApproveWithdrawal 审核通过提现
func ApproveWithdrawal(ctx context.Context, id uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *db.DB) error {
		var withdrawal model.DistributionWithdrawal
		if err := tx.First(&withdrawal, id).Error; err != nil {
			return err
		}

		if withdrawal.Status != "pending" {
			return errors.New("提现状态不正确")
		}

		now := model.Now()
		withdrawal.Status = "approved"
		withdrawal.ProcessedAt = &now

		return tx.Save(&withdrawal).Error
	})
}

// RejectWithdrawal 拒绝提现
func RejectWithdrawal(ctx context.Context, id uint64, reason string) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *db.DB) error {
		var withdrawal model.DistributionWithdrawal
		if err := tx.First(&withdrawal, id).Error; err != nil {
			return err
		}

		if withdrawal.Status != "pending" {
			return errors.New("提现状态不正确")
		}

		now := model.Now()
		withdrawal.Status = "rejected"
		withdrawal.RejectReason = reason
		withdrawal.ProcessedAt = &now

		if err := tx.Save(&withdrawal).Error; err != nil {
			return err
		}

		// 解冻金额
		return UnfreezeAmount(ctx, withdrawal.DistributorID, withdrawal.Amount)
	})
}

// CompleteWithdrawal 完成提现
func CompleteWithdrawal(ctx context.Context, id uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *db.DB) error {
		var withdrawal model.DistributionWithdrawal
		if err := tx.First(&withdrawal, id).Error; err != nil {
			return err
		}

		if withdrawal.Status != "approved" {
			return errors.New("提现状态不正确")
		}

		now := model.Now()
		withdrawal.Status = "completed"
		withdrawal.CompletedAt = &now

		if err := tx.Save(&withdrawal).Error; err != nil {
			return err
		}

		// 扣除冻结金额
		return DeductFrozenAmount(ctx, withdrawal.DistributorID, withdrawal.Amount)
	})
}

// ListWithdrawals 获取提现列表
func ListWithdrawals(ctx context.Context, page, size int, status string, distributorID uint64) ([]model.DistributionWithdrawal, int64, error) {
	var withdrawals []model.DistributionWithdrawal
	var total int64

	query := db.DB.WithContext(ctx).Model(&model.DistributionWithdrawal{})
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
	if err := query.Order("id DESC").Offset(offset).Limit(size).Find(&withdrawals).Error; err != nil {
		return nil, 0, err
	}

	return withdrawals, total, nil
}

// GetWithdrawal 获取提现详情
func GetWithdrawal(ctx context.Context, id uint64) (*model.DistributionWithdrawal, error) {
	var withdrawal model.DistributionWithdrawal
	err := db.DB.WithContext(ctx).First(&withdrawal, id).Error
	return &withdrawal, err
}
