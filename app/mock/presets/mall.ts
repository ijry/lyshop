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
  pcDecor: {
    components: [
      {
        type: 'hero',
        id: 'pc_hero',
        props: {
          badge: '限时秒杀进行中',
          title: '精选好物\n品质生活从这里开始',
          subtitle: '数千款精选商品，正品保障，极速发货，让购物更简单。',
          btn_text: '立即选购',
          btn_link: '/products',
          btn2_text: '查看全部',
          btn2_link: '/products',
          bg_from: '#b91c1c',
          bg_to: '#991b1b',
        },
      },
      {
        type: 'category_nav',
        id: 'pc_cats',
        props: {
          style: 'floating',
          columns: 8,
          items: [
            { title: '手机数码', icon: '', link: '/products?category=1' },
            { title: '家用电器', icon: '', link: '/products?category=2' },
            { title: '服装鞋帽', icon: '', link: '/products?category=3' },
            { title: '美妆护肤', icon: '', link: '/products?category=4' },
            { title: '食品饮料', icon: '', link: '/products?category=5' },
            { title: '家居日用', icon: '', link: '/products?category=6' },
            { title: '运动户外', icon: '', link: '/products?category=7' },
            { title: '图书文具', icon: '', link: '/products?category=8' },
          ],
        },
      },
      {
        type: 'notice',
        id: 'pc_notice',
        props: {
          items: [
            { text: '618年中大促：全场满300减50，限时3天！', link: '' },
            { text: '新用户注册即送50元优惠券，立即领取', link: '/products' },
            { text: '每日10点限时秒杀，爆款低至5折', link: '/products' },
          ],
          color: '#dc2626',
          bgColor: '#fef2f2',
        },
      },
      {
        type: 'image_ad',
        id: 'pc_ad1',
        props: {
          url: 'https://images.unsplash.com/photo-1550009158-9ebf69173e03?auto=format&fit=crop&w=1400&q=80',
          link: '/products',
          height: 200,
        },
      },
      {
        type: 'product_grid',
        id: 'pc_hot',
        props: {
          title: '热销推荐',
          source: 'hot',
          limit: 8,
          columns: 4,
        },
      },
      {
        type: 'marketing_zone',
        id: 'pc_seckill',
        props: {
          title: '限时秒杀',
          subtitle: '限时抢购中',
          bg_from: '#b91c1c',
          bg_to: '#dc2626',
          more_link: '/products',
          products: [
            { product_id: 1, title: '旗舰智能手机 Pro Max', cover: 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?auto=format&fit=crop&w=800&q=80', origin_price: 4999, activity_price: 3999 },
            { product_id: 2, title: '轻薄笔记本电脑 Air', cover: 'https://images.unsplash.com/photo-1496181133206-80ce9b88a853?auto=format&fit=crop&w=800&q=80', origin_price: 6999, activity_price: 5499 },
            { product_id: 3, title: '真无线降噪耳机', cover: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?auto=format&fit=crop&w=800&q=80', origin_price: 899, activity_price: 599 },
            { product_id: 4, title: '智能运动手表', cover: 'https://images.unsplash.com/photo-1546868871-7041f2a55e12?auto=format&fit=crop&w=800&q=80', origin_price: 1299, activity_price: 899 },
          ],
        },
      },
      {
        type: 'spacer',
        id: 'pc_spacer1',
        props: {
          height: 32,
        },
      },
      {
        type: 'product_grid',
        id: 'pc_new',
        props: {
          title: '新品上架',
          source: 'new',
          limit: 8,
          columns: 4,
        },
      },
      {
        type: 'features',
        id: 'pc_features',
        props: {
          columns: 4,
          items: [
            { icon: 'i-carbon-delivery-truck', title: '快递配送', desc: '全国包邮，极速送达' },
            { icon: 'i-carbon-checkmark-outline', title: '正品保障', desc: '品牌授权，假一赔十' },
            { icon: 'i-carbon-renew', title: '无忧退换', desc: '7天无理由退换货' },
            { icon: 'i-carbon-headset', title: '在线客服', desc: '7x24小时在线服务' },
          ],
        },
      },
    ],
  },
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
