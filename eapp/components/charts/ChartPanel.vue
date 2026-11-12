<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'

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
      <view class="chart-panel__title-wrap">
        <view class="chart-panel__bar" />
        <text v-if="title" class="chart-panel__title">{{ title }}</text>
      </view>
      <view class="chart-panel__extra">
        <slot name="extra" />
      </view>
    </view>
    <view class="chart-panel__body">
      <view v-if="loading" class="chart-panel__state">
        <up-loading-icon mode="circle" size="36" />
        <text class="chart-panel__state-text">加载中…</text>
      </view>
      <view v-else-if="error" class="chart-panel__state">
        <AppIcon name="warning" :size="48" color="danger" />
        <text class="chart-panel__state-text chart-panel__state-text--error">{{ error }}</text>
      </view>
      <view v-else-if="empty" class="chart-panel__state">
        <AppIcon name="chart" :size="48" color="faint" :stroke-width="1.4" />
        <text class="chart-panel__state-text">暂无数据</text>
      </view>
      <slot v-else />
    </view>
  </view>
</template>

<style scoped>
.chart-panel {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-5);
  margin-bottom: var(--eapp-space-5);
  box-sizing: border-box;
}

.chart-panel__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
  margin-bottom: var(--eapp-space-4);
}

.chart-panel__title-wrap {
  display: flex;
  align-items: center;
  gap: 12rpx;
  min-width: 0;
}

.chart-panel__bar {
  width: 6rpx;
  height: 28rpx;
  border-radius: 3rpx;
  background: var(--eapp-primary);
  flex-shrink: 0;
}

.chart-panel__title {
  font-size: var(--eapp-fs-md);
  font-weight: 600;
  color: var(--eapp-text);
}

.chart-panel__extra {
  display: flex;
  align-items: center;
  flex-shrink: 0;
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

.chart-panel__state-text {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-faint);
}

.chart-panel__state-text--error {
  color: var(--eapp-danger);
}
</style>
