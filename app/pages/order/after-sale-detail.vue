<template>
  <view class="min-h-screen bg-gray-50 pb-40rpx">
    <u-navbar title="售后详情" :placeholder="true" />

    <view v-if="detail.id" class="p-20rpx">
      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block">售后单号：{{ detail.case_no }}</text>
        <text class="text-24rpx text-gray-600 block mt-8rpx">订单ID：{{ detail.order_id }}</text>
        <text class="text-24rpx text-gray-600 block mt-8rpx">类型：{{ detail.case_type === 'exchange' ? '换货' : '退货' }}</text>
        <text class="text-24rpx text-gray-600 block mt-8rpx">状态：{{ statusLabel(detail.status) }}</text>
        <text class="text-24rpx text-gray-600 block mt-8rpx">原因：{{ detail.reason }}</text>
        <text v-if="detail.apply_content" class="text-24rpx text-gray-600 block mt-8rpx">说明：{{ detail.apply_content }}</text>
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
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">物流轨迹</text>
        <view v-if="detail.shipments?.length" class="space-y-12rpx">
          <view v-for="ship in detail.shipments" :key="ship.id" class="border border-gray-100 rounded-16rpx p-16rpx">
            <text class="text-24rpx text-gray-700 block">{{ shipmentTitle(ship) }}</text>
            <text class="text-22rpx text-gray-500 block mt-6rpx">{{ ship.company || '未填公司' }} · {{ ship.tracking_no || '-' }}</text>
            <text class="text-22rpx text-gray-500 block mt-6rpx">状态：{{ shipmentStatusLabel(ship.logistics_status) }}</text>
            <text v-if="shipmentPrimaryTime(ship)" class="text-22rpx text-gray-500 block mt-6rpx">
              {{ shipmentTimeLabel(ship) }}：{{ formatDate(shipmentPrimaryTime(ship)) }}
            </text>
            <text v-if="ship.after_sale_case_id" class="text-22rpx text-gray-500 block mt-6rpx">关联售后单：#{{ ship.after_sale_case_id }}</text>
            <text v-if="ship.remark" class="text-22rpx text-gray-500 block mt-6rpx">备注：{{ ship.remark }}</text>
          </view>
        </view>
        <text v-else class="text-24rpx text-gray-400">暂无物流轨迹</text>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">进度日志</text>
        <view v-if="detail.logs?.length" class="space-y-12rpx">
          <view v-for="log in detail.logs" :key="log.id" class="border border-gray-100 rounded-16rpx p-16rpx">
            <text class="text-24rpx text-gray-700 block">{{ log.action }}：{{ log.from_status || '-' }} → {{ log.to_status }}</text>
            <text class="text-22rpx text-gray-500 block mt-6rpx">{{ log.content || '-' }}</text>
            <text class="text-22rpx text-gray-400 block mt-6rpx">{{ formatDate(log.created_at) }}</text>
          </view>
        </view>
        <text v-else class="text-24rpx text-gray-400">暂无状态日志</text>
      </view>

      <view v-if="detail.status === 'approved_wait_user_return'" class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">填写回寄物流</text>
        <view class="mb-12rpx">
          <u-input v-model="returnForm.company" border="surround" placeholder="快递公司（选填）" />
        </view>
        <view class="mb-12rpx">
          <u-input v-model="returnForm.tracking_no" border="surround" placeholder="快递单号（必填）" />
        </view>
        <view class="mb-16rpx">
          <u-textarea v-model="returnForm.remark" placeholder="备注（选填）" :auto-height="true" maxlength="200" />
        </view>
        <u-button type="primary" shape="circle" :loading="submittingReturn" :text="submittingReturn ? '提交中...' : '提交回寄物流'" @click="submitReturnShipment" />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { get, post } from '@/utils/request'

const detail = ref<any>({})
const caseID = ref(0)
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
    toast('请填写快递单号')
    return
  }
  submittingReturn.value = true
  try {
    await post(`/api/v1/after-sales/${caseID.value}/return-shipments`, {
      company: String(returnForm.value.company || ''),
      tracking_no: trackingNo,
      remark: String(returnForm.value.remark || ''),
    })
    uni.showToast({ title: '提交成功', icon: 'success' })
    returnForm.value = { company: '', tracking_no: '', remark: '' }
    await load()
  } catch (error: any) {
    toast(error?.message || '提交失败')
  } finally {
    submittingReturn.value = false
  }
}

onMounted(async () => {
  caseID.value = readCaseID()
  await load()
})
</script>
