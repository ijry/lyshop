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
        // WMS
        { path: 'wms/stock', name: '库存管理', component: () => import('@/views/wms/StockList.vue') },
        // Marketing
        { path: 'marketing/coupon',   name: '优惠券管理', component: () => import('@/views/marketing/CouponList.vue') },
        // System
        { path: 'system/config', name: '支付短信配置', component: () => import('@/views/system/PaymentConfig.vue') },
        // IM
        { path: 'im/sessions',  name: '客服会话', component: () => import('@/views/im/SessionList.vue') },
        // AI
        { path: 'ai/tasks',     name: 'AI生图',   component: () => import('@/views/ai/ImageGen.vue') },
        { path: 'ai/models',    name: 'AI模型配置', component: () => import('@/views/ai/ImageGen.vue') },
        // Decor
        { path: 'decor/index',  name: '首页装修', component: () => import('@/views/decor/DecorEditor.vue') },
      ]
    }
  ]
})

router.beforeEach(to => {
  const token = localStorage.getItem('admin_token')
  if (to.meta.requiresAuth && !token) return '/login'
})

export default router
