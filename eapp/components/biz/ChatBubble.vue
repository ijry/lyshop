<script setup lang="ts">
defineProps<{
  message: {
    sender_type: 'user' | 'staff' | number
    content: string
    created_at: string
    type?: string
  }
}>()

function isStaff(senderType: string | number) {
  return senderType === 'staff' || senderType === 2
}

function formatTime(t: string) {
  if (!t) return ''
  return t.slice(11, 16) || ''
}
</script>

<template>
  <view :class="['bubble-wrap', isStaff(message.sender_type) ? 'bubble-right' : 'bubble-left']">
    <view class="sender">{{ isStaff(message.sender_type) ? '客服' : '买家' }}</view>
    <view :class="['bubble', isStaff(message.sender_type) ? 'bubble-staff' : 'bubble-user']">
      <text class="bubble-text">{{ message.content }}</text>
    </view>
    <view class="time">{{ formatTime(message.created_at) }}</view>
  </view>
</template>

<style scoped>
.bubble-wrap { display: flex; flex-direction: column; margin-bottom: 20rpx; max-width: 80%; }
.bubble-left { align-self: flex-start; align-items: flex-start; }
.bubble-right { align-self: flex-end; align-items: flex-end; }
.sender { font-size: 20rpx; color: var(--eapp-text-muted); margin-bottom: 4rpx; }
.bubble { padding: 16rpx 20rpx; border-radius: 16rpx; max-width: 100%; word-break: break-all; }
.bubble-user { background: #f1f5f9; color: #1e293b; border-top-left-radius: 4rpx; }
.bubble-staff { background: #dbeafe; color: #1e40af; border-top-right-radius: 4rpx; }
.bubble-text { font-size: 26rpx; line-height: 1.5; }
.time { font-size: 18rpx; color: var(--eapp-text-muted); margin-top: 4rpx; }
</style>
