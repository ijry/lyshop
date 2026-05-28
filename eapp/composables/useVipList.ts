import { ref } from 'vue'
import { get, post, put, del } from '@/utils/request'

export type VipResource = 'plans' | 'levels' | 'coupon-rules' | 'sku-prices'

export function useVipList(resource: VipResource) {
  const list = ref<any[]>([])
  const loading = ref(false)
  const refreshing = ref(false)

  async function load(params?: any) {
    loading.value = true
    try {
      const res: any = await get(`/vip/${resource}`, params)
      list.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
    } catch {
      list.value = []
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
    await post(`/vip/${resource}`, payload)
  }

  async function update(id: number | string, payload: any) {
    await put(`/vip/${resource}/${id}`, payload)
  }

  async function remove(id: number | string) {
    await del(`/vip/${resource}/${id}`)
  }

  return { list, loading, refreshing, load, refresh, create, update, remove }
}
