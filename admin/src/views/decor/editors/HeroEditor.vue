<template>
  <div class="space-y-4">
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">角标文字</label>
      <input :value="modelValue.badge" @input="update('badge', ($event.target as HTMLInputElement).value)"
        placeholder="如：限时秒杀进行中" class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
    </div>
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">主标题（换行用 \n）</label>
      <input :value="modelValue.title" @input="update('title', ($event.target as HTMLInputElement).value)"
        placeholder="精选好物\n品质生活" class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
    </div>
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">副标题</label>
      <input :value="modelValue.subtitle" @input="update('subtitle', ($event.target as HTMLInputElement).value)"
        class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
    </div>
    <div class="grid grid-cols-2 gap-3">
      <div>
        <label class="block text-xs text-slate-500 mb-1.5">按钮文字</label>
        <input :value="modelValue.btn_text" @input="update('btn_text', ($event.target as HTMLInputElement).value)"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
      </div>
      <div>
        <label class="block text-xs text-slate-500 mb-1.5">按钮链接</label>
        <input :value="modelValue.btn_link" @input="update('btn_link', ($event.target as HTMLInputElement).value)"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
      </div>
    </div>
    <div class="grid grid-cols-2 gap-3">
      <div>
        <label class="block text-xs text-slate-500 mb-1.5">次按钮文字</label>
        <input :value="modelValue.btn2_text" @input="update('btn2_text', ($event.target as HTMLInputElement).value)"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
      </div>
      <div>
        <label class="block text-xs text-slate-500 mb-1.5">次按钮链接</label>
        <input :value="modelValue.btn2_link" @input="update('btn2_link', ($event.target as HTMLInputElement).value)"
          class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
      </div>
    </div>
    <div class="grid grid-cols-2 gap-3">
      <ColorInput :modelValue="modelValue.bg_from || '#b91c1c'" @update:modelValue="update('bg_from', $event)" label="渐变起始" />
      <ColorInput :modelValue="modelValue.bg_to || '#991b1b'" @update:modelValue="update('bg_to', $event)" label="渐变结束" />
    </div>
  </div>
</template>

<script setup lang="ts">
import ColorInput from '../widgets/ColorInput.vue'

const props = defineProps<{ modelValue: any }>()
const emit = defineEmits<{ 'update:modelValue': [value: any] }>()

function update(key: string, value: any) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}
</script>
