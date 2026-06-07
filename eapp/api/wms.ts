import { del, get, post, put } from '@/utils/request'
// Warehouses
export const getWarehouses = (params?: any) => get<any>('/wms/warehouses', params)
export const createWarehouse = (payload: any) => post<any>('/wms/warehouses', payload)
export const updateWarehouse = (id: number, payload: any) => put<any>(`/wms/warehouses/${id}`, payload)
// Stocks
export const getStocks = (params?: any) => get<any>('/wms/stocks', params)
export const getStocksBySkuIds = (skuIds: number[]) => get<any>('/wms/stocks/by-skus', { sku_ids: skuIds.join(',') })
export const updateSafeQty = (id: number, safe_qty: number) => put<any>(`/wms/stocks/${id}/safe-qty`, { safe_qty })
// Docs
export const getDocs = (params?: any) => get<any>('/wms/docs', params)
export const getDocDetail = (id: number) => get<any>(`/wms/docs/${id}`)
export const createDoc = (payload: any) => post<any>('/wms/docs', payload)
export const saveDoc = (id: number, payload: any) => put<any>(`/wms/docs/${id}`, payload)
export const completeDoc = (id: number) => post<any>(`/wms/docs/${id}/complete`)
export const cancelDoc = (id: number) => post<any>(`/wms/docs/${id}/cancel`)
// Movements
export const getMovements = (params?: any) => get<any>('/wms/movements', params)
