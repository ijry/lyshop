<script setup lang="ts">
import { ref } from 'vue'
import { applyDistributor } from '@/api/distribution'

const form = ref({
  parent_id: 0,
  real_name: '',
  id_card: '',
  phone: ''
})
const submitting = ref(false)

async function handleSubmit() {
  if (!form.value.real_name) {
    uni.showToast({ title: '请输入真实姓名', icon: 'none' })
    return
  }
  if (!form.value.id_card) {
    uni.showToast({ title: '请输入身份证号', icon: 'none' })
    return
  }
  if (!form.value.phone) {
    uni.showToast({ title: '请输入手机号', icon: 'none' })
    return
  }

  submitting.value = true
  try {
    await applyDistributor(form.value)
    uni.showToast({ title: '申请成功，等待审核', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (error: any) {
    uni.showToast({ title: error.message || '申请失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <view class="page">
    <view class="form">
      <view class="form-item">
        <text class="label">真实姓名 *</text>
        <input v-model="form.real_name" class="input" placeholder="请输入真实姓名" />
      </view>
      <view class="form-item">
        <text class="label">身份证号 *</text>
        <input v-model="form.id_card" class="input" placeholder="请输入身份证号" />
      </view>
      <view class="form-item">
        <text class="label">手机号 *</text>
        <input v-model="form.phone" class="input" type="number" placeholder="请输入手机号" />
      </view>
      <view class="form-item">
        <text class="label">推荐人ID（可选）</text>
        <input v-model.number="form.parent_id" class="input" type="number" placeholder="请输入推荐人ID" />
      </view>
    </view>

    <button class="submit-btn" :disabled="submitting" @click="handleSubmit">
      {{ submitting ? '提交中...' : '提交申请' }}
    </button>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: #f5f5f5; padding: 20rpx; }
.form { background: #fff; border-radius: 16rpx; padding: 20rpx; margin-bottom: 20rpx; }
.form-item { margin-bottom: 30rpx; }
.label { display: block; font-size: 28rpx; color: #333; margin-bottom: 10rpx; }
.input { width: 100%; height: 80rpx; padding: 0 20rpx; border: 1px solid #e5e5e5; border-radius: 8rpx; font-size: 28rpx; }
.submit-btn { width: 100%; height: 88rpx; line-height: 88rpx; background: #dc2626; color: #fff; border-radius: 16rpx; border: none; font-size: 32rpx; }
.submit-btn[disabled] { opacity: 0.6; }
</style>
