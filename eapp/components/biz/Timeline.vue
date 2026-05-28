<script setup lang="ts">
export interface TimelineItem {
  key: string
  title: string
  time?: string
  desc?: string
  tone?: 'primary' | 'success' | 'warning' | 'danger'
}

defineProps<{
  items: TimelineItem[]
  compact?: boolean
}>()

function toneDotColor(tone?: string) {
  const map: Record<string, string> = {
    primary: 'var(--eapp-primary)',
    success: 'var(--eapp-success)',
    warning: 'var(--eapp-warning)',
    danger: 'var(--eapp-danger)',
  }
  return map[tone || ''] || 'var(--eapp-primary)'
}
</script>

<template>
  <view class="timeline" :class="{ 'timeline--compact': compact }">
    <view
      v-for="(item, idx) in items"
      :key="item.key"
      class="timeline__item"
    >
      <view class="timeline__dot-col">
        <view
          class="timeline__dot"
          :style="{ backgroundColor: toneDotColor(item.tone) }"
        />
        <view
          v-if="idx < items.length - 1"
          class="timeline__line"
        />
      </view>
      <view class="timeline__content">
        <view class="timeline__header">
          <text class="timeline__title">{{ item.title }}</text>
          <text v-if="item.time" class="timeline__time">{{ item.time }}</text>
        </view>
        <text v-if="item.desc" class="timeline__desc">{{ item.desc }}</text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.timeline {
  display: flex;
  flex-direction: column;
}

.timeline__item {
  display: flex;
  gap: 20rpx;
  min-height: 80rpx;
}

.timeline--compact .timeline__item {
  min-height: 56rpx;
}

.timeline__dot-col {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 24rpx;
  flex-shrink: 0;
}

.timeline__dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 999rpx;
  margin-top: 8rpx;
  flex-shrink: 0;
}

.timeline__line {
  width: 2rpx;
  flex: 1;
  background: var(--eapp-border);
  margin-top: 8rpx;
}

.timeline__content {
  flex: 1;
  padding-bottom: 24rpx;
}

.timeline--compact .timeline__content {
  padding-bottom: 12rpx;
}

.timeline__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}

.timeline__title {
  font-size: 26rpx;
  color: var(--eapp-text);
  font-weight: 500;
}

.timeline__time {
  font-size: 22rpx;
  color: var(--eapp-text-muted);
  flex-shrink: 0;
}

.timeline__desc {
  font-size: 24rpx;
  color: var(--eapp-text-muted);
  margin-top: 8rpx;
}
</style>
