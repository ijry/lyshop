import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', name: 'home', component: () => import('@/views/Home.vue') },
    { path: '/products', name: 'products', component: () => import('@/views/ProductList.vue') },
    { path: '/product/:id', name: 'product', component: () => import('@/views/ProductDetail.vue') },
    { path: '/cart', name: 'cart', component: () => import('@/views/Cart.vue') },
    { path: '/orders', name: 'orders', component: () => import('@/views/OrderList.vue') },
    { path: '/user', name: 'user', component: () => import('@/views/UserCenter.vue') },
    { path: '/chat', name: 'chat', component: () => import('@/views/Chat.vue') },
    { path: '/login', name: 'login', component: () => import('@/views/Login.vue') },
  ],
  scrollBehavior: () => ({ top: 0 }),
})

export default router
