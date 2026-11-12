<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'

withDefaults(
  defineProps<{
    title: string
    subtitle?: string
    eyebrow?: string
    transparent?: boolean
    back?: boolean
    border?: boolean
  }>(),
  {
    transparent: false,
    back: false,
    border: true,
  },
)

const emit = defineEmits<{ back: [] }>()

function onBack() {
  emit('back')
  uni.navigateBack({
    fail: () => uni.switchTab({ url: '/pages/dashboard/index' }),
  })
}
</script>

<template>
  <view
    class="page-header"
    :class="{
      'page-header--transparent': transparent,
      'page-header--bordered': border && !transparent,
    }"
  >
    <view v-if="back" class="page-header__back" hover-class="eapp-hover" @click="onBack">
      <AppIcon name="chevron-left" :size="36" color="text" />
    </view>
    <view class="page-header__left">
      <text v-if="eyebrow" class="page-header__eyebrow">{{ eyebrow }}</text>
      <text class="page-header__title">{{ title }}</text>
      <text v-if="subtitle" class="page-header__subtitle">{{ subtitle }}</text>
    </view>
    <view class="page-header__right">
      <slot name="right" />
    </view>
  </view>
</template>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  gap: 16rpx;
  min-height: 96rpx;
  padding: calc(env(safe-area-inset-top) + 20rpx) 28rpx 20rpx;
  background: var(--eapp-card);
  box-sizing: border-box;
}

.page-header--bordered {
  border-bottom: 1rpx solid var(--eapp-border);
}

.page-header--transparent {
  background: transparent;
}

.page-header__back {
  width: 56rpx;
  height: 56rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius);
  background: var(--eapp-surface-subtle);
  border: 1rpx solid var(--eapp-border);
  box-sizing: border-box;
}

.page-header__left {
  display: flex;
  flex-direction: column;
  gap: 2rpx;
  flex: 1;
  min-width: 0;
}

.page-header__eyebrow {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
  letter-spacing: 2rpx;
}

.page-header__title {
  font-size: var(--eapp-fs-xl);
  font-weight: 700;
  letter-spacing: 0.5rpx;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.page-header__subtitle {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.page-header__right {
  display: flex;
  align-items: center;
  gap: 16rpx;
  flex-shrink: 0;
}
</style>
