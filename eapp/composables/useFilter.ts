import { ref, type Ref } from 'vue'

/**
 * Persists filter state to uni storage by key.
 * Supports apply(patch) to merge values and reset() to clear.
 */
export function useFilter<T extends Record<string, any>>(
  storageKey: string,
  defaults: T
): {
  state: Ref<T>
  apply: (patch: Partial<T>) => void
  reset: () => void
} {
  const state = ref<T>(restore()) as Ref<T>

  function persist(val: T) {
    uni.setStorageSync(storageKey, JSON.stringify(val))
  }

  function restore(): T {
    try {
      const raw = uni.getStorageSync(storageKey)
      if (raw) return { ...defaults, ...JSON.parse(raw) }
    } catch {}
    return { ...defaults }
  }

  function apply(patch: Partial<T>) {
    const next = { ...state.value, ...patch }
    state.value = next
    persist(next)
  }

  function reset() {
    state.value = { ...defaults }
    uni.removeStorageSync(storageKey)
  }

  return { state, apply, reset }
}
