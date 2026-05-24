# 评价系统（根评价+追评+后台回复）设计文档

## 1. 背景与目标

当前订单评价能力过于简化，仅通过弹窗提交一段文本，并写入订单 `remark`。该实现无法满足以下需求：

- 评价应为独立页面，支持商品评分、物流评分和文字输入
- 商品详情页需要“评价 Tab”，展示评价与评分统计
- 评价支持“修改原评价”
- 评价支持“追加评论”（仅根评价可追加，追加为根评价子级）
- 后台支持管理评价，并支持管理员单条回复

本次目标是在不破坏现有订单资源语义的前提下，建立结构化评价模型与多端展示闭环。

## 2. 约束与原则

- 优先升级现有接口语义，不重复造平行接口
- 评价存储使用独立表，不继续挤压 `order.remark`
- 根评价唯一：每个订单商品仅允许一条根评价
- 修改行为覆盖根评价，不生成新根记录
- 追加行为仅允许挂在根评价下（不允许多级追加）
- 管理员回复每根评价仅一条，可覆盖更新

## 3. 数据模型设计

### 3.1 `order_reviews`（根评价）

- `id`
- `order_id`
- `order_item_id`（唯一约束，保证根评价唯一）
- `product_id`
- `user_id`
- `merchant_id`
- `product_score`（1-5）
- `logistics_score`（1-5）
- `content`
- `edited_times`
- `created_at`
- `updated_at`

语义：一条记录代表某个订单商品的“主评价”。

### 3.2 `order_review_appends`（追评）

- `id`
- `review_id`（关联根评价）
- `user_id`
- `content`
- `created_at`
- `updated_at`

语义：追评是根评价的子级事件流，仅根评价可拥有追评。

### 3.3 `order_review_replies`（管理员回复）

- `id`
- `review_id`（唯一约束，确保每根评价最多一条管理员回复）
- `admin_id`
- `content`
- `created_at`
- `updated_at`

语义：后台可对每条根评价给出一条官方回复，后续可编辑覆盖。

## 4. 接口设计

### 4.1 升级现有前台评价提交接口

接口：`POST /api/v1/orders/:id/review`

请求体：

```json
{
  "mode": "create|edit|append",
  "logistics_score": 5,
  "items": [
    { "order_item_id": 11, "product_score": 5, "content": "做工很好" }
  ],
  "append_content": "用了两周，续航依旧稳定"
}
```

行为：

- `create`：创建根评价（仅允许未评价订单商品）
- `edit`：更新根评价评分与内容，`edited_times + 1`
- `append`：新增追评子记录（仅根评价存在时允许）

说明：

- 继续复用订单资源下 `review` 动作，不新增平行“提交评价”接口
- `order_id` 由 path 传入并参与校验，防止越权

### 4.2 新增商品评价查询接口（前台）

接口：`GET /api/v1/products/:id/reviews?page=1&size=10`

返回：

```json
{
  "summary": {
    "avg_product_score": 4.8,
    "avg_logistics_score": 4.7,
    "total": 132
  },
  "list": [
    {
      "id": 1001,
      "product_score": 5,
      "logistics_score": 5,
      "content": "质感不错",
      "edited_times": 1,
      "user": { "id": 1, "nickname": "张三", "avatar": "..." },
      "created_at": "2026-05-24T10:00:00Z",
      "updated_at": "2026-05-24T11:00:00Z",
      "appends": [
        { "id": 9001, "content": "追加：用了1周依然满意", "created_at": "2026-05-26T10:00:00Z" }
      ],
      "admin_reply": {
        "id": 8001,
        "content": "感谢支持",
        "created_at": "2026-05-26T12:00:00Z"
      }
    }
  ],
  "total": 132,
  "page": 1,
  "size": 10
}
```

### 4.3 后台评价管理接口

- `GET /admin/api/reviews?product_id=&keyword=&page=&size=`
- `GET /admin/api/reviews/:id`
- `POST /admin/api/reviews/:id/reply`（创建或覆盖单条管理员回复）

## 5. 前端页面设计

### 5.1 H5 评价页

新增页面：`/pages/order/review`

- 顶部显示订单商品列表（可评价商品）
- 每个商品模块包含：
  - `up-rate` 商品评分
  - `textarea` 评价内容
- 订单维度显示物流评分（`up-rate`）
- 提交按钮：
  - 无根评价 -> `mode=create`
  - 有根评价 -> `mode=edit`
- 追加入口：
  - 已有根评价时显示“追加评论”
  - 追加仅提交 `mode=append + append_content`

### 5.2 PC 评价页

新增路由页：`/orders/:id/review`

- 与 H5 同语义：商品评分、物流评分、文本、修改、追加
- 组件风格延续现有 web 页面

### 5.3 商品详情页“评价 Tab”

H5 与 PC 均升级为 Tab 结构：

- `详情`：现有详情 blocks 渲染
- `评价`：请求 `GET /api/v1/products/:id/reviews`
  - 评分摘要卡片
  - 根评价列表
  - 追评时间线
  - 管理员回复块

## 6. 服务层关键流程

### 6.1 创建评价（`mode=create`）

1. 校验订单归属与状态（至少已收货或已完成可评价）
2. 校验 `order_item_id` 属于该订单
3. 校验根评价不存在
4. 写入 `order_reviews`
5. 若本次评价后该订单所有商品均已评价，可将订单状态更新为已完成

### 6.2 修改评价（`mode=edit`）

1. 按 `order_item_id + user_id` 找到根评价
2. 更新评分与内容，`edited_times + 1`

### 6.3 追加评价（`mode=append`）

1. 查找根评价（用户本人）
2. 写入 `order_review_appends`

### 6.4 后台回复

1. 校验评价存在
2. 若已有回复则更新内容；无则创建

## 7. Mock 与真实接口一致性

- app/web mock 同步支持 `mode=create|edit|append`
- 商品详情 mock 加入 `reviews.summary` 与 `reviews.list`
- 后台 mock（如使用）同步支持列表与回复动作

## 8. 风险与对策

- 风险：旧数据无评价结构，商品详情评价 Tab 为空
  - 对策：前端为空态设计 + mock 提供演示数据
- 风险：`order_item_id` 唯一性与历史脏数据冲突
  - 对策：迁移前排查并清洗重复记录（若存在）
- 风险：多端并发提交 create 导致重复
  - 对策：数据库唯一键 + 事务兜底

## 9. 验证范围

- 后端：
  - 评价创建、修改、追加、后台回复
  - 商品评价查询分页与统计
- 前端：
  - H5/PC 评价页提交流程
  - 商品详情评价 Tab 渲染
  - 订单列表“评价”跳转与状态反馈

## 10. 文档同步要求

需同步更新：

- `docs-site/docs/api/order.md`（review 接口升级）
- `docs-site/docs/api/product.md`（新增产品评价查询接口）
- `docs-site/docs/guide/features.md`（评价系统能力说明）
