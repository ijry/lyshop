import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login } from '@/api/auth'
import { parsePermsFromToken } from '@/utils/permission'
import { getStorage, removeStorage, setStorage } from '@/utils/storage'

export const useAuthStore = defineStore('eapp-auth', () => {
  const token = ref(String(getStorage('eapp_admin_token') || ''))
  const username = ref(String(getStorage('eapp_admin_username') || ''))
  const perms = ref<string[]>(parsePermsFromToken(token.value))

  async function loginAction(account: string, password: string) {
    const data = await login(account, password)
    token.value = String(data?.token || '')
    username.value = account
    perms.value = parsePermsFromToken(token.value)

    setStorage('eapp_admin_token', token.value)
    setStorage('eapp_admin_username', account)
    uni.switchTab({ url: '/pages/dashboard/index' })
  }

  function logout() {
    token.value = ''
    username.value = ''
    perms.value = []
    removeStorage('eapp_admin_token')
    removeStorage('eapp_admin_username')
    uni.reLaunch({ url: '/pages/login/index' })
  }

  return {
    token,
    username,
    perms,
    loginAction,
    logout,
  }
})
