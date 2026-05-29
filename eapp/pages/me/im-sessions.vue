<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { getImSessions } from '@/api/im'
import EmptyState from '@/components/biz/EmptyState.vue'

const list = ref<any[]>([])
const loading = ref(false)

async function loadData() {
  loading.value = true
  try {
    const data: any = await getImSessions({ page: 1, size: 50 })
    list.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
  } finally {
    loading.value = false
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

onShow(loadData)
</script>

<template>
  <view class="page">
    <view class="header">
      <text class="header-title">客服会话</text>
      <text class="header-count">{{ list.length }} 个会话</text>
    </view>

    <EmptyState v-if="!loading && !list.length" title="暂无会话" desc="等待买家发起咨询" />

    <view class="list">
      <view v-for="item in list" :key="item.id" class="session" @click="goChat(item.id)">
        <view class="avatar-wrap">
          <image v-if="item.user_avatar" :src="item.user_avatar" mode="aspectFill" class="avatar" />
          <view v-else class="avatar avatar-placeholder">
            <text>{{ (item.user_nickname || 'U')[0] }}</text>
          </view>
          <view v-if="statusTone(item.status) === 'success'" class="online-dot" />
        </view>

        <view class="session-body">
          <view class="session-top">
            <text class="nickname">{{ item.user_nickname || `用户${item.user_id}` }}</text>
            <text class="time">{{ timeAgo(item.updated_at) }}</text>
          </view>
          <view class="session-bottom">
            <text class="last-msg">{{ item.last_msg || '-' }}</text>
            <view v-if="item.unread_count" class="unread">{{ item.unread_count > 99 ? '99+' : item.unread_count }}</view>
          </view>
        </view>

        <view :class="['status-pill', `tone-${statusTone(item.status)}`]">
          <text class="status-text">{{ statusLabel(item.status) }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }

.header { padding: 24rpx 28rpx 16rpx; display: flex; align-items: baseline; justify-content: space-between; }
.header-title { font-size: 34rpx; font-weight: 700; color: var(--eapp-text); }
.header-count { font-size: 22rpx; color: var(--eapp-text-muted); }

.list { padding: 0 20rpx 20rpx; display: grid; gap: 2rpx; }

.session {
  display: flex; align-items: center; gap: 20rpx;
  background: var(--eapp-card); padding: 24rpx 20rpx;
  border-bottom: 1rpx solid var(--eapp-border);
}
.session:first-child { border-radius: 20rpx 20rpx 0 0; }
.session:last-child { border-radius: 0 0 20rpx 20rpx; border-bottom: none; }
.session:only-child { border-radius: 20rpx; }

.avatar-wrap { position: relative; flex-shrink: 0; }
.avatar { width: 88rpx; height: 88rpx; border-radius: 50%; background: var(--eapp-bg); }
.avatar-placeholder { display: flex; align-items: center; justify-content: center; background: var(--eapp-primary-soft); color: var(--eapp-primary); font-size: 30rpx; font-weight: 700; }
.online-dot { position: absolute; bottom: 4rpx; right: 4rpx; width: 18rpx; height: 18rpx; background: var(--eapp-success); border: 3rpx solid var(--eapp-card); border-radius: 50%; }

.session-body { flex: 1; min-width: 0; }
.session-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8rpx; }
.nickname { font-size: 28rpx; font-weight: 600; color: var(--eapp-text); }
.time { font-size: 22rpx; color: var(--eapp-text-faint); flex-shrink: 0; }
.session-bottom { display: flex; align-items: center; justify-content: space-between; gap: 12rpx; }
.last-msg { font-size: 24rpx; color: var(--eapp-text-muted); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; }
.unread { min-width: 36rpx; height: 36rpx; line-height: 36rpx; text-align: center; font-size: 20rpx; background: var(--eapp-danger); color: #fff; border-radius: 999rpx; padding: 0 10rpx; flex-shrink: 0; }

.status-pill { padding: 6rpx 16rpx; border-radius: 999rpx; flex-shrink: 0; }
.status-text { font-size: 20rpx; }
.tone-warning { background: var(--eapp-warning-soft); }
.tone-warning .status-text { color: #d97706; }
.tone-success { background: var(--eapp-success-soft); }
.tone-success .status-text { color: #16a34a; }
.tone-muted { background: #f1f5f9; }
.tone-muted .status-text { color: #94a3b8; }
</style>
