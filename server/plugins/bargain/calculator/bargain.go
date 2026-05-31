package calculator

import (
	"context"
	"fmt"

	"github.com/ijry/lyshop/core/marketing"
	"github.com/ijry/lyshop/plugins/bargain/service"
)

func init() {
	marketing.Register(&BargainCalculator{})
}

type BargainCalculator struct{}

func (c *BargainCalculator) Name() string {
	return "bargain"
}

func (c *BargainCalculator) Priority() int {
	return 10 // 与秒杀、拼团同优先级，排他性活动
}

func (c *BargainCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
	// 获取当前有效的砍价活动
	activities, err := service.GetActiveActivities(context.Background())
	if err != nil {
		return true, err
	}

	if len(activities) == 0 {
		return true, nil
	}

	// 遍历购物车商品
	for i := range ctx.Items {
		item := &ctx.Items[i]
		key := makeKey(item.ProductID, item.SkuID)

		// 查找匹配的砍价商品
		for _, activity := range activities {
			products, _, err := service.ListProducts(context.Background(), activity.ID, 1, 1000)
			if err != nil {
				continue
			}

			for _, product := range products {
				productKey := makeKey(product.ProductID, product.SkuID)
				if productKey == key {
					// 应用砍价底价
					discount := item.Price - product.FloorPrice
					if discount < 0 {
						discount = 0
					}
					item.ActivityPrice = product.FloorPrice
					ctx.AppliedRules = append(ctx.AppliedRules, marketing.AppliedRule{
						Type:     "bargain",
						Name:     activity.Name,
						Discount: discount,
					})
					// 砍价是排他性活动，停止后续计算器
					return false, nil
				}
			}
		}
	}

	return true, nil
}

func makeKey(productID, skuID uint64) string {
	if skuID > 0 {
		return fmt.Sprintf("%d:%d", productID, skuID)
	}
	return fmt.Sprintf("%d:0", productID)
}
