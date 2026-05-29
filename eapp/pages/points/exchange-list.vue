<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePointsList } from '@/composables/usePointsList'
import { shipExchange, completeExchange } from '@/api/points'

const { t } = useI18n()
const h = usePointsList('exchanges')
const statusTab = ref('')
const tabs = [
  { label: '全部', value: '' },
  { label: '待发货', value: 'pending_ship' },
  { label: '已完成', value: 'completed' },
]

const statusLabels: Record<string, string> = { pending_ship: '待发货', shipped: '已发货', completed: '已完成' }
const statusColors: Record<string, string> = { pending_ship: 'pending', shipped: 'info', completed: 'on' }

function doSearch() {
  const params: any = { page: 1, size: 50 }
  if (statusTab.value) params.status = statusTab.value
  h.load(params)
}

function switchTab(val: string) {
  statusTab.value = val
  doSearch()
}

async function doShip(id: number) {
  await shipExchange(id)
  uni.showToast({ title: '已发货', icon: 'success' })
  doSearch()
}

async function doComplete(id: number) {
  await completeExchange(id)
  uni.showToast({ title: '已完成', icon: 'success' })
  doSearch()
}

onShow(() => doSearch())
</script>

<template>
  <view class="page">
    <view class="tabs">
      <view v-for="tab in tabs" :key="tab.value" :class="['tab', statusTab === tab.value ? 'active' : '']" @click="switchTab(tab.value)">{{ tab.label }}</view>
    </view>
    <view v-if="!h.loading.value && !h.list.value.length" class="empty">暂无兑换记录</view>
    <view v-for="item in h.list.value" :key="item.id" class="card">
      <view class="row">
        <text class="name">{{ item.user_nickname }}</text>
        <text :class="['badge', statusColors[item.status] || '']">{{ statusLabels[item.status] || item.status }}</text>
      </view>
      <view class="desc">{{ item.product_title }}</view>
      <view class="desc">{{ item.points_cost }} 积分</view>
      <view v-if="item.status === 'pending_ship'" class="actions">
        <up-button size="mini" type="primary" plain @click="doShip(item.id)">发货</up-button>
        <up-button size="mini" type="success" plain @click="doComplete(item.id)">完成</up-button>
      </view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中…</view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; align-content: start; }
.tabs { display: flex; gap: 10rpx; }
.tab { padding: 8rpx 20rpx; border-radius: 999rpx; font-size: 24rpx; background: #f1f5f9; color: #64748b; }
.tab.active { background: #2563eb; color: #fff; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.row { display: flex; align-items: center; justify-content: space-between; }
.name { font-size: 28rpx; font-weight: 600; }
.badge { font-size: 22rpx; padding: 4rpx 14rpx; border-radius: 999rpx; }
.badge.pending { background: #fef3c7; color: #d97706; }
.badge.info { background: #dbeafe; color: #2563eb; }
.badge.on { background: #dcfce7; color: #16a34a; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.actions { margin-top: 12rpx; display: flex; gap: 10rpx; justify-content: flex-end; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>
