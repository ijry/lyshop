package inventory

import (
	"context"
	"testing"

	"github.com/ijry/lyshop/config"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type stubProvider struct{ name string }

func (s *stubProvider) Name() string                               { return s.name }
func (s *stubProvider) ReserveTx(_ *gorm.DB, _ ReserveInput) error { return nil }
func (s *stubProvider) ConfirmTx(_ *gorm.DB, _, _ string) error    { return nil }
func (s *stubProvider) ReleaseTx(_ *gorm.DB, _, _, _ string) error { return nil }
func (s *stubProvider) DeductTx(_ *gorm.DB, _ DeductInput) error   { return nil }
func (s *stubProvider) RestoreTx(_ *gorm.DB, _ RestoreInput) error { return nil }
func (s *stubProvider) SyncSkuTx(_ *gorm.DB, _ SyncSkuInput) error { return nil }
func (s *stubProvider) GetSellableStock(_ context.Context, _ []uint64) ([]SellableStock, error) {
	return nil, nil
}

func TestCurrentProviderUsesInventoryConfig(t *testing.T) {
	ResetRegistryForTest()
	Register(&stubProvider{name: "local"})
	Register(&stubProvider{name: "builtin_wms"})

	original := config.Global
	t.Cleanup(func() { config.Global = original })
	config.Global.Inventory.Provider = "builtin_wms"

	p, err := CurrentProvider()
	require.NoError(t, err)
	require.Equal(t, "builtin_wms", p.Name())
}

func TestValidateConfigRejectsBuiltinWMSWithoutPlugin(t *testing.T) {
	original := config.Global
	t.Cleanup(func() { config.Global = original })
	config.Global.Inventory.Provider = "builtin_wms"
	config.Global.Plugins.Enabled = []string{"product", "order"}

	err := ValidateConfig()
	require.ErrorContains(t, err, "wms")
}

func TestValidateConfigRejectsExternalWMSWithoutEndpoint(t *testing.T) {
	original := config.Global
	t.Cleanup(func() { config.Global = original })
	config.Global.Inventory.Provider = "external_wms"
	config.Global.Inventory.ExternalMode = "sync"
	config.Global.ExternalWMS.Endpoint = ""

	err := ValidateConfig()
	require.ErrorContains(t, err, "external_wms.endpoint")
}
