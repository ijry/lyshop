<template>
  <view class="login-page">
    <view class="login-shell">
      <view class="brand-block">
        <image src="/static/lyshop-wordmark.svg" mode="aspectFit" class="brand-logo" />
        <text class="brand-slogan">{{ $t('login.slogan') }}</text>
      </view>

      <view class="login-card">
        <view class="field-wrap">
          <text class="field-label">{{ $t('login.phone') }}</text>
          <view class="field-shell">
            <u-input
              v-model="form.phone"
              :placeholder="$t('login.phonePlaceholder')"
              type="number"
              :maxlength="11"
              border="none"
              shape="circle"
              prefixIcon="phone"
            />
          </view>
        </view>

        <view class="field-wrap">
          <text class="field-label">{{ $t('login.verifyCode') }}</text>
          <view class="verify-row">
            <view class="field-shell code-input">
              <u-input
                v-model="form.code"
                :placeholder="$t('login.codePlaceholder')"
                type="number"
                :maxlength="6"
                border="none"
                shape="circle"
                prefixIcon="lock"
              />
            </view>
            <u-button
              class="code-btn"
              size="small"
              :disabled="codeBtnDisabled"
              :text="countdown > 0 ? `${countdown}s` : $t('login.getCode')"
              @click="sendCode"
              type="primary"
              plain
              shape="circle"
              :custom-style="codeButtonStyle"
            />
          </view>
        </view>

        <u-button
          class="submit-btn"
          type="primary"
          :loading="loading"
          :disabled="loginDisabled"
          :text="$t('login.submit')"
          @click="handleLogin"
          shape="circle"
          :custom-style="submitButtonStyle"
        />
      </view>

      <view class="demo-tip">
        <text>{{ $t('login.demoHint') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { post } from '@/utils/request'

const { t } = useI18n()

const form = ref({ phone: '', code: '' })
const loading = ref(false)
const countdown = ref(0)

let countdownTimer: ReturnType<typeof setInterval> | null = null

const isPhoneValid = computed(() => /^1\d{10}$/.test(form.value.phone))
const codeBtnDisabled = computed(() => countdown.value > 0 || !isPhoneValid.value)
const loginDisabled = computed(() => loading.value || !isPhoneValid.value || form.value.code.length !== 6)

const codeButtonStyle = computed(() => ({
  height: '72rpx',
  minWidth: '220rpx',
  fontSize: '24rpx',
  color: codeBtnDisabled.value ? '#94A3B8' : '#2A4DBF',
  borderColor: codeBtnDisabled.value ? '#D5DEED' : '#3F62D9',
  background: codeBtnDisabled.value ? '#EEF2FA' : '#F7F9FF'
}))

const submitButtonStyle = computed(() => ({
  height: '90rpx',
  fontSize: '30rpx',
  color: '#FFFFFF',
  border: 'none',
  background: loginDisabled.value ? '#A9B8DE' : 'linear-gradient(135deg, #3B61DA 0%, #2347B7 100%)',
  boxShadow: loginDisabled.value ? 'none' : '0 12rpx 24rpx rgba(39, 74, 184, 0.28)'
}))

function clearCountdownTimer() {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
}

function assertPhoneValid() {
  if (!isPhoneValid.value) {
    uni.showToast({ title: t('login.phoneInvalid'), icon: 'none' })
    return false
  }
  return true
}

async function sendCode() {
  if (!assertPhoneValid()) return
  if (countdown.value > 0) return
  try {
    const data = await post<{ dev_code: string }>('/api/v1/auth/sms/send', {
      phone: form.value.phone
    })
    if (data?.dev_code) form.value.code = data.dev_code
  } catch {}

  countdown.value = 60
  clearCountdownTimer()
  countdownTimer = setInterval(() => {
    if (countdown.value <= 1) {
      countdown.value = 0
      clearCountdownTimer()
      return
    }
    countdown.value -= 1
  }, 1000)
}

async function handleLogin() {
  if (!assertPhoneValid()) return
  if (form.value.code.length !== 6) {
    uni.showToast({ title: t('login.codeInvalid'), icon: 'none' })
    return
  }
  if (loading.value) return

  loading.value = true
  try {
    const data = await post<{ token: string }>('/api/v1/auth/sms/login', form.value)
    uni.setStorageSync('user_token', data.token)
    uni.switchTab({ url: '/pages/index/index' })
  } catch {} finally {
    loading.value = false
  }
}

onUnmounted(() => {
  clearCountdownTimer()
})
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 56rpx 40rpx;
  background: linear-gradient(180deg, #f8faff 0%, #f2f5fb 100%);
}

.login-shell {
  width: 100%;
  max-width: 680rpx;
}

.brand-block {
  margin-bottom: 68rpx;
  text-align: center;
}

.brand-logo {
  width: 320rpx;
  height: 120rpx;
  margin: 0 auto;
}

.brand-slogan {
  display: block;
  margin-top: 12rpx;
  color: #7384a8;
  font-size: 26rpx;
}

.login-card {
  padding: 40rpx;
  border-radius: 28rpx;
  background: #ffffff;
  box-shadow: 0 20rpx 48rpx rgba(36, 76, 168, 0.08), inset 0 2rpx 0 rgba(255, 255, 255, 0.65);
}

.field-wrap + .field-wrap {
  margin-top: 26rpx;
}

.field-label {
  display: block;
  margin-bottom: 14rpx;
  font-size: 24rpx;
  color: #5d6d90;
}

.field-shell {
  border-radius: 999rpx;
  border: 2rpx solid #d6dfef;
  background: #f9fbff;
  padding: 0 18rpx;
  box-shadow: inset 0 2rpx 6rpx rgba(38, 64, 132, 0.06);
}

.verify-row {
  display: flex;
  align-items: center;
  gap: 14rpx;
}

.code-input {
  flex: 1;
}

.code-btn {
  flex-shrink: 0;
}

.submit-btn {
  margin-top: 40rpx;
}

.demo-tip {
  margin-top: 34rpx;
  text-align: center;
  color: #8795b3;
  font-size: 24rpx;
}

:deep(.field-shell .u-input) {
  height: 74rpx;
}

:deep(.field-shell .u-input__content) {
  height: 74rpx;
}

:deep(.field-shell .u-input__content__field-wrapper__field) {
  font-size: 28rpx;
  color: #243356;
}

:deep(.field-shell .u-input__content__field-wrapper__field::placeholder) {
  color: #9aa8c5;
}

:deep(.field-shell .u-icon__icon) {
  color: #6f82aa !important;
}

:deep(.submit-btn .u-button__text),
:deep(.code-btn .u-button__text) {
  font-weight: 600;
}
</style>
