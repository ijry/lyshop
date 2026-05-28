<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { useVipList } from '@/composables/useVipList'

const h = useVipList('coupon-rules')
const showPopup = ref(false)
const editingID = ref(0)
const form = reactive({ name: '', coupon_name: '', monthly_limit: 1 })

function resetForm() { form.name = ''; form.coupon_name = ''; form.monthly_limit = 1 }

function openCreate() { editingID.value = 0; resetForm(); showPopup.value = true }

function openEdit(item: any) {
  editingID.value = Number(item.id || 0)
  form.name = String(item.name || '')
  form.coupon_name = String(item.coupon_name || '')
  form.monthly_limit = Number(item.monthly_limit || 1)
  showPopup.value = true
}

async function save() {
  if (!form.name.trim()) { uni.showToast({ title: '请输入规则名称', icon: 'none' }); return }
  const payload = { name: form.name.trim(), coupon_name: form.coupon_name.trim(), monthly_limit: Number(form.monthly_limit) }
  if (editingID.value) await h.update(editingID.value, payload)
  else await h.create(payload)
  showPopup.value = false
  await h.load()
}

async function remove(id: number) {
  await h.remove(id)
  await h.load()
}

onShow(() => h.load())
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <text class="title">会员领券规则</text>
      <up-button size="mini" type="primary" @click="openCreate">新增规则</up-button>
    </view>
    <view v-if="!h.loading.value && !h.list.value.length" class="empty">暂无规则</view>
    <view v-for="item in h.list.value" :key="item.id" class="card">
      <view class="row">
        <text class="name">{{ item.name }}</text>
      </view>
      <view class="desc">优惠券：{{ item.coupon_name || '-' }} · 月限 {{ item.monthly_limit }}</view>
      <view class="actions">
        <up-button size="mini" type="primary" plain @click="openEdit(item)">编辑</up-button>
        <up-button size="mini" type="error" plain @click="remove(item.id)">删除</up-button>
      </view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中…</view>

    <up-popup :show="showPopup" mode="bottom" round="16" @close="showPopup = false">
      <view class="popup-body">
        <view class="popup-title">{{ editingID ? '编辑规则' : '新增规则' }}</view>
        <up-input v-model="form.name" placeholder="规则名称" clearable />
        <view class="mt" />
        <up-input v-model="form.coupon_name" placeholder="优惠券名称" clearable />
        <view class="mt" />
        <up-input v-model="form.monthly_limit" type="number" placeholder="每月限领次数" />
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
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.actions { margin-top: 12rpx; display: flex; gap: 10rpx; justify-content: flex-end; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
.popup-body { padding: 24rpx; box-sizing: border-box; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.mt { margin-top: 12rpx; }
.mt-lg { margin-top: 16rpx; }
</style>
