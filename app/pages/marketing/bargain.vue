<template>
  <view style="min-height: 100vh; background: #f5f5f5;">
    <u-navbar :title="$t('bargain.title')" :placeholder="true" />

    <!-- Header -->
    <view style="background: linear-gradient(135deg, #16a34a 0%, #22c55e 100%); padding: 20px 16px;">
      <text style="color: #fff; font-size: 18px; font-weight: 700;">{{ $t('bargain.subtitle') }}</text>
      <text style="color: rgba(255,255,255,0.7); font-size: 12px; display: block; margin-top: 4px;">{{ $t('bargain.description') }}</text>
    </view>

    <!-- Product list -->
    <view style="padding: 12px;">
      <view v-for="p in products" :key="`${p.activity_id}-${p.product_id}-${p.sku_id}`"
        @click="uni.navigateTo({url:`/pages/product/detail?id=${p.product_id}`})"
        style="background: #fff; border-radius: 12px; padding: 16px; margin-bottom: 10px; box-shadow: 0 1px 4px rgba(0,0,0,0.04);">
        <view style="display: flex;">
          <image :src="p.cover" mode="aspectFill" style="width: 90px; height: 90px; border-radius: 10px; flex-shrink: 0;" />
          <view style="flex: 1; margin-left: 12px;">
            <text style="font-size: 14px; color: #111; font-weight: 500; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;">{{ p.title }}</text>
            <view style="display: flex; align-items: baseline; gap: 6px; margin-top: 6px;">
              <text style="font-size: 12px; color: #999;">{{ $t('bargain.originalPrice') }}</text>
              <text style="font-size: 14px; color: #999; text-decoration: line-through;">¥{{ p.origin_price }}</text>
              <text style="font-size: 12px; color: #16a34a;">起砍 ¥{{ p.start_price || p.price }}</text>
              <text style="font-size: 12px; color: #16a34a;">{{ $t('bargain.floorPrice') }} ¥{{ p.floor_price }}</text>
            </view>
          </view>
        </view>
        <!-- Progress -->
        <view style="margin-top: 12px;">
          <view style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px;">
            <text style="font-size: 12px; color: #666;">已售 {{ p.sold_qty || 0 }}/{{ p.total_stock_limit || '-' }}</text>
            <text style="font-size: 12px; color: #16a34a; font-weight: 600;">{{ soldRatio(p) }}%</text>
          </view>
          <view style="height: 8px; background: #dcfce7; border-radius: 4px; overflow: hidden;">
            <view :style="{ width: soldRatio(p) + '%', height: '100%', background: 'linear-gradient(90deg, #16a34a, #22c55e)', borderRadius: '4px' }" />
          </view>
        </view>
        <view style="display: flex; justify-content: flex-end; margin-top: 10px;">
          <view style="background: #16a34a; color: #fff; font-size: 13px; padding: 6px 20px; border-radius: 14px; font-weight: 600;">
            {{ $t('bargain.helpBargain') }}
          </view>
        </view>
      </view>

      <view v-if="!products.length" style="text-align: center; padding: 60px 0; color: #999;">
        {{ $t('bargain.empty') }}
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'

const products = ref<any[]>([])

onMounted(async () => {
  const data = await get<any>('/api/v1/marketing/bargain/products')
  products.value = Array.isArray(data?.list) ? data.list : (Array.isArray(data) ? data : [])
})

function soldRatio(item: any) {
  const total = Number(item?.total_stock_limit || 0)
  const sold = Number(item?.sold_qty || 0)
  if (total <= 0) return 0
  return Math.max(0, Math.min(100, Math.round((sold / total) * 100)))
}
</script>
