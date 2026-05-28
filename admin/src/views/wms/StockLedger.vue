<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('wms.stockLedgerTitle') }}</h2>
    </div>

    <div class="bg-white rounded-xl shadow-sm p-4 mb-4 flex gap-3 border border-slate-100 flex-wrap">
      <select v-model.number="query.warehouse_id" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none">
        <option :value="0">{{ $t('stock.allWarehouse') }}</option>
        <option v-for="item in warehouses" :key="item.id" :value="item.id">{{ item.name }}</option>
      </select>
      <input
        v-model.trim="query.keyword"
        :placeholder="$t('wms.stockKeywordPlaceholder')"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm flex-1 min-w-[220px] focus:outline-none focus:border-blue-400"
      />
      <button class="px-4 py-2 bg-slate-100 rounded-lg text-sm hover:bg-slate-200" @click="onSearch">{{ $t('common.search') }}</button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.warehouse') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.skuId') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.skuName') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('stock.stockQty') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('stock.safetyStock') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.status') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="row in rows" :key="row.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-slate-700">{{ row.warehouse_name || row.warehouse_id }}</td>
            <td class="px-4 py-3 text-slate-700">{{ row.sku_id }}</td>
            <td class="px-4 py-3 text-slate-700">{{ row.sku_name || '-' }}</td>
            <td class="px-4 py-3 font-medium text-slate-800">{{ row.qty }}</td>
            <td class="px-4 py-3">
              <input
                v-model.number="safeQtyMap[row.id]"
                type="number"
                min="0"
                class="w-24 border border-slate-200 rounded-lg px-2 py-1.5 text-sm"
              />
            </td>
            <td class="px-4 py-3">
              <span :class="lowStock(row) ? 'bg-red-50 text-red-500' : 'bg-green-50 text-green-600'" class="px-2 py-1 rounded-full text-xs">
                {{ lowStock(row) ? $t('stock.warning') : $t('stock.normal') }}
              </span>
            </td>
            <td class="px-4 py-3">
              <button class="text-blue-600 hover:underline text-xs" @click="saveSafeQty(row.id)">{{ $t('wms.updateSafetyStock') }}</button>
            </td>
          </tr>
          <tr v-if="!rows.length">
            <td colspan="7" class="px-4 py-12 text-center text-slate-400">{{ $t('stock.noData') }}</td>
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
import { listStockLedger, listWarehouses, updateSafetyStock, type WmsStockLedgerRow, type WmsWarehouse } from '@/api/wms'
import { notify } from '@/utils/notify'

const { t } = useI18n()
const warehouses = ref<WmsWarehouse[]>([])
const rows = ref<WmsStockLedgerRow[]>([])
const total = ref(0)
const safeQtyMap = ref<Record<number, number>>({})
const query = ref({
  warehouse_id: 0,
  keyword: '',
  page: 1,
  size: 20,
})

function lowStock(row: WmsStockLedgerRow) {
  return Number(row.qty || 0) <= Number(row.safe_qty || 0)
}

async function loadWarehouses() {
  const data = await listWarehouses()
  warehouses.value = Array.isArray(data) ? data : []
}

async function loadRows() {
  const data = await listStockLedger({
    warehouse_id: query.value.warehouse_id > 0 ? query.value.warehouse_id : undefined,
    keyword: query.value.keyword || undefined,
    page: query.value.page,
    size: query.value.size,
  })
  rows.value = Array.isArray(data?.list) ? data.list : []
  total.value = Number(data?.total || 0)
  const map: Record<number, number> = {}
  for (const row of rows.value) {
    map[row.id] = Number(row.safe_qty || 0)
  }
  safeQtyMap.value = map
}

function onSearch() {
  query.value.page = 1
  loadRows()
}

async function saveSafeQty(id: number) {
  const safeQty = Number(safeQtyMap.value[id] || 0)
  if (safeQty < 0) {
    notify(t('wms.invalidSafetyStock'))
    return
  }
  await updateSafetyStock(id, safeQty)
  await loadRows()
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
