<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { getWarehouses } from '@/api/wms'
import { getDocDetail, createDoc, saveDoc, completeDoc, cancelDoc } from '@/api/wms'

const docId = ref(0)
const loading = ref(false)
const warehouses = ref<any[]>([])

const form = reactive({
  doc_type: 'inbound',
  warehouse_id: 0,
  remark: '',
  status: '',
  items: [] as any[],
})

const typeOptions = [
  { label: '入库', value: 'inbound' },
  { label: '出库', value: 'outbound' },
]

function addItem() {
  form.items.push({ sku_id: '', sku_name: '', qty: 1, unit_cost: 0 })
}

function removeItem(index: number) {
  form.items.splice(index, 1)
}

async function loadWarehouses() {
  const res: any = await getWarehouses()
  warehouses.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
}

async function loadDoc(id: number) {
  loading.value = true
  try {
    const res: any = await getDocDetail(id)
    if (res) {
      form.doc_type = res.doc_type || 'inbound'
      form.warehouse_id = Number(res.warehouse_id || 0)
      form.remark = String(res.remark || '')
      form.status = String(res.status || '')
      form.items = Array.isArray(res.items) ? res.items.map((it: any) => ({
        sku_id: it.sku_id || '', sku_name: it.sku_name || '', qty: Number(it.qty || 0), unit_cost: Number(it.unit_cost || 0),
      })) : []
    }
  } finally {
    loading.value = false
  }
}

async function onSave() {
  if (!form.warehouse_id) { uni.showToast({ title: '请选择仓库', icon: 'none' }); return }
  const payload = {
    doc_type: form.doc_type,
    warehouse_id: form.warehouse_id,
    remark: form.remark,
    items: form.items.map((it: any) => ({
      sku_id: Number(it.sku_id || 0), sku_name: String(it.sku_name || ''), qty: Number(it.qty || 0), unit_cost: Number(it.unit_cost || 0),
    })),
  }
  if (docId.value) await saveDoc(docId.value, payload)
  else {
    const res: any = await createDoc(payload)
    if (res?.id) docId.value = res.id
  }
  uni.showToast({ title: '已保存', icon: 'success' })
}

async function onComplete() {
  if (!docId.value) return
  await completeDoc(docId.value)
  form.status = 'completed'
  uni.showToast({ title: '已完成', icon: 'success' })
}

async function onCancel() {
  if (!docId.value) return
  await cancelDoc(docId.value)
  form.status = 'cancelled'
  uni.showToast({ title: '已取消', icon: 'success' })
}

onLoad(async (query: any) => {
  await loadWarehouses()
  if (query?.id) {
    docId.value = Number(query.id)
    await loadDoc(docId.value)
  } else {
    addItem()
  }
})
</script>

<template>
  <view class="page">
    <view v-if="loading" class="loading">加载中...</view>
    <template v-else>
      <view class="section">
        <view class="label">单据类型</view>
        <view class="type-row">
          <view v-for="opt in typeOptions" :key="opt.value" :class="['type-tag', form.doc_type === opt.value ? 'active' : '']" @click="form.status !== 'completed' && form.status !== 'cancelled' && (form.doc_type = opt.value)">{{ opt.label }}</view>
        </view>
      </view>

      <view class="section">
        <view class="label">仓库</view>
        <view class="type-row">
          <view v-for="w in warehouses" :key="w.id" :class="['type-tag', form.warehouse_id === w.id ? 'active' : '']" @click="form.status !== 'completed' && form.status !== 'cancelled' && (form.warehouse_id = w.id)">{{ w.name }}</view>
        </view>
      </view>

      <view class="section">
        <view class="label">备注</view>
        <up-textarea v-model="form.remark" placeholder="备注（可选）" :disabled="form.status === 'completed' || form.status === 'cancelled'" />
      </view>

      <view class="section">
        <view class="row">
          <view class="label">商品明细</view>
          <up-button v-if="form.status !== 'completed' && form.status !== 'cancelled'" size="mini" type="primary" @click="addItem">添加</up-button>
        </view>
        <view v-for="(it, idx) in form.items" :key="idx" class="item-card">
          <up-input v-model="it.sku_name" placeholder="SKU名称" :disabled="form.status === 'completed' || form.status === 'cancelled'" />
          <view class="item-row">
            <up-input v-model="it.sku_id" placeholder="SKU ID" type="number" :disabled="form.status === 'completed' || form.status === 'cancelled'" />
            <up-input v-model="it.qty" placeholder="数量" type="number" :disabled="form.status === 'completed' || form.status === 'cancelled'" />
            <up-input v-model="it.unit_cost" placeholder="单价" type="digit" :disabled="form.status === 'completed' || form.status === 'cancelled'" />
          </view>
          <view v-if="form.status !== 'completed' && form.status !== 'cancelled'" class="item-remove" @click="removeItem(idx)">
            <text class="remove-text">删除</text>
          </view>
        </view>
      </view>

      <view class="actions">
        <up-button v-if="form.status !== 'completed' && form.status !== 'cancelled'" type="primary" @click="onSave">保存草稿</up-button>
        <up-button v-if="docId && form.status === 'draft'" type="success" @click="onComplete">完成</up-button>
        <up-button v-if="docId && form.status === 'draft'" type="error" plain @click="onCancel">取消</up-button>
      </view>
    </template>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 16rpx; align-content: start; }
.section { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.label { font-size: 26rpx; font-weight: 600; margin-bottom: 10rpx; }
.type-row { display: flex; gap: 10rpx; flex-wrap: wrap; }
.type-tag { padding: 8rpx 20rpx; font-size: 24rpx; border-radius: 999rpx; background: #f1f5f9; border: 1px solid var(--eapp-border); }
.type-tag.active { background: var(--eapp-primary, #2563eb); color: #fff; border-color: var(--eapp-primary, #2563eb); }
.row { display: flex; align-items: center; justify-content: space-between; }
.item-card { background: #f8fafc; border: 1px solid var(--eapp-border); border-radius: 12rpx; padding: 14rpx; margin-top: 10rpx; display: grid; gap: 8rpx; }
.item-row { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 8rpx; }
.item-remove { text-align: right; }
.remove-text { font-size: 22rpx; color: #dc2626; }
.actions { display: grid; gap: 10rpx; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
</style>
