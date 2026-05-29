import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock spec-templates CRUD', () => {
  it('GET /spec-templates returns 3 seeds', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/spec-templates', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(3)
    expect(r.data.list[0]).toHaveProperty('attrs')
  })
  it('filters by keyword', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/spec-templates', { keyword: '服装', page: 1, size: 20 })
    expect(r.matched).toBe(true)
    for (const t of r.data.list) expect(t.name.toLowerCase()).toContain('服装')
  })
  it('filters by category_id', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/spec-templates', { category_id: 23, page: 1, size: 20 })
    expect(r.matched).toBe(true)
    for (const t of r.data.list) expect(t.category_ids).toContain(23)
  })
  it('GET /spec-templates/:id returns detail', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/spec-templates/1', {})
    expect(r.matched).toBe(true)
    expect(r.data.name).toBe('服装通用')
  })
  it('POST creates a template', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/spec-templates', { name: '测试模板', category_ids: [1], attrs: [{ name: '材质', values: ['棉', '麻'] }], status: 1 })
    expect(r.matched).toBe(true)
    expect(r.data.id).toBeGreaterThan(0)
  })
  it('PUT updates a template', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/spec-templates/2', { name: '数码配件Pro' })
    expect(r.matched).toBe(true)
  })
  it('DELETE removes a template', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const before = matchMock('GET', '/admin/api/spec-templates', { page: 1, size: 50 })
    const countBefore = before.data.list.length
    matchMock('DELETE', '/admin/api/spec-templates/2', {})
    const after = matchMock('GET', '/admin/api/spec-templates', { page: 1, size: 50 })
    expect(after.data.list.length).toBeLessThan(countBefore)
  })
})
