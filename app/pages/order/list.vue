<template>
  <view class="bg-gray-50 min-h-screen">
    <u-navbar title="我的订单" :placeholder="true" />

    <!-- Status tabs -->
    <u-tabs :list="tabs" :current="activeTab" @click="onTab" />

    <view class="p-3">
      <view v-if="!orders.length" class="text-center py-12 text-gray-400">
        <text>暂无订单</text>
      </view>
      <view v-for="o in orders" :key="o.id" class="bg-white rounded-xl p-4 mb-3 shadow-sm">
        <view class="flex justify-between mb-2">
          <text class="text-xs text-gray-400 font-mono">{{ o.order_no }}</text>
          <text :class="statusColor(o.status)" class="text-xs font-medium">{{ statusLabel(o.status) }}</text>
        </view>
        <view class="flex justify-between items-center">
          <text class="text-gray-600 text-sm">{{ o.created_at?.slice(0, 10) }}</text>
          <text class="text-slate-800 font-bold">¥{{ o.total_amount }}</text>
        </view>
        <view class="flex justify-end gap-2 mt-3">
          <u-button v-if="o.status === 1" size="mini" type="primary" text="去付款" @click="toPay(o)" />
          <u-button v-if="o.status === 4" size="mini" type="success" text="评价" plain />
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'

const orders = ref<any[]>([])
const activeTab = ref(0)

const tabs = [
  { name: '全部' }, { name: '待付款' }, { name: '待发货' },
  { name: '待收货' }, { name: '已完成' }
]
const statusValues = [0, 1, 2, 3, 4]

const statusLabels: Record<number, string> = {
  1: '待付款', 2: '待发货', 3: '待收货', 4: '已完成', 5: '售后'
}
const statusColors: Record<number, string> = {
  1: 'text-yellow-500', 2: 'text-blue-500',
  3: 'text-purple-500', 4: 'text-green-500', 5: 'text-red-500'
}
const statusLabel = (s: number) => statusLabels[s] || ''
const statusColor = (s: number) => statusColors[s] || 'text-gray-400'

async function loadOrders() {
  const status = statusValues[activeTab.value]
  const data = await get<any>('/api/v1/orders', { status: status || undefined, page: 1, size: 20 })
  orders.value = data?.list || []
}

function onTab(index: number) {
  activeTab.value = index
  loadOrders()
}

function toPay(order: any) {
  uni.showToast({ title: '支付功能将在 Phase 3 实现', icon: 'none' })
}

onMounted(loadOrders)
</script>
