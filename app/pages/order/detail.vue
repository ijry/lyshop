<template>
  <view class="min-h-screen bg-gray-50 pb-40rpx">
    <view v-if="detail.id" class="p-24rpx">
      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <view class="flex items-center justify-between mb-16rpx">
          <text class="text-28rpx font-600 text-gray-800">{{ $t('orderDetail.title') }} {{ detail.order_no }}</text>
          <text class="text-24rpx" :class="statusColor(detail.status)">{{ statusLabel(detail.status) }}</text>
        </view>
        <view class="space-y-12rpx text-24rpx text-gray-600">
          <view>{{ $t('orderDetail.paymentMethod') }}{{ payLabel(detail.payment_method) }}</view>
          <view>{{ $t('orderDetail.orderTime') }}{{ formatDate(detail.created_at) }}</view>
          <view v-if="detail.paid_at">{{ $t('orderDetail.payTime') }}{{ formatDate(detail.paid_at) }}</view>
          <view v-if="detail.tracking_no">{{ $t('orderDetail.trackingNo') }}{{ detail.tracking_no }}</view>
        </view>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">{{ $t('orderDetail.productDetail') }}</text>
        <view v-if="detail.items?.length" class="space-y-16rpx">
          <view v-for="it in detail.items" :key="it.id" class="flex gap-16rpx">
            <image :src="it.cover" mode="aspectFill" style="width: 120rpx; height: 120rpx; border-radius: 16rpx;" />
            <view class="flex-1 min-w-0">
              <text class="text-26rpx text-gray-800 block line-clamp-2">{{ it.title }}</text>
              <text class="text-22rpx text-gray-400 block mt-8rpx">{{ $t('orderDetail.quantity') }} x{{ it.qty }}</text>
              <text class="text-24rpx text-gray-600 block mt-8rpx">¥{{ money(it.price) }}</text>
            </view>
          </view>
        </view>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">{{ $t('orderDetail.priceBreakdown') }}</text>
        <view class="space-y-12rpx text-24rpx text-gray-600">
          <view class="flex justify-between"><text>{{ $t('orderDetail.productAmount') }}</text><text>¥{{ money(detail.amount_breakdown?.goods_amount ?? detail.goods_amount) }}</text></view>
          <view class="flex justify-between"><text>{{ $t('orderDetail.discountAmount') }}</text><text>-¥{{ money(detail.amount_breakdown?.discount_amount ?? detail.discount_amount) }}</text></view>
          <view class="flex justify-between"><text>{{ $t('orderDetail.shipping') }}</text><text>¥{{ money(detail.amount_breakdown?.freight_amount ?? detail.freight_amount) }}</text></view>
          <view class="flex justify-between text-gray-800 font-600 pt-12rpx border-t border-gray-100"><text>{{ $t('orderDetail.actualPaid') }}</text><text>¥{{ money(detail.amount_breakdown?.payable_amount ?? detail.total_amount) }}</text></view>
        </view>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <view class="flex items-center justify-between mb-16rpx">
          <text class="text-28rpx font-600 text-gray-800">{{ $t('orderDetail.logistics') }}</text>
          <text v-if="detail.latest_shipment?.tracking_no" class="text-22rpx text-gray-400">{{ $t('orderDetail.latest') }}{{ detail.latest_shipment.tracking_no }}</text>
        </view>
        <view v-if="detail.shipments?.length" class="space-y-12rpx">
          <view v-for="ship in detail.shipments" :key="ship.id" class="border border-gray-100 rounded-16rpx p-16rpx">
            <text class="text-24rpx text-gray-700 block">{{ shipmentTitle(ship) }}</text>
            <template v-if="ship.delivery_type === 'local'">
              <text class="text-22rpx text-gray-500 block mt-6rpx">{{ $t('orderDetail.rider') }}{{ ship.rider_name || '-' }}</text>
              <text class="text-22rpx text-gray-500 block mt-6rpx">{{ $t('orderDetail.riderPhone') }}{{ ship.rider_phone || '-' }}</text>
            </template>
            <template v-else>
              <text class="text-22rpx text-gray-500 block mt-6rpx">{{ ship.company || $t('afterSaleDetail.unknownCompany') }} · {{ ship.tracking_no || '-' }}</text>
              <text class="text-22rpx text-gray-500 block mt-6rpx">{{ $t('orderDetail.channel') }}{{ logisticsProviderLabel(ship.channel_provider) }}</text>
            </template>
            <text class="text-22rpx text-gray-500 block mt-6rpx">{{ $t('orderDetail.status') }}{{ shipmentStatusText(ship.logistics_status, ship.logistics_status_label) }}</text>
            <text v-if="shipmentPrimaryTime(ship)" class="text-22rpx text-gray-500 block mt-6rpx">
              {{ shipmentTimeLabel(ship) }}：{{ formatDate(shipmentPrimaryTime(ship)) }}
            </text>
            <text v-if="ship.after_sale_case_id" class="text-22rpx text-gray-500 block mt-6rpx">{{ $t('orderDetail.relatedAfterSale') }}{{ ship.after_sale_case_id }}</text>
            <text v-if="ship.remark" class="text-22rpx text-gray-500 block mt-6rpx">{{ $t('orderDetail.remark') }}{{ ship.remark }}</text>
            <view v-if="ship.delivery_type !== 'local' && tracksMap[ship.id]?.length" class="mt-12rpx border-t border-gray-100 pt-12rpx">
              <text class="text-22rpx text-gray-500 block mb-8rpx">{{ $t('orderDetail.trackingNode') }}</text>
              <text v-for="track in tracksMap[ship.id]" :key="track.id" class="text-22rpx text-gray-500 block mt-6rpx">
                {{ formatDate(track.event_time) }} · {{ track.status_text }}<text v-if="track.location">（{{ track.location }}）</text>
              </text>
            </view>
          </view>
        </view>
        <text v-else class="text-24rpx text-gray-400">{{ $t('orderDetail.noLogistics') }}</text>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">{{ $t('orderDetail.afterSaleService') }}</text>
        <text class="text-24rpx text-gray-600 block">{{ $t('orderDetail.inProgress') }}{{ detail.after_sale_summary?.in_progress_count || 0 }}</text>
        <text v-if="detail.after_sale_summary?.latest_case_id" class="text-24rpx text-gray-600 block mt-8rpx">
          {{ $t('orderDetail.recentAfterSale') }}{{ detail.after_sale_summary.latest_case_id }}（{{ afterSaleSummaryStatusText(detail.after_sale_summary) || '-' }}）
        </text>
        <view class="flex gap-16rpx mt-20rpx">
          <u-button
            v-if="detail.after_sale_summary?.latest_case_id"
            size="mini"
            :text="$t('orderDetail.viewAfterSaleProgress')"
            shape="circle"
            plain
            @click="goAfterSaleDetail(detail.after_sale_summary.latest_case_id)"
          />
          <u-button
            v-if="detail.after_sale_summary?.can_apply !== false"
            size="mini"
            type="warning"
            :text="$t('orderDetail.applyAfterSale')"
            shape="circle"
            plain
            @click="goAfterSaleApply"
          />
        </view>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">{{ $t('orderDetail.reviewInfo') }}</text>
        <view v-if="reviewItems.length" class="space-y-14rpx">
          <view v-for="rv in reviewItems" :key="rv.id" class="border border-gray-100 rounded-16rpx p-16rpx">
            <text class="text-24rpx text-gray-700 block">{{ rv.product_title }}</text>
            <text class="text-22rpx text-gray-400 block mt-6rpx">{{ $t('orderDetail.product') }} {{ rv.product_score }} / {{ $t('orderDetail.logisticsTab') }} {{ rv.logistics_score }}</text>
            <text class="text-24rpx text-gray-600 block mt-10rpx">{{ rv.content || $t('orderDetail.noReviewText') }}</text>
            <view v-if="rv.images?.length" class="flex flex-wrap gap-10rpx mt-12rpx">
              <image
                v-for="(img, idx) in rv.images"
                :key="img + idx"
                :src="img"
                mode="aspectFill"
                style="width: 120rpx; height: 120rpx; border-radius: 12rpx;"
              />
            </view>
            <view v-if="rv.appends?.length" class="mt-12rpx bg-gray-50 rounded-12rpx p-12rpx">
              <view v-for="ap in rv.appends" :key="ap.id" class="mb-10rpx last:mb-0">
                <text class="text-22rpx text-gray-500 block">{{ $t('orderDetail.append') }}{{ ap.content || $t('orderDetail.imageOnlyAppend') }}</text>
                <view v-if="ap.images?.length" class="flex flex-wrap gap-10rpx mt-8rpx">
                  <image
                    v-for="(img, idx) in ap.images"
                    :key="img + idx"
                    :src="img"
                    mode="aspectFill"
                    style="width: 108rpx; height: 108rpx; border-radius: 10rpx;"
                  />
                </view>
              </view>
            </view>
            <view class="flex justify-end mt-12rpx">
              <u-button size="mini" type="warning" :text="$t('orderReview.appendReview')" shape="circle" plain @click="goReview('append', rv.order_item_id)" />
            </view>
          </view>
        </view>
        <text v-else class="text-24rpx text-gray-400">{{ $t('orderDetail.noReview') }}</text>
        <view class="flex gap-16rpx mt-20rpx">
          <u-button
            v-if="hasUnreviewed"
            size="small"
            type="success"
            shape="circle"
            :text="$t('orderDetail.goReview')"
            plain
            @click="goReview('root')"
          />
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { get } from '@/utils/request'
import {
  afterSaleStatusLabel,
  logisticsProviderLabel,
  orderStatusLabel,
  shipmentPrimaryTime,
  shipmentStatusLabel,
  shipmentTimeLabel,
  shipmentTitle,
} from '@/utils/order-status'

const { t } = useI18n()

const detail = ref<any>({})
const tracksMap = ref<Record<number, any[]>>({})
const reviewItems = ref<any[]>([])
const statusColors: Record<number, string> = { 1: 'text-orange-500', 2: 'text-blue-500', 3: 'text-purple-500', 4: 'text-green-500', 5: 'text-red-500' }
const statusLabel = (s: number) => orderStatusLabel(s)
const statusColor = (s: number) => statusColors[s] || 'text-gray-400'
const payLabel = (m: string) => m === 'wechat' ? t('orderDetail.wechatPay') : m === 'alipay' ? t('orderDetail.alipayPay') : t('orderDetail.unpaid')
const money = (v: number) => Number(v || 0).toFixed(2)
const formatDate = (v?: string) => v ? String(v).slice(0, 19).replace('T', ' ') : '-'
const hasReviewed = ref(false)
const hasUnreviewed = ref(false)

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

function refreshReviewFlags() {
  const items = Array.isArray(detail.value?.items) ? detail.value.items : []
  hasReviewed.value = items.some((item: any) => !!item?.review?.id)
  hasUnreviewed.value = items.length > 0 && items.some((item: any) => !item?.review?.id)
}

function goReview(mode: 'root' | 'append', orderItemID?: number) {
  if (!detail.value?.id) return
  const itemQuery = mode === 'append' && Number(orderItemID || 0) > 0 ? `&item_id=${Number(orderItemID)}` : ''
  uni.navigateTo({ url: `/pages/order/review?id=${detail.value.id}&mode=${mode}${itemQuery}` })
}

function goAfterSaleApply() {
  if (!detail.value?.id) return
  uni.navigateTo({ url: `/pages/order/after-sale-apply?id=${detail.value.id}` })
}

function goAfterSaleDetail(id: number) {
  if (!id) return
  uni.navigateTo({ url: `/pages/order/after-sale-detail?id=${id}` })
}

onMounted(async () => {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  detail.value = await get<any>(`/api/v1/orders/${query.id}`)
  tracksMap.value = {}
  for (const shipment of detail.value?.shipments || []) {
    if (!shipment?.id || shipment.delivery_type === 'local') continue
    const rows = await get<any[]>(`/api/v1/orders/${query.id}/shipments/${shipment.id}/tracks`)
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
