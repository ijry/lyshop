package plugin

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// stub is a minimal Plugin for testing.
type stub struct {
	name string
	deps []string
}

func (s *stub) Meta() Meta                             { return Meta{Name: s.name, Depends: s.deps} }
func (s *stub) RegisterRoutes(_, _ *gin.RouterGroup)   {}
func (s *stub) Migrate(_ *gorm.DB) error               { return nil }
func (s *stub) Install() error                         { return nil }
func (s *stub) Uninstall() error                       { return nil }

func resetRegistry() {
	mu.Lock()
	registry = nil
	mu.Unlock()
}

func TestRegisterAndFind(t *testing.T) {
	resetRegistry()
	Register(&stub{name: "product"})
	Register(&stub{name: "order", deps: []string{"product"}})

	assert.NotNil(t, Find("product"))
	assert.NotNil(t, Find("order"))
	assert.Nil(t, Find("nonexistent"))
}

func TestLoad_MissingDependency(t *testing.T) {
	resetRegistry()
	Register(&stub{name: "order", deps: []string{"product"}})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	err := Load([]string{"order"}, nil, r.Group("/api/v1"), r.Group("/admin/api"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "product")
}

func TestLoad_OK(t *testing.T) {
	resetRegistry()
	Register(&stub{name: "product"})
	Register(&stub{name: "order", deps: []string{"product"}})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	err := Load([]string{"product", "order"}, &gorm.DB{}, r.Group("/api/v1"), r.Group("/admin/api"))
	require.NoError(t, err)
}
