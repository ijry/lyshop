package calculator

import (
	"context"
	"fmt"

	"github.com/ijry/lyshop/core/marketing"
	"github.com/ijry/lyshop/plugins/group_buy/service"
)

func init() {
	marketing.Register(&GroupBuyCalculator{})
}

type GroupBuyCalculator struct{}

func (c *GroupBuyCalculator) Name() string {
	return "group_buy"
}

func (c *GroupBuyCalculator) Priority() int {
	return 10 // 与秒杀同优先级，排他性活动
}

func (c *GroupBuyCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
	// 获取当前有效的拼团活动
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

		// 查找匹配的拼团商品
		for _, activity := range activities {
			products, _, err := service.ListProducts(context.Background(), activity.ID, 1, 1000)
			if err != nil {
				continue
			}

			for _, product := range products {
				productKey := makeKey(product.ProductID, product.SkuID)
				if productKey == key {
					// 应用拼团价格
					discount := item.Price - product.GroupPrice
					if discount < 0 {
						discount = 0
					}
					item.ActivityPrice = product.GroupPrice
					ctx.AppliedRules = append(ctx.AppliedRules, marketing.AppliedRule{
						Type:     "group_buy",
						Name:     activity.Name,
						Discount: discount,
					})
					// 拼团是排他性活动，停止后续计算器
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
