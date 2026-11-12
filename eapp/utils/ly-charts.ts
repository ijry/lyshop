/**
 * 图表统一色板与坐标轴样式。
 * 与 App.vue 设计令牌保持一致：主色 + 低饱和语义色 + 中性灰蓝，
 * 避免消费级高饱和配色，适合商家后台长时间阅读。
 *
 * 注意：uni_modules/ly-charts 的 line / bar / pie 组件并不读取 ECharts
 * 顶层的 option.color，而是按 series.color、series.itemStyle.color、
 * 或 pie 数据项的 color / itemStyle.color 取色，未指定时回退到组件内置
 * 的 ECharts 高饱和默认色板（#5470c6 / #91cc75 / #fac858 …）。
 * 因此这里必须在 series 与数据项上显式着色，才能真正确保跨端一致。
 */
export const EAPP_CHART_COLORS = [
  '#1D4ED8', // 品牌主蓝
  '#0E7490', // 青（物流 / 流量）
  '#157F4B', // 绿（成功 / 完成）
  '#A75A08', // 琥珀（待办 / 警告）
  '#B42318', // 红（异常 / 取消）
  '#475467', // 石板
  '#7A8AA5', // 灰蓝
  '#A9B4C4', // 浅灰蓝
]

const AXIS_LABEL_COLOR = '#98A2B3'
const SPLIT_LINE_COLOR = '#EEF0F4'

export function paletteColor(index: number): string {
  return EAPP_CHART_COLORS[((index % EAPP_CHART_COLORS.length) + EAPP_CHART_COLORS.length) % EAPP_CHART_COLORS.length]
}

export type ChartCategoriesSeries = {
  categories: string[]
  series: Array<{ name: string; data: number[] }>
}

export type ChartPie = Array<{ name: string; value: number }>

export function buildLineOption(data: ChartCategoriesSeries, opts?: { area?: boolean }) {
  return {
    color: EAPP_CHART_COLORS,
    legend: { data: data.series.length > 1 ? data.series.map((s) => s.name) : [] },
    xAxis: { type: 'category', data: data.categories, axisLabel: { color: AXIS_LABEL_COLOR, fontSize: 10 }, axisLine: { show: false } },
    yAxis: { type: 'value', splitLine: { lineStyle: { type: 'dashed', color: SPLIT_LINE_COLOR } }, axisLabel: { color: AXIS_LABEL_COLOR, fontSize: 10 } },
    series: data.series.map((s, index) => ({
      name: s.name,
      type: 'line',
      data: s.data,
      smooth: true,
      color: paletteColor(index),
      lineStyle: { color: paletteColor(index), width: 2 },
      itemStyle: { color: paletteColor(index) },
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
      // 每个扇区与图例都按同一色板显式着色，避免回退到组件内置彩虹色
      data: data.map((d, index) => ({
        name: d.name,
        value: d.value,
        color: paletteColor(index),
        itemStyle: { color: paletteColor(index) },
      })),
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
      ? { type: 'value', splitLine: { lineStyle: { type: 'dashed', color: SPLIT_LINE_COLOR } }, axisLabel: { color: AXIS_LABEL_COLOR } }
      : { type: 'category', data: data.categories, axisLabel: { color: AXIS_LABEL_COLOR, fontSize: 10 } },
    yAxis: opts?.horizontal
      ? { type: 'category', data: data.categories, axisLabel: { color: AXIS_LABEL_COLOR, fontSize: 10 } }
      : { type: 'value', splitLine: { lineStyle: { type: 'dashed', color: SPLIT_LINE_COLOR } }, axisLabel: { color: AXIS_LABEL_COLOR } },
    series: data.series.map((s, index) => ({
      name: s.name,
      type: 'bar',
      data: s.data,
      barWidth: 20,
      color: paletteColor(index),
      itemStyle: { color: paletteColor(index), borderRadius: [4, 4, 0, 0] },
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
    series: [{
      type: 'radar',
      color: paletteColor(0),
      itemStyle: { color: paletteColor(0) },
      data: data.data.map((d, index) => ({
        ...d,
        color: paletteColor(index),
        itemStyle: { color: paletteColor(index) },
      })),
    }],
  }
}

export function buildGaugeOption(value: number, label?: string, max?: number) {
  const brand = paletteColor(0)
  return {
    color: EAPP_CHART_COLORS,
    series: [{
      type: 'gauge',
      data: [{ value, name: label || '', itemStyle: { color: brand } }],
      max: max || 100,
      detail: { formatter: '{value}%' },
      itemStyle: { color: brand },
      axisLine: { lineStyle: { color: [[1, SPLIT_LINE_COLOR]], width: 12 } },
      progress: { show: true, width: 12, itemStyle: { color: brand } },
    }],
  }
}
