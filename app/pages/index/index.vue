<template>
  <view class="min-h-screen" style="background: var(--app-page-bg);">
    <!-- Header -->
    <view class="px-30rpx pt-60rpx pb-30rpx" style="background: var(--app-card-bg);">
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
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { get } from '@/utils/request'
import DecorRender from '@/components/decor/DecorRender.vue'

const { t } = useI18n()

const components = ref<any[]>([])
const isPreview = ref(false)

function onPreviewMessage(e: MessageEvent) {
  if (e.data?.type === 'DECOR_PREVIEW_UPDATE' && e.data?.source === 'lyshop-admin') {
    components.value = e.data.components || []
  }
}

onMounted(async () => {
  // Check if in admin preview iframe mode
  const params = new URLSearchParams(window.location.search)
  if (params.get('preview') === '1') {
    isPreview.value = true
    window.addEventListener('message', onPreviewMessage)
    // Signal readiness to admin
    if (window.parent !== window) {
      window.parent.postMessage({ type: 'DECOR_PREVIEW_READY', source: 'lyshop-app' }, '*')
    }
    return
  }

  // Normal mode: fetch from API
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
            { url: '/static/demo/banner-electronics.png', link: '/pages/product/list?category_id=1' },
            { url: '/static/demo/banner-fashion.png', link: '/pages/marketing/coupon?mode=claim' },
            { url: '/static/demo/banner-home.png', link: '/pages/product/list?category_id=3' },
          ],
          height: 340
        }
      },
      {
        type: 'notice',
        id: 'default_notice',
        props: {
          items: [
            { text: t('home.welcome'), link: '/pages/index/index' },
            { text: t('home.newUserCoupon'), link: '/pages/marketing/coupon?mode=claim' },
            { text: t('home.hotProducts'), link: '/pages/product/list' }
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

onUnmounted(() => {
  if (isPreview.value) {
    window.removeEventListener('message', onPreviewMessage)
  }
})
</script>
