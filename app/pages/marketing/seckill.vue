<template>
  <view class="bg-gray-50 min-h-screen">
    <u-navbar title="限时秒杀" :placeholder="true" />

    <view v-if="!activities.length" class="text-center py-12 text-gray-400">
      <text>暂无进行中的秒杀活动</text>
    </view>

    <view class="p-3">
      <view v-for="act in activities" :key="act.id"
        class="bg-white rounded-xl mb-4 overflow-hidden shadow-sm">
        <!-- Activity header -->
        <view class="bg-gradient-to-r from-red-600 to-orange-500 px-4 py-3">
          <text class="text-white font-bold text-base block">{{ act.name }}</text>
          <text class="text-red-100 text-xs">{{ act.start_at?.slice(0,16) }} ~ {{ act.end_at?.slice(0,16) }}</text>
        </view>
        <!-- Products -->
        <view v-if="act.products?.length" class="p-3">
          <scroll-view scroll-x>
            <view class="flex gap-3">
              <view v-for="ap in act.products" :key="ap.id"
                @click="uni.navigateTo({url:`/pages/product/detail?id=${ap.product_id}`})"
                class="w-32 shrink-0">
                <view class="bg-gray-100 rounded-xl" style="height:128px;" />
                <view class="mt-2 px-1">
                  <text class="text-red-600 font-bold block">¥{{ ap.activity_price }}</text>
                  <text class="text-gray-400 text-xs block">库存 {{ ap.activity_stock }}</text>
                </view>
              </view>
            </view>
          </scroll-view>
        </view>
        <view v-else class="px-4 pb-4 text-gray-400 text-sm">活动商品加载中...</view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'

const activities = ref<any[]>([])

onMounted(async () => {
  const data = await get<any[]>('/api/v1/marketing/seckills')
  activities.value = data || []
})
</script>
