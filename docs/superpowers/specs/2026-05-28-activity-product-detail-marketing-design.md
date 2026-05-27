# 活动商品详情补充链路设计（H5/PC）

## 1. 背景与目标

当前 H5 与 PC 商品详情页仅请求标准商品详情接口，无法在活动商品场景下展示活动特有元素（如秒杀倒计时、拼团发起入口、砍价入口与活动库存/限购状态）。

目标：在保持现有商品详情接口兼容的前提下，建立“活动来源驱动”的营销补充链路，实现：

1. 活动入口跳转详情时携带活动商品唯一 ID（`activity_product_id`）。
2. 详情页在标准详情之外，按需请求营销活动商品详情接口并渲染额外元素。
3. 购物车与下单全链路透传 `activity_product_id`，后端强约束活动价格与库存校验。
4. 无来源（无 `activity_product_id`）时按普通商品流程，不自动套活动价。

## 2. 方案对比与选型

### 方案 A：仅详情补充展示，不改下单链路

- 优点：改动小，交付快。
- 缺点：无法保证活动价与活动库存在下单阶段被强约束，不满足业务要求。

### 方案 B：兼容升级现有接口（选型）

- 做法：新增活动商品详情接口；在详情、购物车、下单链路透传 `activity_product_id`；保留旧字段并向后兼容。
- 优点：满足“来源决定活动”与“全链路强约束”，改造面可控。
- 缺点：需要调整购物车键模型与下单请求结构。

### 方案 C：新建营销专用购物车/下单接口

- 优点：语义最清晰，普通链路与营销链路完全隔离。
- 缺点：改造与回归范围大，重复建设成本高。

结论：采用方案 B。

## 3. 关键业务规则

1. 活动来源决定行为：只有携带 `activity_product_id` 的请求才进入活动逻辑。
2. 无来源不自动匹配：未携带 `activity_product_id` 时，始终按普通商品购买流程处理。
3. 活动商品唯一标识：使用 `activity_products.id` 作为 `activity_product_id`。
4. 全链路透传：详情 -> 加购/立即购买 -> 购物车 -> 结算 -> 下单 均保留该 ID。

## 4. 接口设计

## 4.1 商品与营销详情

- 保持：`GET /api/v1/products/:id`（标准商品详情，不变）。
- 新增：`GET /api/v1/marketing/activity-products/:id`
  - `:id` 为 `activity_product_id`。
  - 返回活动商品详情（活动类型、活动名称、时间窗、活动价、限购、库存状态、拼团配置等）。

## 4.2 购物车接口兼容升级

- `POST /api/v1/cart/add`
  - 新增可选：`activity_product_id`，默认 `0`。
- `PUT /api/v1/cart/qty`
  - 新增可选：`activity_product_id`（用于定位同 SKU 的不同活动行）。
- `DELETE /api/v1/cart/:sku_id`
  - 新增可选参数：`activity_product_id`。
- `GET /api/v1/cart`
  - 每项返回 `activity_product_id` 与营销快照字段。

## 4.3 下单接口兼容升级

- `POST /api/v1/orders`
  - 保留：`sku_ids`（旧调用继续可用）。
  - 新增：`items: [{ sku_id, activity_product_id }]`。
  - 规则：若传 `items`，优先使用 `items` 作为结算来源并执行活动强校验。

## 5. 数据模型与存储设计

## 5.1 购物车存储键升级

Redis 购物车 field 从 `sku_id` 升级为复合键：

- `sku_id:activity_product_id`
- 示例：`101:0`（普通）、`101:9001`（活动）

目的：避免同 SKU 的普通价与活动价在购物车相互覆盖。

## 5.2 订单项快照升级

`order_items` 增加字段：

- `activity_product_id`
- `activity_id`
- `activity_type`
- `activity_title`

目的：保障订单可追溯、售后与统计可识别活动来源。

## 6. 前端改造设计（H5/PC）

## 6.1 活动列表到详情的来源透传

- H5：
  - `pages/marketing/seckill.vue`
  - `pages/marketing/group-buy.vue`
  - `pages/marketing/bargain.vue`
- PC：
  - `web/src/views/ActivityProductListBase.vue`

上述页面跳转详情时统一追加 `activity_product_id`。

## 6.2 商品详情页双请求流程

1. 请求标准详情 `GET /api/v1/products/:id`。
2. 若 URL 存在 `activity_product_id`，再请求营销详情接口。
3. 合并渲染：标准信息 + 活动扩展区块。

## 6.3 活动扩展 UI 规则

- 秒杀：倒计时、活动价、库存进度、限购提示、抢购按钮。
- 拼团：拼团价、开团人数（如 `group_size`）、发起拼团按钮。
- 砍价：起砍价、底价、发起砍价/帮砍入口。

## 6.4 购物与下单透传

- 详情页“加入购物车”“立即购买”提交 `activity_product_id`。
- 购物车展示活动标签与营销快照。
- 结算页提交 `items[{sku_id, activity_product_id}]`。

## 7. 后端校验与结算流程

当 `activity_product_id > 0` 时执行强校验：

1. 活动商品记录存在。
2. `activity_product_id` 与 `sku_id` / `product_id` 绑定一致。
3. 活动状态有效（启用、未过期、已开始）。
4. 库存约束通过（活动总库存上限、已售数量）。
5. 限购约束通过（单笔限购）。

价格规则：

- 有 `activity_product_id`：按该活动商品记录定价。
- 无 `activity_product_id`：按普通商品价格。

已售累计规则：

- 按 `activity_product_id` 精确累加，避免同 SKU 多活动串量。

## 8. 失败与降级策略

1. 详情页传入 `activity_product_id` 但营销详情无效：
   - 显示“活动已失效/已结束”提示；
   - 禁用活动购买入口；
   - 不自动切到普通活动价。
2. 下单校验失败：返回明确业务错误码与错误信息，前端按错误类型提示用户刷新或返回活动页。

## 9. 测试方案

## 9.1 后端

- 单元测试：
  - 活动详情查询；
  - 活动强校验分支（时间、库存、限购、关联关系）；
  - 按 `activity_product_id` 累加已售。
- 集成测试：
  - 活动来源下单成功；
  - 活动结束后下单失败；
  - 同 SKU 普通/活动并存购物车互不覆盖。

## 9.2 前端（H5/PC）

- 活动列表跳转是否携带 `activity_product_id`。
- 详情页是否按有无来源触发营销二次请求。
- 活动区块是否按类型正确渲染。
- 加购/立即购买/结算下单是否透传来源 ID。

## 10. 文档同步范围（docs-site）

按协作规则，功能变更需同步 `docs-site`，至少覆盖：

1. 功能说明：活动来源驱动详情补充链路。
2. 接口变化：营销详情接口新增、购物车/下单参数升级。
3. 部署/配置影响：数据库新增字段与兼容说明（如有迁移步骤）。

建议更新文档：

- `docs-site/docs/api/marketing.md`
- `docs-site/docs/api/order.md`
- `docs-site/docs/api/product.md`

## 11. 影响范围

- 后端：marketing、product、order 插件及相关模型。
- 前端 PC：活动列表、商品详情、购物车、结算/下单入口。
- 前端 H5：活动列表、商品详情、购物车、确认订单。
- 文档：docs-site API 文档与功能说明。

## 12. 非目标

1. 本次不引入“无来源自动套活动价”。
2. 本次不拆分全新营销专用下单接口。
3. 本次不处理多活动同时命中的自动仲裁逻辑（由来源 ID 显式指定）。
