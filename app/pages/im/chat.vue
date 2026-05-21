<template>
  <view class="bg-gray-50 min-h-screen flex flex-col" style="height: 100vh;">
    <u-navbar title="在线客服" :placeholder="true" />

    <!-- Messages -->
    <scroll-view scroll-y class="flex-1 px-3 py-2" :scroll-top="scrollTop" scroll-with-animation>
      <view class="space-y-3 pb-4">
        <view v-for="m in messages" :key="m.id"
          :class="m.sender_type === 2 ? 'flex-row-reverse' : ''"
          class="flex gap-2 items-end">
          <!-- Avatar -->
          <view :class="m.sender_type === 2 ? 'bg-blue-700' : 'bg-gray-300'"
            class="w-8 h-8 rounded-full shrink-0 flex items-center justify-center">
            <text class="text-white text-xs">{{ m.sender_type === 2 ? '客' : '我' }}</text>
          </view>
          <!-- Bubble -->
          <view :class="m.sender_type === 2 ? 'bg-blue-700 text-white' : 'bg-white text-slate-800'"
            class="max-w-56 px-4 py-2 rounded-2xl shadow-sm text-sm">
            <text>{{ m.content }}</text>
          </view>
        </view>

        <view v-if="!messages.length" class="text-center py-8">
          <text class="text-gray-400 text-sm">您好！有什么可以帮您的？</text>
        </view>
      </view>
    </scroll-view>

    <!-- Input bar -->
    <view class="bg-white border-t border-gray-100 px-3 py-3 flex gap-2 items-center">
      <input
        v-model="inputText"
        placeholder="输入消息..."
        class="flex-1 bg-gray-100 rounded-full px-4 py-2 text-sm"
        @confirm="sendMsg"
        confirm-type="send"
      />
      <u-button type="primary" size="mini" text="发送" @click="sendMsg" :disabled="!inputText.trim()" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { get, post } from '@/utils/request'

const messages = ref<any[]>([])
const inputText = ref('')
const scrollTop = ref(0)
const sessionID = ref(0)

let ws: any = null
let heartbeat: any = null

onMounted(async () => {
  // Get or create session
  const session = await get<any>('/api/v1/im/session')
  if (session) {
    sessionID.value = session.id
    // Load history
    const data = await get<any>('/api/v1/im/messages', { session_id: session.id, size: 50 })
    messages.value = (data?.list || []).reverse()
    scrollToBottom()
  }

  // Connect WebSocket
  const token = uni.getStorageSync('user_token')
  if (token) {
    connectWS(token)
  }
})

function connectWS(token: string) {
  ws = uni.connectSocket({
    url: `ws://localhost:8080/ws/im?token=${token}`,
    complete: () => {}
  })

  ws.onMessage((res: any) => {
    try {
      const frame = JSON.parse(res.data)
      if (frame.type === 'msg') {
        messages.value.push({
          id: Date.now(),
          sender_type: frame.payload.sender_type || 2,
          content: frame.payload.content,
        })
        scrollToBottom()
      }
    } catch {}
  })

  // Heartbeat
  heartbeat = setInterval(() => {
    ws?.send({ data: JSON.stringify({ type: 'ping' }) })
  }, 30000)
}

function scrollToBottom() {
  nextTick(() => {
    scrollTop.value = 999999
  })
}

async function sendMsg() {
  const text = inputText.value.trim()
  if (!text) return

  messages.value.push({ id: Date.now(), sender_type: 1, content: text })
  inputText.value = ''
  scrollToBottom()

  if (ws) {
    ws.send({
      data: JSON.stringify({
        type: 'msg',
        session_id: sessionID.value,
        payload: { msg_type: 'text', content: text }
      })
    })
  } else {
    // Fallback: HTTP (if WS not available)
    await post('/api/v1/im/message', { session_id: sessionID.value, content: text })
  }
}

onUnmounted(() => {
  if (heartbeat) clearInterval(heartbeat)
  ws?.close({})
})
</script>
