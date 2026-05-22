<template>
  <header class="sticky top-0 z-50 bg-white/80 backdrop-blur-lg border-b border-gray-100">
    <div class="max-w-7xl mx-auto px-6 h-16 flex-between">
      <!-- Logo -->
      <router-link to="/" class="flex items-center gap-2.5 group">
        <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-700 to-blue-500 flex-center">
          <span class="text-white text-sm font-bold">L</span>
        </div>
        <span class="text-lg font-bold text-gray-900 group-hover:text-blue-700 transition-colors">LYShop</span>
      </router-link>

      <!-- Nav -->
      <nav class="hidden md:flex items-center gap-8">
        <router-link v-for="nav in navs" :key="nav.path" :to="nav.path"
          class="text-sm text-gray-600 hover:text-blue-700 transition-colors relative py-1"
          active-class="!text-blue-700 font-medium after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-blue-700 after:rounded-full">
          {{ nav.name }}
        </router-link>
      </nav>

      <!-- Actions -->
      <div class="flex items-center gap-4">
        <!-- Search -->
        <div class="hidden sm:flex items-center bg-gray-50 rounded-full px-4 py-2 w-56 hover:bg-gray-100 transition-colors cursor-pointer group"
          @click="$router.push('/products')">
          <div class="i-carbon-search text-gray-400 mr-2 text-base" />
          <span class="text-sm text-gray-400 group-hover:text-gray-500">搜索商品</span>
        </div>

        <!-- Cart -->
        <router-link to="/cart"
          class="relative p-2 rounded-lg hover:bg-gray-50 transition-colors">
          <div class="i-carbon-shopping-cart text-xl text-gray-600" />
          <span v-if="cartCount > 0"
            class="absolute -top-0.5 -right-0.5 min-w-4.5 h-4.5 bg-red-500 text-white text-xs rounded-full flex-center px-1 font-medium">
            {{ cartCount }}
          </span>
        </router-link>

        <!-- User -->
        <template v-if="auth.isLoggedIn">
          <div class="w-8 h-8 rounded-full bg-blue-100 flex-center cursor-pointer hover:bg-blue-200 transition-colors"
            @click="$router.push('/orders')">
            <div class="i-carbon-user text-blue-700" />
          </div>
        </template>
        <router-link v-else to="/login" class="btn-primary text-xs !px-4 !py-2">登录</router-link>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useCartStore } from '@/stores/cart'
import { useAuthStore } from '@/stores/auth'

const cart = useCartStore()
const auth = useAuthStore()
const cartCount = computed(() => cart.count)

const navs = [
  { name: '首页', path: '/' },
  { name: '全部商品', path: '/products' },
  { name: '我的订单', path: '/orders' },
]
</script>
