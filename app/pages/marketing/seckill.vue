<template>
  <view style="min-height: 100vh; background: #f5f5f5;">
    <u-navbar :title="$t('seckill.title')" :placeholder="true" />

    <!-- Countdown header -->
    <view style="background: linear-gradient(135deg, #dc2626 0%, #ea580c 100%); padding: 20px 16px;">
      <view style="display: flex; align-items: center; justify-content: space-between;">
        <text style="color: #fff; font-size: 18px; font-weight: 700;">{{ $t('seckill.title') }}</text>
        <view style="display: flex; align-items: center; gap: 4px;">
          <text style="color: rgba(255,255,255,0.7); font-size: 12px;">{{ $t('seckill.endsIn') }}</text>
          <view v-for="(t, i) in countdownParts" :key="i"
            style="background: rgba(0,0,0,0.3); color: #fff; font-size: 14px; font-weight: 700; padding: 2px 6px; border-radius: 4px; min-width: 28px; text-align: center;">
            {{ t }}
          </view>
        </view>
      </view>
    </view>

    <!-- Product list -->
    <view style="padding: 12px;">
      <view v-for="p in products" :key="`${p.activity_id}-${p.product_id}-${p.sku_id}`"
        @click="openDetail(p)"
        style="display: flex; background: #fff; border-radius: 12px; padding: 12px; margin-bottom: 10px; box-shadow: 0 1px 4px rgba(0,0,0,0.04);">
        <image :src="p.cover" mode="aspectFill" style="width: 100px; height: 100px; border-radius: 10px; flex-shrink: 0;" />
        <view style="flex: 1; margin-left: 12px; display: flex; flex-direction: column; justify-content: space-between;">
          <text style="font-size: 14px; color: #111; font-weight: 500; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;">{{ p.title }}</text>
          <view>
            <view style="display: flex; align-items: baseline; gap: 6px;">
              <text style="font-size: 20px; color: #dc2626; font-weight: 700;">¥{{ p.activity_price }}</text>
              <text style="font-size: 12px; color: #999; text-decoration: line-through;">¥{{ p.origin_price }}</text>
            </view>
            <view style="display: flex; align-items: center; justify-content: space-between; margin-top: 6px;">
              <!-- Progress bar -->
              <view style="flex: 1; margin-right: 10px;">
                <view style="height: 14px; background: #fee2e2; border-radius: 7px; overflow: hidden; position: relative;">
                  <view :style="{ width: soldRatio(p) + '%', height: '100%', background: 'linear-gradient(90deg, #dc2626, #ef4444)', borderRadius: '7px' }" />
                  <text style="position: absolute; top: 0; left: 0; right: 0; text-align: center; font-size: 10px; color: #fff; line-height: 14px;">{{ soldRatio(p) }}%</text>
                </view>
              </view>
              <view style="background: #dc2626; color: #fff; font-size: 12px; padding: 4px 12px; border-radius: 12px; font-weight: 600;">
                {{ $t('seckill.buyNow') }}
              </view>
            </view>
          </view>
        </view>
      </view>

      <view v-if="!products.length" style="text-align: center; padding: 60px 0; color: #999;">
        {{ $t('seckill.empty') }}
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { get } from '@/utils/request'

const products = ref<any[]>([])
const countdownParts = ref(['00', '00', '00'])
let timer: any = null

onMounted(async () => {
  const data = await get<any>('/api/v1/seckill/products')
  products.value = Array.isArray(data?.list) ? data.list : (Array.isArray(data) ? data : [])

  // Countdown
  const endTimes = products.value
    .map((item: any) => item?.activity_end_at ? new Date(item.activity_end_at).getTime() : 0)
    .filter((ts: number) => ts > Date.now())
  const endAt = endTimes.length ? Math.min(...endTimes) : Date.now() + 86400000
  timer = setInterval(() => {
    const diff = Math.max(0, endAt - Date.now())
    const h = Math.floor(diff / 3600000)
    const m = Math.floor((diff % 3600000) / 60000)
    const s = Math.floor((diff % 60000) / 1000)
    countdownParts.value = [String(h).padStart(2,'0'), String(m).padStart(2,'0'), String(s).padStart(2,'0')]
  }, 1000)
})

onUnmounted(() => timer && clearInterval(timer))

function soldRatio(item: any) {
  const total = Number(item?.total_stock_limit || 0)
  const sold = Number(item?.sold_qty || 0)
  if (total <= 0) return 0
  return Math.max(0, Math.min(100, Math.round((sold / total) * 100)))
}

function openDetail(item: any) {
  const productID = Number(item?.product_id || 0)
  if (!productID) return
  const activityProductID = Number(item?.activity_product_id || 0)
  const url = activityProductID > 0
    ? `/pages/product/detail?id=${productID}&activity_product_id=${activityProductID}`
    : `/pages/product/detail?id=${productID}`
  uni.navigateTo({ url })
}
</script>
