import { get, post } from '@/utils/request'

export const getAfterSales = (params?: any) => get<any>('/after-sales', params)
export const getAfterSaleDetail = (id: number | string) => get<any>(`/after-sales/${id}`)
export const auditAfterSale = (id: number | string, payload: any) => post<any>(`/after-sales/${id}/audit`, payload)
export const receiveAfterSale = (id: number | string) => post<any>(`/after-sales/${id}/receive`)
export const refundAfterSale = (id: number | string, payload: any) => post<any>(`/after-sales/${id}/refund`, payload)
export const completeAfterSale = (id: number | string) => post<any>(`/after-sales/${id}/complete`)
export const closeAfterSale = (id: number | string, payload: any) => post<any>(`/after-sales/${id}/close`, payload)
export const addAfterSaleMessage = (id: number | string, payload: { from: 'merchant'|'user'; content: string; images?: string[] }) =>
  post<any>(`/after-sales/${id}/messages`, payload)
export const addAfterSaleEvidence = (id: number | string, payload: { images: string[]; remark?: string }) =>
  post<any>(`/after-sales/${id}/evidences`, payload)
