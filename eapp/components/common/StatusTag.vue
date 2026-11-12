<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ text: string; type?: string }>()

interface TagTone {
  bg: string
  color: string
  border: string
}

const TONES: Record<string, TagTone> = {
  warning: { bg: 'var(--eapp-warning-soft)', color: 'var(--eapp-warning)', border: '#ecd6b8' },
  primary: { bg: 'var(--eapp-primary-soft)', color: 'var(--eapp-primary)', border: 'var(--eapp-primary-border)' },
  success: { bg: 'var(--eapp-success-soft)', color: 'var(--eapp-success)', border: '#c7e3d2' },
  danger: { bg: 'var(--eapp-danger-soft)', color: 'var(--eapp-danger)', border: '#eec9c4' },
  accent: { bg: 'var(--eapp-accent-soft)', color: 'var(--eapp-accent)', border: '#c4dee6' },
  neutral: { bg: 'var(--eapp-surface-subtle)', color: 'var(--eapp-text-muted)', border: 'var(--eapp-border-strong)' },
}

const styleVars = computed<TagTone>(() => {
  const t = String(props.type || '')
  if (t.includes('pending') || t.includes('applied') || t.includes('warning') || t.includes('warn'))
    return TONES.warning
  if (t.includes('ship') || t.includes('return'))
    return TONES.primary
  if (t.includes('complete') || t.includes('success') || t.includes('refunded') || t.includes('enabled'))
    return TONES.success
  if (t.includes('close') || t.includes('reject') || t === '6')
    return TONES.danger
  if (t.includes('disabled'))
    return TONES.neutral
  if (t.includes('hot') || t.includes('accent'))
    return TONES.accent
  return TONES.neutral
})
</script>

<template>
  <view
    class="status-tag"
    :style="{ backgroundColor: styleVars.bg, borderColor: styleVars.border, color: styleVars.color }"
  >
    <view class="status-tag__dot" :style="{ backgroundColor: styleVars.color }" />
    <text class="status-tag__text">{{ text }}</text>
  </view>
</template>

<style scoped>
.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 6rpx 12rpx;
  border-width: 1rpx;
  border-style: solid;
  border-radius: var(--eapp-radius-sm);
  flex-shrink: 0;
  box-sizing: border-box;
}

.status-tag__dot {
  width: 8rpx;
  height: 8rpx;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-tag__text {
  font-size: var(--eapp-fs-xs);
  font-weight: 600;
  line-height: 1.3;
  white-space: nowrap;
  color: inherit;
}
</style>
