import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios'

// Define custom axios instance type that returns unwrapped data
interface CustomAxiosInstance extends Omit<AxiosInstance, 'get' | 'post' | 'put' | 'delete' | 'patch'> {
  get<T = any>(url: string, config?: any): Promise<T>
  post<T = any>(url: string, data?: any, config?: any): Promise<T>
  put<T = any>(url: string, data?: any, config?: any): Promise<T>
  delete<T = any>(url: string, config?: any): Promise<T>
  patch<T = any>(url: string, data?: any, config?: any): Promise<T>
}

const MOCK_ENABLED = import.meta.env.VITE_MOCK === 'true'

const request = axios.create({
  baseURL: '/admin/api',
  timeout: 30000,
}) as CustomAxiosInstance

request.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem('admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

request.interceptors.response.use(
  (res: AxiosResponse) => {
    const { code, msg, data } = res.data
    if (code !== 0) return Promise.reject(new Error(msg || '请求失败'))
    return data
  },
  err => Promise.reject(err)
)

// Mock interceptor: intercept all requests when VITE_MOCK=true
if (MOCK_ENABLED) {
  request.interceptors.request.use(async (config) => {
    const { matchMock } = await import('@/mock/index')
    const method = (config.method || 'GET').toUpperCase()
    const url = (config.baseURL || '') + (config.url || '')
    let payload: any = config.params
    if (method !== 'GET' && config.data) {
      if (typeof config.data === 'string') {
        try {
          payload = JSON.parse(config.data)
        } catch {
          payload = config.params
        }
      } else {
        payload = config.data
      }
    }
    const result = matchMock(method, url, payload)

    if (result.matched) {
      await new Promise(r => setTimeout(r, 100 + Math.random() * 200))
      config.adapter = () =>
        Promise.resolve({
          data: { code: 0, msg: 'success', data: result.data },
          status: 200,
          statusText: 'OK',
          headers: {},
          config,
        })
    }
    return config
  })
}

export default request
