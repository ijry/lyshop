import { computed, ref } from 'vue'
import { getOrders } from '@/api/order'
import { useBatchSelection } from './useBatchSelection'

export type OrderFilter = {
  status?: string; keyword?: string; time_start?: string; time_end?: string
  amount_min?: number; amount_max?: number; logistics_company?: string
  province?: string; pay_method?: string; has_after_sale?: boolean
}

export function useOrderList(initial: OrderFilter = {}) {
  const filter = ref<OrderFilter>({ ...initial })
  const page = ref(1)
  const size = ref(20)
  const total = ref(0)
  const list = ref<any[]>([])
  const loading = ref(false)
  const refreshing = ref(false)
  const sel = useBatchSelection<number>()

  const selectedIds = computed(() => Array.from(sel.selected.value))

  function applyFilter(patch: Partial<OrderFilter>) { filter.value = { ...filter.value, ...patch }; page.value = 1 }
  function resetFilter() { filter.value = { ...initial }; page.value = 1 }

  async function load() {
    loading.value = true
    try {
      const res: any = await getOrders({ ...filter.value, page: page.value, size: size.value })
      list.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
      total.value = Number(res?.total || list.value.length)
    } catch { list.value = []; total.value = 0 } finally { loading.value = false }
  }
  async function refresh() {
    refreshing.value = true; page.value = 1; sel.clear()
    try { await load() } finally { refreshing.value = false }
  }
  async function loadMore() {
    if (loading.value || list.value.length >= total.value) return
    page.value += 1; loading.value = true
    try {
      const res: any = await getOrders({ ...filter.value, page: page.value, size: size.value })
      list.value = list.value.concat(Array.isArray(res?.list) ? res.list : [])
    } finally { loading.value = false }
  }

  return {
    filter, page, size, total, list, loading, refreshing,
    selectedIds, selectCount: sel.count,
    toggleSelect: sel.toggle, isSelected: sel.isSelected, selectAll: sel.selectAll, clearSelect: sel.clear,
    applyFilter, resetFilter, load, refresh, loadMore,
  }
}
