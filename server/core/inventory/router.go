package inventory

import (
	"fmt"
	"strings"

	"github.com/ijry/lyshop/config"
)

func normalizeProviderName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "local":
		return "local"
	case "wms", "builtin_wms":
		return "builtin_wms"
	case "external", "external_wms":
		return "external_wms"
	default:
		return strings.ToLower(strings.TrimSpace(name))
	}
}

func CurrentProvider() (Provider, error) {
	name := normalizeProviderName(config.Global.Inventory.Provider)
	p := Find(name)
	if p == nil {
		return nil, fmt.Errorf("inventory provider %q not registered", name)
	}
	return p, nil
}

func ValidateConfig() error {
	name := normalizeProviderName(config.Global.Inventory.Provider)
	enabled := make(map[string]struct{}, len(config.Global.Plugins.Enabled))
	for _, pluginName := range config.Global.Plugins.Enabled {
		enabled[strings.ToLower(strings.TrimSpace(pluginName))] = struct{}{}
	}
	if name == "builtin_wms" {
		if _, ok := enabled["wms"]; !ok {
			return fmt.Errorf("inventory provider builtin_wms requires plugin \"wms\"")
		}
	}
	if name == "external_wms" {
		mode := strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode))
		if mode == "" {
			mode = "sync"
		}
		if mode != "sync" && mode != "async" {
			return fmt.Errorf("inventory.external_mode must be sync or async")
		}
		if strings.TrimSpace(config.Global.ExternalWMS.Endpoint) == "" {
			return fmt.Errorf("external_wms.endpoint is required when provider=external_wms")
		}
	}
	return nil
}

func IsAsyncExternalProvider() bool {
	return normalizeProviderName(config.Global.Inventory.Provider) == "external_wms" &&
		strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async"
}

func OrderInventoryStatusAfterReserve() string {
	if IsAsyncExternalProvider() {
		return InventoryStatusPending
	}
	return InventoryStatusReserved
}

func OrderInventoryStatusAfterConfirm() string {
	if IsAsyncExternalProvider() {
		return InventoryStatusPending
	}
	return InventoryStatusConfirmed
}

func OrderInventoryStatusAfterRelease() string {
	if IsAsyncExternalProvider() {
		return InventoryStatusPending
	}
	return InventoryStatusReleased
}
