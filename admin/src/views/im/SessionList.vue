<template>
  <div class="flex h-full gap-0" style="height: calc(100vh - 128px)">
    <!-- Session list -->
    <div class="w-64 border-r border-slate-100 flex flex-col bg-white rounded-l-xl overflow-hidden">
      <div class="px-4 py-3 border-b border-slate-100">
        <h3 class="font-semibold text-slate-700 text-sm">客服会话</h3>
      </div>
      <div class="flex-1 overflow-y-auto">
        <div v-for="s in sessions" :key="s.id"
          @click="selectSession(s)"
          :class="activeSession?.id === s.id ? 'bg-blue-50 border-l-2 border-blue-700' : 'hover:bg-slate-50'"
          class="px-4 py-3 cursor-pointer transition">
          <div class="flex justify-between">
            <span class="text-sm font-medium text-slate-700">用户 #{{ s.user_id }}</span>
            <span :class="statusClass(s.status)" class="text-xs px-2 py-0.5 rounded-full">
              {{ statusLabel(s.status) }}
            </span>
          </div>
          <p class="text-xs text-slate-400 mt-1 truncate">{{ s.last_msg || '暂无消息' }}</p>
        </div>
        <div v-if="!sessions.length" class="text-center py-8 text-slate-400 text-sm">暂无会话</div>
      </div>
    </div>

    <!-- Message area -->
    <div class="flex-1 flex flex-col bg-white rounded-r-xl overflow-hidden">
      <div v-if="!activeSession" class="flex-1 flex items-center justify-center text-slate-400">
        <span>← 选择一个会话开始回复</span>
      </div>
      <template v-else>
        <div class="px-4 py-3 border-b border-slate-100 flex items-center justify-between">
          <span class="font-medium text-slate-700">用户 #{{ activeSession.user_id }}</span>
          <span class="text-xs text-slate-400">会话 #{{ activeSession.id }}</span>
        </div>
        <div ref="msgContainer" class="flex-1 overflow-y-auto p-4 space-y-3">
          <div v-for="m in messages" :key="m.id"
            :class="m.sender_type === 2 ? 'flex-row-reverse' : ''"
            class="flex gap-2 items-end">
            <div :class="m.sender_type === 2 ? 'bg-blue-700 text-white' : 'bg-slate-100 text-slate-800'"
              class="max-w-xs px-4 py-2 rounded-2xl text-sm leading-relaxed">
              {{ m.content }}
            </div>
            <span class="text-xs text-slate-300 shrink-0">{{ m.created_at?.slice(11,16) }}</span>
          </div>
        </div>
        <div class="p-4 border-t border-slate-100 flex gap-3">
          <input v-model="replyText" @keyup.enter="sendReply"
            placeholder="输入回复内容，Enter 发送"
            class="flex-1 border border-slate-200 rounded-xl px-4 py-2 text-sm focus:outline-none focus:border-blue-400" />
          <button @click="sendReply"
            class="px-5 py-2 bg-blue-700 text-white rounded-xl text-sm hover:bg-blue-600 transition">发送</button>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import request from '@/api/request'

const sessions = ref<any[]>([])
const messages = ref<any[]>([])
const activeSession = ref<any>(null)
const replyText = ref('')
const msgContainer = ref<HTMLElement>()

const statusLabels: Record<number, string> = { 1: '等待', 2: '进行中', 3: '已关闭' }
const statusColors: Record<number, string> = {
  1: 'bg-yellow-50 text-yellow-600',
  2: 'bg-green-50 text-green-600',
  3: 'bg-slate-100 text-slate-400',
}
const statusLabel = (s: number) => statusLabels[s] || ''
const statusClass = (s: number) => statusColors[s] || ''

async function loadSessions() {
  const data: any = await request.get('/im/sessions')
  sessions.value = data || []
}

async function selectSession(s: any) {
  activeSession.value = s
  const data: any = await request.get(`/im/sessions/${s.id}/messages`, { params: { size: 50 } })
  messages.value = (data.list || []).reverse()
  await nextTick()
  if (msgContainer.value) msgContainer.value.scrollTop = msgContainer.value.scrollHeight
}

async function sendReply() {
  if (!replyText.value.trim() || !activeSession.value) return
  await request.post(`/im/sessions/${activeSession.value.id}/reply`, { content: replyText.value })
  messages.value.push({
    id: Date.now(), sender_type: 2, content: replyText.value,
    created_at: new Date().toISOString()
  })
  replyText.value = ''
  await nextTick()
  if (msgContainer.value) msgContainer.value.scrollTop = msgContainer.value.scrollHeight
}

onMounted(loadSessions)
</script>
