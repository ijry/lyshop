<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { getDistributors, getCommissions } from '@/api/distribution'

const { t } = useI18n()
const loading = ref(false)
const totalDistributors = ref(0)
const pendingCommission = ref(0)
const settledTotal = ref(0)

async function loadSummary() {
  loading.value = true
  try {
    const [distRes, commRes] = await Promise.all([
      getDistributors({ page: 1, size: 1 }),
      getCommissions({ page: 1, size: 100 }),
    ])
    totalDistributors.value = Number(distRes?.total || 0)
    const commList = Array.isArray(commRes?.list) ? commRes.list : []
    pendingCommission.value = commList.filter((c: any) => Number(c.status) === 1).reduce((s: number, c: any) => s + Number(c.amount || 0), 0)
    settledTotal.value = commList.filter((c: any) => Number(c.status) === 2).reduce((s: number, c: any) => s + Number(c.amount || 0), 0)
  } finally {
    loading.value = false
  }
}

const entries = [
  { key: 'distributorMgmt', path: '/pages/distribution/distributor-list', color: '#2563eb' },
  { key: 'commissionMgmt', path: '/pages/distribution/commission-list', color: '#16a34a' },
  { key: 'distributionConfig', path: '/pages/distribution/config', color: '#f97316' },
]

function go(path: string) { uni.navigateTo({ url: path }) }

onShow(() => loadSummary())
</script>

<template>
  <view class="page">
    <view class="metrics">
      <view class="metric-card">
        <text class="metric-value">{{ totalDistributors }}</text>
        <text class="metric-label">{{ t('distribution.totalDistributors') }}</text>
      </view>
      <view class="metric-card">
        <text class="metric-value">¥{{ pendingCommission.toFixed(2) }}</text>
        <text class="metric-label">{{ t('distribution.pendingCommission') }}</text>
      </view>
      <view class="metric-card">
        <text class="metric-value">¥{{ settledTotal.toFixed(2) }}</text>
        <text class="metric-label">{{ t('distribution.settledTotal') }}</text>
      </view>
    </view>
    <view class="grid">
      <view v-for="entry in entries" :key="entry.key" class="card" @click="go(entry.path)">
        <view class="name" :style="{ color: entry.color }">{{ t(`distribution.${entry.key}`) }}</view>
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
