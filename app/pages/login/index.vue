<template>
  <view class="min-h-screen bg-gray-50 flex flex-col items-center justify-center px-40rpx">
    <view class="w-full" style="max-width: 680rpx;">
      <!-- Logo -->
      <view class="flex flex-col items-center mb-80rpx">
        <view class="w-120rpx h-120rpx rounded-full bg-blue-700 flex items-center justify-center mb-24rpx">
          <text class="text-white text-40rpx font-700">L</text>
        </view>
        <text class="text-36rpx font-700 text-gray-800">LYShop</text>
        <text class="text-26rpx text-gray-400 mt-8rpx">开源插件化商城</text>
      </view>

      <!-- Form -->
      <view class="bg-white rounded-24rpx p-40rpx" style="box-shadow: 0 4rpx 24rpx rgba(0,0,0,0.06);">
        <view class="mb-30rpx">
          <u-input
            v-model="form.phone"
            placeholder="请输入手机号"
            type="number"
            :maxlength="11"
            border="surround"
            shape="circle"
            prefixIcon="phone"
          />
        </view>
        <view class="flex items-center gap-16rpx mb-40rpx">
          <view class="flex-1">
            <u-input
              v-model="form.code"
              placeholder="验证码"
              type="number"
              :maxlength="6"
              border="surround"
              shape="circle"
              prefixIcon="lock"
            />
          </view>
          <u-button
            size="small"
            :disabled="countdown > 0"
            :text="countdown > 0 ? `${countdown}s` : '获取验证码'"
            @click="sendCode"
            type="primary"
            plain
            shape="circle"
          />
        </view>
        <u-button
          type="primary"
          :loading="loading"
          text="登 录"
          @click="handleLogin"
          shape="circle"
          :custom-style="{height: '88rpx', fontSize: '30rpx'}"
        />
      </view>

      <view class="mt-40rpx text-center">
        <text class="text-24rpx text-gray-400">演示模式：输入任意手机号即可体验</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { post } from '@/utils/request'

const form = ref({ phone: '', code: '' })
const loading = ref(false)
const countdown = ref(0)

async function sendCode() {
  if (!form.value.phone || form.value.phone.length !== 11) {
    uni.showToast({ title: '请输入正确手机号', icon: 'none' })
    return
  }
  try {
    const data = await post<{ dev_code: string }>('/api/v1/auth/sms/send', {
      phone: form.value.phone
    })
    if (data?.dev_code) form.value.code = data.dev_code
  } catch {}

  countdown.value = 60
  const t = setInterval(() => {
    if (--countdown.value <= 0) clearInterval(t)
  }, 1000)
}

async function handleLogin() {
  loading.value = true
  try {
    const data = await post<{ token: string }>('/api/v1/auth/sms/login', form.value)
    uni.setStorageSync('user_token', data.token)
    uni.switchTab({ url: '/pages/index/index' })
  } catch {} finally {
    loading.value = false
  }
}
</script>
