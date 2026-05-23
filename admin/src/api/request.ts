import axios from 'axios'

const MOCK_ENABLED = import.meta.env.VITE_MOCK === 'true'

const request = axios.create({
  baseURL: '/admin/api',
  timeout: 30000,
})

request.interceptors.request.use(config => {
  const token = localStorage.getItem('admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

request.interceptors.response.use(
  res => {
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
    const result = matchMock(method, url, config.params)

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
