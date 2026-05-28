# WMS Inventory Engine Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将后台仓储管理升级为产品可用闭环，落地仓库管理、出入库单（多 SKU）、库存台账、库存流水四个能力模块。

**Architecture:** 后端采用“库存流水为事实源 + 库存快照冗余查询”的双模型。所有库存变化通过“完成单据”事务写入 `movement` 并更新 `stock`。前端按业务域拆成仓库、单据、台账、流水四个页面，统一调用新的 `/admin/api/wms/*` 接口。

**Tech Stack:** Go + Gin + GORM + MySQL（测试用 sqlite 内存库）；Vue3 + TypeScript + Vue Router + Axios；VitePress 文档站。

---

## File Structure

### Backend

- Create: `server/plugins/wms/model/inventory.go`
- Create: `server/plugins/wms/service/warehouse.go`
- Create: `server/plugins/wms/service/inventory_stock.go`
- Create: `server/plugins/wms/service/inventory_doc.go`
- Create: `server/plugins/wms/service/inventory_movement.go`
- Create: `server/plugins/wms/service/inventory_doc_complete_test.go`
- Create: `server/plugins/wms/api/admin_test.go`
- Modify: `server/plugins/wms/model/wms.go`（保留旧结构或清理旧结构，最终由新模型承载）
- Modify: `server/plugins/wms/api/admin.go`
- Modify: `server/plugins/wms/plugin.go`
- Modify: `server/plugins/wms/plugin.json`

### Admin

- Create: `admin/src/api/wms.ts`
- Create: `admin/src/views/wms/WarehouseList.vue`
- Create: `admin/src/views/wms/DocList.vue`
- Create: `admin/src/views/wms/DocEditor.vue`
- Create: `admin/src/views/wms/StockLedger.vue`
- Create: `admin/src/views/wms/MovementList.vue`
- Modify: `admin/src/router/index.ts`
- Modify: `admin/src/locales/zh-CN.ts`
- Modify: `admin/src/locales/en.ts`
- Modify: `admin/src/views/wms/StockList.vue`（替换为台账页或删除并迁移）

### Docs

- Create: `docs-site/docs/api/wms.md`
- Modify: `docs-site/docs/api/index.md`
- Modify: `docs-site/docs/guide/features.md`
- Modify: `docs-site/docs/.vitepress/config.mts`

---

### Task 1: 重建 WMS 数据模型与迁移入口

**Files:**
- Create: `server/plugins/wms/model/inventory.go`
- Modify: `server/plugins/wms/plugin.go`
- Test: `server/plugins/wms/model/inventory_model_test.go`

- [ ] **Step 1: 写失败测试（状态常量与单据结构）**

```go
package model

import "testing"

func TestDocStatusConstants(t *testing.T) {
	if DocStatusDraft != "draft" {
		t.Fatalf("unexpected draft status: %s", DocStatusDraft)
	}
	if DocStatusCompleted != "completed" {
		t.Fatalf("unexpected completed status: %s", DocStatusCompleted)
	}
}
```

- [ ] **Step 2: 运行测试，确认失败**

Run: `cd server && go test ./plugins/wms/model -run TestDocStatusConstants -v`  
Expected: FAIL，提示 `undefined: DocStatusDraft`。

- [ ] **Step 3: 写最小实现（新模型 + plugin migrate）**

```go
// server/plugins/wms/model/inventory.go
package model

import "github.com/ijry/lyshop/model"

const (
	DocTypeInbound  = "inbound"
	DocTypeOutbound = "outbound"

	DocStatusDraft     = "draft"
	DocStatusCompleted = "completed"
	DocStatusCanceled  = "canceled"
)

type InventoryDoc struct {
	model.Base
	DocNo       string `gorm:"size:32;uniqueIndex;not null" json:"doc_no"`
	DocType     string `gorm:"size:16;index;not null" json:"doc_type"`
	WarehouseID uint64 `gorm:"index;not null" json:"warehouse_id"`
	Status      string `gorm:"size:16;index;not null" json:"status"`
	Remark      string `gorm:"size:255" json:"remark"`
	CompletedAt *int64 `json:"completed_at"`
	OperatorID  uint64 `gorm:"default:0" json:"operator_id"`
}
```

```go
// server/plugins/wms/plugin.go (Migrate)
return db.AutoMigrate(
	&wmsmodel.Warehouse{},
	&wmsmodel.InventoryStock{},
	&wmsmodel.InventoryMovement{},
	&wmsmodel.InventoryDoc{},
	&wmsmodel.InventoryDocItem{},
)
```

- [ ] **Step 4: 运行测试确认通过**

Run: `cd server && go test ./plugins/wms/model -run TestDocStatusConstants -v`  
Expected: PASS。

- [ ] **Step 5: 提交**

```bash
git add server/plugins/wms/model/inventory.go server/plugins/wms/plugin.go server/plugins/wms/model/inventory_model_test.go
git commit -m "重建WMS核心数据模型与迁移入口" -m "新增库存快照、库存流水、出入库单头与单据行模型；插件迁移入口切换为新模型集合，为统一库存交易引擎落地提供基础。"
```

---

### Task 2: 仓库管理接口（列表/新增/编辑/启停）

**Files:**
- Create: `server/plugins/wms/service/warehouse.go`
- Modify: `server/plugins/wms/api/admin.go`
- Test: `server/plugins/wms/api/admin_test.go`

- [ ] **Step 1: 写失败测试（启停接口）**

```go
func TestUpdateWarehouseStatus(t *testing.T) {
	r := gin.New()
	g := r.Group("/admin/api")
	RegisterAdminRoutes(g)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/api/wms/warehouses/1/status", strings.NewReader(`{"status":0}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatalf("route not registered")
	}
}
```

- [ ] **Step 2: 运行测试，确认失败**

Run: `cd server && go test ./plugins/wms/api -run TestUpdateWarehouseStatus -v`  
Expected: FAIL，返回 404 或 handler 缺失。

- [ ] **Step 3: 实现仓库服务与路由**

```go
// server/plugins/wms/service/warehouse.go
func UpdateWarehouseStatus(ctx context.Context, id uint64, status int8) error {
	return db.DB.WithContext(ctx).
		Model(&wmsmodel.Warehouse{}).
		Where("id = ?", id).
		Update("status", status).Error
}
```

```go
// server/plugins/wms/api/admin.go
g.PUT("/wms/warehouses/:id/status", middleware.RequirePermission("wms:edit"), updateWarehouseStatus)
```

- [ ] **Step 4: 运行测试确认通过**

Run: `cd server && go test ./plugins/wms/api -run TestUpdateWarehouseStatus -v`  
Expected: PASS。

- [ ] **Step 5: 提交**

```bash
git add server/plugins/wms/service/warehouse.go server/plugins/wms/api/admin.go server/plugins/wms/api/admin_test.go
git commit -m "完善仓库管理接口能力" -m "新增仓库启停与编辑能力，补齐仓库列表管理所需后端接口，并通过路由级测试覆盖新增端点。"
```

---

### Task 3: 库存台账查询与安全库存维护

**Files:**
- Create: `server/plugins/wms/service/inventory_stock.go`
- Modify: `server/plugins/wms/api/admin.go`
- Test: `server/plugins/wms/service/inventory_stock_test.go`

- [ ] **Step 1: 写失败测试（预警判定）**

```go
func TestIsStockWarning(t *testing.T) {
	if !isStockWarning(9, 10) {
		t.Fatalf("expected warning")
	}
	if isStockWarning(11, 10) {
		t.Fatalf("expected normal")
	}
}
```

- [ ] **Step 2: 运行测试，确认失败**

Run: `cd server && go test ./plugins/wms/service -run TestIsStockWarning -v`  
Expected: FAIL，提示 `undefined: isStockWarning`。

- [ ] **Step 3: 实现库存台账服务与接口**

```go
func isStockWarning(availableQty, safetyQty int) bool {
	return availableQty <= safetyQty
}

func UpdateSafetyQty(ctx context.Context, id uint64, safetyQty int) error {
	return db.DB.WithContext(ctx).
		Model(&wmsmodel.InventoryStock{}).
		Where("id = ?", id).
		Update("safety_qty", safetyQty).Error
}
```

```go
g.PUT("/wms/stocks/:id/safety", middleware.RequirePermission("wms:edit"), updateSafetyQty)
```

- [ ] **Step 4: 运行测试确认通过**

Run: `cd server && go test ./plugins/wms/service -run TestIsStockWarning -v`  
Expected: PASS。

- [ ] **Step 5: 提交**

```bash
git add server/plugins/wms/service/inventory_stock.go server/plugins/wms/service/inventory_stock_test.go server/plugins/wms/api/admin.go
git commit -m "新增库存台账与安全库存维护接口" -m "提供库存台账分页筛选与预警标识，新增安全库存更新接口，满足库存治理最小可用需求。"
```

---

### Task 4: 出入库单草稿能力（创建/编辑/详情/作废）

**Files:**
- Create: `server/plugins/wms/service/inventory_doc.go`
- Modify: `server/plugins/wms/api/admin.go`
- Test: `server/plugins/wms/service/inventory_doc_test.go`

- [ ] **Step 1: 写失败测试（非草稿禁止编辑）**

```go
func TestEnsureDraftEditable(t *testing.T) {
	err := ensureDraftEditable("completed")
	if err == nil {
		t.Fatalf("expected error for completed doc")
	}
}
```

- [ ] **Step 2: 运行测试，确认失败**

Run: `cd server && go test ./plugins/wms/service -run TestEnsureDraftEditable -v`  
Expected: FAIL，提示 `undefined: ensureDraftEditable`。

- [ ] **Step 3: 实现单据草稿服务与接口**

```go
func ensureDraftEditable(status string) error {
	if status != wmsmodel.DocStatusDraft {
		return fmt.Errorf("only draft doc can be edited")
	}
	return nil
}

func CancelInventoryDoc(ctx context.Context, id uint64, operatorID uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var doc wmsmodel.InventoryDoc
		if err := tx.Where("id = ?", id).First(&doc).Error; err != nil {
			return err
		}
		if err := ensureDraftEditable(doc.Status); err != nil {
			return err
		}
		return tx.Model(&doc).Updates(map[string]any{"status": wmsmodel.DocStatusCanceled, "operator_id": operatorID}).Error
	})
}
```

```go
g.GET("/wms/docs", middleware.RequirePermission("wms:view"), listDocs)
g.POST("/wms/docs", middleware.RequirePermission("wms:edit"), createDoc)
g.GET("/wms/docs/:id", middleware.RequirePermission("wms:view"), getDoc)
g.PUT("/wms/docs/:id", middleware.RequirePermission("wms:edit"), updateDoc)
g.POST("/wms/docs/:id/cancel", middleware.RequirePermission("wms:edit"), cancelDoc)
```

- [ ] **Step 4: 运行测试确认通过**

Run: `cd server && go test ./plugins/wms/service -run TestEnsureDraftEditable -v`  
Expected: PASS。

- [ ] **Step 5: 提交**

```bash
git add server/plugins/wms/service/inventory_doc.go server/plugins/wms/service/inventory_doc_test.go server/plugins/wms/api/admin.go
git commit -m "实现出入库单草稿管理接口" -m "落地单据列表、创建、详情、编辑、作废能力，限制仅草稿可编辑，建立单据状态机基础。"
```

---

### Task 5: 单据完成事务（库存快照 + 流水原子更新）

**Files:**
- Create: `server/plugins/wms/service/inventory_doc_complete_test.go`
- Modify: `server/plugins/wms/service/inventory_doc.go`
- Create: `server/plugins/wms/service/inventory_movement.go`
- Modify: `server/go.mod`（测试依赖 sqlite）

- [ ] **Step 1: 写失败测试（完成出库单库存扣减并写流水）**

```go
func TestCompleteOutboundDoc_UpdatesStockAndMovement(t *testing.T) {
	gdb, err := gorm.Open(sqlite.Open("file:wms_complete_test?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })
	require.NoError(t, gdb.AutoMigrate(
		&wmsmodel.Warehouse{},
		&wmsmodel.InventoryStock{},
		&wmsmodel.InventoryDoc{},
		&wmsmodel.InventoryDocItem{},
		&wmsmodel.InventoryMovement{},
	))

	require.NoError(t, gdb.Create(&wmsmodel.Warehouse{Name: "主仓", Code: "WH-A", Status: 1}).Error)
	require.NoError(t, gdb.Create(&wmsmodel.InventoryStock{WarehouseID: 1, SkuID: 1001, AvailableQty: 20, SafetyQty: 5}).Error)
	doc := wmsmodel.InventoryDoc{DocNo: "OUT-TEST-001", DocType: wmsmodel.DocTypeOutbound, WarehouseID: 1, Status: wmsmodel.DocStatusDraft}
	require.NoError(t, gdb.Create(&doc).Error)
	require.NoError(t, gdb.Create(&wmsmodel.InventoryDocItem{DocID: doc.ID, SkuID: 1001, Qty: 6}).Error)

	require.NoError(t, CompleteInventoryDoc(context.Background(), doc.ID, 9527))

	var stock wmsmodel.InventoryStock
	require.NoError(t, gdb.Where("warehouse_id = ? AND sku_id = ?", 1, 1001).First(&stock).Error)
	require.Equal(t, 14, stock.AvailableQty)

	var movements []wmsmodel.InventoryMovement
	require.NoError(t, gdb.Where("biz_doc_id = ?", doc.ID).Find(&movements).Error)
	require.Len(t, movements, 1)
	require.Equal(t, "out", movements[0].Direction)
	require.Equal(t, 20, movements[0].BeforeQty)
	require.Equal(t, 14, movements[0].AfterQty)
}
```

- [ ] **Step 2: 运行测试，确认失败**

Run: `cd server && go test ./plugins/wms/service -run TestCompleteOutboundDoc_UpdatesStockAndMovement -v`  
Expected: FAIL，提示 `undefined: CompleteInventoryDoc`。

- [ ] **Step 3: 实现完成事务逻辑**

```go
func CompleteInventoryDoc(ctx context.Context, id uint64, operatorID uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var doc wmsmodel.InventoryDoc
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&doc).Error; err != nil {
			return err
		}
		if err := ensureDraftEditable(doc.Status); err != nil {
			return err
		}
		var items []wmsmodel.InventoryDocItem
		if err := tx.Where("doc_id = ?", doc.ID).Find(&items).Error; err != nil {
			return err
		}
		for _, item := range items {
			var stock wmsmodel.InventoryStock
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("warehouse_id = ? AND sku_id = ?", doc.WarehouseID, item.SkuID).
				First(&stock).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) && doc.DocType == wmsmodel.DocTypeInbound {
					stock = wmsmodel.InventoryStock{WarehouseID: doc.WarehouseID, SkuID: item.SkuID, AvailableQty: 0}
					if err := tx.Create(&stock).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}

			before := stock.AvailableQty
			after := before
			direction := "in"
			if doc.DocType == wmsmodel.DocTypeInbound {
				after = before + item.Qty
			} else {
				if before < item.Qty {
					return fmt.Errorf("库存不足 sku=%d current=%d need=%d", item.SkuID, before, item.Qty)
				}
				after = before - item.Qty
				direction = "out"
			}

			if err := tx.Model(&stock).Update("available_qty", after).Error; err != nil {
				return err
			}
			mv := wmsmodel.InventoryMovement{
				WarehouseID: doc.WarehouseID, SkuID: item.SkuID, Direction: direction, Qty: item.Qty,
				BeforeQty: before, AfterQty: after, BizType: doc.DocType + "_doc", BizDocID: doc.ID,
				BizDocNo: doc.DocNo, OperatorID: operatorID,
			}
			if err := tx.Create(&mv).Error; err != nil {
				return err
			}
		}
		return tx.Model(&doc).Updates(map[string]any{
			"status":       wmsmodel.DocStatusCompleted,
			"operator_id":  operatorID,
			"completed_at": time.Now().Unix(),
		}).Error
	})
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `cd server && go test ./plugins/wms/service -run TestCompleteOutboundDoc_UpdatesStockAndMovement -v`  
Expected: PASS。

- [ ] **Step 5: 提交**

```bash
git add server/plugins/wms/service/inventory_doc.go server/plugins/wms/service/inventory_movement.go server/plugins/wms/service/inventory_doc_complete_test.go server/go.mod server/go.sum
git commit -m "实现单据完成事务与库存流水引擎" -m "完成单据时原子更新库存快照并写入流水，加入并发锁与库存不足校验，确保库存一致性与可追溯性。"
```

---

### Task 6: 库存流水查询接口与后端回归

**Files:**
- Modify: `server/plugins/wms/service/inventory_movement.go`
- Modify: `server/plugins/wms/api/admin.go`
- Test: `server/plugins/wms/api/admin_test.go`

- [ ] **Step 1: 写失败测试（流水筛选接口）**

```go
func TestListMovements(t *testing.T) {
	r := gin.New()
	g := r.Group("/admin/api")
	RegisterAdminRoutes(g)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/api/wms/movements?page=1&size=20", nil)
	r.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatalf("movement route missing")
	}
}
```

- [ ] **Step 2: 运行测试，确认失败**

Run: `cd server && go test ./plugins/wms/api -run TestListMovements -v`  
Expected: FAIL，返回 404。

- [ ] **Step 3: 实现流水列表接口**

```go
func ListMovements(ctx context.Context, q MovementQuery) ([]wmsmodel.InventoryMovement, int64, error) {
	tx := db.DB.WithContext(ctx).Model(&wmsmodel.InventoryMovement{})
	if q.WarehouseID > 0 {
		tx = tx.Where("warehouse_id = ?", q.WarehouseID)
	}
	if q.SkuID > 0 {
		tx = tx.Where("sku_id = ?", q.SkuID)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []wmsmodel.InventoryMovement
	err := tx.Order("id desc").Offset((q.Page-1)*q.Size).Limit(q.Size).Find(&list).Error
	return list, total, err
}
```

- [ ] **Step 4: 运行后端测试集合**

Run: `cd server && go test ./plugins/wms/...`  
Expected: PASS，WMS 包内新增测试全部通过。

- [ ] **Step 5: 提交**

```bash
git add server/plugins/wms/service/inventory_movement.go server/plugins/wms/api/admin.go server/plugins/wms/api/admin_test.go
git commit -m "补齐库存流水查询接口" -m "新增库存流水分页筛选能力并完成路由接入，支持按仓库、SKU、业务单号等条件追溯库存变化。"
```

---

### Task 7: 管理后台接入新 WMS API 与页面路由

**Files:**
- Create: `admin/src/api/wms.ts`
- Modify: `admin/src/router/index.ts`
- Modify: `admin/src/locales/zh-CN.ts`
- Modify: `admin/src/locales/en.ts`

- [ ] **Step 1: 写失败检查（路由缺失）**

```ts
// admin/src/router/wms-route-check.ts
export const requiredWmsRoutes = ['/wms/warehouse', '/wms/docs', '/wms/stock-ledger', '/wms/movements']
```

- [ ] **Step 2: 运行类型构建，确认当前不满足**

Run: `cd admin && npm run build`  
Expected: FAIL 或页面路由不存在（手工校验菜单跳转出现空白页）。

- [ ] **Step 3: 实现 API 客户端与路由映射**

```ts
// admin/src/api/wms.ts
import request from './request'

export const listWarehouses = (params?: any) => request.get('/wms/warehouses', { params })
export const createWarehouse = (data: any) => request.post('/wms/warehouses', data)
export const updateWarehouse = (id: number, data: any) => request.put(`/wms/warehouses/${id}`, data)
export const updateWarehouseStatus = (id: number, status: number) => request.put(`/wms/warehouses/${id}/status`, { status })

export const listDocs = (params?: any) => request.get('/wms/docs', { params })
export const createDoc = (data: any) => request.post('/wms/docs', data)
export const updateDoc = (id: number, data: any) => request.put(`/wms/docs/${id}`, data)
export const completeDoc = (id: number) => request.post(`/wms/docs/${id}/complete`)
export const cancelDoc = (id: number) => request.post(`/wms/docs/${id}/cancel`)

export const listStocks = (params?: any) => request.get('/wms/stocks', { params })
export const updateSafetyQty = (id: number, safety_qty: number) => request.put(`/wms/stocks/${id}/safety`, { safety_qty })
export const listMovements = (params?: any) => request.get('/wms/movements', { params })
```

- [ ] **Step 4: 运行前端构建确认通过**

Run: `cd admin && npm run build`  
Expected: PASS，`vue-tsc` 与 `vite build` 通过。

- [ ] **Step 5: 提交**

```bash
git add admin/src/api/wms.ts admin/src/router/index.ts admin/src/locales/zh-CN.ts admin/src/locales/en.ts
git commit -m "后台接入WMS新接口与路由骨架" -m "新增仓库、单据、台账、流水的API客户端与路由入口，补齐中英文文案键，打通页面接入基础。"
```

---

### Task 8: 实现后台四个 WMS 页面

**Files:**
- Create: `admin/src/views/wms/WarehouseList.vue`
- Create: `admin/src/views/wms/DocList.vue`
- Create: `admin/src/views/wms/DocEditor.vue`
- Create: `admin/src/views/wms/StockLedger.vue`
- Create: `admin/src/views/wms/MovementList.vue`
- Modify: `admin/src/views/wms/StockList.vue`（删除旧内容或转发到新页面）

- [ ] **Step 1: 写失败检查（单据页多 SKU 明细）**

```bash
cd admin
rg -n "DocEditor|WarehouseList|MovementList|StockLedger" src/views/wms
```

- [ ] **Step 2: 运行构建并手工确认失败**

Run: `cd admin && npm run build`  
Expected: `rg` 无匹配（FAIL），说明四个页面文件尚未落地。

- [ ] **Step 3: 实现页面能力**

```vue
<!-- DocEditor.vue 核心：多 SKU 明细 -->
<tr v-for="(item, idx) in form.items" :key="idx">
  <td><input v-model.number="item.sku_id" type="number" /></td>
  <td><input v-model.number="item.qty" type="number" min="1" /></td>
  <td><button @click="removeItem(idx)">删除</button></td>
</tr>
<button @click="addItem">+ 添加明细</button>
```

```ts
function addItem() {
  form.items.push({ sku_id: 0, qty: 1, remark: '' })
}
```

- [ ] **Step 4: 构建与冒烟**

Run: `cd admin && npm run build`  
Expected: PASS。  
Run: 本地启动 `npm run dev`，手工验证：
- 仓库可新增/编辑/启停
- 出入库单可创建草稿、多 SKU 保存、完成/作废
- 库存台账可筛选并更新安全库存
- 流水可查询并展示单号

- [ ] **Step 5: 提交**

```bash
git add admin/src/views/wms
git commit -m "落地后台WMS四大页面闭环" -m "实现仓库管理、出入库单、库存台账、库存流水页面及核心交互，支持多SKU单据和完成后只读审计。"
```

---

### Task 9: 同步 docs-site 文档（功能+接口+部署影响）

**Files:**
- Create: `docs-site/docs/api/wms.md`
- Modify: `docs-site/docs/api/index.md`
- Modify: `docs-site/docs/guide/features.md`
- Modify: `docs-site/docs/.vitepress/config.mts`

- [ ] **Step 1: 写失败检查（文档缺 WMS API）**

```bash
cd docs-site
rg -n "/api/wms|仓储管理 `wms`" docs/api docs/guide
```

Expected: 无 `docs/api/wms.md`，索引缺仓储 API 入口。

- [ ] **Step 2: 运行检查命令确认失败**

Run: `cd docs-site && npm run docs:build`  
Expected: 文档可构建，但缺少 WMS API 页面与侧边栏入口（功能不满足）。

- [ ] **Step 3: 写文档实现**

```md
## 仓储管理 `wms`

- 仓库管理：列表、新增、编辑、启停
- 出入库单：草稿、完成、作废，多 SKU 明细
- 库存台账：可用库存/安全库存/预警
- 库存流水：按仓库、SKU、业务单号追溯

接口变化：新增 `/admin/api/wms/docs*`、`/admin/api/wms/movements`、`/admin/api/wms/stocks/:id/safety`。
部署/配置影响：无新增依赖与配置项。
```

- [ ] **Step 4: 构建文档确认通过**

Run: `cd docs-site && npm run docs:build`  
Expected: PASS，侧边栏可访问 `/api/wms`。

- [ ] **Step 5: 提交**

```bash
git add docs-site/docs/api/wms.md docs-site/docs/api/index.md docs-site/docs/guide/features.md docs-site/docs/.vitepress/config.mts
git commit -m "同步WMS重构文档到docs-site" -m "补充仓储最新架构功能说明与接口文档，新增WMS API页面及导航入口，并明确部署配置影响。"
```

---

### Task 10: 全链路回归与发布前检查

**Files:**
- Modify: `docs/superpowers/plans/2026-05-28-wms-inventory-engine-implementation.md`（打勾执行记录，可选）

- [ ] **Step 1: 后端测试回归**

Run: `cd server && go test ./plugins/wms/...`  
Expected: PASS。

- [ ] **Step 2: 前端构建回归**

Run: `cd admin && npm run build`  
Expected: PASS。

- [ ] **Step 3: 文档构建回归**

Run: `cd docs-site && npm run docs:build`  
Expected: PASS。

- [ ] **Step 4: 手工业务回归**

```text
1) 新建启用仓库 -> 新建入库单(2个SKU) -> 完成 -> 台账库存增加 -> 流水出现2条in
2) 新建出库单(同2个SKU) -> 完成 -> 台账库存减少 -> 流水出现2条out
3) 出库数量大于库存 -> 完成失败且库存不变
4) 已完成单据再次编辑/完成 -> 返回状态错误
```

- [ ] **Step 5: 提交**

```bash
git add .
git commit -m "完成WMS库存交易引擎全链路回归" -m "执行后端测试、前端构建、文档构建与关键业务冒烟回归，确认仓储管理闭环达到最小可用产品标准。"
```

---

## Self-Review

### 1. Spec coverage

- 架构与模型：Task 1、Task 5 覆盖
- 仓库管理：Task 2 覆盖
- 库存台账：Task 3 覆盖
- 出入库单（草稿/完成/作废，多 SKU）：Task 4、Task 5、Task 8 覆盖
- 库存流水查询：Task 6、Task 8 覆盖
- docs-site 同步：Task 9 覆盖
- 验收与回归：Task 10 覆盖

### 2. Placeholder scan

- 已检查，计划中无 `TODO/TBD/implement later/similar to` 类占位描述。

### 3. Type consistency

- 单据类型统一为 `inbound/outbound`
- 单据状态统一为 `draft/completed/canceled`
- 库存数量字段统一为 `available_qty/safety_qty`
- 流水方向统一为 `in/out`
