import { getStorage, removeStorage } from './storage'

const API_ROOT = '/admin/api'
const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

function getToken() {
  return String(getStorage('eapp_admin_token') || '')
}

function normalizeUrl(url: string) {
  if (!url.startsWith('/')) return `${API_ROOT}/${url}`
  if (url.startsWith(API_ROOT)) return url
  return `${API_ROOT}${url}`
}

function handleUnauthorized() {
  removeStorage('eapp_admin_token')
  removeStorage('eapp_admin_username')
  const pages = getCurrentPages()
  const route = pages[pages.length - 1]?.route || ''
  if (route !== 'pages/login/index') {
    uni.reLaunch({ url: '/pages/login/index' })
  }
}

export function request<T = any>(options: UniNamespace.RequestOptions): Promise<T> {
  return new Promise((resolve, reject) => {
    const token = getToken()
    uni.request({
      ...options,
      url: BASE_URL + normalizeUrl(String(options.url || '')),
      header: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...(options.header || {}),
      },
      success(res) {
        const status = Number(res.statusCode || 0)
        if (status === 401) {
          handleUnauthorized()
          reject(new Error('未登录或登录已失效'))
          return
        }
        const payload = (res.data || {}) as any
        if (payload.code !== 0) {
          const msg = payload.msg || '请求失败'
          uni.showToast({ title: msg, icon: 'none' })
          reject(new Error(msg))
          return
        }
        resolve(payload.data)
      },
      fail(err) {
        reject(err)
      },
    })
  })
}

export const get = <T>(url: string, data?: any) => request<T>({ url, method: 'GET', data })
export const post = <T>(url: string, data?: any) => request<T>({ url, method: 'POST', data })
export const put = <T>(url: string, data?: any) => request<T>({ url, method: 'PUT', data })
export const del = <T>(url: string, data?: any) => request<T>({ url, method: 'DELETE', data })
