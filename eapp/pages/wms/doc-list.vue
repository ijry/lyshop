<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref, computed } from 'vue'
import { useWmsList } from '@/composables/useWmsList'

const h = useWmsList('docs')
const activeType = ref('')
const activeStatus = ref('')
const currentType = ref(0)
const currentStatus = ref(0)

const typeTabs = [
  { key: '', label: '全部' },
  { key: 'inbound', label: '入库' },
  { key: 'outbound', label: '出库' },
]

const statusTabs = [
  { key: '', label: '全部' },
  { key: 'draft', label: '草稿' },
  { key: 'completed', label: '已完成' },
  { key: 'cancelled', label: '已取消' },
]

const filteredList = computed(() => {
  let list = h.list.value
  if (activeType.value) list = list.filter((r: any) => r.doc_type === activeType.value)
  if (activeStatus.value) list = list.filter((r: any) => r.status === activeStatus.value)
  return list
})

function statusLabel(status: string) {
  if (status === 'draft') return '草稿'
  if (status === 'completed') return '已完成'
  if (status === 'cancelled') return '已取消'
  return status
}

function statusClass(status: string) {
  if (status === 'draft') return 'badge-draft'
  if (status === 'completed') return 'badge-done'
  if (status === 'cancelled') return 'badge-cancel'
  return ''
}

function typeLabel(type: string) {
  return type === 'inbound' ? '入库' : '出库'
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/wms/doc-editor?id=${id}` })
}

function goCreate() {
  uni.navigateTo({ url: '/pages/wms/doc-editor' })
}

onShow(() => h.load())
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <text class="title">出入库单</text>
      <up-button size="mini" type="primary" @click="goCreate">新建</up-button>
    </view>

    <up-tabs
      :list="typeTabs"
      :current="currentType"
      :scrollable="true"
      keyName="label"
      @click="(item) => { currentType = item.index; activeType = typeTabs[item.index].key }"
      :activeStyle="{ color: '#fff', backgroundColor: 'var(--eapp-primary)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :inactiveStyle="{ color: 'var(--eapp-text-muted)', backgroundColor: 'var(--eapp-bg)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :itemStyle="{ padding: '0 4rpx', height: '80rpx' }"
      lineColor="transparent"
    />

    <up-tabs
      :list="statusTabs"
      :current="currentStatus"
      :scrollable="true"
      keyName="label"
      @click="(item) => { currentStatus = item.index; activeStatus = statusTabs[item.index].key }"
      :activeStyle="{ color: '#fff', backgroundColor: 'var(--eapp-primary)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :inactiveStyle="{ color: 'var(--eapp-text-muted)', backgroundColor: 'var(--eapp-bg)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :itemStyle="{ padding: '0 4rpx', height: '80rpx' }"
      lineColor="transparent"
    />

    <view v-if="!h.loading.value && !filteredList.length" class="empty">暂无单据</view>
    <view v-for="item in filteredList" :key="item.id" class="card" @click="goDetail(item.id)">
      <view class="row">
        <text class="name">{{ item.doc_no }}</text>
        <text :class="['badge', statusClass(item.status)]">{{ statusLabel(item.status) }}</text>
      </view>
      <view class="desc">类型: {{ typeLabel(item.doc_type) }} | 仓库: {{ item.warehouse_name || '-' }}</view>
      <view class="desc">总数量: {{ item.total_qty || 0 }} | {{ item.created_at?.slice(0, 10) || '-' }}</view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中...</view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; align-content: start; }
.top-bar { display: flex; align-items: center; justify-content: space-between; }
.title { font-size: 32rpx; font-weight: 700; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.row { display: flex; align-items: center; justify-content: space-between; }
.name { font-size: 28rpx; font-weight: 600; }
.badge { font-size: 22rpx; padding: 4rpx 14rpx; border-radius: 999rpx; }
.badge-draft { background: #e0e7ff; color: #4338ca; }
.badge-done { background: #dcfce7; color: #16a34a; }
.badge-cancel { background: #fee2e2; color: #dc2626; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>
