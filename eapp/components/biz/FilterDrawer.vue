<script setup lang="ts">
defineProps<{
  show: boolean
  title?: string
}>()

const emit = defineEmits<{
  close: []
  reset: []
  confirm: []
}>()
</script>

<template>
  <u-popup :show="show" mode="right" @close="emit('close')">
    <view class="filter-drawer">
      <view class="filter-drawer__header">
        <text class="filter-drawer__title">{{ title || '筛选' }}</text>
        <view class="filter-drawer__close" @click="emit('close')">
          <u-icon name="close" size="36" />
        </view>
      </view>
      <scroll-view class="filter-drawer__body" scroll-y>
        <slot />
      </scroll-view>
      <view class="filter-drawer__footer">
        <view class="filter-drawer__btn filter-drawer__btn--reset" @click="emit('reset')">
          <text class="filter-drawer__btn-text">重置</text>
        </view>
        <view class="filter-drawer__btn filter-drawer__btn--confirm" @click="emit('confirm')">
          <text class="filter-drawer__btn-text filter-drawer__btn-text--confirm">确定</text>
        </view>
      </view>
    </view>
  </u-popup>
</template>

<style scoped>
.filter-drawer {
  display: flex;
  flex-direction: column;
  width: 600rpx;
  height: 100vh;
  background: var(--eapp-bg);
}

.filter-drawer__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx;
  border-bottom: 1rpx solid var(--eapp-border);
  background: var(--eapp-card);
}

.filter-drawer__title {
  font-size: 30rpx;
  font-weight: 600;
  color: var(--eapp-text);
}

.filter-drawer__close {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.filter-drawer__body {
  flex: 1;
  padding: 28rpx;
}

.filter-drawer__footer {
  display: flex;
  gap: 20rpx;
  padding: 20rpx 28rpx;
  padding-bottom: calc(env(safe-area-inset-bottom) + 20rpx);
  background: var(--eapp-card);
  border-top: 1rpx solid var(--eapp-border);
}

.filter-drawer__btn {
  flex: 1;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 22rpx;
  min-height: 44rpx;
}

.filter-drawer__btn--reset {
  background: var(--eapp-bg);
  border: 1rpx solid var(--eapp-border);
}

.filter-drawer__btn--confirm {
  background: var(--eapp-primary);
}

.filter-drawer__btn-text {
  font-size: 28rpx;
  color: var(--eapp-text);
}

.filter-drawer__btn-text--confirm {
  color: #ffffff;
}
</style>
