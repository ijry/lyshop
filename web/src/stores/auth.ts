import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('user_token') || '')
  const isLoggedIn = computed(() => !!token.value)

  function setToken(t: string) {
    token.value = t
    localStorage.setItem('user_token', t)
  }
  function logout() {
    token.value = ''
    localStorage.removeItem('user_token')
  }

  return { token, isLoggedIn, setToken, logout }
})
