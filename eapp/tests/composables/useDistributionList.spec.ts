import { describe, expect, it, vi } from 'vitest'

vi.mock('@/utils/request', () => ({
  get: vi.fn(async () => ({
    list: [
      { id: 1, nickname: '张三', status: 1, level: 1, balance: 580.50 },
      { id: 2, nickname: '李四', status: 1, level: 2, balance: 200.00 },
    ],
    total: 2, page: 1, size: 20,
  })),
  post: vi.fn(async () => null),
  put: vi.fn(async () => null),
  del: vi.fn(async () => null),
}))

import { useDistributionList } from '@/composables/useDistributionList'

describe('useDistributionList', () => {
  it('loads distributors list', async () => {
    const h = useDistributionList('distributors')
    await h.load()
    expect(h.list.value).toHaveLength(2)
    expect(h.loading.value).toBe(false)
  })
  it('refresh resets refreshing state', async () => {
    const h = useDistributionList('commissions')
    await h.refresh()
    expect(h.refreshing.value).toBe(false)
    expect(h.list.value).toHaveLength(2)
  })
  it('load with params calls get with params', async () => {
    const { get } = await import('@/utils/request')
    const h = useDistributionList('distributors')
    await h.load({ keyword: '张三' })
    expect(get).toHaveBeenCalledWith('/distribution/distributors', { keyword: '张三' })
  })
})
