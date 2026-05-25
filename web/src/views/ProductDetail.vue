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
            <span>收藏 {{ product.favorite_count || 0 }}</span>
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
          <button @click="openChat" class="btn-outline !py-3 flex-center gap-2 w-36">
            <div class="i-carbon-chat text-lg" />
            客服
          </button>
          <button @click="toggleFavorite" class="btn-outline !py-3 flex-center gap-2 w-36">
            <div :class="product.is_favorited ? 'i-carbon-favorite-filled text-red-500' : 'i-carbon-favorite'" />
            {{ product.is_favorited ? '已收藏' : '收藏' }}
          </button>
          <button @click="addToCart" class="btn-outline flex-1 !py-3 flex-center gap-2">
            <div class="i-carbon-shopping-cart" />
            加入购物车
          </button>
          <button @click="buyNow" class="btn-primary flex-1 !py-3">
            立即购买
          </button>
        </div>

        <div class="mt-10 pt-8 border-t border-gray-100">
          <div class="flex gap-2 bg-gray-100 rounded-xl p-1 mb-4 w-fit">
            <button
              v-for="tab in tabs"
              :key="tab.value"
              @click="setTab(tab.value)"
              :class="activeTab === tab.value ? 'bg-white shadow-sm text-gray-900' : 'text-gray-500 hover:text-gray-700'"
              class="px-4 py-2 rounded-lg text-sm font-medium transition-all"
            >
              {{ tab.label }}
            </button>
          </div>

          <div v-if="activeTab === 'detail'" class="space-y-4">
            <template v-for="block in detailBlocks" :key="block.id">
              <p v-if="block.type === 'text'" class="text-sm leading-7 text-gray-600 whitespace-pre-wrap">{{ block.props?.text || '' }}</p>
              <img v-else-if="block.type === 'image'" :src="block.props?.url" :alt="block.props?.alt || ''" class="w-full rounded-xl border border-gray-100" />
            </template>
            <p v-if="!detailBlocks.length" class="text-sm text-gray-400">暂无详情</p>
          </div>

          <div v-else class="space-y-4">
            <div class="bg-gray-50 rounded-xl p-4">
              <div class="flex-between text-sm text-gray-600 mb-1">
                <span>商品评分</span><span class="font-semibold text-red-500">{{ reviewSummary.avg_product_score.toFixed(1) }}</span>
              </div>
              <div class="flex-between text-sm text-gray-600">
                <span>物流评分</span><span class="font-semibold text-red-500">{{ reviewSummary.avg_logistics_score.toFixed(1) }}</span>
              </div>
            </div>

            <div v-if="reviews.length" class="space-y-3">
              <div v-for="rv in reviews" :key="rv.id" class="border border-gray-100 rounded-xl p-4">
                <div class="flex-between mb-2">
                  <p class="text-sm font-medium text-gray-800">{{ rv.user_nickname || '匿名用户' }}</p>
                  <p class="text-xs text-gray-400">{{ formatDate(rv.created_at) }}</p>
                </div>
                <p class="text-xs text-gray-500 mb-2">商品 {{ rv.product_score }} / 物流 {{ rv.logistics_score }}</p>
                <p class="text-sm text-gray-600 leading-6">{{ rv.content || '用户未填写评价' }}</p>
                <div v-if="rv.images?.length" class="flex flex-wrap gap-2 mt-3">
                  <img v-for="(img, idx) in rv.images" :key="img + idx" :src="img" class="w-16 h-16 rounded-lg object-cover border border-gray-100" />
                </div>
                <div v-if="rv.appends?.length" class="mt-3 bg-gray-50 rounded-lg p-3">
                  <div v-for="ap in rv.appends" :key="ap.id" class="mb-2 last:mb-0">
                    <p class="text-xs text-gray-500">追加：{{ ap.content || '仅图片追评' }}</p>
                    <div v-if="ap.images?.length" class="flex flex-wrap gap-2 mt-2">
                      <img v-for="(img, idx) in ap.images" :key="img + idx" :src="img" class="w-14 h-14 rounded object-cover border border-gray-100" />
                    </div>
                  </div>
                </div>
                <div v-if="rv.admin_reply" class="mt-3 bg-red-50 rounded-lg p-3 text-xs text-red-600">
                  商家回复：{{ rv.admin_reply.content }}
                </div>
              </div>
            </div>
            <p v-else class="text-sm text-gray-400">暂无评价</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { del, get, post } from '@/api/request'
import { useChatStore } from '@/stores/chat'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const chat = useChatStore()
const auth = useAuthStore()
const product = ref<any>({})
const skus = ref<any[]>([])
const selectedSku = ref<any>(null)
const images = ref<string[]>([])
const mainImage = ref('')
const qty = ref(1)
const activeTab = ref<'detail' | 'review'>('detail')
const tabs: Array<{ label: string; value: 'detail' | 'review' }> = [
  { label: '详情', value: 'detail' },
  { label: '评价', value: 'review' },
]
const reviews = ref<any[]>([])
const reviewSummary = ref<any>({ avg_product_score: 0, avg_logistics_score: 0, total: 0 })

const detailBlocks = computed(() => {
  const detail = product.value?.detail
  if (!detail) return []
  const normalized = typeof detail === 'string' ? (() => {
    try { return JSON.parse(detail) } catch { return null }
  })() : detail
  if (!normalized || !Array.isArray(normalized.blocks)) return []
  return normalized.blocks
})

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

function openChat() {
  chat.open('product_detail')
}

async function toggleFavorite() {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  const id = Number(product.value?.id || 0)
  if (!id) return
  const favored = !!product.value?.is_favorited
  if (favored) {
    await del(`/api/v1/products/${id}/favorite`)
    product.value.is_favorited = false
    product.value.favorite_count = Math.max(0, Number(product.value.favorite_count || 0) - 1)
    return
  }
  await post(`/api/v1/products/${id}/favorite`)
  product.value.is_favorited = true
  product.value.favorite_count = Number(product.value.favorite_count || 0) + 1
}

function formatDate(v: string) {
  return v ? v.slice(0, 19).replace('T', ' ') : '-'
}

function setTab(value: 'detail' | 'review') {
  activeTab.value = value
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

  const reviewData = await get<any>(`/api/v1/products/${route.params.id}/reviews`, { page: 1, size: 20 })
  reviewSummary.value = reviewData?.summary || { avg_product_score: 0, avg_logistics_score: 0, total: 0 }
  reviews.value = reviewData?.list || []
})
</script>
