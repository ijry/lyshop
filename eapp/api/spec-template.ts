import { get, post, put, del } from '@/utils/request'

export const getSpecTemplates = (params?: any) => get<any>('/spec-templates', params)
export const getSpecTemplate = (id: number | string) => get<any>(`/spec-templates/${id}`)
export const createSpecTemplate = (payload: any) => post<any>('/spec-templates', payload)
export const updateSpecTemplate = (id: number | string, payload: any) => put<any>(`/spec-templates/${id}`, payload)
export const deleteSpecTemplate = (id: number | string) => del<any>(`/spec-templates/${id}`)
