<template>
  <view class="min-h-screen bg-gray-50">
    <!-- Header -->
    <view class="bg-white px-30rpx pt-60rpx pb-30rpx">
      <view class="flex items-center justify-between">
        <view class="flex items-center gap-12rpx">
          <image src="/static/lyshop-mark.svg" mode="aspectFit" class="w-48rpx h-48rpx" />
          <text class="text-36rpx font-700 text-gray-800">LYShop</text>
        </view>
        <view class="flex items-center gap-20rpx">
          <u-icon name="search" size="22" color="#666" @click="uni.switchTab({url:'/pages/product/list'})" />
          <u-icon name="chat" size="22" color="#666" @click="uni.navigateTo({url:'/pages/im/chat'})" />
        </view>
      </view>
    </view>

    <!-- Decor render -->
    <DecorRender :components="components" />

    <!-- Float customer service -->
    <view class="fixed right-30rpx bottom-180rpx z-50"
      @click="uni.navigateTo({url:'/pages/im/chat'})">
      <view class="w-96rpx h-96rpx rounded-full flex items-center justify-center"
        style="background: #dc2626; box-shadow: 0 4rpx 20rpx rgba(220,38,38,0.4);">
        <u-icon name="kefu-ermai" size="22" color="#fff" />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'
import DecorRender from '@/components/decor/DecorRender.vue'

const components = ref<any[]>([])

onMounted(async () => {
  const data = await get<any>('/api/v1/index/decor')
  if (data?.components) {
    if (Array.isArray(data.components)) {
      components.value = data.components
    } else if (typeof data.components === 'string') {
      try { components.value = JSON.parse(data.components) } catch {}
    }
  }
  if (!components.value.length) {
    components.value = [
      {
        type: 'banner',
        id: 'default_banner',
        props: {
          images: [
            { url: '/static/lyshop-wordmark.svg', link: '/pages/product/list' },
            { url: '/static/lyshop-mark.svg', link: '/pages/marketing/coupon?mode=claim' },
          ],
          height: 300
        }
      },
      {
        type: 'notice',
        id: 'default_notice',
        props: {
          items: [
            { text: '欢迎来到 LYShop', link: '/pages/index/index' },
            { text: '新人优惠券限时领取', link: '/pages/marketing/coupon?mode=claim' },
            { text: '热卖商品持续上新', link: '/pages/product/list' }
          ],
          color: '#f97316',
          bgColor: '#fff7ed',
          duration: 2500,
          mode: 'link'
        }
      },
      { type: 'product_grid', id: 'default_grid', props: { source: 'hot', limit: 10, columns: 2 } }
    ]
  }
})
</script>
