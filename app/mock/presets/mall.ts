import type { MockPreset, SiteSettings } from './types'
import categories from '../data/categories.json'
import products from '../data/products.json'
import productDetail from '../data/product-detail.json'
import indexDecor from '../data/index-decor.json'
import seckills from '../data/seckills.json'
import groupBuy from '../data/group-buy.json'
import bargain from '../data/bargain.json'
import recommend from '../data/recommend.json'
import cart from '../data/cart.json'
import orders from '../data/orders.json'

export const mall: MockPreset = {
  key: 'mall',
  name: '综合商城',
  categories: categories as any,
  products: products as any,
  productDetail: productDetail as any,
  indexDecor: indexDecor as any,
  siteSettings: {
    site_name: 'LYShop',
    site_logo: '',
    seo_title: 'LYShop - 精选好物 品质生活',
    seo_keywords: '商城,电商,购物,正品',
    seo_description: '综合购物商城，精选品质好物，正品保障，极速发货。',
    icp: '京ICP备2026XXXXXX号',
    hero_badge: '限时秒杀进行中',
    hero_title: '精选好物\\n品质生活从这里开始',
    hero_subtitle: '数千款精选商品，正品保障，极速发货，让购物更简单。',
    hero_btn_text: '立即选购',
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
  },
  seckills: seckills as any,
  groupBuy: groupBuy as any,
  bargain: bargain as any,
  recommend: recommend as any,
  cart: cart as any,
  orders: orders as any,
}
