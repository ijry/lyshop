<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ text: string; type?: string }>()

const styleVars = computed(() => {
  const t = String(props.type || '')
  if (t.includes('pending')) return { bg: '#fef3c7', color: '#92400e' }
  if (t.includes('ship')) return { bg: '#dbeafe', color: '#1d4ed8' }
  if (t.includes('complete')) return { bg: '#dcfce7', color: '#166534' }
  if (t.includes('close') || t.includes('reject')) return { bg: '#fee2e2', color: '#991b1b' }
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
