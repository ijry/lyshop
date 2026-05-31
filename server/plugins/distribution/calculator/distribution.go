package calculator

import (
	"context"

	"github.com/ijry/lyshop/core/marketing"
	"github.com/ijry/lyshop/server/plugins/distribution/service"
)

func init() {
	marketing.Register(&DistributionCalculator{})
}

type DistributionCalculator struct{}

func (c *DistributionCalculator) Name() string {
	return "distribution"
}

func (c *DistributionCalculator) Priority() int {
	return 50 // 最低优先级，在所有折扣之后计算
}

func (c *DistributionCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
	// 分销不修改价格，只记录佣金信息
	// 实际佣金计算在订单创建时进行

	// 获取配置
	config, err := service.GetConfig(context.Background())
	if err != nil || !config.Enabled {
		return true, nil
	}

	// 标记订单可参与分销
	ctx.Metadata["distribution_enabled"] = true

	return true, nil
}
