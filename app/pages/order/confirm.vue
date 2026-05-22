<template>
  <view class="min-h-screen bg-gray-50 pb-140rpx">
    <u-navbar title="确认订单" :placeholder="true" />

    <!-- Address -->
    <view class="bg-white mx-20rpx mt-20rpx rounded-20rpx p-30rpx"
      style="box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);">
      <view class="flex items-center justify-between mb-16rpx">
        <text class="text-28rpx font-600 text-gray-800">收货地址</text>
        <u-icon name="arrow-right" size="16" color="#999" />
      </view>
      <view v-if="!address" class="text-center py-20rpx">
        <u-button size="small" text="添加收货地址" @click="showAddressForm = true" type="primary" plain shape="circle" />
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
      <text class="text-28rpx font-600 text-gray-800 block mb-20rpx">支付方式</text>
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
      <u-input v-model="remark" placeholder="备注（选填）" border="none" />
    </view>

    <!-- Bottom bar -->
    <view class="fixed bottom-0 left-0 right-0 z-100 bg-white border-t-1 border-gray-100 px-30rpx py-20rpx flex items-center justify-between"
      :style="{paddingBottom: 'calc(20rpx + env(safe-area-inset-bottom))'}">
      <view class="flex items-baseline">
        <text class="text-26rpx text-gray-500">应付：</text>
        <text class="text-36rpx text-red-500 font-700">¥--</text>
      </view>
      <u-button type="primary" text="提交订单" :loading="submitting" @click="submit"
        shape="circle" :custom-style="{width: '240rpx'}" />
    </view>

    <!-- Address popup -->
    <u-popup :show="showAddressForm" mode="bottom" round="20" @close="showAddressForm=false">
      <view class="p-40rpx">
        <text class="text-32rpx font-600 text-gray-800 block mb-30rpx">添加收货地址</text>
        <view class="mb-24rpx">
          <u-input v-model="addrForm.name" placeholder="收货人姓名" border="surround" shape="circle" />
        </view>
        <view class="mb-24rpx">
          <u-input v-model="addrForm.phone" placeholder="手机号" type="number" border="surround" shape="circle" />
        </view>
        <view class="flex gap-16rpx mb-24rpx">
          <view class="flex-1"><u-input v-model="addrForm.province" placeholder="省" border="surround" shape="circle" /></view>
          <view class="flex-1"><u-input v-model="addrForm.city" placeholder="市" border="surround" shape="circle" /></view>
          <view class="flex-1"><u-input v-model="addrForm.district" placeholder="区" border="surround" shape="circle" /></view>
        </view>
        <view class="mb-30rpx">
          <u-input v-model="addrForm.detail" placeholder="详细地址" border="surround" shape="circle" />
        </view>
        <u-button type="primary" text="保存地址" @click="saveAddress" shape="circle" />
      </view>
    </u-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get, post } from '@/utils/request'

const address = ref<any>(null)
const payMethod = ref('wechat')
const remark = ref('')
const submitting = ref(false)
const showAddressForm = ref(false)
const addrForm = ref({ name: '', phone: '', province: '', city: '', district: '', detail: '' })

const payMethods = [
  { label: '微信支付', value: 'wechat' },
  { label: '支付宝', value: 'alipay' },
]

let skuIds: number[] = []

onMounted(async () => {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  skuIds = (query.sku_ids || '').split(',').map(Number).filter(Boolean)

  const addrs = await get<any[]>('/api/v1/addresses')
  if (addrs?.length) address.value = addrs.find((a: any) => a.is_default) || addrs[0]
})

async function saveAddress() {
  const saved = await post<any>('/api/v1/addresses', addrForm.value)
  address.value = saved
  showAddressForm.value = false
}

async function submit() {
  if (!address.value) { uni.showToast({ title: '请添加收货地址', icon: 'none' }); return }
  submitting.value = true
  try {
    await post<any>('/api/v1/orders', {
      address_id: address.value.id,
      payment_method: payMethod.value,
      sku_ids: skuIds,
      remark: remark.value
    })
    uni.showToast({ title: '下单成功', icon: 'success' })
    setTimeout(() => uni.switchTab({ url: '/pages/order/list' }), 1000)
  } catch {} finally {
    submitting.value = false
  }
}
</script>
