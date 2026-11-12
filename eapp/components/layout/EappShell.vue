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
  box-sizing: border-box;
}

.eapp-shell__header {
  /* 自定义导航页由 PageHeader 自行处理状态栏安全区 */
  background: var(--eapp-card);
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
  padding: var(--eapp-space-5);
  padding-bottom: calc(env(safe-area-inset-bottom) + 48rpx);
}
</style>
