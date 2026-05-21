<template>
  <view class="bg-gray-50 min-h-screen">
    <u-navbar title="商品列表" :placeholder="true" />

    <!-- Search -->
    <view class="px-3 py-2">
      <u-search v-model="keyword" placeholder="搜索商品" @search="onSearch" @clear="onSearch" />
    </view>

    <!-- Category tabs -->
    <scroll-view scroll-x class="bg-white border-b border-gray-100">
      <view class="flex px-3 py-2 gap-3">
        <view
          v-for="c in categories" :key="c.id"
          @click="categoryID = c.id; loadProducts()"
          :class="categoryID === c.id ? 'bg-blue-700 text-white' : 'bg-gray-100 text-gray-600'"
          class="px-3 py-1 rounded-full text-sm whitespace-nowrap"
        >{{ c.name }}</view>
      </view>
    </scroll-view>

    <!-- Product grid -->
    <view class="p-3">
      <view v-if="loading" class="text-center py-8">
        <u-loading-icon text="加载中..." />
      </view>
      <view v-else class="grid grid-cols-2 gap-3">
        <view
          v-for="p in products" :key="p.id"
          @click="toDetail(p.id)"
          class="bg-white rounded-xl overflow-hidden shadow-sm"
        >
          <image
            :src="p.cover || '/static/placeholder.png'"
            mode="aspectFill"
            class="w-full"
            style="height: 160px;"
          />
          <view class="p-3">
            <text class="text-slate-800 text-sm font-medium line-clamp-2">{{ p.title }}</text>
            <view class="flex items-center justify-between mt-2">
              <text class="text-blue-700 font-bold text-base">¥{{ p.price }}</text>
              <text class="text-gray-400 text-xs">销量 {{ p.sales }}</text>
            </view>
          </view>
        </view>
      </view>
      <view v-if="!loading && !products.length" class="text-center py-12 text-gray-400">
        <text>暂无商品</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'

const products = ref<any[]>([])
const categories = ref<any[]>([{ id: '', name: '全部' }])
const keyword = ref('')
const categoryID = ref<string | number>('')
const loading = ref(false)
const page = ref(1)

async function loadProducts() {
  loading.value = true
  try {
    const data = await get<any>('/api/v1/products', {
      keyword: keyword.value,
      category_id: categoryID.value || undefined,
      page: page.value,
      size: 20
    })
    products.value = data.list || []
  } finally {
    loading.value = false
  }
}

function onSearch() { page.value = 1; loadProducts() }

function toDetail(id: number) {
  uni.navigateTo({ url: `/pages/product/detail?id=${id}` })
}

onMounted(async () => {
  const cats = await get<any[]>('/api/v1/categories')
  if (cats) categories.value = [{ id: '', name: '全部' }, ...cats]
  loadProducts()
})
</script>
