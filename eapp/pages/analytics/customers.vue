<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import DateRangePicker from '@/components/biz/DateRangePicker.vue'
import MetricCard from '@/components/biz/MetricCard.vue'
import RingChart from '@/components/charts/RingChart.vue'
import BarChart from '@/components/charts/BarChart.vue'
import RadarChart from '@/components/charts/RadarChart.vue'
import { useAnalytics } from '@/composables/useAnalytics'
import { getCustomerAnalytics } from '@/api/analytics'

const { t } = useI18n()
const { range, loading, data, load, setRange } = useAnalytics(getCustomerAnalytics)

const summary = computed(() => data.value?.summary || {} as any)

onLoad(() => load())
</script>

<template>
  <view class="page">
    <DateRangePicker :model-value="range" @update:model-value="setRange" />

    <view class="body">
      <view class="metric-row">
        <MetricCard :title="t('analytics.totalCustomers')" :value="String(summary.total_customers || 0)" color="#2563eb" />
        <MetricCard :title="t('analytics.repurchaseRate')" :value="`${Number(summary.repurchase_rate || 0).toFixed(1)}%`" color="#16a34a" />
        <MetricCard :title="t('analytics.avgCustomerValue')" :value="`¥${Number(summary.avg_customer_value || 0).toFixed(0)}`" color="#f59e0b" />
      </view>

      <RingChart :title="t('analytics.newVsReturning')" :loading="loading" :data="data?.new_vs_returning" />
      <BarChart :title="t('analytics.orderValueDist')" :loading="loading" :data="data?.order_value_distribution" />
      <BarChart :title="t('analytics.purchaseFreq')" :loading="loading" :data="data?.purchase_frequency" />
      <RadarChart :title="t('analytics.rfmRadar')" :loading="loading" :data="data?.rfm_radar" />
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.body { padding: 20rpx; }
.metric-row { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16rpx; margin-bottom: 24rpx; }
</style>
