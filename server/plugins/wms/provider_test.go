package wms

import (
	"context"
	"testing"

	inventorycore "github.com/ijry/lyshop/core/inventory"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestBuiltinWMSProviderName(t *testing.T) {
	p := &builtinProvider{}
	require.Equal(t, "builtin_wms", p.Name())
}

func TestBuiltinWMSProviderReserveDelegatesToWMSService(t *testing.T) {
	p := &builtinProvider{}
	originalReserve := reserveStockTxFn
	originalPick := pickDefaultWarehouseIDTxFn
	called := false
	pickDefaultWarehouseIDTxFn = func(_ *gorm.DB) (uint64, error) { return 5, nil }
	reserveStockTxFn = func(_ *gorm.DB, in wmssvcReserveStockInput) error {
		called = true
		require.Equal(t, "order", in.BizType)
		require.Equal(t, "ORD-3", in.BizNo)
		require.Equal(t, uint64(5), in.WarehouseID)
		require.Len(t, in.Items, 1)
		return nil
	}
	t.Cleanup(func() {
		reserveStockTxFn = originalReserve
		pickDefaultWarehouseIDTxFn = originalPick
	})

	err := p.ReserveTx(nil, inventorycore.ReserveInput{
		BizType: "order",
		BizNo:   "ORD-3",
		Items:   []inventorycore.ReserveItem{{SkuID: 1, Qty: 2}},
	})
	require.NoError(t, err)
	require.True(t, called)
}

func TestBuiltinWMSProviderGetSellableStock(t *testing.T) {
	p := &builtinProvider{}
	original := listStocksBySkuIDsFn
	listStocksBySkuIDsFn = func(_ context.Context, skuIDs []uint64) ([]wmsmodel.InventoryStock, error) {
		require.Equal(t, []uint64{7}, skuIDs)
		return []wmsmodel.InventoryStock{{SkuID: 7, Qty: 10, ReservedQty: 3}}, nil
	}
	t.Cleanup(func() { listStocksBySkuIDsFn = original })

	list, err := p.GetSellableStock(context.Background(), []uint64{7})
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, 7, list[0].Sellable)
}
