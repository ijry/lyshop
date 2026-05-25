<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('stock.title') }}</h2>
    </div>

    <!-- Warehouse selector -->
    <div class="bg-white rounded-xl shadow-sm p-4 mb-4 flex gap-3 border border-slate-100">
      <select v-model="warehouseID" @change="loadStocks"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none">
        <option value="">{{ $t('stock.allWarehouse') }}</option>
        <option v-for="w in warehouses" :key="w.id" :value="w.id">{{ w.name }}</option>
      </select>
      <button @click="showInbound = true"
        class="px-4 py-2 bg-green-600 text-white text-sm rounded-lg hover:bg-green-500 transition">{{ $t('stock.inbound') }}</button>
      <button @click="showOutbound = true"
        class="px-4 py-2 bg-orange-500 text-white text-sm rounded-lg hover:bg-orange-400 transition">{{ $t('stock.outbound') }}</button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('stock.warehouseId') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('stock.skuId') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('stock.stockQty') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('stock.safetyStock') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.status') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="s in stocks" :key="s.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-slate-500">{{ s.warehouse_id }}</td>
            <td class="px-4 py-3 text-slate-500">{{ s.sku_id }}</td>
            <td class="px-4 py-3 font-medium text-slate-800">{{ s.qty }}</td>
            <td class="px-4 py-3 text-slate-500">{{ s.safe_qty }}</td>
            <td class="px-4 py-3">
              <span :class="s.qty <= s.safe_qty ? 'bg-red-50 text-red-500' : 'bg-green-50 text-green-600'"
                class="px-2 py-1 rounded-full text-xs">
                {{ s.qty <= s.safe_qty ? $t('stock.warning') : $t('stock.normal') }}
              </span>
            </td>
          </tr>
          <tr v-if="!stocks.length">
            <td colspan="5" class="px-4 py-12 text-center text-slate-400">{{ $t('stock.noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Simple inbound/outbound modal -->
    <div v-if="showInbound || showOutbound"
      class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <div class="bg-white rounded-2xl p-6 w-80 shadow-xl">
        <h3 class="font-semibold text-slate-800 mb-4">{{ showInbound ? $t('stock.inbound') : $t('stock.outbound') }}</h3>
        <div class="space-y-3">
          <input v-model.number="opForm.warehouse_id" type="number" :placeholder="$t('stock.warehouseId')"
            class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          <input v-model.number="opForm.sku_id" type="number" :placeholder="$t('stock.skuId')"
            class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          <input v-model.number="opForm.qty" type="number" :placeholder="$t('stock.quantity')"
            class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
        </div>
        <div class="flex gap-3 mt-4">
          <button @click="submitOp" class="flex-1 bg-blue-700 text-white py-2 rounded-xl text-sm hover:bg-blue-600">{{ $t('common.confirm') }}</button>
          <button @click="showInbound = showOutbound = false"
            class="flex-1 bg-slate-100 text-slate-600 py-2 rounded-xl text-sm hover:bg-slate-200">{{ $t('common.cancel') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getWarehouses, getStocks, doInbound, doOutbound } from '@/api/plugins'

const warehouses = ref<any[]>([])
const stocks = ref<any[]>([])
const warehouseID = ref('')
const showInbound = ref(false)
const showOutbound = ref(false)
const opForm = ref({ warehouse_id: 0, sku_id: 0, qty: 0 })

async function loadStocks() {
  const data: any = await getStocks({ warehouse_id: warehouseID.value || undefined, page: 1, size: 50 })
  stocks.value = data.list || []
}

async function submitOp() {
  if (showInbound.value) await doInbound(opForm.value)
  else await doOutbound(opForm.value)
  showInbound.value = showOutbound.value = false
  opForm.value = { warehouse_id: 0, sku_id: 0, qty: 0 }
  loadStocks()
}

onMounted(() => {
  getWarehouses().then((d: any) => warehouses.value = d || [])
  loadStocks()
})
</script>
