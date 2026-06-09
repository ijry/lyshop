import { defineStore } from 'pinia'
import { i18n } from '@/locales'
import { get, upload } from '@/api/request'

const t = (key: string) => i18n.global.t(key)

let ws: WebSocket | null = null
let heartbeat: any = null
let reconnectTimer: any = null
let reconnectDelay = 3000

export const useChatStore = defineStore('chat', {
  state: () => ({
    show: false,
    source: 'global',
    messages: [] as Array<{ id: number; sender_type: number; content: string; type?: string; extra?: any }>,
    inputText: '',
    connected: false,
    queuePosition: 0,
    sessionID: 0,
    context: {} as Record<string, any>,
    peerDraft: '',
  }),
  actions: {
    setContext(context: Record<string, any>) {
      this.context = { ...this.context, ...(context || {}) }
    },
    async open(source = 'global') {
      this.source = source
      this.show = true
      if (!this.messages.length) {
        this.messages.push({ id: Date.now(), sender_type: 2, content: t('chatStore.welcome') })
      }

      // Initialize session and WebSocket
      if (!this.sessionID) {
        try {
          const session = await get<any>('/api/v1/im/session', buildSessionContextParams(this.context))
          if (session) {
            this.sessionID = session.id
            this.queuePosition = session.queue_position || 0
            const data = await get<any>('/api/v1/im/messages', { session_id: session.id, size: 50 })
            const history = (data?.list || []).reverse()
            if (history.length) {
              this.messages = history
            }
          }

          const token = localStorage.getItem('user_token')
          if (token && !ws) {
            this.connectWS(token)
          }
        } catch (err) {
          console.error('Failed to initialize chat:', err)
        }
      }
    },
    close() {
      this.show = false
    },
    send(text: string) {
      const content = text.trim()
      if (!content) return

      this.messages.push({ id: Date.now(), sender_type: 1, content, type: 'text' })

      if (ws?.readyState === WebSocket.OPEN && this.sessionID) {
        this.sendTypingStop()
        ws.send(JSON.stringify({
          type: 'msg',
          session_id: this.sessionID,
          payload: { msg_type: 'text', content }
        }))
        return
      }

      // Fallback: mock auto-reply
      setTimeout(() => {
        this.messages.push({
          id: Date.now() + 1,
          sender_type: 2,
          content: t('chatStore.followUp')
        })
      }, 400)
    },
    sendTypingDraft(draft: string) {
      if (ws?.readyState !== WebSocket.OPEN || !this.sessionID) return
      ws.send(JSON.stringify({
        type: 'typing_draft',
        session_id: this.sessionID,
        payload: { draft },
      }))
    },
    sendTypingStop() {
      if (ws?.readyState !== WebSocket.OPEN || !this.sessionID) return
      ws.send(JSON.stringify({
        type: 'typing_stop',
        session_id: this.sessionID,
        payload: {},
      }))
    },
    connectWS(token: string) {
      const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
      const host = location.hostname + ':8080'
      ws = new WebSocket(`${protocol}//${host}/ws/im?token=${token}`)

      ws.onopen = () => {
        this.connected = true
        reconnectDelay = 3000
      }

      ws.onmessage = (e) => {
        try {
          const frame = JSON.parse(e.data)
          if (frame.type === 'msg') {
            this.messages.push({
              id: Date.now(),
              sender_type: frame.payload.sender_type ?? 2,
              content: frame.payload.content,
              type: frame.payload.msg_type || 'text',
              extra: frame.payload.extra,
            })
          } else if (frame.type === 'queue') {
            this.queuePosition = frame.payload?.position || 0
          } else if (frame.type === 'assign') {
            if (frame.payload?.action === 'accepted') {
              this.queuePosition = 0
              this.messages.push({
                id: Date.now(),
                sender_type: 0,
                content: t('chatStore.assignedNotice'),
              })
            }
          } else if (frame.type === 'close') {
            this.messages.push({
              id: Date.now(),
              sender_type: 0,
              content: t('chatStore.closedNotice'),
            })
          } else if (frame.type === 'typing_draft') {
            this.peerDraft = frame.payload?.draft || ''
          } else if (frame.type === 'typing_stop') {
            this.peerDraft = ''
          }
        } catch {}
      }

      ws.onclose = () => {
        this.connected = false
        this.scheduleReconnect(token)
      }

      ws.onerror = () => {
        this.connected = false
      }

      if (heartbeat) clearInterval(heartbeat)
      heartbeat = setInterval(() => {
        if (ws?.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify({ type: 'ping' }))
        }
      }, 30000)
    },
    async sendAttachment(file: File) {
      if (!this.sessionID) return
      const info: any = await upload('/api/v1/im/upload', file, { session_id: String(this.sessionID) })
      const extra = {
        file_url: info.url,
        file_path: info.path,
        file_name: info.name,
        file_size: info.size,
        mime: info.mime,
      }
      this.messages.push({
        id: Date.now(),
        sender_type: 1,
        content: info.name,
        type: info.message_type,
        extra,
      })
      if (ws?.readyState === WebSocket.OPEN) {
        this.sendTypingStop()
        ws.send(JSON.stringify({
          type: 'msg',
          session_id: this.sessionID,
          payload: { msg_type: info.message_type, content: info.name, extra },
        }))
      }
    },
    scheduleReconnect(token: string) {
      if (reconnectTimer) clearTimeout(reconnectTimer)
      reconnectTimer = setTimeout(() => {
        reconnectDelay = Math.min(reconnectDelay * 2, 30000)
        this.connectWS(token)
      }, reconnectDelay)
    },
    cleanup() {
      if (heartbeat) clearInterval(heartbeat)
      if (reconnectTimer) clearTimeout(reconnectTimer)
      ws?.close()
      ws = null
    }
  },
})

function buildSessionContextParams(context: Record<string, any>) {
  const params: Record<string, string> = {}
  const keys = [
    'visitor_id',
    'visitor_location',
    'visitor_browser',
    'visitor_os',
    'visitor_language',
    'visitor_referrer',
    'visitor_url',
    'visitor_device',
  ]
  for (const key of keys) {
    if (context?.[key]) params[key] = String(context[key])
  }
  if (context?.visitor_extra) {
    params.visitor_extra = JSON.stringify(context.visitor_extra)
  }
  return params
}
