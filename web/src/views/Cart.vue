<template>
  <div class="max-w-5xl mx-auto px-6 py-8">
    <h1 class="text-xl font-bold text-gray-900 mb-6">购物车</h1>

    <div v-if="!cart.items.length" class="card p-16 text-center">
      <div class="i-carbon-shopping-cart text-6xl text-gray-200 mx-auto mb-4" />
      <p class="text-gray-400 mb-4">购物车还是空的</p>
      <router-link to="/products" class="btn-primary inline-block">去逛逛</router-link>
    </div>

    <div v-else class="flex flex-col lg:flex-row gap-6">
      <!-- Item list -->
      <div class="flex-1">
        <div class="card divide-y divide-gray-50">
          <div v-for="item in cart.items" :key="item.sku_id"
            class="flex items-center gap-5 p-5 hover:bg-gray-50/50 transition-colors">
            <button @click="toggleItem(item.sku_id)"
              class="w-5 h-5 rounded-full border flex-center shrink-0 transition-colors"
              :class="isChecked(item.sku_id) ? 'bg-red-500 border-red-500 text-white' : 'bg-white border-gray-300 text-transparent'">
              <div class="i-carbon-checkmark text-xs" />
            </button>
            <img :src="item.product.cover" class="w-20 h-20 rounded-xl object-cover shrink-0" />
            <div class="flex-1 min-w-0">
              <h3 class="text-sm font-medium text-gray-800 line-clamp-1">{{ item.product.title }}</h3>
              <p class="text-xs text-gray-400 mt-1">{{ parseAttrs(item.sku.attrs) }}</p>
              <div class="flex-between mt-3">
                <span class="text-base font-bold text-red-500">¥{{ item.sku.price }}</span>
                <div class="flex items-center gap-2">
                  <button @click="changeQty(item, -1)"
                    class="w-7 h-7 rounded-md border border-gray-200 flex-center text-gray-500 hover:bg-gray-100 transition text-xs">-</button>
                  <span class="w-8 text-center text-sm">{{ item.qty }}</span>
                  <button @click="changeQty(item, 1)"
                    class="w-7 h-7 rounded-md border border-gray-200 flex-center text-gray-500 hover:bg-gray-100 transition text-xs">+</button>
                  <button @click="removeItem(item.sku_id)"
                    class="ml-3 text-gray-300 hover:text-red-500 transition-colors">
                    <div class="i-carbon-trash-can text-lg" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Summary -->
      <div class="w-72 shrink-0">
        <div class="card p-5 sticky top-24">
          <h3 class="text-sm font-semibold text-gray-800 mb-4">订单摘要</h3>
          <div class="space-y-2 text-sm">
            <div class="flex-between">
              <button @click="toggleCheckAll" class="inline-flex items-center gap-2 text-gray-600 hover:text-gray-800">
                <span class="w-4 h-4 rounded-full border flex-center text-[10px] transition-colors"
                  :class="allChecked ? 'bg-red-500 border-red-500 text-white' : 'bg-white border-gray-300 text-transparent'">
                  <span class="i-carbon-checkmark" />
                </span>
                <span>全选</span>
              </button>
            </div>
            <div class="flex-between text-gray-500">
              <span>已选商品</span><span>{{ selectedCount }} 件</span>
            </div>
            <div class="flex-between text-gray-500">
              <span>运费</span><span class="text-green-600">免运费</span>
            </div>
          </div>
          <div class="border-t border-gray-100 mt-4 pt-4">
            <div class="flex-between">
              <span class="text-sm text-gray-600">合计</span>
              <span class="text-xl font-bold text-red-500">¥{{ selectedTotal.toFixed(2) }}</span>
            </div>
          </div>
          <button class="btn-primary w-full mt-4 !py-3 disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="selectedCount === 0" @click="checkout">
            去结算 ({{ selectedCount }})
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { get } from '@/api/request'
import { useCartStore } from '@/stores/cart'

const cart = useCartStore()
const router = useRouter()
const checkedSkuIds = ref<number[]>([])

const selectedItems = computed(() =>
  cart.items.filter(i => checkedSkuIds.value.includes(i.sku_id))
)
const selectedCount = computed(() => selectedItems.value.length)
const selectedTotal = computed(() =>
  selectedItems.value.reduce((s, i) => s + i.sku.price * i.qty, 0)
)
const allChecked = computed(() =>
  cart.items.length > 0 && checkedSkuIds.value.length === cart.items.length
)

function parseAttrs(attrs: string) {
  try { return JSON.parse(attrs).map((a: any) => a.value).join(' / ') }
  catch { return '默认规格' }
}

function changeQty(item: any, delta: number) {
  const newQty = item.qty + delta
  if (newQty < 1) return
  cart.updateQty(item.sku_id, newQty)
}

function isChecked(skuId: number) {
  return checkedSkuIds.value.includes(skuId)
}

function toggleItem(skuId: number) {
  if (isChecked(skuId)) {
    checkedSkuIds.value = checkedSkuIds.value.filter(id => id !== skuId)
    return
  }
  checkedSkuIds.value.push(skuId)
}

function toggleCheckAll() {
  if (allChecked.value) {
    checkedSkuIds.value = []
    return
  }
  checkedSkuIds.value = cart.items.map(i => i.sku_id)
}

function removeItem(skuId: number) {
  cart.removeItem(skuId)
  checkedSkuIds.value = checkedSkuIds.value.filter(id => id !== skuId)
}

function checkout() {
  if (!selectedCount.value) {
    window.alert('请先勾选商品')
    return
  }
  router.push('/orders')
}

function syncCheckedSkuIds() {
  const currentIds = cart.items.map(i => i.sku_id)
  const checkedSet = new Set(checkedSkuIds.value)
  const normalized = currentIds.filter(id => checkedSet.has(id))
  checkedSkuIds.value = normalized.length ? normalized : [...currentIds]
}

onMounted(async () => {
  const data = await get<any[]>('/api/v1/cart')
  if (data) cart.setItems(data)
  checkedSkuIds.value = cart.items.map(i => i.sku_id)
})

watch(
  () => cart.items.map(i => i.sku_id).join(','),
  () => syncCheckedSkuIds()
)
</script>
