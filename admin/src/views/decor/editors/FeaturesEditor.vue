<template>
  <div class="space-y-4">
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">列数</label>
      <div class="flex gap-2">
        <button v-for="n in [2, 3, 4]" :key="n" @click="update('columns', n)"
          :class="modelValue.columns === n ? 'bg-blue-600 text-white' : 'bg-slate-100 text-slate-700 hover:bg-slate-200'"
          class="px-3 py-1.5 rounded-lg text-xs transition">{{ n }}列</button>
      </div>
    </div>

    <div class="flex items-center justify-between">
      <label class="text-xs text-slate-500">特性项</label>
      <button @click="addItem" class="text-xs text-blue-600 hover:underline">+ 添加</button>
    </div>

    <div class="space-y-2">
      <div v-for="(item, idx) in (modelValue.items || [])" :key="idx"
        class="border border-slate-200 rounded-lg p-3 space-y-2.5">
        <div class="flex items-center justify-between">
          <span class="text-xs text-slate-400">特性 {{ idx + 1 }}</span>
          <button @click="removeItem(idx)" class="text-xs text-red-400 hover:text-red-600">×</button>
        </div>
        <div>
          <label class="block text-xs text-slate-500 mb-1">图标 class</label>
          <input :value="item.icon" @input="updateItem(idx, 'icon', ($event.target as HTMLInputElement).value)"
            placeholder="i-carbon-delivery-truck" class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
        </div>
        <div class="grid grid-cols-2 gap-2">
          <div>
            <label class="block text-xs text-slate-500 mb-1">标题</label>
            <input :value="item.title" @input="updateItem(idx, 'title', ($event.target as HTMLInputElement).value)"
              class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
          </div>
          <div>
            <label class="block text-xs text-slate-500 mb-1">描述</label>
            <input :value="item.desc" @input="updateItem(idx, 'desc', ($event.target as HTMLInputElement).value)"
              class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
          </div>
        </div>
      </div>
    </div>

    <div v-if="!(modelValue.items || []).length"
      class="text-center py-4 text-xs text-slate-300 border-2 border-dashed border-slate-200 rounded-lg">
      暂无特性项
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ modelValue: any }>()
const emit = defineEmits<{ 'update:modelValue': [value: any] }>()

function update(key: string, value: any) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}

function addItem() {
  const items = [...(props.modelValue.items || []), { icon: 'i-carbon-star', title: '', desc: '' }]
  emit('update:modelValue', { ...props.modelValue, items })
}

function removeItem(idx: number) {
  const items = [...(props.modelValue.items || [])]
  items.splice(idx, 1)
  emit('update:modelValue', { ...props.modelValue, items })
}

function updateItem(idx: number, field: string, value: string) {
  const items = (props.modelValue.items || []).map((item: any, i: number) =>
    i === idx ? { ...item, [field]: value } : item
  )
  emit('update:modelValue', { ...props.modelValue, items })
}
</script>
