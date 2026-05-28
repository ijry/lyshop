import { get } from '@/utils/request'
import { post, put } from '@/utils/request'

export const getCoupons = (params?: any) => get<any>('/marketing/coupons', params)
export const createCoupon = (payload: any) => post<any>('/marketing/coupons', payload)

export const getSeckillActivities = (params?: any) => get<any>('/marketing/seckill/activities', params)
export const getGroupBuyActivities = (params?: any) => get<any>('/marketing/group-buy/activities', params)
export const getBargainActivities = (params?: any) => get<any>('/marketing/bargain/activities', params)

export type MarketingKind = 'seckill' | 'group-buy' | 'bargain'

function activityPrefix(kind: MarketingKind) {
  if (kind === 'seckill') return '/marketing/seckill'
  if (kind === 'group-buy') return '/marketing/group-buy'
  return '/marketing/bargain'
}

export const getActivities = (kind: MarketingKind, params?: any) =>
  get<any>(`${activityPrefix(kind)}/activities`, params)

export const createActivity = (kind: MarketingKind, payload: any) =>
  post<any>(`${activityPrefix(kind)}/activities`, payload)

export const updateActivity = (kind: MarketingKind, activityID: number | string, payload: any) =>
  put<any>(`${activityPrefix(kind)}/activities/${activityID}`, payload)

export const getActivityProducts = (kind: MarketingKind, params?: any) =>
  get<any>(`${activityPrefix(kind)}/products`, params)

export const upsertActivityProducts = (kind: MarketingKind, activityID: number | string, rows: any[]) =>
  put<any>(`${activityPrefix(kind)}/activities/${activityID}/products`, rows)
