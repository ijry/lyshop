<template>
  <view class="bg-white min-h-screen pb-24">
    <u-navbar :title="product.title || '商品详情'" :placeholder="true" />

    <!-- Images carousel -->
    <u-swiper :list="images" height="300" v-if="images.length" />
    <view v-else class="bg-gray-100" style="height: 300px;" />

    <!-- Product info -->
    <view class="p-4">
      <view class="flex items-baseline gap-2 mb-2">
        <text class="text-2xl font-bold text-blue-700">¥{{ product.price }}</text>
        <text v-if="product.origin_price" class="text-gray-400 line-through text-sm">
          ¥{{ product.origin_price }}
        </text>
      </view>
      <text class="text-lg font-semibold text-slate-800 block mb-1">{{ product.title }}</text>
      <text class="text-gray-500 text-sm">库存 {{ selectedSku?.stock ?? product.stock }}</text>
    </view>

    <u-divider />

    <!-- SKU selector -->
    <view class="px-4 pb-4" v-if="skus.length">
      <text class="text-sm font-medium text-slate-700 block mb-3">规格选择</text>
      <view class="flex flex-wrap gap-2">
        <view
          v-for="sku in skus" :key="sku.id"
          @click="selectedSku = sku"
          :class="selectedSku?.id === sku.id
            ? 'border-blue-700 text-blue-700 bg-blue-50'
            : 'border-gray-200 text-gray-600'"
          class="px-3 py-1 border rounded-lg text-sm"
        >
          <text v-if="sku.attrs">{{ JSON.parse(sku.attrs).map((a:any)=>a.value).join(' ') }}</text>
          <text v-else>默认规格</text>
        </view>
      </view>
    </view>

    <!-- Detail -->
    <view class="px-4 pt-2">
      <text class="text-sm font-medium text-slate-700 block mb-3">商品详情</text>
      <rich-text :nodes="product.detail || ''" class="text-sm text-gray-600" />
    </view>

    <!-- Bottom bar -->
    <view class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-100 p-3 flex gap-3">
      <u-button type="warning" plain text="加入购物车" @click="addCart" class="flex-1" />
      <u-button type="primary" text="立即购买" @click="buyNow" class="flex-1" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { get, post } from '@/utils/request'

const product = ref<any>({})
const skus = ref<any[]>([])
const selectedSku = ref<any>(null)
const images = ref<string[]>([])

onMounted(async () => {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  const id = query.id
  const data = await get<any>(`/api/v1/products/${id}`)
  if (!data) return
  product.value = data
  skus.value = data.skus || []
  if (skus.value.length) selectedSku.value = skus.value[0]
  // Build image list: cover + album images
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
