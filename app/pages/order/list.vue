<template>
  <view class="min-h-screen bg-gray-50">
    <!-- Status tabs -->
    <view class="bg-white">
      <u-tabs :list="tabs" :current="activeTab" @click="(item:any) => onTab(item.index)" />
    </view>

    <view class="p-20rpx">
      <view v-if="!orders.length" class="flex flex-col items-center py-120rpx">
        <u-icon name="order" size="60" color="#ccc" />
        <text class="text-gray-400 text-28rpx mt-20rpx">暂无订单</text>
      </view>

      <view v-for="o in orders" :key="o.id"
        class="bg-white rounded-20rpx p-30rpx mb-20rpx"
        style="box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);">
        <!-- Header -->
        <view class="flex items-center justify-between mb-20rpx">
          <text class="text-22rpx text-gray-400 font-mono">{{ o.order_no }}</text>
          <text :class="statusColor(o.status)" class="text-24rpx font-500">{{ statusLabel(o.status) }}</text>
        </view>
        <!-- Amount + date -->
        <view class="flex items-center justify-between">
          <text class="text-24rpx text-gray-500">{{ o.created_at?.slice(0, 10) }}</text>
          <text class="text-30rpx text-gray-800 font-700">¥{{ o.total_amount }}</text>
        </view>
        <!-- Actions -->
        <view class="flex justify-end gap-16rpx mt-24rpx" v-if="o.status === 1 || o.status === 4">
          <u-button v-if="o.status === 1" size="mini" type="primary" text="去付款" shape="circle" @click="toPay(o)" />
          <u-button v-if="o.status === 4" size="mini" type="success" text="评价" shape="circle" plain />
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
  1: 'text-orange-500', 2: 'text-blue-500',
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

function toPay(_order: any) {
  uni.showToast({ title: '支付功能开发中', icon: 'none' })
}

onMounted(loadOrders)
</script>
