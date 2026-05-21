<template>
  <view class="min-h-screen bg-gray-50">
    <!-- Search -->
    <view class="px-20rpx py-16rpx bg-white">
      <u-search v-model="keyword" placeholder="搜索商品" @search="onSearch" @clear="onSearch" />
    </view>

    <!-- Category tabs -->
    <scroll-view scroll-x class="bg-white border-b-1 border-gray-100">
      <view class="flex px-20rpx py-16rpx gap-16rpx">
        <view v-for="c in categories" :key="c.id"
          @click="categoryID = c.id; loadProducts()"
          :class="categoryID === c.id
            ? 'bg-blue-700 text-white'
            : 'bg-gray-100 text-gray-600'"
          class="px-24rpx py-10rpx rounded-full text-26rpx whitespace-nowrap flex-shrink-0">
          {{ c.name }}
        </view>
      </view>
    </scroll-view>

    <!-- Product grid -->
    <view class="p-20rpx">
      <view v-if="loading" class="flex justify-center py-80rpx">
        <u-loading-icon text="加载中..." />
      </view>
      <view v-else class="flex flex-wrap">
        <view v-for="p in products" :key="p.id"
          @click="toDetail(p.id)"
          class="w-1/2 p-8rpx">
          <view class="bg-white rounded-20rpx overflow-hidden shadow-sm">
            <image :src="p.cover" mode="aspectFill" class="w-full h-320rpx" />
            <view class="p-20rpx">
              <text class="text-26rpx text-gray-800 font-500 line-clamp-2">{{ p.title }}</text>
              <view class="flex items-center justify-between mt-16rpx">
                <text class="text-30rpx text-blue-700 font-700">¥{{ p.price }}</text>
                <text class="text-22rpx text-gray-400">销量 {{ p.sales }}</text>
              </view>
            </view>
          </view>
        </view>
      </view>
      <view v-if="!loading && !products.length" class="text-center py-120rpx text-gray-400 text-28rpx">
        暂无商品
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
