<template>
  <view class="bg-gray-50 flex flex-col" style="height: 100vh; overflow: hidden;">
    <u-navbar :title="$t('chat.title')" :placeholder="true" />

    <!-- Messages -->
    <scroll-view
      scroll-y
      class="flex-1 px-20rpx py-16rpx"
      :scroll-top="scrollTop"
      scroll-with-animation
      :style="{ paddingBottom: '180rpx' }"
    >
      <view class="pb-20rpx">
        <!-- Welcome -->
        <view v-if="!messages.length" class="text-center py-80rpx">
          <view class="w-120rpx h-120rpx rounded-full bg-blue-50 flex items-center justify-center mx-auto mb-20rpx">
            <u-icon name="chat" size="40" color="#1e40af" />
          </view>
          <text class="text-gray-400 text-28rpx">{{ $t('chat.greeting') }}</text>
        </view>

        <view v-for="m in messages" :key="m.id" class="mb-20rpx">
          <!-- Message row -->
          <view :class="m.sender_type === 1 ? 'flex-row-reverse' : ''"
            class="flex items-end gap-12rpx">
            <!-- Avatar -->
            <view :class="m.sender_type === 2 ? 'bg-blue-700' : 'bg-gray-300'"
              class="w-64rpx h-64rpx rounded-full flex-shrink-0 flex items-center justify-center">
              <text class="text-white text-22rpx">{{ m.sender_type === 2 ? $t('chat.serviceAvatar') : $t('chat.userAvatar') }}</text>
            </view>
            <!-- Bubble -->
            <view :class="m.sender_type === 1
              ? 'bg-blue-700 text-white rounded-tl-24rpx rounded-tr-8rpx rounded-bl-24rpx rounded-br-24rpx'
              : 'bg-white text-gray-800 rounded-tl-8rpx rounded-tr-24rpx rounded-bl-24rpx rounded-br-24rpx'"
              class="max-w-450rpx px-28rpx py-18rpx text-28rpx leading-40rpx"
              style="box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.06); word-break: break-all;">
              {{ m.content }}
            </view>
          </view>
        </view>
      </view>
    </scroll-view>

    <!-- Input bar -->
    <view
      class="bg-white border-t-1 border-gray-100 px-20rpx py-16rpx flex items-center gap-16rpx"
      style="position: fixed; left: 0; right: 0; bottom: 0; z-index: 30;"
      :style="{paddingBottom: 'calc(16rpx + env(safe-area-inset-bottom))'}"
    >
      <view class="flex-1">
        <u-input v-model="inputText" :placeholder="$t('chat.inputPlaceholder')" border="surround" shape="circle"
          :maxlength="500"
          @confirm="sendMsg"
          confirmType="send"
          clearable
        />
      </view>
      <u-button type="primary" :text="$t('chat.send')" size="small" shape="circle"
        :disabled="!inputText.trim()" @click="sendMsg"
        :custom-style="{width: '120rpx'}" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { get } from '@/utils/request'

const { t } = useI18n()

const messages = ref<any[]>([])
const inputText = ref('')
const scrollTop = ref(0)
const sessionID = ref(0)
const connected = ref(false)

let ws: any = null
let heartbeat: any = null
let reconnectTimer: any = null
let reconnectDelay = 3000

onMounted(async () => {
  const session = await get<any>('/api/v1/im/session')
  if (session) {
    sessionID.value = session.id
    const data = await get<any>('/api/v1/im/messages', { session_id: session.id, size: 50 })
    messages.value = (data?.list || []).reverse()
    scrollToBottom()
  }

  const token = uni.getStorageSync('user_token')
  if (token) connectWS(token)
})

function connectWS(token: string) {
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
        messages.value.push({
          id: Date.now(),
          sender_type: frame.payload.sender_type || 2,
          content: frame.payload.content,
        })
        scrollToBottom()
      }
    } catch {}
  })

  ws.onClose(() => {
    connected.value = false
    scheduleReconnect(token)
  })

  ws.onError(() => {
    connected.value = false
  })

  if (heartbeat) clearInterval(heartbeat)
  heartbeat = setInterval(() => {
    ws?.send({ data: JSON.stringify({ type: 'ping' }) })
  }, 30000)
}

function scheduleReconnect(token: string) {
  if (reconnectTimer) clearTimeout(reconnectTimer)
  reconnectTimer = setTimeout(() => {
    reconnectDelay = Math.min(reconnectDelay * 2, 30000)
    connectWS(token)
  }, reconnectDelay)
}

function scrollToBottom() {
  nextTick(() => { scrollTop.value = 999999 })
}

function sendMsg() {
  const text = inputText.value.trim()
  if (!text) return

  messages.value.push({ id: Date.now(), sender_type: 1, content: text })
  inputText.value = ''
  scrollToBottom()

  if (ws && connected.value) {
    ws.send({
      data: JSON.stringify({
        type: 'msg',
        session_id: sessionID.value,
        payload: { msg_type: 'text', content: text }
      })
    })
    return
  }
  scheduleLocalReply(text)
}

function scheduleLocalReply(content: string) {
  const replies = [
    t('chat.autoReply1'),
    t('chat.autoReply2'),
    t('chat.autoReply3'),
    t('chat.autoReply4', { topic: content.slice(0, 8) }),
  ]
  setTimeout(() => {
    messages.value.push({
      id: Date.now() + 1,
      sender_type: 2,
      content: replies[Math.floor(Math.random() * replies.length)],
    })
    scrollToBottom()
  }, 400)
}

onUnmounted(() => {
  if (heartbeat) clearInterval(heartbeat)
  if (reconnectTimer) clearTimeout(reconnectTimer)
  ws?.close({})
})
</script>
