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

const orderListSource = Array.isArray((orders as any)?.list)
  ? JSON.parse(JSON.stringify((orders as any).list))
  : []
const addressListSource = Array.isArray(addresses)
  ? JSON.parse(JSON.stringify(addresses))
  : []
const productDetailSource = JSON.parse(JSON.stringify(productDetail as any))
const productReviewMap = new Map<number, any[]>()
const reviewIndexMap = new Map<number, any>()
let reviewSeq = 5000
let appendSeq = 9000
let replySeq = 10000

function clone<T>(v: T): T {
  return JSON.parse(JSON.stringify(v))
}

function nextReviewID() {
  reviewSeq += 1
  return reviewSeq
}

function nextAppendID() {
  appendSeq += 1
  return appendSeq
}

function nextReplyID() {
	replySeq += 1
	return replySeq
}

function nextUploadURL() {
	const n = Math.floor(Math.random() * 1000) + 100
	return `https://picsum.photos/640/640?random=${n}`
}

function formatNowISO() {
  return new Date().toISOString()
}

function getOrderByID(id: number) {
  return orderListSource.find((item: any) => Number(item.id) === id)
}

function normalizeMode(mode: string | undefined) {
  if (mode === 'edit' || mode === 'append' || mode === 'create') return mode
  return 'create'
}

function ensureProductReviewBucket(productID: number) {
  if (!productReviewMap.has(productID)) {
    productReviewMap.set(productID, [])
  }
  return productReviewMap.get(productID) as any[]
}

function hydrateReviewsFromOrders() {
  productReviewMap.clear()
  reviewIndexMap.clear()
  for (const order of orderListSource) {
    const orderNo = String(order.order_no || '')
    const items = Array.isArray(order.items) ? order.items : []
    for (const item of items) {
      const rv = item?.review
      if (!rv || Number(rv.id || 0) <= 0) continue
      const review = {
        id: Number(rv.id),
        order_id: Number(order.id),
        order_no: orderNo,
        order_item_id: Number(item.id),
        product_id: Number(item.product_id || 0),
        product_score: Number(rv.product_score || 5),
        logistics_score: Number(rv.logistics_score || 5),
        content: String(rv.content || ''),
        images: Array.isArray(rv.images) ? rv.images.map((u: any) => String(u || '')) : [],
        edited_times: Number(rv.edited_times || 0),
        user_id: 1,
        user_nickname: '演示用户',
        user_avatar: '',
        created_at: rv.created_at || order.created_at || formatNowISO(),
        updated_at: rv.updated_at || rv.created_at || order.created_at || formatNowISO(),
        appends: Array.isArray(rv.appends) ? rv.appends.map((ap: any) => ({
          id: Number(ap.id || nextAppendID()),
          review_id: Number(rv.id),
          content: String(ap.content || ''),
          images: Array.isArray(ap.images) ? ap.images.map((u: any) => String(u || '')) : [],
          created_at: ap.created_at || formatNowISO(),
          updated_at: ap.updated_at || ap.created_at || formatNowISO(),
        })) : [],
        admin_reply: rv.admin_reply ? {
          id: Number(rv.admin_reply.id || nextReplyID()),
          review_id: Number(rv.id),
          content: String(rv.admin_reply.content || ''),
          created_at: rv.admin_reply.created_at || formatNowISO(),
          updated_at: rv.admin_reply.updated_at || rv.admin_reply.created_at || formatNowISO(),
          admin_id: Number(rv.admin_reply.admin_id || 1),
        } : null,
      }
      reviewSeq = Math.max(reviewSeq, review.id)
      for (const ap of review.appends) appendSeq = Math.max(appendSeq, Number(ap.id || 0))
      if (review.admin_reply) replySeq = Math.max(replySeq, Number(review.admin_reply.id || 0))

      reviewIndexMap.set(review.id, review)
      ensureProductReviewBucket(review.product_id).push(review)
    }
  }
  for (const [productID, list] of productReviewMap.entries()) {
    list.sort((a: any, b: any) => Number(b.id) - Number(a.id))
    productReviewMap.set(productID, list)
  }
}

hydrateReviewsFromOrders()

function buildReviewSummary(productID: number) {
  const list = ensureProductReviewBucket(productID)
  if (!list.length) {
    return { avg_product_score: 0, avg_logistics_score: 0, total: 0 }
  }
  const productSum = list.reduce((sum: number, item: any) => sum + Number(item.product_score || 0), 0)
  const logisticsSum = list.reduce((sum: number, item: any) => sum + Number(item.logistics_score || 0), 0)
  return {
    avg_product_score: Number((productSum / list.length).toFixed(1)),
    avg_logistics_score: Number((logisticsSum / list.length).toFixed(1)),
    total: list.length,
  }
}

function buildProductReviews(productID: number, page = 1, size = 20) {
  const list = ensureProductReviewBucket(productID)
  const safePage = Math.max(1, Number(page) || 1)
  const safeSize = Math.max(1, Number(size) || 20)
  const offset = (safePage - 1) * safeSize
  const pageList = list.slice(offset, offset + safeSize)
  return {
    summary: buildReviewSummary(productID),
    list: clone(pageList),
    total: list.length,
    page: safePage,
    size: safeSize,
  }
}

function buildOrderReviewMeta(orderID: number) {
  const order = getOrderByID(orderID)
  if (!order) return null
  const items = Array.isArray(order.items) ? order.items : []
  const options = items.map((item: any) => {
    const rv = item.review || null
    return {
      order_item_id: Number(item.id || 0),
      review_id: Number(rv?.id || 0),
      has_review: Number(rv?.id || 0) > 0,
      product_id: Number(item.product_id || 0),
      product_title: String(item.title || ''),
      product_cover: String(item.cover || ''),
      product_score: Number(rv?.product_score || 5),
      logistics_score: Number(rv?.logistics_score || 5),
      content: String(rv?.content || ''),
      images: Array.isArray(rv?.images) ? rv.images.map((u: any) => String(u || '')) : [],
    }
  })
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

function applyReviewToOrderItem(orderID: number, payload: any) {
  const order = getOrderByID(orderID)
  if (!order) return { ok: false, msg: '订单不存在' }
  const mode = normalizeMode(String(payload?.mode || 'create'))
  const logisticsScore = Math.max(1, Math.min(5, Number(payload?.logistics_score || 5)))
  const sourceItems = Array.isArray(payload?.items) ? payload.items : []
  const items = sourceItems.length
    ? sourceItems
    : (Array.isArray(order.items) ? order.items.map((item: any) => ({
      order_item_id: Number(item.id || 0),
      product_score: 5,
      content: String(payload?.content || ''),
    })) : [])
  const appendContent = String(payload?.append_content || payload?.content || '').trim()
  const appendImages = Array.isArray(payload?.append_images)
    ? payload.append_images.map((u: any) => String(u || '').trim()).filter(Boolean).slice(0, 9)
    : []
  if (mode === 'append' && !appendContent && !appendImages.length) {
    return { ok: false, msg: '追评内容或图片不能为空' }
  }

  let touched = false
  for (const reqItem of items) {
    const orderItemID = Number(reqItem?.order_item_id || 0)
    const target = Array.isArray(order.items)
      ? order.items.find((item: any) => Number(item.id) === orderItemID)
      : null
    if (!target) continue
    const productID = Number(target.product_id || 0)
    if (mode === 'append') {
      if (!target.review) continue
      const appendRow = {
        id: nextAppendID(),
        review_id: Number(target.review.id),
        content: appendContent,
        images: appendImages.slice(),
        created_at: formatNowISO(),
        updated_at: formatNowISO(),
      }
      if (!Array.isArray(target.review.appends)) target.review.appends = []
      target.review.appends.push(appendRow)
      const review = reviewIndexMap.get(Number(target.review.id))
      if (review) {
        if (!Array.isArray(review.appends)) review.appends = []
        review.appends.push(appendRow)
        review.updated_at = formatNowISO()
      }
      touched = true
      continue
    }

    const score = Math.max(1, Math.min(5, Number(reqItem?.product_score || 5)))
    const content = String(reqItem?.content || payload?.content || '').trim()
    const images = Array.isArray(reqItem?.images)
      ? reqItem.images.map((u: any) => String(u || '').trim()).filter(Boolean).slice(0, 9)
      : []
    const now = formatNowISO()
    if (!target.review) {
      if (mode !== 'create') continue
      const reviewID = nextReviewID()
      const review = {
        id: reviewID,
        order_id: Number(order.id),
        order_no: order.order_no,
        order_item_id: orderItemID,
        product_id: productID,
        product_score: score,
        logistics_score: logisticsScore,
        content,
        images: images.slice(),
        edited_times: 0,
        user_id: 1,
        user_nickname: '演示用户',
        user_avatar: '',
        created_at: now,
        updated_at: now,
        appends: [],
        admin_reply: null,
      }
      target.review = {
        id: reviewID,
        review_id: reviewID,
        has_review: true,
        product_score: score,
        logistics_score: logisticsScore,
        content,
        images: images.slice(),
        edited_times: 0,
        created_at: now,
        updated_at: now,
        appends: [],
        admin_reply: null,
      }
      ensureProductReviewBucket(productID).unshift(review)
      reviewIndexMap.set(reviewID, review)
      touched = true
      continue
    }

    // edit
    target.review.product_score = score
    target.review.logistics_score = logisticsScore
    target.review.content = content
    target.review.images = images.slice()
    target.review.edited_times = Number(target.review.edited_times || 0) + 1
    target.review.updated_at = now
    const review = reviewIndexMap.get(Number(target.review.id))
    if (review) {
      review.product_score = score
      review.logistics_score = logisticsScore
      review.content = content
      review.images = images.slice()
      review.edited_times = Number(review.edited_times || 0) + 1
      review.updated_at = now
    }
    touched = true
  }

  if (!touched) {
    return { ok: false, msg: mode === 'append' ? '只能对根评价追加' : '未找到可评价商品' }
  }

  const allReviewed = Array.isArray(order.items) && order.items.length > 0
    ? order.items.every((item: any) => !!item.review?.id)
    : false
  if (allReviewed) order.status = 4

  return { ok: true, msg: '' }
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
    if (payload.is_default === 1) {
      addressListSource.forEach((item: any) => { item.is_default = 0 })
    }
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
  'GET /api/v1/index/decor': indexDecor,
  'GET /api/v1/categories': categories,
  'GET /api/v1/products': products,
  'GET /api/v1/products/recommend': recommend,
  'GET /api/v1/products/': productDetailSource,
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
      { id: 4, type: 1, points: 2200, remark: '活动签到奖励', created_at: '2026-04-01T09:00:00Z' },
    ],
    total: 4,
  },
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
  'GET /api/v1/checkin/status': { checked_today: false, consecutive_days: 3, month_dates: ['2026-05-17', '2026-05-20', '2026-05-21'], month_count: 3, month_points: 40 },
  'GET /api/v1/checkin/rules': [{ day: 0, points: 10 }, { day: 3, points: 20 }, { day: 7, points: 50 }],
  'POST /api/v1/checkin': { points: 20, consecutive_days: 4 },
  'GET /api/v1/im/session': { id: 1, user_id: 1, status: 2 },
  'GET /api/v1/im/messages': { list: [], total: 0 },
  'GET /api/v1/messages/unread': { system: 2, order: 1, marketing: 3, im: 0 },
  'GET /api/v1/messages': {
    list: [
      { id: 1, group: 'system', title: '系统升级通知', content: 'LYShop 已升级至 2.0 版本', is_read: 0, created_at: '2026-05-22T10:00:00Z' },
      { id: 2, group: 'order', title: '订单发货通知', content: '您的订单已发货', is_read: 0, created_at: '2026-05-21T14:00:00Z' },
      { id: 3, group: 'marketing', title: '618大促即将开始', content: '全场满300减50', is_read: 0, created_at: '2026-05-20T08:00:00Z' },
      { id: 4, group: 'marketing', title: '优惠券到账', content: '满100减20优惠券已到账', is_read: 1, created_at: '2026-05-18T12:00:00Z' },
      { id: 5, group: 'system', title: '欢迎注册', content: '新人专享优惠等你来领', is_read: 1, created_at: '2026-05-15T10:00:00Z' },
    ],
    total: 5,
  },
}

/**
 * Try to match a mock route. Returns { matched: true, data } or { matched: false }.
 */
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
    const meta = buildOrderReviewMeta(id)
    if (!meta) return { matched: true, data: null }
    return { matched: true, data: meta }
  }
  if (upperMethod === 'POST' && path.startsWith('/api/v1/orders/') && path.endsWith('/review')) {
    const id = Number(path.split('/')[4] || 0)
    const result = applyReviewToOrderItem(id, params || {})
    if (!result.ok) {
      // mock mode has no custom error envelope, return null to keep dev flow smooth
      return { matched: true, data: null }
    }
    return { matched: true, data: null }
  }
  if (upperMethod === 'POST' && path === '/api/v1/upload') {
    const url = nextUploadURL()
    return {
      matched: true,
      data: { path: `demo/${Date.now()}.jpg`, url, size: 10240, mime: 'image/jpeg' },
    }
  }
  if (upperMethod === 'GET' && path.startsWith('/api/v1/products/') && path.endsWith('/reviews')) {
    const productID = Number(path.split('/')[4] || 0)
    const page = Number(query.page || 1)
    const size = Number(query.size || 20)
    return { matched: true, data: buildProductReviews(productID, page, size) }
  }

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
    return { matched: true, data }
  }

  // Prefix match (for routes with path params like /products/:id)
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

  // Default: return empty success for unmatched POST/PUT/DELETE
  if (['POST', 'PUT', 'DELETE'].includes(upperMethod)) {
    return { matched: true, data: null }
  }

  return { matched: false }
}

export const isMockEnabled = true
