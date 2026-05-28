import { describe, expect, it, beforeEach } from 'vitest'
import { useFilter } from '@/composables/useFilter'

const storage = new Map<string, string>()
;(globalThis as any).uni = {
  getStorageSync: (k: string) => storage.get(k) || '',
  setStorageSync: (k: string, v: string) => { storage.set(k, String(v)) },
  removeStorageSync: (k: string) => { storage.delete(k) },
}

describe('useFilter', () => {
  beforeEach(() => {
    storage.clear()
  })

  it('returns default value when no stored data', () => {
    const { state } = useFilter('test_key', { status: 'all', page: 1 })
    expect(state.value).toEqual({ status: 'all', page: 1 })
  })

  it('apply persists patch to storage', () => {
    const { state, apply } = useFilter('test_key', { status: 'all', page: 1 })
    apply({ status: 'active' })
    expect(state.value.status).toBe('active')
    expect(state.value.page).toBe(1)
    expect(storage.has('test_key')).toBe(true)
  })

  it('reset clears storage and restores defaults', () => {
    const { state, apply, reset } = useFilter('test_key', { status: 'all', page: 1 })
    apply({ status: 'closed' })
    reset()
    expect(state.value).toEqual({ status: 'all', page: 1 })
    expect(storage.has('test_key')).toBe(false)
  })

  it('restores previously stored data', () => {
    storage.set('restore_key', JSON.stringify({ status: 'active', page: 3 }))
    const { state } = useFilter('restore_key', { status: 'all', page: 1 })
    expect(state.value.status).toBe('active')
    expect(state.value.page).toBe(3)
  })
})
