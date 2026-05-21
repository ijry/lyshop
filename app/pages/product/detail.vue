<template>
  <view class="min-h-screen bg-white pb-120rpx">
    <!-- Images carousel -->
    <u-swiper :list="images" height="375" v-if="images.length" radius="0" />
    <view v-else class="bg-gray-100 h-750rpx" />

    <!-- Price + Title -->
    <view class="p-30rpx">
      <view class="flex items-baseline gap-12rpx mb-16rpx">
        <text class="text-48rpx font-700 text-red-500">¥{{ product.price }}</text>
        <text v-if="product.origin_price" class="text-gray-400 line-through text-26rpx">
          ¥{{ product.origin_price }}
        </text>
        <view v-if="product.origin_price" class="bg-red-50 text-red-500 text-20rpx px-10rpx py-2rpx rounded-6rpx ml-8rpx">
          省¥{{ (product.origin_price - product.price).toFixed(0) }}
        </view>
      </view>
      <text class="text-32rpx font-600 text-gray-800 block leading-44rpx">{{ product.title }}</text>
      <text v-if="product.subtitle" class="text-26rpx text-gray-500 mt-8rpx block">{{ product.subtitle }}</text>
      <view class="flex items-center mt-16rpx gap-20rpx">
        <text class="text-24rpx text-gray-400">库存 {{ selectedSku?.stock ?? product.stock }}</text>
        <text class="text-24rpx text-gray-400">销量 {{ product.sales }}</text>
      </view>
    </view>

    <!-- Divider -->
    <view class="h-16rpx bg-gray-50" />

    <!-- SKU selector -->
    <view class="p-30rpx" v-if="skus.length">
      <text class="text-28rpx font-600 text-gray-800 block mb-20rpx">规格选择</text>
      <view class="flex flex-wrap gap-16rpx">
        <view v-for="sku in skus" :key="sku.id"
          @click="selectedSku = sku"
          :class="selectedSku?.id === sku.id
            ? 'border-blue-700 text-blue-700 bg-blue-50'
            : 'border-gray-200 text-gray-600 bg-white'"
          class="px-24rpx py-14rpx border-1 rounded-12rpx text-26rpx">
          <text v-if="sku.attrs">{{ parseAttrs(sku.attrs) }}</text>
          <text v-else>默认规格</text>
        </view>
      </view>
    </view>

    <!-- Divider -->
    <view class="h-16rpx bg-gray-50" />

    <!-- Detail -->
    <view class="p-30rpx">
      <text class="text-28rpx font-600 text-gray-800 block mb-20rpx">商品详情</text>
      <rich-text :nodes="product.detail || ''" class="text-26rpx text-gray-600" />
    </view>

    <!-- Bottom bar -->
    <view class="fixed bottom-0 left-0 right-0 z-100 bg-white border-t-1 border-gray-100 flex items-center px-20rpx py-16rpx"
      :style="{paddingBottom: 'calc(16rpx + env(safe-area-inset-bottom))'}">
      <!-- Icons -->
      <view class="flex items-center gap-30rpx mr-20rpx">
        <view class="flex flex-col items-center" @click="uni.switchTab({url:'/pages/index/index'})">
          <u-icon name="home" size="22" color="#666" />
          <text class="text-20rpx text-gray-500 mt-4rpx">首页</text>
        </view>
        <view class="flex flex-col items-center" @click="uni.switchTab({url:'/pages/cart/index'})">
          <u-icon name="shopping-cart" size="22" color="#666" />
          <text class="text-20rpx text-gray-500 mt-4rpx">购物车</text>
        </view>
      </view>
      <!-- Buttons -->
      <view class="flex-1 flex gap-16rpx">
        <u-button type="warning" text="加入购物车" @click="addCart" class="flex-1" shape="circle" />
        <u-button type="primary" text="立即购买" @click="buyNow" class="flex-1" shape="circle" />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get, post } from '@/utils/request'

const product = ref<any>({})
const skus = ref<any[]>([])
const selectedSku = ref<any>(null)
const images = ref<string[]>([])

function parseAttrs(attrs: string) {
  try { return JSON.parse(attrs).map((a: any) => a.value).join(' ') }
  catch { return '默认' }
}

onMounted(async () => {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  const id = query.id
  const data = await get<any>(`/api/v1/products/${id}`)
  if (!data) return
  product.value = data
  skus.value = data.skus || []
  if (skus.value.length) selectedSku.value = skus.value[0]
  const imgList: string[] = []
  if (data.cover) imgList.push(data.cover)
  if (data.images) imgList.push(...data.images.map((i: any) => i.url))
  images.value = imgList
})

async function addCart() {
  const skuID = selectedSku.value?.id
  if (!skuID) { uni.showToast({ title: '请选择规格', icon: 'none' }); return }
  await post('/api/v1/cart/add', { sku_id: skuID, qty: 1 })
  uni.showToast({ title: '已加入购物车', icon: 'success' })
}

function buyNow() {
  const skuID = selectedSku.value?.id
  if (!skuID) { uni.showToast({ title: '请选择规格', icon: 'none' }); return }
  uni.navigateTo({ url: `/pages/order/confirm?sku_ids=${skuID}` })
}
</script>
