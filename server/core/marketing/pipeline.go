package marketing

import "sort"

// PriceCalculator is the interface every pricing step must implement.
type PriceCalculator interface {
	// Name returns a human-readable identifier.
	Name() string
	// Priority controls execution order (lower = earlier).
	Priority() int
	// Calculate mutates ctx. Return continueNext=false to stop the pipeline
	// (e.g. an exclusive activity that forbids further discounts).
	Calculate(ctx *PriceContext) (continueNext bool, err error)
}

var calculators []PriceCalculator

// Register adds a calculator to the global pipeline.
// Call this in each calculator package's init().
func Register(c PriceCalculator) {
	calculators = append(calculators, c)
	sort.Slice(calculators, func(i, j int) bool {
		return calculators[i].Priority() < calculators[j].Priority()
	})
}

// Calculate runs the full pricing pipeline.
func Calculate(ctx *PriceContext) error {
	// 1. Compute goods subtotal
	ctx.GoodsAmount = 0
	for _, item := range ctx.Items {
		ctx.GoodsAmount += item.Price * float64(item.Qty)
	}
	ctx.FinalAmount = ctx.GoodsAmount

	// 2. Execute calculators in priority order
	for _, calc := range calculators {
		cont, err := calc.Calculate(ctx)
		if err != nil {
			return err
		}
		// Recompute FinalAmount after each step
		ctx.FinalAmount = ctx.GoodsAmount -
			ctx.ActivityDiscount -
			ctx.FullReduceDiscount -
			ctx.CouponDiscount -
			ctx.PointsDiscount
		if !cont {
			break
		}
	}

	// 3. Floor at 0.01 (never free, unless GoodsAmount is 0)
	if ctx.FinalAmount < 0 {
		ctx.FinalAmount = 0
	}
	return nil
}

// All returns the registered calculators (for debugging / admin display).
func All() []PriceCalculator {
	out := make([]PriceCalculator, len(calculators))
	copy(out, calculators)
	return out
}
