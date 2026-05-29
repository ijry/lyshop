export const EAPP_CHART_COLORS = ['#2563EB', '#16A34A', '#F59E0B', '#DC2626', '#F97316', '#8B5CF6', '#06B6D4', '#EC4899']

export type ChartCategoriesSeries = {
  categories: string[]
  series: Array<{ name: string; data: number[] }>
}

export type ChartPie = Array<{ name: string; value: number }>

export function buildLineOption(data: ChartCategoriesSeries, opts?: { area?: boolean }) {
  return {
    color: EAPP_CHART_COLORS,
    legend: { data: data.series.length > 1 ? data.series.map((s) => s.name) : [] },
    xAxis: { type: 'category', data: data.categories, axisLabel: { color: '#94a3b8', fontSize: 10 }, axisLine: { show: false } },
    yAxis: { type: 'value', splitLine: { lineStyle: { type: 'dashed', color: '#e2e8f0' } }, axisLabel: { color: '#94a3b8', fontSize: 10 } },
    series: data.series.map((s) => ({
      name: s.name,
      type: 'line',
      data: s.data,
      smooth: true,
      ...(opts?.area ? { areaStyle: { opacity: 0.15 } } : {}),
    })),
  }
}

export function buildAreaOpts(data: ChartCategoriesSeries) {
  return buildLineOption(data, { area: true })
}

export function buildPieOption(data: ChartPie, opts?: { ring?: boolean }) {
  return {
    color: EAPP_CHART_COLORS,
    legend: { orient: 'vertical', right: 10, top: 'center', data: data.map((d) => d.name) },
    series: [{
      type: 'pie',
      radius: opts?.ring ? ['45%', '70%'] : '65%',
      center: ['35%', '50%'],
      data: data.map((d) => ({ name: d.name, value: d.value })),
      label: { show: false },
    }],
  }
}

export function buildRingOpts() {
  return { ring: true }
}

export function buildBarOption(data: ChartCategoriesSeries, opts?: { horizontal?: boolean }) {
  return {
    color: EAPP_CHART_COLORS,
    xAxis: opts?.horizontal
      ? { type: 'value', splitLine: { lineStyle: { type: 'dashed', color: '#e2e8f0' } }, axisLabel: { color: '#94a3b8' } }
      : { type: 'category', data: data.categories, axisLabel: { color: '#94a3b8', fontSize: 10 } },
    yAxis: opts?.horizontal
      ? { type: 'category', data: data.categories, axisLabel: { color: '#94a3b8', fontSize: 10 } }
      : { type: 'value', splitLine: { lineStyle: { type: 'dashed', color: '#e2e8f0' } }, axisLabel: { color: '#94a3b8' } },
    series: data.series.map((s) => ({
      name: s.name,
      type: 'bar',
      data: s.data,
      barWidth: 20,
      itemStyle: { borderRadius: [4, 4, 0, 0] },
    })),
  }
}

export function buildBarOpts(horizontal = false) {
  return { horizontal }
}

export type ChartRadar = {
  indicator: Array<{ name: string; max: number }>
  data: Array<{ name: string; value: number[] }>
}

export function buildRadarOption(data: ChartRadar) {
  return {
    color: EAPP_CHART_COLORS,
    radar: { indicator: data.indicator, shape: 'polygon' },
    series: [{ type: 'radar', data: data.data }],
  }
}

export function buildGaugeOption(value: number, label?: string, max?: number) {
  return {
    series: [{
      type: 'gauge',
      data: [{ value, name: label || '' }],
      max: max || 100,
      detail: { formatter: '{value}%' },
    }],
  }
}
