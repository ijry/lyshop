<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('afterSale.list.title') }}</h2>
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
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('afterSale.list.caseNo') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('afterSale.list.type') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('afterSale.list.order') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.status') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('afterSale.list.applyTime') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="item in list" :key="item.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-xs font-mono text-slate-600">{{ item.case_no }}</td>
            <td class="px-4 py-3 text-slate-700">{{ typeLabel(item.case_type) }}</td>
            <td class="px-4 py-3 text-slate-600">#{{ item.order_id }}</td>
            <td class="px-4 py-3"><span class="px-2 py-1 rounded-full text-xs bg-slate-100 text-slate-600">{{ statusText(item.status, item.status_label) }}</span></td>
            <td class="px-4 py-3 text-xs text-slate-400">{{ formatDate(item.created_at) }}</td>
            <td class="px-4 py-3"><button class="text-blue-600 hover:underline text-xs" @click="goDetail(item.id)">{{ $t('common.view') }}</button></td>
          </tr>
          <tr v-if="!list.length">
            <td colspan="6" class="px-4 py-12 text-center text-slate-400">{{ $t('afterSale.list.noData') }}</td>
          </tr>
        </tbody>
      </table>
      <div class="px-4 py-3 flex items-center justify-between border-t border-slate-100 text-sm text-slate-500">
        <span>{{ $t('common.totalCount', { total }) }}</span>
        <div class="flex gap-2">
          <button :disabled="query.page <= 1" @click="prevPage"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40">{{ $t('common.prevPage') }}</button>
          <button :disabled="query.page * query.size >= total" @click="nextPage"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40">{{ $t('common.nextPage') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getAfterSales } from '@/api/plugins'
import { afterSaleStatusLabel } from '@/utils/order-status'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const list = ref<any[]>([])
const activeStatus = ref('')
const total = ref(0)
const query = ref({ page: 1, size: 20 })

const tabs = computed(() => [
  { label: t('afterSale.list.all'), value: '' },
  { label: t('afterSale.list.applied'), value: 'applied' },
  { label: t('afterSale.list.waitReturn'), value: 'approved_wait_user_return' },
  { label: t('afterSale.list.returning'), value: 'user_returning' },
  { label: t('afterSale.list.waitRefund'), value: 'refund_pending' },
  { label: t('afterSale.list.waitReship'), value: 'reship_pending' },
  { label: t('afterSale.list.completed'), value: 'completed' },
])

const typeLabel = (v: string) => v === 'exchange' ? t('afterSale.list.exchange') : t('afterSale.list.return')
const statusLabel = (v: string) => afterSaleStatusLabel(v)
const formatDate = (v?: string) => v ? String(v).slice(0, 19).replace('T', ' ') : '-'

function statusText(status: string, label?: string) {
  const value = String(label || '').trim()
  if (value) return value
  return statusLabel(status)
}

function goDetail(id: number) {
  router.push(`/order/after-sale/detail/${id}`)
}

async function load() {
  const params: any = { page: query.value.page, size: query.value.size }
  if (route.query.order_id) params.order_id = Number(route.query.order_id)
  if (activeStatus.value) params.status = activeStatus.value
  const data: any = await getAfterSales(params)
  list.value = data?.list || []
  total.value = Number(data?.total || 0)
}

function onTabChange(status: string) {
  activeStatus.value = status
  query.value.page = 1
  load()
}

function prevPage() {
  if (query.value.page <= 1) return
  query.value.page -= 1
  load()
}

function nextPage() {
  if (query.value.page * query.value.size >= total.value) return
  query.value.page += 1
  load()
}

onMounted(load)
</script>
