<template>
  <view class="min-h-screen bg-gray-50 pb-40rpx">
    <view v-if="detail.id" class="p-24rpx">
      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <view class="flex items-center justify-between mb-16rpx">
          <text class="text-28rpx font-600 text-gray-800">订单 {{ detail.order_no }}</text>
          <text class="text-24rpx" :class="statusColor(detail.status)">{{ statusLabel(detail.status) }}</text>
        </view>
        <view class="space-y-12rpx text-24rpx text-gray-600">
          <view>支付方式：{{ payLabel(detail.payment_method) }}</view>
          <view>下单时间：{{ formatDate(detail.created_at) }}</view>
          <view v-if="detail.paid_at">支付时间：{{ formatDate(detail.paid_at) }}</view>
          <view v-if="detail.tracking_no">物流单号：{{ detail.tracking_no }}</view>
        </view>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">商品明细</text>
        <view v-if="detail.items?.length" class="space-y-16rpx">
          <view v-for="it in detail.items" :key="it.id" class="flex gap-16rpx">
            <image :src="it.cover" mode="aspectFill" style="width: 120rpx; height: 120rpx; border-radius: 16rpx;" />
            <view class="flex-1 min-w-0">
              <text class="text-26rpx text-gray-800 block line-clamp-2">{{ it.title }}</text>
              <text class="text-22rpx text-gray-400 block mt-8rpx">数量 x{{ it.qty }}</text>
              <text class="text-24rpx text-gray-600 block mt-8rpx">¥{{ money(it.price) }}</text>
            </view>
          </view>
        </view>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">价格体系</text>
        <view class="space-y-12rpx text-24rpx text-gray-600">
          <view class="flex justify-between"><text>商品金额</text><text>¥{{ money(detail.amount_breakdown?.goods_amount ?? detail.goods_amount) }}</text></view>
          <view class="flex justify-between"><text>优惠金额</text><text>-¥{{ money(detail.amount_breakdown?.discount_amount ?? detail.discount_amount) }}</text></view>
          <view class="flex justify-between"><text>运费</text><text>¥{{ money(detail.amount_breakdown?.freight_amount ?? detail.freight_amount) }}</text></view>
          <view class="flex justify-between text-gray-800 font-600 pt-12rpx border-t border-gray-100"><text>实付金额</text><text>¥{{ money(detail.amount_breakdown?.payable_amount ?? detail.total_amount) }}</text></view>
        </view>
      </view>

      <view class="bg-white rounded-20rpx p-24rpx mb-20rpx">
        <text class="text-28rpx font-600 text-gray-800 block mb-16rpx">评价信息</text>
        <view v-if="reviewItems.length" class="space-y-14rpx">
          <view v-for="rv in reviewItems" :key="rv.id" class="border border-gray-100 rounded-16rpx p-16rpx">
            <text class="text-24rpx text-gray-700 block">{{ rv.product_title }}</text>
            <text class="text-22rpx text-gray-400 block mt-6rpx">商品 {{ rv.product_score }} / 物流 {{ rv.logistics_score }}</text>
            <text class="text-24rpx text-gray-600 block mt-10rpx">{{ rv.content || '未填写评价' }}</text>
          </view>
        </view>
        <text v-else class="text-24rpx text-gray-400">暂无评价</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'

const detail = ref<any>({})
const reviewItems = ref<any[]>([])
const statusLabels: Record<number, string> = { 1: '待付款', 2: '待发货', 3: '待收货', 4: '已完成', 5: '售后' }
const statusColors: Record<number, string> = { 1: 'text-orange-500', 2: 'text-blue-500', 3: 'text-purple-500', 4: 'text-green-500', 5: 'text-red-500' }
const statusLabel = (s: number) => statusLabels[s] || ''
const statusColor = (s: number) => statusColors[s] || 'text-gray-400'
const payLabel = (m: string) => m === 'wechat' ? '微信支付' : m === 'alipay' ? '支付宝支付' : '未支付'
const money = (v: number) => Number(v || 0).toFixed(2)
const formatDate = (v: string) => v ? v.slice(0, 19).replace('T', ' ') : '-'

onMounted(async () => {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  detail.value = await get<any>(`/api/v1/orders/${query.id}`)
  reviewItems.value = []
  for (const item of detail.value?.items || []) {
    if (item.review) {
      reviewItems.value.push({
        ...item.review,
        product_title: item.title,
      })
    }
  }
})
</script>
