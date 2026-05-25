import type { MockPreset } from './types'

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
      { id: 201, title: '智利车厘子5斤装', subtitle: '进口JJ级大果，新鲜直达', cover: 'https://picsum.photos/400/400?random=300', price: 168.00, origin_price: 238.00, stock: 300, sales: 4560, status: 1, category_id: 201 },
      { id: 202, title: '有机西蓝花2颗装', subtitle: '有机认证，翠绿新鲜', cover: 'https://picsum.photos/400/400?random=301', price: 15.90, origin_price: 22.00, stock: 800, sales: 3210, status: 1, category_id: 202 },
      { id: 203, title: '鲜活大闸蟹礼盒', subtitle: '阳澄湖直发，只只肥美', cover: 'https://picsum.photos/400/400?random=302', price: 399.00, origin_price: 558.00, stock: 150, sales: 1280, status: 1, category_id: 203 },
      { id: 204, title: '澳洲安格斯牛排200g×4', subtitle: '原切谷饲，入口即化', cover: 'https://picsum.photos/400/400?random=303', price: 259.00, origin_price: 358.00, stock: 200, sales: 2350, status: 1, category_id: 204 },
      { id: 205, title: '土鸡蛋30枚', subtitle: '散养走地鸡，蛋黄饱满', cover: 'https://picsum.photos/400/400?random=304', price: 39.90, origin_price: 55.00, stock: 1000, sales: 8900, status: 1, category_id: 204 },
      { id: 206, title: '挪威三文鱼刺身200g', subtitle: '冰鲜空运，刺身级品质', cover: 'https://picsum.photos/400/400?random=305', price: 89.90, origin_price: 128.00, stock: 400, sales: 3670, status: 1, category_id: 203 },
      { id: 207, title: '海南金煌芒果10斤', subtitle: '树上熟大果，香甜多汁', cover: 'https://picsum.photos/400/400?random=306', price: 59.90, origin_price: 89.00, stock: 600, sales: 6540, status: 1, category_id: 201 },
      { id: 208, title: '有机菠菜500g', subtitle: '有机种植，鲜嫩翠绿', cover: 'https://picsum.photos/400/400?random=307', price: 9.90, origin_price: 15.00, stock: 1500, sales: 5120, status: 1, category_id: 202 },
    ],
    total: 8,
    page: 1,
    size: 20,
  },

  productDetail: {
    id: 201,
    title: '智利车厘子5斤装',
    subtitle: '进口JJ级大果，果径28-30mm，新鲜直达',
    cover: 'https://picsum.photos/750/750?random=320',
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
        { id: 'b2', type: 'image', props: { url: 'https://picsum.photos/750/420?random=330', alt: '车厘子实拍' } },
        { id: 'b3', type: 'text', props: { text: '收货后请冷藏保存，建议3天内食用完毕。如有破损可申请售后理赔。' } },
      ],
    },
    skus: [
      { id: 201, product_id: 201, attrs: '[{"name":"重量","value":"3斤装"}]', price: 108.00, stock: 100, sku_code: 'CHERRY-3' },
      { id: 202, product_id: 201, attrs: '[{"name":"重量","value":"5斤装"}]', price: 168.00, stock: 150, sku_code: 'CHERRY-5' },
      { id: 203, product_id: 201, attrs: '[{"name":"重量","value":"10斤装"}]', price: 318.00, stock: 50, sku_code: 'CHERRY-10' },
    ],
    images: [
      { id: 201, product_id: 201, url: 'https://picsum.photos/750/750?random=320', sort: 0 },
      { id: 202, product_id: 201, url: 'https://picsum.photos/750/750?random=321', sort: 1 },
      { id: 203, product_id: 201, url: 'https://picsum.photos/750/750?random=322', sort: 2 },
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

  seckills: {
    list: [
      {
        id: 201, type: 'seckill', name: '生鲜限时秒杀',
        start_at: '2026-05-20T00:00:00Z', end_at: '2026-06-20T23:59:59Z', status: 1,
        products: [
          { product_id: 201, title: '智利车厘子5斤装', cover: 'https://picsum.photos/400/400?random=300', origin_price: 168.00, activity_price: 128.00, activity_stock: 80 },
          { product_id: 207, title: '海南金煌芒果10斤', cover: 'https://picsum.photos/400/400?random=306', origin_price: 59.90, activity_price: 39.90, activity_stock: 150 },
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
          { product_id: 204, title: '澳洲安格斯牛排200g×4', cover: 'https://picsum.photos/400/400?random=303', origin_price: 259.00, group_price: 199.00, group_stock: 100 },
          { product_id: 205, title: '土鸡蛋30枚', cover: 'https://picsum.photos/400/400?random=304', origin_price: 39.90, group_price: 29.90, group_stock: 300 },
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
          { product_id: 206, title: '挪威三文鱼刺身200g', cover: 'https://picsum.photos/400/400?random=305', origin_price: 128.00, floor_price: 0.01, max_helpers: 20, current_helpers: 9 },
        ],
      },
    ],
  },

  recommend: [
    { product_id: 201, title: '智利车厘子5斤装', cover: 'https://picsum.photos/400/400?random=300', price: 168.00, origin_price: 238.00, sales: 4560 },
    { product_id: 205, title: '土鸡蛋30枚', cover: 'https://picsum.photos/400/400?random=304', price: 39.90, origin_price: 55.00, sales: 8900 },
    { product_id: 207, title: '海南金煌芒果10斤', cover: 'https://picsum.photos/400/400?random=306', price: 59.90, origin_price: 89.00, sales: 6540 },
    { product_id: 206, title: '挪威三文鱼刺身200g', cover: 'https://picsum.photos/400/400?random=305', price: 89.90, origin_price: 128.00, sales: 3670 },
  ],

  cart: [
    {
      sku_id: 202,
      qty: 1,
      product: { id: 201, title: '智利车厘子5斤装', cover: 'https://picsum.photos/200/200?random=340', price: 168.00 },
      sku: { id: 202, product_id: 201, attrs: '[{"name":"重量","value":"5斤装"}]', price: 168.00, stock: 150 },
    },
    {
      sku_id: 205,
      qty: 2,
      product: { id: 205, title: '土鸡蛋30枚', cover: 'https://picsum.photos/200/200?random=344', price: 39.90 },
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
            id: 2011, product_id: 201, title: '智利车厘子5斤装', cover: 'https://picsum.photos/200/200?random=340', qty: 1, price: 128.00,
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
          { id: 2012, product_id: 208, title: '有机菠菜500g', cover: 'https://picsum.photos/200/200?random=347', qty: 2, price: 9.90 },
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
            id: 2021, product_id: 204, title: '澳洲安格斯牛排200g×4', cover: 'https://picsum.photos/200/200?random=343', qty: 1, price: 199.00,
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
          { id: 2031, product_id: 206, title: '挪威三文鱼刺身200g', cover: 'https://picsum.photos/200/200?random=345', qty: 1, price: 89.90 },
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
          { id: 2041, product_id: 206, title: '挪威三文鱼刺身200g', cover: 'https://picsum.photos/200/200?random=345', qty: 1, price: 128.00 },
        ],
        created_at: '2026-05-21T16:00:00Z', paid_at: '2026-05-21T16:02:00Z',
      },
    ],
    total: 4,
    page: 1,
    size: 20,
  },
}
