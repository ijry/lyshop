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

function clone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value))
}

const orderListSource = Array.isArray((orders as any)?.list)
  ? JSON.parse(JSON.stringify((orders as any).list))
  : []
const addressListSource = Array.isArray(addresses)
  ? JSON.parse(JSON.stringify(addresses))
  : []
let afterSaleSeq = 8000
const afterSalesSource: any[] = []

function nextUploadURL() {
  const n = Math.floor(Math.random() * 1000) + 200
  return `https://picsum.photos/640/640?random=${n}`
}

function getOrderByID(id: number) {
  return orderListSource.find((item: any) => Number(item.id) === id)
}

function nextAfterSaleCaseNo() {
  afterSaleSeq += 1
  return `AS${Date.now()}${afterSaleSeq}`
}

function statusOpen(status: string) {
  return !['completed', 'rejected', 'closed'].includes(status)
}

function buildAfterSaleSummary(orderID: number) {
  const rows = afterSalesSource
    .filter((row: any) => Number(row.order_id) === Number(orderID))
    .sort((a: any, b: any) => Number(b.id) - Number(a.id))
  if (!rows.length) {
    return { in_progress_count: 0, has_open_case: false, latest_status: '', latest_case_id: 0, can_apply: true }
  }
  const openCount = rows.filter((row: any) => statusOpen(String(row.status || ''))).length
  return {
    in_progress_count: openCount,
    has_open_case: openCount > 0,
    latest_status: String(rows[0].status || ''),
    latest_case_id: Number(rows[0].id || 0),
    can_apply: openCount === 0,
  }
}

function ensureOrderExt(order: any) {
  if (!order) return
  if (!Array.isArray(order.shipments)) {
    order.shipments = order.tracking_no ? [{
      id: Number(order.id) * 10 + 1,
      order_id: Number(order.id),
      after_sale_case_id: 0,
      direction: 'outbound',
      biz_type: 'initial',
      company: '顺丰',
      tracking_no: order.tracking_no,
      logistics_status: 'shipped',
      remark: '',
      created_by_type: 'admin',
      created_by_id: 0,
      created_at: order.paid_at || order.created_at,
    }] : []
  }
  order.latest_shipment = order.shipments?.[0] || null
  order.after_sale_summary = buildAfterSaleSummary(Number(order.id))
}

for (const order of orderListSource) ensureOrderExt(order)
if (orderListSource[1]) {
  const seedCaseID = ++afterSaleSeq
  afterSalesSource.push({
    id: seedCaseID,
    order_id: Number(orderListSource[1].id),
    user_id: Number(orderListSource[1].user_id || 1),
    case_no: nextAfterSaleCaseNo(),
    case_type: 'return',
    status: 'approved_wait_user_return',
    reason: '尺寸不合适',
    apply_content: '试穿后不合适',
    apply_images: ['https://picsum.photos/200/200?random=991'],
    items: [{
      id: Math.floor(Math.random() * 100000),
      case_id: seedCaseID,
      order_item_id: Number(orderListSource[1].items?.[0]?.id || 0),
      qty: 1,
    }],
    logs: [{
      id: Math.floor(Math.random() * 100000),
      case_id: seedCaseID,
      from_status: '',
      to_status: 'applied',
      action: 'apply',
      operator_type: 'user',
      operator_id: Number(orderListSource[1].user_id || 1),
      content: '提交售后申请',
      created_at: new Date().toISOString(),
    }, {
      id: Math.floor(Math.random() * 100000),
      case_id: seedCaseID,
      from_status: 'applied',
      to_status: 'approved_wait_user_return',
      action: 'audit',
      operator_type: 'admin',
      operator_id: 1,
      content: '售后审核通过',
      created_at: new Date().toISOString(),
    }],
    shipments: [],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  })
}
for (const order of orderListSource) touchOrderAfterSaleSummary(Number(order.id))

function getAfterSaleByID(id: number) {
  return afterSalesSource.find((row: any) => Number(row.id) === id)
}

function listAfterSalesByOrder(orderID: number) {
  return afterSalesSource
    .filter((row: any) => Number(row.order_id) === Number(orderID))
    .sort((a: any, b: any) => Number(b.id) - Number(a.id))
}

function touchOrderAfterSaleSummary(orderID: number) {
  const order = getOrderByID(orderID)
  if (!order) return
  order.after_sale_summary = buildAfterSaleSummary(orderID)
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
        images: Array.isArray(item.review.images) ? item.review.images.map((u: any) => String(u || '')) : [],
        edited_times: Number(item.review.edited_times || 0),
        user_nickname: '演示用户',
        user_avatar: '',
        created_at: item.review.created_at || order.created_at,
        updated_at: item.review.updated_at || item.review.created_at || order.created_at,
        appends: Array.isArray(item.review.appends) ? item.review.appends.map((ap: any) => ({
          id: Number(ap.id || 0),
          content: String(ap.content || ''),
          images: Array.isArray(ap.images) ? ap.images.map((u: any) => String(u || '')) : [],
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
    images: Array.isArray(item.review?.images) ? item.review.images.map((u: any) => String(u || '')) : [],
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
  for (const item of list) ensureOrderExt(item)
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
      const appendImages = Array.isArray(params?.append_images)
        ? params.append_images.map((u: any) => String(u || '').trim()).filter(Boolean).slice(0, 9)
        : []
      if (mode === 'append') {
        const appendContent = String(params?.append_content || '').trim()
        if (!appendContent && !appendImages.length) return { matched: true, data: null }
        for (const item of items) {
          const targetItem = target.items.find((row: any) => Number(row.id) === Number(item.order_item_id))
          if (targetItem?.review) {
            targetItem.review.appends = targetItem.review.appends || []
            targetItem.review.appends.push({
              id: Math.floor(Math.random() * 100000),
              content: appendContent,
              images: appendImages.slice(),
              created_at: new Date().toISOString(),
            })
          }
        }
      } else {
        for (const item of items) {
          const targetItem = target.items.find((row: any) => Number(row.id) === Number(item.order_item_id))
          if (!targetItem) continue
          const now = new Date().toISOString()
          const images = Array.isArray((item as any)?.images)
            ? (item as any).images.map((u: any) => String(u || '').trim()).filter(Boolean).slice(0, 9)
            : []
          if (!targetItem.review) {
            targetItem.review = {
              id: Math.floor(Math.random() * 100000),
              product_score: Number(item.product_score || 5),
              logistics_score: logistics,
              content: String(item.content || ''),
              images: images.slice(),
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
            targetItem.review.images = images.slice()
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
  if (upperMethod === 'POST' && path === '/api/v1/upload') {
    const url = nextUploadURL()
    return { matched: true, data: { path: `demo/${Date.now()}.jpg`, url, size: 10240, mime: 'image/jpeg' } }
  }
  if (upperMethod === 'POST' && path.startsWith('/api/v1/orders/') && path.endsWith('/after-sales')) {
    const orderID = Number(path.split('/')[4] || 0)
    const order = getOrderByID(orderID)
    if (!order) return { matched: true, data: null }
    const body = params || {}
    const items = Array.isArray(body.items) ? body.items : []
    const caseType = String(body.case_type || 'return') === 'exchange' ? 'exchange' : 'return'
    const row = {
      id: ++afterSaleSeq,
      order_id: orderID,
      user_id: Number(order.user_id || 1),
      case_no: nextAfterSaleCaseNo(),
      case_type: caseType,
      status: 'applied',
      reason: String(body.reason || ''),
      apply_content: String(body.apply_content || ''),
      apply_images: Array.isArray(body.apply_images) ? body.apply_images.map((u: any) => String(u || '')).filter(Boolean) : [],
      items: items.map((item: any) => ({
        id: Math.floor(Math.random() * 100000),
        case_id: afterSaleSeq,
        order_item_id: Number(item.order_item_id || 0),
        qty: Math.max(1, Number(item.qty || 1)),
      })),
      logs: [{
        id: Math.floor(Math.random() * 100000),
        case_id: afterSaleSeq,
        from_status: '',
        to_status: 'applied',
        action: 'apply',
        operator_type: 'user',
        operator_id: Number(order.user_id || 1),
        content: '提交售后申请',
        created_at: new Date().toISOString(),
      }],
      shipments: [],
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }
    afterSalesSource.unshift(row)
    order.status = 5
    touchOrderAfterSaleSummary(orderID)
    return { matched: true, data: { id: row.id } }
  }
  if (upperMethod === 'GET' && path.startsWith('/api/v1/after-sales/')) {
    const caseID = Number(path.split('/').pop() || 0)
    const row = getAfterSaleByID(caseID)
    return { matched: true, data: row ? JSON.parse(JSON.stringify(row)) : null }
  }
  if (upperMethod === 'POST' && path.startsWith('/api/v1/after-sales/') && path.endsWith('/return-shipments')) {
    const caseID = Number(path.split('/')[4] || 0)
    const row = getAfterSaleByID(caseID)
    if (!row) return { matched: true, data: null }
    const body = params || {}
    const now = new Date().toISOString()
    const shipment = {
      id: Math.floor(Math.random() * 100000),
      order_id: Number(row.order_id),
      after_sale_case_id: caseID,
      direction: 'inbound',
      biz_type: 'return',
      company: String(body.company || ''),
      tracking_no: String(body.tracking_no || ''),
      logistics_status: 'shipped',
      remark: String(body.remark || ''),
      created_by_type: 'user',
      created_by_id: Number(row.user_id || 1),
      created_at: now,
    }
    row.shipments = Array.isArray(row.shipments) ? row.shipments : []
    row.shipments.unshift(shipment)
    row.status = 'user_returning'
    row.logs = Array.isArray(row.logs) ? row.logs : []
    row.logs.push({
      id: Math.floor(Math.random() * 100000),
      case_id: caseID,
      from_status: 'approved_wait_user_return',
      to_status: 'user_returning',
      action: 'return_ship',
      operator_type: 'user',
      operator_id: Number(row.user_id || 1),
      content: '用户提交回寄物流',
      created_at: now,
    })
    const order = getOrderByID(Number(row.order_id))
    if (order) {
      order.shipments = Array.isArray(order.shipments) ? order.shipments : []
      order.shipments.unshift(shipment)
      order.latest_shipment = order.shipments[0] || null
      touchOrderAfterSaleSummary(Number(row.order_id))
    }
    return { matched: true, data: null }
  }
  if (upperMethod === 'GET' && path.startsWith('/api/v1/orders/') && !path.endsWith('/review')) {
    const id = Number(path.split('/').pop() || 0)
    const detail = getOrderByID(id) || null
    if (!detail) return { matched: true, data: null }
    ensureOrderExt(detail)
    return { matched: true, data: clone(detail) }
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
      return { matched: true, data: routes[pattern] }
    }
  }

  if (['POST', 'PUT', 'DELETE'].includes(upperMethod)) return { matched: true, data: null }
  return { matched: false }
}
