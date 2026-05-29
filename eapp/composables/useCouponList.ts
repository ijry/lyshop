import { ref } from 'vue'
import { getCoupons } from '@/api/marketing'

export type CouponFilter = { keyword?: string; status?: number | ''; type?: number | '' }

export function useCouponList(initial: CouponFilter = {}) {
  const filter = ref<CouponFilter>({ ...initial })
  const list = ref<any[]>([])
  const loading = ref(false)
  const refreshing = ref(false)

  async function load() {
    loading.value = true
    try {
      const params: any = { page: 1, size: 200 }
      if (filter.value.keyword) params.keyword = filter.value.keyword
      if (filter.value.status !== undefined && filter.value.status !== '') params.status = filter.value.status
      if (filter.value.type !== undefined && filter.value.type !== '') params.type = filter.value.type
      const res: any = await getCoupons(params)
      list.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
    } catch {
      list.value = []
    } finally {
      loading.value = false
    }
  }

  async function refresh() {
    refreshing.value = true
    try { await load() } finally { refreshing.value = false }
  }

  function applyFilter(patch: Partial<CouponFilter>) {
    filter.value = { ...filter.value, ...patch }
  }

  return { filter, list, loading, refreshing, load, refresh, applyFilter }
}
