<template>
  <view class="bg-gray-50 min-h-screen">
    <u-navbar title="lyshop" :placeholder="true" />
    <DecorRender :components="components" />
    <!-- Float customer service button -->
    <view class="fixed bottom-20 right-4 z-50">
      <view @click="uni.navigateTo({url:'/pages/im/chat'})"
        class="w-12 h-12 bg-blue-700 rounded-full flex items-center justify-center shadow-lg">
        <text class="text-white text-xs">客服</text>
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
    try { components.value = JSON.parse(data.components) } catch {}
  }
  if (!components.value.length) {
    // Default components when no decoration configured
    components.value = [
      { type: 'banner', id: 'default_banner', props: { images: [], height: 300 } },
      { type: 'product_grid', id: 'default_grid', props: { source: 'hot', limit: 10, columns: 2 } }
    ]
  }
})
</script>
