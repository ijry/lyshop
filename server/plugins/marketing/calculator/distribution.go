package calculator

import (
	"encoding/json"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/marketing"
	usermodel "github.com/ijry/lyshop/model"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
)

// DistributionCalculator computes 2-level rebate commissions.
// It does NOT modify FinalAmount — commissions are settled after payment.
type DistributionCalculator struct{}

func init() { marketing.Register(&DistributionCalculator{}) }

func (c *DistributionCalculator) Name() string  { return "distribution" }
func (c *DistributionCalculator) Priority() int { return 50 }

func (c *DistributionCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
	if db.DB == nil {
		return true, nil
	}

	// Load distribution config from configs table
	var cfgKV usermodel.ConfigKV
	if err := db.DB.Where("plugin = ? AND `key` = ?", "marketing", "distribution_config").First(&cfgKV).Error; err != nil {
		return true, nil // distribution not configured
	}
	var cfg mktmodel.DistributionConfig
	if err := json.Unmarshal([]byte(cfgKV.Value), &cfg); err != nil {
		return true, nil
	}
	if cfg.Level1Rate <= 0 {
		return true, nil
	}

	// Find the buyer's distributor chain
	var dist mktmodel.Distributor
	if err := db.DB.Where("user_id = ?", ctx.UserID).First(&dist).Error; err != nil {
		return true, nil // buyer is not under any distributor
	}

	orderAmount := ctx.FinalAmount

	// Level 1: direct parent
	if dist.ParentID > 0 {
		l1Amount := orderAmount * cfg.Level1Rate
		if l1Amount > 0.01 {
			ctx.Commissions = append(ctx.Commissions, marketing.Commission{
				DistributorID: dist.ParentID, Level: 1, Amount: l1Amount,
			})
		}

		// Level 2: parent's parent
		if cfg.Level2Rate > 0 {
			var parentDist mktmodel.Distributor
			if err := db.DB.Where("user_id = ?", dist.ParentID).First(&parentDist).Error; err == nil && parentDist.ParentID > 0 {
				l2Amount := orderAmount * cfg.Level2Rate
				if l2Amount > 0.01 {
					ctx.Commissions = append(ctx.Commissions, marketing.Commission{
						DistributorID: parentDist.ParentID, Level: 2, Amount: l2Amount,
					})
				}
			}
		}
	}

	return true, nil
}
