# 营销接口

## 功能说明

营销模块用于优惠券、秒杀、拼团、砍价与积分等能力。  
本次调整新增了 PC 端三类活动商品列表（秒杀/拼团/砍价）、活动商品详情补充接口，以及后台六菜单管理能力（活动管理 + 商品管理）。

## 接口变化

### C 端（公开）

新增独立活动商品列表接口：

- `GET /api/v1/marketing/seckill/products`
- `GET /api/v1/marketing/group-buy/products`
- `GET /api/v1/marketing/bargain/products`
- `GET /api/v1/marketing/activity-products/:id`

其中 `GET /api/v1/marketing/activity-products/:id` 用于活动商品详情补充（`id=activity_product_id`），面向 H5/PC 商品详情页在标准商品详情外追加营销元素渲染（如秒杀倒计时、拼团入口、砍价入口）。
核心返回字段包含：

- `activity_product_id`、`activity_id`、`activity_type`、`activity_name`
- `activity_status`、`activity_start_at`、`activity_end_at`
- `activity_price`、`start_price`、`floor_price`、`price`、`origin_price`
- `limit_per_order`、`total_stock_limit`、`sold_qty`

通用查询参数：

- `activity_id`：可选，指定活动批次。
- `category_id`：可选，按分类筛选。
- `keyword`：可选，按商品标题模糊搜索。
- `min_price` / `max_price`：可选，价格区间筛选。
- `sort_by`：可选，`price` 或 `sales`。
- `sort_order`：可选，`asc` 或 `desc`。
- `page` / `size`：分页参数。

返回包含：活动批次信息、商品信息、活动价格信息（秒杀/拼团活动价、砍价起始价/最低价）、每单限购、活动总库存上限、已售数量。
活动列表项兼容返回 `activity_product_id`，用于详情页与下单链路透传来源 ID。

### C 端（需登录）

既有接口：

- `GET /api/v1/coupons`
- `POST /api/v1/coupons/:id/claim`
- `GET /api/v1/user/coupons`
- `GET /api/v1/user/points/logs`

其中 `GET /api/v1/user/coupons` 已在原接口上兼容升级：返回项新增 `coupon` 对象（优惠券快照），包含 `id/name/type/min_amount/discount/start_at/end_at/status` 等字段，前端可直接渲染券面信息，避免再次按 `coupon_id` 补查。旧字段 `coupon_id/status/used_at/order_id` 保持不变。

### 管理端（需 `marketing:view` 或 `marketing:edit`）

#### 秒杀

- `GET /admin/api/marketing/seckill/activities`
- `POST /admin/api/marketing/seckill/activities`
- `PUT /admin/api/marketing/seckill/activities/:id`
- `GET /admin/api/marketing/seckill/products`
- `PUT /admin/api/marketing/seckill/activities/:id/products`

#### 拼团

- `GET /admin/api/marketing/group-buy/activities`
- `POST /admin/api/marketing/group-buy/activities`
- `PUT /admin/api/marketing/group-buy/activities/:id`
- `GET /admin/api/marketing/group-buy/products`
- `PUT /admin/api/marketing/group-buy/activities/:id/products`

#### 砍价

- `GET /admin/api/marketing/bargain/activities`
- `POST /admin/api/marketing/bargain/activities`
- `PUT /admin/api/marketing/bargain/activities/:id`
- `GET /admin/api/marketing/bargain/products`
- `PUT /admin/api/marketing/bargain/activities/:id/products`

商品管理请求体为 SKU 维度明细数组：

- `product_id`、`sku_id`
- `activity_price`（秒杀/拼团）
- `start_price`、`floor_price`（砍价）
- `limit_per_order`
- `total_stock_limit`

## 业务规则

- 同类型活动（秒杀/拼团/砍价）时间区间不可重叠，但可提前创建未来批次。
- 以 `活动 + SKU` 为管理粒度，支持每单限购与活动总库存上限。
- 共享商品 SKU 库存；下单时额外校验活动库存上限。
- 下单成功后会累计活动已售 `sold_qty`，用于后续库存上限校验。

## 部署或配置影响

- 无新增环境变量。
- 无新增配置项。
- 数据库 `activity_products` 增加字段：
  - `start_price`
  - `floor_price`
  - `limit_per_order`
  - `total_stock_limit`
  - `sold_qty`
- 后台菜单新增：
  - 秒杀活动管理、秒杀商品管理
  - 拼团活动管理、拼团商品管理
  - 砍价活动管理、砍价商品管理
