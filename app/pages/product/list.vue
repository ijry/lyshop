<template>
  <view class="min-h-screen bg-gray-50">
    <view class="px-20rpx py-16rpx bg-white">
      <u-search v-model="keyword" placeholder="搜索商品" @search="onSearch" @clear="onSearch" />
    </view>

    <view class="bg-white border-b-1 border-gray-100 px-8rpx">
      <u-tabs
        :list="tabList"
        :current="activeTab"
        lineColor="#dc2626"
        activeStyle="color:#dc2626;font-weight:600;"
        @click="onTab"
      />
    </view>

    <view class="p-16rpx">
      <view v-if="loading" class="flex justify-center py-80rpx">
        <u-loading-icon text="加载中..." />
      </view>

      <up-waterfall
        v-else
        :key="waterfallVersion"
        :modelValue="products"
        :addTime="30"
      >
        <template #left="{ leftList }">
          <view v-for="p in leftList" :key="`left-${p.id}`" class="px-8rpx pb-16rpx">
            <view @click="toDetail(p.id)" class="bg-white rounded-20rpx overflow-hidden shadow-sm">
              <image :src="p.cover" mode="aspectFill" class="w-full h-320rpx" />
              <view class="p-20rpx">
                <text class="text-26rpx text-gray-800 font-500 line-clamp-2">{{ p.title }}</text>
                <view class="flex items-center justify-between mt-16rpx">
                  <text class="text-30rpx text-red-600 font-700">¥{{ p.price }}</text>
                  <text class="text-22rpx text-gray-400">销量 {{ p.sales }}</text>
                </view>
              </view>
            </view>
          </view>
        </template>
        <template #right="{ rightList }">
          <view v-for="p in rightList" :key="`right-${p.id}`" class="px-8rpx pb-16rpx">
            <view @click="toDetail(p.id)" class="bg-white rounded-20rpx overflow-hidden shadow-sm">
              <image :src="p.cover" mode="aspectFill" class="w-full h-320rpx" />
              <view class="p-20rpx">
                <text class="text-26rpx text-gray-800 font-500 line-clamp-2">{{ p.title }}</text>
                <view class="flex items-center justify-between mt-16rpx">
                  <text class="text-30rpx text-red-600 font-700">¥{{ p.price }}</text>
                  <text class="text-22rpx text-gray-400">销量 {{ p.sales }}</text>
                </view>
              </view>
            </view>
          </view>
        </template>
      </up-waterfall>

      <view v-if="!loading && !products.length" class="text-center py-120rpx text-gray-400 text-28rpx">
        暂无商品
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { get } from '@/utils/request'

const products = ref<any[]>([])
const categories = ref<any[]>([{ id: '', name: '全部' }])
const keyword = ref('')
const categoryID = ref<string | number>('')
const activeTab = ref(0)
const loading = ref(false)
const page = ref(1)
const waterfallVersion = ref(0)

const tabList = computed(() => categories.value.map((item) => ({ name: item.name })))

function syncActiveTab() {
  const index = categories.value.findIndex((item) => String(item.id) === String(categoryID.value))
  activeTab.value = index >= 0 ? index : 0
}

async function loadProducts() {
  loading.value = true
  try {
    const data = await get<any>('/api/v1/products', {
      keyword: keyword.value,
      category_id: categoryID.value || undefined,
      page: page.value,
      size: 20
    })
    products.value = data?.list || []
    waterfallVersion.value += 1
  } finally {
    loading.value = false
  }
}

function onSearch() {
  page.value = 1
  loadProducts()
}

function onTab(item: any) {
  activeTab.value = item.index
  categoryID.value = categories.value[item.index]?.id ?? ''
  page.value = 1
  loadProducts()
}

function toDetail(id: number) {
  uni.navigateTo({ url: `/pages/product/detail?id=${id}` })
}

onMounted(async () => {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any)?.options || {}
  categoryID.value = query.category_id || ''

  const cats = await get<any[]>('/api/v1/categories')
  if (Array.isArray(cats) && cats.length) {
    categories.value = [{ id: '', name: '全部' }, ...cats]
  }
  syncActiveTab()
  loadProducts()
})
</script>
