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
        { path: 'dashboard', component: () => import('@/views/Dashboard.vue') },
        // Product
        { path: 'product/list',    name: '商品列表', component: () => import('@/views/product/ProductList.vue') },
        { path: 'product/form',    name: '新增商品', component: () => import('@/views/product/ProductForm.vue') },
        { path: 'product/form/:id',name: '编辑商品', component: () => import('@/views/product/ProductForm.vue') },
        // Order
        { path: 'order/list', name: '订单列表', component: () => import('@/views/order/OrderList.vue') },
        { path: 'order/detail/:id', name: '订单详情', component: () => import('@/views/order/OrderDetail.vue') },
        { path: 'order/after-sale/list', name: '售后列表', component: () => import('@/views/order/AfterSaleList.vue') },
        { path: 'order/after-sale/detail/:id', name: '售后详情', component: () => import('@/views/order/AfterSaleDetail.vue') },
        { path: 'review/list', name: '评价列表', component: () => import('@/views/review/ReviewList.vue') },
        // WMS
        { path: 'wms/stock', name: '库存管理', component: () => import('@/views/wms/StockList.vue') },
        // Marketing
        { path: 'marketing/coupon',   name: '优惠券管理', component: () => import('@/views/marketing/CouponList.vue') },
        // System
        { path: 'system/config', name: '配置中心', component: () => import('@/views/system/PaymentConfig.vue') },
        { path: 'system/admins', name: '管理员管理', component: () => import('@/views/system/AdminList.vue') },
        { path: 'system/roles',  name: '角色管理', component: () => import('@/views/system/RoleList.vue') },
        // IM
        { path: 'im/sessions',  name: '客服会话', component: () => import('@/views/im/SessionList.vue') },
        // AI
        { path: 'ai/tasks',     name: 'AI生图',   component: () => import('@/views/ai/ImageGen.vue') },
        { path: 'ai/models',    name: 'AI模型配置', component: () => import('@/views/ai/ImageGen.vue') },
        // Decor
        { path: 'decor/index',  name: '首页装修', component: () => import('@/views/decor/DecorEditor.vue') },
        // Checkin
        { path: 'checkin/rules', name: '签到规则', component: () => import('@/views/checkin/CheckinRules.vue') },
        { path: 'checkin/logs',  name: '签到记录', component: () => import('@/views/checkin/CheckinLogs.vue') },
        // Message
        { path: 'message/list', name: '消息列表', component: () => import('@/views/message/MessageList.vue') },
        { path: 'message/send', name: '发送消息', component: () => import('@/views/message/MessageSend.vue') },
      ]
    }
  ]
})

router.beforeEach(to => {
  const token = localStorage.getItem('admin_token')
  if (to.meta.requiresAuth && !token) return '/login'
})

export default router
