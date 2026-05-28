# 商家移动端 eapp

`eapp` 是 LYShop 的独立商家移动端工程，面向单店铺后台运营场景，支持 `H5 + 微信小程序 + App`。

在线演示：<a href="https://ijry.github.io/lyshop/eapp-demo/index.html#/pages/dashboard/index" target="_blank" rel="noreferrer">eapp 商家端 H5 演示（Mock 免登录）</a>

## 功能说明

- 五个主导航：工作台、订单、商品、营销、我的
- 工作台：经营概览、待办事项、快捷入口
- 订单：列表、详情、发货弹窗（首次/补发、快递/同城）、物流轨迹同步与节点展示
- 售后：列表、详情、审核通过/拒绝、确认收货、登记退款、完结、关闭
- 商品：列表、上下架、商品编辑（高频字段）
- 营销：优惠券、秒杀/拼团/砍价活动管理、活动商品管理
- 我的：消息中心、IM 会话、店铺设置、管理员、角色权限

## 接口变化

- `eapp` 复用现有后台接口：`/admin/api/*`
- 无新增商家专属接口前缀与独立后端服务
- 鉴权沿用后台登录与 token 语义，移动端使用独立本地存储键隔离会话
- 商品接口按后台语义提交：
  - 新增商品：`POST /products`，body 为 `{ product, skus?, images? }`
  - 编辑商品：`PUT /products/:id`，body 为 `{ product, skus?, images? }`
- 订单发货接口使用统一发货参数：
  - `PUT /orders/:id/ship` 支持 `delivery_type`、`ship_type`、`after_sale_case_id`、物流/骑手字段
  - `POST /orders/:id/shipments/:shipment_id/sync` + `GET /orders/:id/shipments/:shipment_id/tracks` 用于轨迹同步与查询
- 售后动作复用后台接口：
  - 审核 `/after-sales/:id/audit`
  - 收货 `/after-sales/:id/receive`
  - 退款 `/after-sales/:id/refund`
  - 完结 `/after-sales/:id/complete`
  - 关闭 `/after-sales/:id/close`

## 部署与配置影响

- 新增 `eapp/` 前端工程构建与发布
- 后端无新增环境变量与配置项
- 网关仅需新增 `eapp` 静态资源托管路径（H5 部署时）
- 小程序与 App 端仅新增前端构建产物，不涉及后端部署拓扑调整

## 联调验收清单

- 登录态
  - 使用后台管理员账号登录后可进入工作台
  - token 过期后自动回到登录页并清理本地会话
- 订单与物流
  - 订单列表按状态筛选可用
  - 订单详情可执行首次发货（快递/同城）
  - 补发模式可提交 `after_sale_case_id`
  - 快递单可执行轨迹同步并展示节点
- 售后流程
  - 审核通过/拒绝仅在 `applied` 状态可用
  - 确认收货仅在 `user_returning` 状态可用
  - 登记退款仅在 `refund_pending` 状态可用
  - 完结仅在 `refunded`、`reshipped` 状态可用
  - 关闭仅在可关闭状态集合可用
- 商品与营销
  - 商品上下架调用 `PUT /products/:id` 并传 `{ product: { status } }`
  - 商品新增/编辑调用包装体 `{ product, skus?, images? }`
  - 秒杀/拼团/砍价页面可创建活动并维护活动商品
- 多端构建
  - `npm run build:h5`
  - `npm run build:mp-weixin`
  - `npm run build:app`
