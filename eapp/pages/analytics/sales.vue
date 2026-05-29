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
import { getSalesAnalytics } from '@/api/analytics'

const { t } = useI18n()
const { range, loading, data, load, setRange } = useAnalytics(getSalesAnalytics)

const summary = computed(() => data.value?.summary || {} as any)
const hourlyData = computed(() => data.value?.hourly_distribution || null)

onLoad(() => load())
</script>

<template>
  <view class="page">
    <DateRangePicker :model-value="range" @update:model-value="setRange" />

    <view class="body">
      <view class="metric-row">
        <MetricCard :title="t('analytics.gmv')" :value="`¥${Number(summary.gmv || 0).toFixed(0)}`" color="#2563eb" />
        <MetricCard :title="t('analytics.orderCount')" :value="String(summary.order_count || 0)" color="#16a34a" />
        <MetricCard :title="t('analytics.avgOrderValue')" :value="`¥${Number(summary.avg_order_value || 0).toFixed(0)}`" color="#f59e0b" />
      </view>

      <AreaChart :title="t('analytics.revenueChart')" :loading="loading" :data="data?.revenue_trend" />
      <AreaChart :title="t('analytics.orderChart')" :loading="loading" :data="data?.order_trend" />
      <RingChart :title="t('analytics.payMethod')" :loading="loading" :data="data?.pay_method_distribution" />
      <BarChart :title="t('analytics.hourlyChart')" :loading="loading" :data="hourlyData" />
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.body { padding: 20rpx; }
.metric-row { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16rpx; margin-bottom: 24rpx; }
</style>
