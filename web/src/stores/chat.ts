import { defineStore } from 'pinia'

export const useChatStore = defineStore('chat', {
  state: () => ({
    show: false,
    source: 'global',
    messages: [] as Array<{ id: number; sender_type: number; content: string }>,
    inputText: '',
  }),
  actions: {
    open(source = 'global') {
      this.source = source
      this.show = true
      if (!this.messages.length) {
        this.messages.push({ id: Date.now(), sender_type: 2, content: '您好，有什么可以帮您？' })
      }
    },
    close() {
      this.show = false
    },
    send(text: string) {
      const content = text.trim()
      if (!content) return
      this.messages.push({ id: Date.now(), sender_type: 1, content })
      setTimeout(() => {
        this.messages.push({
          id: Date.now() + 1,
          sender_type: 2,
          content: '已收到，我这边继续帮您跟进。'
        })
      }, 400)
    },
  },
})
