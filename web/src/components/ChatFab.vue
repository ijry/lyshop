<template>
  <button
    @click="chat.open('fab')"
    class="fixed bottom-6 right-6 z-50 w-14 h-14 rounded-full bg-red-600 flex-center shadow-lg cursor-pointer hover:bg-red-700 hover:scale-105 transition-all group"
    style="box-shadow: 0 4px 20px rgba(220,38,38,0.35);"
  >
    <div class="i-carbon-chat text-xl text-white" />
    <div class="absolute -top-1 -right-1 w-3 h-3 bg-green-500 rounded-full border-2 border-white" />
    <div class="absolute right-full mr-3 px-3 py-1.5 bg-gray-800 text-white text-xs rounded-lg whitespace-nowrap opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none">
      {{ $t('chatFab.onlineService') }}
    </div>
  </button>
</template>

<script setup lang="ts">
import { useChatStore } from '@/stores/chat'

const chat = useChatStore()
</script>
