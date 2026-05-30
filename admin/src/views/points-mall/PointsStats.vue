<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">积分统计</h2>
    </div>

    <div class="grid grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-5">
        <div class="text-slate-500 text-sm mb-2">累计发放</div>
        <div class="text-2xl font-bold text-blue-600">{{ stats.total_issued || 0 }}</div>
      </div>
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-5">
        <div class="text-slate-500 text-sm mb-2">累计消耗</div>
        <div class="text-2xl font-bold text-orange-600">{{ stats.total_consumed || 0 }}</div>
      </div>
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-5">
        <div class="text-slate-500 text-sm mb-2">当前余额</div>
        <div class="text-2xl font-bold text-green-600">{{ stats.total_balance || 0 }}</div>
      </div>
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-5">
        <div class="text-slate-500 text-sm mb-2">商品数量</div>
        <div class="text-2xl font-bold text-purple-600">{{ stats.product_count || 0 }}</div>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4 mb-6">
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-5">
        <div class="text-slate-700 font-medium mb-3">今日积分</div>
        <div class="flex justify-between items-center">
          <div>
            <div class="text-slate-500 text-sm">今日发放</div>
            <div class="text-xl font-bold text-blue-600">{{ stats.today_issued || 0 }}</div>
          </div>
          <div>
            <div class="text-slate-500 text-sm">今日消耗</div>
            <div class="text-xl font-bold text-orange-600">{{ stats.today_consumed || 0 }}</div>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-5">
        <div class="text-slate-700 font-medium mb-3">兑换统计</div>
        <div class="flex justify-between items-center">
          <div>
            <div class="text-slate-500 text-sm">兑换次数</div>
            <div class="text-xl font-bold text-green-600">{{ stats.exchange_count || 0 }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import request from '@/api/request'

const stats = ref<any>({})

async function load() {
  const data: any = await request.get('/points/stats')
  stats.value = data
}

onMounted(load)
</script>
