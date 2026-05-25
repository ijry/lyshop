<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <label class="text-xs text-slate-500">{{ $t('decor.categoryNav.items') }}</label>
      <button @click="addItem" class="text-xs text-blue-600 hover:underline">{{ $t('decor.categoryNav.add') }}</button>
    </div>

    <div class="space-y-2">
      <div v-for="(item, idx) in modelValue.items" :key="idx"
        class="border border-slate-200 rounded-lg p-3 space-y-2.5">
        <div class="flex items-center justify-between">
          <span class="text-xs text-slate-400">{{ $t('decor.categoryNav.itemLabel', { idx: idx + 1 }) }}</span>
          <div class="flex gap-0.5">
            <button @click="moveItem(idx, -1)" :disabled="idx === 0"
              class="px-1.5 py-1 text-xs text-slate-400 hover:text-slate-600 disabled:opacity-30">↑</button>
            <button @click="moveItem(idx, 1)" :disabled="idx === modelValue.items.length - 1"
              class="px-1.5 py-1 text-xs text-slate-400 hover:text-slate-600 disabled:opacity-30">↓</button>
            <button @click="removeItem(idx)"
              class="px-1.5 py-1 text-xs text-red-400 hover:text-red-600">×</button>
          </div>
        </div>

        <div>
          <label class="block text-xs text-slate-500 mb-1">{{ $t('decor.categoryNav.title') }}</label>
          <input :value="item.title" @input="updateItem(idx, 'title', ($event.target as HTMLInputElement).value)"
            :placeholder="$t('decor.categoryNav.titlePlaceholder')" class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
        </div>

        <ImageUpload :modelValue="item.icon" @update:modelValue="updateItem(idx, 'icon', $event)" :label="$t('decor.categoryNav.icon')" />
        <LinkPicker :modelValue="item.link" @update:modelValue="updateItem(idx, 'link', $event)" :label="$t('decor.categoryNav.link')" />
      </div>
    </div>

    <div v-if="!modelValue.items.length" class="text-center py-4 text-xs text-slate-300 border-2 border-dashed border-slate-200 rounded-lg">
      {{ $t('decor.categoryNav.empty') }}
    </div>
  </div>
</template>

<script setup lang="ts">
import type { CategoryNavProps } from '@/types/decor'
import ImageUpload from '../widgets/ImageUpload.vue'
import LinkPicker from '../widgets/LinkPicker.vue'

const props = defineProps<{ modelValue: CategoryNavProps }>()
const emit = defineEmits<{ 'update:modelValue': [value: CategoryNavProps] }>()

function addItem() {
  emit('update:modelValue', {
    items: [...props.modelValue.items, { title: '', icon: '', link: '' }],
  })
}

function removeItem(idx: number) {
  const items = [...props.modelValue.items]
  items.splice(idx, 1)
  emit('update:modelValue', { items })
}

function moveItem(idx: number, delta: number) {
  const target = idx + delta
  if (target < 0 || target >= props.modelValue.items.length) return
  const items = [...props.modelValue.items]
  ;[items[idx], items[target]] = [items[target], items[idx]]
  emit('update:modelValue', { items })
}

function updateItem(idx: number, field: 'title' | 'icon' | 'link', value: string) {
  const items = props.modelValue.items.map((item, i) =>
    i === idx ? { ...item, [field]: value } : item
  )
  emit('update:modelValue', { items })
}
</script>
