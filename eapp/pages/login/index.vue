<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const authStore = useAuthStore()

const loading = ref(false)
const form = reactive({ username: '', password: '' })

async function onSubmit() {
  if (!form.username.trim()) {
    uni.showToast({ title: t('login.usernameRequired'), icon: 'none' })
    return
  }
  if (!form.password.trim()) {
    uni.showToast({ title: t('login.passwordRequired'), icon: 'none' })
    return
  }
  loading.value = true
  try {
    await authStore.loginAction(form.username.trim(), form.password)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <view class="login-page">
    <view class="login-card">
      <view class="title">{{ t('login.title') }}</view>
      <up-input v-model="form.username" :placeholder="t('login.username')" clearable />
      <view class="mt-20rpx" />
      <up-input v-model="form.password" type="password" password :placeholder="t('login.password')" clearable />
      <view class="mt-32rpx" />
      <up-button type="primary" shape="circle" :loading="loading" @click="onSubmit">{{ t('login.submit') }}</up-button>
    </view>
  </view>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32rpx;
  box-sizing: border-box;
  background: linear-gradient(180deg, #eff6ff, #f8fafc);
}

.login-card {
  width: 100%;
  background: #fff;
  border-radius: 28rpx;
  padding: 32rpx;
  box-sizing: border-box;
  border: 1px solid #dbeafe;
}

.title {
  font-size: 34rpx;
  font-weight: 700;
  margin-bottom: 28rpx;
}
</style>
