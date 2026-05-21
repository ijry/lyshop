# lyshop 2.0 Phase 1: Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** GitHub repo, Go backend foundation (config, DB, Redis, plugin system, 5 driver interfaces, JWT auth APIs), Docker Compose, Vue3 admin skeleton, uni-app skeleton — everything needed before any plugin can be built.

**Architecture:** Single Go binary with a plugin registry. Core wires Gin router, GORM, Redis, and JWT middleware. Plugins self-register on startup via `init()`. Admin (Vue3+Vite) and frontend (uni-app) deploy as static files behind Nginx.

**Tech Stack:** Go 1.22 · Gin v1.10 · GORM v2 · MySQL 8.0 · go-redis/v9 · golang-jwt/v5 · Viper · Vue3 + Vite 5 + TailwindCSS 3 + shadcn-vue · uni-app Vue3 + uview-plus 3.x · Docker Compose v2 · gh CLI

**Spec:** `docs/superpowers/specs/2026-05-21-lyshop2-design.md`

---

## File Map

```
lyshop/
├── .gitignore
├── .github/workflows/ci.yml
├── docker-compose.yml
├── nginx.conf
├── config.example.yaml
│
├── server/                            # Go backend
│   ├── go.mod
│   ├── main.go
│   ├── config/
│   │   └── config.go                  # Config struct + Viper loader
│   ├── core/
│   │   ├── app/app.go                 # App lifecycle: Init + Run
│   │   ├── db/db.go                   # GORM MySQL connection
│   │   ├── cache/cache.go             # Redis client
│   │   ├── response/response.go       # Unified JSON response helpers
│   │   ├── middleware/
│   │   │   ├── auth.go                # JWT Claims, RequireAuth, RequireAdmin
│   │   │   ├── cors.go                # CORS middleware
│   │   │   └── logger.go              # Request logger middleware
│   │   ├── plugin/
│   │   │   ├── plugin.go              # Plugin interface + Meta/MenuItem structs
│   │   │   ├── registry.go            # Global plugin registry (Register/Find/All)
│   │   │   └── loader.go              # Load(): dep-check + migrate + routes + install
│   │   └── driver/
│   │       ├── payment/payment.go     # payment.Driver interface + registry
│   │       ├── sms/sms.go             # sms.Driver interface + registry
│   │       ├── oauth/oauth.go         # oauth.Driver interface + registry
│   │       ├── storage/storage.go     # storage.Driver interface + registry
│   │       └── ai/ai.go               # ai.Driver interface + registry
│   ├── model/
│   │   ├── base.go                    # Base struct (ID, CreatedAt, UpdatedAt)
│   │   ├── user.go
│   │   ├── admin.go
│   │   ├── role.go
│   │   └── config_kv.go
│   ├── service/auth/
│   │   ├── user_auth.go               # SMS code flow + user token
│   │   └── admin_auth.go              # Password login + admin token
│   ├── api/auth/
│   │   ├── front.go                   # POST /api/v1/auth/sms/send, /login
│   │   └── admin.go                   # POST /admin/api/auth/login
│   └── plugins/                       # Empty — populated in Phase 2+
│       └── .gitkeep
│
├── admin/                             # Vue3 + TailwindCSS admin
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   ├── index.html
│   ├── Dockerfile
│   └── src/
│       ├── main.ts
│       ├── App.vue
│       ├── router/index.ts
│       ├── stores/auth.ts
│       ├── api/
│       │   ├── request.ts             # Axios instance + interceptors
│       │   └── auth.ts                # login / logout API calls
│       ├── layouts/AdminLayout.vue    # Sidebar + topbar shell
│       └── views/
│           ├── Login.vue
│           └── Dashboard.vue
│
└── app/                               # uni-app + uview-plus 3.x
    ├── package.json
    ├── manifest.json
    ├── pages.json
    ├── uni.scss
    ├── App.vue
    ├── main.ts
    ├── pages/
    │   ├── index/index.vue
    │   └── login/index.vue
    └── utils/request.ts
```

---

## Task 1: GitHub Repository Setup

**Files:** none (commands only)

- [ ] **Step 1: Verify gh CLI auth**

```bash
gh auth status
```
Expected: shows authenticated account. If not: `gh auth login`

- [ ] **Step 2: Create repo and push**

```bash
cd /d/Repos/xyito/open/lyshop
gh repo create lyshop --public \
  --description "lyshop 2.0 - 开源插件化商城 | Go + Vue3 + uni-app" \
  --source=. --push
```

Expected output: `✓ Created repository <user>/lyshop on GitHub` and `✓ Pushed commits to github.com/<user>/lyshop`

- [ ] **Step 3: Verify**

```bash
gh repo view lyshop --web
```

---

## Task 2: Root Config Files

**Files:**
- Create: `.gitignore`
- Create: `config.example.yaml`

- [ ] **Step 1: Create .gitignore**

```
# Go
server/bin/
server/*.exe
server/lyshop-server

# Config (contains secrets)
config.yaml

# Data volumes
data/

# Node
admin/node_modules/
admin/dist/
app/node_modules/
app/dist/
app/unpackage/

# IDE
.idea/
.vscode/
*.swp
```

- [ ] **Step 2: Create config.example.yaml**

```yaml
server:
  port: 8080
  mode: debug   # debug | release

database:
  dsn: "root:password@tcp(mysql:3306)/lyshop?charset=utf8mb4&parseTime=True&loc=Local"
  max_open: 100
  max_idle: 10

redis:
  addr: "redis:6379"
  password: ""
  db: 0

jwt:
  secret: "change-this-to-a-random-string"
  expire_hours: 168   # 7 days

plugins:
  enabled:
    - storage_local
```

- [ ] **Step 3: Commit**

```bash
git add .gitignore config.example.yaml
git commit -m "chore: add gitignore and config template"
```

---

## Task 3: Go Module + Dependencies

**Files:**
- Create: `server/go.mod`

- [ ] **Step 1: Initialize Go module**

```bash
cd server
go mod init github.com/jry/lyshop
```

- [ ] **Step 2: Install dependencies**

```bash
go get github.com/gin-gonic/gin@v1.10.0
go get gorm.io/gorm@v1.25.10
go get gorm.io/driver/mysql@v1.5.7
go get github.com/redis/go-redis/v9@v9.5.3
go get github.com/golang-jwt/jwt/v5@v5.2.1
go get github.com/spf13/viper@v1.19.0
go get golang.org/x/crypto@v0.22.0
go get github.com/google/uuid@v1.6.0
go get github.com/stretchr/testify@v1.9.0
go mod tidy
```

- [ ] **Step 3: Commit**

```bash
cd ..
git add server/go.mod server/go.sum
git commit -m "chore: init Go module github.com/jry/lyshop"
```

---

## Task 4: Configuration System

**Files:**
- Create: `server/config/config.go`

- [ ] **Step 1: Write config.go**

```go
package config

import "github.com/spf13/viper"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Plugins  PluginsConfig  `mapstructure:"plugins"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	DSN     string `mapstructure:"dsn"`
	MaxOpen int    `mapstructure:"max_open"`
	MaxIdle int    `mapstructure:"max_idle"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type PluginsConfig struct {
	Enabled []string `mapstructure:"enabled"`
}

// Global is the loaded config, available after Load().
var Global Config

// Load reads the YAML file at path and unmarshals it into Global.
func Load(path string) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return viper.Unmarshal(&Global)
}
```

- [ ] **Step 2: Write test**

Create `server/config/config_test.go`:

```go
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
```

- [ ] **Step 3: Run test**

```bash
cd server
go test ./config/... -v
```
Expected: `PASS`

- [ ] **Step 4: Commit**

```bash
cd ..
git add server/config/
git commit -m "feat(core): config system with viper"
```

---

## Task 5: Database + Redis Clients

**Files:**
- Create: `server/core/db/db.go`
- Create: `server/core/cache/cache.go`

- [ ] **Step 1: Write db.go**

```go
package db

import (
	"fmt"
	"time"

	"github.com/jry/lyshop/server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global GORM instance, available after Init().
var DB *gorm.DB

func Init() error {
	cfg := config.Global.Database
	logLevel := logger.Silent
	if config.Global.Server.Mode == "debug" {
		logLevel = logger.Info
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("db connect: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return nil
}
```

- [ ] **Step 2: Write cache.go**

```go
package cache

import (
	"context"
	"fmt"

	"github.com/jry/lyshop/server/config"
	"github.com/redis/go-redis/v9"
)

// Client is the global Redis client, available after Init().
var Client *redis.Client

func Init() error {
	cfg := config.Global.Redis
	c := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := c.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("redis connect: %w", err)
	}
	Client = c
	return nil
}
```

- [ ] **Step 3: Commit**

```bash
git add server/core/db/ server/core/cache/
git commit -m "feat(core): DB and Redis client initializers"
```

---

## Task 6: Unified Response Helper

**Files:**
- Create: `server/core/response/response.go`

- [ ] **Step 1: Write response.go**

```go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type R struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// OK writes a successful response.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, R{Code: 0, Msg: "success", Data: data})
}

// Fail writes a business error response (HTTP 200, non-zero code).
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, R{Code: code, Msg: msg, Data: nil})
}

// Err returns a response struct (for use with AbortWithStatusJSON).
func Err(code int, msg string) R {
	return R{Code: code, Msg: msg, Data: nil}
}

// PageData wraps list + total for paginated responses.
type PageData struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}
```

- [ ] **Step 2: Write test**

Create `server/core/response/response_test.go`:

```go
package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	OK(c, map[string]string{"key": "value"})

	var r R
	json.Unmarshal(w.Body.Bytes(), &r)
	assert.Equal(t, 0, r.Code)
	assert.Equal(t, "success", r.Msg)
}

func TestFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Fail(c, 401, "未登录")

	assert.Equal(t, http.StatusOK, w.Code)
	var r R
	json.Unmarshal(w.Body.Bytes(), &r)
	assert.Equal(t, 401, r.Code)
}
```

- [ ] **Step 3: Run test**

```bash
cd server && go test ./core/response/... -v
```
Expected: `PASS`

- [ ] **Step 4: Commit**

```bash
cd .. && git add server/core/response/
git commit -m "feat(core): unified JSON response helper"
```

---

## Task 7: Plugin System

**Files:**
- Create: `server/core/plugin/plugin.go`
- Create: `server/core/plugin/registry.go`
- Create: `server/core/plugin/loader.go`

- [ ] **Step 1: Write plugin.go (interface + meta structs)**

```go
package plugin

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Meta holds data from a plugin's plugin.json.
type Meta struct {
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Author      string     `json:"author"`
	Depends     []string   `json:"depends"`
	Menus       []MenuItem `json:"menus"`
	Permissions []string   `json:"permissions"`
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
```

- [ ] **Step 2: Write registry.go**

```go
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
```

- [ ] **Step 3: Write loader.go**

```go
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
```

- [ ] **Step 4: Write test**

Create `server/core/plugin/plugin_test.go`:

```go
package plugin

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// stub is a minimal Plugin for testing.
type stub struct{ name string; deps []string }

func (s *stub) Meta() Meta          { return Meta{Name: s.name, Depends: s.deps} }
func (s *stub) RegisterRoutes(_, _ *gin.RouterGroup) {}
func (s *stub) Migrate(_ *gorm.DB) error              { return nil }
func (s *stub) Install() error                        { return nil }
func (s *stub) Uninstall() error                      { return nil }

func TestRegisterAndFind(t *testing.T) {
	// Reset registry for test isolation
	mu.Lock()
	registry = nil
	mu.Unlock()

	Register(&stub{name: "product"})
	Register(&stub{name: "order", deps: []string{"product"}})

	assert.NotNil(t, Find("product"))
	assert.NotNil(t, Find("order"))
	assert.Nil(t, Find("nonexistent"))
}

func TestLoad_MissingDependency(t *testing.T) {
	mu.Lock()
	registry = nil
	mu.Unlock()

	Register(&stub{name: "order", deps: []string{"product"}})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	front := r.Group("/api/v1")
	admin := r.Group("/admin/api")

	err := Load([]string{"order"}, nil, front, admin)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "product")
}

func TestLoad_OK(t *testing.T) {
	mu.Lock()
	registry = nil
	mu.Unlock()

	Register(&stub{name: "product"})
	Register(&stub{name: "order", deps: []string{"product"}})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	front := r.Group("/api/v1")
	admin := r.Group("/admin/api")

	err := Load([]string{"product", "order"}, &gorm.DB{}, front, admin)
	require.NoError(t, err)
}
```

- [ ] **Step 5: Run tests**

```bash
cd server && go test ./core/plugin/... -v
```
Expected: all 3 tests PASS

- [ ] **Step 6: Commit**

```bash
cd .. && git add server/core/plugin/
git commit -m "feat(core): plugin interface, registry, and loader"
```

---

## Task 8: Driver Abstraction Interfaces

**Files:**
- Create: `server/core/driver/payment/payment.go`
- Create: `server/core/driver/sms/sms.go`
- Create: `server/core/driver/oauth/oauth.go`
- Create: `server/core/driver/storage/storage.go`
- Create: `server/core/driver/ai/ai.go`

- [ ] **Step 1: Write payment/payment.go**

```go
package payment

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// OrderParams carries the data needed to create a payment order.
// Drivers read only the fields relevant to their platform.
type OrderParams struct {
	OrderNo     string
	Amount      int64  // cents (e.g. 9900 = ¥99.00)
	Description string
	NotifyURL   string
	// JSAPI (WeChat mini-program)
	OpenID string
	// H5 / App
	ClientIP string
}

// OrderResult is returned to the frontend to trigger payment.
type OrderResult struct {
	// PrepayID or equivalent pre-pay token
	PrepayID string
	// PayParams are passed directly to the frontend SDK (e.g. wx.requestPayment)
	PayParams map[string]string
}

type QueryResult struct {
	OutTradeNo string
	TradeNo    string // platform trade number
	Status     string // "paid" | "unpaid" | "refunded"
	Amount     int64
}

type RefundParams struct {
	OrderNo     string
	RefundNo    string
	Amount      int64 // refund amount in cents
	TotalAmount int64 // original order amount
	Reason      string
}

type RefundResult struct {
	RefundID string
}

type NotifyResult struct {
	OrderNo string
	Amount  int64
	Paid    bool
}

// Driver is the interface all payment plugins must implement.
type Driver interface {
	Name() string
	CreateOrder(ctx context.Context, p *OrderParams) (*OrderResult, error)
	QueryOrder(ctx context.Context, tradeNo string) (*QueryResult, error)
	Refund(ctx context.Context, p *RefundParams) (*RefundResult, error)
	HandleNotify(ctx context.Context, r *http.Request) (*NotifyResult, error)
}

var (
	mu      sync.RWMutex
	drivers = map[string]Driver{}
)

func Register(d Driver) { mu.Lock(); drivers[d.Name()] = d; mu.Unlock() }

func Get(name string) (Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	d, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("payment driver %q not registered", name)
	}
	return d, nil
}

// Names returns the names of all registered payment drivers.
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(drivers))
	for n := range drivers {
		out = append(out, n)
	}
	return out
}
```

- [ ] **Step 2: Write sms/sms.go**

```go
package sms

import (
	"context"
	"fmt"
	"sync"
)

// Driver is the interface all SMS plugins must implement.
type Driver interface {
	Name() string
	// Send sends templateCode with params to phone.
	// params is a map of template variable names to values.
	Send(ctx context.Context, phone, templateCode string, params map[string]string) error
}

var (
	mu     sync.RWMutex
	active Driver
)

func Register(d Driver) { mu.Lock(); active = d; mu.Unlock() }

func Get() (Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	if active == nil {
		return nil, fmt.Errorf("no SMS driver registered")
	}
	return active, nil
}
```

- [ ] **Step 3: Write oauth/oauth.go**

```go
package oauth

import (
	"context"
	"fmt"
	"sync"
)

// UserInfo is the normalized user profile returned by any OAuth provider.
type UserInfo struct {
	OpenID   string
	UnionID  string
	Nickname string
	Avatar   string
}

// Driver is the interface all OAuth login plugins must implement.
type Driver interface {
	Name() string
	// GetAuthURL returns the OAuth redirect URL for web flows.
	GetAuthURL(state string) string
	// HandleCallback exchanges the authorization code for a UserInfo.
	HandleCallback(ctx context.Context, code string) (*UserInfo, error)
	// GetUserInfo fetches the user profile using an access token.
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)
}

var (
	mu      sync.RWMutex
	drivers = map[string]Driver{}
)

func Register(d Driver) { mu.Lock(); drivers[d.Name()] = d; mu.Unlock() }

func Get(name string) (Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	d, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("oauth driver %q not registered", name)
	}
	return d, nil
}
```

- [ ] **Step 4: Write storage/storage.go**

```go
package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"sync"
)

// UploadResult holds the stored path and public URL.
type UploadResult struct {
	Path string // relative storage path (used for deletion)
	URL  string // full public URL
	Size int64
	Mime string
}

// Driver is the interface all file storage plugins must implement.
type Driver interface {
	Name() string
	Upload(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error)
	Delete(ctx context.Context, path string) error
	// GetURL converts a stored path to its full public URL.
	GetURL(path string) string
}

var (
	mu     sync.RWMutex
	active Driver
)

func Register(d Driver) { mu.Lock(); active = d; mu.Unlock() }

func Get() (Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	if active == nil {
		return nil, fmt.Errorf("no storage driver registered")
	}
	return active, nil
}
```

- [ ] **Step 5: Write ai/ai.go**

```go
package ai

import (
	"context"
	"fmt"
	"sync"
)

// GenerateParams are passed to the image generation model.
type GenerateParams struct {
	Prompt    string
	NegPrompt string
	Width     int
	Height    int
	Count     int    // number of images to generate (1-5)
	Style     string // e.g. "ecommerce", "realistic", "illustration"
}

// GenerateResult holds the output image URLs.
type GenerateResult struct {
	URLs []string
}

// Driver is the interface all AI image generation model drivers must implement.
type Driver interface {
	Name() string
	Generate(ctx context.Context, p *GenerateParams) (*GenerateResult, error)
}

var (
	mu      sync.RWMutex
	drivers = map[string]Driver{}
	def     string // default driver name
)

func Register(d Driver, isDefault bool) {
	mu.Lock()
	defer mu.Unlock()
	drivers[d.Name()] = d
	if isDefault || def == "" {
		def = d.Name()
	}
}

func Get(name string) (Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	if name == "" {
		name = def
	}
	d, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("ai driver %q not registered", name)
	}
	return d, nil
}

// Names returns all registered AI driver names.
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(drivers))
	for n := range drivers {
		out = append(out, n)
	}
	return out
}
```

- [ ] **Step 6: Commit**

```bash
git add server/core/driver/
git commit -m "feat(core): payment, sms, oauth, storage, ai driver interfaces"
```

---

## Task 9: Core Data Models

**Files:**
- Create: `server/model/base.go`
- Create: `server/model/user.go`
- Create: `server/model/admin.go`
- Create: `server/model/role.go`
- Create: `server/model/config_kv.go`

- [ ] **Step 1: Write base.go**

```go
package model

import "time"

// Base provides common fields for all models.
// Use gorm:"autoCreateTime;autoUpdateTime" behavior.
type Base struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

- [ ] **Step 2: Write user.go**

```go
package model

// User is the end-customer account.
type User struct {
	Base
	MerchantID uint64 `gorm:"not null;default:0;index" json:"merchant_id"` // reserved
	Phone      string `gorm:"size:20;uniqueIndex"      json:"phone"`
	Nickname   string `gorm:"size:64"                  json:"nickname"`
	Avatar     string `gorm:"size:500"                 json:"avatar"`
	Points     int    `gorm:"not null;default:0"       json:"points"`
	Status     int8   `gorm:"not null;default:1"       json:"status"` // 1=active 0=banned
}
```

- [ ] **Step 3: Write admin.go**

```go
package model

// Admin is a back-office user.
type Admin struct {
	Base
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null"            json:"-"` // bcrypt hash
	RoleID   uint64 `gorm:"not null"                     json:"role_id"`
	Status   int8   `gorm:"not null;default:1"           json:"status"`
}
```

- [ ] **Step 4: Write role.go**

```go
package model

import "encoding/json"

// Role defines a set of permission strings for admins.
type Role struct {
	Base
	Name        string          `gorm:"size:64;not null" json:"name"`
	Permissions json.RawMessage `gorm:"type:json"        json:"permissions"` // []string
}
```

- [ ] **Step 5: Write config_kv.go**

```go
package model

// ConfigKV stores plugin-namespaced key-value configuration.
type ConfigKV struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Plugin string `gorm:"size:64;not null;uniqueIndex:uk_plugin_key" json:"plugin"`
	Key    string `gorm:"size:128;not null;uniqueIndex:uk_plugin_key" json:"key"`
	Value  string `gorm:"type:text"                                  json:"value"`
}

func (ConfigKV) TableName() string { return "configs" }
```

- [ ] **Step 6: Commit**

```bash
git add server/model/
git commit -m "feat(model): core GORM models (user, admin, role, config_kv)"
```

---

## Task 10: Middleware

**Files:**
- Create: `server/core/middleware/auth.go`
- Create: `server/core/middleware/cors.go`
- Create: `server/core/middleware/logger.go`

- [ ] **Step 1: Write auth.go**

```go
package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jry/lyshop/server/config"
	"github.com/jry/lyshop/server/core/response"
)

// Claims is the JWT payload used for both user and admin tokens.
type Claims struct {
	UserID uint64   `json:"user_id"`
	Role   string   `json:"role"` // "user" | "admin"
	Perms  []string `json:"perms,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken signs a JWT for userID with the given role and permissions.
func GenerateToken(userID uint64, role string, perms []string) (string, error) {
	expiry := time.Duration(config.Global.JWT.ExpireHours) * time.Hour
	claims := Claims{
		UserID: userID,
		Role:   role,
		Perms:  perms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(config.Global.JWT.Secret))
}

// ParseToken validates tokenStr and returns the Claims.
func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims,
		func(_ *jwt.Token) (any, error) {
			return []byte(config.Global.JWT.Secret), nil
		},
	)
	return claims, err
}

// RequireAuth aborts with 401 if the request has no valid JWT.
func RequireAuth(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusOK, response.Err(401, "请先登录"))
		return
	}
	claims, err := ParseToken(strings.TrimPrefix(auth, "Bearer "))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, response.Err(401, "Token无效或已过期"))
		return
	}
	c.Set("user_id", claims.UserID)
	c.Set("role", claims.Role)
	c.Set("perms", claims.Perms)
	c.Next()
}

// RequireAdmin calls RequireAuth, then additionally checks role == "admin".
func RequireAdmin(c *gin.Context) {
	RequireAuth(c)
	if c.IsAborted() {
		return
	}
	if role, _ := c.Get("role"); role != "admin" {
		c.AbortWithStatusJSON(http.StatusOK, response.Err(403, "无权限"))
		return
	}
	c.Next()
}
```

- [ ] **Step 2: Write cors.go**

```go
package middleware

import (
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers",
			"Content-Type,Authorization,X-Requested-With")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
```

- [ ] **Step 3: Write logger.go**

```go
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"log/slog"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		slog.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
			"ip", c.ClientIP(),
		)
	}
}
```

- [ ] **Step 4: Commit**

```bash
git add server/core/middleware/
git commit -m "feat(core): JWT auth, CORS, and logger middleware"
```

---

## Task 11: Auth Service + API Handlers

**Files:**
- Create: `server/service/auth/user_auth.go`
- Create: `server/service/auth/admin_auth.go`
- Create: `server/api/auth/front.go`
- Create: `server/api/auth/admin.go`

- [ ] **Step 1: Write user_auth.go**

```go
package auth

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jry/lyshop/server/core/cache"
	"github.com/jry/lyshop/server/core/db"
	"github.com/jry/lyshop/server/core/middleware"
	"github.com/jry/lyshop/server/model"
	"gorm.io/gorm"
)

const smsCodeTTL = 5 * time.Minute
const smsCodeKeyPrefix = "sms:code:"

// SendSMSCode generates a 6-digit code and stores it in Redis for 5 min.
// In production the SMS driver sends it; here we just store it (dev mode).
func SendSMSCode(ctx context.Context, phone string) (string, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1000000)) //nolint:gosec
	key := smsCodeKeyPrefix + phone
	return code, cache.Client.Set(ctx, key, code, smsCodeTTL).Err()
}

// VerifySMSCode checks code against the stored value and deletes it.
func VerifySMSCode(ctx context.Context, phone, code string) error {
	key := smsCodeKeyPrefix + phone
	stored, err := cache.Client.Get(ctx, key).Result()
	if err != nil {
		return errors.New("验证码已过期")
	}
	if stored != code {
		return errors.New("验证码错误")
	}
	cache.Client.Del(ctx, key)
	return nil
}

// SMSLogin verifies the code, then finds or creates the user, and returns a JWT.
func SMSLogin(ctx context.Context, phone, code string) (string, error) {
	if err := VerifySMSCode(ctx, phone, code); err != nil {
		return "", err
	}

	var user model.User
	err := db.DB.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = model.User{
			Phone:    phone,
			Nickname: "用户" + phone[len(phone)-4:],
			Status:   1,
		}
		if err2 := db.DB.WithContext(ctx).Create(&user).Error; err2 != nil {
			return "", err2
		}
	} else if err != nil {
		return "", err
	}

	if user.Status == 0 {
		return "", errors.New("账号已被禁用")
	}

	return middleware.GenerateToken(user.ID, "user", nil)
}
```

- [ ] **Step 2: Write admin_auth.go**

```go
package auth

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jry/lyshop/server/core/db"
	"github.com/jry/lyshop/server/core/middleware"
	"github.com/jry/lyshop/server/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminLogin validates username/password and returns a JWT with permissions.
func AdminLogin(ctx context.Context, username, password string) (string, error) {
	var admin model.Admin
	err := db.DB.WithContext(ctx).Where("username = ?", username).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", errors.New("用户名或密码错误")
	}
	if err != nil {
		return "", err
	}
	if admin.Status == 0 {
		return "", errors.New("账号已被禁用")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	var role model.Role
	if err = db.DB.WithContext(ctx).First(&role, admin.RoleID).Error; err != nil {
		return "", err
	}
	var perms []string
	json.Unmarshal(role.Permissions, &perms) //nolint:errcheck

	return middleware.GenerateToken(admin.ID, "admin", perms)
}

// HashPassword returns the bcrypt hash of plain.
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}
```

- [ ] **Step 3: Write front.go (user auth API)**

```go
package auth

import (
	"github.com/gin-gonic/gin"
	authsvc "github.com/jry/lyshop/server/service/auth"
	"github.com/jry/lyshop/server/core/response"
)

// RegisterFrontRoutes adds user auth routes to the front-end router group.
func RegisterFrontRoutes(g *gin.RouterGroup) {
	g.POST("/auth/sms/send", sendSMSCode)
	g.POST("/auth/sms/login", smsLogin)
}

func sendSMSCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required,len=11"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "手机号格式错误")
		return
	}
	// In dev mode, return the code in the response for easy testing.
	code, err := authsvc.SendSMSCode(c.Request.Context(), req.Phone)
	if err != nil {
		response.Fail(c, 500, "发送失败: "+err.Error())
		return
	}
	response.OK(c, gin.H{"dev_code": code}) // remove dev_code in production
}

func smsLogin(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code"  binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	token, err := authsvc.SMSLogin(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		response.Fail(c, 10001, err.Error())
		return
	}
	response.OK(c, gin.H{"token": token})
}
```

- [ ] **Step 4: Write admin.go (admin auth API)**

```go
package auth

import (
	"github.com/gin-gonic/gin"
	authsvc "github.com/jry/lyshop/server/service/auth"
	"github.com/jry/lyshop/server/core/response"
)

// RegisterAdminRoutes adds admin auth routes (no JWT required for login).
func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.POST("/auth/login", adminLogin)
}

func adminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	token, err := authsvc.AdminLogin(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}
	response.OK(c, gin.H{"token": token})
}
```

- [ ] **Step 5: Commit**

```bash
git add server/service/ server/api/
git commit -m "feat(auth): SMS user login and admin password login"
```

---

## Task 12: Router + App + main.go

**Files:**
- Create: `server/core/app/app.go`
- Create: `server/main.go`

- [ ] **Step 1: Write app.go**

```go
package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	apiauth "github.com/jry/lyshop/server/api/auth"
	"github.com/jry/lyshop/server/config"
	"github.com/jry/lyshop/server/core/cache"
	"github.com/jry/lyshop/server/core/db"
	"github.com/jry/lyshop/server/core/middleware"
	"github.com/jry/lyshop/server/core/plugin"
)

// Init wires config → db → redis → plugin system.
func Init(cfgPath string) error {
	if err := config.Load(cfgPath); err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	gin.SetMode(config.Global.Server.Mode)

	if err := db.Init(); err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	if err := cache.Init(); err != nil {
		return fmt.Errorf("init cache: %w", err)
	}
	return nil
}

// Run builds the Gin engine, loads plugins, and starts the HTTP server.
func Run() error {
	r := gin.New()
	r.Use(middleware.Logger(), middleware.CORS(), gin.Recovery())

	// Route groups
	front := r.Group("/api/v1")
	admin := r.Group("/admin/api")
	admin.Use(middleware.RequireAdmin)

	// Core auth routes (no auth required)
	apiauth.RegisterFrontRoutes(front)
	adminPublic := r.Group("/admin/api")
	apiauth.RegisterAdminRoutes(adminPublic)

	// Admin: menu endpoint (returns menus for enabled plugins)
	admin.GET("/menus", func(c *gin.Context) {
		menus := plugin.EnabledMenus(config.Global.Plugins.Enabled)
		c.JSON(200, menus)
	})

	// Load plugins
	if err := plugin.Load(config.Global.Plugins.Enabled, db.DB, front, admin); err != nil {
		return fmt.Errorf("load plugins: %w", err)
	}

	addr := fmt.Sprintf(":%d", config.Global.Server.Port)
	return r.Run(addr)
}
```

- [ ] **Step 2: Write main.go**

```go
package main

import (
	"flag"
	"log"

	"github.com/jry/lyshop/server/core/app"
	// Blank-import enabled plugins so their init() registers them.
	// Add one line per plugin as you build Phase 2+:
	// _ "github.com/jry/lyshop/server/plugins/product"
)

func main() {
	cfg := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	if err := app.Init(*cfg); err != nil {
		log.Fatalf("init: %v", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}
```

- [ ] **Step 3: Create plugins placeholder**

```bash
mkdir -p server/plugins
touch server/plugins/.gitkeep
```

- [ ] **Step 4: Verify it compiles**

```bash
cd server
go build ./...
```
Expected: no errors (exits 0)

- [ ] **Step 5: Commit**

```bash
cd ..
git add server/core/app/ server/main.go server/plugins/.gitkeep
git commit -m "feat(core): app lifecycle, Gin router, and main entry point"
```

---

## Task 13: Docker Compose + Dockerfile + Nginx

**Files:**
- Create: `server/Dockerfile`
- Create: `docker-compose.yml`
- Create: `nginx.conf`

- [ ] **Step 1: Write server/Dockerfile**

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o lyshop-server .

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/lyshop-server .
EXPOSE 8080
CMD ["./lyshop-server", "--config", "/app/config.yaml"]
```

- [ ] **Step 2: Write nginx.conf**

```nginx
events { worker_connections 1024; }

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
    sendfile      on;
    gzip          on;
    gzip_types    text/plain text/css application/json application/javascript;

    # Upstream for Go backend
    upstream backend { server server:8080; }

    server {
        listen 80;
        client_max_body_size 50m;

        # API + WebSocket → Go backend
        location /api/  { proxy_pass http://backend; proxy_set_header Host $host; proxy_set_header X-Real-IP $remote_addr; }
        location /admin/api/ { proxy_pass http://backend; proxy_set_header Host $host; }
        location /notify/ { proxy_pass http://backend; }
        location /ws/ {
            proxy_pass http://backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
        }

        # Admin Vue3 SPA
        location /admin {
            proxy_pass http://admin:80;
            proxy_set_header Host $host;
        }

        # uni-app H5 SPA (catch-all)
        location / {
            proxy_pass http://app_h5:80;
            proxy_set_header Host $host;
        }
    }
}
```

- [ ] **Step 3: Write docker-compose.yml**

```yaml
version: '3.8'

services:
  nginx:
    image: nginx:1.25-alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - server
      - admin
      - app_h5
    restart: unless-stopped

  server:
    build: ./server
    volumes:
      - ./config.yaml:/app/config.yaml:ro
      - ./data/uploads:/app/uploads
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped

  admin:
    build: ./admin
    restart: unless-stopped

  app_h5:
    build: ./app
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD:-lyshop123}
      MYSQL_DATABASE: lyshop
      MYSQL_CHARSET: utf8mb4
    volumes:
      - ./data/mysql:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - ./data/redis:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5
    restart: unless-stopped
```

- [ ] **Step 4: Commit**

```bash
git add server/Dockerfile docker-compose.yml nginx.conf
git commit -m "feat(deploy): Docker Compose + Nginx + server Dockerfile"
```

---

## Task 14: Admin Frontend Scaffold (Vue3 + TailwindCSS)

**Files:**
- Create: `admin/package.json`
- Create: `admin/vite.config.ts`
- Create: `admin/tailwind.config.js`
- Create: `admin/index.html`
- Create: `admin/Dockerfile`
- Create: `admin/src/main.ts`
- Create: `admin/src/App.vue`
- Create: `admin/src/router/index.ts`
- Create: `admin/src/stores/auth.ts`
- Create: `admin/src/api/request.ts`
- Create: `admin/src/api/auth.ts`
- Create: `admin/src/layouts/AdminLayout.vue`
- Create: `admin/src/views/Login.vue`
- Create: `admin/src/views/Dashboard.vue`

- [ ] **Step 1: Create admin/package.json**

```json
{
  "name": "lyshop-admin",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.3.0",
    "pinia": "^2.1.0",
    "axios": "^1.7.0",
    "@vueuse/core": "^10.10.0",
    "lucide-vue-next": "^0.378.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.0",
    "typescript": "^5.4.0",
    "vite": "^5.2.0",
    "vue-tsc": "^2.0.0",
    "tailwindcss": "^3.4.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0"
  }
}
```

- [ ] **Step 2: Create admin/vite.config.ts**

```ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { '@': resolve(__dirname, 'src') }
  },
  server: {
    port: 9527,
    proxy: {
      '/admin/api': { target: 'http://localhost:8080', changeOrigin: true },
      '/api':       { target: 'http://localhost:8080', changeOrigin: true }
    }
  }
})
```

- [ ] **Step 3: Create admin/tailwind.config.js**

```js
/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,ts,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: { DEFAULT: '#1e40af', light: '#3b82f6' }
      }
    }
  },
  plugins: []
}
```

- [ ] **Step 4: Create admin/index.html**

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>lyshop 管理后台</title>
</head>
<body>
  <div id="app"></div>
  <script type="module" src="/src/main.ts"></script>
</body>
</html>
```

- [ ] **Step 5: Create admin/src/main.ts**

```ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './style.css'  // TailwindCSS base

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
```

Create `admin/src/style.css`:
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

- [ ] **Step 6: Create admin/src/App.vue**

```vue
<template>
  <router-view />
</template>
```

- [ ] **Step 7: Create admin/src/api/request.ts**

```ts
import axios from 'axios'

const request = axios.create({
  baseURL: '/admin/api',
  timeout: 30000,
})

request.interceptors.request.use(config => {
  const token = localStorage.getItem('admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

request.interceptors.response.use(
  res => {
    const { code, msg, data } = res.data
    if (code !== 0) return Promise.reject(new Error(msg || '请求失败'))
    return data
  },
  err => Promise.reject(err)
)

export default request
```

- [ ] **Step 7: Create admin/src/api/auth.ts**

```ts
import request from './request'

export const login = (username: string, password: string) =>
  request.post<never, { token: string }>('/auth/login', { username, password })

export const getMenus = () =>
  request.get<never, any[]>('/menus')
```

- [ ] **Step 8: Create admin/src/stores/auth.ts**

```ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginAPI } from '@/api/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')

  async function loginAction(username: string, password: string) {
    const data = await loginAPI(username, password)
    token.value = data.token
    localStorage.setItem('admin_token', data.token)
    router.push('/dashboard')
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('admin_token')
    router.push('/login')
  }

  return { token, loginAction, logout }
})
```

- [ ] **Step 9: Create admin/src/router/index.ts**

```ts
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory('/admin'),
  routes: [
    { path: '/login', component: () => import('@/views/Login.vue') },
    {
      path: '/',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', component: () => import('@/views/Dashboard.vue') },
      ]
    }
  ]
})

router.beforeEach(to => {
  const token = localStorage.getItem('admin_token')
  if (to.meta.requiresAuth && !token) return '/login'
})

export default router
```

- [ ] **Step 10: Create admin/src/views/Login.vue**

```vue
<template>
  <div class="min-h-screen bg-slate-100 flex items-center justify-center">
    <div class="bg-white rounded-2xl shadow-sm p-10 w-96">
      <h1 class="text-2xl font-bold text-slate-800 mb-8 text-center">lyshop 管理后台</h1>
      <form @submit.prevent="handleLogin" class="space-y-4">
        <input
          v-model="form.username"
          type="text"
          placeholder="用户名"
          class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-primary"
        />
        <input
          v-model="form.password"
          type="password"
          placeholder="密码"
          class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-primary"
        />
        <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-primary text-white rounded-xl py-3 text-sm font-medium hover:bg-primary-light transition"
        >
          {{ loading ? '登录中...' : '登 录' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const form = ref({ username: '', password: '' })
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.loginAction(form.value.username, form.value.password)
  } catch (e: any) {
    error.value = e.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>
```

- [ ] **Step 11: Create admin/src/layouts/AdminLayout.vue**

```vue
<template>
  <div class="flex h-screen bg-slate-50">
    <!-- Sidebar -->
    <aside class="w-64 bg-slate-900 text-slate-100 flex flex-col shrink-0">
      <div class="h-16 flex items-center px-6 border-b border-slate-800">
        <span class="text-lg font-bold text-white">lyshop</span>
      </div>
      <nav class="flex-1 overflow-y-auto py-4 space-y-1 px-3">
        <router-link
          v-for="item in menus"
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm hover:bg-slate-800 transition"
          active-class="bg-primary text-white"
        >
          <span>{{ item.title }}</span>
        </router-link>
      </nav>
    </aside>

    <!-- Main area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Topbar -->
      <header class="h-16 bg-white border-b border-slate-100 flex items-center justify-between px-6 shrink-0">
        <div class="text-sm text-slate-500 breadcrumb">{{ $route.name }}</div>
        <button @click="auth.logout()" class="text-sm text-slate-500 hover:text-slate-800">退出</button>
      </header>
      <!-- Content -->
      <main class="flex-1 overflow-y-auto p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { getMenus } from '@/api/auth'

const auth = useAuthStore()
const menus = ref<any[]>([])

onMounted(async () => {
  try { menus.value = await getMenus() } catch {}
})
</script>
```

- [ ] **Step 12: Create admin/src/views/Dashboard.vue**

```vue
<template>
  <div>
    <h2 class="text-xl font-semibold text-slate-800 mb-6">Dashboard</h2>
    <div class="grid grid-cols-4 gap-4">
      <div v-for="card in cards" :key="card.label"
        class="bg-white rounded-xl shadow-sm p-6">
        <p class="text-sm text-slate-500">{{ card.label }}</p>
        <p class="text-2xl font-bold text-slate-800 mt-1">{{ card.value }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const cards = [
  { label: '今日订单', value: '—' },
  { label: '今日销售额', value: '—' },
  { label: '待处理售后', value: '—' },
  { label: '在线客服会话', value: '—' },
]
</script>
```

- [ ] **Step 13: Create admin/Dockerfile**

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:1.25-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

Create `admin/nginx.conf`:
```nginx
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

- [ ] **Step 14: Install deps and verify build**

```bash
cd admin
npm install
npm run build
```
Expected: `dist/` directory created with no errors.

- [ ] **Step 15: Commit**

```bash
cd ..
git add admin/
git commit -m "feat(admin): Vue3 + TailwindCSS admin skeleton with login + layout"
```

---

## Task 15: uni-app Frontend Scaffold

**Files:**
- Create: `app/package.json`
- Create: `app/manifest.json`
- Create: `app/pages.json`
- Create: `app/uni.scss`
- Create: `app/App.vue`
- Create: `app/main.ts`
- Create: `app/utils/request.ts`
- Create: `app/pages/index/index.vue`
- Create: `app/pages/login/index.vue`
- Create: `app/Dockerfile`

- [ ] **Step 1: Create app/package.json**

```json
{
  "name": "lyshop-app",
  "version": "1.0.0",
  "scripts": {
    "dev:h5": "uni",
    "build:h5": "uni build",
    "dev:mp-weixin": "uni -p mp-weixin",
    "build:mp-weixin": "uni build -p mp-weixin"
  },
  "dependencies": {
    "@dcloudio/uni-app": "*",
    "@dcloudio/uni-h5": "*",
    "@dcloudio/uni-mp-weixin": "*",
    "pinia": "^2.1.0",
    "uview-plus": "3.1.31"
  },
  "devDependencies": {
    "@dcloudio/types": "*",
    "typescript": "^5.4.0",
    "vite": "^5.2.0"
  }
}
```

- [ ] **Step 2: Create app/manifest.json**

```json
{
  "name": "lyshop",
  "appid": "__UNI__LYSHOP",
  "description": "lyshop 商城",
  "versionName": "1.0.0",
  "versionCode": "100",
  "transformPx": false,
  "mp-weixin": {
    "appid": "",
    "setting": { "urlCheck": false },
    "usingComponents": true
  },
  "h5": {
    "devServer": { "port": 9528 },
    "router": { "mode": "history", "base": "/" }
  },
  "app-plus": {
    "usingComponents": true,
    "splashscreen": { "alwaysShowBeforeRender": true }
  }
}
```

- [ ] **Step 3: Create app/pages.json**

```json
{
  "pages": [
    {
      "path": "pages/index/index",
      "style": { "navigationBarTitleText": "首页" }
    },
    {
      "path": "pages/login/index",
      "style": { "navigationBarTitleText": "登录" }
    }
  ],
  "tabBar": {
    "color": "#666",
    "selectedColor": "#1e40af",
    "list": [
      { "pagePath": "pages/index/index",  "text": "首页" },
      { "pagePath": "pages/login/index",  "text": "我的" }
    ]
  },
  "globalStyle": {
    "navigationBarTextStyle": "black",
    "navigationBarBackgroundColor": "#FFFFFF",
    "backgroundColor": "#F8FAFC"
  }
}
```

- [ ] **Step 4: Create app/uni.scss**

```scss
// uview-plus theme variables
$u-primary: #1e40af;
$u-primary-dark: #1e3a8a;
$u-primary-disabled: #93c5fd;
$u-primary-light: #eff6ff;

// Import uview-plus styles
@import 'uview-plus/index.scss';
```

- [ ] **Step 5: Create app/App.vue**

```vue
<script>
export default {
  onLaunch() {},
  onShow() {},
  onHide() {}
}
</script>
<style lang="scss">
@import './uni.scss';
</style>
```

- [ ] **Step 6: Create app/main.ts**

```ts
import { createSSRApp } from 'vue'
import App from './App.vue'
import uviewPlus from 'uview-plus'

export function createApp() {
  const app = createSSRApp(App)
  app.use(uviewPlus)
  return { app }
}
```

- [ ] **Step 7: Create app/utils/request.ts**

```ts
const BASE_URL = 'http://localhost:8080'

function getToken() {
  return uni.getStorageSync('user_token') || ''
}

export function request<T = any>(options: UniNamespace.RequestOptions): Promise<T> {
  return new Promise((resolve, reject) => {
    uni.request({
      ...options,
      url: BASE_URL + options.url,
      header: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${getToken()}`,
        ...(options.header || {})
      },
      success(res) {
        const data = res.data as any
        if (data.code !== 0) {
          uni.showToast({ title: data.msg || '请求失败', icon: 'none' })
          reject(new Error(data.msg))
        } else {
          resolve(data.data)
        }
      },
      fail(err) { reject(err) }
    })
  })
}

export const get = <T>(url: string, data?: any) =>
  request<T>({ url, method: 'GET', data })

export const post = <T>(url: string, data?: any) =>
  request<T>({ url, method: 'POST', data })
```

- [ ] **Step 8: Create app/pages/index/index.vue**

```vue
<template>
  <view class="container">
    <u-navbar title="lyshop" :placeholder="true" />
    <view class="p-4">
      <u-skeleton :loading="true" :rows="3" v-if="false" />
      <view class="text-center py-20 text-gray-400">
        <text>首页装修加载中...</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
// Decor plugin will populate this page in Phase 4
</script>
```

- [ ] **Step 9: Create app/pages/login/index.vue**

```vue
<template>
  <view class="min-h-screen bg-gray-50 flex flex-col items-center justify-center px-6">
    <view class="w-full max-w-sm">
      <text class="text-2xl font-bold text-slate-800 block text-center mb-8">登录</text>

      <u-form :model="form" ref="formRef">
        <u-form-item label="手机号" prop="phone">
          <u-input v-model="form.phone" placeholder="请输入手机号" type="number" maxlength="11" />
        </u-form-item>
        <u-form-item label="验证码" prop="code">
          <u-input v-model="form.code" placeholder="请输入验证码" type="number" maxlength="6">
            <template #suffix>
              <u-button size="mini" :disabled="countdown > 0" @click="sendCode">
                {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
              </u-button>
            </template>
          </u-input>
        </u-form-item>
      </u-form>

      <u-button type="primary" block class="mt-6" @click="handleLogin" :loading="loading">
        登 录
      </u-button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { get, post } from '@/utils/request'

const form = ref({ phone: '', code: '' })
const loading = ref(false)
const countdown = ref(0)

async function sendCode() {
  if (!form.value.phone || form.value.phone.length !== 11) {
    return uni.showToast({ title: '请输入正确手机号', icon: 'none' })
  }
  const data = await get<{ dev_code: string }>(`/api/v1/auth/sms/send`,
    { phone: form.value.phone })
  // dev only: auto fill the code
  if (data?.dev_code) form.value.code = data.dev_code

  countdown.value = 60
  const t = setInterval(() => {
    if (--countdown.value <= 0) clearInterval(t)
  }, 1000)
}

async function handleLogin() {
  loading.value = true
  try {
    const data = await post<{ token: string }>('/api/v1/auth/sms/login', form.value)
    uni.setStorageSync('user_token', data.token)
    uni.switchTab({ url: '/pages/index/index' })
  } catch {} finally {
    loading.value = false
  }
}
</script>
```

- [ ] **Step 10: Create app/Dockerfile**

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build:h5

FROM nginx:1.25-alpine
COPY --from=builder /app/dist/build/h5 /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

Create `app/nginx.conf`:
```nginx
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

- [ ] **Step 11: Commit**

```bash
git add app/
git commit -m "feat(app): uni-app + uview-plus 3.x scaffold (H5 + mini + App)"
```

---

## Task 16: Push to GitHub + Verify

- [ ] **Step 1: Push all commits**

```bash
git push
```

Expected: all commits pushed to `github.com/<user>/lyshop`

- [ ] **Step 2: Create GitHub Actions CI**

Create `.github/workflows/ci.yml`:

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  go-test:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: server
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache-dependency-path: server/go.sum
      - run: go test ./... -v
      - run: go build ./...

  admin-build:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: admin
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: admin/package-lock.json
      - run: npm install
      - run: npm run build
```

- [ ] **Step 3: Commit and push CI**

```bash
git add .github/
git commit -m "ci: GitHub Actions for Go tests and admin build"
git push
```

- [ ] **Step 4: Verify CI passes**

```bash
gh run list --limit 3
```
Expected: latest run shows green ✓

---

## Phase 1 Complete ✓

After this plan, the repo has:
- GitHub repo `<user>/lyshop` with CI
- Go backend: config, DB, Redis, plugin system, 5 driver interfaces, JWT auth APIs
- Docker Compose: one-command deployment
- Admin frontend: Vue3 login + layout skeleton running on `:9527`
- uni-app: skeleton running as H5 on `:9528` and buildable for WeChat mini-program + App

## Next Plans

| Plan | File | Covers |
|------|------|--------|
| Phase 2 | `2026-05-21-lyshop2-phase2-product-order-wms.md` | product, order, wms plugins + admin pages |
| Phase 3 | `2026-05-21-lyshop2-phase3-commerce.md` | marketing, wechat_pay, alipay, sms, wechat_auth |
| Phase 4 | `2026-05-21-lyshop2-phase4-im-ai-decor.md` | im (WebSocket), ai_image, decor |
| Phase 5 | `2026-05-21-lyshop2-phase5-storage.md` | storage_local, storage_oss, storage_cos, storage_qiniu |
| Phase 6 | `2026-05-21-lyshop2-phase6-admin-pages.md` | All Vue3 admin pages per plugin |
| Phase 7 | `2026-05-21-lyshop2-phase7-app-pages.md` | All uni-app pages per plugin |
