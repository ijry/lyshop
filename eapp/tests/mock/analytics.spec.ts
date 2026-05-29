import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock analytics routes', () => {
  it('GET /analytics/sales returns expected structure', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/analytics/sales', { days: 7 })
    expect(r.matched).toBe(true)
    const d = r.data
    expect(d.summary.gmv).toBeTypeOf('number')
    expect(d.summary.order_count).toBeTypeOf('number')
    expect(d.summary.avg_order_value).toBeTypeOf('number')
    expect(d.summary.refund_rate).toBeTypeOf('number')
    expect(d.revenue_trend.categories).toHaveLength(7)
    expect(d.order_trend.categories).toHaveLength(7)
    expect(d.pay_method_distribution).toHaveLength(3)
    expect(d.hourly_distribution.categories).toHaveLength(24)
    expect(d.hourly_distribution.series[0].data).toHaveLength(24)
  })

  it('GET /analytics/products returns expected structure', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/analytics/products', { days: 7 })
    expect(r.matched).toBe(true)
    const d = r.data
    expect(d.summary.total_sku).toBeTypeOf('number')
    expect(d.summary.active_sku_rate).toBeTypeOf('number')
    expect(d.sku_ranking.length).toBeGreaterThan(0)
    expect(d.sku_ranking[0].rank).toBe(1)
    expect(d.category_sales.length).toBeGreaterThan(0)
    expect(d.price_range).toHaveLength(5)
    expect(d.inventory_status).toHaveLength(3)
  })

  it('GET /analytics/customers returns expected structure', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/analytics/customers', { days: 7 })
    expect(r.matched).toBe(true)
    const d = r.data
    expect(d.summary.total_customers).toBeTypeOf('number')
    expect(d.summary.repurchase_rate).toBeTypeOf('number')
    expect(d.new_vs_returning).toHaveLength(2)
    expect(d.order_value_distribution.categories).toHaveLength(5)
    expect(d.purchase_frequency.categories).toHaveLength(4)
    expect(d.rfm_radar.indicator).toHaveLength(3)
    expect(d.rfm_radar.data[0].value).toHaveLength(3)
  })

  it('GET /analytics/traffic returns expected structure', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/analytics/traffic', { days: 7 })
    expect(r.matched).toBe(true)
    const d = r.data
    expect(d.summary.total_pv).toBeTypeOf('number')
    expect(d.summary.total_uv).toBeTypeOf('number')
    expect(d.summary.avg_stay_seconds).toBeTypeOf('number')
    expect(d.summary.bounce_rate).toBeTypeOf('number')
    expect(d.pv_uv_trend.series).toHaveLength(2)
    expect(d.channel_distribution).toHaveLength(5)
    expect(d.device_distribution).toHaveLength(4)
    expect(d.page_stay.categories).toHaveLength(5)
  })

  it('GET /analytics/conversion returns expected structure', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/analytics/conversion', { days: 7 })
    expect(r.matched).toBe(true)
    const d = r.data
    expect(d.summary.visit_to_cart_rate).toBeTypeOf('number')
    expect(d.summary.overall_rate).toBeTypeOf('number')
    expect(d.funnel).toHaveLength(4)
    expect(d.funnel[0].step).toBe('访问')
    expect(d.funnel[0].rate).toBe(100)
    expect(d.step_rates_trend.series).toHaveLength(3)
  })

  it('analytics data is deterministic', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const a = matchMock('GET', '/admin/api/analytics/sales', { days: 7 }).data
    const b = matchMock('GET', '/admin/api/analytics/sales', { days: 7 }).data
    expect(a.summary.gmv).toBe(b.summary.gmv)
    expect(a.revenue_trend.series[0].data).toEqual(b.revenue_trend.series[0].data)
  })

  it('respects days parameter', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r7 = matchMock('GET', '/admin/api/analytics/sales', { days: 7 }).data
    const r30 = matchMock('GET', '/admin/api/analytics/sales', { days: 30 }).data
    expect(r7.revenue_trend.categories).toHaveLength(7)
    expect(r30.revenue_trend.categories).toHaveLength(30)
  })
})
