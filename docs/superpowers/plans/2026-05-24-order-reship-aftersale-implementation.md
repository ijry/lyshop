# 订单补发物流与售后流程增强 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在不破坏现有订单链路的前提下，完成订单补发物流轨迹、退货/换货标准售后流程及 `admin + web + app` 三端展示与操作闭环。

**Architecture:** 以 `order` 插件为中心，新增 `order_shipments` 与 `after_sale_*` 数据模型，并升级现有 `orders`/`ship` 接口语义；新增长动作接口用于申请、审核、收货、退款与完结。查询层聚合物流与售后摘要，前端增量接入新字段与新动作。

**Tech Stack:** Go 1.22 · Gin · GORM · Vue3 + TailwindCSS（admin/web）· uni-app（app）

---

## 文件结构与责任边界

### 后端

- Modify: `server/plugins/order/model/order.go`
  - 扩展 `OrderRefund`（关联售后单）
- Create: `server/plugins/order/model/shipment.go`
  - `OrderShipment` 模型与常量
- Create: `server/plugins/order/model/after_sale.go`
  - `AfterSaleCase` / `AfterSaleCaseItem` / `AfterSaleLog` 模型与状态枚举
- Modify: `server/plugins/order/plugin.go`
  - AutoMigrate 注册新增模型
- Modify: `server/plugins/order/service/order.go`
  - 升级 `ShipOrder` 与订单查询聚合结构
- Create: `server/plugins/order/service/shipment.go`
  - 物流轨迹写入/查询封装
- Create: `server/plugins/order/service/after_sale.go`
  - 售后申请、审核、回寄、收货、退款、完结、关闭核心逻辑
- Modify: `server/plugins/order/api/admin.go`
  - 升级 `ship` 入参 + 新增 admin 售后接口
- Modify: `server/plugins/order/api/front.go`
  - 新增前台售后接口

### 管理端

- Modify: `admin/src/api/plugins.ts`
  - 售后相关 API 封装
- Modify: `admin/src/views/order/OrderList.vue`
  - 列表展示售后标识与最近物流
- Modify: `admin/src/views/order/OrderDetail.vue`
  - 展示物流轨迹 + 补发操作
- Create: `admin/src/views/order/AfterSaleList.vue`
  - 售后列表页
- Create: `admin/src/views/order/AfterSaleDetail.vue`
  - 售后详情操作页
- Modify: `admin/src/router/index.ts`
  - 注册售后管理路由
- Modify: `server/plugins/order/plugin.json`
  - 菜单与权限补充（若采用插件菜单驱动）

### Web 前台

- Modify: `web/src/views/OrderDetail.vue`
  - 物流时间线与售后入口
- Create: `web/src/views/AfterSaleApply.vue`
  - 售后申请页
- Create: `web/src/views/AfterSaleDetail.vue`
  - 售后详情页
- Modify: `web/src/router/index.ts`
  - 新增售后路由

### App 前台

- Modify: `app/pages/order/detail.vue`
  - 物流时间线与售后入口
- Create: `app/pages/order/after-sale-apply.vue`
  - 售后申请页
- Create: `app/pages/order/after-sale-detail.vue`
  - 售后详情页
- Modify: `app/pages.json`
  - 注册页面路径

### 文档

- Modify: `docs-site/docs/api/order.md`
  - 接口升级与新增说明
- Modify: `docs-site/docs/guide/features.md`
  - 功能说明更新

---

### Task 1: 后端模型迁移与状态枚举落地

**Files:**
- Create: `server/plugins/order/model/shipment.go`
- Create: `server/plugins/order/model/after_sale.go`
- Modify: `server/plugins/order/model/order.go`
- Modify: `server/plugins/order/plugin.go`
- Test: `server/plugins/order/service/order_test.go`（若无则创建）

- [ ] **Step 1: Write the failing test**

```go
func TestOrderPluginMigrate_AfterSaleTablesExist(t *testing.T) {
	// arrange
	db := setupTestDB(t)
	p := &orderPlugin{}

	// act
	err := p.Migrate(db)
	require.NoError(t, err)

	// assert
	for _, table := range []string{"order_shipments", "after_sale_cases", "after_sale_case_items", "after_sale_logs"} {
		require.True(t, db.Migrator().HasTable(table), table)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server; go test ./plugins/order/... -run TestOrderPluginMigrate_AfterSaleTablesExist -v`
Expected: FAIL（缺少新增模型与迁移注册）

- [ ] **Step 3: Write minimal implementation**

```go
// shipment.go
package model

type ShipmentDirection string
const (
	ShipmentDirectionOutbound ShipmentDirection = "outbound"
	ShipmentDirectionInbound  ShipmentDirection = "inbound"
)

type ShipmentBizType string
const (
	ShipmentBizTypeInitial ShipmentBizType = "initial"
	ShipmentBizTypeReship  ShipmentBizType = "reship"
	ShipmentBizTypeReturn  ShipmentBizType = "return"
)

type OrderShipment struct {
	model.Base
	OrderID          uint64 `gorm:"not null;index"`
	AfterSaleCaseID  uint64 `gorm:"index"`
	Direction        string `gorm:"size:16;not null;index"`
	BizType          string `gorm:"size:16;not null;index"`
	Company          string `gorm:"size:64"`
	TrackingNo       string `gorm:"size:128;index"`
	LogisticsStatus  string `gorm:"size:32;not null;default:'pending'"`
	Remark           string `gorm:"size:255"`
	ShippedAt        *time.Time
	SignedAt         *time.Time
	CreatedByType    string `gorm:"size:16;not null;default:'admin'"`
	CreatedByID      uint64 `gorm:"not null;default:0"`
}
```

```go
// plugin.go (excerpt)
return db.AutoMigrate(
	&ordermodel.Address{},
	&ordermodel.Order{},
	&ordermodel.OrderItem{},
	&ordermodel.OrderPayment{},
	&ordermodel.OrderRefund{},
	&ordermodel.OrderReview{},
	&ordermodel.OrderReviewAppend{},
	&ordermodel.OrderReviewReply{},
	&ordermodel.OrderShipment{},
	&ordermodel.AfterSaleCase{},
	&ordermodel.AfterSaleCaseItem{},
	&ordermodel.AfterSaleLog{},
)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server; go test ./plugins/order/... -run TestOrderPluginMigrate_AfterSaleTablesExist -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/plugins/order/model/order.go server/plugins/order/model/shipment.go server/plugins/order/model/after_sale.go server/plugins/order/plugin.go
git commit -m "订单插件新增物流轨迹与售后模型\n\n补充 order_shipments、after_sale_cases、after_sale_case_items、after_sale_logs 数据结构并纳入迁移；为后续补发物流和退换货流程提供状态承载。影响范围为订单插件模型层与数据库迁移。"
```

---

### Task 2: 升级发货逻辑为“首发/补发 + 轨迹写入”

**Files:**
- Modify: `server/plugins/order/service/order.go`
- Create: `server/plugins/order/service/shipment.go`
- Modify: `server/plugins/order/api/admin.go`
- Test: `server/plugins/order/service/order_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestShipOrder_CreateShipmentAndKeepCompatibility(t *testing.T) {
	ctx, db := setupOrderServiceCase(t)
	order := seedPaidOrder(t, db)

	err := ShipOrder(ctx, order.ID, ShipOrderReq{TrackingNo: "SF100", ShipType: "initial"})
	require.NoError(t, err)

	err = ShipOrder(ctx, order.ID, ShipOrderReq{TrackingNo: "SF200", ShipType: "reship"})
	require.NoError(t, err)

	var shipments []ordermodel.OrderShipment
	require.NoError(t, db.Where("order_id = ?", order.ID).Order("id asc").Find(&shipments).Error)
	require.Len(t, shipments, 2)
	require.Equal(t, "initial", shipments[0].BizType)
	require.Equal(t, "reship", shipments[1].BizType)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server; go test ./plugins/order/service -run TestShipOrder_CreateShipmentAndKeepCompatibility -v`
Expected: FAIL（当前仅更新 `orders.tracking_no`）

- [ ] **Step 3: Write minimal implementation**

```go
// admin.go request
var req struct {
	TrackingNo       string `json:"tracking_no"`
	ShipType         string `json:"ship_type"`
	AfterSaleCaseID  uint64 `json:"after_sale_case_id"`
	Company          string `json:"company"`
	Remark           string `json:"remark"`
}

// order.go
func ShipOrder(ctx context.Context, orderID uint64, req ShipOrderReq) error {
	if req.ShipType == "" { req.ShipType = "initial" }
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := createShipmentTx(tx, orderID, req); err != nil { return err }
		updates := map[string]any{"tracking_no": req.TrackingNo}
		if req.ShipType == "initial" { updates["status"] = ordermodel.OrderStatusShipped }
		return tx.Model(&ordermodel.Order{}).Where("id = ?", orderID).Updates(updates).Error
	})
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server; go test ./plugins/order/service -run TestShipOrder_CreateShipmentAndKeepCompatibility -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/plugins/order/service/order.go server/plugins/order/service/shipment.go server/plugins/order/api/admin.go
git commit -m "订单发货升级支持补发与物流轨迹\n\n在现有 /orders/:id/ship 语义内扩展 ship_type 与补发参数，写入 order_shipments 并保持 tracking_no 兼容；首发行为保持原有状态流转。影响范围为订单服务与后台发货接口。"
```

---

### Task 3: 实现售后核心服务（申请/审核/回寄/收货/退款/完结）

**Files:**
- Create: `server/plugins/order/service/after_sale.go`
- Modify: `server/plugins/order/service/order.go`
- Modify: `server/plugins/order/model/order.go`
- Test: `server/plugins/order/service/after_sale_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestAfterSale_ReturnFlow_HappyPath(t *testing.T) {
	ctx, env := setupAfterSaleCase(t)
	caseID, err := CreateAfterSale(ctx, CreateAfterSaleReq{OrderID: env.OrderID, UserID: env.UserID, CaseType: "return", Items: []AfterSaleItemReq{{OrderItemID: env.OrderItemID, Qty: 1}}, Reason: "尺寸不合适"})
	require.NoError(t, err)

	require.NoError(t, AuditAfterSale(ctx, caseID, AuditAfterSaleReq{Approve: true}))
	require.NoError(t, SubmitReturnShipment(ctx, caseID, SubmitReturnShipmentReq{Company: "SF", TrackingNo: "RT100"}))
	require.NoError(t, ReceiveAfterSale(ctx, caseID))
	require.NoError(t, MarkRefund(ctx, caseID, MarkRefundReq{Amount: 100}))
	require.NoError(t, CompleteAfterSale(ctx, caseID))

	got, err := GetAfterSale(ctx, caseID)
	require.NoError(t, err)
	require.Equal(t, "completed", got.Status)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server; go test ./plugins/order/service -run TestAfterSale_ReturnFlow_HappyPath -v`
Expected: FAIL（售后服务缺失）

- [ ] **Step 3: Write minimal implementation**

```go
func CreateAfterSale(ctx context.Context, req CreateAfterSaleReq) (uint64, error) {
	// validate order/user/status
	// create after_sale_cases + items + log(applied)
}

func AuditAfterSale(ctx context.Context, caseID uint64, req AuditAfterSaleReq) error {
	// transition applied -> approved_wait_user_return / rejected
}

func SubmitReturnShipment(ctx context.Context, caseID uint64, req SubmitReturnShipmentReq) error {
	// create inbound shipment + transition to user_returning
}

func ReceiveAfterSale(ctx context.Context, caseID uint64) error {
	// transition to warehouse_received then refund_pending/reship_pending
}

func MarkRefund(ctx context.Context, caseID uint64, req MarkRefundReq) error {
	// write order_refunds with after_sale_case_id + transition refunded
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server; go test ./plugins/order/service -run TestAfterSale_ReturnFlow_HappyPath -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/plugins/order/service/after_sale.go server/plugins/order/service/order.go server/plugins/order/model/order.go server/plugins/order/service/after_sale_test.go
git commit -m "订单插件新增售后核心状态机与退款记录关联\n\n实现退货/换货申请、审核、回寄、收货、退款、完结流程，并写入 after_sale_logs；扩展 order_refunds 关联售后单，形成业务闭环。影响范围为订单服务层与退款台账模型。"
```

---

### Task 4: 开放前后台售后接口并升级订单查询返回

**Files:**
- Modify: `server/plugins/order/api/front.go`
- Modify: `server/plugins/order/api/admin.go`
- Modify: `server/plugins/order/service/order.go`
- Test: `server/plugins/order/api/order_api_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestOrderDetailResponse_ContainsShipmentsAndAfterSaleSummary(t *testing.T) {
	r := setupRouter(t)
	w := performGet(t, r, "/api/v1/orders/1", userToken)
	require.Equal(t, 200, w.Code)
	body := decodeJSON(t, w.Body)
	require.NotNil(t, body["data"].(map[string]any)["shipments"])
	require.NotNil(t, body["data"].(map[string]any)["after_sale_summary"])
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server; go test ./plugins/order/api -run TestOrderDetailResponse_ContainsShipmentsAndAfterSaleSummary -v`
Expected: FAIL（返回结构缺字段）

- [ ] **Step 3: Write minimal implementation**

```go
// order view
 type OrderView struct {
	ordermodel.Order
	Items             []OrderItemView      `json:"items"`
	AmountBreakdown   AmountBreakdown      `json:"amount_breakdown"`
	Shipments         []ordermodel.OrderShipment `json:"shipments"`
	AfterSaleSummary  *AfterSaleSummary    `json:"after_sale_summary,omitempty"`
	LatestShipment    *ordermodel.OrderShipment `json:"latest_shipment,omitempty"`
 }

// front/admin routes add after-sale endpoints
auth.POST("/orders/:id/after-sales", createAfterSale)
auth.GET("/after-sales/:id", getAfterSale)
auth.POST("/after-sales/:id/return-shipments", submitReturnShipment)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server; go test ./plugins/order/api -run TestOrderDetailResponse_ContainsShipmentsAndAfterSaleSummary -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/plugins/order/api/front.go server/plugins/order/api/admin.go server/plugins/order/service/order.go server/plugins/order/api/order_api_test.go
git commit -m "订单接口升级返回物流轨迹与售后摘要\n\n在前后台订单列表与详情中聚合 shipments/after_sale_summary/latest_shipment；新增用户售后申请与回寄接口以及后台售后动作接口。影响范围为订单 API 层与查询聚合结构。"
```

---

### Task 5: 管理端订单页接入补发与售后管理页

**Files:**
- Modify: `admin/src/api/plugins.ts`
- Modify: `admin/src/views/order/OrderList.vue`
- Modify: `admin/src/views/order/OrderDetail.vue`
- Create: `admin/src/views/order/AfterSaleList.vue`
- Create: `admin/src/views/order/AfterSaleDetail.vue`
- Modify: `admin/src/router/index.ts`

- [ ] **Step 1: Write the failing test**

```ts
it('renders shipment timeline and after-sale badge in order detail', async () => {
  const wrapper = mount(OrderDetail, { global: mockOrderDetailWithShipments })
  expect(wrapper.text()).toContain('物流轨迹')
  expect(wrapper.text()).toContain('补发')
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd admin; npm run test -- OrderDetail.spec.ts`
Expected: FAIL（页面未接入新字段）

- [ ] **Step 3: Write minimal implementation**

```ts
// plugins.ts
export const getAfterSales = (params?: any) => request.get('/after-sales', { params })
export const getAfterSaleDetail = (id: number) => request.get(`/after-sales/${id}`)
export const auditAfterSale = (id: number, data: any) => request.post(`/after-sales/${id}/audit`, data)

// OrderDetail.vue (excerpt)
<div v-if="detail.shipments?.length" class="...">
  <h3>物流轨迹</h3>
  <div v-for="s in detail.shipments" :key="s.id">{{ s.biz_type }} / {{ s.tracking_no }}</div>
</div>
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd admin; npm run test -- OrderDetail.spec.ts`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add admin/src/api/plugins.ts admin/src/views/order/OrderList.vue admin/src/views/order/OrderDetail.vue admin/src/views/order/AfterSaleList.vue admin/src/views/order/AfterSaleDetail.vue admin/src/router/index.ts
git commit -m "后台订单页接入补发轨迹并新增售后管理页\n\n订单列表与详情展示售后标记和物流时间线，支持补发动作；新增售后列表与详情处理页。影响范围为后台订单运营链路。"
```

---

### Task 6: Web 端订单详情接入售后申请与进度

**Files:**
- Modify: `web/src/views/OrderDetail.vue`
- Create: `web/src/views/AfterSaleApply.vue`
- Create: `web/src/views/AfterSaleDetail.vue`
- Modify: `web/src/router/index.ts`
- Modify: `web/src/api/request.ts`（若需统一错误映射）

- [ ] **Step 1: Write the failing test**

```ts
it('shows apply after-sale button when order item can apply', async () => {
  const wrapper = mount(OrderDetail, { global: mockOrderDetailWithAfterSaleSummary })
  expect(wrapper.text()).toContain('申请售后')
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd web; npm run test -- OrderDetail.spec.ts`
Expected: FAIL（按钮和路由不存在）

- [ ] **Step 3: Write minimal implementation**

```ts
// router
{ path: '/orders/:id/after-sale/apply', component: () => import('@/views/AfterSaleApply.vue') }
{ path: '/after-sales/:id', component: () => import('@/views/AfterSaleDetail.vue') }

// OrderDetail.vue
<button v-if="detail.after_sale_summary?.can_apply" @click="goApply">申请售后</button>
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd web; npm run test -- OrderDetail.spec.ts`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add web/src/views/OrderDetail.vue web/src/views/AfterSaleApply.vue web/src/views/AfterSaleDetail.vue web/src/router/index.ts
git commit -m "Web订单详情新增售后申请与进度页面\n\n接入订单物流轨迹与售后摘要，支持发起退换货申请并查看售后进度。影响范围为 Web 订单详情与售后交互页。"
```

---

### Task 7: App 端订单详情接入售后申请与进度

**Files:**
- Modify: `app/pages/order/detail.vue`
- Create: `app/pages/order/after-sale-apply.vue`
- Create: `app/pages/order/after-sale-detail.vue`
- Modify: `app/pages.json`

- [ ] **Step 1: Write the failing test**

```ts
it('navigates to after-sale apply page from order detail', async () => {
  const page = mountOrderDetailWithSummary()
  await page.find('[data-test="apply-after-sale"]').trigger('click')
  expect(mockNavigateTo).toHaveBeenCalledWith(expect.objectContaining({ url: expect.stringContaining('/pages/order/after-sale-apply') }))
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd app; npm run test -- order-detail.spec.ts`
Expected: FAIL（页面与路由不存在）

- [ ] **Step 3: Write minimal implementation**

```vue
<!-- detail.vue -->
<u-button v-if="detail.after_sale_summary?.can_apply" data-test="apply-after-sale" text="申请售后" @click="goAfterSaleApply" />
```

```json
// pages.json
{ "path": "pages/order/after-sale-apply", "style": { "navigationBarTitleText": "申请售后" } },
{ "path": "pages/order/after-sale-detail", "style": { "navigationBarTitleText": "售后详情" } }
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd app; npm run test -- order-detail.spec.ts`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add app/pages/order/detail.vue app/pages/order/after-sale-apply.vue app/pages/order/after-sale-detail.vue app/pages.json
git commit -m "App订单详情接入售后申请与进度展示\n\n新增售后申请页与售后详情页，订单详情展示物流轨迹并提供申请入口。影响范围为 uni-app 订单链路。"
```

---

### Task 8: 同步 docs-site 文档与全链路回归

**Files:**
- Modify: `docs-site/docs/api/order.md`
- Modify: `docs-site/docs/guide/features.md`

- [ ] **Step 1: Write the failing doc check**

```bash
rg -n "after-sales|补发|物流轨迹|退货|换货" docs-site/docs/api/order.md docs-site/docs/guide/features.md
```

Expected: 关键字缺失或不完整

- [ ] **Step 2: Update docs with concrete API changes**

```md
- 升级 `PUT /admin/api/orders/:id/ship`：新增 `ship_type/company/remark/after_sale_case_id`
- 新增 `POST /api/v1/orders/:id/after-sales`
- 新增 `GET /api/v1/after-sales/:id`
- 新增 `POST /admin/api/after-sales/:id/audit|receive|refund|complete|close`
```

- [ ] **Step 3: Run targeted verification**

Run: `rg -n "after-sales|ship_type|补发|物流轨迹|退货|换货" docs-site/docs/api/order.md docs-site/docs/guide/features.md`
Expected: 命中全部关键项

- [ ] **Step 4: Run core regression suites**

Run: 
- `cd server; go test ./plugins/order/...`
- `cd admin; npm run build`
- `cd web; npm run build`
- `cd app; npm run build:h5`

Expected: 全部通过；如存在与本任务无关的历史失败，记录并隔离说明。

- [ ] **Step 5: Commit**

```bash
git add docs-site/docs/api/order.md docs-site/docs/guide/features.md
git commit -m "同步订单补发与售后流程文档\n\n更新订单 API 与功能说明，覆盖接口升级、新增售后动作、兼容策略及多端展示影响，确保功能变更与文档一致。"
```

---

## 交付检查清单

- [ ] 补发物流可在后台创建并保留历史轨迹
- [ ] `web + app` 订单详情显示完整 `shipments`
- [ ] 用户可发起退货/换货并查看进度
- [ ] 后台可完成审核、收货、退款、补发、完结闭环
- [ ] 订单列表与详情返回兼容旧字段 `tracking_no`
- [ ] `docs-site` 完成功能说明、接口变化、影响范围更新

