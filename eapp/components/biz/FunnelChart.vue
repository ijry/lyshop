<script setup lang="ts">
import { computed } from 'vue'
import ChartPanel from '@/components/charts/ChartPanel.vue'

export interface FunnelStep {
  step: string
  count: number
  rate: number
}

const props = defineProps<{
  title?: string
  loading?: boolean
  error?: string
  steps?: FunnelStep[]
}>()

const isEmpty = computed(() => !props.steps?.length)
const maxCount = computed(() => {
  if (!props.steps?.length) return 1
  return Math.max(1, props.steps[0].count)
})
</script>

<template>
  <ChartPanel :title="title" :loading="loading" :error="error" :empty="isEmpty">
    <view v-if="!isEmpty" class="funnel">
      <view
        v-for="(item, idx) in steps"
        :key="idx"
        class="funnel-bar"
      >
        <view
          class="funnel-bar__fill"
          :style="{
            width: `${Math.max(20, (item.count / maxCount) * 100)}%`,
            opacity: 1 - idx * 0.15,
          }"
        >
          <text class="funnel-bar__label">{{ item.step }}</text>
          <text class="funnel-bar__value">{{ item.count.toLocaleString() }} ({{ item.rate }}%)</text>
        </view>
      </view>
    </view>
  </ChartPanel>
</template>

<style scoped>
.funnel {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}
.funnel-bar {
  width: 100%;
  display: flex;
  justify-content: center;
}
.funnel-bar__fill {
  background: var(--eapp-primary);
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: width 0.3s;
}
.funnel-bar__label {
  font-size: 24rpx;
  color: #fff;
  font-weight: 600;
}
.funnel-bar__value {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.85);
}
</style>
