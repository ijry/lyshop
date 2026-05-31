package service

import (
	"context"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/plugins/distribution/model"
)

// GetConfig 获取分销配置
func GetConfig(ctx context.Context) (*model.DistributionConfig, error) {
	var config model.DistributionConfig
	err := db.DB.WithContext(ctx).First(&config).Error
	if err != nil {
		// 如果不存在，创建默认配置
		if err == db.ErrRecordNotFound {
			config = model.DistributionConfig{
				Enabled:         true,
				Level:           2,
				Level1Rate:      10,
				Level2Rate:      5,
				Level3Rate:      2,
				MinWithdraw:     100,
				WithdrawFeeRate: 0,
				AutoApprove:     false,
				RequireRealName: true,
			}
			if err := db.DB.WithContext(ctx).Create(&config).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &config, nil
}

// UpdateConfig 更新分销配置
func UpdateConfig(ctx context.Context, updates map[string]interface{}) error {
	config, err := GetConfig(ctx)
	if err != nil {
		return err
	}
	return db.DB.WithContext(ctx).Model(config).Updates(updates).Error
}

// GetDistributor 获取分销商信息
func GetDistributor(ctx context.Context, userID uint64) (*model.Distributor, error) {
	var distributor model.Distributor
	err := db.DB.WithContext(ctx).Where("user_id = ?", userID).First(&distributor).Error
	return &distributor, err
}

// GetDistributorByID 根据ID获取分销商
func GetDistributorByID(ctx context.Context, id uint64) (*model.Distributor, error) {
	var distributor model.Distributor
	err := db.DB.WithContext(ctx).First(&distributor, id).Error
	return &distributor, err
}

// CreateDistributor 创建分销商
func CreateDistributor(ctx context.Context, distributor *model.Distributor) error {
	return db.DB.WithContext(ctx).Create(distributor).Error
}

// UpdateDistributor 更新分销商信息
func UpdateDistributor(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&model.Distributor{}).Where("id = ?", id).Updates(updates).Error
}

// ListDistributors 获取分销商列表
func ListDistributors(ctx context.Context, page, size int, status string) ([]model.Distributor, int64, error) {
	var distributors []model.Distributor
	var total int64

	query := db.DB.WithContext(ctx).Model(&model.Distributor{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("id DESC").Offset(offset).Limit(size).Find(&distributors).Error; err != nil {
		return nil, 0, err
	}

	return distributors, total, nil
}

// GetDistributorChain 获取分销商上级链
func GetDistributorChain(ctx context.Context, userID uint64, maxLevel int) ([]model.Distributor, error) {
	var chain []model.Distributor
	currentUserID := userID

	for i := 0; i < maxLevel; i++ {
		var distributor model.Distributor
		err := db.DB.WithContext(ctx).Where("user_id = ?", currentUserID).First(&distributor).Error
		if err != nil {
			break
		}

		if distributor.ParentID == 0 {
			break
		}

		var parent model.Distributor
		err = db.DB.WithContext(ctx).First(&parent, distributor.ParentID).Error
		if err != nil {
			break
		}

		chain = append(chain, parent)
		currentUserID = parent.UserID
	}

	return chain, nil
}

// AddEarnings 增加分销商收益
func AddEarnings(ctx context.Context, distributorID uint64, amount float64) error {
	return db.DB.WithContext(ctx).Model(&model.Distributor{}).
		Where("id = ?", distributorID).
		Updates(map[string]interface{}{
			"total_earnings":   db.Expr("total_earnings + ?", amount),
			"available_amount": db.Expr("available_amount + ?", amount),
		}).Error
}

// FreezeAmount 冻结金额
func FreezeAmount(ctx context.Context, distributorID uint64, amount float64) error {
	return db.DB.WithContext(ctx).Model(&model.Distributor{}).
		Where("id = ? AND available_amount >= ?", distributorID, amount).
		Updates(map[string]interface{}{
			"available_amount": db.Expr("available_amount - ?", amount),
			"frozen_amount":    db.Expr("frozen_amount + ?", amount),
		}).Error
}

// UnfreezeAmount 解冻金额
func UnfreezeAmount(ctx context.Context, distributorID uint64, amount float64) error {
	return db.DB.WithContext(ctx).Model(&model.Distributor{}).
		Where("id = ?", distributorID).
		Updates(map[string]interface{}{
			"available_amount": db.Expr("available_amount + ?", amount),
			"frozen_amount":    db.Expr("frozen_amount - ?", amount),
		}).Error
}

// DeductFrozenAmount 扣除冻结金额
func DeductFrozenAmount(ctx context.Context, distributorID uint64, amount float64) error {
	return db.DB.WithContext(ctx).Model(&model.Distributor{}).
		Where("id = ? AND frozen_amount >= ?", distributorID, amount).
		Updates(map[string]interface{}{
			"frozen_amount":    db.Expr("frozen_amount - ?", amount),
			"withdrawn_amount": db.Expr("withdrawn_amount + ?", amount),
		}).Error
}
