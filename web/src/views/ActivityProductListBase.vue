<template>
  <div class="max-w-7xl mx-auto px-6 py-8">
    <div class="mb-6">
      <h1 class="text-2xl font-semibold text-gray-900 mb-2">{{ pageTitle }}</h1>
      <p class="text-sm text-gray-500">{{ pageDesc }}</p>
    </div>

    <div class="bg-white rounded-xl border border-gray-100 p-4 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-3">
        <select v-model="query.category_id" class="input-base">
          <option :value="''">{{ $t('orderList.all') }}</option>
          <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>
        <input v-model="query.min_price" type="number" min="0" :placeholder="$t('activityList.minPrice')" class="input-base" />
        <input v-model="query.max_price" type="number" min="0" :placeholder="$t('activityList.maxPrice')" class="input-base" />
        <input v-model="query.keyword" :placeholder="$t('productList.searchPlaceholder')" class="input-base" @keyup.enter="search" />
      </div>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-3 mt-3">
        <select v-model="query.sort_by" class="input-base">
          <option value="price">{{ $t('activityList.sortPrice') }}</option>
          <option value="sales">{{ $t('activityList.sortSales') }}</option>
        </select>
        <select v-model="query.sort_order" class="input-base">
          <option value="asc">{{ $t('activityList.sortAsc') }}</option>
          <option value="desc">{{ $t('activityList.sortDesc') }}</option>
        </select>
        <button class="btn-primary" @click="search">{{ $t('common.search') }}</button>
      </div>
    </div>

    <div v-if="loading" class="flex-center py-20">
      <div class="w-6 h-6 border-2 border-red-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <template v-else>
      <div v-if="products.length" class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
        <div
          v-for="p in products"
          :key="`${p.activity_id}_${p.sku_id}`"
          class="bg-white rounded-xl border border-gray-100 overflow-hidden cursor-pointer hover:shadow-lg hover:border-gray-200 transition-all"
          @click="openDetail(p)"
        >
          <div class="aspect-square bg-gray-50 overflow-hidden">
            <img :src="p.cover" :alt="p.title" class="w-full h-full object-cover" />
          </div>
          <div class="p-4">
            <h3 class="text-sm font-medium text-gray-800 line-clamp-2 mb-2">{{ p.title }}</h3>
            <p v-if="p.subtitle" class="text-xs text-gray-400 line-clamp-1 mb-2">{{ p.subtitle }}</p>
            <div class="flex items-center gap-2 mb-2">
              <span class="text-lg font-bold text-red-500">¥{{ p.price }}</span>
              <span class="text-xs text-gray-300 line-through">¥{{ p.origin_price }}</span>
            </div>
            <div class="text-xs text-gray-500 flex justify-between">
              <span>{{ $t('productDetail.sold') }} {{ p.sales }}</span>
              <span>{{ $t('productDetail.stock') }} {{ p.stock }}</span>
            </div>
            <div class="text-xs text-gray-500 mt-1">
              <span>{{ $t('activityList.limitPerOrder') }} {{ p.limit_per_order || '-' }}</span>
              <span class="mx-1">|</span>
              <span>{{ $t('activityList.stockLimit') }} {{ p.total_stock_limit || '-' }}</span>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="flex flex-col items-center py-20">
        <div class="i-carbon-search text-5xl text-gray-200 mb-4" />
        <p class="text-gray-400 text-sm">{{ $t('productList.empty') }}</p>
      </div>

      <div class="flex items-center justify-between mt-6 text-sm text-gray-500" v-if="total > 0">
        <span>{{ $t('activityList.total', { total }) }}</span>
        <div class="flex items-center gap-2">
          <button class="px-3 py-1 rounded-lg border border-gray-200 disabled:opacity-40" :disabled="query.page <= 1" @click="prevPage">
            {{ $t('common.prevPage') }}
          </button>
          <span>{{ query.page }}</span>
          <button class="px-3 py-1 rounded-lg border border-gray-200 disabled:opacity-40" :disabled="query.page * query.size >= total" @click="nextPage">
            {{ $t('common.nextPage') }}
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { get } from '@/api/request'

const props = defineProps<{ activityType: 'seckill' | 'group-buy' | 'bargain' }>()
const { t } = useI18n()
const router = useRouter()

const categories = ref<any[]>([])
const products = ref<any[]>([])
const total = ref(0)
const loading = ref(false)
const query = reactive({
  category_id: '' as string | number,
  keyword: '',
  min_price: '' as string | number,
  max_price: '' as string | number,
  sort_by: 'price',
  sort_order: 'asc',
  page: 1,
  size: 20,
})

const pageTitle = computed(() => {
  if (props.activityType === 'seckill') return t('activityList.seckillTitle')
  if (props.activityType === 'group-buy') return t('activityList.groupBuyTitle')
  return t('activityList.bargainTitle')
})
const pageDesc = computed(() => {
  if (props.activityType === 'seckill') return t('activityList.seckillDesc')
  if (props.activityType === 'group-buy') return t('activityList.groupBuyDesc')
  return t('activityList.bargainDesc')
})

const endpoint = computed(() => {
  if (props.activityType === 'seckill') return '/api/v1/seckill/products'
  if (props.activityType === 'group-buy') return '/api/v1/marketing/group-buy/products'
  return '/api/v1/marketing/bargain/products'
})

async function loadCategories() {
  const data = await get<any[]>('/api/v1/categories')
  categories.value = data || []
}

async function loadProducts() {
  loading.value = true
  try {
    const data = await get<any>(endpoint.value, {
      category_id: query.category_id || undefined,
      keyword: query.keyword || undefined,
      min_price: query.min_price || undefined,
      max_price: query.max_price || undefined,
      sort_by: query.sort_by,
      sort_order: query.sort_order,
      page: query.page,
      size: query.size,
    })
    products.value = data?.list || []
    total.value = Number(data?.total || 0)
  } finally {
    loading.value = false
  }
}

function search() {
  query.page = 1
  loadProducts()
}

function prevPage() {
  if (query.page <= 1) return
  query.page -= 1
  loadProducts()
}

function nextPage() {
  if (query.page * query.size >= total.value) return
  query.page += 1
  loadProducts()
}

function openDetail(item: any) {
  const productID = Number(item?.product_id || 0)
  if (!productID) return
  const activityProductID = Number(item?.activity_product_id || 0)
  const path = activityProductID > 0
    ? `/product/${productID}?activity_product_id=${activityProductID}`
    : `/product/${productID}`
  router.push(path)
}

onMounted(async () => {
  await loadCategories()
  await loadProducts()
})
</script>
