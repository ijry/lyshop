# eapp 商家移动端 P1 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 `eapp/` 由 10 个薄页面升级为成熟移动端商家工作台（P1：核心交易闭环 + Dashboard + 商品 + 售后 + 高阶组件基础设施），同步刷新 mock 与 docs-site。

**Architecture:** 保留 uni-app + uview-plus + admin 共享 mock 架构；新增 `components/charts/`（封装 ly-charts）与 `components/biz/`（商家高阶组件库）；composables 沉淀分页/筛选/批量逻辑；mock 在 `admin/src/mock/index.ts` 内增量。

**Tech Stack:** uni-app (Vue3) / TypeScript / Pinia / uview-plus 3.x / ly-charts (qiun-data-charts uni_modules) / UnoCSS / vitest / vue-i18n / VitePress (docs)

**Spec:** `docs/superpowers/specs/2026-05-28-eapp-merchant-overhaul-design.md`

**预期 commit 划分**：详见 spec 7.5，共 8 个，本计划任务按 commit 分组（A-H）。

---

## 文件结构总览

**新建（约 26 个组件 + 9 个 composables/utils + 2 个 pages）：**

```
eapp/uni_modules/qiun-data-charts/          # ly-charts 组件（uni_modules 标准结构）
eapp/utils/ly-charts.ts                     # ly-charts 主题适配
eapp/components/charts/
  ├── ChartPanel.vue
  ├── AreaChart.vue
  ├── RingChart.vue
  └── BarChart.vue
eapp/components/biz/
  ├── PageHeader.vue
  ├── MetricCard.vue
  ├── ActionGrid.vue
  ├── TodoCenter.vue
  ├── AnnouncementBar.vue
  ├── FilterDrawer.vue
  ├── BatchBar.vue
  ├── BatchResultPopup.vue
  ├── Timeline.vue
  ├── EmptyState.vue
  ├── OrderCard.vue
  ├── ProductCard.vue
  ├── AfterSaleCard.vue
  ├── SkuMatrixEditor.vue
  ├── CategoryTreePicker.vue
  ├── RichTextEditor.vue
  ├── ShipPopup.vue
  ├── RepricingPopup.vue
  └── RemarkPopup.vue
eapp/composables/
  ├── useDashboard.ts
  ├── useOrderList.ts
  ├── useProductList.ts
  ├── useAfterSaleList.ts
  ├── useBatchSelection.ts
  ├── useFilter.ts
  └── useRequest.ts
eapp/api/
  ├── after-sale.ts                         # 从 order.ts 抽出
  └── category.ts                           # 新增
eapp/stores/shop.ts                         # 当前店铺缓存
eapp/pages/order/batch.vue                  # 批量操作集中页
eapp/pages/order/print-preview.vue          # 面单预览
eapp/pages/product/category-tree.vue        # 分类树管理
```

**修改：**

```
eapp/pages/dashboard/index.vue               # 重写
eapp/pages/order/list.vue                    # 重写
eapp/pages/order/detail.vue                  # 增强
eapp/pages/order/after-sale-list.vue         # 增强
eapp/pages/order/after-sale-detail.vue       # 增强
eapp/pages/product/list.vue                  # 重写
eapp/pages/product/edit.vue                  # 重写
eapp/pages.json                              # 新增 3 页注册
eapp/api/dashboard.ts                        # 扩展返回类型
eapp/api/order.ts                            # 扩展 + re-export after-sale
eapp/api/product.ts                          # 扩展 batch
eapp/api/system.ts                           # 扩展 shops/current + announcements
eapp/components/common/StatusTag.vue         # 色板扩展
eapp/components/layout/EappShell.vue         # safe-area 增强
eapp/stores/badge.ts                         # 接入新字段
eapp/locales/zh-CN.ts                        # i18n 补全
eapp/locales/en.ts                           # i18n 补全
eapp/App.vue                                 # CSS 变量扩展
admin/src/mock/index.ts                      # 新增 25 条路由 + 数据扩充
app/mock/data/orders.json                    # 订单数据扩展（如需）
docs-site/docs/guide/eapp-merchant.md        # 重写
docs-site/docs/api/order.md                  # 追加新接口
docs-site/docs/api/product.md                # 追加新接口
docs-site/docs/api/admin.md                  # 追加新接口
docs-site/docs/guide/features.md             # 移动端商家工作台一节
```

**测试（新增）：**

```
eapp/tests/composables/useBatchSelection.spec.ts
eapp/tests/composables/useFilter.spec.ts
eapp/tests/composables/useOrderList.spec.ts
eapp/tests/composables/useProductList.spec.ts
eapp/tests/composables/useDashboard.spec.ts
eapp/tests/mock/dashboard.spec.ts
eapp/tests/mock/orders.spec.ts
eapp/tests/mock/products.spec.ts
eapp/tests/mock/categories.spec.ts
eapp/tests/mock/after-sales.spec.ts
```

---

## Phase A — 引入 ly-charts 与商家高阶组件骨架（→ commit 1）

### Task A1：拉取 ly-charts uni_modules

**Files:**
- Create: `eapp/uni_modules/qiun-data-charts/`（含 `package.json` + `components/qiun-data-charts/` + 内部依赖目录）

- [ ] **Step 1：从 npm 安装 ly-charts，复制 dist 到 uni_modules**

```bash
cd eapp
npm install ly-charts@26.1.4 --save --legacy-peer-deps
# 将 node_modules/ly-charts/uni_modules/qiun-data-charts 整体复制到 eapp/uni_modules/
node -e "require('fs').cpSync('node_modules/ly-charts/uni_modules/qiun-data-charts', 'uni_modules/qiun-data-charts', {recursive:true})"
```

预期：`eapp/uni_modules/qiun-data-charts/package.json` 存在；目录包含 `components/qiun-data-charts/qiun-data-charts.vue`。

- [ ] **Step 2：验证依赖与编译**

```bash
cd eapp && npm run dev:h5 -- --mode demo
```

预期：dev server 起来，无 ly-charts 相关错误（即便页面尚未使用）。Ctrl+C 退出。

- [ ] **Step 3：暂不 commit，等 A4 全部完成统一 commit**

### Task A2：创建 ly-charts 主题适配

**Files:**
- Create: `eapp/utils/ly-charts.ts`

- [ ] **Step 1：写主题适配**

```ts
// eapp/utils/ly-charts.ts
export const EAPP_CHART_COLORS = ['#2563EB', '#F97316', '#16A34A', '#F59E0B', '#DC2626', '#8B5CF6']

export type ChartCategoriesSeries = {
  categories: string[]
  series: Array<{ name: string; data: number[] }>
}

export type ChartPie = Array<{ name: string; value: number }>

export function buildAreaOpts(data: ChartCategoriesSeries) {
  return {
    color: EAPP_CHART_COLORS,
    padding: [12, 12, 0, 0],
    enableScroll: false,
    legend: { show: data.series.length > 1 },
    xAxis: { disableGrid: true, axisLine: false, fontColor: '#94A3B8' },
    yAxis: { gridType: 'dash', dashLength: 4, splitNumber: 4, fontColor: '#94A3B8' },
    extra: { area: { type: 'curve', opacity: 0.16, addLine: true, lineMode: 'smooth' } },
  }
}

export function buildRingOpts() {
  return {
    color: EAPP_CHART_COLORS,
    padding: [4, 4, 4, 4],
    rotate: false,
    enableScroll: false,
    legend: { show: true, position: 'right', lineHeight: 24, fontColor: '#475569' },
    extra: { ring: { ringWidth: 32, activeOpacity: 0.6, activeRadius: 4, offsetAngle: 0, labelWidth: 8, border: false } },
  }
}

export function buildBarOpts(horizontal = false) {
  return {
    color: EAPP_CHART_COLORS,
    padding: [12, 12, 0, 0],
    enableScroll: false,
    legend: { show: false },
    xAxis: { disableGrid: !horizontal, axisLine: false, fontColor: '#94A3B8' },
    yAxis: { gridType: 'dash', dashLength: 4, splitNumber: 4, fontColor: '#94A3B8' },
    extra: { bar: { type: horizontal ? 'group' : 'group', width: 24, meterBorde: 0, meterFillColor: '#FFFFFF', barBorderRadius: [4, 4, 4, 4] } },
    rotate: horizontal,
  }
}
```

预期：文件创建。

### Task A3：编写并通过 ChartPanel 单元逻辑测试，再补 UI

**Files:**
- Test: `eapp/tests/utils/ly-charts.spec.ts`
- Create: `eapp/components/charts/ChartPanel.vue`、`AreaChart.vue`、`RingChart.vue`、`BarChart.vue`

- [ ] **Step 1：写 utils 测试**

```ts
// eapp/tests/utils/ly-charts.spec.ts
import { describe, expect, it } from 'vitest'
import { EAPP_CHART_COLORS, buildAreaOpts, buildBarOpts, buildRingOpts } from '@/utils/ly-charts'

describe('ly-charts opts builders', () => {
  it('uses brand palette', () => {
    expect(EAPP_CHART_COLORS[0]).toBe('#2563EB')
  })

  it('builds area opts that hide legend for single series', () => {
    const opts = buildAreaOpts({ categories: ['1'], series: [{ name: 'a', data: [1] }] })
    expect(opts.legend.show).toBe(false)
  })

  it('builds area opts that show legend for multi series', () => {
    const opts = buildAreaOpts({
      categories: ['1'],
      series: [{ name: 'a', data: [1] }, { name: 'b', data: [2] }],
    })
    expect(opts.legend.show).toBe(true)
  })

  it('builds ring opts with right legend', () => {
    const opts = buildRingOpts()
    expect(opts.legend.position).toBe('right')
  })

  it('builds bar opts with rotate flag', () => {
    expect(buildBarOpts(false).rotate).toBe(false)
    expect(buildBarOpts(true).rotate).toBe(true)
  })
})
```

- [ ] **Step 2：跑测试验证失败**

```bash
cd eapp && npx vitest run tests/utils/ly-charts.spec.ts
```

预期：FAIL — `@/utils/ly-charts` 不存在（Task A2 后会 PASS，若已 A2 完则直接 PASS）。

- [ ] **Step 3：跑测试验证通过（A2 已实现）**

```bash
cd eapp && npx vitest run tests/utils/ly-charts.spec.ts
```

预期：5 个测试全部 PASS。

- [ ] **Step 4：写 ChartPanel.vue（loading / empty / error / 工具栏槽）**

```vue
<!-- eapp/components/charts/ChartPanel.vue -->
<script setup lang="ts">
import { computed } from 'vue'
defineProps<{
  title?: string
  loading?: boolean
  empty?: boolean
  error?: string
  height?: number
}>()
const wrapHeight = computed(() => {
  // 默认 240rpx
})
</script>

<template>
  <view class="chart-panel">
    <view v-if="title || $slots.extra" class="chart-header">
      <text class="chart-title">{{ title }}</text>
      <view class="chart-extra"><slot name="extra" /></view>
    </view>
    <view class="chart-body" :style="{ height: (height || 320) + 'rpx' }">
      <view v-if="loading" class="chart-state">加载中...</view>
      <view v-else-if="error" class="chart-state error">{{ error }}</view>
      <view v-else-if="empty" class="chart-state">暂无数据</view>
      <slot v-else />
    </view>
  </view>
</template>

<style scoped>
.chart-panel { background: #fff; border: 1px solid var(--eapp-border); border-radius: 24rpx; padding: 20rpx; }
.chart-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12rpx; }
.chart-title { font-size: 28rpx; font-weight: 700; color: var(--eapp-text); }
.chart-extra { display: flex; gap: 8rpx; }
.chart-body { width: 100%; box-sizing: border-box; }
.chart-state { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: var(--eapp-text-muted); font-size: 24rpx; }
.chart-state.error { color: var(--eapp-danger); }
</style>
```

- [ ] **Step 5：写 AreaChart.vue（薄封装）**

```vue
<!-- eapp/components/charts/AreaChart.vue -->
<script setup lang="ts">
import { computed } from 'vue'
import ChartPanel from './ChartPanel.vue'
import { buildAreaOpts, type ChartCategoriesSeries } from '@/utils/ly-charts'

const props = defineProps<{
  title?: string
  data: ChartCategoriesSeries
  height?: number
  loading?: boolean
}>()

const opts = computed(() => buildAreaOpts(props.data))
const isEmpty = computed(() => !props.data?.categories?.length || !props.data?.series?.length)
</script>

<template>
  <ChartPanel :title="title" :loading="loading" :empty="isEmpty" :height="height">
    <template #extra><slot name="extra" /></template>
    <qiun-data-charts
      v-if="!isEmpty"
      type="area"
      :opts="opts"
      :chartData="{ categories: data.categories, series: data.series }"
      :animation="true"
      background="none"
    />
  </ChartPanel>
</template>
```

- [ ] **Step 6：写 RingChart.vue**

```vue
<!-- eapp/components/charts/RingChart.vue -->
<script setup lang="ts">
import { computed } from 'vue'
import ChartPanel from './ChartPanel.vue'
import { buildRingOpts, type ChartPie } from '@/utils/ly-charts'

const props = defineProps<{
  title?: string
  data: ChartPie
  height?: number
  loading?: boolean
}>()

const opts = computed(() => buildRingOpts())
const isEmpty = computed(() => !props.data?.length || props.data.every((row) => Number(row.value || 0) === 0))
</script>

<template>
  <ChartPanel :title="title" :loading="loading" :empty="isEmpty" :height="height">
    <qiun-data-charts
      v-if="!isEmpty"
      type="ring"
      :opts="opts"
      :chartData="{ series: [{ data }] }"
      :animation="true"
      background="none"
    />
  </ChartPanel>
</template>
```

- [ ] **Step 7：写 BarChart.vue**

```vue
<!-- eapp/components/charts/BarChart.vue -->
<script setup lang="ts">
import { computed } from 'vue'
import ChartPanel from './ChartPanel.vue'
import { buildBarOpts, type ChartCategoriesSeries } from '@/utils/ly-charts'

const props = defineProps<{
  title?: string
  data: ChartCategoriesSeries
  height?: number
  loading?: boolean
  horizontal?: boolean
}>()

const opts = computed(() => buildBarOpts(!!props.horizontal))
const isEmpty = computed(() => !props.data?.categories?.length || !props.data?.series?.length)
</script>

<template>
  <ChartPanel :title="title" :loading="loading" :empty="isEmpty" :height="height">
    <qiun-data-charts
      v-if="!isEmpty"
      :type="horizontal ? 'bar' : 'column'"
      :opts="opts"
      :chartData="{ categories: data.categories, series: data.series }"
      :animation="true"
      background="none"
    />
  </ChartPanel>
</template>
```

### Task A4：composables 核心三件套（TDD）

**Files:**
- Create: `eapp/composables/useBatchSelection.ts`、`useFilter.ts`、`useRequest.ts`
- Test: `eapp/tests/composables/useBatchSelection.spec.ts`、`useFilter.spec.ts`

- [ ] **Step 1：写 useBatchSelection 测试**

```ts
// eapp/tests/composables/useBatchSelection.spec.ts
import { describe, expect, it } from 'vitest'
import { useBatchSelection } from '@/composables/useBatchSelection'

describe('useBatchSelection', () => {
  it('starts empty', () => {
    const s = useBatchSelection<number>()
    expect(s.selected.value).toEqual([])
    expect(s.count.value).toBe(0)
  })

  it('toggles ids idempotently', () => {
    const s = useBatchSelection<number>()
    s.toggle(1)
    s.toggle(2)
    s.toggle(1)
    expect(s.selected.value.sort()).toEqual([2])
  })

  it('selectAll & clear', () => {
    const s = useBatchSelection<number>()
    s.selectAll([1, 2, 3])
    expect(s.count.value).toBe(3)
    s.clear()
    expect(s.count.value).toBe(0)
  })

  it('isSelected works', () => {
    const s = useBatchSelection<number>()
    s.toggle(5)
    expect(s.isSelected(5)).toBe(true)
    expect(s.isSelected(6)).toBe(false)
  })

  it('respects max limit (silent ignore)', () => {
    const s = useBatchSelection<number>({ max: 2 })
    s.toggle(1); s.toggle(2); s.toggle(3)
    expect(s.count.value).toBe(2)
    expect(s.isSelected(3)).toBe(false)
  })
})
```

- [ ] **Step 2：跑测试验证失败**

```bash
cd eapp && npx vitest run tests/composables/useBatchSelection.spec.ts
```

预期：FAIL — `useBatchSelection` 不存在。

- [ ] **Step 3：实现 useBatchSelection**

```ts
// eapp/composables/useBatchSelection.ts
import { computed, ref } from 'vue'

export function useBatchSelection<T>(opts?: { max?: number }) {
  const selected = ref<T[]>([])
  const max = Number(opts?.max || 0)

  function isSelected(id: T) {
    return selected.value.includes(id as any)
  }
  function toggle(id: T) {
    const idx = selected.value.indexOf(id as any)
    if (idx >= 0) {
      selected.value.splice(idx, 1)
      return
    }
    if (max > 0 && selected.value.length >= max) return
    selected.value.push(id as any)
  }
  function selectAll(ids: T[]) {
    selected.value = max > 0 ? ids.slice(0, max) : ids.slice()
  }
  function clear() {
    selected.value = []
  }
  const count = computed(() => selected.value.length)

  return { selected, count, toggle, isSelected, selectAll, clear }
}
```

- [ ] **Step 4：跑测试验证通过**

```bash
cd eapp && npx vitest run tests/composables/useBatchSelection.spec.ts
```

预期：5 个测试 PASS。

- [ ] **Step 5：写 useFilter 测试**

```ts
// eapp/tests/composables/useFilter.spec.ts
import { describe, expect, it, vi, beforeEach } from 'vitest'
import { useFilter } from '@/composables/useFilter'

const storage = new Map<string, string>()
;(globalThis as any).uni = {
  getStorageSync: (k: string) => storage.get(k) || '',
  setStorageSync: (k: string, v: string) => { storage.set(k, String(v)) },
  removeStorageSync: (k: string) => { storage.delete(k) },
}

describe('useFilter', () => {
  beforeEach(() => storage.clear())

  it('default value applied when nothing in storage', () => {
    const f = useFilter('test-key', { status: '' })
    expect(f.filter.value.status).toBe('')
  })

  it('apply merges and persists', () => {
    const f = useFilter('test-key', { status: '' })
    f.apply({ status: '2' })
    expect(f.filter.value.status).toBe('2')
    expect(storage.get('test-key')).toContain('"status":"2"')
  })

  it('reset goes back to defaults and clears storage', () => {
    const f = useFilter('test-key', { status: '' })
    f.apply({ status: '2' })
    f.reset()
    expect(f.filter.value.status).toBe('')
    expect(storage.has('test-key')).toBe(false)
  })

  it('restore reads back persisted', () => {
    storage.set('test-key', JSON.stringify({ status: '3' }))
    const f = useFilter('test-key', { status: '' })
    expect(f.filter.value.status).toBe('3')
  })
})
```

- [ ] **Step 6：跑测试验证失败**

```bash
cd eapp && npx vitest run tests/composables/useFilter.spec.ts
```

预期：FAIL — `useFilter` 不存在。

- [ ] **Step 7：实现 useFilter**

```ts
// eapp/composables/useFilter.ts
import { ref } from 'vue'

export function useFilter<T extends Record<string, any>>(storageKey: string, defaults: T) {
  const filter = ref<T>({ ...defaults })

  function readStorage(): Partial<T> | null {
    try {
      const raw = String(uni.getStorageSync(storageKey) || '')
      if (!raw) return null
      return JSON.parse(raw)
    } catch {
      return null
    }
  }

  const persisted = readStorage()
  if (persisted) filter.value = { ...defaults, ...persisted }

  function apply(patch: Partial<T>) {
    filter.value = { ...filter.value, ...patch }
    uni.setStorageSync(storageKey, JSON.stringify(filter.value))
  }
  function reset() {
    filter.value = { ...defaults }
    uni.removeStorageSync(storageKey)
  }

  return { filter, apply, reset }
}
```

- [ ] **Step 8：跑测试验证通过**

```bash
cd eapp && npx vitest run tests/composables/useFilter.spec.ts
```

预期：4 个测试 PASS。

- [ ] **Step 9：实现 useRequest（无独立测试，仅薄封装供页面用）**

```ts
// eapp/composables/useRequest.ts
import { ref } from 'vue'

export function useRequest<TArgs extends any[], TRet>(fn: (...args: TArgs) => Promise<TRet>) {
  const loading = ref(false)
  const error = ref<string>('')
  const data = ref<TRet | null>(null)

  async function run(...args: TArgs): Promise<TRet | null> {
    loading.value = true
    error.value = ''
    try {
      const ret = await fn(...args)
      data.value = ret as any
      return ret
    } catch (e: any) {
      error.value = String(e?.message || e || '请求失败')
      return null
    } finally {
      loading.value = false
    }
  }

  return { loading, error, data, run }
}
```

### Task A5：基础高阶组件骨架（无业务联动，先建文件）

**Files:**
- Create: `eapp/components/biz/PageHeader.vue` / `EmptyState.vue` / `MetricCard.vue` / `ActionGrid.vue` / `TodoCenter.vue` / `AnnouncementBar.vue` / `FilterDrawer.vue` / `BatchBar.vue` / `BatchResultPopup.vue` / `Timeline.vue`
- Modify: `eapp/components/common/StatusTag.vue`、`eapp/components/layout/EappShell.vue`、`eapp/App.vue`

- [ ] **Step 1：扩展 CSS 变量（App.vue）**

修改 `eapp/App.vue` 的 `:root, page` 块新增：

```scss
:root,
page {
  --eapp-primary: #2563eb;
  --eapp-primary-soft: #eff6ff;
  --eapp-success: #16a34a;
  --eapp-success-soft: #dcfce7;
  --eapp-warning: #f59e0b;
  --eapp-warning-soft: #fef3c7;
  --eapp-danger: #dc2626;
  --eapp-danger-soft: #fee2e2;
  --eapp-accent: #f97316;
  --eapp-accent-soft: #ffedd5;
  --eapp-bg: #f8fafc;
  --eapp-card: #ffffff;
  --eapp-text: #1e293b;
  --eapp-text-muted: #64748b;
  --eapp-text-faint: #94a3b8;
  --eapp-border: #e2e8f0;
  --eapp-border-strong: #cbd5e1;
}
```

- [ ] **Step 2：扩展 StatusTag 色板**

修改 `eapp/components/common/StatusTag.vue` 的 `styleVars` 计算属性：

```ts
const styleVars = computed(() => {
  const t = String(props.type || '')
  if (t.includes('pending') || t === '1' || t === 'applied') return { bg: 'var(--eapp-warning-soft)', color: '#92400e' }
  if (t.includes('ship') || t === '2' || t === '3' || t.includes('return')) return { bg: 'var(--eapp-primary-soft)', color: 'var(--eapp-primary)' }
  if (t.includes('complete') || t === '4' || t.includes('success') || t.includes('refunded')) return { bg: 'var(--eapp-success-soft)', color: '#166534' }
  if (t.includes('close') || t.includes('reject') || t === '5') return { bg: 'var(--eapp-danger-soft)', color: '#991b1b' }
  if (t.includes('enabled')) return { bg: 'var(--eapp-success-soft)', color: '#166534' }
  if (t.includes('disabled')) return { bg: '#e2e8f0', color: '#475569' }
  if (t.includes('warning') || t.includes('warn')) return { bg: 'var(--eapp-warning-soft)', color: '#92400e' }
  if (t.includes('hot') || t.includes('accent')) return { bg: 'var(--eapp-accent-soft)', color: '#c2410c' }
  return { bg: '#e2e8f0', color: '#334155' }
})
```

- [ ] **Step 3：EappShell 增强（支持顶部安全区与白底头部）**

修改 `eapp/components/layout/EappShell.vue`：

```vue
<script setup lang="ts">
defineProps<{ padded?: boolean; headerSticky?: boolean }>()
</script>
<template>
  <view class="eapp-shell">
    <view :class="['eapp-shell__header', headerSticky ? 'is-sticky' : '']"><slot name="header" /></view>
    <view :class="['eapp-shell__body', padded === false ? '' : 'is-padded']"><slot /></view>
  </view>
</template>
<style scoped>
.eapp-shell { min-height: 100vh; background: var(--eapp-bg); }
.eapp-shell__header { padding-top: env(safe-area-inset-top); background: var(--eapp-card); }
.eapp-shell__header.is-sticky { position: sticky; top: 0; z-index: 10; }
.eapp-shell__body.is-padded { padding: 20rpx; box-sizing: border-box; }
</style>
```

- [ ] **Step 4：写 PageHeader.vue**

```vue
<!-- eapp/components/biz/PageHeader.vue -->
<script setup lang="ts">
defineProps<{ title?: string; subtitle?: string; transparent?: boolean }>()
</script>
<template>
  <view :class="['page-header', transparent ? 'is-transparent' : '']">
    <view class="left">
      <view class="title">{{ title }}</view>
      <view v-if="subtitle" class="subtitle">{{ subtitle }}</view>
    </view>
    <view class="right"><slot name="right" /></view>
  </view>
</template>
<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; padding: 20rpx 24rpx; padding-top: calc(20rpx + env(safe-area-inset-top)); background: var(--eapp-card); }
.page-header.is-transparent { background: transparent; }
.title { font-size: 36rpx; font-weight: 700; color: var(--eapp-text); }
.subtitle { margin-top: 4rpx; font-size: 22rpx; color: var(--eapp-text-muted); }
.right { display: flex; align-items: center; gap: 16rpx; }
</style>
```

- [ ] **Step 5：写 EmptyState.vue**

```vue
<!-- eapp/components/biz/EmptyState.vue -->
<script setup lang="ts">
defineProps<{ title?: string; desc?: string; icon?: string }>()
defineEmits<{ (e: 'action'): void }>()
</script>
<template>
  <view class="empty">
    <view class="icon">{{ icon || '✦' }}</view>
    <view class="title">{{ title || '暂无数据' }}</view>
    <view v-if="desc" class="desc">{{ desc }}</view>
    <view v-if="$slots.default" class="action"><slot /></view>
  </view>
</template>
<style scoped>
.empty { padding: 80rpx 0; text-align: center; }
.icon { font-size: 60rpx; color: var(--eapp-border-strong); }
.title { margin-top: 16rpx; color: var(--eapp-text); font-size: 28rpx; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.action { margin-top: 24rpx; display: flex; justify-content: center; }
</style>
```

- [ ] **Step 6：写 MetricCard.vue**

```vue
<!-- eapp/components/biz/MetricCard.vue -->
<script setup lang="ts">
import { computed } from 'vue'
const props = defineProps<{ title: string; value: string | number; unit?: string; compare?: number; color?: string }>()
const compareText = computed(() => {
  if (props.compare === undefined || props.compare === null) return ''
  const v = Number(props.compare)
  const sign = v >= 0 ? '+' : ''
  return `${sign}${(v * 100).toFixed(1)}%`
})
const compareTone = computed(() => {
  if (props.compare === undefined || props.compare === null) return ''
  return Number(props.compare) >= 0 ? 'up' : 'down'
})
defineEmits<{ (e: 'click'): void }>()
</script>
<template>
  <view class="metric-card" @click="$emit('click')">
    <view class="title-row">
      <text class="title">{{ title }}</text>
      <text v-if="compareText" :class="['compare', compareTone]">{{ compareText }}</text>
    </view>
    <view class="value-row">
      <text class="value" :style="{ color: color || 'var(--eapp-text)' }">{{ value }}</text>
      <text v-if="unit" class="unit">{{ unit }}</text>
    </view>
  </view>
</template>
<style scoped>
.metric-card { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 24rpx; padding: 24rpx; }
.title-row { display: flex; align-items: center; justify-content: space-between; }
.title { font-size: 24rpx; color: var(--eapp-text-muted); }
.compare { font-size: 20rpx; padding: 4rpx 10rpx; border-radius: 10rpx; }
.compare.up { background: var(--eapp-success-soft); color: #166534; }
.compare.down { background: var(--eapp-danger-soft); color: #991b1b; }
.value-row { margin-top: 12rpx; display: flex; align-items: baseline; gap: 6rpx; }
.value { font-size: 44rpx; font-weight: 700; }
.unit { font-size: 22rpx; color: var(--eapp-text-muted); }
</style>
```

- [ ] **Step 7：写 ActionGrid.vue**

```vue
<!-- eapp/components/biz/ActionGrid.vue -->
<script setup lang="ts">
defineProps<{ items: Array<{ key: string; label: string; icon?: string; color?: string; badge?: number; soon?: boolean }>; columns?: number }>()
defineEmits<{ (e: 'click', key: string): void }>()
</script>
<template>
  <view class="action-grid" :style="{ gridTemplateColumns: `repeat(${columns || 4}, 1fr)` }">
    <view v-for="item in items" :key="item.key" class="cell" @click="$emit('click', item.key)">
      <view class="icon" :style="{ background: item.color || 'var(--eapp-primary-soft)' }">{{ item.icon || '◆' }}</view>
      <text class="label">{{ item.label }}</text>
      <text v-if="item.badge" class="badge">{{ item.badge > 99 ? '99+' : item.badge }}</text>
      <text v-if="item.soon" class="soon">即将上线</text>
    </view>
  </view>
</template>
<style scoped>
.action-grid { display: grid; gap: 16rpx; padding: 20rpx; background: var(--eapp-card); border-radius: 24rpx; border: 1px solid var(--eapp-border); }
.cell { display: flex; flex-direction: column; align-items: center; gap: 8rpx; padding: 14rpx 8rpx; position: relative; }
.icon { width: 72rpx; height: 72rpx; border-radius: 22rpx; display: flex; align-items: center; justify-content: center; font-size: 32rpx; color: var(--eapp-primary); }
.label { font-size: 22rpx; color: var(--eapp-text); }
.badge { position: absolute; top: 8rpx; right: 20rpx; min-width: 32rpx; height: 32rpx; padding: 0 8rpx; background: var(--eapp-danger); color: #fff; font-size: 18rpx; border-radius: 16rpx; display: flex; align-items: center; justify-content: center; }
.soon { position: absolute; top: 8rpx; right: 8rpx; padding: 2rpx 8rpx; background: #fef3c7; color: #92400e; font-size: 18rpx; border-radius: 6rpx; }
</style>
```

- [ ] **Step 8：写 TodoCenter.vue**

```vue
<!-- eapp/components/biz/TodoCenter.vue -->
<script setup lang="ts">
defineProps<{ items: Array<{ key: string; title: string; value: number; tone?: 'warn'|'primary'|'danger'|'normal' }> }>()
defineEmits<{ (e: 'click', key: string): void }>()
</script>
<template>
  <view class="todo-center">
    <view v-for="item in items" :key="item.key" class="row" @click="$emit('click', item.key)">
      <text class="title">{{ item.title }}</text>
      <text :class="['value', item.tone || 'normal']">{{ item.value }}</text>
    </view>
  </view>
</template>
<style scoped>
.todo-center { display: grid; grid-template-columns: 1fr 1fr; gap: 12rpx; background: var(--eapp-card); border-radius: 24rpx; padding: 16rpx; border: 1px solid var(--eapp-border); }
.row { display: flex; align-items: center; justify-content: space-between; padding: 18rpx 20rpx; border-radius: 16rpx; background: var(--eapp-bg); }
.title { font-size: 24rpx; color: var(--eapp-text); }
.value { font-size: 32rpx; font-weight: 700; }
.value.warn { color: var(--eapp-warning); }
.value.primary { color: var(--eapp-primary); }
.value.danger { color: var(--eapp-danger); }
.value.normal { color: var(--eapp-text); }
</style>
```

- [ ] **Step 9：写 AnnouncementBar.vue**

```vue
<!-- eapp/components/biz/AnnouncementBar.vue -->
<script setup lang="ts">
defineProps<{ items: Array<{ id: number; title: string; type?: string }> }>()
defineEmits<{ (e: 'click', id: number): void }>()
</script>
<template>
  <scroll-view v-if="items.length" scroll-x class="ann-bar">
    <view class="ann-inner">
      <view v-for="item in items" :key="item.id" class="ann-item" @click="$emit('click', item.id)">
        <text class="badge">{{ item.type === 'urgent' ? '紧急' : '公告' }}</text>
        <text class="title">{{ item.title }}</text>
      </view>
    </view>
  </scroll-view>
</template>
<style scoped>
.ann-bar { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 16rpx; padding: 14rpx 0; white-space: nowrap; }
.ann-inner { display: inline-flex; gap: 24rpx; padding: 0 20rpx; }
.ann-item { display: inline-flex; align-items: center; gap: 10rpx; }
.badge { padding: 4rpx 10rpx; background: var(--eapp-warning-soft); color: #92400e; border-radius: 6rpx; font-size: 20rpx; }
.title { font-size: 24rpx; color: var(--eapp-text); }
</style>
```

- [ ] **Step 10：写 FilterDrawer.vue**

```vue
<!-- eapp/components/biz/FilterDrawer.vue -->
<script setup lang="ts">
defineProps<{ show: boolean; title?: string }>()
defineEmits<{ (e: 'close'): void; (e: 'reset'): void; (e: 'confirm'): void }>()
</script>
<template>
  <up-popup :show="show" mode="right" round="0" @close="$emit('close')">
    <view class="drawer">
      <view class="drawer-header">
        <text class="drawer-title">{{ title || '筛选' }}</text>
        <text class="drawer-close" @click="$emit('close')">关闭</text>
      </view>
      <scroll-view scroll-y class="drawer-body"><slot /></scroll-view>
      <view class="drawer-footer">
        <up-button plain @click="$emit('reset')">重置</up-button>
        <up-button type="primary" @click="$emit('confirm')">确认</up-button>
      </view>
    </view>
  </up-popup>
</template>
<style scoped>
.drawer { width: 600rpx; height: 100vh; display: flex; flex-direction: column; background: var(--eapp-bg); }
.drawer-header { display: flex; align-items: center; justify-content: space-between; padding: 24rpx; padding-top: calc(24rpx + env(safe-area-inset-top)); background: var(--eapp-card); }
.drawer-title { font-size: 30rpx; font-weight: 700; }
.drawer-close { color: var(--eapp-text-muted); font-size: 26rpx; }
.drawer-body { flex: 1; padding: 20rpx; }
.drawer-footer { padding: 20rpx; background: var(--eapp-card); display: grid; grid-template-columns: 1fr 1fr; gap: 16rpx; padding-bottom: calc(20rpx + env(safe-area-inset-bottom)); }
</style>
```

- [ ] **Step 11：写 BatchBar.vue**

```vue
<!-- eapp/components/biz/BatchBar.vue -->
<script setup lang="ts">
defineProps<{ count: number; actions: Array<{ key: string; label: string; tone?: 'primary'|'warning'|'danger' }> }>()
defineEmits<{ (e: 'action', key: string): void; (e: 'cancel'): void }>()
</script>
<template>
  <view v-if="count > 0" class="batch-bar">
    <view class="left">
      <text>已选 {{ count }} 项</text>
      <text class="cancel" @click="$emit('cancel')">取消</text>
    </view>
    <view class="actions">
      <up-button
        v-for="act in actions"
        :key="act.key"
        size="mini"
        :type="act.tone || 'primary'"
        plain
        @click="$emit('action', act.key)"
      >{{ act.label }}</up-button>
    </view>
  </view>
</template>
<style scoped>
.batch-bar { position: fixed; left: 0; right: 0; bottom: 0; background: var(--eapp-card); border-top: 1px solid var(--eapp-border); padding: 16rpx 20rpx; padding-bottom: calc(16rpx + env(safe-area-inset-bottom)); display: flex; align-items: center; justify-content: space-between; gap: 16rpx; z-index: 30; }
.left { display: flex; align-items: center; gap: 16rpx; }
.cancel { color: var(--eapp-text-muted); font-size: 24rpx; }
.actions { display: flex; gap: 10rpx; flex-wrap: wrap; justify-content: flex-end; }
</style>
```

- [ ] **Step 12：写 BatchResultPopup.vue**

```vue
<!-- eapp/components/biz/BatchResultPopup.vue -->
<script setup lang="ts">
defineProps<{ show: boolean; success: number[]; fails: Array<{ id: number; reason: string }> }>()
defineEmits<{ (e: 'close'): void }>()
</script>
<template>
  <up-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="result">
      <view class="result-title">批量操作结果</view>
      <view class="row">成功：<text class="ok">{{ success.length }}</text>　失败：<text class="bad">{{ fails.length }}</text></view>
      <scroll-view v-if="fails.length" scroll-y class="fail-list">
        <view v-for="f in fails" :key="f.id" class="fail-row">#{{ f.id }} — {{ f.reason }}</view>
      </scroll-view>
      <up-button type="primary" @click="$emit('close')">知道了</up-button>
    </view>
  </up-popup>
</template>
<style scoped>
.result { padding: 24rpx; }
.result-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.row { font-size: 26rpx; color: var(--eapp-text); margin-bottom: 14rpx; }
.ok { color: var(--eapp-success); font-weight: 700; }
.bad { color: var(--eapp-danger); font-weight: 700; }
.fail-list { max-height: 400rpx; margin-bottom: 16rpx; }
.fail-row { padding: 10rpx 0; color: var(--eapp-text-muted); font-size: 23rpx; border-bottom: 1px dashed var(--eapp-border); }
</style>
```

- [ ] **Step 13：写 Timeline.vue**

```vue
<!-- eapp/components/biz/Timeline.vue -->
<script setup lang="ts">
defineProps<{ items: Array<{ key: string | number; title: string; time?: string; desc?: string; tone?: 'primary'|'success'|'warn'|'muted' }>; compact?: boolean }>()
</script>
<template>
  <view class="timeline">
    <view v-for="(item, idx) in items" :key="item.key" :class="['tl-item', compact ? 'compact' : '']">
      <view class="dot-col">
        <view :class="['dot', item.tone || 'muted']" />
        <view v-if="idx < items.length - 1" class="bar" />
      </view>
      <view class="content">
        <view class="title">{{ item.title }}</view>
        <view v-if="item.desc" class="desc">{{ item.desc }}</view>
        <view v-if="item.time" class="time">{{ item.time }}</view>
      </view>
    </view>
  </view>
</template>
<style scoped>
.timeline { padding: 12rpx 4rpx; }
.tl-item { display: flex; gap: 16rpx; }
.tl-item.compact .content { padding-bottom: 14rpx; }
.dot-col { display: flex; flex-direction: column; align-items: center; }
.dot { width: 16rpx; height: 16rpx; border-radius: 50%; margin-top: 12rpx; }
.dot.primary { background: var(--eapp-primary); }
.dot.success { background: var(--eapp-success); }
.dot.warn { background: var(--eapp-warning); }
.dot.muted { background: var(--eapp-border-strong); }
.bar { flex: 1; width: 2rpx; background: var(--eapp-border); margin: 8rpx 0; }
.content { flex: 1; padding-bottom: 24rpx; }
.title { font-size: 26rpx; color: var(--eapp-text); }
.desc { margin-top: 4rpx; font-size: 23rpx; color: var(--eapp-text-muted); }
.time { margin-top: 6rpx; font-size: 22rpx; color: var(--eapp-text-faint); }
</style>
```

- [ ] **Step 14：跑全部测试确保仍绿**

```bash
cd eapp && npx vitest run
```

预期：之前已有 + 新加的 utils/ly-charts + composables 测试全部 PASS。

- [ ] **Step 15：提交 commit 1**

```bash
git -C 'D:\Repos\xyito\open\lyshop' add eapp package.json eapp/package-lock.json
git -C 'D:\Repos\xyito\open\lyshop' commit -m "$(cat <<'EOF'
eapp: 引入 ly-charts 与商家高阶组件骨架

- 通过 npm 与 uni_modules 双通道集成 ly-charts（qiun-data-charts）
- 新增 components/charts/（ChartPanel + AreaChart + RingChart + BarChart）封装统一空态/加载态/主题色
- 新增 components/biz/ 基础组件：PageHeader / EmptyState / MetricCard / ActionGrid / TodoCenter / AnnouncementBar / FilterDrawer / BatchBar / BatchResultPopup / Timeline
- 扩展 StatusTag 色板与 App.vue 全局 CSS 变量；EappShell 增加 safe-area 与 sticky 支持
- composables 沉淀 useBatchSelection / useFilter / useRequest，带 vitest 测试
EOF
)"
```

---

## Phase B — 重写工作台 dashboard 与待办中心（→ commit 2）

### Task B1：mock /dashboard 升级 + /shops/current + /announcements

**Files:**
- Modify: `admin/src/mock/index.ts`
- Test: `eapp/tests/mock/dashboard.spec.ts`

- [ ] **Step 1：写 dashboard mock 测试**

```ts
// eapp/tests/mock/dashboard.spec.ts
import { describe, expect, it } from 'vitest'
import { matchMock } from '../../../admin/src/mock/index'

describe('mock /dashboard upgraded', () => {
  it('returns trend/status_distribution/hot_products/announcements/stock_warning_list', () => {
    const r = matchMock('GET', '/admin/api/dashboard', {})
    expect(r.matched).toBe(true)
    const d = r.data
    expect(d.trend.revenue_7d.categories).toHaveLength(7)
    expect(d.trend.revenue_7d.series[0].data).toHaveLength(7)
    expect(d.trend.revenue_30d.categories).toHaveLength(30)
    expect(d.status_distribution.length).toBeGreaterThanOrEqual(5)
    expect(d.hot_products.length).toBe(5)
    expect(d.announcements.length).toBeGreaterThan(0)
    expect(d.compare.revenue_yoy).toBeTypeOf('number')
    expect(typeof d.today_avg_price).toBe('number')
  })

  it('returns stable trend across calls (deterministic)', () => {
    const a = matchMock('GET', '/admin/api/dashboard', {}).data
    const b = matchMock('GET', '/admin/api/dashboard', {}).data
    expect(a.trend.revenue_7d.series[0].data).toEqual(b.trend.revenue_7d.series[0].data)
  })
})

describe('mock /shops/current', () => {
  it('returns current shop', () => {
    const r = matchMock('GET', '/admin/api/shops/current', {})
    expect(r.matched).toBe(true)
    expect(r.data.id).toBeGreaterThan(0)
    expect(r.data.name).toBeTruthy()
  })
})

describe('mock /announcements', () => {
  it('returns list', () => {
    const r = matchMock('GET', '/admin/api/announcements', {})
    expect(r.matched).toBe(true)
    expect(Array.isArray(r.data.list)).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(3)
  })
})
```

- [ ] **Step 2：跑测试验证失败**

```bash
cd eapp && npx vitest run tests/mock/dashboard.spec.ts
```

预期：FAIL（trend / shops / announcements 字段缺失）。

- [ ] **Step 3：在 admin/src/mock/index.ts 顶部 source 区追加 dashboard helpers + shops + announcements**

在 `function clone<T>` 前后插入：

```ts
function seedRandom(seed: number) {
  const x = Math.sin(seed) * 10000
  return x - Math.floor(x)
}

function buildRevenueTrend(days: number): { categories: string[]; series: Array<{ name: string; data: number[] }> } {
  const categories: string[] = []
  const data: number[] = []
  const today = new Date('2026-05-28T00:00:00Z')
  for (let i = days - 1; i >= 0; i--) {
    const d = new Date(today.getTime() - i * 86400000)
    const mm = String(d.getUTCMonth() + 1).padStart(2, '0')
    const dd = String(d.getUTCDate()).padStart(2, '0')
    categories.push(`${mm}-${dd}`)
    const seed = d.getUTCFullYear() * 10000 + (d.getUTCMonth() + 1) * 100 + d.getUTCDate()
    const dayOfWeek = d.getUTCDay()
    const base = 4800 + Math.sin(seed / 7) * 1200
    const weekendBoost = (dayOfWeek === 0 || dayOfWeek === 6) ? 600 : 0
    const noise = (seedRandom(seed) - 0.5) * 800
    data.push(Math.max(800, Math.round(base + weekendBoost + noise)))
  }
  return { categories, series: [{ name: '营收', data }] }
}

function buildOrderTrend(days: number): { categories: string[]; series: Array<{ name: string; data: number[] }> } {
  const r = buildRevenueTrend(days)
  return { categories: r.categories, series: [{ name: '订单', data: r.series[0].data.map((v) => Math.max(2, Math.round(v / 120))) }] }
}

const shopsCurrentSource: any = {
  id: 1,
  name: '示范品牌旗舰店',
  logo: 'https://picsum.photos/200/200?random=shop1',
  owner: 'admin',
  decor_status: 'published',
}

const announcementsSource: any[] = [
  { id: 1, title: '平台 2026 年 6 月例行升级通知', content: '6 月 3 日 02:00-03:00 短暂维护', type: 'normal', created_at: '2026-05-27T10:00:00Z' },
  { id: 2, title: '618 大促招商进行中', content: '前往营销中心报名参与', type: 'urgent', created_at: '2026-05-26T09:00:00Z' },
  { id: 3, title: '电子面单服务费下调', content: '6 月起部分快递公司面单费下调 10%', type: 'normal', created_at: '2026-05-25T14:00:00Z' },
]
```

- [ ] **Step 4：在 matchMock 中替换/补充 dashboard 分支**

定位现有 `key === 'GET /admin/api/dashboard'` 处（如已存在则改写；若无则新增）。新增/改写为：

```ts
if (key === 'GET /admin/api/dashboard') {
  const trendRevenue7d = buildRevenueTrend(7)
  const trendRevenue30d = buildRevenueTrend(30)
  const trendOrder7d = buildOrderTrend(7)
  const todayRevenue = trendRevenue7d.series[0].data[trendRevenue7d.series[0].data.length - 1]
  const todayOrders = trendOrder7d.series[0].data[trendOrder7d.series[0].data.length - 1]
  const yesterdayRevenue = trendRevenue7d.series[0].data[trendRevenue7d.series[0].data.length - 2]
  const yesterdayOrders = trendOrder7d.series[0].data[trendOrder7d.series[0].data.length - 2]
  const orderList = clone(orderListSource)
  const statusDistribution = [
    { name: '待付款', value: orderList.filter((o: any) => String(o.status) === '1').length },
    { name: '待发货', value: orderList.filter((o: any) => String(o.status) === '2').length },
    { name: '已发货', value: orderList.filter((o: any) => String(o.status) === '3').length },
    { name: '已完成', value: orderList.filter((o: any) => String(o.status) === '4').length },
    { name: '已关闭', value: orderList.filter((o: any) => String(o.status) === '5').length },
    { name: '售后中', value: orderList.filter((o: any) => o.has_after_sale).length },
  ]
  const hotProducts = clone(productListSource).slice(0, 5).map((p: any, i: number) => ({
    id: p.id, title: p.title, cover: p.cover || `https://picsum.photos/120/120?random=${200 + i}`,
    sold_qty: 80 + i * 17,
  }))
  return {
    matched: true,
    data: {
      today_orders: todayOrders,
      today_sales: todayRevenue,
      today_avg_price: todayOrders > 0 ? Math.round((todayRevenue / todayOrders) * 100) / 100 : 0,
      pending_ship: orderList.filter((o: any) => String(o.status) === '2').length,
      pending_after_sale: 3,
      unread_message: 6,
      stock_warning: 4,
      compare: {
        revenue_yoy: 0.182,
        revenue_mom: yesterdayRevenue > 0 ? (todayRevenue - yesterdayRevenue) / yesterdayRevenue : 0,
        order_yoy: 0.094,
        order_mom: yesterdayOrders > 0 ? (todayOrders - yesterdayOrders) / yesterdayOrders : 0,
      },
      trend: { revenue_7d: trendRevenue7d, revenue_30d: trendRevenue30d, order_7d: trendOrder7d },
      status_distribution: statusDistribution,
      hot_products: hotProducts,
      announcements: clone(announcementsSource),
      stock_warning_list: [
        { product_id: 101, sku_id: 1011, title: '蓝牙耳机标准版', stock: 8, threshold: 20 },
        { product_id: 102, sku_id: 1021, title: '无线键盘', stock: 5, threshold: 30 },
        { product_id: 205, sku_id: 2051, title: '运动手表', stock: 4, threshold: 10 },
        { product_id: 6, sku_id: 10061, title: '羊绒大衣女款', stock: 7, threshold: 15 },
      ],
    },
  }
}
if (key === 'GET /admin/api/shops/current') {
  return { matched: true, data: clone(shopsCurrentSource) }
}
if (key === 'GET /admin/api/announcements') {
  return { matched: true, data: { list: clone(announcementsSource), total: announcementsSource.length, page: 1, size: 20 } }
}
```

- [ ] **Step 5：跑测试验证通过**

```bash
cd eapp && npx vitest run tests/mock/dashboard.spec.ts
```

预期：5 个测试 PASS。

### Task B2：useDashboard composable + api 扩展 + stores/shop

**Files:**
- Modify: `eapp/api/dashboard.ts`、`eapp/api/system.ts`、`eapp/stores/badge.ts`
- Create: `eapp/composables/useDashboard.ts`、`eapp/stores/shop.ts`
- Test: `eapp/tests/composables/useDashboard.spec.ts`、`eapp/tests/stores/shop.spec.ts`

- [ ] **Step 1：扩展 api/dashboard.ts 与 api/system.ts**

```ts
// eapp/api/dashboard.ts
import { get } from '@/utils/request'

export type DashboardData = {
  today_orders: number
  today_sales: number
  today_avg_price: number
  pending_ship: number
  pending_after_sale: number
  unread_message: number
  stock_warning: number
  compare: { revenue_yoy: number; revenue_mom: number; order_yoy: number; order_mom: number }
  trend: {
    revenue_7d: { categories: string[]; series: Array<{ name: string; data: number[] }> }
    revenue_30d: { categories: string[]; series: Array<{ name: string; data: number[] }> }
    order_7d: { categories: string[]; series: Array<{ name: string; data: number[] }> }
  }
  status_distribution: Array<{ name: string; value: number }>
  hot_products: Array<{ id: number; title: string; cover: string; sold_qty: number }>
  announcements: Array<{ id: number; title: string; content: string; type: string; created_at: string }>
  stock_warning_list: Array<{ product_id: number; sku_id: number; title: string; stock: number; threshold: number }>
}

export const getDashboard = () => get<DashboardData>('/dashboard')
```

```ts
// eapp/api/system.ts （在原文件追加）
export const getCurrentShop = () => get<any>('/shops/current')
export const getAnnouncements = () => get<any>('/announcements')
```

- [ ] **Step 2：实现 stores/shop**

```ts
// eapp/stores/shop.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCurrentShop } from '@/api/system'

export const useShopStore = defineStore('eapp-shop', () => {
  const id = ref(0)
  const name = ref('')
  const logo = ref('')

  async function load() {
    const data: any = await getCurrentShop()
    id.value = Number(data?.id || 0)
    name.value = String(data?.name || '')
    logo.value = String(data?.logo || '')
  }

  return { id, name, logo, load }
})
```

- [ ] **Step 3：写 stores/shop 测试**

```ts
// eapp/tests/stores/shop.spec.ts
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

vi.mock('@/api/system', () => ({
  getCurrentShop: vi.fn(async () => ({ id: 1, name: '示范品牌旗舰店', logo: 'http://x.png' })),
}))

import { useShopStore } from '@/stores/shop'

describe('shop store', () => {
  beforeEach(() => setActivePinia(createPinia()))

  it('loads current shop', async () => {
    const s = useShopStore()
    await s.load()
    expect(s.id).toBe(1)
    expect(s.name).toBe('示范品牌旗舰店')
  })
})
```

- [ ] **Step 4：跑 stores/shop 测试**

```bash
cd eapp && npx vitest run tests/stores/shop.spec.ts
```

预期：PASS。

- [ ] **Step 5：实现 useDashboard composable**

```ts
// eapp/composables/useDashboard.ts
import { ref } from 'vue'
import { getDashboard, type DashboardData } from '@/api/dashboard'

const empty: DashboardData = {
  today_orders: 0, today_sales: 0, today_avg_price: 0,
  pending_ship: 0, pending_after_sale: 0, unread_message: 0, stock_warning: 0,
  compare: { revenue_yoy: 0, revenue_mom: 0, order_yoy: 0, order_mom: 0 },
  trend: {
    revenue_7d: { categories: [], series: [] },
    revenue_30d: { categories: [], series: [] },
    order_7d: { categories: [], series: [] },
  },
  status_distribution: [],
  hot_products: [],
  announcements: [],
  stock_warning_list: [],
}

export function useDashboard() {
  const loading = ref(false)
  const data = ref<DashboardData>({ ...empty })

  async function load() {
    loading.value = true
    try {
      const ret = await getDashboard()
      if (ret) data.value = ret
    } catch {
      data.value = { ...empty }
    } finally {
      loading.value = false
    }
  }

  return { loading, data, load }
}
```

- [ ] **Step 6：写 useDashboard 测试**

```ts
// eapp/tests/composables/useDashboard.spec.ts
import { describe, expect, it, vi } from 'vitest'

vi.mock('@/api/dashboard', () => ({
  getDashboard: vi.fn(async () => ({
    today_orders: 5, today_sales: 800, today_avg_price: 160,
    pending_ship: 1, pending_after_sale: 0, unread_message: 2, stock_warning: 0,
    compare: { revenue_yoy: 0.1, revenue_mom: 0.05, order_yoy: 0.2, order_mom: 0.1 },
    trend: { revenue_7d: { categories: [], series: [] }, revenue_30d: { categories: [], series: [] }, order_7d: { categories: [], series: [] } },
    status_distribution: [], hot_products: [], announcements: [], stock_warning_list: [],
  })),
}))

import { useDashboard } from '@/composables/useDashboard'

describe('useDashboard', () => {
  it('loads data and sets loading flag', async () => {
    const h = useDashboard()
    expect(h.loading.value).toBe(false)
    const p = h.load()
    expect(h.loading.value).toBe(true)
    await p
    expect(h.loading.value).toBe(false)
    expect(h.data.value.today_orders).toBe(5)
  })
})
```

- [ ] **Step 7：跑测试**

```bash
cd eapp && npx vitest run tests/composables/useDashboard.spec.ts
```

预期：PASS。

- [ ] **Step 8：更新 badge store 接入新字段**

修改 `eapp/stores/badge.ts` 的 `syncFromDashboard`：

```ts
function syncFromDashboard(data: any) {
  orderBadge.value = Number(data?.pending_ship || 0) + Number(data?.pending_after_sale || 0)
  messageBadge.value = Number(data?.unread_message || 0)
  syncTabBarBadge(1, orderBadge.value)
  syncTabBarBadge(4, messageBadge.value)
}
```

（结构保持兼容，新字段已被纳入聚合。）

### Task B3：dashboard/index.vue 重写

**Files:**
- Modify: `eapp/pages/dashboard/index.vue`

- [ ] **Step 1：重写 dashboard 页面**

```vue
<!-- eapp/pages/dashboard/index.vue -->
<script setup lang="ts">
import { onLoad, onShow } from '@dcloudio/uni-app'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import PageHeader from '@/components/biz/PageHeader.vue'
import MetricCard from '@/components/biz/MetricCard.vue'
import ActionGrid from '@/components/biz/ActionGrid.vue'
import TodoCenter from '@/components/biz/TodoCenter.vue'
import AnnouncementBar from '@/components/biz/AnnouncementBar.vue'
import AreaChart from '@/components/charts/AreaChart.vue'
import RingChart from '@/components/charts/RingChart.vue'
import BarChart from '@/components/charts/BarChart.vue'
import { useDashboard } from '@/composables/useDashboard'
import { useShopStore } from '@/stores/shop'
import { useBadgeStore } from '@/stores/badge'
import { setStorage } from '@/utils/storage'

const { t } = useI18n()
const dash = useDashboard()
const shop = useShopStore()
const badge = useBadgeStore()
const trendRange = ref<'7d' | '30d'>('7d')

const trendData = computed(() => trendRange.value === '7d' ? dash.data.value.trend.revenue_7d : dash.data.value.trend.revenue_30d)

const todoItems = computed(() => [
  { key: 'pending_ship', title: '待发货', value: dash.data.value.pending_ship, tone: 'warn' as const },
  { key: 'pending_after_sale', title: '待审售后', value: dash.data.value.pending_after_sale, tone: 'danger' as const },
  { key: 'unread_message', title: '待回消息', value: dash.data.value.unread_message, tone: 'primary' as const },
  { key: 'stock_warning', title: '库存预警', value: dash.data.value.stock_warning, tone: 'warn' as const },
  { key: 'pending_invoice', title: '待开发票', value: 2, tone: 'normal' as const },
  { key: 'pending_refund', title: '待处理退款', value: 1, tone: 'danger' as const },
])

const actions = [
  { key: 'scan_ship', label: '扫码发货', icon: '⌨', color: 'var(--eapp-primary-soft)' },
  { key: 'product_new', label: '新建商品', icon: '＋', color: 'var(--eapp-success-soft)' },
  { key: 'coupon', label: '优惠券', icon: '券', color: 'var(--eapp-warning-soft)' },
  { key: 'decor', label: '店铺装修', icon: '装', color: 'var(--eapp-accent-soft)', soon: true },
  { key: 'analytics', label: '数据分析', icon: '图', color: 'var(--eapp-primary-soft)', soon: true },
  { key: 'customer', label: '客户管理', icon: '客', color: 'var(--eapp-success-soft)', soon: true },
  { key: 'finance', label: '资金财务', icon: '财', color: 'var(--eapp-warning-soft)', soon: true },
  { key: 'wms', label: '仓储管理', icon: '仓', color: 'var(--eapp-accent-soft)', soon: true },
]

const hotBar = computed(() => ({
  categories: dash.data.value.hot_products.map((p) => p.title.slice(0, 6)),
  series: [{ name: '销量', data: dash.data.value.hot_products.map((p) => p.sold_qty) }],
}))

async function loadAll() {
  await Promise.all([dash.load(), shop.load()])
  badge.syncFromDashboard(dash.data.value)
}

function onTodo(key: string) {
  if (key === 'pending_ship') {
    setStorage('eapp_order_status_filter', '2')
    uni.switchTab({ url: '/pages/order/list' })
  } else if (key === 'pending_after_sale') {
    uni.navigateTo({ url: '/pages/order/after-sale-list' })
  } else if (key === 'unread_message') {
    uni.switchTab({ url: '/pages/me/index' })
  } else if (key === 'stock_warning') {
    setStorage('eapp_product_status_filter', 'warning')
    uni.switchTab({ url: '/pages/product/list' })
  } else {
    uni.showToast({ title: '该功能即将上线', icon: 'none' })
  }
}

function onAction(key: string) {
  if (key === 'scan_ship') {
    uni.scanCode({
      success: (res) => {
        setStorage('eapp_batch_ship_seed', String(res.result || ''))
        uni.navigateTo({ url: '/pages/order/batch?mode=ship' })
      },
      fail: () => uni.showToast({ title: '扫码已取消', icon: 'none' }),
    })
  } else if (key === 'product_new') {
    uni.navigateTo({ url: '/pages/product/edit?id=0' })
  } else if (key === 'coupon') {
    uni.navigateTo({ url: '/pages/marketing/coupon' })
  } else {
    uni.showToast({ title: '该功能即将上线，请到管理后台操作', icon: 'none' })
  }
}

onLoad(loadAll)
onShow(loadAll)
</script>

<template>
  <view class="page">
    <PageHeader :title="shop.name || '工作台'" subtitle="当前店铺">
      <template #right>
        <text class="icon-btn" @click="onAction('scan_ship')">⌨</text>
        <text class="icon-btn" @click="uni.switchTab({ url: '/pages/me/index' })">★</text>
      </template>
    </PageHeader>

    <view class="body">
      <view class="metric-grid">
        <MetricCard :title="t('dashboard.todaySales')" :value="`¥${Number(dash.data.value.today_sales).toFixed(2)}`" :compare="dash.data.value.compare.revenue_mom" color="var(--eapp-primary)" />
        <MetricCard :title="t('dashboard.todayOrders')" :value="String(dash.data.value.today_orders)" :compare="dash.data.value.compare.order_mom" />
        <MetricCard title="客单价" :value="`¥${Number(dash.data.value.today_avg_price).toFixed(2)}`" />
      </view>

      <AreaChart :title="`营收趋势（近 ${trendRange === '7d' ? 7 : 30} 日）`" :data="trendData" :loading="dash.loading.value" :height="380">
        <template #extra>
          <text :class="['tab', trendRange === '7d' ? 'active' : '']" @click="trendRange = '7d'">7 日</text>
          <text :class="['tab', trendRange === '30d' ? 'active' : '']" @click="trendRange = '30d'">30 日</text>
        </template>
      </AreaChart>

      <RingChart title="订单状态分布" :data="dash.data.value.status_distribution" :loading="dash.loading.value" :height="380" />

      <AnnouncementBar :items="dash.data.value.announcements" @click="(id) => uni.showModal({ title: '公告', content: dash.data.value.announcements.find((x) => x.id === id)?.content || '', showCancel: false })" />

      <view class="section-title">待办</view>
      <TodoCenter :items="todoItems" @click="onTodo" />

      <view class="section-title">快捷入口</view>
      <ActionGrid :items="actions" :columns="4" @click="onAction" />

      <BarChart title="商品销量排行 Top5" :data="hotBar" :loading="dash.loading.value" :height="380" />

      <view v-if="dash.data.value.stock_warning_list.length" class="section-title">库存预警</view>
      <view class="warning-list">
        <view v-for="row in dash.data.value.stock_warning_list" :key="row.sku_id" class="warning-row" @click="onTodo('stock_warning')">
          <text class="warn-title">{{ row.title }}</text>
          <text class="warn-stock">{{ row.stock }} / {{ row.threshold }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.body { padding: 20rpx; display: grid; gap: 20rpx; }
.metric-grid { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16rpx; }
.icon-btn { font-size: 36rpx; padding: 0 12rpx; }
.tab { font-size: 22rpx; padding: 6rpx 14rpx; border-radius: 20rpx; background: var(--eapp-bg); color: var(--eapp-text-muted); }
.tab.active { background: var(--eapp-primary); color: #fff; }
.section-title { margin-top: 4rpx; font-size: 30rpx; font-weight: 700; }
.warning-list { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 16rpx; display: grid; gap: 12rpx; }
.warning-row { display: flex; justify-content: space-between; padding: 12rpx 16rpx; border-radius: 16rpx; background: var(--eapp-bg); }
.warn-title { font-size: 26rpx; }
.warn-stock { color: var(--eapp-warning); font-weight: 700; }
</style>
```

- [ ] **Step 2：在 locales 中追加少量缺失键（其余等 Phase H 一次性补）**

修改 `eapp/locales/zh-CN.ts` 的 `dashboard` 段（如键已存在则跳过）：

```ts
dashboard: {
  title: '工作台',
  todayOrders: '今日订单',
  todaySales: '今日营收',
  pendingShip: '待发货',
  pendingAfterSale: '待审售后',
  unreadMessage: '未读消息',
  stockWarning: '库存预警',
},
```

`en.ts` 对应追加。

- [ ] **Step 3：dev 起服务并人工验证**

```bash
cd eapp && npm run dev:h5 -- --mode demo
```

人工验证：
- 浏览器打开 dev server URL（一般 http://localhost:5173）
- dashboard 首屏：3 指标卡 / 7 日趋势图 / 环形 / 公告 / 待办 / 九宫格 / 销量 Top5 全部渲染
- 点击 7/30 切换；点击待办项按预期跳转
- 控制台无 ly-charts 报错

- [ ] **Step 4：跑全部测试**

```bash
cd eapp && npx vitest run
```

预期：全部 PASS。

- [ ] **Step 5：commit 2**

```bash
git -C 'D:\Repos\xyito\open\lyshop' add eapp admin/src/mock/index.ts
git -C 'D:\Repos\xyito\open\lyshop' commit -m "$(cat <<'EOF'
eapp: 重写工作台 dashboard 与待办中心

- mock /dashboard 升级返回 trend / status_distribution / hot_products / announcements / stock_warning_list / compare / today_avg_price，30 日趋势按日期 seed 生成保证刷新稳定
- 新增 mock /shops/current 与 /announcements 接口
- api/dashboard.ts 增 DashboardData 类型与扩展字段；api/system.ts 增 getCurrentShop / getAnnouncements
- 新增 stores/shop 与 composables/useDashboard
- dashboard 页面接入 3 指标卡（含同环比）+ 营收趋势图（7/30 切换）+ 订单状态环形 + 公告带 + 6 待办 + 8 快捷入口 + 销量 Top5 + 库存预警列表
EOF
)"
```

---

## 计划文件较长，剩余 Phase C-H 详见下方

**Part 2（订单 / 售后 / 商品）**：`2026-05-28-eapp-merchant-overhaul-p1-orders-product.md`
- Phase C：commit 3 — 订单列表与详情升级，新增批量与改价流程
- Phase D：commit 4 — 售后协商时间线与凭证上传
- Phase E：commit 5 — 商品列表与编辑（多规格 SKU + 分类）

**Part 3（Mock 补齐 / 文档 / i18n）**：`2026-05-28-eapp-merchant-overhaul-p1-finalize.md`
- Phase F：commit 6 — mock 路由与示例数据增量补齐
- Phase G：commit 7 — docs-site 同步
- Phase H：commit 8 — 国际化补全与最终单测

执行顺序：本文件（A+B）→ Part 2（C+D+E）→ Part 3（F+G+H）。三份文件中所有 Task 必须按编号顺序执行；每个 commit 应当让仓库处于「编译通过 + vitest 通过」的状态。
