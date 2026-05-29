import { del, get, post, put } from '@/utils/request'
export const getImSessions = (params?: any) => get<any>('/im/sessions', params)
export const getImMessages = (sessionId: number) => get<any>(`/im/sessions/${sessionId}/messages`)
export const sendImMessage = (sessionId: number, content: string) => post<any>(`/im/sessions/${sessionId}/messages`, { content, sender_type: 'staff' })
export const getAutoReplies = () => get<any>('/im/auto-replies')
export const createAutoReply = (payload: any) => post<any>('/im/auto-replies', payload)
export const updateAutoReply = (id: number, payload: any) => put<any>(`/im/auto-replies/${id}`, payload)
export const deleteAutoReply = (id: number) => del<any>(`/im/auto-replies/${id}`)
