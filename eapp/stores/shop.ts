import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCurrentShop } from '@/api/system'

export const useShopStore = defineStore('eapp-shop', () => {
  const id = ref(0)
  const name = ref('')
  const logo = ref('')

  async function load() {
    const data: any = await getCurrentShop()
    id.value = Number(data?.id || 0)
    name.value = String(data?.name || '')
    logo.value = String(data?.logo || '')
  }

  return { id, name, logo, load }
})
