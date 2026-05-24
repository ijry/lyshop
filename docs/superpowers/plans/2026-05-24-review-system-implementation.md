# 评价系统（根评价+追评+后台回复） Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现“独立评价表 + 评价页 + 商品详情评价Tab + 后台单条回复”全链路，并保持现有订单接口语义升级。

**Architecture:** 采用三表模型（根评价、追评、后台回复），升级 `POST /api/v1/orders/:id/review` 支持 `create/edit/append` 模式；商品详情评价通过产品接口新增查询能力聚合返回。前端新增 H5/PC 评价页面，商品详情页改为“详情/评价”双 Tab。

**Tech Stack:** Go + Gin + GORM, Vue3 + Vite (web), uni-app + uview-plus (app), 本地 mock 路由

---

## File Structure (Design Lock)

- `server/plugins/order/model/review.go`（新建）：评价三表模型
- `server/plugins/order/plugin.go`：挂载评价模型迁移
- `server/plugins/order/service/review.go`（新建）：评价创建/编辑/追加/后台回复/查询逻辑
- `server/plugins/order/service/order.go`：迁移旧 `ReviewOrder` 调用到新服务
- `server/plugins/order/api/front.go`：升级评价请求 DTO、绑定 `mode/items/logistics_score/append_content`
- `server/plugins/order/api/admin.go`：新增后台评价管理与回复 API
- `server/plugins/product/service/review_bridge.go`（新建或并入 product.go）：产品维度评价查询聚合
- `server/plugins/product/api/front.go`：新增产品评价查询路由
- `app/pages/order/review.vue`（新建）：H5 评价页
- `app/pages/order/list.vue`：评价按钮跳转新页面
- `app/pages/product/detail.vue`：详情页新增“详情/评价”Tab
- `app/pages.json`：注册 H5 评价页路由
- `app/mock/index.ts`、`app/mock/data/product-detail.json`、`app/mock/data/orders.json`：mock 升级
- `web/src/views/OrderReview.vue`（新建）：PC 评价页
- `web/src/views/OrderList.vue`：评价按钮跳转
- `web/src/views/ProductDetail.vue`：详情页新增“详情/评价”Tab
- `web/src/router/index.ts`：新增 PC 评价页路由
- `web/src/mock/index.ts`：mock 升级
- `admin/src/views/review/ReviewList.vue`（新建）：后台评价管理页
- `admin/src/router/index.ts`、`admin/src/layouts/AdminLayout.vue`、`admin/src/api/plugins.ts`：后台菜单/路由/API
- `docs-site/docs/api/order.md`、`docs-site/docs/api/product.md`、`docs-site/docs/guide/features.md`：文档同步

---

### Task 1: 后端评价模型与迁移接入

**Files:**
- Create: `server/plugins/order/model/review.go`
- Modify: `server/plugins/order/plugin.go`

- [ ] **Step 1: 新建评价模型文件，定义三张表结构与 GORM 约束。**
- [ ] **Step 2: 在 `order` 插件迁移里追加评价模型 AutoMigrate。**
- [ ] **Step 3: 编译校验。**
Run: `go test ./plugins/order/...`
Expected: PASS

- [ ] **Step 4: Commit**
```bash
git add server/plugins/order/model/review.go server/plugins/order/plugin.go
git commit -m "新增订单评价三表模型并接入迁移" -m "- 新增根评价、追评、管理员回复三张模型表。\n- order插件迁移增加评价模型AutoMigrate。"
```

### Task 2: 后端前台评价写入能力（create/edit/append）

**Files:**
- Create: `server/plugins/order/service/review.go`
- Modify: `server/plugins/order/service/order.go`
- Modify: `server/plugins/order/api/front.go`

- [ ] **Step 1: 在 `review.go` 实现 `SubmitOrderReview`，支持 `mode=create|edit|append`。**
- [ ] **Step 2: 在 `front.go` 升级 `reviewOrder` 请求体 DTO，透传到服务层。**
- [ ] **Step 3: 移除旧 remark 评价拼接逻辑，改为调用新服务。**
- [ ] **Step 4: 后端编译验证。**
Run: `go test ./plugins/order/...`
Expected: PASS

- [ ] **Step 5: Commit**
```bash
git add server/plugins/order/service/review.go server/plugins/order/service/order.go server/plugins/order/api/front.go
git commit -m "升级前台订单评价为创建编辑追加模式" -m "- 复用/orders/:id/review接口支持create/edit/append。\n- 评分与文本写入独立评价表，移除remark拼接实现。"
```

### Task 3: 后端产品评价查询接口

**Files:**
- Modify: `server/plugins/product/api/front.go`
- Modify or Create: `server/plugins/product/service/product.go` / `server/plugins/product/service/review_bridge.go`

- [ ] **Step 1: 新增 `GET /api/v1/products/:id/reviews` 路由。**
- [ ] **Step 2: 实现产品评价分页查询，返回 summary + list(appends/admin_reply)。**
- [ ] **Step 3: 后端全量编译验证。**
Run: `go test ./...`
Expected: PASS

- [ ] **Step 4: Commit**
```bash
git add server/plugins/product/api/front.go server/plugins/product/service/product.go server/plugins/product/service/review_bridge.go
git commit -m "新增商品评价查询接口与评分统计" -m "- 新增products/:id/reviews接口。\n- 返回评价分页、追评子级、管理员回复与评分摘要。"
```

### Task 4: 后台评价管理与单条回复接口

**Files:**
- Modify: `server/plugins/order/api/admin.go`
- Modify: `server/plugins/order/service/review.go`

- [ ] **Step 1: 新增后台评价列表、详情、回复（创建或覆盖）API。**
- [ ] **Step 2: 权限命名沿用 `order:view`，回复使用 `order:ship` 同级新权限 `order:review-reply`（若权限体系可扩展）。**
- [ ] **Step 3: 后端验证。**
Run: `go test ./plugins/order/...`
Expected: PASS

- [ ] **Step 4: Commit**
```bash
git add server/plugins/order/api/admin.go server/plugins/order/service/review.go
git commit -m "补齐后台评价管理与单条回复接口" -m "- 新增后台评价列表和详情查询。\n- 支持每条根评价单条管理员回复并可覆盖更新。"
```

### Task 5: H5 评价页与商品详情评价Tab

**Files:**
- Create: `app/pages/order/review.vue`
- Modify: `app/pages/order/list.vue`
- Modify: `app/pages/product/detail.vue`
- Modify: `app/pages.json`

- [ ] **Step 1: 新增 H5 评价页，支持商品评分(`up-rate`)、物流评分(`up-rate`)、文本输入、create/edit/append 切换。**
- [ ] **Step 2: 订单列表评价按钮改为跳转评价页。**
- [ ] **Step 3: 商品详情页改双 Tab（详情/评价），评价 Tab 拉取并渲染评价列表。**
- [ ] **Step 4: H5 构建验证。**
Run: `npm run build:h5:demo`
Expected: PASS

- [ ] **Step 5: Commit**
```bash
git add app/pages/order/review.vue app/pages/order/list.vue app/pages/product/detail.vue app/pages.json
git commit -m "新增H5独立评价页并接入商品详情评价Tab" -m "- 评价页支持评分、修改与追加评论。\n- 商品详情页新增评价Tab展示评分摘要与评价时间线。"
```

### Task 6: H5 mock 升级

**Files:**
- Modify: `app/mock/index.ts`
- Modify: `app/mock/data/product-detail.json`
- Modify: `app/mock/data/orders.json`

- [ ] **Step 1: mock 支持 `review` 的 `mode=create|edit|append`。**
- [ ] **Step 2: 产品详情 mock 增加评价摘要与列表结构。**
- [ ] **Step 3: H5 演示构建验证。**
Run: `npm run build:h5:demo`
Expected: PASS

- [ ] **Step 4: Commit**
```bash
git add app/mock/index.ts app/mock/data/product-detail.json app/mock/data/orders.json
git commit -m "升级H5 mock评价结构与评价动作" -m "- review接口支持创建编辑追加模式。\n- 商品详情mock补齐评价Tab所需数据。"
```

### Task 7: PC 评价页与商品详情评价Tab

**Files:**
- Create: `web/src/views/OrderReview.vue`
- Modify: `web/src/views/OrderList.vue`
- Modify: `web/src/views/ProductDetail.vue`
- Modify: `web/src/router/index.ts`
- Modify: `web/src/mock/index.ts`

- [ ] **Step 1: 新增 PC 评价页，支持与 H5 同语义的评分、修改、追加。**
- [ ] **Step 2: 订单列表评价按钮改跳转评价页。**
- [ ] **Step 3: 商品详情页改为详情/评价 Tab 并渲染评价列表。**
- [ ] **Step 4: mock 同步升级支持评价结构。**
- [ ] **Step 5: web 构建验证。**
Run: `npm run build`
Expected: PASS

- [ ] **Step 6: Commit**
```bash
git add web/src/views/OrderReview.vue web/src/views/OrderList.vue web/src/views/ProductDetail.vue web/src/router/index.ts web/src/mock/index.ts
git commit -m "新增PC评价页并完善商品详情评价Tab" -m "- 评价功能从弹窗改为独立页面。\n- 商品详情新增评价Tab并展示追评与后台回复。"
```

### Task 8: 后台评价管理页面

**Files:**
- Create: `admin/src/views/review/ReviewList.vue`
- Modify: `admin/src/router/index.ts`
- Modify: `admin/src/layouts/AdminLayout.vue`
- Modify: `admin/src/api/plugins.ts`

- [ ] **Step 1: 新增后台评价管理页面，支持筛选、详情查看、单条管理员回复。**
- [ ] **Step 2: 接入后台评价管理 API，菜单加入后台导航。**
- [ ] **Step 3: admin 构建验证。**
Run: `npm run build`
Expected: PASS

- [ ] **Step 4: Commit**
```bash
git add admin/src/views/review/ReviewList.vue admin/src/router/index.ts admin/src/layouts/AdminLayout.vue admin/src/api/plugins.ts
git commit -m "新增后台评价管理与回复页面" -m "- 后台支持评价列表、详情和管理员单条回复。\n- 菜单与路由接入评价管理入口。"
```

### Task 9: 文档同步与总验证

**Files:**
- Modify: `docs-site/docs/api/order.md`
- Create or Modify: `docs-site/docs/api/product.md`
- Modify: `docs-site/docs/guide/features.md`

- [ ] **Step 1: 文档同步接口升级与评价能力说明。**
- [ ] **Step 2: 全量验证。**
Run: `go test ./...` (workdir `server`)
Expected: PASS
Run: `npm run build` (workdir `admin`)
Expected: PASS
Run: `npm run build` (workdir `web`)
Expected: PASS
Run: `npm run build:h5:demo` (workdir `app`)
Expected: PASS
Run: `npm run docs:build` (workdir `docs-site`)
Expected: PASS

- [ ] **Step 3: Commit**
```bash
git add docs-site/docs/api/order.md docs-site/docs/api/product.md docs-site/docs/guide/features.md
git commit -m "同步评价系统接口与功能文档" -m "- 更新订单评价接口为create/edit/append模式。\n- 新增商品评价查询接口文档与前台评价Tab说明。"
```
