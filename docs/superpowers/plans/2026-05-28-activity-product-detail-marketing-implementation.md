# 活动商品详情营销补充链路 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为 H5 与 PC 商品详情建立“活动来源驱动”的营销补充链路，并实现 `activity_product_id` 在购物车和下单链路的强约束透传。

**Architecture:** 保持标准商品详情接口不变，新增营销活动商品详情接口；活动列表进入详情时携带 `activity_product_id`，详情页按需发起二次营销请求。后端将购物车键升级为 `sku_id:activity_product_id` 复合键，并在下单时按来源 ID 做活动有效性、库存、限购和价格强校验。

**Tech Stack:** Go + Gin + GORM + Redis、Vue3（web）、uni-app（app）、TypeScript、VitePress docs-site。

---

## 文件结构与职责

- `server/plugins/marketing/api/front.go`：新增活动商品详情接口路由与 handler。
- `server/plugins/marketing/service/activity.go`：补充活动商品详情查询、来源校验、按活动商品维度累加已售。
- `server/plugins/order/model/order.go`：订单项扩展活动快照字段。
- `server/plugins/order/service/cart.go`：购物车复合键编码/解析、活动来源透传、购物车项结构升级。
- `server/plugins/order/service/order.go`：下单请求 `items[{sku_id,activity_product_id}]` 支持与活动强校验落单。
- `server/plugins/order/api/front.go`：购物车与下单接口参数兼容升级。
- `web/src/views/ActivityProductListBase.vue`：活动列表跳详情透传来源 ID。
- `web/src/views/ProductDetail.vue`：详情双请求、活动区块、加购透传。
- `web/src/views/Cart.vue`：购物车展示活动标签与活动来源分组结算。
- `app/pages/marketing/seckill.vue`、`app/pages/marketing/group-buy.vue`、`app/pages/marketing/bargain.vue`：透传来源 ID。
- `app/pages/product/detail.vue`：详情双请求、活动区块、加购/立即购买透传。
- `app/pages/cart/index.vue`、`app/pages/order/confirm.vue`：结算透传 `items[{sku_id,activity_product_id}]`。
- `app/mock/index.ts`、`web/src/mock/index.ts`：mock 支持新接口与新参数。
- `docs-site/docs/api/marketing.md`、`docs-site/docs/api/order.md`、`docs-site/docs/api/product.md`：功能说明、接口变化、部署影响同步。

---

### Task 1: 营销活动商品详情接口（后端）

**Files:**
- Modify: `server/plugins/marketing/api/front.go`
- Modify: `server/plugins/marketing/service/activity.go`
- Test: `server/plugins/marketing/service/activity_source_test.go`

- [ ] **Step 1: 先写失败测试（活动商品详情查询 + 来源校验）**

```go
func TestValidateActivityProductSource(t *testing.T) {
    // case1: activity_product_id=0 -> normal flow
    // case2: id exists but sku mismatch -> expect error
    // case3: id exists and valid -> expect detail
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./plugins/marketing/service -run TestValidateActivityProductSource -v`
Expected: FAIL，提示缺少 `ValidateActivityProductSource` 或断言不通过。

- [ ] **Step 3: 实现活动商品详情 DTO 与查询方法**

```go
type FrontActivityProductDetail struct {
    ActivityProductID uint64 `json:"activity_product_id"`
    ActivityID        uint64 `json:"activity_id"`
    ActivityType      string `json:"activity_type"`
    ActivityName      string `json:"activity_name"`
    ProductID         uint64 `json:"product_id"`
    SkuID             uint64 `json:"sku_id"`
    ActivityPrice     float64 `json:"activity_price"`
    StartPrice        float64 `json:"start_price"`
    FloorPrice        float64 `json:"floor_price"`
    LimitPerOrder     int `json:"limit_per_order"`
    TotalStockLimit   int `json:"total_stock_limit"`
    SoldQty           int `json:"sold_qty"`
    ActivityStartAt   *time.Time `json:"activity_start_at"`
    ActivityEndAt     *time.Time `json:"activity_end_at"`
}

func GetFrontActivityProductDetail(ctx context.Context, activityProductID uint64) (*FrontActivityProductDetail, error)
func ValidateActivityProductSource(ctx context.Context, activityProductID, skuID, productID uint64) (*FrontActivityProductDetail, error)
```

- [ ] **Step 4: 暴露前台接口路由**

```go
g.GET("/marketing/activity-products/:id", getFrontActivityProductDetail)
```

- [ ] **Step 5: 运行测试确认通过并提交**

Run: `go test ./plugins/marketing/service -run TestValidateActivityProductSource -v`
Expected: PASS。

```bash
git add server/plugins/marketing/api/front.go server/plugins/marketing/service/activity.go server/plugins/marketing/service/activity_source_test.go
git commit -m "新增营销活动商品详情与来源校验能力" -m "主要改动：新增活动商品详情查询与来源校验方法，并开放前台详情接口。\n\n原因：详情页需要按activity_product_id补充秒杀/拼团/砍价元素，且后续下单需强约束来源。\n\n影响范围：marketing前台接口与服务层。"
```

### Task 2: 活动列表返回补充 `activity_product_id`

**Files:**
- Modify: `server/plugins/marketing/service/activity.go`
- Modify: `app/mock/index.ts`
- Modify: `web/src/mock/index.ts`

- [ ] **Step 1: 先写失败测试（列表行应包含 activity_product_id）**

```go
func TestListFrontActivityProducts_HasActivityProductID(t *testing.T) {
    // arrange rows
    // assert list[0].ActivityProductID > 0
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./plugins/marketing/service -run TestListFrontActivityProducts_HasActivityProductID -v`
Expected: FAIL，字段缺失。

- [ ] **Step 3: 扩展返回结构并赋值 `ap.id`**

```go
type FrontActivityProduct struct {
    ActivityProductID uint64 `json:"activity_product_id"`
    // existing fields...
}

// SELECT 增加 ap.id AS activity_product_id
```

- [ ] **Step 4: 同步 app/web mock 列表数据**

```ts
activity_product_id: Number(item?.id || item?.activity_product_id || 0)
```

- [ ] **Step 5: 运行验证并提交**

Run:
- `go test ./plugins/marketing/service -run TestListFrontActivityProducts_HasActivityProductID -v`
- `npm run -s lint`（若仓库无 lint 命令则跳过）

Expected: Go test PASS，前端脚本无语法错误。

```bash
git add server/plugins/marketing/service/activity.go app/mock/index.ts web/src/mock/index.ts
git commit -m "活动列表返回补充活动商品唯一ID" -m "主要改动：活动商品列表增加activity_product_id并同步双端mock。\n\n原因：活动入口进入详情需要携带唯一来源ID。\n\n影响范围：marketing列表返回与前端mock数据。"
```

### Task 3: 购物车复合键与接口参数兼容升级

**Files:**
- Modify: `server/plugins/order/service/cart.go`
- Modify: `server/plugins/order/api/front.go`
- Test: `server/plugins/order/service/cart_field_test.go`

- [ ] **Step 1: 先写失败测试（复合键编码/解析）**

```go
func TestCartFieldEncodeDecode(t *testing.T) {
    field := formatCartField(101, 9001)
    skuID, apID, err := parseCartField(field)
    require.NoError(t, err)
    require.Equal(t, uint64(101), skuID)
    require.Equal(t, uint64(9001), apID)
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./plugins/order/service -run TestCartFieldEncodeDecode -v`
Expected: FAIL，缺少 `formatCartField/parseCartField`。

- [ ] **Step 3: 实现复合键与购物车结构升级**

```go
type CartItem struct {
    SkuID             uint64 `json:"sku_id"`
    ActivityProductID uint64 `json:"activity_product_id"`
    Qty               int    `json:"qty"`
    // Product, Sku...
}

func AddToCart(ctx context.Context, userID, skuID, activityProductID uint64, qty int) error
func UpdateCartQty(ctx context.Context, userID, skuID, activityProductID uint64, qty int) error
func RemoveFromCart(ctx context.Context, userID, skuID, activityProductID uint64) error
```

- [ ] **Step 4: 升级前台购物车 API 参数绑定**

```go
var req struct {
    SkuID             uint64 `json:"sku_id" binding:"required"`
    ActivityProductID uint64 `json:"activity_product_id"`
    Qty               int    `json:"qty"`
}
```

- [ ] **Step 5: 运行测试确认通过并提交**

Run: `go test ./plugins/order/service -run TestCartFieldEncodeDecode -v`
Expected: PASS。

```bash
git add server/plugins/order/service/cart.go server/plugins/order/api/front.go server/plugins/order/service/cart_field_test.go
git commit -m "升级购物车为活动来源复合键" -m "主要改动：购物车键升级为sku_id:activity_product_id并扩展购物车接口参数。\n\n原因：同SKU普通价与活动价需要并存，避免互相覆盖。\n\n影响范围：order购物车服务与前台购物车接口。"
```

### Task 4: 下单请求升级与活动强校验落单

**Files:**
- Modify: `server/plugins/order/model/order.go`
- Modify: `server/plugins/order/service/order.go`
- Modify: `server/plugins/order/api/front.go`
- Modify: `server/plugins/marketing/service/activity.go`
- Test: `server/plugins/order/service/order_activity_test.go`

- [ ] **Step 1: 先写失败测试（有来源ID时必须校验活动）**

```go
func TestCreateOrder_WithActivityItem_ValidateSource(t *testing.T) {
    // mock req.Items with activity_product_id
    // expect source mismatch returns error
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./plugins/order/service -run TestCreateOrder_WithActivityItem_ValidateSource -v`
Expected: FAIL。

- [ ] **Step 3: 扩展请求结构与订单项快照字段**

```go
type CreateOrderReq struct {
    UserID    uint64 `json:"user_id"`
    AddressID uint64 `json:"address_id"`
    Items     []struct {
        SkuID             uint64 `json:"sku_id"`
        ActivityProductID uint64 `json:"activity_product_id"`
    } `json:"items"`
    SkuIDs []uint64 `json:"sku_ids"`
    // ...
}

type OrderItem struct {
    ActivityProductID uint64 `json:"activity_product_id"`
    ActivityID        uint64 `json:"activity_id"`
    ActivityType      string `json:"activity_type"`
    ActivityTitle     string `json:"activity_title"`
}
```

- [ ] **Step 4: 实现创建订单时活动来源强校验与定价覆盖**

```go
if reqItem.ActivityProductID > 0 {
    detail, err := marketingsvc.ValidateActivityProductSource(ctx, reqItem.ActivityProductID, ci.SkuID, ci.Product.ID)
    if err != nil { return nil, err }
    item.ActivityProductID = detail.ActivityProductID
    item.ActivityID = detail.ActivityID
    item.ActivityType = detail.ActivityType
    item.ActivityTitle = detail.ActivityName
    item.Price = resolveActivityPrice(detail, ci.Sku.Price)
}
```

- [ ] **Step 5: 改为按活动商品维度累加已售并提交**

```go
func IncreaseSoldQtyByActivityProductTx(tx *gorm.DB, activityProductID uint64, qty int) error
```

Run: `go test ./plugins/order/service -run TestCreateOrder_WithActivityItem_ValidateSource -v`
Expected: PASS。

```bash
git add server/plugins/order/model/order.go server/plugins/order/service/order.go server/plugins/order/api/front.go server/plugins/marketing/service/activity.go server/plugins/order/service/order_activity_test.go
git commit -m "下单支持活动来源透传与强校验" -m "主要改动：下单请求新增items来源结构，订单项落活动快照，并按activity_product_id强校验和累加已售。\n\n原因：活动商品必须按来源ID保障活动价、限购和库存一致性。\n\n影响范围：order下单链路与marketing已售统计逻辑。"
```

### Task 5: H5 活动入口透传与详情营销区块

**Files:**
- Modify: `app/pages/marketing/seckill.vue`
- Modify: `app/pages/marketing/group-buy.vue`
- Modify: `app/pages/marketing/bargain.vue`
- Modify: `app/pages/product/detail.vue`
- Modify: `app/locales/zh-CN.ts`
- Modify: `app/locales/en.ts`

- [ ] **Step 1: 先写页面行为校验清单（手工失败用例）**

```md
1. 活动列表点击商品，URL 需包含 activity_product_id
2. 详情页有来源时需出现活动区块
3. 无来源时不显示活动区块
```

- [ ] **Step 2: 改造三类活动列表跳转参数**

```ts
uni.navigateTo({
  url: `/pages/product/detail?id=${p.product_id}&activity_product_id=${p.activity_product_id}`,
})
```

- [ ] **Step 3: 详情页实现双请求与活动区块状态**

```ts
const activityProductID = Number(query.activity_product_id || 0)
if (activityProductID > 0) {
  marketingDetail.value = await get(`/api/v1/marketing/activity-products/${activityProductID}`)
}
```

- [ ] **Step 4: 加购与立即购买透传来源 ID**

```ts
await post('/api/v1/cart/add', {
  sku_id: skuID,
  qty: 1,
  activity_product_id: currentActivityProductID.value,
})
uni.navigateTo({
  url: `/pages/order/confirm?items=${encodeURIComponent(JSON.stringify([{ sku_id: skuID, activity_product_id: currentActivityProductID.value }]))}`,
})
```

- [ ] **Step 5: 构建验证并提交**

Run: `npm run build:h5:demo`
Expected: PASS。

```bash
git add app/pages/marketing/seckill.vue app/pages/marketing/group-buy.vue app/pages/marketing/bargain.vue app/pages/product/detail.vue app/locales/zh-CN.ts app/locales/en.ts
git commit -m "H5活动详情接入营销补充与来源透传" -m "主要改动：活动列表跳详情携带activity_product_id，详情页按来源二次请求营销详情并透传加购/立即购买。\n\n原因：活动商品需展示秒杀倒计时、拼团入口等活动元素并保持来源一致。\n\n影响范围：H5营销列表、商品详情与多语言文案。"
```

### Task 6: H5 购物车与确认订单透传 `items`

**Files:**
- Modify: `app/pages/cart/index.vue`
- Modify: `app/pages/order/confirm.vue`
- Modify: `app/mock/index.ts`

- [ ] **Step 1: 先写失败用例（确认页仅 sku_ids 时无法识别活动来源）**

```md
在活动详情加入购物车后进入确认页，若请求体仅含sku_ids则活动来源丢失（当前失败基线）。
```

- [ ] **Step 2: 购物车勾选结算改为透传 items JSON**

```ts
const payload = checkedItems.map((item: any) => ({
  sku_id: item.sku_id,
  activity_product_id: Number(item.activity_product_id || 0),
}))
uni.navigateTo({ url: `/pages/order/confirm?items=${encodeURIComponent(JSON.stringify(payload))}` })
```

- [ ] **Step 3: 确认页优先解析 `items` 并提交下单**

```ts
const sourceItems = JSON.parse(decodeURIComponent(query.items || '[]'))
await post('/api/v1/orders', {
  address_id: address.value.id,
  payment_method: payMethod.value,
  items: sourceItems,
  sku_ids: sourceItems.map((it: any) => it.sku_id),
})
```

- [ ] **Step 4: mock 支持 cart/order 新结构**

```ts
if (upperMethod === 'POST' && path === '/api/v1/orders') {
  // accept items[{sku_id, activity_product_id}]
}
```

- [ ] **Step 5: 构建验证并提交**

Run: `npm run build:h5:demo`
Expected: PASS。

```bash
git add app/pages/cart/index.vue app/pages/order/confirm.vue app/mock/index.ts
git commit -m "H5结算链路透传活动来源items" -m "主要改动：购物车到确认页改为items透传，下单请求优先提交items并兼容sku_ids。\n\n原因：确保活动来源在下单阶段可被后端强校验。\n\n影响范围：H5购物车、确认订单与mock下单逻辑。"
```

### Task 7: PC 活动入口透传、详情营销区块与加购透传

**Files:**
- Modify: `web/src/views/ActivityProductListBase.vue`
- Modify: `web/src/views/ProductDetail.vue`
- Modify: `web/src/views/Cart.vue`
- Modify: `web/src/mock/index.ts`
- Modify: `web/src/locales/zh-CN.ts`
- Modify: `web/src/locales/en.ts`

- [ ] **Step 1: 先写失败用例（PC 活动详情无来源）**

```md
活动列表进入详情后 URL 不含 activity_product_id，详情页无法请求营销补充接口（当前失败基线）。
```

- [ ] **Step 2: 活动列表跳详情追加 query 参数**

```ts
$router.push(`/product/${p.product_id}?activity_product_id=${p.activity_product_id}`)
```

- [ ] **Step 3: 商品详情页实现双请求与活动按钮行为**

```ts
const activityProductID = computed(() => Number(route.query.activity_product_id || 0))
if (activityProductID.value > 0) {
  marketingDetail.value = await get(`/api/v1/marketing/activity-products/${activityProductID.value}`)
}
```

- [ ] **Step 4: 加购透传来源并在购物车展示活动标签**

```ts
await post('/api/v1/cart/add', {
  sku_id: selectedSku.value.id,
  qty: qty.value,
  activity_product_id: activityProductID.value,
})
```

- [ ] **Step 5: 构建验证并提交**

Run: `npm run build`
Expected: PASS。

```bash
git add web/src/views/ActivityProductListBase.vue web/src/views/ProductDetail.vue web/src/views/Cart.vue web/src/mock/index.ts web/src/locales/zh-CN.ts web/src/locales/en.ts
git commit -m "PC详情接入活动营销补充并透传来源" -m "主要改动：活动列表透传activity_product_id，详情页按来源请求营销详情并在加购时透传来源。\n\n原因：PC活动商品详情需补充活动元素且保持来源一致。\n\n影响范围：PC营销列表、详情、购物车与mock文案。"
```

### Task 8: docs-site 文档同步

**Files:**
- Modify: `docs-site/docs/api/marketing.md`
- Modify: `docs-site/docs/api/order.md`
- Modify: `docs-site/docs/api/product.md`

- [ ] **Step 1: 更新 marketing 文档（新增活动商品详情接口）**

```md
- GET /api/v1/marketing/activity-products/:id
- 说明：按activity_product_id返回活动商品详情与活动状态字段
```

- [ ] **Step 2: 更新 order 文档（购物车/下单参数升级）**

```md
- cart/add 与 cart/qty 新增 activity_product_id（可选）
- POST /api/v1/orders 新增 items[{sku_id,activity_product_id}]，传 items 时优先按 items 下单
```

- [ ] **Step 3: 更新 product 文档（详情页双请求说明）**

```md
- 商品详情标准接口保持不变
- 活动来源场景需配合营销活动商品详情接口补充渲染
```

- [ ] **Step 4: 文档构建验证并提交**

Run: `npm run docs:build`
Expected: PASS。

```bash
git add docs-site/docs/api/marketing.md docs-site/docs/api/order.md docs-site/docs/api/product.md
git commit -m "同步活动来源详情链路API文档" -m "主要改动：补充marketing活动详情接口、order购物车与下单参数升级、product详情双请求说明。\n\n原因：功能变更需同步文档并覆盖接口变化与部署影响。\n\n影响范围：docs-site API 文档。"
```

### Task 9: 全链路回归与发布前检查

**Files:**
- Modify: `server/plugins/*`（仅按回归结果最小修复）
- Modify: `app/*`、`web/*`（仅按回归结果最小修复）

- [ ] **Step 1: 后端回归测试**

Run: `go test ./...`
Expected: PASS。

- [ ] **Step 2: 前端构建回归测试**

Run:
- `cd app && npm run build:h5:demo`
- `cd web && npm run build`
- `cd docs-site && npm run docs:build`

Expected: 全部 PASS。

- [ ] **Step 3: 手工验收脚本执行**

```md
1) 秒杀/拼团/砍价列表进入详情，URL有activity_product_id
2) 详情页显示对应活动元素（倒计时/发起拼团/砍价入口）
3) 加购后购物车同SKU普通项与活动项分行展示
4) 确认订单提交items后，后端可按activity_product_id校验活动状态
5) 无activity_product_id的普通详情页不自动套活动价
```

- [ ] **Step 4: 汇总变更并一次性提交**

```bash
git status
git log --oneline -n 10
```

Expected: 仅包含本计划相关改动，提交信息全部为中文 head + body。

---

## Spec 覆盖自检

- 已覆盖“详情页双请求（标准 + 营销）”：Task 1、Task 5、Task 7。
- 已覆盖“来源携带 activity_product_id”：Task 2、Task 5、Task 7。
- 已覆盖“购物车与下单全链路透传”：Task 3、Task 4、Task 6。
- 已覆盖“无来源不自动套活动价”：Task 4（活动强校验分支仅在 `activity_product_id > 0` 触发）。
- 已覆盖“docs-site 同步”：Task 8。

## Placeholder 扫描

- 无 `TODO/TBD/后续补充/类似上一步` 等占位语句。
- 每个代码变更步骤提供了明确片段或函数签名。

## 类型与命名一致性检查

- 统一字段名：`activity_product_id`。
- 统一下单来源项结构：`items[{sku_id,activity_product_id}]`。
- 统一后端方法命名：`ValidateActivityProductSource`、`IncreaseSoldQtyByActivityProductTx`。

---

## 回归执行结果（2026-05-28）

### 执行结论

- Task 1 ~ Task 8 已完成并提交。
- Task 9 回归已完成，核心链路通过。
- 过程中发现并修复 1 个构建阻塞问题：H5 构建在 `pages-json-js` 虚拟模块上的 glob 路径校验失败。

### 回归命令与结果

- `cd server && go test ./...`：PASS
- `cd web && npm run -s build`：PASS
- `cd docs-site && npm run -s docs:build`：PASS
- `cd app && npm run -s build:h5:demo`：初次 FAIL，修复后 PASS

### H5 构建失败根因与修复

- 失败现象：
  - 错误：`[vite:import-glob] In virtual modules, all globs must start with '/'`
  - 位置：`pages-json-js`
- 根因分析：
  - uni 生成的 `pages-json-js` 虚拟模块中使用了相对 glob：`import.meta.glob('./locale/*.json', { eager: true })`。
  - Vite 对虚拟模块的 glob 路径要求为绝对路径（以 `/` 开头），因此构建失败。
- 修复方案：
  - 在 `app/vite.config.ts` 增加预处理插件 `patch-uni-pages-json-import-glob`。
  - 仅针对 `pages-json-js` 模块，将 `./locale/*.json` 定向改写为 `/locale/*.json`。
  - 该修复仅影响构建阶段转换，不改变业务运行时逻辑。

### 回归后状态

- H5、PC、后端与 docs-site 构建/测试结果全部为 PASS。
- 工作区仅保留计划文档未跟踪文件（按协作约定忽略），无业务代码残留未提交变更。
