package service

import (
	"testing"
	"time"

	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	"github.com/stretchr/testify/require"
)

func TestReserveConfirmReleaseFlow(t *testing.T) {
	testDB := setupWmsTestDB(t)
	ctx := t.Context()

	warehouse := wmsmodel.Warehouse{Code: "WH-RSV", Name: "预占仓", Status: wmsmodel.WarehouseStatusEnabled}
	require.NoError(t, testDB.Create(&warehouse).Error)
	require.NoError(t, testDB.Create(&wmsmodel.InventoryStock{
		WarehouseID: warehouse.ID,
		SkuID:       1001,
		Qty:         20,
		ReservedQty: 0,
	}).Error)

	expireAt := time.Now().Add(15 * time.Minute)
	err := ReserveStock(ctx, ReserveStockInput{
		BizType:     "order",
		BizNo:       "ORD1001",
		WarehouseID: warehouse.ID,
		Items: []ReservationItemInput{
			{SkuID: 1001, Qty: 5},
		},
		ExpiredAt: &expireAt,
	})
	require.NoError(t, err)

	var stock wmsmodel.InventoryStock
	require.NoError(t, testDB.Where("warehouse_id = ? AND sku_id = ?", warehouse.ID, 1001).First(&stock).Error)
	require.Equal(t, 20, stock.Qty)
	require.Equal(t, 5, stock.ReservedQty)

	require.NoError(t, ConfirmReservation(ctx, "order", "ORD1001"))
	require.NoError(t, testDB.Where("warehouse_id = ? AND sku_id = ?", warehouse.ID, 1001).First(&stock).Error)
	require.Equal(t, 15, stock.Qty)
	require.Equal(t, 0, stock.ReservedQty)

	var reservations []wmsmodel.InventoryReservation
	require.NoError(t, testDB.Where("biz_type = ? AND biz_no = ?", "order", "ORD1001").Find(&reservations).Error)
	require.Len(t, reservations, 1)
	require.Equal(t, wmsmodel.ReservationStatusConfirmed, reservations[0].Status)
}

func TestReleaseReservationIdempotent(t *testing.T) {
	testDB := setupWmsTestDB(t)
	ctx := t.Context()

	warehouse := wmsmodel.Warehouse{Code: "WH-REL", Name: "释放仓", Status: wmsmodel.WarehouseStatusEnabled}
	require.NoError(t, testDB.Create(&warehouse).Error)
	require.NoError(t, testDB.Create(&wmsmodel.InventoryStock{
		WarehouseID: warehouse.ID,
		SkuID:       2001,
		Qty:         10,
		ReservedQty: 0,
	}).Error)

	require.NoError(t, ReserveStock(ctx, ReserveStockInput{
		BizType:     "order",
		BizNo:       "ORD2001",
		WarehouseID: warehouse.ID,
		Items: []ReservationItemInput{
			{SkuID: 2001, Qty: 3},
		},
	}))

	require.NoError(t, ReleaseReservation(ctx, "order", "ORD2001", "cancel"))
	require.NoError(t, ReleaseReservation(ctx, "order", "ORD2001", "cancel"))

	var stock wmsmodel.InventoryStock
	require.NoError(t, testDB.Where("warehouse_id = ? AND sku_id = ?", warehouse.ID, 2001).First(&stock).Error)
	require.Equal(t, 10, stock.Qty)
	require.Equal(t, 0, stock.ReservedQty)

	var reservations []wmsmodel.InventoryReservation
	require.NoError(t, testDB.Where("biz_type = ? AND biz_no = ?", "order", "ORD2001").Find(&reservations).Error)
	require.Len(t, reservations, 1)
	require.Equal(t, wmsmodel.ReservationStatusReleased, reservations[0].Status)
}
