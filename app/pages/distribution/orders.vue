<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getMyOrders } from '@/api/distribution'

const loading = ref(false)
const list = ref<any[]>([])
const page = ref(1)
const size = ref(20)
const hasMore = ref(true)

async function loadData() {
  if (loading.value) return
  loading.value = true
  try {
    const res = await getMyOrders({ page: page.value, size: size.value })
    const newList = res.list || []
    if (page.value === 1) {
      list.value = newList
    } else {
      list.value.push(...newList)
    }
    hasMore.value = newList.length >= size.value
  } catch (error: any) {
    uni.showToast({ title: error.message || '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

function onReachBottom() {
  if (hasMore.value && !loading.value) {
    page.value++
    loadData()
  }
}

function onPullDownRefresh() {
  page.value = 1
  loadData().then(() => {
    uni.stopPullDownRefresh()
  })
}

const statusMap: Record<string, string> = {
  pending: '待结算',
  settled: '已结算',
  cancelled: '已取消'
}

onMounted(() => loadData())
</script>

<template>
  <view class="page" @scrolltolower="onReachBottom">
    <view v-if="list.length === 0 && !loading" class="empty">暂无订单</view>

    <view v-for="item in list" :key="item.id" class="order-item">
      <view class="order-header">
        <text class="order-id">订单 #{{ item.order_id }}</text>
        <text class="status" :class="item.status">{{ statusMap[item.status] }}</text>
      </view>
      <view class="order-body">
        <view class="info-row">
          <text class="label">分销层级：</text>
          <text class="value">{{ item.level }}级</text>
        </view>
        <view class="info-row">
          <text class="label">订单金额：</text>
          <text class="value">¥{{ (item.order_amount || 0).toFixed(2) }}</text>
        </view>
        <view class="info-row">
          <text class="label">佣金比例：</text>
          <text class="value">{{ item.commission_rate }}%</text>
        </view>
        <view class="info-row">
          <text class="label">佣金金额：</text>
          <text class="value commission">¥{{ (item.commission || 0).toFixed(2) }}</text>
        </view>
        <view class="info-row">
          <text class="label">创建时间：</text>
          <text class="value">{{ item.created_at }}</text>
        </view>
      </view>
    </view>

    <view v-if="loading" class="loading">加载中...</view>
    <view v-if="!hasMore && list.length > 0" class="no-more">没有更多了</view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: #f5f5f5; padding: 20rpx; }
.empty { padding: 100rpx 40rpx; text-align: center; color: #999; }

.order-item { background: #fff; border-radius: 16rpx; padding: 24rpx; margin-bottom: 20rpx; }
.order-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20rpx; padding-bottom: 20rpx; border-bottom: 1px solid #f0f0f0; }
.order-id { font-size: 28rpx; font-weight: bold; color: #333; }
.status { font-size: 24rpx; padding: 6rpx 16rpx; border-radius: 12rpx; }
.status.pending { background: #fff7ed; color: #f97316; }
.status.settled { background: #dcfce7; color: #16a34a; }
.status.cancelled { background: #f3f4f6; color: #6b7280; }

.order-body { }
.info-row { display: flex; justify-content: space-between; margin-bottom: 16rpx; font-size: 26rpx; }
.label { color: #999; }
.value { color: #333; }
.value.commission { color: #dc2626; font-weight: bold; }

.loading, .no-more { text-align: center; padding: 20rpx; color: #999; font-size: 24rpx; }
</style>
