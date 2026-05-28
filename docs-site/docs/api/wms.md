# 仓储接口

## 说明

仓储模块围绕“仓库管理、出入库单、库存台账、库存流水”提供后台能力。库存变更以单据完成动作为准，所有变更可通过流水追溯到业务单据。
后台管理菜单默认入口为 `/wms/warehouse`、`/wms/docs`、`/wms/stock`、`/wms/movements`。

## 功能说明

- 仓库管理：支持仓库分页查询、新建、编辑和启停用。
- 出入库单：统一单据模型，支持入库/出库、多 SKU 明细、草稿保存。
- 单据流转：仅草稿可编辑，支持 `draft -> completed` 与 `draft -> canceled`。
- 库存台账：按仓库、SKU、预警状态查询库存数量与安全库存。
- 安全库存：支持按库存记录维护 `safe_qty`，用于低库存预警识别。
- 库存流水追溯：按仓库、SKU、业务类型、单号、时间区间查询变更明细。

## 接口变化

当前后台仓储接口统一归类在 `/admin/api/wms/*`，沿用既有 `/admin/api` 前缀与权限体系。

### 仓库管理（`warehouses*`）

- `GET /admin/api/wms/warehouses?keyword=&status=&page=&size=`
- `POST /admin/api/wms/warehouses`
- `PUT /admin/api/wms/warehouses/:id`
- `PUT /admin/api/wms/warehouses/:id/status`

`status` 取值：`1` 启用，`0` 停用。

仓库新建/编辑请求示例：

```json
{
  "code": "WH-SH-001",
  "name": "上海一仓",
  "address": "上海市浦东新区示例路 1 号",
  "contact": "仓管员",
  "phone": "13800000000",
  "status": 1
}
```

启停用请求示例：

```json
{
  "status": 0
}
```

### 库存台账与安全库存

- `GET /admin/api/wms/stocks?warehouse_id=&sku_id=&warning_only=&page=&size=`
- `PUT /admin/api/wms/stocks/:id/safety`

台账列表核心字段：`qty`、`safe_qty`、`is_warning`。

安全库存请求示例：

```json
{
  "safe_qty": 20
}
```

### 出入库单

- `GET /admin/api/wms/docs?warehouse_id=&doc_type=&status=&doc_no=&start_at=&end_at=&page=&size=`
- `POST /admin/api/wms/docs`
- `GET /admin/api/wms/docs/:id`
- `PUT /admin/api/wms/docs/:id`
- `POST /admin/api/wms/docs/:id/complete`
- `POST /admin/api/wms/docs/:id/cancel`

`doc_type`：`inbound | outbound`  
`status`：`draft | completed | canceled`

出入库单请求示例（支持多 SKU）：

```json
{
  "warehouse_id": 1,
  "doc_type": "outbound",
  "remark": "门店补货",
  "items": [
    { "sku_id": 1001, "qty": 3, "remark": "常规出库" },
    { "sku_id": 1002, "qty": 1, "remark": "加急" }
  ]
}
```

### 库存流水

- `GET /admin/api/wms/movements?warehouse_id=&sku_id=&biz_type=&doc_no=&start_at=&end_at=&page=&size=`

流水核心字段：`biz_type`、`change_qty`、`before_qty`、`after_qty`、`doc_no`、`occurred_at`。

## 部署与配置影响

- 无新增第三方依赖。
- 无新增环境变量与配置项。
- `wms` 插件在服务启动时自动迁移仓储相关表结构（仓库、库存、单据、单据明细、流水）。
- 权限沿用后台既有鉴权，读写分别使用 `wms:view`、`wms:edit`。
