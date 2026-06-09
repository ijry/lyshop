package service

import (
	"context"
	"testing"

	inventorycore "github.com/ijry/lyshop/core/inventory"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spyPointsInventoryProvider struct {
	deductCalls  int
	restoreCalls int
}

func (s *spyPointsInventoryProvider) Name() string { return "spy" }
func (s *spyPointsInventoryProvider) ReserveTx(_ *gorm.DB, _ inventorycore.ReserveInput) error {
	return nil
}
func (s *spyPointsInventoryProvider) ConfirmTx(_ *gorm.DB, _, _ string) error {
	return nil
}
func (s *spyPointsInventoryProvider) ReleaseTx(_ *gorm.DB, _, _, _ string) error {
	return nil
}
func (s *spyPointsInventoryProvider) DeductTx(_ *gorm.DB, _ inventorycore.DeductInput) error {
	s.deductCalls++
	return nil
}
func (s *spyPointsInventoryProvider) RestoreTx(_ *gorm.DB, _ inventorycore.RestoreInput) error {
	s.restoreCalls++
	return nil
}
func (s *spyPointsInventoryProvider) SyncSkuTx(_ *gorm.DB, _ inventorycore.SyncSkuInput) error {
	return nil
}
func (s *spyPointsInventoryProvider) GetSellableStock(_ context.Context, _ []uint64) ([]inventorycore.SellableStock, error) {
	return nil, nil
}

func TestPointsMallUsesInventoryProviderHooks(t *testing.T) {
	original := getInventoryProviderForPointsMallFn
	spy := &spyPointsInventoryProvider{}
	getInventoryProviderForPointsMallFn = func() (inventorycore.Provider, error) { return spy, nil }
	t.Cleanup(func() { getInventoryProviderForPointsMallFn = original })

	err := deductPointsMallInventory(nil, 1001, 2)
	require.NoError(t, err)
	require.Equal(t, 1, spy.deductCalls)

	err = restorePointsMallInventory(nil, 1001, 2, "cancel")
	require.NoError(t, err)
	require.Equal(t, 1, spy.restoreCalls)
}
