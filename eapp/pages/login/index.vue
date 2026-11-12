<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppIcon from '@/components/common/AppIcon.vue'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const authStore = useAuthStore()

const loading = ref(false)
const remember = ref(!!uni.getStorageSync('eapp_login_remember'))
const form = reactive({
  username: String(uni.getStorageSync('eapp_login_username') || ''),
  password: '',
})

const inputCustomStyle = {
  backgroundColor: 'transparent',
  border: 'none',
  padding: '0',
  height: '72rpx',
}

function toggleRemember() {
  remember.value = !remember.value
  if (!remember.value) {
    uni.removeStorageSync('eapp_login_username')
    uni.removeStorageSync('eapp_login_remember')
  }
}

async function onSubmit() {
  if (!form.username.trim()) {
    uni.showToast({ title: t('login.usernameRequired'), icon: 'none' })
    return
  }
  if (!form.password.trim()) {
    uni.showToast({ title: t('login.passwordRequired'), icon: 'none' })
    return
  }
  if (loading.value) return
  loading.value = true
  try {
    if (remember.value) {
      uni.setStorageSync('eapp_login_username', form.username.trim())
      uni.setStorageSync('eapp_login_remember', '1')
    }
    await authStore.loginAction(form.username.trim(), form.password)
  } finally {
    loading.value = false
  }
}

function onHelp() {
  uni.showModal({
    title: '登录帮助',
    content: '请使用管理后台分配的商家管理员账号。若忘记密码，请联系超级管理员重置。',
    showCancel: false,
    confirmText: '我知道了',
  })
}
</script>

<template>
  <view class="login-page">
    <view class="login-hero">
      <view class="login-hero__mark">
        <AppIcon name="store" :size="56" color="white" :stroke-width="1.5" />
      </view>
      <text class="login-hero__title">LYShop 商家工作台</text>
      <text class="login-hero__sub">MERCHANT CONSOLE</text>
      <text class="login-hero__desc">订单履约 · 商品运营 · 数据洞察</text>
    </view>

    <view class="login-card">
      <view class="login-card__head">
        <text class="login-card__title">账号登录</text>
        <text class="login-card__hint">请使用商家管理员账号登录</text>
      </view>

      <view class="login-field">
        <text class="login-field__label">{{ t('login.username') }}</text>
        <view class="login-field__box">
          <AppIcon name="user" :size="32" color="muted" />
          <view class="login-field__input">
            <up-input
              v-model="form.username"
              :placeholder="`请输入${t('login.username')}`"
              :custom-style="inputCustomStyle"
              font-size="28rpx"
              clearable
              @confirm="onSubmit"
            />
          </view>
        </view>
      </view>

      <view class="login-field">
        <text class="login-field__label">{{ t('login.password') }}</text>
        <view class="login-field__box">
          <AppIcon name="lock" :size="32" color="muted" />
          <view class="login-field__input">
            <up-input
              v-model="form.password"
              type="password"
              :placeholder="`请输入${t('login.password')}`"
              :custom-style="inputCustomStyle"
              font-size="28rpx"
              clearable
              @confirm="onSubmit"
            />
          </view>
        </view>
      </view>

      <view class="login-extra">
        <view class="login-remember" hover-class="eapp-hover" @click="toggleRemember">
          <view :class="['login-check', remember && 'login-check--on']">
            <AppIcon v-if="remember" name="check" :size="18" color="white" :stroke-width="2.4" />
          </view>
          <text class="login-remember__text">记住账号</text>
        </view>
        <text class="login-help" @click="onHelp">登录帮助</text>
      </view>

      <view class="login-submit">
        <up-button
          type="primary"
          :loading="loading"
          :disabled="loading"
          :custom-style="{ height: '88rpx', borderRadius: '12rpx', fontSize: '30rpx', fontWeight: '600', letterSpacing: '2rpx' }"
          @click="onSubmit"
        >
          {{ t('login.submit') }}
        </up-button>
      </view>

      <view class="login-card__footer">
        <AppIcon name="shield" :size="22" color="faint" />
        <text class="login-card__footer-text">账号信息加密传输，仅用于商家后台登录</text>
      </view>
    </view>

    <text class="login-page__foot">© LYShop · 商家管理后台</text>
  </view>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  background: var(--eapp-bg);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 40rpx calc(env(safe-area-inset-bottom) + 40rpx);
  box-sizing: border-box;
}

.login-hero {
  width: 100%;
  padding: calc(env(safe-area-inset-top) + 88rpx) 40rpx 120rpx;
  margin: 0 -40rpx;
  background: var(--eapp-primary-deep);
  display: flex;
  flex-direction: column;
  align-items: center;
  box-sizing: border-box;
}

.login-hero__mark {
  width: 104rpx;
  height: 104rpx;
  border-radius: var(--eapp-radius-xl);
  border: 1rpx solid rgba(255, 255, 255, 0.35);
  background: rgba(255, 255, 255, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 28rpx;
  box-sizing: border-box;
}

.login-hero__title {
  font-size: 40rpx;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 1rpx;
}

.login-hero__sub {
  margin-top: 6rpx;
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.65);
  letter-spacing: 4rpx;
}

.login-hero__desc {
  margin-top: 24rpx;
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
  letter-spacing: 1rpx;
}

.login-card {
  width: 100%;
  margin-top: -72rpx;
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-xl);
  box-shadow: var(--eapp-shadow-md);
  padding: 40rpx 36rpx 28rpx;
  box-sizing: border-box;
}

.login-card__head {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
  margin-bottom: 36rpx;
}

.login-card__title {
  font-size: 36rpx;
  font-weight: 700;
  color: var(--eapp-text);
}

.login-card__hint {
  font-size: 24rpx;
  color: var(--eapp-text-muted);
}

.login-field {
  margin-bottom: 28rpx;
}

.login-field__label {
  display: block;
  font-size: 26rpx;
  font-weight: 600;
  color: var(--eapp-text-secondary);
  margin-bottom: 12rpx;
}

.login-field__box {
  display: flex;
  align-items: center;
  gap: 16rpx;
  height: 88rpx;
  padding: 0 24rpx;
  border: 1rpx solid var(--eapp-border-strong);
  border-radius: var(--eapp-radius);
  background: var(--eapp-surface-subtle);
  box-sizing: border-box;
}

.login-field__input {
  flex: 1;
  min-width: 0;
}

.login-extra {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 4rpx 0 32rpx;
}

.login-remember {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.login-check {
  width: 32rpx;
  height: 32rpx;
  border: 1rpx solid var(--eapp-border-strong);
  border-radius: 6rpx;
  background: var(--eapp-card);
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
}

.login-check--on {
  background: var(--eapp-primary);
  border-color: var(--eapp-primary);
}

.login-remember__text,
.login-help {
  font-size: 24rpx;
  color: var(--eapp-text-muted);
}

.login-help {
  color: var(--eapp-primary);
}

.login-submit {
  width: 100%;
}

.login-card__footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  margin-top: 32rpx;
  padding-top: 24rpx;
  border-top: 1rpx solid var(--eapp-divider);
}

.login-card__footer-text {
  font-size: 20rpx;
  color: var(--eapp-text-faint);
}

.login-page__foot {
  margin-top: auto;
  padding-top: 48rpx;
  font-size: 20rpx;
  color: var(--eapp-text-faint);
  letter-spacing: 1rpx;
}
</style>
