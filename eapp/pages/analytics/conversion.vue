<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import DateRangePicker from '@/components/biz/DateRangePicker.vue'
import MetricCard from '@/components/biz/MetricCard.vue'
import FunnelChart from '@/components/biz/FunnelChart.vue'
import AreaChart from '@/components/charts/AreaChart.vue'
import GaugeChart from '@/components/charts/GaugeChart.vue'
import { useAnalytics } from '@/composables/useAnalytics'
import { getConversionAnalytics } from '@/api/analytics'

const { t } = useI18n()
const { range, loading, data, load, setRange } = useAnalytics(getConversionAnalytics)

const summary = computed(() => data.value?.summary || {} as any)

onLoad(() => load())
</script>

<template>
  <view class="page">
    <DateRangePicker :model-value="range" @update:model-value="setRange" />

    <view class="body">
      <view class="metric-row">
        <MetricCard :title="t('analytics.overallRate')" :value="`${Number(summary.overall_rate || 0).toFixed(2)}%`" color="#2563eb" />
        <MetricCard :title="t('analytics.abandonedCart')" :value="`${Number(summary.abandoned_cart_rate || 0).toFixed(1)}%`" color="#ef4444" />
      </view>

      <FunnelChart :title="t('analytics.funnelChart')" :loading="loading" :steps="data?.funnel" />
      <AreaChart :title="t('analytics.stepRatesTrend')" :loading="loading" :data="data?.step_rates_trend" />
      <GaugeChart
        :title="t('analytics.overallRate')"
        :loading="loading"
        :value="Number(summary.overall_rate || 0)"
        label="转化率"
      />
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.body { padding: 20rpx; }
.metric-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16rpx; margin-bottom: 24rpx; }
</style>
