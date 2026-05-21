const BASE_URL = 'http://localhost:8080'

function getToken(): string {
  return uni.getStorageSync('user_token') || ''
}

export function request<T = any>(options: UniNamespace.RequestOptions): Promise<T> {
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

export const get = <T>(url: string, data?: any) =>
  request<T>({ url, method: 'GET', data })

export const post = <T>(url: string, data?: any) =>
  request<T>({ url, method: 'POST', data })
