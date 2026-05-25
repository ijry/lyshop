package service

import (
	"context"
	"errors"
	"time"

	"github.com/ijry/lyshop/core/db"
	vipmodel "github.com/ijry/lyshop/plugins/vip/model"
	"gorm.io/gorm"
)

func ListPlans(ctx context.Context, page, size int) ([]vipmodel.Plan, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&vipmodel.Plan{})
	var total int64
	tx.Count(&total)
	var list []vipmodel.Plan
	err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func CreatePlan(ctx context.Context, row *vipmodel.Plan) error {
	if row.DurationMonths <= 0 {
		row.DurationMonths = 12
	}
	return db.DB.WithContext(ctx).Create(row).Error
}

func ListLevels(ctx context.Context, page, size int) ([]vipmodel.Level, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&vipmodel.Level{})
	var total int64
	tx.Count(&total)
	var list []vipmodel.Level
	err := tx.Order("sort asc, id asc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func CreateLevel(ctx context.Context, row *vipmodel.Level) error {
	return db.DB.WithContext(ctx).Create(row).Error
}

func ListCouponRules(ctx context.Context, page, size int) ([]vipmodel.CouponRule, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&vipmodel.CouponRule{})
	var total int64
	tx.Count(&total)
	var list []vipmodel.CouponRule
	err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func CreateCouponRule(ctx context.Context, row *vipmodel.CouponRule) error {
	if row.MonthlyQuota <= 0 {
		row.MonthlyQuota = 1
	}
	if row.ClaimMode == "" {
		row.ClaimMode = "manual"
	}
	return db.DB.WithContext(ctx).Create(row).Error
}

func ListSkuPrices(ctx context.Context, page, size int) ([]vipmodel.SkuPrice, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&vipmodel.SkuPrice{})
	var total int64
	tx.Count(&total)
	var list []vipmodel.SkuPrice
	err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func CreateSkuPrice(ctx context.Context, row *vipmodel.SkuPrice) error {
	if row.VipPrice <= 0 && row.VipDiscountRate <= 0 {
		return errors.New("会员价或折扣率至少填写一项")
	}
	return db.DB.WithContext(ctx).Create(row).Error
}

func OpenMembership(ctx context.Context, userID, planID uint64, now time.Time) (*vipmodel.UserAsset, error) {
	var plan vipmodel.Plan
	if err := db.DB.WithContext(ctx).Where("id = ? AND status = 1", planID).First(&plan).Error; err != nil {
		return nil, errors.New("会员计划不存在")
	}
	var asset vipmodel.UserAsset
	if err := db.DB.WithContext(ctx).Where("user_id = ?", userID).First(&asset).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		asset.UserID = userID
	}

	startAt := now
	if asset.VipEndAt != nil && asset.VipEndAt.After(now) {
		startAt = *asset.VipEndAt
	}
	endAt := startAt.AddDate(0, plan.DurationMonths, 0)

	asset.CurrentPlanID = plan.ID
	asset.VipStartAt = &startAt
	asset.VipEndAt = &endAt
	asset.Status = 1
	if err := db.DB.WithContext(ctx).Save(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func GetActiveAsset(ctx context.Context, userID uint64, now time.Time) (*vipmodel.UserAsset, error) {
	var asset vipmodel.UserAsset
	if err := db.DB.WithContext(ctx).Where("user_id = ?", userID).First(&asset).Error; err != nil {
		return nil, err
	}
	if asset.Status != 1 || asset.VipEndAt == nil || now.After(*asset.VipEndAt) {
		if asset.Status != 0 {
			asset.Status = 0
			_ = db.DB.WithContext(ctx).Model(&vipmodel.UserAsset{}).Where("id = ?", asset.ID).Update("status", 0).Error
		}
		return nil, errors.New("会员已失效")
	}
	return &asset, nil
}

type Profile struct {
	IsVIP       bool       `json:"is_vip"`
	LevelName   string     `json:"level_name"`
	GrowthValue int64      `json:"growth_value"`
	ExpireAt    *time.Time `json:"expire_at"`
}

func GetProfile(ctx context.Context, userID uint64, now time.Time) (Profile, error) {
	var asset vipmodel.UserAsset
	if err := db.DB.WithContext(ctx).Where("user_id = ?", userID).First(&asset).Error; err != nil {
		return Profile{IsVIP: false, GrowthValue: 0}, nil
	}
	profile := Profile{
		IsVIP:       asset.Status == 1 && asset.VipEndAt != nil && !now.After(*asset.VipEndAt),
		GrowthValue: asset.GrowthValue,
		ExpireAt:    asset.VipEndAt,
	}
	if asset.CurrentLevelID > 0 {
		var level vipmodel.Level
		if err := db.DB.WithContext(ctx).Where("id = ?", asset.CurrentLevelID).First(&level).Error; err == nil {
			profile.LevelName = level.Name
		}
	}
	return profile, nil
}

func ListGrowthLogs(ctx context.Context, userID uint64, page, size int) ([]vipmodel.GrowthLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&vipmodel.GrowthLog{}).Where("user_id = ?", userID)
	var total int64
	tx.Count(&total)
	var list []vipmodel.GrowthLog
	err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}
