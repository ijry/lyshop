<template>
  <div class="max-w-5xl mx-auto px-6 py-8">
    <h1 class="text-xl font-bold text-gray-900 mb-6">{{ $t('orderList.title') }}</h1>

    <div class="flex gap-1 bg-gray-100 rounded-xl p-1 mb-6 w-fit">
      <button
        v-for="(tab, i) in tabs"
        :key="tab"
        @click="activeTab = i; loadOrders()"
        :class="activeTab === i ? 'bg-white shadow-sm text-gray-900' : 'text-gray-500 hover:text-gray-700'"
        class="px-5 py-2 rounded-lg text-sm font-medium transition-all"
      >
        {{ tab }}
      </button>
    </div>

    <div v-if="!orders.length" class="card p-16 text-center">
      <div class="i-carbon-document text-6xl text-gray-200 mx-auto mb-4" />
      <p class="text-gray-400">{{ $t('orderList.empty') }}</p>
    </div>

    <div v-else class="space-y-4">
      <div v-for="o in orders" :key="o.id" class="card p-5 hover:shadow-md transition-shadow">
        <div class="flex-between mb-4">
          <div class="flex items-center gap-3">
            <span class="text-xs font-mono text-gray-400 bg-gray-50 px-2 py-1 rounded">{{ o.order_no }}</span>
            <span class="text-xs text-gray-400">{{ o.created_at?.slice(0, 10) }}</span>
          </div>
          <span :class="statusColor(o.status)" class="text-xs font-medium px-2.5 py-1 rounded-full">
            {{ statusLabel(o.status) }}
          </span>
        </div>

        <div v-if="o.items?.length" class="space-y-2 mb-4">
          <div v-for="it in o.items.slice(0, 2)" :key="it.id" class="flex items-center gap-3">
            <img :src="it.cover" class="w-12 h-12 rounded-lg object-cover" />
            <div class="flex-1 min-w-0">
              <p class="text-sm text-gray-700 truncate">{{ it.title }}</p>
              <p class="text-xs text-gray-400">x{{ it.qty }}</p>
            </div>
          </div>
        </div>

        <div v-if="hasShipmentSummary(o)" class="flex items-center flex-wrap gap-2 mb-4">
          <span v-if="hasReship(o)" class="text-[11px] px-2 py-0.5 rounded-full bg-orange-50 text-orange-600">{{ $t('orderList.includeReship') }}</span>
          <span v-if="o.latest_shipment?.tracking_no" class="text-xs text-gray-500">
            {{ $t('orderList.latestLogistics') }}{{ shipmentStatusText(o.latest_shipment?.logistics_status, o.latest_shipment?.logistics_status_label) }} · {{ o.latest_shipment?.tracking_no }}
          </span>
          <span v-if="shipmentPrimaryTime(o.latest_shipment)" class="text-xs text-gray-400">
            {{ formatDate(shipmentPrimaryTime(o.latest_shipment)) }}
          </span>
        </div>
        <div v-if="o.after_sale_summary?.latest_case_id" class="flex items-center flex-wrap gap-2 mb-4">
          <span v-if="o.after_sale_summary?.has_open_case" class="text-[11px] px-2 py-0.5 rounded-full bg-red-50 text-red-500">{{ $t('orderList.afterSaleInProgress') }}</span>
          <span class="text-xs text-gray-500">
            {{ $t('orderList.latestAfterSale') }}#{{ o.after_sale_summary.latest_case_id }}（{{ afterSaleSummaryStatusText(o.after_sale_summary) || '-' }}）
          </span>
        </div>

        <div class="flex-between">
          <span class="text-sm text-gray-500">{{ payLabel(o.payment_method) }}</span>
          <span class="text-lg font-bold text-gray-900">¥{{ money(o.amount_breakdown?.payable_amount ?? o.total_amount) }}</span>
        </div>

        <div class="flex justify-end mt-3 gap-2">
          <button class="btn-outline !px-5 text-xs" @click="goDetail(o.id)">{{ $t('orderList.viewDetail') }}</button>
          <button
            class="btn-outline !px-5 text-xs"
            v-if="o.after_sale_summary?.latest_case_id"
            @click="goAfterSaleDetail(o.after_sale_summary.latest_case_id)"
          >
            {{ $t('orderList.afterSaleProgress') }}
          </button>
          <button class="btn-primary !px-6 text-xs" v-if="o.status === 1" :disabled="actioningID === o.id" @click="pay(o)">{{ $t('orderList.pay') }}</button>
          <button class="btn-outline !px-5 text-xs" v-if="canReview(o) && hasUnreviewed(o)" :disabled="actioningID === o.id" @click="review(o, 'root')">{{ $t('orderList.review') }}</button>
          <button class="btn-outline !px-5 text-xs" v-if="canReview(o) && hasReviewed(o)" :disabled="actioningID === o.id" @click="review(o, 'append')">{{ $t('orderList.appendReview') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { get, post } from '@/api/request'
import { afterSaleStatusLabel, hasReshipShipment, orderStatusLabel, shipmentPrimaryTime, shipmentStatusLabel } from '@/utils/order-status'

const { t } = useI18n()
const router = useRouter()
const orders = ref<any[]>([])
const activeTab = ref(0)
const actioningID = ref(0)
const tabs = computed(() => [t('orderList.all'), t('orderList.unpaid'), t('orderList.unshipped'), t('orderList.unreceived'), t('orderList.completed')])
const statusValues = [0, 1, 2, 3, 4]

const statusColors: Record<number, string> = {
  1: 'bg-orange-50 text-orange-600',
  2: 'bg-red-50 text-red-600',
  3: 'bg-purple-50 text-purple-600',
  4: 'bg-green-50 text-green-600',
  5: 'bg-red-50 text-red-500',
}
const statusLabel = (s: number) => orderStatusLabel(s)
const statusColor = (s: number) => statusColors[s] || 'bg-gray-50 text-gray-400'
const payLabel = (m: string) => m === 'wechat' ? t('orderList.wechatPay') : m === 'alipay' ? t('orderList.alipay') : t('orderList.unpaidStatus')
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
  const data: any = await get('/api/v1/orders', { status: status || undefined })
  orders.value = data?.list || []
}

function goDetail(id: number) {
  router.push(`/orders/${id}`)
}

function goAfterSaleDetail(id: number) {
  if (!id) return
  router.push(`/after-sales/${id}`)
}

function canReview(order: any) {
  const status = Number(order?.status || 0)
  return status === 3 || status === 4
}

function hasReviewed(order: any) {
  const items = Array.isArray(order?.items) ? order.items : []
  return items.some((item: any) => !!item?.review?.id)
}

function hasUnreviewed(order: any) {
  const items = Array.isArray(order?.items) ? order.items : []
  return items.length > 0 && items.some((item: any) => !item?.review?.id)
}

async function pay(order: any) {
  const id = Number(order?.id || 0)
  if (!id || actioningID.value) return
  actioningID.value = id
  try {
    await post(`/api/v1/orders/${id}/pay`)
    await loadOrders()
  } catch (error: any) {
    alert(error?.message || t('orderReview.submitFailed'))
  } finally {
    actioningID.value = 0
  }
}

async function review(order: any, mode: 'root' | 'append' = 'root') {
  const id = Number(order?.id || 0)
  if (!id || actioningID.value) return
  router.push(`/orders/${id}/review?mode=${mode}`)
}

onMounted(loadOrders)
</script>
