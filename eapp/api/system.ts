import { get, put } from '@/utils/request'

export const getSiteSettings = () => get<any>('/site-settings')
export const updateSiteSettings = (payload: any) => put<any>('/site-settings', payload)

export const getAdmins = () => get<any>('/admins')
export const getRoles = () => get<any>('/roles')
