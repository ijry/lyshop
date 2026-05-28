import { get, post, put, del } from '@/utils/request'

// Plans
export const getVipPlans = (params?: any) => get<any>('/vip/plans', params)
export const getVipPlan = (id: number | string) => get<any>(`/vip/plans/${id}`)
export const createVipPlan = (payload: any) => post<any>('/vip/plans', payload)
export const updateVipPlan = (id: number | string, payload: any) => put<any>(`/vip/plans/${id}`, payload)
export const deleteVipPlan = (id: number | string) => del<any>(`/vip/plans/${id}`)

// Levels
export const getVipLevels = (params?: any) => get<any>('/vip/levels', params)
export const createVipLevel = (payload: any) => post<any>('/vip/levels', payload)
export const updateVipLevel = (id: number | string, payload: any) => put<any>(`/vip/levels/${id}`, payload)
export const deleteVipLevel = (id: number | string) => del<any>(`/vip/levels/${id}`)

// Coupon rules
export const getVipCouponRules = (params?: any) => get<any>('/vip/coupon-rules', params)
export const createVipCouponRule = (payload: any) => post<any>('/vip/coupon-rules', payload)
export const updateVipCouponRule = (id: number | string, payload: any) => put<any>(`/vip/coupon-rules/${id}`, payload)
export const deleteVipCouponRule = (id: number | string) => del<any>(`/vip/coupon-rules/${id}`)

// SKU prices
export const getVipSkuPrices = (params?: any) => get<any>('/vip/sku-prices', params)
export const createVipSkuPrice = (payload: any) => post<any>('/vip/sku-prices', payload)
export const updateVipSkuPrice = (id: number | string, payload: any) => put<any>(`/vip/sku-prices/${id}`, payload)
export const deleteVipSkuPrice = (id: number | string) => del<any>(`/vip/sku-prices/${id}`)
