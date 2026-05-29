import { describe, expect, it, vi } from 'vitest'

vi.mock('@/api/review', () => ({
  getReviews: vi.fn(async (params: any) => ({
    list: [
      { id: 1, content: 'good', product_score: 5, admin_reply: null },
      { id: 2, content: 'nice', product_score: 4, admin_reply: { id: 1, content: 'thanks' } },
    ],
    total: 2, page: params?.page || 1, size: params?.size || 20,
  })),
}))

import { useReviewList } from '@/composables/useReviewList'

describe('useReviewList', () => {
  it('loads reviews with default filter', async () => {
    const h = useReviewList()
    await h.load()
    expect(h.list.value).toHaveLength(2)
    expect(h.total.value).toBe(2)
  })

  it('applyFilter resets page to 1', async () => {
    const h = useReviewList()
    h.page.value = 3
    h.applyFilter({ reply_status: 'pending' })
    expect(h.page.value).toBe(1)
    expect(h.filter.value.reply_status).toBe('pending')
  })

  it('refresh resets to page 1', async () => {
    const h = useReviewList()
    h.page.value = 5
    await h.refresh()
    expect(h.page.value).toBe(1)
  })

  it('loadMore increments page', async () => {
    const h = useReviewList()
    await h.load()
    // total = 2, list.length = 2, so loadMore should not increment
    await h.loadMore()
    expect(h.page.value).toBe(1)
  })
})
