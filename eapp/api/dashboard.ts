import { get } from '@/utils/request'

export const getDashboard = () => get<any>('/dashboard')
