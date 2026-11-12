<script setup lang="ts">
/** 列表加载骨架：统一 shimmer 效果，替代各页面手写骨架 */
withDefaults(defineProps<{ variant?: 'order' | 'default' }>(), { variant: 'order' })
</script>

<template>
  <view class="skeleton-card">
    <view class="sk-head">
      <view class="sk-block sk-head__no" />
      <view class="sk-block sk-head__tag" />
    </view>

    <template v-if="variant === 'order'">
      <view class="sk-products">
        <view class="sk-block sk-thumb" />
        <view class="sk-block sk-thumb" />
        <view class="sk-col">
          <view class="sk-block sk-col__line" />
          <view class="sk-block sk-col__line sk-col__line--short" />
          <view class="sk-block sk-col__line sk-col__line--price" />
        </view>
      </view>
      <view class="sk-buyer">
        <view class="sk-block sk-buyer__line" />
      </view>
    </template>
    <template v-else>
      <view class="sk-col sk-col--solo">
        <view class="sk-block sk-col__line" />
        <view class="sk-block sk-col__line sk-col__line--short" />
        <view class="sk-block sk-col__line sk-col__line--price" />
      </view>
    </template>

    <view class="sk-divider" />
    <view class="sk-actions">
      <view class="sk-block sk-btn" />
      <view class="sk-block sk-btn" />
      <view class="sk-block sk-btn" />
    </view>
  </view>
</template>

<style scoped>
.skeleton-card {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  padding: var(--eapp-space-4) var(--eapp-space-5) var(--eapp-space-3);
  box-sizing: border-box;
}

.sk-block {
  background: linear-gradient(90deg, #f0f2f6 25%, #e4e8ef 50%, #f0f2f6 75%);
  background-size: 200% 100%;
  border-radius: var(--eapp-radius-sm);
  animation: sk-shimmer 1.5s ease-in-out infinite;
}

@keyframes sk-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.sk-head,
.sk-products,
.sk-actions {
  display: flex;
  align-items: center;
}

.sk-head {
  justify-content: space-between;
}

.sk-head__no { width: 180rpx; height: 28rpx; }
.sk-head__tag { width: 96rpx; height: 32rpx; border-radius: var(--eapp-radius-sm); }

.sk-products {
  gap: 10rpx;
  margin-top: var(--eapp-space-3);
}

.sk-thumb { width: 84rpx; height: 84rpx; flex-shrink: 0; }

.sk-col {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  padding-left: 6rpx;
}

.sk-col--solo { margin-top: var(--eapp-space-3); padding-left: 0; }

.sk-col__line { width: 90%; height: 24rpx; }
.sk-col__line--short { width: 55%; height: 20rpx; }
.sk-col__line--price { width: 40%; height: 28rpx; margin-top: 4rpx; }

.sk-buyer { margin-top: var(--eapp-space-3); }
.sk-buyer__line { width: 46%; height: 22rpx; }

.sk-divider {
  height: 1rpx;
  background: var(--eapp-divider);
  margin: var(--eapp-space-3) 0 var(--eapp-space-2);
}

.sk-actions { gap: var(--eapp-space-2); }
.sk-btn { width: 108rpx; height: 48rpx; }
</style>
