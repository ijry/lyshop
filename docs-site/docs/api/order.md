# 订单接口

## 说明

订单模块覆盖购物车、创建订单、地址、订单查询与状态流转。

## 典型接口

- `POST /api/order/create`
- `GET /api/order/list`
- `GET /api/order/detail`
- `POST /api/order/cancel`

## 说明

- 下单链路建议先校验库存与活动信息
- 支付成功后需同步更新订单状态
