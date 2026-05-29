<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref, computed } from 'vue'
import { useWmsList } from '@/composables/useWmsList'
import { getWarehouses, updateSafeQty } from '@/api/wms'

const h = useWmsList('stocks')
const warehouses = ref<any[]>([])
const warehouseId = ref(0)
const keyword = ref('')

const filteredList = computed(() => {
  let list = h.list.value
  if (warehouseId.value) list = list.filter((r: any) => Number(r.warehouse_id) === warehouseId.value)
  if (keyword.value.trim()) {
    const kw = keyword.value.trim().toLowerCase()
    list = list.filter((r: any) => String(r.sku_name || '').toLowerCase().includes(kw))
  }
  return list
})

async function loadWarehouses() {
  const res: any = await getWarehouses()
  warehouses.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
}

async function onSafeQtyChange(item: any, val: string) {
  const qty = Number(val)
  if (isNaN(qty) || qty < 0) return
  await updateSafeQty(item.id, qty)
  item.safe_qty = qty
  uni.showToast({ title: '已更新', icon: 'success' })
}

function selectWarehouse(id: number) {
  warehouseId.value = warehouseId.value === id ? 0 : id
}

onShow(async () => {
  await Promise.all([h.load(), loadWarehouses()])
})
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <text class="title">库存台账</text>
    </view>

    <view class="filters">
      <scroll-view scroll-x class="wh-scroll">
        <view class="wh-tags">
          <view :class="['wh-tag', !warehouseId ? 'active' : '']" @click="selectWarehouse(0)">全部</view>
          <view v-for="w in warehouses" :key="w.id" :class="['wh-tag', warehouseId === w.id ? 'active' : '']" @click="selectWarehouse(w.id)">{{ w.name }}</view>
        </view>
      </scroll-view>
      <up-input v-model="keyword" placeholder="搜索SKU名称" clearable />
    </view>

    <view v-if="!h.loading.value && !filteredList.length" class="empty">暂无库存记录</view>
    <view v-for="item in filteredList" :key="item.id" class="card">
      <view class="row">
        <text class="name">{{ item.sku_name }}</text>
        <view v-if="item.qty <= item.safe_qty" class="warning-badge">预警</view>
      </view>
      <view class="desc">仓库: {{ item.warehouse_name }} | 数量: {{ item.qty }}</view>
      <view class="safe-row">
        <text class="safe-label">安全库存:</text>
        <input class="safe-input" type="number" :value="String(item.safe_qty)" @blur="(e: any) => onSafeQtyChange(item, e.detail.value)" />
      </view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中...</view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; align-content: start; }
.top-bar { display: flex; align-items: center; justify-content: space-between; }
.title { font-size: 32rpx; font-weight: 700; }
.filters { display: grid; gap: 10rpx; }
.wh-scroll { white-space: nowrap; }
.wh-tags { display: inline-flex; gap: 10rpx; }
.wh-tag { display: inline-block; padding: 8rpx 20rpx; font-size: 24rpx; border-radius: 999rpx; background: #fff; border: 1px solid var(--eapp-border); }
.wh-tag.active { background: var(--eapp-primary, #2563eb); color: #fff; border-color: var(--eapp-primary, #2563eb); }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.row { display: flex; align-items: center; justify-content: space-between; }
.name { font-size: 28rpx; font-weight: 600; }
.warning-badge { font-size: 20rpx; padding: 4rpx 12rpx; border-radius: 999rpx; background: #fef3c7; color: #d97706; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.safe-row { margin-top: 10rpx; display: flex; align-items: center; gap: 10rpx; }
.safe-label { font-size: 24rpx; color: var(--eapp-text-muted); }
.safe-input { width: 120rpx; font-size: 24rpx; border: 1px solid var(--eapp-border); border-radius: 8rpx; padding: 6rpx 12rpx; text-align: center; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>
