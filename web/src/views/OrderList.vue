<template>
  <div class="max-w-5xl mx-auto px-6 py-8">
    <h1 class="text-xl font-bold text-gray-900 mb-6">我的订单</h1>

    <!-- Tabs -->
    <div class="flex gap-1 bg-gray-100 rounded-xl p-1 mb-6 w-fit">
      <button v-for="(tab, i) in tabs" :key="tab"
        @click="activeTab = i; loadOrders()"
        :class="activeTab === i ? 'bg-white shadow-sm text-gray-900' : 'text-gray-500 hover:text-gray-700'"
        class="px-5 py-2 rounded-lg text-sm font-medium transition-all">
        {{ tab }}
      </button>
    </div>

    <!-- Empty -->
    <div v-if="!orders.length" class="card p-16 text-center">
      <div class="i-carbon-document text-6xl text-gray-200 mx-auto mb-4" />
      <p class="text-gray-400">暂无订单</p>
    </div>

    <!-- Order list -->
    <div v-else class="space-y-4">
      <div v-for="o in orders" :key="o.id" class="card p-5 hover:shadow-md transition-shadow">
        <div class="flex-between mb-4">
          <div class="flex items-center gap-3">
            <span class="text-xs font-mono text-gray-400 bg-gray-50 px-2 py-1 rounded">{{ o.order_no }}</span>
            <span class="text-xs text-gray-400">{{ o.created_at?.slice(0, 10) }}</span>
          </div>
          <span :class="statusColor(o.status)"
            class="text-xs font-medium px-2.5 py-1 rounded-full">
            {{ statusLabel(o.status) }}
          </span>
        </div>
        <div class="flex-between">
          <span class="text-sm text-gray-500">
            {{ o.payment_method === 'wechat' ? '微信支付' : o.payment_method === 'alipay' ? '支付宝' : '待支付' }}
          </span>
          <span class="text-lg font-bold text-gray-900">¥{{ o.total_amount }}</span>
        </div>
        <div v-if="o.status === 1" class="flex justify-end mt-3">
          <button class="btn-primary !px-6 text-xs">去付款</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/api/request'

const orders = ref<any[]>([])
const activeTab = ref(0)
const tabs = ['全部', '待付款', '待发货', '待收货', '已完成']
const statusValues = [0, 1, 2, 3, 4]

const statusLabels: Record<number, string> = { 1: '待付款', 2: '待发货', 3: '待收货', 4: '已完成', 5: '售后' }
const statusColors: Record<number, string> = {
  1: 'bg-orange-50 text-orange-600',
  2: 'bg-blue-50 text-blue-600',
  3: 'bg-purple-50 text-purple-600',
  4: 'bg-green-50 text-green-600',
  5: 'bg-red-50 text-red-500',
}
const statusLabel = (s: number) => statusLabels[s] || ''
const statusColor = (s: number) => statusColors[s] || 'bg-gray-50 text-gray-400'

async function loadOrders() {
  const status = statusValues[activeTab.value]
  const data = await get<any>('/api/v1/orders', { status: status || undefined })
  orders.value = data?.list || []
}

onMounted(loadOrders)
</script>
