# 仓储管理重构设计（库存交易引擎版）

## 1. 背景与目标

当前后台仓储管理仅覆盖基础库存列表与手工入/出库，缺少产品级常用能力：

- 仓库管理能力不完整（缺列表管理、编辑、启停）
- 出入库单缺失（无法按单据管理多 SKU 明细与状态）
- 库存流水查询缺失（审计追溯能力弱）

本次采用 **方案 B：统一库存交易引擎**，按最新架构重建，不保留旧兼容层。  
目标是先交付“最小可用产品版”闭环：**仓库管理 + 出入库单 + 库存台账 + 库存流水**。

---

## 2. 设计原则

1. **流水为事实源**：所有库存变化必须先写流水，再更新库存快照。
2. **单据驱动库存**：可用库存不允许人工直接修改，仅通过业务单据生效。
3. **状态可审计**：单据仅允许 `draft -> completed/canceled`，完成后只读。
4. **事务一致性**：单据完成时，库存快照和流水在同一事务内提交。
5. **接口清晰边界**：仓库、库存、单据、流水四组接口职责明确。

---

## 3. 方案对比与决策

### 方案 A（兼容升级）
- 在既有 `/wms/inbound`、`/wms/outbound` 基础上扩展单据。
- 优点：改动较小；缺点：长期维护成本高，模型不统一。

### 方案 B（本次采用）
- 统一为“库存交易引擎 + 业务单据”的最新模型。
- 优点：语义清晰、审计能力完整、后续可扩展盘点/调拨。
- 成本：一次性改造更大，但总成本更低。

### 方案 C（前端先行、后端补丁）
- 先堆页面再补后端。
- 风险：返工高、数据一致性不可控。

**决策：采用方案 B。**

---

## 4. 数据模型设计

### 4.1 仓库主数据 `wms_warehouses`

- `id`
- `name`
- `code`（仓库编码）
- `contact`
- `phone`
- `address`
- `status`（`1=启用`，`0=停用`）
- `created_at/updated_at/deleted_at`

约束：
- `code` 唯一
- 停用仓库不可新建单据

### 4.2 库存快照 `wms_inventory_stocks`

- `id`
- `warehouse_id`
- `sku_id`
- `available_qty`（可用库存）
- `safety_qty`（安全库存）
- `updated_at`

约束：
- 唯一键：`(warehouse_id, sku_id)`

### 4.3 库存流水 `wms_inventory_movements`

- `id`
- `warehouse_id`
- `sku_id`
- `direction`（`in/out`）
- `qty`（正数）
- `before_qty`
- `after_qty`
- `biz_type`（`inbound_doc/outbound_doc/adjust`）
- `biz_doc_id`
- `biz_doc_no`
- `remark`
- `operator_id`
- `created_at`

索引建议：
- `warehouse_id + sku_id + created_at`
- `biz_type + biz_doc_id`
- `biz_doc_no`

### 4.4 单据头 `wms_inventory_docs`

- `id`
- `doc_no`
- `doc_type`（`inbound/outbound`）
- `warehouse_id`
- `status`（`draft/completed/canceled`）
- `remark`
- `completed_at`
- `operator_id`
- `created_at/updated_at`

### 4.5 单据行 `wms_inventory_doc_items`

- `id`
- `doc_id`
- `sku_id`
- `qty`
- `unit_cost`（可空，先保留字段）
- `remark`

约束：
- 同一单据禁止重复 SKU（通过业务校验）

---

## 5. 状态机设计

入库单和出库单共用同一状态模型：

- `draft -> completed`
- `draft -> canceled`

约束：
- 仅 `draft` 可编辑、完成、作废
- `completed/canceled` 后不可修改

---

## 6. 接口设计

### 6.1 仓库管理

- `GET /admin/api/wms/warehouses`
  - 能力：分页、按名称/状态筛选
- `POST /admin/api/wms/warehouses`
  - 能力：新建仓库
- `PUT /admin/api/wms/warehouses/:id`
  - 能力：编辑仓库资料
- `PUT /admin/api/wms/warehouses/:id/status`
  - 能力：启停用

### 6.2 库存台账

- `GET /admin/api/wms/stocks`
  - 能力：按仓库/SKU/预警状态筛选
  - 返回：`available_qty/safety_qty/is_warning`
- `PUT /admin/api/wms/stocks/:id/safety`
  - 能力：仅更新安全库存

### 6.3 出入库单

- `GET /admin/api/wms/docs`
  - 能力：按 `doc_type/status/doc_no/date_range` 筛选
- `POST /admin/api/wms/docs`
  - 能力：创建草稿单，支持多 SKU 明细
- `GET /admin/api/wms/docs/:id`
  - 能力：单据详情
- `PUT /admin/api/wms/docs/:id`
  - 能力：仅草稿可编辑头+明细
- `POST /admin/api/wms/docs/:id/complete`
  - 能力：完成单据并生效库存
- `POST /admin/api/wms/docs/:id/cancel`
  - 能力：作废草稿单

### 6.4 库存流水

- `GET /admin/api/wms/movements`
  - 能力：按仓库/SKU/业务类型/单号/时间区间筛选
  - 返回：方向、变更量、变更前后、业务单号、操作人、时间

---

## 7. 核心事务流程（完成单据）

`complete` 的执行步骤（事务内）：

1. 锁定单据头（防并发重复完成）。
2. 校验单据状态必须为 `draft`。
3. 校验仓库状态为启用。
4. 遍历明细行并锁定对应库存快照（`warehouse_id + sku_id`）。
5. 计算变更后库存：
   - 入库：`after = before + qty`
   - 出库：先校验 `before >= qty`，再 `after = before - qty`
6. 写入库存流水（记录 `before/after/biz_doc_no/operator_id`）。
7. 更新库存快照 `available_qty`。
8. 更新单据状态为 `completed`，写入 `completed_at`。
9. 提交事务。

失败回滚：
- 任一 SKU 出库不足、状态非法、并发冲突，整单回滚。

---

## 8. 后台页面信息架构

菜单：

1. `仓储管理 / 仓库管理`
   - 列表、新建、编辑、启停用
2. `仓储管理 / 出入库单`
   - 入库单/出库单 Tab
   - 草稿编辑（多 SKU）
   - 完成/作废
3. `仓储管理 / 库存台账`
   - 按仓库/SKU/预警筛选
   - 仅允许维护安全库存
4. `仓储管理 / 库存流水`
   - 多条件筛选
   - 支持跳转业务单据详情

交互关键点：
- 完成单据前需二次确认
- 已完成/已作废进入只读态
- 错误提示直出业务原因（库存不足、状态非法、仓库停用）

---

## 9. 错误处理与并发策略

1. 状态非法：返回 `409`（仅草稿可流转）。
2. 库存不足：返回业务错误，包含 `sku_id`、当前库存、请求数量。
3. 并发冲突：单据完成使用行锁，重复完成直接失败。
4. 参数错误：空明细、数量小于等于 0、重复 SKU 一律拒绝。

---

## 10. 测试与验收

### 10.1 服务层单测

- 入库完成：库存增加且写流水
- 出库完成：库存减少且写流水
- 出库不足：整单回滚
- 重复完成：失败

### 10.2 API 集成测试

- 出入库单创建/编辑/完成/作废全链路
- 仓库启停后新建单据限制
- 流水查询筛选正确性（仓库、SKU、单号、时间）

### 10.3 前端冒烟测试

- 仓库管理 CRUD + 启停
- 出入库单多 SKU 草稿编辑与完成
- 库存台账预警筛选与安全库存调整
- 库存流水查询与单据跳转

### 10.4 验收标准

- 任意库存变化都可追溯到业务单据与操作人
- 可用库存仅能通过单据变更
- 出入库单支持多 SKU 且状态流转正确
- 后台形成仓储操作闭环

---

## 11. 对 docs-site 的同步要求

本次属于系统功能变更，实施时必须同步更新 `docs-site`，至少覆盖：

- 功能说明：仓库管理、出入库单、库存台账、库存流水
- 接口变化：`/admin/api/wms/*` 最新接口分组与字段
- 部署/配置影响：如无新增依赖需明确标注“无新增配置项”
