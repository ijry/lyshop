package marketing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubCalc is a test calculator that subtracts a fixed amount.
type stubCalc struct {
	name     string
	prio     int
	discount float64
	stop     bool
}

func (s *stubCalc) Name() string  { return s.name }
func (s *stubCalc) Priority() int { return s.prio }
func (s *stubCalc) Calculate(ctx *PriceContext) (bool, error) {
	ctx.ActivityDiscount += s.discount
	ctx.AppliedRules = append(ctx.AppliedRules, AppliedRule{
		Type: "activity", Name: s.name, Discount: s.discount,
	})
	return !s.stop, nil
}

func TestPipeline(t *testing.T) {
	// Reset global state
	calculators = nil

	Register(&stubCalc{name: "seckill", prio: 10, discount: 100})
	Register(&stubCalc{name: "coupon", prio: 30, discount: 20})

	ctx := &PriceContext{
		Items: []OrderItem{
			{ProductID: 1, SkuID: 1, Price: 500, Qty: 1},
		},
	}
	err := Calculate(ctx)
	require.NoError(t, err)
	assert.Equal(t, 500.0, ctx.GoodsAmount)
	assert.Equal(t, 120.0, ctx.ActivityDiscount)
	assert.Equal(t, 380.0, ctx.FinalAmount)
	assert.Len(t, ctx.AppliedRules, 2)
}

func TestPipeline_ExclusiveStop(t *testing.T) {
	calculators = nil

	Register(&stubCalc{name: "exclusive-seckill", prio: 10, discount: 200, stop: true})
	Register(&stubCalc{name: "coupon", prio: 30, discount: 50})

	ctx := &PriceContext{
		Items: []OrderItem{
			{ProductID: 1, SkuID: 1, Price: 500, Qty: 1},
		},
	}
	err := Calculate(ctx)
	require.NoError(t, err)
	// Coupon should NOT have been executed because seckill stopped the pipeline
	assert.Equal(t, 200.0, ctx.ActivityDiscount)
	assert.Equal(t, 300.0, ctx.FinalAmount)
	assert.Len(t, ctx.AppliedRules, 1)
}

func TestPipeline_FloorAtZero(t *testing.T) {
	calculators = nil

	Register(&stubCalc{name: "mega-discount", prio: 10, discount: 9999})

	ctx := &PriceContext{
		Items: []OrderItem{
			{ProductID: 1, SkuID: 1, Price: 100, Qty: 1},
		},
	}
	Calculate(ctx)
	assert.Equal(t, 0.0, ctx.FinalAmount)
}
