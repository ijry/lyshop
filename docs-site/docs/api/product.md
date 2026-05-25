# 商品接口

## 功能说明

商品模块负责分类、商品、SKU、详情与评价能力，并支持 **SPU 维度商品收藏**：

- 用户可在商品详情页执行收藏/取消收藏（需登录）。
- 用户可在个人中心查看“我的收藏”列表。
- 商品列表、商品详情均返回收藏态与收藏数。
- 后台商品列表可直接读取收藏数字段展示。

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
- `GET /admin/api/products`
  - 列表项新增：
    - `favorite_count: number`

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
- 数据库新增迁移：
  - 新表 `product_favorites`
  - `products` 表新增字段 `favorite_count`
- 依赖现有自动迁移流程，无需额外配置项。
