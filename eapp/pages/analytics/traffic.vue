<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import DateRangePicker from '@/components/biz/DateRangePicker.vue'
import MetricCard from '@/components/biz/MetricCard.vue'
import AreaChart from '@/components/charts/AreaChart.vue'
import RingChart from '@/components/charts/RingChart.vue'
import BarChart from '@/components/charts/BarChart.vue'
import { useAnalytics } from '@/composables/useAnalytics'
import { getTrafficAnalytics } from '@/api/analytics'

const { t } = useI18n()
const { range, loading, data, load, setRange } = useAnalytics(getTrafficAnalytics)

const summary = computed(() => data.value?.summary || {} as any)

onLoad(() => load())
</script>

<template>
  <view class="page">
    <DateRangePicker :model-value="range" @update:model-value="setRange" />

    <view class="body">
      <view class="metric-row">
        <MetricCard :title="t('analytics.totalPv')" :value="String(summary.total_pv || 0)" color="#2563eb" />
        <MetricCard :title="t('analytics.totalUv')" :value="String(summary.total_uv || 0)" color="#16a34a" />
        <MetricCard :title="t('analytics.avgStay')" :value="`${summary.avg_stay_seconds || 0}s`" color="#f59e0b" />
        <MetricCard :title="t('analytics.bounceRate')" :value="`${Number(summary.bounce_rate || 0).toFixed(1)}%`" color="#ef4444" />
      </view>

      <AreaChart :title="t('analytics.pvUvTrend')" :loading="loading" :data="data?.pv_uv_trend" />
      <RingChart :title="t('analytics.channelDist')" :loading="loading" :data="data?.channel_distribution" />
      <RingChart :title="t('analytics.deviceDist')" :loading="loading" :data="data?.device_distribution" />
      <BarChart :title="t('analytics.pageStay')" :loading="loading" :data="data?.page_stay" />
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.body { padding: 20rpx; }
.metric-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16rpx; margin-bottom: 24rpx; }
</style>
