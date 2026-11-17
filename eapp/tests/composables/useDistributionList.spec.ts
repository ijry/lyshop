import { describe, expect, it, vi } from 'vitest'

vi.mock('@/utils/request', () => ({
  get: vi.fn(async () => ({
    list: [
      { id: 1, real_name: '张三', status: 'active', level: 1, available_amount: 580.50, total_earnings: 1280.50 },
      { id: 2, real_name: '李四', status: 'pending', level: 2, available_amount: 200.00, total_earnings: 420.00 },
    ],
    total: 2, page: 1, size: 20,
  })),
  post: vi.fn(async () => null),
  put: vi.fn(async () => null),
  del: vi.fn(async () => null),
}))

import { useDistributionList } from '@/composables/useDistributionList'
import { updateDistributorStatus, type DistributorStatus } from '@/api/distribution'

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
    await h.load({ keyword: '张三', status: 'active' })
    expect(get).toHaveBeenCalledWith('/distribution/distributors', { keyword: '张三', status: 'active' })
  })
})

describe('分销商状态 API', () => {
  it.each<DistributorStatus>(['pending', 'active', 'disabled'])('通过分销商更新接口提交 %s 状态', async (status) => {
    const { put } = await import('@/utils/request')
    await updateDistributorStatus(1, status)
    expect(put).toHaveBeenLastCalledWith('/distribution/distributors/1', { status })
  })
})
