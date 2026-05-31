import request from './request'

// 分销配置
export const getDistributionConfig = () => request.get('/admin/api/distribution/config')
export const updateDistributionConfig = (data: any) => request.put('/admin/api/distribution/config', data)

// 分销商管理
export const getDistributors = (params?: any) => request.get('/admin/api/distribution/distributors', { params })
export const getDistributor = (id: number) => request.get(`/admin/api/distribution/distributors/${id}`)
export const updateDistributor = (id: number, data: any) => request.put(`/admin/api/distribution/distributors/${id}`, data)

// 分销订单
export const getDistributionOrders = (params?: any) => request.get('/admin/api/distribution/orders', { params })
export const getDistributionOrder = (id: number) => request.get(`/admin/api/distribution/orders/${id}`)

// 提现管理
export const getWithdrawals = (params?: any) => request.get('/admin/api/distribution/withdrawals', { params })
export const getWithdrawal = (id: number) => request.get(`/admin/api/distribution/withdrawals/${id}`)
export const approveWithdrawal = (id: number) => request.post(`/admin/api/distribution/withdrawals/${id}/approve`)
export const rejectWithdrawal = (id: number, reason: string) => request.post(`/admin/api/distribution/withdrawals/${id}/reject`, { reason })
export const completeWithdrawal = (id: number) => request.post(`/admin/api/distribution/withdrawals/${id}/complete`)
