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

// Seed 3-level categories (overwrite simplified data)
categoriesSource.length = 0
const seedCats = [
  { id: 1, parent_id: 0, name: '电子', sort: 1 },
  { id: 2, parent_id: 0, name: '服饰', sort: 2 },
  { id: 3, parent_id: 0, name: '家居', sort: 3 },
  { id: 4, parent_id: 0, name: '食品', sort: 4 },
  { id: 5, parent_id: 0, name: '美妆', sort: 5 },
  { id: 11, parent_id: 1, name: '手机数码', sort: 1 },
  { id: 12, parent_id: 1, name: '电脑办公', sort: 2 },
  { id: 13, parent_id: 1, name: '智能穿戴', sort: 3 },
  { id: 21, parent_id: 2, name: '女装', sort: 1 },
  { id: 22, parent_id: 2, name: '男装', sort: 2 },
  { id: 23, parent_id: 2, name: '鞋包配饰', sort: 3 },
  { id: 31, parent_id: 3, name: '家具', sort: 1 },
  { id: 32, parent_id: 3, name: '日用百货', sort: 2 },
  { id: 41, parent_id: 4, name: '生鲜', sort: 1 },
  { id: 42, parent_id: 4, name: '零食饮料', sort: 2 },
  { id: 51, parent_id: 5, name: '彩妆', sort: 1 },
  { id: 52, parent_id: 5, name: '护肤', sort: 2 },
  { id: 111, parent_id: 11, name: '智能手机', sort: 1 },
  { id: 112, parent_id: 11, name: '配件', sort: 2 },
  { id: 211, parent_id: 21, name: '连衣裙', sort: 1 },
  { id: 212, parent_id: 21, name: '外套', sort: 2 },
]
for (const c of seedCats) categoriesSource.push({ ...c, product_count: 0 })

let categorySeq = Math.max(0, ...categoriesSource.map((item: any) => Number(item?.id || 0)))

// Enrich first 6 products with SKU/category data
for (let i = 0; i < Math.min(6, productListSource.length); i++) {
  const p: any = productListSource[i]
  p.category_id = [111, 112, 211, 211, 31, 51][i] || 111
  p.category_path_name = ['电子 / 手机数码 / 智能手机', '电子 / 手机数码 / 配件', '服饰 / 女装 / 连衣裙', '服饰 / 女装 / 连衣裙', '家居 / 家具', '美妆 / 彩妆'][i] || ''
  p.sales_count = 80 + i * 23
  p.low_stock_threshold = 10
  p.skus = [
    { id: p.id * 10 + 1, attrs: [{ name: '颜色', value: '黑色' }, { name: '版本', value: '标准版' }], price: Number(p.price || 0), stock: 50 },
    { id: p.id * 10 + 2, attrs: [{ name: '颜色', value: '白色' }, { name: '版本', value: '标准版' }], price: Number(p.price || 0), stock: 30 },
    { id: p.id * 10 + 3, attrs: [{ name: '颜色', value: '黑色' }, { name: '版本', value: '尊享版' }], price: Math.round(Number(p.price || 0) * 1.2 * 100) / 100, stock: 20 },
    { id: p.id * 10 + 4, attrs: [{ name: '颜色', value: '白色' }, { name: '版本', value: '尊享版' }], price: Math.round(Number(p.price || 0) * 1.2 * 100) / 100, stock: 15 },
  ]
}
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
const couponsSource: any[] = [
  { id: 1, name: '新人满100减20', type: 1, min_amount: 100, discount: 20, total_count: 1000, per_limit: 1, status: 1, used_count: 120, start_at: '2026-05-01T00:00:00Z', end_at: '2026-06-30T23:59:59Z', description: '新用户专享', stack_rule: 'exclusive', target_type: 'new_user', target_value: '' },
  { id: 2, name: '全场9折券', type: 2, min_amount: 0, discount: 0.9, total_count: 500, per_limit: 1, status: 1, used_count: 85, start_at: '2026-05-01T00:00:00Z', end_at: '2026-07-31T23:59:59Z', description: '全场通用折扣', stack_rule: 'same_type', target_type: 'all', target_value: '' },
  { id: 3, name: '无门槛5元券', type: 3, min_amount: 0, discount: 5, total_count: 0, per_limit: 3, status: 1, used_count: 300, start_at: '2026-05-01T00:00:00Z', end_at: '2026-12-31T23:59:59Z', description: '无门槛使用', stack_rule: 'cross_type', target_type: 'all', target_value: '' },
]
let couponSeq = 3
const specTemplatesSource: any[] = [
  { id: 1, name: '服装通用', category_ids: [21, 22], attrs: [{ name: '颜色', values: ['黑色', '白色', '红色', '蓝色'] }, { name: '尺码', values: ['S', 'M', 'L', 'XL', 'XXL'] }], status: 1 },
  { id: 2, name: '数码配件', category_ids: [11, 12], attrs: [{ name: '颜色', values: ['黑色', '白色', '银色'] }, { name: '版本', values: ['标准版', '高配版'] }], status: 1 },
  { id: 3, name: '鞋类', category_ids: [23], attrs: [{ name: '颜色', values: ['黑色', '棕色', '白色'] }, { name: '尺码', values: ['36', '37', '38', '39', '40', '41', '42', '43', '44'] }], status: 1 },
]
let specTemplateSeq = 3
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
const wmsWarehousesSource: any[] = [
  { id: 1, name: '主仓库', code: 'WH-SH-01', address: '上海市浦东新区', contact: '张三', phone: '13800001111', status: 1, created_at: '2026-05-20T08:00:00Z', updated_at: '2026-05-27T10:00:00Z' },
  { id: 2, name: '华南仓', code: 'WH-GZ-01', address: '广州市天河区', contact: '李四', phone: '13800002222', status: 1, created_at: '2026-05-20T08:00:00Z', updated_at: '2026-05-27T10:00:00Z' },
]
const wmsStocksSource: any[] = [
  { id: 1, warehouse_id: 1, warehouse_name: '主仓库', sku_id: 101, sku_name: '蓝牙耳机标准版', qty: 120, safe_qty: 20, updated_at: '2026-05-28T09:10:00Z' },
  { id: 2, warehouse_id: 1, warehouse_name: '主仓库', sku_id: 102, sku_name: '无线键盘', qty: 35, safe_qty: 30, updated_at: '2026-05-28T09:12:00Z' },
  { id: 3, warehouse_id: 2, warehouse_name: '华南仓', sku_id: 101, sku_name: '蓝牙耳机标准版', qty: 60, safe_qty: 18, updated_at: '2026-05-28T09:00:00Z' },
  { id: 4, warehouse_id: 2, warehouse_name: '华南仓', sku_id: 205, sku_name: '运动手表', qty: 12, safe_qty: 10, updated_at: '2026-05-28T08:52:00Z' },
]
const wmsDocsSource: any[] = [
  {
    id: 1,
    doc_no: 'IN202605280001',
    doc_type: 'inbound',
    status: 'draft',
    warehouse_id: 1,
    warehouse_name: '主仓库',
    remark: '供应商到货待上架',
    items: [
      { id: 1, sku_id: 101, sku_name: '蓝牙耳机标准版', qty: 30, unit_cost: 199 },
      { id: 2, sku_id: 205, sku_name: '运动手表', qty: 5, unit_cost: 459 },
    ],
    total_qty: 35,
    created_at: '2026-05-28T07:30:00Z',
    updated_at: '2026-05-28T07:30:00Z',
  },
  {
    id: 2,
    doc_no: 'OUT202605280001',
    doc_type: 'outbound',
    status: 'completed',
    warehouse_id: 1,
    warehouse_name: '主仓库',
    remark: '电商订单批量出库',
    items: [
      { id: 3, sku_id: 101, sku_name: '蓝牙耳机标准版', qty: 15, unit_cost: 0 },
      { id: 4, sku_id: 102, sku_name: '无线键盘', qty: 10, unit_cost: 0 },
    ],
    total_qty: 25,
    created_at: '2026-05-28T08:10:00Z',
    updated_at: '2026-05-28T08:25:00Z',
  },
]
const wmsMovementsSource: any[] = [
  { id: 1, doc_id: 2, doc_no: 'OUT202605280001', biz_type: 'outbound', warehouse_id: 1, warehouse_name: '主仓库', sku_id: 101, sku_name: '蓝牙耳机标准版', change_qty: -15, before_qty: 135, after_qty: 120, occurred_at: '2026-05-28T08:21:00Z' },
  { id: 2, doc_id: 2, doc_no: 'OUT202605280001', biz_type: 'outbound', warehouse_id: 1, warehouse_name: '主仓库', sku_id: 102, sku_name: '无线键盘', change_qty: -10, before_qty: 45, after_qty: 35, occurred_at: '2026-05-28T08:22:00Z' },
]
let marketingActivitySeq = 100
let marketingProductSeq = 5000
let wmsWarehouseSeq = Math.max(0, ...wmsWarehousesSource.map((row: any) => Number(row.id || 0)))
let wmsStockSeq = Math.max(0, ...wmsStocksSource.map((row: any) => Number(row.id || 0)))
let wmsDocSeq = Math.max(0, ...wmsDocsSource.map((row: any) => Number(row.id || 0)))
let wmsDocItemSeq = Math.max(0, ...wmsDocsSource.flatMap((doc: any) => (Array.isArray(doc.items) ? doc.items : [])).map((row: any) => Number(row.id || 0)))
let wmsMovementSeq = Math.max(0, ...wmsMovementsSource.map((row: any) => Number(row.id || 0)))

const shopsCurrentSource: any = {
  id: 1, name: '示范品牌旗舰店', logo: 'https://picsum.photos/200/200?random=shop1', owner: 'admin', decor_status: 'published',
}

const announcementsSource: any[] = [
  { id: 1, title: '平台 2026 年 6 月例行升级通知', content: '6 月 3 日 02:00-03:00 短暂维护', type: 'normal', created_at: '2026-05-27T10:00:00Z' },
  { id: 2, title: '618 大促招商进行中', content: '前往营销中心报名参与', type: 'urgent', created_at: '2026-05-26T09:00:00Z' },
  { id: 3, title: '电子面单服务费下调', content: '6 月起部分快递公司面单费下调 10%', type: 'normal', created_at: '2026-05-25T14:00:00Z' },
]

function seedRandom(seed: number) {
  const x = Math.sin(seed) * 10000
  return x - Math.floor(x)
}

function buildRevenueTrend(days: number): { categories: string[]; series: Array<{ name: string; data: number[] }> } {
  const categories: string[] = []
  const data: number[] = []
  const today = new Date('2026-05-28T00:00:00Z')
  for (let i = days - 1; i >= 0; i--) {
    const d = new Date(today.getTime() - i * 86400000)
    const mm = String(d.getUTCMonth() + 1).padStart(2, '0')
    const dd = String(d.getUTCDate()).padStart(2, '0')
    categories.push(`${mm}-${dd}`)
    const seed = d.getUTCFullYear() * 10000 + (d.getUTCMonth() + 1) * 100 + d.getUTCDate()
    const dayOfWeek = d.getUTCDay()
    const base = 4800 + Math.sin(seed / 7) * 1200
    const weekendBoost = (dayOfWeek === 0 || dayOfWeek === 6) ? 600 : 0
    const noise = (seedRandom(seed) - 0.5) * 800
    data.push(Math.max(800, Math.round(base + weekendBoost + noise)))
  }
  return { categories, series: [{ name: '营收', data }] }
}

function buildOrderTrend(days: number): { categories: string[]; series: Array<{ name: string; data: number[] }> } {
  const r = buildRevenueTrend(days)
  return { categories: r.categories, series: [{ name: '订单', data: r.series[0].data.map((v) => Math.max(2, Math.round(v / 120))) }] }
}

function clone<T>(v: T): T {
  return JSON.parse(JSON.stringify(v))
}

function wmsNowISO() {
  return new Date().toISOString()
}

function wmsWarehouseName(warehouseID: number) {
  const found = wmsWarehousesSource.find((row: any) => Number(row.id || 0) === Number(warehouseID))
  return String(found?.name || '')
}

function wmsDocNoPrefix(type: string) {
  return type === 'inbound' ? 'IN' : 'OUT'
}

function nextWmsDocNo(type: string) {
  const prefix = wmsDocNoPrefix(type)
  const dateKey = new Date().toISOString().slice(0, 10).replace(/-/g, '')
  const currentCount = wmsDocsSource.filter((row: any) => String(row.doc_no || '').startsWith(`${prefix}${dateKey}`)).length
  const serial = String(currentCount + 1).padStart(4, '0')
  return `${prefix}${dateKey}${serial}`
}

function normalizeWmsItems(rows: any[]): any[] {
  const list = Array.isArray(rows) ? rows : []
  return list
    .map((row: any) => {
      const qty = Math.max(0, Number(row?.qty || 0))
      const skuID = Number(row?.sku_id || 0)
      if (!skuID || !qty) return null
      wmsDocItemSeq += 1
      return {
        id: Number(row?.id || 0) || wmsDocItemSeq,
        sku_id: skuID,
        sku_name: String(row?.sku_name || `SKU-${skuID}`),
        qty,
        unit_cost: Math.max(0, Number(row?.unit_cost || 0)),
        remark: String(row?.remark || ''),
      }
    })
    .filter(Boolean)
}

function wmsDocTotalQty(items: any[]) {
  return (Array.isArray(items) ? items : []).reduce((sum: number, row: any) => sum + Math.max(0, Number(row?.qty || 0)), 0)
}

function ensureWmsStock(warehouseID: number, skuID: number, skuName: string) {
  const found = wmsStocksSource.find((row: any) => Number(row.warehouse_id || 0) === warehouseID && Number(row.sku_id || 0) === skuID)
  if (found) return found
  wmsStockSeq += 1
  const created = {
    id: wmsStockSeq,
    warehouse_id: warehouseID,
    warehouse_name: wmsWarehouseName(warehouseID),
    sku_id: skuID,
    sku_name: skuName || `SKU-${skuID}`,
    qty: 0,
    safe_qty: 10,
    updated_at: wmsNowISO(),
  }
  wmsStocksSource.push(created)
  return created
}

function applyWmsDocStockChange(doc: any) {
  const docType = String(doc?.doc_type || 'inbound')
  const warehouseID = Number(doc?.warehouse_id || 0)
  const now = wmsNowISO()
  const items = Array.isArray(doc?.items) ? doc.items : []
  type StockCalc = {
    row: any | null
    warehouse_id: number
    sku_id: number
    sku_name: string
    current_qty: number
  }
  const touched = new Map<string, StockCalc>()
  const pendingMovements: any[] = []

  for (const item of items) {
    const skuID = Number(item?.sku_id || 0)
    const qty = Math.max(0, Number(item?.qty || 0))
    if (!skuID || !qty) continue
    const delta = docType === 'inbound' ? qty : -qty
    const stockKey = `${warehouseID}_${skuID}`
    let entry = touched.get(stockKey)
    if (!entry) {
      const found = wmsStocksSource.find((row: any) => Number(row.warehouse_id || 0) === warehouseID && Number(row.sku_id || 0) === skuID) || null
      entry = {
        row: found,
        warehouse_id: warehouseID,
        sku_id: skuID,
        sku_name: String(item?.sku_name || found?.sku_name || `SKU-${skuID}`),
        current_qty: Number(found?.qty || 0),
      }
      touched.set(stockKey, entry)
    } else if (String(item?.sku_name || '').trim()) {
      entry.sku_name = String(item.sku_name)
    }
    const beforeQty = Number(entry.current_qty || 0)
    const afterQty = beforeQty + delta
    if (afterQty < 0) {
      throw new Error(`SKU ${skuID} 库存不足，单据无法完成`)
    }
    entry.current_qty = afterQty
    pendingMovements.push({
      doc_id: Number(doc.id || 0),
      doc_no: String(doc.doc_no || ''),
      biz_type: docType,
      warehouse_id: warehouseID,
      warehouse_name: wmsWarehouseName(warehouseID),
      sku_id: skuID,
      sku_name: entry.sku_name,
      change_qty: delta,
      before_qty: beforeQty,
      after_qty: afterQty,
      occurred_at: now,
    })
  }

  for (const entry of touched.values()) {
    const stockRow = entry.row || ensureWmsStock(entry.warehouse_id, entry.sku_id, entry.sku_name)
    stockRow.qty = entry.current_qty
    stockRow.warehouse_name = wmsWarehouseName(entry.warehouse_id)
    stockRow.sku_name = entry.sku_name
    stockRow.updated_at = now
  }

  for (const row of pendingMovements) {
    wmsMovementSeq += 1
    wmsMovementsSource.unshift({
      id: wmsMovementSeq,
      ...row,
    })
  }
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

// Expand orders to 40 covering all statuses
function expandOrderListTo40() {
  if (orderListSource.length >= 40) return
  const baseCount = orderListSource.length
  const need = 40 - baseCount
  const statuses = ['1', '2', '3', '4', '5']
  const provinces = ['上海市浦东新区', '北京市朝阳区', '广东省广州市', '浙江省杭州市', '四川省成都市', '湖北省武汉市', '辽宁省沈阳市']
  const companies = ['SF', 'ZTO', 'YTO', 'STO', 'YD', 'JD', 'EMS']
  const pickFromProducts = (i: number) => {
    const arr = (productListSource as any[]).slice(0, 6)
    return arr[i % arr.length]
  }
  for (let i = 0; i < need; i++) {
    const idx = baseCount + i + 1
    const id = 10000 + idx
    const status = statuses[i % statuses.length]
    const product = pickFromProducts(i)
    const qty = 1 + (i % 3)
    const price = Number(product?.price || 199)
    const goods = +(price * qty).toFixed(2)
    const discount = +((goods * (i % 3 === 0 ? 0.1 : 0)).toFixed(2))
    const pay = +(goods - discount).toFixed(2)
    const dateOffset = 86400000 * (i % 14)
    const created = new Date(Date.parse('2026-05-28T08:00:00Z') - dateOffset).toISOString()
    const order: any = {
      id,
      status,
      status_label: { '1': '待付款', '2': '待发货', '3': '已发货', '4': '已完成', '5': '已关闭' }[status],
      pay_method: ['wechat', 'alipay', 'wechat'][i % 3],
      user_nickname: ['张三', '李四', '王五', '赵六', '陈七', '何八'][i % 6],
      receiver_name: ['张先生', '李女士', '王先生', '赵女士'][i % 4],
      receiver_phone: '138****' + String(1000 + i).slice(-4),
      receiver_address: provinces[i % provinces.length] + ` ${i + 1}号`,
      goods_amount: goods,
      discount_amount: discount,
      pay_amount: pay,
      total_amount: pay,
      amount_breakdown: { goods_amount: goods, discount_amount: discount, payable_amount: pay },
      created_at: created,
      items: [{ id: idx * 100 + 1, product_id: product?.id, title: product?.title, cover: product?.cover || `https://picsum.photos/120/120?random=${1000 + i}`, qty, price }],
      shipments: status === '3' || status === '4' ? [{ id: idx * 10, company: companies[i % companies.length], tracking_no: `SF${1000000 + idx}`, delivery_type: 'express', logistics_status: status === '4' ? 'signed' : 'in_transit', logistics_status_label: status === '4' ? '已签收' : '运输中', created_at: created }] : [],
      logs: [],
      notes: [],
      has_after_sale: i % 7 === 0,
    }
    orderListSource.push(order)
  }
}
expandOrderListTo40()

// Add 2 more announcements (total 5)
announcementsSource.push(
  { id: 4, title: '物流时效升级', content: '部分城市顺丰提供「夜间寄」服务', type: 'normal', created_at: '2026-05-24T10:00:00Z' },
  { id: 5, title: '商家服务热线变更', content: '客服热线已升级 7x24，新号码 400-000-XXXX', type: 'urgent', created_at: '2026-05-23T18:00:00Z' },
)

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

// Seed 3 after-sale records for eapp
afterSalesSource.push(
  { id: 8001, order_id: 1, status: 'applied', status_label: '待审核', case_type: 'refund_only', case_type_label: '仅退款', reason: '少件', refund_amount: 99, created_at: '2026-05-27T10:00:00Z', messages: [], evidences: [], logs: [{ id: 1, action: 'apply', action_label: '申请', from_status: '', to_status: 'applied', to_status_label: '待审核', content: '买家提交申请', created_at: '2026-05-27T10:00:00Z' }], shipments: [] },
  { id: 8002, order_id: 2, status: 'user_returning', status_label: '退货中', case_type: 'return_refund', case_type_label: '退货退款', reason: '不喜欢', refund_amount: 199, created_at: '2026-05-26T15:30:00Z', messages: [
    { id: 1, from: 'user', content: '我已寄回，请查收', images: [], created_at: '2026-05-26T18:00:00Z' },
    { id: 2, from: 'merchant', content: '好的，签收后会立即处理退款', images: [], created_at: '2026-05-26T19:00:00Z' },
  ], evidences: [], logs: [], shipments: [{ id: 1, company: 'SF', tracking_no: 'SF889977', direction: 'inbound', direction_label: '回寄', biz_type: 'return', biz_type_label: '回寄', logistics_status: 'in_transit', logistics_status_label: '运输中', created_at: '2026-05-26T17:30:00Z' }] },
  { id: 8003, order_id: 3, status: 'refund_pending', status_label: '退款中', case_type: 'refund_only', case_type_label: '仅退款', reason: '商品损坏', refund_amount: 459, created_at: '2026-05-25T09:00:00Z', messages: [], evidences: [{ id: 1, images: ['https://picsum.photos/200/200?random=evi1'], remark: '商品破损照片', created_at: '2026-05-25T09:10:00Z' }], logs: [], shipments: [] },
)

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
        title_key: 'menu.wms',
        icon: 'warehouse',
        sort: 50,
        menus: [
          { title: '仓储管理', title_key: 'menu.wms', icon: 'warehouse', path: '/wms', sort: 10, children: [
            { title: '库存台账', title_key: 'menu.wms.stock', path: '/wms/stock' },
            { title: '仓库管理', title_key: 'menu.wms.warehouse', path: '/wms/warehouse' },
            { title: '出入库单', title_key: 'menu.wms.docs', path: '/wms/docs' },
            { title: '库存流水', title_key: 'menu.wms.movements', path: '/wms/movements' },
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

  // Dashboard — handled dynamically in matchMock

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
  'GET /admin/api/wms/warehouses': wmsWarehousesSource,
  'GET /admin/api/wms/stocks': {
    list: wmsStocksSource,
    total: wmsStocksSource.length,
    page: 1,
    size: 20,
  },

  // Marketing
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
  const methodUpper = method.toUpperCase()

  // Dashboard — dynamic
  if (key === 'GET /admin/api/dashboard') {
    const pendingShip = orderListSource.filter((o: any) => Number(o.status) === 2).length
    const pendingAfterSale = afterSalesSource.filter((a: any) => !['completed', 'closed', 'refunded', 'rejected'].includes(String(a.status || ''))).length
    const todayOrders = 56
    const todaySales = 28960.50
    const todayAvgPrice = todayOrders > 0 ? Math.round(todaySales / todayOrders * 100) / 100 : 0
    const statusMap: Record<string, number> = {}
    for (const o of orderListSource) {
      const s = String(o.status || '0')
      statusMap[s] = (statusMap[s] || 0) + 1
    }
    const statusLabels: Record<string, string> = { '0': '待付款', '1': '待确认', '2': '待发货', '3': '已发货', '4': '已完成', '5': '已关闭' }
    const statusDistribution = Object.entries(statusMap).map(([k, v]) => ({ name: statusLabels[k] || `状态${k}`, value: v }))
    const hotProducts = productListSource.slice().sort((a: any, b: any) => Number(b.sold_qty || b.sales || 0) - Number(a.sold_qty || a.sales || 0)).slice(0, 5).map((p: any) => ({
      id: Number(p.id || 0),
      title: String(p.title || p.name || ''),
      cover: String(p.cover || p.image || ''),
      sold_qty: Number(p.sold_qty || p.sales || 0),
    }))
    const stockWarningList = productListSource.filter((p: any) => Number(p.stock || 0) < 10 && Number(p.stock || 0) >= 0).slice(0, 5).map((p: any) => ({
      product_id: Number(p.id || 0),
      sku_id: Number(p.sku_id || p.id || 0),
      title: String(p.title || p.name || ''),
      stock: Number(p.stock || 0),
      threshold: 10,
    }))
    return {
      matched: true,
      data: {
        today_orders: todayOrders,
        today_sales: todaySales,
        today_avg_price: todayAvgPrice,
        pending_ship: pendingShip,
        pending_after_sale: pendingAfterSale,
        unread_message: 3,
        stock_warning: stockWarningList.length,
        compare: { revenue_yoy: 12.5, revenue_mom: 5.3, order_yoy: 8.2, order_mom: -2.1 },
        trend: {
          revenue_7d: buildRevenueTrend(7),
          revenue_30d: buildRevenueTrend(30),
          order_7d: buildOrderTrend(7),
        },
        status_distribution: statusDistribution,
        hot_products: hotProducts,
        announcements: clone(announcementsSource),
        stock_warning_list: stockWarningList,
      },
    }
  }

  // Shop info
  if (key === 'GET /admin/api/shops/current') {
    return { matched: true, data: clone(shopsCurrentSource) }
  }

  // Announcements
  if (key === 'GET /admin/api/announcements') {
    const page = Math.max(1, Number(query.page || 1))
    const size = Math.max(1, Number(query.size || 20))
    const offset = (page - 1) * size
    const list = announcementsSource.slice()
    return { matched: true, data: { list: list.slice(offset, offset + size), total: list.length, page, size } }
  }

  if (key === 'GET /admin/api/wms/warehouses') {
    const keyword = String(query.keyword || '').trim().toLowerCase()
    const statusQuery = query.status
    const hasStatus = statusQuery !== undefined && statusQuery !== null && String(statusQuery).trim() !== ''
    const status = Number(statusQuery || 0)
    let list = wmsWarehousesSource.slice()
    if (keyword) {
      list = list.filter((row: any) =>
        String(row.name || '').toLowerCase().includes(keyword) ||
        String(row.code || '').toLowerCase().includes(keyword),
      )
    }
    if (hasStatus) {
      list = list.filter((row: any) => Number(row.status || 0) === status)
    }
    list = list.sort((a: any, b: any) => Number(a.id || 0) - Number(b.id || 0))
    return { matched: true, data: toPageData(clone(list), Number(query.page || 1), Number(query.size || 20)) }
  }
  if (key === 'POST /admin/api/wms/warehouses') {
    const payload: any = params || {}
    wmsWarehouseSeq += 1
    const now = wmsNowISO()
    const row = {
      id: wmsWarehouseSeq,
      name: String(payload.name || `仓库${wmsWarehouseSeq}`),
      code: String(payload.code || `WH-${wmsWarehouseSeq}`),
      address: String(payload.address || ''),
      contact: String(payload.contact || ''),
      phone: String(payload.phone || ''),
      status: payload.status === undefined ? 1 : (Number(payload.status || 0) === 1 ? 1 : 0),
      created_at: now,
      updated_at: now,
    }
    wmsWarehousesSource.push(row)
    return { matched: true, data: clone(row) }
  }
  if (methodUpper === 'PUT' && /\/admin\/api\/wms\/warehouses\/\d+$/.test(url)) {
    const id = Number(url.split('/').pop() || 0)
    const payload: any = params || {}
    const target = wmsWarehousesSource.find((row: any) => Number(row.id || 0) === id)
    if (!target) return { matched: true, data: null }
    if (payload.name !== undefined) target.name = String(payload.name || target.name)
    if (payload.code !== undefined) target.code = String(payload.code || target.code || '')
    if (payload.address !== undefined) target.address = String(payload.address || '')
    if (payload.contact !== undefined) target.contact = String(payload.contact || '')
    if (payload.phone !== undefined) target.phone = String(payload.phone || '')
    if (payload.status !== undefined) target.status = Number(payload.status || 0) === 1 ? 1 : 0
    target.updated_at = wmsNowISO()
    for (const stock of wmsStocksSource) {
      if (Number(stock.warehouse_id || 0) === id) stock.warehouse_name = target.name
    }
    for (const docRow of wmsDocsSource) {
      if (Number(docRow.warehouse_id || 0) === id) docRow.warehouse_name = target.name
    }
    return { matched: true, data: null }
  }

  if (key === 'GET /admin/api/wms/stocks') {
    const warehouseID = Number(query.warehouse_id || 0)
    const skuID = Number(query.sku_id || 0)
    const keyword = String(query.keyword || '').trim().toLowerCase()
    let list = wmsStocksSource.slice()
    if (warehouseID > 0) list = list.filter((row: any) => Number(row.warehouse_id || 0) === warehouseID)
    if (skuID > 0) list = list.filter((row: any) => Number(row.sku_id || 0) === skuID)
    if (keyword) {
      list = list.filter((row: any) =>
        String(row.sku_name || '').toLowerCase().includes(keyword) ||
        String(row.sku_id || '').includes(keyword),
      )
    }
    list = list.sort((a: any, b: any) => Number(a.id || 0) - Number(b.id || 0))
    return { matched: true, data: toPageData(clone(list), Number(query.page || 1), Number(query.size || 20)) }
  }
  if (methodUpper === 'PUT' && /\/admin\/api\/wms\/stocks\/\d+\/safety$/.test(url)) {
    const id = Number(url.split('/')[5] || 0)
    const safeQty = Math.max(0, Number((params as any)?.safe_qty || 0))
    const target = wmsStocksSource.find((row: any) => Number(row.id || 0) === id)
    if (!target) return { matched: true, data: null }
    target.safe_qty = safeQty
    target.updated_at = wmsNowISO()
    return { matched: true, data: null }
  }

  if (key === 'GET /admin/api/wms/docs') {
    const docType = String(query.doc_type || '')
    const status = String(query.status || '')
    const warehouseID = Number(query.warehouse_id || 0)
    const docNo = String(query.doc_no || '').trim().toLowerCase()
    let list = wmsDocsSource.slice()
    if (docType) list = list.filter((row: any) => String(row.doc_type || '') === docType)
    if (status) list = list.filter((row: any) => String(row.status || '') === status)
    if (warehouseID > 0) list = list.filter((row: any) => Number(row.warehouse_id || 0) === warehouseID)
    if (docNo) list = list.filter((row: any) => String(row.doc_no || '').toLowerCase().includes(docNo))
    list = list
      .map((row: any) => ({
        ...row,
        warehouse_name: wmsWarehouseName(Number(row.warehouse_id || 0)),
        total_qty: wmsDocTotalQty(row.items),
      }))
      .sort((a: any, b: any) => Number(b.id || 0) - Number(a.id || 0))
    return { matched: true, data: toPageData(clone(list), Number(query.page || 1), Number(query.size || 20)) }
  }
  if (key === 'POST /admin/api/wms/docs') {
    const payload: any = params || {}
    const now = wmsNowISO()
    wmsDocSeq += 1
    const docType = String(payload.doc_type || 'inbound') === 'outbound' ? 'outbound' : 'inbound'
    const warehouseID = Number(payload.warehouse_id || wmsWarehousesSource[0]?.id || 0)
    const items = normalizeWmsItems(Array.isArray(payload.items) ? payload.items : [])
    if (!items.length) {
      throw new Error('单据明细不能为空')
    }
    const row = {
      id: wmsDocSeq,
      doc_no: nextWmsDocNo(docType),
      doc_type: docType,
      status: 'draft',
      warehouse_id: warehouseID,
      warehouse_name: wmsWarehouseName(warehouseID),
      remark: String(payload.remark || ''),
      items,
      total_qty: wmsDocTotalQty(items),
      created_at: now,
      updated_at: now,
    }
    wmsDocsSource.push(row)
    return { matched: true, data: clone(row) }
  }
  if (methodUpper === 'GET' && /\/admin\/api\/wms\/docs\/\d+$/.test(url)) {
    const id = Number(url.split('/').pop() || 0)
    const target = wmsDocsSource.find((row: any) => Number(row.id || 0) === id)
    if (!target) return { matched: true, data: null }
    target.warehouse_name = wmsWarehouseName(Number(target.warehouse_id || 0))
    target.total_qty = wmsDocTotalQty(target.items)
    return {
      matched: true,
      data: {
        doc: clone({
          ...target,
          items: undefined,
        }),
        items: clone(target.items || []),
      },
    }
  }
  if (methodUpper === 'PUT' && /\/admin\/api\/wms\/docs\/\d+$/.test(url)) {
    const id = Number(url.split('/').pop() || 0)
    const payload: any = params || {}
    const target = wmsDocsSource.find((row: any) => Number(row.id || 0) === id)
    if (!target) return { matched: true, data: null }
    if (String(target.status || '') !== 'draft') {
      throw new Error('已完成或已取消单据不可编辑')
    }
    target.doc_type = String(payload.doc_type || target.doc_type || 'inbound') === 'outbound' ? 'outbound' : 'inbound'
    target.warehouse_id = Number(payload.warehouse_id || target.warehouse_id || 0)
    target.warehouse_name = wmsWarehouseName(Number(target.warehouse_id || 0))
    target.remark = String(payload.remark || '')
    if (Array.isArray(payload.items)) {
      target.items = normalizeWmsItems(payload.items)
    }
    target.total_qty = wmsDocTotalQty(target.items)
    target.updated_at = wmsNowISO()
    return { matched: true, data: null }
  }
  if (methodUpper === 'POST' && /\/admin\/api\/wms\/docs\/\d+\/complete$/.test(url)) {
    const id = Number(url.split('/')[5] || 0)
    const target = wmsDocsSource.find((row: any) => Number(row.id || 0) === id)
    if (!target) return { matched: true, data: null }
    if (String(target.status || '') !== 'draft') return { matched: true, data: null }
    if (!Array.isArray(target.items) || !target.items.length) {
      throw new Error('单据明细不能为空')
    }
    applyWmsDocStockChange(target)
    target.status = 'completed'
    target.total_qty = wmsDocTotalQty(target.items)
    target.updated_at = wmsNowISO()
    return { matched: true, data: null }
  }
  if (methodUpper === 'POST' && /\/admin\/api\/wms\/docs\/\d+\/cancel$/.test(url)) {
    const id = Number(url.split('/')[5] || 0)
    const target = wmsDocsSource.find((row: any) => Number(row.id || 0) === id)
    if (!target) return { matched: true, data: null }
    target.status = 'canceled'
    target.updated_at = wmsNowISO()
    return { matched: true, data: null }
  }

  if (key === 'GET /admin/api/wms/movements') {
    const warehouseID = Number(query.warehouse_id || 0)
    const skuID = Number(query.sku_id || 0)
    const bizType = String(query.biz_type || '')
    const docNo = String(query.doc_no || '').trim().toLowerCase()
    let list = wmsMovementsSource.slice()
    if (warehouseID > 0) list = list.filter((row: any) => Number(row.warehouse_id || 0) === warehouseID)
    if (skuID > 0) list = list.filter((row: any) => Number(row.sku_id || 0) === skuID)
    if (bizType) list = list.filter((row: any) => String(row.biz_type || '') === bizType)
    if (docNo) list = list.filter((row: any) => String(row.doc_no || '').toLowerCase().includes(docNo))
    list = list.sort((a: any, b: any) => Number(b.id || 0) - Number(a.id || 0))
    return { matched: true, data: toPageData(clone(list), Number(query.page || 1), Number(query.size || 20)) }
  }

  // Dynamic coupon CRUD
  if (key === 'GET /admin/api/marketing/coupons') {
    const keyword = String(query.keyword || '').trim().toLowerCase()
    const status = query.status !== undefined && query.status !== null && String(query.status).trim() !== '' ? Number(query.status) : null
    const type = query.type !== undefined && query.type !== null && String(query.type).trim() !== '' ? Number(query.type) : null
    let list = clone(couponsSource)
    if (keyword) list = list.filter((c: any) => String(c.name || '').toLowerCase().includes(keyword))
    if (status !== null) list = list.filter((c: any) => Number(c.status || 0) === status)
    if (type !== null) list = list.filter((c: any) => Number(c.type || 0) === type)
    return { matched: true, data: toPageData(list, Number(query.page || 1), Number(query.size || 20)) }
  }
  if (key === 'POST /admin/api/marketing/coupons') {
    couponSeq += 1
    const payload: any = params || {}
    couponsSource.push({
      id: couponSeq,
      name: String(payload.name || ''),
      type: Number(payload.type || 1),
      min_amount: Number(payload.min_amount || 0),
      discount: Number(payload.discount || 0),
      total_count: Number(payload.total_count || 0),
      per_limit: Number(payload.per_limit || 1),
      status: Number(payload.status || 0) === 1 ? 1 : 0,
      used_count: 0,
      start_at: String(payload.start_at || ''),
      end_at: String(payload.end_at || ''),
      description: String(payload.description || ''),
      stack_rule: String(payload.stack_rule || 'exclusive'),
      target_type: String(payload.target_type || 'all'),
      target_value: String(payload.target_value || ''),
    })
    return { matched: true, data: { id: couponSeq } }
  }
  if (methodUpper === 'PUT' && /\/admin\/api\/marketing\/coupons\/\d+$/.test(url)) {
    const id = Number(url.split('/').pop() || 0)
    const target = couponsSource.find((c: any) => Number(c.id || 0) === id)
    if (target) {
      const payload: any = params || {}
      if (payload.name !== undefined) target.name = String(payload.name)
      if (payload.type !== undefined) target.type = Number(payload.type)
      if (payload.min_amount !== undefined) target.min_amount = Number(payload.min_amount)
      if (payload.discount !== undefined) target.discount = Number(payload.discount)
      if (payload.total_count !== undefined) target.total_count = Number(payload.total_count)
      if (payload.per_limit !== undefined) target.per_limit = Number(payload.per_limit)
      if (payload.status !== undefined) target.status = Number(payload.status) === 1 ? 1 : 0
      if (payload.start_at !== undefined) target.start_at = String(payload.start_at)
      if (payload.end_at !== undefined) target.end_at = String(payload.end_at)
      if (payload.description !== undefined) target.description = String(payload.description)
      if (payload.stack_rule !== undefined) target.stack_rule = String(payload.stack_rule)
      if (payload.target_type !== undefined) target.target_type = String(payload.target_type)
      if (payload.target_value !== undefined) target.target_value = String(payload.target_value)
    }
    return { matched: true, data: null }
  }
  if (methodUpper === 'DELETE' && /\/admin\/api\/marketing\/coupons\/\d+$/.test(url)) {
    const id = Number(url.split('/').pop() || 0)
    const idx = couponsSource.findIndex((c: any) => Number(c.id || 0) === id)
    if (idx >= 0) couponsSource.splice(idx, 1)
    return { matched: true, data: null }
  }
  if (methodUpper === 'POST' && /\/admin\/api\/marketing\/coupons\/\d+\/send$/.test(url)) {
    const id = Number(url.split('/')[5] || 0)
    const target = couponsSource.find((c: any) => Number(c.id || 0) === id)
    const count = Number((params as any)?.count || 1)
    if (target) target.used_count = (target.used_count || 0) + count
    return { matched: true, data: { sent_count: count } }
  }

  // Spec templates CRUD
  if (key === 'GET /admin/api/spec-templates') {
    const keyword = String(query.keyword || '').trim().toLowerCase()
    const categoryID = Number(query.category_id || 0)
    let list = clone(specTemplatesSource)
    if (keyword) list = list.filter((t: any) => String(t.name || '').toLowerCase().includes(keyword))
    if (categoryID > 0) list = list.filter((t: any) => Array.isArray(t.category_ids) && t.category_ids.includes(categoryID))
    return { matched: true, data: toPageData(list, Number(query.page || 1), Number(query.size || 20)) }
  }
  if (methodUpper === 'GET' && /\/admin\/api\/spec-templates\/\d+$/.test(url)) {
    const id = Number(url.split('/').pop() || 0)
    const target = specTemplatesSource.find((t: any) => Number(t.id || 0) === id)
    return { matched: true, data: target ? clone(target) : null }
  }
  if (key === 'POST /admin/api/spec-templates') {
    specTemplateSeq += 1
    const payload: any = params || {}
    specTemplatesSource.push({
      id: specTemplateSeq,
      name: String(payload.name || ''),
      category_ids: Array.isArray(payload.category_ids) ? payload.category_ids : [],
      attrs: Array.isArray(payload.attrs) ? payload.attrs : [],
      status: Number(payload.status || 0) === 1 ? 1 : 0,
    })
    return { matched: true, data: { id: specTemplateSeq } }
  }
  if (methodUpper === 'PUT' && /\/admin\/api\/spec-templates\/\d+$/.test(url)) {
    const id = Number(url.split('/').pop() || 0)
    const target = specTemplatesSource.find((t: any) => Number(t.id || 0) === id)
    if (target) {
      const payload: any = params || {}
      if (payload.name !== undefined) target.name = String(payload.name)
      if (payload.category_ids !== undefined) target.category_ids = Array.isArray(payload.category_ids) ? payload.category_ids : target.category_ids
      if (payload.attrs !== undefined) target.attrs = Array.isArray(payload.attrs) ? payload.attrs : target.attrs
      if (payload.status !== undefined) target.status = Number(payload.status) === 1 ? 1 : 0
    }
    return { matched: true, data: null }
  }
  if (methodUpper === 'DELETE' && /\/admin\/api\/spec-templates\/\d+$/.test(url)) {
    const id = Number(url.split('/').pop() || 0)
    const idx = specTemplatesSource.findIndex((t: any) => Number(t.id || 0) === id)
    if (idx >= 0) specTemplatesSource.splice(idx, 1)
    return { matched: true, data: null }
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

  // Category tree helper
  if (key === 'GET /admin/api/categories/tree') {
    const list = clone(categoriesSource) as any[]
    const map = new Map<number, any>()
    for (const c of list) { c.children = []; map.set(Number(c.id), c) }
    const roots: any[] = []
    for (const c of list) {
      const parent = Number(c.parent_id || 0)
      if (parent && map.get(parent)) map.get(parent).children.push(c)
      else roots.push(c)
    }
    return { matched: true, data: roots }
  }

  if (key === 'GET /admin/api/categories') {
    return { matched: true, data: clone(categoriesSource) }
  }
  if (key === 'POST /admin/api/categories') {
    const payload: any = params || {}
    categorySeq += 1
    const now = new Date().toISOString()
    const row = {
      id: categorySeq,
      parent_id: Number(payload.parent_id || 0),
      name: String(payload.name || `分类${categorySeq}`),
      icon: String(payload.icon || ''),
      sort: Number(payload.sort || 0),
      status: Number(payload.status || 0) === 1 ? 1 : 0,
      product_count: 0,
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

  // Enhanced product query
  if (key === 'GET /admin/api/products') {
    const keyword = String(query.keyword || '').trim().toLowerCase()
    const status = query.status === undefined || query.status === '' ? null : Number(query.status)
    const categoryID = Number(query.category_id || 0)
    const lowStock = query.low_stock === true || String(query.low_stock || '') === 'true' || query.low_stock === 1
    const sortBy = String(query.sort_by || '')

    let list = clone(productListSource) as any[]
    if (keyword) list = list.filter((p: any) => String(p.title || '').toLowerCase().includes(keyword))
    if (status !== null) list = list.filter((p: any) => Number(p.status || 0) === status)
    if (categoryID) list = list.filter((p: any) => Number(p.category_id || 0) === categoryID)
    if (lowStock) list = list.filter((p: any) => Number(p.stock || 0) <= Number(p.low_stock_threshold || 10))
    if (sortBy === 'sales') list.sort((a: any, b: any) => Number(b.sales_count || 0) - Number(a.sales_count || 0))
    else if (sortBy === 'stock') list.sort((a: any, b: any) => Number(a.stock || 0) - Number(b.stock || 0))
    else if (sortBy === 'price_asc') list.sort((a: any, b: any) => Number(a.price || 0) - Number(b.price || 0))
    else if (sortBy === 'price_desc') list.sort((a: any, b: any) => Number(b.price || 0) - Number(a.price || 0))
    else if (sortBy === 'created') list.sort((a: any, b: any) => String(b.created_at || '').localeCompare(String(a.created_at || '')))

    return { matched: true, data: toPageData(list, Number(query.page || 1), Number(query.size || 20)) }
  }

  // Batch product routes
  if (methodUpper === 'PUT' && url === '/admin/api/products/batch/status') {
    const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
    const status = Number(params?.status || 0) === 1 ? 1 : 0
    const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
    for (const pid of ids) {
      const t = (productListSource as any[]).find((p: any) => Number(p.id || 0) === Number(pid))
      if (!t) { fail.push({ id: Number(pid), reason: '商品不存在' }); continue }
      t.status = status; success_ids.push(Number(pid))
    }
    return { matched: true, data: { success_ids, fail } }
  }

  if (methodUpper === 'PUT' && url === '/admin/api/products/batch/category') {
    const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
    const cid = Number(params?.category_id || 0)
    const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
    for (const pid of ids) {
      const t = (productListSource as any[]).find((p: any) => Number(p.id || 0) === Number(pid))
      if (!t) { fail.push({ id: Number(pid), reason: '商品不存在' }); continue }
      t.category_id = cid; success_ids.push(Number(pid))
    }
    return { matched: true, data: { success_ids, fail } }
  }

  if (methodUpper === 'PUT' && url === '/admin/api/products/batch/price') {
    const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
    const adj = params?.adjustment || { type: 'set', value: 0 }
    const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
    for (const pid of ids) {
      const t = (productListSource as any[]).find((p: any) => Number(p.id || 0) === Number(pid))
      if (!t) { fail.push({ id: Number(pid), reason: '商品不存在' }); continue }
      const orig = Number(t.price || 0)
      let next = orig
      if (adj.type === 'set') next = Number(adj.value || 0)
      else if (adj.type === 'percent') next = orig * (1 + Number(adj.value || 0) / 100)
      else if (adj.type === 'amount') next = orig + Number(adj.value || 0)
      t.price = Math.max(0, Math.round(next * 100) / 100); success_ids.push(Number(pid))
    }
    return { matched: true, data: { success_ids, fail } }
  }

  if (key === 'GET /admin/api/orders') {
    const keyword = String(query.keyword || '').trim().toLowerCase()
    const status = String(query.status || '')
    const amountMin = Number(query.amount_min || 0)
    const amountMax = Number(query.amount_max || 0)
    const company = String(query.logistics_company || '')
    const province = String(query.province || '')
    const payMethod = String(query.pay_method || '')
    const hasAfterSale = query.has_after_sale === true || String(query.has_after_sale || '') === 'true' || query.has_after_sale === 1
    const timeStart = String(query.time_start || '')
    const timeEnd = String(query.time_end || '')

    let list = clone(orderListSource)
    if (status) list = list.filter((o: any) => String(o.status) === status)
    if (keyword) {
      list = list.filter((o: any) =>
        String(o.id || '').includes(keyword) ||
        String(o.user_nickname || '').toLowerCase().includes(keyword) ||
        String(o.receiver_name || '').toLowerCase().includes(keyword) ||
        (Array.isArray(o.items) && o.items.some((it: any) => String(it.title || '').toLowerCase().includes(keyword))))
    }
    if (amountMin > 0) list = list.filter((o: any) => Number(o.pay_amount || o.total_amount || 0) >= amountMin)
    if (amountMax > 0) list = list.filter((o: any) => Number(o.pay_amount || o.total_amount || 0) <= amountMax)
    if (company) list = list.filter((o: any) => Array.isArray(o.shipments) && o.shipments.some((s: any) => String(s.company || '') === company))
    if (province) list = list.filter((o: any) => String(o.receiver_address || '').includes(province))
    if (payMethod) list = list.filter((o: any) => String(o.pay_method || '') === payMethod)
    if (hasAfterSale) list = list.filter((o: any) => o.has_after_sale === true || (Array.isArray(o.after_sales) && o.after_sales.length > 0))
    if (timeStart) list = list.filter((o: any) => String(o.created_at || '') >= timeStart)
    if (timeEnd) list = list.filter((o: any) => String(o.created_at || '') <= timeEnd)

    return { matched: true, data: toPageData(list, Number(query.page || 1), Number(query.size || 20)) }
  }

  // Order action routes
  if (methodUpper === 'POST' && /\/admin\/api\/orders\/\d+\/repricing$/.test(url)) {
    const id = Number(url.split('/')[4] || 0)
    const target = orderListSource.find((o: any) => Number(o.id || 0) === id)
    if (!target) return { matched: true, data: null }
    const items: any[] = Array.isArray(params?.items) ? params.items : []
    let goodsTotal = 0
    for (const it of items) {
      const oi = (target.items || []).find((x: any) => Number(x.id || x.item_id || 0) === Number(it.item_id || 0))
      if (oi) { oi.price = Number(it.price || 0); goodsTotal += Number(it.price || 0) * Number(oi.qty || 1) }
    }
    if (!goodsTotal) goodsTotal = Number(target.goods_amount || target.total_amount || 0)
    target.amount_breakdown = {
      goods_amount: goodsTotal,
      discount_amount: Number(target.discount_amount || 0),
      payable_amount: Math.max(0, goodsTotal - Number(target.discount_amount || 0)),
    }
    target.pay_amount = target.amount_breakdown.payable_amount
    target.notes = target.notes || []
    target.notes.push({ id: Date.now(), content: `改价：${params?.remark || ''}`, created_at: new Date().toISOString() })
    return { matched: true, data: { id, amount_breakdown: clone(target.amount_breakdown) } }
  }

  if (methodUpper === 'POST' && /\/admin\/api\/orders\/\d+\/notes$/.test(url)) {
    const id = Number(url.split('/')[4] || 0)
    const target = orderListSource.find((o: any) => Number(o.id || 0) === id)
    if (!target) return { matched: true, data: null }
    target.notes = target.notes || []
    target.notes.push({ id: Date.now(), content: String(params?.content || ''), created_at: new Date().toISOString(), visible_to: String(params?.visible_to || 'merchant_only') })
    return { matched: true, data: { id, notes: clone(target.notes) } }
  }

  if (methodUpper === 'POST' && /\/admin\/api\/orders\/\d+\/remind-pay$/.test(url)) {
    return { matched: true, data: { sent_at: new Date().toISOString(), channel: String(params?.channel || 'sms') } }
  }

  if (methodUpper === 'GET' && /\/admin\/api\/orders\/\d+\/print-template$/.test(url)) {
    const id = Number(url.split('/')[4] || 0)
    const target = orderListSource.find((o: any) => Number(o.id || 0) === id)
    const name = target?.receiver_name || '--'
    const addr = target?.receiver_address || '--'
    const tpl = `<div style="font-family:sans-serif;padding:12px;"><h3>电子面单 #${id}</h3><div>收件人：${name}</div><div>地址：${addr}</div><div>商品：${(target?.items || []).map((x: any) => x.title).join('，')}</div><div style="margin-top:10px;font-weight:bold;">[ 顺丰速运 ] SF${1000000 + id}</div></div>`
    return { matched: true, data: { template: tpl } }
  }

  if (methodUpper === 'GET' && /\/admin\/api\/orders\/\d+\/timeline$/.test(url)) {
    const id = Number(url.split('/')[4] || 0)
    const target = orderListSource.find((o: any) => Number(o.id || 0) === id)
    const createdAt = String(target?.created_at || '2026-05-25T08:00:00Z')
    const stages = [
      { stage: 'created', status: '已下单', time: createdAt, content: '买家下单' },
      { stage: 'paid', status: '已支付', time: createdAt.replace('08:00:00', '08:05:00'), content: '微信支付完成' },
      { stage: 'shipped', status: '已发货', time: createdAt.replace('08:00:00', '12:00:00'), content: '顺丰已揽收' },
      { stage: 'received', status: '已签收', time: createdAt.replace('08:00:00', '23:00:00'), content: '客户已签收' },
      { stage: 'completed', status: '已完成', time: createdAt.replace('08:00:00', '23:30:00'), content: '订单完成' },
    ]
    const currentStatus = String(target?.status || '1')
    const cutoff: Record<string, number> = { '1': 1, '2': 2, '3': 3, '4': 5, '5': 1 }
    return { matched: true, data: stages.slice(0, cutoff[currentStatus] || stages.length) }
  }

  if (methodUpper === 'POST' && url === '/admin/api/orders/batch/ship') {
    const rows: any[] = Array.isArray(params) ? params : []
    const success_ids: number[] = []
    const fail: Array<{ id: number; reason: string }> = []
    for (const row of rows) {
      const oid = Number(row?.order_id || 0)
      const target = orderListSource.find((o: any) => Number(o.id || 0) === oid)
      if (!target) { fail.push({ id: oid, reason: '订单不存在' }); continue }
      if (String(target.status) === '5') { fail.push({ id: oid, reason: '订单已关闭' }); continue }
      target.shipments = target.shipments || []
      target.shipments.push({ id: Date.now() + oid, company: String(row.company || 'SF'), tracking_no: String(row.tracking_no || ''), delivery_type: 'express', logistics_status: 'created', logistics_status_label: '已下单', created_at: new Date().toISOString() })
      target.status = '3'; target.status_label = '已发货'
      success_ids.push(oid)
    }
    return { matched: true, data: { success_ids, fail } }
  }

  if (methodUpper === 'POST' && url === '/admin/api/orders/batch/notes') {
    const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
    const content = String(params?.content || '')
    const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
    for (const oid of ids) {
      const t = orderListSource.find((o: any) => Number(o.id || 0) === Number(oid))
      if (!t) { fail.push({ id: Number(oid), reason: '订单不存在' }); continue }
      t.notes = t.notes || []; t.notes.push({ id: Date.now(), content, created_at: new Date().toISOString() })
      success_ids.push(Number(oid))
    }
    return { matched: true, data: { success_ids, fail } }
  }

  if (methodUpper === 'POST' && url === '/admin/api/orders/batch/repricing') {
    const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
    const adj = params?.adjustment || { type: 'percent', value: -10 }
    const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
    for (const oid of ids) {
      const t = orderListSource.find((o: any) => Number(o.id || 0) === Number(oid))
      if (!t) { fail.push({ id: Number(oid), reason: '订单不存在' }); continue }
      if (String(t.status) !== '1') { fail.push({ id: Number(oid), reason: '当前状态不可改价' }); continue }
      const orig = Number(t.pay_amount || t.total_amount || 0)
      const next = adj.type === 'percent' ? orig * (1 + Number(adj.value || 0) / 100) : orig + Number(adj.value || 0)
      t.pay_amount = Math.max(0, Math.round(next * 100) / 100)
      t.amount_breakdown = { goods_amount: t.pay_amount, discount_amount: 0, payable_amount: t.pay_amount }
      success_ids.push(Number(oid))
    }
    return { matched: true, data: { success_ids, fail } }
  }

  if (methodUpper === 'POST' && url === '/admin/api/orders/batch/close') {
    const ids: number[] = Array.isArray(params?.ids) ? params.ids : []
    const reason = String(params?.reason || '')
    const success_ids: number[] = []; const fail: Array<{ id: number; reason: string }> = []
    for (const oid of ids) {
      const t = orderListSource.find((o: any) => Number(o.id || 0) === Number(oid))
      if (!t) { fail.push({ id: Number(oid), reason: '订单不存在' }); continue }
      if (String(t.status) === '4') { fail.push({ id: Number(oid), reason: '已完成订单不能关闭' }); continue }
      t.status = '5'; t.status_label = '已关闭'; t.close_reason = reason
      success_ids.push(Number(oid))
    }
    return { matched: true, data: { success_ids, fail } }
  }
  if (key === 'GET /admin/api/reviews') {
    return { matched: true, data: listReviews(query) }
  }
  if (key === 'GET /admin/api/after-sales') {
    const status = String(query.status || '').trim()
    const type = String(query.type || '')
    const orderID = Number(query.order_id || 0)
    const timeStart = String(query.time_start || '')
    const timeEnd = String(query.time_end || '')
    const page = Math.max(1, Number(query.page || 1))
    const size = Math.max(1, Number(query.size || 20))
    let list = afterSalesSource.slice()
    if (status) list = list.filter((item: any) => String(item.status || '') === status)
    if (type) list = list.filter((item: any) => String(item.case_type || '') === type)
    if (orderID > 0) list = list.filter((item: any) => Number(item.order_id || 0) === orderID)
    if (timeStart) list = list.filter((item: any) => String(item.created_at || '') >= timeStart)
    if (timeEnd) list = list.filter((item: any) => String(item.created_at || '') <= timeEnd)
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
    // Messages route
    if (methodUpper === 'POST' && /\/admin\/api\/after-sales\/\d+\/messages$/.test(url)) {
      const id = Number(url.split('/')[4] || 0)
      let target = afterSalesSource.find((row: any) => Number(row.id || 0) === id)
      if (!target) {
        target = { id, status: 'applied', case_type: 'refund_only', order_id: 1, messages: [], evidences: [], logs: [], shipments: [], created_at: new Date().toISOString() }
        afterSalesSource.push(target)
      }
      target.messages = target.messages || []
      target.messages.push({ id: Date.now(), from: String(params?.from || 'merchant'), content: String(params?.content || ''), images: Array.isArray(params?.images) ? params.images : [], created_at: new Date().toISOString() })
      return { matched: true, data: { id, messages: clone(target.messages) } }
    }

    // Evidences route
    if (methodUpper === 'POST' && /\/admin\/api\/after-sales\/\d+\/evidences$/.test(url)) {
      const id = Number(url.split('/')[4] || 0)
      const target = afterSalesSource.find((row: any) => Number(row.id || 0) === id)
      if (!target) return { matched: true, data: null }
      target.evidences = target.evidences || []
      target.evidences.push({ id: Date.now(), images: Array.isArray(params?.images) ? params.images : [], remark: String(params?.remark || ''), created_at: new Date().toISOString() })
      return { matched: true, data: { id, evidences: clone(target.evidences) } }
    }
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
