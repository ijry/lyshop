# 订单接口

## 说明

订单模块覆盖购物车、创建订单、地址、订单查询与状态流转。

## 典型接口

- `GET /api/v1/addresses`
- `POST /api/v1/addresses`
- `PUT /api/v1/addresses/:id`
- `DELETE /api/v1/addresses/:id`
- `GET /api/v1/orders?status=<0|1|2|3|4>&page=1&size=20`
- `GET /api/v1/orders/:id`
- `GET /api/v1/orders/:id/review`
- `POST /api/v1/orders/:id/pay`
- `POST /api/v1/orders/:id/review`
- `GET /admin/api/orders?status=<0|1|2|3|4|5>&page=1&size=20`
- `GET /admin/api/orders/:id`
- `PUT /admin/api/orders/:id/ship`
- `GET /admin/api/reviews?product_id=&keyword=&page=&size=`
- `GET /admin/api/reviews/:id`
- `POST /admin/api/reviews/:id/reply`

## 说明

- 下单链路建议先校验库存与活动信息
- 支付成功后需同步更新订单状态
- 地址新增、编辑、删除统一复用 `/api/v1/addresses` 资源语义，不新增平行接口。
- 前台“去付款/评价”动作统一复用订单资源子动作：`/pay`、`/review`。
- 前端订单 tab 切换复用同一查询接口，通过 `status` 参数过滤结果，不新增额外接口。
- 列表与详情统一返回 `items` 和 `amount_breakdown`，用于展示商品明细与价格体系。

## 评价接口

`POST /api/v1/orders/:id/review`

请求体示例：

```json
{
  "mode": "create",
  "logistics_score": 5,
  "items": [
    { "order_item_id": 11, "product_score": 5, "content": "做工很好" }
  ],
  "append_content": ""
}
```

- `mode=create` 创建根评价
- `mode=edit` 覆盖原评价
- `mode=append` 追加到根评价子级
- 订单评价页会先加载 `GET /api/v1/orders/:id/review` 获取当前评价状态

## 部署与配置影响

- 本次仅为接口与前端交互补齐，无新增部署步骤，无新增环境变量。
