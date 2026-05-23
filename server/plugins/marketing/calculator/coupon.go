package calculator

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/marketing"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
)

// CouponCalculator applies user-selected coupons.
// Respects Stackable flag: non-stackable coupons are mutually exclusive (pick best).
type CouponCalculator struct{}

func init() { marketing.Register(&CouponCalculator{}) }

func (c *CouponCalculator) Name() string  { return "coupon" }
func (c *CouponCalculator) Priority() int { return 30 }

func (c *CouponCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
	if db.DB == nil || len(ctx.CouponIDs) == 0 {
		return true, nil
	}

	now := time.Now()
	subtotal := ctx.GoodsAmount - ctx.ActivityDiscount - ctx.FullReduceDiscount

	// Load coupon_users + coupons
	var couponUsers []mktmodel.CouponUser
	db.DB.Where("id IN ? AND user_id = ? AND status = 1", ctx.CouponIDs, ctx.UserID).Find(&couponUsers)

	var couponIDs []uint64
	for _, cu := range couponUsers {
		couponIDs = append(couponIDs, cu.CouponID)
	}
	var coupons []mktmodel.Coupon
	db.DB.Where("id IN ? AND status = 1", couponIDs).Find(&coupons)

	couponMap := map[uint64]mktmodel.Coupon{}
	for _, cp := range coupons {
		couponMap[cp.ID] = cp
	}

	// Separate stackable vs non-stackable
	type couponCandidate struct {
		couponUser mktmodel.CouponUser
		coupon     mktmodel.Coupon
		discount   float64
	}
	var stackable, exclusive []couponCandidate

	for _, cu := range couponUsers {
		cp, ok := couponMap[cu.CouponID]
		if !ok {
			continue
		}
		// Check validity
		if cp.StartAt != nil && now.Before(*cp.StartAt) {
			continue
		}
		if cp.EndAt != nil && now.After(*cp.EndAt) {
			continue
		}
		// Check min amount
		if subtotal < cp.MinAmount {
			continue
		}
		// Check scope
		if cp.ExcludeActivity && ctx.ActivityDiscount > 0 {
			continue
		}
		if cp.Scope != "" && cp.Scope != "all" {
			var scopeIDs []uint64
			json.Unmarshal([]byte(cp.ScopeIDs), &scopeIDs)
			if len(scopeIDs) > 0 && !itemsMatchScope(ctx.Items, cp.Scope, scopeIDs) {
				continue
			}
		}

		// Calculate discount
		var disc float64
		switch cp.Type {
		case mktmodel.CouponTypeFullReduce:
			disc = cp.Discount
		case mktmodel.CouponTypeDiscount:
			disc = subtotal * (1 - cp.Discount) // e.g. 0.9 → 10% off
		case mktmodel.CouponTypeFree:
			disc = cp.Discount
		}
		if disc <= 0 {
			continue
		}

		cc := couponCandidate{couponUser: cu, coupon: cp, discount: disc}
		if cp.Stackable {
			stackable = append(stackable, cc)
		} else {
			exclusive = append(exclusive, cc)
		}
	}

	// Non-stackable: pick the best one
	var bestExclusive *couponCandidate
	for i := range exclusive {
		if bestExclusive == nil || exclusive[i].discount > bestExclusive.discount {
			bestExclusive = &exclusive[i]
		}
	}

	totalDiscount := 0.0
	if bestExclusive != nil {
		totalDiscount += bestExclusive.discount
		ctx.AppliedRules = append(ctx.AppliedRules, marketing.AppliedRule{
			Type: "coupon", Name: bestExclusive.coupon.Name, Discount: bestExclusive.discount,
		})
	}
	for _, sc := range stackable {
		totalDiscount += sc.discount
		ctx.AppliedRules = append(ctx.AppliedRules, marketing.AppliedRule{
			Type: "coupon", Name: sc.coupon.Name, Discount: sc.discount,
		})
	}

	ctx.CouponDiscount += totalDiscount
	return true, nil
}

func itemsMatchScope(items []marketing.OrderItem, scope string, scopeIDs []uint64) bool {
	idSet := map[uint64]bool{}
	for _, id := range scopeIDs {
		idSet[id] = true
	}
	for _, item := range items {
		switch scope {
		case "product":
			if idSet[item.ProductID] {
				return true
			}
		}
	}
	return false
}
