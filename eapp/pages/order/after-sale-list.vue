<script setup lang="ts">
import { onLoad, onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import AfterSaleCard from '@/components/biz/AfterSaleCard.vue'
import EmptyState from '@/components/biz/EmptyState.vue'
import { getAfterSales } from '@/api/after-sale'

const tabs = [
  { name: '全部', status: '' },
  { name: '待审核', status: 'applied' },
  { name: '退货中', status: 'user_returning' },
  { name: '退款中', status: 'refund_pending' },
  { name: '已完成', status: 'refunded' },
  { name: '已关闭', status: 'closed' },
]
const typeNames = ['全部', '仅退款', '退货退款', '换货']
const typeValues = ['', 'refund_only', 'return_refund', 'exchange']

const current = ref(0)
const list = ref<any[]>([])
const loading = ref(false)
const filterType = ref('')
const showTypePicker = ref(false)
const typeOptions = typeNames.map((n, i) => ({ label: n, value: typeValues[i] }))

async function load() {
  loading.value = true
  try {
    const res: any = await getAfterSales({ page: 1, size: 50, status: tabs[current.value]?.status || undefined, type: filterType.value || undefined })
    list.value = Array.isArray(res?.list) ? res.list : []
  } finally { loading.value = false }
}

function onTab(idx: number) { current.value = idx; load() }
function onType(e: any) { filterType.value = typeValues[Number(e.detail.value)] || ''; load() }
function goDetail(id: number) { uni.navigateTo({ url: `/pages/order/after-sale-detail?id=${id}` }) }

onLoad(load); onShow(load)
</script>

<template>
  <view class="page">
    <up-tabs
      :list="tabs"
      :current="current"
      :scrollable="true"
      keyName="name"
      @click="(item) => onTab(item.index)"
      :activeStyle="{ color: '#fff', backgroundColor: 'var(--eapp-primary)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :inactiveStyle="{ color: 'var(--eapp-text-muted)', backgroundColor: 'var(--eapp-bg)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :itemStyle="{ padding: '0 4rpx', height: '80rpx' }"
      lineColor="transparent"
    />
    <view class="filter">
      <view class="picker" @click="showTypePicker = true">类型：{{ typeNames[typeValues.indexOf(filterType)] }}</view>
      <up-picker :show="showTypePicker" :columns="[typeOptions]" keyName="label" @confirm="(e: any) => { filterType = e.value[0].value; showTypePicker = false; load() }" @cancel="showTypePicker = false" @close="showTypePicker = false" />
    </view>
    <view class="list">
      <EmptyState v-if="!loading && !list.length" title="暂无售后" />
      <AfterSaleCard v-for="row in list" :key="row.id" :row="row" @click="goDetail(row.id)" />
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.filter { padding: 16rpx 20rpx; }
.picker { display: inline-flex; height: 60rpx; align-items: center; padding: 0 18rpx; border: 1px solid var(--eapp-border); border-radius: 12rpx; background: var(--eapp-card); }
.list { padding: 20rpx; display: grid; gap: 14rpx; }
</style>
