<template>
  <div class="space-y-4">
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">{{ $t('decor.productGrid.title') }}</label>
      <input type="text" :value="modelValue.title || ''" @input="update('title', ($event.target as HTMLInputElement).value)"
        :placeholder="$t('decor.productGrid.titlePlaceholder')"
        class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400" />
    </div>

    <div>
      <label class="block text-xs text-slate-500 mb-1.5">{{ $t('decor.productGrid.source') }}</label>
      <select :value="modelValue.source" @change="update('source', ($event.target as HTMLSelectElement).value)"
        class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400">
        <option value="hot">{{ $t('decor.productGrid.hot') }}</option>
        <option value="new">{{ $t('decor.productGrid.latest') }}</option>
        <option value="recommend">{{ $t('decor.productGrid.recommend') }}</option>
      </select>
    </div>

    <div>
      <label class="block text-xs text-slate-500 mb-1.5">{{ $t('decor.productGrid.limit') }}</label>
      <div class="flex items-center gap-3">
        <input type="range" :value="modelValue.limit" @input="update('limit', Number(($event.target as HTMLInputElement).value))"
          min="2" max="50" step="2" class="flex-1 accent-blue-600" />
        <input type="number" :value="modelValue.limit" @change="update('limit', Number(($event.target as HTMLInputElement).value))"
          min="2" max="50" class="w-16 border border-slate-200 rounded-lg px-2 py-1.5 text-xs text-center focus:outline-none focus:border-blue-400" />
      </div>
    </div>

    <div>
      <label class="block text-xs text-slate-500 mb-1.5">{{ $t('decor.productGrid.columns') }}</label>
      <div class="flex gap-2">
        <button v-for="n in [2, 3]" :key="n" @click="update('columns', n)"
          :class="modelValue.columns === n ? 'bg-blue-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'"
          class="flex-1 py-2 rounded-lg text-sm transition">
          {{ n }} {{ $t('decor.productGrid.columnUnit') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ProductGridProps } from '@/types/decor'

const props = defineProps<{ modelValue: ProductGridProps }>()
const emit = defineEmits<{ 'update:modelValue': [value: ProductGridProps] }>()

function update<K extends keyof ProductGridProps>(key: K, value: ProductGridProps[K]) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}
</script>
