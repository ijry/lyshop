# 订单接口

## 说明

订单模块覆盖购物车、创建订单、地址、订单查询与状态流转。

## 典型接口

- `GET /api/v1/cart`
- `POST /api/v1/cart/add`
- `PUT /api/v1/cart/qty`
- `DELETE /api/v1/cart/:sku_id`
- `GET /api/v1/addresses`
- `POST /api/v1/addresses`
- `PUT /api/v1/addresses/:id`
- `DELETE /api/v1/addresses/:id`
- `GET /api/v1/orders?status=<0|1|2|3|4|5|6>&page=1&size=20`
- `GET /api/v1/orders/:id`
- `GET /api/v1/orders/:id/review`
- `POST /api/v1/orders/:id/pay`
- `POST /api/v1/orders/:id/cancel`
- `POST /api/v1/orders/:id/review`
- `POST /api/v1/orders/:id/after-sales`
- `GET /api/v1/after-sales/:id`
- `POST /api/v1/after-sales/:id/return-shipments`
- `POST /api/v1/upload`
- `GET /admin/api/orders?status=<0|1|2|3|4|5|6>&page=1&size=20`
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
- `after_sale_summary` 兼容返回 `latest_status_label`，用于“最近售后单”优先直显中文状态；无该字段时可回退 `latest_status` 本地映射。
- `shipments` 单条轨迹包含方向（`direction`）、业务类型（`biz_type`）、物流状态（`logistics_status`）、时间字段（`shipped_at/signed_at/created_at`）、备注（`remark`）和关联售后单（`after_sale_case_id`），并兼容返回 `direction_label`、`biz_type_label`、`logistics_status_label`。
- `latest_shipment` 与 `shipments` 字段口径一致，同样兼容返回上述标签字段，便于列表快速展示中文物流摘要。
- 前台与后台订单列表均可直接使用 `latest_shipment` 展示最新物流状态，并通过 `shipments[*].biz_type=reship` 标注补发摘要。
- 前台与后台建议统一维护状态文案映射：订单状态、售后状态、物流状态与日志状态流转文案应保持一致。
- 订单状态当前定义：`1待付款`、`2待发货`、`3待收货`、`4已完成`、`5售后`、`6已取消`。
- 待付款订单支持用户主动取消：`POST /api/v1/orders/:id/cancel`。
- 订单库存以 WMS 为真源：下单预占、支付确认、取消释放，详见 [库存预占规则](./stock-reservation)。
- 售后相关查询接口在保留枚举字段的同时，兼容返回可读标签字段（`*_label`），用于前端直接展示中文文案。
- 后台商品列表“管理评价”弹窗复用现有评价接口（`GET /admin/api/reviews`、`GET /admin/api/reviews/:id`、`POST /admin/api/reviews/:id/reply`），通过 `product_id` + `keyword` 组合筛选当前商品评价，无新增接口。
- 管理端“管理评价”按钮按权限显隐：优先校验 JWT 权限码 `order:view`，并兼容以菜单 `/review/list` 作为兜底判断，避免无权限账号误触。
- 评价“回复”操作按 `order:review-reply` 权限显隐与拦截，无该权限时不展示回复入口且不可提交回复。
- 评价详情弹窗在无 `order:review-reply` 权限时会展示只读提示文案，明确当前账号仅可查看。
- 若通过非常规方式触发回复动作（如旧状态页面），前端会弹出“无评价回复权限”提示并阻断提交。

## 活动来源下单兼容升级

为支持活动商品来源强约束（`activity_product_id`），购物车与下单接口做了兼容升级：

- `POST /api/v1/cart/add`
  - 新增可选参数：`activity_product_id`
- `PUT /api/v1/cart/qty`
  - 新增可选参数：`activity_product_id`
- `DELETE /api/v1/cart/:sku_id`
  - 新增可选参数：`activity_product_id`
- `POST /api/v1/orders`
  - 新增 `items`：`[{ sku_id, activity_product_id }]`
  - 当请求同时带 `items` 与 `sku_ids` 时，后端优先按 `items` 下单；`sku_ids` 继续兼容旧调用

`POST /api/v1/orders` 请求体示例（活动来源）：

```json
{
  "address_id": 1,
  "payment_method": "wechat",
  "items": [
    { "sku_id": 101, "activity_product_id": 9001 },
    { "sku_id": 102, "activity_product_id": 0 }
  ],
  "sku_ids": [101, 102],
  "remark": "请尽快发货"
}
```

说明：

- `activity_product_id > 0`：走活动来源校验与活动价结算。
- `activity_product_id = 0`：按普通商品流程结算，不自动匹配活动价。

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
- `company` 建议传标准快递公司编码（如 `SF/ZTO/YTO/STO/YD/JD/EMS/DBL/JT`），后台发货页使用下拉字典维护，避免手填错误。
- 每次发货都会写入 `shipments` 轨迹，`tracking_no` 仍保留最近单号以兼容旧端。

## 物流驱动化升级（快递100 / 快递鸟）

- 系统新增物流驱动机制，支持 `kuaidi100` 与 `kdniao` 两个渠道。
- 运单首次同步按“主驱动 -> 备驱动”顺序尝试；首次成功后会写入 `channel_provider` 并固定后续同步渠道。
- 运单同步同时维护：
  - `order_shipments` 最新状态字段（如 `logistics_status`、`signed_at`）
  - `order_shipment_tracks` 轨迹节点明细（时间线展示与审计）

### 后台手动刷新接口

- `POST /admin/api/orders/:id/shipments/:shipment_id/sync`
  - 作用：按运单触发一次立即同步

### 轨迹查询接口

- `GET /admin/api/orders/:id/shipments/:shipment_id/tracks`
  - 作用：后台查看完整轨迹节点
- `GET /api/v1/orders/:id/shipments/:shipment_id/tracks`
  - 作用：用户侧查看完整轨迹节点

轨迹节点返回示例：

```json
[
  {
    "id": 101,
    "shipment_id": 12,
    "provider": "kuaidi100",
    "track_hash": "8b8a2d...",
    "status_code": "in_transit",
    "status_text": "快件已到达杭州转运中心",
    "event_time": "2026-05-25T12:00:00+08:00",
    "location": "杭州",
    "raw_payload": {}
  }
]
```

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
- 订单库存交易依赖 WMS 预占模型，需确保 WMS 迁移已执行（`reserved_qty` 与 `inventory_reservation`）。
- 活动来源链路新增后，`order_items` 会新增活动快照字段（`activity_product_id`、`activity_id`、`activity_type`、`activity_title`），需依赖服务启动自动迁移。
- Redis 购物车键由单一 `sku_id` 兼容升级为 `sku_id:activity_product_id` 复合键（旧键可兼容读取），无新增配置项。
- 后台商品列表评价弹窗仅涉及管理端前端交互改造，服务端接口、数据库结构与配置项均无变化。
- 购物车勾选结算仅涉及前端交互与参数组装，无后端部署、迁移与配置变更。
- 数据库会新增售后与物流相关表，需执行服务启动时的插件迁移。
- 文件上传仍复用统一 `POST /api/v1/upload`（前台）与 `POST /admin/api/upload`（后台）入口；若启用云存储，需先在后台完成对应插件配置（Endpoint/Region/Bucket/密钥/域名）。
- 系统支持同时启用多个存储驱动，并在 `storage_router` 插件配置默认驱动；请求也可通过 `driver` 参数临时指定驱动（`local/oss/cos/qiniu`，兼容 `aliyun_oss/qcloud_cos/qiniu_kodo`）。
- 新增物流路由插件配置项（`logistics_router`）：
  - `enabled`：物流路由总开关
  - `primary_driver` / `secondary_driver`：主备驱动
  - `polling_enabled`：自动轮询开关
  - `polling_interval_seconds`：轮询频率（秒）

## 商家移动端订单增强接口

> 用于商家移动端 eapp 的批量与高频操作，统一前缀 `/admin/api`。

### POST /orders/{id}/repricing

订单改价（仅 status=1 可用）。

请求：`{ items: [{ item_id, price }], remark }`
返回：`{ id, amount_breakdown: { goods_amount, discount_amount, payable_amount } }`

### POST /orders/{id}/notes

订单备注。

请求：`{ content, visible_to?: 'merchant_only' }`
返回：`{ id, notes: [...] }`

### POST /orders/{id}/remind-pay

催付款（短信 / 微信）。

请求：`{ channel: 'sms' | 'wx' }`
返回：`{ sent_at, channel }`

### GET /orders/{id}/print-template

电子面单模板。

返回：`{ template: '<html>...</html>' }`

### GET /orders/{id}/timeline

订单时间线。

返回：`[{ stage, status, time, content }]`

### POST /orders/batch/ship | notes | repricing | close

批量操作系列。`ship` 接收数组 `[{order_id,company,tracking_no}]`；其他统一 `{ ids[], ... }`。

返回：`{ success_ids: number[], fail: [{ id, reason }] }`

### GET /orders?keyword=&time_start=&time_end=&amount_min=&amount_max=&logistics_company=&province=&pay_method=&has_after_sale=

列表接口扩展查询参数。其它字段与原接口一致。
  - `polling_batch_size`：单批处理数量
