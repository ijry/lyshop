package service

import (
	"context"
	"testing"

	inventorycore "github.com/ijry/lyshop/core/inventory"
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
