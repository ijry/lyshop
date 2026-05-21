package service

import (
	"context"
	"errors"

	"github.com/ijry/lyshop/core/db"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
)

func ListCoupons(ctx context.Context, page, size int) ([]mktmodel.Coupon, int64, error) {
	if page <= 0 { page = 1 }
	if size <= 0 || size > 100 { size = 20 }
	var total int64
	db.DB.WithContext(ctx).Model(&mktmodel.Coupon{}).Count(&total)
	var list []mktmodel.Coupon
	err := db.DB.WithContext(ctx).Order("id desc").
		Offset((page-1)*size).Limit(size).Find(&list).Error
	return list, total, err
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

// ListUserCoupons returns a user's coupons.
func ListUserCoupons(ctx context.Context, userID uint64) ([]mktmodel.CouponUser, error) {
	var list []mktmodel.CouponUser
	err := db.DB.WithContext(ctx).Where("user_id = ?", userID).
		Order("id desc").Find(&list).Error
	return list, err
}
