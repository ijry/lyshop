import type { MockPreset, SiteSettings } from './types'

export const supermarket: MockPreset = {
  key: 'supermarket',
  name: '商超便利',

  categories: [
    { id: 101, parent_id: 0, name: '零食饮料', icon: '', sort: 1, status: 1 },
    { id: 102, parent_id: 0, name: '日用百货', icon: '', sort: 2, status: 1 },
    { id: 103, parent_id: 0, name: '个人护理', icon: '', sort: 3, status: 1 },
    { id: 104, parent_id: 0, name: '生鲜食品', icon: '', sort: 4, status: 1 },
    { id: 105, parent_id: 0, name: '粮油调味', icon: '', sort: 5, status: 1 },
    { id: 106, parent_id: 101, name: '碳酸饮料', icon: '', sort: 1, status: 1 },
    { id: 107, parent_id: 103, name: '洗护用品', icon: '', sort: 1, status: 1 },
  ],

  products: {
    list: [
      { id: 101, title: '可口可乐24罐装', subtitle: '经典口味，畅爽一夏', cover: 'https://picsum.photos/400/400?random=200', price: 49.90, origin_price: 69.90, stock: 2000, sales: 8560, status: 1, category_id: 101 },
      { id: 102, title: '蓝月亮洗衣液4L', subtitle: '深层洁净，持久留香', cover: 'https://picsum.photos/400/400?random=201', price: 39.90, origin_price: 59.90, stock: 1500, sales: 6230, status: 1, category_id: 102 },
      { id: 103, title: '维达抽纸10包装', subtitle: '柔韧亲肤，家庭必备', cover: 'https://picsum.photos/400/400?random=202', price: 29.90, origin_price: 45.00, stock: 3000, sales: 12400, status: 1, category_id: 102 },
      { id: 104, title: '乐事薯片大礼包', subtitle: '多口味混合装，追剧必备', cover: 'https://picsum.photos/400/400?random=203', price: 59.90, origin_price: 89.90, stock: 800, sales: 4320, status: 1, category_id: 101 },
      { id: 105, title: '舒肤佳沐浴露套装', subtitle: '温和清洁，滋润保湿', cover: 'https://picsum.photos/400/400?random=204', price: 35.90, origin_price: 55.00, stock: 1200, sales: 3780, status: 1, category_id: 103 },
      { id: 106, title: '康师傅方便面整箱', subtitle: '红烧牛肉面，24袋装', cover: 'https://picsum.photos/400/400?random=205', price: 44.90, origin_price: 59.90, stock: 1000, sales: 7650, status: 1, category_id: 101 },
      { id: 107, title: '特仑苏纯牛奶12盒', subtitle: '优质奶源，醇厚口感', cover: 'https://picsum.photos/400/400?random=206', price: 59.90, origin_price: 79.90, stock: 1800, sales: 9120, status: 1, category_id: 104 },
      { id: 108, title: '立白洗洁精1.5kg', subtitle: '强效去油，不伤手', cover: 'https://picsum.photos/400/400?random=207', price: 9.90, origin_price: 15.90, stock: 5000, sales: 15600, status: 1, category_id: 102 },
    ],
    total: 8,
    page: 1,
    size: 20,
  },

  productDetail: {
    id: 101,
    title: '可口可乐24罐装',
    subtitle: '经典口味，畅爽一夏，330ml×24罐',
    cover: 'https://picsum.photos/750/750?random=220',
    price: 49.90,
    origin_price: 69.90,
    stock: 2000,
    sales: 8560,
    status: 1,
    category_id: 101,
    detail: {
      version: 1,
      blocks: [
        { id: 'b1', type: 'text', props: { text: '可口可乐经典原味，330ml×24罐装。冰镇后饮用口感更佳，适合聚会、家庭囤货。' } },
        { id: 'b2', type: 'image', props: { url: 'https://picsum.photos/750/420?random=230', alt: '产品展示' } },
        { id: 'b3', type: 'text', props: { text: '保质期12个月，常温保存即可。建议冰镇至4-8℃饮用，口感最佳。' } },
      ],
    },
    skus: [
      { id: 101, product_id: 101, attrs: '[{"name":"规格","value":"24罐装"}]', price: 49.90, stock: 1500, sku_code: 'COLA-24' },
      { id: 102, product_id: 101, attrs: '[{"name":"规格","value":"12罐装"}]', price: 26.90, stock: 500, sku_code: 'COLA-12' },
    ],
    images: [
      { id: 101, product_id: 101, url: 'https://picsum.photos/750/750?random=220', sort: 0 },
      { id: 102, product_id: 101, url: 'https://picsum.photos/750/750?random=221', sort: 1 },
      { id: 103, product_id: 101, url: 'https://picsum.photos/750/750?random=222', sort: 2 },
    ],
  },

  indexDecor: {
    components: [
      {
        type: 'banner',
        id: 'demo_banner',
        props: {
          images: [
            { url: '/static/demo/banner-supermarket.png', link: '/pages/product/list?category_id=101' },
            { url: '/static/demo/banner-daily.png', link: '/pages/product/list?category_id=102' },
            { url: '/static/demo/banner-fresh.png', link: '/pages/product/list?category_id=104' },
          ],
          height: 340,
        },
      },
      {
        type: 'category_nav',
        id: 'demo_nav',
        props: {
          items: [
            { title: '零食', icon: '', link: '/pages/product/list?category_id=101' },
            { title: '百货', icon: '', link: '/pages/product/list?category_id=102' },
            { title: '护理', icon: '', link: '/pages/product/list?category_id=103' },
            { title: '生鲜', icon: '', link: '/pages/product/list?category_id=104' },
            { title: '粮油', icon: '', link: '/pages/product/list?category_id=105' },
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
            { text: '会员日全场满99减10，限时抢购中！', link: '' },
            { text: '新人首单立减5元，生鲜蔬果天天特价', link: '/pages/marketing/coupon?mode=claim' },
            { text: '整箱囤货更划算，满199包邮到家', link: '/pages/product/list' },
          ],
          color: '#f97316',
          bgColor: '#fff7ed',
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
        type: 'hero',
        id: 'pc_hero',
        props: {
          badge: '会员日全场满99减10',
          title: '新鲜到家\\n实惠每天',
          subtitle: '日用百货一站购齐，品质保障，1小时极速达。',
          btn_text: '去逛逛',
          btn_link: '/products',
          btn2_text: '全部分类',
          btn2_link: '/products',
          bg_from: '#ea580c',
          bg_to: '#c2410c',
        },
      },
      {
        type: 'category_nav',
        id: 'pc_cats',
        props: {
          style: 'floating',
          columns: 5,
          items: [
            { title: '零食饮料', link: '/products?category=101' },
            { title: '日用百货', link: '/products?category=102' },
            { title: '个护清洁', link: '/products?category=103' },
            { title: '生鲜果蔬', link: '/products?category=104' },
            { title: '粮油调味', link: '/products?category=105' },
          ],
        },
      },
      {
        type: 'grid',
        id: 'pc_grid',
        props: {
          columns: 4,
          items: [
            { title: '秒杀', icon: '⚡', bg: '#fef2f2', link: '/products' },
            { title: '拼团', icon: '👥', bg: '#eff6ff', link: '/products' },
            { title: '优惠券', icon: '🏷️', bg: '#fff7ed', link: '/products' },
            { title: '签到', icon: '📅', bg: '#faf5ff', link: '/products' },
          ],
        },
      },
      {
        type: 'notice',
        id: 'pc_notice',
        props: {
          items: [
            { text: '会员日满99减10' },
            { text: '新人立减5元' },
            { text: '整箱满199包邮' },
          ],
          color: '#f97316',
          bgColor: '#fff7ed',
        },
      },
      {
        type: 'product_grid',
        id: 'pc_hot',
        props: {
          title: '天天特价',
          source: 'hot',
          limit: 8,
          columns: 4,
        },
      },
      {
        type: 'image_ad',
        id: 'pc_ad1',
        props: {
          url: 'https://picsum.photos/1200/300?random=110',
          link: '/products',
          height: 200,
        },
      },
      {
        type: 'product_grid',
        id: 'pc_new',
        props: {
          title: '囤货专区',
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
            { title: '1小时达', icon: 'i-carbon-delivery-truck', desc: '急速配送' },
            { title: '品质保证', icon: 'i-carbon-checkmark-outline', desc: '正品行货' },
            { title: '无忧退换', icon: 'i-carbon-renew', desc: '7天退换' },
            { title: '在线客服', icon: 'i-carbon-headset', desc: '随时响应' },
          ],
        },
      },
    ],
  },

  siteSettings: {
    site_name: '惠选超市',
    site_logo: '',
    seo_title: '惠选超市 - 新鲜到家 实惠每天',
    seo_keywords: '超市,日用百货,零食饮料,生活用品',
    seo_description: '线上商超，日用百货一站购齐，新鲜到家。',
    icp: '',
    hero_badge: '会员日全场满99减10',
    hero_title: '新鲜到家\\n实惠每天',
    hero_subtitle: '日用百货一站购齐，品质保障，1小时极速达。',
    hero_btn_text: '去逛逛',
    hero_btn_link: '/products',
    color_primary: '#f97316',
    color_primary_light: '#fb923c',
    color_primary_dark: '#ea580c',
    color_bg_page: '#fffbf5',
    color_bg_header: 'rgba(255,255,255,0.8)',
    color_bg_footer: '#fffbf5',
    color_price: '#f97316',
    color_hero_from: '#ea580c',
    color_hero_to: '#c2410c',
  },

  seckills: {
    list: [
      {
        id: 101, type: 'seckill', name: '超市限时秒杀',
        start_at: '2026-05-20T00:00:00Z', end_at: '2026-06-20T23:59:59Z', status: 1,
        products: [
          { product_id: 101, title: '可口可乐24罐装', cover: 'https://picsum.photos/400/400?random=200', origin_price: 49.90, activity_price: 35.90, activity_stock: 200 },
          { product_id: 106, title: '康师傅方便面整箱', cover: 'https://picsum.photos/400/400?random=205', origin_price: 44.90, activity_price: 29.90, activity_stock: 150 },
        ],
      },
    ],
    end_at: '2026-06-20T23:59:59Z',
  },

  groupBuy: {
    list: [
      {
        id: 101, type: 'group_buy', name: '3人拼团特惠',
        group_size: 3, expire_hours: 24,
        start_at: '2026-05-01T00:00:00Z', end_at: '2026-07-01T23:59:59Z', status: 1,
        products: [
          { product_id: 107, title: '特仑苏纯牛奶12盒', cover: 'https://picsum.photos/400/400?random=206', origin_price: 59.90, group_price: 45.90, group_stock: 300 },
          { product_id: 104, title: '乐事薯片大礼包', cover: 'https://picsum.photos/400/400?random=203', origin_price: 59.90, group_price: 39.90, group_stock: 200 },
        ],
      },
    ],
  },

  bargain: {
    list: [
      {
        id: 101, type: 'bargain', name: '砍价免费拿',
        start_at: '2026-05-01T00:00:00Z', end_at: '2026-07-01T23:59:59Z', status: 1,
        products: [
          { product_id: 105, title: '舒肤佳沐浴露套装', cover: 'https://picsum.photos/400/400?random=204', origin_price: 55.00, floor_price: 0.01, max_helpers: 15, current_helpers: 6 },
        ],
      },
    ],
  },

  recommend: [
    { product_id: 101, title: '可口可乐24罐装', cover: 'https://picsum.photos/400/400?random=200', price: 49.90, origin_price: 69.90, sales: 8560 },
    { product_id: 103, title: '维达抽纸10包装', cover: 'https://picsum.photos/400/400?random=202', price: 29.90, origin_price: 45.00, sales: 12400 },
    { product_id: 107, title: '特仑苏纯牛奶12盒', cover: 'https://picsum.photos/400/400?random=206', price: 59.90, origin_price: 79.90, sales: 9120 },
    { product_id: 108, title: '立白洗洁精1.5kg', cover: 'https://picsum.photos/400/400?random=207', price: 9.90, origin_price: 15.90, sales: 15600 },
  ],

  cart: [
    {
      sku_id: 101,
      qty: 2,
      product: { id: 101, title: '可口可乐24罐装', cover: 'https://picsum.photos/200/200?random=240', price: 49.90 },
      sku: { id: 101, product_id: 101, attrs: '[{"name":"规格","value":"24罐装"}]', price: 49.90, stock: 1500 },
    },
    {
      sku_id: 103,
      qty: 1,
      product: { id: 103, title: '维达抽纸10包装', cover: 'https://picsum.photos/200/200?random=242', price: 29.90 },
      sku: { id: 103, product_id: 103, attrs: '[{"name":"规格","value":"10包装"}]', price: 29.90, stock: 3000 },
    },
  ],

  orders: {
    list: [
      {
        id: 101, order_no: 'SM20260501001', user_id: 1, status: 4,
        activity_type: 'seckill', activity_name: '超市秒杀',
        payment_method: 'wechat', goods_amount: 94.80, discount_amount: 14.00,
        freight_amount: 0, total_amount: 80.80, remark: '',
        tracking_no: 'YT1234567890',
        amount_breakdown: { goods_amount: 94.80, discount_amount: 14.00, freight_amount: 0.00, payable_amount: 80.80 },
        items: [
          {
            id: 1011, product_id: 101, title: '可口可乐24罐装', cover: 'https://picsum.photos/200/200?random=240', qty: 1, price: 35.90,
            review: {
              id: 1101, review_id: 1101, has_review: true, product_score: 5, logistics_score: 5,
              content: '价格实惠，日期很新鲜，冰镇后口感很好！', edited_times: 0,
              appends: [
                { id: 2101, content: '回购第三箱了，全家都爱喝，性价比超高。', created_at: '2026-05-03T10:00:00Z' },
              ],
              admin_reply: { id: 3101, content: '感谢您的支持，欢迎再次选购！', created_at: '2026-05-04T12:00:00Z' },
              created_at: '2026-05-01T10:31:00Z', updated_at: '2026-05-03T10:00:00Z',
            },
          },
          { id: 1012, product_id: 106, title: '康师傅方便面整箱', cover: 'https://picsum.photos/200/200?random=245', qty: 1, price: 44.90 },
        ],
        created_at: '2026-05-01T10:30:00Z', paid_at: '2026-05-01T10:31:00Z',
      },
      {
        id: 102, order_no: 'SM20260510002', user_id: 1, status: 3,
        activity_type: 'group_buy', activity_name: '3人拼团',
        payment_method: 'alipay', goods_amount: 59.90, discount_amount: 14.00,
        freight_amount: 0, total_amount: 45.90, remark: '尽快发货',
        tracking_no: 'SF6677889900',
        amount_breakdown: { goods_amount: 59.90, discount_amount: 14.00, freight_amount: 0.00, payable_amount: 45.90 },
        items: [
          {
            id: 1021, product_id: 107, title: '特仑苏纯牛奶12盒', cover: 'https://picsum.photos/200/200?random=246', qty: 1, price: 45.90,
            review: {
              id: 1102, review_id: 1102, has_review: true, product_score: 4, logistics_score: 4,
              content: '口感醇厚，日期新鲜。', edited_times: 0,
              appends: [], admin_reply: null,
              created_at: '2026-05-12T14:20:00Z', updated_at: '2026-05-12T14:20:00Z',
            },
          },
        ],
        created_at: '2026-05-10T14:20:00Z', paid_at: '2026-05-10T14:22:00Z',
      },
      {
        id: 103, order_no: 'SM20260520003', user_id: 1, status: 1,
        activity_type: '', activity_name: '',
        payment_method: '', goods_amount: 89.70, discount_amount: 0,
        freight_amount: 0, total_amount: 89.70, remark: '',
        amount_breakdown: { goods_amount: 89.70, discount_amount: 0.00, freight_amount: 0.00, payable_amount: 89.70 },
        items: [
          { id: 1031, product_id: 103, title: '维达抽纸10包装', cover: 'https://picsum.photos/200/200?random=242', qty: 3, price: 29.90 },
        ],
        created_at: '2026-05-20T09:00:00Z',
      },
      {
        id: 104, order_no: 'SM20260521004', user_id: 1, status: 2,
        activity_type: 'bargain', activity_name: '砍价免费拿',
        payment_method: 'wechat', goods_amount: 55.00, discount_amount: 54.99,
        freight_amount: 0, total_amount: 0.01, remark: '',
        tracking_no: '',
        amount_breakdown: { goods_amount: 55.00, discount_amount: 54.99, freight_amount: 0.00, payable_amount: 0.01 },
        items: [
          { id: 1041, product_id: 105, title: '舒肤佳沐浴露套装', cover: 'https://picsum.photos/200/200?random=244', qty: 1, price: 55.00 },
        ],
        created_at: '2026-05-21T16:00:00Z', paid_at: '2026-05-21T16:02:00Z',
      },
    ],
    total: 4,
    page: 1,
    size: 20,
  },
}
