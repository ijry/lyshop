# 商品接口

## 说明

商品模块负责商品、分类、SKU、图片等基础交易数据。

## 典型接口

- `GET /api/product/list`
- `GET /api/product/detail`
- `GET /api/product/category/list`
- `GET /api/product/sku/list`
- `GET /api/v1/products/:id/reviews?page=1&size=10`

## 说明

- 前台接口用于商品浏览与详情展示
- 后台接口用于商品维护、上架与库存信息管理
- 商品详情页的“评价 Tab”通过 `GET /api/v1/products/:id/reviews` 获取评分摘要、评价列表、追评与商家回复
