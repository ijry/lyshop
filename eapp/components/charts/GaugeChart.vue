<script setup lang="ts">
import { computed } from 'vue'
import ChartPanel from './ChartPanel.vue'
import { buildGaugeOption } from '@/utils/ly-charts'

const props = defineProps<{
  title?: string
  loading?: boolean
  error?: string
  value?: number
  label?: string
  max?: number
  height?: string
}>()

const isEmpty = computed(() => props.value == null)
const option = computed(() => props.value != null ? buildGaugeOption(props.value, props.label, props.max) : {})
</script>

<template>
  <ChartPanel :title="title" :loading="loading" :error="error" :empty="isEmpty">
    <template #extra><slot name="extra" /></template>
    <ly-charts-gauge v-if="!isEmpty" :option="option" width="100%" :height="height || '260px'" />
  </ChartPanel>
</template>
