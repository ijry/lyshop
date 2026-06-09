<template>
  <div class="space-y-5">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">{{ $t('imAnalytics.title') }}</h1>
        <p class="text-sm text-slate-400 mt-1">{{ $t('imAnalytics.subtitle') }}</p>
      </div>
      <button class="px-4 py-2 rounded-lg bg-red-600 text-white text-sm hover:bg-red-700 transition-colors cursor-pointer disabled:opacity-50" :disabled="loading" @click="load">
        {{ loading ? $t('common.loading') : $t('common.refresh') }}
      </button>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
      <div v-for="item in cards" :key="item.key" class="bg-white rounded-xl border border-slate-100 p-4">
        <div class="text-xs text-slate-500">{{ item.label }}</div>
        <div class="text-2xl font-semibold text-slate-900 mt-2">{{ item.value }}</div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-slate-100 overflow-hidden">
      <div class="px-4 py-3 border-b border-slate-100 text-sm font-medium text-slate-700">{{ $t('imAnalytics.trend') }}</div>
      <table class="w-full text-sm">
        <thead class="bg-slate-50 text-slate-500">
          <tr>
            <th class="p-3 text-left">{{ $t('imAnalytics.date') }}</th>
            <th v-for="item in cards" :key="item.key" class="p-3 text-left">{{ item.label }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in trend" :key="row.date" class="border-t border-slate-100 hover:bg-slate-50">
            <td class="p-3 text-slate-700">{{ row.date }}</td>
            <td v-for="item in cards" :key="item.key" class="p-3 text-slate-600">{{ row[item.key] || 0 }}</td>
          </tr>
          <tr v-if="!trend.length">
            <td :colspan="cards.length + 1" class="p-8 text-center text-slate-400">{{ $t('common.noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import request from '@/api/request'

const summary = ref<Record<string, number>>({})
const trend = ref<any[]>([])
const loading = ref(false)

const labels: Record<string, string> = {
  sessions: '会话',
  messages: '消息',
  ai_replies: 'AI回复',
  ai_failed: 'AI失败',
  rag_hits: 'RAG命中',
  to_human: '转人工',
  accepts: '接入',
  closes: '关闭',
  transfers: '转接',
  files: '文件',
}

const cards = computed(() => Object.keys(labels).map((key) => ({
  key,
  label: labels[key],
  value: summary.value[key] || 0,
})))

async function load() {
  loading.value = true
  try {
    const data: any = await request.get('/im/analytics')
    summary.value = data?.summary || {}
    trend.value = data?.trend || []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
