<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { getPointsSummary } from '@/api/points'

const { t } = useI18n()
const loading = ref(false)
const totalIssued = ref(0)
const totalConsumed = ref(0)
const productCount = ref(0)

async function loadSummary() {
  loading.value = true
  try {
    const res = await getPointsSummary()
    totalIssued.value = Number(res?.total_issued || 0)
    totalConsumed.value = Number(res?.total_consumed || 0)
    productCount.value = Number(res?.product_count || 0)
  } finally {
    loading.value = false
  }
}

const entries = [
  { key: 'productMgmt', path: '/pages/points/product-list', color: '#2563eb' },
  { key: 'exchangeRecords', path: '/pages/points/exchange-list', color: '#16a34a' },
  { key: 'checkinMgmt', path: '/pages/checkin/rules', color: '#f97316' },
]

function go(path: string) { uni.navigateTo({ url: path }) }

onShow(() => loadSummary())
</script>

<template>
  <view class="page">
    <view class="metrics">
      <view class="metric-card">
        <text class="metric-value">{{ totalIssued }}</text>
        <text class="metric-label">{{ t('points.totalIssued') }}</text>
      </view>
      <view class="metric-card">
        <text class="metric-value">{{ totalConsumed }}</text>
        <text class="metric-label">{{ t('points.totalConsumed') }}</text>
      </view>
      <view class="metric-card">
        <text class="metric-value">{{ productCount }}</text>
        <text class="metric-label">{{ t('points.productCount') }}</text>
      </view>
    </view>
    <view class="grid">
      <view v-for="entry in entries" :key="entry.key" class="card" @click="go(entry.path)">
        <view class="name" :style="{ color: entry.color }">{{ t(`points.${entry.key}`) }}</view>
        <view class="desc">点击进入管理</view>
      </view>
    </view>
    <view v-if="loading" class="loading">加载中…</view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; }
.metrics { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 12rpx; margin-bottom: 20rpx; }
.metric-card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 16rpx; padding: 20rpx; text-align: center; }
.metric-value { display: block; font-size: 32rpx; font-weight: 700; color: #1e293b; }
.metric-label { display: block; margin-top: 6rpx; font-size: 22rpx; color: var(--eapp-text-muted); }
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 26rpx 20rpx; }
.name { font-size: 30rpx; font-weight: 700; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 22rpx; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>
