import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginAPI } from '@/api/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')

  async function loginAction(username: string, password: string) {
    const data = await loginAPI(username, password)
    token.value = data.token
    localStorage.setItem('admin_token', data.token)
    router.push('/dashboard')
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('admin_token')
    router.push('/login')
  }

  return { token, loginAction, logout }
})
