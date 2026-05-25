<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">{{ $t('common.back') }}</button>
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('afterSale.detail.title') }}</h2>
    </div>

    <div v-if="detail" class="grid grid-cols-1 xl:grid-cols-[2fr_1fr] gap-6">
      <div class="space-y-4">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-4">{{ $t('afterSale.detail.info') }}</h3>
          <div class="space-y-2 text-sm text-slate-600">
            <p>{{ $t('afterSale.detail.caseNo') }}<span class="font-mono">{{ detail.case_no }}</span></p>
            <p>{{ $t('afterSale.detail.orderId') }}{{ detail.order_id }}</p>
            <p>{{ $t('afterSale.detail.type') }}{{ typeLabel(detail.case_type) }}</p>
            <p>{{ $t('afterSale.detail.status') }}{{ statusText(detail.status, detail.status_label) }}</p>
            <p>{{ $t('afterSale.detail.reason') }}{{ detail.reason }}</p>
            <p>{{ $t('afterSale.detail.note') }}{{ detail.apply_content || '-' }}</p>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-4">{{ $t('afterSale.detail.statusLog') }}</h3>
          <div v-if="detail.logs?.length" class="space-y-3">
            <div v-for="log in detail.logs" :key="log.id" class="border border-slate-100 rounded-lg p-3">
              <p class="text-sm text-slate-700">
                {{ actionLabel(log.action) }}：{{ statusLabelOrDash(log.from_status, log.from_status_label) }} → {{ statusText(log.to_status, log.to_status_label) }}
              </p>
              <p class="text-xs text-slate-400 mt-1">{{ log.content || '-' }}</p>
              <p class="text-xs text-slate-400 mt-1">{{ formatDate(log.created_at) }}</p>
            </div>
          </div>
          <p v-else class="text-slate-400 text-sm">{{ $t('common.noData') }}</p>
        </div>
      </div>

      <div class="space-y-4">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-3">{{ $t('common.action') }}</h3>
          <div class="space-y-3">
            <button class="w-full px-4 py-2 rounded-lg bg-blue-700 text-white text-sm" @click="audit(true)">{{ $t('afterSale.detail.approve') }}</button>
            <button class="w-full px-4 py-2 rounded-lg bg-slate-100 text-slate-700 text-sm" @click="audit(false)">{{ $t('afterSale.detail.reject') }}</button>
            <button class="w-full px-4 py-2 rounded-lg bg-emerald-600 text-white text-sm" @click="receive">{{ $t('afterSale.detail.confirmReceive') }}</button>
            <button class="w-full px-4 py-2 rounded-lg bg-orange-600 text-white text-sm" @click="refund">{{ $t('afterSale.detail.registerRefund') }}</button>
            <button class="w-full px-4 py-2 rounded-lg bg-purple-600 text-white text-sm" @click="complete">{{ $t('afterSale.detail.complete') }}</button>
            <button class="w-full px-4 py-2 rounded-lg bg-red-600 text-white text-sm" @click="close">{{ $t('afterSale.detail.closeCase') }}</button>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-3">{{ $t('afterSale.detail.logistics') }}</h3>
          <div v-if="detail.shipments?.length" class="space-y-2">
            <div v-for="ship in detail.shipments" :key="ship.id" class="border border-slate-100 rounded-lg p-3 text-sm">
              <p>{{ shipmentTitle(ship) }}</p>
              <p class="text-xs text-slate-400 mt-1">{{ ship.company }} · {{ ship.tracking_no }}</p>
              <p class="text-xs text-slate-400 mt-1">{{ $t('afterSale.detail.status') }}{{ shipmentStatusText(ship.logistics_status, ship.logistics_status_label) }}</p>
              <p v-if="shipmentPrimaryTime(ship)" class="text-xs text-slate-400 mt-1">{{ shipmentTimeLabel(ship) }}：{{ formatDate(shipmentPrimaryTime(ship)) }}</p>
              <p v-if="ship.after_sale_case_id" class="text-xs text-slate-400 mt-1">{{ $t('order.detail.relatedAfterSale') }}#{{ ship.after_sale_case_id }}</p>
              <p v-if="ship.remark" class="text-xs text-slate-400 mt-1">{{ $t('order.detail.remark') }}{{ ship.remark }}</p>
            </div>
          </div>
          <p v-else class="text-slate-400 text-sm">{{ $t('afterSale.detail.noLogistics') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { auditAfterSale, closeAfterSale, completeAfterSale, getAfterSaleDetail, receiveAfterSale, refundAfterSale } from '@/api/plugins'
import { afterSaleStatusLabel, shipmentPrimaryTime, shipmentStatusLabel, shipmentTimeLabel, shipmentTitle } from '@/utils/order-status'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const detail = ref<any>(null)

const typeLabel = (v: string) => v === 'exchange' ? t('afterSale.list.exchange') : t('afterSale.list.return')
const statusLabel = (v: string) => afterSaleStatusLabel(v)
const formatDate = (v?: string) => v ? String(v).slice(0, 19).replace('T', ' ') : '-'
const actionLabelsMap: Record<string, string> = {
  apply: 'afterSale.action.apply',
  audit: 'afterSale.action.review',
  return_ship: 'afterSale.action.returnShip',
  receive: 'afterSale.action.confirmReceive',
  refund: 'afterSale.action.refund',
  reship: 'afterSale.action.reship',
  complete: 'afterSale.action.complete',
  close: 'afterSale.action.close',
}

function actionLabel(action: string) {
  const key = actionLabelsMap[String(action || '')]
  return key ? t(key) : action || '-'
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

function statusLabelOrDash(status: string, label?: string) {
  const value = String(label || status || '')
  if (!value) return '-'
  return statusText(status, label)
}

async function load() {
  detail.value = await getAfterSaleDetail(Number(route.params.id))
}

async function audit(approve: boolean) {
  await auditAfterSale(Number(route.params.id), { approve, audit_remark: '' })
  await load()
}
async function receive() { await receiveAfterSale(Number(route.params.id)); await load() }
async function refund() { await refundAfterSale(Number(route.params.id), { amount: detail.value?.refund_amount || 0, reason: detail.value?.reason || '' }); await load() }
async function complete() { await completeAfterSale(Number(route.params.id)); await load() }
async function close() { await closeAfterSale(Number(route.params.id), { reason: t('afterSale.detail.closeCase') }); await load() }

onMounted(load)
</script>
