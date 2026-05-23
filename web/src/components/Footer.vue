<template>
  <footer class="bg-gray-50 border-t border-gray-100 mt-16">
    <div class="max-w-7xl mx-auto px-6 py-10">
      <div class="flex flex-col md:flex-row justify-between gap-8">
        <div>
          <div class="flex items-center gap-2 mb-3">
            <div class="w-7 h-7 rounded-lg bg-gradient-to-br from-red-600 to-red-500 flex-center">
              <span class="text-white text-xs font-bold">L</span>
            </div>
            <span class="text-base font-bold text-gray-900">LYShop</span>
          </div>
          <p class="text-sm text-gray-500 max-w-xs leading-relaxed">
            开源插件化商城系统，Go + Vue3 + uni-app，支持多端部署。
          </p>
        </div>
        <div class="flex gap-16">
          <div>
            <h4 class="text-sm font-semibold text-gray-800 mb-3">帮助</h4>
            <div class="flex flex-col gap-2">
              <button @click="openChat" class="text-sm text-left text-gray-500 hover:text-red-600 transition-colors bg-transparent border-none p-0">联系客服</button>
              <a href="#" class="text-sm text-gray-500 hover:text-red-600 transition-colors">配送说明</a>
              <a href="#" class="text-sm text-gray-500 hover:text-red-600 transition-colors">退换政策</a>
            </div>
          </div>
          <div>
            <h4 class="text-sm font-semibold text-gray-800 mb-3">关于</h4>
            <div class="flex flex-col gap-2">
              <a href="https://github.com/ijry/lyshop" target="_blank"
                class="text-sm text-gray-500 hover:text-red-600 transition-colors flex items-center gap-1">
                <div class="i-carbon-logo-github text-sm" /> GitHub
              </a>
              <a href="https://ijry.github.io/lyshop/" target="_blank"
                class="text-sm text-gray-500 hover:text-red-600 transition-colors">文档</a>
            </div>
          </div>
        </div>
      </div>
      <div class="border-t border-gray-200 mt-8 pt-6 text-center">
        <p class="text-xs text-gray-400">Copyright &copy; 2026 LYShop. Released under the MIT License.</p>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { useChatStore } from '@/stores/chat'

const chat = useChatStore()

function openChat() {
  chat.open('footer')
}
</script>
