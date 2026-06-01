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

const showStartPicker = ref(false)
const showEndPicker = ref(false)
const showActivityPicker = ref(false)
const showStatusPicker = ref(false)
const showProductPicker = ref(false)
const showSkuPicker = ref(false)

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

const priceHeader = computed(() => (props.kind === 'bargain' ? '起砍价 / 底价' : '活动价'))
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

    <!-- ── Header card ─────────────────────────────── -->
    <view class="header-card">
      <!-- Title row -->
      <view class="header-top">
        <view class="page-title">{{ title }}</view>
        <view class="btn-new-act" @click="openCreateActivity">新增活动</view>
      </view>

      <!-- Search row -->
      <view class="search-row">
        <view class="search-input-wrap">
          <up-input
            v-model="keyword"
            placeholder="搜索商品关键词"
            clearable
            :custom-style="{ height: '72rpx', fontSize: '26rpx' }"
          />
        </view>
        <view class="btn-search" @click="loadRows">搜索</view>
      </view>

      <!-- Activity tabs -->
      <view v-if="activities.length">
        <scroll-view scroll-x class="tabs-scroll">
          <view class="tabs-row">
            <view
              v-for="act in activities"
              :key="act.id"
              class="tab-item"
              :class="{ 'tab-active': Number(act.id) === selectedActivityID }"
              @click="selectedActivityID = Number(act.id); loadRows()"
            >
              <view class="tab-name">{{ act.name }}</view>
              <view class="tab-meta">
                <view class="status-dot" :class="Number(act.status) === 1 ? 'dot-on' : 'dot-off'" />
                <text class="tab-date">{{ formatTime(act.start_at).slice(0, 10) }}</text>
                <text class="tab-edit-link" @click.stop="openEditActivity(act)">编辑</text>
              </view>
            </view>
          </view>
        </scroll-view>
      </view>
      <view v-else class="no-activity-hint">暂无活动，请先新增活动</view>
    </view>

    <!-- ── Product list ─────────────────────────────── -->
    <view class="list-area">
      <!-- Loading -->
      <view v-if="loading" class="state-box">
        <view class="spinner" />
        <text class="state-text">加载中...</text>
      </view>

      <!-- Empty -->
      <view v-else-if="!rows.length" class="state-box">
        <view class="empty-icon-wrap">
          <view class="empty-icon" />
        </view>
        <text class="state-text">暂无活动商品</text>
        <text class="state-sub">{{ canEditRows ? '点击底部按钮新增商品' : '请先选择活动' }}</text>
      </view>

      <!-- Product cards -->
      <view v-else class="row-list">
        <view v-for="item in rows" :key="item.id" class="row-card">
          <!-- Cover image -->
          <image
            v-if="item.product_cover"
            class="product-cover"
            :src="item.product_cover"
            mode="aspectFill"
          />
          <view v-else class="product-cover cover-placeholder" />

          <!-- Info -->
          <view class="card-body">
            <text class="card-title">{{ item.product_title || `商品#${item.product_id}` }}</text>
            <text class="card-sku">{{ parseSkuLabel(item.sku_attrs, Number(item.sku_id || 0)) }}</text>

            <!-- Price -->
            <view class="price-row">
              <template v-if="props.kind !== 'bargain'">
                <view class="price-block">
                  <text class="price-label">活动价</text>
                  <text class="price-val primary">¥{{ Number(item.activity_price || 0).toFixed(2) }}</text>
                </view>
                <view v-if="item.sku_price" class="price-block">
                  <text class="price-label">原价</text>
                  <text class="price-val muted line-through">¥{{ Number(item.sku_price || 0).toFixed(2) }}</text>
                </view>
              </template>
              <template v-else>
                <view class="price-block">
                  <text class="price-label">起砍价</text>
                  <text class="price-val primary">¥{{ Number(item.start_price || 0).toFixed(2) }}</text>
                </view>
                <view class="price-block">
                  <text class="price-label">底价</text>
                  <text class="price-val accent">¥{{ Number(item.floor_price || 0).toFixed(2) }}</text>
                </view>
              </template>
            </view>

            <!-- Stats row -->
            <view class="stats-row">
              <view class="stat-chip">限购 {{ item.limit_per_order || '不限' }}</view>
              <view class="stat-chip">库存 {{ item.total_stock_limit || '不限' }}</view>
              <view class="stat-chip sold">已售 {{ item.sold_qty || 0 }}</view>
            </view>

            <!-- Action buttons -->
            <view class="card-btns">
              <view class="card-btn edit" @click="openEditRow(item)">编辑</view>
              <view class="card-btn del" @click="removeRow(item.id)">删除</view>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- ── Sticky bottom bar ────────────────────────── -->
    <view class="sticky-bar">
      <view
        class="btn-add-product"
        :class="{ 'btn-disabled': !canEditRows }"
        @click="openCreateRow"
      >新增活动商品</view>
    </view>

    <!-- ── Activity popup ──────────────────────────── -->
    <up-popup :show="showActivityPopup" mode="bottom" round="16" @close="showActivityPopup = false">
      <view class="popup-body">
        <view class="popup-header">
          <text class="popup-title">{{ editingActivityID ? '编辑活动' : '新增活动' }}</text>
          <view class="popup-close" @click="showActivityPopup = false">✕</view>
        </view>

        <view class="form-label">活动名称</view>
        <up-input v-model="activityForm.name" placeholder="请输入活动名称" clearable />

        <view class="form-label mt">开始时间</view>
        <view class="picker-row" @click="showStartPicker = true">
          <text :class="activityForm.start_at ? 'picker-val' : 'picker-placeholder'">
            {{ activityForm.start_at ? activityForm.start_at.slice(0, 16).replace('T', ' ') : '请选择开始时间' }}
          </text>
          <text class="picker-arrow">›</text>
        </view>
        <up-datetime-picker
          :show="showStartPicker"
          v-model="activityForm.start_at"
          mode="datetime"
          @confirm="(e: any) => { activityForm.start_at = new Date(e.value).toISOString(); showStartPicker = false }"
          @cancel="showStartPicker = false"
          @close="showStartPicker = false"
        />

        <view class="form-label mt">结束时间</view>
        <view class="picker-row" @click="showEndPicker = true">
          <text :class="activityForm.end_at ? 'picker-val' : 'picker-placeholder'">
            {{ activityForm.end_at ? activityForm.end_at.slice(0, 16).replace('T', ' ') : '请选择结束时间' }}
          </text>
          <text class="picker-arrow">›</text>
        </view>
        <up-datetime-picker
          :show="showEndPicker"
          v-model="activityForm.end_at"
          mode="datetime"
          @confirm="(e: any) => { activityForm.end_at = new Date(e.value).toISOString(); showEndPicker = false }"
          @cancel="showEndPicker = false"
          @close="showEndPicker = false"
        />

        <view class="form-label mt">状态</view>
        <view class="picker-row" @click="showStatusPicker = true">
          <text class="picker-val">{{ Number(activityForm.status || 0) === 1 ? '启用' : '禁用' }}</text>
          <text class="picker-arrow">›</text>
        </view>
        <up-picker
          :show="showStatusPicker"
          :columns="[[{ label: '启用', value: 1 }, { label: '禁用', value: 0 }]]"
          keyName="label"
          @confirm="(e: any) => { activityForm.status = e.value[0].value; showStatusPicker = false }"
          @cancel="showStatusPicker = false"
          @close="showStatusPicker = false"
        />

        <view class="mt-btn">
          <up-button type="primary" @click="saveActivity">保存活动</up-button>
        </view>
      </view>
    </up-popup>

    <!-- ── Row popup ───────────────────────────────── -->
    <up-popup :show="showRowPopup" mode="bottom" round="16" @close="showRowPopup = false">
      <view class="popup-body">
        <view class="popup-header">
          <text class="popup-title">{{ editingRowID ? '编辑活动商品' : '新增活动商品' }}</text>
          <view class="popup-close" @click="showRowPopup = false">✕</view>
        </view>

        <view class="form-label">选择商品</view>
        <view class="picker-row" @click="showProductPicker = true">
          <text :class="rowForm.product_id ? 'picker-val' : 'picker-placeholder'">
            {{ productOptions.find((x) => Number(x.id) === Number(rowForm.product_id || 0))?.title || '请选择商品' }}
          </text>
          <text class="picker-arrow">›</text>
        </view>
        <up-picker
          :show="showProductPicker"
          :columns="[productOptions]"
          keyName="title"
          @confirm="async (e: any) => { rowForm.product_id = Number(e.value[0]?.id || 0); showProductPicker = false; await onSelectProduct() }"
          @cancel="showProductPicker = false"
          @close="showProductPicker = false"
        />

        <view class="form-label mt">选择 SKU</view>
        <view class="picker-row" @click="showSkuPicker = true">
          <text :class="rowForm.sku_id ? 'picker-val' : 'picker-placeholder'">
            {{ skuOptions.find((x) => Number(x.id) === Number(rowForm.sku_id || 0))?.label || '请选择SKU' }}
          </text>
          <text class="picker-arrow">›</text>
        </view>
        <up-picker
          :show="showSkuPicker"
          :columns="[skuOptions]"
          keyName="label"
          @confirm="(e: any) => { rowForm.sku_id = Number(e.value[0]?.id || 0); showSkuPicker = false }"
          @cancel="showSkuPicker = false"
          @close="showSkuPicker = false"
        />

        <template v-if="props.kind !== 'bargain'">
          <view class="form-label mt">活动价（元）</view>
          <up-input v-model="rowForm.activity_price" type="digit" inputmode="decimal" placeholder="请输入活动价" />
        </template>
        <template v-else>
          <view class="form-label mt">起砍价（元）</view>
          <up-input v-model="rowForm.start_price" type="digit" inputmode="decimal" placeholder="请输入起砍价" />
          <view class="form-label mt">底价（元）</view>
          <up-input v-model="rowForm.floor_price" type="digit" inputmode="decimal" placeholder="请输入底价" />
        </template>

        <view class="form-row-2 mt">
          <view class="form-col">
            <view class="form-label">限购数量</view>
            <up-input v-model="rowForm.limit_per_order" type="number" inputmode="numeric" placeholder="0为不限" />
          </view>
          <view class="form-col">
            <view class="form-label">库存上限</view>
            <up-input v-model="rowForm.total_stock_limit" type="number" inputmode="numeric" placeholder="0为不限" />
          </view>
        </view>

        <view class="mt-btn">
          <up-button type="primary" @click="saveRow">保存活动商品</up-button>
        </view>
      </view>
    </up-popup>

  </view>
</template>

<style scoped>
/* ── Page ───────────────────────────────────────────── */
.page {
  min-height: 100vh;
  background: var(--eapp-bg, #f4f6f9);
  padding: 20rpx 20rpx 160rpx;
  box-sizing: border-box;
}

/* ── Header card ────────────────────────────────────── */
.header-card {
  background: #fff;
  border-radius: 24rpx;
  padding: 24rpx 24rpx 16rpx;
  box-shadow: 0 2rpx 16rpx rgba(0, 0, 0, 0.06);
}

.header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18rpx;
}

.page-title {
  font-size: 34rpx;
  font-weight: 700;
  color: var(--eapp-text, #1a1a2e);
  letter-spacing: -0.5rpx;
}

.btn-new-act {
  height: 60rpx;
  padding: 0 24rpx;
  background: var(--eapp-primary, #3b82f6);
  color: #fff;
  border-radius: 30rpx;
  font-size: 24rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: opacity 0.2s;
}
.btn-new-act:active { opacity: 0.8; }

/* Search row */
.search-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 16rpx;
}

.search-input-wrap {
  flex: 1;
}

.btn-search {
  height: 72rpx;
  min-width: 100rpx;
  padding: 0 20rpx;
  background: var(--eapp-primary, #3b82f6);
  color: #fff;
  border-radius: 14rpx;
  font-size: 26rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: opacity 0.2s;
}
.btn-search:active { opacity: 0.8; }

/* Activity tabs */
.tabs-scroll {
  white-space: nowrap;
  overflow: hidden;
  margin: 0 -4rpx;
}

.tabs-row {
  display: flex;
  gap: 12rpx;
  padding: 4rpx 4rpx 6rpx;
}

.tab-item {
  display: inline-flex;
  flex-direction: column;
  min-width: 240rpx;
  max-width: 300rpx;
  padding: 14rpx 18rpx;
  border: 1.5rpx solid var(--eapp-border, #e8edf5);
  border-radius: 18rpx;
  background: #f8fafc;
  cursor: pointer;
  transition: all 0.2s;
  box-sizing: border-box;
}

.tab-item:active { opacity: 0.85; }

.tab-active {
  border-color: var(--eapp-primary, #3b82f6);
  background: #eff6ff;
}

.tab-name {
  font-size: 26rpx;
  font-weight: 600;
  color: var(--eapp-text, #1a1a2e);
  margin-bottom: 8rpx;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tab-active .tab-name {
  color: var(--eapp-primary, #3b82f6);
}

.tab-meta {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.status-dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
  flex-shrink: 0;
}
.dot-on  { background: #22c55e; }
.dot-off { background: #cbd5e1; }

.tab-date {
  font-size: 20rpx;
  color: var(--eapp-text-muted, #94a3b8);
  flex: 1;
}

.tab-edit-link {
  font-size: 22rpx;
  color: var(--eapp-primary, #3b82f6);
  cursor: pointer;
  padding: 4rpx 0;
}

.no-activity-hint {
  font-size: 26rpx;
  color: var(--eapp-text-muted, #94a3b8);
  text-align: center;
  padding: 24rpx 0 8rpx;
}

/* ── List area ──────────────────────────────────────── */
.list-area {
  margin-top: 20rpx;
}

/* State box (loading / empty) */
.state-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 100rpx 40rpx 60rpx;
}

.spinner {
  width: 60rpx;
  height: 60rpx;
  border: 4rpx solid #e2e8f0;
  border-top-color: var(--eapp-primary, #3b82f6);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 20rpx;
}

@keyframes spin { to { transform: rotate(360deg); } }

.empty-icon-wrap {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background: #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20rpx;
}

.empty-icon {
  width: 60rpx;
  height: 60rpx;
  border: 4rpx solid #cbd5e1;
  border-radius: 10rpx;
  position: relative;
}
.empty-icon::before {
  content: '';
  position: absolute;
  top: 12rpx;
  left: 10rpx;
  right: 10rpx;
  height: 4rpx;
  background: #cbd5e1;
  border-radius: 2rpx;
  box-shadow: 0 14rpx 0 #cbd5e1, 0 28rpx 0 #cbd5e1;
}

.state-text {
  font-size: 28rpx;
  color: var(--eapp-text, #1a1a2e);
  font-weight: 600;
  margin-bottom: 8rpx;
}

.state-sub {
  font-size: 24rpx;
  color: var(--eapp-text-muted, #94a3b8);
}

/* Row list */
.row-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.row-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 20rpx;
  display: flex;
  gap: 18rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
  transition: box-shadow 0.2s;
}

.row-card:active {
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);
}

.product-cover {
  width: 120rpx;
  height: 120rpx;
  border-radius: 14rpx;
  flex-shrink: 0;
  object-fit: cover;
}

.cover-placeholder {
  background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%);
}

.card-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.card-title {
  font-size: 28rpx;
  font-weight: 600;
  color: var(--eapp-text, #1a1a2e);
  line-height: 1.4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-sku {
  font-size: 22rpx;
  color: var(--eapp-text-muted, #94a3b8);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.price-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  flex-wrap: wrap;
  margin-top: 2rpx;
}

.price-block {
  display: flex;
  align-items: baseline;
  gap: 6rpx;
}

.price-label {
  font-size: 20rpx;
  color: var(--eapp-text-muted, #94a3b8);
}

.price-val {
  font-size: 28rpx;
  font-weight: 700;
}

.price-val.primary { color: var(--eapp-primary, #3b82f6); }
.price-val.accent  { color: #f97316; }
.price-val.muted   { color: var(--eapp-text-muted, #94a3b8); font-size: 22rpx; font-weight: 400; }
.line-through      { text-decoration: line-through; }

.stats-row {
  display: flex;
  gap: 8rpx;
  flex-wrap: wrap;
  margin-top: 4rpx;
}

.stat-chip {
  font-size: 20rpx;
  color: #475569;
  background: #f1f5f9;
  border-radius: 8rpx;
  padding: 4rpx 10rpx;
}

.stat-chip.sold { background: #fef3c7; color: #92400e; }

.card-btns {
  display: flex;
  gap: 10rpx;
  margin-top: 8rpx;
}

.card-btn {
  height: 56rpx;
  padding: 0 22rpx;
  border-radius: 12rpx;
  font-size: 24rpx;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: opacity 0.2s;
}
.card-btn:active { opacity: 0.75; }

.card-btn.edit {
  border: 1.5rpx solid var(--eapp-primary, #3b82f6);
  color: var(--eapp-primary, #3b82f6);
  background: #eff6ff;
}

.card-btn.del {
  border: 1.5rpx solid #fca5a5;
  color: #ef4444;
  background: #fff5f5;
}

/* ── Sticky bar ─────────────────────────────────────── */
.sticky-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16rpx 30rpx calc(16rpx + env(safe-area-inset-bottom));
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-top: 1rpx solid var(--eapp-border, #e8edf5);
  z-index: 100;
}

.btn-add-product {
  height: 88rpx;
  background: var(--eapp-primary, #3b82f6);
  color: #fff;
  border-radius: 22rpx;
  font-size: 30rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: opacity 0.2s;
  letter-spacing: 0.5rpx;
}
.btn-add-product:active { opacity: 0.85; }

.btn-add-product.btn-disabled {
  background: #cbd5e1;
  cursor: not-allowed;
}

/* ── Popups ─────────────────────────────────────────── */
.popup-body {
  padding: 28rpx 28rpx calc(28rpx + env(safe-area-inset-bottom));
  box-sizing: border-box;
}

.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24rpx;
}

.popup-title {
  font-size: 32rpx;
  font-weight: 700;
  color: var(--eapp-text, #1a1a2e);
}

.popup-close {
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: #f1f5f9;
  font-size: 26rpx;
  color: #64748b;
  cursor: pointer;
}
.popup-close:active { background: #e2e8f0; }

.form-label {
  font-size: 24rpx;
  font-weight: 600;
  color: #475569;
  margin-bottom: 10rpx;
}
.form-label.mt { margin-top: 20rpx; }

.picker-row {
  min-height: 80rpx;
  border: 1.5rpx solid var(--eapp-border, #e8edf5);
  border-radius: 14rpx;
  padding: 0 20rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #f8fafc;
  cursor: pointer;
  transition: border-color 0.2s;
}
.picker-row:active { border-color: var(--eapp-primary, #3b82f6); }

.picker-val {
  font-size: 26rpx;
  color: var(--eapp-text, #1a1a2e);
}

.picker-placeholder {
  font-size: 26rpx;
  color: #94a3b8;
}

.picker-arrow {
  font-size: 32rpx;
  color: #cbd5e1;
  line-height: 1;
}

.form-row-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16rpx;
}

.form-row-2.mt { margin-top: 20rpx; }

.mt-btn {
  margin-top: 30rpx;
}
</style>
