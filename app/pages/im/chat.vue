<template>
  <view class="bg-gray-50 flex flex-col" style="height: 100vh;">
    <u-navbar title="在线客服" :placeholder="true" />

    <!-- Messages -->
    <scroll-view scroll-y class="flex-1 px-20rpx py-16rpx" :scroll-top="scrollTop" scroll-with-animation>
      <view class="pb-20rpx">
        <!-- Welcome -->
        <view v-if="!messages.length" class="text-center py-80rpx">
          <view class="w-120rpx h-120rpx rounded-full bg-blue-50 flex items-center justify-center mx-auto mb-20rpx">
            <u-icon name="chat" size="40" color="#1e40af" />
          </view>
          <text class="text-gray-400 text-28rpx">您好！有什么可以帮您的？</text>
        </view>

        <view v-for="m in messages" :key="m.id" class="mb-20rpx">
          <!-- Message row -->
          <view :class="m.sender_type === 1 ? 'flex-row-reverse' : ''"
            class="flex items-end gap-12rpx">
            <!-- Avatar -->
            <view :class="m.sender_type === 2 ? 'bg-blue-700' : 'bg-gray-300'"
              class="w-64rpx h-64rpx rounded-full flex-shrink-0 flex items-center justify-center">
              <text class="text-white text-22rpx">{{ m.sender_type === 2 ? '客' : '我' }}</text>
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
    <view class="bg-white border-t-1 border-gray-100 px-20rpx py-16rpx flex items-center gap-16rpx"
      :style="{paddingBottom: 'calc(16rpx + env(safe-area-inset-bottom))'}">
      <view class="flex-1">
        <u-input v-model="inputText" placeholder="输入消息..." border="surround" shape="circle"
          @confirm="sendMsg" confirmType="send" />
      </view>
      <u-button type="primary" text="发送" size="small" shape="circle"
        :disabled="!inputText.trim()" @click="sendMsg"
        :custom-style="{width: '120rpx'}" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { get } from '@/utils/request'

const messages = ref<any[]>([])
const inputText = ref('')
const scrollTop = ref(0)
const sessionID = ref(0)

let ws: any = null
let heartbeat: any = null

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

  heartbeat = setInterval(() => {
    ws?.send({ data: JSON.stringify({ type: 'ping' }) })
  }, 30000)
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

  if (ws) {
    ws.send({
      data: JSON.stringify({
        type: 'msg',
        session_id: sessionID.value,
        payload: { msg_type: 'text', content: text }
      })
    })
  }
}

onUnmounted(() => {
  if (heartbeat) clearInterval(heartbeat)
  ws?.close({})
})
</script>
