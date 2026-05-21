<template>
  <view class="bg-gray-50 min-h-screen pb-24">
    <u-navbar title="购物车" :placeholder="true" />

    <view v-if="!items.length" class="text-center py-16 text-gray-400">
      <text class="block mb-4">购物车是空的</text>
      <u-button text="去逛逛" size="small" type="primary" @click="uni.switchTab({url:'/pages/product/list'})" />
    </view>

    <view v-else>
      <u-swipe-action v-for="item in items" :key="item.sku_id" class="mb-2">
        <view class="bg-white p-4 flex gap-3 items-center">
          <image :src="item.product?.cover || '/static/placeholder.png'"
            mode="aspectFill" style="width:72px;height:72px;border-radius:8px;" />
          <view class="flex-1">
            <text class="text-slate-800 text-sm font-medium block">{{ item.product?.title }}</text>
            <text class="text-gray-500 text-xs mt-1 block">{{ skuLabel(item) }}</text>
            <view class="flex items-center justify-between mt-2">
              <text class="text-blue-700 font-bold">¥{{ item.sku?.price }}</text>
              <u-number-box v-model="item.qty" :min="1" :max="99"
                @change="(v:number)=>updateQty(item.sku_id, v)" />
            </view>
          </view>
        </view>
        <template #right>
          <view @click="remove(item.sku_id)"
            class="bg-red-500 text-white flex items-center justify-center px-5 h-full">
            <text class="text-sm">删除</text>
          </view>
        </template>
      </u-swipe-action>
    </view>

    <!-- Bottom checkout bar -->
    <view v-if="items.length"
      class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-100 p-3 flex items-center justify-between">
      <view>
        <text class="text-gray-500 text-sm">合计：</text>
        <text class="text-blue-700 font-bold text-lg">¥{{ total.toFixed(2) }}</text>
      </view>
      <u-button type="primary" :text="`结算(${items.length})`" @click="checkout" />
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
  await post(`/api/v1/cart/${skuID}`, {})
  loadCart()
}

function checkout() {
  const ids = items.value.map(i => i.sku_id).join(',')
  uni.navigateTo({ url: `/pages/order/confirm?sku_ids=${ids}` })
}

onMounted(loadCart)
</script>
