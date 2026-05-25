<template>
  <footer class="border-t border-gray-100 mt-16"
    :style="{ background: 'var(--color-bg-footer, #f9fafb)' }">
    <div class="max-w-7xl mx-auto px-6 py-10">
      <div class="flex flex-col md:flex-row justify-between gap-8">
        <div>
          <div class="flex items-center gap-2 mb-3">
            <img v-if="site.settings.site_logo" :src="site.settings.site_logo" class="w-7 h-7 rounded-lg object-cover" />
            <div v-else class="w-7 h-7 rounded-lg flex-center text-white text-xs font-bold"
              :style="{ background: 'var(--color-primary, #dc2626)' }">
              {{ (site.settings.site_name || 'L').charAt(0) }}
            </div>
            <span class="text-base font-bold text-gray-900">{{ site.settings.site_name || 'LYShop' }}</span>
          </div>
          <p class="text-sm text-gray-500 max-w-xs leading-relaxed">
            {{ site.settings.seo_description || '开源插件化商城系统，Go + Vue3 + uni-app，支持多端部署。' }}
          </p>
        </div>
        <div class="flex gap-16">
          <div>
            <h4 class="text-sm font-semibold text-gray-800 mb-3">帮助</h4>
            <div class="flex flex-col gap-2">
              <button @click="openChat" class="footer-link text-sm text-left text-gray-500 transition-colors bg-transparent border-none p-0">联系客服</button>
              <a href="#" class="footer-link text-sm text-gray-500 transition-colors">配送说明</a>
              <a href="#" class="footer-link text-sm text-gray-500 transition-colors">退换政策</a>
            </div>
          </div>
          <div>
            <h4 class="text-sm font-semibold text-gray-800 mb-3">关于</h4>
            <div class="flex flex-col gap-2">
              <a href="https://github.com/ijry/lyshop" target="_blank"
                class="footer-link text-sm text-gray-500 transition-colors flex items-center gap-1">
                <div class="i-carbon-logo-github text-sm" /> GitHub
              </a>
              <a href="https://ijry.github.io/lyshop/" target="_blank"
                class="footer-link text-sm text-gray-500 transition-colors">文档</a>
            </div>
          </div>
        </div>
      </div>
      <div class="border-t border-gray-200 mt-8 pt-6 text-center">
        <p class="text-xs text-gray-400">
          Copyright &copy; {{ new Date().getFullYear() }} {{ site.settings.site_name || 'LYShop' }}. Released under the MIT License.
        </p>
        <p v-if="site.settings.icp" class="text-xs text-gray-400 mt-1">
          <a href="https://beian.miit.gov.cn/" target="_blank" class="hover:text-gray-500 transition-colors">
            {{ site.settings.icp }}
          </a>
        </p>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { useChatStore } from '@/stores/chat'
import { useSiteStore } from '@/stores/site'

const chat = useChatStore()
const site = useSiteStore()

function openChat() {
  chat.open('footer')
}
</script>

<style scoped>
.footer-link:hover {
  color: var(--color-primary, #dc2626);
}
</style>
