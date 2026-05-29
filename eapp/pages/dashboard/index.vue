<script setup lang="ts">
import { onShow, onLoad } from '@dcloudio/uni-app'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import EappShell from '@/components/layout/EappShell.vue'
import PageHeader from '@/components/biz/PageHeader.vue'
import MetricCard from '@/components/biz/MetricCard.vue'
import ActionGrid from '@/components/biz/ActionGrid.vue'
import TodoCenter from '@/components/biz/TodoCenter.vue'
import AnnouncementBar from '@/components/biz/AnnouncementBar.vue'
import AreaChart from '@/components/charts/AreaChart.vue'
import RingChart from '@/components/charts/RingChart.vue'
import BarChart from '@/components/charts/BarChart.vue'
import { useDashboard } from '@/composables/useDashboard'
import { useShopStore } from '@/stores/shop'
import { useBadgeStore } from '@/stores/badge'
import { setStorage } from '@/utils/storage'
import type { ChartCategoriesSeries, ChartPie } from '@/utils/ly-charts'

const { t } = useI18n()
const shopStore = useShopStore()
const badgeStore = useBadgeStore()
const { loading, data, load: loadDashboard } = useDashboard()

/* ---------- trend period toggle ---------- */
const trendRange = ref<'7d' | '30d'>('7d')
const trendData = computed<ChartCategoriesSeries>(() => {
  return trendRange.value === '30d' ? data.value.trend.revenue_30d : data.value.trend.revenue_7d
})
function toggleTrend() {
  trendRange.value = trendRange.value === '7d' ? '30d' : '7d'
}

/* ---------- ring chart data ---------- */
const statusPie = computed<ChartPie>(() => data.value.status_distribution)

/* ---------- hot products bar data ---------- */
const hotBarData = computed<ChartCategoriesSeries>(() => {
  const items = data.value.hot_products || []
  return {
    categories: items.map((p) => p.title.length > 6 ? p.title.slice(0, 6) + '...' : p.title),
    series: [{ name: '销量', data: items.map((p) => p.sold_qty) }],
  }
})

/* ---------- announcements ---------- */
const announcements = computed(() =>
  (data.value.announcements || []).map((a) => ({
    id: a.id,
    title: a.title,
    type: (a.type === 'urgent' ? 'urgent' : 'info') as 'info' | 'urgent',
  })),
)

/* ---------- todo items ---------- */
const todoItems = computed(() => [
  { key: 'pending_ship', title: t('dashboard.pendingShip'), value: data.value.pending_ship, tone: 'warning' as const },
  { key: 'pending_after_sale', title: t('dashboard.pendingAfterSale'), value: data.value.pending_after_sale, tone: 'danger' as const },
  { key: 'unread_message', title: t('dashboard.unreadMessage'), value: data.value.unread_message, tone: 'primary' as const },
  { key: 'stock_warning', title: t('dashboard.stockWarning'), value: data.value.stock_warning, tone: 'warning' as const },
  { key: 'today_orders', title: t('dashboard.todayOrders'), value: data.value.today_orders, tone: 'primary' as const },
  { key: 'today_sales', title: t('dashboard.todaySales'), value: `¥${Number(data.value.today_sales).toFixed(0)}`, tone: 'success' as const },
])

/* ---------- action grid ---------- */
const actions = [
  { key: 'order', label: '订单管理', icon: '📦' },
  { key: 'product', label: '商品管理', icon: '🏷' },
  { key: 'after_sale', label: '售后处理', icon: '🔧' },
  { key: 'review', label: '评价管理', icon: '⭐' },
  { key: 'decor', label: '店铺装修', icon: '🎨' },
  { key: 'marketing', label: '营销中心', icon: '📣', soon: true },
  { key: 'data', label: '数据报表', icon: '📊', soon: true },
  { key: 'wms', label: '仓储管理', icon: '🏭', soon: true },
]

/* ---------- lifecycle ---------- */
async function loadAll() {
  await Promise.all([
    shopStore.load().catch(() => {}),
    loadDashboard(),
  ])
  badgeStore.syncFromDashboard(data.value)
}

function goOrder(status = '') {
  setStorage('eapp_order_status_filter', status)
  uni.switchTab({ url: '/pages/order/list' })
}

function onTodoClick(key: string) {
  if (key === 'pending_ship') goOrder('2')
  else if (key === 'pending_after_sale') uni.navigateTo({ url: '/pages/order/after-sale-list' })
  else if (key === 'unread_message') uni.switchTab({ url: '/pages/me/index' })
  else if (key === 'stock_warning') uni.switchTab({ url: '/pages/product/list' })
  else if (key === 'today_orders') goOrder('')
}

function onActionClick(key: string) {
  if (key === 'order') uni.switchTab({ url: '/pages/order/list' })
  else if (key === 'product') uni.switchTab({ url: '/pages/product/list' })
  else if (key === 'after_sale') uni.navigateTo({ url: '/pages/order/after-sale-list' })
  else if (key === 'review') uni.navigateTo({ url: '/pages/review/list' })
  else if (key === 'decor') uni.navigateTo({ url: '/pages/decor/index' })
}

function onScan() {
  uni.scanCode({ success: () => {} })
}

function onMessage() {
  uni.switchTab({ url: '/pages/me/index' })
}

onLoad(loadAll)
onShow(loadAll)
</script>

<template>
  <EappShell :padded="false" header-sticky>
    <template #header>
      <PageHeader :title="shopStore.name || t('dashboard.title')">
        <template #right>
          <view class="header-icons">
            <view class="header-icon" @click="onScan">
              <text class="header-icon__text">📷</text>
            </view>
            <view class="header-icon" @click="onMessage">
              <text class="header-icon__text">💬</text>
              <view v-if="data.unread_message" class="header-icon__dot" />
            </view>
          </view>
        </template>
      </PageHeader>
    </template>

    <view class="dashboard-body">
      <!-- Metric Cards -->
      <view class="metric-row">
        <MetricCard
          :title="t('dashboard.todaySales')"
          :value="`¥${Number(data.today_sales).toFixed(2)}`"
          :compare="data.compare.revenue_mom"
          color="#2563eb"
        />
        <MetricCard
          :title="t('dashboard.todayOrders')"
          :value="String(data.today_orders)"
          :compare="data.compare.order_mom"
          color="#16a34a"
        />
        <MetricCard
          title="客单价"
          :value="`¥${Number(data.today_avg_price).toFixed(0)}`"
          color="#f59e0b"
        />
      </view>

      <!-- Revenue Trend -->
      <AreaChart :title="'营收趋势'" :loading="loading" :data="trendData">
        <template #extra>
          <view class="trend-toggle" @click="toggleTrend">
            <text class="trend-toggle__text">{{ trendRange === '7d' ? '近7天' : '近30天' }}</text>
          </view>
        </template>
      </AreaChart>

      <!-- Order Status Ring -->
      <RingChart title="订单状态分布" :loading="loading" :data="statusPie" />

      <!-- Announcements -->
      <view v-if="announcements.length" class="section">
        <view class="section__title">公告通知</view>
        <AnnouncementBar :items="announcements" />
      </view>

      <!-- Todo Center -->
      <view class="section">
        <view class="section__title">待办中心</view>
        <TodoCenter :items="todoItems" @click="onTodoClick" />
      </view>

      <!-- Action Grid -->
      <view class="section">
        <view class="section__title">快捷入口</view>
        <ActionGrid :items="actions" @click="onActionClick" />
      </view>

      <!-- Hot Products Bar -->
      <BarChart
        title="销量 Top5"
        :loading="loading"
        :data="hotBarData"
        horizontal
      />

      <!-- Stock Warning List -->
      <view v-if="data.stock_warning_list.length" class="section">
        <view class="section__title">库存预警</view>
        <view class="stock-list">
          <view
            v-for="item in data.stock_warning_list"
            :key="item.product_id"
            class="stock-item"
          >
            <view class="stock-item__info">
              <text class="stock-item__title">{{ item.title }}</text>
              <text class="stock-item__sub">阈值: {{ item.threshold }}</text>
            </view>
            <text class="stock-item__qty" :class="{ 'stock-item__qty--danger': item.stock < item.threshold }">
              {{ item.stock }}
            </text>
          </view>
        </view>
      </view>

      <view v-if="loading" class="loading-hint">{{ t('common.loading') }}</view>
    </view>
  </EappShell>
</template>

<style scoped>
.dashboard-body {
  padding: 24rpx;
}

.header-icons {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.header-icon {
  position: relative;
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-icon__text {
  font-size: 36rpx;
}

.header-icon__dot {
  position: absolute;
  top: 4rpx;
  right: 4rpx;
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  background: var(--eapp-danger);
}

.metric-row {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 16rpx;
  margin-bottom: 24rpx;
}

.trend-toggle {
  padding: 6rpx 16rpx;
  border-radius: 999rpx;
  background: var(--eapp-bg);
}

.trend-toggle__text {
  font-size: 22rpx;
  color: var(--eapp-primary);
}

.section {
  margin-bottom: 24rpx;
}

.section__title {
  font-size: 28rpx;
  font-weight: 600;
  color: var(--eapp-text);
  margin-bottom: 16rpx;
}

.stock-list {
  background: var(--eapp-card);
  border-radius: 22rpx;
  overflow: hidden;
}

.stock-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 24rpx;
  border-bottom: 1rpx solid var(--eapp-border);
}

.stock-item:last-child {
  border-bottom: none;
}

.stock-item__info {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  flex: 1;
  min-width: 0;
}

.stock-item__title {
  font-size: 26rpx;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stock-item__sub {
  font-size: 22rpx;
  color: var(--eapp-text-muted);
}

.stock-item__qty {
  font-size: 32rpx;
  font-weight: 700;
  color: var(--eapp-text);
  margin-left: 16rpx;
}

.stock-item__qty--danger {
  color: var(--eapp-danger);
}

.loading-hint {
  text-align: center;
  padding: 24rpx 0;
  font-size: 24rpx;
  color: var(--eapp-text-muted);
}
</style>
