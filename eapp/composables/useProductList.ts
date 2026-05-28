import { computed, ref } from 'vue'
import { getProducts } from '@/api/product'
import { useBatchSelection } from './useBatchSelection'

export type ProductFilter = { keyword?: string; status?: 0|1|''; category_id?: number; sort_by?: string; low_stock?: boolean }

export function useProductList(initial: ProductFilter = {}) {
  const filter = ref<ProductFilter>({ ...initial })
  const page = ref(1); const size = ref(20); const total = ref(0)
  const list = ref<any[]>([])
  const loading = ref(false); const refreshing = ref(false)
  const sel = useBatchSelection<number>()

  const selectedIds = computed(() => Array.from(sel.selected.value))

  function applyFilter(patch: Partial<ProductFilter>) { filter.value = { ...filter.value, ...patch }; page.value = 1 }
  async function load() {
    loading.value = true
    try {
      const res: any = await getProducts({ ...filter.value, page: page.value, size: size.value })
      list.value = Array.isArray(res?.list) ? res.list : []
      total.value = Number(res?.total || list.value.length)
    } finally { loading.value = false }
  }
  async function refresh() { refreshing.value = true; page.value = 1; sel.clear(); try { await load() } finally { refreshing.value = false } }
  async function loadMore() {
    if (loading.value || list.value.length >= total.value) return
    page.value += 1; loading.value = true
    try {
      const res: any = await getProducts({ ...filter.value, page: page.value, size: size.value })
      list.value = list.value.concat(Array.isArray(res?.list) ? res.list : [])
    } finally { loading.value = false }
  }

  return {
    filter, page, total, list, loading, refreshing,
    selectedIds, selectCount: sel.count,
    toggleSelect: sel.toggle, isSelected: sel.isSelected, selectAll: sel.selectAll, clearSelect: sel.clear,
    applyFilter, load, refresh, loadMore,
  }
}
