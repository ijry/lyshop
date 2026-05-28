<script setup lang="ts">
export interface BatchAction {
  key: string
  label: string
  tone?: 'primary' | 'danger' | 'warning'
}

defineProps<{
  count: number
  actions: BatchAction[]
}>()

const emit = defineEmits<{
  action: [key: string]
  cancel: []
}>()

function toneBg(tone?: string) {
  const map: Record<string, string> = {
    primary: 'var(--eapp-primary)',
    danger: 'var(--eapp-danger)',
    warning: 'var(--eapp-warning)',
  }
  return map[tone || ''] || 'var(--eapp-primary)'
}
</script>

<template>
  <view class="batch-bar">
    <view class="batch-bar__info">
      <text class="batch-bar__count">已选 {{ count }} 项</text>
      <view class="batch-bar__cancel" @click="emit('cancel')">
        <text class="batch-bar__cancel-text">取消</text>
      </view>
    </view>
    <view class="batch-bar__actions">
      <view
        v-for="act in actions"
        :key="act.key"
        class="batch-bar__btn"
        :style="{ backgroundColor: toneBg(act.tone) }"
        @click="emit('action', act.key)"
      >
        <text class="batch-bar__btn-text">{{ act.label }}</text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.batch-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 28rpx;
  padding-bottom: calc(env(safe-area-inset-bottom) + 16rpx);
  background: var(--eapp-card);
  border-top: 1rpx solid var(--eapp-border);
  z-index: 100;
}

.batch-bar__info {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.batch-bar__count {
  font-size: 26rpx;
  color: var(--eapp-text);
  font-weight: 600;
}

.batch-bar__cancel {
  min-height: 44rpx;
  display: flex;
  align-items: center;
}

.batch-bar__cancel-text {
  font-size: 26rpx;
  color: var(--eapp-text-muted);
}

.batch-bar__actions {
  display: flex;
  gap: 16rpx;
}

.batch-bar__btn {
  padding: 14rpx 28rpx;
  border-radius: 22rpx;
  min-height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.batch-bar__btn-text {
  font-size: 26rpx;
  color: #ffffff;
}
</style>
