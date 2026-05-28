<script setup lang="ts">
withDefaults(defineProps<{
  padded?: boolean
  headerSticky?: boolean
}>(), {
  padded: true,
  headerSticky: false,
})
</script>

<template>
  <view class="eapp-shell">
    <view
      class="eapp-shell__header"
      :class="{ 'eapp-shell__header--sticky': headerSticky }"
    >
      <slot name="header" />
    </view>
    <view class="eapp-shell__body" :class="{ 'eapp-shell__body--padded': padded }">
      <slot />
    </view>
  </view>
</template>

<style scoped>
.eapp-shell {
  min-height: 100vh;
  background: var(--eapp-bg);
}

.eapp-shell__header {
  padding-top: env(safe-area-inset-top);
}

.eapp-shell__header--sticky {
  position: sticky;
  top: 0;
  z-index: 50;
  background: var(--eapp-card);
}

.eapp-shell__body {
  box-sizing: border-box;
}

.eapp-shell__body--padded {
  padding: 24rpx;
}
</style>
