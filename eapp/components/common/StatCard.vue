<script setup lang="ts">
defineProps<{ title: string; value: string; color?: string }>()
</script>

<template>
  <view class="stat-card">
    <view class="stat-title">{{ title }}</view>
    <view class="stat-value" :style="{ color: color || 'var(--eapp-text)' }">{{ value }}</view>
  </view>
</template>

<style scoped>
.stat-card {
  background: var(--eapp-card);
  border-radius: 24rpx;
  padding: 24rpx;
  border: 1px solid var(--eapp-border);
}

.stat-title {
  font-size: 24rpx;
  color: var(--eapp-text-muted);
}

.stat-value {
  margin-top: 12rpx;
  font-size: 40rpx;
  font-weight: 700;
}
</style>
