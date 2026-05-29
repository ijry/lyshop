import { get, post } from '@/utils/request'

export const getReviews = (params?: any) => get<any>('/reviews', params)
export const getReviewDetail = (id: number | string) => get<any>(`/reviews/${id}`)
export const replyReview = (id: number | string, content: string) => post<any>(`/reviews/${id}/reply`, { content })
