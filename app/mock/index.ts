/**
 * Mock data registry.
 * Key = "METHOD URL" (e.g. "GET /api/v1/products"), value = mock response data.
 * Data is the `data` field inside { code: 0, msg: "success", data: ... }.
 */

import indexDecor from './data/index-decor.json'
import categories from './data/categories.json'
import products from './data/products.json'
import productDetail from './data/product-detail.json'
import cart from './data/cart.json'
import orders from './data/orders.json'
import userCoupons from './data/user-coupons.json'
import seckills from './data/seckills.json'
import groupBuy from './data/group-buy.json'
import bargain from './data/bargain.json'
import recommend from './data/recommend.json'
import userProfile from './data/user-profile.json'
import addresses from './data/addresses.json'

function parseQuery(url: string) {
  const queryIndex = url.indexOf('?')
  if (queryIndex < 0) return {}
  return Object.fromEntries(new URLSearchParams(url.slice(queryIndex + 1)).entries())
}

// URL pattern → mock data (supports both exact and prefix match)
const routes: Record<string, any> = {
  'GET /api/v1/index/decor': indexDecor,
  'GET /api/v1/categories': categories,
  'GET /api/v1/products': products,
  'GET /api/v1/products/recommend': recommend,
  'GET /api/v1/products/': productDetail, // prefix match for /products/:id
  'GET /api/v1/cart': cart,
  'GET /api/v1/orders': orders,
  'GET /api/v1/orders/': (orders as any)?.list?.[0] || null,
  'GET /api/v1/user/coupons': userCoupons,
  'GET /api/v1/user/profile': userProfile,
  'GET /api/v1/user/points/logs': { list: [
    { id: 1, type: 1, points: 100, remark: '购买商品奖励', created_at: '2026-05-20T10:00:00Z' },
    { id: 2, type: 1, points: 580, remark: '首单奖励', created_at: '2026-05-15T08:30:00Z' },
    { id: 3, type: 2, points: -200, remark: '积分兑换优惠券', created_at: '2026-05-10T14:20:00Z' },
    { id: 4, type: 1, points: 2200, remark: '活动签到奖励', created_at: '2026-04-01T09:00:00Z' },
  ], total: 4 },
  'GET /api/v1/addresses': addresses,
  'GET /api/v1/marketing/seckills': seckills,
  'GET /api/v1/marketing/group-buy': groupBuy,
  'GET /api/v1/marketing/bargain': bargain,
  'POST /api/v1/marketing/group-buy/join': { group_order_id: 1, status: 1 },
  'POST /api/v1/marketing/bargain/help': { cut_amount: 58.50, current_price: 1200 },
  'POST /api/v1/cart/add': null,
  'POST /api/v1/orders': { order_no: 'DEMO202600001', id: 1, status: 1 },
  'POST /api/v1/auth/sms/send': { dev_code: '123456' },
  'POST /api/v1/auth/sms/login': { token: 'demo_token_mock' },
  'POST /api/v1/addresses': { id: 3 },
'GET /api/v1/checkin/status': { checked_today: false, consecutive_days: 3, month_dates: ['2026-05-17','2026-05-20','2026-05-21'], month_count: 3, month_points: 40 },  'GET /api/v1/checkin/rules': [{ day: 0, points: 10 }, { day: 3, points: 20 }, { day: 7, points: 50 }],  'POST /api/v1/checkin': { points: 20, consecutive_days: 4 },
  'GET /api/v1/im/session': { id: 1, user_id: 1, status: 2 },
  'GET /api/v1/im/messages': { list: [], total: 0 },
'GET /api/v1/messages/unread': { system: 2, order: 1, marketing: 3, im: 0 },  'GET /api/v1/messages': { list: [    { id: 1, group: 'system', title: '系统升级通知', content: 'LYShop 已升级至 2.0 版本', is_read: 0, created_at: '2026-05-22T10:00:00Z' },    { id: 2, group: 'order', title: '订单发货通知', content: '您的订单已发货', is_read: 0, created_at: '2026-05-21T14:00:00Z' },    { id: 3, group: 'marketing', title: '618大促即将开始', content: '全场满300减50', is_read: 0, created_at: '2026-05-20T08:00:00Z' },    { id: 4, group: 'marketing', title: '优惠券到账', content: '满100减20优惠券已到账', is_read: 1, created_at: '2026-05-18T12:00:00Z' },    { id: 5, group: 'system', title: '欢迎注册', content: '新人专享优惠等你来领', is_read: 1, created_at: '2026-05-15T10:00:00Z' },  ], total: 5 },
}

/**
 * Try to match a mock route. Returns { matched: true, data } or { matched: false }.
 */
export function matchMock(method: string, url: string, params?: Record<string, any>): { matched: boolean; data?: any } {
  const upperMethod = method.toUpperCase()
  const [path] = url.split('?')
  const key = `${upperMethod} ${path}`
  const query = { ...parseQuery(url), ...(params || {}) }

  // Exact match
  if (key in routes) {
    const data = routes[key]
    if (upperMethod === 'GET' && path === '/api/v1/products') {
      const sourceList = Array.isArray(data?.list) ? data.list : []
      const keyword = String(query.keyword || '').trim().toLowerCase()
      const categoryID = Number(query.category_id || 0)
      const page = Number(query.page || 1)
      const size = Number(query.size || 20)
      let list = sourceList.slice()
      if (keyword) {
        list = list.filter((item: any) => String(item.title || '').toLowerCase().includes(keyword))
      }
      if (categoryID > 0) {
        list = list.filter((item: any) => Number(item.category_id) === categoryID)
      }
      const offset = Math.max(page - 1, 0) * Math.max(size, 1)
      const pageList = list.slice(offset, offset + size)
      return { matched: true, data: { ...data, list: pageList, total: list.length, page, size } }
    }
    if (upperMethod === 'GET' && path === '/api/v1/orders' && query.status) {
      const status = Number(query.status)
      const list = Array.isArray(data?.list)
        ? data.list.filter((item: any) => Number(item.status) === status)
        : []
      return { matched: true, data: { ...data, list, total: list.length } }
    }
    return { matched: true, data }
  }

  // Prefix match (for routes with path params like /products/:id)
  for (const pattern of Object.keys(routes)) {
    if (key.startsWith(pattern) && pattern.endsWith('/')) {
      if (pattern === 'GET /api/v1/orders/') {
        const id = Number(path.split('/').pop() || 0)
        const source = Array.isArray((orders as any)?.list) ? (orders as any).list : []
        const detail = source.find((item: any) => Number(item.id) === id) || null
        return { matched: true, data: detail }
      }
      return { matched: true, data: routes[pattern] }
    }
  }

  // Default: return empty success for unmatched POST/PUT/DELETE
  if (['POST', 'PUT', 'DELETE'].includes(method.toUpperCase())) {
    return { matched: true, data: null }
  }

  return { matched: false }
}

export const isMockEnabled = true
