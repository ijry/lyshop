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
  doc_id?: number
  sku_id: number
  qty: number
  remark?: string
  // UI helper fields, not part of backend contract
  sku_name?: string
  unit_cost?: number
}

export interface WmsDoc {
  id: number
  doc_no: string
  type: WmsDocType
  status: WmsDocStatus
  warehouse_id: number
  remark?: string
  items: WmsDocItem[]
  total_qty?: number
  // UI helper field, not part of backend contract
  warehouse_name?: string
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
  biz_type?: WmsDocType
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
  biz_type: WmsDocType
  warehouse_id: number
  warehouse_name?: string
  sku_id: number
  sku_name?: string
  change_qty: number
  before_qty: number
  after_qty: number
  occurred_at?: string
}

interface WmsDocServer {
  id: number
  doc_no: string
  doc_type?: WmsDocType
  status?: WmsDocStatus | string
  warehouse_id: number
  remark?: string
  total_qty?: number
  created_at?: string
  updated_at?: string
}

function normalizeDocType(raw: any): WmsDocType {
  return String(raw?.doc_type || 'inbound') === 'outbound' ? 'outbound' : 'inbound'
}

function normalizeDocItems(items: any[] | undefined): WmsDocItem[] {
  if (!Array.isArray(items)) return []
  return items.map((row) => ({
    id: Number(row?.id || 0) || undefined,
    doc_id: Number(row?.doc_id || 0) || undefined,
    sku_id: Number(row?.sku_id || 0),
    qty: Number(row?.qty || 0),
    remark: String(row?.remark || ''),
    sku_name: String(row?.sku_name || ''),
    unit_cost: row?.unit_cost === undefined ? undefined : Number(row.unit_cost || 0),
  }))
}

function normalizeDoc(raw: any, items?: WmsDocItem[]): WmsDoc {
  const list = normalizeDocItems(Array.isArray(items) ? items : raw?.items)
  const rawTotalQty = raw?.total_qty
  return {
    id: Number(raw?.id || 0),
    doc_no: String(raw?.doc_no || ''),
    type: normalizeDocType(raw),
    status: String(raw?.status || 'draft') as WmsDocStatus,
    warehouse_id: Number(raw?.warehouse_id || 0),
    remark: String(raw?.remark || ''),
    items: list,
    total_qty: typeof rawTotalQty === 'number' ? rawTotalQty : undefined,
    warehouse_name: String(raw?.warehouse_name || ''),
    created_at: raw?.created_at,
    updated_at: raw?.updated_at,
  }
}

type WmsDocSubmitPayload = {
  type: WmsDocType
  warehouse_id: number
  remark?: string
  items: Array<{
    id?: number
    sku_id: number
    qty: number
    remark?: string
  }>
}

function toDocSubmitPayload(payload: Partial<WmsDoc>): WmsDocSubmitPayload {
  const type = String(payload?.type || 'inbound') === 'outbound' ? 'outbound' : 'inbound'
  const items = Array.isArray(payload?.items)
    ? payload.items.map((row) => ({
      id: Number(row?.id || 0) || undefined,
      sku_id: Number(row?.sku_id || 0),
      qty: Number(row?.qty || 0),
      remark: String(row?.remark || ''),
    }))
    : []
  return {
    type,
    warehouse_id: Number(payload?.warehouse_id || 0),
    remark: String(payload?.remark || ''),
    items,
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
  const submitPayload = toDocSubmitPayload(payload)
  const data = await request.post<never, WmsDocServer | null>('/wms/docs', {
    doc_type: submitPayload.type,
    warehouse_id: submitPayload.warehouse_id,
    remark: submitPayload.remark,
    items: submitPayload.items,
  })
  if (!data) return null
  return normalizeDoc(data)
}

export const saveDoc = (id: number, payload: Partial<WmsDoc>) => {
  const submitPayload = toDocSubmitPayload(payload)
  return request.put<never, null>(`/wms/docs/${id}`, {
    doc_type: submitPayload.type,
    warehouse_id: submitPayload.warehouse_id,
    remark: submitPayload.remark,
    items: submitPayload.items,
  })
}

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
