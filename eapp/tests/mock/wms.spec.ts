import { describe, expect, it } from 'vitest'

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

describe('mock WMS warehouses', () => {
  it('GET /wms/warehouses returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/wms/warehouses', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
  })
  it('POST /wms/warehouses creates a new warehouse', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/wms/warehouses', { name: '测试仓库', code: 'WH-TEST', address: '测试地址', contact: '张三', phone: '13800000000', status: 1 })
    expect(r.matched).toBe(true)
  })
})

describe('mock WMS stocks', () => {
  it('GET /wms/stocks returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/wms/stocks', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
  })
})

describe('mock WMS docs', () => {
  it('GET /wms/docs returns seed data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/wms/docs', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
  })
  it('POST /wms/docs creates a new doc', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/wms/docs', { doc_type: 'inbound', warehouse_id: 1, items: [{ sku_id: 101, sku_name: '测试商品', qty: 10, unit_cost: 50 }] })
    expect(r.matched).toBe(true)
  })
  it('POST /wms/docs/:id/complete completes a doc', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('POST', '/admin/api/wms/docs/1/complete', {})
    expect(r.matched).toBe(true)
  })
})

describe('mock WMS movements', () => {
  it('GET /wms/movements returns data', async () => {
    const { matchMock } = await import('../../../admin/src/mock/index')
    const r = matchMock('GET', '/admin/api/wms/movements', { page: 1, size: 20 })
    expect(r.matched).toBe(true)
    expect(r.data.list.length).toBeGreaterThanOrEqual(1)
  })
})
