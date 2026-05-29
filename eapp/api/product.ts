import { del, get, post, put } from '@/utils/request'

export const getProducts = (params?: any) => get<any>('/products', params)
export const getProductDetail = (id: number | string) => get<any>(`/products/${id}`)

export type ProductUpsertPayload = {
  product: Record<string, any>
  skus?: Array<Record<string, any>>
  images?: Array<Record<string, any>>
  spec_schema?: Array<{ name: string; values: string[] }>
  sku_overrides?: Array<{ sku_key: string; sku_code?: string; price?: number; stock?: number }>
  sku_generation_mode?: 'auto' | 'manual'
}

export const createProduct = (payload: ProductUpsertPayload) => post<any>('/products', payload)
export const updateProduct = (id: number | string, payload: ProductUpsertPayload) => put<any>(`/products/${id}`, payload)
export const deleteProduct = (id: number | string) => del<any>(`/products/${id}`)

export const batchUpdateProductStatus = (payload: { ids: number[]; status: 0 | 1 }) => put<any>('/products/batch/status', payload)
export const batchUpdateProductCategory = (payload: { ids: number[]; category_id: number }) => put<any>('/products/batch/category', payload)
export const batchUpdateProductPrice = (payload: { ids: number[]; adjustment: { type: 'set'|'percent'|'amount'; value: number; scope?: 'all'|'main_sku' } }) => put<any>('/products/batch/price', payload)
