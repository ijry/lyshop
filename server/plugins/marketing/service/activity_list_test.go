package service

import (
	"context"
	"testing"
	"time"

	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
	"github.com/stretchr/testify/require"
)

func TestListFrontActivityProducts_HasActivityProductID(t *testing.T) {
	original := listFrontActivityProductsLoadRowsFn
	t.Cleanup(func() {
		listFrontActivityProductsLoadRowsFn = original
	})

	listFrontActivityProductsLoadRowsFn = func(ctx context.Context, normalizedType string, q ActivityProductListQuery, now time.Time) ([]frontActivityProductJoinedRow, error) {
		return []frontActivityProductJoinedRow{
			{
				ActivityProductID: 88,
				ActivityID:        11,
				ActivityType:      normalizedType,
				ActivityName:      "测试活动",
				ActivityStartAt:   &now,
				ActivityEndAt:     &now,
				ProductID:         1001,
				SkuID:             2001,
				Title:             "测试商品",
				Cover:             "https://example.com/a.jpg",
				Sales:             12,
				Stock:             99,
				ProductPrice:      100,
				ActivityPrice:     80,
				LimitPerOrder:     2,
				TotalStockLimit:   50,
				SoldQty:           10,
			},
		}, nil
	}

	list, total, err := ListFrontActivityProducts(context.Background(), mktmodel.ActivityTypeSeckill, ActivityProductListQuery{Page: 1, Size: 20})
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	require.Greater(t, list[0].ActivityProductID, uint64(0))
	require.Equal(t, uint64(88), list[0].ActivityProductID)
}
