<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">{{ $t('common.back') }}</button>
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('order.detail.title') }}</h2>
    </div>

    <div v-if="detail" class="grid grid-cols-1 xl:grid-cols-[2fr_1fr] gap-6">
      <div class="space-y-4">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-4">{{ $t('order.detail.items') }}</h3>
          <div v-if="detail.items?.length" class="space-y-4">
            <div v-for="it in detail.items" :key="it.id" class="flex items-center gap-3 p-3 rounded-lg border border-slate-100">
              <img :src="it.cover" class="w-16 h-16 rounded object-cover border border-slate-200" />
              <div class="flex-1 min-w-0">
                <p class="text-sm text-slate-700 truncate">{{ it.title }}</p>
                <p class="text-xs text-slate-400 mt-1">x{{ it.qty }}</p>
              </div>
              <p class="text-sm font-medium text-slate-800">¥{{ money(it.price) }}</p>
            </div>
          </div>
          <p v-else class="text-slate-400 text-sm">{{ $t('order.detail.noItems') }}</p>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-slate-700">{{ $t('order.detail.logistics') }}</h3>
            <button class="text-xs text-blue-600 hover:underline" @click="showShipDialog = true">{{ $t('order.detail.ship') }}</button>
          </div>
          <div v-if="detail.shipments?.length" class="space-y-3">
            <div v-for="ship in detail.shipments" :key="ship.id" class="border border-slate-100 rounded-lg p-3">
              <div class="flex items-center justify-between gap-3">
                <p class="text-sm text-slate-700">{{ shipmentTitle(ship) }}</p>
                <button v-if="ship.delivery_type !== 'local'" data-test="sync-shipment" class="text-xs text-blue-600 hover:underline" @click="syncShip(ship.id)">{{ $t('common.refresh') }}</button>
              </div>
              <template v-if="ship.delivery_type === 'local'">
                <p class="text-xs text-slate-400 mt-1">{{ $t('order.list.rider') }}{{ ship.rider_name || '-' }}</p>
                <p class="text-xs text-slate-400 mt-1">{{ $t('order.detail.riderPhone') }}：{{ ship.rider_phone || '-' }}</p>
              </template>
              <template v-else>
                <p class="text-xs text-slate-400 mt-1">{{ ship.company || '-' }} · {{ ship.tracking_no }}</p>
                <p class="text-xs text-slate-400 mt-1">{{ $t('order.detail.channelProvider') }}{{ ship.channel_provider || '-' }}</p>
              </template>
              <p class="text-xs text-slate-400 mt-1">{{ $t('order.detail.orderStatus') }}{{ shipmentStatusText(ship.logistics_status, ship.logistics_status_label) }}</p>
              <p v-if="shipmentPrimaryTime(ship)" class="text-xs text-slate-400 mt-1">{{ shipmentTimeLabel(ship) }}：{{ formatDate(shipmentPrimaryTime(ship)) }}</p>
              <p v-if="ship.after_sale_case_id" class="text-xs text-slate-400 mt-1">{{ $t('order.detail.relatedAfterSale') }}#{{ ship.after_sale_case_id }}</p>
              <p v-if="ship.remark" class="text-xs text-slate-400 mt-1">{{ $t('order.detail.remark') }}{{ ship.remark }}</p>
              <div v-if="ship.delivery_type !== 'local' && tracksMap[ship.id]?.length" class="mt-3 border-t border-slate-100 pt-3 space-y-2">
                <p class="text-xs text-slate-500">{{ $t('order.detail.trackingNodes') }}</p>
                <div v-for="track in tracksMap[ship.id]" :key="track.id" class="text-xs text-slate-500">
                  <span>{{ formatDate(track.event_time) }}</span>
                  <span class="mx-1">·</span>
                  <span>{{ track.status_text }}</span>
                  <span v-if="track.location">（{{ track.location }}）</span>
                </div>
              </div>
            </div>
          </div>
          <p v-else class="text-slate-400 text-sm">{{ $t('order.detail.noLogistics') }}</p>
        </div>
      </div>

      <div class="space-y-4">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-3">{{ $t('order.detail.orderInfo') }}</h3>
          <div class="space-y-2 text-sm text-slate-600">
            <p>{{ $t('order.detail.orderNo') }}<span class="font-mono">{{ detail.order_no }}</span></p>
            <p>{{ $t('order.detail.userId') }}{{ detail.user_id }}</p>
            <p>{{ $t('order.detail.paymentMethod') }}{{ payLabel(detail.payment_method) }}</p>
            <p>{{ $t('order.detail.orderStatus') }}{{ statusLabel(detail.status) }}</p>
            <p>{{ $t('order.detail.orderTime') }}{{ formatDate(detail.created_at) }}</p>
            <p v-if="detail.paid_at">{{ $t('order.detail.payTime') }}{{ formatDate(detail.paid_at) }}</p>
            <p v-if="detail.tracking_no">{{ $t('order.detail.trackingNoLabel') }}{{ detail.tracking_no }}</p>
            <p v-if="detail.after_sale_summary?.has_open_case" class="text-red-500">{{ $t('order.list.afterSaleActive') }}</p>
            <p v-if="detail.after_sale_summary?.latest_case_id">
              {{ $t('order.detail.latestAfterSale') }}#{{ detail.after_sale_summary.latest_case_id }}（{{ afterSaleSummaryStatusText(detail.after_sale_summary) || '-' }}）
              <button class="text-blue-600 hover:underline text-xs ml-2" @click="goAfterSaleDetail(detail.after_sale_summary.latest_case_id)">{{ $t('common.view') }}</button>
            </p>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-3">{{ $t('order.detail.priceSystem') }}</h3>
          <div class="space-y-2 text-sm text-slate-600">
            <div class="flex items-center justify-between"><span>{{ $t('order.detail.productAmount') }}</span><span>¥{{ money(detail.amount_breakdown?.goods_amount ?? detail.goods_amount) }}</span></div>
            <div class="flex items-center justify-between"><span>{{ $t('order.detail.discountAmount') }}</span><span>-¥{{ money(detail.amount_breakdown?.discount_amount ?? detail.discount_amount) }}</span></div>
            <div class="flex items-center justify-between"><span>{{ $t('order.detail.freightAmount') }}</span><span>¥{{ money(detail.amount_breakdown?.freight_amount ?? detail.freight_amount) }}</span></div>
            <div class="flex items-center justify-between text-slate-800 font-semibold pt-2 border-t border-slate-100"><span>{{ $t('order.detail.paidAmount') }}</span><span>¥{{ money(detail.amount_breakdown?.payable_amount ?? detail.total_amount) }}</span></div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="bg-white rounded-xl shadow-sm border border-slate-100 p-10 text-center text-slate-400">
      {{ $t('order.detail.loading') }}
    </div>

    <div v-if="showShipDialog" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <div class="bg-white rounded-2xl p-6 w-[420px] shadow-xl">
        <h3 class="font-semibold text-slate-800 mb-4">{{ $t('order.detail.ship') }}</h3>
        <div class="space-y-3">
          <select v-model="shipForm.ship_type" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm">
            <option value="initial">{{ $t('order.detail.initialShip') }}</option>
            <option value="reship">{{ $t('order.detail.reShip') }}</option>
          </select>
          <div v-if="storeDeliveryMode === 'both'" class="flex gap-3">
            <label class="flex items-center gap-1.5 text-sm text-slate-600 cursor-pointer">
              <input type="radio" v-model="shipForm.delivery_type" value="express" class="accent-blue-600" />
              {{ $t('order.detail.expressDelivery') }}
            </label>
            <label class="flex items-center gap-1.5 text-sm text-slate-600 cursor-pointer">
              <input type="radio" v-model="shipForm.delivery_type" value="local" class="accent-blue-600" />
              {{ $t('order.detail.localDelivery') }}
            </label>
          </div>
          <template v-if="shipForm.delivery_type === 'express'">
            <select v-model="shipForm.company" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm">
              <option value="">{{ $t('order.detail.selectExpress') }}</option>
              <option v-for="item in courierOptions" :key="item.code" :value="item.code">{{ item.name }}（{{ item.code }}）</option>
            </select>
            <input v-model="shipForm.tracking_no" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" :placeholder="$t('order.detail.trackingNo')" />
          </template>
          <template v-else>
            <input v-model="shipForm.rider_name" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" :placeholder="$t('order.detail.riderName')" />
            <input v-model="shipForm.rider_phone" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" :placeholder="$t('order.detail.riderPhone')" />
          </template>
          <input v-if="shipForm.ship_type === 'reship'" v-model="shipForm.after_sale_case_id" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" :placeholder="$t('order.detail.afterSaleId')" />
        </div>
        <div class="flex gap-3 mt-5">
          <button class="flex-1 px-4 py-2 bg-slate-100 rounded-lg text-sm" @click="showShipDialog = false">{{ $t('common.cancel') }}</button>
          <button class="flex-1 px-4 py-2 bg-blue-700 text-white rounded-lg text-sm" @click="submitShip">{{ $t('common.confirm') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getDeliveryMode, getOrderDetail, getShipmentTracks, shipOrder, syncShipment } from '@/api/plugins'
import { afterSaleStatusLabel, orderStatusLabel, shipmentPrimaryTime, shipmentStatusLabel, shipmentTimeLabel, shipmentTitle } from '@/utils/order-status'
import { notify } from '@/utils/notify'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const detail = ref<any>(null)
const tracksMap = ref<Record<number, any[]>>({})
const showShipDialog = ref(false)
const storeDeliveryMode = ref('express')
const shipForm = ref<any>({
  ship_type: 'initial',
  delivery_type: 'express',
  company: 'SF',
  tracking_no: '',
  rider_name: '',
  rider_phone: '',
  after_sale_case_id: '',
})
const courierOptions = [
  { code: 'SF', name: '顺丰速运' },
  { code: 'ZTO', name: '中通快递' },
  { code: 'YTO', name: '圆通速递' },
  { code: 'STO', name: '申通快递' },
  { code: 'YD', name: '韵达速递' },
  { code: 'JD', name: '京东物流' },
  { code: 'EMS', name: '中国邮政 EMS' },
  { code: 'DBL', name: '德邦快递' },
  { code: 'JT', name: '极兔速递' },
]

const statusLabel = (s: number) => orderStatusLabel(s) || t('common.noData')
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

function goAfterSaleDetail(id: number) {
  if (!id) return
  router.push(`/order/after-sale/detail/${id}`)
}

async function loadDeliveryMode() {
  try {
    const data: any = await getDeliveryMode()
    const mode = data?.mode || 'express'
    storeDeliveryMode.value = mode
    shipForm.value.delivery_type = mode === 'both' ? 'express' : mode
  } catch {
    storeDeliveryMode.value = 'express'
  }
}

async function loadDetail() {
  detail.value = await getOrderDetail(Number(route.params.id))
  tracksMap.value = {}
  const shipments = Array.isArray(detail.value?.shipments) ? detail.value.shipments : []
  for (const shipment of shipments) {
    if (!shipment?.id || shipment.delivery_type === 'local') continue
    const rows: any = await getShipmentTracks(Number(route.params.id), Number(shipment.id))
    tracksMap.value[Number(shipment.id)] = Array.isArray(rows) ? rows : []
  }
}

async function syncShip(shipmentID: number) {
  await syncShipment(Number(route.params.id), shipmentID)
  await loadDetail()
}

async function submitShip() {
  const form = shipForm.value
  if (form.delivery_type === 'express') {
    if (!form.company) {
      notify(t('order.detail.expressRequired'))
      return
    }
    if (!form.tracking_no?.trim()) {
      notify(t('order.detail.trackingRequired'))
      return
    }
  } else {
    if (!form.rider_name?.trim()) {
      notify(t('order.detail.riderNameRequired'))
      return
    }
    if (!form.rider_phone?.trim()) {
      notify(t('order.detail.riderPhoneRequired'))
      return
    }
  }
  await shipOrder(Number(route.params.id), {
    delivery_type: form.delivery_type,
    ship_type: form.ship_type,
    company: form.delivery_type === 'express' ? form.company : undefined,
    tracking_no: form.delivery_type === 'express' ? form.tracking_no : undefined,
    rider_name: form.delivery_type === 'local' ? form.rider_name : undefined,
    rider_phone: form.delivery_type === 'local' ? form.rider_phone : undefined,
    after_sale_case_id: form.after_sale_case_id ? Number(form.after_sale_case_id) : undefined,
  })
  showShipDialog.value = false
  await loadDetail()
}

onMounted(() => {
  loadDeliveryMode()
  loadDetail()
})
</script>
