<template>
  <div class="flex h-screen bg-slate-50">
    <!-- Sidebar -->
    <aside class="w-64 bg-slate-900 text-slate-100 flex flex-col shrink-0">
      <div class="h-16 flex items-center px-6 border-b border-slate-800">
        <span class="text-lg font-bold text-white">lyshop</span>
      </div>
      <nav class="flex-1 overflow-y-auto py-4 space-y-1 px-3">
        <router-link
          v-for="item in menus"
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 hover:text-white transition"
          active-class="bg-blue-700 text-white"
        >
          <span>{{ item.title }}</span>
        </router-link>
      </nav>
    </aside>

    <!-- Main area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Topbar -->
      <header class="h-16 bg-white border-b border-slate-100 flex items-center justify-between px-6 shrink-0 shadow-sm">
        <span class="text-sm text-slate-500">{{ $route.name }}</span>
        <button @click="auth.logout()" class="text-sm text-slate-500 hover:text-slate-800 transition">
          退出登录
        </button>
      </header>
      <!-- Content -->
      <main class="flex-1 overflow-y-auto p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { getMenus } from '@/api/auth'

const auth = useAuthStore()
const menus = ref<any[]>([])

onMounted(async () => {
  try {
    menus.value = await getMenus()
  } catch {}
})
</script>
