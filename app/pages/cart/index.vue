<template>
  <view class="min-h-screen pb-160rpx" style="background: var(--app-page-bg)">
    <!-- Empty state -->
    <view v-if="!items.length" class="flex flex-col items-center pt-200rpx">
      <u-icon name="shopping-cart" size="60" color="#ccc" />
      <text class="text-gray-400 text-28rpx mt-24rpx mb-24rpx">{{ $t('cart.empty') }}</text>
      <u-button :text="$t('cart.goShopping')" size="small" type="primary"
        @click="uni.switchTab({url:'/pages/product/list'})" />
    </view>

    <!-- Cart items -->
    <view v-else class="p-20rpx">
      <view v-for="item in items" :key="itemKey(item)"
        class="flex items-center rounded-20rpx p-24rpx mb-20rpx"
        :style="{ background: 'var(--app-card-bg)', boxShadow: 'var(--app-shadow-sm)' }">
        <!-- Select -->
        <view class="pr-16rpx flex-shrink-0" @click.stop="toggleItem(item)">
          <view
            class="w-36rpx h-36rpx rounded-full border-2 flex items-center justify-center"
            :class="isChecked(item) ? 'border-blue-700 bg-blue-700' : 'border-gray-300 bg-white'"
          >
            <u-icon v-if="isChecked(item)" name="checkmark" size="14" color="#fff" />
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
              @change="(v:any) => updateQty(item, v.value)" />
          </view>
        </view>

        <!-- Delete -->
        <view class="ml-16rpx p-12rpx flex-shrink-0" @click="remove(item)">
          <u-icon name="trash" size="18" color="#f56c6c" />
        </view>
      </view>
    </view>

    <!-- Recommend -->
    <view style="padding: 12px 16px 20px;">
      <text :style="{ fontSize: '15px', fontWeight: '700', color: 'var(--app-text-primary)', display: 'block', marginBottom: '12px' }">{{ $t('cart.guessYouLike') }}</text>
      <view style="display: flex; flex-wrap: wrap; margin: 0 -5px;">
        <view v-for="p in recommends" :key="p.product_id"
          @click="uni.navigateTo({url:`/pages/product/detail?id=${p.product_id}`})"
          style="width: 50%; padding: 5px; box-sizing: border-box;">
          <view :style="{ background: 'var(--app-card-bg)', borderRadius: '12px', overflow: 'hidden', boxShadow: 'var(--app-shadow-sm)' }">
            <image :src="p.cover" mode="aspectFill" style="width: 100%; height: 150px; display: block;" />
            <view style="padding: 8px 10px 12px;">
              <text :style="{ fontSize: '12px', color: 'var(--app-text-secondary)', fontWeight: '500', display: '-webkit-box', '-webkit-line-clamp': '1', '-webkit-box-orient': 'vertical', overflow: 'hidden' }">{{ p.title }}</text>
              <view style="display: flex; align-items: baseline; gap: 4px; margin-top: 4px;">
                <text :style="{ fontSize: '15px', color: 'var(--app-price-color)', fontWeight: '700' }">¥{{ p.price }}</text>
                <text v-if="p.origin_price" :style="{ fontSize: '10px', color: 'var(--app-text-placeholder)', textDecoration: 'line-through' }">¥{{ p.origin_price }}</text>
              </view>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- Bottom checkout bar -->
    <view v-if="items.length"
      class="fixed bottom-0 left-0 right-0 z-100 px-30rpx py-20rpx flex items-center justify-between"
      :style="{ background: 'var(--app-card-bg)', borderTop: '1px solid var(--app-border-color)', paddingBottom: 'calc(20rpx + env(safe-area-inset-bottom))' }">
      <view class="flex items-center gap-24rpx">
        <view class="flex items-center" @click="toggleCheckAll">
          <view
            class="w-36rpx h-36rpx rounded-full border-2 flex items-center justify-center mr-10rpx"
            :class="allChecked ? 'border-blue-700 bg-blue-700' : 'border-gray-300 bg-white'"
          >
            <u-icon v-if="allChecked" name="checkmark" size="14" color="#fff" />
          </view>
          <text class="text-24rpx text-gray-600">{{ $t('cart.selectAll') }}</text>
        </view>
        <view class="flex items-baseline">
          <text class="text-26rpx text-gray-500">{{ $t('cart.total') }}</text>
          <text class="text-36rpx text-red-500 font-700 ml-4rpx">¥{{ selectedTotal.toFixed(2) }}</text>
        </view>
      </view>
      <u-button type="primary" :text="$t('cart.checkout') + selectedCount + ')'" shape="circle"
        :disabled="selectedCount === 0" :custom-style="{width: '220rpx'}" @click="checkout" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { del, get, post } from '@/utils/request'

const { t } = useI18n()

const items = ref<any[]>([])
const recommends = ref<any[]>([])
const checkedItemKeys = ref<string[]>([])

function itemKey(item: any) {
  return `${Number(item?.sku_id || 0)}:${Number(item?.activity_product_id || 0)}`
}

const selectedItems = computed(() =>
  items.value.filter(i => checkedItemKeys.value.includes(itemKey(i)))
)
const selectedCount = computed(() => selectedItems.value.length)
const selectedTotal = computed(() =>
  selectedItems.value.reduce((s, i) => s + (i.sku?.price || 0) * i.qty, 0)
)
const allChecked = computed(() =>
  items.value.length > 0 && checkedItemKeys.value.length === items.value.length
)

function skuLabel(item: any) {
  if (!item.sku?.attrs) return t('cart.defaultSpec')
  try {
    const attrs = JSON.parse(item.sku.attrs)
    return attrs.map((a: any) => a.value).join(' / ')
  } catch { return '' }
}

async function loadCart() {
  const data = await get<any[]>('/api/v1/cart')
  items.value = data || []
  checkedItemKeys.value = items.value.map(itemKey)
  // Load recommendations
  const rec = await get<any[]>('/api/v1/products/recommend')
  recommends.value = rec || []
}

async function updateQty(item: any, qty: number) {
  await post('/api/v1/cart/qty', {
    sku_id: Number(item?.sku_id || 0),
    activity_product_id: Number(item?.activity_product_id || 0),
    qty,
  })
}

async function remove(item: any) {
  const key = itemKey(item)
  items.value = items.value.filter(i => itemKey(i) !== key)
  checkedItemKeys.value = checkedItemKeys.value.filter(v => v !== key)
  await del(`/api/v1/cart/${Number(item?.sku_id || 0)}`, {
    activity_product_id: Number(item?.activity_product_id || 0),
  })
}

function isChecked(item: any) {
  return checkedItemKeys.value.includes(itemKey(item))
}

function toggleItem(item: any) {
  const key = itemKey(item)
  if (checkedItemKeys.value.includes(key)) {
    checkedItemKeys.value = checkedItemKeys.value.filter(v => v !== key)
    return
  }
  checkedItemKeys.value.push(key)
}

function toggleCheckAll() {
  if (allChecked.value) {
    checkedItemKeys.value = []
    return
  }
  checkedItemKeys.value = items.value.map(itemKey)
}

function checkout() {
  if (!checkedItemKeys.value.length) {
    uni.showToast({ title: t('cart.selectFirst'), icon: 'none' })
    return
  }
  const selected = selectedItems.value.map((item: any) => ({
    sku_id: Number(item?.sku_id || 0),
    activity_product_id: Number(item?.activity_product_id || 0),
  }))
  const ids = selected.map((item: any) => item.sku_id).join(',')
  const encodedItems = encodeURIComponent(JSON.stringify(selected))
  uni.navigateTo({ url: `/pages/order/confirm?items=${encodedItems}&sku_ids=${ids}` })
}

onMounted(loadCart)
</script>
