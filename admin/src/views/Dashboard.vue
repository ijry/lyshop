<template>
  <div>
    <h2 class="text-xl font-semibold text-slate-800 mb-6">控制台</h2>
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div
        v-for="card in cards"
        :key="card.label"
        class="bg-white rounded-xl shadow-sm p-6 border border-slate-100"
      >
        <p class="text-sm text-slate-500 mb-1">{{ card.label }}</p>
        <p class="text-2xl font-bold text-slate-800">{{ card.value }}</p>
      </div>
    </div>
    <div class="bg-white rounded-xl shadow-sm p-6 border border-slate-100">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-base font-semibold text-slate-800">近 7 天销售趋势</h3>
        <button
          class="text-xs px-3 py-1.5 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50"
          @click="loadDashboard"
        >
          刷新
        </button>
      </div>

      <p v-if="loading" class="text-slate-400 text-center py-12">数据加载中...</p>

      <div v-else-if="trend.length" ref="chartRef" class="h-72 w-full" />

      <p v-else class="text-slate-400 text-center py-12">暂无趋势数据</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { EChartsType } from 'echarts/core'
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { getDashboard } from '@/api/plugins'

type TrendPoint = {
  date: string
  orders: number
  sales: number
}

type DashboardData = {
  today_orders: number
  today_sales: number
  pending_refunds: number
  online_sessions: number
  sales_trend: TrendPoint[]
}

const loading = ref(false)
const chartRef = ref<HTMLElement | null>(null)
let chart: EChartsType | null = null
let chartInit: ((dom: HTMLElement) => EChartsType) | null = null

const dashboard = ref<DashboardData>({
  today_orders: 0,
  today_sales: 0,
  pending_refunds: 0,
  online_sessions: 0,
  sales_trend: [],
})

const cards = computed(() => [
  { label: '今日订单', value: String(dashboard.value.today_orders || 0) },
  { label: '今日销售额', value: `¥${money(dashboard.value.today_sales || 0)}` },
  { label: '待处理售后', value: String(dashboard.value.pending_refunds || 0) },
  { label: '在线客服会话', value: String(dashboard.value.online_sessions || 0) },
])

const trend = computed(() => dashboard.value.sales_trend || [])
const money = (v: number) => Number(v || 0).toFixed(2)
const shortDate = (v: string) => (v || '').slice(5)

async function ensureChartRuntime() {
  if (chartInit) return
  const [core, charts, components, renderers] = await Promise.all([
    import('echarts/core'),
    import('echarts/charts'),
    import('echarts/components'),
    import('echarts/renderers'),
  ])
  core.use([
    charts.BarChart,
    charts.LineChart,
    components.GridComponent,
    components.TooltipComponent,
    components.LegendComponent,
    renderers.CanvasRenderer,
  ])
  chartInit = core.init
}

async function renderChart() {
  if (!chartRef.value || !trend.value.length) return
  await ensureChartRuntime()
  if (!chartInit) return
  if (!chart) chart = chartInit(chartRef.value)

  chart.setOption({
    animation: true,
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' },
    },
    legend: {
      top: 0,
      right: 0,
      data: ['销售额', '订单量'],
    },
    grid: {
      left: 56,
      right: 56,
      top: 40,
      bottom: 32,
    },
    xAxis: {
      type: 'category',
      boundaryGap: true,
      data: trend.value.map(item => shortDate(item.date)),
      axisLine: { lineStyle: { color: '#cbd5e1' } },
      axisLabel: { color: '#64748b' },
    },
    yAxis: [
      {
        type: 'value',
        name: '销售额(元)',
        axisLabel: { color: '#64748b' },
        splitLine: { lineStyle: { color: '#f1f5f9' } },
      },
      {
        type: 'value',
        name: '订单量(单)',
        axisLabel: { color: '#64748b' },
        splitLine: { show: false },
      },
    ],
    series: [
      {
        name: '销售额',
        type: 'bar',
        yAxisIndex: 0,
        barMaxWidth: 28,
        data: trend.value.map(item => Number(item.sales || 0)),
        itemStyle: {
          color: '#ef4444',
          borderRadius: [4, 4, 0, 0],
        },
      },
      {
        name: '订单量',
        type: 'line',
        yAxisIndex: 1,
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        data: trend.value.map(item => Number(item.orders || 0)),
        lineStyle: { width: 2, color: '#2563eb' },
        itemStyle: { color: '#2563eb' },
      },
    ],
  })
}

function onResize() {
  chart?.resize()
}

async function loadDashboard() {
  loading.value = true
  try {
    const data: any = await getDashboard()
    if (!data) return
    dashboard.value = {
      today_orders: Number(data.today_orders || 0),
      today_sales: Number(data.today_sales || 0),
      pending_refunds: Number(data.pending_refunds || 0),
      online_sessions: Number(data.online_sessions || 0),
      sales_trend: Array.isArray(data.sales_trend) ? data.sales_trend : [],
    }
  } finally {
    loading.value = false
  }
  await nextTick()
  if (!trend.value.length) {
    chart?.dispose()
    chart = null
    return
  }
  await renderChart()
}

onMounted(() => {
  loadDashboard()
  window.addEventListener('resize', onResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
  chart?.dispose()
  chart = null
})
</script>
