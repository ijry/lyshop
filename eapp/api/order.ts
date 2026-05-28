import { get, post, put } from '@/utils/request'

export const getDeliveryMode = () => get<any>('/delivery/mode')
export const getOrders = (params?: any) => get<any>('/orders', params)
export const getOrderDetail = (id: number | string) => get<any>(`/orders/${id}`)
export const shipOrder = (id: number | string, payload: any) => put<any>(`/orders/${id}/ship`, payload)
export const syncShipment = (orderID: number | string, shipmentID: number | string) =>
  post<any>(`/orders/${orderID}/shipments/${shipmentID}/sync`)
export const getShipmentTracks = (orderID: number | string, shipmentID: number | string) =>
  get<any>(`/orders/${orderID}/shipments/${shipmentID}/tracks`)

export const getAfterSales = (params?: any) => get<any>('/after-sales', params)
export const getAfterSaleDetail = (id: number | string) => get<any>(`/after-sales/${id}`)
export const auditAfterSale = (id: number | string, payload: any) => post<any>(`/after-sales/${id}/audit`, payload)
export const receiveAfterSale = (id: number | string) => post<any>(`/after-sales/${id}/receive`)
export const refundAfterSale = (id: number | string, payload: any) => post<any>(`/after-sales/${id}/refund`, payload)
export const completeAfterSale = (id: number | string) => post<any>(`/after-sales/${id}/complete`)
export const closeAfterSale = (id: number | string, payload: any) => post<any>(`/after-sales/${id}/close`, payload)

export const repriceOrder = (id: number | string, payload: { items: Array<{ item_id: number; price: number }>; remark?: string }) =>
  post<any>(`/orders/${id}/repricing`, payload)
export const addOrderNote = (id: number | string, payload: { content: string; visible_to?: string }) =>
  post<any>(`/orders/${id}/notes`, payload)
export const remindPay = (id: number | string, payload: { channel: 'sms' | 'wx' }) =>
  post<any>(`/orders/${id}/remind-pay`, payload)
export const getPrintTemplate = (id: number | string) => get<{ template: string }>(`/orders/${id}/print-template`)
export const getOrderTimeline = (id: number | string) => get<any[]>(`/orders/${id}/timeline`)
export const batchShipOrders = (rows: Array<{ order_id: number; company: string; tracking_no: string }>) =>
  post<{ success_ids: number[]; fail: Array<{ id: number; reason: string }> }>('/orders/batch/ship', rows)
export const batchNoteOrders = (payload: { ids: number[]; content: string }) =>
  post<any>('/orders/batch/notes', payload)
export const batchRepriceOrders = (payload: { ids: number[]; adjustment: { type: 'percent' | 'amount'; value: number } }) =>
  post<any>('/orders/batch/repricing', payload)
export const batchCloseOrders = (payload: { ids: number[]; reason: string }) =>
  post<any>('/orders/batch/close', payload)
