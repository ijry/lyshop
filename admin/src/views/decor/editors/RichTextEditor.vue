<template>
  <div class="space-y-2">
    <label class="block text-xs text-slate-500">HTML 内容</label>
    <textarea :value="modelValue.content"
      @input="emit('update:modelValue', { content: ($event.target as HTMLTextAreaElement).value })"
      rows="8" placeholder="输入 HTML 内容..."
      class="w-full border border-slate-200 rounded-lg px-3 py-2 text-xs font-mono resize-none focus:outline-none focus:border-blue-400" />
  </div>
</template>

<script setup lang="ts">
import type { RichTextProps } from '@/types/decor'

defineProps<{ modelValue: RichTextProps }>()
const emit = defineEmits<{ 'update:modelValue': [value: RichTextProps] }>()
</script>
