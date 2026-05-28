# eapp 商家移动端 P1 实施计划 — 订单 / 售后 / 商品（Part 2）

> 本文件接续 `2026-05-28-eapp-merchant-overhaul-p1.md`（已含 Phase A + B）。任务 ID、commit 顺序与总 spec 一致。

## Phase C — 订单列表与详情升级，新增批量与改价流程（→ commit 3）

### Task C1：订单 mock 扩展（查询参数 + 批量 + repricing + notes + timeline + print + remind-pay）

**Files:** 修改 `admin/src/mock/index.ts`；新增 `eapp/tests/mock/orders.spec.ts`

- [ ] **Step 1：写订单 mock 测试**

```ts
// eapp/tests/mock/orders.spec.ts
import { describe, expect, it } from 'vitest'
import { matchMock } from '../../../admin/src/mock/index'

describe('mock /orders extended query', () => {
  it('filters by status', () => {
    const r = matchMock('GET', '/admin/api/orders', { page: 1, size: 50, status: '2' })
    expect(r.matched).toBe(true)
    for (const o of r.data.list) expect(String(o.status)).toBe('2')
  })
  it('filters by amount range', () => {
    const r = matchMock('GET', '/admin/api/orders', { page: 1, size: 50, amount_min: 100, amount_max: 10000 })
    expect(r.matched).toBe(true)
    for (const o of r.data.list) {
      const amt = Number(o.pay_amount || o.total_amount || 0)
      expect(amt).toBeGreaterThanOrEqual(100); expect(amt).toBeLessThanOrEqual(10000)
    }
  })
})

describe('mock order action routes', () => {
  it('POST /orders/{id}/repricing returns updated breakdown', () => {
    const r = matchMock('POST', '/admin/api/orders/1/repricing', { items: [{ item_id: 1, price: 100 }], remark: 't' })
    expect(r.matched).toBe(true); expect(r.data.amount_breakdown).toBeTruthy()
  })
  it('POST /orders/{id}/notes pushes a note', () => {
    const r = matchMock('POST', '/admin/api/orders/1/notes', { content: 'hello' })
    expect(r.matched).toBe(true); expect(r.data.notes.length).toBeGreaterThan(0)
  })
  it('POST /orders/{id}/remind-pay returns sent_at', () => {
    const r = matchMock('POST', '/admin/api/orders/1/remind-pay', { channel: 'sms' })
    expect(r.matched).toBe(true); expect(r.data.sent_at).toBeTruthy()
  })
  it('GET /orders/{id}/print-template returns html', () => {
    const r = matchMock('GET', '/admin/api/orders/1/print-template', {})
    expect(r.matched).toBe(true); expect(typeof r.data.template).toBe('string')
  })
  it('GET /orders/{id}/timeline returns stages', () => {
    const r = matchMock('GET', '/admin/api/orders/1/timeline', {})
    expect(r.matched).toBe(true); expect(r.data.length).toBeGreaterThanOrEqual(3)
  })
})

describe('mock batch order routes', () => {
  it('POST /orders/batch/ship returns success and fail', () => {
    const r = matchMock('POST', '/admin/api/orders/batch/ship', [
      { order_id: 1, company: 'SF', tracking_no: 'SF1' },
      { order_id: 99999, company: 'SF', tracking_no: 'SF3' },
    ])
    expect(r.matched).toBe(true)
    expect(r.data.success_ids.length).toBeGreaterThan(0)
    expect(r.data.fail.length).toBeGreaterThan(0)
  })
  it('POST /orders/batch/notes works', () => {
    const r = matchMock('POST', '/admin/api/orders/batch/notes', { ids: [1, 2], content: '批量备注' })
    expect(r.matched).toBe(true)
  })
  it('POST /orders/batch/close honors reason', () => {
    const r = matchMock('POST', '/admin/api/orders/batch/close', { ids: [1, 2], reason: '不要了' })
    expect(r.matched).toBe(true)
  })
})
```

- [ ] **Step 2：跑测试验证失败**

```bash
cd eapp && npx vitest run tests/mock/orders.spec.ts
```

预期：FAIL（路由缺失）。

- [ ] **Step 3：扩展 `GET /admin/api/orders` 查询逻辑**

定位 matchMock 中 `'GET /admin/api/orders'` 分支（若已存在则改写；若无则在订单详情分支前插入）：

```ts
if (key === 'GET /admin/api/orders') {
  const keyword = String(query.keyword || '').trim().toLowerCase()
  const status = String(query.status || '')
  const amountMin = Number(query.amount_min || 0)
  const amountMax = Number(query.amount_max || 0)
  const company = String(query.logistics_company || '')
  const province = String(query.province || '')
  const payMethod = String(query.pay_method || '')
  const hasAfterSale = query.has_after_sale === true || String(query.has_after_sale || '') === 'true' || query.has_after_sale === 1
  const timeStart = String(query.time_start || '')
  const timeEnd = String(query.time_end || '')

  let list = clone(orderListSource)
  if (status) list = list.filter((o: any) => String(o.status) === status)
  if (keyword) {
    list = list.filter((o: any) =>
      String(o.id || '').includes(keyword) ||
      String(o.user_nickname || '').toLowerCase().includes(keyword) ||
      String(o.receiver_name || '').toLowerCase().includes(keyword) ||
      (Array.isArray(o.items) && o.items.some((it: any) => String(it.title || '').toLowerCase().includes(keyword))))
  }
  if (amountMin > 0) list = list.filter((o: any) => Number(o.pay_amount || o.total_amount || 0) >= amountMin)
  if (amountMax > 0) list = list.filter((o: any) => Number(o.pay_amount || o.total_amount || 0) <= amountMax)
  if (company) list = list.filter((o: any) => Array.isArray(o.shipments) && o.shipments.some((s: any) => String(s.company || '') === company))
  if (province) list = list.filter((o: any) => String(o.receiver_address || '').includes(province))
  if (payMethod) list = list.filter((o: any) => String(o.pay_method || '') === payMethod)
  if (hasAfterSale) list = list.filter((o: any) => o.has_after_sale === true || (Array.isArray(o.after_sales) && o.after_sales.length > 0))
  if (timeStart) list = list.filter((o: any) => String(o.created_at || '') >= timeStart)
  if (timeEnd) list = list.filter((o: any) => String(o.created_at || '') <= timeEnd)

  return { matched: true, data: toPageData(list, Number(query.page || 1), Number(query.size || 20)) }
}
```

- [ ] **Step 4：在订单分支之后新增 action / batch 路由（共 9 条）**

```ts
if (methodUpper === 'POST' && /\/admin\/api\/orders\/\d+\/repricing$/.test(url)) {
  const id = Number(url.split('/')[4] || 0)
  const target = orderListSource.find((o: any) => Number(o.id || 0) === id)
  if (!target) return { matched: true, data: null }
  const items: any[] = Array.isArray(params?.items) ? params.items : []
  let goodsTotal = 0
  for (const it of items) {
    const oi = (target.items || []).find((x: any) => Number(x.id || x.item_id || 0) === Number(it.item_id || 0))
    if (oi) { oi.price = Number(it.price || 0); goodsTotal += Number(it.price || 0) * Number(oi.qty || 1) }
  }
  if (!goodsTotal) goodsTotal = Number(target.goods_amount || target.total_amount || 0)
  target.amount_breakdown = {
    goods_amount: goodsTotal,
    discount_amount: Number(target.discount_amount || 0),
    payable_amount: Math.max(0, goodsTotal - Number(target.discount_amount || 0)),
  }
  target.pay_amount = target.amount_breakdown.payable_amount
  target.notes = target.notes || []
  target.notes.push({ id: Date.now(), content: `改价：${params?.remark || ''}`, created_at: new Date().toISOString() })
  return { matched: true, data: { id, amount_breakdown: clone(target.amount_breakdown) } }
}

if (methodUpper === 'POST' && /\/admin\/api\/orders\/\d+\/notes$/.test(url)) {
  const id = Number(url.split('/')[4] || 0)
  const target = orderListSource.find((o: any) => Number(o.id || 0) === id)
  if (!target) return { matched: true, data: null }
  target.notes = target.notes || []
  target.notes.push({ id: Date.now(), content: String(params?.content || ''), created_at: new Date().toISOString(), visible_to: String(params?.visible_to || 'merchant_only') })
  return { matched: true, data: { id, notes: clone(target.notes) } }
}

if (methodUpper === 'POST' && /\/admin\/api\/orders\/\d+\/remind-pay$/.test(url)) {
  return { matched: true, data: { sent_at: new Date().toISOString(), channel: String(params?.channel || 'sms') } }
}

if (methodUpper === 'GET' && /\/admin\/api\/orders\/\d+\/print-template$/.test(url)) {
  const id = Number(url.split('/')[4] || 0)
  const target = orderListSource.find((o: any) => Number(o.id || 0) === id)
  const name = target?.receiver_name || '--'
  const addr = target?.receiver_address || '--'
  const tpl = `<div style="font-family:sans-serif;padding:12px;"><h3>电子面单 #${id}</h3><div>收件人：${name}</div><div>地址：${addr}</div><div>商品：${(target?.items || []).map((x: any) => x.title).join('，')}</div><div style="margin-top:10px;font-weight:bold;">[ 顺丰速运 ] SF${1000000 + id}</div></div>`
  return { matched: true, data: { template: tpl } }
}

if (methodUpper === 'GET' && /\/admin\/api\/orders\/\d+\/timeline$/.test(url)) {
  const id = Number(url.split('/')[4] || 0)
  const target = orderListSource.find((o: any) => Number(o.id || 0) === id)
  const createdAt = String(target?.created_at || '2026-05-25T08:00:00Z')
  const stages = [
    { stage: 'created', status: '已下单', time: createdAt, content: '买家下单' },
    { stage: 'paid', status: '已支付', time: createdAt.replace('08:00:00', '08:05:00'), content: '微信支付完成' },
    { stage: 'shipped', status: '已发货', time: createdAt.replace('08:00:00', '12:00:00'), content: '顺丰已揽收' },
    { stage: 'received', status: '已签收', time: createdAt.replace('08:00:00', '23:00:00'), content: '客户已签收' },
    { stage: 'completed', status: '已完成', time: createdAt.replace('08:00:00', '23:30:00'), content: '订单完成' },
  ]
  const currentStatus = String(target?.status || '1')
  const cutoff: Record<string, number> = { '1': 1, '2': 2, '3': 3, '4': 5, '5': 1 }
  return { matched: true, data: stages.slice(0, cutoff[currentStatus] || stages.length) }
}

if (methodUpper === 'POST' && url === '/admin/api/orders/batch/ship') {
  const rows: any[] = Array.isArray(params) ? params : []
  const success_ids: number[] = []
  const fail: Array<{ id: number; reason: string }> = []
  for (const row of rows) {
    const oid = Number(row?.order_id || 0)
    const target = orderListSource.find((o: any) => Number(o.id || 0) === oid)
    if (!target) { fail.push({ id: oid, reason: '订单不存在' }); continue }
    if (String(target.status) === '5') { fail.push({ id: oid, reason: '订单已关闭' }); continue }
    target.shipments = target.shipments || []
    target.shipments.push({ id: Date.now() + oid, company: String(row.company || 'SF'), tracking_no: String(row.tracking_no || ''), delivery_type: 'express', logistics_status: 'created', logistics_status_label: '已下单', created_at: new Date().toISOString() })
    target.status = '3'; target.status_label = '已发货'
    success_ids.push(oid)
  }
  return { matched: true, data: { success_ids, fail } }
}

if (methodUpper === 'POST' && url === '/admin/api/orders/batch/notes') {
  const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
  const content = String(params?.content || '')
  const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
  for (const oid of ids) {
    const t = orderListSource.find((o: any) => Number(o.id || 0) === Number(oid))
    if (!t) { fail.push({ id: Number(oid), reason: '订单不存在' }); continue }
    t.notes = t.notes || []; t.notes.push({ id: Date.now(), content, created_at: new Date().toISOString() })
    success_ids.push(Number(oid))
  }
  return { matched: true, data: { success_ids, fail } }
}

if (methodUpper === 'POST' && url === '/admin/api/orders/batch/repricing') {
  const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
  const adj = params?.adjustment || { type: 'percent', value: -10 }
  const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
  for (const oid of ids) {
    const t = orderListSource.find((o: any) => Number(o.id || 0) === Number(oid))
    if (!t) { fail.push({ id: Number(oid), reason: '订单不存在' }); continue }
    if (String(t.status) !== '1') { fail.push({ id: Number(oid), reason: '当前状态不可改价' }); continue }
    const orig = Number(t.pay_amount || t.total_amount || 0)
    const next = adj.type === 'percent' ? orig * (1 + Number(adj.value || 0) / 100) : orig + Number(adj.value || 0)
    t.pay_amount = Math.max(0, Math.round(next * 100) / 100)
    t.amount_breakdown = { goods_amount: t.pay_amount, discount_amount: 0, payable_amount: t.pay_amount }
    success_ids.push(Number(oid))
  }
  return { matched: true, data: { success_ids, fail } }
}

if (methodUpper === 'POST' && url === '/admin/api/orders/batch/close') {
  const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
  const reason = String(params?.reason || '')
  const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
  for (const oid of ids) {
    const t = orderListSource.find((o: any) => Number(o.id || 0) === Number(oid))
    if (!t) { fail.push({ id: Number(oid), reason: '订单不存在' }); continue }
    if (String(t.status) === '4') { fail.push({ id: Number(oid), reason: '已完成订单不能关闭' }); continue }
    t.status = '5'; t.status_label = '已关闭'; t.close_reason = reason
    success_ids.push(Number(oid))
  }
  return { matched: true, data: { success_ids, fail } }
}
```

- [ ] **Step 5：跑测试验证通过**

```bash
cd eapp && npx vitest run tests/mock/orders.spec.ts
```

预期：全部 PASS。

---

### Task C2：api/order 扩展 + useOrderList composable

**Files:** 修改 `eapp/api/order.ts`；新增 `eapp/composables/useOrderList.ts` + `eapp/tests/composables/useOrderList.spec.ts`

- [ ] **Step 1：在 eapp/api/order.ts 末尾追加**

```ts
export const repriceOrder = (id: number | string, payload: { items: Array<{ item_id: number; price: number }>; remark?: string }) =>
  post<any>(`/orders/${id}/repricing`, payload)
export const addOrderNote = (id: number | string, payload: { content: string; visible_to?: string }) =>
  post<any>(`/orders/${id}/notes`, payload)
export const remindPay = (id: number | string, payload: { channel: 'sms' | 'wx' }) =>
  post<any>(`/orders/${id}/remind-pay`, payload)
export const getPrintTemplate = (id: number | string) => get<{ template: string }>(`/orders/${id}/print-template`)
export const getOrderTimeline = (id: number | string) => get<any[]>(`/orders/${id}/timeline`)
export const batchShipOrders = (rows: Array<{ order_id: number; company: string; tracking_no: string }>) =>
  post<{ success_ids: number[]; fail: Array<{ id: number; reason: string }> }>('/orders/batch/ship', rows)
export const batchNoteOrders = (payload: { ids: number[]; content: string }) =>
  post<any>('/orders/batch/notes', payload)
export const batchRepriceOrders = (payload: { ids: number[]; adjustment: { type: 'percent' | 'amount'; value: number } }) =>
  post<any>('/orders/batch/repricing', payload)
export const batchCloseOrders = (payload: { ids: number[]; reason: string }) =>
  post<any>('/orders/batch/close', payload)
```

- [ ] **Step 2：写 useOrderList 测试**

```ts
// eapp/tests/composables/useOrderList.spec.ts
import { describe, expect, it, vi } from 'vitest'

vi.mock('@/api/order', () => ({
  getOrders: vi.fn(async (params: any) => ({
    list: [
      { id: 1, status: params.status || '2', pay_amount: 199 },
      { id: 2, status: params.status || '2', pay_amount: 99 },
    ],
    total: 2, page: params.page || 1, size: params.size || 20,
  })),
}))

import { useOrderList } from '@/composables/useOrderList'

describe('useOrderList', () => {
  it('loads first page with filter', async () => {
    const h = useOrderList()
    h.applyFilter({ status: '2' })
    await h.load()
    expect(h.list.value).toHaveLength(2)
    expect(h.total.value).toBe(2)
  })
  it('refresh resets to page 1', async () => {
    const h = useOrderList(); h.page.value = 5
    await h.refresh()
    expect(h.page.value).toBe(1)
  })
  it('selection helpers', async () => {
    const h = useOrderList(); await h.load()
    h.toggleSelect(1); expect(h.selectedIds.value).toEqual([1])
    h.clearSelect(); expect(h.selectedIds.value).toEqual([])
  })
})
```

- [ ] **Step 3：跑测试验证失败；实现 composable；再跑验证通过**

```bash
cd eapp && npx vitest run tests/composables/useOrderList.spec.ts
```

实现 `eapp/composables/useOrderList.ts`：

```ts
import { ref } from 'vue'
import { getOrders } from '@/api/order'
import { useBatchSelection } from './useBatchSelection'

export type OrderFilter = {
  status?: string; keyword?: string; time_start?: string; time_end?: string
  amount_min?: number; amount_max?: number; logistics_company?: string
  province?: string; pay_method?: string; has_after_sale?: boolean
}

export function useOrderList(initial: OrderFilter = {}) {
  const filter = ref<OrderFilter>({ ...initial })
  const page = ref(1)
  const size = ref(20)
  const total = ref(0)
  const list = ref<any[]>([])
  const loading = ref(false)
  const refreshing = ref(false)
  const sel = useBatchSelection<number>()

  function applyFilter(patch: Partial<OrderFilter>) { filter.value = { ...filter.value, ...patch }; page.value = 1 }
  function resetFilter() { filter.value = { ...initial }; page.value = 1 }

  async function load() {
    loading.value = true
    try {
      const res: any = await getOrders({ ...filter.value, page: page.value, size: size.value })
      list.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
      total.value = Number(res?.total || list.value.length)
    } catch { list.value = []; total.value = 0 } finally { loading.value = false }
  }
  async function refresh() {
    refreshing.value = true; page.value = 1; sel.clear()
    try { await load() } finally { refreshing.value = false }
  }
  async function loadMore() {
    if (loading.value || list.value.length >= total.value) return
    page.value += 1; loading.value = true
    try {
      const res: any = await getOrders({ ...filter.value, page: page.value, size: size.value })
      list.value = list.value.concat(Array.isArray(res?.list) ? res.list : [])
    } finally { loading.value = false }
  }

  return {
    filter, page, size, total, list, loading, refreshing,
    selectedIds: sel.selected, selectCount: sel.count,
    toggleSelect: sel.toggle, isSelected: sel.isSelected, selectAll: sel.selectAll, clearSelect: sel.clear,
    applyFilter, resetFilter, load, refresh, loadMore,
  }
}
```

```bash
cd eapp && npx vitest run tests/composables/useOrderList.spec.ts
```

预期：3 PASS。

---

### Task C3：订单卡片与表单弹层（OrderCard / ShipPopup / RepricingPopup / RemarkPopup）

**Files:** 新增 `eapp/components/biz/OrderCard.vue`、`ShipPopup.vue`、`RepricingPopup.vue`、`RemarkPopup.vue`

- [ ] **Step 1：OrderCard.vue**

```vue
<script setup lang="ts">
import { computed } from 'vue'
import StatusTag from '@/components/common/StatusTag.vue'

const props = defineProps<{ order: any; selectable?: boolean; selected?: boolean }>()
defineEmits<{ (e: 'click'): void; (e: 'toggle'): void; (e: 'action', key: string): void }>()
const thumbs = computed(() => (Array.isArray(props.order?.items) ? props.order.items.slice(0, 3) : []))
const remain = computed(() => Math.max(0, (props.order?.items?.length || 0) - 3))
function money(v: any) { return Number(v || 0).toFixed(2) }
</script>

<template>
  <view :class="['order-card', selected ? 'is-selected' : '']" @click="$emit('click')">
    <view class="head">
      <view class="left">
        <view v-if="selectable" class="check" @click.stop="$emit('toggle')"><text>{{ selected ? '☑' : '☐' }}</text></view>
        <text class="no">#{{ order.id }}</text>
      </view>
      <StatusTag :text="order.status_label || order.status || '-'" :type="order.status" />
    </view>
    <view class="items">
      <view v-for="it in thumbs" :key="it.id" class="item">
        <image v-if="it.cover" :src="it.cover" mode="aspectFill" class="cover" />
        <view v-else class="cover placeholder">商品</view>
      </view>
      <view v-if="remain > 0" class="more">+{{ remain }}</view>
    </view>
    <view class="meta">
      <text class="user">{{ order.user_nickname || order.receiver_name || '匿名买家' }}</text>
      <text class="amount">¥{{ money(order.pay_amount || order.total_amount) }}</text>
    </view>
    <view class="actions" @click.stop>
      <up-button size="mini" plain @click="$emit('action', 'detail')">详情</up-button>
      <up-button v-if="String(order.status) === '1'" size="mini" type="warning" plain @click="$emit('action', 'reprice')">改价</up-button>
      <up-button v-if="String(order.status) === '1'" size="mini" type="primary" plain @click="$emit('action', 'remind_pay')">催付</up-button>
      <up-button v-if="String(order.status) === '2'" size="mini" type="primary" plain @click="$emit('action', 'ship')">发货</up-button>
      <up-button size="mini" plain @click="$emit('action', 'note')">备注</up-button>
      <up-button size="mini" plain @click="$emit('action', 'print')">打单</up-button>
    </view>
  </view>
</template>

<style scoped>
.order-card { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; }
.order-card.is-selected { border-color: var(--eapp-primary); background: var(--eapp-primary-soft); }
.head { display: flex; align-items: center; justify-content: space-between; }
.left { display: flex; align-items: center; gap: 12rpx; }
.check { width: 44rpx; height: 44rpx; display: flex; align-items: center; justify-content: center; font-size: 28rpx; }
.no { font-size: 22rpx; color: var(--eapp-text-muted); }
.items { margin-top: 12rpx; display: flex; align-items: center; gap: 8rpx; }
.cover { width: 96rpx; height: 96rpx; border-radius: 12rpx; background: var(--eapp-bg); display: flex; align-items: center; justify-content: center; color: var(--eapp-text-faint); font-size: 22rpx; }
.more { padding: 0 12rpx; height: 96rpx; display: flex; align-items: center; color: var(--eapp-text-muted); font-size: 24rpx; }
.meta { margin-top: 12rpx; display: flex; align-items: center; justify-content: space-between; }
.user { font-size: 24rpx; color: var(--eapp-text-muted); }
.amount { font-size: 32rpx; color: var(--eapp-primary); font-weight: 700; }
.actions { margin-top: 14rpx; display: flex; flex-wrap: wrap; gap: 10rpx; justify-content: flex-end; }
</style>
```

- [ ] **Step 2：ShipPopup.vue（抽出现有 detail.vue 内的发货表单）**

```vue
<script setup lang="ts">
import { reactive, watch } from 'vue'
const props = defineProps<{ show: boolean; deliveryMode?: 'express'|'local'|'both'; loading?: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'submit', payload: any): void }>()

const companyOptions = [
  { code: 'SF', name: '顺丰速运' }, { code: 'ZTO', name: '中通快递' }, { code: 'YTO', name: '圆通速递' },
  { code: 'STO', name: '申通快递' }, { code: 'YD', name: '韵达速递' }, { code: 'JD', name: '京东物流' },
  { code: 'EMS', name: '中国邮政 EMS' }, { code: 'DBL', name: '德邦快递' }, { code: 'JT', name: '极兔速递' },
]
const shipTypeOptions = [{ label: '首次发货', value: 'initial' }, { label: '补发', value: 'reship' }]
const deliveryTypeOptions = [{ label: '快递发货', value: 'express' }, { label: '同城配送', value: 'local' }]

const form = reactive({ ship_type: 'initial', delivery_type: 'express', company: 'SF', tracking_no: '', rider_name: '', rider_phone: '', after_sale_case_id: '', remark: '' })

watch(() => props.show, (v) => {
  if (v) {
    form.ship_type = 'initial'
    form.delivery_type = props.deliveryMode === 'both' ? 'express' : (props.deliveryMode || 'express')
    form.company = 'SF'; form.tracking_no = ''; form.rider_name = ''; form.rider_phone = ''
    form.after_sale_case_id = ''; form.remark = ''
  }
})

function onPick(field: keyof typeof form, options: any[], key: string, e: any) {
  ;(form as any)[field] = String(options[Number(e?.detail?.value || 0)]?.[key] || '')
}
function onSubmit() {
  if (form.ship_type === 'reship' && Number(form.after_sale_case_id || 0) <= 0) { uni.showToast({ title: '补发需填写售后单 ID', icon: 'none' }); return }
  if (form.delivery_type === 'local' && (!form.rider_name.trim() || !form.rider_phone.trim())) { uni.showToast({ title: '请填写骑手信息', icon: 'none' }); return }
  if (form.delivery_type === 'express' && (!form.company || !form.tracking_no.trim())) { uni.showToast({ title: '请填写快递与运单号', icon: 'none' }); return }
  emit('submit', {
    ship_type: form.ship_type, delivery_type: form.delivery_type,
    company: form.delivery_type === 'express' ? form.company : undefined,
    tracking_no: form.delivery_type === 'express' ? form.tracking_no.trim() : undefined,
    rider_name: form.delivery_type === 'local' ? form.rider_name.trim() : undefined,
    rider_phone: form.delivery_type === 'local' ? form.rider_phone.trim() : undefined,
    after_sale_case_id: form.ship_type === 'reship' ? Number(form.after_sale_case_id || 0) : undefined,
    remark: form.remark.trim() || undefined,
  })
}
</script>

<template>
  <up-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="popup">
      <view class="title">订单发货</view>
      <picker mode="selector" :range="shipTypeOptions" range-key="label" @change="(e) => onPick('ship_type', shipTypeOptions, 'value', e)">
        <view class="picker">{{ shipTypeOptions.find((x) => x.value === form.ship_type)?.label }}</view>
      </picker>
      <picker v-if="deliveryMode === 'both'" mode="selector" :range="deliveryTypeOptions" range-key="label" @change="(e) => onPick('delivery_type', deliveryTypeOptions, 'value', e)">
        <view class="picker mt">{{ deliveryTypeOptions.find((x) => x.value === form.delivery_type)?.label }}</view>
      </picker>
      <template v-if="form.delivery_type === 'express'">
        <picker mode="selector" :range="companyOptions" range-key="name" @change="(e) => onPick('company', companyOptions, 'code', e)">
          <view class="picker mt">{{ companyOptions.find((x) => x.code === form.company)?.name }}</view>
        </picker>
        <up-input v-model="form.tracking_no" placeholder="运单号" class="mt" />
      </template>
      <template v-else>
        <up-input v-model="form.rider_name" placeholder="骑手姓名" class="mt" />
        <up-input v-model="form.rider_phone" placeholder="骑手电话" class="mt" />
      </template>
      <up-input v-if="form.ship_type === 'reship'" v-model="form.after_sale_case_id" type="number" placeholder="售后单 ID" class="mt" />
      <up-input v-model="form.remark" placeholder="备注（可选）" class="mt" />
      <up-button type="primary" :loading="loading" class="mt-l" @click="onSubmit">确认发货</up-button>
    </view>
  </up-popup>
</template>

<style scoped>
.popup { padding: 24rpx; }
.title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.picker { min-height: 76rpx; border: 1px solid var(--eapp-border); border-radius: 12rpx; padding: 0 20rpx; display: flex; align-items: center; }
.mt { margin-top: 12rpx; }
.mt-l { margin-top: 20rpx; }
</style>
```

- [ ] **Step 3：RepricingPopup.vue**

```vue
<script setup lang="ts">
import { reactive, watch } from 'vue'
const props = defineProps<{ show: boolean; items: Array<{ id: number; title: string; price: number; qty: number }>; loading?: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'submit', payload: any): void }>()
const form = reactive<any>({ items: [], remark: '' })
watch(() => props.show, (v) => {
  if (v) {
    form.items = (props.items || []).map((it) => ({ item_id: it.id, price: Number(it.price || 0) }))
    form.remark = ''
  }
})
function onSubmit() {
  if (form.items.some((x: any) => Number(x.price) < 0)) { uni.showToast({ title: '价格不能为负', icon: 'none' }); return }
  emit('submit', { items: form.items, remark: form.remark.trim() })
}
</script>

<template>
  <up-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="popup">
      <view class="title">订单改价</view>
      <view v-for="(it, i) in props.items" :key="it.id" class="row">
        <text class="row-title">{{ it.title }}</text>
        <up-input v-model="form.items[i].price" type="digit" class="row-input" />
      </view>
      <up-input v-model="form.remark" placeholder="改价备注（必填）" class="mt" />
      <up-button type="primary" :loading="loading" class="mt-l" @click="onSubmit">保存改价</up-button>
    </view>
  </up-popup>
</template>

<style scoped>
.popup { padding: 24rpx; }
.title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.row { display: flex; gap: 16rpx; align-items: center; padding: 10rpx 0; }
.row-title { flex: 1; font-size: 24rpx; }
.row-input { width: 220rpx; }
.mt { margin-top: 12rpx; }
.mt-l { margin-top: 20rpx; }
</style>
```

- [ ] **Step 4：RemarkPopup.vue**

```vue
<script setup lang="ts">
import { reactive, watch } from 'vue'
const props = defineProps<{ show: boolean; placeholder?: string; loading?: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'submit', content: string): void }>()
const form = reactive({ content: '' })
watch(() => props.show, (v) => { if (v) form.content = '' })
function onSubmit() {
  if (!form.content.trim()) { uni.showToast({ title: '请输入内容', icon: 'none' }); return }
  emit('submit', form.content.trim())
}
</script>

<template>
  <up-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="popup">
      <view class="title">添加备注</view>
      <up-textarea v-model="form.content" :placeholder="placeholder || '请输入备注内容'" />
      <up-button type="primary" :loading="loading" class="mt-l" @click="onSubmit">保存</up-button>
    </view>
  </up-popup>
</template>

<style scoped>
.popup { padding: 24rpx; }
.title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.mt-l { margin-top: 20rpx; }
</style>
```

---

### Task C4：订单列表重写 + 新增 batch.vue / print-preview.vue + pages.json

**Files:** 修改 `eapp/pages/order/list.vue`、`eapp/pages.json`；新增 `eapp/pages/order/batch.vue`、`eapp/pages/order/print-preview.vue`

- [ ] **Step 1：pages.json 注册新页**

在 `pages` 数组 `pages/order/after-sale-detail` 后追加：

```json
{ "path": "pages/order/batch", "style": { "navigationBarTitleText": "批量操作" } },
{ "path": "pages/order/print-preview", "style": { "navigationBarTitleText": "面单预览" } },
```

在 `pages/product/edit` 后追加：

```json
{ "path": "pages/product/category-tree", "style": { "navigationBarTitleText": "分类管理" } },
```

- [ ] **Step 2：重写 pages/order/list.vue 完整代码 — 见 Part 1 plan 的 Task C4 Step 2 块**（与本文件 Phase D 共享）

完整代码见 `2026-05-28-eapp-merchant-overhaul-p1.md` 同 task 的 Step 2 — **此处不重复粘贴以避免上下文冗余**。

**重要提示**：执行此 step 时，请打开主 plan 文件 Task C4 Step 2 复制完整代码。

- [ ] **Step 3：实现 pages/order/batch.vue（完整代码）**

```vue
<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import BatchResultPopup from '@/components/biz/BatchResultPopup.vue'
import { batchShipOrders } from '@/api/order'
import { getStorage, removeStorage } from '@/utils/storage'

const ids = ref<number[]>([])
const rows = reactive<any[]>([])
const company = ref('SF')
const loading = ref(false)
const showResult = ref(false)
const result = ref<{ success_ids: number[]; fail: Array<{ id: number; reason: string }> }>({ success_ids: [], fail: [] })

function syncRows() {
  rows.length = 0
  for (const id of ids.value) rows.push({ order_id: id, company: company.value, tracking_no: '' })
}

function applyCompany(c: string) { company.value = c; for (const r of rows) r.company = c }

async function submit() {
  if (!rows.length) return
  if (rows.some((r) => !String(r.tracking_no).trim())) { uni.showToast({ title: '请填写所有运单号', icon: 'none' }); return }
  loading.value = true
  try {
    const r: any = await batchShipOrders(rows.map((x) => ({ order_id: Number(x.order_id), company: String(x.company || 'SF'), tracking_no: String(x.tracking_no || '').trim() })))
    result.value = r; showResult.value = true
  } finally { loading.value = false }
}

onLoad(() => {
  const raw = String(getStorage('eapp_batch_ship_ids') || '[]')
  try { ids.value = JSON.parse(raw) } catch { ids.value = [] }
  removeStorage('eapp_batch_ship_ids')
  const seed = String(getStorage('eapp_batch_ship_seed') || '')
  if (seed && !ids.value.length) {
    uni.showToast({ title: `已扫码：${seed}`, icon: 'none' })
    removeStorage('eapp_batch_ship_seed')
  }
  syncRows()
})
</script>

<template>
  <view class="page">
    <view class="head">
      <view>选中订单：{{ ids.length }} 单</view>
      <view class="company">
        <text>统一快递：</text>
        <picker mode="selector" :range="['SF','ZTO','YTO','STO','YD','JD','EMS','DBL','JT']" @change="(e) => applyCompany(['SF','ZTO','YTO','STO','YD','JD','EMS','DBL','JT'][Number(e.detail.value)])">
          <view class="picker">{{ company }}</view>
        </picker>
      </view>
    </view>
    <view v-for="row in rows" :key="row.order_id" class="row">
      <text class="no">#{{ row.order_id }}</text>
      <up-input v-model="row.company" placeholder="快递代码" class="col-c" />
      <up-input v-model="row.tracking_no" placeholder="运单号" class="col-t" />
    </view>
    <up-button type="primary" :loading="loading" class="mt" @click="submit">提交批量发货</up-button>
    <BatchResultPopup :show="showResult" :success="result.success_ids" :fails="result.fail" @close="showResult = false; uni.navigateBack()" />
  </view>
</template>

<style scoped>
.page { min-height: 100vh; padding: 20rpx; background: var(--eapp-bg); display: grid; gap: 14rpx; }
.head { display: flex; align-items: center; justify-content: space-between; background: var(--eapp-card); border-radius: 18rpx; padding: 18rpx; border: 1px solid var(--eapp-border); }
.company { display: flex; align-items: center; gap: 10rpx; font-size: 24rpx; }
.picker { border: 1px solid var(--eapp-border); border-radius: 10rpx; padding: 0 16rpx; height: 60rpx; display: flex; align-items: center; }
.row { display: flex; align-items: center; gap: 12rpx; background: var(--eapp-card); padding: 14rpx; border-radius: 18rpx; border: 1px solid var(--eapp-border); }
.no { font-size: 24rpx; color: var(--eapp-text-muted); min-width: 120rpx; }
.col-c { width: 200rpx; }
.col-t { flex: 1; }
.mt { margin-top: 12rpx; }
</style>
```

- [ ] **Step 4：实现 pages/order/print-preview.vue（完整代码）**

```vue
<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { getPrintTemplate } from '@/api/order'

const tpl = ref('')
async function load(id: number) {
  const data: any = await getPrintTemplate(id)
  tpl.value = String(data?.template || '')
}
onLoad((opts) => load(Number(opts?.id || 0)))
function onPrint() { uni.showToast({ title: '已发送至默认打印机', icon: 'success' }) }
</script>

<template>
  <view class="page">
    <view class="card">
      <rich-text v-if="tpl" :nodes="tpl" />
      <view v-else class="empty">加载中…</view>
    </view>
    <up-button type="primary" class="mt" @click="onPrint">发送至打印机</up-button>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; padding: 20rpx; background: var(--eapp-bg); }
.card { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; min-height: 480rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 60rpx 0; }
.mt { margin-top: 20rpx; }
</style>
```

- [ ] **Step 5：完整 list.vue 代码（复制粘贴）**

```vue
<script setup lang="ts">
import { onLoad, onPullDownRefresh, onReachBottom, onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import PageHeader from '@/components/biz/PageHeader.vue'
import OrderCard from '@/components/biz/OrderCard.vue'
import EmptyState from '@/components/biz/EmptyState.vue'
import FilterDrawer from '@/components/biz/FilterDrawer.vue'
import BatchBar from '@/components/biz/BatchBar.vue'
import BatchResultPopup from '@/components/biz/BatchResultPopup.vue'
import RemarkPopup from '@/components/biz/RemarkPopup.vue'
import RepricingPopup from '@/components/biz/RepricingPopup.vue'
import ShipPopup from '@/components/biz/ShipPopup.vue'
import { useOrderList } from '@/composables/useOrderList'
import { addOrderNote, batchCloseOrders, batchNoteOrders, getDeliveryMode, remindPay, repriceOrder, shipOrder } from '@/api/order'
import { getStorage, removeStorage, setStorage } from '@/utils/storage'

const h = useOrderList()
const tabs = [
  { name: '全部', status: '' },
  { name: '待付款', status: '1' },
  { name: '待发货', status: '2' },
  { name: '已发货', status: '3' },
  { name: '已完成', status: '4' },
  { name: '已关闭', status: '5' },
  { name: '售后中', status: '', has_after_sale: true },
  { name: '待评价', status: '4', extra: 'pending_review' },
  { name: '待开票', status: '4', extra: 'pending_invoice' },
]
const current = ref(0)
const showFilter = ref(false); const showShip = ref(false); const showReprice = ref(false); const showNote = ref(false); const showResult = ref(false)
const actionLoading = ref(false)
const activeOrder = ref<any>(null)
const deliveryMode = ref<'express'|'local'|'both'>('express')
const batchMode = ref<'ship'|'notes'|'close'|null>(null)
const result = ref<{ success_ids: number[]; fail: Array<{ id: number; reason: string }> }>({ success_ids: [], fail: [] })
const filterDraft = ref<any>({ keyword: '', amount_min: '', amount_max: '', logistics_company: '', province: '', pay_method: '' })

function onTabChange(idx: number) {
  current.value = idx
  const tab = tabs[idx]
  h.applyFilter({ status: tab.status, has_after_sale: !!tab.has_after_sale })
  h.load()
}

async function loadDelivery() { try { const d: any = await getDeliveryMode(); deliveryMode.value = (d?.mode === 'local' || d?.mode === 'both') ? d.mode : 'express' } catch {} }

function openFilter() { filterDraft.value = { ...h.filter.value }; showFilter.value = true }
function applyFilter() {
  h.applyFilter({
    keyword: filterDraft.value.keyword || undefined,
    amount_min: Number(filterDraft.value.amount_min || 0) || undefined,
    amount_max: Number(filterDraft.value.amount_max || 0) || undefined,
    logistics_company: filterDraft.value.logistics_company || undefined,
    province: filterDraft.value.province || undefined,
    pay_method: filterDraft.value.pay_method || undefined,
  })
  showFilter.value = false; h.load()
}
function resetFilter() { filterDraft.value = { keyword: '', amount_min: '', amount_max: '', logistics_company: '', province: '', pay_method: '' } }

async function onCardAction(order: any, key: string) {
  activeOrder.value = order
  if (key === 'detail') uni.navigateTo({ url: `/pages/order/detail?id=${order.id}` })
  else if (key === 'ship') showShip.value = true
  else if (key === 'reprice') showReprice.value = true
  else if (key === 'note') { batchMode.value = null; showNote.value = true }
  else if (key === 'print') uni.navigateTo({ url: `/pages/order/print-preview?id=${order.id}` })
  else if (key === 'remind_pay') {
    actionLoading.value = true
    try { await remindPay(order.id, { channel: 'sms' }); uni.showToast({ title: '已发送催付', icon: 'success' }) } finally { actionLoading.value = false }
  }
}

async function submitShip(payload: any) {
  if (!activeOrder.value) return
  actionLoading.value = true
  try { await shipOrder(activeOrder.value.id, payload); uni.showToast({ title: '发货成功', icon: 'success' }); showShip.value = false; await h.refresh() } finally { actionLoading.value = false }
}
async function submitReprice(payload: any) {
  if (!activeOrder.value) return
  actionLoading.value = true
  try { await repriceOrder(activeOrder.value.id, payload); uni.showToast({ title: '改价成功', icon: 'success' }); showReprice.value = false; await h.refresh() } finally { actionLoading.value = false }
}
async function submitNote(content: string) {
  if (batchMode.value === 'notes') return onBatchNote(content)
  if (!activeOrder.value) return
  actionLoading.value = true
  try { await addOrderNote(activeOrder.value.id, { content }); uni.showToast({ title: '已添加备注', icon: 'success' }); showNote.value = false } finally { actionLoading.value = false }
}

function startBatch(mode: 'ship'|'notes'|'close') {
  if (!h.selectCount.value) { uni.showToast({ title: '请先勾选订单', icon: 'none' }); return }
  batchMode.value = mode
  if (mode === 'ship') {
    setStorage('eapp_batch_ship_ids', JSON.stringify(h.selectedIds.value))
    uni.navigateTo({ url: '/pages/order/batch?mode=ship' })
  } else if (mode === 'notes') {
    showNote.value = true
  } else if (mode === 'close') {
    uni.showModal({ title: '批量关闭', editable: true, placeholderText: '请输入关闭原因', success: async (res) => {
      if (!res.confirm || !res.content) return
      actionLoading.value = true
      try {
        const r: any = await batchCloseOrders({ ids: h.selectedIds.value, reason: res.content })
        result.value = r; showResult.value = true; h.clearSelect(); await h.refresh()
      } finally { actionLoading.value = false }
    } })
  }
}
async function onBatchNote(content: string) {
  actionLoading.value = true
  try {
    const r: any = await batchNoteOrders({ ids: h.selectedIds.value, content })
    result.value = r; showResult.value = true; showNote.value = false; h.clearSelect(); await h.refresh()
  } finally { actionLoading.value = false }
}

onLoad(async () => {
  const status = String(getStorage('eapp_order_status_filter') || '')
  if (status) {
    removeStorage('eapp_order_status_filter')
    const idx = tabs.findIndex((tab) => tab.status === status)
    if (idx >= 0) current.value = idx
  }
  onTabChange(current.value)
  await loadDelivery()
})
onShow(() => h.load())
onPullDownRefresh(async () => { await h.refresh(); uni.stopPullDownRefresh() })
onReachBottom(() => h.loadMore())
</script>

<template>
  <view class="page">
    <PageHeader title="订单">
      <template #right>
        <text class="icon-btn" @click="openFilter">⚲</text>
      </template>
    </PageHeader>
    <scroll-view scroll-x class="tabs">
      <view class="tabs-inner">
        <text v-for="(tab, idx) in tabs" :key="idx" :class="['tab', current === idx ? 'active' : '']" @click="onTabChange(idx)">{{ tab.name }}</text>
      </view>
    </scroll-view>
    <view class="list">
      <EmptyState v-if="!h.loading.value && !h.list.value.length" title="暂无订单" desc="切换状态或调整筛选条件" />
      <OrderCard
        v-for="o in h.list.value"
        :key="o.id"
        :order="o"
        :selectable="true"
        :selected="h.isSelected(o.id)"
        @toggle="h.toggleSelect(o.id)"
        @click="onCardAction(o, 'detail')"
        @action="(key) => onCardAction(o, key)"
      />
      <view v-if="h.loading.value" class="loading">加载中…</view>
    </view>

    <BatchBar
      :count="h.selectCount.value"
      :actions="[
        { key: 'ship', label: '批量发货', tone: 'primary' },
        { key: 'notes', label: '批量备注', tone: 'primary' },
        { key: 'close', label: '批量关闭', tone: 'danger' },
      ]"
      @action="startBatch"
      @cancel="h.clearSelect()"
    />

    <FilterDrawer :show="showFilter" title="订单筛选" @close="showFilter = false" @reset="resetFilter" @confirm="applyFilter">
      <up-input v-model="filterDraft.keyword" placeholder="订单号 / 买家昵称 / 商品名" />
      <view class="row mt">
        <up-input v-model="filterDraft.amount_min" type="digit" placeholder="最低金额" />
        <text class="dash">~</text>
        <up-input v-model="filterDraft.amount_max" type="digit" placeholder="最高金额" />
      </view>
      <up-input v-model="filterDraft.logistics_company" placeholder="物流公司代码（SF/ZTO/...）" class="mt" />
      <up-input v-model="filterDraft.province" placeholder="收货省份" class="mt" />
      <up-input v-model="filterDraft.pay_method" placeholder="支付方式（wechat/alipay/...）" class="mt" />
    </FilterDrawer>

    <ShipPopup :show="showShip" :delivery-mode="deliveryMode" :loading="actionLoading" @close="showShip = false" @submit="submitShip" />
    <RepricingPopup :show="showReprice" :items="activeOrder?.items || []" :loading="actionLoading" @close="showReprice = false" @submit="submitReprice" />
    <RemarkPopup :show="showNote" :loading="actionLoading" @close="showNote = false" @submit="submitNote" />
    <BatchResultPopup :show="showResult" :success="result.success_ids" :fails="result.fail" @close="showResult = false" />
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.icon-btn { font-size: 36rpx; padding: 0 12rpx; }
.tabs { background: var(--eapp-card); position: sticky; top: 0; z-index: 20; white-space: nowrap; }
.tabs-inner { display: inline-flex; padding: 16rpx 12rpx; gap: 8rpx; }
.tab { padding: 10rpx 20rpx; border-radius: 999rpx; font-size: 24rpx; color: var(--eapp-text-muted); background: var(--eapp-bg); }
.tab.active { background: var(--eapp-primary); color: #fff; }
.list { padding: 20rpx; display: grid; gap: 16rpx; padding-bottom: 200rpx; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
.row { display: flex; align-items: center; gap: 12rpx; }
.dash { color: var(--eapp-text-muted); }
.mt { margin-top: 12rpx; }
</style>
```

### Task C5：order/detail.vue 增强（Timeline + 操作菜单）

**Files:** 修改 `eapp/pages/order/detail.vue`

- [ ] **Step 1：在 `<script setup>` 顶部追加 imports 与 helper**

```ts
import Timeline from '@/components/biz/Timeline.vue'
import { addOrderNote, getOrderTimeline, remindPay, repriceOrder } from '@/api/order'

const timelineItems = ref<Array<{ key: string; title: string; time?: string; tone?: 'primary'|'success'|'warn'|'muted' }>>([])

async function loadTimeline() {
  if (!orderID.value) return
  try {
    const rows: any = await getOrderTimeline(orderID.value)
    timelineItems.value = (Array.isArray(rows) ? rows : []).map((r: any) => ({
      key: r.stage, title: r.status, time: formatDate(r.time),
      tone: r.stage === 'completed' ? 'success' : 'primary',
    }))
  } catch {}
}

function openActionSheet() {
  uni.showActionSheet({
    itemList: ['改价', '添加备注', '打印面单', '催付款'],
    success: async (res) => {
      if (res.tapIndex === 0) {
        if (String(detail.status) !== '1') { uni.showToast({ title: '当前状态不可改价', icon: 'none' }); return }
        uni.showModal({ title: '改价', editable: true, placeholderText: '输入新的支付金额', success: async (m) => {
          if (!m.confirm || !m.content) return
          await repriceOrder(detail.id, { items: detail.items.map((it: any) => ({ item_id: it.id, price: Number(m.content) / (detail.items.length || 1) })), remark: '操作菜单改价' })
          await loadData()
        } })
      } else if (res.tapIndex === 1) {
        uni.showModal({ title: '添加备注', editable: true, success: async (m) => {
          if (!m.confirm || !m.content) return
          await addOrderNote(detail.id, { content: m.content })
          uni.showToast({ title: '已添加', icon: 'success' })
        } })
      } else if (res.tapIndex === 2) {
        uni.navigateTo({ url: `/pages/order/print-preview?id=${detail.id}` })
      } else if (res.tapIndex === 3) {
        await remindPay(detail.id, { channel: 'sms' })
        uni.showToast({ title: '已催付', icon: 'success' })
      }
    },
  })
}
```

并在 `loadData` 末尾追加 `await loadTimeline()`。

- [ ] **Step 2：在 `<template>` 中插入 Timeline 卡片与操作菜单按钮**

紧跟 `<template v-else>` 内的第一个 `<view class="card">` 之前插入：

```vue
<view class="card">
  <view class="section-title">订单进度</view>
  <Timeline :items="timelineItems" />
</view>
```

并在第一个 `<view class="head">` 中（`<text>订单 #{{ detail.id }}</text>` + StatusTag 之后）追加：

```vue
<text class="op" @click="openActionSheet">操作</text>
```

样式追加：

```css
.op { color: var(--eapp-primary); font-size: 24rpx; padding-left: 16rpx; }
```

- [ ] **Step 3：dev 验证 + 跑测试**

```bash
cd eapp && npm run dev:h5 -- --mode demo
# 浏览器验证：订单列表/详情/批量/面单可用
# Ctrl+C 退出
cd eapp && npx vitest run
```

预期：所有测试 PASS；浏览器流程可演示。

- [ ] **Step 4：commit 3**

```bash
git -C 'D:\Repos\xyito\open\lyshop' add eapp admin/src/mock/index.ts
git -C 'D:\Repos\xyito\open\lyshop' commit -m "eapp: 订单列表与详情升级，新增批量与改价流程"
```

commit body：

```
- mock /orders 扩展查询参数（keyword/amount_range/logistics/province/pay_method/has_after_sale/time）
- mock 新增 /orders/{id}/repricing、notes、remind-pay、print-template、timeline 与 4 条 batch 接口；含部分失败示例
- api/order.ts 新增对应封装；composables/useOrderList 沉淀分页/筛选/批量逻辑（含 vitest 测试）
- 新增 OrderCard、ShipPopup、RepricingPopup、RemarkPopup
- order/list 重写：9 状态 tabs、筛选抽屉、勾选批量、操作菜单、空态与下拉/触底
- order/detail 增强：订单进度 Timeline + 操作 ActionSheet（改价/备注/打单/催付）
- 新增 order/batch（批量发货）与 order/print-preview（面单 rich-text）；pages.json 注册新页
```

---

## Phase D — 售后协商时间线与凭证上传（→ commit 4）

### Task D1：mock /after-sales 扩展 + messages + evidences

**Files:** 修改 `admin/src/mock/index.ts`；新增 `eapp/tests/mock/after-sales.spec.ts`

- [ ] **Step 1：写测试**

```ts
// eapp/tests/mock/after-sales.spec.ts
import { describe, expect, it } from 'vitest'
import { matchMock } from '../../../admin/src/mock/index'

describe('mock /after-sales extended', () => {
  it('filters by status and type', () => {
    const r = matchMock('GET', '/admin/api/after-sales', { page: 1, size: 50, status: 'applied' })
    expect(r.matched).toBe(true)
    for (const row of r.data.list) expect(String(row.status)).toBe('applied')
  })
})

describe('mock after-sale messages / evidences', () => {
  it('POST messages appends', () => {
    const r = matchMock('POST', '/admin/api/after-sales/8001/messages', { from: 'merchant', content: 'hello' })
    expect(r.matched).toBe(true)
    expect(r.data.messages.length).toBeGreaterThan(0)
  })
  it('POST evidences appends', () => {
    const r = matchMock('POST', '/admin/api/after-sales/8001/evidences', { images: ['a.jpg'], remark: 'evi' })
    expect(r.matched).toBe(true)
    expect(r.data.evidences.length).toBeGreaterThan(0)
  })
})
```

- [ ] **Step 2：跑测试验证失败；实现 mock**

实现要点（在 matchMock 内合适位置追加；若已有 `GET /admin/api/after-sales` 分支则用以下版本替换）：

```ts
if (key === 'GET /admin/api/after-sales') {
  const status = String(query.status || '')
  const type = String(query.type || '')
  const timeStart = String(query.time_start || '')
  const timeEnd = String(query.time_end || '')
  let list = clone(afterSalesSource)
  if (status) list = list.filter((row: any) => String(row.status) === status)
  if (type) list = list.filter((row: any) => String(row.case_type || '') === type)
  if (timeStart) list = list.filter((row: any) => String(row.created_at || '') >= timeStart)
  if (timeEnd) list = list.filter((row: any) => String(row.created_at || '') <= timeEnd)
  return { matched: true, data: toPageData(list, Number(query.page || 1), Number(query.size || 20)) }
}

if (methodUpper === 'GET' && /\/admin\/api\/after-sales\/\d+$/.test(url)) {
  const id = Number(url.split('/').pop() || 0)
  const target = afterSalesSource.find((row: any) => Number(row.id || 0) === id)
  return { matched: true, data: target ? clone(target) : null }
}

if (methodUpper === 'POST' && /\/admin\/api\/after-sales\/\d+\/messages$/.test(url)) {
  const id = Number(url.split('/')[4] || 0)
  let target = afterSalesSource.find((row: any) => Number(row.id || 0) === id)
  if (!target) {
    target = { id, status: 'applied', case_type: 'refund_only', order_id: 1, messages: [], evidences: [], logs: [], shipments: [], created_at: new Date().toISOString() }
    afterSalesSource.push(target)
  }
  target.messages = target.messages || []
  target.messages.push({ id: Date.now(), from: String(params?.from || 'merchant'), content: String(params?.content || ''), images: Array.isArray(params?.images) ? params.images : [], created_at: new Date().toISOString() })
  return { matched: true, data: { id, messages: clone(target.messages) } }
}

if (methodUpper === 'POST' && /\/admin\/api\/after-sales\/\d+\/evidences$/.test(url)) {
  const id = Number(url.split('/')[4] || 0)
  const target = afterSalesSource.find((row: any) => Number(row.id || 0) === id)
  if (!target) return { matched: true, data: null }
  target.evidences = target.evidences || []
  target.evidences.push({ id: Date.now(), images: Array.isArray(params?.images) ? params.images : [], remark: String(params?.remark || ''), created_at: new Date().toISOString() })
  return { matched: true, data: { id, evidences: clone(target.evidences) } }
}
```

并在 `afterSalesSource` 初始化后追加 3 条示例（含完整 case_type / messages / evidences 字段，确保与现有 spec 4.3.2 一致）：

```ts
afterSalesSource.push(
  { id: 8001, order_id: 1, status: 'applied', status_label: '待审核', case_type: 'refund_only', case_type_label: '仅退款', reason: '少件', refund_amount: 99, created_at: '2026-05-27T10:00:00Z', messages: [], evidences: [], logs: [{ id: 1, action: 'apply', action_label: '申请', from_status: '', to_status: 'applied', to_status_label: '待审核', content: '买家提交申请', created_at: '2026-05-27T10:00:00Z' }], shipments: [] },
  { id: 8002, order_id: 2, status: 'user_returning', status_label: '退货中', case_type: 'return_refund', case_type_label: '退货退款', reason: '不喜欢', refund_amount: 199, created_at: '2026-05-26T15:30:00Z', messages: [
    { id: 1, from: 'user', content: '我已寄回，请查收', images: [], created_at: '2026-05-26T18:00:00Z' },
    { id: 2, from: 'merchant', content: '好的，签收后会立即处理退款', images: [], created_at: '2026-05-26T19:00:00Z' },
  ], evidences: [], logs: [], shipments: [{ id: 1, company: 'SF', tracking_no: 'SF889977', direction: 'inbound', direction_label: '回寄', biz_type: 'return', biz_type_label: '回寄', logistics_status: 'in_transit', logistics_status_label: '运输中', created_at: '2026-05-26T17:30:00Z' }] },
  { id: 8003, order_id: 3, status: 'refund_pending', status_label: '退款中', case_type: 'refund_only', case_type_label: '仅退款', reason: '商品损坏', refund_amount: 459, created_at: '2026-05-25T09:00:00Z', messages: [], evidences: [{ id: 1, images: ['https://picsum.photos/200/200?random=evi1'], remark: '商品破损照片', created_at: '2026-05-25T09:10:00Z' }], logs: [], shipments: [] },
)
```

跑测试：

```bash
cd eapp && npx vitest run tests/mock/after-sales.spec.ts
```

预期：PASS。

### Task D2：api/after-sale.ts 抽出

**Files:** 新增 `eapp/api/after-sale.ts`；修改 `eapp/api/order.ts`

- [ ] **Step 1：新建 api/after-sale.ts**

```ts
import { get, post } from '@/utils/request'

export const getAfterSales = (params?: any) => get<any>('/after-sales', params)
export const getAfterSaleDetail = (id: number | string) => get<any>(`/after-sales/${id}`)
export const auditAfterSale = (id: number | string, payload: any) => post<any>(`/after-sales/${id}/audit`, payload)
export const receiveAfterSale = (id: number | string) => post<any>(`/after-sales/${id}/receive`)
export const refundAfterSale = (id: number | string, payload: any) => post<any>(`/after-sales/${id}/refund`, payload)
export const completeAfterSale = (id: number | string) => post<any>(`/after-sales/${id}/complete`)
export const closeAfterSale = (id: number | string, payload: any) => post<any>(`/after-sales/${id}/close`, payload)
export const addAfterSaleMessage = (id: number | string, payload: { from: 'merchant'|'user'; content: string; images?: string[] }) =>
  post<any>(`/after-sales/${id}/messages`, payload)
export const addAfterSaleEvidence = (id: number | string, payload: { images: string[]; remark?: string }) =>
  post<any>(`/after-sales/${id}/evidences`, payload)
```

- [ ] **Step 2：order.ts re-export 兼容**

在 `eapp/api/order.ts` 顶部追加：

```ts
export {
  getAfterSales, getAfterSaleDetail, auditAfterSale, receiveAfterSale,
  refundAfterSale, completeAfterSale, closeAfterSale,
} from './after-sale'
```

并删除该文件中重复的 after-sale 函数定义（保留 order 相关导出）。

### Task D3：AfterSaleCard + 列表/详情增强

**Files:** 新增 `eapp/components/biz/AfterSaleCard.vue`；修改 `eapp/pages/order/after-sale-list.vue`、`after-sale-detail.vue`

- [ ] **Step 1：AfterSaleCard.vue**

```vue
<script setup lang="ts">
import StatusTag from '@/components/common/StatusTag.vue'
defineProps<{ row: any }>()
defineEmits<{ (e: 'click'): void }>()
</script>

<template>
  <view class="card" @click="$emit('click')">
    <view class="head">
      <text class="no">售后 #{{ row.id }}</text>
      <StatusTag :text="row.status_label || row.status || '-'" :type="row.status" />
    </view>
    <view class="line">订单：#{{ row.order_id }} · 类型：{{ row.case_type_label || row.case_type || '-' }}</view>
    <view class="line">原因：{{ row.reason || '-' }}</view>
    <view class="line">退款金额：¥{{ Number(row.refund_amount || 0).toFixed(2) }}</view>
  </view>
</template>

<style scoped>
.card { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; }
.head { display: flex; justify-content: space-between; align-items: center; }
.no { font-size: 26rpx; font-weight: 600; }
.line { margin-top: 6rpx; color: var(--eapp-text-muted); font-size: 23rpx; }
</style>
```

- [ ] **Step 2：after-sale-list.vue 全量重写**

```vue
<script setup lang="ts">
import { onLoad, onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import AfterSaleCard from '@/components/biz/AfterSaleCard.vue'
import EmptyState from '@/components/biz/EmptyState.vue'
import { getAfterSales } from '@/api/after-sale'

const tabs = [
  { name: '全部', status: '' },
  { name: '待审核', status: 'applied' },
  { name: '退货中', status: 'user_returning' },
  { name: '退款中', status: 'refund_pending' },
  { name: '已完成', status: 'refunded' },
  { name: '已关闭', status: 'closed' },
]
const typeNames = ['全部', '仅退款', '退货退款', '换货']
const typeValues = ['', 'refund_only', 'return_refund', 'exchange']

const current = ref(0)
const list = ref<any[]>([])
const loading = ref(false)
const filterType = ref('')

async function load() {
  loading.value = true
  try {
    const res: any = await getAfterSales({ page: 1, size: 50, status: tabs[current.value]?.status || undefined, type: filterType.value || undefined })
    list.value = Array.isArray(res?.list) ? res.list : []
  } finally { loading.value = false }
}

function onTab(idx: number) { current.value = idx; load() }
function onType(e: any) { filterType.value = typeValues[Number(e.detail.value)] || ''; load() }
function goDetail(id: number) { uni.navigateTo({ url: `/pages/order/after-sale-detail?id=${id}` }) }

onLoad(load); onShow(load)
</script>

<template>
  <view class="page">
    <scroll-view scroll-x class="tabs">
      <view class="tabs-inner">
        <text v-for="(tab, idx) in tabs" :key="idx" :class="['tab', current === idx ? 'active' : '']" @click="onTab(idx)">{{ tab.name }}</text>
      </view>
    </scroll-view>
    <view class="filter">
      <picker mode="selector" :range="typeNames" @change="onType">
        <view class="picker">类型：{{ typeNames[typeValues.indexOf(filterType)] }}</view>
      </picker>
    </view>
    <view class="list">
      <EmptyState v-if="!loading && !list.length" title="暂无售后" />
      <AfterSaleCard v-for="row in list" :key="row.id" :row="row" @click="goDetail(row.id)" />
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.tabs { background: var(--eapp-card); position: sticky; top: 0; z-index: 20; white-space: nowrap; }
.tabs-inner { display: inline-flex; padding: 16rpx 12rpx; gap: 8rpx; }
.tab { padding: 10rpx 20rpx; border-radius: 999rpx; font-size: 24rpx; color: var(--eapp-text-muted); background: var(--eapp-bg); }
.tab.active { background: var(--eapp-primary); color: #fff; }
.filter { padding: 16rpx 20rpx; }
.picker { display: inline-flex; height: 60rpx; align-items: center; padding: 0 18rpx; border: 1px solid var(--eapp-border); border-radius: 12rpx; background: var(--eapp-card); }
.list { padding: 20rpx; display: grid; gap: 14rpx; }
</style>
```

- [ ] **Step 3：after-sale-detail.vue 在现有结构中追加进度 + 协商 + 凭证**

`<script setup>` 顶部追加：

```ts
import { computed } from 'vue'
import { addAfterSaleEvidence, addAfterSaleMessage } from '@/api/after-sale'

const steps = [
  { key: 'applied', label: '申请' },
  { key: 'approved_wait_user_return', label: '审核' },
  { key: 'user_returning', label: '寄回' },
  { key: 'warehouse_received', label: '收货' },
  { key: 'refunded', label: '退款' },
]
function stepActive(idx: number) {
  const order = ['applied', 'approved_wait_user_return', 'user_returning', 'warehouse_received', 'refund_pending', 'refunded', 'reshipped']
  const cur = order.indexOf(String(detail.status))
  return cur >= idx
}
const messages = computed(() => (Array.isArray(detail.messages) ? detail.messages : []))
const evidences = computed(() => (Array.isArray(detail.evidences) ? detail.evidences : []))
const msgInput = ref('')
const sendingMsg = ref(false)

async function onSendMsg() {
  if (!msgInput.value.trim()) return
  sendingMsg.value = true
  try {
    await addAfterSaleMessage(detail.id, { from: 'merchant', content: msgInput.value.trim() })
    msgInput.value = ''; await loadData()
  } finally { sendingMsg.value = false }
}

function onUploadEvidence() {
  uni.chooseImage({ count: 3, success: async (res) => {
    const images = (res.tempFilePaths || []).map((_, i) => `https://picsum.photos/200/200?random=evi-${Date.now()}-${i}`)
    await addAfterSaleEvidence(detail.id, { images, remark: '商家凭证' })
    await loadData()
  } })
}
```

`<template>` 内：在 `<view v-if="loading">` 之后插入：

```vue
<view v-if="!loading" class="progress">
  <view v-for="(step, idx) in steps" :key="step.key" :class="['step', stepActive(idx) ? 'active' : '']">
    <view class="dot">{{ idx + 1 }}</view>
    <text class="label">{{ step.label }}</text>
  </view>
</view>
```

在「状态日志」card 之后追加：

```vue
<view class="card">
  <view class="section-title">协商沟通</view>
  <view v-if="!messages.length" class="empty-row">暂无沟通记录</view>
  <view v-for="m in messages" :key="m.id" :class="['msg', m.from === 'merchant' ? 'right' : '']">
    <view class="msg-body">
      <text class="who">{{ m.from === 'merchant' ? '商家' : '买家' }}</text>
      <text class="text">{{ m.content }}</text>
      <view v-if="m.images && m.images.length" class="imgs">
        <image v-for="img in m.images" :key="img" :src="img" mode="aspectFill" class="img" />
      </view>
      <text class="time">{{ formatDate(m.created_at) }}</text>
    </view>
  </view>
  <view class="msg-input">
    <up-input v-model="msgInput" placeholder="回复买家" />
    <up-button type="primary" size="mini" :loading="sendingMsg" @click="onSendMsg">发送</up-button>
  </view>
</view>

<view class="card">
  <view class="section-title">
    <text>商家凭证</text>
    <up-button size="mini" plain @click="onUploadEvidence">上传凭证</up-button>
  </view>
  <view v-if="!evidences.length" class="empty-row">暂未上传凭证</view>
  <view v-for="e in evidences" :key="e.id" class="ev">
    <view class="ev-images">
      <image v-for="img in e.images" :key="img" :src="img" mode="aspectFill" class="img" />
    </view>
    <text class="time">{{ formatDate(e.created_at) }} · {{ e.remark || '' }}</text>
  </view>
</view>
```

样式追加：

```css
.progress { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; display: flex; justify-content: space-between; align-items: center; }
.step { display: flex; flex-direction: column; align-items: center; gap: 6rpx; flex: 1; }
.step .dot { width: 44rpx; height: 44rpx; border-radius: 50%; background: var(--eapp-border); color: #fff; display: flex; align-items: center; justify-content: center; }
.step.active .dot { background: var(--eapp-primary); }
.step .label { font-size: 22rpx; color: var(--eapp-text-muted); }
.step.active .label { color: var(--eapp-text); }
.msg { display: flex; margin-top: 10rpx; }
.msg.right { justify-content: flex-end; }
.msg-body { max-width: 70%; background: var(--eapp-bg); padding: 12rpx 14rpx; border-radius: 16rpx; }
.msg.right .msg-body { background: var(--eapp-primary-soft); }
.who { display: block; font-size: 20rpx; color: var(--eapp-text-muted); }
.text { font-size: 24rpx; color: var(--eapp-text); }
.imgs { display: flex; gap: 8rpx; margin-top: 6rpx; }
.img { width: 120rpx; height: 120rpx; border-radius: 10rpx; }
.time { display: block; font-size: 20rpx; color: var(--eapp-text-faint); margin-top: 4rpx; }
.msg-input { margin-top: 12rpx; display: flex; gap: 10rpx; align-items: center; }
.ev { margin-top: 10rpx; padding-top: 10rpx; border-top: 1px dashed var(--eapp-border); }
.ev-images { display: flex; gap: 8rpx; flex-wrap: wrap; }
```

- [ ] **Step 4：跑测试 + commit 4**

```bash
cd eapp && npx vitest run
git -C 'D:\Repos\xyito\open\lyshop' add eapp admin/src/mock/index.ts
git -C 'D:\Repos\xyito\open\lyshop' commit -m "eapp: 售后协商时间线与凭证上传"
```

commit body：

```
- mock /after-sales 扩展按 status/type/time 筛选；新增 /after-sales/{id}/messages 与 /after-sales/{id}/evidences
- afterSalesSource 预置 3 条示例（含 applied / user_returning / refund_pending），含协商记录与凭证
- 抽出 api/after-sale.ts；order.ts 通过 re-export 保持兼容
- 新增 AfterSaleCard；after-sale-list 加状态 tabs 与类型筛选
- after-sale-detail 增加 5 步进度条 + 商家/买家协商时间线（含回复）+ 凭证上传（chooseImage）
```

---

## Phase E — 商品列表与编辑（→ commit 5）

### Task E1：mock /products 扩展查询 + /products/batch/* + /categories CRUD + /categories/tree

**Files:** 修改 `admin/src/mock/index.ts`；新增 `eapp/tests/mock/products.spec.ts`、`eapp/tests/mock/categories.spec.ts`

- [ ] **Step 1：写测试**

```ts
// eapp/tests/mock/products.spec.ts
import { describe, expect, it } from 'vitest'
import { matchMock } from '../../../admin/src/mock/index'

describe('mock /products extended', () => {
  it('filters by status and low_stock', () => {
    const r = matchMock('GET', '/admin/api/products', { page: 1, size: 100, low_stock: true })
    expect(r.matched).toBe(true)
  })
  it('filters by category_id', () => {
    const r = matchMock('GET', '/admin/api/products', { page: 1, size: 100, category_id: 1 })
    expect(r.matched).toBe(true)
  })
  it('sorts by sales', () => {
    const r = matchMock('GET', '/admin/api/products', { page: 1, size: 100, sort_by: 'sales' })
    expect(r.matched).toBe(true)
  })
})

describe('mock /products/batch/*', () => {
  it('PUT /products/batch/status', () => {
    const r = matchMock('PUT', '/admin/api/products/batch/status', { ids: [1, 2], status: 0 })
    expect(r.matched).toBe(true)
    expect(r.data.success_ids.length).toBeGreaterThan(0)
  })
  it('PUT /products/batch/price percent', () => {
    const r = matchMock('PUT', '/admin/api/products/batch/price', { ids: [1, 2], adjustment: { type: 'percent', value: -10 } })
    expect(r.matched).toBe(true)
  })
})
```

```ts
// eapp/tests/mock/categories.spec.ts
import { describe, expect, it } from 'vitest'
import { matchMock } from '../../../admin/src/mock/index'

describe('mock /categories', () => {
  it('GET /categories/tree returns 3 levels', () => {
    const r = matchMock('GET', '/admin/api/categories/tree', {})
    expect(r.matched).toBe(true)
    expect(r.data.length).toBeGreaterThan(0)
    expect(r.data[0].children).toBeTruthy()
  })
  it('POST /categories creates', () => {
    const r = matchMock('POST', '/admin/api/categories', { name: '新分类', parent_id: 0 })
    expect(r.matched).toBe(true)
    expect(r.data.id).toBeGreaterThan(0)
  })
})
```

- [ ] **Step 2：跑测试验证失败；实现 mock**

定位 `categoriesSource` 与 `productListSource`，在 matchMock 内追加：

```ts
// 分类树辅助
function buildCategoryTree(): any[] {
  const list = clone(categoriesSource) as any[]
  const map = new Map<number, any>()
  for (const c of list) { c.children = []; map.set(Number(c.id), c) }
  const roots: any[] = []
  for (const c of list) {
    const parent = Number(c.parent_id || 0)
    if (parent && map.get(parent)) map.get(parent).children.push(c)
    else roots.push(c)
  }
  return roots
}

if (key === 'GET /admin/api/categories/tree') {
  return { matched: true, data: buildCategoryTree() }
}

if (key === 'POST /admin/api/categories') {
  categorySeq += 1
  const row = { id: categorySeq, parent_id: Number(params?.parent_id || 0), name: String(params?.name || `分类${categorySeq}`), sort: Number(params?.sort || 0), product_count: 0 }
  categoriesSource.push(row)
  return { matched: true, data: clone(row) }
}

if (methodUpper === 'PUT' && /\/admin\/api\/categories\/\d+$/.test(url)) {
  const id = Number(url.split('/').pop() || 0)
  const target = (categoriesSource as any[]).find((row: any) => Number(row.id || 0) === id)
  if (!target) return { matched: true, data: null }
  if (params?.name !== undefined) target.name = String(params.name || target.name)
  if (params?.sort !== undefined) target.sort = Number(params.sort || 0)
  return { matched: true, data: clone(target) }
}

if (methodUpper === 'DELETE' && /\/admin\/api\/categories\/\d+$/.test(url)) {
  const id = Number(url.split('/').pop() || 0)
  const idx = (categoriesSource as any[]).findIndex((row: any) => Number(row.id || 0) === id)
  if (idx >= 0) categoriesSource.splice(idx, 1)
  return { matched: true, data: { id } }
}

if (key === 'GET /admin/api/products') {
  const keyword = String(query.keyword || '').trim().toLowerCase()
  const status = query.status === undefined || query.status === '' ? null : Number(query.status)
  const categoryID = Number(query.category_id || 0)
  const lowStock = query.low_stock === true || String(query.low_stock || '') === 'true' || query.low_stock === 1
  const sortBy = String(query.sort_by || '')

  let list = clone(productListSource) as any[]
  if (keyword) list = list.filter((p: any) => String(p.title || '').toLowerCase().includes(keyword))
  if (status !== null) list = list.filter((p: any) => Number(p.status || 0) === status)
  if (categoryID) list = list.filter((p: any) => Number(p.category_id || 0) === categoryID)
  if (lowStock) list = list.filter((p: any) => Number(p.stock || 0) <= Number(p.low_stock_threshold || 10))
  if (sortBy === 'sales') list.sort((a: any, b: any) => Number(b.sales_count || 0) - Number(a.sales_count || 0))
  else if (sortBy === 'stock') list.sort((a: any, b: any) => Number(a.stock || 0) - Number(b.stock || 0))
  else if (sortBy === 'price_asc') list.sort((a: any, b: any) => Number(a.price || 0) - Number(b.price || 0))
  else if (sortBy === 'price_desc') list.sort((a: any, b: any) => Number(b.price || 0) - Number(a.price || 0))
  else if (sortBy === 'created') list.sort((a: any, b: any) => String(b.created_at || '').localeCompare(String(a.created_at || '')))

  return { matched: true, data: toPageData(list, Number(query.page || 1), Number(query.size || 20)) }
}

if (methodUpper === 'PUT' && url === '/admin/api/products/batch/status') {
  const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
  const status = Number(params?.status || 0) === 1 ? 1 : 0
  const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
  for (const pid of ids) {
    const t = (productListSource as any[]).find((p: any) => Number(p.id || 0) === Number(pid))
    if (!t) { fail.push({ id: Number(pid), reason: '商品不存在' }); continue }
    t.status = status; success_ids.push(Number(pid))
  }
  return { matched: true, data: { success_ids, fail } }
}

if (methodUpper === 'PUT' && url === '/admin/api/products/batch/category') {
  const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
  const cid = Number(params?.category_id || 0)
  const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
  for (const pid of ids) {
    const t = (productListSource as any[]).find((p: any) => Number(p.id || 0) === Number(pid))
    if (!t) { fail.push({ id: Number(pid), reason: '商品不存在' }); continue }
    t.category_id = cid; success_ids.push(Number(pid))
  }
  return { matched: true, data: { success_ids, fail } }
}

if (methodUpper === 'PUT' && url === '/admin/api/products/batch/price') {
  const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
  const adj = params?.adjustment || { type: 'set', value: 0 }
  const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
  for (const pid of ids) {
    const t = (productListSource as any[]).find((p: any) => Number(p.id || 0) === Number(pid))
    if (!t) { fail.push({ id: Number(pid), reason: '商品不存在' }); continue }
    const orig = Number(t.price || 0)
    let next = orig
    if (adj.type === 'set') next = Number(adj.value || 0)
    else if (adj.type === 'percent') next = orig * (1 + Number(adj.value || 0) / 100)
    else if (adj.type === 'amount') next = orig + Number(adj.value || 0)
    t.price = Math.max(0, Math.round(next * 100) / 100); success_ids.push(Number(pid))
  }
  return { matched: true, data: { success_ids, fail } }
}
```

补 `categoriesSource` 三级分类示例（约 20 条）；为 `productListSource` 6 个商品补 `category_id`、`category_path_name`、`sales_count`、`low_stock_threshold` 与 `skus[]`：

```ts
// 在 categoriesSource 初始化后追加（覆盖现有简化版本）
categoriesSource.length = 0
const seedCats = [
  { id: 1, parent_id: 0, name: '电子', sort: 1 },
  { id: 2, parent_id: 0, name: '服饰', sort: 2 },
  { id: 3, parent_id: 0, name: '家居', sort: 3 },
  { id: 4, parent_id: 0, name: '食品', sort: 4 },
  { id: 5, parent_id: 0, name: '美妆', sort: 5 },
  { id: 11, parent_id: 1, name: '手机数码', sort: 1 },
  { id: 12, parent_id: 1, name: '电脑办公', sort: 2 },
  { id: 13, parent_id: 1, name: '智能穿戴', sort: 3 },
  { id: 21, parent_id: 2, name: '女装', sort: 1 },
  { id: 22, parent_id: 2, name: '男装', sort: 2 },
  { id: 23, parent_id: 2, name: '鞋包配饰', sort: 3 },
  { id: 31, parent_id: 3, name: '家具', sort: 1 },
  { id: 32, parent_id: 3, name: '日用百货', sort: 2 },
  { id: 41, parent_id: 4, name: '生鲜', sort: 1 },
  { id: 42, parent_id: 4, name: '零食饮料', sort: 2 },
  { id: 51, parent_id: 5, name: '彩妆', sort: 1 },
  { id: 52, parent_id: 5, name: '护肤', sort: 2 },
  { id: 111, parent_id: 11, name: '智能手机', sort: 1 },
  { id: 112, parent_id: 11, name: '配件', sort: 2 },
  { id: 211, parent_id: 21, name: '连衣裙', sort: 1 },
  { id: 212, parent_id: 21, name: '外套', sort: 2 },
]
for (const c of seedCats) categoriesSource.push({ ...c, product_count: 0 })
categorySeq = Math.max(...categoriesSource.map((c) => Number(c.id || 0)))

// 为前 6 个商品补字段
for (let i = 0; i < Math.min(6, productListSource.length); i++) {
  const p: any = productListSource[i]
  p.category_id = [111, 113, 211, 211, 31, 51][i] || 111
  p.category_path_name = ['电子 / 手机数码 / 智能手机', '电子 / 智能穿戴', '服饰 / 女装 / 连衣裙', '服饰 / 女装 / 连衣裙', '家居 / 家具', '美妆 / 彩妆'][i] || ''
  p.sales_count = 80 + i * 23
  p.low_stock_threshold = 10
  p.skus = [
    { id: p.id * 10 + 1, attrs: [{ name: '颜色', value: '黑色' }, { name: '版本', value: '标准版' }], price: Number(p.price || 0), stock: 50 },
    { id: p.id * 10 + 2, attrs: [{ name: '颜色', value: '白色' }, { name: '版本', value: '标准版' }], price: Number(p.price || 0), stock: 30 },
    { id: p.id * 10 + 3, attrs: [{ name: '颜色', value: '黑色' }, { name: '版本', value: '尊享版' }], price: Math.round(Number(p.price || 0) * 1.2 * 100) / 100, stock: 20 },
    { id: p.id * 10 + 4, attrs: [{ name: '颜色', value: '白色' }, { name: '版本', value: '尊享版' }], price: Math.round(Number(p.price || 0) * 1.2 * 100) / 100, stock: 15 },
  ]
}
```

跑测试：

```bash
cd eapp && npx vitest run tests/mock/products.spec.ts tests/mock/categories.spec.ts
```

预期：PASS。

### Task E2：api/category.ts + api/product 扩展 + useProductList

**Files:** 新增 `eapp/api/category.ts`、`eapp/composables/useProductList.ts`、`eapp/tests/composables/useProductList.spec.ts`；修改 `eapp/api/product.ts`

- [ ] **Step 1：api/category.ts**

```ts
// eapp/api/category.ts
import { del, get, post, put } from '@/utils/request'

export const getCategoriesTree = () => get<any>('/categories/tree')
export const createCategory = (payload: { name: string; parent_id?: number; sort?: number }) => post<any>('/categories', payload)
export const updateCategory = (id: number, payload: { name?: string; sort?: number }) => put<any>(`/categories/${id}`, payload)
export const deleteCategory = (id: number) => del<any>(`/categories/${id}`)
```

- [ ] **Step 2：api/product.ts 追加**

```ts
export const batchUpdateProductStatus = (payload: { ids: number[]; status: 0 | 1 }) => put<any>('/products/batch/status', payload)
export const batchUpdateProductCategory = (payload: { ids: number[]; category_id: number }) => put<any>('/products/batch/category', payload)
export const batchUpdateProductPrice = (payload: { ids: number[]; adjustment: { type: 'set'|'percent'|'amount'; value: number; scope?: 'all'|'main_sku' } }) => put<any>('/products/batch/price', payload)
```

- [ ] **Step 3：useProductList TDD（流程同 useOrderList，省略测试代码以保持简洁；按 useOrderList 模板写）**

```ts
// eapp/composables/useProductList.ts
import { ref } from 'vue'
import { getProducts } from '@/api/product'
import { useBatchSelection } from './useBatchSelection'

export type ProductFilter = { keyword?: string; status?: 0|1|''; category_id?: number; sort_by?: string; low_stock?: boolean }

export function useProductList(initial: ProductFilter = {}) {
  const filter = ref<ProductFilter>({ ...initial })
  const page = ref(1); const size = ref(20); const total = ref(0)
  const list = ref<any[]>([])
  const loading = ref(false); const refreshing = ref(false)
  const sel = useBatchSelection<number>()

  function applyFilter(patch: Partial<ProductFilter>) { filter.value = { ...filter.value, ...patch }; page.value = 1 }
  async function load() {
    loading.value = true
    try {
      const res: any = await getProducts({ ...filter.value, page: page.value, size: size.value })
      list.value = Array.isArray(res?.list) ? res.list : []
      total.value = Number(res?.total || list.value.length)
    } finally { loading.value = false }
  }
  async function refresh() { refreshing.value = true; page.value = 1; sel.clear(); try { await load() } finally { refreshing.value = false } }
  async function loadMore() {
    if (loading.value || list.value.length >= total.value) return
    page.value += 1; loading.value = true
    try {
      const res: any = await getProducts({ ...filter.value, page: page.value, size: size.value })
      list.value = list.value.concat(Array.isArray(res?.list) ? res.list : [])
    } finally { loading.value = false }
  }

  return {
    filter, page, total, list, loading, refreshing,
    selectedIds: sel.selected, selectCount: sel.count,
    toggleSelect: sel.toggle, isSelected: sel.isSelected, selectAll: sel.selectAll, clearSelect: sel.clear,
    applyFilter, load, refresh, loadMore,
  }
}
```

测试（结构与 useOrderList 类似）：

```ts
// eapp/tests/composables/useProductList.spec.ts
import { describe, expect, it, vi } from 'vitest'
vi.mock('@/api/product', () => ({
  getProducts: vi.fn(async (params: any) => ({ list: [{ id: 1, title: 'A', status: params.status ?? 1, stock: 5 }], total: 1, page: 1, size: 20 })),
}))
import { useProductList } from '@/composables/useProductList'
describe('useProductList', () => {
  it('loads with filter', async () => {
    const h = useProductList()
    h.applyFilter({ status: 1 }); await h.load()
    expect(h.list.value[0].id).toBe(1)
  })
})
```

```bash
cd eapp && npx vitest run tests/composables/useProductList.spec.ts
```

预期：PASS。

### Task E3：商品相关组件（ProductCard / SkuMatrixEditor / CategoryTreePicker / RichTextEditor）

**Files:** 新增 4 个组件文件

- [ ] **Step 1：ProductCard.vue**

```vue
<script setup lang="ts">
import StatusTag from '@/components/common/StatusTag.vue'
defineProps<{ product: any; selectable?: boolean; selected?: boolean }>()
defineEmits<{ (e: 'click'): void; (e: 'toggle'): void; (e: 'action', key: string): void }>()
</script>

<template>
  <view :class="['product-card', selected ? 'is-selected' : '']" @click="$emit('click')">
    <view v-if="selectable" class="check" @click.stop="$emit('toggle')">{{ selected ? '☑' : '☐' }}</view>
    <image v-if="product.cover" :src="product.cover" mode="aspectFill" class="cover" />
    <view v-else class="cover placeholder">图</view>
    <view class="body">
      <view class="title-row">
        <text class="title">{{ product.title }}</text>
        <StatusTag :text="Number(product.status || 0) === 1 ? '在售' : '仓库'" :type="Number(product.status || 0) === 1 ? 'enabled' : 'disabled'" />
      </view>
      <view class="meta">¥{{ Number(product.price || 0).toFixed(2) }} · 库存 {{ product.stock || 0 }} · 销量 {{ product.sales_count || 0 }}</view>
      <view class="cat">{{ product.category_path_name || '未分类' }}</view>
      <view class="actions" @click.stop>
        <up-button size="mini" plain @click="$emit('action', 'edit')">编辑</up-button>
        <up-button size="mini" type="warning" plain @click="$emit('action', 'toggle_sale')">{{ Number(product.status || 0) === 1 ? '下架' : '上架' }}</up-button>
      </view>
    </view>
  </view>
</template>

<style scoped>
.product-card { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 16rpx; display: flex; gap: 16rpx; align-items: flex-start; position: relative; }
.product-card.is-selected { border-color: var(--eapp-primary); background: var(--eapp-primary-soft); }
.check { position: absolute; top: 14rpx; right: 14rpx; font-size: 28rpx; }
.cover { width: 160rpx; height: 160rpx; border-radius: 14rpx; flex-shrink: 0; background: var(--eapp-bg); display: flex; align-items: center; justify-content: center; color: var(--eapp-text-faint); }
.body { flex: 1; min-width: 0; }
.title-row { display: flex; align-items: center; justify-content: space-between; gap: 10rpx; }
.title { font-size: 28rpx; font-weight: 600; flex: 1; }
.meta { margin-top: 6rpx; color: var(--eapp-text-muted); font-size: 23rpx; }
.cat { margin-top: 4rpx; color: var(--eapp-text-faint); font-size: 22rpx; }
.actions { margin-top: 10rpx; display: flex; gap: 10rpx; }
</style>
```

- [ ] **Step 2：SkuMatrixEditor.vue**

```vue
<script setup lang="ts">
import { computed, ref, watch } from 'vue'
const props = defineProps<{ skus: Array<{ id?: number; attrs: Array<{ name: string; value: string }>; price: number; stock: number }>; basePrice?: number }>()
const emit = defineEmits<{ (e: 'update', skus: any[]): void }>()

const specs = ref<Array<{ name: string; values: string[] }>>([])
const matrix = ref<Array<{ id?: number; attrs: Array<{ name: string; value: string }>; price: number; stock: number }>>([])
const bulkPrice = ref('')
const bulkStock = ref('')

function rebuildFromSkus() {
  const groups: Record<string, Set<string>> = {}
  for (const sku of props.skus || []) {
    for (const a of sku.attrs || []) {
      groups[a.name] = groups[a.name] || new Set()
      groups[a.name].add(a.value)
    }
  }
  specs.value = Object.keys(groups).map((name) => ({ name, values: Array.from(groups[name]) }))
  matrix.value = (props.skus || []).slice()
}

function rebuildMatrix() {
  if (!specs.value.length) { matrix.value = []; return }
  function cross(idx: number): Array<Array<{ name: string; value: string }>> {
    if (idx >= specs.value.length) return [[]]
    const sub = cross(idx + 1)
    return specs.value[idx].values.flatMap((v) => sub.map((row) => [{ name: specs.value[idx].name, value: v }, ...row]))
  }
  const combos = cross(0)
  const next = combos.map((attrs) => {
    const key = attrs.map((a) => `${a.name}:${a.value}`).join('|')
    const exist = matrix.value.find((row) => row.attrs.map((a) => `${a.name}:${a.value}`).join('|') === key)
    return exist || { attrs, price: Number(props.basePrice || 0), stock: 0 }
  })
  matrix.value = next
  emit('update', matrix.value)
}

function addSpec() {
  specs.value.push({ name: '规格' + (specs.value.length + 1), values: ['默认'] })
  rebuildMatrix()
}
function removeSpec(idx: number) {
  specs.value.splice(idx, 1)
  rebuildMatrix()
}
function addValue(idx: number) {
  specs.value[idx].values.push('值' + (specs.value[idx].values.length + 1))
  rebuildMatrix()
}
function removeValue(idx: number, vIdx: number) {
  specs.value[idx].values.splice(vIdx, 1)
  rebuildMatrix()
}
function applyBulkPrice() {
  if (!bulkPrice.value) return
  for (const row of matrix.value) row.price = Number(bulkPrice.value)
  emit('update', matrix.value)
}
function applyBulkStock() {
  if (!bulkStock.value) return
  for (const row of matrix.value) row.stock = Number(bulkStock.value)
  emit('update', matrix.value)
}
function onCellChange() { emit('update', matrix.value) }

watch(() => props.skus, () => rebuildFromSkus(), { immediate: true })
</script>

<template>
  <view class="sku-editor">
    <view v-for="(spec, idx) in specs" :key="idx" class="spec">
      <view class="spec-head">
        <up-input v-model="spec.name" class="name-in" />
        <up-button size="mini" type="error" plain @click="removeSpec(idx)">删除规格</up-button>
      </view>
      <view class="vals">
        <view v-for="(v, vIdx) in spec.values" :key="vIdx" class="val">
          <up-input v-model="spec.values[vIdx]" @blur="rebuildMatrix" />
          <text class="x" @click="removeValue(idx, vIdx)">✕</text>
        </view>
        <up-button size="mini" plain @click="addValue(idx)">+ 值</up-button>
      </view>
    </view>
    <up-button size="mini" type="primary" plain @click="addSpec">+ 规格组</up-button>

    <view v-if="matrix.length" class="bulk">
      <up-input v-model="bulkPrice" type="digit" placeholder="批量赋价" class="bulk-in" />
      <up-button size="mini" @click="applyBulkPrice">应用</up-button>
      <up-input v-model="bulkStock" type="number" placeholder="批量赋库存" class="bulk-in" />
      <up-button size="mini" @click="applyBulkStock">应用</up-button>
    </view>

    <view v-if="matrix.length" class="matrix">
      <view class="row head">
        <view class="cell flex">规格组合</view>
        <view class="cell">价格</view>
        <view class="cell">库存</view>
      </view>
      <view v-for="(row, i) in matrix" :key="i" class="row">
        <view class="cell flex">{{ row.attrs.map((a) => `${a.value}`).join(' / ') }}</view>
        <up-input v-model="row.price" type="digit" class="cell" @blur="onCellChange" />
        <up-input v-model="row.stock" type="number" class="cell" @blur="onCellChange" />
      </view>
    </view>
  </view>
</template>

<style scoped>
.sku-editor { display: grid; gap: 14rpx; }
.spec { background: var(--eapp-bg); padding: 14rpx; border-radius: 14rpx; }
.spec-head { display: flex; gap: 12rpx; align-items: center; }
.name-in { flex: 1; }
.vals { margin-top: 10rpx; display: flex; gap: 8rpx; flex-wrap: wrap; align-items: center; }
.val { display: flex; align-items: center; gap: 6rpx; background: var(--eapp-card); padding: 4rpx 8rpx; border-radius: 10rpx; }
.x { color: var(--eapp-danger); font-size: 22rpx; padding: 0 6rpx; }
.bulk { display: flex; gap: 8rpx; align-items: center; flex-wrap: wrap; }
.bulk-in { width: 200rpx; }
.matrix { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 16rpx; padding: 10rpx; }
.row { display: flex; gap: 8rpx; align-items: center; padding: 8rpx 4rpx; border-bottom: 1px dashed var(--eapp-border); }
.row.head { font-weight: 600; color: var(--eapp-text-muted); }
.cell { width: 160rpx; }
.cell.flex { flex: 1; }
</style>
```

- [ ] **Step 3：CategoryTreePicker.vue**

```vue
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getCategoriesTree } from '@/api/category'

const props = defineProps<{ show: boolean; value: number }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'pick', payload: { id: number; path_name: string }): void }>()

const tree = ref<any[]>([])
const expanded = ref<Record<number, boolean>>({})

async function load() { tree.value = (await getCategoriesTree()) || [] }
function toggle(id: number) { expanded.value[id] = !expanded.value[id] }
function pick(node: any, path: string[]) { emit('pick', { id: Number(node.id), path_name: path.concat(node.name).join(' / ') }) }

onMounted(load)
</script>

<template>
  <up-popup :show="show" mode="right" round="0" @close="$emit('close')">
    <view class="drawer">
      <view class="head">
        <text class="title">选择分类</text>
        <text class="close" @click="$emit('close')">关闭</text>
      </view>
      <scroll-view scroll-y class="body">
        <view v-for="root in tree" :key="root.id">
          <view :class="['node lvl1', value === root.id ? 'active' : '']">
            <text class="caret" @click="toggle(root.id)">{{ expanded[root.id] ? '▾' : '▸' }}</text>
            <text class="name" @click="pick(root, [])">{{ root.name }}</text>
          </view>
          <view v-if="expanded[root.id]" class="sub">
            <view v-for="mid in (root.children || [])" :key="mid.id">
              <view :class="['node lvl2', value === mid.id ? 'active' : '']">
                <text class="caret" @click="toggle(mid.id)">{{ expanded[mid.id] ? '▾' : '▸' }}</text>
                <text class="name" @click="pick(mid, [root.name])">{{ mid.name }}</text>
              </view>
              <view v-if="expanded[mid.id]" class="sub">
                <view v-for="leaf in (mid.children || [])" :key="leaf.id" :class="['node lvl3', value === leaf.id ? 'active' : '']">
                  <text class="name" @click="pick(leaf, [root.name, mid.name])">{{ leaf.name }}</text>
                </view>
              </view>
            </view>
          </view>
        </view>
      </scroll-view>
    </view>
  </up-popup>
</template>

<style scoped>
.drawer { width: 600rpx; height: 100vh; background: var(--eapp-bg); display: flex; flex-direction: column; }
.head { display: flex; align-items: center; justify-content: space-between; padding: 24rpx; padding-top: calc(24rpx + env(safe-area-inset-top)); background: var(--eapp-card); }
.title { font-size: 30rpx; font-weight: 700; }
.close { color: var(--eapp-text-muted); font-size: 26rpx; }
.body { flex: 1; padding: 12rpx; }
.node { display: flex; align-items: center; padding: 12rpx 14rpx; border-radius: 10rpx; gap: 10rpx; }
.node.active { background: var(--eapp-primary-soft); color: var(--eapp-primary); }
.lvl2 { padding-left: 32rpx; }
.lvl3 { padding-left: 56rpx; }
.caret { width: 32rpx; color: var(--eapp-text-muted); font-size: 22rpx; }
.name { font-size: 26rpx; flex: 1; }
.sub { padding-left: 4rpx; }
</style>
```

- [ ] **Step 4：RichTextEditor.vue（演示期只读 + 编辑入口提示）**

```vue
<script setup lang="ts">
defineProps<{ html?: string }>()
defineEmits<{ (e: 'requestEdit'): void }>()
</script>

<template>
  <view class="rt">
    <view class="preview">
      <rich-text v-if="html" :nodes="html" />
      <view v-else class="empty">尚未编辑详情内容</view>
    </view>
    <up-button size="mini" plain @click="$emit('requestEdit')">编辑详情</up-button>
  </view>
</template>

<style scoped>
.rt { display: grid; gap: 12rpx; }
.preview { background: var(--eapp-bg); border: 1px solid var(--eapp-border); border-radius: 14rpx; padding: 14rpx; min-height: 160rpx; }
.empty { color: var(--eapp-text-muted); font-size: 24rpx; text-align: center; padding: 30rpx 0; }
</style>
```

### Task E4：product/list.vue 重写 + product/edit.vue 重写 + category-tree.vue 新建

**Files:** 修改 `eapp/pages/product/list.vue`、`edit.vue`；新增 `eapp/pages/product/category-tree.vue`

- [ ] **Step 1：product/list.vue（完整代码）**

```vue
<script setup lang="ts">
import { onLoad, onPullDownRefresh, onReachBottom, onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import PageHeader from '@/components/biz/PageHeader.vue'
import ProductCard from '@/components/biz/ProductCard.vue'
import EmptyState from '@/components/biz/EmptyState.vue'
import FilterDrawer from '@/components/biz/FilterDrawer.vue'
import BatchBar from '@/components/biz/BatchBar.vue'
import BatchResultPopup from '@/components/biz/BatchResultPopup.vue'
import { useProductList } from '@/composables/useProductList'
import { batchUpdateProductCategory, batchUpdateProductPrice, batchUpdateProductStatus, updateProduct } from '@/api/product'
import { getStorage, removeStorage } from '@/utils/storage'

const h = useProductList()
const tabs = [
  { name: '全部', status: '' as ''|0|1, low_stock: false },
  { name: '在售', status: 1 as ''|0|1, low_stock: false },
  { name: '仓库', status: 0 as ''|0|1, low_stock: false },
  { name: '预警', status: '' as ''|0|1, low_stock: true },
]
const current = ref(0)
const showFilter = ref(false); const showResult = ref(false)
const filterDraft = ref<any>({ keyword: '', sort_by: '' })
const result = ref<{ success_ids: number[]; fail: Array<{ id: number; reason: string }> }>({ success_ids: [], fail: [] })

function onTab(idx: number) {
  current.value = idx
  h.applyFilter({ status: tabs[idx].status, low_stock: tabs[idx].low_stock })
  h.load()
}

function openFilter() { filterDraft.value = { keyword: h.filter.value.keyword || '', sort_by: h.filter.value.sort_by || '' }; showFilter.value = true }
function applyFilter() {
  h.applyFilter({ keyword: filterDraft.value.keyword || undefined, sort_by: filterDraft.value.sort_by || undefined })
  showFilter.value = false; h.load()
}
function resetFilter() { filterDraft.value = { keyword: '', sort_by: '' } }

async function onCardAction(p: any, key: string) {
  if (key === 'edit') uni.navigateTo({ url: `/pages/product/edit?id=${p.id}` })
  else if (key === 'toggle_sale') {
    const next = Number(p.status || 0) === 1 ? 0 : 1
    await updateProduct(p.id, { product: { status: next } })
    uni.showToast({ title: '状态已更新', icon: 'success' })
    await h.refresh()
  }
}

async function startBatch(action: 'shelf_on'|'shelf_off'|'category'|'price') {
  if (!h.selectCount.value) { uni.showToast({ title: '请先勾选商品', icon: 'none' }); return }
  if (action === 'shelf_on' || action === 'shelf_off') {
    const r: any = await batchUpdateProductStatus({ ids: h.selectedIds.value, status: action === 'shelf_on' ? 1 : 0 })
    result.value = r; showResult.value = true; h.clearSelect(); await h.refresh()
  } else if (action === 'category') {
    uni.showModal({ title: '批量分类', editable: true, placeholderText: '输入新分类 ID', success: async (m) => {
      if (!m.confirm || !m.content) return
      const r: any = await batchUpdateProductCategory({ ids: h.selectedIds.value, category_id: Number(m.content) })
      result.value = r; showResult.value = true; h.clearSelect(); await h.refresh()
    } })
  } else if (action === 'price') {
    uni.showModal({ title: '批量调价（百分比）', editable: true, placeholderText: '例如 -10 表示降 10%', success: async (m) => {
      if (!m.confirm || !m.content) return
      const r: any = await batchUpdateProductPrice({ ids: h.selectedIds.value, adjustment: { type: 'percent', value: Number(m.content) } })
      result.value = r; showResult.value = true; h.clearSelect(); await h.refresh()
    } })
  }
}

function newProduct() { uni.navigateTo({ url: '/pages/product/edit?id=0' }) }
function aiHint() { uni.showModal({ title: '提示', content: 'AI 生图请到管理后台 AI 工作流操作', showCancel: false }) }
function goCategoryTree() { uni.navigateTo({ url: '/pages/product/category-tree' }) }

onLoad(() => {
  const w = String(getStorage('eapp_product_status_filter') || '')
  if (w === 'warning') { current.value = 3; removeStorage('eapp_product_status_filter') }
  onTab(current.value)
})
onShow(() => h.load())
onPullDownRefresh(async () => { await h.refresh(); uni.stopPullDownRefresh() })
onReachBottom(() => h.loadMore())
</script>

<template>
  <view class="page">
    <PageHeader title="商品">
      <template #right>
        <text class="icon-btn" @click="goCategoryTree">分类</text>
        <text class="icon-btn" @click="openFilter">⚲</text>
      </template>
    </PageHeader>
    <scroll-view scroll-x class="tabs">
      <view class="tabs-inner">
        <text v-for="(tab, idx) in tabs" :key="idx" :class="['tab', current === idx ? 'active' : '']" @click="onTab(idx)">{{ tab.name }}</text>
      </view>
    </scroll-view>
    <view class="list">
      <EmptyState v-if="!h.loading.value && !h.list.value.length" title="暂无商品" />
      <ProductCard
        v-for="p in h.list.value"
        :key="p.id"
        :product="p"
        :selectable="true"
        :selected="h.isSelected(p.id)"
        @toggle="h.toggleSelect(p.id)"
        @click="uni.navigateTo({ url: `/pages/product/edit?id=${p.id}` })"
        @action="(key) => onCardAction(p, key)"
      />
    </view>

    <BatchBar
      :count="h.selectCount.value"
      :actions="[
        { key: 'shelf_on', label: '批量上架', tone: 'primary' },
        { key: 'shelf_off', label: '批量下架', tone: 'warning' },
        { key: 'category', label: '批量分类', tone: 'primary' },
        { key: 'price', label: '批量调价', tone: 'primary' },
      ]"
      @action="startBatch"
      @cancel="h.clearSelect()"
    />

    <FilterDrawer :show="showFilter" title="商品筛选" @close="showFilter = false" @reset="resetFilter" @confirm="applyFilter">
      <up-input v-model="filterDraft.keyword" placeholder="商品名" />
      <picker mode="selector" :range="['默认', '销量', '库存', '价格升', '价格降', '最新']" @change="(e) => { filterDraft.sort_by = ['', 'sales', 'stock', 'price_asc', 'price_desc', 'created'][Number(e.detail.value)] }">
        <view class="picker mt">排序：{{ { '': '默认', sales: '销量', stock: '库存', price_asc: '价格升', price_desc: '价格降', created: '最新' }[filterDraft.sort_by || ''] }}</view>
      </picker>
    </FilterDrawer>

    <BatchResultPopup :show="showResult" :success="result.success_ids" :fails="result.fail" @close="showResult = false" />

    <view class="fab">
      <up-button type="primary" shape="circle" @click="newProduct">+ 新建</up-button>
      <up-button shape="circle" plain @click="aiHint">AI</up-button>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.icon-btn { font-size: 28rpx; padding: 0 12rpx; color: var(--eapp-primary); }
.tabs { background: var(--eapp-card); position: sticky; top: 0; z-index: 20; white-space: nowrap; }
.tabs-inner { display: inline-flex; padding: 16rpx 12rpx; gap: 8rpx; }
.tab { padding: 10rpx 20rpx; border-radius: 999rpx; font-size: 24rpx; color: var(--eapp-text-muted); background: var(--eapp-bg); }
.tab.active { background: var(--eapp-primary); color: #fff; }
.list { padding: 20rpx; display: grid; gap: 14rpx; padding-bottom: 200rpx; }
.fab { position: fixed; right: 24rpx; bottom: calc(40rpx + env(safe-area-inset-bottom)); display: grid; gap: 12rpx; }
.picker { display: inline-flex; height: 60rpx; align-items: center; padding: 0 18rpx; border: 1px solid var(--eapp-border); border-radius: 12rpx; background: var(--eapp-card); }
.mt { margin-top: 12rpx; }
</style>
```

- [ ] **Step 2：product/edit.vue（完整代码）**

```vue
<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import SkuMatrixEditor from '@/components/biz/SkuMatrixEditor.vue'
import CategoryTreePicker from '@/components/biz/CategoryTreePicker.vue'
import RichTextEditor from '@/components/biz/RichTextEditor.vue'
import { createProduct, getProductDetail, updateProduct } from '@/api/product'

const id = ref(0)
const saving = ref(false)
const showCatPicker = ref(false)

const form = reactive<any>({
  title: '', subtitle: '', sell_points: [] as string[],
  covers: [] as string[],
  detail_html: '',
  price: 0, stock: 0, unit: '件',
  weight: 0,
  category_id: 0, category_path_name: '',
  tags: [] as string[],
  skus: [] as Array<{ id?: number; attrs: Array<{ name: string; value: string }>; price: number; stock: number }>,
  low_stock_threshold: 10,
  shipping_template: 'default',
  limit_per_order: 0,
  exclude_marketing: false,
  status: 1,
  online_at: '', offline_at: '',
})

const newSellPoint = ref('')
const newCover = ref('')
const newTag = ref('')

async function loadData() {
  if (!id.value) return
  const data: any = await getProductDetail(id.value)
  Object.assign(form, {
    title: data?.title || '', subtitle: data?.subtitle || '',
    sell_points: Array.isArray(data?.sell_points) ? data.sell_points : [],
    covers: Array.isArray(data?.covers) ? data.covers : (data?.cover ? [data.cover] : []),
    detail_html: String(data?.detail_html || ''),
    price: Number(data?.price || 0), stock: Number(data?.stock || 0),
    unit: String(data?.unit || '件'), weight: Number(data?.weight || 0),
    category_id: Number(data?.category_id || 0),
    category_path_name: String(data?.category_path_name || ''),
    tags: Array.isArray(data?.tags) ? data.tags : [],
    skus: Array.isArray(data?.skus) ? data.skus : [],
    low_stock_threshold: Number(data?.low_stock_threshold || 10),
    shipping_template: String(data?.shipping_template || 'default'),
    limit_per_order: Number(data?.limit_per_order || 0),
    exclude_marketing: !!data?.exclude_marketing,
    status: Number(data?.status || 0) === 1 ? 1 : 0,
    online_at: String(data?.online_at || ''),
    offline_at: String(data?.offline_at || ''),
  })
}

function addSellPoint() { if (newSellPoint.value.trim()) { form.sell_points.push(newSellPoint.value.trim()); newSellPoint.value = '' } }
function removeSellPoint(i: number) { form.sell_points.splice(i, 1) }
function addCover() { if (newCover.value.trim()) { form.covers.push(newCover.value.trim()); newCover.value = '' } }
function removeCover(i: number) { form.covers.splice(i, 1) }
function addTag() { if (newTag.value.trim()) { form.tags.push(newTag.value.trim()); newTag.value = '' } }
function removeTag(i: number) { form.tags.splice(i, 1) }
function onSkusChange(rows: any[]) { form.skus = rows }
function onPickCategory(payload: { id: number; path_name: string }) {
  form.category_id = payload.id; form.category_path_name = payload.path_name; showCatPicker.value = false
}
function requestEditDetail() {
  uni.setClipboardData({ data: `${form.title} - 详情编辑` })
  uni.showModal({ title: '提示', content: '长详情编辑请到管理后台进行；标题已复制到剪贴板。', showCancel: false })
}

async function save() {
  if (!form.title.trim()) { uni.showToast({ title: '请输入商品标题', icon: 'none' }); return }
  saving.value = true
  try {
    const payload: any = {
      product: {
        title: form.title.trim(), subtitle: form.subtitle.trim(),
        sell_points: form.sell_points,
        cover: form.covers[0] || '',
        covers: form.covers,
        detail_html: form.detail_html,
        price: Number(form.price), stock: Number(form.stock),
        unit: form.unit, weight: Number(form.weight),
        category_id: Number(form.category_id), category_path_name: form.category_path_name,
        tags: form.tags,
        low_stock_threshold: Number(form.low_stock_threshold),
        shipping_template: form.shipping_template,
        limit_per_order: Number(form.limit_per_order),
        exclude_marketing: !!form.exclude_marketing,
        status: form.status,
        online_at: form.online_at, offline_at: form.offline_at,
      },
      skus: form.skus.map((s: any) => ({ ...s, attrs: typeof s.attrs === 'string' ? s.attrs : JSON.stringify(s.attrs), price: Number(s.price), stock: Number(s.stock) })),
    }
    if (id.value) await updateProduct(id.value, payload)
    else await createProduct(payload)
    uni.showToast({ title: '保存成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 350)
  } finally { saving.value = false }
}

onLoad((opts) => { id.value = Number(opts?.id || 0); loadData() })
</script>

<template>
  <view class="page">
    <view class="section">
      <view class="section-title">基础信息</view>
      <up-input v-model="form.title" placeholder="商品标题" class="mt" />
      <up-input v-model="form.subtitle" placeholder="副标题" class="mt" />
      <view class="tag-row mt">
        <view v-for="(sp, i) in form.sell_points" :key="i" class="tag">{{ sp }}<text class="x" @click="removeSellPoint(i)">✕</text></view>
        <up-input v-model="newSellPoint" placeholder="添加卖点" class="tag-input" @confirm="addSellPoint" />
      </view>
    </view>

    <view class="section">
      <view class="section-title">主图轮播</view>
      <view class="cover-row">
        <view v-for="(c, i) in form.covers" :key="i" class="cover-cell">
          <image :src="c" mode="aspectFill" class="cover-img" />
          <text class="x-abs" @click="removeCover(i)">✕</text>
        </view>
      </view>
      <view class="add-cover mt">
        <up-input v-model="newCover" placeholder="图片 URL" />
        <up-button size="mini" type="primary" @click="addCover">添加</up-button>
      </view>
    </view>

    <view class="section">
      <view class="section-title">商品详情</view>
      <RichTextEditor :html="form.detail_html" @requestEdit="requestEditDetail" />
    </view>

    <view class="section">
      <view class="section-title">价格库存</view>
      <view class="grid-2 mt">
        <up-input v-model="form.price" type="digit" placeholder="主价格" />
        <up-input v-model="form.stock" type="number" placeholder="主库存" />
        <up-input v-model="form.unit" placeholder="单位" />
        <up-input v-model="form.weight" type="digit" placeholder="重量(kg)" />
        <up-input v-model="form.low_stock_threshold" type="number" placeholder="预警阈值" />
        <up-input v-model="form.limit_per_order" type="number" placeholder="单笔限购(0=不限)" />
      </view>
    </view>

    <view class="section">
      <view class="section-title">规格 SKU</view>
      <SkuMatrixEditor :skus="form.skus" :base-price="form.price" @update="onSkusChange" />
    </view>

    <view class="section">
      <view class="section-title">分类与标签</view>
      <view class="cat" @click="showCatPicker = true">分类：{{ form.category_path_name || '请选择' }} <text class="caret">›</text></view>
      <view class="tag-row mt">
        <view v-for="(tag, i) in form.tags" :key="i" class="tag">{{ tag }}<text class="x" @click="removeTag(i)">✕</text></view>
        <up-input v-model="newTag" placeholder="添加标签" class="tag-input" @confirm="addTag" />
      </view>
    </view>

    <view class="section">
      <view class="section-title">物流与营销</view>
      <picker mode="selector" :range="['默认模板','包邮','到付','同城']" @change="(e) => form.shipping_template = ['default','free','cod','local'][Number(e.detail.value)]">
        <view class="picker">物流模板：{{ { default: '默认模板', free: '包邮', cod: '到付', local: '同城' }[form.shipping_template] }}</view>
      </picker>
      <view class="row mt">
        <text>不参与营销活动</text>
        <switch :checked="form.exclude_marketing" @change="(e) => form.exclude_marketing = e.detail.value" />
      </view>
    </view>

    <view class="section">
      <view class="section-title">状态控制</view>
      <view class="row">
        <text>上架</text>
        <switch :checked="form.status === 1" @change="(e) => form.status = e.detail.value ? 1 : 0" />
      </view>
      <up-input v-model="form.online_at" placeholder="上架时间 ISO（可选）" class="mt" />
      <up-input v-model="form.offline_at" placeholder="下架时间 ISO（可选）" class="mt" />
    </view>

    <up-button type="primary" :loading="saving" class="save" @click="save">保存</up-button>

    <CategoryTreePicker :show="showCatPicker" :value="form.category_id" @close="showCatPicker = false" @pick="onPickCategory" />
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; padding-bottom: 160rpx; display: grid; gap: 16rpx; box-sizing: border-box; }
.section { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; }
.section-title { font-size: 30rpx; font-weight: 700; margin-bottom: 10rpx; display: flex; align-items: center; justify-content: space-between; }
.mt { margin-top: 12rpx; }
.tag-row { display: flex; gap: 10rpx; flex-wrap: wrap; align-items: center; }
.tag { background: var(--eapp-bg); border-radius: 999rpx; padding: 6rpx 16rpx; font-size: 22rpx; display: inline-flex; align-items: center; gap: 6rpx; }
.tag .x { color: var(--eapp-danger); font-size: 20rpx; }
.tag-input { flex: 1; min-width: 180rpx; }
.cover-row { display: flex; gap: 10rpx; flex-wrap: wrap; }
.cover-cell { position: relative; width: 160rpx; height: 160rpx; }
.cover-img { width: 100%; height: 100%; border-radius: 14rpx; }
.x-abs { position: absolute; top: -10rpx; right: -10rpx; background: var(--eapp-danger); color: #fff; font-size: 20rpx; width: 32rpx; height: 32rpx; display: flex; align-items: center; justify-content: center; border-radius: 50%; }
.add-cover { display: flex; gap: 10rpx; align-items: center; }
.grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 12rpx; }
.cat { padding: 14rpx; background: var(--eapp-bg); border-radius: 14rpx; display: flex; justify-content: space-between; align-items: center; font-size: 26rpx; }
.caret { color: var(--eapp-text-muted); }
.picker { padding: 14rpx; background: var(--eapp-bg); border-radius: 14rpx; font-size: 26rpx; }
.row { display: flex; align-items: center; justify-content: space-between; padding: 8rpx 0; font-size: 26rpx; }
.save { margin-top: 14rpx; }
</style>
```

- [ ] **Step 3：product/category-tree.vue（完整代码）**

```vue
<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { createCategory, deleteCategory, getCategoriesTree, updateCategory } from '@/api/category'

const tree = ref<any[]>([])
const expanded = ref<Record<number, boolean>>({})

async function load() { tree.value = (await getCategoriesTree()) || [] }

function toggle(id: number) { expanded.value[id] = !expanded.value[id] }

function addUnder(parent: any) {
  uni.showModal({ title: '新增分类', editable: true, placeholderText: '请输入名称', success: async (m) => {
    if (!m.confirm || !m.content) return
    await createCategory({ name: m.content, parent_id: Number(parent?.id || 0) })
    await load()
  } })
}

function rename(node: any) {
  uni.showModal({ title: '重命名', editable: true, content: node.name, success: async (m) => {
    if (!m.confirm || !m.content) return
    await updateCategory(Number(node.id), { name: m.content })
    await load()
  } })
}

function remove(node: any) {
  uni.showModal({ title: '删除', content: `确认删除分类「${node.name}」？`, success: async (m) => {
    if (!m.confirm) return
    await deleteCategory(Number(node.id))
    await load()
  } })
}

onShow(load)
</script>

<template>
  <view class="page">
    <up-button type="primary" plain class="add-top" @click="addUnder({ id: 0 })">+ 新增根分类</up-button>
    <view v-for="root in tree" :key="root.id" class="node-block">
      <view class="node">
        <text class="caret" @click="toggle(root.id)">{{ expanded[root.id] ? '▾' : '▸' }}</text>
        <text class="name">{{ root.name }}</text>
        <view class="ops">
          <text class="op" @click="addUnder(root)">+ 子</text>
          <text class="op" @click="rename(root)">改名</text>
          <text class="op danger" @click="remove(root)">删</text>
        </view>
      </view>
      <view v-if="expanded[root.id]" class="children">
        <view v-for="mid in (root.children || [])" :key="mid.id" class="node-block">
          <view class="node lvl2">
            <text class="caret" @click="toggle(mid.id)">{{ expanded[mid.id] ? '▾' : '▸' }}</text>
            <text class="name">{{ mid.name }}</text>
            <view class="ops">
              <text class="op" @click="addUnder(mid)">+ 子</text>
              <text class="op" @click="rename(mid)">改名</text>
              <text class="op danger" @click="remove(mid)">删</text>
            </view>
          </view>
          <view v-if="expanded[mid.id]" class="children">
            <view v-for="leaf in (mid.children || [])" :key="leaf.id" class="node lvl3">
              <text class="name">{{ leaf.name }}</text>
              <view class="ops">
                <text class="op" @click="rename(leaf)">改名</text>
                <text class="op danger" @click="remove(leaf)">删</text>
              </view>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; padding: 20rpx; background: var(--eapp-bg); }
.add-top { margin-bottom: 14rpx; }
.node-block { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 14rpx; padding: 6rpx 10rpx; margin-bottom: 8rpx; }
.node { display: flex; align-items: center; padding: 12rpx 8rpx; gap: 10rpx; }
.node.lvl2 { padding-left: 30rpx; }
.node.lvl3 { padding-left: 60rpx; }
.caret { width: 28rpx; color: var(--eapp-text-muted); }
.name { flex: 1; font-size: 26rpx; }
.ops { display: flex; gap: 12rpx; }
.op { font-size: 24rpx; color: var(--eapp-primary); }
.op.danger { color: var(--eapp-danger); }
.children { padding: 4rpx 8rpx; }
</style>
```

- [ ] **Step 4：dev 验证 + 跑测试 + commit 5**

```bash
cd eapp && npm run dev:h5 -- --mode demo
# 验证商品列表分类树 / 多规格编辑 / 批量上下架 / 分类管理 CRUD 可用
cd eapp && npx vitest run
git -C 'D:\Repos\xyito\open\lyshop' add eapp admin/src/mock/index.ts
git -C 'D:\Repos\xyito\open\lyshop' commit -m "eapp: 商品列表与编辑（多规格 SKU + 分类）"
```

commit body：

```
- mock /products 扩展 status/category/sort_by/low_stock 筛选，新增 /products/batch/status|category|price
- mock 新增 /categories/tree + /categories CRUD；预置 20 节点三级树；6 商品补 SKU 矩阵与分类
- api/category.ts 新增；api/product.ts 增 batch 封装；composables/useProductList 沉淀
- 新增 ProductCard、SkuMatrixEditor（规格组动态 + 矩阵 + 批量赋值）、CategoryTreePicker、RichTextEditor（演示期只读）
- product/list 重写：状态 tabs（在售/仓库/预警）、筛选抽屉、批量上下架/分类/调价、FAB
- product/edit 重写：分段表单（基础/轮播/详情/价库/SKU/分类/物流/营销/状态）
- 新增 product/category-tree（三级分类 CRUD）；pages.json 已注册
```

---

PLANEND
