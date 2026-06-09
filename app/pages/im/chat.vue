<template>
  <view class="bg-gray-50 flex flex-col w-full" style="height: 100vh; overflow: hidden;">
    <u-navbar :title="$t('chat.title')" :placeholder="true" />

    <!-- Queue notice -->
    <view v-if="queuePosition > 0" class="bg-yellow-50 border-b-1 border-yellow-100 px-24rpx py-20rpx flex items-center gap-16rpx">
      <u-icon name="clock" size="20" color="#d97706" />
      <text class="text-yellow-700 text-26rpx flex-1">{{ $t('chat.queueNotice', { position: queuePosition }) }}</text>
    </view>

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
          <!-- System message -->
          <view v-if="m.sender_type === 0" class="flex justify-center mb-20rpx">
            <text class="text-gray-400 text-24rpx bg-gray-100 px-24rpx py-8rpx rounded-full">{{ m.content }}</text>
          </view>
          <!-- Message row -->
          <view v-else :class="m.sender_type === 1 ? 'flex-row-reverse' : ''"
            class="flex items-end gap-12rpx">
            <!-- Avatar -->
            <view :class="m.sender_type === 2 ? 'bg-blue-700' : m.sender_type === 3 ? 'bg-indigo-500' : 'bg-gray-300'"
              class="w-64rpx h-64rpx rounded-full flex-shrink-0 flex items-center justify-center">
              <text class="text-white text-22rpx">{{ m.sender_type === 2 ? $t('chat.serviceAvatar') : m.sender_type === 3 ? $t('chat.aiAvatar') : $t('chat.userAvatar') }}</text>
            </view>
            <!-- Bubble -->
            <view class="flex flex-col" :class="m.sender_type === 1 ? 'items-end' : 'items-start'">
              <text v-if="m.sender_type === 3" class="text-indigo-400 text-20rpx mb-4rpx ml-8rpx">{{ $t('chat.aiName') }}</text>
              <view :class="m.sender_type === 1
                ? 'bg-blue-700 text-white rounded-tl-24rpx rounded-tr-8rpx rounded-bl-24rpx rounded-br-24rpx'
                : 'bg-white text-gray-800 rounded-tl-8rpx rounded-tr-24rpx rounded-bl-24rpx rounded-br-24rpx'"
                class="max-w-450rpx px-28rpx py-18rpx text-28rpx leading-40rpx"
                style="box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.06); word-break: break-all;">
                <image v-if="m.type === 'image'" :src="fileUrl(m)" mode="widthFix" class="chat-image" @click="previewImage(m)" />
                <view v-else-if="m.type === 'file'" class="file-card" @click="openFile(m)">
                  <u-icon name="file-text" size="18" :color="m.sender_type === 1 ? '#ffffff' : '#64748b'" />
                  <text>{{ fileName(m) }}</text>
                </view>
                <template v-else>{{ m.content }}</template>
              </view>
              <!-- Thumbs feedback — only on AI replies -->
              <view v-if="m.sender_type === 3 && !m.rated" class="flex items-center gap-16rpx mt-8rpx ml-8rpx">
                <text class="text-gray-300 text-20rpx">{{ $t('chat.ratePrompt') }}</text>
                <text class="text-24rpx" @click="rate(m, 1)">👍</text>
                <text class="text-24rpx" @click="rate(m, -1)">👎</text>
              </view>
              <view v-if="m.sender_type === 3 && m.rated" class="mt-8rpx ml-8rpx">
                <text class="text-gray-300 text-20rpx">{{ $t('chat.rated') }}</text>
              </view>
            </view>
          </view>
        </view>

        <!-- AI typing indicator -->
        <view v-if="aiTyping" class="flex items-end gap-12rpx mb-20rpx">
          <view class="w-64rpx h-64rpx rounded-full flex-shrink-0 flex items-center justify-center bg-indigo-500">
            <text class="text-white text-22rpx">{{ $t('chat.aiAvatar') }}</text>
          </view>
          <view class="bg-white text-gray-400 rounded-tl-8rpx rounded-tr-24rpx rounded-bl-24rpx rounded-br-24rpx px-28rpx py-18rpx text-28rpx"
            style="box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.06);">{{ $t('chat.aiThinking') }}</view>
        </view>
      </view>
    </scroll-view>

    <!-- Transfer-to-human quick action (only while AI is serving) -->
    <view v-if="mode === 'ai'"
      class="bg-white border-t-1 border-gray-100 px-20rpx py-12rpx flex items-center gap-12rpx"
      style="position: fixed; left: 0; right: 0; z-index: 29;"
      :style="{ bottom: 'calc(96rpx + env(safe-area-inset-bottom))' }">
      <u-icon name="kefu-ermai" size="18" color="#6366f1" />
      <text class="text-gray-400 text-24rpx flex-1">{{ $t('chat.aiHint') }}</text>
      <u-button :text="$t('chat.transferToHuman')" size="mini" shape="circle" type="info" plain
        @click="requestHuman" :custom-style="{ height: '52rpx' }" />
    </view>

    <!-- Input bar -->
    <view
      class="bg-white border-t-1 border-gray-100 px-20rpx py-16rpx flex items-center gap-16rpx"
      style="position: fixed; left: 0; right: 0; bottom: 0; z-index: 30;"
      :style="{paddingBottom: 'calc(16rpx + env(safe-area-inset-bottom))'}"
    >
      <u-button :text="$t('chat.attachment')" size="small" shape="circle" plain
        :loading="uploading"
        @click="chooseImage"
        :custom-style="{width: '120rpx'}" />
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
import { get, post, upload } from '@/utils/request'

const { t } = useI18n()

const messages = ref<any[]>([])
const inputText = ref('')
const scrollTop = ref(0)
const sessionID = ref(0)
const connected = ref(false)
const queuePosition = ref(0)
const mode = ref<'ai' | 'human'>('human') // 'ai' = 大模型应答中
const aiTyping = ref(false)
const uploading = ref(false)

let ws: any = null
let heartbeat: any = null
let reconnectTimer: any = null
let reconnectDelay = 3000

onMounted(async () => {
  const session = await get<any>('/api/v1/im/session')
  if (session) {
    sessionID.value = session.id
    queuePosition.value = session.queue_position || 0
    mode.value = session.mode === 'ai' ? 'ai' : 'human'
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
        aiTyping.value = false
        messages.value.push({
          id: Date.now(),
          sender_type: frame.payload.sender_type ?? 2,
          content: frame.payload.content,
          type: frame.payload.msg_type || 'text',
          extra: frame.payload.extra,
        })
        scrollToBottom()
      } else if (frame.type === 'typing') {
        aiTyping.value = true
        scrollToBottom()
      } else if (frame.type === 'queue') {
        mode.value = 'human'
        aiTyping.value = false
        queuePosition.value = frame.payload?.position || 0
      } else if (frame.type === 'assign') {
        if (frame.payload?.action === 'accepted') {
          mode.value = 'human'
          queuePosition.value = 0
          messages.value.push({
            id: Date.now(),
            sender_type: 0,
            content: t('chat.assignedNotice'),
          })
          scrollToBottom()
        } else if (frame.payload?.action === 'transfer') {
          messages.value.push({
            id: Date.now(),
            sender_type: 0,
            content: t('chat.transferNotice'),
          })
          scrollToBottom()
        }
      } else if (frame.type === 'close') {
        messages.value.push({
          id: Date.now(),
          sender_type: 0,
          content: t('chat.closedNotice'),
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

  messages.value.push({ id: Date.now(), sender_type: 1, content: text, type: 'text' })
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

function parseExtra(message: any) {
  if (!message?.extra) return {}
  if (typeof message.extra === 'object') return message.extra
  try {
    return JSON.parse(message.extra)
  } catch {
    return {}
  }
}

function fileUrl(message: any) {
  const extra = parseExtra(message)
  return extra.file_url || extra.url || message.content || ''
}

function fileName(message: any) {
  const extra = parseExtra(message)
  return extra.file_name || extra.name || message.content || '附件'
}

function previewImage(message: any) {
  const url = fileUrl(message)
  if (url) uni.previewImage({ urls: [url] })
}

function openFile(message: any) {
  const url = fileUrl(message)
  if (!url) return
  // H5 can open directly; native/miniprogram platforms can add file preview later.
  // #ifdef H5
  window.open(url, '_blank')
  // #endif
}

function sendAttachmentFrame(info: any) {
  const extra = {
    file_url: info.url,
    file_path: info.path,
    file_name: info.name,
    file_size: info.size,
    mime: info.mime,
  }
  messages.value.push({
    id: Date.now(),
    sender_type: 1,
    content: info.name,
    type: info.message_type,
    extra,
  })
  scrollToBottom()
  if (ws && connected.value) {
    ws.send({
      data: JSON.stringify({
        type: 'msg',
        session_id: sessionID.value,
        payload: { msg_type: info.message_type, content: info.name, extra }
      })
    })
  }
}

function chooseImage() {
  if (!sessionID.value || uploading.value) return
  uni.chooseImage({
    count: 1,
    success: async (res) => {
      const filePath = res.tempFilePaths?.[0]
      if (!filePath) return
      uploading.value = true
      try {
        const info = await upload<any>('/api/v1/im/upload', filePath, 'file', { session_id: String(sessionID.value) })
        sendAttachmentFrame(info)
      } finally {
        uploading.value = false
      }
    },
  })
}

// requestHuman asks the backend to hand off from AI to a human agent via a
// dedicated frame (locale-independent; the server runs SwitchToHuman and emits
// the system notice + queue/assign frames).
function requestHuman() {
  if (mode.value !== 'ai') return
  aiTyping.value = false
  mode.value = 'human'
  if (ws && connected.value) {
    ws.send({
      data: JSON.stringify({ type: 'to_human', session_id: sessionID.value })
    })
  }
}

// Submit a 👍/👎 rating for an AI message.
async function rate(msg: any, rating: 1 | -1) {
  msg.rated = true
  try {
    await post('/api/v1/im/feedback', {
      session_id: sessionID.value,
      rating,
      query: messages.value.find(m => m.sender_type === 1 && m.id < msg.id)?.content ?? '',
      answer: msg.content,
    })
  } catch { /* best-effort, ignore errors */ }
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

<style scoped>
.chat-image { max-width: 360rpx; border-radius: 16rpx; display: block; }
.file-card { display: flex; align-items: center; gap: 12rpx; min-width: 240rpx; }
</style>
