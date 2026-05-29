<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePointsList } from '@/composables/usePointsList'

const { t } = useI18n()
const h = usePointsList('products')
const typeTab = ref('')
const typeTabs = [
  { label: '全部', value: '' },
  { label: '实物', value: 'physical' },
  { label: '虚拟', value: 'virtual' },
  { label: '优惠券', value: 'coupon' },
]

const showPopup = ref(false)
const editingID = ref(0)
const form = reactive({ title: '', type: 'coupon', points_price: '', stock: '', cover: '', description: '', status: 1 })

function resetForm() { form.title = ''; form.type = 'coupon'; form.points_price = ''; form.stock = ''; form.cover = ''; form.description = ''; form.status = 1 }

function openCreate() { editingID.value = 0; resetForm(); showPopup.value = true }

function openEdit(item: any) {
  editingID.value = Number(item.id || 0)
  form.title = String(item.title || '')
  form.type = String(item.type || 'coupon')
  form.points_price = String(item.points_price || '')
  form.stock = String(item.stock || '')
  form.cover = String(item.cover || '')
  form.description = String(item.description || '')
  form.status = Number(item.status || 0) === 1 ? 1 : 0
  showPopup.value = true
}

async function save() {
  if (!form.title.trim()) { uni.showToast({ title: '请输入商品名称', icon: 'none' }); return }
  const payload = { title: form.title.trim(), type: form.type, points_price: Number(form.points_price), stock: Number(form.stock), cover: form.cover.trim(), description: form.description.trim(), status: form.status }
  if (editingID.value) await h.update(editingID.value, payload)
  else await h.create(payload)
  showPopup.value = false
  doSearch()
}

async function remove(id: number) {
  await h.remove(id)
  doSearch()
}

function doSearch() {
  const params: any = { page: 1, size: 50 }
  if (typeTab.value) params.type = typeTab.value
  h.load(params)
}

function switchTab(val: string) {
  typeTab.value = val
  doSearch()
}

const typeLabels: Record<string, string> = { coupon: '优惠券', physical: '实物', virtual: '虚拟' }

onShow(() => doSearch())
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <text class="title">{{ t('points.productMgmt') }}</text>
      <up-button size="mini" type="primary" @click="openCreate">新增商品</up-button>
    </view>
    <view class="tabs">
      <view v-for="tab in typeTabs" :key="tab.value" :class="['tab', typeTab === tab.value ? 'active' : '']" @click="switchTab(tab.value)">{{ tab.label }}</view>
    </view>
    <view v-if="!h.loading.value && !h.list.value.length" class="empty">暂无积分商品</view>
    <view v-for="item in h.list.value" :key="item.id" class="card">
      <view class="row">
        <image v-if="item.cover" :src="item.cover" class="cover" mode="aspectFill" />
        <view class="info">
          <text class="name">{{ item.title }}</text>
          <text class="type-badge">{{ typeLabels[item.type] || item.type }}</text>
        </view>
        <text :class="['badge', item.status === 1 ? 'on' : 'off']">{{ item.status === 1 ? '上架' : '下架' }}</text>
      </view>
      <view class="desc">{{ item.points_price }} 积分 · 库存 {{ item.stock }}</view>
      <view class="actions">
        <up-button size="mini" type="primary" plain @click="openEdit(item)">编辑</up-button>
        <up-button size="mini" type="error" plain @click="remove(item.id)">删除</up-button>
      </view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中…</view>

    <up-popup :show="showPopup" mode="bottom" round="16" @close="showPopup = false">
      <view class="popup-body">
        <view class="popup-title">{{ editingID ? '编辑商品' : '新增商品' }}</view>
        <up-input v-model="form.title" placeholder="商品名称" clearable />
        <view class="mt" />
        <view class="row">
          <text class="label">类型</text>
          <picker :value="['coupon','physical','virtual'].indexOf(form.type)" :range="['优惠券','实物','虚拟']" @change="(e: any) => form.type = ['coupon','physical','virtual'][e.detail.value]">
            <text class="picker-text">{{ typeLabels[form.type] || form.type }}</text>
          </picker>
        </view>
        <view class="mt" />
        <up-input v-model="form.points_price" type="number" placeholder="积分价格" />
        <view class="mt" />
        <up-input v-model="form.stock" type="number" placeholder="库存" />
        <view class="mt" />
        <up-input v-model="form.cover" placeholder="封面图 URL" clearable />
        <view class="mt" />
        <up-input v-model="form.description" placeholder="描述" clearable />
        <view class="mt" />
        <view class="row">
          <text>上架</text>
          <switch :checked="form.status === 1" @change="(e: any) => form.status = e.detail.value ? 1 : 0" />
        </view>
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
.tabs { display: flex; gap: 10rpx; }
.tab { padding: 8rpx 20rpx; border-radius: 999rpx; font-size: 24rpx; background: #f1f5f9; color: #64748b; }
.tab.active { background: #2563eb; color: #fff; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.row { display: flex; align-items: center; justify-content: space-between; }
.cover { width: 80rpx; height: 80rpx; border-radius: 12rpx; margin-right: 14rpx; flex-shrink: 0; }
.info { flex: 1; }
.name { font-size: 28rpx; font-weight: 600; display: block; }
.type-badge { font-size: 20rpx; padding: 2rpx 10rpx; border-radius: 999rpx; background: #e0e7ff; color: #4338ca; }
.badge { font-size: 22rpx; padding: 4rpx 14rpx; border-radius: 999rpx; }
.badge.on { background: #dcfce7; color: #16a34a; }
.badge.off { background: #fee2e2; color: #dc2626; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.actions { margin-top: 12rpx; display: flex; gap: 10rpx; justify-content: flex-end; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
.popup-body { padding: 24rpx; box-sizing: border-box; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.label { font-size: 26rpx; color: #475569; }
.picker-text { font-size: 26rpx; color: #2563eb; }
.mt { margin-top: 12rpx; }
.mt-lg { margin-top: 16rpx; }
</style>
