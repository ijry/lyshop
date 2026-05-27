package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateActivityProductSource(t *testing.T) {
	now := time.Now()
	original := getFrontActivityProductDetailFn
	t.Cleanup(func() { getFrontActivityProductDetailFn = original })

	base := &FrontActivityProductDetail{
		ActivityProductID: 9,
		ActivityID:        1,
		ActivityType:      "seckill",
		ActivityName:      "测试秒杀",
		ProductID:         100,
		SkuID:             200,
		ActivityStatus:    1,
		ActivityStartAt:   ptrTime(now.Add(-time.Hour)),
		ActivityEndAt:     ptrTime(now.Add(time.Hour)),
		TotalStockLimit:   10,
		SoldQty:           2,
	}

	t.Run("valid source", func(t *testing.T) {
		getFrontActivityProductDetailFn = func(ctx context.Context, activityProductID uint64) (*FrontActivityProductDetail, error) {
			copy := *base
			return &copy, nil
		}
		detail, err := ValidateActivityProductSource(context.Background(), 9, 200, 100)
		require.NoError(t, err)
		require.NotNil(t, detail)
		require.Equal(t, uint64(9), detail.ActivityProductID)
	})

	t.Run("sku mismatch", func(t *testing.T) {
		getFrontActivityProductDetailFn = func(ctx context.Context, activityProductID uint64) (*FrontActivityProductDetail, error) {
			copy := *base
			copy.SkuID = 201
			return &copy, nil
		}
		_, err := ValidateActivityProductSource(context.Background(), 9, 200, 100)
		require.ErrorContains(t, err, "SKU")
	})

	t.Run("product mismatch", func(t *testing.T) {
		getFrontActivityProductDetailFn = func(ctx context.Context, activityProductID uint64) (*FrontActivityProductDetail, error) {
			copy := *base
			copy.ProductID = 101
			return &copy, nil
		}
		_, err := ValidateActivityProductSource(context.Background(), 9, 200, 100)
		require.ErrorContains(t, err, "商品")
	})

	t.Run("inactive activity", func(t *testing.T) {
		getFrontActivityProductDetailFn = func(ctx context.Context, activityProductID uint64) (*FrontActivityProductDetail, error) {
			copy := *base
			copy.ActivityStatus = 0
			return &copy, nil
		}
		_, err := ValidateActivityProductSource(context.Background(), 9, 200, 100)
		require.ErrorContains(t, err, "失效")
	})

	t.Run("sold out", func(t *testing.T) {
		getFrontActivityProductDetailFn = func(ctx context.Context, activityProductID uint64) (*FrontActivityProductDetail, error) {
			copy := *base
			copy.TotalStockLimit = 2
			copy.SoldQty = 2
			return &copy, nil
		}
		_, err := ValidateActivityProductSource(context.Background(), 9, 200, 100)
		require.ErrorContains(t, err, "库存")
	})
}

func ptrTime(v time.Time) *time.Time {
	return &v
}
