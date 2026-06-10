# 商品接口

## 功能说明

商品模块负责分类、商品、SKU、详情与评价能力，并支持 **SPU 维度商品收藏**：

- 用户可在商品详情页执行收藏/取消收藏（需登录）。
- 用户可在个人中心查看“我的收藏”列表。
- 商品列表、商品详情均返回收藏态与收藏数。
- 后台商品列表可直接读取收藏数字段展示。
- 活动来源场景下，商品详情页采用“双请求”模式：先请求标准商品详情，再按 `activity_product_id` 追加请求营销活动商品详情接口。
- 商品编辑统一采用后端规格引擎：前端提交 `sku_generation_mode=auto + spec_schema + sku_overrides`，后端生成并维护 SKU 集合并返回 `sku_diff`。
- 后台商品编辑的封面图与轮播图使用图片上传组件维护：上传复用 `POST /admin/api/upload`，也支持直接填写图片 URL，保存时仍提交商品 `cover` 字段与 `images[]` 轮播图数组。
- 规格模板为正式后台资源：`admin` 与 `eapp` 共用 `/admin/api/spec-templates*` 接口。
- 后台商品编辑支持在 SKU 区直接配置会员价（需启用 `vip` 插件）。
- 当前最新架构下，SKU 是商品销售、库存占用、库存扣减、库存回补、库存查询的统一粒度。

## SKU 与库存关系

- `product` 负责商品、SPU、SKU、详情、规格生成与展示能力。
- `inventory` 负责统一库存交易与库存来源选择，不再默认由 `wms` 承担全部库存语义。
- `product_skus.stock` 在 `local` 模式下直接作为库存交易值使用。
- 当 `inventory.provider=builtin_wms` 时：
  - SKU 仍然是订单明细和库存交易的统一粒度
  - 实际库存真源切换为内置 `wms` 插件库存表
- 当 `inventory.provider=external_wms` 时：
  - SKU 仍然是商城与企业 WMS 交互的统一库存粒度
  - 商品详情页返回的 SKU 库存可以由外部可售库存查询结果覆盖

说明：

- `products.stock` 适合作为商品聚合展示值，不作为统一库存交易接口。
- 统一库存动作始终围绕 `sku_id` 进行，而不是围绕 SPU 或商品总库存直接操作。

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
  - SKU 库存口径说明：
    - `inventory.provider=local`：直接返回 `product_skus.stock`
    - `inventory.provider=builtin_wms`：SKU 详情仍按商城模型返回，库存交易由内置 WMS 处理
    - `inventory.provider=external_wms`：当外部库存查询成功时，SKU 库存可由外部可售库存覆盖
  - 标准商品详情接口保持不变；当页面带有活动来源（`activity_product_id`）时，前端应追加请求：
    - `GET /api/v1/marketing/activity-products/:id`
  - 用途：补充渲染活动扩展元素（秒杀倒计时、拼团/砍价入口、活动库存与限购等）
- `GET /admin/api/products`
  - 列表项新增：
    - `favorite_count: number`
- `PUT /admin/api/products/:id`
  - 推荐请求模式：`sku_generation_mode=auto`
    - `spec_schema`: 规格组定义（属性名 + 值集合）
    - `sku_overrides`: 覆盖项（按 `sku_key` 指定价格/库存/编码）
    - `product.cover`: 商品封面图片 URL，可由管理端上传组件写入
    - `images`: 商品轮播图数组，元素结构为 `{ url, sort }`，可由管理端上传组件写入
  - `sku_key` 唯一约束为 `(product_id, sku_key)`，确保“同商品内唯一”。
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

- `GET /admin/api/spec-templates?page=1&size=20&keyword=&category_id=`
  - 说明：规格模板分页列表（支持关键词与分类筛选）
- `GET /admin/api/spec-templates/:id`
  - 说明：规格模板详情
- `POST /admin/api/spec-templates`
  - 请求：`{ name, category_ids, attrs, status, sort? }`
  - `attrs` 结构：`[{ name, values[] }]`
- `PUT /admin/api/spec-templates/:id`
  - 请求：同上字段，按需传递
- `DELETE /admin/api/spec-templates/:id`
  - 说明：删除规格模板

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
- `GET /admin/api/spec-templates`
- `POST /admin/api/upload`
- `POST /admin/api/spec-templates`
- `PUT /admin/api/spec-templates/:id`
- `DELETE /admin/api/spec-templates/:id`

## 部署与配置影响

- 无新增环境变量。
- 无新增中间件或外部依赖。
- 商品封面和轮播图上传沿用后台通用上传接口与当前 storage driver；本地部署需保持上传目录可被前端访问。
- SKU 模型新增字段（自动迁移）：
  - `product_skus.sku_key`
  - `product_skus.status`（`active|inactive`）
- SKU 唯一索引为 `(product_id, sku_key)`（自动迁移会重建 `uk_product_sku_key`）。
- 统一库存默认模式为 `inventory.provider=local`，此时 `product_skus.stock` 直接参与库存交易。
- 当切换到 `inventory.provider=builtin_wms` 或 `external_wms` 时：
  - 商品模块仍负责维护 SKU 主数据
  - 库存交易与可售库存来源由统一 `inventory` provider 决定
- 新增库存相关配置说明见：
  - [库存预占交易规则](./stock-reservation)
  - [仓储接口](./wms)
- 新增表：`spec_templates`（规格模板）。
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
