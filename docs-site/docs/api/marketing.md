# 营销接口

## 功能说明

营销模块用于优惠券、秒杀、拼团、砍价与积分等能力。  
当前接口采用“活动插件独立路由”规范：秒杀、拼团、砍价均使用各自命名空间；营销主插件保留优惠券与活动详情补充能力。

## 接口变化

### C 端（公开）

新增独立活动商品列表接口：

- `GET /api/v1/seckill/products`
- `GET /api/v1/group-buy/products`
- `GET /api/v1/bargain/products`
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

### 管理端

#### 秒杀

- `GET /admin/api/seckill/activities`
- `POST /admin/api/seckill/activities`
- `PUT /admin/api/seckill/activities/:id`
- `GET /admin/api/seckill/products`
- `PUT /admin/api/seckill/activities/:id/products`

权限：`seckill:view`、`seckill:edit`

#### 拼团

- `GET /admin/api/group-buy/activities`
- `POST /admin/api/group-buy/activities`
- `PUT /admin/api/group-buy/activities/:id`
- `GET /admin/api/group-buy/products`
- `PUT /admin/api/group-buy/activities/:id/products`

权限：`group_buy:view`、`group_buy:edit`

#### 砍价

- `GET /admin/api/bargain/activities`
- `POST /admin/api/bargain/activities`
- `PUT /admin/api/bargain/activities/:id`
- `GET /admin/api/bargain/products`
- `PUT /admin/api/bargain/activities/:id/products`

权限：`bargain:view`、`bargain:edit`

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
- 需启用对应插件：
  - `marketing`（优惠券与活动详情补充）
  - `seckill`
  - `group_buy`
  - `bargain`
- 活动数据按插件分表存储：
  - 秒杀：`seckill_activities`、`seckill_products`
  - 拼团：`group_buy_activities`、`group_buy_products`、`group_buy_orders`、`group_buy_members`
  - 砍价：`bargain_activities`、`bargain_products`、`bargain_orders`、`bargain_helpers`
- 后台菜单新增：
  - 秒杀活动管理、秒杀商品管理
  - 拼团活动管理、拼团商品管理
  - 砍价活动管理、砍价商品管理

## 优惠券 CRUD (P2 新增)

### 管理端

#### 列表查询

`GET /admin/api/marketing/coupons`

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 按名称模糊搜索 |
| status | number | 0=禁用, 1=启用 |
| type | number | 1=满减, 2=折扣, 3=立减 |
| page | number | 页码 |
| size | number | 每页条数 |

返回字段（新增）：
- `used_count` — 已使用数量
- `start_at` / `end_at` — 有效期
- `description` — 说明
- `stack_rule` — 叠加规则：`exclusive` / `same_type` / `cross_type`
- `target_type` — 目标用户：`all` / `vip_level` / `new_user`
- `target_value` — 目标值（如会员等级 ID）

#### 创建

`POST /admin/api/marketing/coupons`

请求体包含上述所有字段（除 `used_count`）。

#### 更新

`PUT /admin/api/marketing/coupons/:id`

请求体为需要更新的字段子集。

#### 删除

`DELETE /admin/api/marketing/coupons/:id`

#### 定向发券

`POST /admin/api/marketing/coupons/:id/send`

| 参数 | 类型 | 说明 |
|------|------|------|
| count | number | 发送数量 |

返回：`{ sent_count: number }`

发券后自动累加优惠券的 `used_count`。
