import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock /dashboard upgraded', () => {
  it('returns trend and new fields', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/dashboard', {})
    expect(r.matched).toBe(true)
    const d = r.data
    expect(d.trend.revenue_7d.categories).toHaveLength(7)
    expect(d.trend.revenue_7d.series[0].data).toHaveLength(7)
    expect(d.trend.revenue_30d.categories).toHaveLength(30)
    expect(d.status_distribution.length).toBeGreaterThanOrEqual(3)
    expect(d.hot_products.length).toBeLessThanOrEqual(5)
    expect(d.announcements.length).toBeGreaterThan(0)
    expect(d.compare.revenue_yoy).toBeTypeOf('number')
    expect(typeof d.today_avg_price).toBe('number')
  })
  it('deterministic trend data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const a = matchMock('GET', '/admin/api/dashboard', {}).data
    const b = matchMock('GET', '/admin/api/dashboard', {}).data
    expect(a.trend.revenue_7d.series[0].data).toEqual(b.trend.revenue_7d.series[0].data)
  })
})

describe('mock /shops/current', () => {
  it('returns shop', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/shops/current', {})
    expect(r.matched).toBe(true)
    expect(r.data.name).toBeTruthy()
  })
})

describe('mock /announcements', () => {
  it('returns list', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/announcements', {})
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(3)
  })
})
