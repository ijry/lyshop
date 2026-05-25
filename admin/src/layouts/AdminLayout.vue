<template>
  <div class="flex h-screen bg-slate-50">
    <aside class="w-80 bg-slate-900 text-slate-100 flex flex-col shrink-0">
      <div class="h-16 flex items-center px-6 gap-3 border-b border-slate-800">
        <img src="/lyshop-mark.svg" alt="LYShop" class="h-9 w-9 shrink-0" />
        <span class="text-lg font-bold text-white tracking-wide">LYShop</span>
      </div>
      <template v-if="sortedGroupedMenus.length">
        <div class="flex-1 min-h-0 flex">
          <div
            class="w-24 border-r border-slate-800 px-2 py-3 overflow-y-auto sidebar-scroll"
            @mouseleave="handleGroupMouseLeave"
          >
            <button
              v-for="group in sortedGroupedMenus"
              :key="group.key"
              class="w-full mb-2 px-2 py-2 rounded-lg text-xs transition"
              :class="activeGroupKey === group.key
                ? 'bg-slate-700 text-white'
                : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
              @mouseenter="previewGroup(group.key)"
              @click="lockGroup(group.key)"
            >
              {{ group.title }}
            </button>
          </div>
          <nav class="flex-1 overflow-y-auto py-4 px-3 sidebar-scroll">
            <router-link
              :to="dashboardMenu.path"
              class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 hover:text-white transition mb-2"
              active-class="!bg-red-600 !text-white"
            >
              <span>{{ dashboardMenu.title }}</span>
            </router-link>
            <template v-if="activeGroupMenus.length">
              <div v-for="menu in activeGroupMenus" :key="menu.path" class="mb-2">
                <div v-if="menu.children?.length" class="px-3 py-2 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                  {{ menu.title }}
                </div>
                <template v-if="menu.children?.length">
                  <router-link
                    v-for="child in menu.children"
                    :key="child.path"
                    :to="child.path"
                    class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 hover:text-white transition"
                    active-class="!bg-red-600 !text-white"
                  >
                    <span>{{ child.title }}</span>
                  </router-link>
                </template>
                <router-link
                  v-else
                  :to="menu.path"
                  class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 hover:text-white transition"
                  active-class="!bg-red-600 !text-white"
                >
                  <span>{{ menu.title }}</span>
                </router-link>
              </div>
            </template>
            <div v-else class="px-3 py-2 text-xs text-slate-500">暂无菜单</div>
          </nav>
        </div>
      </template>
      <nav v-else class="flex-1 overflow-y-auto py-4 px-3 sidebar-scroll">
        <router-link
          :to="dashboardMenu.path"
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 hover:text-white transition mb-2"
          active-class="!bg-red-600 !text-white"
        >
          <span>{{ dashboardMenu.title }}</span>
        </router-link>
        <div v-for="group in visibleLegacyMenus" :key="group.path" class="mb-2">
          <div class="px-3 py-2 text-xs font-semibold text-slate-500 uppercase tracking-wider">
            {{ group.title }}
          </div>
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

    <div class="flex-1 flex flex-col overflow-hidden">
      <header class="h-16 bg-white border-b border-slate-100 flex items-center justify-between px-6 shrink-0 shadow-sm">
        <span class="text-sm text-slate-500">{{ $route.name || '' }}</span>
        <button @click="auth.logout()" class="text-sm text-slate-500 hover:text-slate-800 transition">
          退出登录
        </button>
      </header>
      <main class="flex-1 overflow-y-auto p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  getMenus,
  type AdminMenuGroup,
  type AdminMenuGroupedResponse,
  type AdminMenuItem,
  type AdminMenuResponse,
} from '@/api/auth'

const auth = useAuthStore()
const route = useRoute()
const groupedMenus = ref<AdminMenuGroup[]>([])
const legacyMenus = ref<AdminMenuItem[]>([])
const dashboardMenu = ref({ title: 'Dashboard', path: '/dashboard' })
const hoverGroupKey = ref('')
const lockedGroupKey = ref('')

function sortBySortAndPath<T extends { sort?: number; path?: string }>(list: T[]): T[] {
  return [...list].sort((left, right) => {
    const diff = Number(left.sort || 0) - Number(right.sort || 0)
    if (diff !== 0) return diff
    return String(left.path || '').localeCompare(String(right.path || ''))
  })
}

function sortMenuTree(list: AdminMenuItem[]): AdminMenuItem[] {
  return sortBySortAndPath(list).map((item) => ({
    ...item,
    children: Array.isArray(item.children) ? sortMenuTree(item.children) : undefined,
  }))
}

function isGroupedResponse(data: AdminMenuResponse): data is AdminMenuGroupedResponse {
  return !!data && !Array.isArray(data) && Array.isArray((data as AdminMenuGroupedResponse).groups)
}

function menuContainsPath(menu: AdminMenuItem, path: string): boolean {
  if (menu.path === path) return true
  return Array.isArray(menu.children) && menu.children.some((child) => menuContainsPath(child, path))
}

const sortedGroupedMenus = computed(() => {
  return [...groupedMenus.value].sort((left, right) => {
    const diff = Number(left.sort || 0) - Number(right.sort || 0)
    if (diff !== 0) return diff
    return String(left.key || '').localeCompare(String(right.key || ''))
  })
})

const visibleLegacyMenus = computed(() =>
  legacyMenus.value.filter((item) => item.path !== dashboardMenu.value.path),
)

const routeMatchedGroupKey = computed(() => {
  const currentPath = route.path
  for (const group of sortedGroupedMenus.value) {
    if (group.menus.some((menu) => menuContainsPath(menu, currentPath))) return group.key
  }
  return ''
})

const activeGroupKey = computed(() =>
  hoverGroupKey.value || lockedGroupKey.value || routeMatchedGroupKey.value || sortedGroupedMenus.value[0]?.key || '',
)

const activeGroupMenus = computed(() => {
  const found = sortedGroupedMenus.value.find((group) => group.key === activeGroupKey.value)
  return found?.menus || []
})

function previewGroup(groupKey: string) {
  hoverGroupKey.value = groupKey
}

function lockGroup(groupKey: string) {
  lockedGroupKey.value = groupKey
}

function handleGroupMouseLeave() {
  hoverGroupKey.value = ''
}

onMounted(async () => {
  try {
    const data = await getMenus()
    if (isGroupedResponse(data)) {
      dashboardMenu.value = data.dashboard || dashboardMenu.value
      groupedMenus.value = [...(data.groups || [])]
        .sort((left, right) => Number(left.sort || 0) - Number(right.sort || 0))
        .map((group) => ({
          ...group,
          menus: sortMenuTree(Array.isArray(group.menus) ? group.menus : []),
        }))
      legacyMenus.value = []
      return
    }
    groupedMenus.value = []
    legacyMenus.value = Array.isArray(data) ? sortMenuTree(data) : []
  } catch {
    groupedMenus.value = []
    legacyMenus.value = []
  }
})
</script>

<style scoped>
.sidebar-scroll {
  scrollbar-color: #475569 #0f172a;
}

.sidebar-scroll::-webkit-scrollbar {
  width: 8px;
}

.sidebar-scroll::-webkit-scrollbar-track {
  background: #0f172a;
}

.sidebar-scroll::-webkit-scrollbar-thumb {
  background: #475569;
  border-radius: 9999px;
}

.sidebar-scroll::-webkit-scrollbar-thumb:hover {
  background: #64748b;
}
</style>
