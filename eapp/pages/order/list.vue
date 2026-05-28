<script setup lang="ts">
import { onLoad, onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import StatusTag from '@/components/common/StatusTag.vue'
import { getOrders } from '@/api/order'
import { getStorage, removeStorage } from '@/utils/storage'

const { t } = useI18n()
const current = ref(0)
const list = ref<any[]>([])
const loading = ref(false)

const tabs = [
  { name: t('order.all'), status: '' },
  { name: t('order.pendingPay'), status: '1' },
  { name: t('order.pendingShip'), status: '2' },
  { name: t('order.shipped'), status: '3' },
  { name: t('order.completed'), status: '4' },
  { name: t('order.closed'), status: '5' },
]

async function loadData() {
  loading.value = true
  try {
    const status = tabs[current.value]?.status || ''
    const data: any = await getOrders({ page: 1, size: 20, status })
    list.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

function onTabChange(item: any) {
  current.value = Number(item?.index || 0)
  loadData()
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/order/detail?id=${id}` })
}

onLoad((opts) => {
  const status = String(opts?.status || '')
  const idx = tabs.findIndex((tab) => tab.status === status)
  if (idx >= 0) current.value = idx
})
onShow(() => {
  const pendingStatus = String(getStorage('eapp_order_status_filter') || '')
  if (pendingStatus) {
    const idx = tabs.findIndex((tab) => tab.status === pendingStatus)
    if (idx >= 0) current.value = idx
    removeStorage('eapp_order_status_filter')
  }
  loadData()
})
</script>

<template>
  <view class="page">
    <up-tabs :list="tabs" :current="current" @click="onTabChange" />
    <view class="list">
      <view v-if="loading" class="empty">{{ t('common.loading') }}</view>
      <view v-else-if="!list.length" class="empty">{{ t('common.empty') }}</view>
      <view v-for="item in list" :key="item.id" class="card" @click="goDetail(item.id)">
        <view class="row">
          <text class="no">#{{ item.id }}</text>
          <StatusTag :text="item.status_label || item.status || '-'" :type="item.status" />
        </view>
        <view class="desc">{{ item.user_nickname || item.receiver_name || '用户' }}</view>
        <view class="amount">¥{{ Number(item.pay_amount || item.total_amount || 0).toFixed(2) }}</view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }
.list { padding: 20rpx; display: grid; gap: 16rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; }
.row { display: flex; align-items: center; justify-content: space-between; }
.no { font-size: 24rpx; color: var(--eapp-text-muted); }
.desc { margin-top: 10rpx; font-size: 26rpx; }
.amount { margin-top: 8rpx; color: var(--eapp-primary); font-size: 32rpx; font-weight: 700; }
.empty { padding: 80rpx 0; text-align: center; color: var(--eapp-text-muted); }
</style>
