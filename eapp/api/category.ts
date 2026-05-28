import { del, get, post, put } from '@/utils/request'

export const getCategoriesTree = () => get<any>('/categories/tree')
export const createCategory = (payload: { name: string; parent_id?: number; sort?: number }) => post<any>('/categories', payload)
export const updateCategory = (id: number, payload: { name?: string; sort?: number }) => put<any>(`/categories/${id}`, payload)
export const deleteCategory = (id: number) => del<any>(`/categories/${id}`)
