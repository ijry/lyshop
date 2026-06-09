// Mock mode: controlled by compile-time env variable VITE_MOCK
// Build with: npm run build:h5:demo (sets VITE_MOCK=true)
const MOCK_ENABLED = import.meta.env.VITE_MOCK === 'true'

const BASE_URL = MOCK_ENABLED ? '' : (import.meta.env.VITE_API_URL || 'http://localhost:8080')

function getToken(): string {
  if (MOCK_ENABLED) return 'demo_token'
  return uni.getStorageSync('user_token') || ''
}

export function request<T = any>(options: UniNamespace.RequestOptions): Promise<T> {
  // Mock intercept
  if (MOCK_ENABLED) {
    return mockRequest<T>(options.method || 'GET', options.url!, options.data)
  }

  return new Promise((resolve, reject) => {
    uni.request({
      ...options,
      url: BASE_URL + options.url,
      header: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${getToken()}`,
        ...(options.header || {})
      },
      success(res) {
        const data = res.data as any
        if (data.code !== 0) {
          uni.showToast({ title: data.msg || '请求失败', icon: 'none' })
          reject(new Error(data.msg))
        } else {
          resolve(data.data)
        }
      },
      fail(err) { reject(err) }
    })
  })
}

// Mock request handler — returns local data with simulated delay
async function mockRequest<T>(method: string, url: string, _data?: any): Promise<T> {
  const { matchMock } = await import('@/mock/index')
  const result = matchMock(method, url, _data)

  // Simulate network delay (100-300ms)
  await new Promise(r => setTimeout(r, 100 + Math.random() * 200))

  if (result.matched) {
    return (result.data ?? null) as T
  }

  // Unmatched GET: return empty
  console.warn(`[Mock] No mock data for: ${method} ${url}`)
  return null as T
}

export const get = <T>(url: string, data?: any) =>
  request<T>({ url, method: 'GET', data })

export const post = <T>(url: string, data?: any) =>
  request<T>({ url, method: 'POST', data })

export const put = <T>(url: string, data?: any) =>
  request<T>({ url, method: 'PUT', data })

export const del = <T>(url: string, data?: any) =>
  request<T>({ url, method: 'DELETE', data })

export function upload<T = any>(url: string, filePath: string, name = 'file', formData?: Record<string, any>): Promise<T> {
  if (MOCK_ENABLED) {
    return mockRequest<T>('POST', url, { name, filePath, ...(formData || {}) })
  }
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: BASE_URL + url,
      filePath,
      name,
      formData,
      header: {
        Authorization: `Bearer ${getToken()}`,
      },
      success(res) {
        try {
          const data = JSON.parse(String(res.data || '{}'))
          if (data.code !== 0) {
            uni.showToast({ title: data.msg || '上传失败', icon: 'none' })
            reject(new Error(data.msg || '上传失败'))
            return
          }
          resolve(data.data as T)
        } catch (err) {
          reject(err)
        }
      },
      fail(err) {
        reject(err)
      },
    })
  })
}
