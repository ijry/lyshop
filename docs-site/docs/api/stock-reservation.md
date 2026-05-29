# 库存预占交易规则

## 功能说明

当前库存真源为 `wms`。订单库存链路采用三段式交易：

1. `Reserve`（下单预占）
2. `Confirm`（支付确认扣减）
3. `Release`（取消/超时释放）

商品侧 `product_skus.stock` 不再作为扣减真源，订单库存一致性由 WMS 负责。

## 交易接口边界

- 用户下单：`POST /api/v1/orders`
  - 服务端在创建订单时触发 WMS 预占。
- 用户支付：`POST /api/v1/orders/:id/pay`
  - 服务端在支付成功后确认预占并扣减在手库存。
- 用户取消：`POST /api/v1/orders/:id/cancel`
  - 服务端释放该订单对应预占库存。

> 预占/确认/释放为服务内能力（order -> wms service），不新增对外公开 REST 端点。

## 单仓优先策略

- 默认按“单仓优先”选仓完成预占。
- 订单明细记录仓库快照，后续确认/释放在同仓执行，避免跨仓不一致。

## 数据模型

- `inventory_stock`
  - `qty`：在手库存
  - `reserved_qty`：预占库存
- `inventory_reservation`
  - `biz_type`、`biz_no`、`sku_id`、`warehouse_id`、`qty`、`status`
  - `status`：`reserved | confirmed | released`

## 部署与配置影响

- 需执行 WMS 自动迁移，确保新增字段与表已创建：
  - `inventory_stock.reserved_qty`
  - `inventory_reservation`
- 无新增环境变量。
