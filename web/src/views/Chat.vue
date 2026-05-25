<template>
  <div class="max-w-3xl mx-auto px-6 py-8">
    <h1 class="text-xl font-bold text-gray-900 mb-6">{{ $t('chat.title') }}</h1>
    <div class="card flex flex-col" style="height: 600px;">
      <!-- Header -->
      <div class="px-5 py-3 border-b border-gray-100 flex-between">
        <div class="flex items-center gap-2">
          <div class="w-2 h-2 bg-green-500 rounded-full" />
          <span class="text-sm font-medium text-gray-700">{{ $t('chat.serviceOnline') }}</span>
        </div>
        <span class="text-xs text-gray-400">{{ $t('chat.sessionName') }}</span>
      </div>

      <!-- Messages -->
      <div ref="msgBox" class="flex-1 overflow-y-auto p-5 space-y-4">
        <!-- Welcome -->
        <div v-if="!messages.length" class="flex flex-col items-center py-16">
          <div class="w-16 h-16 rounded-full bg-red-50 flex-center mb-3">
            <div class="i-carbon-chat text-3xl text-red-500" />
          </div>
          <p class="text-gray-400 text-sm">{{ $t('chat.welcomeMessage') }}</p>
        </div>

        <div v-for="m in messages" :key="m.id"
          :class="m.sender_type === 1 ? 'flex-row-reverse' : ''"
          class="flex items-end gap-3">
          <div :class="m.sender_type === 2 ? 'bg-red-600' : 'bg-gray-200'"
            class="w-8 h-8 rounded-full flex-center shrink-0">
            <span :class="m.sender_type === 2 ? 'text-white' : 'text-gray-600'" class="text-xs font-medium">
              {{ m.sender_type === 2 ? $t('chat.agent') : $t('chat.me') }}
            </span>
          </div>
          <div :class="m.sender_type === 1
            ? 'bg-red-600 text-white rounded-tl-2xl rounded-tr-sm rounded-bl-2xl rounded-br-2xl'
            : 'bg-gray-100 text-gray-800 rounded-tl-sm rounded-tr-2xl rounded-bl-2xl rounded-br-2xl'"
            class="max-w-sm px-4 py-2.5 text-sm leading-relaxed">
            {{ m.content }}
          </div>
        </div>
      </div>

      <!-- Input -->
      <div class="px-5 py-3 border-t border-gray-100 flex gap-3">
        <input v-model="inputText" @keyup.enter="send"
          :placeholder="$t('chat.inputPlaceholder')"
          class="input-base flex-1 !rounded-xl" />
        <button @click="send" :disabled="!inputText.trim()" class="btn-primary !px-6">{{ $t('chat.send') }}</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'

const messages = ref<any[]>([])
const inputText = ref('')
const msgBox = ref<HTMLElement>()

function send() {
  const text = inputText.value.trim()
  if (!text) return
  messages.value.push({ id: Date.now(), sender_type: 1, content: text })
  inputText.value = ''
  scrollBottom()

  // Mock auto-reply
  setTimeout(() => {
    const replies = [
      '您好，感谢您的咨询！请问有什么可以帮您的？',
      '这款商品目前有货，支持全国包邮哦~',
      '好的，我帮您查询一下，请稍等。',
      '您可以在"我的订单"页面查看物流信息。',
    ]
    messages.value.push({
      id: Date.now() + 1,
      sender_type: 2,
      content: replies[Math.floor(Math.random() * replies.length)],
    })
    scrollBottom()
  }, 800 + Math.random() * 1200)
}

function scrollBottom() {
  nextTick(() => {
    if (msgBox.value) msgBox.value.scrollTop = msgBox.value.scrollHeight
  })
}
</script>
