<template>
  <div class="max-w-4xl mx-auto px-6 py-8" v-if="detail">
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">← {{ $t('orderDetail.back') }}</button>
      <h1 class="text-xl font-bold text-gray-900">{{ $t('afterSaleDetail.title') }}</h1>
    </div>

    <div class="card p-5 mb-4">
      <p class="text-sm text-gray-600">{{ $t('afterSaleDetail.afterSaleNo') }}<span class="font-mono">{{ detail.case_no }}</span></p>
      <p class="text-sm text-gray-600 mt-1">{{ $t('afterSaleDetail.orderId') }}{{ detail.order_id }}</p>
      <p class="text-sm text-gray-600 mt-1">{{ $t('afterSaleDetail.type') }}{{ caseTypeText(detail.case_type, detail.case_type_label) }}</p>
      <p class="text-sm text-gray-600 mt-1">{{ $t('afterSaleDetail.status') }}{{ statusText(detail.status, detail.status_label) }}</p>
      <p class="text-sm text-gray-600 mt-1">{{ $t('afterSaleDetail.reason') }}{{ detail.reason }}</p>
      <p v-if="detail.apply_content" class="text-sm text-gray-600 mt-1">{{ $t('afterSaleDetail.description') }}{{ detail.apply_content }}</p>
      <div v-if="detail.apply_images?.length" class="flex flex-wrap gap-2 mt-3">
        <img
          v-for="(img, idx) in detail.apply_images"
          :key="img + idx"
          :src="img"
          class="w-20 h-20 rounded-lg object-cover border border-gray-100 cursor-pointer"
          @click="previewImage(img)"
        />
      </div>
    </div>

    <div class="card p-5 mb-4">
      <h3 class="text-sm font-semibold text-gray-800 mb-3">{{ $t('afterSaleDetail.logistics') }}</h3>
      <div v-if="detail.shipments?.length" class="space-y-3">
        <div v-for="ship in detail.shipments" :key="ship.id" class="border border-gray-100 rounded-lg p-3">
          <p class="text-sm text-gray-700">{{ shipmentTitle(ship) }}</p>
          <p class="text-xs text-gray-500 mt-1">{{ ship.company || t('afterSaleDetail.unknownCompany') }} · {{ ship.tracking_no || '-' }}</p>
          <p class="text-xs text-gray-500 mt-1">{{ $t('afterSaleDetail.status') }}{{ shipmentStatusText(ship.logistics_status, ship.logistics_status_label) }}</p>
          <p v-if="shipmentPrimaryTime(ship)" class="text-xs text-gray-500 mt-1">
            {{ shipmentTimeLabelText(ship) }}：{{ formatDate(shipmentPrimaryTime(ship)) }}
          </p>
          <p v-if="ship.after_sale_case_id" class="text-xs text-gray-500 mt-1">{{ $t('afterSaleDetail.relatedAfterSale') }}#{{ ship.after_sale_case_id }}</p>
          <p v-if="ship.remark" class="text-xs text-gray-500 mt-1">{{ $t('afterSaleDetail.remark') }}{{ ship.remark }}</p>
        </div>
      </div>
      <p v-else class="text-sm text-gray-400">{{ $t('afterSaleDetail.noLogistics') }}</p>
    </div>

    <div class="card p-5 mb-4">
      <h3 class="text-sm font-semibold text-gray-800 mb-3">{{ $t('afterSaleDetail.progressLog') }}</h3>
      <div v-if="detail.logs?.length" class="space-y-3">
        <div v-for="log in detail.logs" :key="log.id" class="border border-gray-100 rounded-lg p-3">
          <p class="text-sm text-gray-700">{{ actionText(log.action, log.action_label) }}：{{ statusLabelOrDash(log.from_status, log.from_status_label) }} → {{ statusText(log.to_status, log.to_status_label) }}</p>
          <p class="text-xs text-gray-500 mt-1">{{ log.content || '-' }}</p>
          <p class="text-xs text-gray-400 mt-1">{{ formatDate(log.created_at) }}</p>
        </div>
      </div>
      <p v-else class="text-sm text-gray-400">{{ $t('afterSaleDetail.noLog') }}</p>
    </div>

    <div v-if="detail.status === 'approved_wait_user_return'" class="card p-5">
      <h3 class="text-sm font-semibold text-gray-800 mb-3">{{ $t('afterSaleDetail.returnShipping') }}</h3>
      <div class="space-y-3">
        <input v-model="returnForm.company" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm" :placeholder="$t('afterSaleDetail.expressCompany')" />
        <input v-model="returnForm.tracking_no" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm" :placeholder="$t('afterSaleDetail.trackingNo')" />
        <textarea v-model="returnForm.remark" class="w-full min-h-[72px] border border-gray-200 rounded-lg px-3 py-2 text-sm" :placeholder="$t('afterSaleDetail.remarkOptional')" />
      </div>
      <button class="btn-primary w-full !py-3 mt-4" :disabled="submittingReturn" @click="submitReturnShipment">
        {{ submittingReturn ? $t('afterSaleDetail.submitting') : $t('afterSaleDetail.submitReturn') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { get, post } from '@/api/request'
import {
  afterSaleStatusLabel,
  shipmentPrimaryTime,
  shipmentStatusLabel,
  shipmentTimeLabel,
  shipmentTitle,
} from '@/utils/order-status'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const detail = ref<any>(null)
const submittingReturn = ref(false)
const returnForm = ref<any>({
  company: '',
  tracking_no: '',
  remark: '',
})
const statusLabelFn = (status: string) => afterSaleStatusLabel(status)
const formatDate = (v?: string) => (v ? String(v).slice(0, 19).replace('T', ' ') : '-')
const caseTypeLabels = () => ({
  return: t('afterSaleDetail.typeRefund'),
  exchange: t('afterSaleDetail.typeExchange'),
}) as Record<string, string>
const actionLabels = () => ({
  apply: t('afterSaleDetail.actionApply'),
  audit: t('afterSaleDetail.actionReview'),
  return_ship: t('afterSaleDetail.actionReturn'),
  receive: t('afterSaleDetail.actionConfirm'),
  refund: t('afterSaleDetail.actionRefund'),
  reship: t('afterSaleDetail.actionReship'),
  complete: t('afterSaleDetail.actionComplete'),
  close: t('afterSaleDetail.actionClose'),
}) as Record<string, string>

function actionLabel(action: string) {
  return actionLabels()[String(action || '')] || action || '-'
}

function caseTypeText(caseType: string, label?: string) {
  const mapped = String(label || '').trim()
  if (mapped) return mapped
  const value = String(caseType || '')
  return caseTypeLabels()[value] || value
}

function statusText(status: string, label?: string) {
  const mapped = String(label || '').trim()
  if (mapped) return mapped
  return statusLabelFn(status)
}

function shipmentStatusText(status: string, label?: string) {
  const mapped = String(label || '').trim()
  if (mapped) return mapped
  return shipmentStatusLabel(status)
}

function shipmentTimeLabelText(ship: any) {
  return shipmentTimeLabel(ship)
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

function alertMsg(msg: string) {
  alert(msg)
}

async function load() {
  detail.value = await get<any>(`/api/v1/after-sales/${route.params.id}`)
}

function previewImage(url: string) {
  if (!url) return
  window.open(url, '_blank')
}

async function submitReturnShipment() {
  if (submittingReturn.value) return
  const trackingNo = String(returnForm.value.tracking_no || '').trim()
  if (!trackingNo) {
    alertMsg(t('afterSaleDetail.trackingNoRequired'))
    return
  }
  submittingReturn.value = true
  try {
    await post(`/api/v1/after-sales/${route.params.id}/return-shipments`, {
      company: String(returnForm.value.company || ''),
      tracking_no: trackingNo,
      remark: String(returnForm.value.remark || ''),
    })
    alertMsg(t('afterSaleDetail.submitSuccess'))
    returnForm.value = { company: '', tracking_no: '', remark: '' }
    await load()
  } catch (error: any) {
    alertMsg(error?.message || t('afterSaleDetail.submitFailed'))
  } finally {
    submittingReturn.value = false
  }
}

onMounted(load)
</script>
