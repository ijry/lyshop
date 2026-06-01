/**
 * Mock data registry.
 * Key = "METHOD URL" (e.g. "GET /api/v1/products"), value = mock response data.
 * Data is the `data` field inside { code: 0, msg: "success", data: ... }.
 */

import { getPreset } from './presets/index'
import userCoupons from './data/user-coupons.json'
import userProfile from './data/user-profile.json'
import addresses from './data/addresses.json'

const preset = getPreset()
const indexDecor = preset.indexDecor
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

const orderListSource = Array.isArray((orders as any)?.list)
  ? JSON.parse(JSON.stringify((orders as any).list))
  : []
const userCouponSource = Array.isArray(userCoupons)
  ? JSON.parse(JSON.stringify(userCoupons))
  : []
const claimableCouponSource = [
  { id: 1, name: '新人满100减20', type: 1, min_amount: 100, discount: 20, per_limit: 1, total_count: 1000, claimed_count: 120, claimed_by_me: 1, can_claim: false, status: 1 },
  { id: 2, name: '全场9折券', type: 2, min_amount: 0, discount: 0.9, per_limit: 1, total_count: 500, claimed_count: 80, claimed_by_me: 0, can_claim: true, status: 1 },
  { id: 3, name: '无门槛5元券', type: 3, min_amount: 0, discount: 5, per_limit: 2, total_count: 2000, claimed_count: 650, claimed_by_me: 0, can_claim: true, status: 1 },
]
const addressListSource = Array.isArray(addresses)
  ? JSON.parse(JSON.stringify(addresses))
  : []
let afterSaleSeq = 8000
const afterSalesSource: any[] = []
const productDetailSource = JSON.parse(JSON.stringify(productDetail as any))
const productListSource = Array.isArray((products as any)?.list)
  ? JSON.parse(JSON.stringify((products as any).list))
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
const favoriteAtMap = new Map<number, string>()
const productFavoriteCountMap = new Map<number, number>()
const productReviewMap = new Map<number, any[]>()
const reviewIndexMap = new Map<number, any>()
let reviewSeq = 5000
let appendSeq = 9000
let replySeq = 10000
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

function clone<T>(v: T): T {
  return JSON.parse(JSON.stringify(v))
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

function enrichFavoriteFields(row: any) {
  const id = Number(row?.id || 0)
  const count = productFavoriteCountMap.get(id) ?? Number(row?.favorite_count || 0)
  return {
    ...row,
    favorite_count: Math.max(0, count),
    is_favorited: favoriteAtMap.has(id),
  }
}

function findProductByID(id: number) {
  if (id <= 0) return null
  const hit = productListSource.find((item: any) => Number(item?.id || 0) === id)
  if (hit) return hit
  if (Number(productDetailSource?.id || 0) === id) return productDetailSource
  return null
}

function pickDetailTheme(id: number) {
  const themes = [
    {
      specName: '规格',
      values: ['标准装', '升级装', '礼盒装'],
      highlights: ['官方正品', '极速发货', '7天无忧退换'],
      scenes: ['居家使用', '办公通勤', '礼赠场景'],
    },
    {
      specName: '颜色',
      values: ['曜石黑', '珍珠白', '星雾灰'],
      highlights: ['人气热销', '新品尝鲜', '质感升级'],
      scenes: ['日常搭配', '商务出行', '旅行便携'],
    },
    {
      specName: '版本',
      values: ['经典款', '高配款', '旗舰款'],
      highlights: ['性能稳定', '配置均衡', '体验进阶'],
      scenes: ['入门首选', '家庭常备', '进阶用户'],
    },
  ]
  return themes[Math.abs(id) % themes.length]
}

function normalizeImageURL(value: any) {
  if (!value) return ''
  if (typeof value === 'string') return value
  if (typeof value?.url === 'string') return value.url
  if (typeof value?.src === 'string') return value.src
  return ''
}

function collectProductImagePool(product: any) {
  const id = Number(product?.id || 0)
  const urls: string[] = []
  const pushURL = (value: any) => {
    const url = normalizeImageURL(value)
    if (url && !urls.includes(url)) urls.push(url)
  }

  pushURL(product?.cover)

  if (Array.isArray(product?.images)) {
    for (const item of product.images) pushURL(item)
  }

  const detailBlocks = Array.isArray(product?.detail?.blocks) ? product.detail.blocks : []
  for (const block of detailBlocks) {
    if (block?.type === 'image') pushURL(block?.props?.url)
  }

  if (!urls.length && Number(productDetailSource?.id || 0) === id) {
    pushURL(productDetailSource?.cover)
    for (const item of productDetailSource?.images || []) pushURL(item)
    for (const block of productDetailSource?.detail?.blocks || []) {
      if (block?.type === 'image') pushURL(block?.props?.url)
    }
  }

  return urls
}

function buildDetailImages(product: any) {
  const id = Number(product?.id || 0)
  const pool = collectProductImagePool(product)
  const fallback = normalizeImageURL(product?.cover) || normalizeImageURL(productDetailSource?.cover)
  const seeds = pool.length ? pool : [fallback].filter(Boolean)
  const unique = Array.from(new Set(seeds.filter(Boolean)))

  while (unique.length < 3) {
    unique.push(unique[unique.length - 1] || unique[0] || 'https://picsum.photos/750/750?random=20')
  }

  return unique.slice(0, Math.max(3, unique.length)).map((url, index) => ({
    id: id * 10 + index + 1,
    product_id: id,
    url,
    sort: index,
  }))
}

function buildDetailSkus(product: any) {
  const id = Number(product?.id || 0)
  const basePrice = Number(product?.price || 0)
  const baseStock = Math.max(10, Number(product?.stock || 0))
  const theme = pickDetailTheme(id)
  return theme.values.map((value, index) => ({
    id: id * 100 + index + 1,
    product_id: id,
    attrs: JSON.stringify([{ name: theme.specName, value }]),
    price: Number((basePrice + (index - 1) * Math.max(basePrice * 0.08, 3)).toFixed(2)),
    stock: Math.max(1, baseStock - index * Math.max(Math.floor(baseStock / 5), 5)),
    sku_code: `DEMO-${id}-${index + 1}`,
  }))
}

function buildDetailBlocks(product: any) {
  const id = Number(product?.id || 0)
  const title = String(product?.title || '演示商品')
  const subtitle = String(product?.subtitle || '')
  const theme = pickDetailTheme(id)
  const images = buildDetailImages(product)
  return [
    {
      id: `text-${id}-1`,
      type: 'text',
      props: {
        text: `${title}${subtitle ? `，${subtitle}` : ''}。${theme.highlights.join('，')}。`,
      },
    },
    {
      id: `image-${id}-1`,
      type: 'image',
      props: {
        url: images[1]?.url || images[0]?.url || '',
      },
    },
    {
      id: `text-${id}-2`,
      type: 'text',
      props: {
        text: `适用场景：${theme.scenes.join('、')}。演示数据已补齐商品轮播、规格库存、详情图文与评价链路。`,
      },
    },
    {
      id: `image-${id}-2`,
      type: 'image',
      props: {
        url: images[2]?.url || images[0]?.url || '',
      },
    },
  ]
}

function buildFallbackReviewSeed(product: any, index: number) {
  const id = Number(product?.id || 0)
  const now = new Date(Date.now() - index * 86400000).toISOString()
  const users = ['演示用户A', '演示用户B', '回头客', '精选买家']
  const comments = [
    '包装完整，实物和页面展示一致，体验比较稳定。',
    '价格和品质匹配，日常使用没有问题，会考虑回购。',
    '发货很快，细节做得不错，整体满意。',
    '作为演示商品数据，规格和详情信息已经比较完整。',
  ]
  const imagePool = buildDetailImages(product).map((item: any) => item.url)
  return {
    id: id * 1000 + index + 1,
    order_id: 900000 + id * 10 + index,
    order_no: `DEMO${id}${index + 1}`,
    order_item_id: id * 100 + index + 1,
    product_id: id,
    product_score: 5 - (index % 2),
    logistics_score: 5,
    content: comments[index % comments.length],
    images: index === 0 ? imagePool.slice(0, 2) : [],
    edited_times: 0,
    user_id: index + 1,
    user_nickname: users[index % users.length],
    user_avatar: '',
    created_at: now,
    updated_at: now,
    appends: index === 0 ? [{
      id: id * 10000 + index + 1,
      review_id: id * 1000 + index + 1,
      content: '补充体验：连续使用一段时间后，表现依然稳定。',
      images: imagePool.slice(2, 3),
      created_at: now,
      updated_at: now,
    }] : [],
    admin_reply: index === 1 ? {
      id: id * 20000 + index + 1,
      review_id: id * 1000 + index + 1,
      content: '感谢反馈，演示环境会持续补充更完整的数据体验。',
      created_at: now,
      updated_at: now,
      admin_id: 1,
    } : null,
  }
}

function ensureProductReviews(product: any) {
  const id = Number(product?.id || 0)
  const bucket = ensureProductReviewBucket(id)
  if (bucket.length) return bucket
  const list = [0, 1, 2].map((index) => buildFallbackReviewSeed(product, index))
  productReviewMap.set(id, list)
  for (const item of list) {
    reviewSeq = Math.max(reviewSeq, Number(item.id || 0))
    reviewIndexMap.set(Number(item.id), item)
    for (const ap of item.appends || []) appendSeq = Math.max(appendSeq, Number(ap.id || 0))
    if (item.admin_reply) replySeq = Math.max(replySeq, Number(item.admin_reply.id || 0))
  }
  return list
}

function buildProductDetailPayload(source: any) {
  const base = clone(source || {})
  const id = Number(base?.id || 0)
  if (Number(productDetailSource?.id || 0) === id) {
    const detail = clone(productDetailSource)
    const existingBlocks = Array.isArray(detail?.detail?.blocks) ? detail.detail.blocks : []
    const merged = {
      ...base,
      ...detail,
      subtitle: detail.subtitle || base.subtitle || '',
      sales: Number(detail.sales || base.sales || 0),
      stock: Number(detail.stock || base.stock || 0),
      detail: {
        version: 1,
        blocks: existingBlocks.length ? existingBlocks : buildDetailBlocks({ ...base, ...detail }),
      },
    }
    merged.images = buildDetailImages(merged)
    ensureProductReviews(merged)
    return enrichFavoriteFields(merged)
  }

  const images = buildDetailImages(base)
  const skus = buildDetailSkus(base)
  const merged = {
    ...base,
    subtitle: String(base.subtitle || `${base.title || '商品'} 演示详情已补齐`),
    sales: Number(base.sales || 0),
    stock: Number(base.stock || 0),
    images,
    skus,
    detail: {
      version: 1,
      blocks: buildDetailBlocks(base),
    },
  }
  ensureProductReviews(merged)
  return enrichFavoriteFields(merged)
}

function favoriteProductByID(id: number) {
  if (id <= 0) return
  if (favoriteAtMap.has(id)) return
  favoriteAtMap.set(id, new Date().toISOString())
  const prev = productFavoriteCountMap.get(id) ?? Number(findProductByID(id)?.favorite_count || 0)
  productFavoriteCountMap.set(id, prev + 1)
}

function unfavoriteProductByID(id: number) {
  if (id <= 0) return
  if (!favoriteAtMap.has(id)) return
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

for (const order of orderListSource) ensureOrderExt(order)

function getAfterSaleByID(id: number) {
  return afterSalesSource.find((row: any) => Number(row.id) === id)
}

function touchOrderAfterSaleSummary(orderID: number) {
  const order = getOrderByID(orderID)
  if (!order) return
  order.after_sale_summary = buildAfterSaleSummary(orderID)
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

buildProductIndex()
hydrateReviewsFromOrders()
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
  'GET /api/v1/products': { list: productListSource, total: productListSource.length, page: 1, size: 20 },
  'GET /api/v1/products/recommend': recommend,
  'GET /api/v1/cart': cart,
  'GET /api/v1/orders': orders,
  'GET /api/v1/orders/': orderListSource[0] || null,
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
  'GET /api/v1/seckill/products': { list: buildActivityRows('seckill'), total: buildActivityRows('seckill').length, page: 1, size: 20 },
  'GET /api/v1/group-buy/products': { list: buildActivityRows('group_buy'), total: buildActivityRows('group_buy').length, page: 1, size: 20 },
  'GET /api/v1/bargain/products': { list: buildActivityRows('bargain'), total: buildActivityRows('bargain').length, page: 1, size: 20 },
  'GET /api/v1/marketing/seckill/products': { list: buildActivityRows('seckill'), total: buildActivityRows('seckill').length, page: 1, size: 20 },
  'GET /api/v1/marketing/group-buy/products': { list: buildActivityRows('group_buy'), total: buildActivityRows('group_buy').length, page: 1, size: 20 },
  'GET /api/v1/marketing/bargain/products': { list: buildActivityRows('bargain'), total: buildActivityRows('bargain').length, page: 1, size: 20 },
  'POST /api/v1/marketing/group-buy/join': { group_order_id: 1, status: 1 },
  'POST /api/v1/marketing/bargain/help': { cut_amount: 58.50, current_price: 1200 },
  'POST /api/v1/cart/add': null,
  'POST /api/v1/orders': { order_no: 'DEMO202600001', id: 1, status: 1 },
  'POST /api/v1/auth/sms/send': { dev_code: '123456' },
  'POST /api/v1/auth/sms/login': { token: 'demo_token_mock' },
  'GET /api/v1/checkin/status': { checked_today: false, consecutive_days: 3, month_dates: ['2026-05-17', '2026-05-20', '2026-05-21'], month_count: 3, month_points: 40 },
  'GET /api/v1/checkin/rules': [{ day: 0, points: 10 }, { day: 3, points: 20 }, { day: 7, points: 50 }],
  'POST /api/v1/checkin': { points: 20, consecutive_days: 4 },
  'GET /api/v1/im/session': { id: 1, user_id: 1, staff_id: 0, mode: 'ai', status: 2, queue_position: 0, last_msg: '', unread_count: 0 },
  'GET /api/v1/im/messages': { list: [
    { id: 1, session_id: 1, sender_type: 0, sender_id: 0, type: 'system', content: 'AI 智能助手已为您服务，如需人工客服请输入「人工」或点击转人工按钮', is_read: 1, created_at: '2026-05-28T10:00:00Z' },
    { id: 2, session_id: 1, sender_type: 1, sender_id: 1, type: 'text', content: '你们退货政策是什么？', is_read: 1, created_at: '2026-05-28T10:01:00Z' },
    { id: 3, session_id: 1, sender_type: 3, sender_id: 0, type: 'text', content: '您好！根据知识库资料，本商城支持7天无理由退货，商品须保持完好原状。退货运费由买家承担，商品质量问题由商家承担，退款1-3个工作日到账。', is_read: 1, created_at: '2026-05-28T10:01:00Z' },
    { id: 4, session_id: 1, sender_type: 1, sender_id: 1, type: 'text', content: '一般几天发货？', is_read: 1, created_at: '2026-05-28T10:03:00Z' },
    { id: 5, session_id: 1, sender_type: 3, sender_id: 0, type: 'text', content: '付款成功后通常1-2个工作日内发货，节假日会顺延。发货后会推送物流信息，可在订单详情页实时查看物流状态。', is_read: 1, created_at: '2026-05-28T10:03:00Z' },
  ], total: 5 },
  'POST /api/v1/im/feedback': { success: true },
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
  if (upperMethod === 'GET' && path === '/api/v1/coupons') {
    return { matched: true, data: claimableCouponSource.map((item) => ({ ...item })) }
  }
  if (upperMethod === 'POST' && /^\/api\/v1\/coupons\/\d+\/claim$/.test(path)) {
    const couponID = Number(path.split('/')[4] || 0)
    const target = claimableCouponSource.find((item) => Number(item.id) === couponID)
    if (target && target.can_claim) {
      target.claimed_count = Number(target.claimed_count || 0) + 1
      target.claimed_by_me = Number(target.claimed_by_me || 0) + 1
      const limitReached = target.per_limit > 0 && Number(target.claimed_by_me || 0) >= target.per_limit
      const soldOut = target.total_count > 0 && Number(target.claimed_count || 0) >= target.total_count
      target.can_claim = !(limitReached || soldOut)
      const nextID = Math.max(0, ...userCouponSource.map((item: any) => Number(item.id || 0))) + 1
      userCouponSource.unshift({
        id: nextID,
        coupon_id: couponID,
        user_id: 1,
        status: 1,
        used_at: null,
        order_id: 0,
        created_at: new Date().toISOString(),
        coupon: {
          id: target.id,
          name: target.name,
          type: target.type,
          min_amount: target.min_amount,
          discount: target.discount,
          end_at: null,
        },
      })
    }
    return { matched: true, data: null }
  }
  if (upperMethod === 'GET' && path === '/api/v1/user/coupons') {
    return { matched: true, data: userCouponSource.map((item: any) => ({ ...item })) }
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

  if (upperMethod === 'POST' && path === '/api/v1/orders') {
    const body = params || {}
    const fromItems = Array.isArray(body.items) ? body.items : []
    const fromSkuIDs = Array.isArray(body.sku_ids)
      ? body.sku_ids.map((v: any) => Number(v || 0)).filter((v: number) => v > 0)
      : String(body.sku_ids || '').split(',').map((v) => Number(v || 0)).filter((v) => v > 0)
    const requestItems = fromItems.length
      ? fromItems
          .map((item: any) => ({
            sku_id: Number(item?.sku_id || 0),
            activity_product_id: Number(item?.activity_product_id || 0),
          }))
          .filter((item: any) => item.sku_id > 0)
      : fromSkuIDs.map((skuID: number) => ({ sku_id: skuID, activity_product_id: 0 }))

    if (!requestItems.length) {
      return { matched: true, data: null }
    }

    const now = new Date().toISOString()
    const nextID = Math.max(0, ...orderListSource.map((item: any) => Number(item.id || 0))) + 1
    const orderNo = `DEMO${Date.now()}${String(nextID).padStart(3, '0')}`
    const orderItems = requestItems.map((item: any, idx: number) => {
      const source = productListSource.find((p: any) => Number(p?.id || 0) === Number(item.sku_id || 0)) || productListSource[idx] || {}
      return {
        id: nextID * 100 + idx + 1,
        order_id: nextID,
        product_id: Number(source?.id || item.sku_id || 0),
        sku_id: Number(item.sku_id || 0),
        activity_product_id: Number(item.activity_product_id || 0),
        title: String(source?.title || `商品${item.sku_id}`),
        cover: String(source?.cover || ''),
        attrs: '[]',
        price: Number(source?.price || 0),
        qty: 1,
      }
    })
    const goodsAmount = Number(orderItems.reduce((sum, item) => sum + Number(item.price || 0) * Number(item.qty || 0), 0).toFixed(2))
    const newOrder: any = {
      id: nextID,
      user_id: 1,
      order_no: orderNo,
      status: 1,
      payment_method: String(body.payment_method || 'wechat'),
      goods_amount: goodsAmount,
      discount_amount: 0,
      freight_amount: 0,
      total_amount: goodsAmount,
      created_at: now,
      updated_at: now,
      paid_at: null,
      tracking_no: '',
      items: orderItems,
      remark: String(body.remark || ''),
      activity_type: Number(orderItems[0]?.activity_product_id || 0) > 0 ? 'activity' : '',
    }
    ensureOrderExt(newOrder)
    orderListSource.unshift(newOrder)
    return { matched: true, data: { order_no: orderNo, id: nextID, status: 1 } }
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
  if (upperMethod === 'POST' && path.startsWith('/api/v1/orders/') && path.endsWith('/cancel')) {
    const id = Number(path.split('/')[4] || 0)
    const target = getOrderByID(id)
    if (target && Number(target.status) === 1) {
      target.status = 6
      target.updated_at = new Date().toISOString()
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
    return { matched: true, data: JSON.parse(JSON.stringify(detail)) }
  }
  if (upperMethod === 'GET' && /^\/api\/v1\/products\/\d+$/.test(path)) {
    const id = Number(path.split('/').pop() || 0)
    const source = findProductByID(id)
    if (!source) return { matched: true, data: null }
    return { matched: true, data: buildProductDetailPayload(source) }
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
      if (keyword) {
        list = list.filter((item: any) => String(item.title || '').toLowerCase().includes(keyword))
      }
      if (categoryID > 0) {
        list = list.filter((item: any) => Number(item.category_id) === categoryID)
      }
      const offset = Math.max(page - 1, 0) * Math.max(size, 1)
      const pageList = list.slice(offset, offset + size).map((item: any) => enrichFavoriteFields(clone(item)))
      return { matched: true, data: { list: pageList, total: list.length, page, size } }
    }
    return { matched: true, data }
  }

  // Prefix match (for routes with path params like /products/:id)
  for (const pattern of Object.keys(routes)) {
    if (key.startsWith(pattern) && pattern.endsWith('/')) {
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
