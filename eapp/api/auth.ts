import { get, post } from '@/utils/request'

export type LoginResp = { token: string }

export const login = (username: string, password: string) =>
  post<LoginResp>('/auth/login', { username, password })

export const getMenus = () => get<any>('/menus')
