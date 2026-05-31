<template>
  <view class="min-h-screen pb-safe" style="background: #f5f5f5;">
    <u-navbar :title="$t('groupBuy.orderDetail')" :placeholder="true" />

    <!-- 拼团状态卡片 -->
    <view class="m-24rpx p-32rpx bg-white rounded-24rpx">
      <view class="flex items-center justify-between mb-24rpx">
        <text class="text-32rpx font-700">{{ statusLabel }}</text>
        <view :class="statusClass" class="px-20rpx py-8rpx rounded-999rpx text-24rpx">
          {{ statusText }}
        </view>
      </view>

      <!-- 倒计时 -->
      <view v-if="groupOrder.status === 'pending'" class="bg-gradient-to-r from-blue-50 to-blue-100 rounded-16rpx p-24rpx mb-24rpx">
        <view class="flex items-center justify-between">
          <text class="text-26rpx text-gray-700">{{ $t('groupBuy.timeRemaining') }}</text>
          <view class="flex items-center gap-8rpx">
            <view v-for="(t, i) in countdownParts" :key="i"
              class="bg-white text-blue-600 text-28rpx font-700 px-12rpx py-6rpx rounded-8rpx min-w-56rpx text-center">
              {{ t }}
            </view>
          </view>
        </view>
      </view>

      <!-- 进度条 -->
      <view class="mb-24rpx">
        <view class="flex items-center justify-between mb-12rpx">
          <text class="text-26rpx text-gray-600">{{ $t('groupBuy.progress') }}</text>
          <text class="text-26rpx text-blue-600 font-600">{{ groupOrder.joined_count }}/{{ groupOrder.group_size }}</text>
        </view>
        <view class="h-16rpx bg-gray-100 rounded-999rpx overflow-hidden">
          <view class="h-full bg-gradient-to-r from-blue-500 to-blue-600 rounded-999rpx transition-all"
            :style="{ width: progressPercent + '%' }" />
        </view>
      </view>

      <!-- 商品信息 -->
      <view class="flex gap-20rpx p-20rpx bg-gray-50 rounded-16rpx">
        <image :src="product.cover" mode="aspectFill" class="w-120rpx h-120rpx rounded-12rpx flex-shrink-0" />
        <view class="flex-1">
          <text class="text-28rpx font-600 text-gray-800 line-clamp-2">{{ product.title }}</text>
          <view class="flex items-baseline gap-12rpx mt-12rpx">
            <text class="text-36rpx font-700 text-red-500">¥{{ groupOrder.group_price || product.price }}</text>
            <text class="text-24rpx text-gray-400 line-through">¥{{ product.origin_price }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 成员列表 -->
    <view class="m-24rpx p-32rpx bg-white rounded-24rpx">
      <text class="text-28rpx font-700 text-gray-800 block mb-24rpx">{{ $t('groupBuy.members') }}</text>

      <view v-if="members.length" class="space-y-20rpx">
        <view v-for="member in members" :key="member.id" class="flex items-center justify-between">
          <view class="flex items-center gap-20rpx">
            <view class="w-80rpx h-80rpx rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center">
              <text class="text-32rpx text-white">{{ member.user_nickname?.charAt(0) || '?' }}</text>
            </view>
            <view>
              <text class="text-26rpx text-gray-800 block">{{ member.user_nickname || $t('groupBuy.anonymous') }}</text>
              <text class="text-22rpx text-gray-400 block mt-4rpx">{{ formatTime(member.created_at) }}</text>
            </view>
          </view>
          <view v-if="member.is_leader" class="bg-gradient-to-r from-yellow-400 to-orange-400 text-white text-22rpx px-16rpx py-6rpx rounded-999rpx">
            {{ $t('groupBuy.leader') }}
          </view>
        </view>
      </view>

      <!-- 空位占位符 -->
      <view v-if="groupOrder.status === 'pending' && emptySlots > 0" class="mt-20rpx space-y-20rpx">
        <view v-for="i in emptySlots" :key="'empty-' + i" class="flex items-center gap-20rpx opacity-40">
          <view class="w-80rpx h-80rpx rounded-full border-2 border-dashed border-gray-300 flex items-center justify-center">
            <u-icon name="plus" size="24" color="#999" />
          </view>
          <text class="text-26rpx text-gray-400">{{ $t('groupBuy.waitingForJoin') }}</text>
        </view>
      </view>
    </view>

    <!-- 底部操作栏 -->
    <view class="fixed bottom-0 left-0 right-0 z-100 flex items-center gap-20rpx px-24rpx py-20rpx bg-white"
      style="padding-bottom: calc(20rpx + env(safe-area-inset-bottom)); border-top: 1px solid #f0f0f0;">
      <u-button v-if="groupOrder.status === 'pending'" type="primary" :text="$t('groupBuy.inviteFriends')"
        @click="shareGroup" class="flex-1" shape="circle" />
      <u-button v-else-if="groupOrder.status === 'success'" type="success" :text="$t('groupBuy.viewOrder')"
        @click="viewOrder" class="flex-1" shape="circle" />
      <u-button v-else type="default" :text="$t('groupBuy.backToList')"
        @click="backToList" class="flex-1" shape="circle" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { get } from '@/utils/request'

const { t } = useI18n()

const groupOrder = ref<any>({})
const members = ref<any[]>([])
const product = ref<any>({})
const countdownParts = ref(['00', '00', '00'])
let timer: any = null

const progressPercent = computed(() => {
  const total = groupOrder.value.group_size || 1
  const joined = groupOrder.value.joined_count || 0
  return Math.min(100, Math.round((joined / total) * 100))
})

const emptySlots = computed(() => {
  return Math.max(0, (groupOrder.value.group_size || 0) - (groupOrder.value.joined_count || 0))
})

const statusLabel = computed(() => {
  const status = groupOrder.value.status
  if (status === 'success') return t('groupBuy.groupSuccess')
  if (status === 'failed') return t('groupBuy.groupFailed')
  return t('groupBuy.grouping')
})

const statusText = computed(() => {
  const status = groupOrder.value.status
  if (status === 'success') return t('groupBuy.completed')
  if (status === 'failed') return t('groupBuy.expired')
  return t('groupBuy.inProgress')
})

const statusClass = computed(() => {
  const status = groupOrder.value.status
  if (status === 'success') return 'bg-green-100 text-green-600'
  if (status === 'failed') return 'bg-gray-100 text-gray-600'
  return 'bg-blue-100 text-blue-600'
})

function formatTime(time: string) {
  if (!time) return ''
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return t('groupBuy.justNow')
  if (minutes < 60) return t('groupBuy.minutesAgo', { n: minutes })
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return t('groupBuy.hoursAgo', { n: hours })
  return date.toLocaleDateString()
}

function startCountdown() {
  if (!groupOrder.value.expire_at) return

  timer = setInterval(() => {
    const expireTime = new Date(groupOrder.value.expire_at).getTime()
    const now = Date.now()
    const diff = Math.max(0, expireTime - now)

    if (diff === 0) {
      clearInterval(timer)
      groupOrder.value.status = 'failed'
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

  // 加载拼团订单
  const orderData = await get<any>(`/api/v1/group-buy/orders/${id}`)
  if (orderData) {
    groupOrder.value = orderData
    if (orderData.status === 'pending') {
      startCountdown()
    }
  }

  // 加载成员列表
  const membersData = await get<any>(`/api/v1/group-buy/orders/${id}/members`)
  if (Array.isArray(membersData)) {
    members.value = membersData
  }

  // 加载商品信息（简化版，实际应该从订单中获取）
  if (orderData?.product_id) {
    const productData = await get<any>(`/api/v1/products/${orderData.product_id}`)
    if (productData) {
      product.value = productData
    }
  }
}

function shareGroup() {
  uni.showToast({ title: t('groupBuy.shareFeatureComingSoon'), icon: 'none' })
}

function viewOrder() {
  uni.navigateTo({ url: '/pages/order/list' })
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
.space-y-20rpx > view + view {
  margin-top: 20rpx;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
