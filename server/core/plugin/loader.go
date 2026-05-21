package plugin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Load migrates, registers routes, and installs each enabled plugin in order.
// It validates that all declared dependencies are also enabled.
func Load(enabled []string, db *gorm.DB, front, admin *gin.RouterGroup) error {
	enabledSet := make(map[string]bool, len(enabled))
	for _, n := range enabled {
		enabledSet[n] = true
	}

	// Validate existence and dependencies first
	for _, name := range enabled {
		p := Find(name)
		if p == nil {
			return fmt.Errorf(
				"plugin %q is in plugins.enabled but not registered; "+
					"add its blank import to main.go", name)
		}
		for _, dep := range p.Meta().Depends {
			if !enabledSet[dep] {
				return fmt.Errorf(
					"plugin %q requires plugin %q, "+
						"but %q is not in plugins.enabled", name, dep, dep)
			}
		}
	}

	// Load in config order (dependency order is the caller's responsibility)
	for _, name := range enabled {
		p := Find(name)
		if err := p.Migrate(db); err != nil {
			return fmt.Errorf("plugin %q Migrate: %w", name, err)
		}
		p.RegisterRoutes(front, admin)
		if err := p.Install(); err != nil {
			return fmt.Errorf("plugin %q Install: %w", name, err)
		}
	}
	return nil
}
