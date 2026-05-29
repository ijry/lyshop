<script setup lang="ts">
import type { DateRange } from '@/composables/useAnalytics'

const props = defineProps<{ modelValue: DateRange }>()
const emit = defineEmits<{ 'update:modelValue': [value: DateRange] }>()

const options: Array<{ label: string; value: DateRange }> = [
  { label: '今日', value: '1' },
  { label: '7日', value: '7' },
  { label: '30日', value: '30' },
  { label: '90日', value: '90' },
]

function onSelect(val: DateRange) {
  emit('update:modelValue', val)
}
</script>

<template>
  <view class="date-range-picker">
    <text
      v-for="opt in options"
      :key="opt.value"
      :class="['pill', modelValue === opt.value ? 'active' : '']"
      @click="onSelect(opt.value)"
    >{{ opt.label }}</text>
  </view>
</template>

<style scoped>
.date-range-picker {
  display: flex;
  gap: 8rpx;
  padding: 16rpx 20rpx;
  background: var(--eapp-card);
  position: sticky;
  top: 0;
  z-index: 20;
}
.pill {
  padding: 10rpx 24rpx;
  border-radius: 999rpx;
  font-size: 24rpx;
  color: var(--eapp-text-muted);
  background: var(--eapp-bg);
}
.pill.active {
  background: var(--eapp-primary);
  color: #fff;
}
</style>
