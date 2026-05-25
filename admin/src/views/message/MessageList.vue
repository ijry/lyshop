<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('message.list.title') }}</h2>
      <select v-model="filterGroup" @change="load" class="border border-slate-200 rounded-lg px-3 py-1.5 text-sm">
        <option value="">{{ $t('message.list.allGroup') }}</option>
        <option value="system">{{ $t('message.group.system') }}</option>
        <option value="order">{{ $t('message.group.order') }}</option>
        <option value="marketing">{{ $t('message.group.marketing') }}</option>
        <option value="im">{{ $t('message.group.service') }}</option>
      </select>
    </div>
    <div class="bg-white rounded-xl shadow-sm overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 text-slate-500">
          <tr>
            <th class="px-4 py-3 text-left">{{ $t('common.id') }}</th>
            <th class="px-4 py-3 text-left">{{ $t('message.list.group') }}</th>
            <th class="px-4 py-3 text-left">{{ $t('message.list.msgTitle') }}</th>
            <th class="px-4 py-3 text-left">{{ $t('message.list.recipient') }}</th>
            <th class="px-4 py-3 text-left">{{ $t('message.list.time') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="m in messages" :key="m.id" class="border-t border-slate-50 hover:bg-slate-50/50">
            <td class="px-4 py-3 text-slate-400">{{ m.id }}</td>
            <td class="px-4 py-3"><span :class="groupClass(m.group)" class="px-2 py-0.5 rounded text-xs">{{ groupLabel(m.group) }}</span></td>
            <td class="px-4 py-3 text-slate-800">{{ m.title }}</td>
            <td class="px-4 py-3 text-slate-500">{{ m.user_id === 0 ? $t('message.list.allUsers') : `#${m.user_id}` }}</td>
            <td class="px-4 py-3 text-slate-400">{{ m.created_at?.slice(0, 16) }}</td>
          </tr>
          <tr v-if="!messages.length"><td colspan="5" class="px-4 py-8 text-center text-slate-400">{{ $t('message.list.noData') }}</td></tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'

const { t } = useI18n()

const messages = ref<any[]>([])
const filterGroup = ref('')

const groupLabelMap: Record<string, string> = {
  system: t('message.group.system'),
  order: t('message.group.order'),
  marketing: t('message.group.marketing'),
  im: t('message.group.service'),
}
const groupClasses: Record<string, string> = {
  system: 'bg-red-50 text-red-600', order: 'bg-blue-50 text-blue-600',
  marketing: 'bg-orange-50 text-orange-600', im: 'bg-green-50 text-green-600',
}
const groupLabel = (g: string) => groupLabelMap[g] || g
const groupClass = (g: string) => groupClasses[g] || 'bg-slate-50 text-slate-500'

async function load() {
  const data: any = await request.get('/messages', { params: { group: filterGroup.value } })
  messages.value = data?.list || []
}

onMounted(load)
</script>
