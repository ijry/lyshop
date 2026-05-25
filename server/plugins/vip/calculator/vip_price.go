package calculator

import (
	"context"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/marketing"
	vipsvc "github.com/ijry/lyshop/plugins/vip/service"
)

type VipPriceCalculator struct{}

func init() { marketing.Register(&VipPriceCalculator{}) }

func (c *VipPriceCalculator) Name() string  { return "vip_price" }
func (c *VipPriceCalculator) Priority() int { return 15 }

func (c *VipPriceCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
	if db.DB == nil || !ctx.IsVIP || ctx.VIPLevelID == 0 {
		return true, nil
	}
	for i := range ctx.Items {
		item := &ctx.Items[i]
		if item.ActivityPrice > 0 {
			continue
		}
		priceRow, err := vipsvc.GetVipSkuPrice(context.Background(), ctx.VIPLevelID, item.ProductID, item.SkuID)
		if err != nil || priceRow == nil {
			continue
		}
		vipUnitPrice := priceRow.VipPrice
		if vipUnitPrice <= 0 && priceRow.VipDiscountRate > 0 {
			vipUnitPrice = item.Price * priceRow.VipDiscountRate
		}
		if vipUnitPrice <= 0 || vipUnitPrice >= item.Price {
			continue
		}
		discount := (item.Price - vipUnitPrice) * float64(item.Qty)
		ctx.VipDiscount += discount
		if ctx.ItemVIPDiscount == nil {
			ctx.ItemVIPDiscount = map[uint64]float64{}
		}
		ctx.ItemVIPDiscount[item.SkuID] += discount
		ctx.AppliedRules = append(ctx.AppliedRules, marketing.AppliedRule{
			Type: "vip", Name: "会员价", Discount: discount,
		})
	}
	return true, nil
}
