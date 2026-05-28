export default {
  common: {
    refresh: '刷新', loading: '加载中…', empty: '暂无数据', save: '保存', cancel: '取消', confirm: '确认', reset: '重置',
  },
  login: {
    title: '商家工作台登录',
    username: '账号', password: '密码', submit: '登录',
    usernameRequired: '请输入账号', passwordRequired: '请输入密码',
  },
  dashboard: {
    title: '工作台',
    todayOrders: '今日订单', todaySales: '今日营收', todayAvgPrice: '客单价',
    pendingShip: '待发货', pendingAfterSale: '待审售后', unreadMessage: '未读消息',
    stockWarning: '库存预警', pendingInvoice: '待开发票', pendingRefund: '待处理退款',
    trend7: '近 7 日', trend30: '近 30 日',
    statusDistribution: '订单状态分布', hotProducts: '商品销量 Top5',
    quickEntries: '快捷入口', announcements: '公告',
  },
  order: {
    all: '全部', pendingPay: '待付款', pendingShip: '待发货', shipped: '已发货',
    completed: '已完成', closed: '已关闭', hasAfterSale: '售后中', pendingReview: '待评价', pendingInvoice: '待开票',
    filterTitle: '订单筛选',
    actions: { detail: '详情', reprice: '改价', remindPay: '催付', ship: '发货', note: '备注', print: '打单' },
    batch: { ship: '批量发货', notes: '批量备注', close: '批量关闭' },
  },
  afterSale: {
    all: '全部', applied: '待审核', returning: '退货中', refunding: '退款中', refunded: '已完成', closed: '已关闭',
    types: { all: '全部', refundOnly: '仅退款', returnRefund: '退货退款', exchange: '换货' },
    progress: { applied: '申请', approved: '审核', returning: '寄回', received: '收货', refunded: '退款' },
    chatPlaceholder: '回复买家',
    evidenceUpload: '上传凭证',
  },
  product: {
    search: '搜索商品', edit: '编辑', onSale: '上架', offSale: '下架',
    tabs: { all: '全部', onSale: '在售', off: '仓库', warning: '预警' },
    sortBy: { default: '默认', sales: '销量', stock: '库存', priceAsc: '价格升', priceDesc: '价格降', created: '最新' },
    batch: { shelfOn: '批量上架', shelfOff: '批量下架', category: '批量分类', price: '批量调价' },
    editSections: {
      base: '基础信息', covers: '主图轮播', detail: '商品详情', pricing: '价格库存',
      sku: '规格 SKU', category: '分类与标签', shipping: '物流与营销', status: '状态控制',
    },
  },
  marketing: { coupon: '优惠券', seckill: '秒杀', groupBuy: '拼团', bargain: '砍价', vip: 'VIP 会员' },
  vip: {
    plans: '会员套餐', levels: '会员等级', couponRules: '领券规则', skuPrices: 'SKU专属价',
    name: '名称', months: '时长（月）', price: '价格', status: '状态',
    growthMin: '最低成长值', discountRate: '折扣率',
    couponName: '优惠券名称', monthlyLimit: '每月限领',
    productId: '商品', skuId: 'SKU', levelId: '等级', vipPrice: '会员价',
  },
  me: {
    title: '我的', messages: '消息中心', sessions: '客服会话',
    siteSettings: '店铺设置', admins: '管理员', roles: '角色权限', logout: '退出登录',
  },
  biz: {
    selected: '已选 {count} 项',
    batchResult: '批量操作结果',
    selectFirst: '请先勾选',
    soon: '即将上线',
  },
}
