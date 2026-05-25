<template>
  <div class="min-h-screen bg-slate-100 flex items-center justify-center">
    <div class="bg-white rounded-2xl shadow-sm p-10 w-96">
      <div class="flex flex-col items-center mb-8">
        <img src="/lyshop-wordmark.svg" alt="LYShop" class="h-16 w-auto" />
      </div>
      <form @submit.prevent="handleLogin" class="space-y-4">
        <input
          v-model="form.username"
          type="text"
          :placeholder="$t('login.username')"
          class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500"
        />
        <input
          v-model="form.password"
          type="password"
          :placeholder="$t('login.password')"
          class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500"
        />
        <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-blue-700 hover:bg-blue-600 text-white rounded-xl py-3 text-sm font-medium transition disabled:opacity-60"
        >
          {{ loading ? $t('login.loggingIn') : $t('login.login') }}
        </button>
      </form>
      <p v-if="isMock" class="text-center text-xs text-gray-400 mt-4">
        {{ $t('login.demoHint') }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const isMock = import.meta.env.VITE_MOCK === 'true'
const auth = useAuthStore()
const form = ref({ username: '', password: '' })
const error = ref('')
const loading = ref(false)

onMounted(() => {
  if (isMock) {
    form.value.username = 'admin'
    form.value.password = 'admin123'
  }
})

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.loginAction(form.value.username, form.value.password)
  } catch (e: any) {
    error.value = e.message || t('login.loginFailed')
  } finally {
    loading.value = false
  }
}
</script>
