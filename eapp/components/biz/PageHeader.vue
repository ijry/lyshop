<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
  transparent?: boolean
}>()
</script>

<template>
  <view class="page-header" :class="{ 'page-header--transparent': transparent }">
    <view class="page-header__left">
      <text class="page-header__title">{{ title }}</text>
      <text v-if="subtitle" class="page-header__subtitle">{{ subtitle }}</text>
    </view>
    <view class="page-header__right">
      <slot name="right" />
    </view>
  </view>
</template>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: env(safe-area-inset-top) 28rpx 20rpx 28rpx;
  padding-top: calc(env(safe-area-inset-top) + 20rpx);
  background: var(--eapp-card);
}

.page-header--transparent {
  background: transparent;
}

.page-header__left {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.page-header__title {
  font-size: 36rpx;
  font-weight: 700;
  color: var(--eapp-text);
}

.page-header__subtitle {
  font-size: 24rpx;
  color: var(--eapp-text-muted);
}

.page-header__right {
  display: flex;
  align-items: center;
  gap: 16rpx;
}
</style>
