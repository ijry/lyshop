import type { MarketingKind } from '@/api/marketing'

/**
 * 营销活动域元数据与展示逻辑的唯一事实来源。
 *
 * 页面 / 工作台 / 组件都从这里取 kind 对应的标题、图标、路由与价格语义，
 * 避免「工作台一份路由表、营销中心一份路由表」的双份维护。
 */

export type { MarketingKind }

export interface MarketingMeta {
  kind: MarketingKind
  title: string
  eyebrow: string
  desc: string
  iconName: string
  route: string
  /** 主价格字段与标签（秒杀/拼团为活动价，砍价为起砍价） */
  priceField: 'activity_price' | 'start_price'
  priceLabel: string
  /** 次价格字段与标签（砍价显示底价，其余显示原价） */
  secondaryField?: string
  secondaryLabel?: string
  /** 价格输入下方的规则提示 */
  priceHint: string
}

export const ACTIVITY_META: Record<MarketingKind, MarketingMeta> = {
  seckill: {
    kind: 'seckill',
    title: '秒杀活动',
    eyebrow: 'FLASH SALE',
    desc: '限时特价、库存限量与限购',
    iconName: 'flash',
    route: '/pages/marketing/seckill',
    priceField: 'activity_price',
    priceLabel: '秒杀价',
    secondaryLabel: '原价',
    priceHint: '秒杀价需低于商品原价，0 表示不限制',
  },
  'group-buy': {
    kind: 'group-buy',
    title: '拼团活动',
    eyebrow: 'GROUP BUY',
    desc: '成团价、限购与库存上限',
    iconName: 'users',
    route: '/pages/marketing/group-buy',
    priceField: 'activity_price',
    priceLabel: '拼团价',
    secondaryLabel: '原价',
    priceHint: '拼团价需低于商品原价，0 表示不限制',
  },
  bargain: {
    kind: 'bargain',
    title: '砍价活动',
    eyebrow: 'BARGAIN',
    desc: '起砍价、底价与砍价库存',
    iconName: 'money',
    route: '/pages/marketing/bargain',
    priceField: 'start_price',
    priceLabel: '起砍价',
    secondaryField: 'floor_price',
    secondaryLabel: '底价',
    priceHint: '底价不得高于起砍价，且均需大于 0',
  },
}

export const ACTIVITY_KINDS: MarketingKind[] = ['seckill', 'group-buy', 'bargain']

/** 非活动类营销入口路由，工作台与营销中心共用 */
export const MARKETING_ROUTES: Record<string, string> = {
  center: '/pages/marketing/index',
  coupon: '/pages/marketing/coupon',
  seckill: '/pages/marketing/seckill',
  groupBuy: '/pages/marketing/group-buy',
  bargain: '/pages/marketing/bargain',
  vip: '/pages/vip/plan-list',
  messageBroadcast: '/pages/me/messages',
  decor: '/pages/decor/editor',
  specTemplate: '/pages/product/spec-templates',
  distribution: '/pages/distribution/index',
  pointsMall: '/pages/points/index',
}

export const ACTIVITY_DETAIL_ROUTE = '/pages/marketing/activity-detail'

export function isActivityKind(value: unknown): value is MarketingKind {
  return ACTIVITY_KINDS.includes(String(value) as MarketingKind)
}

export function activityMeta(kind: unknown): MarketingMeta {
  return isActivityKind(kind) ? ACTIVITY_META[kind] : ACTIVITY_META.seckill
}

export function activityDetailUrl(kind: MarketingKind, id: number | string, name?: string) {
  const title = encodeURIComponent(String(name || activityMeta(kind).title))
  return `${ACTIVITY_DETAIL_ROUTE}?kind=${kind}&id=${id}&name=${title}`
}

/* -------------------------------------------------------------------------- */
/* 状态模型                                                                    */
/* -------------------------------------------------------------------------- */

export type ActivityPhase = 'running' | 'upcoming' | 'ended' | 'disabled'

export interface ActivityState {
  phase: ActivityPhase
  label: string
  /** StatusTag 的 type */
  tone: string
}

export const ACTIVITY_FILTERS: Array<{ key: 'all' | ActivityPhase; label: string }> = [
  { key: 'all', label: '全部' },
  { key: 'running', label: '进行中' },
  { key: 'upcoming', label: '未开始' },
  { key: 'ended', label: '已结束' },
  { key: 'disabled', label: '已停用' },
]

/** 状态由「启用标记 + 当前时间」共同决定，前端不依赖后端下发的状态文案 */
export function activityState(activity: any, now = Date.now()): ActivityState {
  if (Number(activity?.status ?? 0) !== 1) {
    return { phase: 'disabled', label: '已停用', tone: 'neutral' }
  }
  const start = new Date(String(activity?.start_at || '')).getTime()
  const end = new Date(String(activity?.end_at || '')).getTime()
  if (!start || !end || Number.isNaN(start) || Number.isNaN(end)) {
    return { phase: 'upcoming', label: '待排期', tone: 'neutral' }
  }
  if (now < start) return { phase: 'upcoming', label: '未开始', tone: 'primary' }
  if (now > end) return { phase: 'ended', label: '已结束', tone: 'neutral' }
  return { phase: 'running', label: '进行中', tone: 'success' }
}

/** 活动周期进度，用于活动卡与详情页的细进度条 */
export function activityProgress(activity: any, now = Date.now()): number {
  const start = new Date(String(activity?.start_at || '')).getTime()
  const end = new Date(String(activity?.end_at || '')).getTime()
  if (!start || !end || end <= start) return 0
  if (now <= start) return 0
  if (now >= end) return 100
  return Math.min(100, Math.max(0, Math.round(((now - start) / (end - start)) * 100)))
}

/* -------------------------------------------------------------------------- */
/* 展示格式化                                                                  */
/* -------------------------------------------------------------------------- */

function pad(value: number) {
  return String(value).padStart(2, '0')
}

/** ISO / 本地字符串 → `MM-DD HH:mm`（withYear 时 `YYYY-MM-DD HH:mm`） */
export function formatActivityTime(value?: string, withYear = false): string {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return String(value).slice(0, 16).replace('T', ' ')
  const y = date.getFullYear()
  const m = pad(date.getMonth() + 1)
  const d = pad(date.getDate())
  const time = `${pad(date.getHours())}:${pad(date.getMinutes())}`
  return withYear ? `${y}-${m}-${d} ${time}` : `${m}-${d} ${time}`
}

export function activityTimeRange(activity: any, withYear = false): string {
  return `${formatActivityTime(activity?.start_at, withYear)} — ${formatActivityTime(activity?.end_at, withYear)}`
}

export function formatMoney(value: unknown): string {
  return `¥${Number(value || 0).toFixed(2)}`
}

/** SKU 属性可能是 JSON 字符串或数组，统一渲染为「颜色:黑 / 尺码:M」 */
export function parseSkuLabel(attrs: any, skuID = 0): string {
  try {
    const list = typeof attrs === 'string' ? JSON.parse(attrs) : attrs
    if (Array.isArray(list) && list.length) {
      return list.map((item: any) => `${item.name}:${item.value}`).join(' / ')
    }
  } catch {
    // 非 JSON 时退回 SKU 编号
  }
  return skuID ? `SKU #${skuID}` : '默认规格'
}
