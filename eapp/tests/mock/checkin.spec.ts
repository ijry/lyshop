import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock checkin rules', () => {
  it('GET /checkin/rules returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/checkin/rules', {})
    expect(r.matched).toBe(true)
    expect(Array.isArray(r.data)).toBe(true)
    expect(r.data.length).toBeGreaterThanOrEqual(3)
    expect(r.data[0]).toHaveProperty('day')
    expect(r.data[0]).toHaveProperty('points')
  })
  it('PUT /checkin/rules accepts update (static route returns matched)', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/checkin/rules', [{ day: 0, points: 15 }, { day: 5, points: 30 }])
    expect(r.matched).toBe(true)
  })
})

describe('mock checkin logs', () => {
  it('GET /checkin/logs returns paginated data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/checkin/logs', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
    expect(r.data.list[0]).toHaveProperty('user_id')
    expect(r.data.list[0]).toHaveProperty('consecutive_days')
  })
})
