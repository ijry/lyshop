package logistics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type stubDriver struct{ name string }

func (d *stubDriver) Name() string { return d.name }

func (d *stubDriver) Query(context.Context, QueryReq) (*TrackResult, error) {
	return &TrackResult{Provider: d.name}, nil
}

func resetForTest() {
	mu.Lock()
	defer mu.Unlock()
	drivers = map[string]Driver{}
	defaultPrimary = "kuaidi100"
	defaultBackup = "kdniao"
}

func TestGetByName_NotFound(t *testing.T) {
	resetForTest()
	_, err := GetByName("missing")
	require.EqualError(t, err, `logistics driver "missing" not registered`)
}

func TestResolveByPinnedOrFallback_DefaultPrimary(t *testing.T) {
	resetForTest()
	Register(&stubDriver{name: "kuaidi100"})
	Register(&stubDriver{name: "kdniao"})

	driver, provider, err := ResolveByPinnedOrFallback("")
	require.NoError(t, err)
	require.Equal(t, "kuaidi100", provider)
	require.Equal(t, "kuaidi100", driver.Name())
}

func TestResolveByPinnedOrFallback_FallbackBackup(t *testing.T) {
	resetForTest()
	Register(&stubDriver{name: "kdniao"})

	driver, provider, err := ResolveByPinnedOrFallback("")
	require.NoError(t, err)
	require.Equal(t, "kdniao", provider)
	require.Equal(t, "kdniao", driver.Name())
}
