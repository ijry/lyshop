# 商品接口

## 功能说明

商品模块负责分类、商品、SKU、详情与评价能力，并支持 **SPU 维度商品收藏**：

- 用户可在商品详情页执行收藏/取消收藏（需登录）。
- 用户可在个人中心查看“我的收藏”列表。
- 商品列表、商品详情均返回收藏态与收藏数。
- 后台商品列表可直接读取收藏数字段展示。
- 活动来源场景下，商品详情页采用“双请求”模式：先请求标准商品详情，再按 `activity_product_id` 追加请求营销活动商品详情接口。
- 后台商品编辑支持维护 SKU 列表，并在启用 `vip` 插件时在 SKU 区直接配置会员价。

## 接口变化

### 已有接口（兼容升级）

- `GET /api/v1/products`
  - 列表项新增：
    - `favorite_count: number`
    - `is_favorited: boolean`（未登录时为 `false`）
- `GET /api/v1/products/:id`
  - 详情新增：
    - `favorite_count: number`
    - `is_favorited: boolean`（未登录时为 `false`）
  - 标准商品详情接口保持不变；当页面带有活动来源（`activity_product_id`）时，前端应追加请求：
    - `GET /api/v1/marketing/activity-products/:id`
  - 用途：补充渲染活动扩展元素（秒杀倒计时、拼团/砍价入口、活动库存与限购等）
- `GET /admin/api/products`
  - 列表项新增：
    - `favorite_count: number`
- `PUT /admin/api/products/:id`
  - 请求体兼容升级：支持可选字段 `skus: ProductSku[]`
  - 语义：提交时按商品维度替换 SKU 集合（用于商品管理页统一维护 SKU）
  - 当前推荐模式：`sku_generation_mode=auto`
    - `spec_schema`: 规格组定义（属性名 + 值集合）
    - `sku_overrides`: 覆盖项（按 `sku_key` 指定价格/库存/编码）
  - 返回 `sku_diff`：`{ added, kept, inactivated }`
  - 旧 SKU 采用软删除（`status=inactive`），历史订单可继续读取。

### 新增接口

- `POST /api/v1/products/:id/favorite`
  - 说明：收藏商品（需登录）
  - 返回：`data = null`

- `DELETE /api/v1/products/:id/favorite`
  - 说明：取消收藏（需登录）
  - 返回：`data = null`

- `GET /api/v1/user/favorites?page=1&size=20`
  - 说明：我的收藏分页列表（需登录）
  - 返回：`{ list, total, page, size }`
  - 列表项包含商品基础字段及：
    - `is_favorited: true`
    - `favorited_at: string`

## 典型接口清单

- `GET /api/v1/categories`
- `GET /api/v1/products`
- `GET /api/v1/products/:id`
- `GET /api/v1/products/recommend?limit=8`
- `GET /api/v1/products/:id/reviews?page=1&size=10`
- `POST /api/v1/products/:id/favorite`
- `DELETE /api/v1/products/:id/favorite`
- `GET /api/v1/user/favorites?page=1&size=20`
- `GET /admin/api/categories`
- `POST /admin/api/categories`
- `PUT /admin/api/categories/:id`
- `DELETE /admin/api/categories/:id`
- `GET /admin/api/products`

## 部署与配置影响

- 无新增环境变量。
- 无新增中间件或外部依赖。
- SKU 模型新增字段（自动迁移）：
  - `product_skus.sku_key`
  - `product_skus.status`（`active|inactive`）
- 活动来源详情双请求仅涉及前端调用链路与营销接口协同，无新增配置项。
- 数据库新增迁移：
  - 新表 `product_favorites`
  - `products` 表新增字段 `favorite_count`

## 商家移动端商品增强接口

### GET /categories/tree

返回 3 级分类树：`[{ id, parent_id, name, sort, product_count, children: [...] }]`

### POST /categories | PUT /categories/{id} | DELETE /categories/{id}

CRUD：POST `{ name, parent_id?, sort? }`、PUT `{ name?, sort? }`、DELETE 无 body。

### PUT /products/batch/status | category | price

批量上下架 / 批量分类 / 批量调价。

请求：
- status: `{ ids: number[], status: 0|1 }`
- category: `{ ids: number[], category_id }`
- price: `{ ids: number[], adjustment: { type: 'set'|'percent'|'amount', value, scope?: 'all'|'main_sku' } }`

返回：`{ success_ids[], fail: [{id, reason}] }`

### GET /products?status=&category_id=&sort_by=&low_stock=

列表扩展查询参数。`sort_by` 取值 `sales|stock|price_asc|price_desc|created`。
- 依赖现有自动迁移流程，无需额外配置项。
