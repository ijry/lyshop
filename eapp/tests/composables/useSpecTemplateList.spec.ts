import { describe, expect, it, vi } from 'vitest'

vi.mock('@/api/spec-template', () => ({
  getSpecTemplates: vi.fn(async () => ({
    list: [
      { id: 1, name: '服装通用', category_ids: [21, 22], attrs: [{ name: '颜色', values: ['黑色', '白色'] }], status: 1 },
    ],
    total: 1, page: 1, size: 20,
  })),
}))

import { useSpecTemplateList } from '@/composables/useSpecTemplateList'

describe('useSpecTemplateList', () => {
  it('loads list', async () => {
    const h = useSpecTemplateList()
    await h.load()
    expect(h.list.value).toHaveLength(1)
    expect(h.loading.value).toBe(false)
  })
  it('refresh resets refreshing', async () => {
    const h = useSpecTemplateList()
    await h.refresh()
    expect(h.refreshing.value).toBe(false)
  })
})
