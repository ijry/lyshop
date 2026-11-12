<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'

defineProps<{ title: string; value: number | string; desc?: string }>()
defineEmits<{ (e: 'click'): void }>()
</script>

<template>
  <view class="todo-card" hover-class="eapp-hover" @click="$emit('click')">
    <view class="todo-card__main">
      <view class="todo-card__title-row">
        <AppIcon name="clipboard" :size="28" color="primary" />
        <text class="todo-card__title">{{ title }}</text>
      </view>
      <text v-if="desc" class="todo-card__desc">{{ desc }}</text>
    </view>
    <view class="todo-card__side">
      <text class="todo-card__value eapp-num">{{ value }}</text>
      <AppIcon name="chevron-right" :size="26" color="faint" />
    </view>
  </view>
</template>

<style scoped>
.todo-card {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-4) var(--eapp-space-5);
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16rpx;
}

.todo-card__main {
  min-width: 0;
}

.todo-card__title-row {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.todo-card__title {
  font-size: var(--eapp-fs-md);
  font-weight: 600;
  color: var(--eapp-text);
}

.todo-card__desc {
  display: block;
  margin-top: 6rpx;
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
}

.todo-card__side {
  display: flex;
  align-items: center;
  gap: 8rpx;
  flex-shrink: 0;
}

.todo-card__value {
  font-size: var(--eapp-fs-lg);
  font-weight: 700;
  color: var(--eapp-primary);
}
</style>
