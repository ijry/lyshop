<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <button class="text-slate-400 hover:text-slate-600 text-sm" @click="goBack">{{ $t('common.back') }}</button>
        <h2 class="text-xl font-semibold text-slate-800">{{ $t('wms.docEditorTitle') }}</h2>
      </div>
      <div class="flex items-center gap-2">
        <button
          v-if="editable"
          class="px-4 py-2 bg-slate-100 text-slate-700 text-sm rounded-lg hover:bg-slate-200"
          @click="handleSave"
        >
          {{ $t('wms.saveDraft') }}
        </button>
        <button
          v-if="editable"
          class="px-4 py-2 bg-green-600 text-white text-sm rounded-lg hover:bg-green-500"
          @click="handleComplete"
        >
          {{ $t('wms.completeDoc') }}
        </button>
        <button
          v-if="editable"
          class="px-4 py-2 bg-red-600 text-white text-sm rounded-lg hover:bg-red-500"
          @click="handleVoid"
        >
          {{ $t('wms.voidDoc') }}
        </button>
      </div>
    </div>

    <div v-if="doc" class="space-y-4">
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-5">
        <div class="grid grid-cols-1 lg:grid-cols-4 gap-3">
          <div class="text-sm text-slate-600">
            <p class="text-xs text-slate-400 mb-1">{{ $t('wms.docNo') }}</p>
            <p class="font-mono">{{ doc.doc_no || '-' }}</p>
          </div>
          <div>
            <p class="text-xs text-slate-400 mb-1">{{ $t('wms.docType') }}</p>
            <select v-model="doc.type" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" :disabled="!editable">
              <option value="inbound">{{ $t('wms.docTypeInbound') }}</option>
              <option value="outbound">{{ $t('wms.docTypeOutbound') }}</option>
            </select>
          </div>
          <div>
            <p class="text-xs text-slate-400 mb-1">{{ $t('wms.warehouse') }}</p>
            <select v-model.number="doc.warehouse_id" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" :disabled="!editable">
              <option v-for="item in warehouses" :key="item.id" :value="item.id">{{ item.name }}</option>
            </select>
          </div>
          <div class="text-sm text-slate-600">
            <p class="text-xs text-slate-400 mb-1">{{ $t('common.status') }}</p>
            <span :class="statusClass(doc.status)" class="px-2 py-1 rounded-full text-xs">
              {{ statusLabel(doc.status) }}
            </span>
          </div>
        </div>
        <div class="mt-3">
          <p class="text-xs text-slate-400 mb-1">{{ $t('wms.remark') }}</p>
          <textarea
            v-model="doc.remark"
            :disabled="!editable"
            rows="2"
            class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm resize-none"
          />
        </div>
        <p class="text-xs text-slate-400 mt-2">{{ $t('wms.updatedAt') }}: {{ formatDate(doc.updated_at || doc.created_at) }}</p>
      </div>

      <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
        <div class="px-4 py-3 border-b border-slate-100 flex items-center justify-between">
          <p class="font-medium text-slate-700">{{ $t('wms.docItems') }}</p>
          <button
            v-if="editable"
            class="px-3 py-1.5 text-xs bg-slate-100 text-slate-700 rounded-lg hover:bg-slate-200"
            @click="addItem"
          >
            {{ $t('wms.addItem') }}
          </button>
        </div>
        <table class="w-full text-sm">
          <thead class="bg-slate-50 border-b border-slate-100">
            <tr>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.skuId') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.skuName') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('stock.quantity') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.unitCost') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-50">
            <tr v-for="(item, idx) in doc.items" :key="idx" class="hover:bg-slate-50">
              <td class="px-4 py-3">
                <input v-model.number="item.sku_id" :disabled="!editable" type="number" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" />
              </td>
              <td class="px-4 py-3">
                <input v-model.trim="item.sku_name" :disabled="!editable" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" />
              </td>
              <td class="px-4 py-3">
                <input v-model.number="item.qty" :disabled="!editable" type="number" min="1" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" />
              </td>
              <td class="px-4 py-3">
                <input v-model.number="item.unit_cost" :disabled="!editable" type="number" min="0" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm" />
              </td>
              <td class="px-4 py-3">
                <button v-if="editable" class="text-red-500 hover:underline text-xs" @click="removeItem(idx)">{{ $t('common.delete') }}</button>
              </td>
            </tr>
            <tr v-if="!doc.items.length">
              <td colspan="5" class="px-4 py-10 text-center text-slate-400">{{ $t('wms.noItems') }}</td>
            </tr>
          </tbody>
        </table>
        <div class="px-4 py-3 border-t border-slate-100 text-sm text-slate-600">
          {{ $t('wms.totalQty') }}: <span class="font-semibold text-slate-800">{{ totalQty }}</span>
        </div>
      </div>

      <div v-if="!editable" class="text-sm text-slate-500 bg-slate-100 rounded-lg px-4 py-3">
        {{ $t('wms.docReadonlyHint') }}
      </div>
    </div>

    <div v-else class="bg-white rounded-xl shadow-sm border border-slate-100 p-10 text-center text-slate-400">
      {{ $t('common.loading') }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { completeDoc, getDocDetail, listWarehouses, saveDoc, voidDoc, type WmsDoc, type WmsDocItem, type WmsDocStatus, type WmsWarehouse } from '@/api/wms'
import { notify } from '@/utils/notify'
import { confirmAction } from '@/utils/dialog'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const doc = ref<WmsDoc | null>(null)
const warehouses = ref<WmsWarehouse[]>([])

const docID = computed(() => Number(route.params.id || 0))
const isViewMode = computed(() => String(route.query.mode || '') === 'view')
const editable = computed(() => {
  if (!doc.value) return false
  if (isViewMode.value) return false
  return doc.value.status === 'draft'
})
const totalQty = computed(() => {
  if (!doc.value) return 0
  return doc.value.items.reduce((sum, row) => sum + Math.max(0, Number(row.qty || 0)), 0)
})

function statusLabel(status: WmsDocStatus) {
  if (status === 'draft') return t('wms.docStatusDraft')
  if (status === 'completed') return t('wms.docStatusCompleted')
  return t('wms.docStatusVoided')
}

function statusClass(status: WmsDocStatus) {
  if (status === 'draft') return 'bg-yellow-50 text-yellow-700'
  if (status === 'completed') return 'bg-green-50 text-green-600'
  return 'bg-slate-100 text-slate-500'
}

function formatDate(value?: string) {
  return value ? String(value).slice(0, 19).replace('T', ' ') : '-'
}

function normalizeItems(items: WmsDocItem[] | undefined): WmsDocItem[] {
  if (!Array.isArray(items)) return []
  return items.map((row) => ({
    id: row.id,
    sku_id: Number(row.sku_id || 0),
    sku_name: String(row.sku_name || ''),
    qty: Math.max(0, Number(row.qty || 0)),
    unit_cost: Math.max(0, Number(row.unit_cost || 0)),
    note: String(row.note || ''),
  }))
}

async function loadWarehouses() {
  const rows = await listWarehouses()
  warehouses.value = Array.isArray(rows) ? rows : []
}

async function loadDoc() {
  if (!docID.value) {
    notify(t('wms.docNotFound'))
    return
  }
  const data = await getDocDetail(docID.value)
  if (!data) {
    notify(t('wms.docNotFound'))
    return
  }
  doc.value = {
    ...data,
    items: normalizeItems(data.items),
  }
}

function addItem() {
  if (!doc.value || !editable.value) return
  doc.value.items.push({ sku_id: 0, sku_name: '', qty: 1, unit_cost: 0 })
}

function removeItem(index: number) {
  if (!doc.value || !editable.value) return
  doc.value.items.splice(index, 1)
}

function validateBeforeSave() {
  if (!doc.value) return false
  if (!doc.value.warehouse_id) {
    notify(t('wms.warehouseRequired'))
    return false
  }
  if (!doc.value.items.length) {
    notify(t('wms.docItemsRequired'))
    return false
  }
  for (const row of doc.value.items) {
    if (!Number(row.sku_id || 0)) {
      notify(t('wms.skuIdRequired'))
      return false
    }
    if (Number(row.qty || 0) <= 0) {
      notify(t('wms.qtyRequired'))
      return false
    }
  }
  return true
}

async function handleSave() {
  if (!doc.value || !editable.value) return
  if (!validateBeforeSave()) return
  const payload: Partial<WmsDoc> = {
    type: doc.value.type,
    warehouse_id: Number(doc.value.warehouse_id || 0),
    remark: String(doc.value.remark || ''),
    items: doc.value.items.map((row) => ({
      id: row.id,
      sku_id: Number(row.sku_id || 0),
      sku_name: String(row.sku_name || ''),
      qty: Number(row.qty || 0),
      unit_cost: Number(row.unit_cost || 0),
      note: String(row.note || ''),
    })),
  }
  const saved = await saveDoc(doc.value.id, payload)
  doc.value = { ...saved, items: normalizeItems(saved.items) }
}

async function handleComplete() {
  if (!doc.value || !editable.value) return
  if (!validateBeforeSave()) return
  if (!confirmAction(t('wms.confirmCompleteDoc'))) return
  await handleSave()
  const updated = await completeDoc(doc.value.id)
  doc.value = { ...updated, items: normalizeItems(updated.items) }
}

async function handleVoid() {
  if (!doc.value || !editable.value) return
  if (!confirmAction(t('wms.confirmVoidDoc'))) return
  const updated = await voidDoc(doc.value.id)
  doc.value = { ...updated, items: normalizeItems(updated.items) }
}

function goBack() {
  router.back()
}

onMounted(async () => {
  await loadWarehouses()
  await loadDoc()
})
</script>
