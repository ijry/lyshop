import type { MockPreset, SiteSettings } from './types'

const motherImageMap: Record<number, { cover: string; gallery: string[] }> = {
  601: {
    cover: 'https://images.unsplash.com/photo-1582735689369-4fe89db7114c?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1582735689369-4fe89db7114c?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1515488042361-ee00e0ddd4e4?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1590080877777-95f6bd7696cf?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  602: {
    cover: 'https://images.unsplash.com/photo-1622290291468-a28f7a7dc6a8?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1622290291468-a28f7a7dc6a8?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1584362917165-526a968579e8?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1519689680058-324335c77eba?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  603: {
    cover: 'https://images.unsplash.com/photo-1515488042361-ee00e0ddd4e4?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1515488042361-ee00e0ddd4e4?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1566576912321-d58ddd7a6088?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1545558014-8692077e9b5c?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  604: {
    cover: 'https://images.unsplash.com/photo-1587854692152-cbe660dbde88?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1587854692152-cbe660dbde88?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1580281657527-47d4d54f5d0c?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1578496781379-7dcfb995293d?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  605: {
    cover: 'https://images.unsplash.com/photo-1519345182560-3f2917c472ef?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1519345182560-3f2917c472ef?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1514090458221-65bb69cf63e6?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1512436991641-6745cdb1723f?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  606: {
    cover: 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1542291026-7eec264c27ff?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1460353581641-37baddab0fa2?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1525966222134-fcfa99b8ae77?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  607: {
    cover: 'https://images.unsplash.com/photo-1514986888952-8cd320577b68?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1514986888952-8cd320577b68?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1473093295043-cdd812d0e601?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1505576399279-565b52d4ac71?auto=format&fit=crop&w=1200&q=80',
    ],
  },
  608: {
    cover: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?auto=format&fit=crop&w=800&q=80',
    gallery: [
      'https://images.unsplash.com/photo-1517841905240-472988babdf9?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1503454537195-1dcabb73ffb9?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1516627145497-ae6968895b74?auto=format&fit=crop&w=1200&q=80',
    ],
  },
}

function motherCover(id: number) {
  return motherImageMap[id]?.cover || motherImageMap[601].cover
}

function motherGallery(id: number, index = 0) {
  const list = motherImageMap[id]?.gallery || motherImageMap[601].gallery
  return list[index] || list[0]
}

export const mother: MockPreset = {
  key: 'mother',
  name: '母婴亲子',

  categories: [
    { id: 601, parent_id: 0, name: '奶粉辅食', icon: '', sort: 1, status: 1 },
    { id: 602, parent_id: 0, name: '纸尿裤', icon: '', sort: 2, status: 1 },
    { id: 603, parent_id: 0, name: '玩具早教', icon: '', sort: 3, status: 1 },
    { id: 604, parent_id: 0, name: '孕产用品', icon: '', sort: 4, status: 1 },
    { id: 605, parent_id: 0, name: '童装童鞋', icon: '', sort: 5, status: 1 },
    { id: 606, parent_id: 601, name: '进口奶粉', icon: '', sort: 1, status: 1 },
    { id: 607, parent_id: 603, name: '益智玩具', icon: '', sort: 1, status: 1 },
  ],

  products: {
    list: [
      {
        id: 601, title: '爱他美卓萃3段奶粉900g', subtitle: '德国原装进口，含天然乳脂',
        cover: motherCover(601),
        price: 298.00, origin_price: 368.00, stock: 800, sales: 6520, status: 1, category_id: 601,
      },
      {
        id: 602, title: '花王妙而舒纸尿裤L码', subtitle: '日本进口，柔软透气不闷热',
        cover: motherCover(602),
        price: 129.00, origin_price: 168.00, stock: 2000, sales: 12800, status: 1, category_id: 602,
      },
      {
        id: 603, title: '费雪益智积木80粒', subtitle: '安全无毒材质，激发宝宝创造力',
        cover: motherCover(603),
        price: 89.00, origin_price: 128.00, stock: 500, sales: 3450, status: 1, category_id: 603,
      },
      {
        id: 604, title: '斯利安叶酸片90片', subtitle: '备孕孕期必备，每片0.4mg',
        cover: motherCover(604),
        price: 59.00, origin_price: 79.00, stock: 3000, sales: 9870, status: 1, category_id: 604,
      },
      {
        id: 605, title: '纯棉婴儿连体衣', subtitle: 'A类纯棉面料，亲肤舒适',
        cover: motherCover(605),
        price: 49.90, origin_price: 79.00, stock: 1500, sales: 7620, status: 1, category_id: 605,
      },
      {
        id: 606, title: '基诺浦学步鞋', subtitle: '专业学步设计，保护宝宝小脚',
        cover: motherCover(606),
        price: 169.00, origin_price: 228.00, stock: 400, sales: 2890, status: 1, category_id: 605,
      },
      {
        id: 607, title: '小白熊辅食机', subtitle: '蒸煮搅一体，轻松制作营养辅食',
        cover: motherCover(607),
        price: 259.00, origin_price: 358.00, stock: 300, sales: 1560, status: 1, category_id: 601,
      },
      {
        id: 608, title: '好孩子轻便婴儿推车', subtitle: '一键折叠，轻便出行',
        cover: motherCover(608),
        price: 599.00, origin_price: 899.00, stock: 200, sales: 980, status: 1, category_id: 603,
      },
    ],
    total: 8,
    page: 1,
    size: 20,
  },

  productDetail: {
    id: 601,
    title: '爱他美卓萃3段奶粉900g',
    subtitle: '德国原装进口，天然乳脂配方，助力宝宝全面发育',
    cover: motherCover(601),
    price: 298.00,
    origin_price: 368.00,
    stock: 800,
    sales: 6520,
    status: 1,
    category_id: 601,
    detail: {
      version: 1,
      blocks: [
        {
          id: 'b1',
          type: 'text',
          props: {
            text: '爱他美卓萃精萃天然乳脂，含有DHA、ARA等关键营养元素，科学配比助力宝宝大脑和视力发育。独特的GOS/FOS益生元组合，呵护宝宝娇嫩肠道。',
          },
        },
        {
          id: 'b2',
          type: 'image',
          props: {
            url: motherGallery(601, 1),
            alt: '奶粉营养成分展示',
          },
        },
        {
          id: 'b3',
          type: 'text',
          props: {
            text: '德国原装原罐进口，全程冷链运输。每罐均有溯源码，扫码即可查询奶源地、生产日期及进口信息，品质安心有保障。',
          },
        },
      ],
    },
    skus: [
      { id: 601, product_id: 601, attrs: '[{"name":"段数","value":"1段(0-6月)"}]', price: 318.00, stock: 200, sku_code: 'MILK-APT-1' },
      { id: 602, product_id: 601, attrs: '[{"name":"段数","value":"2段(6-12月)"}]', price: 308.00, stock: 300, sku_code: 'MILK-APT-2' },
      { id: 603, product_id: 601, attrs: '[{"name":"段数","value":"3段(1-3岁)"}]', price: 298.00, stock: 300, sku_code: 'MILK-APT-3' },
    ],
    images: [
      { id: 601, product_id: 601, url: motherGallery(601, 0), sort: 0 },
      { id: 602, product_id: 601, url: motherGallery(601, 1), sort: 1 },
      { id: 603, product_id: 601, url: motherGallery(601, 2), sort: 2 },
      { id: 604, product_id: 601, url: motherGallery(602, 0), sort: 3 },
    ],
  },

  indexDecor: {
    components: [
      {
        type: 'banner',
        id: 'demo_banner',
        props: {
          images: [
            { url: motherGallery(601, 0), link: '/pages/product/list?category_id=601' },
            { url: motherGallery(602, 0), link: '/pages/marketing/coupon?mode=claim' },
            { url: motherGallery(603, 0), link: '/pages/product/list?category_id=603' },
          ],
          height: 340,
        },
      },
      {
        type: 'category_nav',
        id: 'demo_nav',
        props: {
          items: [
            { title: '奶粉辅食', icon: '', link: '/pages/product/list?category_id=601' },
            { title: '纸尿裤', icon: '', link: '/pages/product/list?category_id=602' },
            { title: '玩具早教', icon: '', link: '/pages/product/list?category_id=603' },
            { title: '孕产用品', icon: '', link: '/pages/product/list?category_id=604' },
            { title: '童装童鞋', icon: '', link: '/pages/product/list?category_id=605' },
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
            { text: '正品溯源，每罐可查', link: '' },
            { text: '新用户首单满199减30', link: '/pages/marketing/coupon?mode=claim' },
            { text: '纸尿裤囤货季，整箱购更优惠', link: '/pages/product/list?category_id=602' },
          ],
          color: '#7c3aed',
          bgColor: '#faf5ff',
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
        props: { source: 'hot', limit: 8, columns: 2, title: '宝妈精选' },
      },
    ],
  },
  pcDecor: {
    components: [
      {
        type: 'hero', id: 'pc_hero',
        props: {
          badge: '大牌正品，安心之选', title: '安心之选\n呵护成长',
          subtitle: '全球母婴大牌精选，品质保障，呵护宝宝每一步成长。',
          btn_text: '为宝宝选购', btn_link: '/products', btn2_text: '查看全部', btn2_link: '/products',
          bg_from: '#6d28d9', bg_to: '#5b21b6',
        },
      },
      {
        type: 'category_nav', id: 'pc_cats',
        props: {
          style: 'floating', columns: 5,
          items: [
            { title: '奶粉辅食', icon: '', link: '/products?category=601' },
            { title: '纸尿裤', icon: '', link: '/products?category=602' },
            { title: '玩具早教', icon: '', link: '/products?category=603' },
            { title: '童装童鞋', icon: '', link: '/products?category=604' },
            { title: '孕妈用品', icon: '', link: '/products?category=605' },
          ],
        },
      },
      {
        type: 'grid', id: 'pc_grid',
        props: {
          columns: 4,
          items: [
            { title: '秒杀', icon: '⚡', bg: '#faf5ff', link: '/products' },
            { title: '拼团', icon: '👥', bg: '#eff6ff', link: '/products' },
            { title: '优惠券', icon: '🏷️', bg: '#fff7ed', link: '/products' },
            { title: '签到', icon: '📅', bg: '#fdf2f8', link: '/products' },
          ],
        },
      },
      {
        type: 'notice', id: 'pc_notice',
        props: {
          items: [
            { text: '大牌正品安心购，全球母婴品牌直采', link: '' },
            { text: '新人专享：满199减30，首单包邮', link: '/products' },
            { text: '纸尿裤整箱囤货更划算，满299包邮', link: '/products' },
          ],
          color: '#7c3aed', bgColor: '#faf5ff',
        },
      },
      { type: 'product_grid', id: 'pc_hot', props: { title: '宝妈精选', source: 'hot', limit: 8, columns: 4 } },
      { type: 'image_ad', id: 'pc_ad1', props: { url: motherGallery(601, 0), link: '/products?category=601', height: 220 } },
      { type: 'product_grid', id: 'pc_new', props: { title: '新品上市', source: 'new', limit: 4, columns: 4 } },
      {
        type: 'features', id: 'pc_features',
        props: {
          columns: 4,
          items: [
            { icon: 'i-carbon-delivery-truck', title: '当日达', desc: '极速配送' },
            { icon: 'i-carbon-checkmark-outline', title: '正品保障', desc: '全球直采' },
            { icon: 'i-carbon-renew', title: '无忧退换', desc: '7天无理由' },
            { icon: 'i-carbon-headset', title: '育儿顾问', desc: '专业咨询' },
          ],
        },
      },
    ],
  },

  siteSettings: {
    site_name: '宝贝优选',
    site_logo: '',
    seo_title: '宝贝优选 - 安心之选 呵护成长',
    seo_keywords: '母婴,奶粉,纸尿裤,婴儿用品,童装',
    seo_description: '母婴亲子购物平台，全球大牌精选，安心呵护宝宝成长。',
    icp: '',
    hero_badge: '大牌正品，安心之选',
    hero_title: '安心之选\\n呵护成长',
    hero_subtitle: '全球母婴大牌精选，品质保障，呵护宝宝每一步成长。',
    hero_btn_text: '为宝宝选购',
    hero_btn_link: '/products',
    color_primary: '#7c3aed',
    color_primary_light: '#8b5cf6',
    color_primary_dark: '#6d28d9',
    color_bg_page: '#faf5ff',
    color_bg_header: 'rgba(255,255,255,0.8)',
    color_bg_footer: '#faf5ff',
    color_price: '#7c3aed',
    color_hero_from: '#6d28d9',
    color_hero_to: '#5b21b6',
  },

  seckills: {
    list: [
      {
        id: 1, type: 'seckill', name: '母婴限时秒杀',
        start_at: '2026-05-20T00:00:00Z', end_at: '2026-06-20T23:59:59Z', status: 1,
        products: [
          { product_id: 602, title: '花王妙而舒纸尿裤L码', cover: motherCover(602), origin_price: 129, activity_price: 89, activity_stock: 200 },
          { product_id: 605, title: '纯棉婴儿连体衣', cover: motherCover(605), origin_price: 49.90, activity_price: 29.90, activity_stock: 150 },
          { product_id: 603, title: '费雪益智积木80粒', cover: motherCover(603), origin_price: 89, activity_price: 59, activity_stock: 80 },
        ],
      },
    ],
    end_at: '2026-06-20T23:59:59Z',
  },

  groupBuy: {
    list: [
      {
        id: 1, type: 'group_buy', name: '宝妈拼团',
        group_size: 3, expire_hours: 24,
        start_at: '2026-05-01T00:00:00Z', end_at: '2026-07-01T23:59:59Z', status: 1,
        products: [
          { product_id: 601, title: '爱他美卓萃3段奶粉900g', cover: motherCover(601), origin_price: 298, group_price: 248, group_stock: 100 },
          { product_id: 607, title: '小白熊辅食机', cover: motherCover(607), origin_price: 259, group_price: 199, group_stock: 50 },
          { product_id: 606, title: '基诺浦学步鞋', cover: motherCover(606), origin_price: 169, group_price: 139, group_stock: 60 },
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
          { product_id: 603, title: '费雪益智积木80粒', cover: motherCover(603), origin_price: 89, floor_price: 0.01, max_helpers: 15, current_helpers: 7 },
          { product_id: 605, title: '纯棉婴儿连体衣', cover: motherCover(605), origin_price: 49.90, floor_price: 5, max_helpers: 10, current_helpers: 3 },
        ],
      },
    ],
  },

  recommend: [
    { product_id: 601, title: '爱他美卓萃3段奶粉900g', cover: motherCover(601), price: 298, origin_price: 368, sales: 6520 },
    { product_id: 602, title: '花王妙而舒纸尿裤L码', cover: motherCover(602), price: 129, origin_price: 168, sales: 12800 },
    { product_id: 606, title: '基诺浦学步鞋', cover: motherCover(606), price: 169, origin_price: 228, sales: 2890 },
    { product_id: 607, title: '小白熊辅食机', cover: motherCover(607), price: 259, origin_price: 358, sales: 1560 },
  ],

  cart: [
    {
      sku_id: 603,
      qty: 2,
      product: {
        id: 601, title: '爱他美卓萃3段奶粉900g', cover: motherCover(601), price: 298.00,
      },
      sku: {
        id: 603, product_id: 601, attrs: '[{"name":"段数","value":"3段(1-3岁)"}]', price: 298.00, stock: 300,
      },
    },
    {
      sku_id: 605,
      qty: 1,
      product: {
        id: 602, title: '花王妙而舒纸尿裤L码', cover: motherCover(602), price: 129.00,
      },
      sku: {
        id: 605, product_id: 602, attrs: '[{"name":"尺码","value":"L码(9-14kg)"}]', price: 129.00, stock: 2000,
      },
    },
  ],

  orders: {
    list: [
      {
        id: 601, order_no: 'MOM20260501001', user_id: 1, status: 4,
        activity_type: 'seckill', activity_name: '母婴秒杀',
        payment_method: 'wechat', goods_amount: 427.00, discount_amount: 50,
        freight_amount: 0, total_amount: 377.00, remark: '',
        tracking_no: 'SF1234567901',
        amount_breakdown: { goods_amount: 427.00, discount_amount: 50.00, freight_amount: 0.00, payable_amount: 377.00 },
        items: [
          {
            id: 6011, product_id: 601, title: '爱他美卓萃3段奶粉900g', cover: motherCover(601), qty: 1, price: 298.00,
            review: {
              id: 7002, review_id: 7002, has_review: true, product_score: 5, logistics_score: 5,
              content: '宝宝很爱喝，溶解性好，没有结块，正品有溯源码验证过了。',
              edited_times: 0,
              appends: [
                { id: 8002, content: '喝了一个月了，宝宝消化吸收都很好，准备继续囤货。', created_at: '2026-05-03T10:00:00Z' },
              ],
              admin_reply: { id: 9002, content: '感谢宝妈信赖，祝宝宝健康成长！', created_at: '2026-05-04T12:00:00Z' },
              created_at: '2026-05-01T18:00:00Z',
              updated_at: '2026-05-03T10:00:00Z',
            },
          },
          {
            id: 6012, product_id: 602, title: '花王妙而舒纸尿裤L码', cover: motherCover(602), qty: 1, price: 129.00,
          },
        ],
        created_at: '2026-05-01T10:30:00Z', paid_at: '2026-05-01T10:31:00Z',
      },
      {
        id: 602, order_no: 'MOM20260510002', user_id: 1, status: 3,
        activity_type: 'group_buy', activity_name: '宝妈拼团',
        payment_method: 'alipay', goods_amount: 348.00, discount_amount: 30,
        freight_amount: 0, total_amount: 318.00, remark: '请包装牢固',
        tracking_no: 'SF9988776602',
        amount_breakdown: { goods_amount: 348.00, discount_amount: 30.00, freight_amount: 0.00, payable_amount: 318.00 },
        items: [
          {
            id: 6021, product_id: 603, title: '费雪益智积木80粒', cover: motherCover(603), qty: 1, price: 89.00,
            review: {
              id: 7003, review_id: 7003, has_review: true, product_score: 4, logistics_score: 4,
              content: '积木颜色鲜艳，宝宝很喜欢玩，材质安全无异味。',
              edited_times: 0,
              appends: [], admin_reply: null,
              created_at: '2026-05-12T14:20:00Z', updated_at: '2026-05-12T14:20:00Z',
            },
          },
          {
            id: 6022, product_id: 607, title: '小白熊辅食机', cover: motherCover(607), qty: 1, price: 259.00,
          },
        ],
        created_at: '2026-05-10T14:20:00Z', paid_at: '2026-05-10T14:22:00Z',
      },
      {
        id: 603, order_no: 'MOM20260520003', user_id: 1, status: 1,
        activity_type: '', activity_name: '',
        payment_method: '', goods_amount: 599.00, discount_amount: 0,
        freight_amount: 0, total_amount: 599.00, remark: '',
        amount_breakdown: { goods_amount: 599.00, discount_amount: 0.00, freight_amount: 0.00, payable_amount: 599.00 },
        items: [
          {
            id: 6031, product_id: 608, title: '好孩子轻便婴儿推车', cover: motherCover(608), qty: 1, price: 599.00,
          },
        ],
        created_at: '2026-05-20T09:00:00Z',
      },
      {
        id: 604, order_no: 'MOM20260521004', user_id: 1, status: 2,
        activity_type: 'bargain', activity_name: '砍价免费拿',
        payment_method: 'wechat', goods_amount: 89.00, discount_amount: 88.99,
        freight_amount: 0, total_amount: 0.01, remark: '',
        tracking_no: '',
        amount_breakdown: { goods_amount: 89.00, discount_amount: 88.99, freight_amount: 0.00, payable_amount: 0.01 },
        items: [
          {
            id: 6041, product_id: 603, title: '费雪益智积木80粒', cover: motherCover(603), qty: 1, price: 89.00,
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
