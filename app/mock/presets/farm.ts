import type { MockPreset, SiteSettings } from './types'

const farmImageMap: Record<number, { cover: string; gallery: string[] }> = {
  401: {
    cover: 'https://images.unsplash.com/photo-1586201375761-83865001e31c?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1586201375761-83865001e31c?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1509440159596-0249088772ff?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1516684669134-de6f7c473a2a?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  402: {
    cover: 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1515823064-d6e0c04616a7?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1464306076886-da185f6a9d05?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  403: {
    cover: 'https://images.unsplash.com/photo-1587049352851-8d4e89133924?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1587049352851-8d4e89133924?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1471943311424-646960669fbc?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1558640476-437a2b9438a2?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  404: {
    cover: 'https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1547592166-23ac45744acd?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1603048297172-c92544798d5a?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  405: {
    cover: 'https://images.unsplash.com/photo-1611080626919-7cf5a9dbab5b?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1611080626919-7cf5a9dbab5b?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1502741338009-cac2772e18bc?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1519096845289-95806ee03a1a?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  406: {
    cover: 'https://images.unsplash.com/photo-1615485290382-441e4d049cb5?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1615485290382-441e4d049cb5?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1604335399105-a0c585fd81a1?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1603048297172-c92544798d5a?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  407: {
    cover: 'https://images.unsplash.com/photo-1502741338009-cac2772e18bc?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1502741338009-cac2772e18bc?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1547592166-23ac45744acd?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1509440159596-0249088772ff?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  408: {
    cover: 'https://images.unsplash.com/photo-1544025162-d76694265947?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1544025162-d76694265947?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1603048297172-c92544798d5a?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&w=1200&q=80',
    ],
  },
}

function farmCover(id: number) {
  return farmImageMap[id]?.cover || farmImageMap[401].cover
}

function farmGallery(id: number, index = 0) {
  const list = farmImageMap[id]?.gallery || farmImageMap[401].gallery
  return list[index] || list[0]
}

export const farm: MockPreset = {
  key: 'farm',
  name: '农产品特产',
  categories: [
    { id: 401, parent_id: 0, name: '粮油米面', icon: '', sort: 1, status: 1 },
    { id: 402, parent_id: 0, name: '茶叶茗茶', icon: '', sort: 2, status: 1 },
    { id: 403, parent_id: 0, name: '蜂蜜滋补', icon: '', sort: 3, status: 1 },
    { id: 404, parent_id: 0, name: '干货腊味', icon: '', sort: 4, status: 1 },
    { id: 405, parent_id: 0, name: '地方特产', icon: '', sort: 5, status: 1 },
    { id: 406, parent_id: 401, name: '有机大米', icon: '', sort: 1, status: 1 },
    { id: 407, parent_id: 402, name: '绿茶', icon: '', sort: 1, status: 1 },
  ],
  products: {
    list: [
      {
        id: 401, title: '五常有机大米10斤', subtitle: '黑龙江五常产区，稻花香2号',
        cover: farmCover(401),
        price: 69.90, origin_price: 98.00, stock: 2000, sales: 8650, status: 1, category_id: 401,
      },
      {
        id: 402, title: '西湖龙井明前茶250g', subtitle: '核心产区，明前头采，鲜爽回甘',
        cover: farmCover(402),
        price: 268.00, origin_price: 388.00, stock: 300, sales: 1260, status: 1, category_id: 402,
      },
      {
        id: 403, title: '秦岭土蜂蜜500g', subtitle: '秦岭深山百花蜜，自然成熟封盖蜜',
        cover: farmCover(403),
        price: 58.00, origin_price: 88.00, stock: 800, sales: 3420, status: 1, category_id: 403,
      },
      {
        id: 404, title: '云南野生菌菇礼盒', subtitle: '松茸牛肝菌鸡枞菌，山珍组合',
        cover: farmCover(404),
        price: 188.00, origin_price: 268.00, stock: 200, sales: 960, status: 1, category_id: 404,
      },
      {
        id: 405, title: '赣南脐橙10斤装', subtitle: '江西赣州脐橙，皮薄多汁，甜度高',
        cover: farmCover(405),
        price: 49.90, origin_price: 69.00, stock: 1500, sales: 12800, status: 1, category_id: 405,
      },
      {
        id: 406, title: '新疆和田大枣500g', subtitle: '和田骏枣，个大肉厚，自然晾晒',
        cover: farmCover(406),
        price: 35.90, origin_price: 55.00, stock: 1000, sales: 5680, status: 1, category_id: 404,
      },
      {
        id: 407, title: '东北黑木耳250g', subtitle: '长白山秋耳，肉厚朵大，口感爽脆',
        cover: farmCover(407),
        price: 29.90, origin_price: 45.00, stock: 1200, sales: 4350, status: 1, category_id: 404,
      },
      {
        id: 408, title: '湖南烟熏腊肉500g', subtitle: '湘西农家柴火烟熏，肥瘦相间',
        cover: farmCover(408),
        price: 45.90, origin_price: 68.00, stock: 600, sales: 2780, status: 1, category_id: 404,
      },
    ],
    total: 8,
    page: 1,
    size: 20,
  },
  productDetail: {
    id: 401,
    title: '五常有机大米10斤',
    subtitle: '黑龙江五常产区，稻花香2号',
    cover: farmCover(401),
    price: 69.90,
    origin_price: 98.00,
    stock: 2000,
    sales: 8650,
    status: 1,
    category_id: 401,
    detail: {
      version: 1,
      blocks: [
        {
          id: 'b1',
          type: 'text',
          props: {
            text: '产自黑龙江五常核心产区，精选稻花香2号品种。通过国家有机认证，全程可溯源，从田间到餐桌安心无忧。',
          },
        },
        {
          id: 'b2',
          type: 'image',
          props: {
            url: farmGallery(401, 1),
            alt: '五常稻田实拍',
          },
        },
        {
          id: 'b3',
          type: 'text',
          props: {
            text: '当季新米现磨现发，颗粒饱满晶莹剔透，煮饭清香扑鼻、软糯弹牙。绿色有机种植，不使用化肥农药。',
          },
        },
      ],
    },
    skus: [
      { id: 401, product_id: 401, attrs: '[{"name":"规格","value":"5斤"}]', price: 39.90, stock: 800, sku_code: 'RICE-WC-5' },
      { id: 402, product_id: 401, attrs: '[{"name":"规格","value":"10斤"}]', price: 69.90, stock: 700, sku_code: 'RICE-WC-10' },
      { id: 403, product_id: 401, attrs: '[{"name":"规格","value":"20斤"}]', price: 129.90, stock: 500, sku_code: 'RICE-WC-20' },
    ],
    images: [
      { id: 401, product_id: 401, url: farmGallery(401, 0), sort: 0 },
      { id: 402, product_id: 401, url: farmGallery(401, 1), sort: 1 },
      { id: 403, product_id: 401, url: farmGallery(401, 2), sort: 2 },
      { id: 404, product_id: 401, url: farmGallery(405, 0), sort: 3 },
    ],
  },
  indexDecor: {
    components: [
      {
        type: 'banner',
        id: 'demo_banner',
        props: {
          images: [
            { url: farmGallery(401, 0), link: '/pages/product/list?category_id=401' },
            { url: farmGallery(402, 0), link: '/pages/marketing/coupon?mode=claim' },
            { url: farmGallery(405, 0), link: '/pages/product/list?category_id=405' },
          ],
          height: 340,
        },
      },
      {
        type: 'category_nav',
        id: 'demo_nav',
        props: {
          items: [
            { title: '粮油', icon: '', link: '/pages/product/list?category_id=401' },
            { title: '茶叶', icon: '', link: '/pages/product/list?category_id=402' },
            { title: '蜂蜜', icon: '', link: '/pages/product/list?category_id=403' },
            { title: '干货', icon: '', link: '/pages/product/list?category_id=404' },
            { title: '特产', icon: '', link: '/pages/product/list?category_id=405' },
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
            { text: '原产地直供，绿色有机认证，每一口都放心', link: '' },
            { text: '新用户首单满99减20，产地直发包邮', link: '/pages/marketing/coupon?mode=claim' },
            { text: '时令鲜品每周上新，从田间到餐桌48小时达', link: '/pages/product/list' },
          ],
          color: '#15803d',
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
        props: { source: 'hot', limit: 8, columns: 2, title: '田间精选' },
      },
    ],
  },
  pcDecor: {
    components: [
      {
        type: 'hero',
        id: 'pc_hero',
        props: {
          badge: '应季新品，产地直发',
          title: '田间到餐桌\\n自然的味道',
          subtitle: '精选各地优质农产品，绿色有机，从田间到您的餐桌。',
          btn_text: '逛逛好货',
          btn_link: '/products',
          btn2_text: '看看特产',
          btn2_link: '/products',
          bg_from: '#166534',
          bg_to: '#14532d',
        },
      },
      {
        type: 'category_nav',
        id: 'pc_cats',
        props: {
          style: 'floating',
          columns: 5,
          items: [
            { title: '粮油米面', link: '/products?category=401' },
            { title: '茶叶茗茶', link: '/products?category=402' },
            { title: '蜂蜜滋补', link: '/products?category=403' },
            { title: '干货腊味', link: '/products?category=404' },
            { title: '地方特产', link: '/products?category=405' },
          ],
        },
      },
      {
        type: 'notice',
        id: 'pc_notice',
        props: {
          items: [
            { text: '应季好货产地直发' },
            { text: '首单满99减20' },
            { text: '绿色有机认证' },
          ],
          color: '#15803d',
          bgColor: '#f0fdf4',
        },
      },
      {
        type: 'product_grid',
        id: 'pc_hot',
        props: {
          title: '田间精选',
          source: 'hot',
          limit: 8,
          columns: 4,
        },
      },
      {
        type: 'image_ad',
        id: 'pc_ad1',
        props: {
          url: farmGallery(401, 0),
          link: '/products?category=401',
          height: 220,
        },
      },
      {
        type: 'rich_text',
        id: 'pc_rich',
        props: {
          content: '<div style="text-align:center;padding:20px 0"><h3 style="font-size:20px;color:#166534;margin-bottom:8px">源头好货 · 品质保障</h3><p style="color:#666;font-size:14px">每一份农产品都经过严格品控，从田间到餐桌全程可追溯</p></div>',
        },
      },
      {
        type: 'product_grid',
        id: 'pc_new',
        props: {
          title: '地方特产',
          source: 'new',
          limit: 4,
          columns: 4,
        },
      },
      {
        type: 'features',
        id: 'pc_features',
        props: {
          columns: 4,
          items: [
            { title: '产地直发', icon: 'i-carbon-delivery-truck', desc: '48小时发货' },
            { title: '有机认证', icon: 'i-carbon-checkmark-outline', desc: '绿色有机' },
            { title: '坏损包赔', icon: 'i-carbon-renew', desc: '售后无忧' },
            { title: '源头客服', icon: 'i-carbon-headset', desc: '产地对接' },
          ],
        },
      },
    ],
  },
  siteSettings: {
    site_name: '田园优选',
    site_logo: '',
    seo_title: '田园优选 - 田间到餐桌 自然的味道',
    seo_keywords: '农产品,特产,有机,土特产,绿色食品',
    seo_description: '田间到餐桌，精选各地优质农产品特产，绿色有机，自然纯正。',
    icp: '',
    hero_badge: '应季新品，产地直发',
    hero_title: '田间到餐桌\\n自然的味道',
    hero_subtitle: '精选各地优质农产品，绿色有机，从田间到您的餐桌。',
    hero_btn_text: '逛逛好货',
    hero_btn_link: '/products',
    color_primary: '#15803d',
    color_primary_light: '#16a34a',
    color_primary_dark: '#166534',
    color_bg_page: '#f0fdf4',
    color_bg_header: 'rgba(255,255,255,0.8)',
    color_bg_footer: '#f0fdf4',
    color_price: '#15803d',
    color_hero_from: '#166534',
    color_hero_to: '#14532d',
  },
  seckills: {
    list: [
      {
        id: 1, type: 'seckill', name: '产地直供秒杀',
        start_at: '2026-05-20T00:00:00Z', end_at: '2026-06-20T23:59:59Z', status: 1,
        products: [
          { product_id: 401, title: '五常有机大米10斤', cover: farmCover(401), origin_price: 69.90, activity_price: 49.90, activity_stock: 200 },
          { product_id: 403, title: '秦岭土蜂蜜500g', cover: farmCover(403), origin_price: 58, activity_price: 38, activity_stock: 100 },
          { product_id: 405, title: '赣南脐橙10斤装', cover: farmCover(405), origin_price: 49.90, activity_price: 29.90, activity_stock: 300 },
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
          { product_id: 402, title: '西湖龙井明前茶250g', cover: farmCover(402), origin_price: 268, group_price: 198, group_stock: 80 },
          { product_id: 404, title: '云南野生菌菇礼盒', cover: farmCover(404), origin_price: 188, group_price: 138, group_stock: 50 },
          { product_id: 408, title: '湖南烟熏腊肉500g', cover: farmCover(408), origin_price: 45.90, group_price: 35.90, group_stock: 100 },
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
          { product_id: 406, title: '新疆和田大枣500g', cover: farmCover(406), origin_price: 35.90, floor_price: 0.01, max_helpers: 10, current_helpers: 3 },
          { product_id: 407, title: '东北黑木耳250g', cover: farmCover(407), origin_price: 29.90, floor_price: 5.90, max_helpers: 8, current_helpers: 5 },
        ],
      },
    ],
  },
  recommend: [
    { product_id: 405, title: '赣南脐橙10斤装', cover: farmCover(405), price: 49.90, origin_price: 69, sales: 12800 },
    { product_id: 401, title: '五常有机大米10斤', cover: farmCover(401), price: 69.90, origin_price: 98, sales: 8650 },
    { product_id: 406, title: '新疆和田大枣500g', cover: farmCover(406), price: 35.90, origin_price: 55, sales: 5680 },
    { product_id: 407, title: '东北黑木耳250g', cover: farmCover(407), price: 29.90, origin_price: 45, sales: 4350 },
  ],
  cart: [
    {
      sku_id: 402,
      qty: 1,
      product: {
        id: 401, title: '五常有机大米10斤', cover: farmCover(401), price: 69.90,
      },
      sku: {
        id: 402, product_id: 401, attrs: '[{"name":"规格","value":"10斤"}]', price: 69.90, stock: 700,
      },
    },
    {
      sku_id: 403,
      qty: 2,
      product: {
        id: 403, title: '秦岭土蜂蜜500g', cover: farmCover(403), price: 58.00,
      },
      sku: {
        id: 403, product_id: 403, attrs: '[{"name":"规格","value":"500g"}]', price: 58.00, stock: 800,
      },
    },
  ],
  orders: {
    list: [
      {
        id: 401, order_no: 'FARM20260501001', user_id: 1, status: 4,
        activity_type: 'seckill', activity_name: '产地直供秒杀',
        payment_method: 'wechat', goods_amount: 69.90, discount_amount: 20,
        freight_amount: 0, total_amount: 49.90, remark: '',
        tracking_no: 'SF1234567901',
        amount_breakdown: { goods_amount: 69.90, discount_amount: 20.00, freight_amount: 0.00, payable_amount: 49.90 },
        items: [
          {
            id: 4011, product_id: 401, title: '五常有机大米10斤', cover: farmCover(401), qty: 1, price: 69.90,
            review: {
              id: 4002, review_id: 4002, has_review: true, product_score: 5, logistics_score: 5,
              content: '米粒饱满晶莹，煮出来的饭特别香，全家都爱吃。', edited_times: 0,
              appends: [
                { id: 5002, content: '又回购了一袋，品质一如既往地好。', created_at: '2026-05-15T08:00:00Z' },
              ],
              admin_reply: { id: 6002, content: '感谢您的支持，我们坚持产地直发，保证新鲜！', created_at: '2026-05-16T09:00:00Z' },
              created_at: '2026-05-05T09:00:00Z', updated_at: '2026-05-15T08:00:00Z',
            },
          },
        ],
        created_at: '2026-05-01T08:30:00Z', paid_at: '2026-05-01T08:31:00Z',
      },
      {
        id: 402, order_no: 'FARM20260510002', user_id: 1, status: 3,
        activity_type: 'group_buy', activity_name: '3人拼团特惠',
        payment_method: 'alipay', goods_amount: 268.00, discount_amount: 70,
        freight_amount: 0, total_amount: 198.00, remark: '请用礼盒包装',
        tracking_no: 'SF9988776602',
        amount_breakdown: { goods_amount: 268.00, discount_amount: 70.00, freight_amount: 0.00, payable_amount: 198.00 },
        items: [
          { id: 4021, product_id: 402, title: '西湖龙井明前茶250g', cover: farmCover(402), qty: 1, price: 268.00 },
        ],
        created_at: '2026-05-10T14:20:00Z', paid_at: '2026-05-10T14:22:00Z',
      },
      {
        id: 403, order_no: 'FARM20260520003', user_id: 1, status: 1,
        activity_type: '', activity_name: '',
        payment_method: '', goods_amount: 188.00, discount_amount: 0,
        freight_amount: 0, total_amount: 188.00, remark: '',
        amount_breakdown: { goods_amount: 188.00, discount_amount: 0.00, freight_amount: 0.00, payable_amount: 188.00 },
        items: [
          { id: 4031, product_id: 404, title: '云南野生菌菇礼盒', cover: farmCover(404), qty: 1, price: 188.00 },
        ],
        created_at: '2026-05-20T09:00:00Z',
      },
      {
        id: 404, order_no: 'FARM20260521004', user_id: 1, status: 2,
        activity_type: 'bargain', activity_name: '砍价免费拿',
        payment_method: 'wechat', goods_amount: 35.90, discount_amount: 35.89,
        freight_amount: 0, total_amount: 0.01, remark: '',
        tracking_no: '',
        amount_breakdown: { goods_amount: 35.90, discount_amount: 35.89, freight_amount: 0.00, payable_amount: 0.01 },
        items: [
          { id: 4041, product_id: 406, title: '新疆和田大枣500g', cover: farmCover(406), qty: 1, price: 35.90 },
        ],
        created_at: '2026-05-21T16:00:00Z', paid_at: '2026-05-21T16:02:00Z',
      },
    ],
    total: 4,
    page: 1,
    size: 20,
  },
}
