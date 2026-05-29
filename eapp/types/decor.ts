export interface DecorComponentDef {
  type: string
  title: string
  icon: string
}

export const componentLib: DecorComponentDef[] = [
  { type: 'banner', title: '轮播图', icon: '🖼' },
  { type: 'category_nav', title: '分类导航', icon: '☰' },
  { type: 'product_grid', title: '商品网格', icon: '▦' },
  { type: 'notice', title: '公告', icon: '📢' },
  { type: 'image_ad', title: '图片广告', icon: '🏷' },
  { type: 'rich_text', title: '富文本', icon: '📝' },
  { type: 'marketing_zone', title: '营销区', icon: '🎯' },
  { type: 'spacer', title: '间距', icon: '↕' },
]

export interface BannerProps { images: Array<{ url: string }>; height: number }
export interface CategoryNavProps { category_ids: number[]; style: string }
export interface ProductGridProps { source: string; limit: number; columns: number }
export interface NoticeProps { text: string; color: string }
export interface ImageAdProps { image_url: string; link: string }
export interface RichTextProps { content: string }
export interface MarketingZoneProps { title: string; type: string }
export interface SpacerProps { height: number }

export interface DecorComponent {
  type: string
  id: string
  props: any
}

export function createDefaultProps(type: string): any {
  switch (type) {
    case 'banner': return { images: [], height: 350 }
    case 'category_nav': return { category_ids: [], style: 'grid' }
    case 'product_grid': return { source: 'hot', limit: 6, columns: 2 }
    case 'notice': return { text: '', color: '#f59e0b' }
    case 'image_ad': return { image_url: '', link: '' }
    case 'rich_text': return { content: '' }
    case 'marketing_zone': return { title: '', type: 'seckill' }
    case 'spacer': return { height: 20 }
    default: return {}
  }
}
