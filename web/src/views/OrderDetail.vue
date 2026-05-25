<template>
  <div class="max-w-5xl mx-auto px-6 py-8" v-if="detail">
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">← {{ $t('orderDetail.back') }}</button>
      <h1 class="text-xl font-bold text-gray-900">{{ $t('orderDetail.title') }}</h1>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-[2fr_1fr] gap-6">
      <div class="space-y-4">
        <div class="card p-5">
          <h3 class="font-semibold text-gray-800 mb-4">{{ $t('orderDetail.productDetail') }}</h3>
          <div class="space-y-4">
            <div v-for="it in detail.items || []" :key="it.id" class="flex items-center gap-3">
              <img :src="it.cover" class="w-14 h-14 rounded-lg object-cover" />
              <div class="flex-1 min-w-0">
                <p class="text-sm text-gray-700 truncate">{{ it.title }}</p>
                <p class="text-xs text-gray-400">{{ $t('orderDetail.quantity') }} x{{ it.qty }}</p>
              </div>
              <p class="text-sm font-medium text-gray-800">¥{{ money(it.price) }}</p>
            </div>
          </div>
        </div>

        <div class="card p-5">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-semibold text-gray-800">{{ $t('orderDetail.logistics') }}</h3>
            <span v-if="detail.latest_shipment?.tracking_no" class="text-xs text-gray-400">{{ $t('orderDetail.latestTrackingNo') }}{{ detail.latest_shipment.tracking_no }}</span>
          </div>
          <div v-if="detail.shipments?.length" class="space-y-3">
            <div v-for="ship in detail.shipments" :key="ship.id" class="border border-gray-100 rounded-lg p-3">
              <p class="text-sm text-gray-700">{{ shipmentTitle(ship) }}</p>
              <p class="text-xs text-gray-500 mt-1">{{ ship.company || t('orderDetail.unknownCompany') }} · {{ ship.tracking_no || '-' }}</p>
              <p class="text-xs text-gray-500 mt-1">{{ $t('orderDetail.channel') }}{{ logisticsProviderLabel(ship.channel_provider) }}</p>
              <p class="text-xs text-gray-500 mt-1">{{ $t('orderDetail.status') }}{{ shipmentStatusText(ship.logistics_status, ship.logistics_status_label) }}</p>
              <p v-if="shipmentPrimaryTime(ship)" class="text-xs text-gray-500 mt-1">
                {{ shipmentTimeLabelText(ship) }}：{{ formatDate(shipmentPrimaryTime(ship)) }}
              </p>
              <p v-if="ship.after_sale_case_id" class="text-xs text-gray-500 mt-1">{{ $t('orderDetail.relatedAfterSale') }}#{{ ship.after_sale_case_id }}</p>
              <p v-if="ship.remark" class="text-xs text-gray-500 mt-1">{{ $t('orderDetail.remark') }}{{ ship.remark }}</p>
              <div v-if="tracksMap[ship.id]?.length" class="mt-3 border-t border-gray-100 pt-3 space-y-1">
                <p class="text-xs text-gray-500">{{ $t('orderDetail.trackingNodes') }}</p>
                <p v-for="track in tracksMap[ship.id]" :key="track.id" class="text-xs text-gray-500">
                  {{ formatDate(track.event_time) }} · {{ track.status_text }}<span v-if="track.location">（{{ track.location }}）</span>
                </p>
              </div>
            </div>
          </div>
          <p v-else class="text-sm text-gray-400">{{ $t('orderDetail.noLogistics') }}</p>
        </div>
      </div>

      <div class="space-y-4">
        <div class="card p-5">
          <h3 class="font-semibold text-gray-800 mb-3">{{ $t('orderDetail.orderInfo') }}</h3>
          <div class="space-y-2 text-sm text-gray-600">
            <p>{{ $t('orderDetail.orderNo') }}{{ detail.order_no }}</p>
            <p>{{ $t('orderDetail.payMethod') }}{{ payLabel(detail.payment_method) }}</p>
            <p>{{ $t('orderDetail.status') }}{{ statusLabel(detail.status) }}</p>
            <p>{{ $t('orderDetail.orderTime') }}{{ formatDate(detail.created_at) }}</p>
            <p v-if="detail.paid_at">{{ $t('orderDetail.payTime') }}{{ formatDate(detail.paid_at) }}</p>
            <p v-if="detail.tracking_no">{{ $t('orderDetail.trackingNo') }}{{ detail.tracking_no }}</p>
          </div>
        </div>

        <div class="card p-5">
          <h3 class="font-semibold text-gray-800 mb-3">{{ $t('orderDetail.priceSection') }}</h3>
          <div class="space-y-2 text-sm text-gray-600">
            <div class="flex justify-between"><span>{{ $t('orderDetail.productAmount') }}</span><span>¥{{ money(detail.amount_breakdown?.goods_amount ?? detail.goods_amount) }}</span></div>
            <div class="flex justify-between"><span>{{ $t('orderDetail.discountAmount') }}</span><span>-¥{{ money(detail.amount_breakdown?.discount_amount ?? detail.discount_amount) }}</span></div>
            <div class="flex justify-between"><span>{{ $t('orderDetail.shippingFee') }}</span><span>¥{{ money(detail.amount_breakdown?.freight_amount ?? detail.freight_amount) }}</span></div>
            <div class="flex justify-between font-semibold text-gray-800 pt-2 border-t border-gray-100"><span>{{ $t('orderDetail.paidAmount') }}</span><span>¥{{ money(detail.amount_breakdown?.payable_amount ?? detail.total_amount) }}</span></div>
          </div>
        </div>

        <div class="card p-5">
          <h3 class="font-semibold text-gray-800 mb-3">{{ $t('orderDetail.afterSaleSection') }}</h3>
          <div class="space-y-2 text-sm text-gray-600">
            <p>{{ $t('orderDetail.inProgress') }}{{ detail.after_sale_summary?.in_progress_count || 0 }}</p>
            <p v-if="detail.after_sale_summary?.latest_case_id">
              {{ $t('orderList.latestAfterSale') }}#{{ detail.after_sale_summary.latest_case_id }}（{{ afterSaleSummaryStatusText(detail.after_sale_summary) || '-' }}）
            </p>
          </div>
          <div class="flex gap-2 mt-4">
            <button
              v-if="detail.after_sale_summary?.latest_case_id"
              class="btn-outline !px-4 !py-2 text-xs"
              @click="goAfterSaleDetail(detail.after_sale_summary.latest_case_id)"
            >
              {{ $t('orderDetail.viewAfterSaleProgress') }}
            </button>
            <button
              v-if="detail.after_sale_summary?.can_apply !== false"
              class="btn-primary !px-4 !py-2 text-xs"
              @click="goAfterSaleApply"
            >
              {{ $t('orderDetail.applyAfterSale') }}
            </button>
          </div>
        </div>

        <div class="card p-5">
          <h3 class="font-semibold text-gray-800 mb-3">{{ $t('orderDetail.reviewSection') }}</h3>
          <div v-if="reviewItems.length" class="space-y-3">
            <div v-for="rv in reviewItems" :key="rv.id" class="border border-gray-100 rounded-lg p-3">
              <p class="text-sm text-gray-700">{{ rv.product_title }}</p>
              <p class="text-xs text-gray-400 mt-1">{{ $t('productDetail.product') }} {{ rv.product_score }} / {{ $t('productDetail.logistics') }} {{ rv.logistics_score }}</p>
              <p class="text-xs text-gray-500 mt-2">{{ rv.content || $t('orderDetail.noReview') }}</p>
              <div v-if="rv.images?.length" class="flex flex-wrap gap-2 mt-2">
                <img v-for="(img, idx) in rv.images" :key="img + idx" :src="img" class="w-14 h-14 rounded-md object-cover border border-gray-100" />
              </div>
              <div v-if="rv.appends?.length" class="mt-2 p-2 rounded bg-gray-50 space-y-2">
                <div v-for="ap in rv.appends" :key="ap.id">
                  <p class="text-xs text-gray-500">{{ $t('orderDetail.appendPrefix') }}{{ ap.content || $t('productDetail.imageOnlyAppend') }}</p>
                  <div v-if="ap.images?.length" class="flex flex-wrap gap-2 mt-1">
                    <img v-for="(img, idx) in ap.images" :key="img + idx" :src="img" class="w-12 h-12 rounded object-cover border border-gray-100" />
                  </div>
                </div>
              </div>
              <div class="flex justify-end mt-3">
                <button class="btn-outline !px-4 !py-2 text-xs" @click="goReview('append', rv.order_item_id)">{{ $t('orderList.appendReview') }}</button>
              </div>
            </div>
          </div>
          <p v-else class="text-sm text-gray-400">{{ $t('productDetail.noReview') }}</p>
          <div class="flex gap-2 mt-4" v-if="hasUnreviewed">
            <button class="btn-primary !px-4 !py-2 text-xs" @click="goReview('root')">{{ $t('orderList.review') }}</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { get } from '@/api/request'
import {
  afterSaleStatusLabel,
  orderStatusLabel,
  logisticsProviderLabel,
  shipmentPrimaryTime,
  shipmentStatusLabel,
  shipmentTimeLabel,
  shipmentTitle,
} from '@/utils/order-status'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const detail = ref<any>(null)
const tracksMap = ref<Record<number, any[]>>({})
const reviewItems = ref<any[]>([])
const hasReviewed = ref(false)
const hasUnreviewed = ref(false)

const statusLabel = (s: number) => orderStatusLabel(s) || String(s)
const payLabel = (m: string) => (m === 'wechat' ? t('orderList.wechatPay') : m === 'alipay' ? t('orderList.alipay') : t('orderList.unpaidStatus'))
const money = (v: number) => Number(v || 0).toFixed(2)
const formatDate = (v?: string) => (v ? String(v).slice(0, 19).replace('T', ' ') : '-')

function shipmentStatusText(status: string, label?: string) {
  const mapped = String(label || '').trim()
  if (mapped) return mapped
  return shipmentStatusLabel(status)
}

function shipmentTimeLabelText(ship: any) {
  return shipmentTimeLabel(ship)
}

function afterSaleSummaryStatusText(summary: any) {
  const label = String(summary?.latest_status_label || '').trim()
  if (label) return label
  return afterSaleStatusLabel(summary?.latest_status)
}

function refreshReviewFlags() {
  const items = Array.isArray(detail.value?.items) ? detail.value.items : []
  hasReviewed.value = items.some((item: any) => !!item?.review?.id)
  hasUnreviewed.value = items.length > 0 && items.some((item: any) => !item?.review?.id)
}

function goReview(mode: 'root' | 'append', orderItemID?: number) {
  if (!detail.value?.id) return
  const itemQuery = mode === 'append' && Number(orderItemID || 0) > 0 ? `&item_id=${Number(orderItemID)}` : ''
  router.push(`/orders/${detail.value.id}/review?mode=${mode}${itemQuery}`)
}

function goAfterSaleApply() {
  if (!detail.value?.id) return
  router.push(`/orders/${detail.value.id}/after-sale/apply`)
}

function goAfterSaleDetail(id: number) {
  if (!id) return
  router.push(`/after-sales/${id}`)
}

onMounted(async () => {
  detail.value = await get<any>(`/api/v1/orders/${route.params.id}`)
  tracksMap.value = {}
  for (const shipment of detail.value?.shipments || []) {
    if (!shipment?.id) continue
    const rows = await get<any[]>(`/api/v1/orders/${route.params.id}/shipments/${shipment.id}/tracks`)
    tracksMap.value[Number(shipment.id)] = Array.isArray(rows) ? rows : []
  }
  reviewItems.value = []
  for (const item of detail.value?.items || []) {
    if (item.review) {
      reviewItems.value.push({
        ...item.review,
        product_title: item.title,
      })
    }
  }
  refreshReviewFlags()
})
</script>
