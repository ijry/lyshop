<script setup lang="ts">
import { computed } from 'vue'
import AppIcon from '@/components/common/AppIcon.vue'

const props = defineProps<{
  title: string
  value: string | number
  unit?: string
  compare?: number
  color?: string
}>()

/** 兼容旧页面传入的高饱和色值，统一收敛到新主题 */
const LEGACY_COLORS: Record<string, string> = {
  '#2563eb': 'var(--eapp-primary)',
  '#16a34a': 'var(--eapp-success)',
  '#f59e0b': 'var(--eapp-warning)',
  '#f97316': 'var(--eapp-accent)',
  '#dc2626': 'var(--eapp-danger)',
}

const dotColor = computed(() =>
  LEGACY_COLORS[props.color || ''] || props.color || 'var(--eapp-primary)',
)

const isUp = computed(() => (props.compare ?? 0) >= 0)
const compareText = computed(() => {
  if (props.compare == null) return ''
  return props.compare >= 0 ? `+${props.compare}%` : `${props.compare}%`
})

/**
 * 指标卡并排三列时单列可用宽度有限，长金额（如 ¥28960.50）在 393px
 * 手机上会被 ellipsis 截断。按字符数收敛字号，保证数值始终完整可读。
 */
const valueFontSize = computed(() => {
  const len = String(props.value ?? '').length + (props.unit ? String(props.unit).length : 0)
  if (len >= 11) return '24rpx'
  if (len >= 9) return '28rpx'
  if (len >= 7) return '32rpx'
  return '36rpx'
})
</script>

<template>
  <view class="metric-card">
    <view class="metric-card__head">
      <view class="metric-card__dot" :style="{ background: dotColor }" />
      <text class="metric-card__title">{{ title }}</text>
    </view>
    <view class="metric-card__value-row">
      <text class="metric-card__value eapp-num" :style="{ fontSize: valueFontSize }">{{ value }}</text>
      <text v-if="unit" class="metric-card__unit">{{ unit }}</text>
    </view>
    <view v-if="compare != null" class="metric-card__foot" :class="isUp ? 'is-up' : 'is-down'">
      <AppIcon :name="isUp ? 'trend-up' : 'trend-down'" :size="22" :color="isUp ? 'success' : 'danger'" />
      <text class="metric-card__compare eapp-num">{{ compareText }}</text>
      <text class="metric-card__foot-label">较上期</text>
    </view>
    <view v-else class="metric-card__foot metric-card__foot--empty" />
  </view>
</template>

<style scoped>
.metric-card {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-4);
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.metric-card__head {
  display: flex;
  align-items: center;
  gap: 10rpx;
  min-width: 0;
}

.metric-card__dot {
  width: 10rpx;
  height: 10rpx;
  border-radius: 50%;
  flex-shrink: 0;
}

.metric-card__title {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.metric-card__value-row {
  display: flex;
  align-items: baseline;
  gap: 4rpx;
  margin-top: 12rpx;
  min-width: 0;
}

.metric-card__value {
  font-size: 36rpx;
  font-weight: 700;
  line-height: 1.2;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.metric-card__unit {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
  flex-shrink: 0;
}

.metric-card__foot {
  display: flex;
  align-items: center;
  gap: 6rpx;
  margin-top: 10rpx;
  min-height: 24rpx;
}

.metric-card__foot--empty {
  opacity: 0;
}

.metric-card__compare {
  font-size: var(--eapp-fs-xs);
  font-weight: 600;
}

.is-up .metric-card__compare {
  color: var(--eapp-success);
}

.is-down .metric-card__compare {
  color: var(--eapp-danger);
}

.metric-card__foot-label {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}
</style>
