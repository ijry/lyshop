<template>
  <view class="min-h-screen pb-safe" style="background: #f5f5f5;">
    <u-navbar :title="$t('bargain.orderDetail')" :placeholder="true" />

    <!-- 砍价状态卡片 -->
    <view class="m-24rpx p-32rpx bg-white rounded-24rpx">
      <view class="flex items-center justify-between mb-24rpx">
        <text class="text-32rpx font-700">{{ statusLabel }}</text>
        <view :class="statusClass" class="px-20rpx py-8rpx rounded-999rpx text-24rpx">
          {{ statusText }}
        </view>
      </view>

      <!-- 倒计时 -->
      <view v-if="bargainOrder.status === 'pending'" class="bg-gradient-to-r from-red-50 to-pink-100 rounded-16rpx p-24rpx mb-24rpx">
        <view class="flex items-center justify-between">
          <text class="text-26rpx text-gray-700">{{ $t('bargain.timeRemaining') }}</text>
          <view class="flex items-center gap-8rpx">
            <view v-for="(t, i) in countdownParts" :key="i"
              class="bg-white text-red-600 text-28rpx font-700 px-12rpx py-6rpx rounded-8rpx min-w-56rpx text-center">
              {{ t }}
            </view>
          </view>
        </view>
      </view>

      <!-- 价格进度 -->
      <view class="mb-24rpx">
        <view class="flex items-center justify-between mb-16rpx">
          <view>
            <text class="text-24rpx text-gray-500 block">{{ $t('bargain.currentPrice') }}</text>
            <text class="text-48rpx font-700 text-red-500 block mt-4rpx">¥{{ bargainOrder.current_price?.toFixed(2) }}</text>
          </view>
          <view class="text-right">
            <text class="text-24rpx text-gray-500 block">{{ $t('bargain.targetPrice') }}</text>
            <text class="text-32rpx font-600 text-gray-700 block mt-4rpx">¥{{ bargainOrder.floor_price?.toFixed(2) }}</text>
          </view>
        </view>

        <view class="h-20rpx bg-gray-100 rounded-999rpx overflow-hidden relative">
          <view class="h-full bg-gradient-to-r from-red-500 to-pink-500 rounded-999rpx transition-all"
            :style="{ width: progressPercent + '%' }" />
          <text class="absolute top-0 left-0 right-0 text-center text-20rpx text-white leading-20rpx font-600">
            {{ progressPercent }}%
          </text>
        </view>

        <view class="flex items-center justify-between mt-12rpx">
          <text class="text-22rpx text-gray-400">{{ $t('bargain.saved') }} ¥{{ savedAmount.toFixed(2) }}</text>
          <text class="text-22rpx text-gray-400">{{ $t('bargain.remaining') }} ¥{{ remainingAmount.toFixed(2) }}</text>
        </view>
      </view>

      <!-- 商品信息 -->
      <view class="flex gap-20rpx p-20rpx bg-gray-50 rounded-16rpx">
        <image :src="product.cover" mode="aspectFill" class="w-120rpx h-120rpx rounded-12rpx flex-shrink-0" />
        <view class="flex-1">
          <text class="text-28rpx font-600 text-gray-800 line-clamp-2">{{ product.title }}</text>
          <view class="flex items-baseline gap-12rpx mt-12rpx">
            <text class="text-24rpx text-gray-500">{{ $t('bargain.startPrice') }}</text>
            <text class="text-28rpx font-600 text-gray-600">¥{{ bargainOrder.start_price?.toFixed(2) }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 助力记录 -->
    <view class="m-24rpx p-32rpx bg-white rounded-24rpx">
      <view class="flex items-center justify-between mb-24rpx">
        <text class="text-28rpx font-700 text-gray-800">{{ $t('bargain.helpRecords') }}</text>
        <text class="text-24rpx text-gray-500">{{ helpers.length }} {{ $t('bargain.people') }}</text>
      </view>

      <view v-if="helpers.length" class="space-y-24rpx">
        <view v-for="helper in helpers" :key="helper.id" class="flex items-center justify-between">
          <view class="flex items-center gap-20rpx">
            <view class="w-80rpx h-80rpx rounded-full bg-gradient-to-br from-red-400 to-pink-600 flex items-center justify-center">
              <text class="text-32rpx text-white">{{ helper.user_nickname?.charAt(0) || '?' }}</text>
            </view>
            <view>
              <text class="text-26rpx text-gray-800 block">{{ helper.user_nickname || $t('bargain.anonymous') }}</text>
              <text class="text-22rpx text-gray-400 block mt-4rpx">{{ formatTime(helper.created_at) }}</text>
            </view>
          </view>
          <view class="text-right">
            <text class="text-28rpx font-700 text-red-500 block">-¥{{ helper.cut_amount?.toFixed(2) }}</text>
            <text class="text-20rpx text-gray-400 block mt-4rpx">{{ $t('bargain.helped') }}</text>
          </view>
        </view>
      </view>

      <view v-else class="text-center py-60rpx">
        <u-icon name="empty-list" size="60" color="#ccc" />
        <text class="text-24rpx text-gray-400 block mt-16rpx">{{ $t('bargain.noHelpersYet') }}</text>
      </view>
    </view>

    <!-- 底部操作栏 -->
    <view class="fixed bottom-0 left-0 right-0 z-100 flex items-center gap-20rpx px-24rpx py-20rpx bg-white"
      style="padding-bottom: calc(20rpx + env(safe-area-inset-bottom)); border-top: 1px solid #f0f0f0;">
      <u-button v-if="bargainOrder.status === 'pending'" type="primary" :text="$t('bargain.inviteHelp')"
        @click="shareBargain" class="flex-1" shape="circle" />
      <u-button v-else-if="bargainOrder.status === 'success'" type="success" :text="$t('bargain.buyNow')"
        @click="buyNow" class="flex-1" shape="circle" />
      <u-button v-else type="default" :text="$t('bargain.backToList')"
        @click="backToList" class="flex-1" shape="circle" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { get } from '@/utils/request'

const { t } = useI18n()

const bargainOrder = ref<any>({})
const helpers = ref<any[]>([])
const product = ref<any>({})
const countdownParts = ref(['00', '00', '00'])
let timer: any = null

const progressPercent = computed(() => {
  const start = bargainOrder.value.start_price || 0
  const floor = bargainOrder.value.floor_price || 0
  const current = bargainOrder.value.current_price || start
  if (start <= floor) return 100
  return Math.min(100, Math.round(((start - current) / (start - floor)) * 100))
})

const savedAmount = computed(() => {
  const start = bargainOrder.value.start_price || 0
  const current = bargainOrder.value.current_price || start
  return Math.max(0, start - current)
})

const remainingAmount = computed(() => {
  const current = bargainOrder.value.current_price || 0
  const floor = bargainOrder.value.floor_price || 0
  return Math.max(0, current - floor)
})

const statusLabel = computed(() => {
  const status = bargainOrder.value.status
  if (status === 'success') return t('bargain.bargainSuccess')
  if (status === 'failed') return t('bargain.bargainFailed')
  return t('bargain.bargaining')
})

const statusText = computed(() => {
  const status = bargainOrder.value.status
  if (status === 'success') return t('bargain.completed')
  if (status === 'failed') return t('bargain.expired')
  return t('bargain.inProgress')
})

const statusClass = computed(() => {
  const status = bargainOrder.value.status
  if (status === 'success') return 'bg-green-100 text-green-600'
  if (status === 'failed') return 'bg-gray-100 text-gray-600'
  return 'bg-red-100 text-red-600'
})

function formatTime(time: string) {
  if (!time) return ''
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return t('bargain.justNow')
  if (minutes < 60) return t('bargain.minutesAgo', { n: minutes })
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return t('bargain.hoursAgo', { n: hours })
  return date.toLocaleDateString()
}

function startCountdown() {
  if (!bargainOrder.value.expire_at) return

  timer = setInterval(() => {
    const expireTime = new Date(bargainOrder.value.expire_at).getTime()
    const now = Date.now()
    const diff = Math.max(0, expireTime - now)

    if (diff === 0) {
      clearInterval(timer)
      bargainOrder.value.status = 'failed'
      return
    }

    const h = Math.floor(diff / 3600000)
    const m = Math.floor((diff % 3600000) / 60000)
    const s = Math.floor((diff % 60000) / 1000)
    countdownParts.value = [
      String(h).padStart(2, '0'),
      String(m).padStart(2, '0'),
      String(s).padStart(2, '0')
    ]
  }, 1000)
}

async function loadData() {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  const id = Number(query.id)

  if (!id) return

  // 加载砍价订单
  const orderData = await get<any>(`/api/v1/bargain/orders/${id}`)
  if (orderData) {
    bargainOrder.value = orderData
    if (orderData.status === 'pending') {
      startCountdown()
    }
  }

  // 加载助力记录
  const helpersData = await get<any>(`/api/v1/bargain/orders/${id}/helpers`)
  if (Array.isArray(helpersData)) {
    helpers.value = helpersData
  }

  // 加载商品信息
  if (orderData?.product_id) {
    const productData = await get<any>(`/api/v1/products/${orderData.product_id}`)
    if (productData) {
      product.value = productData
    }
  }
}

function shareBargain() {
  uni.showToast({ title: t('bargain.shareFeatureComingSoon'), icon: 'none' })
}

function buyNow() {
  const skuID = bargainOrder.value.sku_id
  if (!skuID) return
  const items = encodeURIComponent(JSON.stringify([{ sku_id: skuID }]))
  uni.navigateTo({ url: `/pages/order/confirm?items=${items}&sku_ids=${skuID}` })
}

function backToList() {
  uni.navigateBack()
}

onMounted(() => {
  loadData()
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.space-y-24rpx > view + view {
  margin-top: 24rpx;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
