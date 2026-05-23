package calculator

import (
	"time"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/marketing"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
)

// ActivityCalculator handles seckill, group_buy, bargain, and custom activities.
// All share the same priority (10) — only one activity applies per item.
type ActivityCalculator struct{}

func init() { marketing.Register(&ActivityCalculator{}) }

func (c *ActivityCalculator) Name() string  { return "activity" }
func (c *ActivityCalculator) Priority() int { return 10 }

func (c *ActivityCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
	if db.DB == nil {
		return true, nil // no DB in test mode
	}

	now := time.Now()
	var activities []mktmodel.Activity
	db.DB.Where("status = 1 AND start_at <= ? AND end_at >= ?", now, now).Find(&activities)
	if len(activities) == 0 {
		return true, nil
	}

	// Build activity→products map
	var actIDs []uint64
	for _, a := range activities {
		actIDs = append(actIDs, a.ID)
	}
	var actProducts []mktmodel.ActivityProduct
	db.DB.Where("activity_id IN ?", actIDs).Find(&actProducts)

	actProdMap := map[uint64]map[uint64]mktmodel.ActivityProduct{} // actID → skuID → product
	for _, ap := range actProducts {
		if actProdMap[ap.ActivityID] == nil {
			actProdMap[ap.ActivityID] = map[uint64]mktmodel.ActivityProduct{}
		}
		actProdMap[ap.ActivityID][ap.SkuID] = ap
	}

	exclusive := false
	for i := range ctx.Items {
		item := &ctx.Items[i]
		for _, act := range activities {
			prods := actProdMap[act.ID]
			ap, hit := prods[item.SkuID]
			if !hit {
				ap, hit = prods[0] // sku_id=0 means all SKUs of that product
			}
			if !hit {
				continue
			}
			if ap.ProductID != item.ProductID {
				continue
			}

			// Calculate activity price
			var actPrice float64
			switch act.PriceRule {
			case mktmodel.PriceRuleFixedPrice:
				actPrice = ap.ActivityPrice
				if actPrice == 0 {
					actPrice = act.PriceValue
				}
			case mktmodel.PriceRuleDiscountRate:
				actPrice = item.Price * act.PriceValue
			case mktmodel.PriceRuleReduce:
				actPrice = item.Price - act.PriceValue
			default:
				if ap.ActivityPrice > 0 {
					actPrice = ap.ActivityPrice
				} else {
					continue
				}
			}
			if actPrice < 0 {
				actPrice = 0.01
			}

			discount := (item.Price - actPrice) * float64(item.Qty)
			if discount > 0 {
				item.ActivityPrice = actPrice
				ctx.ActivityDiscount += discount
				ctx.AppliedRules = append(ctx.AppliedRules, marketing.AppliedRule{
					Type: "activity", Name: act.Name, Discount: discount,
				})
				if act.Exclusive {
					exclusive = true
				}
			}
			break // one activity per item
		}
	}

	return !exclusive, nil
}
