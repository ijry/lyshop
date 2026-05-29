import { get, put } from '@/utils/request'

export const getDistributors = (params?: any) => get<any>('/distribution/distributors', params)
export const getDistributorDetail = (id: number) => get<any>(`/distribution/distributors/${id}`)
export const updateDistributorStatus = (id: number, status: number) => put<any>(`/distribution/distributors/${id}/status`, { status })
export const getCommissions = (params?: any) => get<any>('/distribution/commissions', params)
export const settleCommission = (id: number) => put<any>(`/distribution/commissions/${id}/settle`)
export const returnCommission = (id: number) => put<any>(`/distribution/commissions/${id}/return`)
export const getDistributionConfig = () => get<any>('/distribution/config')
export const updateDistributionConfig = (payload: any) => put<any>('/distribution/config', payload)
