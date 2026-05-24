# 订单补发物流与售后流程增强设计

**日期：** 2026-05-24
**范围：** `server` + `admin` + `web` + `app` + `docs-site`
**目标：** 在优先升级现有订单接口语义的前提下，实现完整物流轨迹（含补发）、标准售后流程（退货/换货）、用户侧申请入口与多端一致展示。

---

## 1. 已确认需求

1. 后台订单支持补发快递，并保留完整物流轨迹（非覆盖单号）。
2. `web` 与 `app` 订单详情页均需展示补发物流轨迹。
3. 系统业务逻辑需支持标准售后流程：退货与换货。
4. 退款能力采用“业务闭环 + 退款记录”模式：优先记录与追踪，不强依赖立即打款。
5. 换货采用“售后单内完成”模式，不创建独立换货订单。
6. 用户侧（`web` + `app`）提供申请售后入口，并可查看进度。

---

## 2. 设计原则

1. **优先升级现有接口**：对 `orders` 与 `ship` 语义扩展，避免无边界平行接口。
2. **新增接口有必要性证明**：仅在申请/审核/收货/退款等动作无法被现有资源表达时新增。
3. **轨迹可追溯**：物流与售后状态均保留事件日志，支持审计与回放。
4. **兼容现有前端**：保留 `orders.tracking_no` 字段语义，新端逐步迁移到 `shipments`。
5. **多端一致**：`admin`、`web`、`app` 共享同构核心字段，减少分叉逻辑。

---

## 3. 数据模型设计

### 3.1 新增 `order_shipments`（统一物流轨迹）

用途：统一承载首发、补发、用户退回、换货寄回等物流记录。

建议字段：

- `id`
- `order_id`（索引）
- `after_sale_case_id`（可空，索引）
- `direction`：`outbound|inbound`
- `biz_type`：`initial|reship|return`
- `company`
- `tracking_no`
- `logistics_status`：`pending|shipped|in_transit|signed|exception`
- `remark`
- `shipped_at`（可空）
- `signed_at`（可空）
- `created_by_type`：`admin|user|system`
- `created_by_id`
- `created_at` / `updated_at`

约束：

- 同一订单允许多条 `outbound`（支持补发）
- `tracking_no` 不做全局唯一，避免第三方测试号冲突

### 3.2 新增 `after_sale_cases`（售后主单）

用途：承载退货/换货主流程。

建议字段：

- `id`
- `order_id`（索引）
- `user_id`（索引）
- `merchant_id`
- `case_no`（唯一）
- `case_type`：`return|exchange`
- `status`（流程状态，见状态机）
- `reason`
- `apply_content`
- `apply_images_json`
- `audit_status`：`pending|approved|rejected`
- `audit_remark`
- `refund_amount`（decimal）
- `close_reason`
- `completed_at`（可空）
- `created_at` / `updated_at`

### 3.3 新增 `after_sale_case_items`（售后商品明细）

用途：支持订单部分商品售后。

建议字段：

- `id`
- `case_id`（索引）
- `order_item_id`（索引）
- `qty`
- `created_at` / `updated_at`

约束：

- 同一 `order_item_id` 在“进行中售后”内只能存在一条（避免并发重复售后）

### 3.4 新增 `after_sale_logs`（售后状态日志）

用途：全链路审计与前端进度时间线。

建议字段：

- `id`
- `case_id`（索引）
- `from_status`
- `to_status`
- `action`（apply/audit/return_ship/receive/refund/reship/complete/close）
- `operator_type`：`user|admin|system`
- `operator_id`
- `content`（说明）
- `ext_json`（扩展）
- `created_at`

### 3.5 复用与兼容

- 复用 `order_refunds` 作为退款台账，并新增 `after_sale_case_id` 关联字段。
- 保留 `orders.tracking_no`：继续同步最新一条 `outbound` 物流单号作为兼容展示。

---

## 4. 状态机设计

### 4.1 退货流程

`applied -> approved_wait_user_return -> user_returning -> warehouse_received -> refund_pending -> refunded -> completed`

分支：

- `applied -> rejected`
- 可关闭节点（如 `approved_wait_user_return`、`user_returning`、`refund_pending`）-> `closed`

### 4.2 换货流程

`applied -> approved_wait_user_return -> user_returning -> warehouse_received -> reship_pending -> reshipped -> completed`

分支：

- `applied -> rejected`
- 可关闭节点 -> `closed`

### 4.3 订单状态联动

- 存在进行中售后单：订单状态置为 `售后(5)`。
- 售后单全部完结/关闭后：订单恢复到终态（通常 `已完成(4)`）。

---

## 5. 接口设计

## 5.1 升级现有接口（优先）

### 5.1.1 升级 `PUT /admin/api/orders/:id/ship`

新增可选字段：

- `ship_type`：`initial|reship`（默认 `initial`）
- `after_sale_case_id`（`reship` 时可带）
- `company`
- `tracking_no`
- `remark`

兼容：旧请求仅传 `tracking_no` 时，行为保持首发发货。

执行效果：

- 写入 `order_shipments`
- 更新 `orders.tracking_no` 为最新出库单号
- 必要时推进售后状态（如换货补发后进入 `reshipped`）

### 5.1.2 升级订单查询接口

升级返回结构：

- 前台：`GET /api/v1/orders`、`GET /api/v1/orders/:id`
- 后台：`GET /admin/api/orders`、`GET /admin/api/orders/:id`

新增字段：

- `shipments`: 物流轨迹数组
- `after_sale_summary`: 售后摘要（进行中数量、最新节点、是否可申请）
- `latest_shipment`: 最近物流（列表快捷展示）

## 5.2 新增接口（必要且最小化）

新增的必要性：现有 `orders` 资源无法表达“申请/审核/收货/退款/完结”等动作语义与审计边界。

### 5.2.1 用户侧

- `POST /api/v1/orders/:id/after-sales`：创建售后申请（退货/换货）
- `GET /api/v1/after-sales/:id`：售后详情（状态、日志、物流、退款）
- `POST /api/v1/after-sales/:id/return-shipments`：用户填写回寄物流

### 5.2.2 管理侧

- `GET /admin/api/after-sales`：售后列表
- `GET /admin/api/after-sales/:id`：售后详情
- `POST /admin/api/after-sales/:id/audit`：审核通过/拒绝
- `POST /admin/api/after-sales/:id/receive`：仓库收货确认
- `POST /admin/api/after-sales/:id/refund`：退款登记（写 `order_refunds`）
- `POST /admin/api/after-sales/:id/complete`：人工完结
- `POST /admin/api/after-sales/:id/close`：关闭售后单

补发继续复用升级后的 `PUT /admin/api/orders/:id/ship`，不新增平行补发接口。

---

## 6. 前后端改造范围

### 6.1 `server`

- `plugins/order/model`：新增 4 张表模型 + `order_refunds` 字段扩展
- `plugins/order/service`：
  - 发货逻辑改为写 `order_shipments`
  - 新增售后申请/审核/收货/退款/完结领域服务
  - 订单查询聚合 `shipments + after_sale_summary`
- `plugins/order/api`：升级现有接口 + 新增最小售后动作接口

### 6.2 `admin`

- 订单列表/详情：展示物流轨迹与补发入口
- 新增售后管理页面：列表、详情、审核、收货、退款、完结、关闭
- 权限补充：`order:after-sale-view`、`order:after-sale-audit`、`order:after-sale-refund`

### 6.3 `web` + `app`

- 订单详情新增物流时间线（含补发）
- 新增“申请售后/查看售后进度”入口
- 售后申请页：类型、商品项、数量、原因、图片
- 售后详情页：状态节点、日志、回寄物流填写、退款状态

---

## 7. 关键业务规则

1. 同一 `order_item` 同时只允许一个进行中售后单。
2. 仅允许已支付订单申请售后。
3. 换货必须先“仓库收货确认”再允许补发。
4. 退款金额不得超过订单实付及可退余额。
5. 所有状态流转必须写 `after_sale_logs`。
6. 非法状态跳转（如未审核先收货）直接拒绝。

---

## 8. 错误处理与幂等

- 幂等键建议：申请、审核、退款动作支持前端传 `request_id`。
- 重复提交：服务端检测当前状态与最新动作，重复请求返回成功态或明确错误。
- 并发防护：关键流转在事务中对售后单行加锁（`FOR UPDATE`）。

---

## 9. 测试与验收标准

### 9.1 后端

- 发货接口：首发、补发、兼容旧参数三种路径通过。
- 售后流程：退货/换货主流程 + 拒绝/关闭分支通过。
- 退款记录：`order_refunds` 与售后单关联正确。
- 查询聚合：订单详情可返回完整 `shipments`。

### 9.2 前端

- `admin` 能看到并操作售后闭环。
- `web` 与 `app` 订单详情均展示补发轨迹。
- 用户可发起售后并查看进度。

### 9.3 回归

- 原支付/下单/评价流程不回归。
- 旧端仅读 `tracking_no` 仍可展示最近出库单号。

---

## 10. 文档更新要求（docs-site）

本功能实施后同步更新：

1. `docs-site/docs/api/order.md`
   - 补充升级接口字段
   - 新增售后接口章节
   - 说明兼容策略与状态流
2. `docs-site/docs/guide/features.md`
   - 更新订单与售后能力说明
   - 明确支持补发物流轨迹、退货/换货闭环

---

## 11. 非目标（本阶段不做）

1. 不引入外部物流轨迹抓取服务（仅存业务录入轨迹）。
2. 不实现自动退款打款网关编排（先完成退款记录与状态闭环）。
3. 不实现多轮协商/仲裁流程。

---

## 12. 风险与缓解

1. **风险：** 状态机分支多，易出现非法跳转。
   - **缓解：** 统一状态转移表 + 单点校验函数。
2. **风险：** 旧页面仍依赖 `tracking_no`。
   - **缓解：** 保留并同步 `tracking_no` 到最新出库单。
3. **风险：** 前后端字段升级不同步。
   - **缓解：** 先落后端兼容，再分端灰度展示。
