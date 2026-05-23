package marketing

// OrderItem represents one line-item entering the pricing pipeline.
type OrderItem struct {
	ProductID uint64
	SkuID     uint64
	Title     string
	Price     float64 // original unit price
	Qty       int
	// Filled by activity calculators if hit
	ActivityPrice float64 // 0 means no activity price
}

// AppliedRule records a discount that was applied.
type AppliedRule struct {
	Type     string  `json:"type"`     // activity|coupon|full_reduce|points|distribution
	Name     string  `json:"name"`     // human-readable label
	Discount float64 `json:"discount"` // positive number
}

// Commission records a distribution rebate.
type Commission struct {
	DistributorID uint64  `json:"distributor_id"`
	Level         int     `json:"level"` // 1 or 2
	Amount        float64 `json:"amount"`
}

// PriceContext flows through the entire pipeline.
type PriceContext struct {
	// Input
	UserID     uint64
	Items      []OrderItem
	CouponIDs  []uint64 // coupons the user chose to use
	ActivityID uint64   // specific activity (0 = auto-detect)
	PointsUse  int      // points the user wants to spend

	// Computed (filled step by step)
	GoodsAmount        float64 // sum of original prices
	ActivityDiscount   float64
	FullReduceDiscount float64
	CouponDiscount     float64
	PointsDiscount     float64
	FinalAmount        float64

	AppliedRules []AppliedRule
	Commissions  []Commission // distribution rebates (does not affect FinalAmount)
}
