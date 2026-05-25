<template>
  <view class="min-h-screen pb-120rpx" style="background: var(--app-card-bg);">
    <u-swiper :list="images" height="375" v-if="images.length" radius="0" />
    <view v-else class="bg-gray-100 h-750rpx" />

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
        <text class="text-24rpx text-gray-400">收藏 {{ product.favorite_count || 0 }}</text>
      </view>
    </view>

    <view class="h-16rpx bg-gray-50" />

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

    <view class="h-16rpx bg-gray-50" />
    <view class="bg-white">
      <u-tabs :list="detailTabs" :current="activeTab" @click="onTabClick" />
    </view>

    <view class="p-30rpx" v-if="activeTab === 0">
      <view v-if="detailBlocks.length">
        <view v-for="block in detailBlocks" :key="block.id" class="mb-20rpx">
          <text v-if="block.type === 'text'" class="text-26rpx text-gray-600 leading-42rpx block">{{ block.props?.text || '' }}</text>
          <image v-else-if="block.type === 'image'" :src="block.props?.url || ''" mode="widthFix" style="width:100%; border-radius: 12px;" />
        </view>
      </view>
      <text v-else class="text-24rpx text-gray-400">暂无详情</text>
    </view>

    <view class="p-30rpx" v-else>
      <view class="bg-gray-50 rounded-16rpx p-20rpx mb-20rpx">
        <view class="flex items-center justify-between text-24rpx text-gray-500">
          <text>商品评分</text>
          <text class="text-red-500 font-600">{{ reviewSummary.avg_product_score.toFixed(1) }}</text>
        </view>
        <view class="flex items-center justify-between text-24rpx text-gray-500 mt-8rpx">
          <text>物流评分</text>
          <text class="text-red-500 font-600">{{ reviewSummary.avg_logistics_score.toFixed(1) }}</text>
        </view>
        <view class="text-22rpx text-gray-400 mt-10rpx">共 {{ reviewSummary.total }} 条评价</view>
      </view>

          <view v-if="reviews.length">
        <view v-for="rv in reviews" :key="rv.id" class="bg-white border border-gray-100 rounded-16rpx p-20rpx mb-16rpx">
          <view class="flex items-center justify-between mb-10rpx">
            <text class="text-24rpx text-gray-700">{{ rv.user_nickname || '匿名用户' }}</text>
            <text class="text-22rpx text-gray-400">{{ formatDate(rv.created_at) }}</text>
          </view>
          <view class="flex items-center gap-12rpx mb-10rpx">
            <text class="text-22rpx text-gray-500">商品 {{ rv.product_score }}</text>
            <text class="text-22rpx text-gray-500">物流 {{ rv.logistics_score }}</text>
          </view>
          <text class="text-24rpx text-gray-700 leading-40rpx block">{{ rv.content || '用户未填写文字评价' }}</text>
          <view v-if="rv.images?.length" class="flex flex-wrap gap-10rpx mt-12rpx">
            <image
              v-for="(img, idx) in rv.images"
              :key="img + idx"
              :src="img"
              mode="aspectFill"
              style="width: 136rpx; height: 136rpx; border-radius: 12rpx;"
            />
          </view>
          <view v-if="rv.appends?.length" class="mt-12rpx bg-gray-50 rounded-12rpx p-12rpx">
            <view v-for="ap in rv.appends" :key="ap.id" class="mb-10rpx last:mb-0">
              <view class="text-22rpx text-gray-500 leading-34rpx">追加：{{ ap.content || '仅图片追评' }}</view>
              <view v-if="ap.images?.length" class="flex flex-wrap gap-10rpx mt-8rpx">
                <image
                  v-for="(img, idx) in ap.images"
                  :key="img + idx"
                  :src="img"
                  mode="aspectFill"
                  style="width: 120rpx; height: 120rpx; border-radius: 10rpx;"
                />
              </view>
            </view>
          </view>
          <view v-if="rv.admin_reply" class="mt-12rpx bg-red-50 rounded-12rpx p-12rpx text-22rpx text-red-600">
            商家回复：{{ rv.admin_reply.content }}
          </view>
        </view>
      </view>
      <view v-else class="text-center py-80rpx text-24rpx text-gray-400">暂无评价</view>
    </view>

    <view class="fixed bottom-0 left-0 right-0 z-100 flex items-center px-20rpx py-16rpx"
      :style="{ background: 'var(--app-card-bg)', borderTop: '1px solid var(--app-border-color)', paddingBottom: 'calc(16rpx + env(safe-area-inset-bottom))' }">
      <view class="flex items-center gap-30rpx mr-20rpx">
        <view class="flex flex-col items-center" @click="uni.navigateTo({url:'/pages/im/chat'})">
          <u-icon name="kefu-ermai" size="22" color="#666" />
          <text class="text-20rpx text-gray-500 mt-4rpx">客服</text>
        </view>
        <view class="flex flex-col items-center" @click="uni.switchTab({url:'/pages/index/index'})">
          <u-icon name="home" size="22" color="#666" />
          <text class="text-20rpx text-gray-500 mt-4rpx">首页</text>
        </view>
        <view class="flex flex-col items-center" @click="uni.switchTab({url:'/pages/cart/index'})">
          <u-icon name="shopping-cart" size="22" color="#666" />
          <text class="text-20rpx text-gray-500 mt-4rpx">购物车</text>
        </view>
        <view class="flex flex-col items-center" @click="toggleFavorite">
          <u-icon :name="product.is_favorited ? 'heart-fill' : 'heart'" size="22" :color="product.is_favorited ? '#ef4444' : '#666'" />
          <text class="text-20rpx mt-4rpx" :style="{ color: product.is_favorited ? '#ef4444' : '#6b7280' }">{{ product.is_favorited ? '已收藏' : '收藏' }}</text>
        </view>
      </view>
      <view class="flex-1 flex gap-16rpx">
        <u-button type="warning" text="加入购物车" @click="addCart" class="flex-1" shape="circle" />
        <u-button type="primary" text="立即购买" @click="buyNow" class="flex-1" shape="circle" />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { del, get, post } from '@/utils/request'

const product = ref<any>({})
const skus = ref<any[]>([])
const selectedSku = ref<any>(null)
const images = ref<string[]>([])
const reviews = ref<any[]>([])
const reviewSummary = ref<any>({ avg_product_score: 0, avg_logistics_score: 0, total: 0 })
const activeTab = ref(0)
const detailTabs = [{ name: '详情' }, { name: '评价' }]

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
  try { return JSON.parse(attrs).map((a: any) => a.value).join(' ') }
  catch { return '默认' }
}

function formatDate(v: string) {
  return v ? String(v).slice(0, 19).replace('T', ' ') : '-'
}

async function loadReviews(id: number) {
  const data = await get<any>(`/api/v1/products/${id}/reviews`, { page: 1, size: 20 })
  reviewSummary.value = data?.summary || { avg_product_score: 0, avg_logistics_score: 0, total: 0 }
  reviews.value = Array.isArray(data?.list) ? data.list : []
}

function onTabClick(item: any) {
  activeTab.value = Number(item?.index || 0)
}

onMounted(async () => {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  const id = Number(query.id)
  const data = await get<any>(`/api/v1/products/${id}`)
  if (!data) return
  product.value = data
  skus.value = data.skus || []
  if (skus.value.length) selectedSku.value = skus.value[0]
  const imgList: string[] = []
  if (data.cover) imgList.push(data.cover)
  if (data.images) imgList.push(...data.images.map((i: any) => i.url))
  images.value = imgList
  await loadReviews(id)
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

async function toggleFavorite() {
  const token = uni.getStorageSync('user_token')
  if (!token) {
    uni.navigateTo({ url: '/pages/login/index' })
    return
  }
  const id = Number(product.value?.id || 0)
  if (!id) return
  const favored = !!product.value?.is_favorited
  if (favored) {
    await del(`/api/v1/products/${id}/favorite`)
    product.value.is_favorited = false
    product.value.favorite_count = Math.max(0, Number(product.value.favorite_count || 0) - 1)
    uni.showToast({ title: '已取消收藏', icon: 'none' })
    return
  }
  await post(`/api/v1/products/${id}/favorite`)
  product.value.is_favorited = true
  product.value.favorite_count = Number(product.value.favorite_count || 0) + 1
  uni.showToast({ title: '收藏成功', icon: 'success' })
}
</script>
