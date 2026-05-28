import { describe, expect, it, vi } from 'vitest'
vi.mock('@/api/product', () => ({
  getProducts: vi.fn(async (params: any) => ({ list: [{ id: 1, title: 'A', status: params.status ?? 1, stock: 5 }], total: 1, page: 1, size: 20 })),
}))
import { useProductList } from '@/composables/useProductList'
describe('useProductList', () => {
  it('loads with filter', async () => {
    const h = useProductList()
    h.applyFilter({ status: 1 }); await h.load()
    expect(h.list.value[0].id).toBe(1)
  })
})
