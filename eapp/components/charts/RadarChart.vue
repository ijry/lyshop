<script setup lang="ts">
import { computed } from 'vue'
import ChartPanel from './ChartPanel.vue'
import { buildRadarOption, type ChartRadar } from '@/utils/ly-charts'

const props = defineProps<{
  title?: string
  loading?: boolean
  error?: string
  data?: ChartRadar
  height?: string
}>()

const isEmpty = computed(() => !props.data?.indicator?.length || !props.data?.data?.length)
const option = computed(() => props.data ? buildRadarOption(props.data) : {})
</script>

<template>
  <ChartPanel :title="title" :loading="loading" :error="error" :empty="isEmpty">
    <template #extra><slot name="extra" /></template>
    <ly-charts-radar v-if="!isEmpty" :option="option" width="100%" :height="height || '260px'" />
  </ChartPanel>
</template>
