<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import request from '@/utils/request'
import EmptyState from '@/components/biz/EmptyState.vue'

const list = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingStaff = ref<any>(null)
const form = ref({ admin_id: '', max_load: 5 })

async function loadData() {
  loading.value = true
  try {
    const data: any = await request.get('/im/staff')
    list.value = Array.isArray(data) ? data : []
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingStaff.value = null
  form.value = { admin_id: '', max_load: 5 }
  showModal.value = true
}

function openEditModal(staff: any) {
  editingStaff.value = staff
  form.value = { admin_id: String(staff.admin_id), max_load: staff.max_load }
  showModal.value = true
}

async function submitForm() {
  if (!editingStaff.value && !form.value.admin_id.trim()) {
    uni.showToast({ title: '请输入管理员ID', icon: 'none' })
    return
  }
  if (!form.value.max_load || form.value.max_load < 1) {
    uni.showToast({ title: '最大负载必须大于0', icon: 'none' })
    return
  }

  try {
    if (editingStaff.value) {
      await request.put(`/im/staff/${editingStaff.value.id}`, { max_load: form.value.max_load })
    } else {
      await request.post('/im/staff', { admin_id: Number(form.value.admin_id), max_load: form.value.max_load })
    }
    uni.showToast({ title: '操作成功', icon: 'success' })
    showModal.value = false
    await loadData()
  } catch (err: any) {
    uni.showToast({ title: err.message || '操作失败', icon: 'none' })
  }
}

function confirmDelete(id: number) {
  uni.showModal({
    title: '确认删除',
    content: '确认删除该客服？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await request.delete(`/im/staff/${id}`)
          uni.showToast({ title: '删除成功', icon: 'success' })
          await loadData()
        } catch (err: any) {
          uni.showToast({ title: err.message || '删除失败', icon: 'none' })
        }
      }
    }
  })
}

onShow(loadData)
</script>

<template>
  <view class="page">
    <view class="header">
      <text class="header-title">客服坐席管理</text>
      <up-button size="small" type="primary" @click="openCreateModal">+ 添加客服</up-button>
    </view>

    <EmptyState v-if="!loading && !list.length" title="暂无客服人员" desc="点击右上角添加客服" />

    <view class="list">
      <view v-for="item in list" :key="item.id" class="staff-card">
        <view class="card-header">
          <view class="staff-info">
            <text class="staff-id">ID: {{ item.id }}</text>
            <text class="admin-id">管理员 #{{ item.admin_id }}</text>
          </view>
          <view :class="['status-badge', item.is_online ? 'online' : 'offline']">
            <text class="status-text">{{ item.is_online ? '在线' : '离线' }}</text>
          </view>
        </view>
        <view class="card-body">
          <view class="info-row">
            <text class="info-label">当前负载</text>
            <text class="info-value">{{ item.current_load }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">最大负载</text>
            <text class="info-value">{{ item.max_load }}</text>
          </view>
        </view>
        <view class="card-actions">
          <up-button size="small" plain @click="openEditModal(item)">编辑</up-button>
          <up-button size="small" type="error" plain @click="confirmDelete(item.id)">删除</up-button>
        </view>
      </view>
    </view>

    <!-- Modal -->
    <up-popup :show="showModal" mode="center" round="20" @close="showModal = false">
      <view class="modal">
        <view class="modal-title">{{ editingStaff ? '编辑客服' : '添加客服' }}</view>
        <view class="modal-body">
          <view v-if="!editingStaff" class="form-item">
            <text class="form-label">管理员ID *</text>
            <up-input v-model="form.admin_id" placeholder="输入管理员ID" type="number" />
          </view>
          <view class="form-item">
            <text class="form-label">最大负载 *</text>
            <up-input v-model.number="form.max_load" placeholder="默认5" type="number" />
          </view>
        </view>
        <view class="modal-actions">
          <up-button plain @click="showModal = false">取消</up-button>
          <view style="width: 16rpx;" />
          <up-button type="primary" @click="submitForm">确认</up-button>
        </view>
      </view>
    </up-popup>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding-bottom: 20rpx; }

.header { padding: 24rpx 28rpx 16rpx; display: flex; align-items: center; justify-content: space-between; }
.header-title { font-size: 34rpx; font-weight: 700; color: var(--eapp-text); }

.list { padding: 0 20rpx; display: grid; gap: 16rpx; }

.staff-card { background: var(--eapp-card); border-radius: 20rpx; padding: 24rpx; border: 1rpx solid var(--eapp-border); }
.card-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20rpx; }
.staff-info { display: flex; flex-direction: column; gap: 6rpx; }
.staff-id { font-size: 24rpx; color: var(--eapp-text-muted); }
.admin-id { font-size: 28rpx; font-weight: 600; color: var(--eapp-text); }
.status-badge { padding: 6rpx 16rpx; border-radius: 999rpx; }
.status-badge.online { background: var(--eapp-success-soft); }
.status-badge.offline { background: #f1f5f9; }
.status-badge.online .status-text { color: var(--eapp-success); font-size: 22rpx; font-weight: 600; }
.status-badge.offline .status-text { color: #94a3b8; font-size: 22rpx; font-weight: 600; }

.card-body { display: grid; gap: 12rpx; margin-bottom: 20rpx; }
.info-row { display: flex; justify-content: space-between; align-items: center; }
.info-label { font-size: 26rpx; color: var(--eapp-text-muted); }
.info-value { font-size: 26rpx; font-weight: 600; color: var(--eapp-text); }

.card-actions { display: flex; gap: 12rpx; }

.modal { width: 600rpx; padding: 32rpx 24rpx 24rpx; }
.modal-title { font-size: 32rpx; font-weight: 700; margin-bottom: 24rpx; text-align: center; }
.modal-body { margin-bottom: 24rpx; }
.form-item { margin-bottom: 20rpx; }
.form-label { display: block; font-size: 26rpx; color: var(--eapp-text); margin-bottom: 12rpx; }
.modal-actions { display: flex; gap: 16rpx; }
</style>
