<script setup lang="ts">
export interface ActionItem {
  key: string
  label: string
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
</script>

<template>
  <view
    class="action-grid"
    :style="{ gridTemplateColumns: `repeat(${columns || 4}, 1fr)` }"
  >
    <view
      v-for="item in items"
      :key="item.key"
      class="action-grid__item"
      :class="{ 'action-grid__item--soon': item.soon }"
      @click="!item.soon && emit('click', item.key)"
    >
      <view class="action-grid__icon-wrap">
        <text class="action-grid__icon" :style="item.color ? { color: item.color } : {}">
          {{ item.icon || '' }}
        </text>
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
  gap: 20rpx;
  padding: 20rpx 0;
}

.action-grid__item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10rpx;
  position: relative;
  min-height: 44rpx;
}

.action-grid__item--soon {
  opacity: 0.5;
}

.action-grid__icon-wrap {
  position: relative;
  width: 80rpx;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--eapp-bg);
  border-radius: 22rpx;
}

.action-grid__icon {
  font-size: 40rpx;
}

.action-grid__badge {
  position: absolute;
  top: -8rpx;
  right: -8rpx;
  background: var(--eapp-danger);
  border-radius: 999rpx;
  min-width: 32rpx;
  height: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 8rpx;
}

.action-grid__badge-text {
  font-size: 18rpx;
  color: #fff;
  font-weight: 600;
}

.action-grid__label {
  font-size: 22rpx;
  color: var(--eapp-text);
  text-align: center;
}

.action-grid__soon {
  position: absolute;
  top: 0;
  right: 0;
}

.action-grid__soon-text {
  font-size: 16rpx;
  color: var(--eapp-text-muted);
  background: var(--eapp-border);
  border-radius: 999rpx;
  padding: 2rpx 8rpx;
}
</style>
