# SKU 自动笛卡尔积与 WMS 一体化设计

## 1. 背景

当前商品 SKU 与库存链路存在三类核心问题：

1. SKU 生成主要依赖前端矩阵编辑与前端笛卡尔积，后端以整包替换方式保存，SKU ID 稳定性差。
2. 订单下单仍直接扣减 `product_skus.stock`，与 WMS 库存台账/流水未打通。
3. WMS 现阶段仅支持入/出库单完成后变更库存，缺少订单场景的预占库存能力。

## 2. 目标

1. 支持属性组配置后自动生成 SKU 笛卡尔积，且保证 SKU 主键稳定。
2. 以 WMS 作为库存唯一真源，订单全链路走预占/确认/释放。
3. 保持单仓优先履约策略，确保规则简单、可控。
4. 对既有接口优先做兼容升级，避免无节制新增接口。
5. `app`、`eapp`、`admin` 三端同步改造，统一消费新库存语义与 SKU 结构。

## 3. 业务决策（已确认）

1. 库存真源：`WMS`。
2. 规格变更：旧 SKU 采用软删除（下线）策略，不硬删。
3. 履约策略：单仓优先，不做自动多仓拆单。

## 4. 方案比较

### 4.1 方案 A：前端笛卡尔积 + 后端整包替换

- 优点：改造小，上线快。
- 缺点：SKU ID 不稳定，无法优雅软删除，库存与 WMS 仍弱耦合。
- 结论：不推荐。

### 4.2 方案 B：后端规格引擎 + 差异化 SKU 生成 + WMS 预占（推荐）

- 优点：SKU 稳定、库存一致性强、可防超卖、可持续演进。
- 缺点：需改造商品保存与订单库存链路。
- 结论：推荐。

### 4.3 方案 C：双写最终一致

- 优点：短期落地较快。
- 缺点：串账风险高、对账与补偿复杂。
- 结论：不建议。

## 5. 总体架构（推荐方案）

1. 商品侧录入规格属性组（颜色、尺码等）。
2. 后端规格引擎生成规范化组合键 `sku_key`，执行差异化 Upsert：
   - 保留：同 `sku_key` 复用原 SKU。
   - 新增：新组合创建新 SKU。
   - 下线：缺失组合标记 `inactive`（软删除）。
3. WMS 维护 `on_hand` 与 `reserved`，可售库存统一计算：
   - `sellable = on_hand - reserved`
4. 订单创建时先预占库存；支付成功后确认扣减；取消/超时释放预占。

## 6. 数据模型设计

### 6.1 商品 SKU 表 `product_skus`（升级）

新增字段：

- `sku_key`：规格组合规范化键（唯一，按 `product_id + sku_key` 唯一）。
- `status`：`active|inactive`。
- `deleted_at`：软删除时间（可选，按项目模型习惯）。

说明：

- `attrs` 继续沿用 JSON 存储。
- 历史订单、WMS 流水通过 `sku_id` 可持续追溯。

### 6.2 WMS 库存表 `inventory_stock`（升级）

新增字段：

- `reserved_qty`：预占数量。

含义：

- `qty` 视为在手库存（on hand）。
- 可售库存为 `qty - reserved_qty`。

### 6.3 预占记录表 `inventory_reservations`（新增）

建议字段：

- `id`
- `biz_type`（order）
- `biz_no`（订单号）
- `warehouse_id`
- `sku_id`
- `qty`
- `status`（reserved|confirmed|released）
- `expired_at`
- `created_at/updated_at`

唯一键建议：`(biz_type, biz_no, sku_id)`，用于幂等。

### 6.4 库存流水 `inventory_movement`（扩展）

`biz_type` 增补值：

- `reserve`
- `release`
- `order_outbound`

并保留 `doc_no` 或扩展 `biz_no` 以承载订单号。

## 7. 接口设计

### 7.1 商品接口（优先升级现有）

升级已有接口：

- `PUT /admin/api/products/:id`

新增可选输入：

- `spec_schema`：属性组定义。
- `sku_overrides`：按 `sku_key` 的价格/条码/重量覆盖。
- `sku_generation_mode=auto`。

新增返回摘要：

- `sku_diff`：`added/kept/inactivated`。

说明：不强制新增商品接口，通过升级现有保存接口完成自动生成。

### 7.2 WMS 服务接口（内部服务能力）

新增服务方法：

- `ReserveStock(bizType, bizNo, warehouseID, items, expireAt)`
- `ConfirmReservation(bizType, bizNo)`
- `ReleaseReservation(bizType, bizNo, reason)`

说明：优先在服务层落地，不急于暴露额外管理端 HTTP 接口。

### 7.3 订单服务接口（兼容升级）

保持创建订单接口不变，内部库存处理切换为：

- 下单：调用 WMS 预占。
- 支付回调：确认预占并正式出库。
- 取消/超时：释放预占。

### 7.4 全端接口消费改造（app/eapp/admin）

1. `admin`：
   - 商品编辑页切换为“属性组 + 自动生成 + 差异预览”模式，保留手工覆盖能力。
   - 订单列表/详情增加“库存预占状态”展示（已预占、已确认、已释放）。
   - WMS 台账与流水增加 `reserved_qty` 和预占交易类型筛选。
2. `eapp`：
   - 商品编辑页 `SkuMatrixEditor` 改为以属性组为主，接收后端回传的 `sku_diff`。
   - 订单页新增取消订单动作，触发预占释放。
   - 商品卡片库存展示改为“可售库存”口径（`sellable`）。
3. `app`：
   - 购物车、确认单、下单链路透传库存不足与预占失败错误码。
   - 订单列表支持取消待支付订单并释放预占。
   - 商品详情库存展示使用后端可售库存口径。

## 8. 核心流程

### 8.1 商品保存与 SKU 差异化生成

1. 管理端提交商品 + 规格定义。
2. 后端规格引擎生成组合并计算 `sku_key`。
3. 与当前 `product_skus` 做差异比较。
4. 事务内执行：新增、更新、下线。
5. 返回 `sku_diff` 给前端展示。

### 8.2 下单预占

1. 订单服务按“单仓优先”选仓。
2. 锁定 `inventory_stock(warehouse_id, sku_id)`。
3. 校验 `qty - reserved_qty >= 购买数量`。
4. 增加 `reserved_qty`，写 `inventory_reservations` 与 `reserve` 流水。
5. 订单写入成功后提交事务。

### 8.3 支付确认

1. 按 `biz_type + biz_no` 查询预占记录。
2. 幂等校验：已确认直接成功返回。
3. 扣减 `qty` 与 `reserved_qty`。
4. 记录 `order_outbound` 流水并置预占状态为 `confirmed`。

### 8.4 取消/超时释放

1. 按订单号查询预占。
2. 幂等校验：已释放直接成功返回。
3. 扣减 `reserved_qty`，写 `release` 流水。
4. 预占状态置 `released`。

## 9. 单仓优先规则

1. 先取店铺默认仓（或全局默认仓）。
2. 仅在该仓校验可售库存。
3. 库存不足直接下单失败，不自动拆单。
4. 订单行持久化 `warehouse_id` 快照，保证履约一致。

## 10. 错误与并发控制

1. 统一业务错误码：库存不足、仓库不可用、SKU 已下线、预占过期。
2. 库存操作全部使用行锁（`FOR UPDATE`）。
3. 预占、确认、释放均需幂等约束。
4. 规格组合设置上限（如 >300 组合拒绝生成）避免笛卡尔积爆炸。

## 11. 测试设计

### 11.1 单元测试

1. 规格引擎：组合生成、键规范化、差异判定。
2. WMS 预占服务：预占成功、库存不足、并发冲突、幂等确认/释放。
3. 订单库存桥接：下单预占、支付确认、取消释放。

### 11.2 集成测试

1. 商品编辑后 SKU 新增/保留/下线结果正确。
2. 下单到支付全链路库存变更与流水一致。
3. 超时任务释放预占后可再次下单。

### 11.3 回归测试

1. 现有 WMS 入/出库单功能不回归。
2. 营销活动按 SKU 扣减与库存判断不回归。
3. 历史订单详情与 SKU 显示正常。

### 11.4 多端联调与构建验证

1. `admin`：商品编辑自动 SKU、订单取消、WMS 预占流水筛选可用。
2. `eapp`：商品编辑 `sku_diff` 展示、订单取消、库存展示口径一致。
3. `app`：购物车下单、库存不足提示、订单取消与状态刷新正确。
4. 构建验证：
   - `cd admin && npm run build`
   - `cd eapp && npm run test && npm run build:h5`
   - `cd app && npm run build:h5`

## 12. 部署与配置影响

1. 数据库迁移：
   - `product_skus` 加字段/索引。
   - `inventory_stock` 加 `reserved_qty`。
   - 新建 `inventory_reservations`。
2. 定时任务：新增“订单预占超时释放”任务。
3. 配置项：
   - 预占有效期（默认 15 分钟）。
   - 默认仓库选择策略。

## 13. 文档同步要求（docs-site）

落地实现时同步更新 `docs-site` 至最新架构，至少覆盖：

1. 功能说明：SKU 自动生成、SKU 下线、单仓优先。
2. 接口变化：商品保存接口新增字段、订单库存行为变化。
3. 部署影响：数据库迁移、超时释放任务、相关配置项。
4. 多端页面影响：`app/eapp/admin` 关键页面交互与字段语义更新。

## 14. 成功标准

1. 商品规格变更后 SKU ID 稳定，历史单据可追溯。
2. 任意订单库存变更均通过 WMS 预占与流水记录。
3. 系统无超卖，取消/超时可正确释放库存。
4. 管理端可直观看到 SKU 生成差异结果。
5. `app/eapp/admin` 三端库存展示与订单状态保持一致。
