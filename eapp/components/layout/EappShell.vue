<template>
  <view class="eapp-shell">
    <slot name="header" />
    <view class="eapp-shell__body">
      <slot />
    </view>
  </view>
</template>

<style scoped>
.eapp-shell {
  min-height: 100vh;
  background: var(--eapp-bg);
}

.eapp-shell__body {
  padding: 24rpx;
  box-sizing: border-box;
}
</style>
