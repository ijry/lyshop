<template>
  <view class="min-h-screen bg-gray-50 pb-160rpx">
    <!-- Empty state -->
    <view v-if="!items.length" class="flex flex-col items-center pt-200rpx">
      <u-icon name="shopping-cart" size="60" color="#ccc" />
      <text class="text-gray-400 text-28rpx mt-24rpx mb-24rpx">购物车是空的</text>
      <u-button text="去逛逛" size="small" type="primary"
        @click="uni.switchTab({url:'/pages/product/list'})" />
    </view>

    <!-- Cart items -->
    <view v-else class="p-20rpx">
      <view v-for="item in items" :key="item.sku_id"
        class="flex items-center bg-white rounded-20rpx p-24rpx mb-20rpx shadow-sm">
        <!-- Product image -->
        <image :src="item.product?.cover" mode="aspectFill"
          class="w-160rpx h-160rpx rounded-16rpx flex-shrink-0" />

        <!-- Info -->
        <view class="flex-1 ml-24rpx overflow-hidden">
          <text class="text-28rpx text-gray-800 font-500 line-clamp-2">{{ item.product?.title }}</text>
          <view class="mt-8rpx">
            <text class="text-22rpx text-gray-400 bg-gray-100 px-12rpx py-4rpx rounded-6rpx">
              {{ skuLabel(item) }}
            </text>
          </view>
          <view class="flex items-center justify-between mt-16rpx">
            <text class="text-32rpx text-blue-700 font-700">¥{{ item.sku?.price }}</text>
            <u-number-box v-model="item.qty" :min="1" :max="99"
              @change="(v:any) => updateQty(item.sku_id, v.value)" />
          </view>
        </view>

        <!-- Delete -->
        <view class="ml-16rpx p-12rpx flex-shrink-0" @click="remove(item.sku_id)">
          <u-icon name="trash" size="18" color="#f56c6c" />
        </view>
      </view>
    </view>

    <!-- Bottom checkout bar -->
    <view v-if="items.length"
      class="fixed bottom-0 left-0 right-0 z-100 bg-white border-t-1 border-gray-100 px-30rpx py-20rpx flex items-center justify-between"
      :style="{paddingBottom: 'calc(20rpx + env(safe-area-inset-bottom))'}">
      <view class="flex items-baseline">
        <text class="text-26rpx text-gray-500">合计：</text>
        <text class="text-36rpx text-blue-700 font-700 ml-4rpx">¥{{ total.toFixed(2) }}</text>
      </view>
      <u-button type="primary" :text="`结算(${items.length})`" shape="circle"
        :custom-style="{width: '220rpx'}" @click="checkout" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { get, post } from '@/utils/request'

const items = ref<any[]>([])

const total = computed(() =>
  items.value.reduce((s, i) => s + (i.sku?.price || 0) * i.qty, 0)
)

function skuLabel(item: any) {
  if (!item.sku?.attrs) return '默认规格'
  try {
    const attrs = JSON.parse(item.sku.attrs)
    return attrs.map((a: any) => a.value).join(' / ')
  } catch { return '' }
}

async function loadCart() {
  const data = await get<any[]>('/api/v1/cart')
  items.value = data || []
}

async function updateQty(skuID: number, qty: number) {
  await post('/api/v1/cart/qty', { sku_id: skuID, qty })
}

async function remove(skuID: number) {
  items.value = items.value.filter(i => i.sku_id !== skuID)
}

function checkout() {
  const ids = items.value.map(i => i.sku_id).join(',')
  uni.navigateTo({ url: `/pages/order/confirm?sku_ids=${ids}` })
}

onMounted(loadCart)
</script>
