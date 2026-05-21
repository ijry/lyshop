<template>
  <view>
    <template v-for="comp in components" :key="comp.id">
      <!-- Banner -->
      <u-swiper
        v-if="comp.type === 'banner'"
        :list="comp.props?.images?.map((i: any) => i.url || i)"
        :height="comp.props?.height || 350"
      />

      <!-- Category nav -->
      <view v-else-if="comp.type === 'category_nav'"
        class="bg-white px-3 py-3">
        <scroll-view scroll-x>
          <view class="flex gap-4">
            <view v-for="item in comp.props?.items" :key="item.title"
              @click="uni.navigateTo({url: item.link})"
              class="flex flex-col items-center gap-1 shrink-0">
              <image v-if="item.icon" :src="item.icon" style="width:44px;height:44px;border-radius:12px;" />
              <view v-else class="bg-blue-50 flex items-center justify-center" style="width:44px;height:44px;border-radius:12px;">
                <text class="text-blue-700 text-xs">{{ item.title?.slice(0,2) }}</text>
              </view>
              <text class="text-slate-700 text-xs">{{ item.title }}</text>
            </view>
          </view>
        </scroll-view>
      </view>

      <!-- Product grid -->
      <view v-else-if="comp.type === 'product_grid'" class="px-3 py-2">
        <view class="grid grid-cols-2 gap-2">
          <view v-for="p in gridProducts[comp.id]" :key="p.id"
            @click="uni.navigateTo({url:`/pages/product/detail?id=${p.id}`})"
            class="bg-white rounded-xl overflow-hidden shadow-sm">
            <image :src="p.cover || ''" mode="aspectFill" style="width:100%;height:160px;" />
            <view class="p-2">
              <text class="text-slate-800 text-xs font-medium block truncate">{{ p.title }}</text>
              <text class="text-blue-700 text-sm font-bold mt-1 block">¥{{ p.price }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Notice -->
      <view v-else-if="comp.type === 'notice'" class="bg-yellow-50 px-4 py-2">
        <u-notice-bar :text="comp.props?.text" />
      </view>

      <!-- Image ad -->
      <view v-else-if="comp.type === 'image_ad'" class="px-3 py-2">
        <image v-if="comp.props?.url" :src="comp.props.url" mode="widthFix"
          class="w-full rounded-xl" @click="comp.props.link && uni.navigateTo({url:comp.props.link})" />
      </view>

      <!-- Rich text -->
      <view v-else-if="comp.type === 'rich_text'" class="px-4 py-2">
        <rich-text :nodes="comp.props?.content || ''" />
      </view>

      <!-- Spacer -->
      <view v-else-if="comp.type === 'spacer'"
        :style="{ height: (comp.props?.height || 16) + 'px', background: comp.props?.background || '#f8fafc' }" />

      <!-- Marketing zone placeholder -->
      <view v-else-if="comp.type === 'marketing_zone'" class="px-3 py-2">
        <view class="bg-gradient-to-r from-red-600 to-orange-500 rounded-xl px-4 py-3">
          <text class="text-white font-bold">限时秒杀</text>
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
      if (comp.props?.source === 'hot') params.sort = 'sales'
      if (comp.props?.category_id) params.category_id = comp.props.category_id
      const data = await get<any>('/api/v1/products', params)
      gridProducts.value[comp.id] = data?.list || []
    }
  }
}

onMounted(loadGridProducts)
watch(() => props.components, loadGridProducts, { deep: true })
</script>
