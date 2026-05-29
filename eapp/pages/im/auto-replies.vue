<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { getAutoReplies, createAutoReply, updateAutoReply, deleteAutoReply } from '@/api/im'

const list = ref<any[]>([])
const loading = ref(false)
const showPopup = ref(false)
const editingID = ref(0)
const form = reactive({ keyword: '', match_type: 'contains', reply: '', status: 1 })

const matchTypes = [
  { label: '精确匹配', value: 'exact' },
  { label: '包含匹配', value: 'contains' },
]

function resetForm() {
  form.keyword = ''; form.match_type = 'contains'; form.reply = ''; form.status = 1
}

async function loadData() {
  loading.value = true
  try {
    const res: any = await getAutoReplies()
    list.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingID.value = 0; resetForm(); showPopup.value = true
}

function openEdit(item: any) {
  editingID.value = Number(item.id || 0)
  form.keyword = String(item.keyword || '')
  form.match_type = String(item.match_type || 'contains')
  form.reply = String(item.reply || '')
  form.status = Number(item.status || 0) === 1 ? 1 : 0
  showPopup.value = true
}

async function save() {
  if (!form.keyword.trim() || !form.reply.trim()) {
    uni.showToast({ title: '请填写关键词和回复内容', icon: 'none' }); return
  }
  const payload = {
    keyword: form.keyword.trim(), match_type: form.match_type,
    reply: form.reply.trim(), status: form.status,
  }
  if (editingID.value) await updateAutoReply(editingID.value, payload)
  else await createAutoReply(payload)
  showPopup.value = false
  await loadData()
}

async function remove(id: number) {
  uni.showModal({
    title: '确认删除', content: '确定删除该自动回复规则吗？',
    success: async (res) => {
      if (res.confirm) {
        await deleteAutoReply(id)
        await loadData()
      }
    },
  })
}

function matchLabel(type: string) {
  return type === 'exact' ? '精确' : '包含'
}

onShow(loadData)
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <text class="title">自动回复</text>
      <up-button size="mini" type="primary" @click="openCreate">新增规则</up-button>
    </view>

    <view v-if="loading" class="loading">加载中...</view>
    <view v-else-if="!list.length" class="empty">暂无自动回复规则</view>
    <view v-for="item in list" :key="item.id" class="card">
      <view class="row">
        <text class="name">{{ item.keyword }}</text>
        <text :class="['badge', item.status === 1 ? 'on' : 'off']">{{ item.status === 1 ? '启用' : '禁用' }}</text>
      </view>
      <view class="desc">匹配方式: {{ matchLabel(item.match_type) }}</view>
      <view class="desc reply-text">回复: {{ item.reply }}</view>
      <view class="actions">
        <up-button size="mini" type="primary" plain @click="openEdit(item)">编辑</up-button>
        <up-button size="mini" type="error" plain @click="remove(item.id)">删除</up-button>
      </view>
    </view>

    <up-popup :show="showPopup" mode="bottom" round="16" @close="showPopup = false">
      <view class="popup-body">
        <view class="popup-title">{{ editingID ? '编辑规则' : '新增规则' }}</view>
        <up-input v-model="form.keyword" placeholder="触发关键词" clearable />
        <view class="mt" />
        <view class="label">匹配方式</view>
        <view class="type-row">
          <view v-for="opt in matchTypes" :key="opt.value" :class="['type-tag', form.match_type === opt.value ? 'active' : '']" @click="form.match_type = opt.value">{{ opt.label }}</view>
        </view>
        <view class="mt" />
        <up-textarea v-model="form.reply" placeholder="回复内容" />
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
.reply-text { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.actions { margin-top: 12rpx; display: flex; gap: 10rpx; justify-content: flex-end; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.popup-body { padding: 24rpx; box-sizing: border-box; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.label { font-size: 26rpx; font-weight: 600; margin-bottom: 8rpx; }
.type-row { display: flex; gap: 10rpx; }
.type-tag { padding: 8rpx 20rpx; font-size: 24rpx; border-radius: 999rpx; background: #f1f5f9; border: 1px solid var(--eapp-border); }
.type-tag.active { background: var(--eapp-primary, #2563eb); color: #fff; border-color: var(--eapp-primary, #2563eb); }
.mt { margin-top: 12rpx; }
.mt-lg { margin-top: 16rpx; }
</style>
