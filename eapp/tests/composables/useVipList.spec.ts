import { describe, expect, it, vi } from 'vitest'

vi.mock('@/utils/request', () => ({
  get: vi.fn(async () => ({
    list: [
      { id: 1, name: '月卡', months: 1, price: 19.9, status: 1 },
      { id: 2, name: '季卡', months: 3, price: 49.9, status: 1 },
    ],
    total: 2, page: 1, size: 20,
  })),
  post: vi.fn(async () => null),
  put: vi.fn(async () => null),
  del: vi.fn(async () => null),
}))

import { useVipList } from '@/composables/useVipList'

describe('useVipList', () => {
  it('loads list for a given resource', async () => {
    const h = useVipList('plans')
    await h.load()
    expect(h.list.value).toHaveLength(2)
    expect(h.loading.value).toBe(false)
  })
  it('refresh resets loading state', async () => {
    const h = useVipList('levels')
    await h.refresh()
    expect(h.refreshing.value).toBe(false)
    expect(h.list.value).toHaveLength(2)
  })
  it('create calls post', async () => {
    const { post } = await import('@/utils/request')
    const h = useVipList('plans')
    await h.create({ name: 'test' })
    expect(post).toHaveBeenCalled()
  })
  it('remove calls del', async () => {
    const { del } = await import('@/utils/request')
    const h = useVipList('coupon-rules')
    await h.remove(1)
    expect(del).toHaveBeenCalled()
  })
})
