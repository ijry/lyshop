<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getMyTeam } from '@/api/distribution'

const loading = ref(false)
const team = ref<any[]>([])

async function loadData() {
  loading.value = true
  try {
    const res = await getMyTeam()
    team.value = res.team || []
  } catch (error: any) {
    uni.showToast({ title: error.message || '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

onMounted(() => loadData())
</script>

<template>
  <view class="page">
    <view v-if="loading" class="loading">加载中...</view>
    <view v-else-if="team.length === 0" class="empty">暂无团队成员</view>

    <view v-for="item in team" :key="item.id" class="member-item">
      <view class="member-header">
        <text class="name">{{ item.real_name }}</text>
        <text class="level">等级 {{ item.level }}</text>
      </view>
      <view class="member-body">
        <view class="info-row">
          <text class="label">累计收益：</text>
          <text class="value">¥{{ (item.total_earnings || 0).toFixed(2) }}</text>
        </view>
        <view class="info-row">
          <text class="label">客户数：</text>
          <text class="value">{{ item.total_customers || 0 }}</text>
        </view>
        <view class="info-row">
          <text class="label">订单数：</text>
          <text class="value">{{ item.total_orders || 0 }}</text>
        </view>
        <view class="info-row">
          <text class="label">加入时间：</text>
          <text class="value">{{ item.created_at }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: #f5f5f5; padding: 20rpx; }
.loading, .empty { padding: 100rpx 40rpx; text-align: center; color: #999; }

.member-item { background: #fff; border-radius: 16rpx; padding: 24rpx; margin-bottom: 20rpx; }
.member-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20rpx; padding-bottom: 20rpx; border-bottom: 1px solid #f0f0f0; }
.name { font-size: 30rpx; font-weight: bold; color: #333; }
.level { font-size: 24rpx; padding: 6rpx 16rpx; border-radius: 12rpx; background: #dbeafe; color: #2563eb; }

.member-body { }
.info-row { display: flex; justify-content: space-between; margin-bottom: 16rpx; font-size: 26rpx; }
.label { color: #999; }
.value { color: #333; }
</style>
