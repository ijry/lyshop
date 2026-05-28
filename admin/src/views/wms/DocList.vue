<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ pageTitle }}</h2>
      <button
        class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition"
        @click="createDraftDoc"
      >
        {{ $t('wms.createDraftDoc') }}
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm p-4 mb-4 flex gap-3 border border-slate-100 flex-wrap">
      <select
        v-model="query.type"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none"
        :disabled="!!fixedType"
      >
        <option value="">{{ $t('wms.allDocTypes') }}</option>
        <option value="inbound">{{ $t('wms.docTypeInbound') }}</option>
        <option value="outbound">{{ $t('wms.docTypeOutbound') }}</option>
      </select>
      <select v-model="query.status" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none">
        <option value="">{{ $t('wms.allDocStatuses') }}</option>
        <option value="draft">{{ $t('wms.docStatusDraft') }}</option>
        <option value="completed">{{ $t('wms.docStatusCompleted') }}</option>
        <option value="canceled">{{ $t('wms.docStatusCanceled') }}</option>
      </select>
      <select v-model.number="query.warehouse_id" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none">
        <option :value="0">{{ $t('stock.allWarehouse') }}</option>
        <option v-for="item in warehouses" :key="item.id" :value="item.id">{{ item.name }}</option>
      </select>
      <input
        v-model.trim="query.doc_no"
        :placeholder="$t('wms.docNoPlaceholder')"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm flex-1 min-w-[180px] focus:outline-none focus:border-blue-400"
      />
      <button class="px-4 py-2 bg-slate-100 rounded-lg text-sm hover:bg-slate-200" @click="onSearch">{{ $t('common.search') }}</button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.docNo') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.docType') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.warehouse') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.totalQty') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.status') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.updatedAt') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="row in docs" :key="row.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 font-mono text-xs text-slate-600">{{ row.doc_no }}</td>
            <td class="px-4 py-3 text-slate-700">{{ typeLabel(row.type) }}</td>
            <td class="px-4 py-3 text-slate-700">{{ row.warehouse_name || row.warehouse_id || '-' }}</td>
            <td class="px-4 py-3 text-slate-700">{{ row.total_qty }}</td>
            <td class="px-4 py-3">
              <span :class="statusClass(row.status)" class="px-2 py-1 rounded-full text-xs">
                {{ statusLabel(row.status) }}
              </span>
            </td>
            <td class="px-4 py-3 text-slate-500">{{ formatDate(row.updated_at || row.created_at) }}</td>
            <td class="px-4 py-3">
              <button class="text-blue-600 hover:underline text-xs mr-3" @click="goView(row.id)">{{ $t('common.view') }}</button>
              <button
                v-if="row.status === 'draft'"
                class="text-emerald-600 hover:underline text-xs"
                @click="goEdit(row.id)"
              >
                {{ $t('common.edit') }}
              </button>
            </td>
          </tr>
          <tr v-if="!docs.length">
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
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { listDocs, listWarehouses, type WmsDoc, type WmsDocStatus, type WmsDocType, type WmsWarehouse } from '@/api/wms'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const docs = ref<WmsDoc[]>([])
const warehouses = ref<WmsWarehouse[]>([])
const total = ref(0)
const prevFixedType = ref<WmsDocType | ''>('')
const query = ref({
  type: '',
  status: '',
  warehouse_id: 0,
  doc_no: '',
  page: 1,
  size: 20,
})

const fixedType = computed<WmsDocType | ''>(() => {
  if (route.path === '/wms/inbound') return 'inbound'
  if (route.path === '/wms/outbound') return 'outbound'
  return ''
})

const pageTitle = computed(() => {
  if (fixedType.value === 'inbound') return t('nav.inboundDocs')
  if (fixedType.value === 'outbound') return t('nav.outboundDocs')
  return t('wms.docListTitle')
})

function applyFixedType() {
  if (fixedType.value) {
    query.value.type = fixedType.value
  } else if (prevFixedType.value) {
    query.value.type = ''
  }
  prevFixedType.value = fixedType.value
}

function statusLabel(status: WmsDocStatus) {
  if (status === 'draft') return t('wms.docStatusDraft')
  if (status === 'completed') return t('wms.docStatusCompleted')
  return t('wms.docStatusCanceled')
}

function statusClass(status: WmsDocStatus) {
  if (status === 'draft') return 'bg-yellow-50 text-yellow-700'
  if (status === 'completed') return 'bg-green-50 text-green-600'
  return 'bg-slate-100 text-slate-500'
}

function typeLabel(type: WmsDocType) {
  return type === 'inbound' ? t('wms.docTypeInbound') : t('wms.docTypeOutbound')
}

function formatDate(value?: string) {
  return value ? String(value).slice(0, 19).replace('T', ' ') : '-'
}

async function loadWarehouses() {
  const data = await listWarehouses()
  warehouses.value = Array.isArray(data?.list) ? data.list : []
}

async function loadDocs() {
  const data = await listDocs({
    doc_type: (query.value.type as WmsDocType) || undefined,
    status: (query.value.status as WmsDocStatus) || undefined,
    warehouse_id: query.value.warehouse_id > 0 ? Number(query.value.warehouse_id) : undefined,
    doc_no: query.value.doc_no || undefined,
    page: query.value.page,
    size: query.value.size,
  })
  docs.value = Array.isArray(data?.list) ? data.list : []
  total.value = Number(data?.total || 0)
}

function onSearch() {
  query.value.page = 1
  loadDocs()
}

async function createDraftDoc() {
  const warehouseID = query.value.warehouse_id || Number(warehouses.value[0]?.id || 0)
  const type = (fixedType.value || query.value.type || 'inbound') as WmsDocType
  router.push({
    path: '/wms/docs/new',
    query: {
      type,
      warehouse_id: warehouseID || undefined,
      back_to: route.fullPath,
    },
  })
}

function goEdit(id: number) {
  router.push({
    path: `/wms/docs/${id}`,
    query: { back_to: route.fullPath },
  })
}

function goView(id: number) {
  router.push({
    path: `/wms/docs/${id}`,
    query: { mode: 'view', back_to: route.fullPath },
  })
}

function prevPage() {
  if (query.value.page <= 1) return
  query.value.page -= 1
  loadDocs()
}

function nextPage() {
  if (query.value.page * query.value.size >= total.value) return
  query.value.page += 1
  loadDocs()
}

watch(
  () => route.path,
  () => {
    query.value.page = 1
    applyFixedType()
    loadDocs()
  },
)

onMounted(async () => {
  applyFixedType()
  await loadWarehouses()
  await loadDocs()
})
</script>
