<script setup lang="ts">
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useI18n } from 'vue-i18n'
import PageHeader from '@/components/biz/PageHeader.vue'
import SectionHeader from '@/components/biz/SectionHeader.vue'
import EntryList, { type EntryItem } from '@/components/biz/EntryList.vue'
import AppIcon from '@/components/common/AppIcon.vue'
import { getActivities, getCoupons } from '@/api/marketing'
import {
  ACTIVITY_KINDS,
  ACTIVITY_META,
  MARKETING_ROUTES,
  activityState,
} from '@/utils/marketing'

const { t } = useI18n()

const overview = ref({ running: 0, upcoming: 0, total: 0, coupons: 0 })
const activityCounts = ref<Record<string, { total: number; running: number }>>({})

const activityTiles = computed(() => [
  ...ACTIVITY_KINDS.map((kind) => {
    const meta = ACTIVITY_META[kind]
    const count = activityCounts.value[kind] || { total: 0, running: 0 }
    return {
      key: kind,
      title: meta.title,
      desc: meta.desc,
      iconName: meta.iconName,
      path: meta.route,
      foot: `${count.running} / ${count.total} 进行中`,
    }
  }),
  {
    key: 'coupon',
    title: '优惠券',
    desc: '满减 / 折扣券与定向发券',
    iconName: 'coupon',
    path: MARKETING_ROUTES.coupon,
    foot: `${overview.value.coupons} 张可用`,
  },
])

const promotions: EntryItem[] = [
  {
    key: 'coupon',
    title: t('marketing.coupon'),
    desc: '满减 / 折扣券与定向发券',
    iconName: 'coupon',
    path: MARKETING_ROUTES.coupon,
  },
  ...ACTIVITY_KINDS.map((kind) => ({
    key: kind,
    title: ACTIVITY_META[kind].title,
    desc: ACTIVITY_META[kind].desc,
    iconName: ACTIVITY_META[kind].iconName,
    path: ACTIVITY_META[kind].route,
  })),
]

const operations: EntryItem[] = [
  { key: 'vip', title: t('marketing.vip'), desc: '会员套餐 / 等级 / 专属价', iconName: 'crown', path: MARKETING_ROUTES.vip },
  { key: 'distribution', title: t('marketing.distribution'), desc: '分销商、佣金与比例配置', iconName: 'share', path: MARKETING_ROUTES.distribution },
  { key: 'pointsMall', title: t('marketing.pointsMall'), desc: '积分商品与兑换记录', iconName: 'diamond', path: MARKETING_ROUTES.pointsMall },
  { key: 'messageBroadcast', title: '消息群发', desc: '按分组与优先级触达客户', iconName: 'megaphone', path: MARKETING_ROUTES.messageBroadcast },
]

const tools: EntryItem[] = [
  { key: 'decor', title: t('marketing.decor'), desc: '移动端首页组件编排', iconName: 'decor', path: MARKETING_ROUTES.decor },
  { key: 'specTemplate', title: t('marketing.specTemplate'), desc: '预定义属性组与 SKU 生成', iconName: 'clipboard', path: MARKETING_ROUTES.specTemplate },
]

async function load() {
  const [seckill, groupBuy, bargain, coupons] = await Promise.all([
    getActivities('seckill', { page: 1, size: 200 }).catch(() => null),
    getActivities('group-buy', { page: 1, size: 200 }).catch(() => null),
    getActivities('bargain', { page: 1, size: 200 }).catch(() => null),
    getCoupons({ page: 1, size: 200 }).catch(() => null),
  ]) as any[]

  const groups = [seckill, groupBuy, bargain]
  const counts: Record<string, { total: number; running: number }> = {}
  let running = 0
  let upcoming = 0
  let total = 0
  ACTIVITY_KINDS.forEach((kind, index) => {
    const list = Array.isArray(groups[index]?.list) ? groups[index].list : []
    const live = list.filter((item: any) => activityState(item).phase === 'running').length
    const planned = list.filter((item: any) => activityState(item).phase === 'upcoming').length
    counts[kind] = { total: list.length, running: live }
    running += live
    upcoming += planned
    total += list.length
  })
  const couponList = Array.isArray(coupons?.list) ? coupons.list : Array.isArray(coupons) ? coupons : []
  activityCounts.value = counts
  overview.value = { running, upcoming, total, coupons: couponList.length }
}

function go(item: EntryItem) {
  if (item.path) uni.navigateTo({ url: item.path })
}

function goTile(tile: { path: string }) {
  uni.navigateTo({ url: tile.path })
}

onShow(load)
</script>

<template>
  <view class="marketing-page">
    <PageHeader title="营销中心" eyebrow="MARKETING CENTER" subtitle="活动、优惠与客户运营" back />

    <view class="marketing-page__body">
      <view class="hero">
        <view class="hero__main">
          <text class="hero__eyebrow">MARKETING OVERVIEW</text>
          <view class="hero__headline">
            <text class="hero__number eapp-num">{{ overview.running }}</text>
            <text class="hero__unit">场活动进行中</text>
          </view>
          <view class="hero__metrics">
            <view class="hero__metric">
              <text class="hero__metric-value eapp-num">{{ overview.upcoming }}</text>
              <text class="hero__metric-label">待开始</text>
            </view>
            <view class="hero__metric">
              <text class="hero__metric-value eapp-num">{{ overview.coupons }}</text>
              <text class="hero__metric-label">优惠券</text>
            </view>
            <view class="hero__metric">
              <text class="hero__metric-value eapp-num">{{ overview.total }}</text>
              <text class="hero__metric-label">活动总数</text>
            </view>
          </view>
        </view>
        <view class="hero__badge">
          <AppIcon name="target" :size="64" color="white" :stroke-width="1.3" />
        </view>
      </view>

      <SectionHeader title="活动管理" :subtitle="`${overview.total} 场活动 · ${overview.coupons} 张券`" />
      <view class="tiles">
        <view
          v-for="tile in activityTiles"
          :key="tile.key"
          class="tile"
          hover-class="eapp-hover"
          @click="goTile(tile)"
        >
          <view class="tile__icon">
            <AppIcon :name="tile.iconName" :size="38" color="primary" />
          </view>
          <text class="tile__title">{{ tile.title }}</text>
          <text class="tile__desc">{{ tile.desc }}</text>
          <view class="tile__foot">
            <text class="tile__count eapp-num">{{ tile.foot }}</text>
            <AppIcon name="chevron-right" :size="24" color="faint" />
          </view>
        </view>
      </view>

      <view class="group">
        <SectionHeader title="全部功能" />
        <EntryList :items="promotions" @click="go" />
      </view>

      <view class="group">
        <SectionHeader title="会员与运营" />
        <EntryList :items="operations" @click="go" />
      </view>

      <view class="group">
        <SectionHeader title="内容与工具" />
        <EntryList :items="tools" @click="go" />
      </view>
    </view>
  </view>
</template>

<style scoped>
.marketing-page {
  min-height: 100vh;
  background: var(--eapp-bg);
}
.marketing-page__body {
  padding: var(--eapp-space-5);
  padding-bottom: calc(env(safe-area-inset-bottom) + 48rpx);
  box-sizing: border-box;
}
.hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--eapp-space-4);
  padding: var(--eapp-space-5);
  margin-bottom: var(--eapp-space-6);
  border-radius: var(--eapp-radius-lg);
  background: var(--eapp-primary-deep);
  box-sizing: border-box;
}
.hero__main {
  flex: 1;
  min-width: 0;
}
.hero__eyebrow {
  font-size: var(--eapp-fs-xs);
  letter-spacing: 2rpx;
  color: #f4b8b1;
}
.hero__headline {
  display: flex;
  align-items: baseline;
  gap: var(--eapp-space-2);
  margin-top: var(--eapp-space-2);
}
.hero__number {
  font-size: 72rpx;
  font-weight: 700;
  line-height: 1.1;
  color: #ffffff;
}
.hero__unit {
  font-size: var(--eapp-fs-sm);
  color: #f9d8d4;
}
.hero__metrics {
  display: flex;
  gap: var(--eapp-space-6);
  margin-top: var(--eapp-space-4);
  padding-top: var(--eapp-space-3);
  border-top: 1rpx solid rgba(255, 255, 255, 0.16);
}
.hero__metric {
  display: flex;
  flex-direction: column;
  gap: 2rpx;
}
.hero__metric-value {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: #ffffff;
}
.hero__metric-label {
  font-size: var(--eapp-fs-xs);
  color: #f4b8b1;
}
.hero__badge {
  width: 120rpx;
  height: 120rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  border: 1rpx solid rgba(255, 255, 255, 0.18);
  background: rgba(255, 255, 255, 0.08);
}
.group {
  margin-bottom: var(--eapp-space-6);
}
.tiles {
  display: flex;
  flex-wrap: wrap;
  gap: var(--eapp-space-3);
  margin-bottom: var(--eapp-space-6);
}
.tile {
  width: calc((100% - var(--eapp-space-3)) / 2);
  padding: var(--eapp-space-4);
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  box-sizing: border-box;
}
.tile__icon {
  width: 68rpx;
  height: 68rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary-soft);
  border: 1rpx solid var(--eapp-primary-border);
  box-sizing: border-box;
}
.tile__title {
  display: block;
  margin-top: var(--eapp-space-3);
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
}
.tile__desc {
  display: block;
  margin-top: 2rpx;
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-muted);
  line-height: 1.5;
}
.tile__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: var(--eapp-space-3);
  padding-top: var(--eapp-space-2);
  border-top: 1rpx solid var(--eapp-divider);
}
.tile__count {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-primary);
}
</style>
