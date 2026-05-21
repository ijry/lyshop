<template>
  <view class="min-h-screen bg-gray-50 flex flex-col items-center justify-center px-6">
    <view class="w-full" style="max-width: 360px;">
      <text class="text-2xl font-bold text-slate-800 block text-center mb-8">登录</text>

      <view class="bg-white rounded-2xl p-6 shadow-sm">
        <view class="mb-4">
          <u-input
            v-model="form.phone"
            placeholder="请输入手机号"
            type="number"
            :maxlength="11"
            border="surround"
          />
        </view>
        <view class="mb-6 flex gap-2">
          <view class="flex-1">
            <u-input
              v-model="form.code"
              placeholder="验证码"
              type="number"
              :maxlength="6"
              border="surround"
            />
          </view>
          <u-button
            size="small"
            :disabled="countdown > 0"
            :text="countdown > 0 ? `${countdown}s` : '获取验证码'"
            @click="sendCode"
            type="primary"
            plain
          />
        </view>
        <u-button
          type="primary"
          :loading="loading"
          text="登 录"
          @click="handleLogin"
        />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { get, post } from '@/utils/request'

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
