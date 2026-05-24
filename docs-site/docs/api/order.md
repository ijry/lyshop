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
- `POST /api/v1/upload`
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
- `shipments` 单条轨迹包含方向（`direction`）、业务类型（`biz_type`）、物流状态（`logistics_status`）、时间字段（`shipped_at/signed_at/created_at`）、备注（`remark`）和关联售后单（`after_sale_case_id`）。
- 前台与后台订单列表均可直接使用 `latest_shipment` 展示最新物流状态，并通过 `shipments[*].biz_type=reship` 标注补发摘要。
- 前台与后台建议统一维护状态文案映射：订单状态、售后状态、物流状态与日志状态流转文案应保持一致。
- 售后相关查询接口在保留枚举字段的同时，兼容返回可读标签字段（`*_label`），用于前端直接展示中文文案。

## 购物车勾选结算（兼容升级）

- 购物车结算继续复用既有确认页入口：`/pages/order/confirm?sku_ids=...`，不新增接口。
- `sku_ids` 语义升级为“购物车勾选商品集合”，不再默认等于购物车全量商品。
- 购物车列表初始化默认全选；用户取消全选或未选中任何商品时，不允许进入结算流程。
- PC/H5/小程序/App 前端未勾选商品时，应在客户端拦截并提示后再发起下单。

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
- `GET /api/v1/after-sales/:id` 中的 `shipments` 与订单详情轨迹字段口径一致，可直接展示方向、业务类型、状态、时间、备注与关联售后单。
- `GET /api/v1/after-sales/:id` 兼容返回以下标签字段：
  - 售后单：`status_label`、`case_type_label`
  - 日志：`from_status_label`、`to_status_label`、`action_label`
  - 物流：`direction_label`、`biz_type_label`、`logistics_status_label`

### 后台售后动作

- `POST /admin/api/after-sales/:id/audit`：审核通过/拒绝
- `POST /admin/api/after-sales/:id/receive`：仓库确认收货
- `POST /admin/api/after-sales/:id/refund`：登记退款
- `POST /admin/api/after-sales/:id/complete`：售后完结
- `POST /admin/api/after-sales/:id/close`：关闭售后
- `GET /admin/api/after-sales` 与 `GET /admin/api/after-sales/:id` 会兼容返回同口径标签字段（`status_label`、`case_type_label`、`from_status_label`、`to_status_label`、`action_label`、`direction_label`、`biz_type_label`、`logistics_status_label`）。

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
- 购物车勾选结算仅涉及前端交互与参数组装，无后端部署、迁移与配置变更。
- 数据库会新增售后与物流相关表，需执行服务启动时的插件迁移。
- 文件上传仍复用统一 `POST /api/v1/upload`（前台）与 `POST /admin/api/upload`（后台）入口；若启用云存储，需先在后台完成对应插件配置（Endpoint/Region/Bucket/密钥/域名）。
- 系统支持同时启用多个存储驱动，并在 `storage_router` 插件配置默认驱动；请求也可通过 `driver` 参数临时指定驱动（`local/oss/cos/qiniu`，兼容 `aliyun_oss/qcloud_cos/qiniu_kodo`）。
