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
  { name: '已取消', status: '6' },
  { name: '售后中', status: '', has_after_sale: true },
]
const quickKeyword = ref('')
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
  h.applyFilter({ status: tab.status, has_after_sale: !!tab.has_after_sale, keyword: quickKeyword.value || undefined })
  h.load()
}

function onQuickSearch() {
  h.applyFilter({ keyword: quickKeyword.value || undefined })
  h.load()
}

function clearQuickSearch() {
  quickKeyword.value = ''
  h.applyFilter({ keyword: undefined })
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
    <PageHeader title="订单管理">
      <template #right>
        <view class="filter-btn" @click="openFilter">
          <svg width="40" height="40" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M3 4.5h18M7 12h10M10.5 19.5h3" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
          </svg>
        </view>
      </template>
    </PageHeader>

    <!-- Tab bar -->
    <view class="tab-wrap">
      <up-tabs
        :list="tabs"
        :current="current"
        :scrollable="true"
        keyName="name"
        @click="(item: any) => onTabChange(item.index)"
        :activeStyle="{ color: 'var(--eapp-primary)', fontWeight: '700', fontSize: '28rpx' }"
        :inactiveStyle="{ color: 'var(--eapp-text-muted)', fontSize: '26rpx' }"
        :itemStyle="{ height: '80rpx', padding: '0 8rpx' }"
        lineColor="var(--eapp-primary)"
        lineWidth="40rpx"
        lineHeight="4rpx"
        :lineStyle="{ borderRadius: '2rpx' }"
      />
    </view>

    <!-- Quick search bar -->
    <view class="search-bar">
      <up-search
        v-model="quickKeyword"
        placeholder="搜索订单号 / 买家 / 商品"
        :showAction="!!quickKeyword"
        actionText="清空"
        shape="round"
        :inputStyle="{ fontSize: '26rpx' }"
        @search="onQuickSearch"
        @custom="clearQuickSearch"
      />
    </view>

    <!-- Order count hint -->
    <view v-if="h.total && h.total.value > 0" class="list-meta">
      <text class="list-meta-text">共 {{ h.total.value }} 个订单</text>
    </view>

    <!-- List -->
    <view class="list">
      <!-- Skeleton loading -->
      <template v-if="h.loading.value && !h.list.value.length">
        <view v-for="i in 4" :key="i" class="skeleton-card">
          <view class="sk-accent" />
          <view class="sk-body">
            <view class="sk-row">
              <view class="sk-block" style="width: 160rpx; height: 28rpx;" />
              <view class="sk-block" style="width: 80rpx; height: 28rpx;" />
            </view>
            <view class="sk-row" style="margin-top: 16rpx; gap: 8rpx;">
              <view class="sk-block" style="width: 84rpx; height: 84rpx; border-radius: 10rpx;" />
              <view class="sk-block" style="width: 84rpx; height: 84rpx; border-radius: 10rpx;" />
              <view style="flex: 1; display: flex; flex-direction: column; gap: 8rpx; padding-top: 4rpx;">
                <view class="sk-block" style="width: 80%; height: 24rpx;" />
                <view class="sk-block" style="width: 50%; height: 24rpx;" />
              </view>
            </view>
            <view class="sk-row" style="margin-top: 12rpx;">
              <view class="sk-block" style="width: 120rpx; height: 22rpx;" />
              <view class="sk-block" style="width: 100rpx; height: 36rpx;" />
            </view>
          </view>
        </view>
      </template>

      <EmptyState
        v-else-if="!h.loading.value && !h.list.value.length"
        title="暂无订单"
        desc="切换状态或调整筛选条件试试"
      />

      <OrderCard
        v-for="o in h.list.value"
        :key="o.id"
        :order="o"
        :selectable="true"
        :selected="h.isSelected(o.id)"
        @toggle="h.toggleSelect(o.id)"
        @click="onCardAction(o, 'detail')"
        @action="(key: string) => onCardAction(o, key)"
      />
      <view v-if="h.loading.value && h.list.value.length" class="loading">
        <up-loading-icon mode="circle" size="32" />
        <text class="loading-text">加载中…</text>
      </view>
      <view v-if="!h.loading.value && h.list.value.length && h.list.value.length >= h.total.value && h.total.value > 0" class="no-more">
        · 没有更多了 ·
      </view>
    </view>

    <!-- Batch action bar -->
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

    <!-- Filter drawer -->
    <FilterDrawer :show="showFilter" title="高级筛选" @close="showFilter = false" @reset="resetFilter" @confirm="applyFilter">
      <up-input v-model="filterDraft.keyword" placeholder="订单号 / 买家昵称 / 商品名" />
      <view class="row mt">
        <up-input v-model="filterDraft.amount_min" type="digit" placeholder="最低金额" />
        <text class="dash">—</text>
        <up-input v-model="filterDraft.amount_max" type="digit" placeholder="最高金额" />
      </view>
      <up-input v-model="filterDraft.logistics_company" placeholder="物流公司（SF / ZTO / ...）" class="mt" />
      <up-input v-model="filterDraft.province" placeholder="收货省份" class="mt" />
      <up-input v-model="filterDraft.pay_method" placeholder="支付方式（wechat / alipay / ...）" class="mt" />
    </FilterDrawer>

    <ShipPopup :show="showShip" :delivery-mode="deliveryMode" :loading="actionLoading" @close="showShip = false" @submit="submitShip" />
    <RepricingPopup :show="showReprice" :items="activeOrder?.items || []" :loading="actionLoading" @close="showReprice = false" @submit="submitReprice" />
    <RemarkPopup :show="showNote" :loading="actionLoading" @close="showNote = false" @submit="submitNote" />
    <BatchResultPopup :show="showResult" :success="result.success_ids" :fails="result.fail" @close="showResult = false" />
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.filter-btn { padding: 4rpx 8rpx; color: var(--eapp-primary); display: flex; align-items: center; }

/* Tabs */
.tab-wrap { background: var(--eapp-card); border-bottom: 1rpx solid var(--eapp-border); }

/* Search bar */
.search-bar { background: var(--eapp-card); padding: 12rpx 24rpx 16rpx; border-bottom: 1rpx solid var(--eapp-border); }

/* List meta */
.list-meta { padding: 12rpx 24rpx 4rpx; }
.list-meta-text { font-size: 22rpx; color: var(--eapp-text-faint); }

/* List */
.list { padding: 16rpx 20rpx; display: grid; gap: 16rpx; padding-bottom: 200rpx; }

/* Loading inline */
.loading { display: flex; align-items: center; justify-content: center; gap: 12rpx; padding: 24rpx 0; }
.loading-text { font-size: 26rpx; color: var(--eapp-text-muted); }

/* No more */
.no-more { text-align: center; font-size: 24rpx; color: var(--eapp-text-faint); padding: 16rpx 0; letter-spacing: 2rpx; }

/* Filter drawer helpers */
.row { display: flex; align-items: center; gap: 12rpx; }
.dash { color: var(--eapp-text-muted); font-size: 28rpx; flex-shrink: 0; }
.mt { margin-top: 16rpx; }

/* Skeleton */
.skeleton-card { display: flex; flex-direction: row; background: var(--eapp-card); border-radius: 20rpx; overflow: hidden; }
.sk-accent { width: 8rpx; background: var(--eapp-border); }
.sk-body { flex: 1; padding: 20rpx 22rpx; }
.sk-row { display: flex; align-items: center; justify-content: space-between; }
.sk-block { background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%); background-size: 200% 100%; border-radius: 8rpx; animation: shimmer 1.5s infinite; }
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }
</style>
