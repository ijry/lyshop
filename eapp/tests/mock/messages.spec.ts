import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock messages with priority', () => {
  it('GET /messages returns 8 messages with priority field', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/messages', { page: 1, size: 50 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBe(8)
    for (const msg of r.data.list) {
      expect(['normal', 'important', 'urgent']).toContain(msg.priority)
    }
  })

  it('filters by group', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/messages', { group: 'system', page: 1, size: 50 })
    expect(r.matched).toBe(true)
    for (const msg of r.data.list) {
      expect(msg.group).toBe('system')
    }
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
  })

  it('filters by priority', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/messages', { priority: 'urgent', page: 1, size: 50 })
    expect(r.matched).toBe(true)
    for (const msg of r.data.list) {
      expect(msg.priority).toBe('urgent')
    }
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
  })

  it('filters by both group and priority', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/messages', { group: 'order', priority: 'important', page: 1, size: 50 })
    expect(r.matched).toBe(true)
    for (const msg of r.data.list) {
      expect(msg.group).toBe('order')
      expect(msg.priority).toBe('important')
    }
  })

  it('covers all group types', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/messages', { page: 1, size: 50 })
    const groups = new Set(r.data.list.map((m: any) => m.group))
    expect(groups.has('system')).toBe(true)
    expect(groups.has('order')).toBe(true)
    expect(groups.has('marketing')).toBe(true)
  })

  it('covers all priority levels', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/messages', { page: 1, size: 50 })
    const priorities = new Set(r.data.list.map((m: any) => m.priority))
    expect(priorities.has('normal')).toBe(true)
    expect(priorities.has('important')).toBe(true)
    expect(priorities.has('urgent')).toBe(true)
  })
})
