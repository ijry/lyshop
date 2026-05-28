<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { useVipList } from '@/composables/useVipList'

const h = useVipList('levels')
const showPopup = ref(false)
const editingID = ref(0)
const form = reactive({ name: '', growth_min: 0, discount_rate: '' })

function resetForm() { form.name = ''; form.growth_min = 0; form.discount_rate = '' }

function openCreate() { editingID.value = 0; resetForm(); showPopup.value = true }

function openEdit(item: any) {
  editingID.value = Number(item.id || 0)
  form.name = String(item.name || '')
  form.growth_min = Number(item.growth_min || 0)
  form.discount_rate = String(item.discount_rate || '')
  showPopup.value = true
}

async function save() {
  if (!form.name.trim()) { uni.showToast({ title: '请输入等级名称', icon: 'none' }); return }
  const payload = { name: form.name.trim(), growth_min: Number(form.growth_min), discount_rate: Number(form.discount_rate) }
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
      <text class="title">会员等级</text>
      <up-button size="mini" type="primary" @click="openCreate">新增等级</up-button>
    </view>
    <view v-if="!h.loading.value && !h.list.value.length" class="empty">暂无等级</view>
    <view v-for="item in h.list.value" :key="item.id" class="card">
      <view class="row">
        <text class="name">{{ item.name }}</text>
      </view>
      <view class="desc">成长值 ≥ {{ item.growth_min }} · 折扣率 {{ item.discount_rate }}</view>
      <view class="actions">
        <up-button size="mini" type="primary" plain @click="openEdit(item)">编辑</up-button>
        <up-button size="mini" type="error" plain @click="remove(item.id)">删除</up-button>
      </view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中…</view>

    <up-popup :show="showPopup" mode="bottom" round="16" @close="showPopup = false">
      <view class="popup-body">
        <view class="popup-title">{{ editingID ? '编辑等级' : '新增等级' }}</view>
        <up-input v-model="form.name" placeholder="等级名称" clearable />
        <view class="mt" />
        <up-input v-model="form.growth_min" type="number" placeholder="最低成长值" />
        <view class="mt" />
        <up-input v-model="form.discount_rate" type="digit" placeholder="折扣率（如 0.95）" />
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
