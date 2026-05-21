<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">订单列表</h2>
    </div>

    <!-- Status tabs -->
    <div class="flex gap-2 mb-4">
      <button v-for="tab in tabs" :key="tab.value"
        @click="activeStatus = tab.value; loadOrders()"
        :class="activeStatus === tab.value
          ? 'bg-blue-700 text-white'
          : 'bg-white text-slate-600 border border-slate-200 hover:bg-slate-50'"
        class="px-4 py-2 rounded-xl text-sm transition">
        {{ tab.label }}
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">订单号</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">用户ID</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">金额</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">状态</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">下单时间</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="o in orders" :key="o.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 font-mono text-xs text-slate-500">{{ o.order_no }}</td>
            <td class="px-4 py-3 text-slate-600">{{ o.user_id }}</td>
            <td class="px-4 py-3 font-medium text-slate-800">¥{{ o.total_amount }}</td>
            <td class="px-4 py-3">
              <span :class="statusClass(o.status)" class="px-2 py-1 rounded-full text-xs">
                {{ statusLabel(o.status) }}
              </span>
            </td>
            <td class="px-4 py-3 text-slate-400 text-xs">{{ o.created_at?.slice(0,10) }}</td>
            <td class="px-4 py-3">
              <button v-if="o.status === 2" @click="ship(o.id)"
                class="text-blue-600 hover:underline text-xs">发货</button>
            </td>
          </tr>
          <tr v-if="!orders.length">
            <td colspan="6" class="px-4 py-12 text-center text-slate-400">暂无订单</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getOrders, shipOrder } from '@/api/plugins'

const orders = ref<any[]>([])
const activeStatus = ref(0)

const tabs = [
  { label: '全部', value: 0 },
  { label: '待付款', value: 1 },
  { label: '待发货', value: 2 },
  { label: '待收货', value: 3 },
  { label: '已完成', value: 4 },
]

const statusLabels: Record<number, string> = {
  1: '待付款', 2: '待发货', 3: '待收货', 4: '已完成', 5: '售后'
}
const statusColors: Record<number, string> = {
  1: 'bg-yellow-50 text-yellow-600',
  2: 'bg-blue-50 text-blue-600',
  3: 'bg-purple-50 text-purple-600',
  4: 'bg-green-50 text-green-600',
  5: 'bg-red-50 text-red-500',
}
const statusLabel = (s: number) => statusLabels[s] || '未知'
const statusClass = (s: number) => statusColors[s] || 'bg-slate-50 text-slate-400'

async function loadOrders() {
  const data: any = await getOrders({ status: activeStatus.value || undefined, page: 1, size: 20 })
  orders.value = data.list || []
}

async function ship(id: number) {
  const no = prompt('请输入快递单号')
  if (!no) return
  await shipOrder(id, no)
  loadOrders()
}

onMounted(loadOrders)
</script>
