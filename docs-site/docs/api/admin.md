# 后台概览接口

## 功能说明

后台首页（`/dashboard`）改为读取真实接口数据，不再使用静态占位文案。  
接口返回四个核心指标与近 7 天趋势，前端据此渲染概览卡片和 ECharts 双轴组合图（销售额柱状 + 订单量折线）。

后台管理端交互提示已统一接入 `notify` 通道，并默认通过全局 Toast 组件展示，便于后续替换为自定义消息系统。

## 接口变化

- 本次仅涉及后台前端提示通道改造，无新增/变更后台 API。

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
- `notify` 默认使用前端全局 Toast 渲染，无需新增服务端配置。
- 依赖订单、售后、IM 表数据；若对应插件未启用，接口会返回可用默认值（0 和空趋势）。

## 商家移动端接口复用说明

- 商家移动端 `eapp` 与管理后台统一复用 `/admin/api/*` 资源。
- 本次无新增 `merchant` 前缀接口，无需额外网关转发规则。

## 商家移动端工作台与基础接口

### GET /dashboard（升级）

在原有 today_* 与待办字段基础上，追加返回：

- `today_avg_price`
- `compare`：`{ revenue_yoy, revenue_mom, order_yoy, order_mom }` （0.18 表示 +18%）
- `trend`：`{ revenue_7d, revenue_30d, order_7d }`，每项为 `{ categories: string[], series: [{ name, data }] }`
- `status_distribution`：`[{ name, value }]`
- `hot_products`：`[{ id, title, cover, sold_qty }]`
- `announcements`：`[{ id, title, content, type, created_at }]`
- `stock_warning_list`：`[{ product_id, sku_id, title, stock, threshold }]`

### GET /shops/current

当前店铺：`{ id, name, logo, owner, decor_status }`

### GET /announcements

平台公告列表：标准 `{ list, total, page, size }` 分页结构。

### POST /after-sales/{id}/messages

售后协商消息：`{ from: 'merchant'|'user', content, images? }` → `{ id, messages: [...] }`

### POST /after-sales/{id}/evidences

商家凭证：`{ images: string[], remark? }` → `{ id, evidences: [...] }`
- eapp 侧仅增加客户端请求封装与会话存储键，不改变服务端接口语义。
