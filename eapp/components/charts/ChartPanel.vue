<script setup lang="ts">
defineProps<{
  title?: string
  loading?: boolean
  empty?: boolean
  error?: string
}>()
</script>

<template>
  <view class="chart-panel">
    <view v-if="title || $slots.extra" class="chart-panel__header">
      <text v-if="title" class="chart-panel__title">{{ title }}</text>
      <view class="chart-panel__extra">
        <slot name="extra" />
      </view>
    </view>
    <view class="chart-panel__body">
      <view v-if="loading" class="chart-panel__state">
        <u-loading-icon size="48" />
        <text class="chart-panel__state-text">加载中...</text>
      </view>
      <view v-else-if="error" class="chart-panel__state chart-panel__state--error">
        <text class="chart-panel__state-text">{{ error }}</text>
      </view>
      <view v-else-if="empty" class="chart-panel__state">
        <text class="chart-panel__state-text">暂无数据</text>
      </view>
      <slot v-else />
    </view>
  </view>
</template>

<style scoped>
.chart-panel {
  background: var(--eapp-card);
  border-radius: 24rpx;
  padding: 28rpx;
  margin-bottom: 24rpx;
}

.chart-panel__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.chart-panel__title {
  font-size: 28rpx;
  font-weight: 600;
  color: var(--eapp-text);
}

.chart-panel__extra {
  display: flex;
  align-items: center;
}

.chart-panel__body {
  min-height: 200rpx;
}

.chart-panel__state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 200rpx;
  gap: 16rpx;
}

.chart-panel__state--error .chart-panel__state-text {
  color: var(--eapp-danger);
}

.chart-panel__state-text {
  font-size: 24rpx;
  color: var(--eapp-text-muted);
}
</style>
