import request from './request'

export interface AdminMenuItem {
  title: string
  path: string
  icon?: string
  sort?: number
  children?: AdminMenuItem[]
}

export interface AdminMenuGroup {
  key: string
  title: string
  icon?: string
  sort?: number
  menus: AdminMenuItem[]
}

export interface AdminMenuGroupedResponse {
  dashboard: {
    title: string
    path: string
  }
  groups: AdminMenuGroup[]
}

export type AdminMenuLegacyResponse = AdminMenuItem[]
export type AdminMenuResponse = AdminMenuGroupedResponse | AdminMenuLegacyResponse

export const login = (username: string, password: string) =>
  request.post<never, { token: string }>('/auth/login', { username, password })

export const getMenus = () =>
  request.get<never, AdminMenuResponse>('/menus')
