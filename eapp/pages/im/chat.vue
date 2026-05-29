<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { ref, nextTick } from 'vue'
import { getImMessages, sendImMessage } from '@/api/im'
import ChatBubble from '@/components/biz/ChatBubble.vue'

const sessionId = ref(0)
const messages = ref<any[]>([])
const inputText = ref('')
const loading = ref(false)
const scrollId = ref('')

async function loadMessages() {
  if (!sessionId.value) return
  loading.value = true
  try {
    const res: any = await getImMessages(sessionId.value)
    messages.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
    scrollToBottom()
  } finally {
    loading.value = false
  }
}

function scrollToBottom() {
  nextTick(() => {
    scrollId.value = ''
    nextTick(() => { scrollId.value = 'msg-bottom' })
  })
}

async function onSend() {
  const text = inputText.value.trim()
  if (!text || !sessionId.value) return
  inputText.value = ''
  messages.value.push({
    id: Date.now(),
    session_id: sessionId.value,
    sender_type: 'staff',
    content: text,
    type: 'text',
    created_at: new Date().toISOString(),
  })
  scrollToBottom()
  await sendImMessage(sessionId.value, text)
}

onLoad((query: any) => {
  sessionId.value = Number(query?.session_id || 0)
  if (sessionId.value) loadMessages()
})
</script>

<template>
  <view class="page">
    <scroll-view class="msg-list" scroll-y :scroll-into-view="scrollId">
      <view v-if="loading" class="loading">加载中...</view>
      <view v-else-if="!messages.length" class="empty">暂无消息</view>
      <ChatBubble v-for="msg in messages" :key="msg.id" :message="msg" />
      <view id="msg-bottom" />
    </scroll-view>

    <view class="input-bar">
      <input class="msg-input" v-model="inputText" placeholder="输入消息..." @confirm="onSend" />
      <up-button size="mini" type="primary" @click="onSend">发送</up-button>
    </view>
  </view>
</template>

<style scoped>
.page { height: 100vh; display: flex; flex-direction: column; background: var(--eapp-bg); }
.msg-list { flex: 1; padding: 20rpx; box-sizing: border-box; display: flex; flex-direction: column; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 40rpx 0; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.input-bar { display: flex; align-items: center; gap: 12rpx; padding: 16rpx 20rpx; background: #fff; border-top: 1px solid var(--eapp-border); padding-bottom: calc(16rpx + env(safe-area-inset-bottom, 0)); }
.msg-input { flex: 1; height: 72rpx; padding: 0 20rpx; border: 1px solid var(--eapp-border); border-radius: 36rpx; font-size: 26rpx; background: #f8fafc; }
</style>
