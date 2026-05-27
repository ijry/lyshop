<template>
  <div class="space-y-4">
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">列数</label>
      <div class="flex gap-2">
        <button v-for="n in [3, 4, 5]" :key="n" @click="update('columns', n)"
          :class="modelValue.columns === n ? 'bg-blue-600 text-white' : 'bg-slate-100 text-slate-700 hover:bg-slate-200'"
          class="px-3 py-1.5 rounded-lg text-xs transition">{{ n }}列</button>
      </div>
    </div>

    <div class="flex items-center justify-between">
      <label class="text-xs text-slate-500">快捷入口</label>
      <button @click="addItem" class="text-xs text-blue-600 hover:underline">+ 添加</button>
    </div>

    <div class="space-y-2">
      <div v-for="(item, idx) in (modelValue.items || [])" :key="idx"
        class="border border-slate-200 rounded-lg p-3 space-y-2.5">
        <div class="flex items-center justify-between">
          <span class="text-xs text-slate-400">项 {{ idx + 1 }}</span>
          <div class="flex gap-0.5">
            <button @click="moveItem(idx, -1)" :disabled="idx === 0"
              class="px-1.5 py-1 text-xs text-slate-400 hover:text-slate-600 disabled:opacity-30">↑</button>
            <button @click="moveItem(idx, 1)" :disabled="idx === (modelValue.items || []).length - 1"
              class="px-1.5 py-1 text-xs text-slate-400 hover:text-slate-600 disabled:opacity-30">↓</button>
            <button @click="removeItem(idx)" class="px-1.5 py-1 text-xs text-red-400 hover:text-red-600">×</button>
          </div>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <div>
            <label class="block text-xs text-slate-500 mb-1">标题</label>
            <input :value="item.title" @input="updateItem(idx, 'title', ($event.target as HTMLInputElement).value)"
              class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
          </div>
          <div>
            <label class="block text-xs text-slate-500 mb-1">图标 (emoji)</label>
            <input :value="item.icon" @input="updateItem(idx, 'icon', ($event.target as HTMLInputElement).value)"
              class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <ColorInput :modelValue="item.bg || '#f5f5f5'" @update:modelValue="updateItem(idx, 'bg', $event)" label="背景色" />
          <div>
            <label class="block text-xs text-slate-500 mb-1">链接</label>
            <input :value="item.link" @input="updateItem(idx, 'link', ($event.target as HTMLInputElement).value)"
              class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
          </div>
        </div>
      </div>
    </div>

    <div v-if="!(modelValue.items || []).length"
      class="text-center py-4 text-xs text-slate-300 border-2 border-dashed border-slate-200 rounded-lg">
      暂无入口项
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

function addItem() {
  const items = [...(props.modelValue.items || []), { title: '', icon: '⭐', bg: '#f5f5f5', link: '' }]
  emit('update:modelValue', { ...props.modelValue, items })
}

function removeItem(idx: number) {
  const items = [...(props.modelValue.items || [])]
  items.splice(idx, 1)
  emit('update:modelValue', { ...props.modelValue, items })
}

function moveItem(idx: number, delta: number) {
  const items = [...(props.modelValue.items || [])]
  const target = idx + delta
  if (target < 0 || target >= items.length) return
  ;[items[idx], items[target]] = [items[target], items[idx]]
  emit('update:modelValue', { ...props.modelValue, items })
}

function updateItem(idx: number, field: string, value: string) {
  const items = (props.modelValue.items || []).map((item: any, i: number) =>
    i === idx ? { ...item, [field]: value } : item
  )
  emit('update:modelValue', { ...props.modelValue, items })
}
</script>
