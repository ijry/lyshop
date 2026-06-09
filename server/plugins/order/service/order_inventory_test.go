package service

import (
	"context"
	"testing"

	inventorycore "github.com/ijry/lyshop/core/inventory"
	ordermodel "github.com/ijry/lyshop/plugins/order/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spyInventoryProvider struct {
	reserveCalls int
	confirmCalls int
	releaseCalls int
	lastReserve  inventorycore.ReserveInput
}

func (s *spyInventoryProvider) Name() string { return "spy" }

func (s *spyInventoryProvider) ReserveTx(_ *gorm.DB, in inventorycore.ReserveInput) error {
	s.reserveCalls++
	s.lastReserve = in
	return nil
}

func (s *spyInventoryProvider) ConfirmTx(_ *gorm.DB, _, _ string) error {
	s.confirmCalls++
	return nil
}

func (s *spyInventoryProvider) ReleaseTx(_ *gorm.DB, _, _, _ string) error {
	s.releaseCalls++
	return nil
}

func (s *spyInventoryProvider) DeductTx(_ *gorm.DB, _ inventorycore.DeductInput) error {
	return nil
}

func (s *spyInventoryProvider) RestoreTx(_ *gorm.DB, _ inventorycore.RestoreInput) error {
	return nil
}

func (s *spyInventoryProvider) SyncSkuTx(_ *gorm.DB, _ inventorycore.SyncSkuInput) error {
	return nil
}

func (s *spyInventoryProvider) GetSellableStock(_ context.Context, _ []uint64) ([]inventorycore.SellableStock, error) {
	return nil, nil
}

func TestReserveOrderInventoryUsesInventoryProvider(t *testing.T) {
	original := getInventoryProviderFn
	spy := &spyInventoryProvider{}
	getInventoryProviderFn = func() (inventorycore.Provider, error) { return spy, nil }
	t.Cleanup(func() { getInventoryProviderFn = original })

	items := []ordermodel.OrderItem{{SkuID: 101, Qty: 2}}
	err := reserveOrderInventory(nil, "ORD-TEST", items)
	require.NoError(t, err)
	require.Equal(t, 1, spy.reserveCalls)
	require.Equal(t, "ORD-TEST", spy.lastReserve.BizNo)
	require.Len(t, spy.lastReserve.Items, 1)
	require.Equal(t, 2, spy.lastReserve.Items[0].Qty)
}
