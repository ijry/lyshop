<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { useWmsList } from '@/composables/useWmsList'

const h = useWmsList('warehouses')
const showPopup = ref(false)
const editingID = ref(0)
const form = reactive({ name: '', code: '', address: '', contact: '', phone: '', status: 1 })

function resetForm() {
  form.name = ''; form.code = ''; form.address = ''; form.contact = ''; form.phone = ''; form.status = 1
}

function openCreate() { editingID.value = 0; resetForm(); showPopup.value = true }

function openEdit(item: any) {
  editingID.value = Number(item.id || 0)
  form.name = String(item.name || '')
  form.code = String(item.code || '')
  form.address = String(item.address || '')
  form.contact = String(item.contact || '')
  form.phone = String(item.phone || '')
  form.status = Number(item.status || 0) === 1 ? 1 : 0
  showPopup.value = true
}

async function save() {
  if (!form.name.trim()) { uni.showToast({ title: '请输入仓库名称', icon: 'none' }); return }
  const payload = {
    name: form.name.trim(), code: form.code.trim(), address: form.address.trim(),
    contact: form.contact.trim(), phone: form.phone.trim(), status: form.status,
  }
  if (editingID.value) await h.update(editingID.value, payload)
  else await h.create(payload)
  showPopup.value = false
  await h.load()
}

onShow(() => h.load())
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <text class="title">仓库管理</text>
      <up-button size="mini" type="primary" @click="openCreate">新增仓库</up-button>
    </view>
    <view v-if="!h.loading.value && !h.list.value.length" class="empty">暂无仓库</view>
    <view v-for="item in h.list.value" :key="item.id" class="card">
      <view class="row">
        <text class="name">{{ item.name }}</text>
        <text :class="['badge', item.status === 1 ? 'on' : 'off']">{{ item.status === 1 ? '启用' : '禁用' }}</text>
      </view>
      <view class="desc">编码: {{ item.code || '-' }} | 联系人: {{ item.contact || '-' }}</view>
      <view class="desc">地址: {{ item.address || '-' }}</view>
      <view class="actions">
        <up-button size="mini" type="primary" plain @click="openEdit(item)">编辑</up-button>
      </view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中...</view>

    <up-popup :show="showPopup" mode="bottom" round="16" @close="showPopup = false">
      <view class="popup-body">
        <view class="popup-title">{{ editingID ? '编辑仓库' : '新增仓库' }}</view>
        <up-input v-model="form.name" placeholder="仓库名称" clearable />
        <view class="mt" />
        <up-input v-model="form.code" placeholder="仓库编码" clearable />
        <view class="mt" />
        <up-input v-model="form.address" placeholder="仓库地址" clearable />
        <view class="mt" />
        <up-input v-model="form.contact" placeholder="联系人" clearable />
        <view class="mt" />
        <up-input v-model="form.phone" placeholder="联系电话" clearable />
        <view class="mt" />
        <view class="row">
          <text>启用</text>
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
.mt { margin-top: 12rpx; }
.mt-lg { margin-top: 16rpx; }
</style>
