<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'

withDefaults(
  defineProps<{
    title: string
    subtitle?: string
    /** 右侧文字，例如「查看全部」 */
    extra?: string
    /** 是否显示右侧箭头 */
    arrow?: boolean
  }>(),
  { arrow: false },
)

const emit = defineEmits<{ click: [] }>()
</script>

<template>
  <view class="section-header" @click="emit('click')">
    <view class="section-header__main">
      <view class="section-header__bar" />
      <text class="section-header__title">{{ title }}</text>
      <text v-if="subtitle" class="section-header__subtitle">{{ subtitle }}</text>
    </view>
    <view class="section-header__side">
      <slot name="extra">
        <template v-if="extra || arrow">
          <text v-if="extra" class="section-header__extra">{{ extra }}</text>
          <AppIcon v-if="arrow" name="chevron-right" :size="24" color="faint" />
        </template>
      </slot>
    </view>
  </view>
</template>

<style scoped>
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16rpx;
  margin-bottom: var(--eapp-space-3);
  min-height: 40rpx;
}

.section-header__main {
  display: flex;
  align-items: center;
  gap: 12rpx;
  min-width: 0;
}

.section-header__bar {
  width: 6rpx;
  height: 28rpx;
  border-radius: 3rpx;
  background: var(--eapp-primary);
  flex-shrink: 0;
}

.section-header__title {
  font-size: var(--eapp-fs-md);
  font-weight: 600;
  color: var(--eapp-text);
  letter-spacing: 0.5rpx;
}

.section-header__subtitle {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.section-header__side {
  display: flex;
  align-items: center;
  gap: 4rpx;
  flex-shrink: 0;
}

.section-header__extra {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
}
</style>
