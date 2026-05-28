import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock VIP plans CRUD', () => {
  it('GET /vip/plans returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/vip/plans', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(2)
  })
  it('POST /vip/plans creates a new plan', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/vip/plans', { name: '年卡', months: 12, price: 168, status: 1 })
    expect(r.matched).toBe(true)
  })
  it('PUT /vip/plans/1 updates a plan', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/vip/plans/1', { name: '月卡Plus' })
    expect(r.matched).toBe(true)
  })
  it('DELETE /vip/plans/1 removes a plan', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const before = matchMock('GET', '/admin/api/vip/plans', { page: 1, size: 50 })
    const countBefore = before.data.list.length
    matchMock('DELETE', '/admin/api/vip/plans/1', {})
    const after = matchMock('GET', '/admin/api/vip/plans', { page: 1, size: 50 })
    expect(after.data.list.length).toBeLessThan(countBefore)
  })
})

describe('mock VIP levels CRUD', () => {
  it('GET /vip/levels returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/vip/levels', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
  })
  it('POST /vip/levels creates a new level', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/vip/levels', { name: '钻石卡', growth_min: 5000, discount_rate: 0.9 })
    expect(r.matched).toBe(true)
  })
})

describe('mock VIP coupon-rules CRUD', () => {
  it('GET /vip/coupon-rules returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/vip/coupon-rules', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
  })
})

describe('mock VIP sku-prices CRUD', () => {
  it('GET /vip/sku-prices returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/vip/sku-prices', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
  })
  it('filters by product_id', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/vip/sku-prices', { product_id: 1, page: 1, size: 20 })
    expect(r.matched).toBe(true)
    for (const item of r.data.list) expect(Number(item.product_id)).toBe(1)
  })
})
