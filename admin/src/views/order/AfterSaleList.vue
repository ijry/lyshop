<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">售后列表</h2>
    </div>

    <div class="flex gap-2 mb-4 flex-wrap">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        @click="onTabChange(tab.value)"
        :class="activeStatus === tab.value ? 'bg-blue-700 text-white' : 'bg-white text-slate-600 border border-slate-200 hover:bg-slate-50'"
        class="px-4 py-2 rounded-xl text-sm transition"
      >
        {{ tab.label }}
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">售后单号</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">类型</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">订单</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">状态</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">申请时间</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="item in list" :key="item.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-xs font-mono text-slate-600">{{ item.case_no }}</td>
            <td class="px-4 py-3 text-slate-700">{{ typeLabel(item.case_type) }}</td>
            <td class="px-4 py-3 text-slate-600">#{{ item.order_id }}</td>
            <td class="px-4 py-3"><span class="px-2 py-1 rounded-full text-xs bg-slate-100 text-slate-600">{{ statusLabel(item.status) }}</span></td>
            <td class="px-4 py-3 text-xs text-slate-400">{{ formatDate(item.created_at) }}</td>
            <td class="px-4 py-3"><button class="text-blue-600 hover:underline text-xs" @click="goDetail(item.id)">查看</button></td>
          </tr>
          <tr v-if="!list.length">
            <td colspan="6" class="px-4 py-12 text-center text-slate-400">暂无售后单</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getAfterSales } from '@/api/plugins'
import { afterSaleStatusLabel } from '@/utils/order-status'

const route = useRoute()
const router = useRouter()
const list = ref<any[]>([])
const activeStatus = ref('')

const tabs = [
  { label: '全部', value: '' },
  { label: '申请中', value: 'applied' },
  { label: '待用户回寄', value: 'approved_wait_user_return' },
  { label: '用户回寄中', value: 'user_returning' },
  { label: '待退款', value: 'refund_pending' },
  { label: '待补发', value: 'reship_pending' },
  { label: '已完结', value: 'completed' },
]

const typeLabel = (v: string) => v === 'exchange' ? '换货' : '退货'
const statusLabel = (v: string) => afterSaleStatusLabel(v)
const formatDate = (v?: string) => v ? String(v).slice(0, 19).replace('T', ' ') : '-'

function goDetail(id: number) {
  router.push(`/order/after-sale/detail/${id}`)
}

async function load() {
  const params: any = { page: 1, size: 20 }
  if (route.query.order_id) params.order_id = Number(route.query.order_id)
  if (activeStatus.value) params.status = activeStatus.value
  const data: any = await getAfterSales(params)
  list.value = data?.list || []
}

function onTabChange(status: string) {
  activeStatus.value = status
  load()
}

onMounted(load)
</script>
