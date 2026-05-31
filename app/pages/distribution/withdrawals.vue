<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getMyWithdrawals } from '@/api/distribution'

const loading = ref(false)
const list = ref<any[]>([])
const page = ref(1)
const size = ref(20)
const hasMore = ref(true)

async function loadData() {
  if (loading.value) return
  loading.value = true
  try {
    const res = await getMyWithdrawals({ page: page.value, size: size.value })
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
  pending: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
  completed: '已完成'
}

onMounted(() => loadData())
</script>

<template>
  <view class="page" @scrolltolower="onReachBottom">
    <view v-if="list.length === 0 && !loading" class="empty">暂无提现记录</view>

    <view v-for="item in list" :key="item.id" class="withdrawal-item">
      <view class="header">
        <text class="amount">¥{{ (item.amount || 0).toFixed(2) }}</text>
        <text class="status" :class="item.status">{{ statusMap[item.status] }}</text>
      </view>
      <view class="body">
        <view class="info-row">
          <text class="label">手续费：</text>
          <text class="value">¥{{ (item.fee || 0).toFixed(2) }}</text>
        </view>
        <view class="info-row">
          <text class="label">实际到账：</text>
          <text class="value">¥{{ (item.actual_amount || 0).toFixed(2) }}</text>
        </view>
        <view class="info-row">
          <text class="label">银行：</text>
          <text class="value">{{ item.bank_name }}</text>
        </view>
        <view class="info-row">
          <text class="label">账号：</text>
          <text class="value">{{ item.bank_account }}</text>
        </view>
        <view class="info-row">
          <text class="label">户名：</text>
          <text class="value">{{ item.account_name }}</text>
        </view>
        <view v-if="item.reject_reason" class="info-row">
          <text class="label">拒绝原因：</text>
          <text class="value error">{{ item.reject_reason }}</text>
        </view>
        <view class="info-row">
          <text class="label">申请时间：</text>
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

.withdrawal-item { background: #fff; border-radius: 16rpx; padding: 24rpx; margin-bottom: 20rpx; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20rpx; padding-bottom: 20rpx; border-bottom: 1px solid #f0f0f0; }
.amount { font-size: 36rpx; font-weight: bold; color: #dc2626; }
.status { font-size: 24rpx; padding: 6rpx 16rpx; border-radius: 12rpx; }
.status.pending { background: #fff7ed; color: #f97316; }
.status.approved { background: #dbeafe; color: #2563eb; }
.status.rejected { background: #fee2e2; color: #dc2626; }
.status.completed { background: #dcfce7; color: #16a34a; }

.body { }
.info-row { display: flex; justify-content: space-between; margin-bottom: 16rpx; font-size: 26rpx; }
.label { color: #999; }
.value { color: #333; }
.value.error { color: #dc2626; }

.loading, .no-more { text-align: center; padding: 20rpx; color: #999; font-size: 24rpx; }
</style>
