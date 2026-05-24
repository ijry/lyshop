# Admin AI融合与订单链路增强 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 完成“商品编辑内AI生图（支持参考图）+ 商品详情JSON block全链路改造 + 多端订单列表/详情增强 + 后端真实接口与mock一致”改造。

**Architecture:** 在现有 `ai/generate` 与 `orders` 资源语义内扩展字段，避免平行新增接口；商品详情采用 `version+blocks` JSON schema 并在 admin/web/app 全链路统一消费。AI模型能力通过 `supports_ref_image` 显式声明，前端根据能力禁用参考图输入并提示。

**Tech Stack:** Go + Gin + GORM, Vue3 + Vite + Tailwind (admin/web), uni-app + uview-plus (app), 本地 mock 路由

---

## File Structure (Design Lock)

- `server/plugins/ai_image/model/ai_image.go`: AI 模型能力字段与任务业务字段扩展
- `server/core/driver/ai/ai.go`: 驱动参数增加参考图字段
- `server/plugins/ai_image/service/ai_image.go`: 生成参数透传、能力校验、任务落库扩展
- `server/plugins/ai_image/api/admin.go`: 生成请求 DTO 升级
- `server/plugins/order/service/order.go`: 列表返回订单+items+金额明细；新增详情服务
- `server/plugins/order/api/front.go`: 新增前台订单详情接口，升级列表返回
- `server/plugins/order/api/admin.go`: 新增后台订单详情接口，升级列表返回
- `server/plugins/product/model/product.go`: 商品 detail 存储改为 JSON schema
- `server/plugins/product/service/product.go`: 商品详情 JSON 结构返回与兜底
- `admin/src/views/product/ProductForm.vue`: 商品编辑内 AI 助手 + JSON block 编辑器
- `admin/src/views/order/OrderList.vue`: 表格增强、状态切换、跳转详情
- `admin/src/views/order/OrderDetail.vue`: 新增后台详情页
- `admin/src/layouts/AdminLayout.vue`: Dashboard 固定第一项
- `admin/src/router/index.ts`, `admin/src/api/plugins.ts`: 后台详情路由与接口
- `admin/src/mock/index.ts`: AI能力字段、订单列表过滤与详情 mock
- `web/src/views/OrderList.vue`, `web/src/views/OrderDetail.vue`, `web/src/router/index.ts`, `web/src/api/request.ts`, `web/src/mock/index.ts`: PC订单链路增强
- `app/pages/order/list.vue`, `app/pages/order/detail.vue`, `app/pages.json`, `app/pages/user/index.vue`, `app/mock/index.ts`, `app/mock/data/orders.json`: H5订单链路增强与图标修复
- `README.md`, `docs-site/docs/guide/features.md`, `docs-site/docs/api/order.md`, `docs-site/docs/api/im.md`: 文档同步

---

### Task 1: 后端 AI 模型能力与任务扩展

**Files:**
- Modify: `server/plugins/ai_image/model/ai_image.go`
- Modify: `server/core/driver/ai/ai.go`
- Modify: `server/plugins/ai_image/api/admin.go`
- Modify: `server/plugins/ai_image/service/ai_image.go`

- [ ] **Step 1: 先写/更新后端单测骨架（若目录无测试则新建）覆盖生成请求扩展字段绑定。**
Run: `go test ./plugins/ai_image/... -run TestGenerateRequest -v`
Expected: 初次失败（缺少字段或校验逻辑）。

- [ ] **Step 2: 增加模型字段 `supports_ref_image` 与任务字段 `biz_type/ref_image_url/target_product_id`，并完成 GORM 标签。**
Run: `go test ./plugins/ai_image/... -v`
Expected: 编译通过。

- [ ] **Step 3: 扩展 `GenerateParams` 增加 `RefImageURL`，服务层透传并在模型不支持时返回业务错误。**
Run: `go test ./plugins/ai_image/... -v`
Expected: 测试通过，包含“模型不支持参考图”场景。

- [ ] **Step 4: 升级 `POST /admin/api/ai/generate` DTO，接收 `biz_type/ref_image_url/target_product_id`。**
Run: `go test ./plugins/ai_image/... -v`
Expected: 接口绑定测试通过。

- [ ] **Step 5: Commit**
```bash
git add server/core/driver/ai/ai.go server/plugins/ai_image/model/ai_image.go server/plugins/ai_image/api/admin.go server/plugins/ai_image/service/ai_image.go
git commit -m "扩展AI生图接口支持参考图与业务类型" -m "- ai模型增加supports_ref_image能力字段。\n- ai任务增加biz_type/ref_image_url/target_product_id用于审计与回放。\n- 生成请求与驱动参数透传参考图URL，不支持模型返回明确错误。"
```

### Task 2: 后端订单列表/详情模型增强

**Files:**
- Modify: `server/plugins/order/service/order.go`
- Modify: `server/plugins/order/api/front.go`
- Modify: `server/plugins/order/api/admin.go`

- [ ] **Step 1: 写失败用例，验证列表返回包含 items 与 amount_breakdown。**
Run: `go test ./plugins/order/... -run TestListOrdersWithItems -v`
Expected: FAIL（当前仅返回 Order）。

- [ ] **Step 2: 在 service 增加订单聚合结构（order + items + amount_breakdown），并批量查询 order_items 避免 N+1。**
Run: `go test ./plugins/order/... -v`
Expected: 列表聚合测试通过。

- [ ] **Step 3: 新增订单详情 service（前台按 user_id 过滤，后台不按 user_id 过滤）。**
Run: `go test ./plugins/order/... -run TestGetOrderDetail -v`
Expected: 详情测试通过。

- [ ] **Step 4: 暴露 `GET /api/v1/orders/:id` 与 `GET /admin/api/orders/:id`，升级列表接口返回新结构。**
Run: `go test ./plugins/order/... -v`
Expected: API层测试通过。

- [ ] **Step 5: Commit**
```bash
git add server/plugins/order/service/order.go server/plugins/order/api/front.go server/plugins/order/api/admin.go
git commit -m "升级订单接口返回结构并新增详情接口" -m "- 列表接口返回订单商品明细与金额分解，支持多端统一渲染。\n- 新增前后台订单详情接口，补齐价格体系与地址快照展示能力。"
```

### Task 3: 管理后台商品编辑页集成 AI 助手

**Files:**
- Modify: `admin/src/views/product/ProductForm.vue`
- Modify: `admin/src/api/plugins.ts` (若需AI相关函数可放独立api文件)

- [ ] **Step 1: 为商品编辑页添加 AI 助手区域（目标类型、模型、prompt、参考图、结果面板）。**
- [ ] **Step 2: 接入上传接口与 `ai/generate`，加载模型能力并实现“能力禁用+提示文案”。**
- [ ] **Step 3: 实现结果应用：封面赋值、轮播追加、详情图插入当前 block 位置（提示“将插入到当前编辑位置”）。**
- [ ] **Step 4: 本地自测编辑页完整流程（上传参考图→生成→应用）。**
Run: `npm run build` (workdir `admin`)
Expected: build 通过。

- [ ] **Step 5: Commit**
```bash
git add admin/src/views/product/ProductForm.vue admin/src/api/plugins.ts
git commit -m "商品编辑页集成AI生图助手" -m "- 支持封面/轮播/详情/介绍图目标类型。\n- 支持参考图上传并按模型能力禁用不支持场景。\n- 详情图支持插入光标位置并给出明确提示。"
```

### Task 4: 管理后台 Dashboard置顶与订单页面正规化

**Files:**
- Modify: `admin/src/layouts/AdminLayout.vue`
- Modify: `admin/src/views/order/OrderList.vue`
- Create: `admin/src/views/order/OrderDetail.vue`
- Modify: `admin/src/router/index.ts`
- Modify: `admin/src/api/plugins.ts`

- [ ] **Step 1: 侧边栏菜单渲染调整为 Dashboard 第一项固定显示。**
- [ ] **Step 2: 后台订单列表增强字段（商品摘要、优惠、支付方式、状态动作）并接入详情跳转。**
- [ ] **Step 3: 新增后台订单详情页展示商品明细与价格体系。**
- [ ] **Step 4: 构建验证。**
Run: `npm run build` (workdir `admin`)
Expected: build 通过。

- [ ] **Step 5: Commit**
```bash
git add admin/src/layouts/AdminLayout.vue admin/src/views/order/OrderList.vue admin/src/views/order/OrderDetail.vue admin/src/router/index.ts admin/src/api/plugins.ts
git commit -m "完善后台导航与订单管理页面" -m "- Dashboard固定首位。\n- 订单列表补齐商品与金额信息，状态切换有效。\n- 新增订单详情页展示完整价格体系与订单信息。"
```

### Task 4.5: 商品详情 JSON schema 全链路改造

**Files:**
- Modify: `server/plugins/product/model/product.go`
- Modify: `server/plugins/product/service/product.go`
- Modify: `server/plugins/product/api/admin.go`
- Modify: `server/plugins/product/api/front.go`
- Modify: `admin/src/views/product/ProductForm.vue`
- Modify: `web/src/views/ProductDetail.vue`
- Modify: `app/pages/product/detail.vue`
- Modify: `admin/src/mock/index.ts`
- Modify: `web/src/mock/index.ts`
- Modify: `app/mock/data/product-detail.json`

- [ ] **Step 1: 定义 `detail` JSON schema（version + blocks）并改造后端 product model/service 返回结构。**
- [ ] **Step 2: 管理后台实现 JSON block 编辑器（text/image），支持增删排序和当前位置插入。**
- [ ] **Step 3: PC/H5 商品详情页改为 block 渲染，不再直接渲染 HTML。**
- [ ] **Step 4: 三端 mock 与真实结构对齐。**
- [ ] **Step 5: 构建验证。**
Run: `go test ./plugins/product/... -v` (workdir `server`)
Expected: PASS
Run: `npm run build` (workdir `admin`)
Expected: PASS
Run: `npm run build` (workdir `web`)
Expected: PASS
Run: `npm run build:h5:demo` (workdir `app`)
Expected: PASS

- [ ] **Step 6: Commit**
```bash
git add server/plugins/product/model/product.go server/plugins/product/service/product.go server/plugins/product/api/admin.go server/plugins/product/api/front.go admin/src/views/product/ProductForm.vue web/src/views/ProductDetail.vue app/pages/product/detail.vue admin/src/mock/index.ts web/src/mock/index.ts app/mock/data/product-detail.json
git commit -m "切换商品详情为JSON block全链路渲染" -m "- 商品detail统一为version+blocks结构。\\n- 管理后台提供自定义block编辑器并支持AI图片插入。\\n- web/app商品详情页改为结构化渲染。"
```

### Task 5: admin mock 与真实接口对齐

**Files:**
- Modify: `admin/src/mock/index.ts`

- [ ] **Step 1: mock `ai/models` 增加 `supports_ref_image`，`ai/generate` 支持扩展字段。**
- [ ] **Step 2: mock `orders` 支持 status 过滤与详情接口返回。**
- [ ] **Step 3: 验证 `npm run dev:demo` 下后台新页面流程可跑通。**

- [ ] **Step 4: Commit**
```bash
git add admin/src/mock/index.ts
git commit -m "补齐后台mock能力与真实接口结构一致" -m "- AI模型能力字段与生成请求扩展字段同步。\n- 订单列表过滤与详情返回结构同步真实接口。"
```

### Task 6: PC 订单链路增强

**Files:**
- Modify: `web/src/views/OrderList.vue`
- Create: `web/src/views/OrderDetail.vue`
- Modify: `web/src/router/index.ts`
- Modify: `web/src/api/request.ts`
- Modify: `web/src/mock/index.ts`

- [ ] **Step 1: 升级 web mock request 支持 query 参数透传，订单status切换生效。**
- [ ] **Step 2: 订单列表展示商品行与金额分解摘要。**
- [ ] **Step 3: 新增订单详情页并接路由。**
- [ ] **Step 4: 构建验证。**
Run: `npm run build` (workdir `web`)
Expected: build 通过。

- [ ] **Step 5: Commit**
```bash
git add web/src/views/OrderList.vue web/src/views/OrderDetail.vue web/src/router/index.ts web/src/api/request.ts web/src/mock/index.ts
git commit -m "增强PC订单列表并新增订单详情页" -m "- 订单列表补齐商品明细与价格摘要。\n- 新增订单详情页展示完整订单信息。\n- mock参数透传后状态tab切换生效。"
```

### Task 7: H5 订单链路增强与图标/按钮修复

**Files:**
- Modify: `app/pages/order/list.vue`
- Create: `app/pages/order/detail.vue`
- Modify: `app/pages.json`
- Modify: `app/pages/user/index.vue`
- Modify: `app/mock/index.ts`
- Modify: `app/mock/data/orders.json`

- [ ] **Step 1: 升级订单mock数据，补齐 `items` 与 `amount_breakdown`。**
- [ ] **Step 2: 订单列表显示商品行，操作按钮套容器避免按钮 100% 占满。**
- [ ] **Step 3: 新增订单详情页并配置路由。**
- [ ] **Step 4: 修复个人中心客服图标名称。**
- [ ] **Step 5: 构建验证。**
Run: `npm run build:h5:demo` (workdir `app`)
Expected: build 通过。

- [ ] **Step 6: Commit**
```bash
git add app/pages/order/list.vue app/pages/order/detail.vue app/pages.json app/pages/user/index.vue app/mock/index.ts app/mock/data/orders.json
git commit -m "完善H5订单展示并修复按钮与客服图标" -m "- 订单列表新增商品明细展示并优化操作按钮布局。\n- 新增订单详情页展示价格体系。\n- 修复个人中心客服图标显示异常。"
```

### Task 8: 文档同步与全量验证

**Files:**
- Modify: `README.md`
- Modify: `docs-site/docs/guide/features.md`
- Modify: `docs-site/docs/api/order.md`
- Modify: `docs-site/docs/api/im.md`

- [ ] **Step 1: 补充 AI 参考图能力、订单详情接口与多端订单展示变更说明。**
- [ ] **Step 2: 执行构建校验。**
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
git add README.md docs-site/docs/guide/features.md docs-site/docs/api/order.md docs-site/docs/api/im.md
git commit -m "同步AI与订单链路增强文档" -m "- 更新README与docs-site功能/API文档。\n- 覆盖AI参考图能力、订单详情接口与多端展示变化。"
```
