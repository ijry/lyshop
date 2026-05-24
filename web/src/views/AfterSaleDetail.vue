<template>
  <div class="max-w-4xl mx-auto px-6 py-8" v-if="detail">
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">← 返回</button>
      <h1 class="text-xl font-bold text-gray-900">售后详情</h1>
    </div>

    <div class="card p-5 mb-4">
      <p class="text-sm text-gray-600">售后单号：<span class="font-mono">{{ detail.case_no }}</span></p>
      <p class="text-sm text-gray-600 mt-1">订单ID：{{ detail.order_id }}</p>
      <p class="text-sm text-gray-600 mt-1">类型：{{ detail.case_type === 'exchange' ? '换货' : '退货' }}</p>
      <p class="text-sm text-gray-600 mt-1">状态：{{ statusLabel(detail.status) }}</p>
      <p class="text-sm text-gray-600 mt-1">原因：{{ detail.reason }}</p>
      <p v-if="detail.apply_content" class="text-sm text-gray-600 mt-1">说明：{{ detail.apply_content }}</p>
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
      <h3 class="text-sm font-semibold text-gray-800 mb-3">物流轨迹</h3>
      <div v-if="detail.shipments?.length" class="space-y-3">
        <div v-for="ship in detail.shipments" :key="ship.id" class="border border-gray-100 rounded-lg p-3">
          <p class="text-sm text-gray-700">{{ shipmentTitle(ship) }}</p>
          <p class="text-xs text-gray-500 mt-1">{{ ship.company || '未填公司' }} · {{ ship.tracking_no || '-' }}</p>
          <p class="text-xs text-gray-500 mt-1">状态：{{ shipmentStatusLabel(ship.logistics_status) }}</p>
          <p v-if="shipmentPrimaryTime(ship)" class="text-xs text-gray-500 mt-1">
            {{ shipmentTimeLabel(ship) }}：{{ formatDate(shipmentPrimaryTime(ship)) }}
          </p>
          <p v-if="ship.after_sale_case_id" class="text-xs text-gray-500 mt-1">关联售后单：#{{ ship.after_sale_case_id }}</p>
          <p v-if="ship.remark" class="text-xs text-gray-500 mt-1">备注：{{ ship.remark }}</p>
        </div>
      </div>
      <p v-else class="text-sm text-gray-400">暂无物流轨迹</p>
    </div>

    <div class="card p-5 mb-4">
      <h3 class="text-sm font-semibold text-gray-800 mb-3">进度日志</h3>
      <div v-if="detail.logs?.length" class="space-y-3">
        <div v-for="log in detail.logs" :key="log.id" class="border border-gray-100 rounded-lg p-3">
          <p class="text-sm text-gray-700">{{ log.action }}：{{ log.from_status || '-' }} → {{ log.to_status }}</p>
          <p class="text-xs text-gray-500 mt-1">{{ log.content || '-' }}</p>
          <p class="text-xs text-gray-400 mt-1">{{ formatDate(log.created_at) }}</p>
        </div>
      </div>
      <p v-else class="text-sm text-gray-400">暂无状态日志</p>
    </div>

    <div v-if="detail.status === 'approved_wait_user_return'" class="card p-5">
      <h3 class="text-sm font-semibold text-gray-800 mb-3">填写回寄物流</h3>
      <div class="space-y-3">
        <input v-model="returnForm.company" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm" placeholder="快递公司（选填）" />
        <input v-model="returnForm.tracking_no" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm" placeholder="快递单号（必填）" />
        <textarea v-model="returnForm.remark" class="w-full min-h-[72px] border border-gray-200 rounded-lg px-3 py-2 text-sm" placeholder="备注（选填）" />
      </div>
      <button class="btn-primary w-full !py-3 mt-4" :disabled="submittingReturn" @click="submitReturnShipment">
        {{ submittingReturn ? '提交中...' : '提交回寄物流' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get, post } from '@/api/request'

const route = useRoute()
const router = useRouter()
const detail = ref<any>(null)
const submittingReturn = ref(false)
const returnForm = ref<any>({
  company: '',
  tracking_no: '',
  remark: '',
})
const shipmentBizLabels: Record<string, string> = { initial: '首发', reship: '补发', return: '回寄' }
const shipmentStatusLabels: Record<string, string> = {
  pending: '待揽收',
  shipped: '已发货',
  in_transit: '运输中',
  signed: '已签收',
  exception: '物流异常',
}

const statusMap: Record<string, string> = {
  applied: '已申请',
  approved_wait_user_return: '待回寄',
  user_returning: '回寄中',
  warehouse_received: '仓库已收货',
  refund_pending: '待退款',
  refunded: '已退款',
  reship_pending: '待补发',
  reshipped: '已补发',
  completed: '已完结',
  rejected: '已拒绝',
  closed: '已关闭',
}

const statusLabel = (status: string) => statusMap[status] || status
const formatDate = (v?: string) => (v ? String(v).slice(0, 19).replace('T', ' ') : '-')

function shipmentDirectionLabel(direction: string) {
  return direction === 'inbound' ? '回寄' : '寄出'
}

function shipmentBizLabel(bizType: string) {
  return shipmentBizLabels[bizType] || bizType || '-'
}

function shipmentStatusLabel(status: string) {
  return shipmentStatusLabels[status] || status || '-'
}

function shipmentTitle(ship: any) {
  return `${shipmentDirectionLabel(String(ship.direction || ''))} · ${shipmentBizLabel(String(ship.biz_type || ''))}`
}

function shipmentPrimaryTime(ship: any) {
  return String(ship.signed_at || ship.shipped_at || ship.created_at || '')
}

function shipmentTimeLabel(ship: any) {
  if (ship.signed_at) return '签收时间'
  if (ship.shipped_at) return '发货时间'
  return '记录时间'
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
    alertMsg('请填写快递单号')
    return
  }
  submittingReturn.value = true
  try {
    await post(`/api/v1/after-sales/${route.params.id}/return-shipments`, {
      company: String(returnForm.value.company || ''),
      tracking_no: trackingNo,
      remark: String(returnForm.value.remark || ''),
    })
    alertMsg('回寄物流提交成功')
    returnForm.value = { company: '', tracking_no: '', remark: '' }
    await load()
  } catch (error: any) {
    alertMsg(error?.message || '提交失败')
  } finally {
    submittingReturn.value = false
  }
}

onMounted(load)
</script>
