import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginAPI } from '@/api/auth'
import router from '@/router'

type AdminAccount = {
  username: string
  token: string
  avatar: string
}

const ACCOUNT_LIST_KEY = 'admin_accounts'

function readAccounts(): AdminAccount[] {
  try {
    const rows = JSON.parse(localStorage.getItem(ACCOUNT_LIST_KEY) || '[]')
    return Array.isArray(rows) ? rows : []
  } catch {
    return []
  }
}

function persistAccounts(accounts: AdminAccount[]) {
  localStorage.setItem(ACCOUNT_LIST_KEY, JSON.stringify(accounts))
}

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
  const currentUsername = ref(localStorage.getItem('admin_username') || '')
  const accounts = ref<AdminAccount[]>(readAccounts())

  function buildAccountProfile(username: string, tokenValue: string): AdminAccount {
    const normalized = String(username || 'admin').trim() || 'admin'
    return {
      username: normalized,
      token: tokenValue,
      avatar: `https://api.dicebear.com/9.x/initials/svg?seed=${encodeURIComponent(normalized)}`,
    }
  }

  function syncAccount(profile: AdminAccount) {
    const next = accounts.value.filter((item) => item.username !== profile.username)
    next.unshift(profile)
    accounts.value = next.slice(0, 8)
    persistAccounts(accounts.value)
  }

  async function loginAction(username: string, password: string) {
    const data = await loginAPI(username, password)
    token.value = data.token
    perms.value = parsePermsFromToken(data.token)
    currentUsername.value = username
    localStorage.setItem('admin_token', data.token)
    localStorage.setItem('admin_username', username)
    syncAccount(buildAccountProfile(username, data.token))
    router.push('/dashboard')
  }

  function switchAccount(username: string) {
    const target = accounts.value.find((item) => item.username === username)
    if (!target) return
    token.value = target.token
    perms.value = parsePermsFromToken(target.token)
    currentUsername.value = target.username
    localStorage.setItem('admin_token', target.token)
    localStorage.setItem('admin_username', target.username)
    syncAccount(target)
    router.push('/dashboard')
  }

  function hasPermission(permission: string): boolean {
    if (!permission) return true
    return perms.value.includes('*') || perms.value.includes(permission)
  }

  function logout() {
    token.value = ''
    perms.value = []
    currentUsername.value = ''
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_username')
    router.push('/login')
  }

  return { token, perms, currentUsername, accounts, hasPermission, loginAction, switchAccount, logout }
})
