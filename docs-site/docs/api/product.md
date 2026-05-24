# 商品接口

## 说明

商品模块负责商品、分类、SKU、图片等基础交易数据。

## 典型接口

- `GET /api/product/list`
- `GET /api/product/detail`
- `GET /api/product/category/list`
- `GET /api/product/sku/list`
- `GET /admin/api/categories`
- `POST /admin/api/categories`
- `PUT /admin/api/categories/:id`
- `DELETE /admin/api/categories/:id`
- `GET /api/v1/products/recommend?limit=8`
- `GET /api/v1/products/:id/reviews?page=1&size=10`

## 说明

- 前台接口用于商品浏览与详情展示
- 后台接口用于商品维护、上架与库存信息管理
- 后台“商品分类”页面复用现有 `/admin/api/categories` 资源接口完成增删改与启停，不新增平行接口
- 后台分类列表返回全部分类（含停用），前台分类接口仍仅返回启用分类（`status=1`）
- `GET /api/v1/products/recommend` 返回上架商品推荐列表，默认 8 条，支持 `limit` 参数（最大 50）
- 商品详情页的“评价 Tab”通过 `GET /api/v1/products/:id/reviews` 获取评分摘要、评价列表、追评与商家回复（含评价图片与追评图片）

## 部署与配置影响

- 无新增环境变量。
- 无新增中间件或外部依赖。
- 仅调整后台分类列表查询口径（含停用）与后台页面路由，无需额外部署步骤。
