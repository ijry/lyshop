<template>
  <view class="min-h-screen bg-gray-50 pb-140rpx">
    <u-navbar :title="$t('orderConfirm.title')" :placeholder="true" />

    <!-- Address -->
    <view class="bg-white mx-20rpx mt-20rpx rounded-20rpx p-30rpx"
      style="box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);">
      <view class="flex items-center justify-between mb-16rpx">
        <text class="text-28rpx font-600 text-gray-800">{{ $t('orderConfirm.shippingAddress') }}</text>
        <view class="flex items-center gap-20rpx">
          <text v-if="address" class="text-22rpx text-blue-500" @click="openAddressEditor(address)">{{ $t('orderConfirm.edit') }}</text>
          <u-icon name="arrow-right" size="16" color="#999" />
        </view>
      </view>
      <view v-if="!address" class="text-center py-20rpx">
        <u-button size="small" :text="$t('orderConfirm.addAddress')" @click="openAddressEditor()" type="primary" plain shape="circle" />
      </view>
      <view v-else>
        <view class="flex items-center gap-16rpx">
          <text class="text-28rpx font-600 text-gray-800">{{ address.name }}</text>
          <text class="text-26rpx text-gray-500">{{ address.phone }}</text>
        </view>
        <text class="text-24rpx text-gray-400 mt-8rpx block">
          {{ address.province }}{{ address.city }}{{ address.district }} {{ address.detail }}
        </text>
      </view>
    </view>

    <!-- Payment method -->
    <view class="bg-white mx-20rpx mt-20rpx rounded-20rpx p-30rpx"
      style="box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);">
      <text class="text-28rpx font-600 text-gray-800 block mb-20rpx">{{ $t('orderConfirm.paymentMethod') }}</text>
      <view class="flex gap-20rpx">
        <view v-for="m in payMethods" :key="m.value"
          @click="payMethod = m.value"
          :class="payMethod === m.value
            ? 'border-blue-700 text-blue-700 bg-blue-50'
            : 'border-gray-200 text-gray-600'"
          class="flex-1 border-1 rounded-16rpx py-24rpx text-center text-26rpx">
          {{ m.label }}
        </view>
      </view>
    </view>

    <!-- Remark -->
    <view class="bg-white mx-20rpx mt-20rpx rounded-20rpx p-30rpx"
      style="box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);">
      <u-input v-model="remark" :placeholder="$t('orderConfirm.remark')" border="none" />
    </view>

    <!-- Bottom bar -->
    <view class="fixed bottom-0 left-0 right-0 z-100 bg-white border-t-1 border-gray-100 px-30rpx py-20rpx flex items-center justify-between"
      :style="{paddingBottom: 'calc(20rpx + env(safe-area-inset-bottom))'}">
      <view class="flex items-baseline">
        <text class="text-26rpx text-gray-500">{{ $t('orderConfirm.amountDue') }}</text>
        <text class="text-36rpx text-red-500 font-700">¥--</text>
      </view>
      <u-button type="primary" :text="$t('orderConfirm.submitOrder')" :loading="submitting" @click="submit"
        shape="circle" :custom-style="{width: '240rpx'}" />
    </view>

    <!-- Address popup -->
    <u-popup :show="showAddressForm" mode="bottom" round="20" @close="closeAddressEditor">
      <view class="p-40rpx">
        <text class="text-32rpx font-600 text-gray-800 block mb-30rpx">{{ editingAddressID ? $t('orderConfirm.editAddress') : $t('orderConfirm.addNewAddress') }}</text>
        <view class="mb-24rpx">
          <u-input v-model="addrForm.name" :placeholder="$t('orderConfirm.recipientName')" border="surround" shape="circle" />
        </view>
        <view class="mb-24rpx">
          <u-input v-model="addrForm.phone" :placeholder="$t('orderConfirm.phone')" type="number" border="surround" shape="circle" />
        </view>
        <view class="mb-24rpx">
          <view @click="showRegionPicker = true">
            <u-input :modelValue="regionText" :placeholder="$t('orderConfirm.region')" border="surround" shape="circle" readonly />
          </view>
        </view>
        <view class="mb-30rpx">
          <u-input v-model="addrForm.detail" :placeholder="$t('orderConfirm.detailAddress')" border="surround" shape="circle" />
        </view>
        <view class="flex items-center justify-between mb-30rpx">
          <text class="text-24rpx text-gray-500">{{ $t('orderConfirm.setDefault') }}</text>
          <u-switch v-model="isDefault" />
        </view>
        <u-button type="primary" :text="$t('orderConfirm.saveAddress')" @click="saveAddress" shape="circle" />
      </view>
    </u-popup>

    <up-cascader
      v-model:show="showRegionPicker"
      :data="regionOptions"
      v-model="regionValues"
      valueKey="value"
      labelKey="label"
      childrenKey="children"
      @cancel="showRegionPicker = false"
      @confirm="onRegionConfirm"
    />
  </view>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { get, post, put } from '@/utils/request'

const { t } = useI18n()

const address = ref<any>(null)
const payMethod = ref('wechat')
const remark = ref('')
const submitting = ref(false)
const showAddressForm = ref(false)
const showRegionPicker = ref(false)
const editingAddressID = ref(0)
const isDefault = ref(false)
const addrForm = ref({ name: '', phone: '', province: '', city: '', district: '', detail: '' })
const regionValues = ref<string[]>([])
const regionOptions = [
  {
    label: '北京市', value: '北京市', children: [
      { label: '北京市', value: '北京市', children: [{ label: '朝阳区', value: '朝阳区' }, { label: '海淀区', value: '海淀区' }] },
    ],
  },
  {
    label: '上海市', value: '上海市', children: [
      { label: '上海市', value: '上海市', children: [{ label: '浦东新区', value: '浦东新区' }, { label: '闵行区', value: '闵行区' }] },
    ],
  },
  {
    label: '广东省', value: '广东省', children: [
      { label: '广州市', value: '广州市', children: [{ label: '天河区', value: '天河区' }, { label: '海珠区', value: '海珠区' }] },
      { label: '深圳市', value: '深圳市', children: [{ label: '南山区', value: '南山区' }, { label: '福田区', value: '福田区' }] },
    ],
  },
]
const regionText = computed(() => [addrForm.value.province, addrForm.value.city, addrForm.value.district].filter(Boolean).join(' / '))

const payMethods = computed(() => [
  { label: t('orderConfirm.wechatPay'), value: 'wechat' },
  { label: t('orderConfirm.alipay'), value: 'alipay' },
])

let skuIds: number[] = []
let orderItems: Array<{ sku_id: number; activity_product_id: number }> = []

onMounted(async () => {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  skuIds = (query.sku_ids || '').split(',').map(Number).filter(Boolean)
  const rawItems = String(query.items || '').trim()
  if (rawItems) {
    try {
      const parsed = JSON.parse(decodeURIComponent(rawItems))
      if (Array.isArray(parsed)) {
        orderItems = parsed
          .map((item: any) => ({
            sku_id: Number(item?.sku_id || 0),
            activity_product_id: Number(item?.activity_product_id || 0),
          }))
          .filter((item: any) => item.sku_id > 0)
      }
    } catch {
      orderItems = []
    }
  }
  if (!orderItems.length) {
    orderItems = skuIds.map((skuID) => ({ sku_id: skuID, activity_product_id: 0 }))
  }
  if (!skuIds.length) {
    skuIds = orderItems.map((item) => item.sku_id)
  }

  await loadAddresses()
})

async function loadAddresses() {
  const addrs = await get<any[]>('/api/v1/addresses')
  if (addrs?.length) {
    address.value = addrs.find((a: any) => Number(a.is_default || 0) === 1) || addrs[0]
  } else {
    address.value = null
  }
}

function resetAddressForm() {
  editingAddressID.value = 0
  isDefault.value = false
  addrForm.value = { name: '', phone: '', province: '', city: '', district: '', detail: '' }
  regionValues.value = []
}

function openAddressEditor(addr?: any) {
  if (!addr) {
    resetAddressForm()
    showAddressForm.value = true
    return
  }
  editingAddressID.value = Number(addr.id || 0)
  isDefault.value = Number(addr.is_default || 0) === 1
  addrForm.value = {
    name: addr.name || '',
    phone: addr.phone || '',
    province: addr.province || '',
    city: addr.city || '',
    district: addr.district || '',
    detail: addr.detail || '',
  }
  regionValues.value = [addrForm.value.province, addrForm.value.city, addrForm.value.district].filter(Boolean)
  showAddressForm.value = true
}

function closeAddressEditor() {
  showAddressForm.value = false
  resetAddressForm()
}

function onRegionConfirm(values: string[]) {
  showRegionPicker.value = false
  if (!Array.isArray(values) || values.length < 3) return
  addrForm.value.province = values[0] || ''
  addrForm.value.city = values[1] || ''
  addrForm.value.district = values[2] || ''
}

function validateAddress() {
  if (!addrForm.value.name.trim()) return t('orderConfirm.nameRequired')
  if (!/^1\d{10}$/.test(addrForm.value.phone)) return t('orderConfirm.phoneInvalid')
  if (!addrForm.value.province || !addrForm.value.city || !addrForm.value.district) return t('orderConfirm.regionRequired')
  if (!addrForm.value.detail.trim()) return t('orderConfirm.addressRequired')
  return ''
}

async function saveAddress() {
  const error = validateAddress()
  if (error) {
    uni.showToast({ title: error, icon: 'none' })
    return
  }
  const payload = { ...addrForm.value, is_default: isDefault.value ? 1 : 0 }
  try {
    if (editingAddressID.value) {
      await put<any>(`/api/v1/addresses/${editingAddressID.value}`, payload)
    } else {
      await post<any>('/api/v1/addresses', payload)
    }
    await loadAddresses()
    closeAddressEditor()
  } catch {
    uni.showToast({ title: t('orderConfirm.saveFailed'), icon: 'none' })
  }
}

async function submit() {
  if (!address.value) { uni.showToast({ title: t('orderConfirm.addressRequired2'), icon: 'none' }); return }
  if (!orderItems.length && !skuIds.length) {
    uni.showToast({ title: t('cart.selectFirst'), icon: 'none' })
    return
  }
  submitting.value = true
  try {
    await post<any>('/api/v1/orders', {
      address_id: address.value.id,
      payment_method: payMethod.value,
      items: orderItems,
      sku_ids: skuIds,
      remark: remark.value
    })
    uni.showToast({ title: t('orderConfirm.orderSuccess'), icon: 'success' })
    setTimeout(() => uni.switchTab({ url: '/pages/order/list' }), 1000)
  } catch {} finally {
    submitting.value = false
  }
}
</script>
