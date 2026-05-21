package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	content := `
server:
  port: 8080
  mode: debug
database:
  dsn: "root:pw@tcp(localhost:3306)/lyshop"
  max_open: 10
  max_idle: 5
redis:
  addr: "localhost:6379"
jwt:
  secret: "test-secret"
  expire_hours: 168
plugins:
  enabled:
    - product
    - order
`
	f, err := os.CreateTemp("", "lyshop-cfg-*.yaml")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	f.WriteString(content)
	f.Close()

	err = Load(f.Name())
	require.NoError(t, err)
	assert.Equal(t, 8080, Global.Server.Port)
	assert.Equal(t, "debug", Global.Server.Mode)
	assert.Equal(t, []string{"product", "order"}, Global.Plugins.Enabled)
}
