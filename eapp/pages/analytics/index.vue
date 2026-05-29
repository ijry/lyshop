<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const cards = [
  { key: 'sales', path: '/pages/analytics/sales', color: '#2563eb' },
  { key: 'products', path: '/pages/analytics/products', color: '#16a34a' },
  { key: 'customers', path: '/pages/analytics/customers', color: '#f59e0b' },
  { key: 'traffic', path: '/pages/analytics/traffic', color: '#8b5cf6' },
  { key: 'conversion', path: '/pages/analytics/conversion', color: '#ef4444' },
]

function go(path: string) {
  uni.navigateTo({ url: path })
}
</script>

<template>
  <view class="page">
    <view class="grid">
      <view v-for="card in cards" :key="card.key" class="card" @click="go(card.path)">
        <view class="name" :style="{ color: card.color }">{{ t(`analytics.${card.key}`) }}</view>
        <view class="desc">点击查看详情</view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; }
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 26rpx 20rpx; }
.name { font-size: 30rpx; font-weight: 700; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 22rpx; }
</style>
