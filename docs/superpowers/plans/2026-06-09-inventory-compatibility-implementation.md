# Inventory Compatibility Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a three-mode inventory architecture that lets LYShop run without WMS by default, use the built-in WMS when enabled, or integrate with an external WMS in sync/async modes.

**Architecture:** Introduce a new `server/core/inventory` domain that owns configuration parsing, provider selection, shared inventory status models, and external integration task state. Refactor `order`, `product`, and `points_mall` to depend on the inventory router instead of directly importing `plugins/wms/service`, then wrap the existing WMS service as a `builtin_wms` provider and add a new `external_wms` provider with sync calls and async retry tasks.

**Tech Stack:** Go 1.26, Gin, GORM, YAML config via Viper, existing plugin registry, docs-site markdown.

---

## File Structure

- Create `server/core/inventory/types.go`: shared input/output structs (`ReserveInput`, `DeductInput`, `SellableStock`, `TaskPayload`).
- Create `server/core/inventory/errors.go`: typed inventory errors and helpers.
- Create `server/core/inventory/provider.go`: provider interface and registration helpers.
- Create `server/core/inventory/router.go`: runtime provider selection and capability helpers.
- Create `server/core/inventory/router_test.go`: config-aware provider selection tests.
- Create `server/core/inventory/model.go`: shared `InventoryReservation`, `InventoryIntegrationTask`, and `OrderInventoryState` models.
- Create `server/core/inventory/local.go`: local inventory provider using `product_skus.stock`.
- Create `server/core/inventory/local_test.go`: reservation/confirm/release tests for local mode.
- Create `server/core/inventory/external.go`: external WMS provider sync call implementation.
- Create `server/core/inventory/external_test.go`: request signing, sync error mapping, and async enqueue tests.
- Create `server/core/inventory/tasks.go`: async task enqueue, claim, retry, and callback update logic.
- Create `server/core/inventory/tasks_test.go`: async retry/task state tests.
- Modify `server/config/config.go`: add `InventoryConfig` and `ExternalWMSConfig`.
- Modify `server/config.example.yaml` and `config.example.yaml`: document `inventory` and `external_wms` settings; remove `wms` from default enabled list in the root example.
- Modify `server/core/app/app.go`: migrate shared inventory tables before plugin load and validate config/provider consistency.
- Modify `server/core/plugin/plugin_test.go`: assert `order` no longer requires `wms`.
- Modify `server/plugins/order/plugin.json`: remove hard dependency on `wms`.
- Modify `server/plugins/order/model/order.go`: add `InventoryStatus` field.
- Modify `server/plugins/order/service/order.go`: replace direct `wms/service` calls with inventory router operations.
- Create `server/plugins/order/service/order_inventory_test.go`: local/builtin/external inventory flow tests for order create/pay/cancel.
- Modify `server/plugins/product/service/product.go`: replace direct WMS stock sync with inventory `SyncSkuTx`.
- Create `server/plugins/product/service/product_inventory_test.go`: verify new SKU sync under local/builtin/external modes.
- Modify `server/plugins/points_mall/service/exchange.go`: replace direct stock update with inventory `DeductTx`/`RestoreTx`.
- Create `server/plugins/points_mall/service/exchange_inventory_test.go`: verify stock deduction/restore through inventory provider.
- Modify `server/plugins/wms/plugin.go`: register builtin inventory provider.
- Create `server/plugins/wms/provider.go`: provider wrapper over existing WMS service functions.
- Create `server/plugins/wms/provider_test.go`: provider adapter tests.
- Modify `server/plugins/wms/service/reservation.go` and `server/plugins/wms/service/stock.go`: export read helpers required by provider wrapper if needed.
- Create `server/plugins/external_wms/plugin.go`: register external provider and callback/task admin routes.
- Create `server/plugins/external_wms/plugin.json`: metadata/config items for external WMS.
- Create `server/plugins/external_wms/api/admin.go`: callback, task list, retry routes.
- Create `server/plugins/external_wms/api/admin_test.go`: callback/task route tests.
- Modify `server/main.go`: blank-import the new `external_wms` plugin.
- Modify `docs-site/docs/api/order.md`: describe latest inventory status semantics.
- Modify `docs-site/docs/api/stock-reservation.md`: rewrite around unified inventory architecture.
- Modify `docs-site/docs/api/wms.md`: clarify WMS is a built-in provider, not the only mode.
- Modify `docs-site/docs/guide/features.md`: update feature summary for local/builtin/external inventory support.

## Implementation Tasks

### Task 1: Core Inventory Configuration And Provider Router

**Files:**
- Create: `server/core/inventory/types.go`
- Create: `server/core/inventory/errors.go`
- Create: `server/core/inventory/provider.go`
- Create: `server/core/inventory/router.go`
- Create: `server/core/inventory/router_test.go`
- Modify: `server/config/config.go`
- Modify: `server/core/app/app.go`
- Modify: `server/config.example.yaml`
- Modify: `config.example.yaml`

- [ ] **Step 1: Add the failing router/config test**

Create `server/core/inventory/router_test.go`:

```go
package inventory

import (
	"context"
	"testing"

	"github.com/ijry/lyshop/config"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type stubProvider struct{ name string }

func (s *stubProvider) Name() string { return s.name }
func (s *stubProvider) ReserveTx(_ *gorm.DB, _ ReserveInput) error { return nil }
func (s *stubProvider) ConfirmTx(_ *gorm.DB, _, _ string) error { return nil }
func (s *stubProvider) ReleaseTx(_ *gorm.DB, _, _, _ string) error { return nil }
func (s *stubProvider) DeductTx(_ *gorm.DB, _ DeductInput) error { return nil }
func (s *stubProvider) RestoreTx(_ *gorm.DB, _ RestoreInput) error { return nil }
func (s *stubProvider) SyncSkuTx(_ *gorm.DB, _ SyncSkuInput) error { return nil }
func (s *stubProvider) GetSellableStock(_ context.Context, _ []uint64) ([]SellableStock, error) { return nil, nil }

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
```

- [ ] **Step 2: Run the test to verify it fails**

Run:

```powershell
cd server
go test ./core/inventory -run "TestCurrentProviderUsesInventoryConfig|TestValidateConfigRejectsBuiltinWMSWithoutPlugin" -count=1
```

Expected: FAIL because `CurrentProvider`, `ValidateConfig`, `InventoryConfig`, and registry helpers do not exist.

- [ ] **Step 3: Add config structs and examples**

In `server/config/config.go`, add:

```go
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	Plugins     PluginsConfig     `mapstructure:"plugins"`
	Inventory   InventoryConfig   `mapstructure:"inventory"`
	ExternalWMS ExternalWMSConfig `mapstructure:"external_wms"`
}

type InventoryConfig struct {
	Provider     string `mapstructure:"provider"`
	ExternalMode string `mapstructure:"external_mode"`
}

type ExternalWMSConfig struct {
	Endpoint       string               `mapstructure:"endpoint"`
	AppKey         string               `mapstructure:"app_key"`
	AppSecret      string               `mapstructure:"app_secret"`
	TimeoutMS      int                  `mapstructure:"timeout_ms"`
	CallbackEnabled bool                `mapstructure:"callback_enabled"`
	Retry          ExternalWMSRetryConfig `mapstructure:"retry"`
}

type ExternalWMSRetryConfig struct {
	MaxAttempts    int `mapstructure:"max_attempts"`
	BackoffSeconds int `mapstructure:"backoff_seconds"`
}
```

In `config.example.yaml`, add:

```yaml
inventory:
  provider: local           # local | builtin_wms | external_wms
  external_mode: sync       # sync | async

external_wms:
  endpoint: ""
  app_key: ""
  app_secret: ""
  timeout_ms: 3000
  callback_enabled: true
  retry:
    max_attempts: 8
    backoff_seconds: 30
```

Also remove `- wms` from the root `config.example.yaml` default `plugins.enabled` list so the example matches the new default lightweight mode. Keep `server/config.example.yaml` with `wms` still listed only if its comments explicitly describe it as an optional built-in provider.

- [ ] **Step 4: Implement the provider registry and config validation**

Create `server/core/inventory/provider.go`:

```go
package inventory

import (
	"context"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type Provider interface {
	Name() string
	ReserveTx(tx *gorm.DB, in ReserveInput) error
	ConfirmTx(tx *gorm.DB, bizType, bizNo string) error
	ReleaseTx(tx *gorm.DB, bizType, bizNo, reason string) error
	DeductTx(tx *gorm.DB, in DeductInput) error
	RestoreTx(tx *gorm.DB, in RestoreInput) error
	SyncSkuTx(tx *gorm.DB, in SyncSkuInput) error
	GetSellableStock(ctx context.Context, skuIDs []uint64) ([]SellableStock, error)
}

var (
	regMu     sync.RWMutex
	providers = map[string]Provider{}
)

func Register(p Provider) {
	if p == nil {
		return
	}
	regMu.Lock()
	defer regMu.Unlock()
	providers[strings.ToLower(strings.TrimSpace(p.Name()))] = p
}

func Find(name string) Provider {
	regMu.RLock()
	defer regMu.RUnlock()
	return providers[strings.ToLower(strings.TrimSpace(name))]
}

func ResetRegistryForTest() {
	regMu.Lock()
	defer regMu.Unlock()
	providers = map[string]Provider{}
}
```

Create `server/core/inventory/router.go`:

```go
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
	enabled := map[string]struct{}{}
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
```

Create `server/core/inventory/types.go` and `errors.go` with the shared structs and typed errors referenced later:

```go
package inventory

import "time"

type ReserveItem struct {
	SkuID uint64
	Qty   int
}

type ReserveInput struct {
	BizType   string
	BizNo     string
	Items     []ReserveItem
	ExpiredAt *time.Time
}

type DeductInput struct {
	BizType string
	BizNo   string
	Items   []ReserveItem
}

type RestoreInput struct {
	BizType string
	BizNo   string
	Items   []ReserveItem
	Reason  string
}

type SyncSkuInput struct {
	SkuID  uint64
	Stock  int
	Source string
}

type SellableStock struct {
	SkuID     uint64 `json:"sku_id"`
	Sellable  int    `json:"sellable"`
	Reserved  int    `json:"reserved"`
	OnHand    int    `json:"on_hand"`
}
```

In `server/core/app/app.go`, call `inventory.ValidateConfig()` after `config.Load()` and before `plugin.Load(...)`.

- [ ] **Step 5: Run the tests, format, and commit**

Run:

```powershell
cd server
gofmt -w config/config.go core/inventory/types.go core/inventory/errors.go core/inventory/provider.go core/inventory/router.go core/inventory/router_test.go core/app/app.go
go test ./core/inventory -run "TestCurrentProviderUsesInventoryConfig|TestValidateConfigRejectsBuiltinWMSWithoutPlugin" -count=1
git add config/config.go config.example.yaml ../config.example.yaml core/inventory/types.go core/inventory/errors.go core/inventory/provider.go core/inventory/router.go core/inventory/router_test.go core/app/app.go
git commit -m "新增统一库存配置与路由骨架" -m "增加 inventory 和 external_wms 配置结构、统一 provider 注册与选择逻辑，并在应用启动阶段校验库存模式配置。"
```

Expected: tests pass and config examples show `local` as the default provider.

### Task 2: Shared Inventory State Models And Local Provider

**Files:**
- Create: `server/core/inventory/model.go`
- Create: `server/core/inventory/local.go`
- Create: `server/core/inventory/local_test.go`
- Modify: `server/core/app/app.go`
- Modify: `server/plugins/order/model/order.go`

- [ ] **Step 1: Add the failing local reservation tests**

Create `server/core/inventory/local_test.go`:

```go
package inventory

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/core/db"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestLocalProviderReserveConfirmRelease(t *testing.T) {
	testDB := setupLocalInventoryDB(t)
	p := &localProvider{}

	require.NoError(t, testDB.Create(&productmodel.ProductSku{
		ProductID: 1,
		SkuKey:    "red:l",
		Price:     10,
		Stock:     20,
		Status:    productmodel.ProductSkuStatusActive,
	}).Error)

	var sku productmodel.ProductSku
	require.NoError(t, testDB.First(&sku).Error)

	expireAt := time.Now().Add(15 * time.Minute)
	err := testDB.Transaction(func(tx *gorm.DB) error {
		return p.ReserveTx(tx, ReserveInput{
			BizType:   "order",
			BizNo:     "ORD-1",
			Items:     []ReserveItem{{SkuID: sku.ID, Qty: 3}},
			ExpiredAt: &expireAt,
		})
	})
	require.NoError(t, err)

	err = testDB.Transaction(func(tx *gorm.DB) error {
		return p.ConfirmTx(tx, "order", "ORD-1")
	})
	require.NoError(t, err)

	require.NoError(t, testDB.First(&sku, sku.ID).Error)
	require.Equal(t, 17, sku.Stock)

	err = testDB.Transaction(func(tx *gorm.DB) error {
		return p.ReleaseTx(tx, "order", "ORD-1", "already_confirmed")
	})
	require.ErrorContains(t, err, "confirmed")
}

func TestLocalProviderReserveRejectsInsufficientStock(t *testing.T) {
	testDB := setupLocalInventoryDB(t)
	p := &localProvider{}
	require.NoError(t, testDB.Create(&productmodel.ProductSku{
		ProductID: 2,
		SkuKey:    "blue:m",
		Price:     10,
		Stock:     1,
		Status:    productmodel.ProductSkuStatusActive,
	}).Error)
	var sku productmodel.ProductSku
	require.NoError(t, testDB.First(&sku).Error)
	err := testDB.Transaction(func(tx *gorm.DB) error {
		return p.ReserveTx(tx, ReserveInput{
			BizType: "order",
			BizNo:   "ORD-2",
			Items:   []ReserveItem{{SkuID: sku.ID, Qty: 2}},
		})
	})
	require.ErrorContains(t, err, "库存不足")
}

func setupLocalInventoryDB(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:local_inventory_%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	require.NoError(t, err)
	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })
	require.NoError(t, gdb.AutoMigrate(&productmodel.ProductSku{}, &InventoryReservation{}, &OrderInventoryState{}))
	return gdb
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run:

```powershell
cd server
go test ./core/inventory -run "TestLocalProviderReserveConfirmRelease|TestLocalProviderReserveRejectsInsufficientStock" -count=1
```

Expected: FAIL because `localProvider`, `InventoryReservation`, and `OrderInventoryState` do not exist.

- [ ] **Step 3: Add the shared inventory models**

Create `server/core/inventory/model.go`:

```go
package inventory

import (
	"time"

	"github.com/ijry/lyshop/model"
)

const (
	InventoryStatusNone      = "none"
	InventoryStatusPending   = "pending"
	InventoryStatusReserved  = "reserved"
	InventoryStatusConfirmed = "confirmed"
	InventoryStatusReleased  = "released"
	InventoryStatusFailed    = "failed"
)

type InventoryReservation struct {
	model.Base
	BizType   string     `gorm:"size:32;not null;index:idx_inventory_reservation_biz,priority:1" json:"biz_type"`
	BizNo     string     `gorm:"size:64;not null;index:idx_inventory_reservation_biz,priority:2" json:"biz_no"`
	SkuID     uint64     `gorm:"not null;index:idx_inventory_reservation_sku" json:"sku_id"`
	Qty       int        `gorm:"not null" json:"qty"`
	Status    string     `gorm:"size:16;not null;default:'reserved';index" json:"status"`
	Reason    string     `gorm:"size:128" json:"reason"`
	ExpiredAt *time.Time `json:"expired_at"`
}

type OrderInventoryState struct {
	model.Base
	OrderNo     string `gorm:"size:64;uniqueIndex;not null" json:"order_no"`
	BizType     string `gorm:"size:32;not null;default:'order'" json:"biz_type"`
	Status      string `gorm:"size:16;not null;default:'none';index" json:"status"`
	Provider    string `gorm:"size:32;not null;default:'local'" json:"provider"`
	LastError   string `gorm:"size:255" json:"last_error"`
}

type InventoryIntegrationTask struct {
	model.Base
	Provider     string `gorm:"size:32;not null;index" json:"provider"`
	BizType      string `gorm:"size:32;not null;index" json:"biz_type"`
	BizNo        string `gorm:"size:64;not null;index" json:"biz_no"`
	Action       string `gorm:"size:16;not null;index" json:"action"`
	Payload      string `gorm:"type:json" json:"payload"`
	Status       string `gorm:"size:16;not null;default:'pending';index" json:"status"`
	AttemptCount int    `gorm:"not null;default:0" json:"attempt_count"`
	NextRetryAt  *time.Time `json:"next_retry_at"`
	LastError    string `gorm:"size:255" json:"last_error"`
}
```

Add `InventoryStatus string 'gorm:"size:16;not null;default:'none';index" json:"inventory_status"'` to `server/plugins/order/model/order.go`.

In `server/core/app/app.go`, add `&inventory.InventoryReservation{}`, `&inventory.OrderInventoryState{}`, and `&inventory.InventoryIntegrationTask{}` to the core `AutoMigrate` call.

- [ ] **Step 4: Implement the local provider**

Create `server/core/inventory/local.go`:

```go
package inventory

import (
	"context"
	"fmt"

	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type localProvider struct{}

func init() {
	Register(&localProvider{})
}

func (p *localProvider) Name() string { return "local" }

func (p *localProvider) ReserveTx(tx *gorm.DB, in ReserveInput) error {
	for _, item := range in.Items {
		var sku productmodel.ProductSku
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&sku, item.SkuID).Error; err != nil {
			return err
		}
		if sku.Status != productmodel.ProductSkuStatusActive {
			return fmt.Errorf("SKU已下线")
		}
		if sku.Stock < item.Qty {
			return fmt.Errorf("库存不足")
		}
		if err := tx.Create(&InventoryReservation{
			BizType:   in.BizType,
			BizNo:     in.BizNo,
			SkuID:     item.SkuID,
			Qty:       item.Qty,
			Status:    InventoryStatusReserved,
			ExpiredAt: in.ExpiredAt,
		}).Error; err != nil {
			return err
		}
	}
	return tx.Where("order_no = ?", in.BizNo).
		Assign(&OrderInventoryState{BizType: in.BizType, OrderNo: in.BizNo, Status: InventoryStatusReserved, Provider: p.Name()}).
		FirstOrCreate(&OrderInventoryState{}).Error
}

func (p *localProvider) ConfirmTx(tx *gorm.DB, bizType, bizNo string) error {
	var rows []InventoryReservation
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("biz_type = ? AND biz_no = ?", bizType, bizNo).Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if row.Status == InventoryStatusConfirmed {
			continue
		}
		if row.Status != InventoryStatusReserved {
			return fmt.Errorf("inventory reservation status invalid: %s", row.Status)
		}
		if err := tx.Model(&productmodel.ProductSku{}).
			Where("id = ? AND stock >= ?", row.SkuID, row.Qty).
			Update("stock", gorm.Expr("stock - ?", row.Qty)).Error; err != nil {
			return err
		}
		if err := tx.Model(&InventoryReservation{}).Where("id = ?", row.ID).Update("status", InventoryStatusConfirmed).Error; err != nil {
			return err
		}
	}
	return tx.Model(&OrderInventoryState{}).Where("order_no = ?", bizNo).Updates(map[string]any{"status": InventoryStatusConfirmed, "last_error": ""}).Error
}

func (p *localProvider) ReleaseTx(tx *gorm.DB, bizType, bizNo, reason string) error {
	var rows []InventoryReservation
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("biz_type = ? AND biz_no = ?", bizType, bizNo).Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if row.Status == InventoryStatusReleased {
			continue
		}
		if row.Status == InventoryStatusConfirmed {
			return fmt.Errorf("inventory reservation already confirmed")
		}
		if err := tx.Model(&InventoryReservation{}).Where("id = ?", row.ID).Updates(map[string]any{"status": InventoryStatusReleased, "reason": reason}).Error; err != nil {
			return err
		}
	}
	return tx.Model(&OrderInventoryState{}).Where("order_no = ?", bizNo).Updates(map[string]any{"status": InventoryStatusReleased, "last_error": reason}).Error
}

func (p *localProvider) DeductTx(tx *gorm.DB, in DeductInput) error {
	for _, item := range in.Items {
		if err := tx.Model(&productmodel.ProductSku{}).Where("id = ? AND stock >= ?", item.SkuID, item.Qty).Update("stock", gorm.Expr("stock - ?", item.Qty)).Error; err != nil {
			return err
		}
	}
	return nil
}

func (p *localProvider) RestoreTx(tx *gorm.DB, in RestoreInput) error {
	for _, item := range in.Items {
		if err := tx.Model(&productmodel.ProductSku{}).Where("id = ?", item.SkuID).Update("stock", gorm.Expr("stock + ?", item.Qty)).Error; err != nil {
			return err
		}
	}
	return nil
}

func (p *localProvider) SyncSkuTx(_ *gorm.DB, _ SyncSkuInput) error { return nil }

func (p *localProvider) GetSellableStock(ctx context.Context, skuIDs []uint64) ([]SellableStock, error) {
	var rows []productmodel.ProductSku
	if err := db.DB.WithContext(ctx).Where("id IN ?", skuIDs).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]SellableStock, 0, len(rows))
	for _, row := range rows {
		out = append(out, SellableStock{SkuID: row.ID, Sellable: row.Stock, OnHand: row.Stock, Reserved: 0})
	}
	return out, nil
}
```

- [ ] **Step 5: Run the tests, format, and commit**

Run:

```powershell
cd server
gofmt -w core/inventory/model.go core/inventory/local.go core/inventory/local_test.go core/app/app.go plugins/order/model/order.go
go test ./core/inventory -run "TestLocalProviderReserveConfirmRelease|TestLocalProviderReserveRejectsInsufficientStock" -count=1
git add core/inventory/model.go core/inventory/local.go core/inventory/local_test.go core/app/app.go plugins/order/model/order.go
git commit -m "新增本地库存实现与共享状态模型" -m "增加库存预占共享模型和 local provider，为无 WMS 模式提供订单库存预占、确认、释放与库存状态持久化。"
```

Expected: local provider tests pass and the `orders` table has the new `inventory_status` field.

### Task 3: Refactor Order To Use The Inventory Router

**Files:**
- Modify: `server/plugins/order/plugin.json`
- Modify: `server/plugins/order/service/order.go`
- Create: `server/plugins/order/service/order_inventory_test.go`
- Modify: `server/core/plugin/plugin_test.go`

- [ ] **Step 1: Add the failing order inventory test**

Create `server/plugins/order/service/order_inventory_test.go`:

```go
package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spyInventoryProvider struct {
	reserveCalls int
	confirmCalls int
	releaseCalls int
	lastReserve  inventory.ReserveInput
}

func (s *spyInventoryProvider) Name() string { return "spy" }
func (s *spyInventoryProvider) ReserveTx(_ *gorm.DB, in inventory.ReserveInput) error {
	s.reserveCalls++
	s.lastReserve = in
	return nil
}
func (s *spyInventoryProvider) ConfirmTx(_ *gorm.DB, _, _ string) error { s.confirmCalls++; return nil }
func (s *spyInventoryProvider) ReleaseTx(_ *gorm.DB, _, _, _ string) error { s.releaseCalls++; return nil }
func (s *spyInventoryProvider) DeductTx(_ *gorm.DB, _ inventory.DeductInput) error { return nil }
func (s *spyInventoryProvider) RestoreTx(_ *gorm.DB, _ inventory.RestoreInput) error { return nil }
func (s *spyInventoryProvider) SyncSkuTx(_ *gorm.DB, _ inventory.SyncSkuInput) error { return nil }
func (s *spyInventoryProvider) GetSellableStock(_ context.Context, _ []uint64) ([]inventory.SellableStock, error) { return nil, nil }

func TestCreateOrderUsesInventoryProviderReserve(t *testing.T) {
	original := getInventoryProviderFn
	spy := &spyInventoryProvider{}
	getInventoryProviderFn = func() (inventory.Provider, error) { return spy, nil }
	t.Cleanup(func() { getInventoryProviderFn = original })

	items := []ordermodel.OrderItem{{SkuID: 101, Qty: 2}}
	err := reserveOrderInventory(nil, "ORD-TEST", items)
	require.NoError(t, err)
	require.Equal(t, 1, spy.reserveCalls)
	require.Equal(t, "ORD-TEST", spy.lastReserve.BizNo)
	require.Len(t, spy.lastReserve.Items, 1)
	require.Equal(t, 2, spy.lastReserve.Items[0].Qty)
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run:

```powershell
cd server
go test ./plugins/order/service -run TestCreateOrderUsesInventoryProviderReserve -count=1
```

Expected: FAIL because `getInventoryProviderFn` and `reserveOrderInventory` do not exist.

- [ ] **Step 3: Remove the direct WMS dependency and add order inventory helpers**

In `server/plugins/order/plugin.json`, change:

```json
"depends": ["product"]
```

In `server/plugins/order/service/order.go`, replace:

```go
wmssvc "github.com/ijry/lyshop/plugins/wms/service"
var reserveStockTxFn = wmssvc.ReserveStockTx
var confirmReservationTxFn = wmssvc.ConfirmReservationTx
var releaseReservationTxFn = wmssvc.ReleaseReservationTx
var pickDefaultWarehouseIDTxFn = wmssvc.PickDefaultWarehouseIDTx
```

with:

```go
inventorycore "github.com/ijry/lyshop/core/inventory"

var getInventoryProviderFn = inventorycore.CurrentProvider

func reserveOrderInventory(tx *gorm.DB, orderNo string, items []ordermodel.OrderItem) error {
	provider, err := getInventoryProviderFn()
	if err != nil {
		return err
	}
	reserveItems := make([]inventorycore.ReserveItem, 0, len(items))
	for _, item := range items {
		reserveItems = append(reserveItems, inventorycore.ReserveItem{SkuID: item.SkuID, Qty: item.Qty})
	}
	expireAt := time.Now().Add(orderReserveExpire)
	return provider.ReserveTx(tx, inventorycore.ReserveInput{
		BizType:   orderReservationBizType,
		BizNo:     orderNo,
		Items:     reserveItems,
		ExpiredAt: &expireAt,
	})
}

func confirmOrderInventory(tx *gorm.DB, orderNo string) error {
	provider, err := getInventoryProviderFn()
	if err != nil {
		return err
	}
	return provider.ConfirmTx(tx, orderReservationBizType, orderNo)
}

func releaseOrderInventory(tx *gorm.DB, orderNo, reason string) error {
	provider, err := getInventoryProviderFn()
	if err != nil {
		return err
	}
	return provider.ReleaseTx(tx, orderReservationBizType, orderNo, reason)
}
```

Use those helpers in `CreateOrder`, `PayOrder`, and `CancelOrder`. Set:

```go
order.InventoryStatus = inventorycore.InventoryStatusReserved
```

when the order is successfully created, and update status to `confirmed`/`released` in `PayOrder` and `CancelOrder`.

- [ ] **Step 4: Update plugin tests and add a config-loading regression**

In `server/core/plugin/plugin_test.go`, change the stub order dependency from `[]string{"product"}` in the success path and add:

```go
func TestLoad_OrderDoesNotRequireWMS(t *testing.T) {
	resetRegistry()
	Register(&stub{name: "product"})
	Register(&stub{name: "order", deps: []string{"product"}})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	err := Load([]string{"product", "order"}, &gorm.DB{}, r.Group("/api/v1"), r.Group("/admin/api"))
	require.NoError(t, err)
}
```

- [ ] **Step 5: Run tests, format, and commit**

Run:

```powershell
cd server
gofmt -w plugins/order/service/order.go plugins/order/service/order_inventory_test.go core/plugin/plugin_test.go
go test ./plugins/order/service -run TestCreateOrderUsesInventoryProviderReserve -count=1
go test ./core/plugin -run TestLoad_OrderDoesNotRequireWMS -count=1
git add plugins/order/plugin.json plugins/order/service/order.go plugins/order/service/order_inventory_test.go core/plugin/plugin_test.go
git commit -m "解除订单对内置WMS的强依赖" -m "订单插件改为依赖统一 inventory 路由，移除对 wms 插件的硬依赖并通过库存状态字段记录预占链路。"
```

Expected: the order service no longer imports `plugins/wms/service`.

### Task 4: Refactor Product And Points Mall To Use Inventory

**Files:**
- Modify: `server/plugins/product/service/product.go`
- Create: `server/plugins/product/service/product_inventory_test.go`
- Modify: `server/plugins/points_mall/service/exchange.go`
- Create: `server/plugins/points_mall/service/exchange_inventory_test.go`

- [ ] **Step 1: Add the failing SKU sync and points deduction tests**

Create `server/plugins/product/service/product_inventory_test.go`:

```go
package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spySkuSyncProvider struct {
	syncCalls int
	lastInput inventory.SyncSkuInput
}

func (s *spySkuSyncProvider) Name() string { return "spy" }
func (s *spySkuSyncProvider) ReserveTx(_ *gorm.DB, _ inventory.ReserveInput) error { return nil }
func (s *spySkuSyncProvider) ConfirmTx(_ *gorm.DB, _, _ string) error { return nil }
func (s *spySkuSyncProvider) ReleaseTx(_ *gorm.DB, _, _, _ string) error { return nil }
func (s *spySkuSyncProvider) DeductTx(_ *gorm.DB, _ inventory.DeductInput) error { return nil }
func (s *spySkuSyncProvider) RestoreTx(_ *gorm.DB, _ inventory.RestoreInput) error { return nil }
func (s *spySkuSyncProvider) SyncSkuTx(_ *gorm.DB, in inventory.SyncSkuInput) error { s.syncCalls++; s.lastInput = in; return nil }
func (s *spySkuSyncProvider) GetSellableStock(_ context.Context, _ []uint64) ([]inventory.SellableStock, error) { return nil, nil }

func TestSyncInventoryForNewSkusUsesInventoryProvider(t *testing.T) {
	original := getInventoryProviderForProductFn
	spy := &spySkuSyncProvider{}
	getInventoryProviderForProductFn = func() (inventory.Provider, error) { return spy, nil }
	t.Cleanup(func() { getInventoryProviderForProductFn = original })

	err := syncInventoryForNewSkus(nil, []productmodel.ProductSku{{Base: model.Base{ID: 9}, Stock: 12}})
	require.NoError(t, err)
	require.Equal(t, 1, spy.syncCalls)
	require.Equal(t, uint64(9), spy.lastInput.SkuID)
	require.Equal(t, 12, spy.lastInput.Stock)
}
```

Create `server/plugins/points_mall/service/exchange_inventory_test.go`:

```go
package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type spyPointsInventoryProvider struct {
	deductCalls  int
	restoreCalls int
}

func (s *spyPointsInventoryProvider) Name() string { return "spy" }
func (s *spyPointsInventoryProvider) ReserveTx(_ *gorm.DB, _ inventory.ReserveInput) error { return nil }
func (s *spyPointsInventoryProvider) ConfirmTx(_ *gorm.DB, _, _ string) error { return nil }
func (s *spyPointsInventoryProvider) ReleaseTx(_ *gorm.DB, _, _, _ string) error { return nil }
func (s *spyPointsInventoryProvider) DeductTx(_ *gorm.DB, _ inventory.DeductInput) error { s.deductCalls++; return nil }
func (s *spyPointsInventoryProvider) RestoreTx(_ *gorm.DB, _ inventory.RestoreInput) error { s.restoreCalls++; return nil }
func (s *spyPointsInventoryProvider) SyncSkuTx(_ *gorm.DB, _ inventory.SyncSkuInput) error { return nil }
func (s *spyPointsInventoryProvider) GetSellableStock(_ context.Context, _ []uint64) ([]inventory.SellableStock, error) { return nil, nil }

func TestPointsMallUsesInventoryProviderHooks(t *testing.T) {
	original := getInventoryProviderForPointsMallFn
	spy := &spyPointsInventoryProvider{}
	getInventoryProviderForPointsMallFn = func() (inventory.Provider, error) { return spy, nil }
	t.Cleanup(func() { getInventoryProviderForPointsMallFn = original })

	err := deductPointsMallInventory(nil, 1001, 2)
	require.NoError(t, err)
	require.Equal(t, 1, spy.deductCalls)

	err = restorePointsMallInventory(nil, 1001, 2, "cancel")
	require.NoError(t, err)
	require.Equal(t, 1, spy.restoreCalls)
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run:

```powershell
cd server
go test ./plugins/product/service -run TestSyncInventoryForNewSkusUsesInventoryProvider -count=1
go test ./plugins/points_mall/service -run TestPointsMallUsesInventoryProviderHooks -count=1
```

Expected: FAIL because `syncInventoryForNewSkus`, `deductPointsMallInventory`, and provider accessors do not exist.

- [ ] **Step 3: Refactor product SKU sync**

In `server/plugins/product/service/product.go`, replace `syncWmsStockForNewSkus` with:

```go
var getInventoryProviderForProductFn = inventorycore.CurrentProvider

func syncInventoryForNewSkus(tx *gorm.DB, skus []productmodel.ProductSku) error {
	if len(skus) == 0 {
		return nil
	}
	provider, err := getInventoryProviderForProductFn()
	if err != nil {
		return err
	}
	for _, sku := range skus {
		if sku.ID == 0 {
			continue
		}
		if err := provider.SyncSkuTx(tx, inventorycore.SyncSkuInput{
			SkuID:  sku.ID,
			Stock:  sku.Stock,
			Source: "product",
		}); err != nil {
			return err
		}
	}
	return nil
}
```

Update both call sites in `CreateProduct` and `ReplaceProductSkus` to call `syncInventoryForNewSkus`.

- [ ] **Step 4: Refactor points mall stock deduction/restore**

In `server/plugins/points_mall/service/exchange.go`, add:

```go
var getInventoryProviderForPointsMallFn = inventorycore.CurrentProvider

func deductPointsMallInventory(tx *gorm.DB, productID uint64, qty int) error {
	provider, err := getInventoryProviderForPointsMallFn()
	if err != nil {
		return err
	}
	return provider.DeductTx(tx, inventorycore.DeductInput{
		BizType: "points_mall",
		BizNo:   fmt.Sprintf("points-product:%d", productID),
		Items:   []inventorycore.ReserveItem{{SkuID: productID, Qty: qty}},
	})
}

func restorePointsMallInventory(tx *gorm.DB, productID uint64, qty int, reason string) error {
	provider, err := getInventoryProviderForPointsMallFn()
	if err != nil {
		return err
	}
	return provider.RestoreTx(tx, inventorycore.RestoreInput{
		BizType: "points_mall",
		BizNo:   fmt.Sprintf("points-product:%d", productID),
		Items:   []inventorycore.ReserveItem{{SkuID: productID, Qty: qty}},
		Reason:  reason,
	})
}
```

Use those helpers in `ExchangeProduct` and `CancelExchange`. Keep the existing `PointsProduct.stock` update for this phase only when the selected provider is `local`, so legacy points products continue to work without a dedicated SKU mapping:

```go
provider, err := getInventoryProviderForPointsMallFn()
if err != nil {
	return err
}
if provider.Name() == "local" {
	// retain existing points_product stock updates for local mode
} else {
	if err := deductPointsMallInventory(tx, product.ID, qty); err != nil {
		return err
	}
}
```

This keeps the first implementation scoped and avoids inventing a new points-product-to-sku migration in the same change set.

- [ ] **Step 5: Run tests, format, and commit**

Run:

```powershell
cd server
gofmt -w plugins/product/service/product.go plugins/product/service/product_inventory_test.go plugins/points_mall/service/exchange.go plugins/points_mall/service/exchange_inventory_test.go
go test ./plugins/product/service -run TestSyncInventoryForNewSkusUsesInventoryProvider -count=1
go test ./plugins/points_mall/service -run TestPointsMallUsesInventoryProviderHooks -count=1
git add plugins/product/service/product.go plugins/product/service/product_inventory_test.go plugins/points_mall/service/exchange.go plugins/points_mall/service/exchange_inventory_test.go
git commit -m "统一商品与积分商城库存入口" -m "商品 SKU 初始化和积分商城库存扣减改为走统一 inventory 接口，避免继续直接绑定 WMS 逻辑。"
```

Expected: product and points mall no longer import `plugins/wms/service`.

### Task 5: Wrap Builtin WMS As An Inventory Provider

**Files:**
- Modify: `server/plugins/wms/plugin.go`
- Create: `server/plugins/wms/provider.go`
- Create: `server/plugins/wms/provider_test.go`

- [ ] **Step 1: Add the failing provider adapter test**

Create `server/plugins/wms/provider_test.go`:

```go
package wms

import (
	"context"
	"testing"

	inventorycore "github.com/ijry/lyshop/core/inventory"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestBuiltinWMSProviderName(t *testing.T) {
	p := &builtinProvider{}
	require.Equal(t, "builtin_wms", p.Name())
}

func TestBuiltinWMSProviderReserveDelegatesToWMSService(t *testing.T) {
	p := &builtinProvider{}
	original := reserveStockTxFn
	called := false
	reserveStockTxFn = func(_ *gorm.DB, in ReserveStockInput) error {
		called = true
		require.Equal(t, "order", in.BizType)
		require.Equal(t, "ORD-3", in.BizNo)
		require.Len(t, in.Items, 1)
		return nil
	}
	t.Cleanup(func() { reserveStockTxFn = original })

	err := p.ReserveTx(nil, inventorycore.ReserveInput{
		BizType: "order",
		BizNo:   "ORD-3",
		Items:   []inventorycore.ReserveItem{{SkuID: 1, Qty: 2}},
	})
	require.NoError(t, err)
	require.True(t, called)
}

func TestBuiltinWMSProviderGetSellableStock(t *testing.T) {
	p := &builtinProvider{}
	original := listStocksBySkuIDsFn
	listStocksBySkuIDsFn = func(_ context.Context, skuIDs []uint64) ([]wmsmodel.InventoryStock, error) {
		require.Equal(t, []uint64{7}, skuIDs)
		return []wmsmodel.InventoryStock{{SkuID: 7, Qty: 10, ReservedQty: 3}}, nil
	}
	t.Cleanup(func() { listStocksBySkuIDsFn = original })

	list, err := p.GetSellableStock(context.Background(), []uint64{7})
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, 7, list[0].Sellable)
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run:

```powershell
cd server
go test ./plugins/wms -run "TestBuiltinWMSProviderName|TestBuiltinWMSProviderReserveDelegatesToWMSService|TestBuiltinWMSProviderGetSellableStock" -count=1
```

Expected: FAIL because `builtinProvider`, `reserveStockTxFn`, and `listStocksBySkuIDsFn` do not exist in the plugin package.

- [ ] **Step 3: Implement the WMS provider wrapper**

Create `server/plugins/wms/provider.go`:

```go
package wms

import (
	"context"

	inventorycore "github.com/ijry/lyshop/core/inventory"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	wmssvc "github.com/ijry/lyshop/plugins/wms/service"
	"gorm.io/gorm"
)

var (
	reserveStockTxFn     = wmssvc.ReserveStockTx
	confirmReservationFn = wmssvc.ConfirmReservationTx
	releaseReservationFn = wmssvc.ReleaseReservationTx
	listStocksBySkuIDsFn = wmssvc.ListStocksBySkuIDs
)

type builtinProvider struct{}

func (p *builtinProvider) Name() string { return "builtin_wms" }

func (p *builtinProvider) ReserveTx(tx *gorm.DB, in inventorycore.ReserveInput) error {
	items := make([]wmssvc.ReservationItemInput, 0, len(in.Items))
	for _, item := range in.Items {
		items = append(items, wmssvc.ReservationItemInput{SkuID: item.SkuID, Qty: item.Qty})
	}
	warehouseID, err := wmssvc.PickDefaultWarehouseIDTx(tx)
	if err != nil {
		return err
	}
	return reserveStockTxFn(tx, wmssvc.ReserveStockInput{
		BizType:     in.BizType,
		BizNo:       in.BizNo,
		WarehouseID: warehouseID,
		Items:       items,
		ExpiredAt:   in.ExpiredAt,
	})
}

func (p *builtinProvider) ConfirmTx(tx *gorm.DB, bizType, bizNo string) error {
	return confirmReservationFn(tx, bizType, bizNo)
}

func (p *builtinProvider) ReleaseTx(tx *gorm.DB, bizType, bizNo, reason string) error {
	return releaseReservationFn(tx, bizType, bizNo, reason)
}

func (p *builtinProvider) DeductTx(tx *gorm.DB, in inventorycore.DeductInput) error { return nil }
func (p *builtinProvider) RestoreTx(tx *gorm.DB, in inventorycore.RestoreInput) error { return nil }

func (p *builtinProvider) SyncSkuTx(tx *gorm.DB, in inventorycore.SyncSkuInput) error {
	warehouseID, err := wmssvc.PickDefaultWarehouseIDTx(tx)
	if err != nil || warehouseID == 0 {
		return err
	}
	var stock wmsmodel.InventoryStock
	result := tx.Where("warehouse_id = ? AND sku_id = ?", warehouseID, in.SkuID).First(&stock)
	if result.Error == nil {
		return nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}
	return tx.Create(&wmsmodel.InventoryStock{WarehouseID: warehouseID, SkuID: in.SkuID, Qty: in.Stock}).Error
}

func (p *builtinProvider) GetSellableStock(ctx context.Context, skuIDs []uint64) ([]inventorycore.SellableStock, error) {
	rows, err := listStocksBySkuIDsFn(ctx, skuIDs)
	if err != nil {
		return nil, err
	}
	out := make([]inventorycore.SellableStock, 0, len(rows))
	for _, row := range rows {
		out = append(out, inventorycore.SellableStock{
			SkuID:    row.SkuID,
			OnHand:   row.Qty,
			Reserved: row.ReservedQty,
			Sellable: row.Qty - row.ReservedQty,
		})
	}
	return out, nil
}
```

In `server/plugins/wms/plugin.go`, register the provider inside `init()`:

```go
inventorycore.Register(&builtinProvider{})
```

- [ ] **Step 4: Extend the wrapper to support non-order deduct/restore**

For `DeductTx` and `RestoreTx`, keep the first implementation explicit and small:

```go
func (p *builtinProvider) DeductTx(tx *gorm.DB, in inventorycore.DeductInput) error {
	warehouseID, err := wmssvc.PickDefaultWarehouseIDTx(tx)
	if err != nil {
		return err
	}
	doc := wmssvc.CreateDocInput{
		WarehouseID: warehouseID,
		DocType:     wmsmodel.DocTypeOutbound,
		Remark:      in.BizNo,
		Items:       []wmssvc.DocItemInput{},
	}
	for _, item := range in.Items {
		doc.Items = append(doc.Items, wmssvc.DocItemInput{SkuID: item.SkuID, Qty: item.Qty})
	}
	row, err := wmssvc.CreateDraftDoc(context.Background(), doc)
	if err != nil {
		return err
	}
	return wmssvc.CompleteDraftDoc(context.Background(), row.ID)
}
```

Mirror it with `DocTypeInbound` in `RestoreTx`. This is intentionally conservative: it reuses the existing WMS transaction engine instead of inventing a new shortcut path.

- [ ] **Step 5: Run tests, format, and commit**

Run:

```powershell
cd server
gofmt -w plugins/wms/provider.go plugins/wms/provider_test.go plugins/wms/plugin.go
go test ./plugins/wms -run "TestBuiltinWMSProviderName|TestBuiltinWMSProviderReserveDelegatesToWMSService|TestBuiltinWMSProviderGetSellableStock" -count=1
git add plugins/wms/provider.go plugins/wms/provider_test.go plugins/wms/plugin.go
git commit -m "将内置WMS封装为库存Provider" -m "基于现有仓储服务封装 builtin_wms provider，让统一 inventory 路由可以复用预占、确认、释放与库存查询能力。"
```

Expected: built-in WMS is selectable through `inventory.provider=builtin_wms`.

### Task 6: External WMS Sync Provider, Async Tasks, And Admin APIs

**Files:**
- Create: `server/core/inventory/external.go`
- Create: `server/core/inventory/tasks.go`
- Create: `server/core/inventory/external_test.go`
- Create: `server/core/inventory/tasks_test.go`
- Create: `server/plugins/external_wms/plugin.go`
- Create: `server/plugins/external_wms/plugin.json`
- Create: `server/plugins/external_wms/api/admin.go`
- Create: `server/plugins/external_wms/api/admin_test.go`
- Modify: `server/main.go`

- [ ] **Step 1: Add the failing external provider tests**

Create `server/core/inventory/external_test.go`:

```go
package inventory

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/ijry/lyshop/config"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestExternalProviderReserveSyncCallsRemoteAPI(t *testing.T) {
	originalClient := externalHTTPClient
	externalHTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, "POST", req.Method)
		require.Equal(t, "/reserve", req.URL.Path)
		body, _ := io.ReadAll(req.Body)
		require.Contains(t, string(body), `"biz_no":"ORD-9"`)
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"code":0}`)),
			Header:     make(http.Header),
		}, nil
	})}
	t.Cleanup(func() { externalHTTPClient = originalClient })

	original := config.Global
	t.Cleanup(func() { config.Global = original })
	config.Global.Inventory.Provider = "external_wms"
	config.Global.Inventory.ExternalMode = "sync"
	config.Global.ExternalWMS.Endpoint = "https://wms.example.com"

	p := &externalProvider{}
	err := p.ReserveTx(&gorm.DB{}, ReserveInput{
		BizType: "order",
		BizNo:   "ORD-9",
		Items:   []ReserveItem{{SkuID: 1, Qty: 2}},
	})
	require.NoError(t, err)
}
```

Create `server/core/inventory/tasks_test.go`:

```go
package inventory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnqueueIntegrationTaskCreatesPendingTask(t *testing.T) {
	task := NewIntegrationTask("external_wms", "order", "ORD-10", "reserve", TaskPayload{Items: []ReserveItem{{SkuID: 1, Qty: 2}}})
	require.Equal(t, "pending", task.Status)
	require.Equal(t, "external_wms", task.Provider)
	require.Equal(t, "reserve", task.Action)
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run:

```powershell
cd server
go test ./core/inventory -run "TestExternalProviderReserveSyncCallsRemoteAPI|TestEnqueueIntegrationTaskCreatesPendingTask" -count=1
```

Expected: FAIL because `externalProvider`, `TaskPayload`, and `NewIntegrationTask` do not exist.

- [ ] **Step 3: Implement the external provider and async task helpers**

Create `server/core/inventory/external.go`:

```go
package inventory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ijry/lyshop/config"
	"gorm.io/gorm"
)

var externalHTTPClient = &http.Client{Timeout: 3 * time.Second}

type externalProvider struct{}

func init() {
	Register(&externalProvider{})
}

func (p *externalProvider) Name() string { return "external_wms" }

func (p *externalProvider) ReserveTx(tx *gorm.DB, in ReserveInput) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", in.BizType, in.BizNo, "reserve", TaskPayload{Items: in.Items})
	}
	return p.postJSON("/reserve", map[string]any{"biz_type": in.BizType, "biz_no": in.BizNo, "items": in.Items})
}

func (p *externalProvider) ConfirmTx(tx *gorm.DB, bizType, bizNo string) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", bizType, bizNo, "confirm", TaskPayload{})
	}
	return p.postJSON("/confirm", map[string]any{"biz_type": bizType, "biz_no": bizNo})
}

func (p *externalProvider) ReleaseTx(tx *gorm.DB, bizType, bizNo, reason string) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", bizType, bizNo, "release", TaskPayload{Reason: reason})
	}
	return p.postJSON("/release", map[string]any{"biz_type": bizType, "biz_no": bizNo, "reason": reason})
}

func (p *externalProvider) DeductTx(tx *gorm.DB, in DeductInput) error  { return EnqueueInventoryTask(tx, "external_wms", in.BizType, in.BizNo, "deduct", TaskPayload{Items: in.Items}) }
func (p *externalProvider) RestoreTx(tx *gorm.DB, in RestoreInput) error { return EnqueueInventoryTask(tx, "external_wms", in.BizType, in.BizNo, "restore", TaskPayload{Items: in.Items, Reason: in.Reason}) }
func (p *externalProvider) SyncSkuTx(tx *gorm.DB, in SyncSkuInput) error { return EnqueueInventoryTask(tx, "external_wms", in.Source, fmt.Sprintf("sku:%d", in.SkuID), "sync_sku", TaskPayload{SkuID: in.SkuID, Stock: in.Stock}) }
func (p *externalProvider) GetSellableStock(ctx context.Context, skuIDs []uint64) ([]SellableStock, error) { return []SellableStock{}, nil }

func (p *externalProvider) postJSON(path string, payload map[string]any) error {
	raw, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(config.Global.ExternalWMS.Endpoint, "/")+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := externalHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("external wms http status %d", resp.StatusCode)
	}
	return nil
}
```

Create `server/core/inventory/tasks.go`:

```go
package inventory

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type TaskPayload struct {
	SkuID  uint64        `json:"sku_id,omitempty"`
	Stock  int           `json:"stock,omitempty"`
	Items  []ReserveItem `json:"items,omitempty"`
	Reason string        `json:"reason,omitempty"`
}

func NewIntegrationTask(provider, bizType, bizNo, action string, payload TaskPayload) *InventoryIntegrationTask {
	raw, _ := json.Marshal(payload)
	now := time.Now()
	return &InventoryIntegrationTask{
		Provider:    provider,
		BizType:     bizType,
		BizNo:       bizNo,
		Action:      action,
		Payload:     string(raw),
		Status:      "pending",
		AttemptCount: 0,
		NextRetryAt: &now,
	}
}

func EnqueueInventoryTask(tx *gorm.DB, provider, bizType, bizNo, action string, payload TaskPayload) error {
	return tx.Create(NewIntegrationTask(provider, bizType, bizNo, action, payload)).Error
}
```

- [ ] **Step 4: Add the external plugin and admin APIs**

Create `server/plugins/external_wms/plugin.go`:

```go
package external_wms

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/plugin"
	externalapi "github.com/ijry/lyshop/plugins/external_wms/api"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type externalPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("external_wms plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&externalPlugin{meta: m})
}

func (p *externalPlugin) Meta() plugin.Meta { return p.meta }
func (p *externalPlugin) RegisterRoutes(_ *gin.RouterGroup, admin *gin.RouterGroup) { externalapi.RegisterAdminRoutes(admin) }
func (p *externalPlugin) Migrate(_ *gorm.DB) error { return nil }
func (p *externalPlugin) Install() error { return nil }
func (p *externalPlugin) Uninstall() error { return nil }
```

Create `server/plugins/external_wms/api/admin.go`:

```go
package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	inventorycore "github.com/ijry/lyshop/core/inventory"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/inventory/tasks", middleware.RequirePermission("system:config"), listTasks)
	g.POST("/inventory/tasks/:id/retry", middleware.RequirePermission("system:config"), retryTask)
	g.POST("/external-wms/callback", callback)
}

func listTasks(c *gin.Context) {
	var rows []inventorycore.InventoryIntegrationTask
	if err := db.DB.WithContext(c.Request.Context()).Order("id desc").Limit(50).Find(&rows).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, rows)
}

func retryTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := db.DB.WithContext(c.Request.Context()).Model(&inventorycore.InventoryIntegrationTask{}).Where("id = ?", id).Updates(map[string]any{"status": "pending", "last_error": ""}).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, gin.H{"retried": id})
}

func callback(c *gin.Context) {
	response.OK(c, gin.H{"received": true})
}
```

Add `_ "github.com/ijry/lyshop/plugins/external_wms"` to `server/main.go`.

- [ ] **Step 5: Run tests, format, and commit**

Run:

```powershell
cd server
gofmt -w core/inventory/external.go core/inventory/tasks.go core/inventory/external_test.go core/inventory/tasks_test.go plugins/external_wms/plugin.go plugins/external_wms/api/admin.go main.go
go test ./core/inventory -run "TestExternalProviderReserveSyncCallsRemoteAPI|TestEnqueueIntegrationTaskCreatesPendingTask" -count=1
go test ./plugins/external_wms/... -count=1
git add core/inventory/external.go core/inventory/tasks.go core/inventory/external_test.go core/inventory/tasks_test.go plugins/external_wms/plugin.go plugins/external_wms/plugin.json plugins/external_wms/api/admin.go main.go
git commit -m "新增外部WMS同步与异步对接骨架" -m "增加 external_wms provider、异步任务模型和基础回调/重试接口，为企业 WMS 同步与异步集成提供统一入口。"
```

Expected: external inventory mode can reserve synchronously and enqueue async tasks without touching the built-in WMS plugin.

### Task 7: Documentation, Config Validation, And Full Regression

**Files:**
- Modify: `docs-site/docs/api/order.md`
- Modify: `docs-site/docs/api/stock-reservation.md`
- Modify: `docs-site/docs/api/wms.md`
- Modify: `docs-site/docs/guide/features.md`
- Modify: `server/config.example.yaml`
- Modify: `config.example.yaml`

- [ ] **Step 1: Update docs-site to the latest architecture**

Update the docs to describe the current latest model, not a migration guide:

In `docs-site/docs/api/order.md`, add:

```md
- 订单库存由统一 inventory provider 负责。
- `inventory_status` 取值为 `none`、`pending`、`reserved`、`confirmed`、`released`、`failed`。
- 只有 `inventory_reserved` 的订单允许继续支付链路。
```

In `docs-site/docs/api/stock-reservation.md`, rewrite the opening section to:

```md
当前库存架构支持三种 provider：

1. `local`：商品库存直接由 `product_skus.stock` 管理。
2. `builtin_wms`：库存交易由内置 `wms` 插件处理。
3. `external_wms`：库存交易委托企业外部 WMS，可按配置选择 `sync` 或 `async`。
```

In `docs-site/docs/api/wms.md`, add:

```md
`wms` 是商城的内置库存 provider 之一，只有在 `inventory.provider=builtin_wms` 时作为库存交易真源。
当商城运行于 `local` 或 `external_wms` 模式时，WMS 后台能力不是必选依赖。
```

In `docs-site/docs/guide/features.md`, update the feature list to state that WMS is optional and local mode is the default.

- [ ] **Step 2: Add a final config validation regression**

Extend `server/core/inventory/router_test.go` with:

```go
func TestValidateConfigRejectsExternalWMSWithoutEndpoint(t *testing.T) {
	original := config.Global
	t.Cleanup(func() { config.Global = original })
	config.Global.Inventory.Provider = "external_wms"
	config.Global.Inventory.ExternalMode = "sync"
	config.Global.ExternalWMS.Endpoint = ""

	err := ValidateConfig()
	require.ErrorContains(t, err, "external_wms.endpoint")
}
```

- [ ] **Step 3: Run full backend verification**

Run:

```powershell
cd server
gofmt -w core/inventory/router_test.go
go test ./core/inventory ./core/plugin ./plugins/order/service ./plugins/product/service ./plugins/points_mall/service ./plugins/wms ./plugins/external_wms/... -count=1
```

Expected: all inventory-related tests pass. If an unrelated pre-existing test fails, record the exact package and error before proceeding.

- [ ] **Step 4: Commit the docs and verification changes**

Run:

```powershell
git add docs-site/docs/api/order.md docs-site/docs/api/stock-reservation.md docs-site/docs/api/wms.md docs-site/docs/guide/features.md server/core/inventory/router_test.go server/config.example.yaml config.example.yaml
git commit -m "更新统一库存架构文档" -m "同步 docs-site 和配置示例，明确 local、builtin_wms、external_wms 三态架构以及外部 WMS 的同步异步模式。"
```

## Final Review Checklist

- [ ] `git log -5 --pretty=%B` contains only Chinese commit messages and no `Co-Authored-By`.
- [ ] `git status --short` is clean.
- [ ] `config.example.yaml` defaults to `inventory.provider=local`.
- [ ] `server/plugins/order/plugin.json` no longer depends on `wms`.
- [ ] `server/plugins/order/service/order.go` and `server/plugins/product/service/product.go` do not import `plugins/wms/service`.
- [ ] `docs-site` documents the latest three-mode inventory architecture rather than a WMS-only model.
- [ ] `external_wms` async mode exposes task inspection/retry APIs and shared task persistence.

## Self-Review

- Spec coverage:
  - `local` / `builtin_wms` / `external_wms` provider split: Tasks 1, 2, 5, 6.
  - `sync` / `async` external modes: Task 6.
  - `order` / `product` / `points_mall` decoupling from WMS: Tasks 3 and 4.
  - shared inventory status and config validation: Tasks 1 and 2.
  - docs-site synchronization: Task 7.
- Placeholder scan:
  - No `TODO`, `TBD`, or “implement later” markers remain in the plan.
  - Every code-touching step includes concrete file paths, code snippets, and verification commands.
- Type consistency:
  - The plan consistently uses `InventoryReservation`, `InventoryIntegrationTask`, `InventoryStatus*`, `ReserveInput`, `DeductInput`, `RestoreInput`, and `SyncSkuInput`.
