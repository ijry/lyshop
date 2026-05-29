import { get, post, put, del } from '@/utils/request'

// Products
export const getPointsProducts = (params?: any) => get<any>('/points/products', params)
export const createPointsProduct = (payload: any) => post<any>('/points/products', payload)
export const updatePointsProduct = (id: number, payload: any) => put<any>(`/points/products/${id}`, payload)
export const deletePointsProduct = (id: number) => del<any>(`/points/products/${id}`)

// Exchanges
export const getPointsExchanges = (params?: any) => get<any>('/points/exchanges', params)
export const shipExchange = (id: number) => put<any>(`/points/exchanges/${id}/ship`)
export const completeExchange = (id: number) => put<any>(`/points/exchanges/${id}/complete`)

// Summary
export const getPointsSummary = () => get<any>('/points/summary')
