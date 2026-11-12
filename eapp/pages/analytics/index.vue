<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import SectionHeader from '@/components/biz/SectionHeader.vue'
import EntryList, { type EntryItem } from '@/components/biz/EntryList.vue'
import AppIcon from '@/components/common/AppIcon.vue'

const { t } = useI18n()

const reports: EntryItem[] = [
  { key: 'sales', title: t('analytics.sales'), desc: 'GMV、订单数、客单价与支付方式', iconName: 'trend-up', path: '/pages/analytics/sales' },
  { key: 'products', title: t('analytics.products'), desc: 'SKU 排行、动销率与库存周转', iconName: 'product', path: '/pages/analytics/products' },
  { key: 'customers', title: t('analytics.customers'), desc: '复购率、客均消费与 RFM 分层', iconName: 'users', path: '/pages/analytics/customers' },
  { key: 'traffic', title: t('analytics.traffic'), desc: 'PV / UV、渠道与设备分布', iconName: 'eye', path: '/pages/analytics/traffic' },
  { key: 'conversion', title: t('analytics.conversion'), desc: '转化漏斗与弃购分析', iconName: 'pie', path: '/pages/analytics/conversion' },
]

function go(item: EntryItem) {
  if (item.path) uni.navigateTo({ url: item.path })
}
</script>

<template>
  <view class="page">
    <view class="intro">
      <view class="intro__icon">
        <AppIcon name="chart" :size="40" color="primary" />
      </view>
      <view class="intro__text">
        <text class="intro__title">数据分析</text>
        <text class="intro__desc">按专题查看经营指标与趋势图表</text>
      </view>
    </view>

    <view class="group">
      <SectionHeader title="经营报表" :subtitle="`共 ${reports.length} 个专题`" />
      <EntryList :items="reports" @click="go" />
    </view>

    <view class="note">
      <AppIcon name="info" :size="26" color="faint" />
      <text class="note__text">各专题页支持日期区间切换，数据按所选区间统计。</text>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background: var(--eapp-bg);
  padding: var(--eapp-space-5);
  padding-bottom: calc(env(safe-area-inset-bottom) + 48rpx);
  box-sizing: border-box;
}

.intro {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
  padding: var(--eapp-space-4) var(--eapp-space-5);
  margin-bottom: var(--eapp-space-6);
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  box-sizing: border-box;
}

.intro__icon {
  width: 72rpx;
  height: 72rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary-soft);
  border: 1rpx solid var(--eapp-primary-border);
  box-sizing: border-box;
}

.intro__text {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  min-width: 0;
}

.intro__title {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
}

.intro__desc {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
}

.group {
  margin-bottom: var(--eapp-space-5);
}

.note {
  display: flex;
  align-items: center;
  gap: 10rpx;
  padding: 0 var(--eapp-space-2);
}

.note__text {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}
</style>
