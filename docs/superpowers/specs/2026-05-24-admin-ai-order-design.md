# 管理后台 AI 融合与订单链路增强设计（前后端联动）

**日期：** 2026-05-24
**范围：** `server` + `admin` + `web` + `app` + `mock`
**目标：** 在不无节制新增接口的前提下，优先升级现有接口语义，完成 AI 进商品编辑、订单列表/详情信息完整化与多端一致展示。

---

## 1. 需求目标（已确认）

1. 管理后台商品编辑页内集成 AI 生图，支持上传实物参考图，并用于详情图/封面/轮播图生成。
2. 详情图生成结果支持“插入到光标位置”，且 UI 有明确提示。
3. 同时支持封面图、轮播图 AI 生成，并预留商品介绍图生成入口。
4. 当模型不支持参考图生图时，前端禁用参考图并提示：`该模型不支持参考图`。
5. 管理后台 Dashboard 在导航上置顶（第一项）。
6. 订单体系补全：后台/H5/PC 列表展示商品明细，支持状态切换有效，新增订单详情页，展示价格体系明细与更多订单信息。
7. 修复 H5 订单列表操作按钮布局与个人中心客服图标。

---

## 2. 设计原则

1. **优先升级现有接口**：在 `orders` 与 `ai/generate` 语义内扩展字段，不平行新增无边界接口。
2. **多端统一数据模型**：后台/H5/PC 共享订单摘要与订单详情结构，减少端侧分叉逻辑。
3. **能力显式化**：AI 模型能力（是否支持参考图）由模型字段声明，前端据此禁用功能。
4. **可回溯**：AI 任务保留业务目标与参考图信息，便于复现与审计。
5. **渐进兼容**：升级返回结构时兼容旧字段，避免一次性破坏当前页面。

---

## 3. 后端接口与数据模型升级

### 3.1 AI 生图能力升级（复用 `POST /admin/api/ai/generate`）

#### 3.1.1 模型字段扩展

在 `ai_models` 增加字段：
- `supports_ref_image`：`tinyint(1)`，默认 `0`

用途：前端据此决定是否可上传参考图并参与生成。

#### 3.1.2 任务字段扩展

在 `ai_image_tasks` 增加字段：
- `biz_type`：`varchar(32)`，值域 `cover|gallery|detail|intro`
- `ref_image_url`：`varchar(500)`，参考图 URL
- `target_product_id`：`bigint unsigned`，目标商品 ID

用途：记录“为哪个商品、哪种图片用途、基于哪张参考图生成”。

#### 3.1.3 生成请求扩展（升级而非新增）

`POST /admin/api/ai/generate` 请求体扩展字段：
- `biz_type`（必填）
- `target_product_id`（可选）
- `ref_image_url`（可选）

保持 `scene/prompt/neg_prompt/params/model_id` 现有兼容。

#### 3.1.4 驱动参数扩展

`GenerateParams` 扩展：
- `RefImageURL string`

驱动策略：
- 支持图生图驱动：读取 `RefImageURL` 并传给上游。
- 不支持驱动：后端不做自动降级（前端已禁用）。若收到非法请求则返回明确错误。

---

### 3.2 订单查询能力升级（复用 `orders`，补充详情资源）

#### 3.2.1 列表接口升级

升级现有接口返回结构：
- 前台：`GET /api/v1/orders`
- 后台：`GET /admin/api/orders`

每个订单增加：
- `items`: `[{ product_id, sku_id, title, cover, attrs, price, qty, subtotal }]`
- `amount_breakdown`: `{ goods_amount, discount_amount, freight_amount, total_amount, activity_discount?, coupon_discount?, points_discount? }`

> 说明：`items` 由 `order_items` 关联查询填充；`amount_breakdown` 使用已有金额字段组装，规则细项可由 remark/规则快照解析补充。

#### 3.2.2 详情接口新增（资源语义清晰且必要）

新增：
- 前台：`GET /api/v1/orders/:id`
- 后台：`GET /admin/api/orders/:id`

返回内容：
- 订单基本信息（号、状态、时间线）
- 商品明细 `items`
- 价格体系 `amount_breakdown`
- 地址快照 `address_snapshot`
- 支付信息（method、paid_at）
- 物流信息（tracking_no，若存在）

> 新增详情接口的必要性：现有列表接口分页场景下不应承载全部详情字段；详情资源语义独立明确。

---

## 4. 管理后台改造

### 4.1 商品编辑页融合 AI（`ProductForm`）

新增 `AI 图片助手` 区块，包含：
- 目标类型：封面 / 轮播 / 详情 / 介绍图（预留）
- 模型选择 + 能力展示（是否支持参考图）
- 参考图上传（受 `supports_ref_image` 控制）
- Prompt 输入
- 生成结果列表与“一键应用”操作

应用行为：
- 封面：写入 `form.cover`
- 轮播：追加到商品图片列表
- 详情：插入到富文本光标位置
  - 明确提示文案：`将插入到当前光标位置`
- 介绍图：先保留为素材位与数据透传，不强制插入详情

### 4.2 Dashboard 首位

调整后台菜单渲染优先级：
- `Dashboard` 固定第一项
- 其他菜单按后端返回顺序/排序继续展示

### 4.3 后台订单列表与详情

- 订单列表表格扩列：用户、商品摘要、支付方式、优惠金额、实付金额、下单时间、状态、操作。
- 状态 tab 与后端 `status` 查询参数联动生效。
- 新增后台订单详情页：完整订单信息 + 商品明细 + 价格体系。

---

## 5. H5 与 PC 改造

### 5.1 H5 订单列表

- 每个订单卡片展示商品行（图、标题、规格、数量、单价）。
- 操作按钮外包一层容器，限制按钮自适应宽度，避免 `u-button` 占满 100%。
- 支持跳转订单详情页。

### 5.2 H5 订单详情页（新增）

- 展示：基本信息、商品明细、价格体系、地址、支付信息、物流信息。

### 5.3 H5 个人中心客服图标

- 替换为 uview-plus 确认可用图标名，避免缺失。

### 5.4 PC 订单列表与详情

- 列表补商品摘要行与价格体系摘要。
- 新增 PC 订单详情页并接路由。

---

## 6. Mock 与演示数据策略

### 6.1 admin mock

- `GET /admin/api/orders` 支持按 `status` 过滤并返回丰富订单项。
- 新增 `GET /admin/api/orders/:id`。
- `GET /admin/api/ai/models` 增加 `supports_ref_image`。
- `POST /admin/api/ai/generate` 支持 `biz_type/ref_image_url/target_product_id`。

### 6.2 app/web mock

- `GET /api/v1/orders` 补订单项与金额明细。
- 新增 `GET /api/v1/orders/:id`。
- `status` 过滤保持可用，确保 tab 切换有效果。

---

## 7. 兼容性与风险控制

1. **兼容旧前端**：后端返回结构新增字段，不移除原字段。
2. **参考图能力门禁**：前端禁用 + 后端兜底校验双保险。
3. **查询性能**：列表页使用批量查询 order_items，避免 N+1。
4. **详情资源边界**：列表只放摘要，详情接口承载完整信息。
5. **Mock 与真实接口一致**：避免演示态与生产态行为偏差。

---

## 8. 测试与验收标准

### 8.1 后端

- `go test ./...` 通过。
- 订单列表/详情接口返回结构正确。
- AI 生成接口可接收扩展字段并正确落任务。

### 8.2 管理后台

- 商品编辑页 AI 生成：
  - 不支持参考图模型 -> 控件禁用且提示正确
  - 支持参考图模型 -> 可上传参考图并生成
  - 详情图可插入光标位置
- Dashboard 在导航第一位。
- 后台订单 tab 切换有效，详情页信息完整。

### 8.3 H5/PC

- 订单列表展示商品行与金额摘要。
- 订单详情页可达，展示价格体系明细。
- H5 按钮布局正常不拉满。
- H5 个人中心客服图标正常显示。

---

## 9. 涉及文件清单（预期）

- `server/core/driver/ai/ai.go`
- `server/plugins/ai_image/model/ai_image.go`
- `server/plugins/ai_image/service/ai_image.go`
- `server/plugins/ai_image/api/admin.go`
- `server/plugins/order/model/order.go`
- `server/plugins/order/service/order.go`
- `server/plugins/order/api/front.go`
- `server/plugins/order/api/admin.go`
- `admin/src/views/product/ProductForm.vue`
- `admin/src/views/order/OrderList.vue`
- `admin/src/views/order/OrderDetail.vue`（新增）
- `admin/src/layouts/AdminLayout.vue`
- `admin/src/router/index.ts`
- `admin/src/api/plugins.ts`
- `admin/src/mock/index.ts`
- `app/pages/order/list.vue`
- `app/pages/order/detail.vue`（新增）
- `app/pages/user/index.vue`
- `app/pages.json`
- `app/mock/index.ts`
- `app/mock/data/orders.json`
- `web/src/views/OrderList.vue`
- `web/src/views/OrderDetail.vue`（新增）
- `web/src/router/index.ts`
- `web/src/api/request.ts`
- `web/src/mock/index.ts`
- `README.md`
- `docs-site/docs/guide/features.md`
- `docs-site/docs/api/order.md`
- `docs-site/docs/api/im.md`

---

## 10. 超范围说明（本轮不做）

1. 不引入新的富文本编辑器库，仅在现有表单能力下实现光标插入。
2. 不重构 AI 驱动实现细节（仅增加参数透传与能力约束）。
3. 不引入订单物流第三方查询，仅预留物流字段展示。

