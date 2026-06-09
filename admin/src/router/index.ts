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
        { path: 'product/spec-templates',name: 'nav.productSpecTemplates', component: () => import('@/views/product/SpecTemplateList.vue'), meta: { titleKey: 'nav.productSpecTemplates' } },
        { path: 'product/form',    name: 'nav.productAdd', component: () => import('@/views/product/ProductForm.vue'), meta: { titleKey: 'nav.productAdd' } },
        { path: 'product/form/:id',name: 'nav.productEdit', component: () => import('@/views/product/ProductForm.vue'), meta: { titleKey: 'nav.productEdit' } },
        // Order
        { path: 'order/list', name: 'nav.orderList', component: () => import('@/views/order/OrderList.vue'), meta: { titleKey: 'nav.orderList' } },
        { path: 'order/detail/:id', name: 'nav.orderDetail', component: () => import('@/views/order/OrderDetail.vue'), meta: { titleKey: 'nav.orderDetail' } },
        { path: 'order/after-sale/list', name: 'nav.afterSaleList', component: () => import('@/views/order/AfterSaleList.vue'), meta: { titleKey: 'nav.afterSaleList' } },
        { path: 'order/after-sale/detail/:id', name: 'nav.afterSaleDetail', component: () => import('@/views/order/AfterSaleDetail.vue'), meta: { titleKey: 'nav.afterSaleDetail' } },
        { path: 'review/list', name: 'nav.reviewList', component: () => import('@/views/review/ReviewList.vue'), meta: { titleKey: 'nav.reviewList' } },
        // WMS
        { path: 'wms/stock', name: 'nav.stockManage', component: () => import('@/views/wms/StockLedger.vue'), meta: { titleKey: 'nav.stockManage' } },
        { path: 'wms/warehouse', name: 'nav.warehouseManage', component: () => import('@/views/wms/WarehouseList.vue'), meta: { titleKey: 'nav.warehouseManage' } },
        { path: 'wms/docs', name: 'nav.wmsDocList', component: () => import('@/views/wms/DocList.vue'), meta: { titleKey: 'nav.wmsDocList' } },
        { path: 'wms/movements', name: 'nav.wmsMovementList', component: () => import('@/views/wms/MovementList.vue'), meta: { titleKey: 'nav.wmsMovementList' } },
        { path: 'wms/docs/new', name: 'nav.wmsDocCreate', component: () => import('@/views/wms/DocEditor.vue'), meta: { titleKey: 'nav.wmsDocEditor' } },
        { path: 'wms/docs/:id', name: 'nav.wmsDocEditor', component: () => import('@/views/wms/DocEditor.vue'), meta: { titleKey: 'nav.wmsDocEditor' } },
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
        // Points Mall
        { path: 'points-mall/products', name: 'nav.pointsProducts', component: () => import('@/views/points-mall/ProductList.vue'), meta: { titleKey: 'menu.pointsProducts' } },
        { path: 'points-mall/exchanges', name: 'nav.pointsExchanges', component: () => import('@/views/points-mall/ExchangeList.vue'), meta: { titleKey: 'menu.pointsExchanges' } },
        { path: 'points-mall/logs', name: 'nav.pointsLogs', component: () => import('@/views/points-mall/PointsLogs.vue'), meta: { titleKey: 'menu.pointsLogs' } },
        { path: 'points-mall/stats', name: 'nav.pointsStats', component: () => import('@/views/points-mall/PointsStats.vue'), meta: { titleKey: 'menu.pointsStats' } },
        { path: 'points-mall/config', name: 'nav.pointsConfig', component: () => import('@/views/points-mall/PointsConfig.vue'), meta: { titleKey: 'menu.pointsConfig' } },
        // Seckill
        { path: 'seckill/activities', name: 'nav.seckillActivities', component: () => import('@/views/seckill/ActivityList.vue'), meta: { titleKey: 'menu.seckillActivities' } },
        { path: 'seckill/products', name: 'nav.seckillProducts', component: () => import('@/views/seckill/ProductManage.vue'), meta: { titleKey: 'menu.seckillProducts' } },
        // Group Buy
        { path: 'group-buy/activities', name: 'nav.groupBuyActivities', component: () => import('@/views/group-buy/ActivityList.vue'), meta: { titleKey: 'menu.groupBuyActivities' } },
        { path: 'group-buy/products', name: 'nav.groupBuyProducts', component: () => import('@/views/group-buy/ProductManage.vue'), meta: { titleKey: 'menu.groupBuyProducts' } },
        // Bargain
        { path: 'bargain/activities', name: 'nav.bargainActivities', component: () => import('@/views/bargain/ActivityList.vue'), meta: { titleKey: 'menu.bargainActivities' } },
        { path: 'bargain/products', name: 'nav.bargainProducts', component: () => import('@/views/bargain/ProductManage.vue'), meta: { titleKey: 'menu.bargainProducts' } },
        // Distribution
        { path: 'distribution/distributors', name: 'nav.distributorManage', component: () => import('@/views/distribution/DistributorList.vue'), meta: { titleKey: 'nav.distributorManage' } },
        { path: 'distribution/orders', name: 'nav.distributionOrders', component: () => import('@/views/distribution/OrderList.vue'), meta: { titleKey: 'nav.distributionOrders' } },
        { path: 'distribution/withdrawals', name: 'nav.distributionWithdrawals', component: () => import('@/views/distribution/WithdrawalList.vue'), meta: { titleKey: 'nav.distributionWithdrawals' } },
        { path: 'distribution/config', name: 'nav.distributionConfig', component: () => import('@/views/distribution/ConfigForm.vue'), meta: { titleKey: 'nav.distributionConfig' } },
        // System
        { path: 'system/site',   name: 'nav.siteSettings', component: () => import('@/views/system/SiteSettings.vue'), meta: { titleKey: 'nav.siteSettings' } },
        { path: 'system/config', name: 'nav.configCenter', component: () => import('@/views/system/PaymentConfig.vue'), meta: { titleKey: 'nav.configCenter' } },
        { path: 'system/admins', name: 'nav.adminManage', component: () => import('@/views/system/AdminList.vue'), meta: { titleKey: 'nav.adminManage' } },
        { path: 'system/roles',  name: 'nav.roleManage', component: () => import('@/views/system/RoleList.vue'), meta: { titleKey: 'nav.roleManage' } },
        // IM
        { path: 'im/sessions',  name: 'nav.imSession', component: () => import('@/views/im/SessionList.vue'), meta: { titleKey: 'nav.imSession' } },
        { path: 'im/analytics', name: 'nav.imAnalytics', component: () => import('@/views/im/Analytics.vue'), meta: { titleKey: 'nav.imAnalytics' } },
        { path: 'im/logs',      name: 'nav.imLogs', component: () => import('@/views/im/EventLogs.vue'), meta: { titleKey: 'nav.imLogs' } },
        { path: 'im/knowledge', name: 'nav.imKnowledge', component: () => import('@/views/im/KnowledgeManage.vue'), meta: { titleKey: 'nav.imKnowledge' } },
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
