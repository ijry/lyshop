<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('wms.movementTitle') }}</h2>
    </div>

    <div class="bg-white rounded-xl shadow-sm p-4 mb-4 flex gap-3 border border-slate-100 flex-wrap">
      <select v-model.number="query.warehouse_id" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none">
        <option :value="0">{{ $t('stock.allWarehouse') }}</option>
        <option v-for="item in warehouses" :key="item.id" :value="item.id">{{ item.name }}</option>
      </select>
      <input
        v-model.number="query.sku_id"
        type="number"
        :placeholder="$t('wms.skuIdPlaceholder')"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm w-40 focus:outline-none focus:border-blue-400"
      />
      <input
        v-model.trim="query.doc_no"
        :placeholder="$t('wms.docNoPlaceholder')"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm w-56 focus:outline-none focus:border-blue-400"
      />
      <button class="px-4 py-2 bg-slate-100 rounded-lg text-sm hover:bg-slate-200" @click="onSearch">{{ $t('common.search') }}</button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.updatedAt') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.docNo') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.docType') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.warehouse') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.skuInfo') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.changeQty') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.balanceQty') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="row in rows" :key="row.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-slate-600">{{ formatDate(row.created_at) }}</td>
            <td class="px-4 py-3 font-mono text-xs text-slate-600">{{ row.doc_no }}</td>
            <td class="px-4 py-3 text-slate-700">{{ typeLabel(row.type) }}</td>
            <td class="px-4 py-3 text-slate-700">{{ row.warehouse_name || row.warehouse_id }}</td>
            <td class="px-4 py-3 text-slate-700">{{ row.sku_name || '-' }} ({{ row.sku_id }})</td>
            <td class="px-4 py-3 font-medium" :class="row.qty >= 0 ? 'text-green-600' : 'text-red-500'">
              {{ row.qty >= 0 ? '+' : '' }}{{ row.qty }}
            </td>
            <td class="px-4 py-3 text-slate-600">{{ row.before_qty }} → {{ row.after_qty }}</td>
          </tr>
          <tr v-if="!rows.length">
            <td colspan="7" class="px-4 py-12 text-center text-slate-400">{{ $t('common.noData') }}</td>
          </tr>
        </tbody>
      </table>

      <div class="px-4 py-3 flex items-center justify-between border-t border-slate-100 text-sm text-slate-500">
        <span>{{ $t('common.totalCount', { total }) }}</span>
        <div class="flex gap-2">
          <button
            :disabled="query.page <= 1"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40"
            @click="prevPage"
          >
            {{ $t('common.prevPage') }}
          </button>
          <button
            :disabled="query.page * query.size >= total"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40"
            @click="nextPage"
          >
            {{ $t('common.nextPage') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { listMovements, listWarehouses, type WmsMovementRow, type WmsWarehouse } from '@/api/wms'

const { t } = useI18n()
const warehouses = ref<WmsWarehouse[]>([])
const rows = ref<WmsMovementRow[]>([])
const total = ref(0)
const query = ref({
  warehouse_id: 0,
  sku_id: 0,
  doc_no: '',
  page: 1,
  size: 20,
})

function formatDate(value?: string) {
  return value ? String(value).slice(0, 19).replace('T', ' ') : '-'
}

function typeLabel(type: string) {
  return type === 'inbound' ? t('wms.docTypeInbound') : t('wms.docTypeOutbound')
}

async function loadWarehouses() {
  const data = await listWarehouses()
  warehouses.value = Array.isArray(data?.list) ? data.list : []
}

async function loadRows() {
  const data = await listMovements({
    warehouse_id: query.value.warehouse_id > 0 ? query.value.warehouse_id : undefined,
    sku_id: query.value.sku_id > 0 ? query.value.sku_id : undefined,
    doc_no: query.value.doc_no || undefined,
    page: query.value.page,
    size: query.value.size,
  })
  rows.value = Array.isArray(data?.list) ? data.list : []
  total.value = Number(data?.total || 0)
}

function onSearch() {
  query.value.page = 1
  loadRows()
}

function prevPage() {
  if (query.value.page <= 1) return
  query.value.page -= 1
  loadRows()
}

function nextPage() {
  if (query.value.page * query.value.size >= total.value) return
  query.value.page += 1
  loadRows()
}

onMounted(async () => {
  await loadWarehouses()
  await loadRows()
})
</script>
