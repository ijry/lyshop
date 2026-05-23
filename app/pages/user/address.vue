<template>
  <view style="min-height: 100vh; background: #f5f5f5;">
    <u-navbar title="收货地址" :placeholder="true" />

    <view style="padding: 12px 16px 96px;">
      <view
        v-for="addr in addresses"
        :key="addr.id"
        style="background: #fff; border-radius: 16px; padding: 16px 20px; margin-bottom: 12px; box-shadow: 0 1px 6px rgba(0,0,0,0.04);"
      >
        <view style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px;">
          <view style="display: flex; align-items: center; gap: 12px;">
            <text style="font-size: 15px; font-weight: 600; color: #111;">{{ addr.name }}</text>
            <text style="font-size: 13px; color: #666;">{{ addr.phone }}</text>
          </view>
          <view
            v-if="addr.is_default"
            style="background: #fef2f2; color: #dc2626; font-size: 11px; padding: 2px 8px; border-radius: 4px;"
          >
            默认
          </view>
        </view>
        <text style="font-size: 13px; color: #999; line-height: 1.5;">
          {{ addr.province }}{{ addr.city }}{{ addr.district }} {{ addr.detail }}
        </text>
      </view>

      <view v-if="!addresses.length" style="text-align: center; padding: 80px 0; color: #999; font-size: 14px;">
        暂无收货地址
      </view>
    </view>

    <view style="padding: 16px; position: fixed; bottom: 0; left: 0; right: 0; background: #f5f5f5;">
      <u-button
        type="primary"
        text="新增收货地址"
        shape="circle"
        :custom-style="{background: '#dc2626', borderColor: '#dc2626'}"
        @click="showEditor = true"
      />
    </view>

    <u-popup :show="showEditor" mode="bottom" round="20" @close="closeEditor">
      <view style="padding: 20px 16px 24px;">
        <text style="display: block; font-size: 16px; font-weight: 700; margin-bottom: 16px;">新增收货地址</text>
        <view style="display: flex; flex-direction: column; gap: 12px;">
          <u-input v-model="form.name" placeholder="收货人姓名" border="surround" />
          <u-input v-model="form.phone" placeholder="手机号" border="surround" maxlength="11" type="number" />
          <u-input v-model="form.province" placeholder="省份" border="surround" />
          <u-input v-model="form.city" placeholder="城市" border="surround" />
          <u-input v-model="form.district" placeholder="区县" border="surround" />
          <u-input v-model="form.detail" placeholder="详细地址" border="surround" />
          <view style="display: flex; align-items: center; justify-content: space-between;">
            <text style="font-size: 13px; color: #666;">设为默认地址</text>
            <u-switch v-model="defaultSwitch" activeColor="#dc2626" />
          </view>
        </view>
        <view style="display: flex; gap: 12px; margin-top: 18px;">
          <u-button text="取消" shape="circle" @click="closeEditor" />
          <u-button
            text="保存"
            type="primary"
            shape="circle"
            :loading="saving"
            :custom-style="{background: '#dc2626', borderColor: '#dc2626'}"
            @click="saveAddress"
          />
        </view>
      </view>
    </u-popup>
  </view>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { get, post } from '@/utils/request'

const IS_DEMO_MODE = import.meta.env.VITE_MOCK === 'true'
const addresses = ref<any[]>([])
const showEditor = ref(false)
const saving = ref(false)
const defaultSwitch = ref(false)
const form = ref({
  name: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
})

function resetForm() {
  form.value = {
    name: '',
    phone: '',
    province: '',
    city: '',
    district: '',
    detail: '',
  }
  defaultSwitch.value = false
}

function closeEditor() {
  showEditor.value = false
  resetForm()
}

function validateForm() {
  if (!form.value.name.trim()) return '请输入收货人姓名'
  if (!/^1\d{10}$/.test(form.value.phone)) return '请输入正确手机号'
  if (!form.value.province.trim() || !form.value.city.trim() || !form.value.district.trim()) return '请填写完整省市区'
  if (!form.value.detail.trim()) return '请输入详细地址'
  return ''
}

async function loadAddresses() {
  const data = await get<any[]>('/api/v1/addresses')
  addresses.value = Array.isArray(data) ? data : []
}

async function saveAddress() {
  const error = validateForm()
  if (error) {
    uni.showToast({ title: error, icon: 'none' })
    return
  }
  if (saving.value) return
  saving.value = true

  const payload = {
    ...form.value,
    is_default: defaultSwitch.value ? 1 : 0,
  }

  let createdID = Date.now()
  try {
    const res = await post<any>('/api/v1/addresses', payload)
    if (res?.id) createdID = Number(res.id)
  } catch (error) {
    if (!IS_DEMO_MODE) {
      uni.showToast({ title: '保存失败，请稍后重试', icon: 'none' })
      saving.value = false
      return
    }
  }

  if (payload.is_default) {
    addresses.value = addresses.value.map((item) => ({ ...item, is_default: 0 }))
  }
  addresses.value.unshift({
    id: createdID,
    ...payload,
  })

  uni.showToast({ title: '新增成功', icon: 'success' })
  saving.value = false
  closeEditor()
}

onMounted(loadAddresses)
</script>
