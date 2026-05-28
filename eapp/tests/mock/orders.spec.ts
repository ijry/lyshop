import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock /orders extended query', () => {
  it('filters by status', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/orders', { page: 1, size: 50, status: '2' })
    expect(r.matched).toBe(true)
    for (const o of r.data.list) expect(String(o.status)).toBe('2')
  })
  it('filters by amount range', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/orders', { page: 1, size: 50, amount_min: 100, amount_max: 10000 })
    expect(r.matched).toBe(true)
    for (const o of r.data.list) {
      const amt = Number(o.pay_amount || o.total_amount || 0)
      expect(amt).toBeGreaterThanOrEqual(100); expect(amt).toBeLessThanOrEqual(10000)
    }
  })
})

describe('mock order action routes', () => {
  it('POST /orders/{id}/repricing returns updated breakdown', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/orders/1/repricing', { items: [{ item_id: 1, price: 100 }], remark: 't' })
    expect(r.matched).toBe(true); expect(r.data.amount_breakdown).toBeTruthy()
  })
  it('POST /orders/{id}/notes pushes a note', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/orders/1/notes', { content: 'hello' })
    expect(r.matched).toBe(true); expect(r.data.notes.length).toBeGreaterThan(0)
  })
  it('POST /orders/{id}/remind-pay returns sent_at', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/orders/1/remind-pay', { channel: 'sms' })
    expect(r.matched).toBe(true); expect(r.data.sent_at).toBeTruthy()
  })
  it('GET /orders/{id}/print-template returns html', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/orders/1/print-template', {})
    expect(r.matched).toBe(true); expect(typeof r.data.template).toBe('string')
  })
  it('GET /orders/{id}/timeline returns stages', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/orders/1/timeline', {})
    expect(r.matched).toBe(true); expect(r.data.length).toBeGreaterThanOrEqual(3)
  })
})

describe('mock batch order routes', () => {
  it('POST /orders/batch/ship returns success and fail', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/orders/batch/ship', [
      { order_id: 1, company: 'SF', tracking_no: 'SF1' },
      { order_id: 99999, company: 'SF', tracking_no: 'SF3' },
    ])
    expect(r.matched).toBe(true)
    expect(r.data.success_ids.length).toBeGreaterThan(0)
    expect(r.data.fail.length).toBeGreaterThan(0)
  })
  it('POST /orders/batch/notes works', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/orders/batch/notes', { ids: [1, 2], content: '批量备注' })
    expect(r.matched).toBe(true)
  })
  it('POST /orders/batch/close honors reason', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/orders/batch/close', { ids: [1, 2], reason: '不要了' })
    expect(r.matched).toBe(true)
  })
})
