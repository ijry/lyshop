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
    <up-tabs
      :list="tabs"
      :current="current"
      :scrollable="true"
      keyName="name"
      @click="(item) => onTabChange(item.index)"
      :activeStyle="{ color: '#fff', backgroundColor: 'var(--eapp-primary)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :inactiveStyle="{ color: 'var(--eapp-text-muted)', backgroundColor: 'var(--eapp-bg)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :itemStyle="{ padding: '0 4rpx', height: '80rpx' }"
      lineColor="transparent"
    />
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
.list { padding: 20rpx; display: grid; gap: 16rpx; padding-bottom: 200rpx; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
.row { display: flex; align-items: center; gap: 12rpx; }
.dash { color: var(--eapp-text-muted); }
.mt { margin-top: 12rpx; }
</style>
