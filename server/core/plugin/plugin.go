package plugin

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Meta holds data from a plugin's plugin.json.
type Meta struct {
	Name        string        `json:"name"`
	Title       string        `json:"title"`
	Version     string        `json:"version"`
	Description string        `json:"description"`
	Author      string        `json:"author"`
	Depends     []string      `json:"depends"`
	Menus       []MenuItem    `json:"menus"`
	Permissions []string      `json:"permissions"`
	ConfigItems []ConfigField `json:"config_items,omitempty"` // declarative config schema
}

// MenuItem describes one entry in the admin sidebar.
type MenuItem struct {
	Title      string     `json:"title"`
	Icon       string     `json:"icon"`
	Path       string     `json:"path"`
	Sort       int        `json:"sort"`
	Permission string     `json:"permission,omitempty"`
	Children   []MenuItem `json:"children,omitempty"`
}

// ConfigField describes one configuration field for admin UI rendering.
type ConfigField struct {
	Key         string `json:"key"`                    // config_kv key
	Label       string `json:"label"`                  // display label
	Type        string `json:"type"`                   // text|password|number|select|switch
	Placeholder string `json:"placeholder,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Options     []struct {                              // for type=select
		Label string `json:"label"`
		Value string `json:"value"`
	} `json:"options,omitempty"`
}

// Plugin is the interface every plugin must implement.
// Plugins self-register in their package's init() function.
type Plugin interface {
	// Meta returns plugin metadata (parsed from embedded plugin.json).
	Meta() Meta
	// RegisterRoutes registers front-end and admin API routes.
	RegisterRoutes(front, admin *gin.RouterGroup)
	// Migrate runs the plugin's DDL against db (idempotent).
	Migrate(db *gorm.DB) error
	// Install is called once after Migrate and RegisterRoutes.
	Install() error
	// Uninstall is called when the plugin is disabled.
	Uninstall() error
}
