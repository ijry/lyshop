<template>
  <div>
    <h2 class="text-xl font-semibold text-slate-800 mb-6">控制台</h2>
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div
        v-for="card in cards"
        :key="card.label"
        class="bg-white rounded-xl shadow-sm p-6 border border-slate-100"
      >
        <p class="text-sm text-slate-500 mb-1">{{ card.label }}</p>
        <p class="text-2xl font-bold text-slate-800">{{ card.value }}</p>
      </div>
    </div>
    <div class="bg-white rounded-xl shadow-sm p-6 border border-slate-100">
      <p class="text-slate-400 text-center py-12">图表数据加载中...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
const cards = [
  { label: '今日订单', value: '—' },
  { label: '今日销售额', value: '—' },
  { label: '待处理售后', value: '—' },
  { label: '在线客服会话', value: '—' },
]
</script>
