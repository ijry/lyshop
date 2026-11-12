<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ title: string; value: string; color?: string }>()

const LEGACY_COLORS: Record<string, string> = {
  '#2563eb': 'var(--eapp-primary)',
  '#16a34a': 'var(--eapp-success)',
  '#f59e0b': 'var(--eapp-warning)',
  '#dc2626': 'var(--eapp-danger)',
}
const dotColor = computed(() =>
  LEGACY_COLORS[props.color || ''] || props.color || 'var(--eapp-primary)',
)

/** 长数值在窄屏卡片里会被 ellipsis 截断，按字符数收敛字号，短数值保持主题字号 */
const valueFontSize = computed(() => {
  const len = String(props.value ?? '').length
  if (len >= 11) return '24rpx'
  if (len >= 9) return '28rpx'
  if (len >= 7) return '32rpx'
  return undefined
})
</script>

<template>
  <view class="stat-card">
    <view class="stat-card__head">
      <view class="stat-card__dot" :style="{ background: dotColor }" />
      <text class="stat-card__title">{{ title }}</text>
    </view>
    <text class="stat-card__value eapp-num" :style="{ fontSize: valueFontSize }">{{ value }}</text>
  </view>
</template>

<style scoped>
.stat-card {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-4);
  box-sizing: border-box;
  min-width: 0;
}

.stat-card__head {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.stat-card__dot {
  width: 10rpx;
  height: 10rpx;
  border-radius: 50%;
  flex-shrink: 0;
}

.stat-card__title {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stat-card__value {
  display: block;
  margin-top: 12rpx;
  font-size: var(--eapp-fs-2xl);
  font-weight: 700;
  line-height: 1.2;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
