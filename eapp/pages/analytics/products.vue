<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import DateRangePicker from '@/components/biz/DateRangePicker.vue'
import MetricCard from '@/components/biz/MetricCard.vue'
import BarChart from '@/components/charts/BarChart.vue'
import RingChart from '@/components/charts/RingChart.vue'
import { useAnalytics } from '@/composables/useAnalytics'
import { getProductAnalytics } from '@/api/analytics'

const { t } = useI18n()
const { range, loading, data, load, setRange } = useAnalytics(getProductAnalytics)

const summary = computed(() => data.value?.summary || {} as any)

const skuBarData = computed(() => {
  const ranking = data.value?.sku_ranking || []
  return {
    categories: ranking.map((r: any) => {
      const title = String(r.title || '')
      return title.length > 6 ? title.slice(0, 6) + '...' : title
    }),
    series: [{ name: '销量', data: ranking.map((r: any) => Number(r.sold_qty || 0)) }],
  }
})

onLoad(() => load())
</script>

<template>
  <view class="page">
    <DateRangePicker :model-value="range" @update:model-value="setRange" />

    <view class="body">
      <view class="metric-row">
        <MetricCard :title="t('analytics.totalSku')" :value="String(summary.total_sku || 0)" color="#2563eb" />
        <MetricCard :title="t('analytics.activeSkuRate')" :value="`${(Number(summary.active_sku_rate || 0) * 100).toFixed(0)}%`" color="#16a34a" />
        <MetricCard :title="t('analytics.avgTurnover')" :value="String(summary.avg_inventory_turnover || 0)" color="#f59e0b" />
      </view>

      <BarChart :title="t('analytics.skuRanking')" :loading="loading" :data="skuBarData" horizontal />
      <RingChart :title="t('analytics.categorySales')" :loading="loading" :data="data?.category_sales" />
      <RingChart :title="t('analytics.priceRange')" :loading="loading" :data="data?.price_range" />
      <RingChart :title="t('analytics.inventoryStatus')" :loading="loading" :data="data?.inventory_status" />
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.body { padding: 20rpx; }
.metric-row { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16rpx; margin-bottom: 24rpx; }
</style>
