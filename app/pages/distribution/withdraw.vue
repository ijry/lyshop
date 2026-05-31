<script setup lang="ts">
import { ref } from 'vue'
import { createWithdrawal } from '@/api/distribution'

const form = ref({
  amount: 0,
  bank_name: '',
  bank_account: '',
  account_name: ''
})
const submitting = ref(false)

async function handleSubmit() {
  if (!form.value.amount || form.value.amount <= 0) {
    uni.showToast({ title: '请输入提现金额', icon: 'none' })
    return
  }
  if (!form.value.bank_name) {
    uni.showToast({ title: '请输入银行名称', icon: 'none' })
    return
  }
  if (!form.value.bank_account) {
    uni.showToast({ title: '请输入银行账号', icon: 'none' })
    return
  }
  if (!form.value.account_name) {
    uni.showToast({ title: '请输入账户名', icon: 'none' })
    return
  }

  submitting.value = true
  try {
    await createWithdrawal(form.value)
    uni.showToast({ title: '提现申请已提交', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (error: any) {
    uni.showToast({ title: error.message || '提交失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <view class="page">
    <view class="form">
      <view class="form-item">
        <text class="label">提现金额 *</text>
        <input v-model.number="form.amount" class="input" type="digit" placeholder="请输入提现金额" />
      </view>
      <view class="form-item">
        <text class="label">银行名称 *</text>
        <input v-model="form.bank_name" class="input" placeholder="例如：中国工商银行" />
      </view>
      <view class="form-item">
        <text class="label">银行账号 *</text>
        <input v-model="form.bank_account" class="input" type="number" placeholder="请输入银行账号" />
      </view>
      <view class="form-item">
        <text class="label">账户名 *</text>
        <input v-model="form.account_name" class="input" placeholder="请输入账户名" />
      </view>
    </view>

    <view class="tips">
      <text class="tips-title">提现说明：</text>
      <text class="tips-text">1. 提现金额需大于最低提现额度</text>
      <text class="tips-text">2. 提现申请提交后需等待审核</text>
      <text class="tips-text">3. 审核通过后将在3-5个工作日内到账</text>
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

.tips { background: #fff7ed; border-radius: 16rpx; padding: 20rpx; margin-bottom: 20rpx; }
.tips-title { display: block; font-size: 28rpx; color: #f97316; font-weight: bold; margin-bottom: 10rpx; }
.tips-text { display: block; font-size: 24rpx; color: #92400e; margin-bottom: 8rpx; }

.submit-btn { width: 100%; height: 88rpx; line-height: 88rpx; background: #dc2626; color: #fff; border-radius: 16rpx; border: none; font-size: 32rpx; }
.submit-btn[disabled] { opacity: 0.6; }
</style>
