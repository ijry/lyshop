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

// URL pattern → mock data (supports both exact and prefix match)
const routes: Record<string, any> = {
  'GET /api/v1/index/decor': indexDecor,
  'GET /api/v1/categories': categories,
  'GET /api/v1/products': products,
  'GET /api/v1/products/recommend': recommend,
  'GET /api/v1/products/': productDetail, // prefix match for /products/:id
  'GET /api/v1/cart': cart,
  'GET /api/v1/orders': orders,
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
  'GET /api/v1/im/session': { id: 1, user_id: 1, status: 2 },
  'GET /api/v1/im/messages': { list: [], total: 0 },
}

/**
 * Try to match a mock route. Returns { matched: true, data } or { matched: false }.
 */
export function matchMock(method: string, url: string): { matched: boolean; data?: any } {
  const key = `${method.toUpperCase()} ${url}`

  // Exact match
  if (key in routes) {
    return { matched: true, data: routes[key] }
  }

  // Prefix match (for routes with path params like /products/:id)
  for (const pattern of Object.keys(routes)) {
    if (key.startsWith(pattern) && pattern.endsWith('/')) {
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
