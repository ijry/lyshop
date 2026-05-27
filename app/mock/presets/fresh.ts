import type { MockPreset } from './types'

const freshImageMap: Record<number, { cover: string; gallery: string[] }> = {
  201: {
    cover: 'https://images.unsplash.com/photo-1619566636858-adf3ef46400b?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1619566636858-adf3ef46400b?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1528825871115-3581a5387919?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1576179635662-9d1983e97e1e?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  202: {
    cover: 'https://images.unsplash.com/photo-1459411621453-7b03977f4bfc?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1459411621453-7b03977f4bfc?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1461354464878-ad92f492a5a0?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1518977676601-b53f82aba655?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  203: {
    cover: 'https://images.unsplash.com/photo-1510130387422-82bed34b37e9?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1510130387422-82bed34b37e9?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1544943910-4c1dc44aab44?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1559737558-2f5a35f4523b?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  204: {
    cover: 'https://images.unsplash.com/photo-1603048297172-c92544798d5a?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1603048297172-c92544798d5a?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1544025162-d76694265947?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1551028150-64b9f398f678?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  205: {
    cover: 'https://images.unsplash.com/photo-1506976785307-8732e854ad03?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1506976785307-8732e854ad03?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1482049016688-2d3e1b311543?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1518569656558-1f25e69d93d7?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  206: {
    cover: 'https://images.unsplash.com/photo-1485921325833-c519f76c4927?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1485921325833-c519f76c4927?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1510130387422-82bed34b37e9?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1579631542720-3a87824fff86?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  207: {
    cover: 'https://images.unsplash.com/photo-1553279768-865429fa0078?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1553279768-865429fa0078?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1605027990121-cbae9e0642df?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1591073113125-e46713c829ed?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  208: {
    cover: 'https://images.unsplash.com/photo-1576045057995-568f588f82fb?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1576045057995-568f588f82fb?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1518977676601-b53f82aba655?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1576045057995-568f588f82fb?auto=format&fit=crop&w=1200&q=80',
    ],
  },
}

function freshCover(id: number) {
  return freshImageMap[id]?.cover || freshImageMap[201].cover
}

function freshGallery(id: number, index = 0) {
  const list = freshImageMap[id]?.gallery || freshImageMap[201].gallery
  return list[index] || list[0]
}

export const fresh: MockPreset = {
  key: 'fresh',
  name: '生鲜果蔬',

  categories: [
    { id: 201, parent_id: 0, name: '新鲜水果', icon: '', sort: 1, status: 1 },
    { id: 202, parent_id: 0, name: '时令蔬菜', icon: '', sort: 2, status: 1 },
    { id: 203, parent_id: 0, name: '海鲜水产', icon: '', sort: 3, status: 1 },
    { id: 204, parent_id: 0, name: '肉禽蛋品', icon: '', sort: 4, status: 1 },
    { id: 205, parent_id: 0, name: '冷冻速食', icon: '', sort: 5, status: 1 },
    { id: 206, parent_id: 201, name: '进口水果', icon: '', sort: 1, status: 1 },
    { id: 207, parent_id: 202, name: '有机蔬菜', icon: '', sort: 1, status: 1 },
  ],

  products: {
    list: [
      { id: 201, title: '智利车厘子5斤装', subtitle: '进口JJ级大果，新鲜直达', cover: freshCover(201), price: 168.00, origin_price: 238.00, stock: 300, sales: 4560, status: 1, category_id: 201 },
      { id: 202, title: '有机西蓝花2颗装', subtitle: '有机认证，翠绿新鲜', cover: freshCover(202), price: 15.90, origin_price: 22.00, stock: 800, sales: 3210, status: 1, category_id: 202 },
      { id: 203, title: '鲜活大闸蟹礼盒', subtitle: '阳澄湖直发，只只肥美', cover: freshCover(203), price: 399.00, origin_price: 558.00, stock: 150, sales: 1280, status: 1, category_id: 203 },
      { id: 204, title: '澳洲安格斯牛排200g×4', subtitle: '原切谷饲，入口即化', cover: freshCover(204), price: 259.00, origin_price: 358.00, stock: 200, sales: 2350, status: 1, category_id: 204 },
      { id: 205, title: '土鸡蛋30枚', subtitle: '散养走地鸡，蛋黄饱满', cover: freshCover(205), price: 39.90, origin_price: 55.00, stock: 1000, sales: 8900, status: 1, category_id: 204 },
      { id: 206, title: '挪威三文鱼刺身200g', subtitle: '冰鲜空运，刺身级品质', cover: freshCover(206), price: 89.90, origin_price: 128.00, stock: 400, sales: 3670, status: 1, category_id: 203 },
      { id: 207, title: '海南金煌芒果10斤', subtitle: '树上熟大果，香甜多汁', cover: freshCover(207), price: 59.90, origin_price: 89.00, stock: 600, sales: 6540, status: 1, category_id: 201 },
      { id: 208, title: '有机菠菜500g', subtitle: '有机种植，鲜嫩翠绿', cover: freshCover(208), price: 9.90, origin_price: 15.00, stock: 1500, sales: 5120, status: 1, category_id: 202 },
    ],
    total: 8,
    page: 1,
    size: 20,
  },

  productDetail: {
    id: 201,
    title: '智利车厘子5斤装',
    subtitle: '进口JJ级大果，果径28-30mm，新鲜直达',
    cover: freshCover(201),
    price: 168.00,
    origin_price: 238.00,
    stock: 300,
    sales: 4560,
    status: 1,
    category_id: 201,
    detail: {
      version: 1,
      blocks: [
        { id: 'b1', type: 'text', props: { text: '智利直采JJ级车厘子，果径28-30mm，颗颗饱满。全程冷链运输，锁住新鲜与甜度。' } },
        { id: 'b2', type: 'image', props: { url: freshGallery(201, 1), alt: '车厘子实拍' } },
        { id: 'b3', type: 'text', props: { text: '收货后请冷藏保存，建议3天内食用完毕。如有破损可申请售后理赔。' } },
      ],
    },
    skus: [
      { id: 201, product_id: 201, attrs: '[{"name":"重量","value":"3斤装"}]', price: 108.00, stock: 100, sku_code: 'CHERRY-3' },
      { id: 202, product_id: 201, attrs: '[{"name":"重量","value":"5斤装"}]', price: 168.00, stock: 150, sku_code: 'CHERRY-5' },
      { id: 203, product_id: 201, attrs: '[{"name":"重量","value":"10斤装"}]', price: 318.00, stock: 50, sku_code: 'CHERRY-10' },
    ],
    images: [
      { id: 201, product_id: 201, url: freshGallery(201, 0), sort: 0 },
      { id: 202, product_id: 201, url: freshGallery(201, 1), sort: 1 },
      { id: 203, product_id: 201, url: freshGallery(201, 2), sort: 2 },
    ],
  },

  indexDecor: {
    components: [
      {
        type: 'banner',
        id: 'demo_banner',
        props: {
          images: [
            { url: '/static/demo/banner-fruit.png', link: '/pages/product/list?category_id=201' },
            { url: '/static/demo/banner-seafood.png', link: '/pages/product/list?category_id=203' },
            { url: '/static/demo/banner-meat.png', link: '/pages/product/list?category_id=204' },
          ],
          height: 340,
        },
      },
      {
        type: 'category_nav',
        id: 'demo_nav',
        props: {
          items: [
            { title: '水果', icon: '', link: '/pages/product/list?category_id=201' },
            { title: '蔬菜', icon: '', link: '/pages/product/list?category_id=202' },
            { title: '海鲜', icon: '', link: '/pages/product/list?category_id=203' },
            { title: '肉禽', icon: '', link: '/pages/product/list?category_id=204' },
            { title: '速食', icon: '', link: '/pages/product/list?category_id=205' },
          ],
        },
      },
      {
        type: 'grid',
        id: 'demo_grid_entry',
        props: {
          columns: 4,
          items: [
            { title: '秒杀', icon: '⚡', bg: '#fef2f2', link: '/pages/marketing/seckill' },
            { title: '拼团', icon: '👥', bg: '#eff6ff', link: '/pages/marketing/group-buy' },
            { title: '砍价', icon: '🔪', bg: '#f0fdf4', link: '/pages/marketing/bargain' },
            { title: '优惠券', icon: '🏷️', bg: '#fff7ed', link: '/pages/marketing/coupon?mode=claim' },
            { title: '积分', icon: '⭐', bg: '#fefce8', link: '/pages/user/points' },
            { title: '签到', icon: '📅', bg: '#faf5ff', link: '/pages/checkin/index' },
            { title: '收藏', icon: '❤️', bg: '#fff1f2', link: '' },
            { title: '客服', icon: '💬', bg: '#ecfdf5', link: '/pages/im/chat' },
          ],
        },
      },
      {
        type: 'notice',
        id: 'demo_notice',
        props: {
          items: [
            { text: '产地直供，冷链配送到家，新鲜看得见！', link: '' },
            { text: '新人首单满88减15，生鲜蔬果天天特价', link: '/pages/marketing/coupon?mode=claim' },
            { text: '每日限量秒杀，进口水果低至5折', link: '/pages/marketing/seckill' },
          ],
          color: '#16a34a',
          bgColor: '#f0fdf4',
          duration: 2500,
          mode: 'link',
        },
      },
      {
        type: 'marketing_zone',
        id: 'demo_seckill',
        props: {},
      },
      {
        type: 'spacer',
        id: 'demo_spacer',
        props: { height: 16, background: '#f5f5f5' },
      },
      {
        type: 'product_grid',
        id: 'demo_grid',
        props: { source: 'hot', limit: 8, columns: 2, title: '热销推荐' },
      },
    ],
  },

  pcDecor: {
    components: [
      {
        type: 'banner',
        id: 'pc_banner',
        props: {
          images: [
            { url: freshGallery(201, 0), link: '/products?category=201' },
            { url: freshGallery(203, 0), link: '/products?category=203' },
            { url: freshGallery(204, 0), link: '/products?category=204' },
          ],
          height: 400,
        },
      },
      {
        type: 'category_nav',
        id: 'pc_cats',
        props: {
          style: 'floating',
          columns: 5,
          items: [
            { title: '新鲜水果', icon: '', link: '/products?category=201' },
            { title: '时令蔬菜', icon: '', link: '/products?category=202' },
            { title: '海鲜水产', icon: '', link: '/products?category=203' },
            { title: '肉禽蛋品', icon: '', link: '/products?category=204' },
            { title: '冷冻速食', icon: '', link: '/products?category=205' },
          ],
        },
      },
      {
        type: 'notice',
        id: 'pc_notice',
        props: {
          items: [
            { text: '产地直供冷链到家', link: '' },
            { text: '新人满88减15', link: '/products' },
            { text: '每日限量秒杀', link: '/products' },
          ],
          color: '#16a34a',
          bgColor: '#f0fdf4',
        },
      },
      {
        type: 'marketing_zone',
        id: 'pc_seckill',
        props: {
          title: '限时秒杀',
          subtitle: '新鲜不等待',
          bg_from: '#15803d',
          bg_to: '#16a34a',
          more_link: '/products',
          products: [
            { product_id: 201, title: '智利车厘子5斤装', cover: freshCover(201), origin_price: 238, activity_price: 128 },
            { product_id: 207, title: '海南金煌芒果10斤', cover: freshCover(207), origin_price: 89, activity_price: 39.9 },
            { product_id: 206, title: '挪威三文鱼刺身200g', cover: freshCover(206), origin_price: 128, activity_price: 69.9 },
            { product_id: 204, title: '澳洲安格斯牛排', cover: freshCover(204), origin_price: 358, activity_price: 199 },
          ],
        },
      },
      {
        type: 'product_grid',
        id: 'pc_hot',
        props: { title: '今日精选', source: 'hot', limit: 8, columns: 4 },
      },
      {
        type: 'spacer',
        id: 'pc_spacer1',
        props: { height: 24 },
      },
      {
        type: 'image_ad',
        id: 'pc_ad1',
        props: { url: freshGallery(207, 0), link: '/products?category=201' },
      },
      {
        type: 'product_grid',
        id: 'pc_new',
        props: { title: '新鲜上市', source: 'new', limit: 4, columns: 4 },
      },
      {
        type: 'features',
        id: 'pc_features',
        props: {
          columns: 4,
          items: [
            { title: '冷链配送', icon: 'i-carbon-delivery-truck', desc: '全程冷链保鲜' },
            { title: '产地直供', icon: 'i-carbon-checkmark-outline', desc: '源头直采' },
            { title: '无忧退换', icon: 'i-carbon-renew', desc: '坏果包赔' },
            { title: '在线客服', icon: 'i-carbon-headset', desc: '7×24在线' },
          ],
        },
      },
    ],
  },

  siteSettings: {
    site_name: '鲜到家',
    site_logo: '',
    seo_title: '鲜到家 - 产地直供 冷链配送',
    seo_keywords: '生鲜,水果,蔬菜,冷链配送,产地直供',
    seo_description: '产地直供生鲜平台，冷链配送到家，新鲜看得见。',
    icp: '',
    hero_badge: '产地直供，冷链配送到家',
    hero_title: '产地直供\\n鲜到你家',
    hero_subtitle: '每日精选时令水果蔬菜，冷链极速配送，新鲜看得见。',
    hero_btn_text: '立即抢鲜',
    hero_btn_link: '/products',
    color_primary: '#16a34a',
    color_primary_light: '#22c55e',
    color_primary_dark: '#15803d',
    color_bg_page: '#f0fdf4',
    color_bg_header: 'rgba(255,255,255,0.8)',
    color_bg_footer: '#f0fdf4',
    color_price: '#16a34a',
    color_hero_from: '#15803d',
    color_hero_to: '#166534',
  },

  seckills: {
    list: [
      {
        id: 201, type: 'seckill', name: '生鲜限时秒杀',
        start_at: '2026-05-20T00:00:00Z', end_at: '2026-06-20T23:59:59Z', status: 1,
        products: [
          { product_id: 201, title: '智利车厘子5斤装', cover: freshCover(201), origin_price: 168.00, activity_price: 128.00, activity_stock: 80 },
          { product_id: 207, title: '海南金煌芒果10斤', cover: freshCover(207), origin_price: 59.90, activity_price: 39.90, activity_stock: 150 },
        ],
      },
    ],
    end_at: '2026-06-20T23:59:59Z',
  },

  groupBuy: {
    list: [
      {
        id: 201, type: 'group_buy', name: '3人拼团特惠',
        group_size: 3, expire_hours: 24,
        start_at: '2026-05-01T00:00:00Z', end_at: '2026-07-01T23:59:59Z', status: 1,
        products: [
          { product_id: 204, title: '澳洲安格斯牛排200g×4', cover: freshCover(204), origin_price: 259.00, group_price: 199.00, group_stock: 100 },
          { product_id: 205, title: '土鸡蛋30枚', cover: freshCover(205), origin_price: 39.90, group_price: 29.90, group_stock: 300 },
        ],
      },
    ],
  },

  bargain: {
    list: [
      {
        id: 201, type: 'bargain', name: '砍价免费拿',
        start_at: '2026-05-01T00:00:00Z', end_at: '2026-07-01T23:59:59Z', status: 1,
        products: [
          { product_id: 206, title: '挪威三文鱼刺身200g', cover: freshCover(206), origin_price: 128.00, floor_price: 0.01, max_helpers: 20, current_helpers: 9 },
        ],
      },
    ],
  },

  recommend: [
    { product_id: 201, title: '智利车厘子5斤装', cover: freshCover(201), price: 168.00, origin_price: 238.00, sales: 4560 },
    { product_id: 205, title: '土鸡蛋30枚', cover: freshCover(205), price: 39.90, origin_price: 55.00, sales: 8900 },
    { product_id: 207, title: '海南金煌芒果10斤', cover: freshCover(207), price: 59.90, origin_price: 89.00, sales: 6540 },
    { product_id: 206, title: '挪威三文鱼刺身200g', cover: freshCover(206), price: 89.90, origin_price: 128.00, sales: 3670 },
  ],

  cart: [
    {
      sku_id: 202,
      qty: 1,
      product: { id: 201, title: '智利车厘子5斤装', cover: freshCover(201), price: 168.00 },
      sku: { id: 202, product_id: 201, attrs: '[{"name":"重量","value":"5斤装"}]', price: 168.00, stock: 150 },
    },
    {
      sku_id: 205,
      qty: 2,
      product: { id: 205, title: '土鸡蛋30枚', cover: freshCover(205), price: 39.90 },
      sku: { id: 205, product_id: 205, attrs: '[{"name":"规格","value":"30枚"}]', price: 39.90, stock: 1000 },
    },
  ],

  orders: {
    list: [
      {
        id: 201, order_no: 'FR20260501001', user_id: 1, status: 4,
        activity_type: 'seckill', activity_name: '生鲜秒杀',
        payment_method: 'wechat', goods_amount: 168.00, discount_amount: 40.00,
        freight_amount: 0, total_amount: 128.00, remark: '',
        tracking_no: 'SF2233445566',
        amount_breakdown: { goods_amount: 168.00, discount_amount: 40.00, freight_amount: 0.00, payable_amount: 128.00 },
        items: [
          {
            id: 2011, product_id: 201, title: '智利车厘子5斤装', cover: freshCover(201), qty: 1, price: 128.00,
            review: {
              id: 1201, review_id: 1201, has_review: true, product_score: 5, logistics_score: 5,
              content: '车厘子又大又甜，冷链配送很新鲜，完全没有坏果！', edited_times: 0,
              appends: [
                { id: 2201, content: '第二天又下单了一箱，家人都说好吃。', created_at: '2026-05-03T10:00:00Z' },
              ],
              admin_reply: { id: 3201, content: '感谢您的好评，我们会继续保证品质！', created_at: '2026-05-04T12:00:00Z' },
              created_at: '2026-05-01T10:31:00Z', updated_at: '2026-05-03T10:00:00Z',
            },
          },
          { id: 2012, product_id: 208, title: '有机菠菜500g', cover: freshCover(208), qty: 2, price: 9.90 },
        ],
        created_at: '2026-05-01T10:30:00Z', paid_at: '2026-05-01T10:31:00Z',
      },
      {
        id: 202, order_no: 'FR20260510002', user_id: 1, status: 3,
        activity_type: 'group_buy', activity_name: '3人拼团',
        payment_method: 'alipay', goods_amount: 259.00, discount_amount: 60.00,
        freight_amount: 0, total_amount: 199.00, remark: '注意保鲜',
        tracking_no: 'SF7788990011',
        amount_breakdown: { goods_amount: 259.00, discount_amount: 60.00, freight_amount: 0.00, payable_amount: 199.00 },
        items: [
          {
            id: 2021, product_id: 204, title: '澳洲安格斯牛排200g×4', cover: freshCover(204), qty: 1, price: 199.00,
            review: {
              id: 1202, review_id: 1202, has_review: true, product_score: 4, logistics_score: 5,
              content: '牛排品质不错，冷链送达很新鲜。', edited_times: 0,
              appends: [], admin_reply: null,
              created_at: '2026-05-12T14:20:00Z', updated_at: '2026-05-12T14:20:00Z',
            },
          },
        ],
        created_at: '2026-05-10T14:20:00Z', paid_at: '2026-05-10T14:22:00Z',
      },
      {
        id: 203, order_no: 'FR20260520003', user_id: 1, status: 1,
        activity_type: '', activity_name: '',
        payment_method: '', goods_amount: 89.90, discount_amount: 0,
        freight_amount: 0, total_amount: 89.90, remark: '',
        amount_breakdown: { goods_amount: 89.90, discount_amount: 0.00, freight_amount: 0.00, payable_amount: 89.90 },
        items: [
          { id: 2031, product_id: 206, title: '挪威三文鱼刺身200g', cover: freshCover(206), qty: 1, price: 89.90 },
        ],
        created_at: '2026-05-20T09:00:00Z',
      },
      {
        id: 204, order_no: 'FR20260521004', user_id: 1, status: 2,
        activity_type: 'bargain', activity_name: '砍价免费拿',
        payment_method: 'wechat', goods_amount: 128.00, discount_amount: 127.99,
        freight_amount: 0, total_amount: 0.01, remark: '',
        tracking_no: '',
        amount_breakdown: { goods_amount: 128.00, discount_amount: 127.99, freight_amount: 0.00, payable_amount: 0.01 },
        items: [
          { id: 2041, product_id: 206, title: '挪威三文鱼刺身200g', cover: freshCover(206), qty: 1, price: 128.00 },
        ],
        created_at: '2026-05-21T16:00:00Z', paid_at: '2026-05-21T16:02:00Z',
      },
    ],
    total: 4,
    page: 1,
    size: 20,
  },
}
