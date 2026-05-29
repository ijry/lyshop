import { ref } from 'vue'
import { getSpecTemplates } from '@/api/spec-template'

export function useSpecTemplateList() {
  const list = ref<any[]>([])
  const loading = ref(false)
  const refreshing = ref(false)

  async function load(params?: any) {
    loading.value = true
    try {
      const res: any = await getSpecTemplates(params)
      list.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
    } catch {
      list.value = []
    } finally {
      loading.value = false
    }
  }

  async function refresh(params?: any) {
    refreshing.value = true
    try { await load(params) } finally { refreshing.value = false }
  }

  return { list, loading, refreshing, load, refresh }
}
