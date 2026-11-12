<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'

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

const TYPE: Record<string, { icon: string; color: string; soft: string; label: string; labelColor: string }> = {
  info: { icon: 'info', color: 'primary', soft: 'var(--eapp-primary-soft)', label: '通知', labelColor: 'var(--eapp-primary)' },
  warning: { icon: 'warning', color: 'warning', soft: 'var(--eapp-warning-soft)', label: '提醒', labelColor: 'var(--eapp-warning)' },
  urgent: { icon: 'close-circle', color: 'danger', soft: 'var(--eapp-danger-soft)', label: '紧急', labelColor: 'var(--eapp-danger)' },
}

function typeOf(type?: string) {
  return TYPE[type || 'info'] || TYPE.info
}
</script>

<template>
  <view class="announcement-bar">
    <view
      v-for="item in items"
      :key="item.id"
      class="announcement-bar__item"
      hover-class="eapp-hover"
      :hover-stay-time="80"
      @click="emit('click', item.id)"
    >
      <view class="announcement-bar__icon" :style="{ background: typeOf(item.type).soft }">
        <AppIcon :name="typeOf(item.type).icon" :size="30" :color="typeOf(item.type).color" />
      </view>
      <view class="announcement-bar__content">
        <text class="announcement-bar__label" :style="{ color: typeOf(item.type).labelColor }">
          {{ typeOf(item.type).label }}
        </text>
        <text class="announcement-bar__text">{{ item.title }}</text>
      </view>
      <AppIcon name="chevron-right" :size="26" color="faint" />
    </view>
  </view>
</template>

<style scoped>
.announcement-bar {
  display: flex;
  flex-direction: column;
  gap: var(--eapp-space-2);
}

.announcement-bar__item {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-3) var(--eapp-space-4);
  box-sizing: border-box;
}

.announcement-bar__icon {
  width: 56rpx;
  height: 56rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius);
  box-sizing: border-box;
}

.announcement-bar__content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2rpx;
}

.announcement-bar__label {
  font-size: 18rpx;
  font-weight: 600;
  letter-spacing: 1rpx;
}

.announcement-bar__text {
  font-size: var(--eapp-fs-base);
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
