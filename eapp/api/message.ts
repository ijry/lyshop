import { get, post, put } from '@/utils/request'

export const getMessages = (params?: any) => get<any>('/messages', params)
export const sendMessage = (payload: any) => post<any>('/messages/send', payload)
export const markMessageRead = (id: number) => put<any>(`/messages/${id}/read`)
export const getUnreadCounts = () => get<any>('/messages/unread-counts')
