import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock GET /reviews with reply_status filter', () => {
  it('returns all reviews without filter', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/reviews', {})
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThan(0)
    expect(r.data.total).toBeTypeOf('number')
  })

  it('filters pending reviews (no admin_reply)', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const all = matchMock('GET', '/admin/api/reviews', {})
    const pending = matchMock('GET', '/admin/api/reviews', { reply_status: 'pending' })
    expect(pending.matched).toBe(true)
    const pendingCount = all.data.list.filter((r: any) => !r.admin_reply).length
    expect(pending.data.list.length).toBe(pendingCount)
    for (const row of pending.data.list) {
      expect(row.admin_reply).toBeFalsy()
    }
  })

  it('filters replied reviews (has admin_reply)', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const all = matchMock('GET', '/admin/api/reviews', {})
    const replied = matchMock('GET', '/admin/api/reviews', { reply_status: 'replied' })
    expect(replied.matched).toBe(true)
    const repliedCount = all.data.list.filter((r: any) => !!r.admin_reply).length
    expect(replied.data.list.length).toBe(repliedCount)
    for (const row of replied.data.list) {
      expect(row.admin_reply).toBeTruthy()
    }
  })
})

describe('mock GET /reviews/:id', () => {
  it('returns review detail by id', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const all = matchMock('GET', '/admin/api/reviews', {})
    const first = all.data.list[0]
    if (!first) return
    const r = matchMock('GET', `/admin/api/reviews/${first.id}`, {})
    expect(r.matched).toBe(true)
    expect(r.data.id).toBe(first.id)
    expect(r.data.content).toBeTypeOf('string')
  })
})

describe('mock POST /reviews/:id/reply', () => {
  it('creates or updates reply', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const all = matchMock('GET', '/admin/api/reviews', {})
    const target = all.data.list[0]
    if (!target) return
    const r = matchMock('POST', `/admin/api/reviews/${target.id}/reply`, { content: '感谢评价' })
    expect(r.matched).toBe(true)
    const after = matchMock('GET', `/admin/api/reviews/${target.id}`, {})
    expect(after.data.admin_reply).toBeTruthy()
    expect(after.data.admin_reply.content).toBe('感谢评价')
  })
})
