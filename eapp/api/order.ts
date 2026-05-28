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
