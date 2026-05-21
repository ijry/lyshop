<template>
  <view class="bg-gray-50 min-h-screen pb-24">
    <u-navbar title="确认订单" :placeholder="true" />

    <!-- Address -->
    <view class="bg-white mx-3 mt-3 rounded-xl p-4 shadow-sm">
      <text class="text-sm font-medium text-slate-700 block mb-2">收货地址</text>
      <view v-if="!address" class="text-center py-2">
        <u-button size="small" text="添加收货地址" @click="showAddressForm = true" />
      </view>
      <view v-else>
        <text class="text-slate-800 font-medium">{{ address.name }}  {{ address.phone }}</text>
        <text class="text-gray-500 text-sm block mt-1">
          {{ address.province }}{{ address.city }}{{ address.district }}{{ address.detail }}
        </text>
      </view>
    </view>

    <!-- Payment method -->
    <view class="bg-white mx-3 mt-3 rounded-xl p-4 shadow-sm">
      <text class="text-sm font-medium text-slate-700 block mb-3">支付方式</text>
      <view class="flex gap-3">
        <view v-for="m in payMethods" :key="m.value"
          @click="payMethod = m.value"
          :class="payMethod === m.value ? 'border-blue-700 text-blue-700 bg-blue-50' : 'border-gray-200 text-gray-600'"
          class="flex-1 border rounded-xl py-3 text-center text-sm">
          {{ m.label }}
        </view>
      </view>
    </view>

    <!-- Remark -->
    <view class="bg-white mx-3 mt-3 rounded-xl p-4 shadow-sm">
      <u-input v-model="remark" placeholder="备注（选填）" border="none" />
    </view>

    <!-- Bottom bar -->
    <view class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-100 p-3 flex items-center justify-between">
      <text class="text-gray-500 text-sm">应付：<text class="text-blue-700 font-bold text-lg">¥--</text></text>
      <u-button type="primary" text="提交订单" :loading="submitting" @click="submit" />
    </view>

    <!-- Quick address form -->
    <u-popup :show="showAddressForm" mode="bottom" round @close="showAddressForm=false">
      <view class="p-5">
        <text class="text-lg font-semibold text-slate-800 block mb-4">添加地址</text>
        <u-form :model="addrForm" class="space-y-3">
          <u-form-item label="姓名"><u-input v-model="addrForm.name" /></u-form-item>
          <u-form-item label="电话"><u-input v-model="addrForm.phone" type="number" /></u-form-item>
          <u-form-item label="省市区"><u-input v-model="addrForm.province" placeholder="省" /><u-input v-model="addrForm.city" placeholder="市" /><u-input v-model="addrForm.district" placeholder="区" /></u-form-item>
          <u-form-item label="详细地址"><u-input v-model="addrForm.detail" /></u-form-item>
        </u-form>
        <u-button type="primary" text="保存" class="mt-4" @click="saveAddress" />
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
    const order = await post<any>('/api/v1/orders', {
      address_id: address.value.id,
      payment_method: payMethod.value,
      sku_ids: skuIds,
      remark: remark.value
    })
    uni.showToast({ title: '下单成功', icon: 'success' })
    setTimeout(() => uni.navigateTo({ url: `/pages/order/list` }), 1000)
  } catch {} finally {
    submitting.value = false
  }
}
</script>
