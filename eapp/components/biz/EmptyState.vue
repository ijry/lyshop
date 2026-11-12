<script setup lang="ts">
import { computed } from 'vue'
import AppIcon from '@/components/common/AppIcon.vue'

const props = defineProps<{
  title?: string
  desc?: string
  /** 线性图标名；传入历史 emoji 时统一降级为 inbox */
  icon?: string
}>()

const KNOWN = new Set([
  'inbox', 'order', 'package', 'product', 'search', 'document', 'chat', 'chart',
  'list', 'star', 'coupon', 'users', 'warehouse', 'warning', 'info',
])

const iconName = computed(() => (props.icon && KNOWN.has(props.icon) ? props.icon : 'inbox'))
</script>

<template>
  <view class="empty-state">
    <view class="empty-state__illustration">
      <AppIcon :name="iconName" :size="64" color="faint" :stroke-width="1.4" />
    </view>
    <text v-if="title" class="empty-state__title">{{ title }}</text>
    <text v-if="desc" class="empty-state__desc">{{ desc }}</text>
    <view v-if="$slots.default" class="empty-state__action">
      <slot />
    </view>
  </view>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 88rpx 40rpx;
  gap: 12rpx;
  box-sizing: border-box;
}

.empty-state__illustration {
  width: 128rpx;
  height: 128rpx;
  border-radius: 50%;
  background: var(--eapp-surface-subtle);
  border: 1rpx solid var(--eapp-border);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12rpx;
  box-sizing: border-box;
}

.empty-state__title {
  font-size: var(--eapp-fs-md);
  font-weight: 600;
  color: var(--eapp-text-secondary);
}

.empty-state__desc {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-faint);
  text-align: center;
  line-height: 1.6;
}

.empty-state__action {
  margin-top: 20rpx;
}
</style>
