<script setup lang="ts">
defineProps<{ title: string; value: number | string; desc?: string }>()
defineEmits<{ (e: 'click'): void }>()
</script>

<template>
  <view class="todo-card" @click="$emit('click')">
    <view>
      <view class="todo-title">{{ title }}</view>
      <view v-if="desc" class="todo-desc">{{ desc }}</view>
    </view>
    <view class="todo-value">{{ value }}</view>
  </view>
</template>

<style scoped>
.todo-card {
  background: var(--eapp-card);
  border: 1px solid var(--eapp-border);
  border-radius: 24rpx;
  padding: 24rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.todo-title {
  font-size: 28rpx;
  font-weight: 600;
}

.todo-desc {
  margin-top: 6rpx;
  font-size: 22rpx;
  color: var(--eapp-text-muted);
}

.todo-value {
  font-size: 34rpx;
  color: var(--eapp-primary);
  font-weight: 700;
}
</style>
