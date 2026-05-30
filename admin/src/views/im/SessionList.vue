<template>
  <div class="flex" style="height: calc(100vh - 128px);">
    <!-- Left: Session list -->
    <div class="w-80 border-r border-slate-100 flex flex-col bg-white shrink-0">
      <!-- Header -->
      <div class="px-4 py-3 border-b border-slate-100 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <div class="w-2 h-2 rounded-full" :class="wsConnected ? 'bg-green-500' : 'bg-red-400'" />
          <span class="text-sm font-semibold text-slate-700">{{ $t('im.title') }}</span>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-xs text-slate-400">{{ $t('im.sessionCount', { count: sessions.length }) }}</span>
          <!-- Online/Offline toggle -->
          <button @click="toggleOnline"
            :class="staffOnline ? 'bg-green-500 hover:bg-green-600' : 'bg-slate-300 hover:bg-slate-400'"
            class="flex items-center gap-1 px-2 py-1 rounded-full text-white text-xs font-medium transition">
            <span class="w-1.5 h-1.5 rounded-full bg-white/80" />
            {{ staffOnline ? $t('im.online') : $t('im.offline') }}
          </button>
        </div>
      </div>
      <!-- Staff load indicator -->
      <div v-if="staffOnline" class="px-4 py-2 bg-green-50 border-b border-green-100 flex items-center justify-between">
        <span class="text-xs text-green-700">{{ $t('im.serving') }}</span>
        <span class="text-xs font-semibold text-green-700">{{ staffStatus.current_load }} / {{ staffStatus.max_load }}</span>
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
              <div>
                <span class="text-sm font-medium text-slate-800">{{ $t('im.userLabel', { id: s.user_id }) }}</span>
                <span v-if="s.queue_position > 0" class="ml-1 text-xs text-yellow-600">排队第{{ s.queue_position }}位</span>
              </div>
            </div>
            <span :class="statusClass(s.status)" class="text-xs px-1.5 py-0.5 rounded-full">
              {{ statusLabel(s.status) }}
            </span>
          </div>
          <div class="flex items-center justify-between pl-10">
            <span class="text-xs text-slate-400 truncate flex-1 mr-2">{{ s.last_msg || $t('im.noMessages') }}</span>
            <div class="flex items-center gap-1.5 shrink-0">
              <span class="text-xs text-slate-300">{{ formatTime(s.updated_at) }}</span>
              <span v-if="s.unread_count > 0"
                class="min-w-4.5 h-4.5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center px-1 font-medium">
                {{ s.unread_count }}
              </span>
            </div>
          </div>
        </div>
        <div v-if="!sessions.length" class="text-center py-12 text-slate-400 text-sm">{{ $t('im.noSessions') }}</div>
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
            <span class="text-sm">{{ $t('im.selectSession') }}</span>
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
              <span class="text-sm font-semibold text-slate-800">{{ $t('im.userLabel', { id: activeSession.user_id }) }}</span>
              <span class="text-xs text-slate-400 ml-2">{{ $t('im.sessionLabel', { id: activeSession.id }) }}</span>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <span :class="statusClass(activeSession.status)" class="text-xs px-2 py-1 rounded-full">
              {{ statusLabel(activeSession.status) }}
            </span>
            <button v-if="activeSession.status === 1" @click="acceptSession"
              class="text-xs px-3 py-1 bg-green-500 text-white rounded-full hover:bg-green-600 transition">
              接入
            </button>
            <button v-if="activeSession.status === 2" @click="showTransferModal = true"
              class="text-xs px-3 py-1 bg-blue-500 text-white rounded-full hover:bg-blue-600 transition">
              转接
            </button>
            <button v-if="activeSession.status === 2" @click="closeSession"
              class="text-xs px-3 py-1 bg-slate-200 text-slate-600 rounded-full hover:bg-slate-300 transition">
              结束
            </button>
          </div>
        </div>

        <!-- Messages -->
        <div ref="msgContainer" class="flex-1 overflow-y-auto px-5 py-4">
          <div v-for="m in messages" :key="m.id" class="mb-4">
            <!-- System message -->
            <div v-if="m.sender_type === 0" class="flex justify-center">
              <span class="text-xs text-slate-400 bg-slate-100 px-3 py-1 rounded-full">{{ m.content }}</span>
            </div>
            <div v-else :class="m.sender_type === 2 ? 'flex-row-reverse' : ''" class="flex items-end gap-2">
              <div :class="m.sender_type === 2 ? 'bg-red-600' : 'bg-slate-200'"
                class="w-7 h-7 rounded-full flex items-center justify-center shrink-0">
                <span :class="m.sender_type === 2 ? 'text-white' : 'text-slate-600'" class="text-xs font-medium">
                  {{ m.sender_type === 2 ? $t('im.avatarStaff') : $t('im.avatarUser') }}
                </span>
              </div>
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
            :placeholder="$t('im.inputPlaceholder')"
            :disabled="activeSession.status !== 2"
            class="flex-1 border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-400 focus:ring-1 focus:ring-red-400/20 transition disabled:bg-slate-50 disabled:text-slate-400" />
          <button @click="sendReply" :disabled="!replyText.trim() || activeSession.status !== 2"
            class="px-6 py-2.5 bg-red-600 text-white rounded-xl text-sm font-medium hover:bg-red-700 transition disabled:opacity-40">
            {{ $t('im.send') }}
          </button>
        </div>
      </template>
    </div>

    <!-- Transfer Modal -->
    <div v-if="showTransferModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showTransferModal = false">
      <div class="bg-white rounded-2xl p-6 w-96 shadow-xl">
        <h3 class="text-lg font-semibold mb-4">{{ $t('im.transfer') }}</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">目标客服ID *</label>
            <input v-model.number="transferForm.toStaffId" type="number" placeholder="输入客服ID"
              class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">转接备注（可选）</label>
            <textarea v-model="transferForm.remark" placeholder="如：专业问题转技术客服" rows="3"
              class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20 resize-none"></textarea>
          </div>
        </div>
        <div class="flex gap-3 mt-6">
          <button @click="showTransferModal = false"
            class="flex-1 px-4 py-2 border border-slate-200 text-slate-600 rounded-lg text-sm font-medium hover:bg-slate-50 transition">
            取消
          </button>
          <button @click="confirmTransfer"
            class="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg text-sm font-medium hover:bg-blue-600 transition">
            确认转接
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'

const { t } = useI18n()

const sessions = ref<any[]>([])
const messages = ref<any[]>([])
const activeSession = ref<any>(null)
const replyText = ref('')
const msgContainer = ref<HTMLElement>()
const wsConnected = ref(false)
const staffOnline = ref(false)
const staffStatus = ref({ current_load: 0, max_load: 5 })
const showTransferModal = ref(false)
const transferForm = ref({ toStaffId: null as number | null, remark: '' })

// Sound notification
const notifySound = new Audio('data:audio/wav;base64,UklGRl9vT19teleVmH0AAAAQAAABAAAQ//8CABAAD//')
notifySound.volume = 0.3

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
        if (activeSession.value && frame.session_id === activeSession.value.id) {
          messages.value.push({
            id: Date.now(),
            sender_type: frame.payload.sender_type ?? 1,
            content: frame.payload.content,
            created_at: new Date().toISOString(),
          })
          scrollToBottom()
        }
        const sess = sessions.value.find(s => s.id === frame.session_id)
        if (sess) {
          sess.last_msg = frame.payload.content
          sess.updated_at = new Date().toISOString()
          if (activeSession.value?.id !== frame.session_id) {
            sess.unread_count = (sess.unread_count || 0) + 1
          }
        }
        try { notifySound.play() } catch {}
      } else if (frame.type === 'assign') {
        const action = frame.payload?.action
        if (action === 'new') {
          // New session assigned — reload list
          loadSessions()
          loadStaffStatus()
        } else if (action === 'transfer_out') {
          // Remove session from list
          sessions.value = sessions.value.filter(s => s.id !== frame.session_id)
          if (activeSession.value?.id === frame.session_id) activeSession.value = null
          loadStaffStatus()
        }
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

async function toggleOnline() {
  const newState = !staffOnline.value
  await request.post('/im/staff/online', { online: newState })
  staffOnline.value = newState
  await loadStaffStatus()
  if (newState) loadSessions()
}

async function loadStaffStatus() {
  const data: any = await request.get('/im/staff/status')
  if (data) {
    staffStatus.value = data
    staffOnline.value = data.is_online === 1
  }
}

async function acceptSession() {
  if (!activeSession.value) return
  await request.post(`/im/sessions/${activeSession.value.id}/accept`)
  activeSession.value.status = 2
  const sess = sessions.value.find(s => s.id === activeSession.value.id)
  if (sess) { sess.status = 2; sess.queue_position = 0 }
  await loadStaffStatus()
}

async function closeSession() {
  if (!activeSession.value) return
  await request.post(`/im/sessions/${activeSession.value.id}/close`)
  activeSession.value.status = 3
  const sess = sessions.value.find(s => s.id === activeSession.value.id)
  if (sess) sess.status = 3
  await loadStaffStatus()
}

async function confirmTransfer() {
  if (!transferForm.value.toStaffId) {
    alert('请输入目标客服ID')
    return
  }
  if (!activeSession.value) return
  try {
    await request.post(`/im/sessions/${activeSession.value.id}/transfer`, {
      to_staff_id: transferForm.value.toStaffId,
      remark: transferForm.value.remark
    })
    showTransferModal.value = false
    transferForm.value = { toStaffId: null, remark: '' }
    sessions.value = sessions.value.filter(s => s.id !== activeSession.value.id)
    activeSession.value = null
    await loadStaffStatus()
  } catch (err: any) {
    alert(err.message || '转接失败')
  }
}

async function submitTransfer() {
  // Deprecated - use confirmTransfer instead
  return confirmTransfer()
}

const statusLabels: Record<number, string> = { 1: t('im.statusWaiting'), 2: t('im.statusActive'), 3: t('im.statusClosed') }
const statusClasses: Record<number, string> = {
  1: 'bg-yellow-50 text-yellow-600',
  2: 'bg-green-50 text-green-600',
  3: 'bg-slate-100 text-slate-400',
}
const statusLabel = (s: number) => statusLabels[s] || ''
const statusClass = (s: number) => statusClasses[s] || ''

function formatTime(ts: string) {
  if (!ts) return ''
  const d = new Date(ts)
  const now = new Date()
  if (d.toDateString() === now.toDateString()) return ts.slice(11, 16)
  return ts.slice(5, 10)
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
  loadStaffStatus()
  connectWS()
})

onUnmounted(() => {
  if (reconnectTimer) clearTimeout(reconnectTimer)
  ws?.close()
})
</script>

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
              <span class="text-sm font-medium text-slate-800">{{ $t('im.userLabel', { id: s.user_id }) }}</span>
            </div>
            <span :class="statusClass(s.status)" class="text-xs px-1.5 py-0.5 rounded-full">
              {{ statusLabel(s.status) }}
            </span>
          </div>
          <div class="flex items-center justify-between pl-10">
            <span class="text-xs text-slate-400 truncate flex-1 mr-2">{{ s.last_msg || $t('im.noMessages') }}</span>
            <div class="flex items-center gap-1.5 shrink-0">
              <span class="text-xs text-slate-300">{{ formatTime(s.updated_at) }}</span>
              <span v-if="s.unread_count > 0"
                class="min-w-4.5 h-4.5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center px-1 font-medium">
                {{ s.unread_count }}
              </span>
            </div>
          </div>
        </div>
        <div v-if="!sessions.length" class="text-center py-12 text-slate-400 text-sm">{{ $t('im.noSessions') }}</div>
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
            <span class="text-sm">{{ $t('im.selectSession') }}</span>
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
              <span class="text-sm font-semibold text-slate-800">{{ $t('im.userLabel', { id: activeSession.user_id }) }}</span>
              <span class="text-xs text-slate-400 ml-2">{{ $t('im.sessionLabel', { id: activeSession.id }) }}</span>
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
                  {{ m.sender_type === 2 ? $t('im.avatarStaff') : $t('im.avatarUser') }}
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
            :placeholder="$t('im.inputPlaceholder')"
            class="flex-1 border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-400 focus:ring-1 focus:ring-red-400/20 transition" />
          <button @click="sendReply" :disabled="!replyText.trim()"
            class="px-6 py-2.5 bg-red-600 text-white rounded-xl text-sm font-medium hover:bg-red-700 transition disabled:opacity-40">
            {{ $t('im.send') }}
          </button>
        </div>
      </template>
    </div>

    <!-- Transfer Modal -->
    <div v-if="showTransferModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showTransferModal = false">
      <div class="bg-white rounded-2xl p-6 w-96 max-w-[90vw]">
        <h3 class="text-lg font-semibold mb-4">转接会话</h3>
        <div class="mb-4">
          <label class="block text-sm text-slate-600 mb-2">目标客服 ID</label>
          <input v-model="transferForm.toStaffId" type="number" placeholder="输入客服 ID"
            class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20" />
        </div>
        <div class="mb-4">
          <label class="block text-sm text-slate-600 mb-2">转接备注（可选）</label>
          <textarea v-model="transferForm.remark" placeholder="例如：专业问题转技术客服"
            class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20"
            rows="3" />
        </div>
        <div class="flex gap-2 justify-end">
          <button @click="showTransferModal = false"
            class="px-4 py-2 text-sm text-slate-600 hover:bg-slate-100 rounded-lg transition">
            取消
          </button>
          <button @click="submitTransfer" :disabled="!transferForm.toStaffId"
            class="px-4 py-2 text-sm bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition disabled:opacity-40 disabled:cursor-not-allowed">
            确认转接
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'

const { t } = useI18n()

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

const statusLabels: Record<number, string> = { 1: t('im.statusWaiting'), 2: t('im.statusActive'), 3: t('im.statusClosed') }
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
