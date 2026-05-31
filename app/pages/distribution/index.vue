<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDistributorInfo } from '@/api/distribution'

const loading = ref(false)
const distributor = ref<any>(null)

async function loadData() {
  loading.value = true
  try {
    const res = await getDistributorInfo()
    distributor.value = res
  } catch (error: any) {
    console.error('加载失败', error)
  } finally {
    loading.value = false
  }
}

function goToApply() {
  uni.navigateTo({ url: '/pages/distribution/apply' })
}

function goToTeam() {
  uni.navigateTo({ url: '/pages/distribution/team' })
}

function goToOrders() {
  uni.navigateTo({ url: '/pages/distribution/orders' })
}

function goToWithdrawals() {
  uni.navigateTo({ url: '/pages/distribution/withdrawals' })
}

function goToWithdraw() {
  uni.navigateTo({ url: '/pages/distribution/withdraw' })
}

onMounted(() => loadData())
</script>

<template>
  <view class="page">
    <view v-if="loading" class="loading">加载中...</view>

    <view v-else-if="!distributor" class="empty">
      <text class="empty-text">您还不是分销商</text>
      <button class="apply-btn" @click="goToApply">立即申请</button>
    </view>

    <view v-else class="content">
      <view class="header">
        <view class="status-badge" :class="distributor.status">
          {{ distributor.status === 'active' ? '已激活' : distributor.status === 'pending' ? '待审核' : '已禁用' }}
        </view>
        <text class="name">{{ distributor.real_name }}</text>
        <text class="level">等级：{{ distributor.level }}</text>
      </view>

      <view class="stats">
        <view class="stat-item">
          <text class="stat-value">¥{{ (distributor.total_earnings || 0).toFixed(2) }}</text>
          <text class="stat-label">累计收益</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">¥{{ (distributor.available_amount || 0).toFixed(2) }}</text>
          <text class="stat-label">可提现</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ distributor.total_customers || 0 }}</text>
          <text class="stat-label">客户数</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ distributor.total_orders || 0 }}</text>
          <text class="stat-label">订单数</text>
        </view>
      </view>

      <view class="actions">
        <view class="action-item" @click="goToTeam">
          <text class="action-icon">👥</text>
          <text class="action-label">我的团队</text>
        </view>
        <view class="action-item" @click="goToOrders">
          <text class="action-icon">📦</text>
          <text class="action-label">分销订单</text>
        </view>
        <view class="action-item" @click="goToWithdrawals">
          <text class="action-icon">💰</text>
          <text class="action-label">提现记录</text>
        </view>
      </view>

      <button v-if="distributor.status === 'active' && distributor.available_amount > 0"
              class="withdraw-btn" @click="goToWithdraw">
        申请提现
      </button>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: #f5f5f5; }
.loading, .empty { padding: 100rpx 40rpx; text-align: center; }
.empty-text { display: block; font-size: 32rpx; color: #999; margin-bottom: 40rpx; }
.apply-btn { width: 400rpx; height: 80rpx; line-height: 80rpx; background: #dc2626; color: #fff; border-radius: 40rpx; border: none; }

.content { padding: 20rpx; }
.header { background: linear-gradient(135deg, #dc2626, #ef4444); padding: 40rpx; border-radius: 20rpx; color: #fff; margin-bottom: 20rpx; }
.status-badge { display: inline-block; padding: 8rpx 20rpx; background: rgba(255,255,255,0.3); border-radius: 20rpx; font-size: 24rpx; margin-bottom: 20rpx; }
.name { display: block; font-size: 40rpx; font-weight: bold; margin-bottom: 10rpx; }
.level { display: block; font-size: 28rpx; opacity: 0.9; }

.stats { display: grid; grid-template-columns: 1fr 1fr 1fr 1fr; gap: 20rpx; margin-bottom: 20rpx; }
.stat-item { background: #fff; padding: 30rpx 20rpx; border-radius: 16rpx; text-align: center; }
.stat-value { display: block; font-size: 32rpx; font-weight: bold; color: #333; margin-bottom: 10rpx; }
.stat-label { display: block; font-size: 24rpx; color: #999; }

.actions { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 20rpx; margin-bottom: 20rpx; }
.action-item { background: #fff; padding: 40rpx 20rpx; border-radius: 16rpx; text-align: center; }
.action-icon { display: block; font-size: 48rpx; margin-bottom: 10rpx; }
.action-label { display: block; font-size: 28rpx; color: #333; }

.withdraw-btn { width: 100%; height: 88rpx; line-height: 88rpx; background: #dc2626; color: #fff; border-radius: 16rpx; border: none; font-size: 32rpx; }
</style>
