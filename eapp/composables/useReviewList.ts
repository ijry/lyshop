import { ref } from 'vue'
import { getReviews } from '@/api/review'

export type ReviewFilter = {
  keyword?: string
  reply_status?: 'all' | 'pending' | 'replied'
}

export function useReviewList(initial: ReviewFilter = {}) {
  const filter = ref<ReviewFilter>({ ...initial })
  const page = ref(1)
  const size = ref(20)
  const total = ref(0)
  const list = ref<any[]>([])
  const loading = ref(false)
  const refreshing = ref(false)

  function applyFilter(patch: Partial<ReviewFilter>) {
    filter.value = { ...filter.value, ...patch }
    page.value = 1
  }

  async function load() {
    loading.value = true
    try {
      const params: any = { ...filter.value, page: page.value, size: size.value }
      if (params.reply_status === 'all') delete params.reply_status
      const res: any = await getReviews(params)
      list.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
      total.value = Number(res?.total || list.value.length)
    } catch {
      list.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }

  async function refresh() {
    refreshing.value = true
    page.value = 1
    try {
      await load()
    } finally {
      refreshing.value = false
    }
  }

  async function loadMore() {
    if (loading.value || list.value.length >= total.value) return
    page.value += 1
    loading.value = true
    try {
      const params: any = { ...filter.value, page: page.value, size: size.value }
      if (params.reply_status === 'all') delete params.reply_status
      const res: any = await getReviews(params)
      list.value = list.value.concat(Array.isArray(res?.list) ? res.list : [])
    } finally {
      loading.value = false
    }
  }

  return {
    filter, page, size, total, list, loading, refreshing,
    applyFilter, load, refresh, loadMore,
  }
}
