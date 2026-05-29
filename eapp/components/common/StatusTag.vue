<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ text: string; type?: string }>()

const styleVars = computed(() => {
  const t = String(props.type || '')
  // warning tones
  if (t.includes('pending') || t.includes('applied') || t.includes('warning') || t.includes('warn'))
    return { bg: '#fef3c7', color: '#92400e' }
  // primary tones
  if (t.includes('ship') || t.includes('return'))
    return { bg: '#dbeafe', color: '#1d4ed8' }
  // success tones
  if (t.includes('complete') || t.includes('success') || t.includes('refunded') || t.includes('enabled'))
    return { bg: '#dcfce7', color: '#166534' }
  // danger tones
  if (t.includes('close') || t.includes('reject') || t === '6')
    return { bg: '#fee2e2', color: '#991b1b' }
  // neutral tone
  if (t.includes('disabled'))
    return { bg: '#e2e8f0', color: '#64748b' }
  // accent tones
  if (t.includes('hot') || t.includes('accent'))
    return { bg: '#ffedd5', color: '#c2410c' }
  return { bg: '#e2e8f0', color: '#334155' }
})
</script>

<template>
  <text class="status-tag" :style="{ backgroundColor: styleVars.bg, color: styleVars.color }">{{ text }}</text>
</template>

<style scoped>
.status-tag {
  font-size: 20rpx;
  border-radius: 999rpx;
  padding: 8rpx 14rpx;
  display: inline-flex;
  align-items: center;
}
</style>
