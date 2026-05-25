# 商品收藏功能 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为 `app`、`web`、`admin` 三端实现 SPU 商品收藏（登录后收藏/取消、我的收藏列表、后台收藏数展示），并保持现有商品接口兼容升级。

**Architecture:** 在 `product` 插件新增 `product_favorites` 关系表并给 `products` 增加 `favorite_count` 冗余字段。前台以鉴权接口提供收藏写操作与收藏列表，商品列表/详情接口扩展 `is_favorited` 与 `favorite_count` 字段。前端三端复用同一接口语义，`app/web` 直接在详情页与个人中心落地，`admin` 仅读取收藏数字段。

**Tech Stack:** Go + Gin + GORM、Vue3 + Pinia + Axios、uni-app、VitePress docs-site。

---

## 文件结构与职责

- `server/plugins/product/model/favorite.go`（新建）：收藏关系模型与索引。
- `server/plugins/product/model/product.go`（修改）：`Product` 增加 `favorite_count` 字段。
- `server/plugins/product/plugin.go`（修改）：加入收藏模型迁移。
- `server/plugins/product/service/favorite.go`（新建）：收藏、取消、收藏列表、收藏态查询。
- `server/plugins/product/service/product.go`（修改）：列表与详情注入 `is_favorited`；结构体增加该字段。
- `server/plugins/product/api/front.go`（修改）：新增收藏相关路由与处理函数。
- `app/pages/product/detail.vue`（修改）：收藏按钮与状态。
- `app/pages/user/index.vue`（修改）：个人中心入口。
- `app/pages/user/favorites.vue`（新建）：我的收藏页面。
- `app/pages.json`（修改）：注册收藏页面路由。
- `web/src/views/ProductDetail.vue`（修改）：收藏按钮与状态。
- `web/src/views/UserCenter.vue`（修改）：收藏菜单与列表视图。
- `admin/src/views/product/ProductList.vue`（修改）：新增收藏数列。
- `app/mock/index.ts`（修改）：收藏相关 mock。
- `web/src/mock/index.ts`（修改）：收藏相关 mock。
- `admin/src/mock/index.ts`（修改）：后台商品返回收藏数字段 mock。
- `docs-site/docs/api/product.md`（修改）：功能说明、接口变化、部署影响。

### Task 1: 后端模型与迁移

**Files:**
- Create: `server/plugins/product/model/favorite.go`
- Modify: `server/plugins/product/model/product.go`
- Modify: `server/plugins/product/plugin.go`

- [ ] **Step 1: 新建收藏模型并声明唯一索引**

```go
package model

type ProductFavorite struct {
    model.Base
    UserID    uint64 `gorm:"not null;index;uniqueIndex:uniq_user_product" json:"user_id"`
    ProductID uint64 `gorm:"not null;index;uniqueIndex:uniq_user_product" json:"product_id"`
}
```

- [ ] **Step 2: 扩展商品模型收藏数字段**

```go
FavoriteCount int `gorm:"not null;default:0" json:"favorite_count"`
```

- [ ] **Step 3: 注册迁移模型**

```go
return db.AutoMigrate(
  &productmodel.ProductCategory{},
  &productmodel.Product{},
  &productmodel.ProductSku{},
  &productmodel.ProductImage{},
  &productmodel.ProductFavorite{},
)
```

- [ ] **Step 4: 运行后端编译验证**

Run: `go test ./plugins/product/...`
Expected: PASS（至少编译通过）。

### Task 2: 收藏服务与商品查询扩展

**Files:**
- Create: `server/plugins/product/service/favorite.go`
- Modify: `server/plugins/product/service/product.go`

- [ ] **Step 1: 新增收藏/取消收藏服务（幂等 + 事务计数）**

```go
func FavoriteProduct(ctx context.Context, userID, productID uint64) error
func UnfavoriteProduct(ctx context.Context, userID, productID uint64) error
func ListUserFavorites(ctx context.Context, userID uint64, page, size int) ([]FavoriteProductItem, int64, error)
```

关键逻辑：
- 收藏时插入关系成功才 `favorite_count + 1`。
- 唯一冲突直接返回 `nil`。
- 取消时删除成功才 `favorite_count - 1`，并用 SQL 下限保护不小于 0。

- [ ] **Step 2: 扩展商品列表/详情返回结构**

```go
type ProductListItem struct {
  productmodel.Product
  IsFavorited bool `json:"is_favorited"`
}

type ProductDetail struct {
  productmodel.Product
  SKUs []productmodel.ProductSku `json:"skus"`
  Images []productmodel.ProductImage `json:"images"`
  IsFavorited bool `json:"is_favorited"`
}
```

- [ ] **Step 3: 在列表和详情注入收藏态**

```go
func ListProducts(ctx context.Context, q ProductListQuery, userID uint64) ([]ProductListItem, int64, error)
func GetProduct(ctx context.Context, id uint64, userID uint64) (*ProductDetail, error)
```

- [ ] **Step 4: 运行后端编译验证**

Run: `go test ./plugins/product/...`
Expected: PASS。

### Task 3: 前台 API 路由与处理器

**Files:**
- Modify: `server/plugins/product/api/front.go`

- [ ] **Step 1: 在现有接口中透传用户身份（可选）**

```go
func currentUserID(c *gin.Context) uint64 {
  auth := c.GetHeader("Authorization")
  if !strings.HasPrefix(auth, "Bearer ") { return 0 }
  claims, err := middleware.ParseToken(strings.TrimPrefix(auth, "Bearer "))
  if err != nil { return 0 }
  return claims.UserID
}
```

- [ ] **Step 2: 调整商品列表/详情调用新服务签名**

```go
uid := currentUserID(c)
list, total, err := productsvc.ListProducts(c.Request.Context(), q, uid)
detail, err := productsvc.GetProduct(c.Request.Context(), id, uid)
```

- [ ] **Step 3: 新增鉴权收藏路由**

```go
auth := g.Group("")
auth.Use(middleware.RequireAuth)
auth.POST("/products/:id/favorite", favoriteProduct)
auth.DELETE("/products/:id/favorite", unfavoriteProduct)
auth.GET("/user/favorites", listUserFavorites)
```

- [ ] **Step 4: 运行后端全量验证**

Run: `go test ./...`
Expected: PASS。

### Task 4: app 端收藏功能

**Files:**
- Modify: `app/pages/product/detail.vue`
- Modify: `app/pages/user/index.vue`
- Create: `app/pages/user/favorites.vue`
- Modify: `app/pages.json`

- [ ] **Step 1: 详情页增加收藏按钮与登录校验**

```ts
async function toggleFavorite() {
  const token = uni.getStorageSync('user_token')
  if (!token) { uni.navigateTo({ url: '/pages/login/index' }); return }
  if (product.value.is_favorited) {
    await del(`/api/v1/products/${product.value.id}/favorite`)
    product.value.is_favorited = false
    product.value.favorite_count = Math.max(0, Number(product.value.favorite_count || 0) - 1)
  } else {
    await post(`/api/v1/products/${product.value.id}/favorite`)
    product.value.is_favorited = true
    product.value.favorite_count = Number(product.value.favorite_count || 0) + 1
  }
}
```

- [ ] **Step 2: 个人中心新增收藏入口**

```ts
{ label: '我的收藏', icon: 'heart', value: '', action: () => uni.navigateTo({ url: '/pages/user/favorites' }) }
```

- [ ] **Step 3: 新建收藏列表页并实现“取消即移除”**

```ts
const data = await get<any>('/api/v1/user/favorites', { page: page.value, size: size })
list.value = data.list || []

async function unfavorite(item: any, index: number) {
  await del(`/api/v1/products/${item.id}/favorite`)
  list.value.splice(index, 1)
}
```

- [ ] **Step 4: 注册页面路由并做 app 构建验证**

Run: `npm run build:h5:demo`
Expected: 构建成功。

### Task 5: web 端收藏功能

**Files:**
- Modify: `web/src/views/ProductDetail.vue`
- Modify: `web/src/views/UserCenter.vue`

- [ ] **Step 1: 详情页增加收藏按钮与未登录跳转**

```ts
if (!auth.isLoggedIn) {
  router.push('/login')
  return
}
```

并按收藏状态调用：
- `POST /api/v1/products/:id/favorite`
- `DELETE /api/v1/products/:id/favorite`

- [ ] **Step 2: 用户中心加入“我的收藏”菜单与列表加载**

```ts
const favorites = ref<any[]>([])
const favoriteQuery = ref({ page: 1, size: 12, total: 0 })
favorites.value = data.list || []
```

- [ ] **Step 3: 取消收藏后立即移除当前项**

```ts
await del(`/api/v1/products/${id}/favorite`)
favorites.value = favorites.value.filter((row: any) => Number(row.id) !== id)
```

- [ ] **Step 4: web 构建验证**

Run: `npm run build`
Expected: 构建成功。

### Task 6: admin 商品列表展示收藏数

**Files:**
- Modify: `admin/src/views/product/ProductList.vue`

- [ ] **Step 1: 表头新增收藏数列**

```html
<th class="px-4 py-3 text-left text-slate-500 font-medium">收藏数</th>
```

- [ ] **Step 2: 行数据渲染收藏数**

```html
<td class="px-4 py-3 text-slate-700">{{ p.favorite_count || 0 }}</td>
```

- [ ] **Step 3: 调整空状态 colspan 与构建验证**

Run: `npm run build`
Expected: 构建成功。

### Task 7: 三端 mock 与 docs-site 文档同步

**Files:**
- Modify: `app/mock/index.ts`
- Modify: `web/src/mock/index.ts`
- Modify: `admin/src/mock/index.ts`
- Modify: `docs-site/docs/api/product.md`

- [ ] **Step 1: app/web mock 增加收藏接口与收藏字段**

新增匹配：
- `POST /api/v1/products/:id/favorite`
- `DELETE /api/v1/products/:id/favorite`
- `GET /api/v1/user/favorites`

并在商品详情/列表返回中补：`favorite_count`、`is_favorited`。

- [ ] **Step 2: admin mock 商品列表补 `favorite_count`**

确保 `GET /admin/api/products` 返回每个商品带 `favorite_count`。

- [ ] **Step 3: 更新 docs-site 商品 API 文档**

补充：
- 功能说明：SPU 收藏、登录门槛、收藏列表行为。
- 接口变化：新增 3 个接口 + 现有 3 个接口字段升级。
- 部署影响：新增表和字段迁移。

- [ ] **Step 4: docs-site 构建验证**

Run: `npm run docs:build`
Expected: 构建成功。

### Task 8: 最终联调与回归检查

**Files:**
- Modify:（按联调结果修复最小差异）

- [ ] **Step 1: 运行后端与三前端最终构建回归**

Run:
- `go test ./...`（`server`）
- `npm run build:h5:demo`（`app`）
- `npm run build`（`web`）
- `npm run build`（`admin`）

Expected: 全部通过。

- [ ] **Step 2: 手工行为回归清单**

- app/web 商品详情：收藏 -> 已收藏，取消 -> 未收藏。
- app/web 我的收藏：取消后立即从列表移除。
- admin 商品列表：收藏数展示正常。

- [ ] **Step 3: 输出变更说明（不执行 git commit）**

说明按仓库约定准备中文提交信息草案（`head + body`），但本次不实际提交。
