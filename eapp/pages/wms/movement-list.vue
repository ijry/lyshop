<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref, computed } from 'vue'
import { useWmsList } from '@/composables/useWmsList'
import { getWarehouses } from '@/api/wms'

const h = useWmsList('movements')
const warehouses = ref<any[]>([])
const warehouseId = ref(0)
const keyword = ref('')
const docNo = ref('')

const filteredList = computed(() => {
  let list = h.list.value
  if (warehouseId.value) list = list.filter((r: any) => Number(r.warehouse_id) === warehouseId.value)
  if (keyword.value.trim()) {
    const kw = keyword.value.trim().toLowerCase()
    list = list.filter((r: any) => String(r.sku_name || '').toLowerCase().includes(kw))
  }
  if (docNo.value.trim()) {
    const dn = docNo.value.trim().toLowerCase()
    list = list.filter((r: any) => String(r.doc_no || '').toLowerCase().includes(dn))
  }
  return list
})

async function loadWarehouses() {
  const res: any = await getWarehouses()
  warehouses.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
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
      <text class="title">库存流水</text>
    </view>

    <view class="filters">
      <scroll-view scroll-x class="wh-scroll">
        <view class="wh-tags">
          <view :class="['wh-tag', !warehouseId ? 'active' : '']" @click="selectWarehouse(0)">全部</view>
          <view v-for="w in warehouses" :key="w.id" :class="['wh-tag', warehouseId === w.id ? 'active' : '']" @click="selectWarehouse(w.id)">{{ w.name }}</view>
        </view>
      </scroll-view>
      <up-input v-model="keyword" placeholder="搜索SKU名称" clearable />
      <up-input v-model="docNo" placeholder="搜索单据号" clearable />
    </view>

    <view v-if="!h.loading.value && !filteredList.length" class="empty">暂无流水记录</view>
    <view v-for="item in filteredList" :key="item.id" class="card">
      <view class="row">
        <text class="name">{{ item.sku_name }}</text>
        <text :class="['qty', item.change_qty >= 0 ? 'qty-green' : 'qty-red']">{{ item.change_qty >= 0 ? '+' : '' }}{{ item.change_qty }}</text>
      </view>
      <view class="desc">单号: {{ item.doc_no }} | 类型: {{ item.biz_type === 'inbound' ? '入库' : '出库' }}</view>
      <view class="desc">{{ item.before_qty }} -> {{ item.after_qty }} | {{ item.occurred_at?.slice(0, 16)?.replace('T', ' ') || '-' }}</view>
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
.qty { font-size: 30rpx; font-weight: 700; }
.qty-green { color: #16a34a; }
.qty-red { color: #dc2626; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>
