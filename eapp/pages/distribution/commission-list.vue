<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDistributionList } from '@/composables/useDistributionList'
import { settleCommission, returnCommission } from '@/api/distribution'

const { t } = useI18n()
const h = useDistributionList('commissions')
const current = ref(0)
const statusTab = ref<number | ''>('')
const tabs = [
  { label: '全部', value: '' as const },
  { label: '待结算', value: 1 },
  { label: '已结算', value: 2 },
  { label: '已退回', value: 3 },
]

const statusLabels: Record<number, string> = { 1: '待结算', 2: '已结算', 3: '已退回' }
const statusColors: Record<number, string> = { 1: 'pending', 2: 'on', 3: 'off' }

function doSearch() {
  const params: any = { page: 1, size: 50 }
  if (statusTab.value !== '') params.status = statusTab.value
  h.load(params)
}

function switchTab(item: any) {
  current.value = item.index
  statusTab.value = tabs[item.index].value
  doSearch()
}

async function doSettle(id: number) {
  await settleCommission(id)
  uni.showToast({ title: '已结算', icon: 'success' })
  doSearch()
}

async function doReturn(id: number) {
  await returnCommission(id)
  uni.showToast({ title: '已退回', icon: 'success' })
  doSearch()
}

onShow(() => doSearch())
</script>

<template>
  <view class="page">
    <up-tabs
      :list="tabs"
      :current="current"
      :scrollable="true"
      keyName="label"
      @click="switchTab"
      :activeStyle="{ color: '#fff', backgroundColor: 'var(--eapp-primary)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :inactiveStyle="{ color: 'var(--eapp-text-muted)', backgroundColor: 'var(--eapp-bg)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :itemStyle="{ padding: '0 4rpx', height: '80rpx' }"
      lineColor="transparent"
    />
    <view v-if="!h.loading.value && !h.list.value.length" class="empty">暂无佣金记录</view>
    <view v-for="item in h.list.value" :key="item.id" class="card">
      <view class="row">
        <text class="order-id">订单 #{{ item.order_id }}</text>
        <text :class="['badge', statusColors[item.status] || '']">{{ statusLabels[item.status] || '未知' }}</text>
      </view>
      <view class="desc">分销商 ID: {{ item.distributor_id }} · {{ item.level === 1 ? '直推' : '间推' }}</view>
      <view class="amount">¥{{ Number(item.amount || 0).toFixed(2) }}</view>
      <view v-if="item.status === 1" class="actions">
        <up-button size="mini" type="primary" plain @click="doSettle(item.id)">结算</up-button>
        <up-button size="mini" type="error" plain @click="doReturn(item.id)">退回</up-button>
      </view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中…</view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; align-content: start; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.row { display: flex; align-items: center; justify-content: space-between; }
.order-id { font-size: 28rpx; font-weight: 600; }
.badge { font-size: 22rpx; padding: 4rpx 14rpx; border-radius: 999rpx; }
.badge.pending { background: #fef3c7; color: #d97706; }
.badge.on { background: #dcfce7; color: #16a34a; }
.badge.off { background: #fee2e2; color: #dc2626; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.amount { margin-top: 6rpx; font-size: 30rpx; font-weight: 700; color: #ef4444; }
.actions { margin-top: 12rpx; display: flex; gap: 10rpx; justify-content: flex-end; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>
