import { describe, expect, it, vi } from 'vitest'
import { useAnalytics } from '@/composables/useAnalytics'

describe('useAnalytics', () => {
  it('loads data with default range', async () => {
    const fetchFn = vi.fn(async (days: number) => ({ total: days * 10 }))
    const { range, loading, data, load } = useAnalytics(fetchFn)
    expect(range.value).toBe('7')
    expect(loading.value).toBe(false)
    expect(data.value).toBeNull()
    await load()
    expect(fetchFn).toHaveBeenCalledWith(7)
    expect(data.value).toEqual({ total: 70 })
    expect(loading.value).toBe(false)
  })

  it('setRange changes range and reloads', async () => {
    const fetchFn = vi.fn(async (days: number) => ({ total: days }))
    const { range, data, setRange } = useAnalytics(fetchFn)
    await setRange('30')
    expect(range.value).toBe('30')
    expect(fetchFn).toHaveBeenCalledWith(30)
    expect(data.value).toEqual({ total: 30 })
  })

  it('handles fetch error gracefully', async () => {
    const fetchFn = vi.fn(async () => { throw new Error('fail') })
    const { data, load } = useAnalytics(fetchFn)
    await load()
    expect(data.value).toBeNull()
  })

  it('loading state toggles correctly', async () => {
    let resolve: Function
    const fetchFn = vi.fn(() => new Promise<{ ok: boolean }>((r) => { resolve = r }))
    const { loading, load } = useAnalytics(fetchFn)
    const promise = load()
    expect(loading.value).toBe(true)
    resolve!({ ok: true })
    await promise
    expect(loading.value).toBe(false)
  })
})
