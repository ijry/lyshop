import { describe, expect, it, vi } from 'vitest'

vi.mock('@/utils/request', () => ({
  get: vi.fn(async () => ({
    list: [
      { id: 1, title: '10元无门槛优惠券', type: 'coupon', points_price: 500, stock: 100, status: 1 },
      { id: 2, title: '定制帆布袋', type: 'physical', points_price: 2000, stock: 50, status: 1 },
    ],
    total: 2, page: 1, size: 20,
  })),
  post: vi.fn(async () => null),
  put: vi.fn(async () => null),
  del: vi.fn(async () => null),
}))

import { usePointsList } from '@/composables/usePointsList'

describe('usePointsList', () => {
  it('loads products list', async () => {
    const h = usePointsList('products')
    await h.load()
    expect(h.list.value).toHaveLength(2)
    expect(h.loading.value).toBe(false)
  })
  it('refresh resets refreshing state', async () => {
    const h = usePointsList('exchanges')
    await h.refresh()
    expect(h.refreshing.value).toBe(false)
    expect(h.list.value).toHaveLength(2)
  })
  it('create calls post', async () => {
    const { post } = await import('@/utils/request')
    const h = usePointsList('products')
    await h.create({ title: 'test' })
    expect(post).toHaveBeenCalled()
  })
  it('remove calls del', async () => {
    const { del } = await import('@/utils/request')
    const h = usePointsList('products')
    await h.remove(1)
    expect(del).toHaveBeenCalled()
  })
})
