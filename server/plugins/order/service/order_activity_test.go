package service

import (
	"context"
	"errors"
	"testing"
	"time"

	marketingsvc "github.com/ijry/lyshop/plugins/marketing/service"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"github.com/stretchr/testify/require"
)

func TestBuildOrderItemsFromCart_UsesItemsPriorityAndActivityPrice(t *testing.T) {
	original := validateActivityProductSourceFn
	t.Cleanup(func() { validateActivityProductSourceFn = original })

	now := time.Now()
	validateActivityProductSourceFn = func(ctx context.Context, activityProductID, skuID, productID uint64) (*marketingsvc.FrontActivityProductDetail, error) {
		if activityProductID != 9001 || skuID != 101 || productID != 11 {
			return nil, errors.New("unexpected source")
		}
		return &marketingsvc.FrontActivityProductDetail{
			ActivityProductID: activityProductID,
			ActivityID:        77,
			ActivityType:      "seckill",
			ActivityName:      "测试秒杀",
			ActivityStatus:    1,
			ActivityStartAt:   &now,
			ActivityEndAt:     &now,
			ProductID:         productID,
			SkuID:             skuID,
			Price:             79.9,
		}, nil
	}

	cartItems := []CartItem{
		{
			SkuID:             101,
			ActivityProductID: 9001,
			Qty:               2,
			Product:           &productmodel.Product{Base: productmodel.Product{}.Base, Title: "A", Cover: "a.jpg"},
			Sku:               &productmodel.ProductSku{Base: productmodel.ProductSku{}.Base, ProductID: 11, Price: 99.9},
		},
		{
			SkuID:             102,
			ActivityProductID: 0,
			Qty:               1,
			Product:           &productmodel.Product{Base: productmodel.Product{}.Base, Title: "B", Cover: "b.jpg"},
			Sku:               &productmodel.ProductSku{Base: productmodel.ProductSku{}.Base, ProductID: 12, Price: 49.9},
		},
	}
	cartItems[0].Product.ID = 11
	cartItems[0].Sku.ID = 101
	cartItems[1].Product.ID = 12
	cartItems[1].Sku.ID = 102

	req := CreateOrderReq{
		Items: []CreateOrderItemReq{
			{SkuID: 101, ActivityProductID: 9001},
		},
		SkuIDs: []uint64{102},
	}

	items, err := buildOrderItemsFromCart(context.Background(), cartItems, req)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, uint64(101), items[0].SkuID)
	require.Equal(t, uint64(9001), items[0].ActivityProductID)
	require.Equal(t, uint64(77), items[0].ActivityID)
	require.Equal(t, "seckill", items[0].ActivityType)
	require.Equal(t, "测试秒杀", items[0].ActivityTitle)
	require.Equal(t, 79.9, items[0].Price)
}

func TestBuildOrderItemsFromCart_ActivitySourceError(t *testing.T) {
	original := validateActivityProductSourceFn
	t.Cleanup(func() { validateActivityProductSourceFn = original })

	validateActivityProductSourceFn = func(ctx context.Context, activityProductID, skuID, productID uint64) (*marketingsvc.FrontActivityProductDetail, error) {
		return nil, errors.New("活动已结束")
	}

	cartItems := []CartItem{
		{
			SkuID:             101,
			ActivityProductID: 9001,
			Qty:               1,
			Product:           &productmodel.Product{Base: productmodel.Product{}.Base, Title: "A", Cover: "a.jpg"},
			Sku:               &productmodel.ProductSku{Base: productmodel.ProductSku{}.Base, ProductID: 11, Price: 99.9},
		},
	}
	cartItems[0].Product.ID = 11
	cartItems[0].Sku.ID = 101

	req := CreateOrderReq{
		Items: []CreateOrderItemReq{
			{SkuID: 101, ActivityProductID: 9001},
		},
	}

	_, err := buildOrderItemsFromCart(context.Background(), cartItems, req)
	require.ErrorContains(t, err, "活动已结束")
}
