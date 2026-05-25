import axios from 'axios'
import type { AxiosRequestConfig } from 'axios'
import { i18n } from '@/locales'

const t = (key: string) => i18n.global.t(key)

const MOCK_ENABLED = import.meta.env.VITE_MOCK === 'true'
const BASE_URL = MOCK_ENABLED ? '' : (import.meta.env.VITE_API_URL || '')

const http = axios.create({ baseURL: BASE_URL, timeout: 30000 })

http.interceptors.request.use(config => {
  const token = localStorage.getItem('user_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

http.interceptors.response.use(
  res => {
    const { code, msg, data } = res.data
    if (code !== 0) return Promise.reject(new Error(msg || t('api.requestFailed')))
    return data
  },
  err => Promise.reject(err)
)

async function mockRequest<T>(method: string, url: string, params?: any): Promise<T> {
  const { matchMock } = await import('@/mock/index')
  const result = matchMock(method, url, params)
  await new Promise(r => setTimeout(r, 100 + Math.random() * 200))
  if (result.matched) return (result.data ?? null) as T
  console.warn(`[Mock] No data for: ${method} ${url}`)
  return null as T
}

export async function get<T = any>(url: string, params?: any): Promise<T> {
  if (MOCK_ENABLED) return mockRequest<T>('GET', url, params)
  return http.get(url, { params }) as Promise<T>
}

export async function post<T = any>(url: string, data?: any): Promise<T> {
  if (MOCK_ENABLED) return mockRequest<T>('POST', url, data)
  return http.post(url, data) as Promise<T>
}

export async function put<T = any>(url: string, data?: any): Promise<T> {
  if (MOCK_ENABLED) return mockRequest<T>('PUT', url, data)
  return http.put(url, data) as Promise<T>
}

export async function del<T = any>(url: string, data?: any): Promise<T> {
  if (MOCK_ENABLED) return mockRequest<T>('DELETE', url, data)
  return http.delete(url, { data }) as Promise<T>
}

export async function upload<T = any>(url: string, file: File): Promise<T> {
  if (MOCK_ENABLED) return mockRequest<T>('POST', url, { name: file.name, size: file.size })
  const form = new FormData()
  form.append('file', file)
  return http.post(url, form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }) as Promise<T>
}
