<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDistributionList } from '@/composables/useDistributionList'
import { updateDistributorStatus } from '@/api/distribution'

const { t } = useI18n()
const h = useDistributionList('distributors')
const keyword = ref('')
const current = ref(0)
const statusTab = ref<number | ''>('')
const tabs = [
  { label: '全部', value: '' as const },
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 },
]

function doSearch() {
  const params: any = { page: 1, size: 50 }
  if (keyword.value.trim()) params.keyword = keyword.value.trim()
  if (statusTab.value !== '') params.status = statusTab.value
  h.load(params)
}

function switchTab(item: any) {
  current.value = item.index
  statusTab.value = tabs[item.index].value
  doSearch()
}

async function toggleStatus(item: any) {
  const newStatus = Number(item.status) === 1 ? 0 : 1
  await updateDistributorStatus(item.id, newStatus)
  doSearch()
}

onShow(() => doSearch())
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <up-input v-model="keyword" :placeholder="t('distribution.searchDistributor')" clearable @clear="doSearch" @confirm="doSearch" />
    </view>
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
    <view v-if="!h.loading.value && !h.list.value.length" class="empty">暂无分销商</view>
    <view v-for="item in h.list.value" :key="item.id" class="card">
      <view class="row">
        <text class="name">{{ item.nickname }}</text>
        <text :class="['badge', item.status === 1 ? 'on' : 'off']">{{ item.status === 1 ? '启用' : '禁用' }}</text>
      </view>
      <view class="desc">{{ item.phone }} · {{ item.level === 1 ? '一级' : '二级' }}</view>
      <view class="desc">余额 ¥{{ Number(item.balance || 0).toFixed(2) }} · 累计 ¥{{ Number(item.total_earn || 0).toFixed(2) }}</view>
      <view class="actions">
        <up-button size="mini" :type="item.status === 1 ? 'error' : 'primary'" plain @click="toggleStatus(item)">{{ item.status === 1 ? '禁用' : '启用' }}</up-button>
      </view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中…</view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; align-content: start; }
.top-bar { display: flex; align-items: center; gap: 10rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.row { display: flex; align-items: center; justify-content: space-between; }
.name { font-size: 28rpx; font-weight: 600; }
.badge { font-size: 22rpx; padding: 4rpx 14rpx; border-radius: 999rpx; }
.badge.on { background: #dcfce7; color: #16a34a; }
.badge.off { background: #fee2e2; color: #dc2626; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.actions { margin-top: 12rpx; display: flex; gap: 10rpx; justify-content: flex-end; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>
