<template>
  <teleport to="body">
    <div v-if="chat.show" class="fixed inset-0 z-[2000] flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/35" @click="chat.close()" />
      <div class="relative w-[92vw] max-w-[760px] h-[78vh] max-h-[720px] bg-white rounded-[20px] overflow-hidden flex flex-col">
        <!-- Header -->
        <div class="px-5 py-3 border-b border-gray-100 flex-between">
          <div class="flex items-center gap-2">
            <div :class="chat.connected ? 'bg-green-500' : 'bg-yellow-500'" class="w-2 h-2 rounded-full" />
            <span class="text-sm font-medium text-gray-700">
              {{ chat.connected ? $t('chatDialog.serviceOnline') : $t('chatDialog.connecting') }}
            </span>
          </div>
          <button class="text-xs text-gray-400 hover:text-gray-600" @click="chat.close()">{{ $t('chatDialog.close') }}</button>
        </div>

        <!-- Queue Notice -->
        <div v-if="chat.queuePosition > 0" class="px-5 py-2.5 bg-yellow-50 border-b border-yellow-100 flex items-center gap-2">
          <div class="i-carbon-time text-yellow-600" />
          <span class="text-sm text-yellow-700">{{ $t('chatDialog.queueNotice', { position: chat.queuePosition }) }}</span>
        </div>

        <!-- Messages -->
        <div ref="msgBox" class="flex-1 overflow-y-auto p-5 space-y-4 bg-gray-50">
          <div v-for="m in chat.messages" :key="m.id">
            <!-- System message -->
            <div v-if="m.sender_type === 0" class="flex justify-center">
              <span class="text-xs text-gray-400 bg-white px-3 py-1 rounded-full shadow-sm">{{ m.content }}</span>
            </div>
            <!-- User/Agent message -->
            <div v-else :class="m.sender_type === 1 ? 'flex-row-reverse' : ''" class="flex items-end gap-3">
              <div :class="m.sender_type === 2 ? 'bg-red-600' : 'bg-gray-200'" class="w-8 h-8 rounded-full flex-center shrink-0">
                <span :class="m.sender_type === 2 ? 'text-white' : 'text-gray-600'" class="text-xs font-medium">
                  {{ m.sender_type === 2 ? $t('chatDialog.agent') : $t('chatDialog.me') }}
                </span>
              </div>
              <div :class="m.sender_type === 1
                ? 'bg-red-600 text-white rounded-tl-2xl rounded-tr-sm rounded-bl-2xl rounded-br-2xl'
                : 'bg-white text-gray-800 rounded-tl-sm rounded-tr-2xl rounded-bl-2xl rounded-br-2xl'"
                class="max-w-sm px-4 py-2.5 text-sm leading-relaxed shadow-sm">
                {{ m.content }}
              </div>
            </div>
          </div>
        </div>

        <!-- Input -->
        <div class="px-5 py-3 border-t border-gray-100 flex gap-3 bg-white">
          <input v-model="chat.inputText" @keyup.enter="send" :placeholder="$t('chatDialog.inputPlaceholder')" class="input-base flex-1 !rounded-xl" />
          <button @click="send" :disabled="!chat.inputText.trim()" class="btn-primary !px-6">{{ $t('chatDialog.send') }}</button>
        </div>
      </div>
    </div>
  </teleport>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { useChatStore } from '@/stores/chat'

const chat = useChatStore()
const msgBox = ref<HTMLElement>()

function send() {
  const text = chat.inputText.trim()
  if (!text) return
  chat.send(text)
  chat.inputText = ''
  scrollBottom()
}

function scrollBottom() {
  nextTick(() => {
    if (msgBox.value) msgBox.value.scrollTop = msgBox.value.scrollHeight
  })
}

watch(() => chat.messages.length, scrollBottom)
</script>
