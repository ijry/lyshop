import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', name: 'home', component: () => import('@/views/Home.vue') },
    { path: '/products', name: 'products', component: () => import('@/views/ProductList.vue') },
    { path: '/products/seckill', name: 'products-seckill', component: () => import('@/views/SeckillProductList.vue') },
    { path: '/products/group-buy', name: 'products-group-buy', component: () => import('@/views/GroupBuyProductList.vue') },
    { path: '/products/bargain', name: 'products-bargain', component: () => import('@/views/BargainProductList.vue') },
    { path: '/product/:id', name: 'product', component: () => import('@/views/ProductDetail.vue') },
    { path: '/cart', name: 'cart', component: () => import('@/views/Cart.vue') },
    { path: '/orders', name: 'orders', component: () => import('@/views/OrderList.vue') },
    { path: '/orders/:id', name: 'order-detail', component: () => import('@/views/OrderDetail.vue') },
    { path: '/orders/:id/after-sale/apply', name: 'after-sale-apply', component: () => import('@/views/AfterSaleApply.vue') },
    { path: '/after-sales/:id', name: 'after-sale-detail', component: () => import('@/views/AfterSaleDetail.vue') },
    { path: '/orders/:id/review', name: 'order-review', component: () => import('@/views/OrderReview.vue') },
    { path: '/user', name: 'user', component: () => import('@/views/UserCenter.vue') },
    { path: '/chat', name: 'chat', component: () => import('@/views/Chat.vue') },
    { path: '/login', name: 'login', component: () => import('@/views/Login.vue') },
  ],
  scrollBehavior: () => ({ top: 0 }),
})

export default router
