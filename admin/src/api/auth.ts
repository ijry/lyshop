import request from './request'

export const login = (username: string, password: string) =>
  request.post<never, { token: string }>('/auth/login', { username, password })

export const getMenus = () =>
  request.get<never, any[]>('/menus')
