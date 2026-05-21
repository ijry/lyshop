<template>
  <div class="min-h-screen bg-slate-100 flex items-center justify-center">
    <div class="bg-white rounded-2xl shadow-sm p-10 w-96">
      <h1 class="text-2xl font-bold text-slate-800 mb-8 text-center">lyshop 管理后台</h1>
      <form @submit.prevent="handleLogin" class="space-y-4">
        <input
          v-model="form.username"
          type="text"
          placeholder="用户名"
          class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500"
        />
        <input
          v-model="form.password"
          type="password"
          placeholder="密码"
          class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500"
        />
        <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-blue-700 hover:bg-blue-600 text-white rounded-xl py-3 text-sm font-medium transition disabled:opacity-60"
        >
          {{ loading ? '登录中...' : '登 录' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const form = ref({ username: '', password: '' })
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.loginAction(form.value.username, form.value.password)
  } catch (e: any) {
    error.value = e.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>
