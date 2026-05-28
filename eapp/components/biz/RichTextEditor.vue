<script setup lang="ts">
defineProps<{ html?: string }>()
defineEmits<{ (e: 'requestEdit'): void }>()
</script>

<template>
  <view class="rt">
    <view class="preview">
      <rich-text v-if="html" :nodes="html" />
      <view v-else class="empty">尚未编辑详情内容</view>
    </view>
    <up-button size="mini" plain @click="$emit('requestEdit')">编辑详情</up-button>
  </view>
</template>

<style scoped>
.rt { display: grid; gap: 12rpx; }
.preview { background: var(--eapp-bg); border: 1px solid var(--eapp-border); border-radius: 14rpx; padding: 14rpx; min-height: 160rpx; }
.empty { color: var(--eapp-text-muted); font-size: 24rpx; text-align: center; padding: 30rpx 0; }
</style>
