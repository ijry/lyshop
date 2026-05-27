import { createRouter, createWebHistory, createWebHashHistory } from 'vue-router'

const isMock = import.meta.env.VITE_MOCK === 'true'

const router = createRouter({
  history: isMock ? createWebHashHistory() : createWebHistory('/admin'),
  routes: [
    { path: '/login', component: () => import('@/views/Login.vue') },
    {
      path: '/',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', name: 'nav.home', component: () => import('@/views/Dashboard.vue'), meta: { titleKey: 'nav.home' } },
        // Product
        { path: 'product/list',    name: 'nav.productList', component: () => import('@/views/product/ProductList.vue'), meta: { titleKey: 'nav.productList' } },
        { path: 'product/category',name: 'nav.productCategory', component: () => import('@/views/product/CategoryList.vue'), meta: { titleKey: 'nav.productCategory' } },
        { path: 'product/form',    name: 'nav.productAdd', component: () => import('@/views/product/ProductForm.vue'), meta: { titleKey: 'nav.productAdd' } },
        { path: 'product/form/:id',name: 'nav.productEdit', component: () => import('@/views/product/ProductForm.vue'), meta: { titleKey: 'nav.productEdit' } },
        // Order
        { path: 'order/list', name: 'nav.orderList', component: () => import('@/views/order/OrderList.vue'), meta: { titleKey: 'nav.orderList' } },
        { path: 'order/detail/:id', name: 'nav.orderDetail', component: () => import('@/views/order/OrderDetail.vue'), meta: { titleKey: 'nav.orderDetail' } },
        { path: 'order/after-sale/list', name: 'nav.afterSaleList', component: () => import('@/views/order/AfterSaleList.vue'), meta: { titleKey: 'nav.afterSaleList' } },
        { path: 'order/after-sale/detail/:id', name: 'nav.afterSaleDetail', component: () => import('@/views/order/AfterSaleDetail.vue'), meta: { titleKey: 'nav.afterSaleDetail' } },
        { path: 'review/list', name: 'nav.reviewList', component: () => import('@/views/review/ReviewList.vue'), meta: { titleKey: 'nav.reviewList' } },
        // WMS
        { path: 'wms/stock', name: 'nav.stockManage', component: () => import('@/views/wms/StockList.vue'), meta: { titleKey: 'nav.stockManage' } },
        // Marketing
        { path: 'marketing/coupon',   name: 'nav.couponManage', component: () => import('@/views/marketing/CouponList.vue'), meta: { titleKey: 'nav.couponManage' } },
        { path: 'marketing/seckill/activity', name: 'nav.seckillActivityManage', component: () => import('@/views/marketing/SeckillActivityList.vue'), meta: { titleKey: 'nav.seckillActivityManage' } },
        { path: 'marketing/seckill/product', name: 'nav.seckillProductManage', component: () => import('@/views/marketing/SeckillProductManage.vue'), meta: { titleKey: 'nav.seckillProductManage' } },
        { path: 'marketing/group-buy/activity', name: 'nav.groupBuyActivityManage', component: () => import('@/views/marketing/GroupBuyActivityList.vue'), meta: { titleKey: 'nav.groupBuyActivityManage' } },
        { path: 'marketing/group-buy/product', name: 'nav.groupBuyProductManage', component: () => import('@/views/marketing/GroupBuyProductManage.vue'), meta: { titleKey: 'nav.groupBuyProductManage' } },
        { path: 'marketing/bargain/activity', name: 'nav.bargainActivityManage', component: () => import('@/views/marketing/BargainActivityList.vue'), meta: { titleKey: 'nav.bargainActivityManage' } },
        { path: 'marketing/bargain/product', name: 'nav.bargainProductManage', component: () => import('@/views/marketing/BargainProductManage.vue'), meta: { titleKey: 'nav.bargainProductManage' } },
        // VIP
        { path: 'vip/plans', name: 'nav.vipPlan', component: () => import('@/views/vip/PlanList.vue'), meta: { titleKey: 'nav.vipPlan' } },
        { path: 'vip/levels', name: 'nav.vipLevel', component: () => import('@/views/vip/LevelList.vue'), meta: { titleKey: 'nav.vipLevel' } },
        { path: 'vip/coupon-rules', name: 'nav.vipCouponRule', component: () => import('@/views/vip/CouponRuleList.vue'), meta: { titleKey: 'nav.vipCouponRule' } },
        { path: 'vip/sku-prices', name: 'nav.vipSkuPrice', component: () => import('@/views/vip/SkuPriceList.vue'), meta: { titleKey: 'nav.vipSkuPrice' } },
        // System
        { path: 'system/site',   name: 'nav.siteSettings', component: () => import('@/views/system/SiteSettings.vue'), meta: { titleKey: 'nav.siteSettings' } },
        { path: 'system/config', name: 'nav.configCenter', component: () => import('@/views/system/PaymentConfig.vue'), meta: { titleKey: 'nav.configCenter' } },
        { path: 'system/admins', name: 'nav.adminManage', component: () => import('@/views/system/AdminList.vue'), meta: { titleKey: 'nav.adminManage' } },
        { path: 'system/roles',  name: 'nav.roleManage', component: () => import('@/views/system/RoleList.vue'), meta: { titleKey: 'nav.roleManage' } },
        // IM
        { path: 'im/sessions',  name: 'nav.imSession', component: () => import('@/views/im/SessionList.vue'), meta: { titleKey: 'nav.imSession' } },
        // AI
        { path: 'ai/tasks',     name: 'nav.aiImageGen',   component: () => import('@/views/ai/ImageGen.vue'), meta: { titleKey: 'nav.aiImageGen' } },
        { path: 'ai/models',    name: 'nav.aiModelConfig', component: () => import('@/views/ai/ImageGen.vue'), meta: { titleKey: 'nav.aiModelConfig' } },
        // Decor
        { path: 'decor/index',  name: 'nav.decorEditor', component: () => import('@/views/decor/DecorEditor.vue'), meta: { titleKey: 'nav.decorEditor' } },
        { path: 'decor/pc',    name: 'PC首页装修', component: () => import('@/views/decor/PcDecorEditor.vue'), meta: { titleKey: 'nav.pcDecorEditor' } },
        // Checkin
        { path: 'checkin/rules', name: 'nav.checkinRules', component: () => import('@/views/checkin/CheckinRules.vue'), meta: { titleKey: 'nav.checkinRules' } },
        { path: 'checkin/logs',  name: 'nav.checkinLogs', component: () => import('@/views/checkin/CheckinLogs.vue'), meta: { titleKey: 'nav.checkinLogs' } },
        // Message
        { path: 'message/list', name: 'nav.messageList', component: () => import('@/views/message/MessageList.vue'), meta: { titleKey: 'nav.messageList' } },
        { path: 'message/send', name: 'nav.messageSend', component: () => import('@/views/message/MessageSend.vue'), meta: { titleKey: 'nav.messageSend' } },
      ]
    }
  ]
})

router.beforeEach(to => {
  const token = localStorage.getItem('admin_token')
  if (to.meta.requiresAuth && !token) return '/login'
})

export default router
