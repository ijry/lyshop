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
          <view style="display:flex; align-items:center; gap:8px;">
            <view
              v-if="addr.is_default"
              style="background: #fef2f2; color: #dc2626; font-size: 11px; padding: 2px 8px; border-radius: 4px;"
            >
              默认
            </view>
            <text style="font-size:12px;color:#3b82f6;" @click="openEditor(addr)">编辑</text>
            <text style="font-size:12px;color:#ef4444;" @click="removeAddress(addr)">删除</text>
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
        @click="openEditor()"
      />
    </view>

    <u-popup :show="showEditor" mode="bottom" round="20" @close="closeEditor">
      <view style="padding: 20px 16px 24px;">
        <text style="display: block; font-size: 16px; font-weight: 700; margin-bottom: 16px;">{{ editingID ? '编辑收货地址' : '新增收货地址' }}</text>
        <view style="display: flex; flex-direction: column; gap: 12px;">
          <u-input v-model="form.name" placeholder="收货人姓名" border="surround" />
          <u-input v-model="form.phone" placeholder="手机号" border="surround" maxlength="11" type="number" />
          <view @click="showRegionPicker = true">
            <u-input :modelValue="regionText" placeholder="省/市/区" border="surround" readonly />
          </view>
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

    <up-cascader
      v-model:show="showRegionPicker"
      :data="regionOptions"
      v-model="regionValues"
      @cancel="showRegionPicker = false"
      @confirm="onRegionConfirm"
      valueKey="value"
      labelKey="label"
      childrenKey="children"
    />
  </view>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { del, get, post, put } from '@/utils/request'

const addresses = ref<any[]>([])
const showEditor = ref(false)
const showRegionPicker = ref(false)
const saving = ref(false)
const defaultSwitch = ref(false)
const editingID = ref<number>(0)

const form = ref({
  name: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
})

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
const regionValues = ref<string[]>([])
const regionText = computed(() => [form.value.province, form.value.city, form.value.district].filter(Boolean).join(' / '))

function resetForm() {
  form.value = { name: '', phone: '', province: '', city: '', district: '', detail: '' }
  defaultSwitch.value = false
  editingID.value = 0
  regionValues.value = []
}

function openEditor(addr?: any) {
  if (!addr) {
    resetForm()
    showEditor.value = true
    return
  }
  editingID.value = Number(addr.id)
  form.value = {
    name: addr.name || '',
    phone: addr.phone || '',
    province: addr.province || '',
    city: addr.city || '',
    district: addr.district || '',
    detail: addr.detail || '',
  }
  defaultSwitch.value = Number(addr.is_default || 0) === 1
  regionValues.value = [form.value.province, form.value.city, form.value.district].filter(Boolean)
  showEditor.value = true
}

function closeEditor() {
  showEditor.value = false
  resetForm()
}

function onRegionConfirm(values: string[]) {
  showRegionPicker.value = false
  if (!Array.isArray(values) || values.length < 3) return
  form.value.province = values[0] || ''
  form.value.city = values[1] || ''
  form.value.district = values[2] || ''
}

function validateForm() {
  if (!form.value.name.trim()) return '请输入收货人姓名'
  if (!/^1\d{10}$/.test(form.value.phone)) return '请输入正确手机号'
  if (!form.value.province.trim() || !form.value.city.trim() || !form.value.district.trim()) return '请选择完整省市区'
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

  try {
    if (editingID.value) {
      await put(`/api/v1/addresses/${editingID.value}`, payload)
      uni.showToast({ title: '修改成功', icon: 'success' })
    } else {
      await post('/api/v1/addresses', payload)
      uni.showToast({ title: '新增成功', icon: 'success' })
    }
    await loadAddresses()
    closeEditor()
  } catch {
    uni.showToast({ title: '保存失败，请稍后重试', icon: 'none' })
  } finally {
    saving.value = false
  }
}

async function removeAddress(addr: any) {
  const id = Number(addr?.id || 0)
  if (!id) return
  try {
    await del(`/api/v1/addresses/${id}`)
    uni.showToast({ title: '删除成功', icon: 'success' })
    await loadAddresses()
  } catch {
    uni.showToast({ title: '删除失败', icon: 'none' })
  }
}

onMounted(loadAddresses)
</script>
