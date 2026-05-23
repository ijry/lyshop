<template>
  <div class="flex h-screen bg-slate-50">
    <!-- Sidebar -->
    <aside class="w-64 bg-slate-900 text-slate-100 flex flex-col shrink-0">
      <div class="h-16 flex items-center px-6 gap-3 border-b border-slate-800">
        <img src="/lyshop-mark.svg" alt="LYShop" class="h-9 w-9 shrink-0" />
        <span class="text-lg font-bold text-white tracking-wide">LYShop</span>
      </div>
      <nav class="flex-1 overflow-y-auto py-4 px-3">
        <router-link
          to="/dashboard"
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 hover:text-white transition mb-2"
          active-class="!bg-red-600 !text-white"
        >
          <span>Dashboard</span>
        </router-link>
        <div v-for="group in visibleMenus" :key="group.path" class="mb-2">
          <!-- Group title -->
          <div class="px-3 py-2 text-xs font-semibold text-slate-500 uppercase tracking-wider">
            {{ group.title }}
          </div>
          <!-- Children -->
          <template v-if="group.children?.length">
            <router-link
              v-for="child in group.children"
              :key="child.path"
              :to="child.path"
              class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 hover:text-white transition"
              active-class="!bg-red-600 !text-white"
            >
              <span>{{ child.title }}</span>
            </router-link>
          </template>
          <!-- No children — link directly -->
          <router-link
            v-else
            :to="group.path"
            class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 hover:text-white transition"
            active-class="!bg-red-600 !text-white"
          >
            <span>{{ group.title }}</span>
          </router-link>
        </div>
      </nav>
    </aside>

    <!-- Main area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Topbar -->
      <header class="h-16 bg-white border-b border-slate-100 flex items-center justify-between px-6 shrink-0 shadow-sm">
        <span class="text-sm text-slate-500">{{ $route.name || '' }}</span>
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
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { getMenus } from '@/api/auth'

const auth = useAuthStore()
const menus = ref<any[]>([])
const visibleMenus = computed(() => menus.value.filter((item) => item.path !== '/dashboard'))

onMounted(async () => {
  try {
    const data = await getMenus()
    menus.value = data || []
  } catch {}
})
</script>
