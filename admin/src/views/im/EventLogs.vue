<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-semibold text-slate-900">{{ $t('imLogs.title') }}</h1>
      <button class="px-4 py-2 rounded-lg bg-red-600 text-white text-sm hover:bg-red-700 transition-colors cursor-pointer disabled:opacity-50" :disabled="loading" @click="load">
        {{ loading ? $t('common.loading') : $t('common.search') }}
      </button>
    </div>

    <div class="bg-white rounded-xl border border-slate-100 p-4 flex flex-wrap gap-3">
      <label class="sr-only" for="im-log-event">event</label>
      <input id="im-log-event" v-model="filters.event" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-red-400" placeholder="event" />
      <label class="sr-only" for="im-log-session">session_id</label>
      <input id="im-log-session" v-model="filters.session_id" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-red-400" placeholder="session_id" />
      <label class="sr-only" for="im-log-source">source</label>
      <input id="im-log-source" v-model="filters.source" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-red-400" placeholder="source" />
      <input v-model="filters.level" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-red-400" placeholder="level" />
      <input v-model="filters.category" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-red-400" placeholder="category" />
      <input v-model="filters.trace_id" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-red-400" placeholder="trace_id" />
      <input v-model="filters.keyword" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-red-400" placeholder="keyword" />
    </div>

    <div class="bg-white rounded-xl border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 text-slate-500">
          <tr>
            <th class="p-3 text-left">ID</th>
            <th class="p-3 text-left">Event</th>
            <th class="p-3 text-left">Level</th>
            <th class="p-3 text-left">Session</th>
            <th class="p-3 text-left">Source</th>
            <th class="p-3 text-left">Message</th>
            <th class="p-3 text-left">Success</th>
            <th class="p-3 text-left">Time</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in list" :key="row.id" class="border-t border-slate-100 hover:bg-slate-50">
            <td class="p-3 text-slate-500">{{ row.id }}</td>
            <td class="p-3 text-slate-800">{{ row.event }}</td>
            <td class="p-3 text-slate-600">{{ row.level || 'info' }}</td>
            <td class="p-3 text-slate-600">{{ row.session_id }}</td>
            <td class="p-3 text-slate-600">{{ row.source }}</td>
            <td class="p-3 text-slate-500 max-w-xs truncate" :title="row.message || row.meta || row.extra">{{ row.message || row.meta || row.extra }}</td>
            <td class="p-3">
              <span :class="row.success === 1 ? 'bg-green-50 text-green-600' : 'bg-red-50 text-red-600'" class="px-2 py-0.5 rounded-full text-xs">
                {{ row.success === 1 ? 'ok' : 'fail' }}
              </span>
            </td>
            <td class="p-3 text-slate-500">{{ row.created_at }}</td>
          </tr>
          <tr v-if="!list.length">
            <td colspan="8" class="p-8 text-center text-slate-400">{{ $t('common.noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import request from '@/api/request'

const filters = ref({ event: '', session_id: '', source: '', level: '', category: '', trace_id: '', keyword: '' })
const list = ref<any[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const params: any = {}
    if (filters.value.event) params.event = filters.value.event
    if (filters.value.session_id) params.session_id = filters.value.session_id
    if (filters.value.source) params.source = filters.value.source
    if (filters.value.level) params.level = filters.value.level
    if (filters.value.category) params.category = filters.value.category
    if (filters.value.trace_id) params.trace_id = filters.value.trace_id
    if (filters.value.keyword) params.keyword = filters.value.keyword
    const data: any = await request.get('/im/logs', { params })
    list.value = data?.list || []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
