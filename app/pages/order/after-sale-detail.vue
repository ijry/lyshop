<template>
  <view class="min-h-screen bg-gray-50 pb-40rpx">
    <u-navbar :title="$t('afterSaleDetail.title')" :placeholder="true" />

    <view v-if="detail.id" class="p-20rpx">
      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block">{{ $t('afterSaleDetail.afterSaleNo') }}{{ detail.case_no }}</text>
        <text class="text-24rpx text-gray-600 block mt-8rpx">{{ $t('afterSaleDetail.orderId') }}{{ detail.order_id }}</text>
        <text class="text-24rpx text-gray-600 block mt-8rpx">{{ $t('afterSaleDetail.type') }}{{ caseTypeText(detail.case_type, detail.case_type_label) }}</text>
        <text class="text-24rpx text-gray-600 block mt-8rpx">{{ $t('afterSaleDetail.status') }}{{ statusText(detail.status, detail.status_label) }}</text>
        <text class="text-24rpx text-gray-600 block mt-8rpx">{{ $t('afterSaleDetail.reason') }}{{ detail.reason }}</text>
        <text v-if="detail.apply_content" class="text-24rpx text-gray-600 block mt-8rpx">{{ $t('afterSaleDetail.description') }}{{ detail.apply_content }}</text>
        <view v-if="detail.apply_images?.length" class="flex flex-wrap gap-12rpx mt-16rpx">
          <image
            v-for="(img, idx) in detail.apply_images"
            :key="img + idx"
            :src="img"
            mode="aspectFill"
            style="width: 150rpx; height: 150rpx; border-radius: 16rpx;"
            @click="previewImages(detail.apply_images, idx)"
          />
        </view>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">{{ $t('afterSaleDetail.logistics') }}</text>
        <view v-if="detail.shipments?.length" class="space-y-12rpx">
          <view v-for="ship in detail.shipments" :key="ship.id" class="border border-gray-100 rounded-16rpx p-16rpx">
            <text class="text-24rpx text-gray-700 block">{{ shipmentTitle(ship) }}</text>
            <text class="text-22rpx text-gray-500 block mt-6rpx">{{ ship.company || $t('afterSaleDetail.unknownCompany') }} · {{ ship.tracking_no || '-' }}</text>
            <text class="text-22rpx text-gray-500 block mt-6rpx">{{ $t('afterSaleDetail.shipmentStatus') }}{{ shipmentStatusText(ship.logistics_status, ship.logistics_status_label) }}</text>
            <text v-if="shipmentPrimaryTime(ship)" class="text-22rpx text-gray-500 block mt-6rpx">
              {{ shipmentTimeLabel(ship) }}：{{ formatDate(shipmentPrimaryTime(ship)) }}
            </text>
            <text v-if="ship.after_sale_case_id" class="text-22rpx text-gray-500 block mt-6rpx">{{ $t('afterSaleDetail.relatedAfterSale') }}{{ ship.after_sale_case_id }}</text>
            <text v-if="ship.remark" class="text-22rpx text-gray-500 block mt-6rpx">{{ $t('afterSaleDetail.remark') }}{{ ship.remark }}</text>
          </view>
        </view>
        <text v-else class="text-24rpx text-gray-400">{{ $t('afterSaleDetail.noLogistics') }}</text>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">{{ $t('afterSaleDetail.progressLog') }}</text>
        <view v-if="detail.logs?.length" class="space-y-12rpx">
          <view v-for="log in detail.logs" :key="log.id" class="border border-gray-100 rounded-16rpx p-16rpx">
            <text class="text-24rpx text-gray-700 block">{{ actionText(log.action, log.action_label) }}：{{ statusLabelOrDash(log.from_status, log.from_status_label) }} → {{ statusText(log.to_status, log.to_status_label) }}</text>
            <text class="text-22rpx text-gray-500 block mt-6rpx">{{ log.content || '-' }}</text>
            <text class="text-22rpx text-gray-400 block mt-6rpx">{{ formatDate(log.created_at) }}</text>
          </view>
        </view>
        <text v-else class="text-24rpx text-gray-400">{{ $t('afterSaleDetail.noStatusLog') }}</text>
      </view>

      <view v-if="detail.status === 'approved_wait_user_return'" class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">{{ $t('afterSaleDetail.fillReturnLogistics') }}</text>
        <view class="mb-12rpx">
          <u-input v-model="returnForm.company" border="surround" :placeholder="$t('afterSaleDetail.expressCompany')" />
        </view>
        <view class="mb-12rpx">
          <u-input v-model="returnForm.tracking_no" border="surround" :placeholder="$t('afterSaleDetail.trackingNo')" />
        </view>
        <view class="mb-16rpx">
          <u-textarea v-model="returnForm.remark" :placeholder="$t('afterSaleDetail.remarkOptional')" :auto-height="true" maxlength="200" />
        </view>
        <u-button type="primary" shape="circle" :loading="submittingReturn" :text="submittingReturn ? $t('afterSaleDetail.submitting') : $t('afterSaleDetail.submitReturnLogistics')" @click="submitReturnShipment" />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { get, post } from '@/utils/request'
import {
  afterSaleStatusLabel,
  shipmentPrimaryTime,
  shipmentStatusLabel,
  shipmentTimeLabel,
  shipmentTitle,
} from '@/utils/order-status'

const { t } = useI18n()

const detail = ref<any>({})
const caseID = ref(0)
const submittingReturn = ref(false)
const returnForm = ref<any>({
  company: '',
  tracking_no: '',
  remark: '',
})

const statusLabel = (status: string) => afterSaleStatusLabel(status)
const formatDate = (v?: string) => (v ? String(v).slice(0, 19).replace('T', ' ') : '-')

const caseTypeLabels = computed(() => ({
  return: t('afterSaleDetail.returnGoods'),
  exchange: t('afterSaleDetail.exchangeGoods'),
}))

const actionLabels = computed(() => ({
  apply: t('afterSaleDetail.submitApplication'),
  audit: t('afterSaleDetail.review'),
  return_ship: t('afterSaleDetail.returnLogistics'),
  receive: t('afterSaleDetail.confirmReceipt'),
  refund: t('afterSaleDetail.refund'),
  reship: t('afterSaleDetail.reship'),
  complete: t('afterSaleDetail.close'),
  close: t('afterSaleDetail.closed'),
}))

function actionLabel(action: string) {
  return (actionLabels.value as any)[String(action || '')] || action || '-'
}

function caseTypeText(caseType: string, label?: string) {
  const mapped = String(label || '').trim()
  if (mapped) return mapped
  const value = String(caseType || '')
  return (caseTypeLabels.value as any)[value] || value
}

function statusText(status: string, label?: string) {
  const mapped = String(label || '').trim()
  if (mapped) return mapped
  return statusLabel(status)
}

function shipmentStatusText(status: string, label?: string) {
  const mapped = String(label || '').trim()
  if (mapped) return mapped
  return shipmentStatusLabel(status)
}

function actionText(action: string, label?: string) {
  const mapped = String(label || '').trim()
  if (mapped) return mapped
  return actionLabel(action)
}

function statusLabelOrDash(status: string, label?: string) {
  const value = String(label || status || '')
  if (!value) return '-'
  return statusText(status, label)
}

function toast(msg: string) {
  uni.showToast({ title: msg, icon: 'none' })
}

function readCaseID() {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  return Number(query.id || 0)
}

function previewImages(urls: string[], index: number) {
  if (!urls.length) return
  uni.previewImage({ urls, current: urls[index] || urls[0] })
}

async function load() {
  detail.value = await get<any>(`/api/v1/after-sales/${caseID.value}`)
}

async function submitReturnShipment() {
  if (submittingReturn.value) return
  const trackingNo = String(returnForm.value.tracking_no || '').trim()
  if (!trackingNo) {
    toast(t('afterSaleDetail.trackingNoRequired'))
    return
  }
  submittingReturn.value = true
  try {
    await post(`/api/v1/after-sales/${caseID.value}/return-shipments`, {
      company: String(returnForm.value.company || ''),
      tracking_no: trackingNo,
      remark: String(returnForm.value.remark || ''),
    })
    uni.showToast({ title: t('afterSaleDetail.submitSuccess'), icon: 'success' })
    returnForm.value = { company: '', tracking_no: '', remark: '' }
    await load()
  } catch (error: any) {
    toast(error?.message || t('afterSaleDetail.submitFailed'))
  } finally {
    submittingReturn.value = false
  }
}

onMounted(async () => {
  caseID.value = readCaseID()
  await load()
})
</script>
