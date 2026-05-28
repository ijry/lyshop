import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock /categories', () => {
  it('GET /categories/tree returns 3 levels', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/categories/tree', {})
    expect(r.matched).toBe(true)
    expect(r.data.length).toBeGreaterThan(0)
    expect(r.data[0].children).toBeTruthy()
  })
  it('POST /categories creates', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/categories', { name: '新分类', parent_id: 0 })
    expect(r.matched).toBe(true)
    expect(r.data.id).toBeGreaterThan(0)
  })
})
