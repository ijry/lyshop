# EApp 商家移动端（完整版）Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 新建独立 `eapp/` 工程，交付支持 H5/微信小程序/App 的商家移动端完整版，并复用现有后台接口完成核心经营闭环。

**Architecture:** 采用 `uni-app + uview-plus` 构建独立移动商家端，按“工作台、订单、商品、营销、我的”5 Tab 组织信息架构。接口统一复用 `/admin/api`，登录态与权限语义对齐现有 `admin`，但本地存储键独立，避免互相污染。通过模块化 API、Store 与通用组件保证跨端一致渲染与后续扩展。

**Tech Stack:** uni-app (Vue3 + TypeScript), uview-plus, pinia, vue-i18n, vitest, docs-site (VitePress)

---

## Scope Check

该需求覆盖多个子系统（鉴权、工作台、订单、商品、营销、消息/权限、文档与部署）。本计划将其拆分为可独立交付的任务序列，每个任务都可单独回归。若执行中需要进一步并行，可再拆成子计划（如“交易域计划”“营销域计划”）。

## File Structure

- Create: `eapp/package.json`
- Create: `eapp/tsconfig.json`
- Create: `eapp/env.d.ts`
- Create: `eapp/vite.config.ts`
- Create: `eapp/vitest.config.ts`
- Create: `eapp/main.ts`
- Create: `eapp/App.vue`
- Create: `eapp/uni.scss`
- Create: `eapp/pages.json`
- Create: `eapp/manifest.json`
- Create: `eapp/uno.config.ts`
- Create: `eapp/locales/index.ts`
- Create: `eapp/locales/zh-CN.ts`
- Create: `eapp/locales/en.ts`
- Create: `eapp/utils/request.ts`
- Create: `eapp/utils/permission.ts`
- Create: `eapp/utils/storage.ts`
- Create: `eapp/api/auth.ts`
- Create: `eapp/api/dashboard.ts`
- Create: `eapp/api/order.ts`
- Create: `eapp/api/product.ts`
- Create: `eapp/api/marketing.ts`
- Create: `eapp/api/message.ts`
- Create: `eapp/api/system.ts`
- Create: `eapp/stores/auth.ts`
- Create: `eapp/stores/badge.ts`
- Create: `eapp/components/layout/EappShell.vue`
- Create: `eapp/components/common/StatCard.vue`
- Create: `eapp/components/common/TodoCard.vue`
- Create: `eapp/components/common/StatusTag.vue`
- Create: `eapp/pages/login/index.vue`
- Create: `eapp/pages/dashboard/index.vue`
- Create: `eapp/pages/order/list.vue`
- Create: `eapp/pages/order/detail.vue`
- Create: `eapp/pages/order/after-sale-list.vue`
- Create: `eapp/pages/order/after-sale-detail.vue`
- Create: `eapp/pages/product/list.vue`
- Create: `eapp/pages/product/edit.vue`
- Create: `eapp/pages/marketing/index.vue`
- Create: `eapp/pages/marketing/coupon.vue`
- Create: `eapp/pages/marketing/seckill.vue`
- Create: `eapp/pages/marketing/group-buy.vue`
- Create: `eapp/pages/marketing/bargain.vue`
- Create: `eapp/pages/me/index.vue`
- Create: `eapp/pages/me/messages.vue`
- Create: `eapp/pages/me/im-sessions.vue`
- Create: `eapp/pages/me/site-settings.vue`
- Create: `eapp/pages/me/admins.vue`
- Create: `eapp/pages/me/roles.vue`
- Create: `eapp/tests/utils/request.spec.ts`
- Create: `eapp/tests/utils/permission.spec.ts`
- Create: `eapp/tests/stores/badge.spec.ts`
- Modify: `docs-site/docs/guide/features.md`
- Create: `docs-site/docs/guide/eapp-merchant.md`
- Modify: `docs-site/docs/api/admin.md`
- Modify: `docs-site/docs/deploy/index.md`

---

### Task 1: 初始化 `eapp/` 工程骨架与构建脚本

**Files:**
- Create: `eapp/package.json`
- Create: `eapp/tsconfig.json`
- Create: `eapp/env.d.ts`
- Create: `eapp/vite.config.ts`
- Create: `eapp/manifest.json`
- Create: `eapp/pages.json`
- Create: `eapp/main.ts`
- Create: `eapp/App.vue`
- Create: `eapp/uni.scss`

- [ ] **Step 1: 先写最小脚手架并定义三端脚本**

```json
{
  "name": "lyshop-eapp",
  "private": true,
  "scripts": {
    "dev:h5": "cross-env UNI_INPUT_DIR=. uni",
    "build:h5": "cross-env UNI_INPUT_DIR=. uni build",
    "dev:mp-weixin": "cross-env UNI_INPUT_DIR=. uni -p mp-weixin",
    "build:mp-weixin": "cross-env UNI_INPUT_DIR=. uni build -p mp-weixin",
    "dev:app": "cross-env UNI_INPUT_DIR=. uni -p app",
    "test": "vitest run"
  }
}
```

- [ ] **Step 2: 在 `pages.json` 建立 5 Tab 基础路由**

```json
{
  "pages": [
    { "path": "pages/dashboard/index", "style": { "navigationStyle": "custom" } },
    { "path": "pages/order/list", "style": { "navigationBarTitleText": "订单" } },
    { "path": "pages/product/list", "style": { "navigationBarTitleText": "商品" } },
    { "path": "pages/marketing/index", "style": { "navigationBarTitleText": "营销" } },
    { "path": "pages/me/index", "style": { "navigationBarTitleText": "我的" } },
    { "path": "pages/login/index", "style": { "navigationBarTitleText": "登录" } }
  ]
}
```

- [ ] **Step 3: 运行依赖安装与最小构建**

Run: `cd eapp && npm install --legacy-peer-deps`  
Expected: install 完成，无 peer 依赖阻塞。

Run: `cd eapp && npm run build:h5`  
Expected: H5 构建成功，生成 dist 目录。

- [ ] **Step 4: Commit**

```bash
git add eapp
git commit -m "初始化商家移动端eapp工程骨架" -m "主要改动：新增eapp独立工程、基础构建脚本、页面路由与三端运行入口。原因：建立商家端独立交付边界，避免与用户端app耦合。影响范围：新增eapp目录，不影响现有前端工程。"
```

---

### Task 2: 建立请求层、鉴权存储与权限工具（含单测）

**Files:**
- Create: `eapp/utils/request.ts`
- Create: `eapp/utils/storage.ts`
- Create: `eapp/utils/permission.ts`
- Create: `eapp/tests/utils/request.spec.ts`
- Create: `eapp/tests/utils/permission.spec.ts`
- Create: `eapp/vitest.config.ts`

- [ ] **Step 1: 先写失败测试（请求响应拆包与权限判断）**

```ts
import { describe, it, expect } from 'vitest'
import { unwrapResponse, hasPermission } from '@/utils/permission'

describe('permission', () => {
  it('supports wildcard permission', () => {
    expect(hasPermission(['*'], 'order.ship')).toBe(true)
  })
})
```

- [ ] **Step 2: 运行测试确认失败**

Run: `cd eapp && npm run test`  
Expected: FAIL，提示方法或模块尚未实现。

- [ ] **Step 3: 实现请求封装与统一错误处理**

```ts
const request = axios.create({ baseURL: '/admin/api', timeout: 30000 })
request.interceptors.request.use((config) => {
  const token = getStorage('eapp_admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})
export function unwrapResponse<T>(raw: any): T {
  const { code, msg, data } = raw || {}
  if (code !== 0) throw new Error(msg || '请求失败')
  return data as T
}
```

- [ ] **Step 4: 实现权限工具并通过测试**

```ts
export function hasPermission(perms: string[], permission: string): boolean {
  if (!permission) return true
  return perms.includes('*') || perms.includes(permission)
}
```

Run: `cd eapp && npm run test`  
Expected: PASS。

- [ ] **Step 5: Commit**

```bash
git add eapp/utils eapp/tests eapp/vitest.config.ts
git commit -m "实现eapp请求层与权限基础能力" -m "主要改动：封装/admin/api请求链路、独立token存储键与权限判断工具，并补充单元测试。原因：保证接口复用后台语义的同时隔离端内登录态。影响范围：eapp基础设施层。"
```

---

### Task 3: 建立 i18n、设计令牌与通用布局组件

**Files:**
- Create: `eapp/locales/index.ts`
- Create: `eapp/locales/zh-CN.ts`
- Create: `eapp/locales/en.ts`
- Create: `eapp/components/layout/EappShell.vue`
- Create: `eapp/components/common/StatCard.vue`
- Create: `eapp/components/common/TodoCard.vue`
- Create: `eapp/components/common/StatusTag.vue`
- Modify: `eapp/App.vue`
- Modify: `eapp/uni.scss`

- [ ] **Step 1: 定义主题变量与卡片样式 Token**

```scss
:root {
  --eapp-primary: #2563EB;
  --eapp-success: #16A34A;
  --eapp-warning: #F59E0B;
  --eapp-danger: #DC2626;
  --eapp-accent: #F97316;
  --eapp-bg: #F8FAFC;
  --eapp-text: #1E293B;
}
```

- [ ] **Step 2: 搭建统一页面壳组件**

```vue
<template>
  <view class="eapp-shell">
    <slot name="header" />
    <view class="eapp-shell__body"><slot /></view>
  </view>
</template>
```

- [ ] **Step 3: 在 `main.ts` 接入 `uview-plus + i18n`**

```ts
const app = createSSRApp(App)
app.use(uviewPlus)
app.use(i18n)
```

- [ ] **Step 4: 构建验证**

Run: `cd eapp && npm run build:h5`  
Expected: PASS。

- [ ] **Step 5: Commit**

```bash
git add eapp/locales eapp/components eapp/App.vue eapp/uni.scss eapp/main.ts
git commit -m "落地eapp设计系统与通用布局组件" -m "主要改动：新增多语言资源、主题令牌与商家端通用页面壳和卡片组件。原因：统一视觉和组件语义，支撑后续模块快速搭建。影响范围：eapp UI基础层。"
```

---

### Task 4: 实现登录流程与鉴权守卫

**Files:**
- Create: `eapp/api/auth.ts`
- Create: `eapp/stores/auth.ts`
- Create: `eapp/pages/login/index.vue`
- Modify: `eapp/App.vue`

- [ ] **Step 1: 先写失败测试（无 token 进入业务页应被拦截）**

```ts
it('redirects to login when token is missing', () => {
  expect(shouldRedirectToLogin('', '/pages/dashboard/index')).toBe(true)
})
```

- [ ] **Step 2: 实现登录 API 与 store**

```ts
export const login = (username: string, password: string) =>
  request.post('/auth/login', { username, password })
```

```ts
const token = ref(getStorage('eapp_admin_token') || '')
function logout() {
  token.value = ''
  removeStorage('eapp_admin_token')
  uni.reLaunch({ url: '/pages/login/index' })
}
```

- [ ] **Step 3: 登录页完成账号密码提交与错误提示**

```vue
<up-input v-model="form.username" placeholder="请输入账号" />
<up-input v-model="form.password" type="password" placeholder="请输入密码" />
<up-button type="primary" :loading="loading" @click="onSubmit">登录</up-button>
```

- [ ] **Step 4: 回归验证**

Run: `cd eapp && npm run build:h5`  
Expected: PASS。

- [ ] **Step 5: Commit**

```bash
git add eapp/api/auth.ts eapp/stores/auth.ts eapp/pages/login/index.vue eapp/App.vue
git commit -m "实现eapp登录与鉴权守卫" -m "主要改动：接入后台登录接口、独立token存储、无登录态跳转登录页。原因：建立商家端安全入口并与后台鉴权语义保持一致。影响范围：eapp鉴权链路。"
```

---

### Task 5: 工作台模块（经营概览 + 待办 + 徽标数据）

**Files:**
- Create: `eapp/api/dashboard.ts`
- Create: `eapp/stores/badge.ts`
- Create: `eapp/tests/stores/badge.spec.ts`
- Create: `eapp/pages/dashboard/index.vue`

- [ ] **Step 1: 先写失败测试（徽标聚合计算）**

```ts
it('aggregates todo counters for tab badges', () => {
  expect(toBadgeCount({ pending_ship: 3, unread_message: 2 })).toBe(5)
})
```

- [ ] **Step 2: 接入工作台统计与待办接口**

```ts
export const getDashboardOverview = () => request.get('/dashboard/overview')
export const getDashboardTodo = () => request.get('/dashboard/todo')
```

- [ ] **Step 3: 实现 `badge` store 与首页联动刷新**

```ts
const badges = ref({ order: 0, message: 0 })
function syncFromTodo(todo: any) {
  badges.value.order = Number(todo?.pending_ship || 0) + Number(todo?.after_sale_processing || 0)
  badges.value.message = Number(todo?.unread_message || 0)
}
```

- [ ] **Step 4: 运行测试与构建**

Run: `cd eapp && npm run test`  
Expected: PASS。

Run: `cd eapp && npm run build:h5`  
Expected: PASS。

- [ ] **Step 5: Commit**

```bash
git add eapp/api/dashboard.ts eapp/stores/badge.ts eapp/tests/stores/badge.spec.ts eapp/pages/dashboard/index.vue
git commit -m "完成eapp工作台与待办徽标联动" -m "主要改动：实现经营概览、待办数据拉取与Tab徽标同步。原因：满足移动商家端任务优先的核心入口需求。影响范围：eapp工作台模块。"
```

---

### Task 6: 订单模块（列表/详情/发货/售后）

**Files:**
- Create: `eapp/api/order.ts`
- Create: `eapp/pages/order/list.vue`
- Create: `eapp/pages/order/detail.vue`
- Create: `eapp/pages/order/after-sale-list.vue`
- Create: `eapp/pages/order/after-sale-detail.vue`

- [ ] **Step 1: 先写失败测试（订单状态映射）**

```ts
it('maps order status to readable labels', () => {
  expect(toOrderStatusLabel('pending_ship')).toBe('待发货')
})
```

- [ ] **Step 2: 接入订单与售后 API**

```ts
export const listOrders = (params: any) => request.get('/orders', { params })
export const getOrderDetail = (id: number) => request.get(`/orders/${id}`)
export const shipOrder = (id: number, payload: any) => request.post(`/orders/${id}/shipments`, payload)
export const listAfterSales = (params: any) => request.get('/orders/after-sales', { params })
```

- [ ] **Step 3: 实现订单页状态筛选与详情跳转**

```vue
<up-tabs :list="tabs" :current="tabIndex" @click="onTabChange" />
<view v-for="order in list" :key="order.id" @click="goDetail(order.id)">
  <StatusTag :type="order.status" :text="order.status_label" />
</view>
```

- [ ] **Step 4: 实现发货弹窗与售后处理页核心动作**

```ts
await shipOrder(orderId.value, { company_code: form.companyCode, tracking_no: form.trackingNo })
```

- [ ] **Step 5: 构建验证**

Run: `cd eapp && npm run build:h5`  
Expected: PASS。

- [ ] **Step 6: Commit**

```bash
git add eapp/api/order.ts eapp/pages/order
git commit -m "完成eapp订单与售后主链路页面" -m "主要改动：实现订单列表、详情、发货和售后列表/详情处理页面并接入后台接口。原因：订单履约是商家端最高频业务链路。影响范围：eapp订单模块。"
```

---

### Task 7: 商品模块（列表/上下架/库存价格快捷改/编辑）

**Files:**
- Create: `eapp/api/product.ts`
- Create: `eapp/pages/product/list.vue`
- Create: `eapp/pages/product/edit.vue`

- [ ] **Step 1: 先写失败测试（价格与库存输入校验）**

```ts
it('rejects negative stock and price', () => {
  expect(validateQuickEdit({ stock: -1, price: 10 })).toBe(false)
})
```

- [ ] **Step 2: 实现商品 API 封装**

```ts
export const listProducts = (params: any) => request.get('/products', { params })
export const updateProductStatus = (id: number, on_sale: boolean) => request.put(`/products/${id}/status`, { on_sale })
export const quickUpdateSku = (id: number, payload: any) => request.put(`/products/${id}/sku-quick`, payload)
```

- [ ] **Step 3: 商品列表页接入搜索、筛选、快捷操作**

```vue
<up-search v-model="keyword" @search="reload" />
<up-switch :modelValue="item.on_sale" @change="onToggleSale(item)" />
<up-button size="mini" @click="openQuickEdit(item)">改库存/价格</up-button>
```

- [ ] **Step 4: 商品编辑页先覆盖高频字段**

```ts
const form = reactive({ title: '', subtitle: '', category_id: 0, cover: '', images: [] as string[] })
```

- [ ] **Step 5: 构建验证并提交**

Run: `cd eapp && npm run build:h5`  
Expected: PASS。

```bash
git add eapp/api/product.ts eapp/pages/product
git commit -m "完成eapp商品管理核心页面" -m "主要改动：实现商品列表、上下架、库存价格快捷改与高频字段编辑页面。原因：满足移动端商品日常运营效率需求。影响范围：eapp商品模块。"
```

---

### Task 8: 营销模块（优惠券/秒杀/拼团/砍价）

**Files:**
- Create: `eapp/api/marketing.ts`
- Create: `eapp/pages/marketing/index.vue`
- Create: `eapp/pages/marketing/coupon.vue`
- Create: `eapp/pages/marketing/seckill.vue`
- Create: `eapp/pages/marketing/group-buy.vue`
- Create: `eapp/pages/marketing/bargain.vue`

- [ ] **Step 1: 先写失败测试（营销入口路由映射）**

```ts
it('maps marketing card key to route', () => {
  expect(resolveMarketingRoute('coupon')).toBe('/pages/marketing/coupon')
})
```

- [ ] **Step 2: 接入营销 API**

```ts
export const listCoupons = (params: any) => request.get('/marketing/coupons', { params })
export const listSeckillActivities = (params: any) => request.get('/marketing/seckill/activities', { params })
export const listGroupBuyActivities = (params: any) => request.get('/marketing/group-buy/activities', { params })
export const listBargainActivities = (params: any) => request.get('/marketing/bargain/activities', { params })
```

- [ ] **Step 3: 营销首页实现聚合入口卡片**

```vue
<TodoCard title="优惠券" desc="查看与管理券规则" @click="go('/pages/marketing/coupon')" />
```

- [ ] **Step 4: 各列表页实现查询与状态展示**

```vue
<StatusTag :type="row.status" :text="row.status_label || row.status" />
```

- [ ] **Step 5: 构建验证并提交**

Run: `cd eapp && npm run build:h5`  
Expected: PASS。

```bash
git add eapp/api/marketing.ts eapp/pages/marketing
git commit -m "完成eapp营销模块入口与列表页面" -m "主要改动：实现营销聚合首页及优惠券、秒杀、拼团、砍价列表管理页面。原因：覆盖商家移动端完整运营模块范围。影响范围：eapp营销模块。"
```

---

### Task 9: 我的模块（消息/IM/店铺/员工/角色）

**Files:**
- Create: `eapp/api/message.ts`
- Create: `eapp/api/system.ts`
- Create: `eapp/pages/me/index.vue`
- Create: `eapp/pages/me/messages.vue`
- Create: `eapp/pages/me/im-sessions.vue`
- Create: `eapp/pages/me/site-settings.vue`
- Create: `eapp/pages/me/admins.vue`
- Create: `eapp/pages/me/roles.vue`

- [ ] **Step 1: 先写失败测试（菜单权限过滤）**

```ts
it('filters me menu by permissions', () => {
  const rows = filterMeMenus(['system.admin.list'], allMenus)
  expect(rows.some((r) => r.key === 'roles')).toBe(false)
})
```

- [ ] **Step 2: 封装消息、IM、系统配置与账号权限 API**

```ts
export const listMessages = (params: any) => request.get('/messages', { params })
export const sendMessage = (payload: any) => request.post('/messages/send', payload)
export const listSessions = (params: any) => request.get('/im/sessions', { params })
export const getSiteSettings = () => request.get('/system/site')
export const listAdmins = () => request.get('/admins')
export const listRoles = () => request.get('/roles')
```

- [ ] **Step 3: 我的首页按权限展示功能入口**

```ts
const menus = computed(() => filterMeMenus(authStore.perms, rawMenus))
```

- [ ] **Step 4: 完成消息发送与设置页保存**

```ts
await sendMessage({ title: form.title, content: form.content, receiver_type: form.receiverType })
await updateSiteSettings(siteForm)
```

- [ ] **Step 5: 构建验证并提交**

Run: `cd eapp && npm run build:h5`  
Expected: PASS。

```bash
git add eapp/api/message.ts eapp/api/system.ts eapp/pages/me
git commit -m "完成eapp我的模块与后台管理入口" -m "主要改动：实现消息、IM会话、店铺设置、管理员与角色页面并接入权限过滤。原因：补齐完整版商家端管理能力边界。影响范围：eapp我的模块。"
```

---

### Task 10: 多端适配与交互细节收口

**Files:**
- Modify: `eapp/pages.json`
- Modify: `eapp/components/layout/EappShell.vue`
- Modify: `eapp/pages/*`（涉及导航栏、安全区、键盘行为）
- Modify: `eapp/manifest.json`

- [ ] **Step 1: 适配安全区与底部 Tab 占位**

```scss
.safe-bottom {
  padding-bottom: calc(24rpx + env(safe-area-inset-bottom));
}
```

- [ ] **Step 2: 平台差异处理（H5/小程序/App）**

```ts
// #ifdef MP-WEIXIN
const navHeight = 88
// #endif
// #ifdef APP-PLUS
const navHeight = 96
// #endif
```

- [ ] **Step 3: 表单输入键盘类型与返回行为校验**

```vue
<up-input inputmode="numeric" v-model="form.stock" />
```

- [ ] **Step 4: 三端构建验证**

Run:
- `cd eapp && npm run build:h5`
- `cd eapp && npm run build:mp-weixin`

Expected: 均 PASS。

- [ ] **Step 5: Commit**

```bash
git add eapp/pages.json eapp/manifest.json eapp/components/layout/EappShell.vue eapp/pages
git commit -m "完成eapp三端适配与交互细节优化" -m "主要改动：补齐安全区、平台条件编译、输入行为与导航交互的一致性处理。原因：确保H5/小程序/App体验稳定可用。影响范围：eapp全局UI交互层。"
```

---

### Task 11: docs-site 文档同步（功能/接口/部署）

**Files:**
- Create: `docs-site/docs/guide/eapp-merchant.md`
- Modify: `docs-site/docs/guide/features.md`
- Modify: `docs-site/docs/api/admin.md`
- Modify: `docs-site/docs/deploy/index.md`

- [ ] **Step 1: 新增 eapp 功能文档**

```md
# 商家移动端 eapp
## 模块
- 工作台
- 订单
- 商品
- 营销
- 我的（消息、店铺、员工与权限）
```

- [ ] **Step 2: 在功能总览补充 eapp 能力项**

```md
- 新增独立商家移动端工程 `eapp`，支持 H5/微信小程序/App 三端。
```

- [ ] **Step 3: 在 API 文档写明接口复用策略**

```md
eapp 复用现有后台接口（`/admin/api/*`），未新增商家专属前缀。
```

- [ ] **Step 4: 在部署文档补充 eapp 构建发布步骤与配置影响**

```md
cd eapp
npm install --legacy-peer-deps
npm run build:h5
```

- [ ] **Step 5: 文档构建验证**

Run: `cd docs-site && npm run docs:build`  
Expected: PASS。

- [ ] **Step 6: Commit**

```bash
git add docs-site/docs/guide/eapp-merchant.md docs-site/docs/guide/features.md docs-site/docs/api/admin.md docs-site/docs/deploy/index.md
git commit -m "同步eapp商家移动端文档" -m "主要改动：新增eapp功能文档并同步功能总览、接口说明与部署说明。原因：功能变更需同步文档并覆盖接口与部署影响。影响范围：docs-site文档站。"
```

---

### Task 12: 全量回归与交付检查

**Files:**
- Modify: 无（仅验证，若失败则最小修复相关文件）

- [ ] **Step 1: eapp 单测回归**

Run: `cd eapp && npm run test`  
Expected: PASS。

- [ ] **Step 2: eapp 多端构建回归**

Run:
- `cd eapp && npm run build:h5`
- `cd eapp && npm run build:mp-weixin`

Expected: PASS。

- [ ] **Step 3: 文档构建回归**

Run: `cd docs-site && npm run docs:build`  
Expected: PASS。

- [ ] **Step 4: 手工验收清单**

```md
1) 未登录访问业务页会跳转登录页
2) 工作台能展示经营概览与待办，Tab徽标可刷新
3) 订单页可筛选状态，详情页可发货，售后页可处理
4) 商品页可上下架和快捷改库存/价格
5) 营销页四类活动列表可浏览
6) 我的页可进入消息、店铺、管理员、角色并受权限控制
7) H5/小程序/App 页面结构一致，无明显错位
```

- [ ] **Step 5: 汇总提交（若本任务产生修复）**

```bash
git status
git log --oneline -n 12
```

---

## Self-Review

### 1) Spec coverage

- 独立工程 `eapp/`：Task 1
- 三端支持：Task 1、Task 10、Task 12
- 复用后台接口：Task 2、Task 4~9
- 完整版功能范围：Task 5~9
- UI 设计系统与交互：Task 3、Task 10
- 文档同步（功能/接口/部署）：Task 11

### 2) Placeholder scan

- 已检查：无占位词、无“稍后实现”类描述。

### 3) Type consistency

- 统一使用 `eapp_admin_token` 作为登录态键。
- 统一 API 前缀 `/admin/api`。
- 权限判断统一走 `hasPermission(perms, permission)`。
- 订单/商品/营销/系统模块 API 命名保持 `list/get/update` 前缀一致。
