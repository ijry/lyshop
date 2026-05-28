<script setup lang="ts">
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
  <u-popup :show="show" mode="bottom" round="24" @close="emit('close')">
    <view class="batch-result">
      <view class="batch-result__header">
        <text class="batch-result__title">批量操作结果</text>
        <view class="batch-result__close" @click="emit('close')">
          <u-icon name="close" size="36" />
        </view>
      </view>
      <scroll-view class="batch-result__body" scroll-y>
        <view v-if="success?.length" class="batch-result__section">
          <text class="batch-result__label batch-result__label--success">
            成功 {{ success.length }} 项
          </text>
        </view>
        <view v-if="fails?.length" class="batch-result__section">
          <text class="batch-result__label batch-result__label--fail">
            失败 {{ fails.length }} 项
          </text>
          <view
            v-for="item in fails"
            :key="item.id"
            class="batch-result__fail-item"
          >
            <text class="batch-result__fail-id">{{ item.id }}</text>
            <text class="batch-result__fail-reason">{{ item.reason }}</text>
          </view>
        </view>
      </scroll-view>
    </view>
  </u-popup>
</template>

<style scoped>
.batch-result {
  padding: 28rpx;
  padding-bottom: calc(env(safe-area-inset-bottom) + 28rpx);
  max-height: 60vh;
  display: flex;
  flex-direction: column;
}

.batch-result__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24rpx;
}

.batch-result__title {
  font-size: 30rpx;
  font-weight: 600;
  color: var(--eapp-text);
}

.batch-result__close {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.batch-result__body {
  flex: 1;
}

.batch-result__section {
  margin-bottom: 20rpx;
}

.batch-result__label {
  font-size: 26rpx;
  font-weight: 600;
  margin-bottom: 12rpx;
  display: block;
}

.batch-result__label--success {
  color: var(--eapp-success);
}

.batch-result__label--fail {
  color: var(--eapp-danger);
}

.batch-result__fail-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 0;
  border-bottom: 1rpx solid var(--eapp-border);
}

.batch-result__fail-id {
  font-size: 24rpx;
  color: var(--eapp-text);
}

.batch-result__fail-reason {
  font-size: 22rpx;
  color: var(--eapp-danger);
}
</style>
