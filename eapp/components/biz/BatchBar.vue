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
  <view v-if="count > 0" class="batch-bar">
    <view class="batch-bar__info">
      <text class="batch-bar__count eapp-num">已选 {{ count }} 项</text>
      <view class="batch-bar__cancel" hover-class="eapp-hover" @click="emit('cancel')">
        <text class="batch-bar__cancel-text">取消</text>
      </view>
    </view>
    <view class="batch-bar__actions">
      <view
        v-for="act in actions"
        :key="act.key"
        class="batch-bar__btn"
        :style="{ backgroundColor: toneBg(act.tone) }"
        hover-class="eapp-hover"
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
  gap: var(--eapp-space-3);
  padding: var(--eapp-space-3) var(--eapp-space-5);
  padding-bottom: calc(env(safe-area-inset-bottom) + 16rpx);
  background: var(--eapp-card);
  border-top: 1rpx solid var(--eapp-border);
  box-shadow: 0 -4rpx 20rpx rgba(16, 24, 40, 0.06);
  z-index: 100;
  box-sizing: border-box;
}

.batch-bar__info {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
  flex-shrink: 0;
}

.batch-bar__count {
  font-size: var(--eapp-fs-base);
  color: var(--eapp-text);
  font-weight: 600;
}

.batch-bar__cancel {
  min-height: 44rpx;
  display: flex;
  align-items: center;
}

.batch-bar__cancel-text {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
}

.batch-bar__actions {
  display: flex;
  gap: var(--eapp-space-2);
  flex-wrap: wrap;
  justify-content: flex-end;
}

.batch-bar__btn {
  padding: 12rpx 24rpx;
  border-radius: var(--eapp-radius);
  min-height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
}

.batch-bar__btn-text {
  font-size: var(--eapp-fs-sm);
  font-weight: 500;
  color: #ffffff;
  white-space: nowrap;
}
</style>
