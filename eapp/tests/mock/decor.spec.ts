import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock decor variant routes', () => {
  it('GET /decor/index/variants returns variants', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/decor/index/variants', {})
    expect(r.matched).toBe(true)
    expect(Array.isArray(r.data)).toBe(true)
    expect(r.data.length).toBeGreaterThanOrEqual(1)
  })
  it('GET /decor/index returns default variant', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/decor/index?variant=default', {})
    expect(r.matched).toBe(true)
    expect(r.data).toBeTruthy()
    expect(r.data.variant_key).toBe('default')
  })
  it('PUT /decor/index updates components', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/decor/index?variant=default', { components: [{ type: 'banner', id: 'c1', props: {} }] })
    expect(r.matched).toBe(true)
  })
  it('POST /decor/index/publish publishes a variant', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/decor/index/publish?variant=default', {})
    expect(r.matched).toBe(true)
  })
  it('POST /decor/index/copies creates a copy', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/decor/index/copies', { from_variant_key: 'default', new_variant_key: 'test_copy', new_variant_name: '测试副本' })
    expect(r.matched).toBe(true)
    const list = matchMock('GET', '/admin/api/decor/index/variants', {})
    expect(list.data.find((v: any) => v.variant_key === 'test_copy')).toBeTruthy()
  })
})
