// Admin mock data — routes keyed by "METHOD /admin/api/..."
// Reuses product/order JSON from app mock where applicable

import products from '../../../app/mock/data/products.json'
import categories from '../../../app/mock/data/categories.json'
import orders from '../../../app/mock/data/orders.json'
import productDetail from '../../../app/mock/data/product-detail.json'

const orderListSource = (orders as any)?.list || []
const toNumber = (v: any) => Number(v || 0)
let replySeq = 20000

function clone<T>(v: T): T {
  return JSON.parse(JSON.stringify(v))
}

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

  // Menus (dynamically generated from enabled plugins)
  'GET /admin/api/menus': [
    { title: '商品管理', icon: 'box', path: '/product', sort: 10, children: [
      { title: '商品列表', path: '/product/list' },
      { title: '新增商品', path: '/product/form' },
    ]},
    { title: '订单管理', icon: 'shopping-cart', path: '/order', sort: 20, children: [
      { title: '订单列表', path: '/order/list' },
    ]},
    { title: '评价管理', icon: 'star', path: '/review', sort: 21, children: [
      { title: '评价列表', path: '/review/list' },
    ]},
    { title: '仓储管理', icon: 'warehouse', path: '/wms', sort: 30, children: [
      { title: '库存管理', path: '/wms/stock' },
    ]},
    { title: '营销管理', icon: 'tag', path: '/marketing', sort: 40, children: [
      { title: '优惠券管理', path: '/marketing/coupon' },
    ]},
    { title: '客服中心', icon: 'message-circle', path: '/im', sort: 50, children: [
      { title: '客服会话', path: '/im/sessions' },
    ]},
    { title: 'AI 工具', icon: 'cpu', path: '/ai', sort: 60, children: [
      { title: 'AI 生图', path: '/ai/tasks' },
      { title: '模型配置', path: '/ai/models' },
    ]},
    { title: '店铺装修', icon: 'layout', path: '/decor', sort: 70, children: [
      { title: '首页装修', path: '/decor/index' },
    ]},
    { title: '签到管理', icon: 'calendar', path: '/checkin', sort: 75, children: [
      { title: '签到规则', path: '/checkin/rules' },
      { title: '签到记录', path: '/checkin/logs' },
    ]},
    { title: '消息管理', icon: 'bell', path: '/message', sort: 80, children: [
      { title: '消息列表', path: '/message/list' },
      { title: '发送消息', path: '/message/send' },
    ]},
    { title: '系统设置', icon: 'settings', path: '/system', sort: 90, children: [
      { title: '配置中心', path: '/system/config' },
      { title: '管理员管理', path: '/system/admins' },
      { title: '角色管理', path: '/system/roles' },
    ]},
  ],

  // Dashboard
  'GET /admin/api/dashboard': {
    today_orders: 56,
    today_sales: 28960.50,
    pending_refunds: 3,
    online_sessions: 2,
  },

  // Products
  'GET /admin/api/products': products,
  'GET /admin/api/products/': productDetail,
  'POST /admin/api/products': { id: 100 },
  'GET /admin/api/categories': categories,
  'POST /admin/api/categories': { id: 10 },

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

  // Decor
  'GET /admin/api/decor/': {
    id: 1, page_key: 'index', merchant_id: 0,
    components: JSON.stringify([
      { type: 'banner', id: 'c1', props: { images: [{ url: 'https://picsum.photos/750/350?random=1' }], height: 350 } },
      { type: 'product_grid', id: 'c2', props: { source: 'hot', limit: 6, columns: 2 } },
    ]),
  },

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
  ],
  'GET /admin/api/config/': { app_id: '', mch_id: '', api_key: '' },

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
  if (key.startsWith('POST /admin/api/reviews/') && key.endsWith('/reply')) {
    const id = Number(url.split('/')[4] || 0)
    const content = String((params as any)?.content || '').trim()
    if (id && content) upsertReviewReply(id, content)
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
