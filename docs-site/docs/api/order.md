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
- `POST /api/v1/orders/:id/after-sales`
- `GET /api/v1/after-sales/:id`
- `POST /api/v1/after-sales/:id/return-shipments`
- `GET /admin/api/orders?status=<0|1|2|3|4|5>&page=1&size=20`
- `GET /admin/api/orders/:id`
- `PUT /admin/api/orders/:id/ship`
- `GET /admin/api/after-sales?status=&case_type=&order_id=&page=&size=`
- `GET /admin/api/after-sales/:id`
- `POST /admin/api/after-sales/:id/audit`
- `POST /admin/api/after-sales/:id/receive`
- `POST /admin/api/after-sales/:id/refund`
- `POST /admin/api/after-sales/:id/complete`
- `POST /admin/api/after-sales/:id/close`
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
- 订单列表/详情升级返回 `shipments`、`latest_shipment`、`after_sale_summary`，并继续保留 `tracking_no` 兼容旧前端。

## 发货接口升级（兼容）

`PUT /admin/api/orders/:id/ship`

请求体示例：

```json
{
  "tracking_no": "SF1234567890",
  "ship_type": "initial",
  "company": "顺丰",
  "remark": "首发",
  "after_sale_case_id": 0
}
```

- `ship_type=initial`：首发，保持原有发货语义并更新订单状态为待收货。
- `ship_type=reship`：补发，需传 `after_sale_case_id`，用于换货售后补发闭环。
- 每次发货都会写入 `shipments` 轨迹，`tracking_no` 仍保留最近单号以兼容旧端。

## 售后接口（退货/换货）

### 用户申请售后

`POST /api/v1/orders/:id/after-sales`

```json
{
  "case_type": "return",
  "reason": "尺寸不合适",
  "apply_content": "试穿后不合适",
  "apply_images": ["https://cdn.example.com/a.jpg"],
  "items": [
    { "order_item_id": 11, "qty": 1 }
  ]
}
```

- `case_type` 支持 `return`（退货退款）和 `exchange`（换货）。
- 同一订单商品存在进行中售后时，不可重复申请。

### 用户查询与回寄物流

- `GET /api/v1/after-sales/:id`：返回售后详情、状态日志、关联物流。
- `POST /api/v1/after-sales/:id/return-shipments`：用户提交回寄物流（快递公司、单号、备注）。

### 后台售后动作

- `POST /admin/api/after-sales/:id/audit`：审核通过/拒绝
- `POST /admin/api/after-sales/:id/receive`：仓库确认收货
- `POST /admin/api/after-sales/:id/refund`：登记退款
- `POST /admin/api/after-sales/:id/complete`：售后完结
- `POST /admin/api/after-sales/:id/close`：关闭售后

状态流：

- 退货：`applied -> approved_wait_user_return -> user_returning -> refund_pending -> refunded -> completed`
- 换货：`applied -> approved_wait_user_return -> user_returning -> reship_pending -> reshipped -> completed`

## 评价接口

`POST /api/v1/orders/:id/review`

请求体示例：

```json
{
  "mode": "create",
  "logistics_score": 5,
  "items": [
    { "order_item_id": 11, "product_score": 5, "content": "做工很好", "images": ["https://cdn.example.com/a.jpg"] }
  ],
  "append_content": "",
  "append_images": []
}
```

- `mode=create` 创建根评价
- `mode=edit` 覆盖原评价
- `mode=append` 追加到根评价子级，必须先存在对应根评价
- 根评价与追加评价都支持 `images`
- 订单评价页会先加载 `GET /api/v1/orders/:id/review` 获取当前评价状态
- 前台上传文件统一使用 `POST /api/v1/upload`

## 部署与配置影响

- 本次仅为订单与售后能力增强，无新增部署步骤，无新增环境变量。
- 数据库会新增售后与物流相关表，需执行服务启动时的插件迁移。
- 文件上传仍复用统一 `POST /api/v1/upload` 入口；若启用 `storage_oss`、`storage_cos` 或 `storage_qiniu`，需先在后台完成对应插件配置（Endpoint/Region/Bucket/密钥/域名），并且 `plugins.enabled` 仅保留一个存储驱动。
- 切换存储驱动后，前台与后台返回的文件 URL 仍保持统一接口语义，不需要改动调用方。
