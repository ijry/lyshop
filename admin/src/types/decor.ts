// ---- Per-component prop interfaces ----

export interface BannerProps {
  images: Array<{ url: string; link: string }>
  height: number
}

export interface CategoryNavProps {
  items: Array<{ title: string; icon: string; link: string }>
}

export interface ProductGridProps {
  source: string
  limit: number
  columns: number
  title?: string
}

export interface NoticeProps {
  items: Array<{ text: string; link: string }>
  color: string
  bgColor: string
  duration: number
  mode: string
}

export interface ImageAdProps {
  url: string
  link: string
}

export interface RichTextProps {
  content: string
}

export interface MarketingZoneProps {
  type: string
}

export interface SpacerProps {
  height: number
  background: string
}

// ---- Discriminated union ----

interface DecorBase {
  id: string
}

export type DecorComponent =
  | (DecorBase & { type: 'banner'; props: BannerProps })
  | (DecorBase & { type: 'category_nav'; props: CategoryNavProps })
  | (DecorBase & { type: 'product_grid'; props: ProductGridProps })
  | (DecorBase & { type: 'notice'; props: NoticeProps })
  | (DecorBase & { type: 'image_ad'; props: ImageAdProps })
  | (DecorBase & { type: 'rich_text'; props: RichTextProps })
  | (DecorBase & { type: 'marketing_zone'; props: MarketingZoneProps })
  | (DecorBase & { type: 'spacer'; props: SpacerProps })

export type DecorComponentType = DecorComponent['type']

// ---- Component library metadata ----

export const componentLib: Array<{ type: DecorComponentType; title: string; icon: string }> = [
  { type: 'banner',          title: '轮播图',   icon: '🖼' },
  { type: 'category_nav',    title: '分类导航', icon: '📂' },
  { type: 'product_grid',    title: '商品列表', icon: '🛍' },
  { type: 'notice',          title: '公告栏',   icon: '📢' },
  { type: 'image_ad',        title: '广告图',   icon: '🎯' },
  { type: 'rich_text',       title: '富文本',   icon: '📝' },
  { type: 'marketing_zone',  title: '营销区块', icon: '🏷' },
  { type: 'spacer',          title: '间距',     icon: '↕' },
]

export const compTitleMap: Record<string, string> = Object.fromEntries(
  componentLib.map(c => [c.type, c.title])
)

// ---- Default props factory ----

export function createDefaultProps(type: DecorComponentType): DecorComponent['props'] {
  const defaults: Record<DecorComponentType, () => DecorComponent['props']> = {
    banner: () => ({ images: [], height: 350 }),
    category_nav: () => ({ items: [] }),
    product_grid: () => ({ source: 'hot', limit: 10, columns: 2 }),
    notice: () => ({
      items: [
        { text: '欢迎来到 LYShop', link: '/pages/index/index' },
        { text: '新人下单立减，优惠券限时领取', link: '/pages/marketing/coupon' },
        { text: '精选好物每日上新，支持多端下单', link: '/pages/product/list' },
      ],
      color: '#f97316',
      bgColor: '#fff7ed',
      duration: 2500,
      mode: 'link',
    }),
    image_ad: () => ({ url: '', link: '' }),
    rich_text: () => ({ content: '' }),
    marketing_zone: () => ({ type: 'seckill' }),
    spacer: () => ({ height: 16, background: '#ffffff' }),
  }
  return defaults[type]()
}

// ---- Preview message protocol ----

export interface DecorPreviewUpdate {
  type: 'DECOR_PREVIEW_UPDATE'
  source: 'lyshop-admin'
  components: DecorComponent[]
}

export interface DecorPreviewReady {
  type: 'DECOR_PREVIEW_READY'
  source: 'lyshop-app'
}
