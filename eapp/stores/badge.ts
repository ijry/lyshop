import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useBadgeStore = defineStore('eapp-badge', () => {
  const orderBadge = ref(0)
  const messageBadge = ref(0)

  function syncTabBarBadge(index: number, count: number) {
    if (typeof uni === 'undefined') return
    if (!count) {
      uni.removeTabBarBadge({ index })
      return
    }
    uni.setTabBarBadge({ index, text: count > 99 ? '99+' : String(count) })
  }

  function syncFromDashboard(data: any) {
    orderBadge.value = Number(data?.pending_ship || 0) + Number(data?.pending_after_sale || 0)
    messageBadge.value = Number(data?.unread_message || 0)
    syncTabBarBadge(1, orderBadge.value)
    syncTabBarBadge(4, messageBadge.value)
  }

  return {
    orderBadge,
    messageBadge,
    syncFromDashboard,
  }
})
