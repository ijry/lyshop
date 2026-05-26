<script setup>
import { ref, onMounted } from 'vue'
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

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const route = currentPage?.route || ''
  const tabIndex = TAB_PAGES.findIndex(tab => tab.path === route)
  isTabPage.value = tabIndex >= 0
  if (tabIndex >= 0) {
    currentTab.value = tabIndex
  }
})

function onTabChange(index) {
  if (index === currentTab.value) return
  uni.switchTab({ url: '/' + TAB_PAGES[index].path })
}
</script>

<template>
  <UpRootView />
  <up-tabbar
    v-if="isTabPage"
    :value="currentTab"
    @change="onTabChange"
    :fixed="true"
    :placeholder="true"
    :safeAreaInsetBottom="true"
    activeColor="#dc2626"
    inactiveColor="#666666"
  >
    <up-tabbar-item
      v-for="(tab, index) in TAB_PAGES"
      :key="tab.path"
      :text="t(tab.textKey)"
      :icon="tab.icon"
      :activeIcon="tab.activeIcon"
      :name="index"
    />
  </up-tabbar>
</template>
