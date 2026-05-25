<template>
  <view style="min-height: 100vh; background: #f5f5f5;">
    <u-navbar :title="$t('address.title')" :placeholder="true" />

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
              {{ $t('address.default') }}
            </view>
            <text style="font-size:12px;color:#3b82f6;" @click="openEditor(addr)">{{ $t('address.edit') }}</text>
            <text style="font-size:12px;color:#ef4444;" @click="removeAddress(addr)">{{ $t('address.delete') }}</text>
          </view>
        </view>
        <text style="font-size: 13px; color: #999; line-height: 1.5;">
          {{ addr.province }}{{ addr.city }}{{ addr.district }} {{ addr.detail }}
        </text>
      </view>

      <view v-if="!addresses.length" style="text-align: center; padding: 80px 0; color: #999; font-size: 14px;">
        {{ $t('address.empty') }}
      </view>
    </view>

    <view style="padding: 16px; position: fixed; bottom: 0; left: 0; right: 0; background: #f5f5f5;">
      <u-button
        type="primary"
        :text="$t('address.addNew')"
        shape="circle"
        :custom-style="{background: '#dc2626', borderColor: '#dc2626'}"
        @click="openEditor()"
      />
    </view>

    <u-popup :show="showEditor" mode="bottom" round="20" @close="closeEditor">
      <view style="padding: 20px 16px 24px;">
        <text style="display: block; font-size: 16px; font-weight: 700; margin-bottom: 16px;">{{ editingID ? $t('address.edit') : $t('address.addNew') }}</text>
        <view style="display: flex; flex-direction: column; gap: 12px;">
          <u-input v-model="form.name" :placeholder="$t('address.recipientName')" border="surround" />
          <u-input v-model="form.phone" :placeholder="$t('address.phone')" border="surround" maxlength="11" type="number" />
          <view @click="showRegionPicker = true">
            <u-input :modelValue="regionText" :placeholder="$t('address.region')" border="surround" readonly />
          </view>
          <u-input v-model="form.detail" :placeholder="$t('address.detailAddress')" border="surround" />
          <view style="display: flex; align-items: center; justify-content: space-between;">
            <text style="font-size: 13px; color: #666;">{{ $t('address.setDefault') }}</text>
            <u-switch v-model="defaultSwitch" activeColor="#dc2626" />
          </view>
        </view>
        <view style="display: flex; gap: 12px; margin-top: 18px;">
          <u-button :text="$t('address.cancel')" shape="circle" @click="closeEditor" />
          <u-button
            :text="$t('address.save')"
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
import { useI18n } from 'vue-i18n'
import { del, get, post, put } from '@/utils/request'

const { t } = useI18n()

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
  if (!form.value.name.trim()) return t('address.nameRequired')
  if (!/^1\d{10}$/.test(form.value.phone)) return t('address.phoneInvalid')
  if (!form.value.province.trim() || !form.value.city.trim() || !form.value.district.trim()) return t('address.regionRequired')
  if (!form.value.detail.trim()) return t('address.addressRequired')
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
      uni.showToast({ title: t('address.updateSuccess'), icon: 'success' })
    } else {
      await post('/api/v1/addresses', payload)
      uni.showToast({ title: t('address.addSuccess'), icon: 'success' })
    }
    await loadAddresses()
    closeEditor()
  } catch {
    uni.showToast({ title: t('address.saveFailed'), icon: 'none' })
  } finally {
    saving.value = false
  }
}

async function removeAddress(addr: any) {
  const id = Number(addr?.id || 0)
  if (!id) return
  try {
    await del(`/api/v1/addresses/${id}`)
    uni.showToast({ title: t('address.deleteSuccess'), icon: 'success' })
    await loadAddresses()
  } catch {
    uni.showToast({ title: t('address.deleteFailed'), icon: 'none' })
  }
}

onMounted(loadAddresses)
</script>
