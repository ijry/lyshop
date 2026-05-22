<template>
  <view>
    <template v-for="comp in components" :key="comp.id">
      <!-- Banner -->
      <u-swiper
        v-if="comp.type === 'banner'"
        :list="(comp.props?.images || []).map((i: any) => i.url || i)"
        :height="comp.props?.height || 350"
        radius="0"
      />

      <!-- Category nav -->
      <view v-else-if="comp.type === 'category_nav'" class="bg-white px-20rpx py-24rpx">
        <scroll-view scroll-x>
          <view class="flex gap-32rpx">
            <view v-for="item in comp.props?.items" :key="item.title"
              @click="item.link && uni.navigateTo({url: item.link})"
              class="flex flex-col items-center gap-8rpx flex-shrink-0">
              <view v-if="item.icon" class="w-88rpx h-88rpx rounded-20rpx overflow-hidden">
                <image :src="item.icon" mode="aspectFill" class="w-full h-full" />
              </view>
              <view v-else class="w-88rpx h-88rpx rounded-20rpx bg-blue-50 flex items-center justify-center">
                <text class="text-blue-700 text-24rpx font-500">{{ item.title?.slice(0,2) }}</text>
              </view>
              <text class="text-24rpx text-gray-700">{{ item.title }}</text>
            </view>
          </view>
        </scroll-view>
      </view>

      <!-- Product grid -->
      <view v-else-if="comp.type === 'product_grid'" class="px-20rpx py-16rpx">
        <view class="flex flex-wrap">
          <view v-for="p in gridProducts[comp.id]" :key="p.id"
            @click="uni.navigateTo({url:`/pages/product/detail?id=${p.id}`})"
            class="w-1/2 p-8rpx">
            <view class="bg-white rounded-20rpx overflow-hidden"
              style="box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);">
              <image :src="p.cover || ''" mode="aspectFill" class="w-full h-320rpx" />
              <view class="p-16rpx">
                <text class="text-26rpx text-gray-800 font-500 line-clamp-2">{{ p.title }}</text>
                <view class="flex items-center justify-between mt-12rpx">
                  <text class="text-30rpx text-red-500 font-700">¥{{ p.price }}</text>
                  <text class="text-20rpx text-gray-400">{{ p.sales }}人付款</text>
                </view>
              </view>
            </view>
          </view>
        </view>
      </view>

      <!-- Notice -->
      <view v-else-if="comp.type === 'notice'" class="bg-orange-50 px-24rpx py-16rpx">
        <u-notice-bar :text="comp.props?.text" />
      </view>

      <!-- Image ad -->
      <view v-else-if="comp.type === 'image_ad'" class="px-20rpx py-12rpx">
        <image v-if="comp.props?.url" :src="comp.props.url" mode="widthFix"
          class="w-full rounded-16rpx" @click="comp.props.link && uni.navigateTo({url:comp.props.link})" />
      </view>

      <!-- Rich text -->
      <view v-else-if="comp.type === 'rich_text'" class="px-30rpx py-16rpx">
        <rich-text :nodes="comp.props?.content || ''" class="text-28rpx text-gray-600" />
      </view>

      <!-- Spacer -->
      <view v-else-if="comp.type === 'spacer'"
        :style="{ height: (comp.props?.height || 16) + 'rpx', background: comp.props?.background || '#f8f8f8' }" />

      <!-- Marketing zone -->
      <view v-else-if="comp.type === 'marketing_zone'" class="px-20rpx py-12rpx">
        <view class="rounded-20rpx px-30rpx py-24rpx"
          style="background: linear-gradient(135deg, #dc2626 0%, #ea580c 100%);">
          <view class="flex items-center justify-between">
            <view>
              <text class="text-white text-32rpx font-700">限时秒杀</text>
              <text class="text-white text-22rpx opacity-80 ml-12rpx">抢购进行中</text>
            </view>
            <text class="text-white text-24rpx opacity-80">更多 ></text>
          </view>
        </view>
      </view>
    </template>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { get } from '@/utils/request'

const props = defineProps<{ components: any[] }>()
const gridProducts = ref<Record<string, any[]>>({})

async function loadGridProducts() {
  for (const comp of props.components) {
    if (comp.type === 'product_grid') {
      const params: any = { page: 1, size: comp.props?.limit || 10 }
      const data = await get<any>('/api/v1/products', params)
      gridProducts.value[comp.id] = data?.list || []
    }
  }
}

onMounted(loadGridProducts)
watch(() => props.components, loadGridProducts, { deep: true })
</script>
