import { get, post } from '@/utils/request'

export const getMessages = (params?: any) => get<any>('/messages', params)
export const sendMessage = (payload: any) => post<any>('/messages/send', payload)
export const getImSessions = (params?: any) => get<any>('/im/sessions', params)
