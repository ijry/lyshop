<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { getImSessions } from '@/api/im'

const list = ref<any[]>([])

async function loadData() {
  const data: any = await getImSessions({ page: 1, size: 50 })
  list.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
}

function statusLabel(status: number) {
  if (status === 1) return '等待中'
  if (status === 2) return '进行中'
  if (status === 3) return '已关闭'
  return ''
}

function statusClass(status: number) {
  if (status === 1) return 'status-waiting'
  if (status === 2) return 'status-active'
  if (status === 3) return 'status-closed'
  return ''
}

function truncate(text: string, max = 30) {
  if (!text) return '-'
  return text.length > max ? text.slice(0, max) + '...' : text
}

function goChat(sessionId: number) {
  uni.navigateTo({ url: `/pages/im/chat?session_id=${sessionId}` })
}

onShow(loadData)
</script>

<template>
  <view class="page">
    <view v-if="!list.length" class="empty">暂无会话</view>
    <view v-for="item in list" :key="item.id" class="card" @click="goChat(item.id)">
      <view class="card-top">
        <view class="title-row">
          <text class="title">{{ item.user_nickname || item.user_id || '-' }}</text>
          <view v-if="item.unread_count" class="unread">{{ item.unread_count > 99 ? '99+' : item.unread_count }}</view>
        </view>
        <text :class="['status-tag', statusClass(item.status)]">{{ statusLabel(item.status) }}</text>
      </view>
      <view class="desc">{{ truncate(item.last_msg || item.last_message || '') }}</view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 12rpx; align-content: start; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.card-top { display: flex; align-items: center; justify-content: space-between; }
.title-row { display: flex; align-items: center; gap: 10rpx; }
.title { font-size: 28rpx; font-weight: 600; }
.unread { min-width: 32rpx; height: 32rpx; line-height: 32rpx; text-align: center; font-size: 20rpx; background: #dc2626; color: #fff; border-radius: 999rpx; padding: 0 8rpx; }
.status-tag { font-size: 20rpx; padding: 4rpx 12rpx; border-radius: 999rpx; }
.status-waiting { background: #fef3c7; color: #d97706; }
.status-active { background: #dcfce7; color: #16a34a; }
.status-closed { background: #f1f5f9; color: #64748b; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
</style>
