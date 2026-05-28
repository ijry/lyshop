import { del, get, post, put } from '@/utils/request'

export const getProducts = (params?: any) => get<any>('/products', params)
export const getProductDetail = (id: number | string) => get<any>(`/products/${id}`)

export type ProductUpsertPayload = {
  product: Record<string, any>
  skus?: Array<Record<string, any>>
  images?: Array<Record<string, any>>
}

export const createProduct = (payload: ProductUpsertPayload) => post<any>('/products', payload)
export const updateProduct = (id: number | string, payload: ProductUpsertPayload) => put<any>(`/products/${id}`, payload)
export const deleteProduct = (id: number | string) => del<any>(`/products/${id}`)
