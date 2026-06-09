package service

import (
	"context"
	"testing"

	"github.com/ijry/lyshop/config"
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

func TestOrderInventoryStatusTransitionsByProviderMode(t *testing.T) {
	original := config.Global
	t.Cleanup(func() { config.Global = original })

	cases := []struct {
		name              string
		provider          string
		externalMode      string
		wantAfterReserve  string
		wantAfterConfirm  string
		wantAfterRelease  string
	}{
		{
			name:             "local",
			provider:         "local",
			wantAfterReserve: inventorycore.InventoryStatusReserved,
			wantAfterConfirm: inventorycore.InventoryStatusConfirmed,
			wantAfterRelease: inventorycore.InventoryStatusReleased,
		},
		{
			name:             "builtin_wms",
			provider:         "builtin_wms",
			wantAfterReserve: inventorycore.InventoryStatusReserved,
			wantAfterConfirm: inventorycore.InventoryStatusConfirmed,
			wantAfterRelease: inventorycore.InventoryStatusReleased,
		},
		{
			name:             "external_sync",
			provider:         "external_wms",
			externalMode:     "sync",
			wantAfterReserve: inventorycore.InventoryStatusReserved,
			wantAfterConfirm: inventorycore.InventoryStatusConfirmed,
			wantAfterRelease: inventorycore.InventoryStatusReleased,
		},
		{
			name:             "external_async",
			provider:         "external_wms",
			externalMode:     "async",
			wantAfterReserve: inventorycore.InventoryStatusPending,
			wantAfterConfirm: inventorycore.InventoryStatusPending,
			wantAfterRelease: inventorycore.InventoryStatusPending,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config.Global.Inventory.Provider = tc.provider
			config.Global.Inventory.ExternalMode = tc.externalMode

			require.Equal(t, tc.wantAfterReserve, inventorycore.OrderInventoryStatusAfterReserve())
			require.Equal(t, tc.wantAfterConfirm, inventorycore.OrderInventoryStatusAfterConfirm())
			require.Equal(t, tc.wantAfterRelease, inventorycore.OrderInventoryStatusAfterRelease())
		})
	}
}
