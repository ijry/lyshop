<template>
  <div class="space-y-4">
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">{{ $t('decor.spacer.height') }}</label>
      <div class="flex items-center gap-3">
        <input type="range" :value="modelValue.height" @input="update('height', Number(($event.target as HTMLInputElement).value))"
          min="4" max="200" step="4" class="flex-1 accent-blue-600" />
        <input type="number" :value="modelValue.height" @change="update('height', Number(($event.target as HTMLInputElement).value))"
          min="4" max="200" class="w-16 border border-slate-200 rounded-lg px-2 py-1.5 text-xs text-center focus:outline-none focus:border-blue-400" />
      </div>
    </div>

    <ColorInput :modelValue="modelValue.background" @update:modelValue="update('background', $event)" :label="$t('decor.spacer.bgColor')" />
  </div>
</template>

<script setup lang="ts">
import type { SpacerProps } from '@/types/decor'
import ColorInput from '../widgets/ColorInput.vue'

const props = defineProps<{ modelValue: SpacerProps }>()
const emit = defineEmits<{ 'update:modelValue': [value: SpacerProps] }>()

function update<K extends keyof SpacerProps>(key: K, value: SpacerProps[K]) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}
</script>
