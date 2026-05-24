# 后台概览接口

## 功能说明

后台首页（`/dashboard`）改为读取真实接口数据，不再使用静态占位文案。  
接口返回四个核心指标与近 7 天趋势，前端据此渲染概览卡片和 ECharts 双轴组合图（销售额柱状 + 订单量折线）。

## 接口变化

### `GET /admin/api/dashboard`

返回示例：

```json
{
  "today_orders": 56,
  "today_sales": 28960.5,
  "pending_refunds": 3,
  "online_sessions": 2,
  "sales_trend": [
    { "date": "2026-05-19", "orders": 42, "sales": 18660 },
    { "date": "2026-05-20", "orders": 38, "sales": 17280 }
  ]
}
```

字段说明：

- `today_orders`：今日创建订单数
- `today_sales`：今日销售额（已支付相关订单）
- `pending_refunds`：待处理售后数量（未关闭/未拒绝/未完结）
- `online_sessions`：在线客服会话数（未关闭）
- `sales_trend`：近 7 天趋势数据（日期、订单数、销售额）

## 部署与配置影响

- 无新增环境变量。
- 无新增外部服务依赖。
- 管理端新增 `echarts` 前端依赖，用于后台首页趋势图渲染。
- 依赖订单、售后、IM 表数据；若对应插件未启用，接口会返回可用默认值（0 和空趋势）。
