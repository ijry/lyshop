import type { MockPreset, SiteSettings } from './types'

export const cake: MockPreset = {
  key: 'cake',
  name: '烘焙甜品',

  categories: [
    { id: 501, parent_id: 0, name: '生日蛋糕', icon: '', sort: 1, status: 1 },
    { id: 502, parent_id: 0, name: '甜品点心', icon: '', sort: 2, status: 1 },
    { id: 503, parent_id: 0, name: '面包烘焙', icon: '', sort: 3, status: 1 },
    { id: 504, parent_id: 0, name: '下午茶套装', icon: '', sort: 4, status: 1 },
    { id: 505, parent_id: 0, name: '节日礼盒', icon: '', sort: 5, status: 1 },
    { id: 506, parent_id: 501, name: '奶油蛋糕', icon: '', sort: 1, status: 1 },
    { id: 507, parent_id: 503, name: '欧式面包', icon: '', sort: 1, status: 1 },
  ],

  products: {
    list: [
      {
        id: 501, title: '草莓奶油生日蛋糕8寸', subtitle: '新鲜草莓搭配轻盈奶油，甜而不腻',
        cover: 'https://picsum.photos/400/400?random=600',
        price: 168.00, origin_price: 228.00, stock: 200, sales: 3260, status: 1, category_id: 501,
      },
      {
        id: 502, title: '法式马卡龙12粒装', subtitle: '进口原料手工制作，口感细腻',
        cover: 'https://picsum.photos/400/400?random=601',
        price: 98.00, origin_price: 138.00, stock: 500, sales: 5120, status: 1, category_id: 502,
      },
      {
        id: 503, title: '北海道吐司面包450g', subtitle: '日式绵密口感，早餐首选',
        cover: 'https://picsum.photos/400/400?random=602',
        price: 25.90, origin_price: 35.00, stock: 1000, sales: 8930, status: 1, category_id: 503,
      },
      {
        id: 504, title: '提拉米苏6寸', subtitle: '意式经典，浓郁咖啡香',
        cover: 'https://picsum.photos/400/400?random=603',
        price: 128.00, origin_price: 168.00, stock: 300, sales: 2180, status: 1, category_id: 502,
      },
      {
        id: 505, title: '蛋黄酥礼盒12枚', subtitle: '层层酥皮包裹咸香蛋黄',
        cover: 'https://picsum.photos/400/400?random=604',
        price: 68.00, origin_price: 98.00, stock: 800, sales: 6540, status: 1, category_id: 505,
      },
      {
        id: 506, title: '肉松小贝8个装', subtitle: '网红爆款，咸甜酥脆',
        cover: 'https://picsum.photos/400/400?random=605',
        price: 38.00, origin_price: 55.00, stock: 600, sales: 4210, status: 1, category_id: 503,
      },
      {
        id: 507, title: '抹茶千层蛋糕6寸', subtitle: '日式宇治抹茶，层层奶油夹心',
        cover: 'https://picsum.photos/400/400?random=606',
        price: 158.00, origin_price: 218.00, stock: 150, sales: 1870, status: 1, category_id: 502,
      },
      {
        id: 508, title: '芒果班戟6个装', subtitle: '新鲜芒果果肉，Q弹薄皮',
        cover: 'https://picsum.photos/400/400?random=607',
        price: 48.00, origin_price: 68.00, stock: 400, sales: 3640, status: 1, category_id: 502,
      },
    ],
    total: 8,
    page: 1,
    size: 20,
  },

  productDetail: {
    id: 501,
    title: '草莓奶油生日蛋糕8寸',
    subtitle: '新鲜草莓搭配轻盈奶油，每日现做现发',
    cover: 'https://picsum.photos/750/750?random=600',
    price: 168.00,
    origin_price: 228.00,
    stock: 200,
    sales: 3260,
    status: 1,
    category_id: 501,
    detail: {
      version: 1,
      blocks: [
        {
          id: 'b1',
          type: 'text',
          props: {
            text: '选用新西兰进口淡奶油与当季新鲜草莓，纯手工裱花制作。蛋糕胚采用低糖配方，口感松软绵密，甜而不腻。',
          },
        },
        {
          id: 'b2',
          type: 'image',
          props: {
            url: 'https://picsum.photos/750/420?random=608',
            alt: '草莓蛋糕制作工艺展示',
          },
        },
        {
          id: 'b3',
          type: 'text',
          props: {
            text: '每日限量供应，下单后2小时内新鲜配送。附赠生日蜡烛、餐盘套装及贺卡，适合生日聚会、纪念日等各种场景。',
          },
        },
      ],
    },
    skus: [
      { id: 501, product_id: 501, attrs: '[{"name":"尺寸","value":"6寸"}]', price: 128.00, stock: 80, sku_code: 'CAKE-STRAW-6' },
      { id: 502, product_id: 501, attrs: '[{"name":"尺寸","value":"8寸"}]', price: 168.00, stock: 80, sku_code: 'CAKE-STRAW-8' },
      { id: 503, product_id: 501, attrs: '[{"name":"尺寸","value":"10寸"}]', price: 228.00, stock: 40, sku_code: 'CAKE-STRAW-10' },
    ],
    images: [
      { id: 501, product_id: 501, url: 'https://picsum.photos/750/750?random=600', sort: 0 },
      { id: 502, product_id: 501, url: 'https://picsum.photos/750/750?random=601', sort: 1 },
      { id: 503, product_id: 501, url: 'https://picsum.photos/750/750?random=602', sort: 2 },
      { id: 504, product_id: 501, url: 'https://picsum.photos/750/750?random=603', sort: 3 },
    ],
  },

  indexDecor: {
    components: [
      {
        type: 'banner',
        id: 'demo_banner',
        props: {
          images: [
            { url: 'https://picsum.photos/750/340?random=600', link: '/pages/product/list?category_id=501' },
            { url: 'https://picsum.photos/750/340?random=601', link: '/pages/marketing/coupon?mode=claim' },
            { url: 'https://picsum.photos/750/340?random=602', link: '/pages/product/list?category_id=505' },
          ],
          height: 340,
        },
      },
      {
        type: 'category_nav',
        id: 'demo_nav',
        props: {
          items: [
            { title: '生日蛋糕', icon: '', link: '/pages/product/list?category_id=501' },
            { title: '甜品点心', icon: '', link: '/pages/product/list?category_id=502' },
            { title: '面包烘焙', icon: '', link: '/pages/product/list?category_id=503' },
            { title: '下午茶套装', icon: '', link: '/pages/product/list?category_id=504' },
            { title: '节日礼盒', icon: '', link: '/pages/product/list?category_id=505' },
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
            { text: '每日新鲜现做，2小时内送达', link: '' },
            { text: '下午茶套餐第二件半价', link: '/pages/product/list?category_id=504' },
            { text: '新品上线：法式千层可颂，限时尝鲜价', link: '/pages/product/list' },
          ],
          color: '#db2777',
          bgColor: '#fdf2f8',
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
        props: { source: 'hot', limit: 8, columns: 2, title: '甜蜜精选' },
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
            { url: 'https://picsum.photos/1200/420?random=610', link: '/products?category=501' },
            { url: 'https://picsum.photos/1200/420?random=611', link: '/products?category=502' },
            { url: 'https://picsum.photos/1200/420?random=612', link: '/products?category=505' },
          ],
          height: 420,
        },
      },
      {
        type: 'category_nav',
        id: 'pc_cats',
        props: {
          columns: 5,
          items: [
            { title: '生日蛋糕', icon: '', link: '/products?category=501' },
            { title: '甜品点心', icon: '', link: '/products?category=502' },
            { title: '面包烘焙', icon: '', link: '/products?category=503' },
            { title: '饮品奶茶', icon: '', link: '/products?category=504' },
            { title: '节日定制', icon: '', link: '/products?category=505' },
          ],
        },
      },
      {
        type: 'notice',
        id: 'pc_notice',
        props: {
          items: [
            { text: '每日现烤新鲜直送', link: '' },
            { text: '生日蛋糕免配送费', link: '' },
            { text: '会员积分换甜品', link: '' },
          ],
          color: '#db2777',
          bgColor: '#fdf2f8',
        },
      },
      {
        type: 'product_grid',
        id: 'pc_hot',
        props: {
          title: '甜蜜精选',
          source: 'hot',
          limit: 8,
          columns: 4,
        },
      },
      {
        type: 'marketing_zone',
        id: 'pc_seckill',
        props: {
          title: '限时优惠',
          subtitle: '甜蜜不打烊',
          bg_from: '#be185d',
          bg_to: '#db2777',
          more_link: '/products',
          products: [
            { product_id: 501, title: '草莓奶油生日蛋糕', cover: 'https://picsum.photos/400/400?random=600', origin_price: 298, activity_price: 198 },
            { product_id: 502, title: '法式马卡龙礼盒', cover: 'https://picsum.photos/400/400?random=601', origin_price: 168, activity_price: 118 },
            { product_id: 503, title: '北海道吐司', cover: 'https://picsum.photos/400/400?random=602', origin_price: 38, activity_price: 25.9 },
            { product_id: 504, title: '提拉米苏6寸', cover: 'https://picsum.photos/400/400?random=603', origin_price: 188, activity_price: 128 },
          ],
        },
      },
      {
        type: 'spacer',
        id: 'pc_spacer1',
        props: { height: 24 },
      },
      {
        type: 'product_grid',
        id: 'pc_new',
        props: {
          title: '新品尝鲜',
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
            { title: '2小时送达', icon: 'i-carbon-delivery-truck', desc: '冷链配送' },
            { title: '手工现做', icon: 'i-carbon-checkmark-outline', desc: '新鲜保证' },
            { title: '不满意重做', icon: 'i-carbon-renew', desc: '售后保障' },
            { title: '在线客服', icon: 'i-carbon-headset', desc: '定制咨询' },
          ],
        },
      },
    ],
  },

  siteSettings: {
    site_name: '甜蜜工坊',
    site_logo: '',
    seo_title: '甜蜜工坊 - 甜蜜手作 幸福味道',
    seo_keywords: '蛋糕,烘焙,甜品,面包,下午茶',
    seo_description: '手工烘焙甜品商城，新鲜现做，甜蜜配送到家。',
    icp: '',
    hero_badge: '每日现烤，新鲜直送',
    hero_title: '甜蜜手作\\n幸福味道',
    hero_subtitle: '精选原料，手工制作，每一口都是幸福的味道。',
    hero_btn_text: '选蛋糕',
    hero_btn_link: '/products',
    color_primary: '#db2777',
    color_primary_light: '#ec4899',
    color_primary_dark: '#be185d',
    color_bg_page: '#fdf2f8',
    color_bg_header: 'rgba(255,255,255,0.8)',
    color_bg_footer: '#fdf2f8',
    color_price: '#db2777',
    color_hero_from: '#be185d',
    color_hero_to: '#9d174d',
  },

  seckills: {
    list: [
      {
        id: 1, type: 'seckill', name: '甜品限时秒杀',
        start_at: '2026-05-20T00:00:00Z', end_at: '2026-06-20T23:59:59Z', status: 1,
        products: [
          { product_id: 501, title: '草莓奶油生日蛋糕8寸', cover: 'https://picsum.photos/400/400?random=600', origin_price: 168, activity_price: 128, activity_stock: 30 },
          { product_id: 502, title: '法式马卡龙12粒装', cover: 'https://picsum.photos/400/400?random=601', origin_price: 98, activity_price: 68, activity_stock: 100 },
          { product_id: 506, title: '肉松小贝8个装', cover: 'https://picsum.photos/400/400?random=605', origin_price: 38, activity_price: 25, activity_stock: 80 },
        ],
      },
    ],
    end_at: '2026-06-20T23:59:59Z',
  },

  groupBuy: {
    list: [
      {
        id: 1, type: 'group_buy', name: '甜蜜拼团',
        group_size: 3, expire_hours: 24,
        start_at: '2026-05-01T00:00:00Z', end_at: '2026-07-01T23:59:59Z', status: 1,
        products: [
          { product_id: 504, title: '提拉米苏6寸', cover: 'https://picsum.photos/400/400?random=603', origin_price: 128, group_price: 98, group_stock: 60 },
          { product_id: 507, title: '抹茶千层蛋糕6寸', cover: 'https://picsum.photos/400/400?random=606', origin_price: 158, group_price: 118, group_stock: 40 },
          { product_id: 508, title: '芒果班戟6个装', cover: 'https://picsum.photos/400/400?random=607', origin_price: 48, group_price: 35, group_stock: 100 },
        ],
      },
    ],
  },

  bargain: {
    list: [
      {
        id: 1, type: 'bargain', name: '砍价免费吃',
        start_at: '2026-05-01T00:00:00Z', end_at: '2026-07-01T23:59:59Z', status: 1,
        products: [
          { product_id: 505, title: '蛋黄酥礼盒12枚', cover: 'https://picsum.photos/400/400?random=604', origin_price: 68, floor_price: 0.01, max_helpers: 15, current_helpers: 6 },
          { product_id: 503, title: '北海道吐司面包450g', cover: 'https://picsum.photos/400/400?random=602', origin_price: 25.90, floor_price: 5, max_helpers: 10, current_helpers: 4 },
        ],
      },
    ],
  },

  recommend: [
    { product_id: 502, title: '法式马卡龙12粒装', cover: 'https://picsum.photos/400/400?random=601', price: 98, origin_price: 138, sales: 5120 },
    { product_id: 507, title: '抹茶千层蛋糕6寸', cover: 'https://picsum.photos/400/400?random=606', price: 158, origin_price: 218, sales: 1870 },
    { product_id: 505, title: '蛋黄酥礼盒12枚', cover: 'https://picsum.photos/400/400?random=604', price: 68, origin_price: 98, sales: 6540 },
    { product_id: 508, title: '芒果班戟6个装', cover: 'https://picsum.photos/400/400?random=607', price: 48, origin_price: 68, sales: 3640 },
  ],

  cart: [
    {
      sku_id: 502,
      qty: 1,
      product: {
        id: 501, title: '草莓奶油生日蛋糕8寸', cover: 'https://picsum.photos/200/200?random=600', price: 168.00,
      },
      sku: {
        id: 502, product_id: 501, attrs: '[{"name":"尺寸","value":"8寸"}]', price: 168.00, stock: 80,
      },
    },
    {
      sku_id: 504,
      qty: 2,
      product: {
        id: 502, title: '法式马卡龙12粒装', cover: 'https://picsum.photos/200/200?random=601', price: 98.00,
      },
      sku: {
        id: 504, product_id: 502, attrs: '[{"name":"口味","value":"经典混合"}]', price: 98.00, stock: 500,
      },
    },
  ],

  orders: {
    list: [
      {
        id: 501, order_no: 'CAKE20260501001', user_id: 1, status: 4,
        activity_type: 'seckill', activity_name: '甜品秒杀',
        payment_method: 'wechat', goods_amount: 266.00, discount_amount: 40,
        freight_amount: 0, total_amount: 226.00, remark: '',
        tracking_no: 'SF1234567801',
        amount_breakdown: { goods_amount: 266.00, discount_amount: 40.00, freight_amount: 0.00, payable_amount: 226.00 },
        items: [
          {
            id: 5011, product_id: 501, title: '草莓奶油生日蛋糕8寸', cover: 'https://picsum.photos/200/200?random=600', qty: 1, price: 168.00,
            review: {
              id: 7001, review_id: 7001, has_review: true, product_score: 5, logistics_score: 5,
              content: '蛋糕非常新鲜，奶油细腻不甜腻，草莓也很新鲜饱满，家人都很喜欢！',
              edited_times: 0,
              appends: [
                { id: 8001, content: '放冰箱第二天吃依然很好吃，奶油不会塌。', created_at: '2026-05-03T10:00:00Z' },
              ],
              admin_reply: { id: 9001, content: '感谢您的好评，期待下次光临！', created_at: '2026-05-04T12:00:00Z' },
              created_at: '2026-05-01T18:00:00Z',
              updated_at: '2026-05-03T10:00:00Z',
            },
          },
          {
            id: 5012, product_id: 502, title: '法式马卡龙12粒装', cover: 'https://picsum.photos/200/200?random=601', qty: 1, price: 98.00,
          },
        ],
        created_at: '2026-05-01T10:30:00Z', paid_at: '2026-05-01T10:31:00Z',
      },
      {
        id: 502, order_no: 'CAKE20260510002', user_id: 1, status: 3,
        activity_type: 'group_buy', activity_name: '甜蜜拼团',
        payment_method: 'alipay', goods_amount: 286.00, discount_amount: 30,
        freight_amount: 0, total_amount: 256.00, remark: '请注意保温配送',
        tracking_no: 'SF9988776601',
        amount_breakdown: { goods_amount: 286.00, discount_amount: 30.00, freight_amount: 0.00, payable_amount: 256.00 },
        items: [
          {
            id: 5021, product_id: 504, title: '提拉米苏6寸', cover: 'https://picsum.photos/200/200?random=603', qty: 1, price: 128.00,
            review: {
              id: 7011, review_id: 7011, has_review: true, product_score: 4, logistics_score: 4,
              content: '口感不错，咖啡味浓郁，就是稍微有点甜。',
              edited_times: 0,
              appends: [], admin_reply: null,
              created_at: '2026-05-12T14:20:00Z', updated_at: '2026-05-12T14:20:00Z',
            },
          },
          {
            id: 5022, product_id: 507, title: '抹茶千层蛋糕6寸', cover: 'https://picsum.photos/200/200?random=606', qty: 1, price: 158.00,
          },
        ],
        created_at: '2026-05-10T14:20:00Z', paid_at: '2026-05-10T14:22:00Z',
      },
      {
        id: 503, order_no: 'CAKE20260520003', user_id: 1, status: 1,
        activity_type: '', activity_name: '',
        payment_method: '', goods_amount: 228.00, discount_amount: 0,
        freight_amount: 10, total_amount: 238.00, remark: '',
        amount_breakdown: { goods_amount: 228.00, discount_amount: 0.00, freight_amount: 10.00, payable_amount: 238.00 },
        items: [
          {
            id: 5031, product_id: 501, title: '草莓奶油生日蛋糕10寸', cover: 'https://picsum.photos/200/200?random=600', qty: 1, price: 228.00,
          },
        ],
        created_at: '2026-05-20T09:00:00Z',
      },
      {
        id: 504, order_no: 'CAKE20260521004', user_id: 1, status: 2,
        activity_type: 'bargain', activity_name: '砍价免费吃',
        payment_method: 'wechat', goods_amount: 68.00, discount_amount: 67.99,
        freight_amount: 0, total_amount: 0.01, remark: '',
        tracking_no: '',
        amount_breakdown: { goods_amount: 68.00, discount_amount: 67.99, freight_amount: 0.00, payable_amount: 0.01 },
        items: [
          {
            id: 5041, product_id: 505, title: '蛋黄酥礼盒12枚', cover: 'https://picsum.photos/200/200?random=604', qty: 1, price: 68.00,
          },
        ],
        created_at: '2026-05-21T16:00:00Z', paid_at: '2026-05-21T16:02:00Z',
      },
    ],
    total: 4,
    page: 1,
    size: 20,
  },
}
