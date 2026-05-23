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

// EnabledMenus returns the merged menu tree for the enabled plugin list,
// filtered by the admin's permissions. Pass nil perms to return all menus.
func EnabledMenus(enabled []string, perms []string) []MenuItem {
	hasPerm := func(p string) bool {
		if perms == nil {
			return true
		}
		for _, pp := range perms {
			if pp == "*" || pp == p {
				return true
			}
		}
		return p == "" // no permission required = always visible
	}

	var menus []MenuItem
	for _, name := range enabled {
		p := Find(name)
		if p == nil {
			continue
		}
		for _, menu := range p.Meta().Menus {
			filtered := filterMenu(menu, hasPerm)
			if filtered != nil {
				menus = append(menus, *filtered)
			}
		}
	}
	return menus
}

func filterMenu(m MenuItem, hasPerm func(string) bool) *MenuItem {
	// If this item has a permission requirement and user lacks it, skip
	if m.Permission != "" && !hasPerm(m.Permission) {
		return nil
	}
	// Filter children
	var children []MenuItem
	for _, child := range m.Children {
		fc := filterMenu(child, hasPerm)
		if fc != nil {
			children = append(children, *fc)
		}
	}
	// If all children were filtered out and this is a parent-only item, skip
	if len(m.Children) > 0 && len(children) == 0 {
		return nil
	}
	result := m
	result.Children = children
	return &result
}

// AllPermissions returns a deduplicated list of all permissions declared by enabled plugins.
func AllPermissions(enabled []string) []string {
	seen := map[string]bool{}
	var perms []string
	for _, name := range enabled {
		p := Find(name)
		if p == nil {
			continue
		}
		for _, perm := range p.Meta().Permissions {
			if !seen[perm] {
				seen[perm] = true
				perms = append(perms, perm)
			}
		}
	}
	return perms
}
