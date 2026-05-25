import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { get } from '@/api/request'
import { i18n } from '@/locales'

const t = (key: string) => i18n.global.t(key)

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

function getDefaults(): SiteSettings {
  return {
    site_name: 'LYShop',
    site_logo: '',
    seo_title: t('site.title'),
    seo_keywords: t('site.keywords'),
    seo_description: t('site.description'),
    icp: '',
    hero_badge: t('site.flashSale'),
    hero_title: t('site.heroTitle'),
    hero_subtitle: t('site.heroSubtitle'),
    hero_btn_text: t('site.heroCta'),
    hero_btn_link: '/products',
    color_primary: '#dc2626',
    color_primary_light: '#ef4444',
    color_primary_dark: '#b91c1c',
    color_bg_page: '#f9fafb',
    color_bg_header: 'rgba(255,255,255,0.8)',
    color_bg_footer: '#f9fafb',
    color_price: '#ef4444',
    color_hero_from: '#b91c1c',
    color_hero_to: '#991b1b',
  }
}

function applyTheme(s: SiteSettings) {
  const root = document.documentElement.style
  root.setProperty('--color-primary', s.color_primary)
  root.setProperty('--color-primary-light', s.color_primary_light)
  root.setProperty('--color-primary-dark', s.color_primary_dark)
  root.setProperty('--color-bg-page', s.color_bg_page)
  root.setProperty('--color-bg-header', s.color_bg_header)
  root.setProperty('--color-bg-footer', s.color_bg_footer)
  root.setProperty('--color-price', s.color_price)
  root.setProperty('--color-hero-from', s.color_hero_from)
  root.setProperty('--color-hero-to', s.color_hero_to)

  // SEO
  document.title = s.seo_title || s.site_name
  const metaKw = document.querySelector('meta[name="keywords"]')
  if (metaKw) metaKw.setAttribute('content', s.seo_keywords)
  else if (s.seo_keywords) {
    const m = document.createElement('meta')
    m.name = 'keywords'
    m.content = s.seo_keywords
    document.head.appendChild(m)
  }
  const metaDesc = document.querySelector('meta[name="description"]')
  if (metaDesc) metaDesc.setAttribute('content', s.seo_description)
  else if (s.seo_description) {
    const m = document.createElement('meta')
    m.name = 'description'
    m.content = s.seo_description
    document.head.appendChild(m)
  }
}

export const useSiteStore = defineStore('site', () => {
  const defaults = getDefaults()
  const settings = ref<SiteSettings>({ ...defaults })
  const loaded = ref(false)

  async function load() {
    const defaults = getDefaults()
    try {
      const data = await get<SiteSettings>('/api/v1/site-settings')
      if (data) settings.value = { ...defaults, ...data }
    } catch { /* use defaults */ }
    applyTheme(settings.value)
    loaded.value = true
  }

  watch(settings, (v) => applyTheme(v), { deep: true })

  return { settings, loaded, load }
})
