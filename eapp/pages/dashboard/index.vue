<script setup lang="ts">
import { onShow, onLoad } from '@dcloudio/uni-app'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import EappShell from '@/components/layout/EappShell.vue'
import PageHeader from '@/components/biz/PageHeader.vue'
import SectionHeader from '@/components/biz/SectionHeader.vue'
import MetricCard from '@/components/biz/MetricCard.vue'
import ActionGrid from '@/components/biz/ActionGrid.vue'
import TodoCenter from '@/components/biz/TodoCenter.vue'
import AnnouncementBar from '@/components/biz/AnnouncementBar.vue'
import AppIcon from '@/components/common/AppIcon.vue'
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

/* ---------- 日期 ---------- */
const todayText = computed(() => {
  const d = new Date()
  const week = ['日', '一', '二', '三', '四', '五', '六'][d.getDay()]
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 · 星期${week}`
})

/* ---------- 趋势切换 ---------- */
const trendRange = ref<'7d' | '30d'>('7d')
const trendData = computed<ChartCategoriesSeries>(() => {
  const trend = data.value.sales_trend || []
  const days = trendRange.value === '7d' ? 7 : 30
  const sliced = trend.slice(-days)
  return {
    categories: sliced.map((item) => item.date.slice(5)),
    series: [{ name: '营收', data: sliced.map((item) => item.sales) }],
  }
})
function toggleTrend() {
  trendRange.value = trendRange.value === '7d' ? '30d' : '7d'
}

/* ---------- 图表数据 ---------- */
const statusPie = computed<ChartPie>(() => data.value.status_distribution)
const hotBarData = computed<ChartCategoriesSeries>(() => {
  const items = data.value.hot_products || []
  return {
    categories: items.map((p) => (p.title.length > 6 ? p.title.slice(0, 6) + '…' : p.title)),
    series: [{ name: '销量', data: items.map((p) => p.sold_qty) }],
  }
})

/* ---------- 公告 ---------- */
const announcements = computed(() =>
  (data.value.announcements || []).map((a) => ({
    id: a.id,
    title: a.title,
    type: (a.type === 'urgent' ? 'urgent' : 'info') as 'info' | 'urgent',
  })),
)

/* ---------- 待办 ---------- */
const todoItems = computed(() => [
  { key: 'pending_ship', title: t('dashboard.pendingShip'), value: data.value.pending_ship, tone: 'warning' as const },
  { key: 'pending_after_sale', title: t('dashboard.pendingAfterSale'), value: data.value.pending_after_sale, tone: 'danger' as const },
  { key: 'unread_message', title: t('dashboard.unreadMessage'), value: data.value.unread_message, tone: 'primary' as const },
  { key: 'stock_warning', title: t('dashboard.stockWarning'), value: data.value.stock_warning, tone: 'warning' as const },
])

/* ---------- 快捷入口 / 营销中心 ---------- */
const actions = [
  { key: 'order', label: '订单管理', iconName: 'order' },
  { key: 'product', label: '商品管理', iconName: 'product' },
  { key: 'after_sale', label: '售后处理', iconName: 'after-sale' },
  { key: 'review', label: '评价管理', iconName: 'star' },
  { key: 'decor', label: '店铺装修', iconName: 'decor' },
  { key: 'data', label: '数据报表', iconName: 'chart' },
  { key: 'wms', label: '仓储管理', iconName: 'warehouse' },
]

const marketingActions = [
  { key: 'coupon', label: t('marketing.coupon'), iconName: 'coupon' },
  { key: 'seckill', label: t('marketing.seckill'), iconName: 'flash' },
  { key: 'groupBuy', label: t('marketing.groupBuy'), iconName: 'users' },
  { key: 'bargain', label: t('marketing.bargain'), iconName: 'money' },
  { key: 'vip', label: t('marketing.vip'), iconName: 'crown' },
  { key: 'messageBroadcast', label: '消息群发', iconName: 'megaphone' },
  { key: 'decor', label: t('marketing.decor'), iconName: 'decor' },
  { key: 'specTemplate', label: t('marketing.specTemplate'), iconName: 'clipboard' },
  { key: 'distribution', label: t('marketing.distribution'), iconName: 'share' },
  { key: 'pointsMall', label: t('marketing.pointsMall'), iconName: 'diamond' },
]

const marketingRoutes: Record<string, string> = {
  coupon: '/pages/marketing/coupon',
  seckill: '/pages/marketing/seckill',
  groupBuy: '/pages/marketing/group-buy',
  bargain: '/pages/marketing/bargain',
  vip: '/pages/vip/plan-list',
  messageBroadcast: '/pages/me/messages',
  decor: '/pages/decor/editor',
  specTemplate: '/pages/product/spec-templates',
  distribution: '/pages/distribution/index',
  pointsMall: '/pages/points/index',
}

function onMarketingClick(key: string) {
  const path = marketingRoutes[key]
  if (path) uni.navigateTo({ url: path })
}

/* ---------- 生命周期 ---------- */
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

function goProductWarning() {
  setStorage('eapp_product_status_filter', 'warning')
  uni.navigateTo({ url: '/pages/product/list' })
}

function onTodoClick(key: string) {
  if (key === 'pending_ship') goOrder('2')
  else if (key === 'pending_after_sale') uni.navigateTo({ url: '/pages/order/after-sale-list' })
  else if (key === 'unread_message') uni.switchTab({ url: '/pages/me/im-sessions' })
  else if (key === 'stock_warning') goProductWarning()
  else if (key === 'today_orders') goOrder('')
}

function onActionClick(key: string) {
  if (key === 'order') uni.switchTab({ url: '/pages/order/list' })
  else if (key === 'product') uni.navigateTo({ url: '/pages/product/list' })
  else if (key === 'after_sale') uni.navigateTo({ url: '/pages/order/after-sale-list' })
  else if (key === 'review') uni.navigateTo({ url: '/pages/review/list' })
  else if (key === 'decor') uni.navigateTo({ url: '/pages/decor/editor' })
  else if (key === 'data') uni.navigateTo({ url: '/pages/analytics/index' })
  else if (key === 'wms') uni.navigateTo({ url: '/pages/wms/index' })
}

function onScan() {
  uni.scanCode({
    success: (res) => {
      setStorage('eapp_order_keyword_filter', String(res.result || ''))
      uni.switchTab({ url: '/pages/order/list' })
    },
    fail: () => uni.showToast({ title: '已取消扫码', icon: 'none' }),
  })
}

function onMessage() {
  uni.switchTab({ url: '/pages/me/im-sessions' })
}

function onAnnouncement(id: string | number) {
  uni.showModal({
    title: '公告通知',
    content: announcements.value.find((a) => a.id === id)?.title || '',
    showCancel: false,
    confirmText: '知道了',
  })
}

function goAnalytics() {
  uni.navigateTo({ url: '/pages/analytics/index' })
}

onLoad(loadAll)
onShow(loadAll)
</script>

<template>
  <EappShell :padded="false" header-sticky>
    <template #header>
      <PageHeader eyebrow="MERCHANT CONSOLE" :title="shopStore.name || t('dashboard.title')">
        <template #right>
          <view class="eapp-icon-btn" hover-class="eapp-hover" @click="onScan">
            <AppIcon name="scan" :size="32" color="secondary" />
          </view>
          <view class="eapp-icon-btn" hover-class="eapp-hover" @click="onMessage">
            <AppIcon name="bell" :size="32" color="secondary" />
            <view v-if="data.unread_message" class="eapp-icon-btn__dot" />
          </view>
        </template>
      </PageHeader>
    </template>

    <view class="dashboard-body">
      <!-- 店铺概览条 -->
      <view class="overview">
        <view class="overview__logo">
          <AppIcon name="store" :size="40" color="white" :stroke-width="1.6" />
        </view>
        <view class="overview__info">
          <view class="overview__name-row">
            <text class="overview__name">{{ shopStore.name || t('dashboard.title') }}</text>
            <view class="overview__status">
              <view class="overview__status-dot" />
              <text class="overview__status-text">营业中</text>
            </view>
          </view>
          <text class="overview__date eapp-num">{{ todayText }}</text>
        </view>
        <view class="overview__link" hover-class="eapp-hover" @click="goAnalytics">
          <text class="overview__link-text">经营报表</text>
          <AppIcon name="chevron-right" :size="24" color="primary" />
        </view>
      </view>

      <!-- 今日指标 -->
      <view class="section">
        <SectionHeader title="今日营业概览" extra="查看报表" arrow @click="goAnalytics" />
        <view class="metric-row">
          <MetricCard
            :title="t('dashboard.todaySales')"
            :value="`¥${Number(data.today_sales || 0).toFixed(2)}`"
            :compare="data.compare?.revenue_mom"
            color="var(--eapp-primary)"
          />
          <MetricCard
            :title="t('dashboard.todayOrders')"
            :value="String(data.today_orders || 0)"
            :compare="data.compare?.order_mom"
            color="var(--eapp-success)"
          />
          <MetricCard
            :title="t('dashboard.todayAvgPrice')"
            :value="`¥${Number(data.today_avg_price || 0).toFixed(0)}`"
            color="var(--eapp-warning)"
          />
        </view>
      </view>

      <!-- 待办中心 -->
      <view class="section">
        <SectionHeader title="待办中心" :subtitle="`${todoItems.length} 项需关注`" />
        <TodoCenter :items="todoItems" @click="onTodoClick" />
      </view>

      <!-- 公告通知 -->
      <view v-if="announcements.length" class="section">
        <SectionHeader title="公告通知" />
        <AnnouncementBar :items="announcements" @click="onAnnouncement" />
      </view>

      <!-- 快捷入口 -->
      <view class="section">
        <SectionHeader title="快捷入口" />
        <view class="panel">
          <ActionGrid :items="actions" :columns="4" @click="onActionClick" />
        </view>
      </view>

      <!-- 营销中心 -->
      <view class="section">
        <SectionHeader title="营销中心" />
        <view class="panel">
          <ActionGrid :items="marketingActions" :columns="5" @click="onMarketingClick" />
        </view>
      </view>

      <!-- 经营分析 -->
      <view class="section">
        <SectionHeader title="经营分析" extra="全部报表" arrow @click="goAnalytics" />
        <AreaChart title="营收趋势" :loading="loading" :data="trendData">
          <template #extra>
            <view class="trend-toggle" hover-class="eapp-hover" @click="toggleTrend">
              <AppIcon name="calendar" :size="22" color="primary" />
              <text class="trend-toggle__text">{{ trendRange === '7d' ? '近7天' : '近30天' }}</text>
              <AppIcon name="chevron-down" :size="20" color="primary" />
            </view>
          </template>
        </AreaChart>
        <RingChart title="订单状态分布" :loading="loading" :data="statusPie" />
        <BarChart title="商品销量 Top5" :loading="loading" :data="hotBarData" horizontal />
      </view>

      <!-- 库存预警 -->
      <view v-if="data.stock_warning_list?.length" class="section">
        <SectionHeader title="库存预警" extra="查看全部" arrow @click="goProductWarning" />
        <view class="stock-list">
          <view
            v-for="item in data.stock_warning_list"
            :key="item.product_id"
            class="stock-item"
            hover-class="eapp-hover"
            @click="goProductWarning"
          >
            <view class="stock-item__icon">
              <AppIcon name="warning" :size="30" color="danger" />
            </view>
            <view class="stock-item__info">
              <text class="stock-item__title">{{ item.title }}</text>
              <text class="stock-item__sub eapp-num">预警阈值 {{ item.threshold }}</text>
            </view>
            <text class="stock-item__qty eapp-num" :class="{ 'stock-item__qty--danger': item.stock < item.threshold }">
              {{ item.stock }}
            </text>
            <AppIcon name="chevron-right" :size="24" color="faint" />
          </view>
        </view>
      </view>

      <view v-if="loading" class="loading-hint">
        <up-loading-icon mode="circle" size="28" />
        <text class="loading-hint__text">正在更新经营数据…</text>
      </view>
    </view>
  </EappShell>
</template>

<style scoped>
.dashboard-body {
  padding: var(--eapp-space-5);
  padding-bottom: calc(env(safe-area-inset-bottom) + 48rpx);
  box-sizing: border-box;
}

.section {
  margin-bottom: var(--eapp-space-6);
}

.panel {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-3) var(--eapp-space-2);
  box-sizing: border-box;
}

/* 店铺概览条 */
.overview {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-4);
  margin-bottom: var(--eapp-space-6);
  box-sizing: border-box;
}

.overview__logo {
  width: 76rpx;
  height: 76rpx;
  flex-shrink: 0;
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary-deep);
  display: flex;
  align-items: center;
  justify-content: center;
}

.overview__info {
  flex: 1;
  min-width: 0;
}

.overview__name-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.overview__name {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.overview__status {
  display: flex;
  align-items: center;
  gap: 6rpx;
  padding: 2rpx 10rpx;
  border-radius: var(--eapp-radius-sm);
  background: var(--eapp-success-soft);
  flex-shrink: 0;
}

.overview__status-dot {
  width: 8rpx;
  height: 8rpx;
  border-radius: 50%;
  background: var(--eapp-success);
}

.overview__status-text {
  font-size: 18rpx;
  font-weight: 600;
  color: var(--eapp-success);
}

.overview__date {
  display: block;
  margin-top: 4rpx;
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}

.overview__link {
  display: flex;
  align-items: center;
  gap: 2rpx;
  flex-shrink: 0;
}

.overview__link-text {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-primary);
}

/* 指标行 */
.metric-row {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--eapp-space-2);
}

/* 趋势切换 */
.trend-toggle {
  display: flex;
  align-items: center;
  gap: 6rpx;
  padding: 6rpx 14rpx;
  border: 1rpx solid var(--eapp-primary-border);
  border-radius: var(--eapp-radius-sm);
  background: var(--eapp-primary-soft);
}

.trend-toggle__text {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-primary);
  font-weight: 600;
}

/* 库存预警 */
.stock-list {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  overflow: hidden;
}

.stock-item {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
  padding: var(--eapp-space-3) var(--eapp-space-4);
  border-bottom: 1rpx solid var(--eapp-divider);
}

.stock-item:last-child {
  border-bottom: none;
}

.stock-item__icon {
  width: 56rpx;
  height: 56rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius);
  background: var(--eapp-danger-soft);
}

.stock-item__info {
  display: flex;
  flex-direction: column;
  gap: 2rpx;
  flex: 1;
  min-width: 0;
}

.stock-item__title {
  font-size: var(--eapp-fs-base);
  font-weight: 500;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stock-item__sub {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}

.stock-item__qty {
  font-size: var(--eapp-fs-lg);
  font-weight: 700;
  color: var(--eapp-text-secondary);
  flex-shrink: 0;
}

.stock-item__qty--danger {
  color: var(--eapp-danger);
}

.loading-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  padding: var(--eapp-space-4) 0;
}

.loading-hint__text {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-faint);
}
</style>
