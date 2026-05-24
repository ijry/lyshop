<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">订单列表</h2>
    </div>

    <div class="flex gap-2 mb-4 flex-wrap">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        @click="onTabChange(tab.value)"
        :class="activeStatus === tab.value ? 'bg-blue-700 text-white' : 'bg-white text-slate-600 border border-slate-200 hover:bg-slate-50'"
        class="px-4 py-2 rounded-xl text-sm transition"
      >
        {{ tab.label }}
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">订单号</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">商品信息</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">金额明细</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">用户/支付</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">状态</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="o in orders" :key="o.id" class="align-top hover:bg-slate-50">
            <td class="px-4 py-3">
              <p class="font-mono text-xs text-slate-600">{{ o.order_no }}</p>
              <p class="text-xs text-slate-400 mt-1">{{ formatDate(o.created_at) }}</p>
              <p v-if="o.tracking_no" class="text-xs text-slate-400 mt-1">物流单号：{{ o.tracking_no }}</p>
              <p v-if="o.after_sale_summary?.has_open_case" class="text-xs text-red-500 mt-1">售后中</p>
              <div v-if="hasShipmentSummary(o)" class="flex items-center flex-wrap gap-1 mt-1">
                <span v-if="hasReship(o)" class="text-[11px] px-2 py-0.5 rounded-full bg-orange-50 text-orange-600">含补发</span>
                <span v-if="o.latest_shipment?.tracking_no" class="text-xs text-slate-500">
                  最新物流：{{ shipmentStatusLabel(o.latest_shipment?.logistics_status) }} · {{ o.latest_shipment?.tracking_no }}
                </span>
                <span v-if="shipmentPrimaryTime(o.latest_shipment)" class="text-xs text-slate-400">
                  {{ formatDate(shipmentPrimaryTime(o.latest_shipment)) }}
                </span>
              </div>
            </td>
            <td class="px-4 py-3">
              <div v-if="o.items?.length" class="space-y-2 min-w-[260px]">
                <div v-for="it in o.items.slice(0, 2)" :key="it.id" class="flex items-center gap-2">
                  <img :src="it.cover" class="w-10 h-10 rounded border border-slate-200 object-cover" />
                  <div class="min-w-0">
                    <p class="text-slate-700 truncate">{{ it.title }}</p>
                    <p class="text-xs text-slate-400">x{{ it.qty }}</p>
                  </div>
                </div>
                <p v-if="o.items.length > 2" class="text-xs text-slate-400">共 {{ o.items.length }} 件商品</p>
              </div>
              <p v-else class="text-slate-400 text-xs">暂无商品明细</p>
            </td>
            <td class="px-4 py-3 text-xs text-slate-600">
              <p>商品：¥{{ money(o.amount_breakdown?.goods_amount ?? o.goods_amount) }}</p>
              <p>优惠：-¥{{ money(o.amount_breakdown?.discount_amount ?? o.discount_amount) }}</p>
              <p>运费：¥{{ money(o.amount_breakdown?.freight_amount ?? o.freight_amount) }}</p>
              <p class="text-slate-800 font-semibold mt-1">实付：¥{{ money(o.amount_breakdown?.payable_amount ?? o.total_amount) }}</p>
            </td>
            <td class="px-4 py-3 text-xs text-slate-600">
              <p>用户ID：{{ o.user_id }}</p>
              <p class="mt-1">{{ payLabel(o.payment_method) }}</p>
            </td>
            <td class="px-4 py-3">
              <span :class="statusClass(o.status)" class="px-2 py-1 rounded-full text-xs">
                {{ statusLabel(o.status) }}
              </span>
            </td>
            <td class="px-4 py-3">
              <div class="flex flex-col gap-2 items-start">
                <button class="text-blue-600 hover:underline text-xs" @click="goDetail(o.id)">查看详情</button>
                <button v-if="o.status === 2" @click="ship(o.id)" class="text-emerald-600 hover:underline text-xs">发货</button>
                <button v-if="o.after_sale_summary?.has_open_case" @click="goAfterSale(o.id)" class="text-red-600 hover:underline text-xs">售后</button>
              </div>
            </td>
          </tr>
          <tr v-if="!orders.length">
            <td colspan="6" class="px-4 py-12 text-center text-slate-400">暂无订单</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getOrders, shipOrder } from '@/api/plugins'

const router = useRouter()
const orders = ref<any[]>([])
const activeStatus = ref(0)

const tabs = [
  { label: '全部', value: 0 },
  { label: '待付款', value: 1 },
  { label: '待发货', value: 2 },
  { label: '待收货', value: 3 },
  { label: '已完成', value: 4 },
  { label: '售后', value: 5 },
]

const statusLabels: Record<number, string> = {
  1: '待付款', 2: '待发货', 3: '待收货', 4: '已完成', 5: '售后'
}
const statusColors: Record<number, string> = {
  1: 'bg-yellow-50 text-yellow-700',
  2: 'bg-blue-50 text-blue-700',
  3: 'bg-purple-50 text-purple-700',
  4: 'bg-green-50 text-green-700',
  5: 'bg-red-50 text-red-600',
}
const shipmentStatusLabels: Record<string, string> = {
  pending: '待揽收',
  shipped: '已发货',
  in_transit: '运输中',
  signed: '已签收',
  exception: '物流异常',
}

const statusLabel = (s: number) => statusLabels[s] || '未知'
const statusClass = (s: number) => statusColors[s] || 'bg-slate-50 text-slate-400'
const payLabel = (m: string) => m === 'wechat' ? '微信支付' : m === 'alipay' ? '支付宝支付' : '未支付'
const money = (v: number) => Number(v || 0).toFixed(2)
const formatDate = (v?: string) => v ? String(v).slice(0, 19).replace('T', ' ') : '-'

function shipmentStatusLabel(status: string) {
  return shipmentStatusLabels[status] || status || '-'
}

function shipmentPrimaryTime(shipment: any) {
  if (!shipment) return ''
  return String(shipment.signed_at || shipment.shipped_at || shipment.created_at || '')
}

function hasReship(order: any) {
  const shipments = Array.isArray(order?.shipments) ? order.shipments : []
  return shipments.some((ship: any) => String(ship?.biz_type || '') === 'reship')
}

function hasShipmentSummary(order: any) {
  return Boolean(hasReship(order) || order?.latest_shipment?.tracking_no)
}

function goDetail(id: number) {
  router.push(`/order/detail/${id}`)
}

function goAfterSale(id: number) {
  router.push(`/order/after-sale/list?order_id=${id}`)
}

function onTabChange(status: number) {
  activeStatus.value = status
  loadOrders()
}

async function loadOrders() {
  const data: any = await getOrders({ status: activeStatus.value || undefined, page: 1, size: 20 })
  orders.value = data?.list || []
}

async function ship(id: number) {
  const no = prompt('请输入快递单号')
  if (!no) return
  await shipOrder(id, { tracking_no: no })
  await loadOrders()
}

onMounted(loadOrders)
</script>
