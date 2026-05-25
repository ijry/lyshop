# Admin Sidebar Grouped Tabs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将管理后台左侧导航升级为“分组 TAB + 菜单双列侧栏”，并由后端 `GET /admin/api/menus` 返回分组结构，`Dashboard` 固定入口不归组。

**Architecture:** 保持现有接口路径不变，仅升级响应结构；前端布局改为双列侧栏并实现 `hover 预览 + click 锁定` 交互；在兼容期保留旧结构降级渲染，避免前后端发布窗口不一致导致不可用。

**Tech Stack:** Vue 3 + TypeScript + Vue Router + Tailwind（admin），VitePress Markdown（docs-site）

---

## File Structure

- Modify: `admin/src/api/auth.ts`
  - 定义菜单类型（新结构 + 旧结构联合），升级 `getMenus` 返回类型。
- Modify: `admin/src/layouts/AdminLayout.vue`
  - 双列侧栏模板、分组状态机、路由归组匹配、深色滚动条样式、旧结构降级。
- Modify: `admin/src/mock/index.ts`
  - `GET /admin/api/menus` mock 升级为 `{ dashboard, groups }`。
- Modify: `docs-site/docs/api/auth.md`
  - 补充 `/admin/api/menus` 新响应结构、兼容说明、字段定义。
- Modify: `docs-site/docs/guide/features.md`
  - 补充管理后台导航功能说明与交互行为。

### Task 1: 定义菜单响应类型并升级 API 声明

**Files:**
- Modify: `admin/src/api/auth.ts`

- [ ] **Step 1: 写类型并保留旧结构兼容**

```ts
import request from './request'

export interface AdminMenuItem {
  title: string
  path: string
  icon?: string
  sort?: number
  children?: AdminMenuItem[]
}

export interface AdminMenuGroup {
  key: string
  title: string
  icon?: string
  sort?: number
  menus: AdminMenuItem[]
}

export interface AdminMenuGroupedResponse {
  dashboard: { title: string; path: string }
  groups: AdminMenuGroup[]
}

export type AdminMenuLegacyResponse = AdminMenuItem[]
export type AdminMenuResponse = AdminMenuGroupedResponse | AdminMenuLegacyResponse

export const login = (username: string, password: string) =>
  request.post<never, { token: string }>('/auth/login', { username, password })

export const getMenus = () =>
  request.get<never, AdminMenuResponse>('/menus')
```

- [ ] **Step 2: 运行类型检查（仅 admin）**

Run: `npm run build`（在 `admin` 目录）  
Expected: 构建通过，`auth.ts` 类型无报错。

### Task 2: 升级 mock 返回分组结构

**Files:**
- Modify: `admin/src/mock/index.ts`

- [ ] **Step 1: 将 `/admin/api/menus` 改为 grouped 响应**

```ts
'GET /admin/api/menus': {
  dashboard: { title: 'Dashboard', path: '/dashboard' },
  groups: [
    {
      key: 'product',
      title: '商品',
      icon: 'box',
      sort: 10,
      menus: [
        { title: '商品管理', icon: 'box', path: '/product', sort: 10, children: [
          { title: '商品列表', path: '/product/list' },
          { title: '商品分类', path: '/product/category' },
          { title: '新增商品', path: '/product/form' },
        ]},
        { title: '评价管理', icon: 'star', path: '/review', sort: 20, children: [
          { title: '评价列表', path: '/review/list' },
        ]},
      ],
    },
    {
      key: 'order',
      title: '订单',
      icon: 'shopping-cart',
      sort: 20,
      menus: [
        { title: '订单管理', icon: 'shopping-cart', path: '/order', sort: 10, children: [
          { title: '订单列表', path: '/order/list' },
          { title: '售后列表', path: '/order/after-sale/list' },
        ]},
      ],
    }
  ]
},
```

- [ ] **Step 2: 保持所有原菜单都被纳入对应 groups**

Run: `rg -n "GET /admin/api/menus" admin/src/mock/index.ts`  
Expected: 只有一个菜单定义且为对象结构。

### Task 3: 重构 AdminLayout 为双列分组侧栏

**Files:**
- Modify: `admin/src/layouts/AdminLayout.vue`

- [ ] **Step 1: 引入类型与状态计算**

```ts
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getMenus, type AdminMenuItem, type AdminMenuGroup, type AdminMenuGroupedResponse } from '@/api/auth'

const route = useRoute()
const groupedMenus = ref<AdminMenuGroup[]>([])
const legacyMenus = ref<AdminMenuItem[]>([])
const dashboardMenu = ref({ title: 'Dashboard', path: '/dashboard' })
const hoverGroupKey = ref('')
const lockedGroupKey = ref('')
```

- [ ] **Step 2: 新增解析函数与路由归组逻辑**

```ts
function isGroupedResponse(data: unknown): data is AdminMenuGroupedResponse {
  return !!data && typeof data === 'object' && Array.isArray((data as any).groups)
}

function menuContainsPath(menu: AdminMenuItem, path: string): boolean {
  if (menu.path === path) return true
  return Array.isArray(menu.children) && menu.children.some((child) => menuContainsPath(child, path))
}

const routeMatchedGroupKey = computed(() => {
  const path = route.path
  for (const group of groupedMenus.value) {
    if (group.menus.some((menu) => menuContainsPath(menu, path))) return group.key
  }
  return ''
})

const activeGroupKey = computed(() => {
  return hoverGroupKey.value || lockedGroupKey.value || routeMatchedGroupKey.value || groupedMenus.value[0]?.key || ''
})
```

- [ ] **Step 3: 双列模板 + hover/click 交互 + Dashboard 固定入口**

```vue
<aside class="w-80 bg-slate-900 text-slate-100 flex shrink-0">
  <div class="w-20 border-r border-slate-800 py-3 overflow-y-auto sidebar-scroll"
       @mouseleave="hoverGroupKey = ''">
    <button
      v-for="group in groupedMenus"
      :key="group.key"
      @mouseenter="hoverGroupKey = group.key"
      @click="lockedGroupKey = group.key"
      class="w-full px-2 py-2 text-xs rounded-md"
      :class="activeGroupKey === group.key ? 'bg-slate-700 text-white' : 'text-slate-300 hover:bg-slate-800'"
    >
      {{ group.title }}
    </button>
  </div>
  <nav class="flex-1 py-4 px-3 overflow-y-auto sidebar-scroll">
    <router-link to="/dashboard" ...>Dashboard</router-link>
    <!-- activeGroup menus 渲染 -->
  </nav>
</aside>
```

- [ ] **Step 4: 旧结构降级渲染分支**

```vue
<template v-if="groupedMenus.length"> ...新双列... </template>
<template v-else> ...保留原单列 legacyMenus 渲染... </template>
```

- [ ] **Step 5: 添加深色滚动条样式**

```vue
<style scoped>
.sidebar-scroll {
  scrollbar-color: #475569 #0f172a;
}
.sidebar-scroll::-webkit-scrollbar { width: 8px; }
.sidebar-scroll::-webkit-scrollbar-track { background: #0f172a; }
.sidebar-scroll::-webkit-scrollbar-thumb { background: #475569; border-radius: 9999px; }
.sidebar-scroll::-webkit-scrollbar-thumb:hover { background: #64748b; }
</style>
```

### Task 4: 对接加载逻辑与异常降级

**Files:**
- Modify: `admin/src/layouts/AdminLayout.vue`

- [ ] **Step 1: 在 onMounted 中区分新旧结构**

```ts
onMounted(async () => {
  try {
    const data = await getMenus()
    if (isGroupedResponse(data)) {
      groupedMenus.value = [...(data.groups || [])].sort((a, b) => Number(a.sort || 0) - Number(b.sort || 0))
      dashboardMenu.value = data.dashboard || dashboardMenu.value
      legacyMenus.value = []
      return
    }
    legacyMenus.value = Array.isArray(data) ? data.filter((item) => item.path !== '/dashboard') : []
    groupedMenus.value = []
  } catch {
    groupedMenus.value = []
    legacyMenus.value = []
  }
})
```

- [ ] **Step 2: 加入空菜单占位文案**

```vue
<div v-if="!activeGroupMenus.length" class="px-3 py-2 text-xs text-slate-500">暂无菜单</div>
```

### Task 5: 同步 docs-site API 与功能文档

**Files:**
- Modify: `docs-site/docs/api/auth.md`
- Modify: `docs-site/docs/guide/features.md`

- [ ] **Step 1: 在 API 文档写 `/menus` 新结构与兼容说明**

```md
### GET /admin/api/menus

返回管理后台菜单。当前版本返回分组结构：

```json
{
  "dashboard": { "title": "Dashboard", "path": "/dashboard" },
  "groups": [{ "key": "product", "title": "商品", "menus": [] }]
}
```

兼容说明：旧版可能返回数组结构，前端可按旧单列方式渲染。
```

- [ ] **Step 2: 在功能文档新增“后台导航分组 TAB”说明**

```md
- 左侧导航采用双列结构：分组 TAB + 菜单列。
- 分组交互支持 hover 预览、click 锁定。
- Dashboard 固定入口，不归属于任一分组。
- 侧栏滚动条采用深色主题，与导航背景统一。
```

### Task 6: 本地验证

**Files:**
- Test: `admin` 构建结果
- Test: `docs-site` 构建结果（可选）

- [ ] **Step 1: Admin 构建验证**

Run: `npm run build`（`admin`）  
Expected: 构建通过，无 TS/Vue 编译错误。

- [ ] **Step 2: 文档构建验证（可选但推荐）**

Run: `npm run docs:build` 或文档项目等效命令（`docs-site`）  
Expected: 文档构建通过，新增章节可被渲染。

- [ ] **Step 3: 手工验收清单**

Run: 启动后台，人工检查：
- hover 分组菜单是否预览。
- click 后是否锁定。
- 刷新后分组高亮是否匹配当前路由。
- Dashboard 是否始终固定可见。
- 深色滚动条是否统一。

## Self-Review

- Spec 覆盖：接口升级、双列分组、交互规则、Dashboard 不归组、滚动条、兼容降级、docs-site 同步均有对应任务。
- Placeholder 扫描：无 TBD/TODO/“后续补充”。
- 类型一致性：`AdminMenuItem/AdminMenuGroup/AdminMenuGroupedResponse` 在 API 与 Layout 中保持一致命名。
