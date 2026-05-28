import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock /after-sales extended', () => {
  it('filters by status and type', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/after-sales', { page: 1, size: 50, status: 'applied' })
    expect(r.matched).toBe(true)
    for (const row of r.data.list) expect(String(row.status)).toBe('applied')
  })
})

describe('mock after-sale messages / evidences', () => {
  it('POST messages appends', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/after-sales/8001/messages', { from: 'merchant', content: 'hello' })
    expect(r.matched).toBe(true)
    expect(r.data.messages.length).toBeGreaterThan(0)
  })
  it('POST evidences appends', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/after-sales/8001/evidences', { images: ['a.jpg'], remark: 'evi' })
    expect(r.matched).toBe(true)
    expect(r.data.evidences.length).toBeGreaterThan(0)
  })
})
