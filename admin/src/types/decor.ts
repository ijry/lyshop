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

// ---- PC-specific prop interfaces ----

export interface HeroProps {
  badge: string
  title: string
  subtitle: string
  btn_text: string
  btn_link: string
  btn2_text: string
  btn2_link: string
  bg_from: string
  bg_to: string
}

export interface FeaturesProps {
  columns: number
  items: Array<{ icon: string; title: string; desc: string }>
}

export interface GridProps {
  columns: number
  items: Array<{ title: string; icon: string; bg: string; link: string }>
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

export type PcDecorComponent =
  | DecorComponent
  | (DecorBase & { type: 'hero'; props: HeroProps })
  | (DecorBase & { type: 'features'; props: FeaturesProps })
  | (DecorBase & { type: 'grid'; props: GridProps })

export type PcDecorComponentType = PcDecorComponent['type']

// ---- Component library metadata ----

export const componentLib: Array<{ type: DecorComponentType; titleKey: string; icon: string }> = [
  { type: 'banner',          titleKey: 'decor.type.banner',         icon: '🖼' },
  { type: 'category_nav',    titleKey: 'decor.type.categoryNav',    icon: '📂' },
  { type: 'product_grid',    titleKey: 'decor.type.productGrid',    icon: '🛍' },
  { type: 'notice',          titleKey: 'decor.type.notice',         icon: '📢' },
  { type: 'image_ad',        titleKey: 'decor.type.imageAd',        icon: '🎯' },
  { type: 'rich_text',       titleKey: 'decor.type.richText',       icon: '📝' },
  { type: 'marketing_zone',  titleKey: 'decor.type.marketingZone',  icon: '🏷' },
  { type: 'spacer',          titleKey: 'decor.type.spacer',         icon: '↕' },
]

export const compTitleKeyMap: Record<string, string> = Object.fromEntries(
  componentLib.map(c => [c.type, c.titleKey])
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

// ---- PC component library ----

export const pcComponentLib: Array<{ type: PcDecorComponentType; title: string; icon: string }> = [
  { type: 'hero',           title: 'Hero 横幅',   icon: '🎯' },
  { type: 'banner',         title: '轮播图',      icon: '🖼' },
  { type: 'category_nav',   title: '分类导航',    icon: '📂' },
  { type: 'grid',           title: '快捷入口',    icon: '⊞' },
  { type: 'product_grid',   title: '商品列表',    icon: '🛍' },
  { type: 'notice',         title: '公告栏',      icon: '📢' },
  { type: 'image_ad',       title: '广告图',      icon: '🎯' },
  { type: 'rich_text',      title: '富文本',      icon: '📝' },
  { type: 'marketing_zone', title: '营销专区',    icon: '🏷' },
  { type: 'features',       title: '特性栏',      icon: '✅' },
  { type: 'spacer',         title: '间距',        icon: '↕' },
]

export const pcCompTitleMap: Record<string, string> = Object.fromEntries(
  pcComponentLib.map(c => [c.type, c.title])
)

export function createPcDefaultProps(type: PcDecorComponentType): any {
  const defaults: Record<string, () => any> = {
    hero: () => ({ badge: '', title: '', subtitle: '', btn_text: '立即选购', btn_link: '/products', btn2_text: '', btn2_link: '', bg_from: '#b91c1c', bg_to: '#991b1b' }),
    banner: () => ({ images: [], height: 400 }),
    category_nav: () => ({ style: 'floating', columns: 8, items: [] }),
    grid: () => ({ columns: 4, items: [] }),
    product_grid: () => ({ source: 'hot', limit: 8, columns: 4, title: '推荐商品' }),
    notice: () => ({ items: [{ text: '欢迎光临', link: '' }], color: '#f97316', bgColor: '#fff7ed' }),
    image_ad: () => ({ url: '', link: '', height: 200 }),
    rich_text: () => ({ content: '' }),
    marketing_zone: () => ({ title: '限时秒杀', subtitle: '限时抢购中', bg_from: '#b91c1c', bg_to: '#dc2626', more_link: '/products', products: [] }),
    features: () => ({ columns: 4, items: [{ icon: 'i-carbon-star', title: '特性', desc: '描述' }] }),
    spacer: () => ({ height: 24, background: 'transparent' }),
  }
  return (defaults[type] || (() => ({})))()
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
