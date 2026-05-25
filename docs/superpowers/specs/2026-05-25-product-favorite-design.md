# 商品收藏功能设计

## 目标

在 `app`、`web`、`admin` 三端实现商品收藏能力，采用 SPU 粒度（商品级收藏）。

- 用户端：商品详情可收藏/取消收藏，未登录必须先登录。
- 用户端：个人中心可查看“我的收藏”列表，取消收藏后立即从列表移除并补位。
- 管理端：商品列表显示收藏数。

## 范围

- 后端：`product` 插件内新增收藏数据模型、服务与前台接口；扩展现有商品查询返回字段。
- app：商品详情页收藏操作、个人中心收藏入口、收藏列表页。
- web：商品详情页收藏操作、用户中心收藏分栏与列表。
- admin：商品列表增加收藏数字段展示。
- docs-site：同步更新商品 API 文档与功能说明。

## 约束与原则

- 优先升级现有接口，避免新增平行查询接口。
- 仅在现有接口无法覆盖“写收藏关系”和“我的收藏列表”时新增最小接口集合。
- 收藏操作幂等：重复收藏/重复取消均返回成功，且不破坏计数。

## 数据模型设计

### 新增表

`product_favorites`

- `id`（主键）
- `user_id`（索引）
- `product_id`（索引）
- `created_at`
- `updated_at`
- 唯一索引：`uniq_user_product(user_id, product_id)`

### 现有表扩展

`products` 新增字段：

- `favorite_count int not null default 0`

用于后台直接展示与前台快速读取，避免高频聚合查询。

## 接口设计

### 兼容升级接口（保留原语义）

1. `GET /api/v1/products`
   - 列表项增加：`favorite_count`、`is_favorited`
   - `is_favorited`：仅登录态计算；未登录或无 token 返回 `false`
2. `GET /api/v1/products/:id`
   - 详情增加：`favorite_count`、`is_favorited`
3. `GET /admin/api/products`
   - 列表项增加：`favorite_count`

### 新增最小接口（必要新增）

1. `POST /api/v1/products/:id/favorite`（需登录）
2. `DELETE /api/v1/products/:id/favorite`（需登录）
3. `GET /api/v1/user/favorites?page=1&size=20`（需登录）

新增必要性：现有接口无法表达“用户与商品的收藏关系写入/删除”以及“按用户维度的收藏列表读取”。

## 服务与数据流

### 收藏

- 校验商品存在且可收藏（可按当前商品可见性规则校验）。
- 事务内执行：
  1) 尝试插入 `product_favorites`
  2) 插入成功才 `products.favorite_count = favorite_count + 1`
- 若唯一冲突，按幂等成功返回，不增计数。

### 取消收藏

- 事务内执行：
  1) 删除 `product_favorites`
  2) 仅在删除行数 > 0 时，`favorite_count` 原子减 1（下限保护为 0）
- 若原本未收藏，按幂等成功返回。

### 查询收藏状态

- 商品列表/详情查询时，若有用户身份，批量查询该用户已收藏商品 ID 集并映射到 `is_favorited`。

### 我的收藏列表

- 从 `product_favorites` 按 `created_at desc` 分页，并关联商品基础信息返回。

## 三端页面改造

### app

- `pages/product/detail`：新增收藏按钮。
- `pages/user/index`：新增“我的收藏”入口。
- 新增 `pages/user/favorites` 页面：分页列表、跳转详情、取消即移除。

### web

- `views/ProductDetail.vue`：新增收藏按钮。
- `views/UserCenter.vue`：新增“我的收藏”菜单与列表视图。
- 列表取消后本地立即移除并补位。

### admin

- `views/product/ProductList.vue`：新增“收藏数”列。

## 错误处理与交互

- 未登录收藏：前端直接跳登录。
- 商品不存在：返回商品不存在业务错误。
- 收藏/取消失败：前端提示并回滚乐观更新。
- 计数保护：任何路径都不允许 `favorite_count` 变负。

## 测试与验证

### 后端

- 收藏新增、重复收藏幂等、取消收藏幂等、计数增减正确。
- 收藏列表分页正确。
- 商品列表/详情 `is_favorited` 与 `favorite_count` 返回正确。

### 前端

- app/web 详情页收藏状态切换。
- app/web 我的收藏页取消即移除。
- admin 商品列表收藏数字段展示。

### mock

- `app/mock`、`web/mock`、`admin/mock` 同步收藏相关路由与字段，保证演示模式可运行。

## 部署与配置影响

- 无新增环境变量、无新增外部依赖。
- 数据库迁移新增：`product_favorites` 表、`products.favorite_count` 字段（通过现有 AutoMigrate 生效）。
- 需同步更新 `docs-site/docs/api/product.md`，覆盖：功能说明、接口变化、部署影响。
