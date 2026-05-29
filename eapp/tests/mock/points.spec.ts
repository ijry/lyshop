import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock points products', () => {
  it('GET /points/products returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/points/products', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(5)
  })
  it('filters by type', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/points/products', { type: 'coupon', page: 1, size: 20 })
    expect(r.matched).toBe(true)
    for (const item of r.data.list) expect(item.type).toBe('coupon')
  })
  it('POST /points/products creates a product', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/points/products', { title: '测试商品', type: 'virtual', points_price: 100, stock: 50, status: 1 })
    expect(r.matched).toBe(true)
    expect(r.data.title).toBe('测试商品')
  })
  it('PUT /points/products/:id updates a product', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/points/products/1', { title: '更新后的商品' })
    expect(r.matched).toBe(true)
  })
  it('DELETE /points/products/:id removes a product', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const before = matchMock('GET', '/admin/api/points/products', { page: 1, size: 50 })
    const countBefore = before.data.list.length
    matchMock('DELETE', '/admin/api/points/products/1', {})
    const after = matchMock('GET', '/admin/api/points/products', { page: 1, size: 50 })
    expect(after.data.list.length).toBeLessThan(countBefore)
  })
})

describe('mock points exchanges', () => {
  it('GET /points/exchanges returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/points/exchanges', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(3)
  })
  it('filters by status', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/points/exchanges', { status: 'completed', page: 1, size: 20 })
    expect(r.matched).toBe(true)
    for (const item of r.data.list) expect(item.status).toBe('completed')
  })
  it('PUT ship changes status', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/points/exchanges/2/ship', {})
    expect(r.matched).toBe(true)
    const after = matchMock('GET', '/admin/api/points/exchanges', { page: 1, size: 50 })
    const item = after.data.list.find((e: any) => e.id === 2)
    expect(item.status).toBe('shipped')
  })
  it('PUT complete changes status', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/points/exchanges/5/complete', {})
    expect(r.matched).toBe(true)
    const after = matchMock('GET', '/admin/api/points/exchanges', { page: 1, size: 50 })
    const item = after.data.list.find((e: any) => e.id === 5)
    expect(item.status).toBe('completed')
  })
})

describe('mock points summary', () => {
  it('GET /points/summary returns computed values', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/points/summary', {})
    expect(r.matched).toBe(true)
    expect(r.data.total_issued).toBeGreaterThan(0)
    expect(r.data.product_count).toBeGreaterThanOrEqual(0)
  })
})
