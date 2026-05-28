import request from './request'

export type WmsDocType = 'inbound' | 'outbound'
export type WmsDocStatus = 'draft' | 'completed' | 'canceled'

export interface WmsWarehouse {
  id: number
  name: string
  code?: string
  address?: string
  contact?: string
  phone?: string
  status: number
  created_at?: string
  updated_at?: string
}

export interface WmsWarehousePayload {
  name: string
  code?: string
  address?: string
  contact?: string
  phone?: string
  status?: number
}

export interface WmsDocItem {
  id?: number
  sku_id: number
  sku_name: string
  qty: number
  unit_cost?: number
  note?: string
}

export interface WmsDoc {
  id: number
  doc_no: string
  type: WmsDocType
  status: WmsDocStatus
  warehouse_id: number
  warehouse_name?: string
  remark?: string
  items: WmsDocItem[]
  total_qty: number
  created_at?: string
  updated_at?: string
}

export interface WmsPageQuery {
  page?: number
  size?: number
}

export interface WmsDocQuery extends WmsPageQuery {
  doc_type?: WmsDocType
  status?: WmsDocStatus
  warehouse_id?: number
  doc_no?: string
}

export interface WmsStockLedgerQuery extends WmsPageQuery {
  warehouse_id?: number
  sku_id?: number
  keyword?: string
}

export interface WmsMovementQuery extends WmsPageQuery {
  warehouse_id?: number
  sku_id?: number
  doc_no?: string
}

export interface WmsPageResult<T> {
  list: T[]
  total: number
  page: number
  size: number
}

export interface WmsStockLedgerRow {
  id: number
  warehouse_id: number
  warehouse_name: string
  sku_id: number
  sku_name: string
  qty: number
  safe_qty: number
  updated_at?: string
}

export interface WmsMovementRow {
  id: number
  doc_id: number
  doc_no: string
  type: WmsDocType
  warehouse_id: number
  warehouse_name: string
  sku_id: number
  sku_name: string
  qty: number
  before_qty: number
  after_qty: number
  created_at?: string
}

type WmsDocServer = Omit<WmsDoc, 'type' | 'items'> & {
  type?: WmsDocType
  doc_type?: WmsDocType
}

function normalizeDocType(raw: any): WmsDocType {
  return String(raw?.doc_type || raw?.type || 'inbound') === 'outbound' ? 'outbound' : 'inbound'
}

function normalizeDoc(raw: any, items?: WmsDocItem[]): WmsDoc {
  const list = Array.isArray(items) ? items : Array.isArray(raw?.items) ? raw.items : []
  const derivedTotal = list.reduce((sum: number, row: WmsDocItem) => sum + Number(row?.qty || 0), 0)
  return {
    ...raw,
    type: normalizeDocType(raw),
    status: String(raw?.status || 'draft') as WmsDocStatus,
    items: list,
    total_qty: Number(raw?.total_qty || derivedTotal),
  }
}

export const listWarehouses = (params?: { keyword?: string; status?: number }) =>
  request.get<never, WmsPageResult<WmsWarehouse>>('/wms/warehouses', { params })

export const createWarehouse = (payload: WmsWarehousePayload) =>
  request.post<never, WmsWarehouse | null>('/wms/warehouses', payload)

export const updateWarehouse = (id: number, payload: WmsWarehousePayload) =>
  request.put<never, null>(`/wms/warehouses/${id}`, payload)

export const listDocs = async (params?: WmsDocQuery) => {
  const data = await request.get<never, WmsPageResult<WmsDocServer>>('/wms/docs', { params })
  return {
    ...data,
    list: Array.isArray(data?.list) ? data.list.map((item) => normalizeDoc(item)) : [],
  } as WmsPageResult<WmsDoc>
}

export const getDocDetail = async (id: number) => {
  const data = await request.get<never, { doc?: WmsDocServer; items?: WmsDocItem[] } | null>(`/wms/docs/${id}`)
  if (!data?.doc) return null
  return normalizeDoc(data.doc, data.items)
}

export const createDoc = async (payload: Partial<WmsDoc>) => {
  const data = await request.post<never, WmsDocServer | null>('/wms/docs', {
    ...payload,
    doc_type: payload.type,
  })
  if (!data) return null
  return normalizeDoc(data)
}

export const saveDoc = (id: number, payload: Partial<WmsDoc>) =>
  request.put<never, null>(`/wms/docs/${id}`, {
    ...payload,
    doc_type: payload.type,
  })

export const completeDoc = (id: number) =>
  request.post<never, null>(`/wms/docs/${id}/complete`)

export const cancelDoc = (id: number) =>
  request.post<never, null>(`/wms/docs/${id}/cancel`)

export const listStockLedger = (params?: WmsStockLedgerQuery) =>
  request.get<never, WmsPageResult<WmsStockLedgerRow>>('/wms/stocks', { params })

export const updateSafetyStock = (id: number, safeQty: number) =>
  request.put<never, null>(`/wms/stocks/${id}/safety`, { safe_qty: safeQty })

export const listMovements = (params?: WmsMovementQuery) =>
  request.get<never, WmsPageResult<WmsMovementRow>>('/wms/movements', { params })
