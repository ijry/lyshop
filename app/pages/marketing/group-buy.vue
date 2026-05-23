<template>
  <view style="min-height: 100vh; background: #f5f5f5;">
    <u-navbar title="拼团特惠" :placeholder="true" />

    <!-- Header -->
    <view style="background: linear-gradient(135deg, #2563eb 0%, #3b82f6 100%); padding: 20px 16px;">
      <text style="color: #fff; font-size: 18px; font-weight: 700;">多人拼团 更便宜</text>
      <text style="color: rgba(255,255,255,0.7); font-size: 12px; display: block; margin-top: 4px;">邀请好友一起拼，低价带走好商品</text>
    </view>

    <!-- Product list -->
    <view style="padding: 12px;">
      <view v-for="p in products" :key="p.product_id"
        @click="uni.navigateTo({url:`/pages/product/detail?id=${p.product_id}`})"
        style="display: flex; background: #fff; border-radius: 12px; padding: 12px; margin-bottom: 10px; box-shadow: 0 1px 4px rgba(0,0,0,0.04);">
        <image :src="p.cover" mode="aspectFill" style="width: 100px; height: 100px; border-radius: 10px; flex-shrink: 0;" />
        <view style="flex: 1; margin-left: 12px; display: flex; flex-direction: column; justify-content: space-between;">
          <text style="font-size: 14px; color: #111; font-weight: 500; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;">{{ p.title }}</text>
          <view>
            <view style="display: flex; align-items: baseline; gap: 6px;">
              <text style="font-size: 12px; color: #2563eb; background: #eff6ff; padding: 1px 6px; border-radius: 4px;">{{ groupSize }}人团</text>
              <text style="font-size: 20px; color: #dc2626; font-weight: 700;">¥{{ p.group_price }}</text>
              <text style="font-size: 12px; color: #999; text-decoration: line-through;">¥{{ p.origin_price }}</text>
            </view>
            <view style="display: flex; align-items: center; justify-content: space-between; margin-top: 8px;">
              <text style="font-size: 11px; color: #999;">已团{{ Math.floor(Math.random()*200+50) }}件</text>
              <view style="background: #2563eb; color: #fff; font-size: 12px; padding: 5px 16px; border-radius: 14px; font-weight: 600;">
                去拼团
              </view>
            </view>
          </view>
        </view>
      </view>

      <view v-if="!products.length" style="text-align: center; padding: 60px 0; color: #999;">
        暂无拼团活动
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'

const products = ref<any[]>([])
const groupSize = ref(3)

onMounted(async () => {
  const data = await get<any>('/api/v1/marketing/group-buy')
  const activity = data?.list?.[0]
  if (activity) {
    products.value = activity.products || []
    groupSize.value = activity.group_size || 3
  }
})
</script>
