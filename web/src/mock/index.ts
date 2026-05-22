import indexDecor from '../../../app/mock/data/index-decor.json'
import categories from '../../../app/mock/data/categories.json'
import products from '../../../app/mock/data/products.json'
import productDetail from '../../../app/mock/data/product-detail.json'
import cart from '../../../app/mock/data/cart.json'
import orders from '../../../app/mock/data/orders.json'
import userCoupons from '../../../app/mock/data/user-coupons.json'
import seckills from '../../../app/mock/data/seckills.json'

const routes: Record<string, any> = {
  'GET /api/v1/index/decor': indexDecor,
  'GET /api/v1/categories': categories,
  'GET /api/v1/products': products,
  'GET /api/v1/products/': productDetail,
  'GET /api/v1/cart': cart,
  'GET /api/v1/orders': orders,
  'GET /api/v1/user/coupons': userCoupons,
  'GET /api/v1/marketing/seckills': seckills,
  'POST /api/v1/cart/add': null,
  'POST /api/v1/orders': { order_no: 'DEMO202600001', id: 1, status: 1 },
  'POST /api/v1/auth/sms/send': { dev_code: '123456' },
  'POST /api/v1/auth/sms/login': { token: 'demo_token_mock' },
}

export function matchMock(method: string, url: string): { matched: boolean; data?: any } {
  const key = `${method.toUpperCase()} ${url}`
  if (key in routes) return { matched: true, data: routes[key] }
  for (const pattern of Object.keys(routes)) {
    if (key.startsWith(pattern) && pattern.endsWith('/')) {
      return { matched: true, data: routes[pattern] }
    }
  }
  if (['POST', 'PUT', 'DELETE'].includes(method.toUpperCase())) {
    return { matched: true, data: null }
  }
  return { matched: false }
}
