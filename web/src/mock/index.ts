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

function parseQuery(url: string) {
  const queryIndex = url.indexOf('?')
  if (queryIndex < 0) return {}
  return Object.fromEntries(new URLSearchParams(url.slice(queryIndex + 1)).entries())
}

const orderListSource = Array.isArray((orders as any)?.list)
  ? JSON.parse(JSON.stringify((orders as any).list))
  : []
const addressListSource = Array.isArray(addresses)
  ? JSON.parse(JSON.stringify(addresses))
  : []

function getOrderByID(id: number) {
  return orderListSource.find((item: any) => Number(item.id) === id)
}

function buildProductReview(productID: number) {
  const list: any[] = []
  for (const order of orderListSource) {
    for (const item of order.items || []) {
      if (Number(item.product_id || 0) !== productID || !item.review) continue
      list.push({
        id: Number(item.review.id),
        order_id: Number(order.id),
        order_no: order.order_no,
        order_item_id: Number(item.id),
        product_id: Number(item.product_id || 0),
        product_score: Number(item.review.product_score || 5),
        logistics_score: Number(item.review.logistics_score || 5),
        content: String(item.review.content || ''),
        edited_times: Number(item.review.edited_times || 0),
        user_nickname: '演示用户',
        user_avatar: '',
        created_at: item.review.created_at || order.created_at,
        updated_at: item.review.updated_at || item.review.created_at || order.created_at,
        appends: Array.isArray(item.review.appends) ? item.review.appends.map((ap: any) => ({
          id: Number(ap.id || 0),
          content: String(ap.content || ''),
          created_at: ap.created_at || order.created_at,
        })) : [],
        admin_reply: item.review.admin_reply ? {
          id: Number(item.review.admin_reply.id || 0),
          content: String(item.review.admin_reply.content || ''),
          created_at: item.review.admin_reply.created_at || order.created_at,
        } : null,
      })
    }
  }
  list.sort((a, b) => Number(b.id) - Number(a.id))
  const summary = list.length
    ? {
        avg_product_score: Number((list.reduce((s, i) => s + Number(i.product_score || 0), 0) / list.length).toFixed(1)),
        avg_logistics_score: Number((list.reduce((s, i) => s + Number(i.logistics_score || 0), 0) / list.length).toFixed(1)),
        total: list.length,
      }
    : { avg_product_score: 0, avg_logistics_score: 0, total: 0 }
  return { summary, list, total: list.length, page: 1, size: 20 }
}

function buildOrderReviewMeta(orderID: number) {
  const order = getOrderByID(orderID)
  if (!order) return null
  const options = (order.items || []).map((item: any) => ({
    order_item_id: Number(item.id || 0),
    review_id: Number(item.review?.id || 0),
    has_review: !!item.review,
    product_id: Number(item.product_id || 0),
    product_title: String(item.title || ''),
    product_cover: String(item.cover || ''),
    product_score: Number(item.review?.product_score || 5),
    logistics_score: Number(item.review?.logistics_score || 5),
    content: String(item.review?.content || ''),
  }))
  const reviewed = options.filter((item: any) => item.has_review)
  return {
    order_id: Number(order.id),
    order_no: order.order_no,
    logistics_score: Number(reviewed[0]?.logistics_score || 5),
    can_create: reviewed.length < options.length,
    can_edit: reviewed.length > 0,
    can_append: reviewed.length > 0,
    options,
  }
}

function listOrders(status: number) {
  const list = status > 0
    ? orderListSource.filter((item: any) => Number(item.status) === status)
    : orderListSource.slice()
  return { list, total: list.length, page: 1, size: 20 }
}

function upsertAddress(data: Record<string, any>, id?: number) {
  const payload = {
    name: String(data.name || '').trim(),
    phone: String(data.phone || '').trim(),
    province: String(data.province || '').trim(),
    city: String(data.city || '').trim(),
    district: String(data.district || '').trim(),
    detail: String(data.detail || '').trim(),
    is_default: Number(data.is_default || 0) === 1 ? 1 : 0,
  }
  if (id) {
    const idx = addressListSource.findIndex((item: any) => Number(item.id) === id)
    if (idx < 0) return null
    if (payload.is_default === 1) addressListSource.forEach((item: any) => { item.is_default = 0 })
    addressListSource[idx] = { ...addressListSource[idx], ...payload }
    return addressListSource[idx]
  }
  const nextID = Math.max(0, ...addressListSource.map((item: any) => Number(item.id || 0))) + 1
  if (payload.is_default === 1 || addressListSource.length === 0) {
    addressListSource.forEach((item: any) => { item.is_default = 0 })
    payload.is_default = 1
  }
  const created = { id: nextID, user_id: 1, ...payload }
  addressListSource.unshift(created)
  return created
}

function removeAddress(id: number) {
  const idx = addressListSource.findIndex((item: any) => Number(item.id) === id)
  if (idx < 0) return
  const removed = addressListSource[idx]
  addressListSource.splice(idx, 1)
  if (Number(removed.is_default) === 1 && addressListSource.length > 0) {
    addressListSource[0].is_default = 1
  }
}

const routes: Record<string, any> = {
  'GET /admin/api/index/decor': indexDecor,
  'GET /api/v1/categories': categories,
  'GET /api/v1/products': products,
  'GET /api/v1/products/recommend': recommend,
  'GET /api/v1/products/': productDetail,
  'GET /api/v1/cart': cart,
  'GET /api/v1/orders': orders,
  'GET /api/v1/orders/': orderListSource[0] || null,
  'GET /api/v1/user/coupons': userCoupons,
  'GET /api/v1/user/profile': userProfile,
  'GET /api/v1/addresses': addresses,
  'GET /api/v1/marketing/seckills': seckills,
  'GET /api/v1/marketing/group-buy': groupBuy,
  'GET /api/v1/marketing/bargain': bargain,
}

export function matchMock(method: string, url: string, params?: Record<string, any>): { matched: boolean; data?: any } {
  const upperMethod = method.toUpperCase()
  const [path] = url.split('?')
  const key = `${upperMethod} ${path}`
  const query = { ...parseQuery(url), ...(params || {}) }

  if (upperMethod === 'GET' && path === '/api/v1/orders') {
    const status = Number(query.status || 0)
    return { matched: true, data: listOrders(status) }
  }
  if (upperMethod === 'POST' && path === '/api/v1/addresses') {
    return { matched: true, data: upsertAddress(params || {}) }
  }
  if (upperMethod === 'PUT' && path.startsWith('/api/v1/addresses/')) {
    const id = Number(path.split('/').pop() || 0)
    return { matched: true, data: upsertAddress(params || {}, id) }
  }
  if (upperMethod === 'DELETE' && path.startsWith('/api/v1/addresses/')) {
    const id = Number(path.split('/').pop() || 0)
    removeAddress(id)
    return { matched: true, data: null }
  }
  if (upperMethod === 'GET' && path === '/api/v1/addresses') {
    return { matched: true, data: addressListSource.slice() }
  }
  if (upperMethod === 'POST' && path.startsWith('/api/v1/orders/') && path.endsWith('/pay')) {
    const id = Number(path.split('/')[4] || 0)
    const target = getOrderByID(id)
    if (target && Number(target.status) === 1) {
      target.status = 2
      target.payment_method = target.payment_method || 'wechat'
      target.paid_at = new Date().toISOString()
    }
    return { matched: true, data: null }
  }
  if (upperMethod === 'GET' && path.startsWith('/api/v1/orders/') && path.endsWith('/review')) {
    const id = Number(path.split('/')[4] || 0)
    return { matched: true, data: buildOrderReviewMeta(id) }
  }
  if (upperMethod === 'POST' && path.startsWith('/api/v1/orders/') && path.endsWith('/review')) {
    const id = Number(path.split('/')[4] || 0)
    const target = getOrderByID(id)
    if (target && Array.isArray(target.items)) {
      const mode = String(params?.mode || 'create')
      const logistics = Number(params?.logistics_score || 5)
      const items = Array.isArray(params?.items) ? params.items : []
      if (mode === 'append') {
        for (const item of items) {
          const targetItem = target.items.find((row: any) => Number(row.id) === Number(item.order_item_id))
          if (targetItem?.review) {
            targetItem.review.appends = targetItem.review.appends || []
            targetItem.review.appends.push({
              id: Math.floor(Math.random() * 100000),
              content: String(params?.append_content || ''),
              created_at: new Date().toISOString(),
            })
          }
        }
      } else {
        for (const item of items) {
          const targetItem = target.items.find((row: any) => Number(row.id) === Number(item.order_item_id))
          if (!targetItem) continue
          const now = new Date().toISOString()
          if (!targetItem.review) {
            targetItem.review = {
              id: Math.floor(Math.random() * 100000),
              product_score: Number(item.product_score || 5),
              logistics_score: logistics,
              content: String(item.content || ''),
              edited_times: 0,
              appends: [],
              admin_reply: null,
              created_at: now,
              updated_at: now,
            }
          } else {
            targetItem.review.product_score = Number(item.product_score || 5)
            targetItem.review.logistics_score = logistics
            targetItem.review.content = String(item.content || '')
            targetItem.review.edited_times = Number(targetItem.review.edited_times || 0) + 1
            targetItem.review.updated_at = now
          }
        }
      }
      if ((target.items || []).every((item: any) => !!item.review)) {
        target.status = 4
      }
    }
    return { matched: true, data: null }
  }
  if (upperMethod === 'GET' && path.startsWith('/api/v1/products/') && path.endsWith('/reviews')) {
    const productID = Number(path.split('/')[4] || 0)
    return { matched: true, data: buildProductReview(productID) }
  }

  if (key in routes) {
    const data = routes[key]
    if (upperMethod === 'GET' && path === '/api/v1/products') {
      const sourceList = Array.isArray(data?.list) ? data.list : []
      const keyword = String(query.keyword || '').trim().toLowerCase()
      const categoryID = Number(query.category_id || 0)
      const page = Number(query.page || 1)
      const size = Number(query.size || 20)
      let list = sourceList.slice()
      if (keyword) list = list.filter((item: any) => String(item.title || '').toLowerCase().includes(keyword))
      if (categoryID > 0) list = list.filter((item: any) => Number(item.category_id) === categoryID)
      const offset = Math.max(page - 1, 0) * Math.max(size, 1)
      const pageList = list.slice(offset, offset + size)
      return { matched: true, data: { ...data, list: pageList, total: list.length, page, size } }
    }
    return { matched: true, data }
  }

  for (const pattern of Object.keys(routes)) {
    if (key.startsWith(pattern) && pattern.endsWith('/')) {
      if (pattern === 'GET /api/v1/orders/') {
        const id = Number(path.split('/').pop() || 0)
        const detail = getOrderByID(id) || null
        return { matched: true, data: detail }
      }
      return { matched: true, data: routes[pattern] }
    }
  }

  if (['POST', 'PUT', 'DELETE'].includes(upperMethod)) return { matched: true, data: null }
  return { matched: false }
}
