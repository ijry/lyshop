<template>
  <view class="bg-gray-50 min-h-screen">
    <u-navbar title="我的优惠券" :placeholder="true" />

    <view v-if="claimMode" class="px-20rpx pt-12rpx">
      <view class="bg-white rounded-24rpx p-24rpx mb-20rpx">
        <view class="flex items-center justify-between mb-16rpx">
          <text class="text-30rpx font-700 text-gray-800">可领取优惠券</text>
          <text class="text-22rpx text-gray-400">领券后自动进入未使用</text>
        </view>
        <view v-if="!claimCoupons.length" class="text-24rpx text-gray-400 py-12rpx">暂无可领取优惠券</view>
        <view v-for="item in claimCoupons" :key="item.id" class="border border-gray-100 rounded-16rpx p-20rpx mb-16rpx">
          <view class="flex items-center justify-between">
            <view class="flex-1 pr-20rpx">
              <text class="text-28rpx font-600 text-gray-800 block">{{ item.name }}</text>
              <text class="text-24rpx text-gray-500 block mt-6rpx">{{ claimDesc(item) }}</text>
              <text class="text-22rpx text-gray-400 block mt-6rpx">
                每人限领{{ item.per_limit > 0 ? item.per_limit : '不限' }}张 · 已领{{ item.claimed_by_me || 0 }}张
              </text>
            </view>
            <u-button
              :text="item.can_claim ? '立即领取' : '已达上限'"
              size="mini"
              type="primary"
              :disabled="!item.can_claim || isClaiming(item.id)"
              @click="onClaim(item.id)"
            />
          </view>
        </view>
      </view>
    </view>

    <!-- Tabs -->
    <u-tabs :list="tabs" :current="activeTab" @click="onTab" />

    <view class="p-3">
      <view v-if="!coupons.length" class="text-center py-12 text-gray-400">
        <text class="block mb-3">暂无优惠券</text>
        <u-button text="去领取" size="small" type="primary" @click="goClaimCoupon" />
      </view>

      <view v-for="c in coupons" :key="c.id" class="bg-white rounded-24rpx mb-24rpx overflow-hidden shadow-sm">
        <view class="flex">
          <!-- Left: coupon value -->
          <view class="bg-blue-700 w-180rpx flex flex-col items-center justify-center py-32rpx flex-shrink-0">
            <text class="text-white text-48rpx font-700 leading-56rpx">
              {{ c.coupon?.type === 2 ? c.coupon?.discount * 10 + '折' : '¥' + c.coupon?.discount }}
            </text>
            <text class="text-blue-200 text-22rpx mt-8rpx leading-30rpx">
              {{ c.coupon?.min_amount > 0 ? '满' + c.coupon?.min_amount + '可用' : '无门槛' }}
            </text>
          </view>
          <!-- Right: info -->
          <view class="flex-1 p-24rpx">
            <text class="text-slate-800 text-30rpx font-600 block">{{ c.coupon?.name || '优惠券' }}</text>
            <text class="text-gray-400 text-24rpx mt-8rpx block">有效期至：{{ c.coupon?.end_at?.slice(0,10) || '长期有效' }}</text>
            <view class="mt-20rpx">
              <text :class="c.status===1 ? 'bg-blue-50 text-blue-700' : 'bg-gray-100 text-gray-400'"
                class="px-16rpx py-8rpx rounded-full text-22rpx inline-block">
                {{ statusLabel(c.status) }}
              </text>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get, post } from '@/utils/request'

const coupons = ref<any[]>([])
const claimCoupons = ref<any[]>([])
const claimMode = ref(false)
const claimingIds = ref<number[]>([])
const activeTab = ref(0)
const tabs = [{ name: '未使用' }, { name: '已使用' }, { name: '已过期' }]
const statusMap = [1, 2, 3]
const statusLabels: Record<number, string> = { 1: '未使用', 2: '已使用', 3: '已过期' }
const statusLabel = (s: number) => statusLabels[s] || ''
const isClaiming = (id: number) => claimingIds.value.includes(Number(id))

function claimDesc(coupon: any) {
  if (coupon?.type === 2) return `${Number(coupon.discount || 0) * 10}折券`
  if (coupon?.type === 3) return '无门槛券'
  return `满${coupon?.min_amount || 0}减${coupon?.discount || 0}`
}

function detectClaimMode() {
  const pages = getCurrentPages() as any[]
  const current = pages[pages.length - 1] as any
  claimMode.value = current?.options?.mode === 'claim'
}

async function loadCoupons() {
  const data = await get<any[]>('/api/v1/user/coupons')
  const allStatus = statusMap[activeTab.value]
  coupons.value = (data || []).filter((c: any) => c.status === allStatus)
}

async function loadClaimCoupons() {
  const data = await get<any[]>('/api/v1/coupons')
  claimCoupons.value = data || []
}

async function onClaim(id: number) {
  const couponID = Number(id)
  if (!couponID || isClaiming(couponID)) return
  claimingIds.value = [...claimingIds.value, couponID]
  try {
    await post(`/api/v1/coupons/${couponID}/claim`)
    uni.showToast({ title: '领取成功', icon: 'success' })
    await Promise.all([loadClaimCoupons(), loadCoupons()])
  } finally {
    claimingIds.value = claimingIds.value.filter((item) => item !== couponID)
  }
}

function onTab(event: number | { index?: number; current?: number }) {
  const index = typeof event === 'number' ? event : Number(event?.index ?? event?.current ?? 0)
  activeTab.value = index
  loadCoupons()
}

function goClaimCoupon() {
  if (claimMode.value) {
    loadClaimCoupons()
    return
  }
  uni.navigateTo({ url: '/pages/marketing/coupon?mode=claim' })
}

onMounted(async () => {
  detectClaimMode()
  await Promise.all([loadCoupons(), claimMode.value ? loadClaimCoupons() : Promise.resolve()])
})
</script>
