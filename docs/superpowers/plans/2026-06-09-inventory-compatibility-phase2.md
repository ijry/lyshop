# Inventory Compatibility Phase 2 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete the external WMS production path so LYShop can run the unified inventory architecture with a real enterprise external WMS in sync or async mode.

**Architecture:** Keep the existing `server/core/inventory` abstraction as the only inventory entry point, then extend the `external_wms` branch with a real async worker loop, signed request/callback protocol, callback-driven task state transitions, and external sellable stock reads. Order state continues to depend on shared inventory status models instead of direct WMS coupling, so phase 2 stays additive to the phase 1 architecture.

**Tech Stack:** Go 1.26, Gin, GORM, YAML config via Viper, existing plugin registry, docs-site markdown, existing inventory shared models.

---

## File Structure

- Modify `server/core/inventory/tasks.go`: add task claim, processing heartbeat metadata, retry scheduling, and final state transitions.
- Create `server/core/inventory/worker.go`: worker tick loop helpers that load due tasks and execute external operations safely.
- Create `server/core/inventory/worker_test.go`: retry, claim, timeout, and duplicate execution tests.
- Modify `server/core/inventory/external.go`: add signed sync request building, async execution entrypoint reuse, callback payload parsing, and external stock query support.
- Modify `server/core/inventory/external_test.go`: cover signing, sellable stock queries, and callback status mapping.
- Modify `server/core/inventory/model.go`: extend task state metadata if needed for lock owner, lock expiry, callback ids, and failure reason.
- Modify `server/config/config.go`: add external WMS signing, callback validation, and worker interval configuration.
- Modify `server/config.example.yaml`: document new phase 2 `external_wms` config.
- Modify `config.example.yaml`: document root deployment defaults for external WMS production mode.
- Modify `server/plugins/external_wms/plugin.go`: register background worker startup and wire new service helpers.
- Create `server/plugins/external_wms/service/worker.go`: start/stop plugin-level task worker if plugin lifecycle patterns require it.
- Create `server/plugins/external_wms/service/worker_test.go`: plugin worker lifecycle tests if service wrapper is added.
- Modify `server/plugins/external_wms/api/admin.go`: add signed callback validation and stronger retry/task detail endpoints where needed.
- Modify `server/plugins/external_wms/api/admin_test.go`: callback idempotency and invalid signature tests.
- Modify `server/plugins/order/service/order.go`: enforce async inventory status transitions in order flow where phase 1 left permissive handling.
- Modify `server/plugins/order/service/order_inventory_test.go`: add async callback completion and final failure scenarios.
- Modify `server/plugins/product/service/product.go`: route inventory stock reads through external provider when configured.
- Modify `server/plugins/product/service/product_inventory_test.go`: verify external sellable stock reads and fallback behavior.
- Modify `docs-site/docs/api/order.md`: document latest async inventory order status semantics.
- Modify `docs-site/docs/api/stock-reservation.md`: document async task states, retry behavior, and callback handling.
- Modify `docs-site/docs/api/wms.md`: document external WMS signed protocol and clarify built-in vs external provider behavior.
- Modify `docs-site/docs/guide/features.md`: update feature matrix for external WMS production readiness.

## Implementation Tasks

### Task 1: Async Inventory Task Execution Loop

**Files:**
- Modify: `server/core/inventory/tasks.go`
- Create: `server/core/inventory/worker.go`
- Create: `server/core/inventory/worker_test.go`
- Modify: `server/core/inventory/model.go`
- Modify: `server/plugins/external_wms/plugin.go`

- [ ] **Step 1: Add the failing worker execution tests**

Create `server/core/inventory/worker_test.go`:

```go
package inventory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClaimDueTaskSkipsLockedRecords(t *testing.T) {
	db := newInventoryTestDB(t)
	now := time.Now()

	task := InventoryIntegrationTask{
		Provider:     "external_wms",
		Operation:    "reserve",
		BizType:      "order",
		BizNo:        "T1001",
		Status:       TaskStatusPending,
		NextRetryAt:  &now,
		LockExpiresAt: ptrTime(now.Add(2 * time.Minute)),
		LockOwner:    "worker-a",
	}
	require.NoError(t, db.Create(&task).Error)

	claimed, err := ClaimDueTask(db, "worker-b", now)
	require.NoError(t, err)
	require.Nil(t, claimed)
}

func TestProcessTaskMarksSuccess(t *testing.T) {
	db := newInventoryTestDB(t)
	provider := &workerStubProvider{processErr: nil}
	now := time.Now()

	task := InventoryIntegrationTask{
		Provider:    "external_wms",
		Operation:   "deduct",
		BizType:     "order",
		BizNo:       "T1002",
		Status:      TaskStatusPending,
		MaxAttempts: 3,
		NextRetryAt: &now,
	}
	require.NoError(t, db.Create(&task).Error)

	err := ProcessTask(db, provider, &task, "worker-a", now)
	require.NoError(t, err)

	var latest InventoryIntegrationTask
	require.NoError(t, db.First(&latest, task.ID).Error)
	require.Equal(t, TaskStatusSuccess, latest.Status)
	require.Equal(t, 1, latest.AttemptCount)
}

func TestProcessTaskSchedulesRetryOnTransientFailure(t *testing.T) {
	db := newInventoryTestDB(t)
	provider := &workerStubProvider{processErr: ErrInventoryBusy}
	now := time.Now()

	task := InventoryIntegrationTask{
		Provider:      "external_wms",
		Operation:     "reserve",
		BizType:       "order",
		BizNo:         "T1003",
		Status:        TaskStatusPending,
		MaxAttempts:   3,
		BackoffSeconds: 30,
		NextRetryAt:   &now,
	}
	require.NoError(t, db.Create(&task).Error)

	err := ProcessTask(db, provider, &task, "worker-a", now)
	require.NoError(t, err)

	var latest InventoryIntegrationTask
	require.NoError(t, db.First(&latest, task.ID).Error)
	require.Equal(t, TaskStatusPending, latest.Status)
	require.Equal(t, 1, latest.AttemptCount)
	require.NotNil(t, latest.NextRetryAt)
	require.True(t, latest.NextRetryAt.After(now))
}
```

- [ ] **Step 2: Run the worker tests to verify they fail**

Run:

```powershell
cd server
go test ./core/inventory -run "TestClaimDueTaskSkipsLockedRecords|TestProcessTaskMarksSuccess|TestProcessTaskSchedulesRetryOnTransientFailure" -count=1
```

Expected: FAIL because `ClaimDueTask`, `ProcessTask`, task lock fields, and worker test helpers do not exist.

- [ ] **Step 3: Add task lock and retry metadata to the shared model**

Update `server/core/inventory/model.go`:

```go
type InventoryIntegrationTask struct {
	model.Model
	Provider        string     `gorm:"size:32;index"`
	Operation       string     `gorm:"size:32;index"`
	BizType         string     `gorm:"size:32;index"`
	BizNo           string     `gorm:"size:64;index"`
	Status          string     `gorm:"size:32;index"`
	RequestID       string     `gorm:"size:64;index"`
	Payload         string     `gorm:"type:longtext"`
	AttemptCount    int        `gorm:"not null;default:0"`
	MaxAttempts     int        `gorm:"not null;default:0"`
	BackoffSeconds  int        `gorm:"not null;default:0"`
	NextRetryAt     *time.Time `gorm:"index"`
	LastError       string     `gorm:"size:255"`
	LockOwner       string     `gorm:"size:64;index"`
	LockExpiresAt   *time.Time `gorm:"index"`
	CompletedAt     *time.Time
	LastCallbackID  string     `gorm:"size:64;index"`
}
```

Also keep the model import list aligned with `time` and the existing shared model package structure.

- [ ] **Step 4: Implement due-task claim and processing helpers**

Create `server/core/inventory/worker.go`:

```go
package inventory

import (
	"time"

	"gorm.io/gorm"
)

type AsyncTaskProcessor interface {
	ProcessTask(tx *gorm.DB, task *InventoryIntegrationTask, now time.Time) error
}

func ClaimDueTask(db *gorm.DB, worker string, now time.Time) (*InventoryIntegrationTask, error) {
	var task InventoryIntegrationTask
	err := db.Transaction(func(tx *gorm.DB) error {
		query := tx.
			Where("status = ?", TaskStatusPending).
			Where("next_retry_at IS NULL OR next_retry_at <= ?", now).
			Where("lock_expires_at IS NULL OR lock_expires_at <= ?", now).
			Order("id ASC").
			First(&task).Error
		if err != nil {
			return err
		}

		lockUntil := now.Add(2 * time.Minute)
		return tx.Model(&InventoryIntegrationTask{}).
			Where("id = ?", task.ID).
			Updates(map[string]any{
				"status":          TaskStatusProcessing,
				"lock_owner":      worker,
				"lock_expires_at": lockUntil,
			}).Error
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	task.Status = TaskStatusProcessing
	task.LockOwner = worker
	expiresAt := now.Add(2 * time.Minute)
	task.LockExpiresAt = &expiresAt
	return &task, nil
}

func ProcessTask(db *gorm.DB, processor AsyncTaskProcessor, task *InventoryIntegrationTask, worker string, now time.Time) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := processor.ProcessTask(tx, task, now); err != nil {
			return MarkTaskRetry(tx, task, err, now)
		}
		return MarkTaskSuccess(tx, task, now)
	})
	if err != nil {
		return err
	}
	task.LockOwner = worker
	return nil
}
```

- [ ] **Step 5: Extend task state transitions for success and retry**

Update `server/core/inventory/tasks.go`:

```go
func MarkTaskSuccess(tx *gorm.DB, task *InventoryIntegrationTask, now time.Time) error {
	return tx.Model(&InventoryIntegrationTask{}).
		Where("id = ?", task.ID).
		Updates(map[string]any{
			"status":          TaskStatusSuccess,
			"attempt_count":   gorm.Expr("attempt_count + 1"),
			"last_error":      "",
			"lock_owner":      "",
			"lock_expires_at": nil,
			"completed_at":    now,
		}).Error
}

func MarkTaskRetry(tx *gorm.DB, task *InventoryIntegrationTask, cause error, now time.Time) error {
	next := now.Add(time.Duration(task.BackoffSeconds) * time.Second)
	updates := map[string]any{
		"attempt_count":   gorm.Expr("attempt_count + 1"),
		"last_error":      cause.Error(),
		"lock_owner":      "",
		"lock_expires_at": nil,
	}
	if task.AttemptCount+1 >= task.MaxAttempts {
		updates["status"] = TaskStatusFailed
		updates["completed_at"] = now
		updates["next_retry_at"] = nil
	} else {
		updates["status"] = TaskStatusPending
		updates["next_retry_at"] = next
	}
	return tx.Model(&InventoryIntegrationTask{}).
		Where("id = ?", task.ID).
		Updates(updates).Error
}
```

- [ ] **Step 6: Run inventory tests to verify the worker logic passes**

Run:

```powershell
cd server
go test ./core/inventory -count=1
```

Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add server/core/inventory/model.go server/core/inventory/tasks.go server/core/inventory/worker.go server/core/inventory/worker_test.go server/plugins/external_wms/plugin.go
git commit -m "补齐外部WMS异步任务执行骨架" -m "新增库存异步任务claim与处理流程\n补充共享任务锁定与重试元数据\n为外部WMS后台执行器预留启动接入点"
```

### Task 2: Callback State Machine And Order Inventory Status

**Files:**
- Modify: `server/core/inventory/tasks.go`
- Modify: `server/core/inventory/external.go`
- Modify: `server/plugins/external_wms/api/admin.go`
- Modify: `server/plugins/external_wms/api/admin_test.go`
- Modify: `server/plugins/order/service/order.go`
- Modify: `server/plugins/order/service/order_inventory_test.go`

- [ ] **Step 1: Add the failing callback idempotency tests**

Append to `server/plugins/external_wms/api/admin_test.go`:

```go
func TestExternalWMSCallbackIgnoresDuplicateCallbackID(t *testing.T) {
	db := newExternalWMSTestDB(t)
	now := time.Now()
	task := inventory.InventoryIntegrationTask{
		Provider:       "external_wms",
		Operation:      "reserve",
		BizType:        "order",
		BizNo:          "O1001",
		Status:         inventory.TaskStatusProcessing,
		RequestID:      "REQ-1",
		LastCallbackID: "CALLBACK-1",
		NextRetryAt:    &now,
	}
	require.NoError(t, db.Create(&task).Error)

	router := newExternalWMSTestRouter(db)
	body := `{"request_id":"REQ-1","callback_id":"CALLBACK-1","status":"success"}`
	req := httptest.NewRequest(http.MethodPost, "/admin/api/external-wms/callback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var latest inventory.InventoryIntegrationTask
	require.NoError(t, db.First(&latest, task.ID).Error)
	require.Equal(t, inventory.TaskStatusProcessing, latest.Status)
}

func TestExternalWMSCallbackMarksOrderInventoryFailed(t *testing.T) {
	db := newExternalWMSTestDB(t)
	now := time.Now()
	task := inventory.InventoryIntegrationTask{
		Provider:    "external_wms",
		Operation:   "deduct",
		BizType:     "order",
		BizNo:       "O1002",
		Status:      inventory.TaskStatusProcessing,
		RequestID:   "REQ-2",
		NextRetryAt: &now,
	}
	require.NoError(t, db.Create(&task).Error)
	require.NoError(t, db.Model(&orderModel.Order{}).Where("order_sn = ?", "O1002").Update("inventory_status", "pending").Error)

	router := newExternalWMSTestRouter(db)
	body := `{"request_id":"REQ-2","callback_id":"CALLBACK-2","status":"failed","message":"stock not enough"}`
	req := httptest.NewRequest(http.MethodPost, "/admin/api/external-wms/callback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	requireOrderInventoryStatus(t, db, "O1002", "failed")
}
```

- [ ] **Step 2: Run the callback tests to verify they fail**

Run:

```powershell
cd server
go test ./plugins/external_wms/api -run "TestExternalWMSCallbackIgnoresDuplicateCallbackID|TestExternalWMSCallbackMarksOrderInventoryFailed" -count=1
```

Expected: FAIL because duplicate callback handling and order inventory status linkage do not exist.

- [ ] **Step 3: Add callback update helpers for task and order state**

Update `server/core/inventory/tasks.go`:

```go
func CompleteTaskByCallback(tx *gorm.DB, requestID, callbackID, status, message string, now time.Time) error {
	var task InventoryIntegrationTask
	if err := tx.Where("request_id = ?", requestID).First(&task).Error; err != nil {
		return err
	}
	if task.LastCallbackID == callbackID {
		return nil
	}
	if task.Status == TaskStatusSuccess || task.Status == TaskStatusFailed {
		return nil
	}

	updates := map[string]any{
		"last_callback_id": callbackID,
		"last_error":       message,
		"lock_owner":       "",
		"lock_expires_at":  nil,
	}
	switch status {
	case "success":
		updates["status"] = TaskStatusSuccess
		updates["completed_at"] = now
	case "failed":
		updates["status"] = TaskStatusFailed
		updates["completed_at"] = now
	default:
		updates["status"] = TaskStatusProcessing
	}
	if err := tx.Model(&InventoryIntegrationTask{}).Where("id = ?", task.ID).Updates(updates).Error; err != nil {
		return err
	}
	return UpdateOrderInventoryStatusByTask(tx, task.BizType, task.BizNo, status)
}
```

- [ ] **Step 4: Add order inventory status mapping helpers**

Update `server/plugins/order/service/order.go`:

```go
func UpdateInventoryStatusTx(tx *gorm.DB, orderSN string, inventoryStatus string) error {
	return tx.Model(&model.Order{}).
		Where("order_sn = ?", orderSN).
		Update("inventory_status", inventoryStatus).Error
}

func InventoryStatusFromTaskStatus(status string) string {
	switch status {
	case "success":
		return "success"
	case "failed":
		return "failed"
	default:
		return "pending"
	}
}
```

Also expose a small helper that `inventory` can call without reintroducing direct WMS coupling. If package boundaries make direct import awkward, move the mapper into `server/core/inventory/router.go` and call `tx.Table("orders")` directly to keep phase 2 scope tight.

- [ ] **Step 5: Wire callback handler to the shared state machine**

Update `server/plugins/external_wms/api/admin.go`:

```go
func callback(c *gin.Context) {
	var req callbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		return inventory.CompleteTaskByCallback(tx, req.RequestID, req.CallbackID, req.Status, req.Message, time.Now())
	})
	if err != nil {
		response.FailWithMessage("回调处理失败", c)
		return
	}

	response.OkWithMessage("ok", c)
}
```

- [ ] **Step 6: Run external WMS and order tests**

Run:

```powershell
cd server
go test ./plugins/external_wms/api ./plugins/order/service -count=1
```

Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add server/core/inventory/tasks.go server/core/inventory/external.go server/plugins/external_wms/api/admin.go server/plugins/external_wms/api/admin_test.go server/plugins/order/service/order.go server/plugins/order/service/order_inventory_test.go
git commit -m "补齐外部WMS回调状态联动" -m "新增异步回调幂等与任务状态机更新\n补充订单库存状态与任务结果映射\n保证异步库存处理结果可回写订单"
```

### Task 3: Signed External WMS Request And Callback Verification

**Files:**
- Modify: `server/config/config.go`
- Modify: `server/config.example.yaml`
- Modify: `config.example.yaml`
- Modify: `server/core/inventory/external.go`
- Modify: `server/core/inventory/external_test.go`
- Modify: `server/plugins/external_wms/api/admin.go`
- Modify: `server/plugins/external_wms/api/admin_test.go`

- [ ] **Step 1: Add the failing signing tests**

Append to `server/core/inventory/external_test.go`:

```go
func TestBuildExternalSignature(t *testing.T) {
	req := signedPayload{
		AppKey:    "demo-key",
		Timestamp: "1717910400",
		Nonce:     "nonce-1",
		Body:      `{"biz_no":"O1001"}`,
	}

	sign := BuildSignature(req, "demo-secret")
	require.Equal(t, "4d8560a2f8a7f97031d2f01cd4f67f6d4f8fd8c9a5bc2c321c5f831c3dbe4f6c", sign)
}

func TestVerifyCallbackSignatureRejectsInvalidSign(t *testing.T) {
	err := VerifyCallbackSignature(callbackEnvelope{
		AppKey:    "demo-key",
		Timestamp: "1717910400",
		Nonce:     "nonce-1",
		Sign:      "bad-sign",
		Body:      `{"request_id":"REQ-1"}`,
	}, "demo-secret", time.Unix(1717910400, 0))
	require.ErrorContains(t, err, "signature")
}
```

- [ ] **Step 2: Run signing tests to verify they fail**

Run:

```powershell
cd server
go test ./core/inventory -run "TestBuildExternalSignature|TestVerifyCallbackSignatureRejectsInvalidSign" -count=1
```

Expected: FAIL because signature helpers do not exist.

- [ ] **Step 3: Add config for signing and worker controls**

Update `server/config/config.go`:

```go
type ExternalWMSConfig struct {
	Endpoint          string                 `mapstructure:"endpoint"`
	AppKey            string                 `mapstructure:"app_key"`
	AppSecret         string                 `mapstructure:"app_secret"`
	TimeoutMS         int                    `mapstructure:"timeout_ms"`
	CallbackEnabled   bool                   `mapstructure:"callback_enabled"`
	SignatureTTL      int                    `mapstructure:"signature_ttl"`
	WorkerIntervalSec int                    `mapstructure:"worker_interval_sec"`
	Retry             ExternalWMSRetryConfig `mapstructure:"retry"`
}
```

Update both example configs:

```yaml
external_wms:
  endpoint: ""
  app_key: ""
  app_secret: ""
  timeout_ms: 3000
  callback_enabled: true
  signature_ttl: 300
  worker_interval_sec: 5
  retry:
    max_attempts: 8
    backoff_seconds: 30
```

- [ ] **Step 4: Implement request signing and callback verification**

Update `server/core/inventory/external.go`:

```go
func BuildSignature(req signedPayload, secret string) string {
	raw := strings.Join([]string{
		req.AppKey,
		req.Timestamp,
		req.Nonce,
		req.Body,
		secret,
	}, "\n")
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func VerifyCallbackSignature(req callbackEnvelope, secret string, now time.Time) error {
	ts, err := strconv.ParseInt(req.Timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %w", err)
	}
	if absDuration(now.Sub(time.Unix(ts, 0))) > time.Duration(config.Global.ExternalWMS.SignatureTTL)*time.Second {
		return fmt.Errorf("signature expired")
	}

	expected := BuildSignature(signedPayload{
		AppKey:    req.AppKey,
		Timestamp: req.Timestamp,
		Nonce:     req.Nonce,
		Body:      req.Body,
	}, secret)
	if !strings.EqualFold(expected, req.Sign) {
		return fmt.Errorf("signature mismatch")
	}
	return nil
}
```

- [ ] **Step 5: Enforce callback signature verification in the API handler**

Update `server/plugins/external_wms/api/admin.go`:

```go
if err := inventory.VerifyCallbackSignature(callbackEnvelopeFromRequest(req), config.Global.ExternalWMS.AppSecret, time.Now()); err != nil {
	response.FailWithMessage("签名校验失败", c)
	return
}
```

- [ ] **Step 6: Run inventory and external WMS API tests**

Run:

```powershell
cd server
go test ./core/inventory ./plugins/external_wms/api -count=1
```

Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add server/config/config.go server/config.example.yaml config.example.yaml server/core/inventory/external.go server/core/inventory/external_test.go server/plugins/external_wms/api/admin.go server/plugins/external_wms/api/admin_test.go
git commit -m "补齐外部WMS签名与验签能力" -m "新增外部WMS请求签名与回调验签实现\n补充签名时效与worker调度配置\n增强企业外部WMS接入安全能力"
```

### Task 4: External Sellable Stock Query Support

**Files:**
- Modify: `server/core/inventory/external.go`
- Modify: `server/core/inventory/external_test.go`
- Modify: `server/plugins/product/service/product.go`
- Modify: `server/plugins/product/service/product_inventory_test.go`

- [ ] **Step 1: Add the failing external stock query tests**

Append to `server/core/inventory/external_test.go`:

```go
func TestExternalProviderGetSellableStock(t *testing.T) {
	server := newExternalWMSTestServer(t, http.StatusOK, `{"code":0,"data":[{"sku_id":11,"sellable_stock":7},{"sku_id":12,"sellable_stock":0}]}`)
	defer server.Close()

	original := config.Global
	t.Cleanup(func() { config.Global = original })
	config.Global.ExternalWMS.Endpoint = server.URL

	provider := NewExternalProvider()
	stocks, err := provider.GetSellableStock(context.Background(), []uint64{11, 12})
	require.NoError(t, err)
	require.Len(t, stocks, 2)
	require.Equal(t, uint64(11), stocks[0].SkuID)
	require.Equal(t, 7, stocks[0].SellableStock)
}

func TestExternalProviderGetSellableStockMapsRemoteError(t *testing.T) {
	server := newExternalWMSTestServer(t, http.StatusBadGateway, `{"code":5001,"msg":"remote error"}`)
	defer server.Close()

	original := config.Global
	t.Cleanup(func() { config.Global = original })
	config.Global.ExternalWMS.Endpoint = server.URL

	provider := NewExternalProvider()
	_, err := provider.GetSellableStock(context.Background(), []uint64{11})
	require.ErrorContains(t, err, "remote error")
}
```

- [ ] **Step 2: Run the external stock tests to verify they fail**

Run:

```powershell
cd server
go test ./core/inventory -run "TestExternalProviderGetSellableStock|TestExternalProviderGetSellableStockMapsRemoteError" -count=1
```

Expected: FAIL because `GetSellableStock` still returns a placeholder result.

- [ ] **Step 3: Implement external stock query support**

Update `server/core/inventory/external.go`:

```go
func (p *ExternalProvider) GetSellableStock(ctx context.Context, skuIDs []uint64) ([]SellableStock, error) {
	reqBody, err := json.Marshal(map[string]any{
		"sku_ids": skuIDs,
	})
	if err != nil {
		return nil, err
	}

	body, err := p.doSignedRequest(ctx, "POST", "/stock/sellable", reqBody)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data []SellableStock `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("external wms stock query failed: %s", resp.Msg)
	}
	return resp.Data, nil
}
```

- [ ] **Step 4: Route product stock reads through the unified inventory provider**

Update `server/plugins/product/service/product.go` in the stock read path:

```go
provider, err := inventory.CurrentProvider()
if err == nil && provider.Name() == "external_wms" {
	stocks, stockErr := provider.GetSellableStock(ctx, skuIDs)
	if stockErr == nil {
		stockMap := make(map[uint64]int, len(stocks))
		for _, item := range stocks {
			stockMap[item.SkuID] = item.SellableStock
		}
		applySellableStockToSkus(list, stockMap)
	}
}
```

Keep the fallback path unchanged so existing local/builtin behavior is not regressed if the remote stock query fails.

- [ ] **Step 5: Run inventory and product tests**

Run:

```powershell
cd server
go test ./core/inventory ./plugins/product/service -count=1
```

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add server/core/inventory/external.go server/core/inventory/external_test.go server/plugins/product/service/product.go server/plugins/product/service/product_inventory_test.go
git commit -m "补齐外部WMS可售库存查询" -m "新增外部WMS可售库存查询实现\n将商品库存读取接入统一库存Provider\n保留本地与内置WMS回退路径"
```

### Task 5: End-To-End Inventory Mode Coverage

**Files:**
- Modify: `server/plugins/order/service/order_inventory_test.go`
- Modify: `server/plugins/product/service/product_inventory_test.go`
- Modify: `server/plugins/points_mall/service/exchange_inventory_test.go`

- [ ] **Step 1: Add end-to-end coverage for all inventory modes**

Append to `server/plugins/order/service/order_inventory_test.go`:

```go
func TestOrderInventoryFlowAcrossProviders(t *testing.T) {
	cases := []struct {
		name         string
		provider     string
		externalMode string
		wantStatus   string
	}{
		{name: "local", provider: "local", wantStatus: "success"},
		{name: "builtin", provider: "builtin_wms", wantStatus: "success"},
		{name: "external-sync", provider: "external_wms", externalMode: "sync", wantStatus: "success"},
		{name: "external-async", provider: "external_wms", externalMode: "async", wantStatus: "pending"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := newOrderInventoryTestContext(t, tc.provider, tc.externalMode)
			orderSN := createPaidOrderThroughInventoryFlow(t, ctx)
			requireOrderInventoryStatus(t, ctx.DB, orderSN, tc.wantStatus)
		})
	}
}
```

- [ ] **Step 2: Run the targeted service tests to verify any missing coverage fails first**

Run:

```powershell
cd server
go test ./plugins/order/service ./plugins/product/service ./plugins/points_mall/service -count=1
```

Expected: FAIL if helper coverage or inventory status expectations are incomplete.

- [ ] **Step 3: Add minimal test helpers to support async completion assertions**

Update the service test files with shared helper patterns already used in phase 1:

```go
func requireOrderInventoryStatus(t *testing.T, db *gorm.DB, orderSN string, want string) {
	t.Helper()
	var order model.Order
	require.NoError(t, db.Where("order_sn = ?", orderSN).First(&order).Error)
	require.Equal(t, want, order.InventoryStatus)
}
```

Also add an async callback simulation helper in `server/plugins/order/service/order_inventory_test.go` if the file does not already have one:

```go
func completeExternalInventoryTask(t *testing.T, db *gorm.DB, requestID string, status string) {
	t.Helper()
	require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
		return inventory.CompleteTaskByCallback(tx, requestID, "callback-"+status, status, "", time.Now())
	}))
}
```

- [ ] **Step 4: Run the expanded service test suite**

Run:

```powershell
cd server
go test ./plugins/order/service ./plugins/product/service ./plugins/points_mall/service ./plugins/external_wms/api ./core/inventory -count=1
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/plugins/order/service/order_inventory_test.go server/plugins/product/service/product_inventory_test.go server/plugins/points_mall/service/exchange_inventory_test.go
git commit -m "补齐统一库存多模式链路测试" -m "新增本地内置WMS外部WMS链路覆盖\n补充异步库存回调完成断言\n降低库存架构改造回归风险"
```

### Task 6: Update Latest-Architecture Docs For External WMS Production Mode

**Files:**
- Modify: `docs-site/docs/api/order.md`
- Modify: `docs-site/docs/api/stock-reservation.md`
- Modify: `docs-site/docs/api/wms.md`
- Modify: `docs-site/docs/guide/features.md`
- Modify: `server/config.example.yaml`
- Modify: `config.example.yaml`

- [ ] **Step 1: Add the failing doc parity checklist**

Create a local checklist in the task notes and verify each item is represented in docs before editing:

```text
[ ] order docs explain inventory_status in async external mode
[ ] stock reservation docs explain task states, retry, callback
[ ] wms docs explain built-in provider vs external provider
[ ] feature guide lists local / builtin_wms / external_wms
[ ] example configs show signature_ttl and worker_interval_sec
```

- [ ] **Step 2: Update order API docs to describe async inventory status**

Add to `docs-site/docs/api/order.md`:

```md
## inventory_status

- `success`：库存已完成确认或扣减
- `pending`：外部 WMS 异步处理中，订单可见但库存结果未最终确认
- `failed`：外部 WMS 返回失败或异步任务最终失败，需要人工处理或业务补偿
```

- [ ] **Step 3: Update stock reservation docs to describe async task states**

Add to `docs-site/docs/api/stock-reservation.md`:

```md
## 异步外部库存任务

统一库存架构在 `external_wms + async` 模式下，会把库存操作写入异步任务表并由后台任务执行。

任务状态：

- `pending`：待执行或待重试
- `processing`：正在执行或等待回调确认
- `success`：库存操作已完成
- `failed`：超过最大重试次数或收到明确失败回调
```

- [ ] **Step 4: Update WMS and feature docs to describe the latest architecture**

Add to `docs-site/docs/api/wms.md` and `docs-site/docs/guide/features.md`:

```md
LYShop 当前库存模式支持：

- `local`：不启用 WMS，直接使用商城本地库存
- `builtin_wms`：启用内置 WMS 插件，使用系统内仓储能力
- `external_wms`：对接企业已有 WMS，支持同步与异步两种集成方式
```

Also add the external signing fields to both example configs so deployment docs stay aligned with runtime behavior.

- [ ] **Step 5: Verify docs and config examples are consistent**

Run:

```powershell
Get-Content docs-site/docs/api/order.md
Get-Content docs-site/docs/api/stock-reservation.md
Get-Content docs-site/docs/api/wms.md
Get-Content docs-site/docs/guide/features.md
Get-Content server/config.example.yaml
Get-Content config.example.yaml
```

Expected: each checklist item from Step 1 is visibly covered with latest-architecture wording.

- [ ] **Step 6: Commit**

```bash
git add docs-site/docs/api/order.md docs-site/docs/api/stock-reservation.md docs-site/docs/api/wms.md docs-site/docs/guide/features.md server/config.example.yaml config.example.yaml
git commit -m "更新外部WMS生产模式文档" -m "补充异步库存状态与任务说明\n更新外部WMS签名和调度配置示例\n保持文档对齐统一库存最新架构"
```

