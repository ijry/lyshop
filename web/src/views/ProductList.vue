<template>
  <div class="max-w-7xl mx-auto px-6 py-8">
    <div class="flex gap-6">
      <!-- Sidebar categories -->
      <aside class="hidden md:block w-52 shrink-0">
        <div class="sticky top-24 card p-4">
          <h3 class="text-sm font-semibold text-gray-800 mb-3">商品分类</h3>
          <div class="flex flex-col gap-0.5">
            <button v-for="c in categories" :key="c.id"
              @click="categoryID = c.id; loadProducts()"
              :class="categoryID === c.id
                ? 'bg-blue-50 text-blue-700 font-medium'
                : 'text-gray-600 hover:bg-gray-50'"
              class="text-left text-sm px-3 py-2 rounded-lg transition-colors">
              {{ c.name }}
            </button>
          </div>
        </div>
      </aside>

      <!-- Main content -->
      <div class="flex-1 min-w-0">
        <!-- Search -->
        <div class="mb-6">
          <div class="relative">
            <div class="i-carbon-search absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" />
            <input v-model="keyword" @keyup.enter="loadProducts"
              placeholder="搜索你想要的商品..."
              class="input-base pl-11 !rounded-xl !py-3" />
          </div>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="flex-center py-20">
          <div class="w-6 h-6 border-2 border-blue-700 border-t-transparent rounded-full animate-spin" />
        </div>

        <!-- Products grid -->
        <div v-else-if="products.length" class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
          <ProductCard v-for="p in products" :key="p.id" :product="p" />
        </div>

        <!-- Empty -->
        <div v-else class="flex flex-col items-center py-20">
          <div class="i-carbon-search text-5xl text-gray-200 mb-4" />
          <p class="text-gray-400 text-sm">暂无商品</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/api/request'
import ProductCard from '@/components/ProductCard.vue'

const products = ref<any[]>([])
const categories = ref<any[]>([{ id: '', name: '全部' }])
const keyword = ref('')
const categoryID = ref<string | number>('')
const loading = ref(false)

async function loadProducts() {
  loading.value = true
  try {
    const data = await get<any>('/api/v1/products', {
      keyword: keyword.value,
      category_id: categoryID.value || undefined,
    })
    products.value = data?.list || []
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const cats = await get<any[]>('/api/v1/categories')
  if (cats) categories.value = [{ id: '', name: '全部' }, ...cats]
  loadProducts()
})
</script>
