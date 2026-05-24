<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">← 返回</button>
      <h2 class="text-xl font-semibold text-slate-800">订单详情</h2>
    </div>

    <div v-if="detail" class="grid grid-cols-1 xl:grid-cols-[2fr_1fr] gap-6">
      <div class="space-y-4">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-4">商品明细</h3>
          <div v-if="detail.items?.length" class="space-y-4">
            <div v-for="it in detail.items" :key="it.id" class="flex items-center gap-3 p-3 rounded-lg border border-slate-100">
              <img :src="it.cover" class="w-16 h-16 rounded object-cover border border-slate-200" />
              <div class="flex-1 min-w-0">
                <p class="text-sm text-slate-700 truncate">{{ it.title }}</p>
                <p class="text-xs text-slate-400 mt-1">数量 x{{ it.qty }}</p>
              </div>
              <p class="text-sm font-medium text-slate-800">¥{{ money(it.price) }}</p>
            </div>
          </div>
          <p v-else class="text-slate-400 text-sm">暂无商品数据</p>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-slate-700">物流轨迹</h3>
            <button class="text-xs text-blue-600 hover:underline" @click="showShipDialog = true">补发/发货</button>
          </div>
          <div v-if="detail.shipments?.length" class="space-y-3">
            <div v-for="ship in detail.shipments" :key="ship.id" class="border border-slate-100 rounded-lg p-3">
              <p class="text-sm text-slate-700">{{ shipmentTitle(ship) }}</p>
              <p class="text-xs text-slate-400 mt-1">{{ ship.company || '未填公司' }} · {{ ship.tracking_no }}</p>
              <p class="text-xs text-slate-400 mt-1">状态：{{ shipmentStatusLabel(ship.logistics_status) }}</p>
              <p v-if="shipmentPrimaryTime(ship)" class="text-xs text-slate-400 mt-1">{{ shipmentTimeLabel(ship) }}：{{ formatDate(shipmentPrimaryTime(ship)) }}</p>
              <p v-if="ship.after_sale_case_id" class="text-xs text-slate-400 mt-1">关联售后单：#{{ ship.after_sale_case_id }}</p>
              <p v-if="ship.remark" class="text-xs text-slate-400 mt-1">备注：{{ ship.remark }}</p>
            </div>
          </div>
          <p v-else class="text-slate-400 text-sm">暂无物流轨迹</p>
        </div>
      </div>

      <div class="space-y-4">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-3">订单信息</h3>
          <div class="space-y-2 text-sm text-slate-600">
            <p>订单号：<span class="font-mono">{{ detail.order_no }}</span></p>
            <p>用户ID：{{ detail.user_id }}</p>
            <p>支付方式：{{ payLabel(detail.payment_method) }}</p>
            <p>状态：{{ statusLabel(detail.status) }}</p>
            <p>下单时间：{{ formatDate(detail.created_at) }}</p>
            <p v-if="detail.paid_at">支付时间：{{ formatDate(detail.paid_at) }}</p>
            <p v-if="detail.tracking_no">物流单号：{{ detail.tracking_no }}</p>
            <p v-if="detail.after_sale_summary?.has_open_case" class="text-red-500">售后中</p>
            <p v-if="detail.after_sale_summary?.latest_case_id">
              最近售后单：#{{ detail.after_sale_summary.latest_case_id }}（{{ afterSaleSummaryStatusText(detail.after_sale_summary) || '-' }}）
            </p>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-3">价格体系</h3>
          <div class="space-y-2 text-sm text-slate-600">
            <div class="flex items-center justify-between"><span>商品金额</span><span>¥{{ money(detail.amount_breakdown?.goods_amount ?? detail.goods_amount) }}</span></div>
            <div class="flex items-center justify-between"><span>优惠金额</span><span>-¥{{ money(detail.amount_breakdown?.discount_amount ?? detail.discount_amount) }}</span></div>
            <div class="flex items-center justify-between"><span>运费</span><span>¥{{ money(detail.amount_breakdown?.freight_amount ?? detail.freight_amount) }}</span></div>
            <div class="flex items-center justify-between text-slate-800 font-semibold pt-2 border-t border-slate-100"><span>实付金额</span><span>¥{{ money(detail.amount_breakdown?.payable_amount ?? detail.total_amount) }}</span></div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="bg-white rounded-xl shadow-sm border border-slate-100 p-10 text-center text-slate-400">
      加载中...
    </div>

    <div v-if="showShipDialog" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <div class="bg-white rounded-2xl p-6 w-[420px] shadow-xl">
        <h3 class="font-semibold text-slate-800 mb-4">补发/发货</h3>
        <div class="space-y-3">
          <select v-model="shipForm.ship_type" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm">
            <option value="initial">首发</option>
            <option value="reship">补发</option>
          </select>
          <input v-model="shipForm.company" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" placeholder="快递公司" />
          <input v-model="shipForm.tracking_no" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" placeholder="快递单号" />
          <input v-if="shipForm.ship_type === 'reship'" v-model="shipForm.after_sale_case_id" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" placeholder="售后单ID" />
        </div>
        <div class="flex gap-3 mt-5">
          <button class="flex-1 px-4 py-2 bg-slate-100 rounded-lg text-sm" @click="showShipDialog = false">取消</button>
          <button class="flex-1 px-4 py-2 bg-blue-700 text-white rounded-lg text-sm" @click="submitShip">确认</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getOrderDetail, shipOrder } from '@/api/plugins'
import { afterSaleStatusLabel, orderStatusLabel, shipmentPrimaryTime, shipmentStatusLabel, shipmentTimeLabel, shipmentTitle } from '@/utils/order-status'

const route = useRoute()
const router = useRouter()
const detail = ref<any>(null)
const showShipDialog = ref(false)
const shipForm = ref<any>({
  ship_type: 'initial',
  company: '',
  tracking_no: '',
  after_sale_case_id: '',
})

const statusLabel = (s: number) => orderStatusLabel(s) || '未知'
const payLabel = (m: string) => m === 'wechat' ? '微信支付' : m === 'alipay' ? '支付宝支付' : '未支付'
const money = (v: number) => Number(v || 0).toFixed(2)
const formatDate = (v?: string) => v ? String(v).slice(0, 19).replace('T', ' ') : '-'

function afterSaleSummaryStatusText(summary: any) {
  const label = String(summary?.latest_status_label || '').trim()
  if (label) return label
  return afterSaleStatusLabel(summary?.latest_status)
}

async function loadDetail() {
  detail.value = await getOrderDetail(Number(route.params.id))
}

async function submitShip() {
  await shipOrder(Number(route.params.id), {
    ship_type: shipForm.value.ship_type,
    company: shipForm.value.company,
    tracking_no: shipForm.value.tracking_no,
    after_sale_case_id: shipForm.value.after_sale_case_id ? Number(shipForm.value.after_sale_case_id) : undefined,
  })
  showShipDialog.value = false
  await loadDetail()
}

onMounted(loadDetail)
</script>
