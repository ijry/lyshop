<script setup lang="ts">
export interface AnnouncementItem {
  id: string | number
  title: string
  type?: 'info' | 'warning' | 'urgent'
}

defineProps<{
  items: AnnouncementItem[]
}>()

const emit = defineEmits<{
  click: [id: string | number]
}>()

function typeColor(type?: string) {
  const map: Record<string, string> = {
    info: 'var(--eapp-primary)',
    warning: 'var(--eapp-warning)',
    urgent: 'var(--eapp-danger)',
  }
  return map[type || ''] || 'var(--eapp-primary)'
}
</script>

<template>
  <scroll-view class="announcement-bar" scroll-x :show-scrollbar="false">
    <view class="announcement-bar__inner">
      <view
        v-for="item in items"
        :key="item.id"
        class="announcement-bar__item"
        :style="{ borderLeftColor: typeColor(item.type) }"
        @click="emit('click', item.id)"
      >
        <text class="announcement-bar__text">{{ item.title }}</text>
      </view>
    </view>
  </scroll-view>
</template>

<style scoped>
.announcement-bar {
  white-space: nowrap;
  width: 100%;
}

.announcement-bar__inner {
  display: inline-flex;
  gap: 16rpx;
  padding: 8rpx 0;
}

.announcement-bar__item {
  display: inline-flex;
  align-items: center;
  background: var(--eapp-card);
  border-radius: 16rpx;
  padding: 16rpx 24rpx;
  border-left: 6rpx solid var(--eapp-primary);
  min-width: 280rpx;
  min-height: 44rpx;
}

.announcement-bar__text {
  font-size: 24rpx;
  color: var(--eapp-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 400rpx;
}
</style>
