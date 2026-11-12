<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'

export interface BatchFailItem {
  id: string | number
  reason: string
}

defineProps<{
  show: boolean
  success?: (string | number)[]
  fails?: BatchFailItem[]
}>()

const emit = defineEmits<{
  close: []
}>()
</script>

<template>
  <up-popup :show="show" mode="bottom" round="20" @close="emit('close')">
    <view class="batch-result">
      <view class="batch-result__header">
        <view class="batch-result__title-wrap">
          <view class="batch-result__bar" />
          <text class="batch-result__title">批量操作结果</text>
        </view>
        <view class="batch-result__close" hover-class="eapp-hover" @click="emit('close')">
          <AppIcon name="close" :size="30" color="muted" />
        </view>
      </view>

      <scroll-view class="batch-result__body" scroll-y>
        <view v-if="success?.length" class="batch-result__summary batch-result__summary--success">
          <AppIcon name="check-circle" :size="36" color="success" />
          <text class="batch-result__summary-text">成功处理 <text class="eapp-num">{{ success.length }}</text> 项</text>
        </view>

        <view v-if="fails?.length" class="batch-result__section">
          <view class="batch-result__summary batch-result__summary--fail">
            <AppIcon name="close-circle" :size="36" color="danger" />
            <text class="batch-result__summary-text">失败 <text class="eapp-num">{{ fails.length }}</text> 项</text>
          </view>
          <view class="batch-result__fails">
            <view
              v-for="item in fails"
              :key="item.id"
              class="batch-result__fail-item"
            >
              <text class="batch-result__fail-id eapp-num">#{{ item.id }}</text>
              <text class="batch-result__fail-reason">{{ item.reason }}</text>
            </view>
          </view>
        </view>
      </scroll-view>

      <view class="batch-result__footer" hover-class="eapp-hover" @click="emit('close')">
        <text class="batch-result__footer-text">完成</text>
      </view>
    </view>
  </up-popup>
</template>

<style scoped>
.batch-result {
  padding: var(--eapp-space-5);
  padding-bottom: calc(env(safe-area-inset-bottom) + 28rpx);
  max-height: 60vh;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}

.batch-result__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--eapp-space-4);
}

.batch-result__title-wrap {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.batch-result__bar {
  width: 6rpx;
  height: 28rpx;
  border-radius: 3rpx;
  background: var(--eapp-primary);
}

.batch-result__title {
  font-size: var(--eapp-fs-lg);
  font-weight: 600;
  color: var(--eapp-text);
}

.batch-result__close {
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius);
  border: 1rpx solid var(--eapp-border);
  background: var(--eapp-surface-subtle);
  box-sizing: border-box;
}

.batch-result__body {
  flex: 1;
  min-height: 80rpx;
}

.batch-result__section {
  margin-top: var(--eapp-space-3);
}

.batch-result__summary {
  display: flex;
  align-items: center;
  gap: 12rpx;
  padding: var(--eapp-space-3) var(--eapp-space-4);
  border-radius: var(--eapp-radius);
  margin-bottom: var(--eapp-space-3);
}

.batch-result__summary--success {
  background: var(--eapp-success-soft);
}

.batch-result__summary--fail {
  background: var(--eapp-danger-soft);
}

.batch-result__summary-text {
  font-size: var(--eapp-fs-base);
  font-weight: 600;
  color: var(--eapp-text);
}

.batch-result__fails {
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius);
  overflow: hidden;
}

.batch-result__fail-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16rpx;
  padding: var(--eapp-space-3);
  border-bottom: 1rpx solid var(--eapp-divider);
}

.batch-result__fail-item:last-child {
  border-bottom: none;
}

.batch-result__fail-id {
  font-size: var(--eapp-fs-sm);
  font-weight: 600;
  color: var(--eapp-text-secondary);
  flex-shrink: 0;
}

.batch-result__fail-reason {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-danger);
  text-align: right;
}

.batch-result__footer {
  height: 80rpx;
  margin-top: var(--eapp-space-4);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary);
  box-sizing: border-box;
}

.batch-result__footer-text {
  font-size: var(--eapp-fs-md);
  font-weight: 600;
  color: #ffffff;
}
</style>
