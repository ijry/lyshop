import indexDecor from '../../../app/mock/data/index-decor.json'
import categories from '../../../app/mock/data/categories.json'
import products from '../../../app/mock/data/products.json'
import productDetail from '../../../app/mock/data/product-detail.json'
import cart from '../../../app/mock/data/cart.json'
import orders from '../../../app/mock/data/orders.json'
import userCoupons from '../../../app/mock/data/user-coupons.json'
import seckills from '../../../app/mock/data/seckills.json'
import groupBuy from '../../../app/mock/data/group-buy.json'
import bargain from '../../../app/mock/data/bargain.json'
import recommend from '../../../app/mock/data/recommend.json'
import userProfile from '../../../app/mock/data/user-profile.json'
import addresses from '../../../app/mock/data/addresses.json'

const orderListSource = (orders as any)?.list || []

const routes: Record<string, any> = {
  'GET /api/v1/index/decor': indexDecor,
  'GET /api/v1/categories': categories,
  'GET /api/v1/products': products,
  'GET /api/v1/products/recommend': recommend,
  'GET /api/v1/products/': productDetail,
  'GET /api/v1/cart': cart,
  'GET /api/v1/orders': orders,
  'GET /api/v1/orders/': orderListSource[0] || null,
  'GET /api/v1/user/coupons': userCoupons,
  'GET /api/v1/user/profile': userProfile,
  'GET /api/v1/user/points/logs': {
    list: [
      { id: 1, type: 1, points: 100, remark: '购买商品奖励', created_at: '2026-05-20T10:00:00Z' },
      { id: 2, type: 1, points: 580, remark: '首单奖励', created_at: '2026-05-15T08:30:00Z' },
      { id: 3, type: 2, points: -200, remark: '积分兑换优惠券', created_at: '2026-05-10T14:20:00Z' },
    ],
    total: 3,
  },
  'GET /api/v1/addresses': addresses,
  'GET /api/v1/marketing/seckills': seckills,
  'GET /api/v1/marketing/group-buy': groupBuy,
  'GET /api/v1/marketing/bargain': bargain,
  'POST /api/v1/cart/add': null,
  'POST /api/v1/orders': { order_no: 'DEMO202600001', id: 1, status: 1 },
  'POST /api/v1/auth/sms/send': { dev_code: '123456' },
  'POST /api/v1/auth/sms/login': { token: 'demo_token_mock' },
  'POST /api/v1/addresses': { id: 3 },
}

export function matchMock(method: string, url: string, params?: Record<string, any>): { matched: boolean; data?: any } {
  const key = `${method.toUpperCase()} ${url}`
  const query = params || {}

  if (key === 'GET /api/v1/orders') {
    const status = Number(query.status || 0)
    const list = status > 0
      ? orderListSource.filter((item: any) => Number(item.status) === status)
      : orderListSource.slice()
    return { matched: true, data: { ...orders, list, total: list.length } }
  }

  if (key in routes) return { matched: true, data: routes[key] }

  for (const pattern of Object.keys(routes)) {
    if (key.startsWith(pattern) && pattern.endsWith('/')) {
      if (pattern === 'GET /api/v1/orders/') {
        const id = Number(url.split('/').pop() || 0)
        const detail = orderListSource.find((item: any) => Number(item.id) === id) || null
        return { matched: true, data: detail }
      }
      return { matched: true, data: routes[pattern] }
    }
  }

  if (['POST', 'PUT', 'DELETE'].includes(method.toUpperCase())) {
    return { matched: true, data: null }
  }
  return { matched: false }
}
