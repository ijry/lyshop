import type { MockPreset, SiteSettings } from './types'

const jewelryImageMap: Record<number, { cover: string; gallery: string[] }> = {
  301: {
    cover: 'https://images.unsplash.com/photo-1617038220319-276d3cfab638?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1617038220319-276d3cfab638?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1611652022419-a9419f74343d?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1515562141207-7a88fb7ce338?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  302: {
    cover: 'https://images.unsplash.com/photo-1605100804763-247f67b3557e?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1605100804763-247f67b3557e?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1599643478518-a784e5dc4c8f?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1543295204-8e6d7d62f8df?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  303: {
    cover: 'https://images.unsplash.com/photo-1617038260897-41a1f14a8ca0?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1617038260897-41a1f14a8ca0?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1611107683227-e9060eccd846?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1617038220319-276d3cfab638?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  304: {
    cover: 'https://images.unsplash.com/photo-1611599537845-1c7aca0091c0?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1611599537845-1c7aca0091c0?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1535632066927-ab7c9ab60908?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1603974372039-adc49044b6bd?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  305: {
    cover: 'https://images.unsplash.com/photo-1535632787350-4e68ef0ac584?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1535632787350-4e68ef0ac584?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1617038260897-41a1f14a8ca0?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1589128777073-263566ae5e4d?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  306: {
    cover: 'https://images.unsplash.com/photo-1617038261661-1d5f1c4ccb20?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1617038261661-1d5f1c4ccb20?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1617038260897-41a1f14a8ca0?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1617038220319-276d3cfab638?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  307: {
    cover: 'https://images.unsplash.com/photo-1617038260897-41a1f14a8ca0?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1617038260897-41a1f14a8ca0?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1617038261661-1d5f1c4ccb20?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1611652022419-a9419f74343d?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  308: {
    cover: 'https://images.unsplash.com/photo-1588444837495-c6cfeb53f32d?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1588444837495-c6cfeb53f32d?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1535632066927-ab7c9ab60908?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1602173574767-37ac01994b2a?auto=format&fit=crop&w=1200&q=80',
    ],
  },
}

function jewelryCover(id: number) {
  return jewelryImageMap[id]?.cover || jewelryImageMap[301].cover
}

function jewelryGallery(id: number, index = 0) {
  const list = jewelryImageMap[id]?.gallery || jewelryImageMap[301].gallery
  return list[index] || list[0]
}

export const jewelry: MockPreset = {
  key: 'jewelry',
  name: '珠宝饰品',
  categories: [
    { id: 301, parent_id: 0, name: '黄金饰品', icon: '', sort: 1, status: 1 },
    { id: 302, parent_id: 0, name: '钻石珠宝', icon: '', sort: 2, status: 1 },
    { id: 303, parent_id: 0, name: '翡翠玉石', icon: '', sort: 3, status: 1 },
    { id: 304, parent_id: 0, name: '珍珠饰品', icon: '', sort: 4, status: 1 },
    { id: 305, parent_id: 0, name: '银饰轻奢', icon: '', sort: 5, status: 1 },
    { id: 306, parent_id: 301, name: '黄金项链', icon: '', sort: 1, status: 1 },
    { id: 307, parent_id: 302, name: '婚戒对戒', icon: '', sort: 1, status: 1 },
  ],
  products: {
    list: [
      {
        id: 301, title: '足金转运珠手链', subtitle: '999足金，精工编织，寓意好运',
        cover: jewelryCover(301),
        price: 1580.00, origin_price: 1980.00, stock: 300, sales: 2150, status: 1, category_id: 301,
      },
      {
        id: 302, title: '1克拉钻石戒指', subtitle: 'GIA认证，D色VS1净度',
        cover: jewelryCover(302),
        price: 29999.00, origin_price: 39999.00, stock: 50, sales: 86, status: 1, category_id: 302,
      },
      {
        id: 303, title: '冰种翡翠吊坠', subtitle: '缅甸A货翡翠，冰润通透',
        cover: jewelryCover(303),
        price: 6800.00, origin_price: 8800.00, stock: 80, sales: 320, status: 1, category_id: 303,
      },
      {
        id: 304, title: '天然淡水珍珠项链', subtitle: '8-9mm正圆强光，925银扣',
        cover: jewelryCover(304),
        price: 899.00, origin_price: 1280.00, stock: 500, sales: 1860, status: 1, category_id: 304,
      },
      {
        id: 305, title: '925纯银耳钉套装', subtitle: '一周耳钉七对装，简约百搭',
        cover: jewelryCover(305),
        price: 99.00, origin_price: 168.00, stock: 2000, sales: 6820, status: 1, category_id: 305,
      },
      {
        id: 306, title: '和田玉平安扣', subtitle: '新疆和田籽料，油润细腻',
        cover: jewelryCover(306),
        price: 3680.00, origin_price: 4980.00, stock: 120, sales: 450, status: 1, category_id: 303,
      },
      {
        id: 307, title: '18K玫瑰金手镯', subtitle: '意大利工艺，优雅时尚',
        cover: jewelryCover(307),
        price: 4280.00, origin_price: 5680.00, stock: 200, sales: 680, status: 1, category_id: 301,
      },
      {
        id: 308, title: '蓝宝石耳环', subtitle: '斯里兰卡蓝宝石，璀璨夺目',
        cover: jewelryCover(308),
        price: 2580.00, origin_price: 3580.00, stock: 150, sales: 390, status: 1, category_id: 302,
      },
    ],
    total: 8,
    page: 1,
    size: 20,
  },
  productDetail: {
    id: 301,
    title: '足金转运珠手链',
    subtitle: '999足金，精工编织，寓意好运',
    cover: jewelryCover(301),
    price: 1580.00,
    origin_price: 1980.00,
    stock: 300,
    sales: 2150,
    status: 1,
    category_id: 301,
    detail: {
      version: 1,
      blocks: [
        {
          id: 'b1',
          type: 'text',
          props: {
            text: '精选999足金，采用古法錾刻工艺，每颗转运珠均为匠人手工打磨。含金量≥99.9%，附带NGTC国检证书。',
          },
        },
        {
          id: 'b2',
          type: 'image',
          props: {
            url: jewelryGallery(301, 1),
            alt: '足金工艺细节展示',
          },
        },
        {
          id: 'b3',
          type: 'text',
          props: {
            text: '红绳采用手工编织五彩金刚结，结实耐用不褪色，寓意转运纳福、吉祥如意。',
          },
        },
      ],
    },
    skus: [
      { id: 301, product_id: 301, attrs: '[{"name":"重量","value":"3g"}]', price: 1580.00, stock: 120, sku_code: 'GOLD-BEAD-3G' },
      { id: 302, product_id: 301, attrs: '[{"name":"重量","value":"5g"}]', price: 2580.00, stock: 100, sku_code: 'GOLD-BEAD-5G' },
      { id: 303, product_id: 301, attrs: '[{"name":"重量","value":"8g"}]', price: 4080.00, stock: 80, sku_code: 'GOLD-BEAD-8G' },
    ],
    images: [
      { id: 301, product_id: 301, url: jewelryGallery(301, 0), sort: 0 },
      { id: 302, product_id: 301, url: jewelryGallery(301, 1), sort: 1 },
      { id: 303, product_id: 301, url: jewelryGallery(301, 2), sort: 2 },
      { id: 304, product_id: 301, url: jewelryGallery(307, 0), sort: 3 },
    ],
  },
  indexDecor: {
    components: [
      {
        type: 'banner',
        id: 'demo_banner',
        props: {
          images: [
            { url: jewelryGallery(301, 0), link: '/pages/product/list?category_id=301' },
            { url: jewelryGallery(305, 0), link: '/pages/marketing/coupon?mode=claim' },
            { url: jewelryGallery(302, 0), link: '/pages/product/list?category_id=302' },
          ],
          height: 340,
        },
      },
      {
        type: 'category_nav',
        id: 'demo_nav',
        props: {
          items: [
            { title: '黄金', icon: '', link: '/pages/product/list?category_id=301' },
            { title: '钻石', icon: '', link: '/pages/product/list?category_id=302' },
            { title: '翡翠', icon: '', link: '/pages/product/list?category_id=303' },
            { title: '珍珠', icon: '', link: '/pages/product/list?category_id=304' },
            { title: '银饰', icon: '', link: '/pages/product/list?category_id=305' },
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
            { text: '每件商品均附GIA/NGTC权威证书，品质保障', link: '' },
            { text: '新会员首单立减200元，黄金饰品免工费', link: '/pages/marketing/coupon?mode=claim' },
            { text: '七天无理由退换，终身免费清洗保养', link: '' },
          ],
          color: '#b45309',
          bgColor: '#fffbeb',
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
        props: { source: 'hot', limit: 8, columns: 2, title: '臻品推荐' },
      },
    ],
  },
  pcDecor: {
    components: [
      {
        type: 'hero', id: 'pc_hero',
        props: {
          badge: '新品上市，限时特惠', title: '臻品珠宝\n璀璨人生',
          subtitle: '匠心工艺，甄选全球珍品，让每一刻都闪耀。',
          btn_text: '探索臻品', btn_link: '/products', btn2_text: '新品首发', btn2_link: '/products',
          bg_from: '#92400e', bg_to: '#78350f',
        },
      },
      {
        type: 'category_nav', id: 'pc_cats',
        props: {
          style: 'floating', columns: 5,
          items: [
            { title: '黄金饰品', icon: '', link: '/products?category=301' },
            { title: '钻石珠宝', icon: '', link: '/products?category=302' },
            { title: '翡翠玉石', icon: '', link: '/products?category=303' },
            { title: '珍珠饰品', icon: '', link: '/products?category=304' },
            { title: '银饰轻奢', icon: '', link: '/products?category=305' },
          ],
        },
      },
      { type: 'image_ad', id: 'pc_ad1', props: { url: jewelryGallery(301, 0), link: '/products?category=301', height: 280 } },
      { type: 'product_grid', id: 'pc_hot', props: { title: '臻品推荐', source: 'hot', limit: 8, columns: 4 } },
      { type: 'spacer', id: 'pc_spacer1', props: { height: 24 } },
      { type: 'image_ad', id: 'pc_ad2', props: { url: jewelryGallery(302, 0), link: '/products?category=302' } },
      { type: 'product_grid', id: 'pc_new', props: { title: '新品首发', source: 'new', limit: 4, columns: 4 } },
      {
        type: 'features', id: 'pc_features',
        props: {
          columns: 4,
          items: [
            { icon: 'i-carbon-delivery-truck', title: '顺丰包邮', desc: '全国顺丰配送' },
            { icon: 'i-carbon-checkmark-outline', title: '正品认证', desc: '附带鉴定证书' },
            { icon: 'i-carbon-renew', title: '30天退换', desc: '无忧售后' },
            { icon: 'i-carbon-headset', title: '专属顾问', desc: '1对1服务' },
          ],
        },
      },
    ],
  },
  siteSettings: {
    site_name: '璀璨珠宝',
    site_logo: '',
    seo_title: '璀璨珠宝 - 臻品珠宝 璀璨人生',
    seo_keywords: '珠宝,黄金,钻石,翡翠,首饰',
    seo_description: '臻品珠宝商城，黄金钻石翡翠，匠心工艺，品质之选。',
    icp: '',
    hero_badge: '新品上市，限时特惠',
    hero_title: '臻品珠宝\\n璀璨人生',
    hero_subtitle: '匠心工艺，甄选全球珍品，让每一刻都闪耀。',
    hero_btn_text: '探索臻品',
    hero_btn_link: '/products',
    color_primary: '#b45309',
    color_primary_light: '#d97706',
    color_primary_dark: '#92400e',
    color_bg_page: '#fffbeb',
    color_bg_header: 'rgba(255,255,255,0.8)',
    color_bg_footer: '#fffbeb',
    color_price: '#b45309',
    color_hero_from: '#92400e',
    color_hero_to: '#78350f',
  },
  seckills: {
    list: [
      {
        id: 1, type: 'seckill', name: '珠宝限时秒杀',
        start_at: '2026-05-20T00:00:00Z', end_at: '2026-06-20T23:59:59Z', status: 1,
        products: [
          { product_id: 301, title: '足金转运珠手链', cover: jewelryCover(301), origin_price: 1580, activity_price: 1280, activity_stock: 30 },
          { product_id: 305, title: '925纯银耳钉套装', cover: jewelryCover(305), origin_price: 99, activity_price: 59, activity_stock: 200 },
          { product_id: 304, title: '天然淡水珍珠项链', cover: jewelryCover(304), origin_price: 899, activity_price: 599, activity_stock: 80 },
        ],
      },
    ],
    end_at: '2026-06-20T23:59:59Z',
  },
  groupBuy: {
    list: [
      {
        id: 1, type: 'group_buy', name: '3人拼团特惠',
        group_size: 3, expire_hours: 24,
        start_at: '2026-05-01T00:00:00Z', end_at: '2026-07-01T23:59:59Z', status: 1,
        products: [
          { product_id: 303, title: '冰种翡翠吊坠', cover: jewelryCover(303), origin_price: 6800, group_price: 5200, group_stock: 30 },
          { product_id: 304, title: '天然淡水珍珠项链', cover: jewelryCover(304), origin_price: 899, group_price: 699, group_stock: 100 },
          { product_id: 308, title: '蓝宝石耳环', cover: jewelryCover(308), origin_price: 2580, group_price: 1980, group_stock: 50 },
        ],
      },
    ],
  },
  bargain: {
    list: [
      {
        id: 1, type: 'bargain', name: '砍价免费拿',
        start_at: '2026-05-01T00:00:00Z', end_at: '2026-07-01T23:59:59Z', status: 1,
        products: [
          { product_id: 305, title: '925纯银耳钉套装', cover: jewelryCover(305), origin_price: 99, floor_price: 0.01, max_helpers: 10, current_helpers: 4 },
          { product_id: 307, title: '18K玫瑰金手镯', cover: jewelryCover(307), origin_price: 4280, floor_price: 1280, max_helpers: 30, current_helpers: 12 },
        ],
      },
    ],
  },
  recommend: [
    { product_id: 305, title: '925纯银耳钉套装', cover: jewelryCover(305), price: 99, origin_price: 168, sales: 6820 },
    { product_id: 301, title: '足金转运珠手链', cover: jewelryCover(301), price: 1580, origin_price: 1980, sales: 2150 },
    { product_id: 304, title: '天然淡水珍珠项链', cover: jewelryCover(304), price: 899, origin_price: 1280, sales: 1860 },
    { product_id: 307, title: '18K玫瑰金手镯', cover: jewelryCover(307), price: 4280, origin_price: 5680, sales: 680 },
  ],
  cart: [
    {
      sku_id: 301,
      qty: 1,
      product: {
        id: 301, title: '足金转运珠手链', cover: jewelryCover(301), price: 1580.00,
      },
      sku: {
        id: 301, product_id: 301, attrs: '[{"name":"重量","value":"3g"}]', price: 1580.00, stock: 120,
      },
    },
    {
      sku_id: 305,
      qty: 2,
      product: {
        id: 305, title: '925纯银耳钉套装', cover: jewelryCover(305), price: 99.00,
      },
      sku: {
        id: 305, product_id: 305, attrs: '[{"name":"款式","value":"经典款"}]', price: 99.00, stock: 2000,
      },
    },
  ],
  orders: {
    list: [
      {
        id: 301, order_no: 'JWL20260501001', user_id: 1, status: 4,
        activity_type: 'seckill', activity_name: '珠宝限时秒杀',
        payment_method: 'wechat', goods_amount: 1580.00, discount_amount: 300,
        freight_amount: 0, total_amount: 1280.00, remark: '',
        tracking_no: 'SF1234567801',
        amount_breakdown: { goods_amount: 1580.00, discount_amount: 300.00, freight_amount: 0.00, payable_amount: 1280.00 },
        items: [
          {
            id: 3011, product_id: 301, title: '足金转运珠手链', cover: jewelryCover(301), qty: 1, price: 1580.00,
            review: {
              id: 4001, review_id: 4001, has_review: true, product_score: 5, logistics_score: 5,
              content: '做工精致，金珠光泽很好，送人很有面子。', edited_times: 0,
              appends: [
                { id: 5001, content: '戴了一个月，没有褪色，非常满意。', created_at: '2026-05-15T10:00:00Z' },
              ],
              admin_reply: { id: 6001, content: '感谢您的认可，祝您佩戴愉快！', created_at: '2026-05-16T09:00:00Z' },
              created_at: '2026-05-03T14:00:00Z', updated_at: '2026-05-15T10:00:00Z',
            },
          },
        ],
        created_at: '2026-05-01T10:30:00Z', paid_at: '2026-05-01T10:31:00Z',
      },
      {
        id: 302, order_no: 'JWL20260510002', user_id: 1, status: 3,
        activity_type: 'group_buy', activity_name: '3人拼团特惠',
        payment_method: 'alipay', goods_amount: 6800.00, discount_amount: 1600,
        freight_amount: 0, total_amount: 5200.00, remark: '请用礼盒包装',
        tracking_no: 'SF9988776601',
        amount_breakdown: { goods_amount: 6800.00, discount_amount: 1600.00, freight_amount: 0.00, payable_amount: 5200.00 },
        items: [
          { id: 3021, product_id: 303, title: '冰种翡翠吊坠', cover: jewelryCover(303), qty: 1, price: 6800.00 },
        ],
        created_at: '2026-05-10T14:20:00Z', paid_at: '2026-05-10T14:22:00Z',
      },
      {
        id: 303, order_no: 'JWL20260520003', user_id: 1, status: 1,
        activity_type: '', activity_name: '',
        payment_method: '', goods_amount: 29999.00, discount_amount: 0,
        freight_amount: 0, total_amount: 29999.00, remark: '',
        amount_breakdown: { goods_amount: 29999.00, discount_amount: 0.00, freight_amount: 0.00, payable_amount: 29999.00 },
        items: [
          { id: 3031, product_id: 302, title: '1克拉钻石戒指', cover: jewelryCover(302), qty: 1, price: 29999.00 },
        ],
        created_at: '2026-05-20T09:00:00Z',
      },
      {
        id: 304, order_no: 'JWL20260521004', user_id: 1, status: 2,
        activity_type: 'bargain', activity_name: '砍价免费拿',
        payment_method: 'wechat', goods_amount: 99.00, discount_amount: 98.99,
        freight_amount: 0, total_amount: 0.01, remark: '',
        tracking_no: '',
        amount_breakdown: { goods_amount: 99.00, discount_amount: 98.99, freight_amount: 0.00, payable_amount: 0.01 },
        items: [
          { id: 3041, product_id: 305, title: '925纯银耳钉套装', cover: jewelryCover(305), qty: 1, price: 99.00 },
        ],
        created_at: '2026-05-21T16:00:00Z', paid_at: '2026-05-21T16:02:00Z',
      },
    ],
    total: 4,
    page: 1,
    size: 20,
  },
}
