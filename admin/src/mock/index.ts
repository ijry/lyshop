// Admin mock data — routes keyed by "METHOD /admin/api/..."
// Reuses product/order JSON from app mock where applicable

import products from '../../../app/mock/data/products.json'
import categories from '../../../app/mock/data/categories.json'
import orders from '../../../app/mock/data/orders.json'
import productDetail from '../../../app/mock/data/product-detail.json'
import { afterSaleStatusLabel, shipmentStatusLabel } from '../utils/order-status'

const orderListSource = (orders as any)?.list || []
const shipmentDirectionLabels: Record<string, string> = {
  outbound: '寄出',
  inbound: '回寄',
}
const shipmentBizTypeLabels: Record<string, string> = {
  initial: '首发',
  reship: '补发',
  return: '回寄',
}
const toNumber = (v: any) => Number(v || 0)
let replySeq = 20000
let afterSaleSeq = 8000
const afterSalesSource: any[] = []
const categoriesSource: any[] = clone(Array.isArray(categories) ? categories : [])
const productListSource: any[] = clone(Array.isArray((products as any)?.list) ? (products as any).list : [])
let categorySeq = Math.max(0, ...categoriesSource.map((item: any) => Number(item?.id || 0)))
const decorVariantsSource: any[] = [
  {
    id: 1,
    merchant_id: 0,
    page_key: 'index',
    variant_key: 'default',
    variant_name: '默认副本',
    components: JSON.stringify([
      { type: 'banner', id: 'c1', props: { images: [{ url: 'https://picsum.photos/750/350?random=1' }], height: 350 } },
      { type: 'product_grid', id: 'c2', props: { source: 'hot', limit: 6, columns: 2 } },
    ]),
    is_published: true,
    published_at: '2026-05-25T10:00:00Z',
  },
]
const pcDecorSource: any = {
  id: 1,
  components: JSON.stringify({
    pageStyle: {
      background: {
        mode: 'solid',
        solidColor: '#f8fafc',
        gradient: {
          angle: 135,
          stops: [{ color: '#f8fafc', position: 0 }, { color: '#eef2ff', position: 100 }],
        },
        image: {
          url: '',
          size: 'cover',
          customSize: '100% auto',
          position: 'center top',
          repeat: 'no-repeat',
          attachment: 'scroll',
        },
        overlay: { enabled: false, color: '#000000', opacity: 0.2 },
      },
      content: { maxWidth: 1280, gutterX: 24, sectionGap: 24 },
      surface: { radius: 12, shadow: 'none' },
    },
    components: [
      { type: 'hero', id: 'pc_hero', props: { badge: '限时秒杀进行中', title: '精选好物\n品质生活从这里开始', subtitle: '数千款精选商品，正品保障，极速发货。', btn_text: '立即选购', btn_link: '/products', btn2_text: '查看全部', btn2_link: '/products', bg_from: '#b91c1c', bg_to: '#991b1b' }, style: {} },
      { type: 'product_grid', id: 'pc_hot', props: { title: '热销推荐', source: 'hot', limit: 8, columns: 4 }, style: {} },
      { type: 'features', id: 'pc_features', props: { columns: 4, items: [{ icon: 'i-carbon-delivery-truck', title: '快递配送', desc: '全国包邮' }, { icon: 'i-carbon-checkmark-outline', title: '正品保障', desc: '假一赔十' }, { icon: 'i-carbon-renew', title: '无忧退换', desc: '7天退换' }, { icon: 'i-carbon-headset', title: '在线客服', desc: '随时响应' }] }, style: {} },
    ],
  }),
  is_published: true,
}

const vipPlansSource: any[] = [
  { id: 1, name: '月卡', months: 1, price: 19.9, status: 1 },
  { id: 2, name: '季卡', months: 3, price: 49.9, status: 1 },
]
const vipLevelsSource: any[] = [
  { id: 1, name: '银卡', growth_min: 0, discount_rate: 0.98 },
  { id: 2, name: '金卡', growth_min: 1000, discount_rate: 0.95 },
]
const vipCouponRulesSource: any[] = [
  { id: 1, name: '银卡月券', coupon_name: '满99减10', monthly_limit: 1 },
  { id: 2, name: '金卡月券', coupon_name: '满199减30', monthly_limit: 2 },
]
const vipSkuPricesSource: any[] = [
  { id: 1, product_id: 1, sku_id: 1, level_id: 1, level_name: '银卡', vip_price: 4888, status: 1 },
  { id: 2, product_id: 1, sku_id: 1, level_id: 2, level_name: '金卡', vip_price: 4699, status: 1 },
  { id: 3, product_id: 1, sku_id: 2, level_id: 2, level_name: '金卡', vip_price: 5399, status: 1 },
]
const seckillActivitiesSource: any[] = [
  { id: 1, type: 'seckill', name: '午场秒杀', start_at: '2026-05-27T10:00:00Z', end_at: '2026-05-27T12:00:00Z', status: 1 },
  { id: 2, type: 'seckill', name: '晚场秒杀', start_at: '2026-05-27T20:00:00Z', end_at: '2026-05-27T22:00:00Z', status: 1 },
]
const groupBuyActivitiesSource: any[] = [
  { id: 11, type: 'group_buy', name: '今日拼团', start_at: '2026-05-27T00:00:00Z', end_at: '2026-05-27T23:59:59Z', status: 1 },
]
const bargainActivitiesSource: any[] = [
  { id: 21, type: 'bargain', name: '砍价专场', start_at: '2026-05-27T00:00:00Z', end_at: '2026-05-28T23:59:59Z', status: 1 },
]
const seckillProductsSource: any[] = [
  { id: 1001, activity_id: 1, product_id: 1, sku_id: 10011, activity_price: 3599, start_price: 0, floor_price: 0, limit_per_order: 1, total_stock_limit: 120, sold_qty: 18, product_title: '旗舰智能手机 Pro Max', product_cover: 'https://picsum.photos/120/120?random=101', sku_price: 4999, sku_stock: 400, sku_attrs: [{ name: '颜色', value: '曜石黑' }] },
]
const groupBuyProductsSource: any[] = [
  { id: 2001, activity_id: 11, product_id: 4, sku_id: 10041, activity_price: 2699, start_price: 0, floor_price: 0, limit_per_order: 2, total_stock_limit: 200, sold_qty: 29, product_title: '智能手表 Ultra', product_cover: 'https://picsum.photos/120/120?random=102', sku_price: 3299, sku_stock: 380, sku_attrs: [{ name: '版本', value: '标准版' }] },
]
const bargainProductsSource: any[] = [
  { id: 3001, activity_id: 21, product_id: 6, sku_id: 10061, activity_price: 0, start_price: 1980, floor_price: 299, limit_per_order: 1, total_stock_limit: 80, sold_qty: 6, product_title: '羊绒大衣女款', product_cover: 'https://picsum.photos/120/120?random=103', sku_price: 1980, sku_stock: 220, sku_attrs: [{ name: '尺码', value: 'M' }] },
]
let marketingActivitySeq = 100
let marketingProductSeq = 5000

function clone<T>(v: T): T {
  return JSON.parse(JSON.stringify(v))
}

for (const item of productListSource) {
  if (item.favorite_count === undefined || item.favorite_count === null) {
    item.favorite_count = Math.max(0, Math.floor(Number(item.sales || 0) / 8))
  }
}

function withAfterSaleLabels(row: any) {
  const data = clone(row || {})
  const logs = Array.isArray(data.logs)
    ? data.logs.map((log: any) => ({
      ...log,
      from_status_label: String(log?.from_status || '').trim() ? afterSaleStatusLabel(log.from_status) : '',
      to_status_label: afterSaleStatusLabel(log?.to_status),
    }))
    : []
  const shipments = Array.isArray(data.shipments)
    ? data.shipments.map((ship: any) => ({
      ...ship,
      logistics_status_label: shipmentStatusLabel(ship?.logistics_status),
    }))
    : []
  return {
    ...data,
    status_label: afterSaleStatusLabel(data?.status),
    logs,
    shipments,
  }
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

function touchOrderAfterSaleSummary(orderID: number) {
  const target = orderListSource.find((item: any) => Number(item.id) === Number(orderID))
  if (!target) return
  target.after_sale_summary = buildAfterSaleSummary(orderID)
}

for (const order of orderListSource) ensureOrderExt(order)

function buildAfterSaleRows() {
  const rows: any[] = []
  const target = orderListSource.find((item: any) => Number(item.id) === 2) || orderListSource[0]
  if (!target) return rows
  const caseID = ++afterSaleSeq
  rows.push({
    id: caseID,
    order_id: Number(target.id || 0),
    user_id: Number(target.user_id || 1),
    case_no: nextAfterSaleCaseNo(),
    case_type: 'return',
    status: 'approved_wait_user_return',
    reason: '尺寸不合适',
    apply_content: '试穿后不合适',
    apply_images: ['https://picsum.photos/200/200?random=991'],
    refund_amount: 0,
    logs: [
      {
        id: Math.floor(Math.random() * 100000),
        case_id: caseID,
        from_status: '',
        to_status: 'applied',
        action: 'apply',
        operator_type: 'user',
        operator_id: Number(target.user_id || 1),
        content: '提交售后申请',
        created_at: new Date().toISOString(),
      },
      {
        id: Math.floor(Math.random() * 100000),
        case_id: caseID,
        from_status: 'applied',
        to_status: 'approved_wait_user_return',
        action: 'audit',
        operator_type: 'admin',
        operator_id: 1,
        content: '售后审核通过',
        created_at: new Date().toISOString(),
      },
    ],
    items: (target.items || []).slice(0, 1).map((item: any) => ({
      id: Math.floor(Math.random() * 100000),
      case_id: caseID,
      order_item_id: Number(item.id || 0),
      qty: 1,
    })),
    shipments: [],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  })
  return rows
}

afterSalesSource.push(...buildAfterSaleRows())
for (const order of orderListSource) touchOrderAfterSaleSummary(Number(order.id))

function buildReviewRows() {
  const rows: any[] = []
  for (const order of orderListSource) {
    const items = Array.isArray(order?.items) ? order.items : []
    for (const item of items) {
      const rv = item?.review
      if (!rv || Number(rv.id || 0) <= 0) continue
      rows.push({
        id: Number(rv.id),
        order_id: Number(order.id || 0),
        order_no: String(order.order_no || ''),
        order_item_id: Number(item.id || 0),
        product_id: Number(item.product_id || 0),
        product_score: Number(rv.product_score || 5),
        logistics_score: Number(rv.logistics_score || 5),
        content: String(rv.content || ''),
        images: Array.isArray(rv.images) ? rv.images.map((u: any) => String(u || '')) : [],
        edited_times: Number(rv.edited_times || 0),
        user_id: Number(order.user_id || 1),
        user_nickname: '演示用户',
        user_avatar: '',
        created_at: rv.created_at || order.created_at,
        updated_at: rv.updated_at || rv.created_at || order.created_at,
        appends: Array.isArray(rv.appends) ? rv.appends.map((ap: any) => ({
          id: Number(ap.id || 0),
          review_id: Number(rv.id),
          content: String(ap.content || ''),
          images: Array.isArray(ap.images) ? ap.images.map((u: any) => String(u || '')) : [],
          created_at: ap.created_at || order.created_at,
          updated_at: ap.updated_at || ap.created_at || order.created_at,
        })) : [],
        admin_reply: rv.admin_reply ? {
          id: Number(rv.admin_reply.id || ++replySeq),
          review_id: Number(rv.id),
          admin_id: Number(rv.admin_reply.admin_id || 1),
          content: String(rv.admin_reply.content || ''),
          created_at: rv.admin_reply.created_at || order.created_at,
          updated_at: rv.admin_reply.updated_at || rv.admin_reply.created_at || order.created_at,
        } : null,
        order_item: {
          id: Number(item.id || 0),
          order_id: Number(order.id || 0),
          product_id: Number(item.product_id || 0),
          title: String(item.title || ''),
          cover: String(item.cover || ''),
          qty: Number(item.qty || 0),
          price: Number(item.price || 0),
        },
        product: {
          id: Number(item.product_id || 0),
          title: String(item.title || ''),
          cover: String(item.cover || ''),
        },
      })
    }
  }
  rows.sort((a, b) => Number(b.id) - Number(a.id))
  return rows
}

const reviewRows = buildReviewRows()

function getReviewByID(id: number) {
  return reviewRows.find((row: any) => Number(row.id) === id)
}

function listReviews(params?: Record<string, any>) {
  const query = params || {}
  const productID = toNumber(query.product_id)
  const keyword = String(query.keyword || '').trim().toLowerCase()
  const page = Math.max(1, toNumber(query.page) || 1)
  const size = Math.max(1, toNumber(query.size) || 20)
  let list = reviewRows.slice()
  if (productID > 0) list = list.filter((item: any) => toNumber(item.product_id) === productID)
  if (keyword) list = list.filter((item: any) => String(item.content || '').toLowerCase().includes(keyword))
  const offset = (page - 1) * size
  const pageList = list.slice(offset, offset + size).map((item: any) => clone(item))
  return { list: pageList, total: list.length, page, size }
}

function upsertReviewReply(id: number, content: string) {
  const target = getReviewByID(id)
  if (!target) return
  const now = new Date().toISOString()
  if (!target.admin_reply) {
    replySeq += 1
    target.admin_reply = {
      id: replySeq,
      review_id: id,
      admin_id: 1,
      content,
      created_at: now,
      updated_at: now,
    }
  } else {
    target.admin_reply.content = content
    target.admin_reply.updated_at = now
  }
  for (const order of orderListSource) {
    const items = Array.isArray(order?.items) ? order.items : []
    const item = items.find((row: any) => toNumber(row?.review?.id) === id)
    if (!item?.review) continue
    item.review.admin_reply = {
      id: target.admin_reply.id,
      content,
      created_at: target.admin_reply.created_at,
      updated_at: target.admin_reply.updated_at,
      admin_id: 1,
    }
    item.review.updated_at = now
  }
}

const routes: Record<string, any> = {
  // Auth
  'POST /admin/api/auth/login': { token: 'demo_admin_token' },

  // Menus (首页提升为一级 TAB，其余按分组返回)
  'GET /admin/api/menus': {
    dashboard: { title: '首页', path: '/dashboard' },
    groups: [
      {
        key: 'product',
        title: '商品',
        icon: 'box',
        sort: 10,
        menus: [
          { title: '商品管理', icon: 'box', path: '/product', sort: 10, children: [
            { title: '商品列表', path: '/product/list' },
            { title: '商品分类', path: '/product/category' },
            { title: '新增商品', path: '/product/form' },
          ]},
          { title: '评价管理', icon: 'star', path: '/review', sort: 20, children: [
            { title: '评价列表', path: '/review/list' },
          ]},
          { title: '店铺装修', icon: 'layout', path: '/decor', sort: 30, children: [
            { title: '移动端装修', path: '/decor/index' },
            { title: 'PC端装修', path: '/decor/pc' },
          ]},
        ],
      },
      {
        key: 'order',
        title: '订单',
        icon: 'shopping-cart',
        sort: 20,
        menus: [
          { title: '订单管理', icon: 'shopping-cart', path: '/order', sort: 10, children: [
            { title: '订单列表', path: '/order/list' },
            { title: '售后列表', path: '/order/after-sale/list' },
          ]},
        ],
      },
      {
        key: 'user',
        title: '用户',
        icon: 'users',
        sort: 30,
        menus: [
          { title: '客服中心', icon: 'message-circle', path: '/im', sort: 10, children: [
            { title: '客服会话', path: '/im/sessions' },
          ]},
          { title: '消息管理', icon: 'bell', path: '/message', sort: 20, children: [
            { title: '消息列表', path: '/message/list' },
            { title: '发送消息', path: '/message/send' },
          ]},
          { title: '签到管理', icon: 'calendar', path: '/checkin', sort: 30, children: [
            { title: '签到规则', path: '/checkin/rules' },
            { title: '签到记录', path: '/checkin/logs' },
          ]},
        ],
      },
      {
        key: 'marketing',
        title: '营销',
        icon: 'tag',
        sort: 40,
        menus: [
          { title: '营销管理', icon: 'tag', path: '/marketing', sort: 10, children: [
            { title: '优惠券管理', path: '/marketing/coupon' },
            { title: '秒杀活动管理', path: '/marketing/seckill/activity' },
            { title: '秒杀商品管理', path: '/marketing/seckill/product' },
            { title: '拼团活动管理', path: '/marketing/group-buy/activity' },
            { title: '拼团商品管理', path: '/marketing/group-buy/product' },
            { title: '砍价活动管理', path: '/marketing/bargain/activity' },
            { title: '砍价商品管理', path: '/marketing/bargain/product' },
          ]},
        ],
      },
      {
        key: 'vip',
        title: '会员',
        icon: 'award',
        sort: 45,
        menus: [
          { title: '会员中心', icon: 'award', path: '/vip', sort: 10, children: [
            { title: '会员套餐', path: '/vip/plans' },
            { title: '会员等级', path: '/vip/levels' },
            { title: '会员领券规则', path: '/vip/coupon-rules' },
          ]},
        ],
      },
      {
        key: 'wms',
        title: '仓储',
        icon: 'warehouse',
        sort: 50,
        menus: [
          { title: '仓储管理', icon: 'warehouse', path: '/wms', sort: 10, children: [
            { title: '库存管理', path: '/wms/stock' },
          ]},
        ],
      },
      {
        key: 'system',
        title: '系统',
        icon: 'settings',
        sort: 60,
        menus: [
          { title: '系统设置', icon: 'settings', path: '/system', sort: 10, children: [
            { title: '站点设置', path: '/system/site' },
            { title: '配置中心', path: '/system/config' },
            { title: '管理员管理', path: '/system/admins' },
            { title: '角色管理', path: '/system/roles' },
          ]},
        ],
      },
    ],
  },

  // Dashboard
  'GET /admin/api/dashboard': {
    today_orders: 56,
    today_sales: 28960.50,
    pending_refunds: 3,
    online_sessions: 2,
    sales_trend: [
      { date: '2026-05-19', orders: 42, sales: 18660.00 },
      { date: '2026-05-20', orders: 38, sales: 17280.00 },
      { date: '2026-05-21', orders: 47, sales: 21340.00 },
      { date: '2026-05-22', orders: 52, sales: 24120.50 },
      { date: '2026-05-23', orders: 49, sales: 22590.00 },
      { date: '2026-05-24', orders: 58, sales: 27630.00 },
      { date: '2026-05-25', orders: 56, sales: 28960.50 },
    ],
  },

  // Products
  'GET /admin/api/products': { list: productListSource, total: productListSource.length, page: 1, size: 20 },
  'GET /admin/api/products/': { ...productDetail, favorite_count: Number((productDetail as any)?.favorite_count || 0) },
  'POST /admin/api/products': { id: 100 },

  // Orders
  'GET /admin/api/orders': orders,
  'GET /admin/api/orders/': orderListSource[0] || null,
  'GET /admin/api/reviews': { list: reviewRows.slice(), total: reviewRows.length, page: 1, size: 20 },
  'GET /admin/api/reviews/': reviewRows[0] || null,
  'POST /admin/api/reviews/': null,

  // WMS
  'GET /admin/api/wms/warehouses': [
    { id: 1, name: '主仓库', address: '上海市浦东新区', contact: '张三', phone: '13800001111', status: 1 },
    { id: 2, name: '华南仓', address: '广州市天河区', contact: '李四', phone: '13800002222', status: 1 },
  ],
  'GET /admin/api/wms/stocks': {
    list: [
      { id: 1, warehouse_id: 1, sku_id: 1, qty: 200, safe_qty: 50 },
      { id: 2, warehouse_id: 1, sku_id: 2, qty: 100, safe_qty: 30 },
      { id: 3, warehouse_id: 2, sku_id: 3, qty: 3000, safe_qty: 100 },
    ],
    total: 3, page: 1, size: 20,
  },

  // Marketing
  'GET /admin/api/marketing/coupons': {
    list: [
      { id: 1, name: '新人满100减20', type: 1, min_amount: 100, discount: 20, total_count: 1000, per_limit: 1, status: 1 },
      { id: 2, name: '全场9折券', type: 2, min_amount: 0, discount: 0.9, total_count: 500, per_limit: 1, status: 1 },
      { id: 3, name: '无门槛5元券', type: 3, min_amount: 0, discount: 5, total_count: 0, per_limit: 3, status: 1 },
    ],
    total: 3, page: 1, size: 20,
  },
  'GET /admin/api/marketing/seckill/activities': { list: seckillActivitiesSource, total: seckillActivitiesSource.length, page: 1, size: 20 },
  'GET /admin/api/marketing/group-buy/activities': { list: groupBuyActivitiesSource, total: groupBuyActivitiesSource.length, page: 1, size: 20 },
  'GET /admin/api/marketing/bargain/activities': { list: bargainActivitiesSource, total: bargainActivitiesSource.length, page: 1, size: 20 },
  'GET /admin/api/marketing/seckill/products': { list: seckillProductsSource, total: seckillProductsSource.length, page: 1, size: 20 },
  'GET /admin/api/marketing/group-buy/products': { list: groupBuyProductsSource, total: groupBuyProductsSource.length, page: 1, size: 20 },
  'GET /admin/api/marketing/bargain/products': { list: bargainProductsSource, total: bargainProductsSource.length, page: 1, size: 20 },
  'GET /admin/api/vip/plans': { list: vipPlansSource, total: vipPlansSource.length, page: 1, size: 20 },
  'GET /admin/api/vip/levels': { list: vipLevelsSource, total: vipLevelsSource.length, page: 1, size: 20 },
  'GET /admin/api/vip/coupon-rules': { list: vipCouponRulesSource, total: vipCouponRulesSource.length, page: 1, size: 20 },
  'GET /admin/api/vip/sku-prices': { list: vipSkuPricesSource, total: vipSkuPricesSource.length, page: 1, size: 20 },
  'GET /api/v1/vip/profile': { is_vip: true, level_name: '金卡', growth_value: 1280, expire_at: '2026-12-31' },
  'GET /api/v1/vip/coupons/monthly': [
    { rule_id: 1, name: '银卡月券', coupon_name: '满99减10', claimed: 0, monthly_limit: 1 },
    { rule_id: 2, name: '金卡月券', coupon_name: '满199减30', claimed: 1, monthly_limit: 2 },
  ],

  // IM
  'GET /admin/api/im/sessions': [
    { id: 1, user_id: 1001, staff_id: 1, status: 2, last_msg: '这款耳机降噪效果怎么样？', unread_count: 1, updated_at: '2026-05-22T10:30:00Z' },
    { id: 2, user_id: 1002, staff_id: 0, status: 1, last_msg: '你好，我想退货', unread_count: 3, updated_at: '2026-05-22T09:15:00Z' },
  ],
  'GET /admin/api/im/sessions/': { list: [
    { id: 1, session_id: 1, sender_type: 1, content: '你好，我想问一下这款耳机', type: 'text', created_at: '2026-05-22T10:28:00Z' },
    { id: 2, session_id: 1, sender_type: 2, content: '您好！请问有什么可以帮您？', type: 'text', created_at: '2026-05-22T10:29:00Z' },
    { id: 3, session_id: 1, sender_type: 1, content: '这款耳机降噪效果怎么样？', type: 'text', created_at: '2026-05-22T10:30:00Z' },
  ], total: 3 },

  // AI
  'GET /admin/api/ai/models': [
    { id: 1, name: '通义万象', driver: 'tongyi', endpoint: '', is_default: 1, status: 1, supports_ref_image: 1 },
    { id: 2, name: 'DALL-E 3', driver: 'openai', endpoint: 'https://api.openai.com', is_default: 0, status: 1, supports_ref_image: 0 },
  ],
  'GET /admin/api/ai/tasks': {
    list: [
      { id: 1, model_id: 1, scene: 'carousel', prompt: '白色简约风手提包', status: 2, result_urls: '["https://picsum.photos/750/750?random=90"]', created_at: '2026-05-22T08:00:00Z' },
    ],
    total: 1, page: 1, size: 20,
  },
  'POST /admin/api/ai/generate': { id: 99, status: 1, model_id: 1, prompt: 'demo', result_urls: null },
  'GET /admin/api/ai/tasks/': { id: 99, status: 2, model_id: 1, prompt: 'demo', result_urls: '["https://picsum.photos/750/750?random=91","https://picsum.photos/750/750?random=92"]' },

  // Upload
  'POST /admin/api/upload': { path: 'demo/mock.jpg', url: 'https://picsum.photos/400/400?random=50', size: 12345, mime: 'image/jpeg' },

  // System config
  'GET /admin/api/system/sms/config': { provider: 'aliyun', access_key: 'demo***', sign_name: 'LYShop' },

  // Config Center
  'GET /admin/api/config/schemas': [
    { plugin: 'wechat_pay', title: '微信支付', fields: [
      { key: 'app_id', label: 'AppID', type: 'text', required: true, placeholder: '微信支付 AppID' },
      { key: 'mch_id', label: '商户号', type: 'text', required: true },
      { key: 'api_key', label: 'API v3 密钥', type: 'password', required: true },
      { key: 'notify_url', label: '异步回调地址', type: 'text' },
    ]},
    { plugin: 'alipay', title: '支付宝支付', fields: [
      { key: 'app_id', label: 'AppID', type: 'text', required: true },
      { key: 'private_key', label: '应用私钥', type: 'textarea', required: true, placeholder: 'RSA2 私钥内容' },
      { key: 'public_key', label: '支付宝公钥', type: 'textarea', required: true, placeholder: '支付宝公钥内容' },
      { key: 'sandbox', label: '沙箱模式', type: 'switch' },
    ]},
    { plugin: 'sms', title: '短信插件', fields: [
      { key: 'provider', label: '短信服务商', type: 'select', required: true, options: [
        { label: '阿里云', value: 'aliyun' }, { label: '腾讯云', value: 'tencent' }
      ]},
      { key: 'access_key', label: 'AccessKey', type: 'text', required: true },
      { key: 'secret_key', label: 'SecretKey', type: 'password', required: true },
      { key: 'sign_name', label: '签名名称', type: 'text', required: true },
    ]},
    { plugin: 'wechat_auth', title: '微信登录', fields: [
      { key: 'mini_app_id', label: '小程序 AppID', type: 'text' },
      { key: 'mini_app_secret', label: '小程序 AppSecret', type: 'password' },
      { key: 'h5_app_id', label: 'H5/App AppID', type: 'text' },
      { key: 'h5_app_secret', label: 'H5/App AppSecret', type: 'password' },
    ]},
    { plugin: 'storage_local', title: '本地存储', fields: [
      { key: 'upload_dir', label: '上传目录', type: 'text', placeholder: 'uploads' },
      { key: 'base_url', label: '访问URL前缀', type: 'text', placeholder: 'http://localhost:8080/uploads' },
    ]},
    { plugin: 'storage_router', title: '存储路由', fields: [
      { key: 'default_driver', label: '默认上传驱动', type: 'select', required: true, options: [
        { label: '本地存储', value: 'local' },
        { label: '阿里云 OSS', value: 'oss' },
        { label: '腾讯云 COS', value: 'cos' },
        { label: '七牛云', value: 'qiniu' },
      ]},
    ]},
  ],
  'GET /admin/api/config/': { app_id: '', mch_id: '', api_key: '' },

  // Site Settings
  'GET /admin/api/site-settings': {
    site_name: 'LYShop', site_logo: '',
    seo_title: 'LYShop - 开源商城', seo_keywords: '商城,电商,开源', seo_description: '开源插件化商城系统',
    icp: '',
    hero_badge: '限时秒杀进行中', hero_title: '精选好物\\n品质生活从这里开始',
    hero_subtitle: '数千款精选商品，正品保障，极速发货，让购物更简单。',
    hero_btn_text: '立即选购', hero_btn_link: '/products',
    color_primary: '#dc2626', color_primary_light: '#ef4444', color_primary_dark: '#b91c1c',
    color_bg_page: '#f9fafb', color_bg_header: 'rgba(255,255,255,0.8)', color_bg_footer: '#f9fafb',
    color_price: '#ef4444', color_hero_from: '#b91c1c', color_hero_to: '#991b1b',
  },
  'PUT /admin/api/site-settings': { success: true },

  // Checkin
  'GET /admin/api/checkin/rules': [
    { id: 1, day: 0, points: 10 },
    { id: 2, day: 3, points: 20 },
    { id: 3, day: 7, points: 50 },
  ],
  'GET /admin/api/checkin/logs': { list: [
    { id: 1, user_id: 1001, date: '2026-05-22', consecutive_days: 5, points: 20 },
    { id: 2, user_id: 1002, date: '2026-05-22', consecutive_days: 1, points: 10 },
    { id: 3, user_id: 1001, date: '2026-05-21', consecutive_days: 4, points: 10 },
  ], total: 3, page: 1, size: 20 },

  // Messages
  'GET /admin/api/messages': { list: [
    { id: 1, user_id: 0, group: 'system', title: '系统升级通知', content: 'LYShop 升级至 2.0', is_read: 0, created_at: '2026-05-22T10:00:00Z' },
    { id: 2, user_id: 1001, group: 'order', title: '订单发货', content: '您的订单已发货', is_read: 1, created_at: '2026-05-21T14:00:00Z' },
    { id: 3, user_id: 0, group: 'marketing', title: '618大促', content: '全场满300减50', is_read: 0, created_at: '2026-05-20T08:00:00Z' },
  ], total: 3, page: 1, size: 20 },

  // RBAC: Admins
  'GET /admin/api/admins': [
    { id: 1, username: 'admin', role_id: 1, status: 1, created_at: '2026-01-01T00:00:00Z' },
    { id: 2, username: 'kefu01', role_id: 2, status: 1, created_at: '2026-03-15T10:00:00Z' },
    { id: 3, username: 'operator', role_id: 3, status: 1, created_at: '2026-04-01T08:00:00Z' },
  ],

  // RBAC: Roles
  'GET /admin/api/roles': [
    { id: 1, name: '超级管理员', permissions: '["*"]' },
    { id: 2, name: '客服', permissions: '["im:view","im:reply"]' },
    { id: 3, name: '运营', permissions: '["product:view","product:create","product:edit","order:view","marketing:view","marketing:edit"]' },
  ],

  // RBAC: All available permissions
  'GET /admin/api/permissions': [
    'system:admin', 'system:config',
    'product:view', 'product:create', 'product:edit', 'product:delete',
    'order:view', 'order:ship', 'order:refund', 'order:review-reply',
    'wms:view', 'wms:edit',
    'marketing:view', 'marketing:edit',
    'im:view', 'im:reply',
    'ai:view', 'ai:generate', 'ai:config',
    'decor:view', 'decor:edit',
  ],
}

export function matchMock(method: string, url: string, params?: Record<string, any>): { matched: boolean; data?: any } {
  const key = `${method.toUpperCase()} ${url}`
  const query = params || {}
  const getActivityStore = (type: 'seckill' | 'group-buy' | 'bargain') => {
    if (type === 'seckill') return seckillActivitiesSource
    if (type === 'group-buy') return groupBuyActivitiesSource
    return bargainActivitiesSource
  }
  const getProductStore = (type: 'seckill' | 'group-buy' | 'bargain') => {
    if (type === 'seckill') return seckillProductsSource
    if (type === 'group-buy') return groupBuyProductsSource
    return bargainProductsSource
  }
  const parseTypeFromPath = (path: string): 'seckill' | 'group-buy' | 'bargain' | null => {
    if (path.includes('/marketing/seckill/')) return 'seckill'
    if (path.includes('/marketing/group-buy/')) return 'group-buy'
    if (path.includes('/marketing/bargain/')) return 'bargain'
    return null
  }
  const hasTimeConflict = (store: any[], startAt: string, endAt: string, excludeID = 0) => {
    const start = new Date(startAt).getTime()
    const end = new Date(endAt).getTime()
    if (!start || !end || start >= end) return true
    return store.some((item: any) => {
      if (Number(item.id || 0) === excludeID) return false
      const left = new Date(String(item.start_at || '')).getTime()
      const right = new Date(String(item.end_at || '')).getTime()
      return left < end && right > start
    })
  }
  const toPageData = (list: any[], page: number, size: number) => {
    const p = Math.max(1, Number(page || 1))
    const s = Math.max(1, Number(size || 20))
    const offset = (p - 1) * s
    return { list: list.slice(offset, offset + s), total: list.length, page: p, size: s }
  }

  const marketingType = parseTypeFromPath(url)
  if (marketingType) {
    const activityStore = getActivityStore(marketingType)
    const productStore = getProductStore(marketingType)
    if (key.endsWith('/activities')) {
      if (method.toUpperCase() === 'GET') {
        const list = activityStore.slice().sort((a: any, b: any) => Number(b.id || 0) - Number(a.id || 0))
        return { matched: true, data: toPageData(list, Number(query.page || 1), Number(query.size || 20)) }
      }
      if (method.toUpperCase() === 'POST') {
        const payload: any = params || {}
        const startAt = String(payload.start_at || '')
        const endAt = String(payload.end_at || '')
        if (hasTimeConflict(activityStore, startAt, endAt)) {
          throw new Error('同类型活动时间段冲突')
        }
        marketingActivitySeq += 1
        const row = {
          id: marketingActivitySeq,
          type: marketingType === 'group-buy' ? 'group_buy' : marketingType,
          name: String(payload.name || `${marketingType}活动`),
          start_at: startAt,
          end_at: endAt,
          status: Number(payload.status || 0) === 1 ? 1 : 0,
        }
        activityStore.push(row)
        return { matched: true, data: row }
      }
    }
    if (method.toUpperCase() === 'PUT' && /\/activities\/\d+$/.test(url)) {
      const id = Number(url.split('/').pop() || 0)
      const target = activityStore.find((item: any) => Number(item.id || 0) === id)
      if (!target) return { matched: true, data: null }
      const payload: any = params || {}
      const startAt = String(payload.start_at || target.start_at || '')
      const endAt = String(payload.end_at || target.end_at || '')
      if (hasTimeConflict(activityStore, startAt, endAt, id)) {
        throw new Error('同类型活动时间段冲突')
      }
      target.name = String(payload.name || target.name)
      target.start_at = startAt
      target.end_at = endAt
      target.status = Number(payload.status ?? target.status) === 1 ? 1 : 0
      return { matched: true, data: null }
    }
    if (key.endsWith('/products')) {
      if (method.toUpperCase() === 'GET') {
        const activityID = Number(query.activity_id || 0)
        const keyword = String(query.keyword || '').trim().toLowerCase()
        const page = Number(query.page || 1)
        const size = Number(query.size || 20)
        let list = productStore.slice()
        if (activityID > 0) list = list.filter((item: any) => Number(item.activity_id || 0) === activityID)
        if (keyword) list = list.filter((item: any) => String(item.product_title || '').toLowerCase().includes(keyword))
        return { matched: true, data: toPageData(list, page, size) }
      }
      if (method.toUpperCase() === 'PUT' && /\/activities\/\d+\/products$/.test(url)) {
        const activityID = Number(url.split('/')[6] || 0)
        const rows = Array.isArray(params) ? params : []
        for (let i = productStore.length - 1; i >= 0; i -= 1) {
          if (Number(productStore[i].activity_id || 0) === activityID) {
            productStore.splice(i, 1)
          }
        }
        for (const row of rows) {
          marketingProductSeq += 1
          productStore.push({
            id: marketingProductSeq,
            activity_id: activityID,
            product_id: Number((row as any).product_id || 0),
            sku_id: Number((row as any).sku_id || 0),
            activity_price: Number((row as any).activity_price || 0),
            start_price: Number((row as any).start_price || 0),
            floor_price: Number((row as any).floor_price || 0),
            limit_per_order: Number((row as any).limit_per_order || 0),
            total_stock_limit: Number((row as any).total_stock_limit || 0),
            sold_qty: 0,
            product_title: String((row as any).product_title || `商品${(row as any).product_id || ''}`),
            product_cover: String((row as any).product_cover || ''),
            sku_price: Number((row as any).sku_price || 0),
            sku_stock: Number((row as any).sku_stock || 0),
            sku_attrs: Array.isArray((row as any).sku_attrs) ? (row as any).sku_attrs : [],
          })
        }
        return { matched: true, data: null }
      }
    }
  }
  const vipCrud = (path: string, source: any[]) => {
    if (key === `GET ${path}`) {
      let list = clone(source)
      if (path === '/admin/api/vip/sku-prices') {
        const productID = Number(query.product_id || 0)
        const skuID = Number(query.sku_id || 0)
        const levelID = Number(query.level_id || 0)
        const hasStatus = query.status !== undefined && query.status !== null && String(query.status).trim() !== ''
        const status = Number(query.status || 0)
        if (productID > 0) list = list.filter((item: any) => Number(item.product_id || 0) === productID)
        if (skuID > 0) list = list.filter((item: any) => Number(item.sku_id || 0) === skuID)
        if (levelID > 0) list = list.filter((item: any) => Number(item.level_id || 0) === levelID)
        if (hasStatus) list = list.filter((item: any) => Number(item.status || 0) === status)
      }
      const page = Math.max(1, Number(query.page || 1))
      const size = Math.max(1, Number(query.size || 20))
      const offset = (page - 1) * size
      return { matched: true, data: { list: list.slice(offset, offset + size), total: list.length, page, size } }
    }
    if (key === `POST ${path}`) {
      const nextID = Math.max(0, ...source.map((item: any) => Number(item.id || 0))) + 1
      source.push({ id: nextID, status: 1, ...(params || {}) })
      return { matched: true, data: null }
    }
    if (key.startsWith(`PUT ${path}/`)) {
      const id = Number(url.split('/').pop() || 0)
      const target = source.find((item: any) => Number(item.id || 0) === id)
      if (target) Object.assign(target, params || {})
      return { matched: true, data: null }
    }
    if (key.startsWith(`DELETE ${path}/`)) {
      const id = Number(url.split('/').pop() || 0)
      const idx = source.findIndex((item: any) => Number(item.id || 0) === id)
      if (idx >= 0) source.splice(idx, 1)
      return { matched: true, data: null }
    }
    return null
  }
  if (key === 'POST /api/v1/vip/coupons/monthly/1/claim' || key === 'POST /api/v1/vip/coupons/monthly/2/claim') {
    const ruleID = Number(url.split('/')[6] || 0)
    const target = (routes['GET /api/v1/vip/coupons/monthly'] as any[]).find((item: any) => Number(item.rule_id) === ruleID)
    if (target && target.claimed < target.monthly_limit) target.claimed += 1
    return { matched: true, data: { success: true } }
  }
  const vipPlanResult = vipCrud('/admin/api/vip/plans', vipPlansSource)
  if (vipPlanResult) return vipPlanResult
  const vipLevelResult = vipCrud('/admin/api/vip/levels', vipLevelsSource)
  if (vipLevelResult) return vipLevelResult
  const vipRuleResult = vipCrud('/admin/api/vip/coupon-rules', vipCouponRulesSource)
  if (vipRuleResult) return vipRuleResult
  const vipSkuResult = vipCrud('/admin/api/vip/sku-prices', vipSkuPricesSource)
  if (vipSkuResult) return vipSkuResult
  if (key === 'GET /admin/api/categories') {
    return { matched: true, data: clone(categoriesSource) }
  }
  if (key === 'POST /admin/api/categories') {
    const payload: any = params || {}
    categorySeq += 1
    const now = new Date().toISOString()
    const row = {
      id: categorySeq,
      parent_id: 0,
      name: String(payload.name || `分类${categorySeq}`),
      icon: String(payload.icon || ''),
      sort: Number(payload.sort || 0),
      status: Number(payload.status || 0) === 1 ? 1 : 0,
      created_at: now,
      updated_at: now,
    }
    categoriesSource.push(row)
    categoriesSource.sort((a: any, b: any) => Number(a.sort || 0) - Number(b.sort || 0) || Number(a.id || 0) - Number(b.id || 0))
    return { matched: true, data: clone(row) }
  }
  if (key.startsWith('PUT /admin/api/categories/')) {
    const id = Number(url.split('/').pop() || 0)
    const target = categoriesSource.find((item: any) => Number(item.id || 0) === id)
    if (!target) return { matched: true, data: null }
    const payload: any = params || {}
    if (payload.name !== undefined) target.name = String(payload.name || target.name)
    if (payload.icon !== undefined) target.icon = String(payload.icon || '')
    if (payload.sort !== undefined) target.sort = Number(payload.sort || 0)
    if (payload.status !== undefined) target.status = Number(payload.status || 0) === 1 ? 1 : 0
    target.updated_at = new Date().toISOString()
    categoriesSource.sort((a: any, b: any) => Number(a.sort || 0) - Number(b.sort || 0) || Number(a.id || 0) - Number(b.id || 0))
    return { matched: true, data: null }
  }
  if (key.startsWith('DELETE /admin/api/categories/')) {
    const id = Number(url.split('/').pop() || 0)
    const idx = categoriesSource.findIndex((item: any) => Number(item.id || 0) === id)
    if (idx >= 0) categoriesSource.splice(idx, 1)
    return { matched: true, data: null }
  }
  if (key === 'GET /admin/api/orders') {
    const status = toNumber(query.status)
    const list = status > 0
      ? orderListSource.filter((item: any) => toNumber(item.status) === status)
      : orderListSource.slice()
    return { matched: true, data: { ...orders, list, total: list.length } }
  }
  if (key === 'GET /admin/api/reviews') {
    return { matched: true, data: listReviews(query) }
  }
  if (key === 'GET /admin/api/after-sales') {
    const status = String(query.status || '').trim()
    const orderID = Number(query.order_id || 0)
    const page = Math.max(1, Number(query.page || 1))
    const size = Math.max(1, Number(query.size || 20))
    let list = afterSalesSource.slice()
    if (status) list = list.filter((item: any) => String(item.status || '') === status)
    if (orderID > 0) list = list.filter((item: any) => Number(item.order_id || 0) === orderID)
    const offset = (page - 1) * size
    const pageList = list.slice(offset, offset + size).map((item: any) => withAfterSaleLabels(item))
    return { matched: true, data: { list: pageList, total: list.length, page, size } }
  }
  if (key.startsWith('GET /admin/api/after-sales/')) {
    const id = Number(url.split('/').pop() || 0)
    const target = afterSalesSource.find((item: any) => Number(item.id) === id) || null
    return { matched: true, data: target ? withAfterSaleLabels(target) : null }
  }
  if (key.startsWith('POST /admin/api/after-sales/')) {
    const id = Number(url.split('/')[4] || 0)
    const action = String(url.split('/')[5] || '')
    const target = afterSalesSource.find((item: any) => Number(item.id) === id)
    if (!target) return { matched: true, data: null }
    const now = new Date().toISOString()
    const pushLog = (fromStatus: string, toStatus: string, actionName: string, content: string) => {
      target.logs = Array.isArray(target.logs) ? target.logs : []
      target.logs.push({
        id: Math.floor(Math.random() * 100000),
        case_id: id,
        from_status: fromStatus,
        to_status: toStatus,
        action: actionName,
        operator_type: 'admin',
        operator_id: 1,
        content,
        created_at: now,
      })
    }
    if (action === 'audit') {
      const approve = !!(params as any)?.approve
      const fromStatus = String(target.status || '')
      target.status = approve ? 'approved_wait_user_return' : 'rejected'
      pushLog(fromStatus, target.status, 'audit', approve ? '售后审核通过' : '售后审核拒绝')
      touchOrderAfterSaleSummary(Number(target.order_id || 0))
      return { matched: true, data: null }
    }
    if (action === 'receive') {
      const fromStatus = String(target.status || '')
      target.status = String(target.case_type || '') === 'exchange' ? 'reship_pending' : 'refund_pending'
      pushLog(fromStatus, target.status, 'receive', '仓库收货确认')
      touchOrderAfterSaleSummary(Number(target.order_id || 0))
      return { matched: true, data: null }
    }
    if (action === 'refund') {
      const fromStatus = String(target.status || '')
      target.status = 'refunded'
      target.refund_amount = Number((params as any)?.amount || target.refund_amount || 0)
      pushLog(fromStatus, target.status, 'refund', '退款登记')
      touchOrderAfterSaleSummary(Number(target.order_id || 0))
      return { matched: true, data: null }
    }
    if (action === 'complete') {
      const fromStatus = String(target.status || '')
      target.status = 'completed'
      pushLog(fromStatus, target.status, 'complete', '售后完结')
      touchOrderAfterSaleSummary(Number(target.order_id || 0))
      return { matched: true, data: null }
    }
    if (action === 'close') {
      const fromStatus = String(target.status || '')
      target.status = 'closed'
      pushLog(fromStatus, target.status, 'close', '关闭售后')
      touchOrderAfterSaleSummary(Number(target.order_id || 0))
      return { matched: true, data: null }
    }
    return { matched: true, data: null }
  }
  if (key.startsWith('POST /admin/api/reviews/') && key.endsWith('/reply')) {
    const id = Number(url.split('/')[4] || 0)
    const content = String((params as any)?.content || '').trim()
    if (id && content) upsertReviewReply(id, content)
    return { matched: true, data: null }
  }
  if (key === 'GET /admin/api/decor/index/variants') {
    return { matched: true, data: clone(decorVariantsSource) }
  }
  if (key.startsWith('GET /admin/api/decor/index')) {
    const parsed = new URL(url, 'https://mock.local')
    const variantKey = String(parsed.searchParams.get('variant') || 'default')
    const item = decorVariantsSource.find((row: any) => String(row.variant_key) === variantKey) || decorVariantsSource[0]
    return { matched: true, data: clone(item) }
  }
  if (key.startsWith('PUT /admin/api/decor/index')) {
    const parsed = new URL(url, 'https://mock.local')
    const variantKey = String(parsed.searchParams.get('variant') || 'default')
    const target = decorVariantsSource.find((row: any) => String(row.variant_key) === variantKey)
    if (target) {
      target.components = JSON.stringify((params as any)?.components || [])
    }
    return { matched: true, data: target ? clone(target) : null }
  }
  if (key.startsWith('POST /admin/api/decor/index/publish')) {
    const parsed = new URL(url, 'https://mock.local')
    const variantKey = String(parsed.searchParams.get('variant') || 'default')
    const now = new Date().toISOString()
    for (const row of decorVariantsSource) {
      row.is_published = false
      row.published_at = null
    }
    const target = decorVariantsSource.find((row: any) => String(row.variant_key) === variantKey)
    if (target) {
      target.is_published = true
      target.published_at = now
    }
    return { matched: true, data: null }
  }
  if (key === 'POST /admin/api/decor/index/copies') {
    const payload: any = params || {}
    const fromVariantKey = String(payload.from_variant_key || 'default')
    const newVariantKey = String(payload.new_variant_key || '').trim()
    const newVariantName = String(payload.new_variant_name || '').trim() || `副本 ${newVariantKey}`
    const source = decorVariantsSource.find((row: any) => String(row.variant_key) === fromVariantKey)
    if (!source || !newVariantKey) return { matched: true, data: null }
    if (!decorVariantsSource.find((row: any) => String(row.variant_key) === newVariantKey)) {
      decorVariantsSource.push({
        ...clone(source),
        id: Math.max(...decorVariantsSource.map((row: any) => Number(row.id || 0))) + 1,
        variant_key: newVariantKey,
        variant_name: newVariantName,
        is_published: false,
        published_at: null,
      })
    }
    return { matched: true, data: null }
  }
  if (key.startsWith('PUT /admin/api/decor/index/variants/')) {
    const variantKey = decodeURIComponent(url.split('/').pop() || '')
    const target = decorVariantsSource.find((row: any) => String(row.variant_key) === variantKey)
    if (target) {
      target.variant_name = String((params as any)?.variant_name || target.variant_name)
    }
    return { matched: true, data: null }
  }
  if (key.startsWith('DELETE /admin/api/decor/index/variants/')) {
    const variantKey = decodeURIComponent(url.split('/').pop() || '')
    if (variantKey !== 'default') {
      const idx = decorVariantsSource.findIndex((row: any) => String(row.variant_key) === variantKey && !row.is_published)
      if (idx >= 0) decorVariantsSource.splice(idx, 1)
    }
    return { matched: true, data: null }
  }
  // PC Decor
  if (key === 'GET /admin/api/decor/pc') {
    return { matched: true, data: clone(pcDecorSource) }
  }
  if (key === 'PUT /admin/api/decor/pc') {
    const payload = (params as any)?.components
    pcDecorSource.components = JSON.stringify(payload || { pageStyle: null, components: [] })
    return { matched: true, data: clone(pcDecorSource) }
  }
  if (key === 'POST /admin/api/decor/pc/publish') {
    pcDecorSource.is_published = true
    return { matched: true, data: null }
  }

  if (key in routes) return { matched: true, data: routes[key] }
  for (const pattern of Object.keys(routes)) {
    if (key.startsWith(pattern) && pattern.endsWith('/')) {
      if (pattern === 'GET /admin/api/orders/') {
        const id = Number(url.split('/').pop() || 0)
        const detail = orderListSource.find((item: any) => Number(item.id) === id) || null
        return { matched: true, data: detail }
      }
      if (pattern === 'GET /admin/api/reviews/') {
        const id = Number(url.split('/').pop() || 0)
        return { matched: true, data: getReviewByID(id) || null }
      }
      return { matched: true, data: routes[pattern] }
    }
  }
  if (['POST', 'PUT', 'DELETE'].includes(method.toUpperCase())) {
    if (key === 'POST /admin/api/reviews/') {
      const id = Number((params as any)?.review_id || 0) || Number((url.split('/').pop() || 0))
      const content = String((params as any)?.content || '').trim()
      if (id && content) {
        upsertReviewReply(id, content)
      }
      return { matched: true, data: null }
    }
    return { matched: true, data: null }
  }
  return { matched: false }
}
