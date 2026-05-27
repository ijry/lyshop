<template>
  <div class="flex h-screen bg-slate-50">
    <aside class="w-64 bg-slate-900 text-slate-100 flex flex-col shrink-0">
      <div class="h-16 flex shrink-0">
        <div class="w-24 shrink-0 flex items-center justify-center border-r border-b border-slate-800">
          <img src="/lyshop-mark.svg" alt="LYShop" class="h-9 w-9" />
        </div>
        <div class="flex-1 flex items-center px-3 bg-slate-100 border-b border-slate-200">
          <span class="text-lg font-bold text-slate-800 tracking-wide">LYShop</span>
        </div>
      </div>
      <template v-if="sortedGroupedMenus.length">
        <div class="flex-1 min-h-0 flex" @mouseleave="clearPreviewGroup">
          <div class="w-24 border-r border-slate-800 px-2 py-3 overflow-y-auto sidebar-scroll">
            <router-link
              :to="dashboardMenu.path"
              class="block w-full mb-2 px-2 py-2 rounded-lg text-xs text-center transition"
              :class="routeMatchedGroupKey === homeTabKey && !hoverGroupKey
                ? 'bg-slate-700 text-white'
                : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
              @mouseenter="previewGroup(homeTabKey)"
            >
              {{ mt(dashboardMenu) }}
            </router-link>
            <button
              v-for="group in sortedGroupedMenus"
              :key="group.key"
              class="w-full mb-2 px-2 py-2 rounded-lg text-xs transition"
              :class="activeGroupKey === group.key
                ? 'bg-slate-700 text-white'
                : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
              @mouseenter="previewGroup(group.key)"
            >
              {{ mt(group) }}
            </button>
          </div>
          <nav class="flex-1 overflow-y-auto py-4 px-3 bg-slate-100 sidebar-scroll-light">
            <template v-if="activeGroupMenus.length">
              <div v-for="menu in activeGroupMenus" :key="menu.path" class="mb-2">
                <div v-if="menu.children?.length" class="px-3 py-2 text-xs font-semibold text-slate-400 uppercase tracking-wider">
                  {{ mt(menu) }}
                </div>
                <template v-if="menu.children?.length">
                  <router-link
                    v-for="child in menu.children"
                    :key="child.path"
                    :to="child.path"
                    class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-600 hover:bg-slate-200 hover:text-slate-900 transition"
                    active-class="!bg-red-600 !text-white"
                  >
                    <span>{{ mt(child) }}</span>
                  </router-link>
                </template>
                <router-link
                  v-else
                  :to="menu.path"
                  class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-600 hover:bg-slate-200 hover:text-slate-900 transition"
                  active-class="!bg-red-600 !text-white"
                >
                  <span>{{ mt(menu) }}</span>
                </router-link>
              </div>
            </template>
            <div v-else class="px-3 py-2 text-xs text-slate-400">{{ $t('layout.noMenu') }}</div>
          </nav>
        </div>
      </template>
      <nav v-else class="flex-1 overflow-y-auto py-4 px-3 sidebar-scroll">
        <router-link
          :to="dashboardMenu.path"
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 hover:text-white transition mb-2"
          active-class="!bg-red-600 !text-white"
        >
          <span>{{ mt(dashboardMenu) }}</span>
        </router-link>
        <div v-for="group in visibleLegacyMenus" :key="group.path" class="mb-2">
          <div class="px-3 py-2 text-xs font-semibold text-slate-500 uppercase tracking-wider">
            {{ mt(group) }}
          </div>
          <template v-if="group.children?.length">
            <router-link
              v-for="child in group.children"
              :key="child.path"
              :to="child.path"
              class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 hover:text-white transition"
              active-class="!bg-red-600 !text-white"
            >
              <span>{{ mt(child) }}</span>
            </router-link>
          </template>
          <router-link
            v-else
            :to="group.path"
            class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-300 hover:bg-slate-800 hover:text-white transition"
            active-class="!bg-red-600 !text-white"
          >
            <span>{{ mt(group) }}</span>
          </router-link>
        </div>
      </nav>
    </aside>

    <div class="flex-1 flex flex-col overflow-hidden">
      <header class="h-16 bg-white border-b border-slate-100 flex items-center justify-between px-6 shrink-0 shadow-sm">
        <span class="text-sm text-slate-500">{{ $route.meta.titleKey ? $t($route.meta.titleKey as string) : ($route.name || '') }}</span>
        <div class="flex items-center gap-4">
          <select
            :value="$i18n.locale"
            class="text-sm border border-slate-200 rounded px-2 py-1 text-slate-600"
            @change="switchLocale(($event.target as HTMLSelectElement).value)"
          >
            <option value="zh-CN">中文</option>
            <option value="en">English</option>
          </select>
          <div ref="accountMenuRef" class="relative">
            <button
              class="flex items-center gap-3 rounded-full border border-slate-200 bg-white px-2 py-1 pr-3 hover:border-slate-300"
              @click="accountMenuOpen = !accountMenuOpen"
            >
              <img :src="currentAccount.avatar" :alt="currentAccount.username" class="h-8 w-8 rounded-full bg-slate-100" />
              <div class="flex flex-col items-start">
                <span class="text-sm font-medium text-slate-700">{{ currentAccount.username }}</span>
                <a
                  :href="currentAccount.github"
                  target="_blank"
                  rel="noreferrer"
                  class="flex items-center gap-1 text-xs text-slate-400 hover:text-slate-600"
                  @click.stop
                >
                  <Github class="h-3.5 w-3.5" />
                  <span>GitHub</span>
                </a>
              </div>
            </button>
            <div v-if="accountMenuOpen" class="absolute right-0 top-12 z-30 w-64 rounded-xl border border-slate-200 bg-white p-2 shadow-xl">
              <div class="px-3 py-2 text-xs font-semibold text-slate-400">{{ $t('layout.accountSwitch') }}</div>
              <button
                v-for="account in auth.accounts"
                :key="account.username"
                class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left hover:bg-slate-50"
                @click="handleSwitchAccount(account.username)"
              >
                <img :src="account.avatar" :alt="account.username" class="h-8 w-8 rounded-full bg-slate-100" />
                <div class="min-w-0 flex-1">
                  <div class="truncate text-sm text-slate-700">{{ account.username }}</div>
                  <div class="truncate text-xs text-slate-400">{{ account.github }}</div>
                </div>
                <span v-if="account.username === auth.currentUsername" class="text-xs text-emerald-600">{{ $t('layout.currentAccount') }}</span>
              </button>
              <button @click="auth.logout()" class="mt-2 flex w-full items-center justify-center rounded-lg bg-slate-100 px-3 py-2 text-sm text-slate-600 hover:bg-slate-200">
                {{ $t('layout.logout') }}
              </button>
            </div>
          </div>
        </div>
      </header>
      <main class="flex-1 overflow-y-auto p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Github } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import {
  getMenus,
  type AdminMenuGroup,
  type AdminMenuGroupedResponse,
  type AdminMenuItem,
  type AdminMenuResponse,
} from '@/api/auth'

const { t, locale } = useI18n()
const auth = useAuthStore()
const route = useRoute()
const accountMenuOpen = ref(false)
const accountMenuRef = ref<HTMLElement | null>(null)

/** Translate menu title: use title_key if available, fall back to raw title */
function mt(item: { title: string; title_key?: string }) {
  return item.title_key ? t(item.title_key) : item.title
}

const currentAccount = computed(() => {
  const found = auth.accounts.find((item) => item.username === auth.currentUsername)
  return found || {
    username: auth.currentUsername || 'admin',
    avatar: `https://api.dicebear.com/9.x/initials/svg?seed=${encodeURIComponent(auth.currentUsername || 'admin')}`,
    github: `https://github.com/${encodeURIComponent(auth.currentUsername || 'admin')}`,
  }
})

function switchLocale(lang: string) {
  locale.value = lang
  localStorage.setItem('locale', lang)
}

function handleSwitchAccount(username: string) {
  accountMenuOpen.value = false
  if (username === auth.currentUsername) return
  auth.switchAccount(username)
}

function handleWindowClick(event: MouseEvent) {
  if (!accountMenuOpen.value) return
  const target = event.target as Node | null
  if (target && accountMenuRef.value?.contains(target)) return
  accountMenuOpen.value = false
}
const groupedMenus = ref<AdminMenuGroup[]>([])
const legacyMenus = ref<AdminMenuItem[]>([])
const dashboardMenu = ref({ title: '', path: '/dashboard' })
const homeTabKey = '__home__'
const hoverGroupKey = ref('')

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
  if (currentPath === dashboardMenu.value.path) return homeTabKey
  for (const group of sortedGroupedMenus.value) {
    if (group.menus.some((menu) => menuContainsPath(menu, currentPath))) return group.key
  }
  return ''
})

const activeGroupKey = computed(() => {
  const key = hoverGroupKey.value || routeMatchedGroupKey.value || sortedGroupedMenus.value[0]?.key || ''
  return key === homeTabKey ? '' : key
})

const activeGroupMenus = computed(() => {
  if (!activeGroupKey.value) return []
  const found = sortedGroupedMenus.value.find((group) => group.key === activeGroupKey.value)
  return found?.menus || []
})

function previewGroup(groupKey: string) {
  hoverGroupKey.value = groupKey
}

function clearPreviewGroup() {
  hoverGroupKey.value = ''
}

onMounted(async () => {
  window.addEventListener('click', handleWindowClick)
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

onBeforeUnmount(() => {
  window.removeEventListener('click', handleWindowClick)
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

.sidebar-scroll-light {
  scrollbar-color: #cbd5e1 #f1f5f9;
}

.sidebar-scroll-light::-webkit-scrollbar {
  width: 8px;
}

.sidebar-scroll-light::-webkit-scrollbar-track {
  background: #f1f5f9;
}

.sidebar-scroll-light::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 9999px;
}

.sidebar-scroll-light::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>
