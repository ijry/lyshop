<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'

/**
 * 营销域统一底部表单弹层：标题 / 关闭 / 内容滚动区 / 取消与提交。
 * 小程序端 view 的 overflow 不生效，因此内容区用 scroll-view 并显式给高度。
 */
withDefaults(
  defineProps<{
    show: boolean
    title: string
    subtitle?: string
    submitText?: string
    submitting?: boolean
    bodyHeight?: string
  }>(),
  {
    submitText: '保存',
    submitting: false,
    bodyHeight: '52vh',
  },
)

const emit = defineEmits<{ close: []; submit: [] }>()
</script>

<template>
  <up-popup :show="show" mode="bottom" :round="20" :close-on-click-overlay="true" @close="emit('close')">
    <view class="sheet">
      <view class="sheet__grabber" />
      <view class="sheet__head">
        <view class="sheet__heading">
          <text class="sheet__title">{{ title }}</text>
          <text v-if="subtitle" class="sheet__subtitle">{{ subtitle }}</text>
        </view>
        <view class="sheet__close" hover-class="eapp-hover" @click="emit('close')">
          <AppIcon name="close" :size="30" color="muted" />
        </view>
      </view>

      <scroll-view scroll-y class="sheet__body" :style="{ height: bodyHeight }">
        <slot />
      </scroll-view>

      <view class="sheet__footer">
        <view class="sheet__btn sheet__btn--ghost" hover-class="eapp-hover" @click="emit('close')">取消</view>
        <view
          class="sheet__btn sheet__btn--primary"
          :class="{ 'sheet__btn--busy': submitting }"
          hover-class="eapp-hover"
          @click="emit('submit')"
        >
          <text>{{ submitting ? '处理中…' : submitText }}</text>
        </view>
      </view>
    </view>
  </up-popup>
</template>

<style scoped>
.sheet {
  background: var(--eapp-card);
  border-radius: var(--eapp-radius-xl) var(--eapp-radius-xl) 0 0;
  padding-top: var(--eapp-space-3);
  box-sizing: border-box;
}

.sheet__grabber {
  width: 72rpx;
  height: 8rpx;
  border-radius: var(--eapp-radius-pill);
  background: var(--eapp-border-strong);
  margin: 0 auto var(--eapp-space-3);
}

.sheet__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--eapp-space-3);
  padding: 0 var(--eapp-space-5) var(--eapp-space-4);
  border-bottom: 1rpx solid var(--eapp-divider);
  box-sizing: border-box;
}

.sheet__heading {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  flex: 1;
  min-width: 0;
}

.sheet__title {
  font-size: var(--eapp-fs-lg);
  font-weight: 700;
  color: var(--eapp-text);
}

.sheet__subtitle {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
}

.sheet__close {
  width: 56rpx;
  height: 56rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--eapp-surface-subtle);
  border: 1rpx solid var(--eapp-border);
  box-sizing: border-box;
}

.sheet__body {
  width: 100%;
  padding: var(--eapp-space-5);
  box-sizing: border-box;
}

.sheet__footer {
  display: flex;
  gap: var(--eapp-space-3);
  padding: var(--eapp-space-4) var(--eapp-space-5)
    calc(var(--eapp-space-4) + env(safe-area-inset-bottom));
  border-top: 1rpx solid var(--eapp-divider);
  box-sizing: border-box;
}

.sheet__btn {
  height: 84rpx;
  border-radius: var(--eapp-radius);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--eapp-fs-md);
  font-weight: 600;
  box-sizing: border-box;
}

.sheet__btn--ghost {
  width: 200rpx;
  flex-shrink: 0;
  color: var(--eapp-text-secondary);
  background: var(--eapp-surface-subtle);
  border: 1rpx solid var(--eapp-border);
}

.sheet__btn--primary {
  flex: 1;
  color: var(--eapp-primary-on);
  background: var(--eapp-primary);
}

.sheet__btn--busy {
  opacity: 0.6;
}
</style>
