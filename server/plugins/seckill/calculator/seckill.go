package calculator

import (
	"time"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/marketing"
	seckillmodel "github.com/ijry/lyshop/plugins/seckill/model"
)

// SeckillCalculator 秒杀价格计算器
type SeckillCalculator struct{}

func init() {
	marketing.Register(&SeckillCalculator{})
}

func (c *SeckillCalculator) Name() string  { return "seckill" }
func (c *SeckillCalculator) Priority() int { return 10 } // 最高优先级

func (c *SeckillCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
	// 查询所有有效的秒杀活动
	var activities []seckillmodel.SeckillActivity
	now := time.Now()
	if err := db.DB.Table("seckill_activities").
		Where("status = 1 AND start_at <= ? AND end_at >= ?", now, now).
		Find(&activities).Error; err != nil {
		return true, nil
	}

	if len(activities) == 0 {
		return true, nil
	}

	// 获取活动ID列表
	activityIDs := make([]uint64, len(activities))
	for i, act := range activities {
		activityIDs[i] = act.ID
	}

	// 查询秒杀商品
	var products []seckillmodel.SeckillProduct
	if err := db.DB.Table("seckill_products").
		Where("activity_id IN ?", activityIDs).
		Find(&products).Error; err != nil {
		return true, nil
	}

	// 构建商品映射：product_id:sku_id -> SeckillProduct
	productMap := make(map[string]*seckillmodel.SeckillProduct)
	for i := range products {
		key := makeKey(products[i].ProductID, products[i].SkuID)
		productMap[key] = &products[i]
	}

	// 应用秒杀价格
	hasApplied := false
	for i := range ctx.Items {
		item := &ctx.Items[i]

		// 尝试匹配 product_id + sku_id
		key := makeKey(item.ProductID, item.SkuID)
		if sp, ok := productMap[key]; ok {
			// 检查库存
			if sp.TotalStockLimit > 0 && sp.SoldQty >= sp.TotalStockLimit {
				continue
			}

			// 应用秒杀价格
			item.ActivityPrice = sp.SeckillPrice
			hasApplied = true
			continue
		}

		// 尝试匹配 product_id + sku_id=0（全部SKU）
		key = makeKey(item.ProductID, 0)
		if sp, ok := productMap[key]; ok {
			if sp.TotalStockLimit > 0 && sp.SoldQty >= sp.TotalStockLimit {
				continue
			}

			item.ActivityPrice = sp.SeckillPrice
			hasApplied = true
		}
	}

	if hasApplied {
		// 秒杀是排他性活动，停止后续计算器（除了分销）
		ctx.AppliedRules = append(ctx.AppliedRules, marketing.AppliedRule{
			Type: "seckill",
			Name: "秒杀活动",
		})
		return false, nil // 停止后续计算器
	}

	return true, nil
}

func makeKey(productID, skuID uint64) string {
	return string(rune(productID))+ ":" + string(rune(skuID))
}
