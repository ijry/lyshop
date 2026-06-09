<script setup lang="ts">
import { onLoad, onUnload } from '@dcloudio/uni-app'
import { ref, nextTick } from 'vue'
import { getImMessages, replyImSession, sendImMessage, uploadImAttachment } from '@/api/im'
import ChatBubble from '@/components/biz/ChatBubble.vue'

const sessionId = ref(0)
const messages = ref<any[]>([])
const inputText = ref('')
const loading = ref(false)
const scrollId = ref('')
const connected = ref(false)

let ws: any = null
let heartbeat: any = null
let reconnectTimer: any = null
let reconnectDelay = 3000

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

function connectWS() {
  const token = uni.getStorageSync('eapp_token')
  if (!token) return

  ws = uni.connectSocket({
    url: `ws://localhost:8080/ws/im?token=${token}`,
    complete: () => {}
  })

  ws.onOpen(() => {
    connected.value = true
    reconnectDelay = 3000
  })

  ws.onMessage((res: any) => {
    try {
      const frame = JSON.parse(res.data)
      if (frame.type === 'msg') {
        if (frame.session_id === sessionId.value) {
          messages.value.push({
            id: Date.now(),
            sender_type: frame.payload.sender_type ?? 1,
            content: frame.payload.content,
            type: frame.payload.msg_type || 'text',
            extra: frame.payload.extra,
            created_at: new Date().toISOString(),
          })
          scrollToBottom()
        }
      } else if (frame.type === 'assign') {
        if (frame.session_id === sessionId.value) {
          const action = frame.payload?.action
          if (action === 'transfer_out') {
            messages.value.push({
              id: Date.now(),
              sender_type: 0,
              content: '会话已转接给其他客服',
              type: 'system',
              created_at: new Date().toISOString(),
            })
            scrollToBottom()
          } else if (action === 'transfer_in') {
            messages.value.push({
              id: Date.now(),
              sender_type: 0,
              content: `会话已转接进来（来自客服 ${frame.payload.from_staff_id || ''}）`,
              type: 'system',
              created_at: new Date().toISOString(),
            })
            scrollToBottom()
          }
        }
      }
    } catch {}
  })

  ws.onClose(() => {
    connected.value = false
    scheduleReconnect()
  })

  ws.onError(() => {
    connected.value = false
  })

  if (heartbeat) clearInterval(heartbeat)
  heartbeat = setInterval(() => {
    ws?.send({ data: JSON.stringify({ type: 'ping' }) })
  }, 30000)
}

function scheduleReconnect() {
  if (reconnectTimer) clearTimeout(reconnectTimer)
  reconnectTimer = setTimeout(() => {
    reconnectDelay = Math.min(reconnectDelay * 2, 30000)
    connectWS()
  }, reconnectDelay)
}

async function onSend() {
  const text = inputText.value.trim()
  if (!text || !sessionId.value) return
  inputText.value = ''

  messages.value.push({
    id: Date.now(),
    session_id: sessionId.value,
    sender_type: 2,
    content: text,
    type: 'text',
    created_at: new Date().toISOString(),
  })
  scrollToBottom()

  if (ws && connected.value) {
    ws.send({
      data: JSON.stringify({
        type: 'msg',
        session_id: sessionId.value,
        payload: { msg_type: 'text', content: text }
      })
    })
  } else {
    await sendImMessage(sessionId.value, text)
  }
}

async function chooseImage() {
  if (!sessionId.value) return
  uni.chooseImage({
    count: 1,
    success: async (res) => {
      const filePath = res.tempFilePaths?.[0]
      if (!filePath) return
      const info: any = await uploadImAttachment(sessionId.value, filePath)
      const extra = {
        file_url: info.url,
        file_path: info.path,
        file_name: info.name,
        file_size: info.size,
        mime: info.mime,
      }
      await replyImSession(sessionId.value, {
        type: info.message_type,
        content: info.name,
        extra: JSON.stringify(extra),
      })
      messages.value.push({
        id: Date.now(),
        session_id: sessionId.value,
        sender_type: 2,
        content: info.name,
        type: info.message_type,
        extra,
        created_at: new Date().toISOString(),
      })
      scrollToBottom()
    },
  })
}

onLoad((query: any) => {
  sessionId.value = Number(query?.session_id || 0)
  if (sessionId.value) {
    loadMessages()
    connectWS()
  }
})

onUnload(() => {
  if (heartbeat) clearInterval(heartbeat)
  if (reconnectTimer) clearTimeout(reconnectTimer)
  ws?.close({})
})
</script>

<template>
  <view class="page">
    <view v-if="connected" class="status-bar online">
      <view class="status-dot" />
      <text class="status-text">已连接</text>
    </view>
    <view v-else class="status-bar offline">
      <view class="status-dot" />
      <text class="status-text">连接中...</text>
    </view>

    <scroll-view class="msg-list" scroll-y :scroll-into-view="scrollId">
      <view v-if="loading" class="loading">加载中...</view>
      <view v-else-if="!messages.length" class="empty">暂无消息</view>
      <ChatBubble v-for="msg in messages" :key="msg.id" :message="msg" />
      <view id="msg-bottom" />
    </scroll-view>

    <view class="input-bar">
      <input class="msg-input" v-model="inputText" placeholder="输入消息..." @confirm="onSend" />
      <up-button size="mini" type="info" plain @click="chooseImage">图片</up-button>
      <up-button size="mini" type="primary" @click="onSend" :disabled="!inputText.trim()">发送</up-button>
    </view>
  </view>
</template>

<style scoped>
.page { height: 100vh; display: flex; flex-direction: column; background: var(--eapp-bg); }
.status-bar { display: flex; align-items: center; gap: 8rpx; padding: 12rpx 20rpx; font-size: 22rpx; }
.status-bar.online { background: #ecfdf5; border-bottom: 1rpx solid #a7f3d0; }
.status-bar.offline { background: #fef3c7; border-bottom: 1rpx solid #fde68a; }
.status-dot { width: 12rpx; height: 12rpx; border-radius: 50%; }
.status-bar.online .status-dot { background: #10b981; }
.status-bar.offline .status-dot { background: #f59e0b; }
.status-bar.online .status-text { color: #065f46; }
.status-bar.offline .status-text { color: #92400e; }
.msg-list { flex: 1; padding: 20rpx; box-sizing: border-box; display: flex; flex-direction: column; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 40rpx 0; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.input-bar { display: flex; align-items: center; gap: 12rpx; padding: 16rpx 20rpx; background: #fff; border-top: 1px solid var(--eapp-border); padding-bottom: calc(16rpx + env(safe-area-inset-bottom, 0)); }
.msg-input { flex: 1; height: 72rpx; padding: 0 20rpx; border: 1px solid var(--eapp-border); border-radius: 36rpx; font-size: 26rpx; background: #f8fafc; }
</style>
