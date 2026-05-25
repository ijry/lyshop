<template>
  <view class="min-h-screen bg-gray-50 p-20rpx">
    <view v-if="loading && !favorites.length" class="flex justify-center py-80rpx">
      <u-loading-icon :text="$t('favorites.loading')" />
    </view>

    <view v-else-if="!favorites.length" class="text-center py-120rpx text-gray-400 text-28rpx">
      {{ $t('favorites.empty') }}
    </view>

    <view v-else>
      <view
        v-for="(item, index) in favorites"
        :key="item.id"
        class="bg-white rounded-16rpx p-20rpx mb-16rpx flex"
      >
        <image
          :src="item.cover"
          mode="aspectFill"
          class="w-180rpx h-180rpx rounded-12rpx"
          @click="toDetail(item.id)"
        />
        <view class="flex-1 ml-16rpx flex flex-col justify-between min-w-0">
          <view @click="toDetail(item.id)">
            <text class="text-28rpx text-gray-800 font-600 line-clamp-2">{{ item.title }}</text>
            <text v-if="item.subtitle" class="text-22rpx text-gray-500 mt-8rpx block line-clamp-1">{{ item.subtitle }}</text>
          </view>
          <view class="flex items-center justify-between mt-12rpx">
            <view>
              <text class="text-red-600 text-30rpx font-700">¥{{ item.price }}</text>
              <text v-if="item.origin_price" class="text-22rpx text-gray-400 line-through ml-10rpx">¥{{ item.origin_price }}</text>
            </view>
            <u-button type="error" plain size="mini" :text="$t('favorites.unfavorite')" @click="unfavorite(item, index)" />
          </view>
        </view>
      </view>
      <view v-if="loading" class="text-center py-20rpx text-24rpx text-gray-400">{{ $t('favorites.loading') }}</view>
      <view v-else-if="finished" class="text-center py-20rpx text-24rpx text-gray-400">{{ $t('favorites.noMore') }}</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { onReachBottom } from '@dcloudio/uni-app'
import { del, get } from '@/utils/request'

const favorites = ref<any[]>([])
const loading = ref(false)
const finished = ref(false)
const page = ref(1)
const size = 10
const total = ref(0)

async function loadFavorites(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    finished.value = false
  }
  if (finished.value) return
  loading.value = true
  try {
    const data = await get<any>('/api/v1/user/favorites', { page: page.value, size })
    const list = Array.isArray(data?.list) ? data.list : []
    total.value = Number(data?.total || 0)
    if (reset) {
      favorites.value = list
    } else {
      favorites.value.push(...list)
    }
    const loaded = favorites.value.length
    finished.value = loaded >= total.value || list.length < size
    if (!finished.value) page.value += 1
  } finally {
    loading.value = false
  }
}

function toDetail(id: number) {
  uni.navigateTo({ url: `/pages/product/detail?id=${id}` })
}

async function unfavorite(item: any, index: number) {
  await del(`/api/v1/products/${item.id}/favorite`)
  favorites.value.splice(index, 1)
  total.value = Math.max(0, total.value - 1)
  await loadFavorites(true)
}

onMounted(() => {
  loadFavorites(true)
})

onReachBottom(() => {
  loadFavorites(false)
})
</script>
