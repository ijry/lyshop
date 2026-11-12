import { describe, expect, it } from 'vitest'

// admin mock 的订单状态模块会在加载时读取 uni 存储，node 测试环境需先打桩
;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
}

const ACTIVITIES = '/admin/api/seckill/activities'
const PRODUCTS = '/admin/api/seckill/products'
const hoursFromNow = (hours: number) => new Date(Date.now() + hours * 3600_000).toISOString()

async function match(method: string, url: string, params?: any) {
  const { matchMock } = await import('../../../admin/src/mock/index')
  return matchMock(method, url, params)
}

function listOf(result: any): any[] {
  return Array.isArray(result?.data?.list) ? result.data.list : []
}

describe('mock marketing activities', () => {
  it('seeds activities covering running / upcoming / ended / disabled phases', async () => {
    const list = listOf(await match('GET', ACTIVITIES, { page: 1, size: 200 }))
    expect(list.length).toBeGreaterThanOrEqual(4)
    const now = Date.now()
    const phases = new Set<string>()
    for (const item of list) {
      if (Number(item.status) !== 1) phases.add('disabled')
      else if (new Date(item.start_at).getTime() > now) phases.add('upcoming')
      else if (new Date(item.end_at).getTime() < now) phases.add('ended')
      else phases.add('running')
    }
    expect([...phases].sort()).toEqual(['disabled', 'ended', 'running', 'upcoming'])
  })

  it('returns a business error for an invalid period instead of throwing', async () => {
    const result = await match('POST', ACTIVITIES, {
      name: '无效时间',
      start_at: hoursFromNow(5),
      end_at: hoursFromNow(1),
      status: 1,
    })
    expect(result.matched).toBe(true)
    expect(result.error).toBe('请填写有效的活动时间区间')
  })

  it('returns a conflict error for overlapping periods', async () => {
    const list = listOf(await match('GET', ACTIVITIES, { page: 1, size: 200 }))
    const running = list.find(
      (item) => new Date(item.start_at).getTime() < Date.now() && new Date(item.end_at).getTime() > Date.now(),
    )
    expect(running).toBeTruthy()
    const result = await match('POST', ACTIVITIES, {
      name: '重叠活动',
      start_at: running.start_at,
      end_at: running.end_at,
      status: 1,
    })
    expect(result.error).toBe('同类型活动时间段冲突')
  })

  it('creates an activity, enriches its products and cascades on delete', async () => {
    const created: any = await match('POST', ACTIVITIES, {
      name: '回归测试活动',
      start_at: hoursFromNow(300),
      end_at: hoursFromNow(304),
      status: 1,
    })
    expect(created.error).toBeUndefined()
    const id = Number(created.data?.id || 0)
    expect(id).toBeGreaterThan(100)

    // 前端只提交业务字段，mock 需要从商品主数据补齐展示字段
    const put = await match('PUT', `${ACTIVITIES}/${id}/products`, [
      { product_id: 1, sku_id: 10011, activity_price: 1999 },
    ])
    expect(put.matched).toBe(true)
    const rows = listOf(await match('GET', PRODUCTS, { activity_id: id, page: 1, size: 200 }))
    expect(rows).toHaveLength(1)
    expect(rows[0].product_title).toBe('旗舰智能手机 Pro Max')
    expect(rows[0].product_cover).toContain('https://')
    expect(rows[0].sku_price).toBeGreaterThan(0)
    expect(rows[0].activity_price).toBe(1999)

    const removed = await match('DELETE', `${ACTIVITIES}/${id}`, {})
    expect(removed.matched).toBe(true)
    const after = listOf(await match('GET', ACTIVITIES, { page: 1, size: 200 }))
    expect(after.find((item) => Number(item.id) === id)).toBeUndefined()
    expect(listOf(await match('GET', PRODUCTS, { activity_id: id, page: 1, size: 200 }))).toHaveLength(0)
  })

  it('keeps display fields and sold_qty when replacing existing activity products', async () => {
    const before = listOf(await match('GET', PRODUCTS, { activity_id: 1, page: 1, size: 200 }))[0]
    expect(before).toBeTruthy()
    await match('PUT', `${ACTIVITIES}/1/products`, [
      {
        product_id: before.product_id,
        sku_id: before.sku_id,
        activity_price: 999,
        limit_per_order: before.limit_per_order,
        total_stock_limit: before.total_stock_limit,
      },
    ])
    const after = listOf(await match('GET', PRODUCTS, { activity_id: 1, page: 1, size: 200 }))[0]
    expect(after.activity_price).toBe(999)
    expect(after.product_title).toBe(before.product_title)
    expect(after.product_cover).toBe(before.product_cover)
    expect(after.sku_price).toBe(before.sku_price)
    expect(after.sold_qty).toBe(before.sold_qty)
  })
})
