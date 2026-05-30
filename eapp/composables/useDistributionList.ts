import { ref } from 'vue'
import { get } from '@/utils/request'

export type DistributionResource = 'distributors' | 'commissions'

export function useDistributionList(resource: DistributionResource) {
  const list = ref<any[]>([])
  const loading = ref(false)
  const refreshing = ref(false)

  async function load(params?: any) {
    loading.value = true
    try {
      const res: any = await get(`/distribution/${resource}`, params)
      const newList = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
      // If page > 1, append to existing list; otherwise replace
      if (params?.page && params.page > 1) {
        list.value = list.value.concat(newList)
      } else {
        list.value = newList
      }
    } catch {
      if (!params?.page || params.page === 1) {
        list.value = []
      }
    } finally {
      loading.value = false
    }
  }

  async function refresh(params?: any) {
    refreshing.value = true
    try {
      await load(params)
    } finally {
      refreshing.value = false
    }
  }

  return { list, loading, refreshing, load, refresh }
}
