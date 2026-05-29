import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock distribution distributors', () => {
  it('GET /distribution/distributors returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/distribution/distributors', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(3)
  })
  it('filters by keyword', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/distribution/distributors', { keyword: '张三', page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
    expect(r.data.list[0].nickname).toBe('张三')
  })
  it('filters by status', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/distribution/distributors', { status: 0, page: 1, size: 20 })
    expect(r.matched).toBe(true)
    for (const item of r.data.list) expect(Number(item.status)).toBe(0)
  })
  it('GET /distribution/distributors/:id returns detail', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/distribution/distributors/1', {})
    expect(r.matched).toBe(true)
    expect(r.data.nickname).toBe('张三')
  })
  it('PUT /distribution/distributors/:id/status toggles status', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/distribution/distributors/1/status', { status: 0 })
    expect(r.matched).toBe(true)
    const detail = matchMock('GET', '/admin/api/distribution/distributors/1', {})
    expect(detail.data.status).toBe(0)
    // restore
    matchMock('PUT', '/admin/api/distribution/distributors/1/status', { status: 1 })
  })
})

describe('mock distribution commissions', () => {
  it('GET /distribution/commissions returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/distribution/commissions', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(5)
  })
  it('filters by status', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/distribution/commissions', { status: 2, page: 1, size: 20 })
    expect(r.matched).toBe(true)
    for (const item of r.data.list) expect(Number(item.status)).toBe(2)
  })
  it('filters by distributor_id', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/distribution/commissions', { distributor_id: 1, page: 1, size: 20 })
    expect(r.matched).toBe(true)
    for (const item of r.data.list) expect(Number(item.distributor_id)).toBe(1)
  })
  it('PUT settle changes status to 2', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/distribution/commissions/6/settle', {})
    expect(r.matched).toBe(true)
    const after = matchMock('GET', '/admin/api/distribution/commissions', { page: 1, size: 50 })
    const item = after.data.list.find((c: any) => c.id === 6)
    expect(item.status).toBe(2)
  })
  it('PUT return changes status to 3', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/distribution/commissions/7/return', {})
    expect(r.matched).toBe(true)
    const after = matchMock('GET', '/admin/api/distribution/commissions', { page: 1, size: 50 })
    const item = after.data.list.find((c: any) => c.id === 7)
    expect(item.status).toBe(3)
  })
})

describe('mock distribution config', () => {
  it('GET /distribution/config returns config', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/distribution/config', {})
    expect(r.matched).toBe(true)
    expect(r.data.level1_rate).toBe(0.10)
    expect(r.data.level2_rate).toBe(0.05)
  })
  it('PUT /distribution/config updates config', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('PUT', '/admin/api/distribution/config', { level1_rate: 0.15, level2_rate: 0.08 })
    expect(r.matched).toBe(true)
    const after = matchMock('GET', '/admin/api/distribution/config', {})
    expect(after.data.level1_rate).toBe(0.15)
    expect(after.data.level2_rate).toBe(0.08)
    // restore
    matchMock('PUT', '/admin/api/distribution/config', { level1_rate: 0.10, level2_rate: 0.05 })
  })
})
