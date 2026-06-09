package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	inventorycore "github.com/ijry/lyshop/core/inventory"
	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/model"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spySkuSyncProvider struct {
	syncCalls int
	lastInput inventorycore.SyncSkuInput
}

func (s *spySkuSyncProvider) Name() string { return "spy" }
func (s *spySkuSyncProvider) ReserveTx(_ *gorm.DB, _ inventorycore.ReserveInput) error {
	return nil
}
func (s *spySkuSyncProvider) ConfirmTx(_ *gorm.DB, _, _ string) error {
	return nil
}
func (s *spySkuSyncProvider) ReleaseTx(_ *gorm.DB, _, _, _ string) error {
	return nil
}
func (s *spySkuSyncProvider) DeductTx(_ *gorm.DB, _ inventorycore.DeductInput) error {
	return nil
}
func (s *spySkuSyncProvider) RestoreTx(_ *gorm.DB, _ inventorycore.RestoreInput) error {
	return nil
}
func (s *spySkuSyncProvider) SyncSkuTx(_ *gorm.DB, in inventorycore.SyncSkuInput) error {
	s.syncCalls++
	s.lastInput = in
	return nil
}
func (s *spySkuSyncProvider) GetSellableStock(_ context.Context, _ []uint64) ([]inventorycore.SellableStock, error) {
	return nil, nil
}

func TestSyncInventoryForNewSkusUsesInventoryProvider(t *testing.T) {
	original := getInventoryProviderForProductFn
	spy := &spySkuSyncProvider{}
	getInventoryProviderForProductFn = func() (inventorycore.Provider, error) { return spy, nil }
	t.Cleanup(func() { getInventoryProviderForProductFn = original })

	err := syncInventoryForNewSkus(nil, []productmodel.ProductSku{{Base: model.Base{ID: 9}, Stock: 12}})
	require.NoError(t, err)
	require.Equal(t, 1, spy.syncCalls)
	require.Equal(t, uint64(9), spy.lastInput.SkuID)
	require.Equal(t, 12, spy.lastInput.Stock)
}

type spySellableStockProvider struct {
	stocks []inventorycore.SellableStock
}

func (s *spySellableStockProvider) Name() string { return "external_wms" }
func (s *spySellableStockProvider) ReserveTx(_ *gorm.DB, _ inventorycore.ReserveInput) error {
	return nil
}
func (s *spySellableStockProvider) ConfirmTx(_ *gorm.DB, _, _ string) error {
	return nil
}
func (s *spySellableStockProvider) ReleaseTx(_ *gorm.DB, _, _, _ string) error {
	return nil
}
func (s *spySellableStockProvider) DeductTx(_ *gorm.DB, _ inventorycore.DeductInput) error {
	return nil
}
func (s *spySellableStockProvider) RestoreTx(_ *gorm.DB, _ inventorycore.RestoreInput) error {
	return nil
}
func (s *spySellableStockProvider) SyncSkuTx(_ *gorm.DB, _ inventorycore.SyncSkuInput) error {
	return nil
}
func (s *spySellableStockProvider) GetSellableStock(_ context.Context, _ []uint64) ([]inventorycore.SellableStock, error) {
	return s.stocks, nil
}

func TestGetProductUsesExternalSellableStock(t *testing.T) {
	testDB := setupProductInventoryDB(t)
	provider := &spySellableStockProvider{
		stocks: []inventorycore.SellableStock{{SkuID: 21, SellableStock: 7}},
	}
	original := getInventoryProviderForProductFn
	getInventoryProviderForProductFn = func() (inventorycore.Provider, error) { return provider, nil }
	t.Cleanup(func() { getInventoryProviderForProductFn = original })

	product := productmodel.Product{
		Base:   model.Base{ID: 1},
		Title:  "demo",
		Price:  10,
		Stock:  99,
		Status: 1,
		Detail: json.RawMessage(`{"version":1,"blocks":[]}`),
	}
	require.NoError(t, testDB.Create(&product).Error)
	sku := productmodel.ProductSku{
		Base:      model.Base{ID: 21},
		ProductID: 1,
		Price:     10,
		Stock:     99,
		SkuKey:    "red:l",
		Status:    productmodel.ProductSkuStatusActive,
	}
	require.NoError(t, testDB.Create(&sku).Error)

	detail, err := GetProduct(context.Background(), 1, 0)
	require.NoError(t, err)
	require.Len(t, detail.SKUs, 1)
	require.Equal(t, 7, detail.SKUs[0].Stock)
}

func setupProductInventoryDB(t *testing.T) *gorm.DB {
	t.Helper()

	gdb, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:product_inventory_%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	require.NoError(t, err)

	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })

	require.NoError(t, gdb.AutoMigrate(&productmodel.Product{}, &productmodel.ProductSku{}, &productmodel.ProductImage{}))
	return gdb
}
