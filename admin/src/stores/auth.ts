import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginAPI } from '@/api/auth'
import router from '@/router'

function parsePermsFromToken(token: string): string[] {
  try {
    const payloadSeg = String(token || '').split('.')[1] || ''
    if (!payloadSeg) return []
    const normalized = payloadSeg.replace(/-/g, '+').replace(/_/g, '/')
    const padded = normalized + '='.repeat((4 - normalized.length % 4) % 4)
    const payload = JSON.parse(atob(padded))
    const perms = Array.isArray(payload?.perms) ? payload.perms : []
    return perms.map((item: any) => String(item || ''))
  } catch {
    return []
  }
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const perms = ref<string[]>(parsePermsFromToken(token.value))

  async function loginAction(username: string, password: string) {
    const data = await loginAPI(username, password)
    token.value = data.token
    perms.value = parsePermsFromToken(data.token)
    localStorage.setItem('admin_token', data.token)
    router.push('/dashboard')
  }

  function hasPermission(permission: string): boolean {
    if (!permission) return true
    return perms.value.includes('*') || perms.value.includes(permission)
  }

  function logout() {
    token.value = ''
    perms.value = []
    localStorage.removeItem('admin_token')
    router.push('/login')
  }

  return { token, perms, hasPermission, loginAction, logout }
})
