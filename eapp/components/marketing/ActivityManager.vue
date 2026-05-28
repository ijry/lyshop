<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import {
  createActivity,
  getActivities,
  getActivityProducts,
  type MarketingKind,
  updateActivity,
  upsertActivityProducts,
} from '@/api/marketing'
import { getProductDetail, getProducts } from '@/api/product'
import StatusTag from '@/components/common/StatusTag.vue'

const props = defineProps<{ kind: MarketingKind }>()

const loading = ref(false)
const activities = ref<any[]>([])
const rows = ref<any[]>([])
const selectedActivityID = ref(0)
const keyword = ref('')
const filterProductID = ref(0)

const productOptions = ref<any[]>([])
const skuOptions = ref<Array<{ id: number; label: string }>>([])

const showActivityPopup = ref(false)
const editingActivityID = ref(0)
const activityForm = reactive({
  name: '',
  start_at: '',
  end_at: '',
  status: 1,
})

const showRowPopup = ref(false)
const editingRowID = ref(0)
const rowForm = reactive({
  product_id: 0,
  sku_id: 0,
  activity_price: 0,
  start_price: 0,
  floor_price: 0,
  limit_per_order: 0,
  total_stock_limit: 0,
})

const title = computed(() => {
  if (props.kind === 'seckill') return '秒杀活动'
  if (props.kind === 'group-buy') return '拼团活动'
  return '砍价活动'
})

const priceHeader = computed(() => (props.kind === 'bargain' ? '起砍价/底价' : '活动价'))
const canEditRows = computed(() => selectedActivityID.value > 0)

function formatTime(v?: string) {
  if (!v) return '-'
  return String(v).slice(0, 19).replace('T', ' ')
}

function parseSkuLabel(attrs: any, skuID: number) {
  try {
    const list = typeof attrs === 'string' ? JSON.parse(attrs) : attrs
    if (Array.isArray(list) && list.length) {
      return list.map((item: any) => `${item.name}:${item.value}`).join(' / ')
    }
  } catch {
    // ignore
  }
  return `SKU#${skuID}`
}

function payloadRows(sourceRows: any[]) {
  return sourceRows.map((item: any) => ({
    product_id: Number(item.product_id || 0),
    sku_id: Number(item.sku_id || 0),
    activity_price: Number(item.activity_price || 0),
    start_price: Number(item.start_price || 0),
    floor_price: Number(item.floor_price || 0),
    limit_per_order: Number(item.limit_per_order || 0),
    total_stock_limit: Number(item.total_stock_limit || 0),
  }))
}

async function loadActivities() {
  const data: any = await getActivities(props.kind, { page: 1, size: 200 })
  activities.value = Array.isArray(data?.list) ? data.list : []
  if (!selectedActivityID.value && activities.value.length) {
    selectedActivityID.value = Number(activities.value[0].id || 0)
  }
}

async function loadProductOptions() {
  const data: any = await getProducts({ page: 1, size: 200 })
  productOptions.value = Array.isArray(data?.list) ? data.list : []
}

async function loadRows() {
  if (!selectedActivityID.value) {
    rows.value = []
    return
  }
  loading.value = true
  try {
    const data: any = await getActivityProducts(props.kind, {
      activity_id: selectedActivityID.value,
      product_id: filterProductID.value || undefined,
      keyword: keyword.value || undefined,
      page: 1,
      size: 200,
    })
    rows.value = Array.isArray(data?.list) ? data.list : []
  } finally {
    loading.value = false
  }
}

async function bootstrap() {
  await Promise.all([loadActivities(), loadProductOptions()])
  await loadRows()
}

function resetActivityForm() {
  activityForm.name = ''
  activityForm.start_at = ''
  activityForm.end_at = ''
  activityForm.status = 1
}

function openCreateActivity() {
  editingActivityID.value = 0
  resetActivityForm()
  showActivityPopup.value = true
}

function openEditActivity(item: any) {
  editingActivityID.value = Number(item?.id || 0)
  activityForm.name = String(item?.name || '')
  activityForm.start_at = String(item?.start_at || '')
  activityForm.end_at = String(item?.end_at || '')
  activityForm.status = Number(item?.status || 0) === 1 ? 1 : 0
  showActivityPopup.value = true
}

async function saveActivity() {
  if (!activityForm.name.trim()) {
    uni.showToast({ title: '请输入活动名称', icon: 'none' })
    return
  }
  if (!activityForm.start_at || !activityForm.end_at) {
    uni.showToast({ title: '请填写开始和结束时间', icon: 'none' })
    return
  }
  const payload = {
    name: activityForm.name.trim(),
    start_at: activityForm.start_at,
    end_at: activityForm.end_at,
    status: Number(activityForm.status || 0) === 1 ? 1 : 0,
  }
  if (editingActivityID.value > 0) {
    await updateActivity(props.kind, editingActivityID.value, payload)
  } else {
    await createActivity(props.kind, payload)
  }
  showActivityPopup.value = false
  await loadActivities()
  if (!selectedActivityID.value && activities.value.length) {
    selectedActivityID.value = Number(activities.value[0].id || 0)
  }
  await loadRows()
}

function resetRowForm() {
  rowForm.product_id = 0
  rowForm.sku_id = 0
  rowForm.activity_price = 0
  rowForm.start_price = 0
  rowForm.floor_price = 0
  rowForm.limit_per_order = 0
  rowForm.total_stock_limit = 0
  skuOptions.value = []
}

function openCreateRow() {
  if (!selectedActivityID.value) {
    uni.showToast({ title: '请先选择活动', icon: 'none' })
    return
  }
  editingRowID.value = 0
  resetRowForm()
  showRowPopup.value = true
}

async function openEditRow(item: any) {
  editingRowID.value = Number(item?.id || 0)
  rowForm.product_id = Number(item?.product_id || 0)
  rowForm.sku_id = Number(item?.sku_id || 0)
  rowForm.activity_price = Number(item?.activity_price || 0)
  rowForm.start_price = Number(item?.start_price || 0)
  rowForm.floor_price = Number(item?.floor_price || 0)
  rowForm.limit_per_order = Number(item?.limit_per_order || 0)
  rowForm.total_stock_limit = Number(item?.total_stock_limit || 0)
  skuOptions.value = [{ id: rowForm.sku_id, label: parseSkuLabel(item?.sku_attrs, rowForm.sku_id) }]
  showRowPopup.value = true
}

async function onSelectProduct() {
  rowForm.sku_id = 0
  skuOptions.value = []
  if (!rowForm.product_id) return
  const detail: any = await getProductDetail(rowForm.product_id)
  const skus = Array.isArray(detail?.skus) ? detail.skus : []
  skuOptions.value = skus.map((item: any) => ({
    id: Number(item?.id || 0),
    label: `${parseSkuLabel(item?.attrs, Number(item?.id || 0))} ｜ ¥${Number(item?.price || 0)} ｜ 库存 ${Number(item?.stock || 0)}`,
  })).filter((item: any) => item.id > 0)
}

function validateRowForm() {
  if (!rowForm.product_id || !rowForm.sku_id) return '请选择商品和 SKU'
  if (props.kind === 'bargain') {
    if (rowForm.start_price <= 0 || rowForm.floor_price <= 0 || rowForm.floor_price > rowForm.start_price) {
      return '请输入有效的起砍价和底价'
    }
  } else if (rowForm.activity_price <= 0) {
    return '请输入有效活动价'
  }
  if (rowForm.limit_per_order < 0 || rowForm.total_stock_limit < 0) {
    return '限购与库存不能小于 0'
  }
  return ''
}

async function saveRow() {
  const msg = validateRowForm()
  if (msg) {
    uni.showToast({ title: msg, icon: 'none' })
    return
  }
  const duplicate = rows.value.find((item: any) =>
    Number(item.id || 0) !== editingRowID.value &&
    Number(item.product_id || 0) === Number(rowForm.product_id || 0) &&
    Number(item.sku_id || 0) === Number(rowForm.sku_id || 0),
  )
  if (duplicate) {
    uni.showToast({ title: '同一活动下 SKU 不能重复', icon: 'none' })
    return
  }
  const nextRows = rows.value.filter((item: any) => Number(item.id || 0) !== editingRowID.value)
  nextRows.push({
    product_id: rowForm.product_id,
    sku_id: rowForm.sku_id,
    activity_price: rowForm.activity_price,
    start_price: rowForm.start_price,
    floor_price: rowForm.floor_price,
    limit_per_order: rowForm.limit_per_order,
    total_stock_limit: rowForm.total_stock_limit,
  })
  await upsertActivityProducts(props.kind, selectedActivityID.value, payloadRows(nextRows))
  showRowPopup.value = false
  await loadRows()
}

async function removeRow(rowID: number) {
  const nextRows = rows.value.filter((item: any) => Number(item.id || 0) !== Number(rowID || 0))
  await upsertActivityProducts(props.kind, selectedActivityID.value, payloadRows(nextRows))
  await loadRows()
}

onLoad(bootstrap)
onShow(loadRows)
</script>

<template>
  <view class="page">
    <view class="top-card">
      <view class="title-wrap">
        <view class="title">{{ title }}</view>
        <up-button size="mini" type="primary" plain @click="openCreateActivity">新增活动</up-button>
      </view>
      <view class="controls">
        <picker
          mode="selector"
          :range="activities"
          range-key="name"
          @change="(e) => { selectedActivityID = Number(activities[e.detail.value]?.id || 0); loadRows() }"
        >
          <view class="picker">{{ activities.find((x) => Number(x.id) === selectedActivityID)?.name || '请选择活动' }}</view>
        </picker>
        <up-input v-model="keyword" placeholder="关键词搜索" clearable />
        <up-button size="mini" type="primary" @click="loadRows">搜索</up-button>
      </view>
      <scroll-view scroll-x class="activity-scroll">
        <view class="activity-list">
          <view
            v-for="act in activities"
            :key="act.id"
            class="activity-item"
            :class="{ 'active-item': Number(act.id) === selectedActivityID }"
            @click="selectedActivityID = Number(act.id); loadRows()"
          >
            <view class="name">{{ act.name }}</view>
            <view class="time">{{ formatTime(act.start_at) }} ~ {{ formatTime(act.end_at) }}</view>
            <view class="row-end">
              <StatusTag :text="Number(act.status || 0) === 1 ? '启用' : '禁用'" :type="Number(act.status || 0) === 1 ? 'enabled' : 'disabled'" />
              <text class="edit-link" @click.stop="openEditActivity(act)">编辑</text>
            </view>
          </view>
        </view>
      </scroll-view>
    </view>

    <view class="op-row">
      <up-button size="mini" type="primary" :disabled="!canEditRows" @click="openCreateRow">新增活动商品</up-button>
    </view>

    <view v-if="loading" class="empty">加载中...</view>
    <view v-else-if="!rows.length" class="empty">暂无活动商品</view>
    <view v-else class="row-list">
      <view v-for="item in rows" :key="item.id" class="row-card">
        <view class="name">{{ item.product_title || `商品#${item.product_id}` }}</view>
        <view class="sub">SKU：{{ parseSkuLabel(item.sku_attrs, Number(item.sku_id || 0)) }}</view>
        <view class="sub" v-if="props.kind !== 'bargain'">{{ priceHeader }}：¥{{ Number(item.activity_price || 0).toFixed(2) }}</view>
        <view class="sub" v-else>{{ priceHeader }}：¥{{ Number(item.start_price || 0).toFixed(2) }} / ¥{{ Number(item.floor_price || 0).toFixed(2) }}</view>
        <view class="sub">限购：{{ item.limit_per_order || 0 }} · 库存上限：{{ item.total_stock_limit || 0 }} · 已售：{{ item.sold_qty || 0 }}</view>
        <view class="btns">
          <up-button size="mini" type="primary" plain @click="openEditRow(item)">编辑</up-button>
          <up-button size="mini" type="error" plain @click="removeRow(item.id)">删除</up-button>
        </view>
      </view>
    </view>

    <up-popup :show="showActivityPopup" mode="bottom" round="16" @close="showActivityPopup = false">
      <view class="popup-body">
        <view class="popup-title">{{ editingActivityID ? '编辑活动' : '新增活动' }}</view>
        <up-input v-model="activityForm.name" placeholder="活动名称" clearable />
        <view class="mt-12rpx" />
        <up-input v-model="activityForm.start_at" placeholder="开始时间（ISO 或 2026-05-28T10:00:00Z）" clearable />
        <view class="mt-12rpx" />
        <up-input v-model="activityForm.end_at" placeholder="结束时间（ISO 或 2026-05-30T10:00:00Z）" clearable />
        <view class="mt-12rpx" />
        <picker mode="selector" :range="[{ label: '启用', value: 1 }, { label: '禁用', value: 0 }]" range-key="label" @change="(e) => { activityForm.status = Number(([1,0][e.detail.value]) || 0) }">
          <view class="picker">{{ Number(activityForm.status || 0) === 1 ? '启用' : '禁用' }}</view>
        </picker>
        <view class="mt-16rpx" />
        <up-button type="primary" @click="saveActivity">保存活动</up-button>
      </view>
    </up-popup>

    <up-popup :show="showRowPopup" mode="bottom" round="16" @close="showRowPopup = false">
      <view class="popup-body">
        <view class="popup-title">{{ editingRowID ? '编辑活动商品' : '新增活动商品' }}</view>
        <picker mode="selector" :range="productOptions" range-key="title" @change="async (e) => { rowForm.product_id = Number(productOptions[e.detail.value]?.id || 0); await onSelectProduct() }">
          <view class="picker">{{ productOptions.find((x) => Number(x.id) === Number(rowForm.product_id || 0))?.title || '请选择商品' }}</view>
        </picker>
        <view class="mt-12rpx" />
        <picker mode="selector" :range="skuOptions" range-key="label" @change="(e) => { rowForm.sku_id = Number(skuOptions[e.detail.value]?.id || 0) }">
          <view class="picker">{{ skuOptions.find((x) => Number(x.id) === Number(rowForm.sku_id || 0))?.label || '请选择SKU' }}</view>
        </picker>
        <view class="mt-12rpx" />
        <up-input v-if="props.kind !== 'bargain'" v-model="rowForm.activity_price" type="digit" inputmode="decimal" placeholder="活动价" />
        <template v-else>
          <up-input v-model="rowForm.start_price" type="digit" inputmode="decimal" placeholder="起砍价" />
          <view class="mt-12rpx" />
          <up-input v-model="rowForm.floor_price" type="digit" inputmode="decimal" placeholder="底价" />
        </template>
        <view class="mt-12rpx" />
        <up-input v-model="rowForm.limit_per_order" type="number" inputmode="numeric" placeholder="限购数量（0为不限）" />
        <view class="mt-12rpx" />
        <up-input v-model="rowForm.total_stock_limit" type="number" inputmode="numeric" placeholder="库存上限（0为不限）" />
        <view class="mt-16rpx" />
        <up-button type="primary" @click="saveRow">保存活动商品</up-button>
      </view>
    </up-popup>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; }
.top-card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 24rpx; padding: 20rpx; }
.title-wrap { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12rpx; }
.title { font-size: 32rpx; font-weight: 700; }
.controls { display: grid; grid-template-columns: 1fr 1fr auto; gap: 10rpx; align-items: center; }
.activity-scroll { margin-top: 12rpx; white-space: nowrap; }
.activity-list { display: flex; gap: 12rpx; }
.activity-item { width: 420rpx; border: 1px solid var(--eapp-border); border-radius: 16rpx; padding: 14rpx; box-sizing: border-box; }
.active-item { border-color: var(--eapp-primary); background: #eff6ff; }
.name { font-size: 26rpx; font-weight: 600; }
.time { margin-top: 6rpx; color: var(--eapp-text-muted); font-size: 22rpx; }
.row-end { margin-top: 8rpx; display: flex; align-items: center; justify-content: space-between; }
.edit-link { color: var(--eapp-primary); font-size: 22rpx; }
.op-row { margin-top: 14rpx; margin-bottom: 12rpx; }
.row-list { display: grid; gap: 12rpx; }
.row-card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 18rpx; }
.sub { margin-top: 6rpx; color: var(--eapp-text-muted); font-size: 23rpx; }
.btns { margin-top: 10rpx; display: flex; gap: 10rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.popup-body { padding: 24rpx; box-sizing: border-box; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.picker { min-height: 76rpx; border: 1px solid var(--eapp-border); border-radius: 12rpx; padding: 0 20rpx; display: flex; align-items: center; color: var(--eapp-text); }
</style>
