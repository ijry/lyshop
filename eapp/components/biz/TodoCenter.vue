<script setup lang="ts">
export interface TodoItem {
  key: string
  title: string
  value: number | string
  tone?: 'warning' | 'danger' | 'primary' | 'success'
}

defineProps<{
  items: TodoItem[]
}>()

const emit = defineEmits<{
  click: [key: string]
}>()

function toneColor(tone?: string) {
  const map: Record<string, string> = {
    warning: 'var(--eapp-warning)',
    danger: 'var(--eapp-danger)',
    primary: 'var(--eapp-primary)',
    success: 'var(--eapp-success)',
  }
  return map[tone || ''] || 'var(--eapp-primary)'
}
</script>

<template>
  <view class="todo-center">
    <view
      v-for="item in items"
      :key="item.key"
      class="todo-center__item"
      @click="emit('click', item.key)"
    >
      <text class="todo-center__value" :style="{ color: toneColor(item.tone) }">
        {{ item.value }}
      </text>
      <text class="todo-center__title">{{ item.title }}</text>
    </view>
  </view>
</template>

<style scoped>
.todo-center {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
}

.todo-center__item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: var(--eapp-card);
  border-radius: 22rpx;
  padding: 28rpx 16rpx;
  gap: 8rpx;
  min-height: 44rpx;
}

.todo-center__value {
  font-size: 40rpx;
  font-weight: 700;
}

.todo-center__title {
  font-size: 24rpx;
  color: var(--eapp-text-muted);
}
</style>
