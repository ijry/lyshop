<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { useVipList } from '@/composables/useVipList'
import { getProducts } from '@/api/product'

const h = useVipList('sku-prices')
const showPopup = ref(false)
const editingID = ref(0)

const filterProductID = ref(0)
const filterLevelID = ref(0)
const productOptions = ref<any[]>([])
const skuOptions = ref<Array<{ id: number; label: string }>>([])
const levelOptions = ref<any[]>([])

const form = reactive({ product_id: 0, sku_id: 0, level_id: 0, vip_price: '', status: 1 })

function resetForm() { form.product_id = 0; form.sku_id = 0; form.level_id = 0; form.vip_price = ''; form.status = 1; skuOptions.value = [] }

async function loadFiltered() {
  const params: any = {}
  if (filterProductID.value) params.product_id = filterProductID.value
  if (filterLevelID.value) params.level_id = filterLevelID.value
  await h.load(params)
}

async function loadProductOptions() {
  const data: any = await getProducts({ page: 1, size: 200 })
  productOptions.value = Array.isArray(data?.list) ? data.list : []
}

async function loadLevelOptions() {
  const lvl = useVipList('levels')
  await lvl.load()
  levelOptions.value = lvl.list.value
}

function parseSkuLabel(attrs: any, skuID: number) {
  try {
    const list = typeof attrs === 'string' ? JSON.parse(attrs) : attrs
    if (Array.isArray(list) && list.length) return list.map((a: any) => `${a.name}:${a.value}`).join(' / ')
  } catch { /* ignore */ }
  return `SKU#${skuID}`
}

async function onSelectProduct() {
  form.sku_id = 0; skuOptions.value = []
  if (!form.product_id) return
  const { getProductDetail } = await import('@/api/product')
  const detail: any = await getProductDetail(form.product_id)
  const skus = Array.isArray(detail?.skus) ? detail.skus : []
  skuOptions.value = skus.map((s: any) => ({ id: Number(s.id || 0), label: parseSkuLabel(s.attrs, Number(s.id || 0)) })).filter((s: any) => s.id > 0)
}

function openCreate() { editingID.value = 0; resetForm(); showPopup.value = true }

function openEdit(item: any) {
  editingID.value = Number(item.id || 0)
  form.product_id = Number(item.product_id || 0)
  form.sku_id = Number(item.sku_id || 0)
  form.level_id = Number(item.level_id || 0)
  form.vip_price = String(item.vip_price || '')
  form.status = Number(item.status || 0) === 1 ? 1 : 0
  skuOptions.value = [{ id: form.sku_id, label: `SKU#${form.sku_id}` }]
  showPopup.value = true
}

async function save() {
  if (!form.product_id || !form.sku_id || !form.level_id) { uni.showToast({ title: '请选择商品、SKU和等级', icon: 'none' }); return }
  const level = levelOptions.value.find((l: any) => Number(l.id) === form.level_id)
  const payload = { product_id: form.product_id, sku_id: form.sku_id, level_id: form.level_id, level_name: level?.name || '', vip_price: Number(form.vip_price), status: form.status }
  if (editingID.value) await h.update(editingID.value, payload)
  else await h.create(payload)
  showPopup.value = false
  await loadFiltered()
}

async function remove(id: number) {
  await h.remove(id)
  await loadFiltered()
}

onShow(async () => {
  await Promise.all([loadProductOptions(), loadLevelOptions()])
  await loadFiltered()
})
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <text class="title">会员SKU价格</text>
      <up-button size="mini" type="primary" @click="openCreate">新增</up-button>
    </view>
    <view class="filter-row">
      <picker mode="selector" :range="[{ id: 0, title: '全部商品' }, ...productOptions]" range-key="title" @change="(e: any) => { filterProductID = Number([{ id: 0 }, ...productOptions][e.detail.value]?.id || 0); loadFiltered() }">
        <view class="picker-sm">{{ productOptions.find((p: any) => Number(p.id) === filterProductID)?.title || '全部商品' }}</view>
      </picker>
      <picker mode="selector" :range="[{ id: 0, name: '全部等级' }, ...levelOptions]" range-key="name" @change="(e: any) => { filterLevelID = Number([{ id: 0 }, ...levelOptions][e.detail.value]?.id || 0); loadFiltered() }">
        <view class="picker-sm">{{ levelOptions.find((l: any) => Number(l.id) === filterLevelID)?.name || '全部等级' }}</view>
      </picker>
    </view>
    <view v-if="!h.loading.value && !h.list.value.length" class="empty">暂无数据</view>
    <view v-for="item in h.list.value" :key="item.id" class="card">
      <view class="row">
        <text class="name">商品#{{ item.product_id }} / SKU#{{ item.sku_id }}</text>
        <text :class="['badge', item.status === 1 ? 'on' : 'off']">{{ item.status === 1 ? '启用' : '禁用' }}</text>
      </view>
      <view class="desc">等级：{{ item.level_name || '-' }} · 会员价 ¥{{ Number(item.vip_price || 0).toFixed(2) }}</view>
      <view class="actions">
        <up-button size="mini" type="primary" plain @click="openEdit(item)">编辑</up-button>
        <up-button size="mini" type="error" plain @click="remove(item.id)">删除</up-button>
      </view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中…</view>

    <up-popup :show="showPopup" mode="bottom" round="16" @close="showPopup = false">
      <view class="popup-body">
        <view class="popup-title">{{ editingID ? '编辑价格' : '新增价格' }}</view>
        <picker mode="selector" :range="productOptions" range-key="title" @change="async (e: any) => { form.product_id = Number(productOptions[e.detail.value]?.id || 0); await onSelectProduct() }">
          <view class="picker">{{ productOptions.find((p: any) => Number(p.id) === form.product_id)?.title || '请选择商品' }}</view>
        </picker>
        <view class="mt" />
        <picker mode="selector" :range="skuOptions" range-key="label" @change="(e: any) => { form.sku_id = Number(skuOptions[e.detail.value]?.id || 0) }">
          <view class="picker">{{ skuOptions.find((s: any) => Number(s.id) === form.sku_id)?.label || '请选择SKU' }}</view>
        </picker>
        <view class="mt" />
        <picker mode="selector" :range="levelOptions" range-key="name" @change="(e: any) => { form.level_id = Number(levelOptions[e.detail.value]?.id || 0) }">
          <view class="picker">{{ levelOptions.find((l: any) => Number(l.id) === form.level_id)?.name || '请选择等级' }}</view>
        </picker>
        <view class="mt" />
        <up-input v-model="form.vip_price" type="digit" placeholder="会员价" />
        <view class="mt-lg" />
        <up-button type="primary" @click="save">保存</up-button>
      </view>
    </up-popup>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; align-content: start; }
.top-bar { display: flex; align-items: center; justify-content: space-between; }
.title { font-size: 32rpx; font-weight: 700; }
.filter-row { display: flex; gap: 12rpx; }
.picker-sm { min-height: 60rpx; border: 1px solid var(--eapp-border); border-radius: 12rpx; padding: 0 16rpx; display: flex; align-items: center; font-size: 24rpx; color: var(--eapp-text); }
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
.popup-body { padding: 24rpx; box-sizing: border-box; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.picker { min-height: 76rpx; border: 1px solid var(--eapp-border); border-radius: 12rpx; padding: 0 20rpx; display: flex; align-items: center; color: var(--eapp-text); }
.mt { margin-top: 12rpx; }
.mt-lg { margin-top: 16rpx; }
</style>
