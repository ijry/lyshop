<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  title: string
  value: string | number
  unit?: string
  compare?: number
  color?: string
}>()

const compareText = computed(() => {
  if (props.compare == null) return ''
  const pct = props.compare
  return pct >= 0 ? `+${pct}%` : `${pct}%`
})

const compareClass = computed(() => {
  if (props.compare == null) return ''
  return props.compare >= 0 ? 'metric-card__compare--up' : 'metric-card__compare--down'
})
</script>

<template>
  <view class="metric-card" :style="color ? { borderTopColor: color } : {}">
    <text class="metric-card__title">{{ title }}</text>
    <view class="metric-card__value-row">
      <text class="metric-card__value" :style="color ? { color } : {}">{{ value }}</text>
      <text v-if="unit" class="metric-card__unit">{{ unit }}</text>
    </view>
    <text v-if="compare != null" class="metric-card__compare" :class="compareClass">
      {{ compareText }}
    </text>
  </view>
</template>

<style scoped>
.metric-card {
  background: var(--eapp-card);
  border-radius: 22rpx;
  padding: 24rpx;
  border-top: 4rpx solid var(--eapp-primary);
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.metric-card__title {
  font-size: 24rpx;
  color: var(--eapp-text-muted);
}

.metric-card__value-row {
  display: flex;
  align-items: baseline;
  gap: 6rpx;
}

.metric-card__value {
  font-size: 40rpx;
  font-weight: 700;
  color: var(--eapp-text);
}

.metric-card__unit {
  font-size: 22rpx;
  color: var(--eapp-text-muted);
}

.metric-card__compare {
  font-size: 22rpx;
}

.metric-card__compare--up {
  color: var(--eapp-success);
}

.metric-card__compare--down {
  color: var(--eapp-danger);
}
</style>
