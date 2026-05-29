import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock coupon CRUD', () => {
  it('GET /marketing/coupons returns expanded data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/marketing/coupons', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(3)
    expect(r.data.list[0]).toHaveProperty('stack_rule')
    expect(r.data.list[0]).toHaveProperty('target_type')
  })
  it('filters by keyword', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/marketing/coupons', { keyword: '新人', page: 1, size: 20 })
    expect(r.matched).toBe(true)
    for (const c of r.data.list) expect(c.name.toLowerCase()).toContain('新人')
  })
  it('POST creates a coupon', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/marketing/coupons', { name: '测试券', type: 1, discount: 10, status: 1, stack_rule: 'exclusive', target_type: 'all' })
    expect(r.matched).toBe(true)
    expect(r.data.id).toBeGreaterThan(0)
  })
  it('PUT updates a coupon', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/marketing/coupons/1', { name: '新人满100减25' })
    expect(r.matched).toBe(true)
  })
  it('DELETE removes a coupon', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const before = matchMock('GET', '/admin/api/marketing/coupons', { page: 1, size: 50 })
    const countBefore = before.data.list.length
    matchMock('DELETE', '/admin/api/marketing/coupons/1', {})
    const after = matchMock('GET', '/admin/api/marketing/coupons', { page: 1, size: 50 })
    expect(after.data.list.length).toBeLessThan(countBefore)
  })
  it('POST send increments used_count', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/marketing/coupons/2/send', { count: 5 })
    expect(r.matched).toBe(true)
    expect(r.data.sent_count).toBe(5)
  })
})
