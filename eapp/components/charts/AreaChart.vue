<script setup lang="ts">
import { computed } from 'vue'
import ChartPanel from './ChartPanel.vue'
import { buildLineOption, type ChartCategoriesSeries } from '@/utils/ly-charts'

const props = defineProps<{
  title?: string
  loading?: boolean
  error?: string
  data?: ChartCategoriesSeries
  height?: string
}>()

const isEmpty = computed(() => !props.data?.categories?.length || !props.data?.series?.length)
const option = computed(() => props.data ? buildLineOption(props.data, { area: true }) : {})
</script>

<template>
  <ChartPanel :title="title" :loading="loading" :error="error" :empty="isEmpty">
    <template #extra><slot name="extra" /></template>
    <ly-charts-line v-if="!isEmpty" :option="option" width="100%" :height="height || '260px'" />
  </ChartPanel>
</template>
