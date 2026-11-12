<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { getImSessions, getStaffStatus, setStaffOnline, acceptSession, closeSession, transferSession } from '@/api/im'
import EmptyState from '@/components/biz/EmptyState.vue'
import PageHeader from '@/components/biz/PageHeader.vue'
import AppIcon from '@/components/common/AppIcon.vue'

const list = ref<any[]>([])
const loading = ref(false)
const staffOnline = ref(false)
const staffStatus = ref({ current_load: 0, max_load: 5 })
const showTransferModal = ref(false)
const transferSessionId = ref(0)
const transferToStaffId = ref('')
const transferRemark = ref('')

async function loadData() {
  loading.value = true
  try {
    const data: any = await getImSessions({ page: 1, size: 50 })
    list.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
  } finally {
    loading.value = false
  }
}

async function loadStaffStatus() {
  try {
    const data: any = await getStaffStatus()
    if (data) {
      staffStatus.value = data
      staffOnline.value = data.is_online === 1
    }
  } catch {}
}

async function toggleOnline() {
  const newState = !staffOnline.value
  try {
    await setStaffOnline({ online: newState })
    staffOnline.value = newState
    await loadStaffStatus()
    if (newState) loadData()
    uni.showToast({ title: newState ? '已上线' : '已下线', icon: 'success' })
  } catch (e: any) {
    uni.showToast({ title: e.message || '操作失败', icon: 'none' })
  }
}

async function onAccept(sessionId: number) {
  try {
    await acceptSession(sessionId)
    uni.showToast({ title: '已接入', icon: 'success' })
    await Promise.all([loadData(), loadStaffStatus()])
  } catch (err: any) {
    uni.showToast({ title: err.message || '接入失败', icon: 'none' })
  }
}

async function onClose(sessionId: number) {
  uni.showModal({
    title: '确认结束',
    content: '确认结束该会话？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await closeSession(sessionId)
          uni.showToast({ title: '已结束', icon: 'success' })
          await Promise.all([loadData(), loadStaffStatus()])
        } catch (err: any) {
          uni.showToast({ title: err.message || '结束失败', icon: 'none' })
        }
      }
    }
  })
}

function openTransferModal(sessionId: number) {
  transferSessionId.value = sessionId
  transferToStaffId.value = ''
  transferRemark.value = ''
  showTransferModal.value = true
}

async function confirmTransfer() {
  if (!transferToStaffId.value.trim()) {
    uni.showToast({ title: '请输入目标客服ID', icon: 'none' })
    return
  }
  try {
    await transferSession(transferSessionId.value, {
      to_staff_id: Number(transferToStaffId.value),
      remark: transferRemark.value
    })
    uni.showToast({ title: '转接成功', icon: 'success' })
    showTransferModal.value = false
    await Promise.all([loadData(), loadStaffStatus()])
  } catch (err: any) {
    uni.showToast({ title: err.message || '转接失败', icon: 'none' })
  }
}

function statusLabel(status: number) {
  if (status === 1) return '等待接入'
  if (status === 2) return '服务中'
  if (status === 3) return '已关闭'
  return ''
}

function statusTone(status: number) {
  if (status === 1) return 'warning'
  if (status === 2) return 'success'
  return 'muted'
}

function statusTagType(status: number): string {
  if (status === 1) return 'warning'
  if (status === 2) return 'success'
  return 'info'
}

function timeAgo(t: string) {
  if (!t) return ''
  const d = new Date(t)
  const now = new Date()
  const diff = Math.floor((now.getTime() - d.getTime()) / 1000)
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return t.slice(5, 16).replace('T', ' ')
}

function goChat(sessionId: number) {
  uni.navigateTo({ url: `/pages/im/chat?session_id=${sessionId}` })
}

onShow(() => {
  loadData()
  loadStaffStatus()
})
</script>

<template>
  <view class="page">
    <PageHeader title="客服会话" eyebrow="CUSTOMER SERVICE" :subtitle="`${list.length} 个会话`">
      <template #right>
        <view class="eapp-icon-btn" hover-class="eapp-hover" @click="toggleOnline">
          <AppIcon :name="staffOnline ? 'check-circle' : 'info'" :size="30" :color="staffOnline ? 'success' : 'muted'" />
        </view>
      </template>
    </PageHeader>

    <view class="status-bar">
      <view class="status-bar__left">
        <view :class="['status-dot', staffOnline && 'status-dot--on']" />
        <text class="status-bar__text">{{ staffOnline ? '接待中' : '已离线' }}</text>
        <text v-if="staffOnline" class="status-bar__load eapp-num">
          当前接待 {{ staffStatus.current_load }}/{{ staffStatus.max_load }}
        </text>
      </view>
      <view class="status-bar__toggle" hover-class="eapp-hover" @click="toggleOnline">
        <text class="status-bar__toggle-text">{{ staffOnline ? '切换为离线' : '上线接客' }}</text>
      </view>
    </view>

    <EmptyState v-if="!loading && !list.length" title="暂无会话" desc="等待买家发起咨询" />

    <view class="list">
      <view v-for="item in list" :key="item.id" class="session" hover-class="eapp-hover" :hover-stay-time="80" @click="goChat(item.id)">
        <view class="avatar-wrap">
          <image v-if="item.user_avatar" :src="item.user_avatar" mode="aspectFill" class="avatar" />
          <view v-else class="avatar avatar-placeholder">
            <text>{{ (item.user_nickname || 'U')[0] }}</text>
          </view>
          <view v-if="statusTone(item.status) === 'success'" class="online-dot" />
        </view>

        <view class="session-body">
          <view class="session-top">
            <view class="nickname-row">
              <text class="nickname">{{ item.user_nickname || `用户${item.user_id}` }}</text>
              <text v-if="item.queue_position > 0" class="queue-badge">排队第{{ item.queue_position }}位</text>
            </view>
            <text class="time">{{ timeAgo(item.updated_at) }}</text>
          </view>
          <view class="session-bottom">
            <text class="last-msg">{{ item.last_msg || '-' }}</text>
            <view v-if="item.unread_count" class="unread">{{ item.unread_count > 99 ? '99+' : item.unread_count }}</view>
          </view>
        </view>

        <view class="session-actions">
          <up-tag :text="statusLabel(item.status)" :type="statusTagType(item.status)" plain size="mini" />
          <up-tag v-if="item.status === 1" text="接入" type="success" size="mini" @click.stop="onAccept(item.id)" />
          <view v-if="item.status === 2" class="action-row">
            <up-tag text="转接" type="primary" size="mini" @click.stop="openTransferModal(item.id)" />
            <up-tag text="结束" type="info" plain size="mini" @click.stop="onClose(item.id)" />
          </view>
        </view>
      </view>
    </view>

    <!-- Transfer Modal -->
    <up-popup :show="showTransferModal" mode="center" round="20" @close="showTransferModal = false">
      <view class="transfer-modal">
        <view class="modal-title">转接会话</view>
        <view class="modal-body">
          <view class="form-item">
            <text class="form-label">目标客服ID</text>
            <up-input v-model="transferToStaffId" placeholder="输入客服ID" type="number" />
          </view>
          <view class="form-item">
            <text class="form-label">转接备注（可选）</text>
            <up-textarea v-model="transferRemark" placeholder="如：专业问题转技术客服" :autoHeight="true" />
          </view>
        </view>
        <view class="modal-actions">
          <up-button plain @click="showTransferModal = false">取消</up-button>
          <view style="width: 16rpx;" />
          <up-button type="primary" @click="confirmTransfer">确认转接</up-button>
        </view>
      </view>
    </up-popup>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }

.status-bar {
  display: flex; align-items: center; justify-content: space-between; gap: 16rpx;
  padding: 20rpx 28rpx; margin-bottom: 8rpx;
  background: var(--eapp-card); border-bottom: 1rpx solid var(--eapp-border);
}
.status-bar__left { display: flex; align-items: center; gap: 10rpx; min-width: 0; }
.status-dot { width: 14rpx; height: 14rpx; border-radius: 50%; background: var(--eapp-border-strong); flex-shrink: 0; }
.status-dot--on { background: var(--eapp-success); }
.status-bar__text { font-size: 24rpx; font-weight: 600; color: var(--eapp-text-secondary); flex-shrink: 0; }
.status-bar__load { font-size: 22rpx; color: var(--eapp-text-faint); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.status-bar__toggle { flex-shrink: 0; padding: 8rpx 18rpx; border: 1rpx solid var(--eapp-primary-border); border-radius: var(--eapp-radius-sm); background: var(--eapp-primary-soft); }
.status-bar__toggle-text { font-size: 22rpx; font-weight: 600; color: var(--eapp-primary); }

.list { padding: 20rpx; padding-bottom: calc(env(safe-area-inset-bottom) + 40rpx); display: grid; gap: 16rpx; }

.session {
  display: flex; align-items: center; gap: 20rpx;
  background: var(--eapp-card); padding: 24rpx 20rpx;
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  box-sizing: border-box;
}

.avatar-wrap { position: relative; flex-shrink: 0; }
.avatar { width: 88rpx; height: 88rpx; border-radius: 50%; background: var(--eapp-bg); }
.avatar-placeholder { display: flex; align-items: center; justify-content: center; background: var(--eapp-primary-soft); color: var(--eapp-primary); font-size: 30rpx; font-weight: 700; }
.online-dot { position: absolute; bottom: 4rpx; right: 4rpx; width: 18rpx; height: 18rpx; background: var(--eapp-success); border: 3rpx solid var(--eapp-card); border-radius: 50%; }

.session-body { flex: 1; min-width: 0; }
.session-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8rpx; }
.nickname-row { display: flex; align-items: center; gap: 8rpx; }
.nickname { font-size: 28rpx; font-weight: 600; color: var(--eapp-text); }
.queue-badge { font-size: 20rpx; color: #d97706; background: var(--eapp-warning-soft); padding: 2rpx 10rpx; border-radius: 999rpx; }
.time { font-size: 22rpx; color: var(--eapp-text-faint); flex-shrink: 0; }
.session-bottom { display: flex; align-items: center; justify-content: space-between; gap: 12rpx; }
.last-msg { font-size: 24rpx; color: var(--eapp-text-muted); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; }
.unread { min-width: 36rpx; height: 36rpx; line-height: 36rpx; text-align: center; font-size: 20rpx; background: var(--eapp-danger); color: #fff; border-radius: 999rpx; padding: 0 10rpx; flex-shrink: 0; }

.session-actions { display: flex; flex-direction: column; gap: 8rpx; align-items: flex-end; flex-shrink: 0; }
.action-row { display: flex; gap: 8rpx; }

.transfer-modal { width: 600rpx; padding: 32rpx 24rpx 24rpx; }
.modal-title { font-size: 32rpx; font-weight: 700; margin-bottom: 24rpx; text-align: center; }
.modal-body { margin-bottom: 24rpx; }
.form-item { margin-bottom: 20rpx; }
.form-label { display: block; font-size: 26rpx; color: var(--eapp-text); margin-bottom: 12rpx; }
.modal-actions { display: flex; gap: 16rpx; }
</style>
