<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { computed, reactive, ref } from 'vue'
import StatusTag from '@/components/common/StatusTag.vue'
import {
  auditAfterSale,
  closeAfterSale,
  completeAfterSale,
  getAfterSaleDetail,
  receiveAfterSale,
  refundAfterSale,
  addAfterSaleMessage,
  addAfterSaleEvidence,
} from '@/api/after-sale'

const loading = ref(false)
const actionLoading = ref(false)
const showAuditPopup = ref(false)
const showRefundPopup = ref(false)
const showClosePopup = ref(false)

const detail = reactive<any>({
  id: 0,
  status: '',
  status_label: '',
  logs: [],
  shipments: [],
})

const auditForm = reactive({
  approve: true,
  audit_remark: '',
})

const refundForm = reactive({
  amount: 0,
  reason: '',
  refund_no: '',
})

const closeForm = reactive({
  reason: '',
})

const logs = computed(() => (Array.isArray(detail.logs) ? detail.logs : []))
const shipments = computed(() => (Array.isArray(detail.shipments) ? detail.shipments : []))
const currentStatus = computed(() => String(detail.status || ''))
const canAudit = computed(() => currentStatus.value === 'applied')
const canReceive = computed(() => currentStatus.value === 'user_returning')
const canRefund = computed(() => currentStatus.value === 'refund_pending')
const canComplete = computed(() => ['refunded', 'reshipped'].includes(currentStatus.value))
const canClose = computed(() =>
  [
    'applied',
    'approved_wait_user_return',
    'user_returning',
    'warehouse_received',
    'refund_pending',
    'reship_pending',
  ].includes(currentStatus.value),
)
const needReshipTip = computed(() => currentStatus.value === 'reship_pending' && detail.case_type === 'exchange')

const steps = [
  { key: 'applied', label: '申请' },
  { key: 'approved_wait_user_return', label: '审核' },
  { key: 'user_returning', label: '寄回' },
  { key: 'warehouse_received', label: '收货' },
  { key: 'refunded', label: '退款' },
]
function stepActive(idx: number) {
  const order = ['applied', 'approved_wait_user_return', 'user_returning', 'warehouse_received', 'refund_pending', 'refunded', 'reshipped']
  const cur = order.indexOf(String(detail.status))
  return cur >= idx
}
const messages = computed(() => (Array.isArray(detail.messages) ? detail.messages : []))
const evidences = computed(() => (Array.isArray(detail.evidences) ? detail.evidences : []))
const msgInput = ref('')
const sendingMsg = ref(false)

async function onSendMsg() {
  if (!msgInput.value.trim()) return
  sendingMsg.value = true
  try {
    await addAfterSaleMessage(detail.id, { from: 'merchant', content: msgInput.value.trim() })
    msgInput.value = ''; await loadData()
  } finally { sendingMsg.value = false }
}

function onUploadEvidence() {
  uni.chooseImage({ count: 3, success: async (res) => {
    const images = (res.tempFilePaths || []).map((_: any, i: number) => `https://picsum.photos/200/200?random=evi-${Date.now()}-${i}`)
    await addAfterSaleEvidence(detail.id, { images, remark: '商家凭证' })
    await loadData()
  } })
}

function formatDate(v?: string) {
  return v ? String(v).slice(0, 19).replace('T', ' ') : '-'
}

function money(v: any) {
  return Number(v || 0).toFixed(2)
}

function caseTypeText() {
  if (detail.case_type_label) return String(detail.case_type_label)
  return detail.case_type === 'exchange' ? '换货' : '退货'
}

function openAudit(approve: boolean) {
  if (!canAudit.value) {
    uni.showToast({ title: '当前状态不可审核', icon: 'none' })
    return
  }
  auditForm.approve = approve
  auditForm.audit_remark = ''
  showAuditPopup.value = true
}

function openRefund() {
  if (!canRefund.value) {
    uni.showToast({ title: '当前状态不可退款', icon: 'none' })
    return
  }
  refundForm.amount = Number(detail.refund_amount || detail.total_amount || 0)
  refundForm.reason = String(detail.reason || '')
  refundForm.refund_no = ''
  showRefundPopup.value = true
}

function openClose() {
  if (!canClose.value) {
    uni.showToast({ title: '当前状态不可关闭', icon: 'none' })
    return
  }
  closeForm.reason = ''
  showClosePopup.value = true
}

function goOrderDetail() {
  if (!detail.order_id) return
  uni.navigateTo({ url: `/pages/order/detail?id=${detail.order_id}` })
}

async function loadData(id = 0) {
  const targetID = Number(id || detail.id || 0)
  if (!targetID) return
  loading.value = true
  try {
    const data: any = await getAfterSaleDetail(targetID)
    Object.assign(detail, {
      id: targetID,
      status: '',
      status_label: '',
      logs: [],
      shipments: [],
      ...(data || {}),
    })
  } finally {
    loading.value = false
  }
}

async function submitAudit() {
  actionLoading.value = true
  try {
    await auditAfterSale(detail.id, {
      approve: auditForm.approve,
      audit_remark: auditForm.audit_remark.trim(),
    })
    showAuditPopup.value = false
    uni.showToast({ title: '审核已提交', icon: 'success' })
    await loadData()
  } finally {
    actionLoading.value = false
  }
}

async function confirmReceive() {
  if (!canReceive.value) {
    uni.showToast({ title: '当前状态不可确认收货', icon: 'none' })
    return
  }
  const res = await new Promise<UniNamespace.ShowModalResolve>((resolve) => {
    uni.showModal({
      title: '确认收货',
      content: '确认已收到用户寄回商品？',
      success: resolve,
    })
  })
  if (!res.confirm) return

  actionLoading.value = true
  try {
    await receiveAfterSale(detail.id)
    uni.showToast({ title: '已确认收货', icon: 'success' })
    await loadData()
  } finally {
    actionLoading.value = false
  }
}

async function submitRefund() {
  if (Number(refundForm.amount || 0) <= 0) {
    uni.showToast({ title: '请输入有效退款金额', icon: 'none' })
    return
  }

  actionLoading.value = true
  try {
    await refundAfterSale(detail.id, {
      amount: Number(refundForm.amount || 0),
      reason: refundForm.reason.trim(),
      refund_no: refundForm.refund_no.trim(),
    })
    showRefundPopup.value = false
    uni.showToast({ title: '退款登记成功', icon: 'success' })
    await loadData()
  } finally {
    actionLoading.value = false
  }
}

async function confirmComplete() {
  if (!canComplete.value) {
    uni.showToast({ title: '当前状态不可完结', icon: 'none' })
    return
  }
  const res = await new Promise<UniNamespace.ShowModalResolve>((resolve) => {
    uni.showModal({
      title: '完结售后',
      content: '确认将当前售后单标记为完结？',
      success: resolve,
    })
  })
  if (!res.confirm) return

  actionLoading.value = true
  try {
    await completeAfterSale(detail.id)
    uni.showToast({ title: '售后已完结', icon: 'success' })
    await loadData()
  } finally {
    actionLoading.value = false
  }
}

async function submitClose() {
  if (!closeForm.reason.trim()) {
    uni.showToast({ title: '请填写关闭原因', icon: 'none' })
    return
  }

  actionLoading.value = true
  try {
    await closeAfterSale(detail.id, { reason: closeForm.reason.trim() })
    showClosePopup.value = false
    uni.showToast({ title: '售后单已关闭', icon: 'success' })
    await loadData()
  } finally {
    actionLoading.value = false
  }
}

onLoad((opts) => {
  detail.id = Number(opts?.id || 0)
  loadData()
})
</script>

<template>
  <view class="page">
    <view v-if="loading" class="empty">加载中...</view>
    <template v-else>
      <view v-if="!loading" class="progress">
        <view v-for="(step, idx) in steps" :key="step.key" :class="['step', stepActive(idx) ? 'active' : '']">
          <view class="dot">{{ idx + 1 }}</view>
          <text class="label">{{ step.label }}</text>
        </view>
      </view>
      <view class="card">
        <view class="head">
          <text>售后 #{{ detail.id }}</text>
          <StatusTag :text="detail.status_label || detail.status || '-'" :type="detail.status" />
        </view>
        <view class="line">订单号：#{{ detail.order_id || '-' }}</view>
        <view class="line">类型：{{ caseTypeText() }}</view>
        <view class="line">原因：{{ detail.reason || '-' }}</view>
        <view class="line">说明：{{ detail.apply_content || detail.description || '-' }}</view>
        <view class="line">申请时间：{{ formatDate(detail.created_at) }}</view>
        <view class="line">退款金额：¥{{ money(detail.refund_amount || 0) }}</view>
      </view>

      <view class="card">
        <view class="section-title">售后动作</view>
        <view v-if="needReshipTip" class="tips">
          <text>当前售后为“待补发”，请到订单详情执行补发。</text>
          <up-button size="mini" type="primary" plain @click="goOrderDetail">去订单补发</up-button>
        </view>
        <view class="action-grid">
          <up-button size="mini" type="primary" :loading="actionLoading" :disabled="!canAudit || actionLoading" @click="openAudit(true)">审核通过</up-button>
          <up-button size="mini" type="warning" :loading="actionLoading" :disabled="!canAudit || actionLoading" @click="openAudit(false)">审核拒绝</up-button>
          <up-button size="mini" type="success" :loading="actionLoading" :disabled="!canReceive || actionLoading" @click="confirmReceive">确认收货</up-button>
          <up-button size="mini" type="error" :loading="actionLoading" :disabled="!canRefund || actionLoading" @click="openRefund">登记退款</up-button>
          <up-button size="mini" type="primary" plain :loading="actionLoading" :disabled="!canComplete || actionLoading" @click="confirmComplete">完结</up-button>
          <up-button size="mini" type="error" plain :loading="actionLoading" :disabled="!canClose || actionLoading" @click="openClose">关闭</up-button>
        </view>
      </view>

      <view class="card">
        <view class="section-title">状态日志</view>
        <view v-if="!logs.length" class="empty-row">暂无日志</view>
        <view v-for="item in logs" :key="item.id" class="log-row">
          <view class="item-sub">{{ item.action_label || item.action || '-' }}</view>
          <view class="item-sub">{{ item.from_status_label || item.from_status || '-' }} → {{ item.to_status_label || item.to_status || '-' }}</view>
          <view class="item-sub">{{ item.content || '-' }}</view>
          <view class="item-sub">{{ formatDate(item.created_at) }}</view>
        </view>
      </view>

      <view class="card">
        <view class="section-title">物流信息</view>
        <view v-if="!shipments.length" class="empty-row">暂无物流信息</view>
        <view v-for="item in shipments" :key="item.id" class="log-row">
          <view class="item-sub">{{ item.company || '-' }} · {{ item.tracking_no || '-' }}</view>
          <view class="item-sub">方向：{{ item.direction_label || item.direction || '-' }} · 业务：{{ item.biz_type_label || item.biz_type || '-' }}</view>
          <view class="item-sub">状态：{{ item.logistics_status_label || item.logistics_status || '-' }}</view>
          <view class="item-sub">创建时间：{{ formatDate(item.created_at) }}</view>
        </view>
      </view>

      <view class="card">
        <view class="section-title">协商沟通</view>
        <view v-if="!messages.length" class="empty-row">暂无沟通记录</view>
        <view v-for="m in messages" :key="m.id" :class="['msg', m.from === 'merchant' ? 'right' : '']">
          <view class="msg-body">
            <text class="who">{{ m.from === 'merchant' ? '商家' : '买家' }}</text>
            <text class="text">{{ m.content }}</text>
            <view v-if="m.images && m.images.length" class="imgs">
              <image v-for="img in m.images" :key="img" :src="img" mode="aspectFill" class="img" />
            </view>
            <text class="time">{{ formatDate(m.created_at) }}</text>
          </view>
        </view>
        <view class="msg-input">
          <up-input v-model="msgInput" placeholder="回复买家" />
          <up-button type="primary" size="mini" :loading="sendingMsg" @click="onSendMsg">发送</up-button>
        </view>
      </view>

      <view class="card">
        <view class="section-title">
          <text>商家凭证</text>
          <up-button size="mini" plain @click="onUploadEvidence">上传凭证</up-button>
        </view>
        <view v-if="!evidences.length" class="empty-row">暂未上传凭证</view>
        <view v-for="e in evidences" :key="e.id" class="ev">
          <view class="ev-images">
            <image v-for="img in e.images" :key="img" :src="img" mode="aspectFill" class="img" />
          </view>
          <text class="time">{{ formatDate(e.created_at) }} · {{ e.remark || '' }}</text>
        </view>
      </view>
    </template>

    <up-popup :show="showAuditPopup" mode="bottom" round="16" @close="showAuditPopup = false">
      <view class="popup-body">
        <view class="popup-title">{{ auditForm.approve ? '审核通过' : '审核拒绝' }}</view>
        <up-input v-model="auditForm.audit_remark" placeholder="审核备注（可选）" clearable />
        <view class="mt-16rpx" />
        <up-button type="primary" :loading="actionLoading" @click="submitAudit">提交审核</up-button>
      </view>
    </up-popup>

    <up-popup :show="showRefundPopup" mode="bottom" round="16" @close="showRefundPopup = false">
      <view class="popup-body">
        <view class="popup-title">登记退款</view>
        <up-input v-model="refundForm.amount" type="digit" inputmode="decimal" placeholder="退款金额" />
        <view class="mt-12rpx" />
        <up-input v-model="refundForm.reason" placeholder="退款原因（可选）" clearable />
        <view class="mt-12rpx" />
        <up-input v-model="refundForm.refund_no" placeholder="退款单号（可选）" clearable />
        <view class="mt-16rpx" />
        <up-button type="primary" :loading="actionLoading" @click="submitRefund">提交退款</up-button>
      </view>
    </up-popup>

    <up-popup :show="showClosePopup" mode="bottom" round="16" @close="showClosePopup = false">
      <view class="popup-body">
        <view class="popup-title">关闭售后单</view>
        <up-input v-model="closeForm.reason" placeholder="请输入关闭原因" clearable />
        <view class="mt-16rpx" />
        <up-button type="error" :loading="actionLoading" @click="submitClose">确认关闭</up-button>
      </view>
    </up-popup>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; }
.head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10rpx; }
.section-title { font-size: 30rpx; font-weight: 700; margin-bottom: 10rpx; }
.line { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.item-sub { margin-top: 6rpx; color: var(--eapp-text-muted); font-size: 23rpx; }
.action-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10rpx; }
.log-row { margin-top: 12rpx; padding: 14rpx; border: 1px solid var(--eapp-border); border-radius: 14rpx; }
.tips { margin-bottom: 12rpx; padding: 14rpx; border: 1px solid #bfdbfe; border-radius: 14rpx; background: #eff6ff; display: flex; align-items: center; justify-content: space-between; gap: 12rpx; color: #1e3a8a; font-size: 23rpx; }
.empty { padding: 100rpx 0; text-align: center; color: var(--eapp-text-muted); }
.empty-row { color: var(--eapp-text-muted); font-size: 24rpx; text-align: center; padding: 30rpx 0; }
.popup-body { padding: 24rpx; box-sizing: border-box; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.progress { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; display: flex; justify-content: space-between; align-items: center; }
.step { display: flex; flex-direction: column; align-items: center; gap: 6rpx; flex: 1; }
.step .dot { width: 44rpx; height: 44rpx; border-radius: 50%; background: var(--eapp-border); color: #fff; display: flex; align-items: center; justify-content: center; }
.step.active .dot { background: var(--eapp-primary); }
.step .label { font-size: 22rpx; color: var(--eapp-text-muted); }
.step.active .label { color: var(--eapp-text); }
.msg { display: flex; margin-top: 10rpx; }
.msg.right { justify-content: flex-end; }
.msg-body { max-width: 70%; background: var(--eapp-bg); padding: 12rpx 14rpx; border-radius: 16rpx; }
.msg.right .msg-body { background: var(--eapp-primary-soft); }
.who { display: block; font-size: 20rpx; color: var(--eapp-text-muted); }
.text { font-size: 24rpx; color: var(--eapp-text); }
.imgs { display: flex; gap: 8rpx; margin-top: 6rpx; }
.img { width: 120rpx; height: 120rpx; border-radius: 10rpx; }
.time { display: block; font-size: 20rpx; color: var(--eapp-text-faint); margin-top: 4rpx; }
.msg-input { margin-top: 12rpx; display: flex; gap: 10rpx; align-items: center; }
.ev { margin-top: 10rpx; padding-top: 10rpx; border-top: 1px dashed var(--eapp-border); }
.ev-images { display: flex; gap: 8rpx; flex-wrap: wrap; }
</style>
