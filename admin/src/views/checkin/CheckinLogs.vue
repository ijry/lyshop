<template>
  <div>
    <h2 class="text-xl font-semibold text-slate-800 mb-6">签到记录</h2>
    <div class="bg-white rounded-xl shadow-sm overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 text-slate-500">
          <tr>
            <th class="px-4 py-3 text-left">ID</th>
            <th class="px-4 py-3 text-left">用户ID</th>
            <th class="px-4 py-3 text-left">日期</th>
            <th class="px-4 py-3 text-left">连续天数</th>
            <th class="px-4 py-3 text-left">获得积分</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="log in logs" :key="log.id" class="border-t border-slate-50 hover:bg-slate-50/50">
            <td class="px-4 py-3 text-slate-400">{{ log.id }}</td>
            <td class="px-4 py-3 text-slate-700">{{ log.user_id }}</td>
            <td class="px-4 py-3 text-slate-600">{{ log.date }}</td>
            <td class="px-4 py-3"><span class="bg-blue-50 text-blue-600 px-2 py-0.5 rounded text-xs">{{ log.consecutive_days }}天</span></td>
            <td class="px-4 py-3 text-green-600 font-medium">+{{ log.points }}</td>
          </tr>
          <tr v-if="!logs.length"><td colspan="5" class="px-4 py-8 text-center text-slate-400">暂无记录</td></tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'

const logs = ref<any[]>([])

onMounted(async () => {
  const data: any = await request.get('/checkin/logs')
  logs.value = data?.list || []
})
</script>
