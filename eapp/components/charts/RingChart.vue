<script setup lang="ts">
import { computed } from 'vue'
import ChartPanel from './ChartPanel.vue'
import { buildPieOption, type ChartPie } from '@/utils/ly-charts'

const props = defineProps<{
  title?: string
  loading?: boolean
  error?: string
  data?: ChartPie
  height?: string
}>()

const isEmpty = computed(() => !props.data?.length || props.data.every((r) => Number(r.value || 0) === 0))
const option = computed(() => props.data ? buildPieOption(props.data, { ring: true }) : {})
</script>

<template>
  <ChartPanel :title="title" :loading="loading" :error="error" :empty="isEmpty">
    <template #extra><slot name="extra" /></template>
    <ly-charts-pie v-if="!isEmpty" :option="option" width="100%" :height="height || '260px'" />
  </ChartPanel>
</template>
