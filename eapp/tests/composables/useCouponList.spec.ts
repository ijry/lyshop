import { describe, expect, it, vi } from 'vitest'

vi.mock('@/api/marketing', () => ({
  getCoupons: vi.fn(async (params: any) => ({
    list: [
      { id: 1, name: '测试券', type: 1, discount: 10, status: 1, stack_rule: 'exclusive', target_type: 'all' },
    ],
    total: 1, page: params?.page || 1, size: params?.size || 20,
  })),
}))

import { useCouponList } from '@/composables/useCouponList'

describe('useCouponList', () => {
  it('loads list', async () => {
    const h = useCouponList()
    await h.load()
    expect(h.list.value).toHaveLength(1)
    expect(h.loading.value).toBe(false)
  })
  it('refresh resets refreshing', async () => {
    const h = useCouponList()
    await h.refresh()
    expect(h.refreshing.value).toBe(false)
  })
  it('applyFilter merges filter', () => {
    const h = useCouponList()
    h.applyFilter({ keyword: 'test' })
    expect(h.filter.value.keyword).toBe('test')
  })
})
