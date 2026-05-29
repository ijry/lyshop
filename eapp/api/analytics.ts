import { get } from '@/utils/request'

export const getSalesAnalytics = (days?: number) => get<any>('/analytics/sales', { days: days || 7 })
export const getProductAnalytics = (days?: number) => get<any>('/analytics/products', { days: days || 7 })
export const getCustomerAnalytics = (days?: number) => get<any>('/analytics/customers', { days: days || 7 })
export const getTrafficAnalytics = (days?: number) => get<any>('/analytics/traffic', { days: days || 7 })
export const getConversionAnalytics = (days?: number) => get<any>('/analytics/conversion', { days: days || 7 })
