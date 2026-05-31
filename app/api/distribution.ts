import { request } from '@/utils/request'

// 分销商信息
export const getDistributorInfo = () => request({ url: '/api/v1/distribution/info', method: 'GET' })
export const applyDistributor = (data: any) => request({ url: '/api/v1/distribution/apply', method: 'POST', data })

// 我的团队
export const getMyTeam = () => request({ url: '/api/v1/distribution/team', method: 'GET' })

// 分销订单
export const getMyOrders = (params?: any) => request({ url: '/api/v1/distribution/orders', method: 'GET', params })

// 提现
export const createWithdrawal = (data: any) => request({ url: '/api/v1/distribution/withdrawals', method: 'POST', data })
export const getMyWithdrawals = (params?: any) => request({ url: '/api/v1/distribution/withdrawals', method: 'GET', params })
