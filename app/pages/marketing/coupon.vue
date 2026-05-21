<template>
  <view class="bg-gray-50 min-h-screen">
    <u-navbar title="我的优惠券" :placeholder="true" />

    <!-- Tabs -->
    <u-tabs :list="tabs" :current="activeTab" @click="onTab" />

    <view class="p-3">
      <view v-if="!coupons.length" class="text-center py-12 text-gray-400">
        <text class="block mb-3">暂无优惠券</text>
        <u-button text="去领取" size="small" type="primary"
          @click="uni.navigateTo({url:'/pages/marketing/coupon_center'})" />
      </view>

      <view v-for="c in coupons" :key="c.id" class="bg-white rounded-xl mb-3 overflow-hidden shadow-sm">
        <view class="flex">
          <!-- Left: coupon value -->
          <view class="bg-blue-700 w-28 flex flex-col items-center justify-center py-5 shrink-0">
            <text class="text-white text-2xl font-bold">
              {{ c.coupon?.type === 2 ? c.coupon?.discount * 10 + '折' : '¥' + c.coupon?.discount }}
            </text>
            <text class="text-blue-200 text-xs mt-1">
              {{ c.coupon?.min_amount > 0 ? '满' + c.coupon?.min_amount + '可用' : '无门槛' }}
            </text>
          </view>
          <!-- Right: info -->
          <view class="flex-1 p-4">
            <text class="text-slate-800 font-medium block">{{ c.coupon?.name || '优惠券' }}</text>
            <text class="text-gray-400 text-xs mt-1 block">有效期至：{{ c.coupon?.end_at?.slice(0,10) || '长期有效' }}</text>
            <view class="mt-3">
              <span :class="c.status===1 ? 'bg-blue-50 text-blue-700' : 'bg-gray-100 text-gray-400'"
                class="px-2 py-1 rounded-full text-xs">
                {{ statusLabel(c.status) }}
              </span>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'

const coupons = ref<any[]>([])
const activeTab = ref(0)
const tabs = [{ name: '未使用' }, { name: '已使用' }, { name: '已过期' }]
const statusMap = [1, 2, 3]
const statusLabels: Record<number, string> = { 1: '未使用', 2: '已使用', 3: '已过期' }
const statusLabel = (s: number) => statusLabels[s] || ''

async function loadCoupons() {
  const data = await get<any[]>('/api/v1/user/coupons')
  const allStatus = statusMap[activeTab.value]
  coupons.value = (data || []).filter((c: any) => c.status === allStatus)
}

function onTab(index: number) {
  activeTab.value = index
  loadCoupons()
}

onMounted(loadCoupons)
</script>
