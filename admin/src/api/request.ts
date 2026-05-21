import axios from 'axios'

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

export default request
