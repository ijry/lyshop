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

// URL pattern → mock data (supports both exact and prefix match)
const routes: Record<string, any> = {
  'GET /api/v1/index/decor': indexDecor,
  'GET /api/v1/categories': categories,
  'GET /api/v1/products': products,
  'GET /api/v1/products/': productDetail, // prefix match for /products/:id
  'GET /api/v1/cart': cart,
  'GET /api/v1/orders': orders,
  'GET /api/v1/user/coupons': userCoupons,
  'GET /api/v1/marketing/seckills': seckills,
  'POST /api/v1/cart/add': null,
  'POST /api/v1/orders': { order_no: 'DEMO202600001', id: 1, status: 1 },
  'POST /api/v1/auth/sms/send': { dev_code: '123456' },
  'POST /api/v1/auth/sms/login': { token: 'demo_token_mock' },
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
