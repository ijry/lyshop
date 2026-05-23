package calculator

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/marketing"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
)

// FullReduceCalculator handles "满减" rules defined as full_save activities.
type FullReduceCalculator struct{}

func init() { marketing.Register(&FullReduceCalculator{}) }

func (c *FullReduceCalculator) Name() string  { return "full_reduce" }
func (c *FullReduceCalculator) Priority() int { return 20 }

type fullReduceRule struct {
	Min    float64 `json:"min"`
	Reduce float64 `json:"reduce"`
}

func (c *FullReduceCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
	if db.DB == nil {
		return true, nil
	}

	now := time.Now()
	var activities []mktmodel.Activity
	db.DB.Where("type = ? AND status = 1 AND start_at <= ? AND end_at >= ?",
		mktmodel.ActivityTypeFullSave, now, now).Find(&activities)

	// Current subtotal after activity discounts
	subtotal := ctx.GoodsAmount - ctx.ActivityDiscount

	for _, act := range activities {
		var rules []fullReduceRule
		json.Unmarshal(act.Config, &rules)

		// Find the best matching rule (highest min that subtotal meets)
		var best *fullReduceRule
		for i := range rules {
			if subtotal >= rules[i].Min {
				if best == nil || rules[i].Min > best.Min {
					best = &rules[i]
				}
			}
		}
		if best != nil && best.Reduce > 0 {
			ctx.FullReduceDiscount += best.Reduce
			ctx.AppliedRules = append(ctx.AppliedRules, marketing.AppliedRule{
				Type: "full_reduce", Name: act.Name, Discount: best.Reduce,
			})
		}
	}

	return true, nil
}
