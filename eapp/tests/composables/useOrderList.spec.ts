import { describe, expect, it, vi } from 'vitest'

vi.mock('@/api/order', () => ({
  getOrders: vi.fn(async (params: any) => ({
    list: [
      { id: 1, status: params.status || '2', pay_amount: 199 },
      { id: 2, status: params.status || '2', pay_amount: 99 },
    ],
    total: 2, page: params.page || 1, size: params.size || 20,
  })),
}))

import { useOrderList } from '@/composables/useOrderList'

describe('useOrderList', () => {
  it('loads first page with filter', async () => {
    const h = useOrderList()
    h.applyFilter({ status: '2' })
    await h.load()
    expect(h.list.value).toHaveLength(2)
    expect(h.total.value).toBe(2)
  })
  it('refresh resets to page 1', async () => {
    const h = useOrderList(); h.page.value = 5
    await h.refresh()
    expect(h.page.value).toBe(1)
  })
  it('selection helpers', async () => {
    const h = useOrderList(); await h.load()
    h.toggleSelect(1); expect(h.selectedIds.value).toEqual([1])
    h.clearSelect(); expect(h.selectedIds.value).toEqual([])
  })
})
