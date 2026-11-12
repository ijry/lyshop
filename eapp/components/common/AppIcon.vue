<script setup lang="ts">
import { computed } from 'vue'
import { iconDataUri } from '@/utils/icons'

/**
 * 跨端线性图标。
 * 用 base64 SVG data URI 渲染 <image>，避免内联 <svg> 在小程序/App 端不可用，
 * 也避免 emoji 在不同系统上字形不一致。
 */
const props = withDefaults(
  defineProps<{
    name: string
    /** 尺寸，单位 rpx */
    size?: number
    /** 语义色键（primary / muted / danger ...）、#hex 或 var(--eapp-*) */
    color?: string
    strokeWidth?: number
  }>(),
  {
    size: 40,
    color: 'secondary',
    strokeWidth: 1.7,
  },
)

const src = computed(() => iconDataUri(props.name, props.color, props.strokeWidth))
</script>

<template>
  <image
    class="app-icon"
    :src="src"
    :style="{ width: `${size}rpx`, height: `${size}rpx` }"
    mode="aspectFit"
  />
</template>

<style scoped>
.app-icon {
  display: block;
  flex-shrink: 0;
}
</style>
