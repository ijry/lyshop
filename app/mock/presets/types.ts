export interface SiteSettings {
  site_name: string
  site_logo: string
  seo_title: string
  seo_keywords: string
  seo_description: string
  icp: string
  hero_badge: string
  hero_title: string
  hero_subtitle: string
  hero_btn_text: string
  hero_btn_link: string
  color_primary: string
  color_primary_light: string
  color_primary_dark: string
  color_bg_page: string
  color_bg_header: string
  color_bg_footer: string
  color_price: string
  color_hero_from: string
  color_hero_to: string
}

export interface PcDecorPagePayload {
  pageStyle: {
    background: {
      mode: 'solid' | 'gradient' | 'image'
      solidColor?: string
      gradient?: {
        angle: number
        stops: Array<{ color: string; position: number }>
      }
      image?: {
        url: string
        size: 'cover' | 'contain' | 'auto' | 'custom'
        customSize?: string
        position: string
        repeat: 'no-repeat' | 'repeat' | 'repeat-x' | 'repeat-y'
        attachment: 'scroll' | 'fixed'
      }
      overlay?: {
        enabled: boolean
        color: string
        opacity: number
      }
    }
    content: {
      maxWidth: number
      gutterX: number
      sectionGap: number
    }
    surface: {
      radius: number
      shadow: 'none' | 'sm' | 'md' | 'lg'
    }
  }
  components: any[]
}

export interface MockPreset {
  key: string
  name: string
  categories: any[]
  products: { list: any[]; total: number; page: number; size: number }
  productDetail: any
  indexDecor: { components: any[] }
  pcDecor: PcDecorPagePayload
  siteSettings: SiteSettings
  seckills: any
  groupBuy: any
  bargain: any
  recommend: any[]
  cart: any[]
  orders: { list: any[]; total: number; page: number; size: number }
}
