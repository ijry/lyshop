<template>
  <view class="min-h-screen bg-gray-50">
    <view class="bg-white">
      <u-tabs :list="tabs" :current="activeTab" @click="(item:any) => onTab(item.index)" />
    </view>

    <view class="p-20rpx">
      <view v-if="!orders.length" class="flex flex-col items-center py-120rpx">
        <u-icon name="order" size="60" color="#ccc" />
        <text class="text-gray-400 text-28rpx mt-20rpx">暂无订单</text>
      </view>

      <view
        v-for="o in orders"
        :key="o.id"
        class="bg-white rounded-20rpx p-30rpx mb-20rpx"
        style="box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);"
      >
        <view class="flex items-center justify-between mb-20rpx">
          <view class="flex items-center" style="gap: 8px;">
            <text class="text-22rpx text-gray-400 font-mono">{{ o.order_no }}</text>
            <text v-if="o.activity_type === 'seckill'" style="font-size: 10px; color: #dc2626; background: #fef2f2; padding: 1px 6px; border-radius: 4px;">秒杀</text>
            <text v-else-if="o.activity_type === 'group_buy'" style="font-size: 10px; color: #2563eb; background: #eff6ff; padding: 1px 6px; border-radius: 4px;">拼团</text>
            <text v-else-if="o.activity_type === 'bargain'" style="font-size: 10px; color: #16a34a; background: #f0fdf4; padding: 1px 6px; border-radius: 4px;">砍价</text>
          </view>
          <text :class="statusColor(o.status)" class="text-24rpx font-500">{{ statusLabel(o.status) }}</text>
        </view>

        <view v-if="o.items?.length" class="mb-16rpx">
          <view v-for="it in o.items.slice(0, 2)" :key="it.id" class="flex items-center mb-12rpx" style="gap: 10px;">
            <image :src="it.cover" mode="aspectFill" style="width: 72rpx; height: 72rpx; border-radius: 10rpx;" />
            <view class="flex-1 min-w-0">
              <text class="text-24rpx text-gray-700 block truncate">{{ it.title }}</text>
              <text class="text-22rpx text-gray-400 mt-4rpx block">x{{ it.qty }}</text>
            </view>
          </view>
        </view>

        <view class="flex items-center justify-between">
          <text class="text-24rpx text-gray-500">{{ o.created_at?.slice(0, 10) }}</text>
          <text class="text-30rpx text-gray-800 font-700">¥{{ money(o.amount_breakdown?.payable_amount ?? o.total_amount) }}</text>
        </view>

        <view class="flex justify-end gap-16rpx mt-24rpx">
          <view class="action-btn-wrap">
            <u-button size="mini" plain text="查看详情" shape="circle" @click="goDetail(o.id)" />
          </view>
          <view class="action-btn-wrap" v-if="o.status === 1">
            <u-button size="mini" type="primary" text="去付款" shape="circle" :loading="actioningID === o.id && actioningType === 'pay'" @click="toPay(o)" />
          </view>
          <view class="action-btn-wrap" v-if="o.status === 3 || o.status === 4">
            <u-button size="mini" type="success" text="评价" shape="circle" plain :loading="actioningID === o.id && actioningType === 'review'" @click="toReview(o)" />
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get, post } from '@/utils/request'

const orders = ref<any[]>([])
const activeTab = ref(0)
const actioningID = ref<number>(0)
const actioningType = ref<'pay' | 'review' | ''>('')

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
const money = (v: number) => Number(v || 0).toFixed(2)

async function loadOrders() {
  const status = statusValues[activeTab.value]
  const data = await get<any>('/api/v1/orders', { status: status || undefined, page: 1, size: 20 })
  orders.value = data?.list || []
}

function onTab(index: number) {
  activeTab.value = index
  loadOrders()
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/order/detail?id=${id}` })
}

async function toPay(order: any) {
  const id = Number(order?.id || 0)
  if (!id || actioningID.value) return
  actioningID.value = id
  actioningType.value = 'pay'
  try {
    await post(`/api/v1/orders/${id}/pay`)
    uni.showToast({ title: '支付成功', icon: 'success' })
    await loadOrders()
  } catch {
    uni.showToast({ title: '支付失败', icon: 'none' })
  } finally {
    actioningID.value = 0
    actioningType.value = ''
  }
}

function toReview(order: any) {
  const id = Number(order?.id || 0)
  if (!id || actioningID.value) return
  uni.showModal({
    title: '订单评价',
    editable: true,
    placeholderText: '请输入评价内容（选填）',
    confirmText: '提交',
    success: async (res: any) => {
      if (!res?.confirm) return
      actioningID.value = id
      actioningType.value = 'review'
      try {
        await post(`/api/v1/orders/${id}/review`, { content: String(res?.content || '').trim() })
        uni.showToast({ title: '评价成功', icon: 'success' })
        await loadOrders()
      } catch {
        uni.showToast({ title: '评价失败', icon: 'none' })
      } finally {
        actioningID.value = 0
        actioningType.value = ''
      }
    },
  } as any)
}

onMounted(loadOrders)
</script>

<style scoped>
.action-btn-wrap {
  display: inline-flex;
}
</style>
