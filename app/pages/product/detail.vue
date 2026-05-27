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
          {{ $t('productDetail.save') }}{{ (product.origin_price - product.price).toFixed(0) }}
        </view>
      </view>
      <text class="text-32rpx font-600 text-gray-800 block leading-44rpx">{{ product.title }}</text>
      <text v-if="product.subtitle" class="text-26rpx text-gray-500 mt-8rpx block">{{ product.subtitle }}</text>
      <view class="flex items-center mt-16rpx gap-20rpx">
        <text class="text-24rpx text-gray-400">{{ $t('productDetail.stock') }} {{ selectedSku?.stock ?? product.stock }}</text>
        <text class="text-24rpx text-gray-400">{{ $t('productDetail.sales') }} {{ product.sales }}</text>
        <text class="text-24rpx text-gray-400">{{ $t('productDetail.favorite') }} {{ product.favorite_count || 0 }}</text>
      </view>
    </view>

    <view class="h-16rpx bg-gray-50" />

    <view v-if="marketingDetail" class="p-30rpx bg-white">
      <view class="rounded-16rpx p-20rpx" :style="marketingCardStyle">
        <view class="flex items-center justify-between">
          <view>
            <text class="text-24rpx font-600 block" :style="{ color: marketingTextColor }">{{ marketingTypeLabel }}</text>
            <text class="text-22rpx block mt-6rpx" :style="{ color: marketingSubTextColor }">{{ marketingDetail.activity_name || '-' }}</text>
          </view>
          <text class="text-22rpx" :style="{ color: marketingTextColor }">{{ marketingStatusLabel }}</text>
        </view>
        <view class="flex items-end gap-12rpx mt-14rpx">
          <text class="text-42rpx font-700" :style="{ color: marketingTextColor }">¥{{ marketingPrice }}</text>
          <text class="text-22rpx text-gray-400 line-through">¥{{ marketingOriginPrice }}</text>
        </view>
        <view class="mt-12rpx">
          <text class="text-22rpx text-gray-500 block">{{ $t('productDetail.activityLimitPerOrder') }} {{ marketingDetail.limit_per_order || '-' }}</text>
          <text class="text-22rpx text-gray-500 block mt-4rpx">{{ $t('productDetail.activityStockProgress') }} {{ marketingDetail.sold_qty || 0 }}/{{ marketingDetail.total_stock_limit || '-' }}</text>
        </view>
        <view v-if="marketingCountdownLabel" class="mt-12rpx">
          <text class="text-22rpx" :style="{ color: marketingTextColor }">{{ marketingCountdownLabel }}</text>
        </view>
        <view class="mt-14rpx">
          <view class="inline-flex items-center px-16rpx py-8rpx rounded-999rpx text-22rpx"
            :style="{ color: '#fff', background: marketingActionBg }">
            {{ marketingActionLabel }}
          </view>
        </view>
      </view>
    </view>

    <view class="h-16rpx bg-gray-50" v-if="marketingDetail" />

    <view class="p-30rpx" v-if="skus.length">
      <text class="text-28rpx font-600 text-gray-800 block mb-20rpx">{{ $t('productDetail.specSelect') }}</text>
      <view class="flex flex-wrap gap-16rpx">
        <view v-for="sku in skus" :key="sku.id"
          @click="selectedSku = sku"
          :class="selectedSku?.id === sku.id
            ? 'border-blue-700 text-blue-700 bg-blue-50'
            : 'border-gray-200 text-gray-600 bg-white'"
          class="px-24rpx py-14rpx border-1 rounded-12rpx text-26rpx">
          <text v-if="sku.attrs">{{ parseAttrs(sku.attrs) }}</text>
          <text v-else>{{ $t('productDetail.defaultSpec') }}</text>
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
          <view v-else-if="block.type === 'rich_text'" class="rich-text-content text-26rpx text-gray-600 leading-42rpx" v-html="block.props?.content || ''" />
        </view>
      </view>
      <text v-else class="text-24rpx text-gray-400">{{ $t('productDetail.noDetail') }}</text>
    </view>

    <view class="p-30rpx" v-else>
      <view class="bg-gray-50 rounded-16rpx p-20rpx mb-20rpx">
        <view class="flex items-center justify-between text-24rpx text-gray-500">
          <text>{{ $t('productDetail.productScore') }}</text>
          <text class="text-red-500 font-600">{{ reviewSummary.avg_product_score.toFixed(1) }}</text>
        </view>
        <view class="flex items-center justify-between text-24rpx text-gray-500 mt-8rpx">
          <text>{{ $t('productDetail.logisticsScore') }}</text>
          <text class="text-red-500 font-600">{{ reviewSummary.avg_logistics_score.toFixed(1) }}</text>
        </view>
        <view class="text-22rpx text-gray-400 mt-10rpx">{{ $t('productDetail.total') }} {{ reviewSummary.total }} {{ $t('productDetail.reviewCount') }}</view>
      </view>

          <view v-if="reviews.length">
        <view v-for="rv in reviews" :key="rv.id" class="bg-white border border-gray-100 rounded-16rpx p-20rpx mb-16rpx">
          <view class="flex items-center justify-between mb-10rpx">
            <text class="text-24rpx text-gray-700">{{ rv.user_nickname || $t('productDetail.anonymous') }}</text>
            <text class="text-22rpx text-gray-400">{{ formatDate(rv.created_at) }}</text>
          </view>
          <view class="flex items-center gap-12rpx mb-10rpx">
            <text class="text-22rpx text-gray-500">{{ $t('productDetail.productScore') }} {{ rv.product_score }}</text>
            <text class="text-22rpx text-gray-500">{{ $t('productDetail.logisticsScore') }} {{ rv.logistics_score }}</text>
          </view>
          <text class="text-24rpx text-gray-700 leading-40rpx block">{{ rv.content || $t('productDetail.noTextReview') }}</text>
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
              <view class="text-22rpx text-gray-500 leading-34rpx">{{ $t('productDetail.append') }}{{ ap.content || $t('productDetail.imageOnlyAppend') }}</view>
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
            {{ $t('productDetail.merchantReply') }}{{ rv.admin_reply.content }}
          </view>
        </view>
      </view>
      <view v-else class="text-center py-80rpx text-24rpx text-gray-400">{{ $t('productDetail.noReview') }}</view>
    </view>

    <view class="fixed bottom-0 left-0 right-0 z-100 flex items-center px-20rpx py-16rpx"
      :style="{ background: 'var(--app-card-bg)', borderTop: '1px solid var(--app-border-color)', paddingBottom: 'calc(16rpx + env(safe-area-inset-bottom))' }">
      <view class="flex items-center gap-30rpx mr-20rpx">
        <view class="flex flex-col items-center" @click="uni.navigateTo({url:'/pages/im/chat'})">
          <u-icon name="kefu-ermai" size="22" color="#666" />
          <text class="text-20rpx text-gray-500 mt-4rpx">{{ $t('productDetail.customerService') }}</text>
        </view>
        <view class="flex flex-col items-center" @click="uni.switchTab({url:'/pages/index/index'})">
          <u-icon name="home" size="22" color="#666" />
          <text class="text-20rpx text-gray-500 mt-4rpx">{{ $t('productDetail.homePage') }}</text>
        </view>
        <view class="flex flex-col items-center" @click="uni.switchTab({url:'/pages/cart/index'})">
          <u-icon name="shopping-cart" size="22" color="#666" />
          <text class="text-20rpx text-gray-500 mt-4rpx">{{ $t('productDetail.cart') }}</text>
        </view>
        <view class="flex flex-col items-center" @click="toggleFavorite">
          <u-icon :name="product.is_favorited ? 'heart-fill' : 'heart'" size="22" :color="product.is_favorited ? '#ef4444' : '#666'" />
          <text class="text-20rpx mt-4rpx" :style="{ color: product.is_favorited ? '#ef4444' : '#6b7280' }">{{ product.is_favorited ? $t('productDetail.favorited') : $t('productDetail.favorite') }}</text>
        </view>
      </view>
      <view class="flex-1 flex gap-16rpx">
        <u-button type="warning" :text="$t('productDetail.addToCart')" @click="addCart" class="flex-1" shape="circle" />
        <u-button type="primary" :text="$t('productDetail.buyNow')" @click="buyNow" class="flex-1" shape="circle" />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { del, get, post } from '@/utils/request'

const { t } = useI18n()

const product = ref<any>({})
const skus = ref<any[]>([])
const selectedSku = ref<any>(null)
const images = ref<string[]>([])
const activityProductID = ref(0)
const marketingDetail = ref<any>(null)
const countdownText = ref('')
let countdownTimer: any = null
const reviews = ref<any[]>([])
const reviewSummary = ref<any>({ avg_product_score: 0, avg_logistics_score: 0, total: 0 })
const activeTab = ref(0)
const detailTabs = computed(() => [{ name: t('productDetail.detail') }, { name: t('productDetail.reviews') }])

const detailBlocks = computed(() => {
  const detail = product.value?.detail
  if (!detail) return []
  const normalized = typeof detail === 'string' ? (() => {
    try { return JSON.parse(detail) } catch { return null }
  })() : detail
  if (!normalized || !Array.isArray(normalized.blocks)) return []
  return normalized.blocks
})

const marketingTypeLabel = computed(() => {
  const type = String(marketingDetail.value?.activity_type || '')
  if (type === 'seckill') return t('productDetail.activityTypeSeckill')
  if (type === 'group_buy') return t('productDetail.activityTypeGroupBuy')
  if (type === 'bargain') return t('productDetail.activityTypeBargain')
  return t('productDetail.activityTypeDefault')
})

const marketingPrice = computed(() => Number(marketingDetail.value?.price || 0).toFixed(2))
const marketingOriginPrice = computed(() => Number(marketingDetail.value?.origin_price || product.value?.price || 0).toFixed(2))

const marketingStatusLabel = computed(() => {
  if (!marketingDetail.value) return ''
  const now = Date.now()
  const startAt = marketingDetail.value?.activity_start_at ? new Date(marketingDetail.value.activity_start_at).getTime() : 0
  const endAt = marketingDetail.value?.activity_end_at ? new Date(marketingDetail.value.activity_end_at).getTime() : 0
  if (marketingDetail.value?.activity_status !== 1) return t('productDetail.activityStatusInactive')
  if (startAt > now) return t('productDetail.activityStatusNotStarted')
  if (endAt > 0 && endAt <= now) return t('productDetail.activityStatusEnded')
  if (Number(marketingDetail.value?.total_stock_limit || 0) > 0 && Number(marketingDetail.value?.sold_qty || 0) >= Number(marketingDetail.value?.total_stock_limit || 0)) {
    return t('productDetail.activityStatusSoldOut')
  }
  return t('productDetail.activityStatusOngoing')
})

const marketingActionLabel = computed(() => {
  const type = String(marketingDetail.value?.activity_type || '')
  if (type === 'seckill') return t('productDetail.activityActionSeckill')
  if (type === 'group_buy') return t('productDetail.activityActionGroupBuy')
  if (type === 'bargain') return t('productDetail.activityActionBargain')
  return t('productDetail.activityActionDefault')
})

const marketingCountdownLabel = computed(() => {
  if (!marketingDetail.value) return ''
  if (String(marketingDetail.value?.activity_type || '') !== 'seckill') return ''
  return countdownText.value ? `${t('productDetail.activityCountdown')}${countdownText.value}` : ''
})

const marketingCardStyle = computed(() => {
  const type = String(marketingDetail.value?.activity_type || '')
  if (type === 'seckill') return 'background: linear-gradient(135deg, #fef2f2, #fff7ed);'
  if (type === 'group_buy') return 'background: linear-gradient(135deg, #eff6ff, #f5f3ff);'
  if (type === 'bargain') return 'background: linear-gradient(135deg, #ecfdf5, #f0fdf4);'
  return 'background: #f8fafc;'
})
const marketingTextColor = computed(() => {
  const type = String(marketingDetail.value?.activity_type || '')
  if (type === 'seckill') return '#dc2626'
  if (type === 'group_buy') return '#2563eb'
  if (type === 'bargain') return '#16a34a'
  return '#334155'
})
const marketingSubTextColor = computed(() => {
  const type = String(marketingDetail.value?.activity_type || '')
  if (type === 'seckill') return '#b91c1c'
  if (type === 'group_buy') return '#1d4ed8'
  if (type === 'bargain') return '#15803d'
  return '#64748b'
})
const marketingActionBg = computed(() => {
  const type = String(marketingDetail.value?.activity_type || '')
  if (type === 'seckill') return '#dc2626'
  if (type === 'group_buy') return '#2563eb'
  if (type === 'bargain') return '#16a34a'
  return '#334155'
})

function parseAttrs(attrs: string) {
  try { return JSON.parse(attrs).map((a: any) => a.value).join(' ') }
  catch { return t('productDetail.defaultSpec') }
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
  activityProductID.value = Number(query.activity_product_id || 0)
  const data = await get<any>(`/api/v1/products/${id}`)
  if (!data) return
  product.value = data
  skus.value = data.skus || []
  if (skus.value.length) selectedSku.value = skus.value[0]
  const imgList: string[] = []
  if (data.cover) imgList.push(data.cover)
  if (data.images) imgList.push(...data.images.map((i: any) => i.url))
  images.value = imgList

  if (activityProductID.value > 0) {
    const detail = await get<any>(`/api/v1/marketing/activity-products/${activityProductID.value}`)
    if (detail?.activity_product_id) {
      marketingDetail.value = detail
      startCountdown()
    } else {
      marketingDetail.value = null
      activityProductID.value = 0
    }
  }
  await loadReviews(id)
})

onUnmounted(() => {
  if (countdownTimer) clearInterval(countdownTimer)
})

function startCountdown() {
  if (countdownTimer) clearInterval(countdownTimer)
  const tick = () => {
    const endAt = marketingDetail.value?.activity_end_at ? new Date(marketingDetail.value.activity_end_at).getTime() : 0
    if (!endAt) {
      countdownText.value = ''
      return
    }
    const diff = Math.max(0, endAt - Date.now())
    const h = Math.floor(diff / 3600000)
    const m = Math.floor((diff % 3600000) / 60000)
    const s = Math.floor((diff % 60000) / 1000)
    countdownText.value = `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
  }
  tick()
  countdownTimer = setInterval(tick, 1000)
}

async function addCart() {
  const skuID = selectedSku.value?.id
  if (!skuID) { uni.showToast({ title: t('productDetail.selectSpec'), icon: 'none' }); return }
  await post('/api/v1/cart/add', { sku_id: skuID, qty: 1, activity_product_id: activityProductID.value || 0 })
  uni.showToast({ title: t('productDetail.addedToCart'), icon: 'success' })
}

function buyNow() {
  const skuID = selectedSku.value?.id
  if (!skuID) { uni.showToast({ title: t('productDetail.selectSpec'), icon: 'none' }); return }
  const items = encodeURIComponent(JSON.stringify([{ sku_id: skuID, activity_product_id: activityProductID.value || 0 }]))
  uni.navigateTo({ url: `/pages/order/confirm?items=${items}&sku_ids=${skuID}` })
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
    uni.showToast({ title: t('productDetail.unfavorited'), icon: 'none' })
    return
  }
  await post(`/api/v1/products/${id}/favorite`)
  product.value.is_favorited = true
  product.value.favorite_count = Number(product.value.favorite_count || 0) + 1
  uni.showToast({ title: t('productDetail.favoriteSuccess'), icon: 'success' })
}
</script>

<style scoped>
.rich-text-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 12rpx;
}

.rich-text-content :deep(p) {
  margin-bottom: 16rpx;
}

.rich-text-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
}

.rich-text-content :deep(td),
.rich-text-content :deep(th) {
  border: 1px solid #e5e7eb;
  padding: 12rpx;
}
</style>
