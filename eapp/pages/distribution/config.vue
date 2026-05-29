<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { getDistributionConfig, updateDistributionConfig } from '@/api/distribution'

const { t } = useI18n()
const loading = ref(false)
const form = reactive({ level1_rate: '', level2_rate: '' })

async function loadConfig() {
  loading.value = true
  try {
    const res = await getDistributionConfig()
    form.level1_rate = String(Number(res?.level1_rate || 0) * 100)
    form.level2_rate = String(Number(res?.level2_rate || 0) * 100)
  } finally {
    loading.value = false
  }
}

async function save() {
  const level1 = Number(form.level1_rate) / 100
  const level2 = Number(form.level2_rate) / 100
  if (isNaN(level1) || isNaN(level2) || level1 < 0 || level2 < 0) {
    uni.showToast({ title: '请输入有效比例', icon: 'none' })
    return
  }
  await updateDistributionConfig({ level1_rate: level1, level2_rate: level2 })
  uni.showToast({ title: '保存成功', icon: 'success' })
}

onShow(() => loadConfig())
</script>

<template>
  <view class="page">
    <view class="card">
      <view class="form-title">{{ t('distribution.distributionConfig') }}</view>
      <view class="form-item">
        <text class="label">{{ t('distribution.level1Rate') }}（%）</text>
        <up-input v-model="form.level1_rate" type="digit" placeholder="例如 10" />
      </view>
      <view class="form-item">
        <text class="label">{{ t('distribution.level2Rate') }}（%）</text>
        <up-input v-model="form.level2_rate" type="digit" placeholder="例如 5" />
      </view>
      <view class="mt-lg" />
      <up-button type="primary" @click="save">{{ t('common.save') }}</up-button>
    </view>
    <view v-if="loading" class="loading">加载中…</view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 24rpx; }
.form-title { font-size: 30rpx; font-weight: 700; margin-bottom: 16rpx; }
.form-item { margin-top: 14rpx; }
.label { font-size: 26rpx; color: #475569; margin-bottom: 6rpx; display: block; }
.mt-lg { margin-top: 20rpx; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>
