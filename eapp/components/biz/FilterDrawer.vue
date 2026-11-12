<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'

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
  <up-popup :show="show" mode="right" @close="emit('close')">
    <view class="filter-drawer">
      <view class="filter-drawer__header">
        <text class="filter-drawer__title">{{ title || '筛选' }}</text>
        <view class="filter-drawer__close" hover-class="eapp-hover" @click="emit('close')">
          <AppIcon name="close" :size="32" color="muted" />
        </view>
      </view>
      <scroll-view class="filter-drawer__body" scroll-y>
        <slot />
      </scroll-view>
      <view class="filter-drawer__footer">
        <view class="filter-drawer__btn filter-drawer__btn--reset" hover-class="eapp-hover" @click="emit('reset')">
          <text class="filter-drawer__btn-text">重置</text>
        </view>
        <view class="filter-drawer__btn filter-drawer__btn--confirm" hover-class="eapp-hover" @click="emit('confirm')">
          <text class="filter-drawer__btn-text filter-drawer__btn-text--confirm">确定</text>
        </view>
      </view>
    </view>
  </up-popup>
</template>

<style scoped>
.filter-drawer {
  display: flex;
  flex-direction: column;
  width: 600rpx;
  height: 100vh;
  background: var(--eapp-bg);
  box-sizing: border-box;
}

.filter-drawer__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--eapp-space-5);
  padding-top: calc(env(safe-area-inset-top) + 28rpx);
  border-bottom: 1rpx solid var(--eapp-border);
  background: var(--eapp-card);
}

.filter-drawer__title {
  font-size: var(--eapp-fs-lg);
  font-weight: 600;
  color: var(--eapp-text);
}

.filter-drawer__close {
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

.filter-drawer__body {
  flex: 1;
  padding: var(--eapp-space-5);
  box-sizing: border-box;
}

.filter-drawer__footer {
  display: flex;
  gap: var(--eapp-space-3);
  padding: var(--eapp-space-4) var(--eapp-space-5);
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
  border-radius: var(--eapp-radius);
  box-sizing: border-box;
}

.filter-drawer__btn--reset {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border-strong);
}

.filter-drawer__btn--confirm {
  background: var(--eapp-primary);
  border: 1rpx solid var(--eapp-primary);
}

.filter-drawer__btn-text {
  font-size: var(--eapp-fs-md);
  font-weight: 500;
  color: var(--eapp-text-secondary);
}

.filter-drawer__btn-text--confirm {
  color: #ffffff;
}
</style>
