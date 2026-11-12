<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'

export interface ActionItem {
  key: string
  label: string
  /** 统一线性图标名，见 utils/icons.ts */
  iconName?: string
  /** 兼容旧数据：可传入图标名，也可传入历史 emoji（会被自动映射） */
  icon?: string
  color?: string
  badge?: number | string
  soon?: boolean
}

defineProps<{
  items: ActionItem[]
  columns?: number
}>()

const emit = defineEmits<{
  click: [key: string]
}>()

/** 历史页面仍传 emoji，这里统一收敛为线性图标，避免逐个改页面 */
const EMOJI_MAP: Record<string, string> = {
  '📦': 'order',
  '🏷': 'product',
  '🔧': 'after-sale',
  '⭐': 'star',
  '🎨': 'decor',
  '📊': 'chart',
  '🏭': 'warehouse',
  '🎫': 'coupon',
  '⚡': 'flash',
  '👥': 'users',
  '💰': 'money',
  '👑': 'crown',
  '📢': 'megaphone',
  '📋': 'clipboard',
  '🔗': 'share',
  '💎': 'diamond',
  '📷': 'scan',
  '💬': 'chat',
  '🏠': 'warehouse',
  '📄': 'document',
  '🔄': 'refresh',
  '🏪': 'store',
  '📈': 'trend-up',
  '📉': 'trend-down',
  '⚙': 'settings',
  '🛒': 'package',
  '🔍': 'search',
}

function iconOf(item: ActionItem): string {
  const raw = item.iconName || item.icon || 'circle'
  return EMOJI_MAP[raw] || raw
}
</script>

<template>
  <view
    class="action-grid"
    :style="{ gridTemplateColumns: `repeat(${columns || 4}, minmax(0, 1fr))` }"
  >
    <view
      v-for="item in items"
      :key="item.key"
      class="action-grid__item"
      :class="{ 'action-grid__item--soon': item.soon }"
      hover-class="eapp-hover"
      :hover-stay-time="80"
      @click="!item.soon && emit('click', item.key)"
    >
      <view class="action-grid__icon-wrap">
        <AppIcon :name="iconOf(item)" :size="40" color="primary" />
        <view v-if="item.badge" class="action-grid__badge">
          <text class="action-grid__badge-text">{{ item.badge }}</text>
        </view>
      </view>
      <text class="action-grid__label">{{ item.label }}</text>
      <view v-if="item.soon" class="action-grid__soon">
        <text class="action-grid__soon-text">即将上线</text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.action-grid {
  display: grid;
  gap: var(--eapp-space-4) var(--eapp-space-2);
}

.action-grid__item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
  position: relative;
  min-width: 0;
  padding: 16rpx 4rpx;
  border-radius: var(--eapp-radius);
  box-sizing: border-box;
}

.action-grid__item--soon {
  opacity: 0.48;
}

.action-grid__icon-wrap {
  position: relative;
  width: 76rpx;
  height: 76rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--eapp-primary-soft);
  border: 1rpx solid var(--eapp-primary-border);
  border-radius: var(--eapp-radius);
  box-sizing: border-box;
}

.action-grid__badge {
  position: absolute;
  top: -10rpx;
  right: -12rpx;
  min-width: 30rpx;
  height: 30rpx;
  padding: 0 8rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--eapp-danger);
  border: 2rpx solid var(--eapp-card);
  border-radius: var(--eapp-radius-pill);
  box-sizing: border-box;
}

.action-grid__badge-text {
  font-size: 18rpx;
  color: #ffffff;
  font-weight: 600;
  line-height: 1;
}

.action-grid__label {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-secondary);
  text-align: center;
  line-height: 1.35;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.action-grid__soon {
  position: absolute;
  top: 4rpx;
  right: 0;
}

.action-grid__soon-text {
  font-size: 16rpx;
  color: var(--eapp-text-muted);
  background: var(--eapp-border);
  border-radius: var(--eapp-radius-pill);
  padding: 2rpx 8rpx;
  white-space: nowrap;
}
</style>
