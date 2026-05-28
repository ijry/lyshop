import { ref, computed } from 'vue'

/**
 * Generic batch selection composable.
 * @param options.max - optional maximum selection count
 */
export function useBatchSelection<T = string>(options?: { max?: number }) {
  const selected = ref<Set<T>>(new Set())

  const count = computed(() => selected.value.size)

  function toggle(id: T) {
    const s = new Set(selected.value)
    if (s.has(id)) {
      s.delete(id)
    } else {
      if (options?.max && s.size >= options.max) return
      s.add(id)
    }
    selected.value = s
  }

  function selectAll(ids: T[]) {
    const s = new Set<T>()
    const limit = options?.max ?? Infinity
    for (const id of ids) {
      if (s.size >= limit) break
      s.add(id)
    }
    selected.value = s
  }

  function clear() {
    selected.value = new Set()
  }

  function isSelected(id: T) {
    return selected.value.has(id)
  }

  return { selected, count, toggle, selectAll, clear, isSelected }
}
