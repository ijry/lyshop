import { getPreset } from '../../../app/mock/presets/index'
import userCoupons from '../../../app/mock/data/user-coupons.json'
import userProfile from '../../../app/mock/data/user-profile.json'
import addresses from '../../../app/mock/data/addresses.json'

const preset = getPreset()
const indexDecor = preset.indexDecor
const pcDecor = preset.pcDecor
const siteSettings = preset.siteSettings
const categories = preset.categories
const products = preset.products
const productDetail = preset.productDetail
const seckills = preset.seckills
const groupBuy = preset.groupBuy
const bargain = preset.bargain
const recommend = preset.recommend
const cart = preset.cart
const orders = preset.orders
import { afterSaleStatusLabel, shipmentStatusLabel } from '../utils/order-status'

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
const cartListSource = Array.isArray(cart)
  ? clone(cart as any[])
  : []
const addressListSource = Array.isArray(addresses)
  ? JSON.parse(JSON.stringify(addresses))
  : []
let afterSaleSeq = 8000
const afterSalesSource: any[] = []
const productDetailSource = clone(productDetail as any)
const productListSource = Array.isArray((products as any)?.list)
  ? clone((products as any).list)
  : []
type ActivityType = 'seckill' | 'group_buy' | 'bargain'

function asListPayload(input: any): any[] {
  if (Array.isArray(input)) return input
  if (Array.isArray(input?.list)) return input.list
  return []
}

function buildActivityRows(type: ActivityType) {
  const data = type === 'seckill' ? seckills : type === 'group_buy' ? groupBuy : bargain
  const rows = asListPayload(data)
  const result: any[] = []
  for (const act of rows) {
    const products = asListPayload(act?.products)
    for (const item of products) {
      const pid = Number(item?.product_id || 0)
      const source = productListSource.find((p: any) => Number(p?.id || 0) === pid)
      const originPrice = Number(item?.origin_price || source?.price || 0)
      const marketPrice = Number(source?.price || originPrice || 0)
      const skuID = Number(item?.sku_id || pid || 0)
      const activityID = Number(act?.id || 0)
      const activityProductID = Number(item?.id || 0) || Number(`${activityID}${skuID}`)
      result.push({
        activity_product_id: activityProductID,
        activity_id: activityID,
        activity_type: type,
        activity_name: String(act?.name || ''),
        activity_start_at: act?.start_at || null,
        activity_end_at: act?.end_at || null,
        product_id: pid,
        sku_id: skuID,
        title: String(item?.title || source?.title || ''),
        subtitle: String(source?.subtitle || ''),
        cover: String(item?.cover || source?.cover || ''),
        category_id: Number(source?.category_id || 0),
        sales: Number(source?.sales || 0),
        stock: Number(source?.stock || 0),
        origin_price: originPrice || marketPrice,
        price: type === 'bargain' ? Number(item?.origin_price || marketPrice) : Number(item?.activity_price || item?.group_price || marketPrice),
        activity_price: Number(item?.activity_price || item?.group_price || 0),
        start_price: type === 'bargain' ? Number(item?.origin_price || marketPrice) : 0,
        floor_price: Number(item?.floor_price || 0),
        limit_per_order: Number(item?.limit_per_order || 0),
        total_stock_limit: Number(item?.activity_stock || item?.group_stock || item?.total_stock_limit || 0),
        sold_qty: Number(item?.sold_qty || 0),
      })
    }
  }
  return result
}

function findActivityRowByID(activityProductID: number) {
  if (activityProductID <= 0) return null
  const rows = [...buildActivityRows('seckill'), ...buildActivityRows('group_buy'), ...buildActivityRows('bargain')]
  return rows.find((item: any) => Number(item?.activity_product_id || 0) === activityProductID) || null
}

function findActivityRowBySkuID(skuID: number) {
  if (skuID <= 0) return null
  const rows = [...buildActivityRows('seckill'), ...buildActivityRows('group_buy'), ...buildActivityRows('bargain')]
  return rows.find((item: any) => Number(item?.sku_id || 0) === skuID) || null
}

function buildMarketingDetailByActivityProductID(activityProductID: number) {
  const row = findActivityRowByID(activityProductID)
  if (!row) return null
  return {
    activity_product_id: Number(row.activity_product_id || 0),
    activity_id: Number(row.activity_id || 0),
    activity_type: String(row.activity_type || ''),
    activity_name: String(row.activity_name || ''),
    activity_status: 1,
    activity_start_at: row.activity_start_at || null,
    activity_end_at: row.activity_end_at || null,
    product_id: Number(row.product_id || 0),
    sku_id: Number(row.sku_id || 0),
    title: String(row.title || ''),
    subtitle: String(row.subtitle || ''),
    cover: String(row.cover || ''),
    category_id: Number(row.category_id || 0),
    sales: Number(row.sales || 0),
    stock: Number(row.stock || 0),
    origin_price: Number(row.origin_price || 0),
    price: Number(row.price || 0),
    activity_price: Number(row.activity_price || 0),
    start_price: Number(row.start_price || 0),
    floor_price: Number(row.floor_price || 0),
    limit_per_order: Number(row.limit_per_order || 0),
    total_stock_limit: Number(row.total_stock_limit || 0),
    sold_qty: Number(row.sold_qty || 0),
  }
}

function normalizeCartItemActivity(item: any, index: number) {
  const row = Number(item?.activity_product_id || 0) > 0
    ? findActivityRowByID(Number(item.activity_product_id || 0))
    : (index === 0 ? findActivityRowBySkuID(Number(item?.sku_id || 0)) : null)
  if (!row) return item
  return {
    ...item,
    activity_product_id: Number(row.activity_product_id || 0),
    activity_snapshot: {
      activity_product_id: Number(row.activity_product_id || 0),
      activity_id: Number(row.activity_id || 0),
      activity_type: String(row.activity_type || ''),
      activity_name: String(row.activity_name || ''),
      price: Number(row.price || 0),
      origin_price: Number(row.origin_price || 0),
      activity_start_at: row.activity_start_at || null,
      activity_end_at: row.activity_end_at || null,
    },
  }
}

function listCartItems() {
  return cartListSource.map((item: any, index: number) => normalizeCartItemActivity(clone(item), index))
}

function addCartItem(payload: Record<string, any>) {
  const skuID = Number(payload?.sku_id || 0)
  const qty = Math.max(1, Number(payload?.qty || 1))
  const activityProductID = Math.max(0, Number(payload?.activity_product_id || 0))
  if (!skuID) return
  const key = `${skuID}:${activityProductID}`
  const index = cartListSource.findIndex((item: any) => `${Number(item?.sku_id || 0)}:${Number(item?.activity_product_id || 0)}` === key)
  if (index >= 0) {
    cartListSource[index].qty = Math.max(1, Number(cartListSource[index].qty || 1) + qty)
    return
  }
  const detailSkus = Array.isArray(productDetailSource?.skus) ? productDetailSource.skus : []
  const matchedSku = detailSkus.find((item: any) => Number(item?.id || 0) === skuID)
  const productID = Number(matchedSku?.product_id || skuID || 0)
  const sourceProduct = findProductByID(productID) || productDetailSource || {}
  cartListSource.push({
    sku_id: skuID,
    qty,
    activity_product_id: activityProductID,
    product: {
      id: Number(sourceProduct?.id || productID || 0),
      title: String(sourceProduct?.title || `商品${productID}`),
      cover: String(sourceProduct?.cover || ''),
      price: Number(sourceProduct?.price || matchedSku?.price || 0),
    },
    sku: {
      id: skuID,
      product_id: productID,
      attrs: String(matchedSku?.attrs || '[]'),
      price: Number(matchedSku?.price || sourceProduct?.price || 0),
      stock: Number(matchedSku?.stock || sourceProduct?.stock || 0),
    },
  })
}
const favoriteAtMap = new Map<number, string>()
const productFavoriteCountMap = new Map<number, number>()
const afterSaleCaseTypeLabels: Record<string, string> = {
  return: '退货',
  exchange: '换货',
}
const afterSaleActionLabels: Record<string, string> = {
  apply: '提交申请',
  audit: '审核',
  return_ship: '回寄物流',
  receive: '确认收货',
  refund: '退款',
  reship: '补发',
  complete: '完结',
  close: '关闭',
}
const shipmentDirectionLabels: Record<string, string> = {
  outbound: '寄出',
  inbound: '回寄',
}
const shipmentBizTypeLabels: Record<string, string> = {
  initial: '首发',
  reship: '补发',
  return: '回寄',
}

function buildProductIndex() {
  for (const item of productListSource) {
    const id = Number(item?.id || 0)
    if (id <= 0) continue
    productFavoriteCountMap.set(id, Number(item?.favorite_count || 0))
  }
  const detailID = Number(productDetailSource?.id || 0)
  if (detailID > 0 && !productFavoriteCountMap.has(detailID)) {
    productFavoriteCountMap.set(detailID, Number(productDetailSource?.favorite_count || 0))
  }
}

function findProductByID(id: number) {
  if (id <= 0) return null
  const hit = productListSource.find((item: any) => Number(item?.id || 0) === id)
  if (hit) return hit
  if (Number(productDetailSource?.id || 0) === id) return productDetailSource
  return null
}

function enrichFavoriteFields(row: any) {
  const id = Number(row?.id || 0)
  const count = productFavoriteCountMap.get(id) ?? Number(row?.favorite_count || 0)
  return {
    ...row,
    favorite_count: Math.max(0, count),
    is_favorited: favoriteAtMap.has(id),
  }
}

function favoriteProductByID(id: number) {
  if (id <= 0 || favoriteAtMap.has(id)) return
  favoriteAtMap.set(id, new Date().toISOString())
  const prev = productFavoriteCountMap.get(id) ?? Number(findProductByID(id)?.favorite_count || 0)
  productFavoriteCountMap.set(id, prev + 1)
}

function unfavoriteProductByID(id: number) {
  if (id <= 0 || !favoriteAtMap.has(id)) return
  favoriteAtMap.delete(id)
  const prev = productFavoriteCountMap.get(id) ?? Number(findProductByID(id)?.favorite_count || 0)
  productFavoriteCountMap.set(id, Math.max(0, prev - 1))
}

function listUserFavorites(page = 1, size = 20) {
  const rows = Array.from(favoriteAtMap.entries())
    .sort((a, b) => String(b[1]).localeCompare(String(a[1])))
    .map(([id, favoritedAt]) => {
      const source = findProductByID(Number(id))
      if (!source) return null
      return {
        ...enrichFavoriteFields(clone(source)),
        favorited_at: favoritedAt,
      }
    })
    .filter(Boolean) as any[]
  const safePage = Math.max(1, Number(page) || 1)
  const safeSize = Math.max(1, Number(size) || 20)
  const offset = (safePage - 1) * safeSize
  return {
    list: rows.slice(offset, offset + safeSize),
    total: rows.length,
    page: safePage,
    size: safeSize,
  }
}

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

function withOrderShipmentLabels(ship: any) {
  return {
    ...ship,
    direction_label: shipmentDirectionLabels[String(ship?.direction || '')] || String(ship?.direction || ''),
    biz_type_label: shipmentBizTypeLabels[String(ship?.biz_type || '')] || String(ship?.biz_type || ''),
    logistics_status_label: shipmentStatusLabel(ship?.logistics_status),
  }
}

function withAfterSaleLabels(row: any) {
  const data = clone(row || {})
  const logs = Array.isArray(data.logs)
    ? data.logs.map((log: any) => ({
      ...log,
      from_status_label: String(log?.from_status || '').trim() ? afterSaleStatusLabel(log.from_status) : '',
      to_status_label: afterSaleStatusLabel(log?.to_status),
      action_label: afterSaleActionLabels[String(log?.action || '')] || String(log?.action || ''),
    }))
    : []
  const shipments = Array.isArray(data.shipments)
    ? data.shipments.map((ship: any) => ({
      ...ship,
      direction_label: shipmentDirectionLabels[String(ship?.direction || '')] || String(ship?.direction || ''),
      biz_type_label: shipmentBizTypeLabels[String(ship?.biz_type || '')] || String(ship?.biz_type || ''),
      logistics_status_label: shipmentStatusLabel(ship?.logistics_status),
    }))
    : []
  return {
    ...data,
    status_label: afterSaleStatusLabel(data?.status),
    case_type_label: afterSaleCaseTypeLabels[String(data?.case_type || '')] || String(data?.case_type || ''),
    logs,
    shipments,
  }
}

function buildAfterSaleSummary(orderID: number) {
  const rows = afterSalesSource
    .filter((row: any) => Number(row.order_id) === Number(orderID))
    .sort((a: any, b: any) => Number(b.id) - Number(a.id))
  if (!rows.length) {
    return {
      in_progress_count: 0,
      has_open_case: false,
      latest_status: '',
      latest_status_label: '',
      latest_case_id: 0,
      can_apply: true,
    }
  }
  const latestStatus = String(rows[0].status || '')
  const openCount = rows.filter((row: any) => statusOpen(String(row.status || ''))).length
  return {
    in_progress_count: openCount,
    has_open_case: openCount > 0,
    latest_status: latestStatus,
    latest_status_label: afterSaleStatusLabel(latestStatus),
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
  order.shipments = order.shipments.map((ship: any) => withOrderShipmentLabels(ship))
  order.latest_shipment = order.shipments?.[0] || null
  order.after_sale_summary = buildAfterSaleSummary(Number(order.id))
}

buildProductIndex()
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
  'GET /api/v1/decor/pc': { components: pcDecor },
  'GET /api/v1/pc/decor': { components: pcDecor },
  'GET /api/v1/site-settings': siteSettings,
  'GET /api/v1/categories': categories,
  'GET /api/v1/products': { list: productListSource, total: productListSource.length, page: 1, size: 20 },
  'GET /api/v1/products/recommend': recommend,
  'GET /api/v1/cart': cartListSource,
  'GET /api/v1/orders': orders,
  'GET /api/v1/orders/': orderListSource[0] || null,
  'GET /api/v1/user/coupons': userCoupons,
  'GET /api/v1/user/profile': userProfile,
  'GET /api/v1/addresses': addresses,
  'GET /api/v1/seckill/products': { list: buildActivityRows('seckill'), total: buildActivityRows('seckill').length, page: 1, size: 20 },
  'GET /api/v1/group-buy/products': { list: buildActivityRows('group_buy'), total: buildActivityRows('group_buy').length, page: 1, size: 20 },
  'GET /api/v1/bargain/products': { list: buildActivityRows('bargain'), total: buildActivityRows('bargain').length, page: 1, size: 20 },
  'GET /api/v1/marketing/seckill/products': { list: buildActivityRows('seckill'), total: buildActivityRows('seckill').length, page: 1, size: 20 },
  'GET /api/v1/marketing/group-buy/products': { list: buildActivityRows('group_buy'), total: buildActivityRows('group_buy').length, page: 1, size: 20 },
  'GET /api/v1/marketing/bargain/products': { list: buildActivityRows('bargain'), total: buildActivityRows('bargain').length, page: 1, size: 20 },
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
  if (upperMethod === 'GET' && path === '/api/v1/cart') {
    return { matched: true, data: listCartItems() }
  }
  if (upperMethod === 'POST' && path === '/api/v1/cart/add') {
    addCartItem(params || {})
    return { matched: true, data: null }
  }
  if (upperMethod === 'POST' && /^\/api\/v1\/products\/\d+\/favorite$/.test(path)) {
    const id = Number(path.split('/')[4] || 0)
    favoriteProductByID(id)
    return { matched: true, data: null }
  }
  if (upperMethod === 'DELETE' && /^\/api\/v1\/products\/\d+\/favorite$/.test(path)) {
    const id = Number(path.split('/')[4] || 0)
    unfavoriteProductByID(id)
    return { matched: true, data: null }
  }
  if (upperMethod === 'GET' && path === '/api/v1/user/favorites') {
    return { matched: true, data: listUserFavorites(Number(query.page || 1), Number(query.size || 20)) }
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
    return { matched: true, data: row ? withAfterSaleLabels(row) : null }
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
      order.shipments.unshift(withOrderShipmentLabels(shipment))
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
  if (upperMethod === 'GET' && /^\/api\/v1\/products\/\d+$/.test(path)) {
    const id = Number(path.split('/').pop() || 0)
    const source = findProductByID(id)
    if (!source) return { matched: true, data: null }
    const detail = Number(productDetailSource?.id || 0) === id
      ? clone(productDetailSource)
      : {
          ...clone(source),
          skus: [],
          images: [],
          detail: { version: 1, blocks: [] },
        }
    return { matched: true, data: enrichFavoriteFields(detail) }
  }
  if (upperMethod === 'GET' && /^\/api\/v1\/marketing\/activity-products\/\d+$/.test(path)) {
    const activityProductID = Number(path.split('/').pop() || 0)
    return { matched: true, data: buildMarketingDetailByActivityProductID(activityProductID) }
  }

  if (key in routes) {
    const data = routes[key]
    if (upperMethod === 'GET' && (path === '/api/v1/seckill/products' || path === '/api/v1/group-buy/products' || path === '/api/v1/bargain/products' || path === '/api/v1/marketing/seckill/products' || path === '/api/v1/marketing/group-buy/products' || path === '/api/v1/marketing/bargain/products')) {
      const sourceList = (path === '/api/v1/seckill/products' || path === '/api/v1/marketing/seckill/products')
        ? buildActivityRows('seckill')
        : (path === '/api/v1/group-buy/products' || path === '/api/v1/marketing/group-buy/products')
          ? buildActivityRows('group_buy')
          : buildActivityRows('bargain')
      const keyword = String(query.keyword || '').trim().toLowerCase()
      const categoryID = Number(query.category_id || 0)
      const minPrice = Number(query.min_price || 0)
      const maxPrice = Number(query.max_price || 0)
      const sortBy = String(query.sort_by || 'price')
      const sortOrder = String(query.sort_order || 'asc')
      const page = Number(query.page || 1)
      const size = Number(query.size || 20)
      let list = sourceList.slice()
      if (keyword) list = list.filter((item: any) => String(item.title || '').toLowerCase().includes(keyword))
      if (categoryID > 0) list = list.filter((item: any) => Number(item.category_id || 0) === categoryID)
      if (minPrice > 0) list = list.filter((item: any) => Number(item.price || 0) >= minPrice)
      if (maxPrice > 0) list = list.filter((item: any) => Number(item.price || 0) <= maxPrice)
      list.sort((a: any, b: any) => {
        const left = sortBy === 'sales' ? Number(a.sales || 0) : Number(a.price || 0)
        const right = sortBy === 'sales' ? Number(b.sales || 0) : Number(b.price || 0)
        return sortOrder === 'desc' ? right - left : left - right
      })
      const offset = Math.max(page - 1, 0) * Math.max(size, 1)
      const pageList = list.slice(offset, offset + size)
      return { matched: true, data: { list: pageList, total: list.length, page, size } }
    }
    if (upperMethod === 'GET' && path === '/api/v1/products') {
      const sourceList = productListSource.slice()
      const keyword = String(query.keyword || '').trim().toLowerCase()
      const categoryID = Number(query.category_id || 0)
      const page = Number(query.page || 1)
      const size = Number(query.size || 20)
      let list = sourceList.slice()
      if (keyword) list = list.filter((item: any) => String(item.title || '').toLowerCase().includes(keyword))
      if (categoryID > 0) list = list.filter((item: any) => Number(item.category_id) === categoryID)
      const offset = Math.max(page - 1, 0) * Math.max(size, 1)
      const pageList = list.slice(offset, offset + size).map((item: any) => enrichFavoriteFields(clone(item)))
      return { matched: true, data: { list: pageList, total: list.length, page, size } }
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
