<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { getSiteSettings, updateSiteSettings } from '@/api/system'

const loading = ref(false)
const saving = ref(false)
const form = reactive<any>({
  shop_name: '',
  service_phone: '',
  order_auto_cancel_minutes: 30,
})

async function loadData() {
  loading.value = true
  try {
    const data: any = await getSiteSettings()
    Object.assign(form, data || {})
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    await updateSiteSettings({ ...form })
    uni.showToast({ title: '保存成功', icon: 'success' })
  } finally {
    saving.value = false
  }
}

onShow(loadData)
</script>

<template>
  <view class="page">
    <view v-if="loading" class="empty">加载中...</view>
    <view v-else class="card">
      <up-form>
        <up-form-item label="店铺名称"><up-input v-model="form.shop_name" /></up-form-item>
        <up-form-item label="客服电话"><up-input v-model="form.service_phone" inputmode="tel" /></up-form-item>
        <up-form-item label="自动取消（分钟）"><up-input v-model="form.order_auto_cancel_minutes" type="number" inputmode="numeric" /></up-form-item>
      </up-form>
      <view class="mt-16rpx" />
      <up-button type="primary" :loading="saving" @click="save">保存</up-button>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
</style>
