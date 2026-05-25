<template>
  <div class="flex items-center gap-2">
    <label v-if="label" class="text-xs text-slate-500 shrink-0">{{ label }}</label>
    <input type="color" :value="modelValue" @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      class="w-8 h-8 rounded-lg border border-slate-200 cursor-pointer p-0.5" />
    <input type="text" :value="modelValue" @change="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      class="w-24 border border-slate-200 rounded-lg px-2 py-1.5 text-xs font-mono focus:outline-none focus:border-blue-400" />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string
  label?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>
