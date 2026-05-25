<template>
  <div class="space-y-4">
    <div>
      <label class="block text-xs text-slate-500 mb-1.5">{{ $t('decor.marketingZone.type') }}</label>
      <select :value="modelValue.type" @change="emit('update:modelValue', { ...modelValue, type: ($event.target as HTMLSelectElement).value })"
        class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400">
        <option value="seckill">{{ $t('decor.marketingZone.seckill') }}</option>
        <option value="group_buy">{{ $t('decor.marketingZone.groupBuy') }}</option>
        <option value="bargain">{{ $t('decor.marketingZone.bargain') }}</option>
      </select>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { MarketingZoneProps } from '@/types/decor'

defineProps<{ modelValue: MarketingZoneProps }>()
const emit = defineEmits<{ 'update:modelValue': [value: MarketingZoneProps] }>()
</script>
