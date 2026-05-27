package service

import (
	"context"
	"errors"
	"time"

	"github.com/ijry/lyshop/core/db"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
)

type ClaimableCoupon struct {
	mktmodel.Coupon
	ClaimedCount int64 `json:"claimed_count"`
	ClaimedByMe  int64 `json:"claimed_by_me"`
	CanClaim     bool  `json:"can_claim"`
}

type UserCoupon struct {
	mktmodel.CouponUser
	Coupon *mktmodel.Coupon `json:"coupon,omitempty"`
}

func ListCoupons(ctx context.Context, page, size int) ([]mktmodel.Coupon, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	var total int64
	db.DB.WithContext(ctx).Model(&mktmodel.Coupon{}).Count(&total)
	var list []mktmodel.Coupon
	err := db.DB.WithContext(ctx).Order("id desc").
		Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func ListClaimableCoupons(ctx context.Context, userID uint64, page, size int) ([]ClaimableCoupon, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	now := time.Now()
	var list []mktmodel.Coupon
	err := db.DB.WithContext(ctx).
		Where("status = 1").
		Where("start_at IS NULL OR start_at <= ?", now).
		Where("end_at IS NULL OR end_at >= ?", now).
		Order("id desc").
		Offset((page - 1) * size).
		Limit(size).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	result := make([]ClaimableCoupon, 0, len(list))
	for _, item := range list {
		var claimedCount int64
		db.DB.WithContext(ctx).Model(&mktmodel.CouponUser{}).
			Where("coupon_id = ?", item.ID).Count(&claimedCount)

		var claimedByMe int64
		db.DB.WithContext(ctx).Model(&mktmodel.CouponUser{}).
			Where("coupon_id = ? AND user_id = ?", item.ID, userID).Count(&claimedByMe)

		canClaim := true
		if item.PerLimit > 0 && int(claimedByMe) >= item.PerLimit {
			canClaim = false
		}
		if item.TotalCount > 0 && int(claimedCount) >= item.TotalCount {
			canClaim = false
		}

		result = append(result, ClaimableCoupon{
			Coupon:       item,
			ClaimedCount: claimedCount,
			ClaimedByMe:  claimedByMe,
			CanClaim:     canClaim,
		})
	}

	return result, nil
}

func CreateCoupon(ctx context.Context, c *mktmodel.Coupon) error {
	return db.DB.WithContext(ctx).Create(c).Error
}

// ClaimCoupon lets a user claim a coupon.
func ClaimCoupon(ctx context.Context, couponID, userID uint64) error {
	var coupon mktmodel.Coupon
	if err := db.DB.WithContext(ctx).First(&coupon, couponID).Error; err != nil {
		return errors.New("优惠券不存在")
	}
	if coupon.Status != 1 {
		return errors.New("优惠券已停用")
	}
	// Check per-user limit
	var count int64
	db.DB.WithContext(ctx).Model(&mktmodel.CouponUser{}).
		Where("coupon_id = ? AND user_id = ?", couponID, userID).Count(&count)
	if coupon.PerLimit > 0 && int(count) >= coupon.PerLimit {
		return errors.New("已达领取上限")
	}
	return db.DB.WithContext(ctx).Create(&mktmodel.CouponUser{
		CouponID: couponID, UserID: userID, Status: 1,
	}).Error
}

// ListUserCoupons returns a user's coupons with coupon details.
func ListUserCoupons(ctx context.Context, userID uint64) ([]UserCoupon, error) {
	var rows []mktmodel.CouponUser
	if err := db.DB.WithContext(ctx).Where("user_id = ?", userID).
		Order("id desc").Find(&rows).Error; err != nil {
		return nil, err
	}

	idSet := map[uint64]struct{}{}
	for _, row := range rows {
		if row.CouponID > 0 {
			idSet[row.CouponID] = struct{}{}
		}
	}

	couponMap := map[uint64]mktmodel.Coupon{}
	if len(idSet) > 0 {
		ids := make([]uint64, 0, len(idSet))
		for id := range idSet {
			ids = append(ids, id)
		}
		var coupons []mktmodel.Coupon
		if err := db.DB.WithContext(ctx).Where("id IN ?", ids).Find(&coupons).Error; err != nil {
			return nil, err
		}
		for _, item := range coupons {
			couponMap[item.ID] = item
		}
	}

	list := make([]UserCoupon, 0, len(rows))
	for _, row := range rows {
		result := UserCoupon{CouponUser: row}
		if c, ok := couponMap[row.CouponID]; ok {
			coupon := c
			result.Coupon = &coupon
		}
		list = append(list, result)
	}
	return list, nil
}
