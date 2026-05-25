<template>
  <view class="min-h-screen bg-gray-50">
    <view class="bg-white">
      <u-tabs :list="tabs" :current="activeTab" @click="(item:any) => onTab(item.index)" />
    </view>

    <view class="p-20rpx">
      <view v-if="!orders.length" class="flex flex-col items-center py-120rpx">
        <u-icon name="order" size="60" color="#ccc" />
        <text class="text-gray-400 text-28rpx mt-20rpx">{{ $t('orderList.empty') }}</text>
      </view>

      <view
        v-for="o in orders"
        :key="o.id"
        class="bg-white rounded-20rpx p-30rpx mb-20rpx"
        style="box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);"
      >
        <view class="flex items-center justify-between mb-20rpx">
          <view class="flex items-center" style="gap: 8px;">
            <text class="text-22rpx text-gray-400 font-mono">{{ o.order_no }}</text>
            <text v-if="o.activity_type === 'seckill'" style="font-size: 10px; color: #dc2626; background: #fef2f2; padding: 1px 6px; border-radius: 4px;">{{ $t('orderList.seckill') }}</text>
            <text v-else-if="o.activity_type === 'group_buy'" style="font-size: 10px; color: #2563eb; background: #eff6ff; padding: 1px 6px; border-radius: 4px;">{{ $t('orderList.groupBuy') }}</text>
            <text v-else-if="o.activity_type === 'bargain'" style="font-size: 10px; color: #16a34a; background: #f0fdf4; padding: 1px 6px; border-radius: 4px;">{{ $t('orderList.bargain') }}</text>
          </view>
          <text :class="statusColor(o.status)" class="text-24rpx font-500">{{ statusLabel(o.status) }}</text>
        </view>

        <view v-if="o.items?.length" class="mb-16rpx">
          <view v-for="it in o.items.slice(0, 2)" :key="it.id" class="flex items-center mb-12rpx" style="gap: 10px;">
            <image :src="it.cover" mode="aspectFill" style="width: 72rpx; height: 72rpx; border-radius: 10rpx;" />
            <view class="flex-1 min-w-0">
              <text class="text-24rpx text-gray-700 block truncate">{{ it.title }}</text>
              <text class="text-22rpx text-gray-400 mt-4rpx block">x{{ it.qty }}</text>
            </view>
          </view>
        </view>

        <view v-if="hasShipmentSummary(o)" class="flex items-center flex-wrap mb-16rpx" style="gap: 8rpx;">
          <text
            v-if="hasReship(o)"
            style="font-size: 10px; color: #c2410c; background: #fff7ed; padding: 1px 6px; border-radius: 4px;"
          >
            {{ $t('orderList.includesReship') }}
          </text>
          <text v-if="o.latest_shipment?.tracking_no" class="text-22rpx text-gray-500">
            {{ $t('orderList.latestLogistics') }}{{ shipmentStatusText(o.latest_shipment?.logistics_status, o.latest_shipment?.logistics_status_label) }} · {{ o.latest_shipment?.tracking_no }}
          </text>
          <text v-if="shipmentPrimaryTime(o.latest_shipment)" class="text-22rpx text-gray-400">
            {{ formatDate(shipmentPrimaryTime(o.latest_shipment)) }}
          </text>
        </view>
        <view v-if="o.after_sale_summary?.latest_case_id" class="flex items-center flex-wrap mb-16rpx" style="gap: 8rpx;">
          <text
            v-if="o.after_sale_summary?.has_open_case"
            style="font-size: 10px; color: #dc2626; background: #fef2f2; padding: 1px 6px; border-radius: 4px;"
          >
            {{ $t('orderList.afterSaleInProgress') }}
          </text>
          <text class="text-22rpx text-gray-500">
            {{ $t('orderList.recentAfterSale') }}{{ o.after_sale_summary.latest_case_id }}（{{ afterSaleSummaryStatusText(o.after_sale_summary) || '-' }}）
          </text>
        </view>

        <view class="flex items-center justify-between">
          <text class="text-24rpx text-gray-500">{{ o.created_at?.slice(0, 10) }}</text>
          <text class="text-30rpx text-gray-800 font-700">¥{{ money(o.amount_breakdown?.payable_amount ?? o.total_amount) }}</text>
        </view>

        <view class="flex justify-end gap-16rpx mt-24rpx">
          <view class="action-btn-wrap">
            <u-button size="mini" plain :text="$t('orderList.viewDetail')" shape="circle" @click="goDetail(o.id)" />
          </view>
          <view class="action-btn-wrap" v-if="o.after_sale_summary?.latest_case_id">
            <u-button size="mini" plain :text="$t('orderList.afterSaleProgress')" shape="circle" @click="goAfterSaleDetail(o.after_sale_summary.latest_case_id)" />
          </view>
          <view class="action-btn-wrap" v-if="o.status === 1">
            <u-button size="mini" type="primary" :text="$t('orderList.goPay')" shape="circle" :loading="actioningID === o.id && actioningType === 'pay'" @click="toPay(o)" />
          </view>
          <view class="action-btn-wrap" v-if="canReview(o) && hasUnreviewed(o)">
            <u-button size="mini" type="success" :text="$t('orderList.review')" shape="circle" plain :loading="actioningID === o.id && actioningType === 'review'" @click="toReview(o, 'root')" />
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { get, post } from '@/utils/request'
import { afterSaleStatusLabel, hasReshipShipment, orderStatusLabel, shipmentPrimaryTime, shipmentStatusLabel } from '@/utils/order-status'

const { t } = useI18n()

const orders = ref<any[]>([])
const activeTab = ref(0)
const actioningID = ref<number>(0)
const actioningType = ref<'pay' | 'review' | ''>('')

const tabs = computed(() => [
  { name: t('orderList.all') }, { name: t('orderList.pendingPayment') }, { name: t('orderList.pendingShipment') },
  { name: t('orderList.pendingReceipt') }, { name: t('orderList.completed') }
])
const statusValues = [0, 1, 2, 3, 4]

const statusColors: Record<number, string> = {
  1: 'text-orange-500', 2: 'text-blue-500',
  3: 'text-purple-500', 4: 'text-green-500', 5: 'text-red-500'
}
const statusLabel = (s: number) => orderStatusLabel(s)
const statusColor = (s: number) => statusColors[s] || 'text-gray-400'
const money = (v: number) => Number(v || 0).toFixed(2)
const formatDate = (v?: string) => (v ? String(v).slice(0, 19).replace('T', ' ') : '-')

function shipmentStatusText(status: string, label?: string) {
  const mapped = String(label || '').trim()
  if (mapped) return mapped
  return shipmentStatusLabel(status)
}

function afterSaleSummaryStatusText(summary: any) {
  const label = String(summary?.latest_status_label || '').trim()
  if (label) return label
  return afterSaleStatusLabel(summary?.latest_status)
}

function hasReship(order: any) {
  return hasReshipShipment(Array.isArray(order?.shipments) ? order.shipments : [])
}

function hasShipmentSummary(order: any) {
  return Boolean(hasReship(order) || order?.latest_shipment?.tracking_no)
}

async function loadOrders() {
  const status = statusValues[activeTab.value]
  const data = await get<any>('/api/v1/orders', { status: status || undefined, page: 1, size: 20 })
  orders.value = data?.list || []
}

function onTab(index: number) {
  activeTab.value = index
  loadOrders()
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/order/detail?id=${id}` })
}

function goAfterSaleDetail(id: number) {
  if (!id) return
  uni.navigateTo({ url: `/pages/order/after-sale-detail?id=${id}` })
}

function canReview(order: any) {
  const status = Number(order?.status || 0)
  return status === 3 || status === 4
}

function hasUnreviewed(order: any) {
  const items = Array.isArray(order?.items) ? order.items : []
  if (!items.length) return false
  return items.some((item: any) => !item?.review?.id)
}

async function toPay(order: any) {
  const id = Number(order?.id || 0)
  if (!id || actioningID.value) return
  actioningID.value = id
  actioningType.value = 'pay'
  try {
    await post(`/api/v1/orders/${id}/pay`)
    uni.showToast({ title: t('orderConfirm.orderSuccess'), icon: 'success' })
    await loadOrders()
  } catch {
    uni.showToast({ title: t('afterSaleApply.submitFailed'), icon: 'none' })
  } finally {
    actioningID.value = 0
    actioningType.value = ''
  }
}

function toReview(order: any, mode: 'root' | 'append' = 'root') {
  const id = Number(order?.id || 0)
  if (!id || actioningID.value) return
  uni.navigateTo({ url: `/pages/order/review?id=${id}&mode=${mode}` })
}

onMounted(loadOrders)
</script>

<style scoped>
.action-btn-wrap {
  display: inline-flex;
}
</style>
