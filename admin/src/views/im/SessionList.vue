<template>
  <div class="flex" style="height: calc(100vh - 128px);">
    <!-- Left: Session list -->
    <div class="w-80 border-r border-slate-100 flex flex-col bg-white shrink-0">
      <!-- Header -->
      <div class="px-4 py-3 border-b border-slate-100 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <div class="w-2 h-2 rounded-full" :class="wsConnected ? 'bg-green-500' : 'bg-red-400'" />
          <span class="text-sm font-semibold text-slate-700">客服会话</span>
        </div>
        <span class="text-xs text-slate-400">{{ sessions.length }} 个会话</span>
      </div>
      <!-- Session list -->
      <div class="flex-1 overflow-y-auto">
        <div v-for="s in sessions" :key="s.id"
          @click="selectSession(s)"
          class="px-4 py-3 cursor-pointer transition border-b border-slate-50 hover:bg-slate-50"
          :class="activeSession?.id === s.id ? 'bg-red-50 border-l-3 border-l-red-600' : ''">
          <div class="flex items-center justify-between mb-1">
            <div class="flex items-center gap-2">
              <div class="w-8 h-8 rounded-full bg-slate-200 flex items-center justify-center shrink-0">
                <span class="text-xs text-slate-600 font-medium">U{{ s.user_id }}</span>
              </div>
              <span class="text-sm font-medium text-slate-800">用户 #{{ s.user_id }}</span>
            </div>
            <span :class="statusClass(s.status)" class="text-xs px-1.5 py-0.5 rounded-full">
              {{ statusLabel(s.status) }}
            </span>
          </div>
          <div class="flex items-center justify-between pl-10">
            <span class="text-xs text-slate-400 truncate flex-1 mr-2">{{ s.last_msg || '暂无消息' }}</span>
            <div class="flex items-center gap-1.5 shrink-0">
              <span class="text-xs text-slate-300">{{ formatTime(s.updated_at) }}</span>
              <span v-if="s.unread_count > 0"
                class="min-w-4.5 h-4.5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center px-1 font-medium">
                {{ s.unread_count }}
              </span>
            </div>
          </div>
        </div>
        <div v-if="!sessions.length" class="text-center py-12 text-slate-400 text-sm">暂无会话</div>
      </div>
    </div>

    <!-- Right: Chat area -->
    <div class="flex-1 flex flex-col bg-white overflow-hidden">
      <template v-if="!activeSession">
        <div class="flex-1 flex items-center justify-center text-slate-300">
          <div class="text-center">
            <svg class="w-16 h-16 mx-auto mb-3 text-slate-200" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1">
              <path d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
            <span class="text-sm">选择左侧会话开始回复</span>
          </div>
        </div>
      </template>
      <template v-else>
        <!-- Chat header -->
        <div class="px-5 py-3 border-b border-slate-100 flex items-center justify-between shrink-0">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-full bg-red-100 flex items-center justify-center">
              <span class="text-xs text-red-600 font-bold">U{{ activeSession.user_id }}</span>
            </div>
            <div>
              <span class="text-sm font-semibold text-slate-800">用户 #{{ activeSession.user_id }}</span>
              <span class="text-xs text-slate-400 ml-2">会话 #{{ activeSession.id }}</span>
            </div>
          </div>
          <span :class="statusClass(activeSession.status)" class="text-xs px-2 py-1 rounded-full">
            {{ statusLabel(activeSession.status) }}
          </span>
        </div>

        <!-- Messages -->
        <div ref="msgContainer" class="flex-1 overflow-y-auto px-5 py-4">
          <div v-for="m in messages" :key="m.id" class="mb-4">
            <div :class="m.sender_type === 2 ? 'flex-row-reverse' : ''" class="flex items-end gap-2">
              <!-- Avatar -->
              <div :class="m.sender_type === 2 ? 'bg-red-600' : 'bg-slate-200'"
                class="w-7 h-7 rounded-full flex items-center justify-center shrink-0">
                <span :class="m.sender_type === 2 ? 'text-white' : 'text-slate-600'" class="text-xs font-medium">
                  {{ m.sender_type === 2 ? '客' : '用' }}
                </span>
              </div>
              <!-- Bubble -->
              <div :class="m.sender_type === 2
                ? 'bg-red-600 text-white rounded-tl-2xl rounded-tr-sm rounded-bl-2xl rounded-br-2xl'
                : 'bg-slate-100 text-slate-800 rounded-tl-sm rounded-tr-2xl rounded-bl-2xl rounded-br-2xl'"
                class="max-w-md px-4 py-2.5 text-sm leading-relaxed">
                {{ m.content }}
              </div>
              <span class="text-xs text-slate-300 shrink-0">{{ m.created_at?.slice(11,16) }}</span>
            </div>
          </div>
        </div>

        <!-- Input -->
        <div class="px-5 py-3 border-t border-slate-100 flex gap-3 shrink-0">
          <input v-model="replyText" @keyup.enter="sendReply"
            placeholder="输入回复内容，Enter 发送"
            class="flex-1 border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-400 focus:ring-1 focus:ring-red-400/20 transition" />
          <button @click="sendReply" :disabled="!replyText.trim()"
            class="px-6 py-2.5 bg-red-600 text-white rounded-xl text-sm font-medium hover:bg-red-700 transition disabled:opacity-40">
            发送
          </button>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import request from '@/api/request'

const sessions = ref<any[]>([])
const messages = ref<any[]>([])
const activeSession = ref<any>(null)
const replyText = ref('')
const msgContainer = ref<HTMLElement>()
const wsConnected = ref(false)

// Sound notification
const notifySound = new Audio('data:audio/wav;base64,UklGRl9vT19teleVmH0AAAAQAAABAAAQ//8CABAAD//')
notifySound.volume = 0.3

// WebSocket reconnection
let ws: WebSocket | null = null
let reconnectTimer: any = null
let reconnectDelay = 3000

function connectWS() {
  const token = localStorage.getItem('admin_token')
  if (!token) return

  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = location.hostname + ':8080'
  ws = new WebSocket(`${protocol}//${host}/ws/im?token=${token}`)

  ws.onopen = () => {
    wsConnected.value = true
    reconnectDelay = 3000
  }

  ws.onmessage = (e) => {
    try {
      const frame = JSON.parse(e.data)
      if (frame.type === 'msg') {
        // Add to messages if current session
        if (activeSession.value && frame.session_id === activeSession.value.id) {
          messages.value.push({
            id: Date.now(),
            sender_type: frame.payload.sender_type || 1,
            content: frame.payload.content,
            created_at: new Date().toISOString(),
          })
          scrollToBottom()
        }
        // Update session list
        const sess = sessions.value.find(s => s.id === frame.session_id)
        if (sess) {
          sess.last_msg = frame.payload.content
          sess.updated_at = new Date().toISOString()
          if (activeSession.value?.id !== frame.session_id) {
            sess.unread_count = (sess.unread_count || 0) + 1
          }
        }
        // Play sound
        try { notifySound.play() } catch {}
      }
    } catch {}
  }

  ws.onclose = () => {
    wsConnected.value = false
    scheduleReconnect()
  }

  ws.onerror = () => {
    wsConnected.value = false
  }
}

function scheduleReconnect() {
  if (reconnectTimer) clearTimeout(reconnectTimer)
  reconnectTimer = setTimeout(() => {
    reconnectDelay = Math.min(reconnectDelay * 2, 30000)
    connectWS()
  }, reconnectDelay)
}

const statusLabels: Record<number, string> = { 1: '等待', 2: '进行中', 3: '已关闭' }
const statusClasses: Record<number, string> = {
  1: 'bg-yellow-50 text-yellow-600',
  2: 'bg-green-50 text-green-600',
  3: 'bg-slate-100 text-slate-400',
}
const statusLabel = (s: number) => statusLabels[s] || ''
const statusClass = (s: number) => statusClasses[s] || ''

function formatTime(t: string) {
  if (!t) return ''
  const d = new Date(t)
  const now = new Date()
  if (d.toDateString() === now.toDateString()) return t.slice(11,16)
  return t.slice(5,10)
}

async function loadSessions() {
  const data: any = await request.get('/im/sessions')
  sessions.value = data || []
}

async function selectSession(s: any) {
  activeSession.value = s
  s.unread_count = 0
  const data: any = await request.get(`/im/sessions/${s.id}/messages`, { params: { size: 100 } })
  messages.value = (data?.list || []).reverse()
  await scrollToBottom()
}

async function sendReply() {
  if (!replyText.value.trim() || !activeSession.value) return
  const text = replyText.value
  replyText.value = ''

  await request.post(`/im/sessions/${activeSession.value.id}/reply`, { content: text })
  messages.value.push({
    id: Date.now(), sender_type: 2, content: text,
    created_at: new Date().toISOString()
  })

  // Update session last msg
  const sess = sessions.value.find(s => s.id === activeSession.value.id)
  if (sess) sess.last_msg = text

  await scrollToBottom()
}

async function scrollToBottom() {
  await nextTick()
  if (msgContainer.value) msgContainer.value.scrollTop = msgContainer.value.scrollHeight
}

onMounted(() => {
  loadSessions()
  connectWS()
})

onUnmounted(() => {
  if (reconnectTimer) clearTimeout(reconnectTimer)
  ws?.close()
})
</script>
