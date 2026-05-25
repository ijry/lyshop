<template>
  <div class="space-y-4">
    <!-- Notice items -->
    <div>
      <div class="flex items-center justify-between mb-2">
        <label class="text-xs text-slate-500">公告列表</label>
        <button @click="addItem" class="text-xs text-blue-600 hover:underline">+ 添加</button>
      </div>
      <div class="space-y-2">
        <div v-for="(item, idx) in modelValue.items" :key="idx"
          class="border border-slate-200 rounded-lg p-2.5 space-y-2">
          <div class="flex items-center gap-2">
            <input :value="item.text" @input="updateItem(idx, 'text', ($event.target as HTMLInputElement).value)"
              placeholder="公告内容" class="flex-1 border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
            <div class="flex gap-0.5 shrink-0">
              <button @click="moveItem(idx, -1)" :disabled="idx === 0"
                class="px-1.5 py-1 text-xs text-slate-400 hover:text-slate-600 disabled:opacity-30">↑</button>
              <button @click="moveItem(idx, 1)" :disabled="idx === modelValue.items.length - 1"
                class="px-1.5 py-1 text-xs text-slate-400 hover:text-slate-600 disabled:opacity-30">↓</button>
              <button @click="removeItem(idx)"
                class="px-1.5 py-1 text-xs text-red-400 hover:text-red-600">×</button>
            </div>
          </div>
          <LinkPicker :modelValue="item.link" @update:modelValue="updateItem(idx, 'link', $event)" label="跳转" />
        </div>
      </div>
    </div>

    <!-- Colors -->
    <ColorInput :modelValue="modelValue.color" @update:modelValue="updateField('color', $event)" label="文字颜色" />
    <ColorInput :modelValue="modelValue.bgColor" @update:modelValue="updateField('bgColor', $event)" label="背景颜色" />

    <!-- Duration -->
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">滚动间隔 (ms)</label>
      <div class="flex items-center gap-3">
        <input type="range" :value="modelValue.duration" @input="updateField('duration', Number(($event.target as HTMLInputElement).value))"
          min="1000" max="10000" step="500" class="flex-1 accent-blue-600" />
        <span class="text-xs text-slate-600 w-12 text-right">{{ modelValue.duration }}</span>
      </div>
    </div>

    <!-- Mode -->
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">样式</label>
      <select :value="modelValue.mode" @change="updateField('mode', ($event.target as HTMLSelectElement).value)"
        class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400">
        <option value="link">可点击跳转</option>
        <option value="closable">可关闭</option>
      </select>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { NoticeProps } from '@/types/decor'
import ColorInput from '../widgets/ColorInput.vue'
import LinkPicker from '../widgets/LinkPicker.vue'

const props = defineProps<{ modelValue: NoticeProps }>()
const emit = defineEmits<{ 'update:modelValue': [value: NoticeProps] }>()

function updateField<K extends keyof NoticeProps>(key: K, value: NoticeProps[K]) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}

function addItem() {
  emit('update:modelValue', {
    ...props.modelValue,
    items: [...props.modelValue.items, { text: '', link: '' }],
  })
}

function removeItem(idx: number) {
  const items = [...props.modelValue.items]
  items.splice(idx, 1)
  emit('update:modelValue', { ...props.modelValue, items })
}

function moveItem(idx: number, delta: number) {
  const target = idx + delta
  if (target < 0 || target >= props.modelValue.items.length) return
  const items = [...props.modelValue.items]
  ;[items[idx], items[target]] = [items[target], items[idx]]
  emit('update:modelValue', { ...props.modelValue, items })
}

function updateItem(idx: number, field: 'text' | 'link', value: string) {
  const items = props.modelValue.items.map((item, i) =>
    i === idx ? { ...item, [field]: value } : item
  )
  emit('update:modelValue', { ...props.modelValue, items })
}
</script>
