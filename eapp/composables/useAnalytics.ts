import { ref, type Ref } from 'vue'

export type DateRange = '1' | '7' | '30' | '90'

export function useAnalytics<T>(fetchFn: (days: number) => Promise<T>) {
  const range = ref<DateRange>('7')
  const loading = ref(false)
  const data = ref<T | null>(null) as Ref<T | null>

  async function load() {
    loading.value = true
    try {
      data.value = await fetchFn(Number(range.value))
    } catch {
      data.value = null
    } finally {
      loading.value = false
    }
  }

  function setRange(r: DateRange) {
    range.value = r
    load()
  }

  return { range, loading, data, load, setRange }
}
