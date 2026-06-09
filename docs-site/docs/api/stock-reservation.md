# 库存预占交易规则

## 功能说明

当前库存架构统一由 `inventory` provider 抽象，支持三种模式：

1. `local`
   - 商品库存直接由 `product_skus.stock` 管理。
   - 适用于不启用 WMS 的轻量商城部署。
2. `builtin_wms`
   - 订单库存交易由内置 `wms` 插件处理。
   - 使用 `inventory_stock`、`reserved_qty` 与预占流水模型。
3. `external_wms`
   - 库存交易委托企业外部 WMS。
   - 可按配置选择 `sync` 或 `async` 两种模式。

订单库存链路统一采用三段式交易：

1. `Reserve`（下单预占）
2. `Confirm`（支付确认扣减）
3. `Release`（取消/超时释放）

## 交易接口边界

- 用户下单：`POST /api/v1/orders`
  - 服务端在创建订单时触发当前 provider 的库存预占。
- 用户支付：`POST /api/v1/orders/:id/pay`
  - 服务端在支付成功后确认预占并扣减库存。
- 用户取消：`POST /api/v1/orders/:id/cancel`
  - 服务端释放该订单对应预占库存。

> 预占/确认/释放为服务内能力（`order -> inventory provider`），不新增对外公开订单库存 REST 接口。

### 外部 WMS 模式

- `inventory.provider=external_wms`
  - `inventory.external_mode=sync`
    - 下单时实时调用企业 WMS；外部预占失败则订单创建失败。
  - `inventory.external_mode=async`
    - 商城先落订单与库存意图，再异步投递外部 WMS。
    - 异步任务记录在 `inventory_integration_tasks`，支持 claim、重试、回调补偿与最终失败落库。

## 异步外部库存任务

当 `inventory.provider=external_wms` 且 `inventory.external_mode=async` 时，统一库存会把库存动作写入异步任务表，再由后台任务执行。

任务状态：

- `pending`
  - 待执行或待重试。
- `processing`
  - 已被 worker 领取，正在执行或等待外部 WMS 回调确认。
- `success`
  - 外部库存动作已完成，订单库存状态会同步收口为 `reserved / confirmed / released`。
- `failed`
  - 外部 WMS 明确返回失败，或任务超过最大重试次数后最终失败。

异步任务关键字段：

- `attempt_count`
- `max_attempts`
- `backoff_seconds`
- `next_retry_at`
- `lock_owner`
- `lock_expires_at`
- `completed_at`
- `last_callback_id`

回调规则：

- 同一 `callback_id` 重复回调会被幂等忽略。
- 已进入最终态（`success / failed`）的任务不会被后续回调回滚。
- 订单侧 `inventory_status` 会根据异步任务最终结果同步更新。

## 单仓优先策略

- `builtin_wms` 默认按“单仓优先”选仓完成预占。
- 订单明细记录仓库快照，后续确认/释放在同仓执行，避免跨仓不一致。
- `local` 模式不涉及仓库维度。
- `external_wms` 是否使用仓库维度由企业 WMS 自身决定。

## 数据模型

- `inventory_reservation`
  - `biz_type`、`biz_no`、`sku_id`、`qty`、`status`
  - `status`：`reserved | confirmed | released`
- `order_inventory_state`
  - `order_no`、`provider`、`status`、`last_error`
- `inventory_integration_tasks`
  - `provider`、`biz_type`、`biz_no`、`action`、`payload`、`status`
  - `attempt_count`、`max_attempts`、`backoff_seconds`、`next_retry_at`
  - `lock_owner`、`lock_expires_at`、`completed_at`、`last_callback_id`
- `inventory_stock`
  - 仅 `builtin_wms` 模式使用
  - `qty`：在手库存
  - `reserved_qty`：预占库存

## 部署与配置影响

- 服务启动时会自动迁移统一 inventory 共享表：
  - `inventory_reservation`
  - `order_inventory_state`
  - `inventory_integration_tasks`
- `builtin_wms` 模式需额外执行 WMS 自动迁移，确保新增字段与表已创建：
  - `inventory_stock.reserved_qty`
  - WMS 预占相关表
- 配置项新增：
  - `inventory.provider`
  - `inventory.external_mode`
  - `external_wms.endpoint/app_key/app_secret`
  - `external_wms.signature_ttl`
  - `external_wms.worker_interval_sec`
