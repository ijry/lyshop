<template>
  <div class="space-y-3 border-t border-slate-100 pt-3">
    <p class="text-xs text-slate-500">组件样式覆盖</p>
    <div class="grid grid-cols-2 gap-3">
      <div>
        <label class="block text-xs text-slate-500 mb-1.5">上边距</label>
        <input
          type="number"
          :value="styleValue.marginTop ?? ''"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
          @input="update('marginTop', Number(($event.target as HTMLInputElement).value || 0))"
        />
      </div>
      <div>
        <label class="block text-xs text-slate-500 mb-1.5">下边距</label>
        <input
          type="number"
          :value="styleValue.marginBottom ?? ''"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
          @input="update('marginBottom', Number(($event.target as HTMLInputElement).value || 0))"
        />
      </div>
      <div>
        <label class="block text-xs text-slate-500 mb-1.5">左右内边距</label>
        <input
          type="number"
          :value="styleValue.paddingX ?? ''"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
          @input="update('paddingX', Number(($event.target as HTMLInputElement).value || 0))"
        />
      </div>
      <div>
        <label class="block text-xs text-slate-500 mb-1.5">上下内边距</label>
        <input
          type="number"
          :value="styleValue.paddingY ?? ''"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
          @input="update('paddingY', Number(($event.target as HTMLInputElement).value || 0))"
        />
      </div>
      <div>
        <label class="block text-xs text-slate-500 mb-1.5">圆角</label>
        <input
          type="number"
          :value="styleValue.borderRadius ?? ''"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
          @input="update('borderRadius', Number(($event.target as HTMLInputElement).value || 0))"
        />
      </div>
      <div>
        <label class="block text-xs text-slate-500 mb-1.5">边框宽度</label>
        <input
          type="number"
          :value="styleValue.borderWidth ?? ''"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
          @input="update('borderWidth', Number(($event.target as HTMLInputElement).value || 0))"
        />
      </div>
    </div>

    <ColorInput
      :modelValue="styleValue.backgroundColor || '#ffffff'"
      label="背景色"
      @update:modelValue="update('backgroundColor', $event)"
    />
    <ColorInput
      :modelValue="styleValue.borderColor || '#e5e7eb'"
      label="边框色"
      @update:modelValue="update('borderColor', $event)"
    />

    <div>
      <label class="block text-xs text-slate-500 mb-1.5">阴影</label>
      <select
        :value="styleValue.shadow || 'none'"
        class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400"
        @change="update('shadow', ($event.target as HTMLSelectElement).value)"
      >
        <option value="none">none</option>
        <option value="sm">sm</option>
        <option value="md">md</option>
        <option value="lg">lg</option>
      </select>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import ColorInput from './ColorInput.vue'

const props = defineProps<{ modelValue?: any }>()
const emit = defineEmits<{ 'update:modelValue': [value: any] }>()

const styleValue = computed(() => props.modelValue || {})

function update(key: string, value: any) {
  emit('update:modelValue', {
    ...styleValue.value,
    [key]: value,
  })
}
</script>
