<template>
  <div class="min-h-[80vh] flex-center">
    <div class="w-full max-w-md mx-auto px-6">
      <!-- Logo -->
      <div class="text-center mb-10">
        <div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-red-600 to-red-500 flex-center mx-auto mb-4 shadow-lg shadow-red-600/20">
          <span class="text-white text-2xl font-bold">L</span>
        </div>
        <h1 class="text-2xl font-bold text-gray-900">{{ $t('login.title') }}</h1>
        <p class="text-sm text-gray-500 mt-2">{{ $t('login.subtitle') }}</p>
      </div>

      <!-- Form -->
      <div class="card p-8">
        <div class="mb-5">
          <label class="text-sm font-medium text-gray-700 mb-1.5 block">{{ $t('login.phone') }}</label>
          <input v-model="form.phone" type="tel" maxlength="11" :placeholder="$t('login.phonePlaceholder')" class="input-base" />
        </div>
        <div class="mb-6">
          <label class="text-sm font-medium text-gray-700 mb-1.5 block">{{ $t('login.code') }}</label>
          <div class="flex gap-3">
            <input v-model="form.code" type="text" maxlength="6" :placeholder="$t('login.codePlaceholder')" class="input-base flex-1" />
            <button @click="sendCode" :disabled="countdown > 0"
              class="btn-outline whitespace-nowrap !px-4 text-xs"
              :class="countdown > 0 ? 'opacity-50 cursor-not-allowed' : ''">
              {{ countdown > 0 ? `${countdown}${$t('login.second')}` : $t('login.getCode') }}
            </button>
          </div>
        </div>
        <button @click="handleLogin" :disabled="loading"
          class="btn-primary w-full !py-3 !text-base">
          {{ loading ? $t('login.submitting') : $t('login.submit') }}
        </button>
      </div>

      <p class="text-center text-xs text-gray-400 mt-6">{{ $t('login.demoTip') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { post } from '@/api/request'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const form = ref({ phone: '', code: '' })
const loading = ref(false)
const countdown = ref(0)

async function sendCode() {
  if (!form.value.phone || form.value.phone.length !== 11) return
  const data = await post<{ dev_code: string }>('/api/v1/auth/sms/send', { phone: form.value.phone })
  if (data?.dev_code) form.value.code = data.dev_code
  countdown.value = 60
  const t = setInterval(() => { if (--countdown.value <= 0) clearInterval(t) }, 1000)
}

async function handleLogin() {
  loading.value = true
  try {
    const data = await post<{ token: string }>('/api/v1/auth/sms/login', form.value)
    auth.setToken(data.token)
    router.push('/')
  } catch {} finally {
    loading.value = false
  }
}
</script>
