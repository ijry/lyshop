import { ref } from 'vue'
import { get, post, put, del } from '@/utils/request'

export type PointsResource = 'products' | 'exchanges'

export function usePointsList(resource: PointsResource) {
  const list = ref<any[]>([])
  const loading = ref(false)
  const refreshing = ref(false)

  async function load(params?: any) {
    loading.value = true
    try {
      const res: any = await get(`/points/${resource}`, params)
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

  async function create(payload: any) {
    await post(`/points/${resource}`, payload)
  }

  async function update(id: number | string, payload: any) {
    await put(`/points/${resource}/${id}`, payload)
  }

  async function remove(id: number | string) {
    await del(`/points/${resource}/${id}`)
  }

  return { list, loading, refreshing, load, refresh, create, update, remove }
}
