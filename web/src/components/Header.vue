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
        <router-link to="/"
          class="nav-link text-sm text-gray-600 transition-colors relative py-1"
          active-class="nav-link-active font-medium">
          {{ $t('header.navHome') }}
        </router-link>

        <div class="relative" @mouseenter="scenesOpen = true" @mouseleave="scenesOpen = false">
          <button type="button"
            class="scene-trigger text-sm text-gray-600 transition-colors relative py-1 flex items-center gap-1.5"
            :class="scenesOpen ? 'text-[var(--color-primary)]' : ''">
            {{ $t('header.navScenes') }}
            <span class="i-carbon-chevron-down text-xs transition-transform" :class="scenesOpen ? 'rotate-180' : ''" />
          </button>
          <div v-if="scenesOpen"
            class="absolute top-[calc(100%+12px)] left-1/2 -translate-x-1/2 w-[980px] z-60">
            <div class="rounded-2xl border border-gray-100 shadow-2xl p-4 bg-white/95 backdrop-blur-md">
              <p class="text-xs text-gray-400 px-1">{{ $t('header.sceneHint') }}</p>
              <div class="grid grid-cols-4 gap-3 mt-3">
                <button v-for="scene in sceneCards" :key="scene.key" type="button"
                  class="scene-card rounded-xl overflow-hidden bg-white text-left transition-all"
                  @click="openScene(scene.key)">
                  <img :src="scene.image" :alt="scene.name" class="w-full h-24 object-cover" />
                  <div class="p-3">
                    <p class="text-sm font-semibold text-gray-900 leading-none">{{ scene.name }}</p>
                    <p class="text-xs text-gray-500 mt-2 leading-5 line-clamp-2">{{ scene.desc }}</p>
                  </div>
                </button>
              </div>
            </div>
          </div>
        </div>

        <router-link v-for="nav in navs" :key="nav.path" :to="nav.path"
          class="nav-link text-sm text-gray-600 transition-colors relative py-1"
          active-class="nav-link-active font-medium">
          {{ nav.name }}
        </router-link>
      </nav>

      <!-- Actions -->
      <div class="flex items-center gap-4">
        <!-- Language Switcher -->
        <select :value="$i18n.locale" @change="switchLocale(($event.target as HTMLSelectElement).value)"
          class="text-xs border border-gray-200 rounded-lg px-2 py-1.5 bg-white text-gray-600 outline-none cursor-pointer">
          <option value="zh-CN">中文</option>
          <option value="en">English</option>
        </select>

        <!-- Search -->
        <div class="hidden sm:flex items-center bg-gray-50 rounded-full px-4 py-2 w-56 hover:bg-gray-100 transition-colors cursor-pointer group"
          @click="$router.push('/products')">
          <div class="i-carbon-search text-gray-400 mr-2 text-base" />
          <span class="text-sm text-gray-400 group-hover:text-gray-500">{{ $t('header.search') }}</span>
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
          {{ $t('header.login') }}
        </router-link>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCartStore } from '@/stores/cart'
import { useAuthStore } from '@/stores/auth'
import { useSiteStore } from '@/stores/site'

const { t, locale } = useI18n()
const cart = useCartStore()
const auth = useAuthStore()
const site = useSiteStore()
const cartCount = computed(() => cart.count)
const scenesOpen = ref(false)

const navs = computed(() => [
  { name: t('header.navProducts'), path: '/products' },
  { name: t('header.navOrders'), path: '/orders' },
])

const sceneCards = computed(() => [
  { key: 'supermarket', image: '/showcase/scenes/supermarket-home.png', name: t('header.scenes.supermarket.name'), desc: t('header.scenes.supermarket.desc') },
  { key: 'cake', image: '/showcase/scenes/cake-home.png', name: t('header.scenes.cake.name'), desc: t('header.scenes.cake.desc') },
  { key: 'fresh', image: '/showcase/scenes/fresh-home.png', name: t('header.scenes.fresh.name'), desc: t('header.scenes.fresh.desc') },
  { key: 'jewelry', image: '/showcase/scenes/jewelry-home.png', name: t('header.scenes.jewelry.name'), desc: t('header.scenes.jewelry.desc') },
  { key: 'farm', image: '/showcase/scenes/farm-home.png', name: t('header.scenes.farm.name'), desc: t('header.scenes.farm.desc') },
  { key: 'mother', image: '/showcase/scenes/mother-home.png', name: t('header.scenes.mother.name'), desc: t('header.scenes.mother.desc') },
  { key: 'mall', image: '/showcase/scenes/mall-home.png', name: t('header.scenes.mall.name'), desc: t('header.scenes.mall.desc') },
])

function switchLocale(lang: string) {
  locale.value = lang
  localStorage.setItem('locale', lang)
}

function openScene(key: string) {
  const url = new URL(window.location.href)
  if (key === 'mall') {
    url.searchParams.delete('demo')
  } else {
    url.searchParams.set('demo', key)
  }
  url.hash = '#/'
  window.location.href = url.toString()
}
</script>

<style scoped>
.nav-link:hover {
  color: var(--color-primary, #dc2626);
}
.scene-trigger:hover {
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
.scene-card {
  border: 1px solid #e5e7eb;
}
.scene-card:hover {
  border-color: var(--color-primary, #dc2626);
  transform: translateY(-1px);
  box-shadow: 0 14px 28px rgba(15, 23, 42, 0.1);
}
</style>
