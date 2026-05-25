<template>
  <div class="min-h-screen flex flex-col" :style="{ background: 'var(--color-bg-page, #f9fafb)' }">
    <Header />
    <main class="flex-1">
      <router-view />
    </main>
    <Footer />
    <ChatFab />
    <ChatDialog />
  </div>
</template>

<script setup lang="ts">
import { useSiteStore } from '@/stores/site'
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import ChatFab from '@/components/ChatFab.vue'
import ChatDialog from '@/components/ChatDialog.vue'

const siteStore = useSiteStore()
siteStore.load()
</script>
