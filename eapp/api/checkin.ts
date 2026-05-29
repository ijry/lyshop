import { get, put } from '@/utils/request'

export const getCheckinRules = () => get<any>('/checkin/rules')
export const saveCheckinRules = (rules: Array<{ day: number; points: number }>) => put<any>('/checkin/rules', rules)
export const getCheckinLogs = (params?: any) => get<any>('/checkin/logs', params)
