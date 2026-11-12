<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'

export interface EntryItem {
  key: string
  title: string
  desc?: string
  /** 线性图标名，见 utils/icons.ts */
  iconName?: string
  badge?: string | number
  path?: string
}

defineProps<{
  items: EntryItem[]
}>()

const emit = defineEmits<{
  click: [item: EntryItem]
}>()
</script>

<template>
  <view class="entry-list">
    <view
      v-for="item in items"
      :key="item.key"
      class="entry-list__item"
      hover-class="eapp-hover"
      :hover-stay-time="80"
      @click="emit('click', item)"
    >
      <view class="entry-list__icon">
        <AppIcon :name="item.iconName || 'circle'" :size="36" color="primary" />
      </view>
      <view class="entry-list__content">
        <text class="entry-list__title">{{ item.title }}</text>
        <text v-if="item.desc" class="entry-list__desc">{{ item.desc }}</text>
      </view>
      <text v-if="item.badge" class="entry-list__badge eapp-num">{{ item.badge }}</text>
      <AppIcon name="chevron-right" :size="28" color="faint" />
    </view>
  </view>
</template>

<style scoped>
.entry-list {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  overflow: hidden;
  box-sizing: border-box;
}

.entry-list__item {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
  padding: var(--eapp-space-4) var(--eapp-space-5);
  border-bottom: 1rpx solid var(--eapp-divider);
  box-sizing: border-box;
}

.entry-list__item:last-child {
  border-bottom: none;
}

.entry-list__icon {
  width: 68rpx;
  height: 68rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary-soft);
  border: 1rpx solid var(--eapp-primary-border);
  box-sizing: border-box;
}

.entry-list__content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.entry-list__title {
  font-size: var(--eapp-fs-md);
  font-weight: 600;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.entry-list__desc {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.entry-list__badge {
  flex-shrink: 0;
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-primary);
  background: var(--eapp-primary-soft);
  border-radius: var(--eapp-radius-sm);
  padding: 4rpx 12rpx;
}
</style>
