<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import BatchResultPopup from '@/components/biz/BatchResultPopup.vue'
import { batchShipOrders } from '@/api/order'
import { getStorage, removeStorage } from '@/utils/storage'

const ids = ref<number[]>([])
const rows = reactive<any[]>([])
const company = ref('SF')
const loading = ref(false)
const showResult = ref(false)
const result = ref<{ success_ids: number[]; fail: Array<{ id: number; reason: string }> }>({ success_ids: [], fail: [] })

function syncRows() {
  rows.length = 0
  for (const id of ids.value) rows.push({ order_id: id, company: company.value, tracking_no: '' })
}

function applyCompany(c: string) { company.value = c; for (const r of rows) r.company = c }

async function submit() {
  if (!rows.length) return
  if (rows.some((r) => !String(r.tracking_no).trim())) { uni.showToast({ title: '请填写所有运单号', icon: 'none' }); return }
  loading.value = true
  try {
    const r: any = await batchShipOrders(rows.map((x) => ({ order_id: Number(x.order_id), company: String(x.company || 'SF'), tracking_no: String(x.tracking_no || '').trim() })))
    result.value = r; showResult.value = true
  } finally { loading.value = false }
}

onLoad(() => {
  const raw = String(getStorage('eapp_batch_ship_ids') || '[]')
  try { ids.value = JSON.parse(raw) } catch { ids.value = [] }
  removeStorage('eapp_batch_ship_ids')
  const seed = String(getStorage('eapp_batch_ship_seed') || '')
  if (seed && !ids.value.length) {
    uni.showToast({ title: `已扫码：${seed}`, icon: 'none' })
    removeStorage('eapp_batch_ship_seed')
  }
  syncRows()
})
</script>

<template>
  <view class="page">
    <view class="head">
      <view>选中订单：{{ ids.length }} 单</view>
      <view class="company">
        <text>统一快递：</text>
        <picker mode="selector" :range="['SF','ZTO','YTO','STO','YD','JD','EMS','DBL','JT']" @change="(e) => applyCompany(['SF','ZTO','YTO','STO','YD','JD','EMS','DBL','JT'][Number(e.detail.value)])">
          <view class="picker">{{ company }}</view>
        </picker>
      </view>
    </view>
    <view v-for="row in rows" :key="row.order_id" class="row">
      <text class="no">#{{ row.order_id }}</text>
      <up-input v-model="row.company" placeholder="快递代码" class="col-c" />
      <up-input v-model="row.tracking_no" placeholder="运单号" class="col-t" />
    </view>
    <up-button type="primary" :loading="loading" class="mt" @click="submit">提交批量发货</up-button>
    <BatchResultPopup :show="showResult" :success="result.success_ids" :fails="result.fail" @close="showResult = false; uni.navigateBack()" />
  </view>
</template>

<style scoped>
.page { min-height: 100vh; padding: 20rpx; background: var(--eapp-bg); display: grid; gap: 14rpx; }
.head { display: flex; align-items: center; justify-content: space-between; background: var(--eapp-card); border-radius: 18rpx; padding: 18rpx; border: 1px solid var(--eapp-border); }
.company { display: flex; align-items: center; gap: 10rpx; font-size: 24rpx; }
.picker { border: 1px solid var(--eapp-border); border-radius: 10rpx; padding: 0 16rpx; height: 60rpx; display: flex; align-items: center; }
.row { display: flex; align-items: center; gap: 12rpx; background: var(--eapp-card); padding: 14rpx; border-radius: 18rpx; border: 1px solid var(--eapp-border); }
.no { font-size: 24rpx; color: var(--eapp-text-muted); min-width: 120rpx; }
.col-c { width: 200rpx; }
.col-t { flex: 1; }
.mt { margin-top: 12rpx; }
</style>
