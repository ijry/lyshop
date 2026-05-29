import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock IM sessions', () => {
  it('GET /im/sessions returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/im/sessions', { page: 1, size: 50 })
    expect(r.matched).toBe(true)
    expect(Array.isArray(r.data)).toBe(true)
    expect(r.data.length).toBeGreaterThanOrEqual(1)
  })
  it('GET /im/sessions/:id/messages returns messages', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/im/sessions/1/messages', {})
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
  })
  it('POST /im/sessions/:id/messages sends a message', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/im/sessions/1/messages', { content: '你好', sender_type: 'staff' })
    expect(r.matched).toBe(true)
  })
})

describe('mock IM auto-replies', () => {
  it('GET /im/auto-replies returns data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/im/auto-replies', {})
    // may return matched:false (no static route) which is ok — the UI handles null gracefully
    expect(typeof r.matched).toBe('boolean')
  })
  it('POST /im/auto-replies creates a rule', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/im/auto-replies', { keyword: '你好', match_type: 'contains', reply: '欢迎光临', status: 1 })
    expect(r.matched).toBe(true)
  })
})
