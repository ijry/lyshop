import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory('/admin'),
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
      ]
    }
  ]
})

router.beforeEach(to => {
  const token = localStorage.getItem('admin_token')
  if (to.meta.requiresAuth && !token) return '/login'
})

export default router
