import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock /products extended', () => {
  it('filters by status and low_stock', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/products', { page: 1, size: 100, low_stock: true })
    expect(r.matched).toBe(true)
  })
  it('filters by category_id', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/products', { page: 1, size: 100, category_id: 1 })
    expect(r.matched).toBe(true)
  })
  it('sorts by sales', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/products', { page: 1, size: 100, sort_by: 'sales' })
    expect(r.matched).toBe(true)
  })
})

describe('mock /products/batch/*', () => {
  it('PUT /products/batch/status', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/products/batch/status', { ids: [1, 2], status: 0 })
    expect(r.matched).toBe(true)
    expect(r.data.success_ids.length).toBeGreaterThan(0)
  })
  it('PUT /products/batch/price percent', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/products/batch/price', { ids: [1, 2], adjustment: { type: 'percent', value: -10 } })
    expect(r.matched).toBe(true)
  })
})
