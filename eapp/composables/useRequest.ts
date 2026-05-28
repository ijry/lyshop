import { ref, type Ref } from 'vue'

/**
 * Thin wrapper around an async function.
 * Provides loading / error / data / run.
 */
export function useRequest<T>(fn: (...args: any[]) => Promise<T>) {
  const loading = ref(false)
  const error = ref<string | null>(null)
  const data = ref<T | null>(null) as Ref<T | null>

  async function run(...args: any[]) {
    loading.value = true
    error.value = null
    try {
      data.value = await fn(...args)
    } catch (e: any) {
      error.value = e?.message || String(e)
    } finally {
      loading.value = false
    }
  }

  return { loading, error, data, run }
}
