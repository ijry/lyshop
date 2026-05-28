<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { computed, reactive, ref } from 'vue'
import StatusTag from '@/components/common/StatusTag.vue'
import Timeline from '@/components/biz/Timeline.vue'
import { addOrderNote, getDeliveryMode, getOrderDetail, getOrderTimeline, getShipmentTracks, remindPay, repriceOrder, shipOrder, syncShipment } from '@/api/order'

const loading = ref(false)
const actionLoading = ref(false)
const trackLoadingID = ref(0)
const orderID = ref(0)
const deliveryMode = ref<'express' | 'local' | 'both'>('express')
const showShipPopup = ref(false)

const timelineItems = ref<Array<{ key: string; title: string; time?: string; tone?: 'primary'|'success'|'warn'|'muted' }>>([])

async function loadTimeline() {
  if (!orderID.value) return
  try {
    const rows: any = await getOrderTimeline(orderID.value)
    timelineItems.value = (Array.isArray(rows) ? rows : []).map((r: any) => ({
      key: r.stage, title: r.status, time: formatDate(r.time),
      tone: r.stage === 'completed' ? 'success' : 'primary',
    }))
  } catch {}
}

function openActionSheet() {
  uni.showActionSheet({
    itemList: ['改价', '添加备注', '打印面单', '催付款'],
    success: async (res) => {
      if (res.tapIndex === 0) {
        if (String(detail.status) !== '1') { uni.showToast({ title: '当前状态不可改价', icon: 'none' }); return }
        uni.showModal({ title: '改价', editable: true, placeholderText: '输入新的支付金额', success: async (m) => {
          if (!m.confirm || !m.content) return
          await repriceOrder(detail.id, { items: detail.items.map((it: any) => ({ item_id: it.id, price: Number(m.content) / (detail.items.length || 1) })), remark: '操作菜单改价' })
          await loadData()
        } })
      } else if (res.tapIndex === 1) {
        uni.showModal({ title: '添加备注', editable: true, success: async (m) => {
          if (!m.confirm || !m.content) return
          await addOrderNote(detail.id, { content: m.content })
          uni.showToast({ title: '已添加', icon: 'success' })
        } })
      } else if (res.tapIndex === 2) {
        uni.navigateTo({ url: `/pages/order/print-preview?id=${detail.id}` })
      } else if (res.tapIndex === 3) {
        await remindPay(detail.id, { channel: 'sms' })
        uni.showToast({ title: '已催付', icon: 'success' })
      }
    },
  })
}

const detail = reactive<any>({
  id: 0,
  status: '',
  status_label: '',
  items: [],
  shipments: [],
  amount_breakdown: {},
})

const tracksMap = reactive<Record<number, any[]>>({})

const shipTypeOptions = [
  { label: '首次发货', value: 'initial' },
  { label: '补发', value: 'reship' },
]

const deliveryTypeOptions = [
  { label: '快递发货', value: 'express' },
  { label: '同城配送', value: 'local' },
]

const companyOptions = [
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

const shipForm = reactive({
  ship_type: 'initial',
  delivery_type: 'express',
  company: 'SF',
  tracking_no: '',
  rider_name: '',
  rider_phone: '',
  after_sale_case_id: '',
  remark: '',
})

const itemList = computed(() => (Array.isArray(detail.items) ? detail.items : []))
const shipmentList = computed(() => (Array.isArray(detail.shipments) ? detail.shipments : []))
const canChooseDeliveryType = computed(() => deliveryMode.value === 'both')
const selectedCompanyName = computed(() => companyOptions.find((x) => x.code === shipForm.company)?.name || '请选择快递公司')

function money(v: any) {
  return Number(v || 0).toFixed(2)
}

function formatDate(v?: string) {
  return v ? String(v).slice(0, 19).replace('T', ' ') : '-'
}

function shipmentTitle(item: any) {
  if (item?.delivery_type === 'local') {
    return `同城配送 · ${item?.rider_name || '-'}`
  }
  return `${item?.company || '快递'} · ${item?.tracking_no || '-'}`
}

function shipmentStatusText(item: any) {
  return String(item?.logistics_status_label || item?.logistics_status || '-')
}

function resetTrackMap() {
  for (const k of Object.keys(tracksMap)) {
    delete tracksMap[Number(k)]
  }
}

function resetShipForm() {
  shipForm.ship_type = 'initial'
  shipForm.delivery_type = deliveryMode.value === 'both' ? 'express' : deliveryMode.value
  shipForm.company = 'SF'
  shipForm.tracking_no = ''
  shipForm.rider_name = ''
  shipForm.rider_phone = ''
  shipForm.after_sale_case_id = ''
  shipForm.remark = ''
}

async function loadDeliveryModeData() {
  try {
    const data: any = await getDeliveryMode()
    const mode = String(data?.mode || 'express')
    if (mode === 'local' || mode === 'both') {
      deliveryMode.value = mode
    } else {
      deliveryMode.value = 'express'
    }
  } catch {
    deliveryMode.value = 'express'
  }
}

async function loadTracks() {
  resetTrackMap()
  const tasks = shipmentList.value
    .filter((item: any) => Number(item?.id || 0) > 0 && item?.delivery_type !== 'local')
    .map(async (item: any) => {
      const rows: any = await getShipmentTracks(orderID.value, item.id)
      tracksMap[Number(item.id)] = Array.isArray(rows) ? rows : []
    })
  await Promise.all(tasks)
}

async function loadData() {
  if (!orderID.value) return
  loading.value = true
  try {
    const data: any = await getOrderDetail(orderID.value)
    Object.assign(detail, {
      id: 0,
      status: '',
      status_label: '',
      items: [],
      shipments: [],
      amount_breakdown: {},
      ...(data || {}),
    })
    await loadTracks()
    await loadTimeline()
  } finally {
    loading.value = false
  }
}

async function syncTrack(shipmentID: number) {
  if (!shipmentID || !orderID.value) return
  trackLoadingID.value = shipmentID
  try {
    await syncShipment(orderID.value, shipmentID)
    await loadData()
  } finally {
    trackLoadingID.value = 0
  }
}

function openShipPopup() {
  resetShipForm()
  showShipPopup.value = true
}

function onPickShipType(e: any) {
  const item = shipTypeOptions[Number(e?.detail?.value || 0)]
  shipForm.ship_type = String(item?.value || 'initial')
}

function onPickDeliveryType(e: any) {
  const item = deliveryTypeOptions[Number(e?.detail?.value || 0)]
  shipForm.delivery_type = String(item?.value || 'express')
}

function onPickCompany(e: any) {
  const item = companyOptions[Number(e?.detail?.value || 0)]
  shipForm.company = String(item?.code || 'SF')
}

async function submitShip() {
  if (shipForm.ship_type === 'reship' && Number(shipForm.after_sale_case_id || 0) <= 0) {
    uni.showToast({ title: '补发需填写售后单 ID', icon: 'none' })
    return
  }
  if (shipForm.delivery_type === 'local') {
    if (!shipForm.rider_name.trim() || !shipForm.rider_phone.trim()) {
      uni.showToast({ title: '请填写骑手姓名和手机号', icon: 'none' })
      return
    }
  } else {
    if (!shipForm.company.trim() || !shipForm.tracking_no.trim()) {
      uni.showToast({ title: '请填写快递公司和物流单号', icon: 'none' })
      return
    }
  }

  actionLoading.value = true
  try {
    await shipOrder(orderID.value, {
      ship_type: shipForm.ship_type,
      delivery_type: shipForm.delivery_type,
      company: shipForm.delivery_type === 'express' ? shipForm.company : undefined,
      tracking_no: shipForm.delivery_type === 'express' ? shipForm.tracking_no.trim() : undefined,
      rider_name: shipForm.delivery_type === 'local' ? shipForm.rider_name.trim() : undefined,
      rider_phone: shipForm.delivery_type === 'local' ? shipForm.rider_phone.trim() : undefined,
      after_sale_case_id: shipForm.ship_type === 'reship' ? Number(shipForm.after_sale_case_id || 0) : undefined,
      remark: shipForm.remark.trim() || undefined,
    })
    uni.showToast({ title: '发货成功', icon: 'success' })
    showShipPopup.value = false
    await loadData()
  } finally {
    actionLoading.value = false
  }
}

onLoad(async (opts) => {
  orderID.value = Number(opts?.id || 0)
  await loadDeliveryModeData()
  await loadData()
})
</script>

<template>
  <view class="page">
    <view v-if="loading" class="empty">加载中...</view>
    <template v-else>
      <view class="card">
        <view class="section-title">订单进度</view>
        <Timeline :items="timelineItems" />
      </view>
      <view class="card">
        <view class="head">
          <text>订单 #{{ detail.id }}</text>
          <StatusTag :text="detail.status_label || detail.status || '-'" :type="detail.status" />
          <text class="op" @click="openActionSheet">操作</text>
        </view>
        <view class="line">下单时间：{{ formatDate(detail.created_at) }}</view>
        <view class="line">收货人：{{ detail.receiver_name || '-' }}</view>
        <view class="line">手机号：{{ detail.receiver_phone || '-' }}</view>
        <view class="line">地址：{{ detail.receiver_address || '-' }}</view>
        <view class="line">商品总额：¥{{ money(detail.amount_breakdown?.goods_amount ?? detail.goods_amount) }}</view>
        <view class="line">优惠金额：¥{{ money(detail.amount_breakdown?.discount_amount ?? detail.discount_amount) }}</view>
        <view class="line">支付金额：¥{{ money(detail.amount_breakdown?.payable_amount ?? detail.pay_amount ?? detail.total_amount) }}</view>
      </view>

      <view class="card">
        <view class="section-title">商品明细</view>
        <view v-if="!itemList.length" class="empty-row">暂无商品</view>
        <view v-for="item in itemList" :key="item.id" class="item-row">
          <view class="item-title">{{ item.title || '-' }}</view>
          <view class="item-sub">数量：{{ item.qty || 0 }} · 单价：¥{{ money(item.price) }}</view>
        </view>
      </view>

      <view class="card">
        <view class="head">
          <view class="section-title">物流信息</view>
          <up-button size="mini" type="primary" plain :loading="actionLoading" @click="openShipPopup">发货</up-button>
        </view>
        <view v-if="!shipmentList.length" class="empty-row">暂无物流记录</view>
        <view v-for="ship in shipmentList" :key="ship.id" class="ship-card">
          <view class="ship-head">
            <text>{{ shipmentTitle(ship) }}</text>
            <up-button
              v-if="ship.delivery_type !== 'local'"
              size="mini"
              type="primary"
              plain
              :loading="trackLoadingID === Number(ship.id)"
              @click="syncTrack(Number(ship.id))"
            >
              同步轨迹
            </up-button>
          </view>
          <view class="item-sub">状态：{{ shipmentStatusText(ship) }}</view>
          <view class="item-sub" v-if="ship.after_sale_case_id">关联售后：#{{ ship.after_sale_case_id }}</view>
          <view class="item-sub" v-if="ship.rider_phone">骑手电话：{{ ship.rider_phone }}</view>
          <view class="item-sub" v-if="ship.remark">备注：{{ ship.remark }}</view>
          <view
            v-for="track in tracksMap[Number(ship.id)] || []"
            :key="track.id"
            class="track-row"
          >
            {{ formatDate(track.event_time) }} · {{ track.status_text }}{{ track.location ? `（${track.location}）` : '' }}
          </view>
        </view>
      </view>
    </template>

    <up-popup :show="showShipPopup" mode="bottom" round="16" @close="showShipPopup = false">
      <view class="popup-body">
        <view class="popup-title">订单发货</view>
        <picker mode="selector" :range="shipTypeOptions" range-key="label" @change="onPickShipType">
          <view class="picker">{{ shipTypeOptions.find((x) => x.value === shipForm.ship_type)?.label || '首次发货' }}</view>
        </picker>
        <view class="mt-12rpx" />

        <picker v-if="canChooseDeliveryType" mode="selector" :range="deliveryTypeOptions" range-key="label" @change="onPickDeliveryType">
          <view class="picker">{{ deliveryTypeOptions.find((x) => x.value === shipForm.delivery_type)?.label || '快递发货' }}</view>
        </picker>
        <view v-else class="picker">{{ shipForm.delivery_type === 'local' ? '同城配送' : '快递发货' }}</view>
        <view class="mt-12rpx" />

        <template v-if="shipForm.delivery_type === 'express'">
          <picker mode="selector" :range="companyOptions" range-key="name" @change="onPickCompany">
            <view class="picker">{{ selectedCompanyName }}</view>
          </picker>
          <view class="mt-12rpx" />
          <up-input v-model="shipForm.tracking_no" placeholder="请输入物流单号" clearable />
        </template>
        <template v-else>
          <up-input v-model="shipForm.rider_name" placeholder="骑手姓名" clearable />
          <view class="mt-12rpx" />
          <up-input v-model="shipForm.rider_phone" placeholder="骑手手机号" clearable />
        </template>

        <view class="mt-12rpx" />
        <up-input v-if="shipForm.ship_type === 'reship'" v-model="shipForm.after_sale_case_id" type="number" inputmode="numeric" placeholder="售后单 ID（补发必填）" />
        <view v-if="shipForm.ship_type === 'reship'" class="mt-12rpx" />
        <up-input v-model="shipForm.remark" placeholder="备注（可选）" clearable />
        <view class="mt-16rpx" />
        <up-button type="primary" :loading="actionLoading" @click="submitShip">确认发货</up-button>
      </view>
    </up-popup>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 24rpx; padding: 22rpx; }
.head { display: flex; justify-content: space-between; align-items: center; gap: 12rpx; }
.section-title { font-size: 30rpx; font-weight: 700; }
.line { margin-top: 10rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.item-row { margin-top: 12rpx; padding: 14rpx; border: 1px solid var(--eapp-border); border-radius: 14rpx; }
.item-title { font-size: 26rpx; font-weight: 600; color: var(--eapp-text); }
.item-sub { margin-top: 6rpx; color: var(--eapp-text-muted); font-size: 23rpx; }
.ship-card { margin-top: 12rpx; padding: 14rpx; border: 1px solid var(--eapp-border); border-radius: 14rpx; }
.ship-head { display: flex; justify-content: space-between; align-items: center; gap: 10rpx; }
.track-row { margin-top: 8rpx; padding-top: 8rpx; border-top: 1px dashed var(--eapp-border); color: var(--eapp-text-muted); font-size: 22rpx; }
.empty { padding: 100rpx 0; text-align: center; color: var(--eapp-text-muted); }
.empty-row { color: var(--eapp-text-muted); font-size: 24rpx; text-align: center; padding: 30rpx 0; }
.popup-body { padding: 24rpx; box-sizing: border-box; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.op { color: var(--eapp-primary); font-size: 24rpx; padding-left: 16rpx; }
.picker { min-height: 76rpx; border: 1px solid var(--eapp-border); border-radius: 12rpx; padding: 0 20rpx; display: flex; align-items: center; color: var(--eapp-text); }
</style>
