# SKU-WMS 一体化库存改造 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现“SKU 自动笛卡尔积 + 旧 SKU 软删除 + WMS 库存真源（预占/确认/释放）+ 单仓优先”闭环，消除 `product_skus.stock` 直接扣减路径。

**Architecture:** 商品侧由后端规格引擎统一生成和维护 SKU（以 `sku_key` 稳定映射），订单库存全部下沉到 WMS 交易能力。下单先预占，支付确认再扣减在手库存，取消或超时释放预占。前台/后台 API 尽量兼容升级，避免新增外部接口。

**Tech Stack:** Go, Gin, GORM, MySQL, Vue3(文档与接口描述), VitePress(docs-site)

---

## File Structure

**Create:**
- `server/plugins/product/service/sku_engine.go`（规格组合生成、`sku_key` 计算、差异结果）
- `server/plugins/product/service/sku_engine_test.go`（规格引擎单测）
- `server/plugins/wms/model/reservation.go`（预占记录模型）
- `server/plugins/wms/service/reservation.go`（预占/确认/释放服务）
- `server/plugins/wms/service/reservation_test.go`（WMS 预占服务单测）
- `server/plugins/order/service/order_wms_test.go`（订单与 WMS 库存链路集成测试）
- `docs-site/docs/api/stock-reservation.md`（库存交易规则页，供 order/wms/product 交叉引用）

**Modify:**
- `server/plugins/product/model/sku.go`
- `server/plugins/product/service/product.go`
- `server/plugins/product/api/admin.go`
- `server/plugins/wms/model/wms.go`
- `server/plugins/wms/plugin.go`
- `server/plugins/wms/service/doc.go`
- `server/plugins/wms/service/movement.go`
- `server/plugins/order/service/order.go`
- `server/plugins/order/api/front.go`
- `server/plugins/order/plugin.go`
- `docs-site/docs/api/product.md`
- `docs-site/docs/api/order.md`
- `docs-site/docs/api/wms.md`
- `docs-site/docs/guide/features.md`

**Test Commands:**
- `go test ./plugins/product/service -run TestSkuEngine -v`
- `go test ./plugins/wms/service -run Reservation -v`
- `go test ./plugins/order/service -run Wms -v`
- `go test ./plugins/order/service ./plugins/wms/service ./plugins/product/service -v`

---

### Task 1: 商品 SKU 模型与规格引擎

**Files:**
- Create: `server/plugins/product/service/sku_engine.go`
- Create: `server/plugins/product/service/sku_engine_test.go`
- Modify: `server/plugins/product/model/sku.go`
- Modify: `server/plugins/product/service/product.go`

- [ ] **Step 1: 写失败测试（组合生成与 `sku_key` 稳定性）**

```go
func TestSkuEngine_GenerateAndDiff(t *testing.T) {
    schema := []SpecGroup{{Name:"颜色", Values:[]string{"红","蓝"}}, {Name:"尺码", Values:[]string{"M","L"}}}
    existing := []productmodel.ProductSku{{SkuKey:"颜色:红|尺码:M"}, {SkuKey:"颜色:蓝|尺码:M"}}
    out := BuildSkuDiff(schema, existing)
    require.Len(t, out.Added, 2)
    require.Len(t, out.Kept, 2)
    require.Len(t, out.Inactivated, 0)
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./plugins/product/service -run TestSkuEngine_GenerateAndDiff -v`  
Expected: FAIL（`BuildSkuDiff` 未定义）

- [ ] **Step 3: 实现最小可用规格引擎**

```go
func CanonicalSkuKey(attrs []model.SkuAttr) string {
    sort.Slice(attrs, func(i, j int) bool { return attrs[i].Name < attrs[j].Name })
    parts := make([]string, 0, len(attrs))
    for _, a := range attrs { parts = append(parts, a.Name+":"+a.Value) }
    return strings.Join(parts, "|")
}
```

- [ ] **Step 4: 修改 `ProductSku` 与保存逻辑**

```go
type ProductSku struct {
    // ...
    SkuKey string `gorm:"size:255;not null;index:uk_product_sku_key,unique" json:"sku_key"`
    Status string `gorm:"size:16;not null;default:'active';index" json:"status"`
}
```

在 `ReplaceProductSkus` 中改为“差异化新增/更新/下线”，不再全删重建。

- [ ] **Step 5: 运行测试并提交**

Run: `go test ./plugins/product/service -run TestSkuEngine -v`  
Expected: PASS

```bash
git add server/plugins/product/model/sku.go server/plugins/product/service/sku_engine.go server/plugins/product/service/sku_engine_test.go server/plugins/product/service/product.go
git commit -m "实现SKU规格引擎与差异化持久化" -m "主要改动点：新增sku_key与状态字段，后端生成规格组合并按差异处理新增/保留/下线。\n\n原因：解决前端生成导致的SKU不稳定和全量替换问题。\n\n影响范围：product插件模型与商品保存链路。"
```

---

### Task 2: 商品管理接口兼容升级（自动笛卡尔积）

**Files:**
- Modify: `server/plugins/product/api/admin.go`
- Modify: `server/plugins/product/service/product.go`
- Test: `server/plugins/product/service/sku_engine_test.go`

- [ ] **Step 1: 增加接口请求结构测试样例（兼容老字段）**

```go
// 旧请求仅传 skus 仍可工作；新请求可传 spec_schema + sku_overrides
```

- [ ] **Step 2: 运行测试确认失败场景**

Run: `go test ./plugins/product/service -run TestProductUpdateSkuAutoMode -v`  
Expected: FAIL（新模式未接入）

- [ ] **Step 3: 接入 `sku_generation_mode=auto`**

```go
if req.Product["sku_generation_mode"] == "auto" {
    // 使用 spec_schema + sku_overrides 走后端引擎
}
```

- [ ] **Step 4: 回包增加 `sku_diff`**

```go
response.OK(c, gin.H{"id": id, "sku_diff": diff})
```

- [ ] **Step 5: 测试与提交**

Run: `go test ./plugins/product/service -v`  
Expected: PASS

```bash
git add server/plugins/product/api/admin.go server/plugins/product/service/product.go
git commit -m "升级商品保存接口支持自动SKU生成" -m "主要改动点：在兼容原有skus提交的基础上，新增自动生成模式并返回sku_diff摘要。\n\n原因：统一SKU生成入口并降低人工维护复杂度。\n\n影响范围：后台商品编辑保存接口。"
```

---

### Task 3: WMS 模型升级（预占能力）

**Files:**
- Create: `server/plugins/wms/model/reservation.go`
- Modify: `server/plugins/wms/model/wms.go`
- Modify: `server/plugins/wms/plugin.go`

- [ ] **Step 1: 写失败测试（迁移后字段与表存在）**

```go
func TestWmsMigrateIncludesReservation(t *testing.T) {
    // 断言 inventory_stock 有 reserved_qty，inventory_reservation 表可创建
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./plugins/wms/model -run Reservation -v`  
Expected: FAIL

- [ ] **Step 3: 添加模型字段与枚举**

```go
type InventoryStock struct {
    Qty         int `json:"qty"`
    ReservedQty int `gorm:"not null;default:0" json:"reserved_qty"`
}
```

```go
type InventoryReservation struct {
    BizType string `gorm:"size:32;not null;uniqueIndex:uk_biz_sku"`
    BizNo   string `gorm:"size:64;not null;uniqueIndex:uk_biz_sku"`
    SkuID   uint64 `gorm:"not null;uniqueIndex:uk_biz_sku"`
}
```

- [ ] **Step 4: 更新 `AutoMigrate` 注册**

```go
return db.AutoMigrate(..., &wmsmodel.InventoryReservation{})
```

- [ ] **Step 5: 测试与提交**

Run: `go test ./plugins/wms/model -v`  
Expected: PASS

```bash
git add server/plugins/wms/model/wms.go server/plugins/wms/model/reservation.go server/plugins/wms/plugin.go
git commit -m "扩展WMS模型支持库存预占" -m "主要改动点：库存表新增reserved_qty，新增预占记录模型并纳入自动迁移。\n\n原因：为订单预占/确认/释放提供持久化基础。\n\n影响范围：wms插件数据模型与建表流程。"
```

---

### Task 4: 实现 WMS 预占/确认/释放服务

**Files:**
- Create: `server/plugins/wms/service/reservation.go`
- Create: `server/plugins/wms/service/reservation_test.go`
- Modify: `server/plugins/wms/service/movement.go`
- Modify: `server/plugins/wms/service/errors.go`

- [ ] **Step 1: 写失败测试（预占、幂等确认、幂等释放）**

```go
func TestReserveConfirmRelease_Idempotent(t *testing.T) {
    // reserve -> confirm -> confirm(幂等) -> release(应拒绝或幂等按状态规则)
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./plugins/wms/service -run Reservation -v`  
Expected: FAIL

- [ ] **Step 3: 实现 `ReserveStock`（单事务+行锁）**

```go
if stock.Qty-stock.ReservedQty < req.Qty { return ConflictError("库存不足") }
stock.ReservedQty += req.Qty
```

- [ ] **Step 4: 实现 `ConfirmReservation` / `ReleaseReservation`**

```go
// confirm: qty -= reserved, reserved -= reserved
// release: reserved -= reserved
// 都写 movement(biz_type=reserve/release/order_outbound)
```

- [ ] **Step 5: 测试与提交**

Run: `go test ./plugins/wms/service -run Reservation -v`  
Expected: PASS

```bash
git add server/plugins/wms/service/reservation.go server/plugins/wms/service/reservation_test.go server/plugins/wms/service/movement.go server/plugins/wms/service/errors.go
git commit -m "实现WMS库存预占确认释放服务" -m "主要改动点：新增预占三段式交易服务与对应流水写入，补齐并发锁与幂等校验。\n\n原因：构建订单场景可防超卖的库存交易能力。\n\n影响范围：wms库存服务层。"
```

---

### Task 5: 订单创建改造为先预占库存

**Files:**
- Modify: `server/plugins/order/service/order.go`
- Create: `server/plugins/order/service/order_wms_test.go`

- [ ] **Step 1: 写失败测试（下单调用 WMS 预占而非扣 `product_skus.stock`）**

```go
func TestCreateOrder_ReserveByWms(t *testing.T) {
    // mock ReserveStock 被调用；product_skus.stock 不再直接 UpdateColumn
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./plugins/order/service -run ReserveByWms -v`  
Expected: FAIL

- [ ] **Step 3: 替换下单库存逻辑**

```go
// 删除：UpdateColumn("stock", stock-qty)
// 新增：wmssvc.ReserveStock(...)
```

- [ ] **Step 4: 记录订单与仓库绑定**

```go
items[i].WarehouseID = selectedWarehouseID // 需先在模型补字段
```

- [ ] **Step 5: 测试与提交**

Run: `go test ./plugins/order/service -run ReserveByWms -v`  
Expected: PASS

```bash
git add server/plugins/order/service/order.go server/plugins/order/service/order_wms_test.go server/plugins/order/model/order.go server/plugins/order/plugin.go
git commit -m "下单流程切换为WMS预占库存" -m "主要改动点：移除直接扣减SKU库存，改为订单创建阶段调用WMS预占并记录仓库快照。\n\n原因：统一库存真源并防止超卖。\n\n影响范围：order插件下单链路与订单行模型。"
```

---

### Task 6: 支付确认与取消释放接入

**Files:**
- Modify: `server/plugins/order/service/order.go`
- Modify: `server/plugins/order/api/front.go`
- Test: `server/plugins/order/service/order_wms_test.go`

- [ ] **Step 1: 写失败测试（支付后确认预占，取消释放预占）**

```go
func TestPayOrder_ConfirmReservation(t *testing.T) {}
func TestCancelOrder_ReleaseReservation(t *testing.T) {}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./plugins/order/service -run "ConfirmReservation|ReleaseReservation" -v`  
Expected: FAIL

- [ ] **Step 3: 在 `PayOrder` 中调用 `ConfirmReservation`**

```go
if err := wmssvc.ConfirmReservation(ctx, "order", order.OrderNo); err != nil { return err }
```

- [ ] **Step 4: 新增取消订单接口并释放预占**

```go
auth.POST("/orders/:id/cancel", cancelOrder)
```

```go
if err := wmssvc.ReleaseReservation(ctx, "order", order.OrderNo, "user_cancel"); err != nil { return err }
```

- [ ] **Step 5: 测试与提交**

Run: `go test ./plugins/order/service -run "ConfirmReservation|ReleaseReservation" -v`  
Expected: PASS

```bash
git add server/plugins/order/service/order.go server/plugins/order/api/front.go server/plugins/order/service/order_wms_test.go

git commit -m "接入支付确认与取消释放库存预占" -m "主要改动点：支付完成确认预占扣减，新增取消订单释放预占路径。\n\n原因：完成订单库存交易闭环。\n\n影响范围：order插件支付与订单状态流转。"
```

---

### Task 7: 过期预占释放机制

**Files:**
- Modify: `server/plugins/order/plugin.go`
- Modify: `server/plugins/order/service/order.go`
- Create: `server/plugins/order/service/reservation_reaper.go`
- Test: `server/plugins/order/service/order_wms_test.go`

- [ ] **Step 1: 写失败测试（过期待支付订单可释放）**

```go
func TestReleaseExpiredReservations(t *testing.T) {}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./plugins/order/service -run ReleaseExpiredReservations -v`  
Expected: FAIL

- [ ] **Step 3: 实现回收函数（先做可调用函数）**

```go
func ReleaseExpiredReservations(ctx context.Context, before time.Time) (int, error)
```

- [ ] **Step 4: 在 `Install()` 启动轻量 ticker（可配置间隔）**

```go
go ordersvc.StartReservationReaper(context.Background(), 30*time.Second)
```

- [ ] **Step 5: 测试与提交**

Run: `go test ./plugins/order/service -run ReleaseExpiredReservations -v`  
Expected: PASS

```bash
git add server/plugins/order/plugin.go server/plugins/order/service/reservation_reaper.go server/plugins/order/service/order.go server/plugins/order/service/order_wms_test.go
git commit -m "新增过期预占库存自动释放任务" -m "主要改动点：增加过期待支付订单预占释放函数与后台定时回收任务。\n\n原因：避免遗留预占长期占用可售库存。\n\n影响范围：order插件运行时任务与库存状态维护。"
```

---

### Task 8: docs-site 文档同步（最新架构口径）

**Files:**
- Create: `docs-site/docs/api/stock-reservation.md`
- Modify: `docs-site/docs/api/product.md`
- Modify: `docs-site/docs/api/order.md`
- Modify: `docs-site/docs/api/wms.md`
- Modify: `docs-site/docs/guide/features.md`
- Modify: `docs-site/docs/api/index.md`

- [ ] **Step 1: 写文档断言清单（功能/接口/部署）**

```md
- 功能说明：SKU自动生成、软删除、单仓优先
- 接口变化：product/order/wms
- 部署影响：新增字段、表、回收任务
```

- [ ] **Step 2: 更新 product/order/wms 接口文档**

Run: 编辑后检查 `docs-site/docs/api/*.md` 中相关段落。

- [ ] **Step 3: 新增库存交易规则页并加入索引**

```md
# 库存预占交易规则
Reserve -> Confirm -> Release
```

- [ ] **Step 4: 更新功能总览页**

Run: 修改 `docs-site/docs/guide/features.md` WMS 与订单库存描述。

- [ ] **Step 5: 本地构建校验并提交**

Run: `cd docs-site && npm run docs:build`  
Expected: BUILD SUCCESS

```bash
git add docs-site/docs/api/index.md docs-site/docs/api/product.md docs-site/docs/api/order.md docs-site/docs/api/wms.md docs-site/docs/api/stock-reservation.md docs-site/docs/guide/features.md
git commit -m "同步SKU与WMS一体化文档到docs-site" -m "主要改动点：补充库存交易规则页并更新商品、订单、仓储接口说明与部署影响。\n\n原因：保证文档默认描述最新接口与架构。\n\n影响范围：docs-site API与功能说明文档。"
```

---

### Task 9: 全链路回归与发布验收

**Files:**
- Modify: `server/plugins/order/service/order_wms_test.go`
- Modify: `server/plugins/wms/service/reservation_test.go`
- Modify: `server/plugins/product/service/sku_engine_test.go`
- Modify: `docs/superpowers/specs/2026-05-29-sku-wms-unified-design.md`（若验收结果有微调）

- [ ] **Step 1: 增加回归场景测试**

```go
// A. SKU变更后历史订单可读
// B. 下单-支付-发货库存一致
// C. 下单后取消释放可售库存
```

- [ ] **Step 2: 运行核心测试集**

Run: `go test ./plugins/product/service ./plugins/wms/service ./plugins/order/service -v`  
Expected: PASS

- [ ] **Step 3: 执行人工验收脚本**

Run:
1) 商品编辑应用规格模板并保存
2) 创建订单（待支付）验证 `reserved_qty` 增加
3) 支付订单验证 `qty` 减少且 `reserved_qty` 归零
4) 新建订单后取消验证释放

- [ ] **Step 4: 检查性能与并发**

Run: 针对同 SKU 并发下单压测（至少 50 并发），确认无负库存与重复确认。

- [ ] **Step 5: 最终提交**

```bash
git add server/plugins/product/service/sku_engine_test.go server/plugins/wms/service/reservation_test.go server/plugins/order/service/order_wms_test.go
git commit -m "完成SKU与WMS库存一体化回归验收" -m "主要改动点：补齐跨插件回归测试并完成并发与流程验收。\n\n原因：确保上线前库存一致性与幂等正确。\n\n影响范围：product/wms/order测试基线。"
```

---

## Self-Review

1. **Spec coverage:** 已覆盖 SKU 自动生成、软删除、WMS 预占三段式、单仓优先、文档同步、部署影响。
2. **Placeholder scan:** 计划未包含 TBD/TODO；每任务有文件、命令、期望结果。
3. **Type consistency:** 统一使用 `sku_key`、`reserved_qty`、`InventoryReservation`、`Reserve/Confirm/ReleaseReservation` 命名。

