package plugin

import "sync"

var (
	mu       sync.RWMutex
	registry []Plugin
)

// Register adds p to the global registry.
// Call this inside each plugin package's init() function.
func Register(p Plugin) {
	mu.Lock()
	defer mu.Unlock()
	registry = append(registry, p)
}

// All returns a snapshot of registered plugins.
func All() []Plugin {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Plugin, len(registry))
	copy(out, registry)
	return out
}

// Find returns the plugin with the given name, or nil.
func Find(name string) Plugin {
	mu.RLock()
	defer mu.RUnlock()
	for _, p := range registry {
		if p.Meta().Name == name {
			return p
		}
	}
	return nil
}

// EnabledMenus returns the merged menu tree for the enabled plugin list.
func EnabledMenus(enabled []string) []MenuItem {
	enabledSet := make(map[string]bool, len(enabled))
	for _, n := range enabled {
		enabledSet[n] = true
	}
	var menus []MenuItem
	for _, name := range enabled {
		p := Find(name)
		if p != nil {
			menus = append(menus, p.Meta().Menus...)
		}
	}
	return menus
}
