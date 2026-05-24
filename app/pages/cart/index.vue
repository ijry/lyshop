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
        <!-- Select -->
        <view class="pr-16rpx flex-shrink-0" @click.stop="toggleItem(item.sku_id)">
          <view
            class="w-36rpx h-36rpx rounded-full border-2 flex items-center justify-center"
            :class="isChecked(item.sku_id) ? 'border-blue-700 bg-blue-700' : 'border-gray-300 bg-white'"
          >
            <u-icon v-if="isChecked(item.sku_id)" name="checkmark" size="14" color="#fff" />
          </view>
        </view>

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

    <!-- Recommend -->
    <view style="padding: 12px 16px 20px;">
      <text style="font-size: 15px; font-weight: 700; color: #111; display: block; margin-bottom: 12px;">猜你喜欢</text>
      <view style="display: flex; flex-wrap: wrap; margin: 0 -5px;">
        <view v-for="p in recommends" :key="p.product_id"
          @click="uni.navigateTo({url:`/pages/product/detail?id=${p.product_id}`})"
          style="width: 50%; padding: 5px; box-sizing: border-box;">
          <view style="background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 1px 6px rgba(0,0,0,0.04);">
            <image :src="p.cover" mode="aspectFill" style="width: 100%; height: 150px; display: block;" />
            <view style="padding: 8px 10px 12px;">
              <text style="font-size: 12px; color: #333; font-weight: 500; display: -webkit-box; -webkit-line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden;">{{ p.title }}</text>
              <view style="display: flex; align-items: baseline; gap: 4px; margin-top: 4px;">
                <text style="font-size: 15px; color: #dc2626; font-weight: 700;">¥{{ p.price }}</text>
                <text v-if="p.origin_price" style="font-size: 10px; color: #ccc; text-decoration: line-through;">¥{{ p.origin_price }}</text>
              </view>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- Bottom checkout bar -->
    <view v-if="items.length"
      class="fixed bottom-0 left-0 right-0 z-100 bg-white border-t-1 border-gray-100 px-30rpx py-20rpx flex items-center justify-between"
      :style="{paddingBottom: 'calc(20rpx + env(safe-area-inset-bottom))'}">
      <view class="flex items-center gap-24rpx">
        <view class="flex items-center" @click="toggleCheckAll">
          <view
            class="w-36rpx h-36rpx rounded-full border-2 flex items-center justify-center mr-10rpx"
            :class="allChecked ? 'border-blue-700 bg-blue-700' : 'border-gray-300 bg-white'"
          >
            <u-icon v-if="allChecked" name="checkmark" size="14" color="#fff" />
          </view>
          <text class="text-24rpx text-gray-600">全选</text>
        </view>
        <view class="flex items-baseline">
          <text class="text-26rpx text-gray-500">合计：</text>
          <text class="text-36rpx text-red-500 font-700 ml-4rpx">¥{{ selectedTotal.toFixed(2) }}</text>
        </view>
      </view>
      <u-button type="primary" :text="`结算(${selectedCount})`" shape="circle"
        :disabled="selectedCount === 0" :custom-style="{width: '220rpx'}" @click="checkout" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { get, post } from '@/utils/request'

const items = ref<any[]>([])
const recommends = ref<any[]>([])
const checkedSkuIds = ref<number[]>([])

const selectedItems = computed(() =>
  items.value.filter(i => checkedSkuIds.value.includes(i.sku_id))
)
const selectedCount = computed(() => selectedItems.value.length)
const selectedTotal = computed(() =>
  selectedItems.value.reduce((s, i) => s + (i.sku?.price || 0) * i.qty, 0)
)
const allChecked = computed(() =>
  items.value.length > 0 && checkedSkuIds.value.length === items.value.length
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
  checkedSkuIds.value = items.value.map(i => i.sku_id)
  // Load recommendations
  const rec = await get<any[]>('/api/v1/products/recommend')
  recommends.value = rec || []
}

async function updateQty(skuID: number, qty: number) {
  await post('/api/v1/cart/qty', { sku_id: skuID, qty })
}

async function remove(skuID: number) {
  items.value = items.value.filter(i => i.sku_id !== skuID)
  checkedSkuIds.value = checkedSkuIds.value.filter(id => id !== skuID)
}

function isChecked(skuID: number) {
  return checkedSkuIds.value.includes(skuID)
}

function toggleItem(skuID: number) {
  if (isChecked(skuID)) {
    checkedSkuIds.value = checkedSkuIds.value.filter(id => id !== skuID)
    return
  }
  checkedSkuIds.value.push(skuID)
}

function toggleCheckAll() {
  if (allChecked.value) {
    checkedSkuIds.value = []
    return
  }
  checkedSkuIds.value = items.value.map(i => i.sku_id)
}

function checkout() {
  if (!checkedSkuIds.value.length) {
    uni.showToast({ title: '请先勾选商品', icon: 'none' })
    return
  }
  const ids = checkedSkuIds.value.join(',')
  uni.navigateTo({ url: `/pages/order/confirm?sku_ids=${ids}` })
}

onMounted(loadCart)
</script>
