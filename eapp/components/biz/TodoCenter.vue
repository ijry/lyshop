<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'

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

const ICON_BY_KEY: Record<string, string> = {
  pending_ship: 'truck',
  pending_after_sale: 'after-sale',
  unread_message: 'chat',
  stock_warning: 'warning',
  today_orders: 'order',
  pending_invoice: 'document',
  pending_refund: 'money',
}

type Tone = NonNullable<TodoItem['tone']>

const TONE: Record<Tone, { color: string; soft: string; icon: string }> = {
  warning: { color: 'var(--eapp-warning)', soft: 'var(--eapp-warning-soft)', icon: 'warning' },
  danger: { color: 'var(--eapp-danger)', soft: 'var(--eapp-danger-soft)', icon: 'danger' },
  primary: { color: 'var(--eapp-primary)', soft: 'var(--eapp-primary-soft)', icon: 'primary' },
  success: { color: 'var(--eapp-success)', soft: 'var(--eapp-success-soft)', icon: 'success' },
}

function toneOf(tone?: string) {
  return TONE[(tone || 'primary') as Tone] || TONE.primary
}

function iconOf(item: TodoItem) {
  return ICON_BY_KEY[item.key] || 'clipboard'
}
</script>

<template>
  <view class="todo-center">
    <view
      v-for="item in items"
      :key="item.key"
      class="todo-center__item"
      hover-class="eapp-hover"
      :hover-stay-time="80"
      @click="emit('click', item.key)"
    >
      <view class="todo-center__top">
        <view class="todo-center__icon" :style="{ background: toneOf(item.tone).soft }">
          <AppIcon :name="iconOf(item)" :size="32" :color="toneOf(item.tone).icon" />
        </view>
        <AppIcon name="chevron-right" :size="24" color="faint" />
      </view>
      <text class="todo-center__value eapp-num" :style="{ color: toneOf(item.tone).color }">
        {{ item.value }}
      </text>
      <text class="todo-center__title">{{ item.title }}</text>
    </view>
  </view>
</template>

<style scoped>
.todo-center {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--eapp-space-3);
}

.todo-center__item {
  display: flex;
  flex-direction: column;
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-4);
  box-sizing: border-box;
  min-width: 0;
}

.todo-center__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12rpx;
}

.todo-center__icon {
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius);
  box-sizing: border-box;
}

.todo-center__value {
  font-size: var(--eapp-fs-2xl);
  font-weight: 700;
  line-height: 1.15;
}

.todo-center__title {
  margin-top: 4rpx;
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
