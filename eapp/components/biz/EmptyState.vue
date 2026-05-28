<script setup lang="ts">
defineProps<{
  title?: string
  desc?: string
  icon?: string
}>()
</script>

<template>
  <view class="empty-state">
    <text v-if="icon" class="empty-state__icon">{{ icon }}</text>
    <text v-if="title" class="empty-state__title">{{ title }}</text>
    <text v-if="desc" class="empty-state__desc">{{ desc }}</text>
    <view v-if="$slots.default" class="empty-state__action">
      <slot />
    </view>
  </view>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80rpx 40rpx;
  gap: 16rpx;
}

.empty-state__icon {
  font-size: 72rpx;
  margin-bottom: 8rpx;
}

.empty-state__title {
  font-size: 30rpx;
  font-weight: 600;
  color: var(--eapp-text);
}

.empty-state__desc {
  font-size: 24rpx;
  color: var(--eapp-text-muted);
  text-align: center;
}

.empty-state__action {
  margin-top: 24rpx;
}
</style>
