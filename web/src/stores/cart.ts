import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface CartItem {
  sku_id: number
  qty: number
  product: { id: number; title: string; cover: string; price: number }
  sku: { id: number; attrs: string; price: number; stock: number }
}

export const useCartStore = defineStore('cart', () => {
  const items = ref<CartItem[]>([])

  const total = computed(() =>
    items.value.reduce((s, i) => s + i.sku.price * i.qty, 0)
  )
  const count = computed(() => items.value.length)

  function setItems(list: CartItem[]) { items.value = list }
  function removeItem(skuId: number) {
    items.value = items.value.filter(i => i.sku_id !== skuId)
  }
  function updateQty(skuId: number, qty: number) {
    const item = items.value.find(i => i.sku_id === skuId)
    if (item) item.qty = qty
  }

  return { items, total, count, setItems, removeItem, updateQty }
})
