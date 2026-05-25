package service

import (
	"context"
	"errors"
	"time"

	"github.com/ijry/lyshop/core/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
	vipmodel "github.com/ijry/lyshop/plugins/vip/model"
)

type MonthlyCouponView struct {
	RuleID       uint64 `json:"rule_id"`
	Name         string `json:"name"`
	CouponName   string `json:"coupon_name"`
	Claimed      int    `json:"claimed"`
	MonthlyLimit int    `json:"monthly_limit"`
}

func ListMonthlyCoupons(ctx context.Context, userID uint64, now time.Time) ([]MonthlyCouponView, error) {
	asset, err := GetActiveAsset(ctx, userID, now)
	if err != nil || asset == nil {
		return []MonthlyCouponView{}, nil
	}

	var rules []vipmodel.CouponRule
	err = db.DB.WithContext(ctx).
		Where("status = 1 AND claim_mode = 'manual' AND (plan_id = 0 OR plan_id = ?) AND (level_id = 0 OR level_id = ?)",
			asset.CurrentPlanID, asset.CurrentLevelID).
		Order("id desc").Find(&rules).Error
	if err != nil || len(rules) == 0 {
		return []MonthlyCouponView{}, err
	}

	period := now.Format("200601")
	ruleIDs := make([]uint64, 0, len(rules))
	couponIDs := make([]uint64, 0, len(rules))
	for _, r := range rules {
		ruleIDs = append(ruleIDs, r.ID)
		couponIDs = append(couponIDs, r.CouponID)
	}

	var claims []vipmodel.CouponClaim
	db.DB.WithContext(ctx).
		Where("user_id = ? AND period_yyyymm = ? AND rule_id IN ?", userID, period, ruleIDs).
		Find(&claims)
	claimMap := map[uint64]int{}
	for _, c := range claims {
		claimMap[c.RuleID] = c.ClaimedCount
	}

	var coupons []mktmodel.Coupon
	db.DB.WithContext(ctx).Where("id IN ? AND status = 1", couponIDs).Find(&coupons)
	couponNameMap := map[uint64]string{}
	for _, cp := range coupons {
		couponNameMap[cp.ID] = cp.Name
	}

	list := make([]MonthlyCouponView, 0, len(rules))
	for _, r := range rules {
		list = append(list, MonthlyCouponView{
			RuleID:       r.ID,
			Name:         "会员月度券",
			CouponName:   couponNameMap[r.CouponID],
			Claimed:      claimMap[r.ID],
			MonthlyLimit: r.MonthlyQuota,
		})
	}
	return list, nil
}

func ClaimMonthlyCoupon(ctx context.Context, userID, ruleID uint64, now time.Time) error {
	asset, err := GetActiveAsset(ctx, userID, now)
	if err != nil || asset == nil {
		return errors.New("会员状态不可领取")
	}

	period := now.Format("200601")
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rule vipmodel.CouponRule
		if err := tx.Where("id = ? AND status = 1", ruleID).First(&rule).Error; err != nil {
			return errors.New("会员券规则不存在")
		}
		if rule.PlanID > 0 && rule.PlanID != asset.CurrentPlanID {
			return errors.New("会员计划不匹配")
		}
		if rule.LevelID > 0 && rule.LevelID != asset.CurrentLevelID {
			return errors.New("会员等级不匹配")
		}

		var claim vipmodel.CouponClaim
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND rule_id = ? AND period_yyyymm = ?", userID, ruleID, period).
			First(&claim).Error
		if err == nil && claim.ClaimedCount >= rule.MonthlyQuota {
			return errors.New("本月已领完")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			claim = vipmodel.CouponClaim{UserID: userID, RuleID: ruleID, PeriodYYYYMM: period}
		} else if err != nil {
			return err
		}

		claim.ClaimedCount++
		last := now
		claim.LastClaimedAt = &last
		if err := tx.Save(&claim).Error; err != nil {
			return err
		}
		return tx.Create(&mktmodel.CouponUser{
			CouponID: rule.CouponID,
			UserID:   userID,
			Status:   1,
		}).Error
	})
}
