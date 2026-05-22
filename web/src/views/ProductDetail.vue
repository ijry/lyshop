<template>
  <div class="max-w-7xl mx-auto px-6 py-8">
    <div v-if="product.id" class="flex flex-col lg:flex-row gap-10">
      <!-- Left: Images -->
      <div class="lg:w-[480px] shrink-0">
        <div class="aspect-square rounded-2xl overflow-hidden bg-gray-50 mb-3">
          <img :src="mainImage" class="w-full h-full object-cover" />
        </div>
        <div class="flex gap-2 overflow-x-auto">
          <div v-for="(img, i) in images" :key="i"
            @click="mainImage = img"
            :class="mainImage === img ? 'ring-2 ring-red-600' : 'ring-1 ring-gray-200'"
            class="w-16 h-16 rounded-lg overflow-hidden shrink-0 cursor-pointer hover:ring-red-300 transition-all">
            <img :src="img" class="w-full h-full object-cover" />
          </div>
        </div>
      </div>

      <!-- Right: Info -->
      <div class="flex-1 min-w-0">
        <h1 class="text-2xl font-bold text-gray-900 leading-8 mb-2">{{ product.title }}</h1>
        <p class="text-sm text-gray-500 mb-6">{{ product.subtitle }}</p>

        <!-- Price -->
        <div class="bg-gradient-to-r from-red-50 to-orange-50 rounded-xl px-6 py-4 mb-6">
          <div class="flex items-baseline gap-3">
            <span class="text-3xl font-bold text-red-500">¥{{ selectedSku?.price ?? product.price }}</span>
            <span v-if="product.origin_price" class="text-base text-gray-400 line-through">¥{{ product.origin_price }}</span>
            <span v-if="product.origin_price"
              class="bg-red-500 text-white text-xs px-2 py-0.5 rounded-md font-medium">
              省¥{{ (product.origin_price - (selectedSku?.price ?? product.price)).toFixed(0) }}
            </span>
          </div>
          <div class="flex gap-4 mt-2 text-xs text-gray-500">
            <span>库存 {{ selectedSku?.stock ?? product.stock }}</span>
            <span>已售 {{ product.sales }}</span>
          </div>
        </div>

        <!-- SKU -->
        <div v-if="skus.length" class="mb-6">
          <h3 class="text-sm font-semibold text-gray-700 mb-3">规格选择</h3>
          <div class="flex flex-wrap gap-2">
            <button v-for="sku in skus" :key="sku.id" @click="selectSku(sku)"
              :class="selectedSku?.id === sku.id
                ? 'border-red-600 bg-red-50 text-red-600'
                : 'border-gray-200 text-gray-600 hover:border-gray-300'"
              class="px-4 py-2 border rounded-lg text-sm transition-colors">
              {{ parseAttrs(sku.attrs) }}
            </button>
          </div>
        </div>

        <!-- Quantity -->
        <div class="mb-8">
          <h3 class="text-sm font-semibold text-gray-700 mb-3">购买数量</h3>
          <div class="flex items-center gap-3">
            <button @click="qty > 1 && qty--"
              class="w-9 h-9 rounded-lg border border-gray-200 flex-center text-gray-600 hover:bg-gray-50 transition-colors">-</button>
            <span class="w-12 text-center text-sm font-medium">{{ qty }}</span>
            <button @click="qty++"
              class="w-9 h-9 rounded-lg border border-gray-200 flex-center text-gray-600 hover:bg-gray-50 transition-colors">+</button>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex gap-3">
          <button @click="addToCart" class="btn-outline flex-1 !py-3 flex-center gap-2">
            <div class="i-carbon-shopping-cart" />
            加入购物车
          </button>
          <button @click="buyNow" class="btn-primary flex-1 !py-3">
            立即购买
          </button>
        </div>

        <!-- Detail -->
        <div class="mt-10 pt-8 border-t border-gray-100">
          <h3 class="text-base font-semibold text-gray-800 mb-4">商品详情</h3>
          <div class="prose prose-sm max-w-none text-gray-600" v-html="product.detail" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get, post } from '@/api/request'

const route = useRoute()
const router = useRouter()
const product = ref<any>({})
const skus = ref<any[]>([])
const selectedSku = ref<any>(null)
const images = ref<string[]>([])
const mainImage = ref('')
const qty = ref(1)

function parseAttrs(attrs: string) {
  try { return JSON.parse(attrs).map((a: any) => a.value).join(' / ') }
  catch { return '默认' }
}

function selectSku(sku: any) {
  selectedSku.value = sku
}

async function addToCart() {
  if (!selectedSku.value) return
  await post('/api/v1/cart/add', { sku_id: selectedSku.value.id, qty: qty.value })
  alert('已加入购物车')
}

function buyNow() {
  if (!selectedSku.value) return
  router.push(`/cart`)
}

onMounted(async () => {
  const data = await get<any>(`/api/v1/products/${route.params.id}`)
  if (!data) return
  product.value = data
  skus.value = data.skus || []
  if (skus.value.length) selectedSku.value = skus.value[0]
  const imgs: string[] = []
  if (data.cover) imgs.push(data.cover)
  if (data.images) imgs.push(...data.images.map((i: any) => i.url))
  images.value = imgs
  mainImage.value = imgs[0] || ''
})
</script>
