import { describe, expect, it } from 'vitest'
import {
  EAPP_CHART_COLORS,
  buildLineOption,
  buildPieOption,
  buildBarOption,
  buildAreaOpts,
  buildRingOpts,
  buildBarOpts,
  buildRadarOption,
  buildGaugeOption,
  paletteColor,
} from '@/utils/ly-charts'

describe('ly-charts utils', () => {
  it('brand palette starts with the merchant primary color', () => {
    expect(EAPP_CHART_COLORS[0]).toBe('#1D4ED8')
    expect(EAPP_CHART_COLORS.length).toBeGreaterThanOrEqual(4)
  })

  it('buildLineOption hides legend for single series', () => {
    const opt = buildLineOption({
      categories: ['Mon', 'Tue'],
      series: [{ name: 'Sales', data: [10, 20] }],
    })
    expect(opt.legend.data).toHaveLength(0)
  })

  it('buildLineOption shows legend for multi series', () => {
    const opt = buildLineOption({
      categories: ['Mon', 'Tue'],
      series: [
        { name: 'Sales', data: [10, 20] },
        { name: 'Returns', data: [2, 3] },
      ],
    })
    expect(opt.legend.data).toEqual(['Sales', 'Returns'])
  })

  it('buildAreaOpts produces area-enabled line option', () => {
    const opt = buildAreaOpts({
      categories: ['Mon'],
      series: [{ name: 'A', data: [1] }],
    })
    expect(opt.series[0].areaStyle).toBeDefined()
    expect(opt.series[0].smooth).toBe(true)
  })

  it('buildPieOption produces ring when ring=true', () => {
    const opt = buildPieOption([{ name: 'A', value: 10 }], { ring: true })
    expect(opt.series[0].radius).toEqual(['45%', '70%'])
  })

  it('buildPieOption produces full pie when ring=false', () => {
    const opt = buildPieOption([{ name: 'A', value: 10 }])
    expect(opt.series[0].radius).toBe('65%')
  })

  it('buildRingOpts returns ring flag', () => {
    expect(buildRingOpts().ring).toBe(true)
  })

  it('buildBarOption sets horizontal axes correctly', () => {
    const v = buildBarOption({ categories: ['A'], series: [{ name: 'x', data: [1] }] })
    expect(v.xAxis.type).toBe('category')
    const h = buildBarOption({ categories: ['A'], series: [{ name: 'x', data: [1] }] }, { horizontal: true })
    expect(h.xAxis.type).toBe('value')
    expect(h.yAxis.type).toBe('category')
  })

  it('buildBarOpts returns horizontal flag', () => {
    expect(buildBarOpts(false).horizontal).toBe(false)
    expect(buildBarOpts(true).horizontal).toBe(true)
  })

  // ly-charts 组件不读取顶层 option.color，必须逐系列 / 逐扇区显式着色
  it('explicitly colors line and bar series with the merchant palette', () => {
    const source = {
      categories: ['A', 'B'],
      series: [
        { name: 'x', data: [1, 2] },
        { name: 'y', data: [3, 4] },
      ],
    }
    const line = buildLineOption(source)
    const bar = buildBarOption(source)
    expect(line.series.map((s) => s.color)).toEqual([EAPP_CHART_COLORS[0], EAPP_CHART_COLORS[1]])
    expect(line.series.map((s) => s.itemStyle.color)).toEqual([EAPP_CHART_COLORS[0], EAPP_CHART_COLORS[1]])
    expect(bar.series.map((s) => s.color)).toEqual([EAPP_CHART_COLORS[0], EAPP_CHART_COLORS[1]])
    expect(bar.series.map((s) => s.itemStyle.color)).toEqual([EAPP_CHART_COLORS[0], EAPP_CHART_COLORS[1]])
  })

  it('explicitly colors every pie sector and legend dot', () => {
    const opt = buildPieOption([
      { name: 'A', value: 1 },
      { name: 'B', value: 2 },
      { name: 'C', value: 3 },
    ], { ring: true })
    expect(opt.series[0].data.map((d) => d.color)).toEqual(EAPP_CHART_COLORS.slice(0, 3))
    expect(opt.series[0].data.map((d) => d.itemStyle.color)).toEqual(EAPP_CHART_COLORS.slice(0, 3))
  })

  it('colors radar and gauge with the brand primary color', () => {
    const radar = buildRadarOption({
      indicator: [{ name: 'A', max: 100 }],
      data: [{ name: 'X', value: [10] }],
    })
    const gauge = buildGaugeOption(60)
    expect(radar.series[0].itemStyle.color).toBe(EAPP_CHART_COLORS[0])
    expect(radar.series[0].data[0].color).toBe(EAPP_CHART_COLORS[0])
    expect(gauge.series[0].itemStyle.color).toBe(EAPP_CHART_COLORS[0])
    expect(gauge.series[0].data[0].itemStyle.color).toBe(EAPP_CHART_COLORS[0])
  })

  it('paletteColor wraps around the palette safely', () => {
    expect(paletteColor(0)).toBe(EAPP_CHART_COLORS[0])
    expect(paletteColor(EAPP_CHART_COLORS.length)).toBe(EAPP_CHART_COLORS[0])
    expect(paletteColor(-1)).toBe(EAPP_CHART_COLORS[EAPP_CHART_COLORS.length - 1])
  })
})
