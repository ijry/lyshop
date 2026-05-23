# 订单接口

## 说明

订单模块覆盖购物车、创建订单、地址、订单查询与状态流转。

## 典型接口

- `GET /api/v1/orders?status=<0|1|2|3|4>&page=1&size=20`
- `GET /api/v1/orders/:id`
- `GET /admin/api/orders?status=<0|1|2|3|4|5>&page=1&size=20`
- `GET /admin/api/orders/:id`
- `PUT /admin/api/orders/:id/ship`

## 说明

- 下单链路建议先校验库存与活动信息
- 支付成功后需同步更新订单状态
- 前端订单 tab 切换复用同一查询接口，通过 `status` 参数过滤结果，不新增额外接口。
- 列表与详情统一返回 `items` 和 `amount_breakdown`，用于展示商品明细与价格体系。
