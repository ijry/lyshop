<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const cards = [
  { key: 'coupon', path: '/pages/marketing/coupon', color: '#2563eb' },
  { key: 'seckill', path: '/pages/marketing/seckill', color: '#ef4444' },
  { key: 'groupBuy', path: '/pages/marketing/group-buy', color: '#16a34a' },
  { key: 'bargain', path: '/pages/marketing/bargain', color: '#f97316' },
  { key: 'vip', path: '/pages/vip/plan-list', color: '#8b5cf6' },
  { key: 'decor', path: '/pages/decor/editor', color: '#0ea5e9' },
  { key: 'specTemplate', path: '/pages/product/spec-templates', color: '#64748b' },
  { key: 'distribution', path: '/pages/distribution/index', color: '#d946ef' },
]

function go(path: string) {
  uni.navigateTo({ url: path })
}
</script>

<template>
  <view class="page">
    <view class="grid">
      <view v-for="card in cards" :key="card.key" class="card" @click="go(card.path)">
        <view class="name" :style="{ color: card.color }">{{ t(`marketing.${card.key}`) }}</view>
        <view class="desc">点击进入管理</view>
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
