package calculator

import (
	"github.com/ijry/lyshop/server/core/db"
	"github.com/ijry/lyshop/server/core/marketing"
	usermodel "github.com/ijry/lyshop/server/model"
)

// PointsCalculator deducts points from the final price.
// 100 points = ¥1.00 (configurable).
type PointsCalculator struct{}

func init() { marketing.Register(&PointsCalculator{}) }

func (c *PointsCalculator) Name() string  { return "points" }
func (c *PointsCalculator) Priority() int { return 40 }

const pointsToYuan = 100.0 // 100 points = ¥1

func (c *PointsCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
	if ctx.PointsUse <= 0 || db.DB == nil {
		return true, nil
	}

	// Check user balance
	var user usermodel.User
	if err := db.DB.First(&user, ctx.UserID).Error; err != nil {
		return true, nil
	}
	use := ctx.PointsUse
	if use > user.Points {
		use = user.Points
	}
	if use <= 0 {
		return true, nil
	}

	discount := float64(use) / pointsToYuan
	// Don't let points discount exceed remaining amount
	remaining := ctx.GoodsAmount - ctx.ActivityDiscount - ctx.FullReduceDiscount - ctx.CouponDiscount
	if discount > remaining {
		discount = remaining
		use = int(discount * pointsToYuan)
	}

	ctx.PointsDiscount += discount
	ctx.PointsUse = use // update to actual used
	ctx.AppliedRules = append(ctx.AppliedRules, marketing.AppliedRule{
		Type: "points", Name: "积分抵扣", Discount: discount,
	})

	return true, nil
}
