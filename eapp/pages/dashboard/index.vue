<script setup lang="ts">
import { onShow, onLoad } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import EappShell from '@/components/layout/EappShell.vue'
import StatCard from '@/components/common/StatCard.vue'
import TodoCard from '@/components/common/TodoCard.vue'
import { getDashboard } from '@/api/dashboard'
import { useBadgeStore } from '@/stores/badge'
import { setStorage } from '@/utils/storage'

const { t } = useI18n()
const badgeStore = useBadgeStore()

const loading = ref(false)
const stats = reactive({ today_orders: 0, today_sales: 0, pending_ship: 0, pending_after_sale: 0, unread_message: 0, stock_warning: 0 })

async function loadData() {
  loading.value = true
  try {
    const data: any = await getDashboard()
    stats.today_orders = Number(data?.today_orders || 0)
    stats.today_sales = Number(data?.today_sales || 0)
    stats.pending_ship = Number(data?.pending_ship || 0)
    stats.pending_after_sale = Number(data?.pending_after_sale || data?.pending_refunds || 0)
    stats.unread_message = Number(data?.unread_message || 0)
    stats.stock_warning = Number(data?.stock_warning || 0)
    badgeStore.syncFromDashboard(stats)
  } catch {
  } finally {
    loading.value = false
  }
}

function goOrder(status = '') {
  setStorage('eapp_order_status_filter', status)
  uni.switchTab({ url: '/pages/order/list' })
}

onLoad(loadData)
onShow(loadData)
</script>

<template>
  <EappShell>
    <template #header>
      <view class="header">
        <view class="title">{{ t('dashboard.title') }}</view>
        <up-button size="mini" type="primary" plain @click="loadData">{{ t('common.refresh') }}</up-button>
      </view>
    </template>

    <view class="grid">
      <StatCard :title="t('dashboard.todayOrders')" :value="String(stats.today_orders)" />
      <StatCard :title="t('dashboard.todaySales')" :value="`¥${Number(stats.today_sales).toFixed(2)}`" color="#2563eb" />
    </view>

    <view class="section-title">待办</view>
    <view class="todo-list">
      <TodoCard :title="t('dashboard.pendingShip')" :value="stats.pending_ship" @click="goOrder('2')" />
      <TodoCard :title="t('dashboard.pendingAfterSale')" :value="stats.pending_after_sale" @click="uni.navigateTo({ url: '/pages/order/after-sale-list' })" />
      <TodoCard :title="t('dashboard.unreadMessage')" :value="stats.unread_message" @click="uni.switchTab({ url: '/pages/me/index' })" />
      <TodoCard :title="t('dashboard.stockWarning')" :value="stats.stock_warning" @click="uni.switchTab({ url: '/pages/product/list' })" />
    </view>

    <view v-if="loading" class="loading">{{ t('common.loading') }}</view>
  </EappShell>
</template>

<style scoped>
.header {
  padding: 24rpx;
  padding-top: calc(24rpx + env(safe-area-inset-top));
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title {
  font-size: 36rpx;
  font-weight: 700;
}

.grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20rpx;
}

.section-title {
  margin-top: 28rpx;
  margin-bottom: 16rpx;
  font-size: 30rpx;
  font-weight: 700;
}

.todo-list {
  display: grid;
  gap: 16rpx;
}

.loading {
  margin-top: 20rpx;
  color: var(--eapp-text-muted);
  font-size: 24rpx;
}
</style>
