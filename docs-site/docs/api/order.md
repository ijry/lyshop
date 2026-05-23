# 订单接口

## 说明

订单模块覆盖购物车、创建订单、地址、订单查询与状态流转。

## 典型接口

- `POST /api/order/create`
- `GET /api/order/list`
- `GET /api/order/detail`
- `POST /api/order/cancel`
- `GET /api/v1/orders?status=<0|1|2|3|4>&page=1&size=20`

## 说明

- 下单链路建议先校验库存与活动信息
- 支付成功后需同步更新订单状态
- 前端订单 tab 切换复用同一查询接口，通过 `status` 参数过滤结果，不新增额外接口。
