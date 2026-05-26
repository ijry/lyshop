<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const TAB_PAGES = [
  { path: 'pages/index/index', icon: 'home', activeIcon: 'home-fill', textKey: 'tabBar.home' },
  { path: 'pages/product/list', icon: 'grid', activeIcon: 'grid-fill', textKey: 'tabBar.products' },
  { path: 'pages/cart/index', icon: 'shopping-cart', activeIcon: 'shopping-cart-fill', textKey: 'tabBar.cart' },
  { path: 'pages/order/list', icon: 'list', activeIcon: 'list-fill', textKey: 'tabBar.orders' },
  { path: 'pages/user/index', icon: 'account', activeIcon: 'account-fill', textKey: 'tabBar.me' },
]

const currentTab = ref(0)
const isTabPage = ref(false)
const upRootStyle = ref({})

const getThemeVars = () => {
  if (typeof uni === 'undefined' || !uni.$u || typeof uni.$u.getThemeVars !== 'function') return {}
  return uni.$u.getThemeVars() || {}
}

const fallbackBgColor = computed(() => {
  if (typeof uni === 'undefined' || !uni.$u || !uni.$u.color) return '#f5f5f5'
  return uni.$u.color.bgColor || '#f5f5f5'
})

const buildRootStyle = () => ({
  ...getThemeVars(),
  minHeight: '100vh',
  backgroundColor: `var(--up-page-bg-color, var(--up-bg-color, ${fallbackBgColor.value}))`,
})

const syncTabState = () => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const route = currentPage?.route || ''
  const tabIndex = TAB_PAGES.findIndex(tab => tab.path === route)
  isTabPage.value = tabIndex >= 0
  if (tabIndex >= 0) {
    currentTab.value = tabIndex
  }
}

const refreshRootStyle = () => {
  upRootStyle.value = buildRootStyle()
  if (typeof uni !== 'undefined' && uni.$u && typeof uni.$u.applyNativeThemeUI === 'function') {
    uni.$u.applyNativeThemeUI()
  }
}

const refreshRootState = () => {
  refreshRootStyle()
  syncTabState()
}

function onTabChange(index) {
  if (index === currentTab.value) return
  uni.switchTab({ url: `/${TAB_PAGES[index].path}` })
}

onMounted(() => {
  refreshRootState()
  if (typeof uni !== 'undefined' && typeof uni.$on === 'function') {
    uni.$on('uThemeChange', refreshRootStyle)
  }
})

onBeforeUnmount(() => {
  if (typeof uni !== 'undefined' && typeof uni.$off === 'function') {
    uni.$off('uThemeChange', refreshRootStyle)
  }
})
</script>

<template>
  <view class="up-root-wrap" :style="upRootStyle">
    <UpRootView />
    <up-tabbar
      v-if="isTabPage"
      :value="currentTab"
      :fixed="true"
      :placeholder="true"
      :safeAreaInsetBottom="true"
      activeColor="#dc2626"
      inactiveColor="#666666"
      @change="onTabChange"
    >
      <up-tabbar-item
        v-for="(tab, index) in TAB_PAGES"
        :key="tab.path"
        :name="index"
        :text="t(tab.textKey)"
        :icon="tab.icon"
        :activeIcon="tab.activeIcon"
      />
    </up-tabbar>
  </view>
</template>
