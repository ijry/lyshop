import { describe, expect, it, vi } from 'vitest'

vi.mock('@/utils/request', () => ({
  get: vi.fn(async () => ({
    list: [
      { id: 1, name: '主仓库', code: 'WH-001', status: 1 },
      { id: 2, name: '华南仓', code: 'WH-002', status: 1 },
    ],
    total: 2, page: 1, size: 20,
  })),
  post: vi.fn(async () => null),
  put: vi.fn(async () => null),
  del: vi.fn(async () => null),
}))

import { useWmsList } from '@/composables/useWmsList'

describe('useWmsList', () => {
  it('loads list for warehouses', async () => {
    const h = useWmsList('warehouses')
    await h.load()
    expect(h.list.value).toHaveLength(2)
    expect(h.loading.value).toBe(false)
  })
  it('refresh resets refreshing state', async () => {
    const h = useWmsList('stocks')
    await h.refresh()
    expect(h.refreshing.value).toBe(false)
    expect(h.list.value).toHaveLength(2)
  })
  it('create calls post', async () => {
    const { post } = await import('@/utils/request')
    const h = useWmsList('warehouses')
    await h.create({ name: '新仓库' })
    expect(post).toHaveBeenCalled()
  })
  it('update calls put', async () => {
    const { put } = await import('@/utils/request')
    const h = useWmsList('docs')
    await h.update(1, { remark: '测试' })
    expect(put).toHaveBeenCalled()
  })
  it('remove calls del', async () => {
    const { del } = await import('@/utils/request')
    const h = useWmsList('movements')
    await h.remove(1)
    expect(del).toHaveBeenCalled()
  })
})
