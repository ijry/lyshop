<template>
  <header class="sticky top-0 z-50 backdrop-blur-lg border-b border-gray-100"
    :style="{ background: 'var(--color-bg-header, rgba(255,255,255,0.8))' }">
    <div class="max-w-7xl mx-auto px-6 h-16 flex-between">
      <!-- Logo -->
      <router-link to="/" class="flex items-center gap-2.5 group">
        <img v-if="site.settings.site_logo" :src="site.settings.site_logo" class="w-8 h-8 rounded-lg object-cover" />
        <div v-else class="w-8 h-8 rounded-lg flex-center text-white text-sm font-bold"
          :style="{ background: 'var(--color-primary, #dc2626)' }">
          {{ (site.settings.site_name || 'L').charAt(0) }}
        </div>
        <span class="text-lg font-bold text-gray-900 group-hover:text-[var(--color-primary)] transition-colors">
          {{ site.settings.site_name || 'LYShop' }}
        </span>
      </router-link>

      <!-- Nav -->
      <nav class="hidden md:flex items-center gap-8">
        <router-link v-for="nav in navs" :key="nav.path" :to="nav.path"
          class="nav-link text-sm text-gray-600 transition-colors relative py-1"
          active-class="nav-link-active font-medium">
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
            class="absolute -top-0.5 -right-0.5 min-w-4.5 h-4.5 text-white text-xs rounded-full flex-center px-1 font-medium"
            :style="{ background: 'var(--color-primary, #dc2626)' }">
            {{ cartCount }}
          </span>
        </router-link>

        <!-- User -->
        <template v-if="auth.isLoggedIn">
          <div class="w-8 h-8 rounded-full flex-center cursor-pointer transition-colors"
            :style="{ background: 'color-mix(in srgb, var(--color-primary) 15%, white)', color: 'var(--color-primary)' }"
            @click="$router.push('/user')">
            <div class="i-carbon-user" />
          </div>
        </template>
        <router-link v-else to="/login"
          class="px-4 py-2 text-xs text-white rounded-lg font-medium transition-colors cursor-pointer"
          :style="{ background: 'var(--color-primary, #dc2626)' }">
          登录
        </router-link>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useCartStore } from '@/stores/cart'
import { useAuthStore } from '@/stores/auth'
import { useSiteStore } from '@/stores/site'

const cart = useCartStore()
const auth = useAuthStore()
const site = useSiteStore()
const cartCount = computed(() => cart.count)

const navs = [
  { name: '首页', path: '/' },
  { name: '全部商品', path: '/products' },
  { name: '我的订单', path: '/orders' },
]
</script>

<style scoped>
.nav-link:hover {
  color: var(--color-primary, #dc2626);
}
.nav-link-active {
  color: var(--color-primary, #dc2626) !important;
}
.nav-link-active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--color-primary, #dc2626);
  border-radius: 9999px;
}
</style>
