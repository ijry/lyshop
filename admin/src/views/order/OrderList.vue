<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('order.list.title') }}</h2>
    </div>

    <div class="flex gap-2 mb-4 flex-wrap">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        @click="onTabChange(tab.value)"
        :class="activeStatus === tab.value ? 'bg-blue-700 text-white' : 'bg-white text-slate-600 border border-slate-200 hover:bg-slate-50'"
        class="px-4 py-2 rounded-xl text-sm transition"
      >
        {{ tab.label }}
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('order.list.orderId') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('order.list.productInfo') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('order.list.amountDetail') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('order.list.userPayment') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.status') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="o in orders" :key="o.id" class="align-top hover:bg-slate-50">
            <td class="px-4 py-3">
              <p class="font-mono text-xs text-slate-600">{{ o.order_no }}</p>
              <p class="text-xs text-slate-400 mt-1">{{ formatDate(o.created_at) }}</p>
              <p v-if="o.tracking_no" class="text-xs text-slate-400 mt-1">{{ $t('order.list.trackingNo') }}{{ o.tracking_no }}</p>
              <p v-if="o.after_sale_summary?.has_open_case" class="text-xs text-red-500 mt-1">{{ $t('order.list.afterSaleActive') }}</p>
              <p v-if="o.after_sale_summary?.latest_case_id" class="text-xs text-slate-500 mt-1">
                {{ $t('order.list.latestAfterSale') }}#{{ o.after_sale_summary.latest_case_id }}（{{ afterSaleSummaryStatusText(o.after_sale_summary) || '-' }}）
              </p>
              <div v-if="hasShipmentSummary(o)" class="flex items-center flex-wrap gap-1 mt-1">
                <span v-if="hasReship(o)" class="text-[11px] px-2 py-0.5 rounded-full bg-orange-50 text-orange-600">{{ $t('order.list.hasReship') }}</span>
                <template v-if="o.latest_shipment?.delivery_type === 'local'">
                  <span class="text-[11px] px-2 py-0.5 rounded-full bg-blue-50 text-blue-600">{{ $t('order.list.localDelivery') }}</span>
                  <span v-if="o.latest_shipment?.rider_name" class="text-xs text-slate-500">
                    {{ $t('order.list.rider') }}{{ o.latest_shipment.rider_name }}
                  </span>
                </template>
                <template v-else>
                  <span v-if="o.latest_shipment?.tracking_no" class="text-xs text-slate-500">
                    {{ $t('order.list.latestLogistics') }}{{ shipmentStatusText(o.latest_shipment?.logistics_status, o.latest_shipment?.logistics_status_label) }} · {{ o.latest_shipment?.tracking_no }}
                  </span>
                </template>
                <span v-if="shipmentPrimaryTime(o.latest_shipment)" class="text-xs text-slate-400">
                  {{ formatDate(shipmentPrimaryTime(o.latest_shipment)) }}
                </span>
              </div>
            </td>
            <td class="px-4 py-3">
              <div v-if="o.items?.length" class="space-y-2 min-w-[260px]">
                <div v-for="it in o.items.slice(0, 2)" :key="it.id" class="flex items-center gap-2">
                  <img :src="it.cover" class="w-10 h-10 rounded border border-slate-200 object-cover" />
                  <div class="min-w-0">
                    <p class="text-slate-700 truncate">{{ it.title }}</p>
                    <p class="text-xs text-slate-400">x{{ it.qty }}</p>
                  </div>
                </div>
                <p v-if="o.items.length > 2" class="text-xs text-slate-400">{{ $t('order.list.itemCount', { count: o.items.length }) }}</p>
              </div>
              <p v-else class="text-slate-400 text-xs">{{ $t('order.list.noItems') }}</p>
            </td>
            <td class="px-4 py-3 text-xs text-slate-600">
              <p>{{ $t('order.list.productAmount') }}¥{{ money(o.amount_breakdown?.goods_amount ?? o.goods_amount) }}</p>
              <p>{{ $t('order.list.discount') }}-¥{{ money(o.amount_breakdown?.discount_amount ?? o.discount_amount) }}</p>
              <p>{{ $t('order.list.freight') }}¥{{ money(o.amount_breakdown?.freight_amount ?? o.freight_amount) }}</p>
              <p class="text-slate-800 font-semibold mt-1">{{ $t('order.list.paidAmount') }}¥{{ money(o.amount_breakdown?.payable_amount ?? o.total_amount) }}</p>
            </td>
            <td class="px-4 py-3 text-xs text-slate-600">
              <p>{{ $t('order.list.userId') }}{{ o.user_id }}</p>
              <p class="mt-1">{{ payLabel(o.payment_method) }}</p>
            </td>
            <td class="px-4 py-3">
              <span :class="statusClass(o.status)" class="px-2 py-1 rounded-full text-xs">
                {{ statusLabel(o.status) }}
              </span>
            </td>
            <td class="px-4 py-3">
              <div class="flex flex-col gap-2 items-start">
                <button class="text-blue-600 hover:underline text-xs" @click="goDetail(o.id)">{{ $t('order.list.viewDetail') }}</button>
                <button v-if="o.status === 2" @click="goDetail(o.id)" class="text-emerald-600 hover:underline text-xs">{{ $t('order.list.ship') }}</button>
                <button
                  v-if="o.after_sale_summary?.latest_case_id"
                  @click="goAfterSaleDetail(o.after_sale_summary.latest_case_id)"
                  class="text-red-600 hover:underline text-xs"
                >
                  {{ $t('order.list.afterSaleProgress') }}
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="!orders.length">
            <td colspan="6" class="px-4 py-12 text-center text-slate-400">{{ $t('order.list.noOrder') }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getOrders } from '@/api/plugins'
import { afterSaleStatusLabel, hasReshipShipment, orderStatusLabel, shipmentPrimaryTime, shipmentStatusLabel } from '@/utils/order-status'

const { t } = useI18n()
const router = useRouter()
const orders = ref<any[]>([])
const activeStatus = ref(0)

const tabs = computed(() => [
  { label: t('order.list.all'), value: 0 },
  { label: t('orderStatus.pending'), value: 1 },
  { label: t('orderStatus.shipped'), value: 2 },
  { label: t('orderStatus.delivering'), value: 3 },
  { label: t('orderStatus.completed'), value: 4 },
  { label: t('orderStatus.afterSale'), value: 5 },
])

const statusColors: Record<number, string> = {
  1: 'bg-yellow-50 text-yellow-700',
  2: 'bg-blue-50 text-blue-700',
  3: 'bg-purple-50 text-purple-700',
  4: 'bg-green-50 text-green-700',
  5: 'bg-red-50 text-red-600',
}
const statusLabel = (s: number) => orderStatusLabel(s) || t('common.noData')
const statusClass = (s: number) => statusColors[s] || 'bg-slate-50 text-slate-400'
const payLabel = (m: string) => m === 'wechat' ? t('order.list.wechatPay') : m === 'alipay' ? t('order.list.alipay') : t('order.list.unpaid')
const money = (v: number) => Number(v || 0).toFixed(2)
const formatDate = (v?: string) => v ? String(v).slice(0, 19).replace('T', ' ') : '-'

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

function goDetail(id: number) {
  router.push(`/order/detail/${id}`)
}

function goAfterSaleDetail(id: number) {
  if (!id) return
  router.push(`/order/after-sale/detail/${id}`)
}

function onTabChange(status: number) {
  activeStatus.value = status
  loadOrders()
}

async function loadOrders() {
  const data: any = await getOrders({ status: activeStatus.value || undefined, page: 1, size: 20 })
  orders.value = data?.list || []
}

onMounted(loadOrders)
</script>
