import { describe, expect, it } from 'vitest'
import {
  EAPP_CHART_COLORS,
  buildLineOption,
  buildPieOption,
  buildBarOption,
  buildAreaOpts,
  buildRingOpts,
  buildBarOpts,
} from '@/utils/ly-charts'

describe('ly-charts utils', () => {
  it('brand palette starts with #2563EB', () => {
    expect(EAPP_CHART_COLORS[0]).toBe('#2563EB')
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
})
