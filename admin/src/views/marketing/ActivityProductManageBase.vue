<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ pageTitle }}</h2>
      <button
        @click="openCreate"
        :disabled="!selectedActivityID"
        class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition disabled:opacity-50"
      >
        {{ $t('common.create') }}
      </button>
    </div>

    <div class="bg-white rounded-xl border border-slate-100 p-4 mb-4">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-3">
        <select v-model.number="selectedActivityID" class="border border-slate-200 rounded-lg px-3 py-2 text-sm" @change="onSelectActivity">
          <option :value="0">{{ $t('marketingProduct.selectActivity') }}</option>
          <option v-for="act in activities" :key="act.id" :value="act.id">
            {{ act.name }} ({{ formatDate(act.start_at) }} ~ {{ formatDate(act.end_at) }})
          </option>
        </select>
        <input v-model="keyword" :placeholder="$t('common.search')" class="border border-slate-200 rounded-lg px-3 py-2 text-sm" />
        <button class="px-4 py-2 bg-slate-100 rounded-lg text-sm hover:bg-slate-200" @click="loadProducts">{{ $t('common.search') }}</button>
      </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.id') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.name') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">SKU</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ priceHeader }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('marketingProduct.limitPerOrder') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('marketingProduct.totalStockLimit') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('marketingProduct.soldQty') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="row in list" :key="row.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-slate-400">{{ row.id }}</td>
            <td class="px-4 py-3 text-slate-700">{{ row.product_title || row.product_id }}</td>
            <td class="px-4 py-3 text-slate-600">{{ formatSku(row.sku_attrs, row.sku_id) }}</td>
            <td class="px-4 py-3 text-slate-600">
              <template v-if="kind !== 'bargain'">¥{{ row.activity_price }}</template>
              <template v-else>¥{{ row.start_price }} / ¥{{ row.floor_price }}</template>
            </td>
            <td class="px-4 py-3 text-slate-600">{{ row.limit_per_order || '-' }}</td>
            <td class="px-4 py-3 text-slate-600">{{ row.total_stock_limit || '-' }}</td>
            <td class="px-4 py-3 text-slate-600">{{ row.sold_qty || 0 }}</td>
            <td class="px-4 py-3">
              <button class="text-blue-600 hover:underline text-xs mr-3" @click="openEdit(row)">{{ $t('common.edit') }}</button>
              <button class="text-red-500 hover:underline text-xs" @click="removeRow(row.id)">{{ $t('common.delete') }}</button>
            </td>
          </tr>
          <tr v-if="!list.length">
            <td colspan="8" class="px-4 py-12 text-center text-slate-400">{{ $t('common.noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showForm" class="fixed inset-0 bg-black/40 flex items-start justify-end z-50">
      <div class="bg-white w-[460px] h-full shadow-2xl p-6 overflow-y-auto">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-lg font-semibold text-slate-800">{{ editingID ? $t('common.edit') : $t('common.create') }}</h3>
          <button @click="showForm = false" class="text-slate-400 hover:text-slate-600">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('marketingProduct.selectProduct') }}</label>
            <div v-if="selectedProductID > 0" class="flex items-center gap-3 p-3 border border-slate-200 rounded-xl mb-2">
              <img v-if="selectedProductCover" :src="selectedProductCover" class="w-10 h-10 rounded-lg object-cover flex-shrink-0" />
              <div v-else class="w-10 h-10 rounded-lg bg-slate-100 flex-shrink-0" />
              <span class="text-sm text-slate-700 flex-1 truncate">{{ selectedProductTitle }}</span>
            </div>
            <button
              type="button"
              :disabled="editingID > 0"
              @click="showProductSelector = true"
              class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm text-left text-slate-500 hover:border-blue-400 hover:text-blue-600 transition disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ selectedProductID > 0 ? $t('marketingProduct.changeProduct') : $t('marketingProduct.selectProductBtn') }}
            </button>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('marketingProduct.selectSku') }}</label>
            <select
              v-model.number="selectedSkuID"
              class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm"
              :disabled="editingID > 0 || !skuOptions.length"
            >
              <option :value="0">{{ $t('marketingProduct.selectSku') }}</option>
              <option v-for="sku in skuOptions" :key="sku.id" :value="sku.id">{{ sku.label }}</option>
            </select>
          </div>
          <div v-if="kind !== 'bargain'">
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('marketingProduct.activityPrice') }}</label>
            <input v-model.number="form.activity_price" type="number" min="0.01" step="0.01" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          </div>
          <template v-else>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('marketingProduct.startPrice') }}</label>
              <input v-model.number="form.start_price" type="number" min="0.01" step="0.01" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('marketingProduct.floorPrice') }}</label>
              <input v-model.number="form.floor_price" type="number" min="0.01" step="0.01" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
            </div>
          </template>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('marketingProduct.limitPerOrder') }}</label>
            <input v-model.number="form.limit_per_order" type="number" min="0" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('marketingProduct.totalStockLimit') }}</label>
            <input v-model.number="form.total_stock_limit" type="number" min="0" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          </div>
          <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
          <button @click="save" class="w-full bg-blue-700 text-white py-3 rounded-xl text-sm font-medium hover:bg-blue-600 transition">
            {{ $t('common.save') }}
          </button>
        </div>
      </div>
    </div>

    <ProductSelectorModal
      :visible="showProductSelector"
      @close="showProductSelector = false"
      @select="onProductSelected"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'
import ProductSelectorModal from '@/components/common/ProductSelectorModal.vue'

const props = defineProps<{ kind: 'seckill' | 'group-buy' | 'bargain' }>()
const { t } = useI18n()

type SkuOption = { id: number; label: string; price: number; stock: number; attrs: any[] }

const activities = ref<any[]>([])
const list = ref<any[]>([])
const selectedActivityID = ref(0)
const keyword = ref('')
const showForm = ref(false)
const editingID = ref(0)
const error = ref('')
const selectedProductID = ref(0)
const selectedProductTitle = ref('')
const selectedProductCover = ref('')
const selectedSkuID = ref(0)
const skuOptions = ref<SkuOption[]>([])
const showProductSelector = ref(false)

const form = reactive<any>({
  product_id: 0,
  sku_id: 0,
  activity_price: 0,
  start_price: 0,
  floor_price: 0,
  limit_per_order: 0,
  total_stock_limit: 0,
})

const pageTitle = computed(() => {
  if (props.kind === 'seckill') return t('nav.seckillProductManage')
  if (props.kind === 'group-buy') return t('nav.groupBuyProductManage')
  return t('nav.bargainProductManage')
})

const priceHeader = computed(() => {
  if (props.kind === 'bargain') return t('marketingProduct.bargainPrice')
  return t('marketingProduct.activityPrice')
})

const activityEndpoint = computed(() => {
  if (props.kind === 'seckill') return '/marketing/seckill/activities'
  if (props.kind === 'group-buy') return '/marketing/group-buy/activities'
  return '/marketing/bargain/activities'
})

const productEndpoint = computed(() => {
  if (props.kind === 'seckill') return '/marketing/seckill/products'
  if (props.kind === 'group-buy') return '/marketing/group-buy/products'
  return '/marketing/bargain/products'
})

function formatDate(value?: string) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

function parseSkuAttrs(raw: any): any[] {
  if (Array.isArray(raw)) return raw
  if (typeof raw !== 'string' || !raw.trim()) return []
  try {
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function formatSku(attrs: any, fallbackID?: number) {
  const list = parseSkuAttrs(attrs)
  if (!list.length) return String(fallbackID || '-')
  return list.map((item: any) => `${item.name}:${item.value}`).join(' / ')
}

function buildSkuLabel(attrs: any[], price: number, stock: number) {
  const attrLabel = attrs.length ? attrs.map((item: any) => `${item.name}:${item.value}`).join(' / ') : '默认规格'
  return `${attrLabel} ｜ ¥${price} ｜ 库存 ${stock}`
}

async function loadActivities() {
  const data: any = await request.get(activityEndpoint.value)
  activities.value = data?.list || []
  if (!selectedActivityID.value && activities.value.length) {
    selectedActivityID.value = Number(activities.value[0].id || 0)
  }
}

async function loadSkuOptions(productID: number) {
  skuOptions.value = []
  selectedSkuID.value = 0
  if (productID <= 0) return
  const detail: any = await request.get(`/products/${productID}`)
  const skus = Array.isArray(detail?.skus) ? detail.skus : []
  skuOptions.value = skus.map((item: any) => {
    const attrs = parseSkuAttrs(item?.attrs)
    return {
      id: Number(item?.id || 0),
      label: buildSkuLabel(attrs, Number(item?.price || 0), Number(item?.stock || 0)),
      price: Number(item?.price || 0),
      stock: Number(item?.stock || 0),
      attrs,
    }
  }).filter((item: SkuOption) => item.id > 0)
}

async function loadProducts() {
  if (!selectedActivityID.value) {
    list.value = []
    return
  }
  const data: any = await request.get(productEndpoint.value, {
    params: {
      activity_id: selectedActivityID.value,
      keyword: keyword.value || undefined,
    },
  })
  list.value = data?.list || []
}

function onSelectActivity() {
  loadProducts()
}

async function onProductSelected(product: any) {
  selectedProductID.value = Number(product?.id || 0)
  selectedProductTitle.value = String(product?.title || '')
  selectedProductCover.value = String(product?.cover || '')
  form.product_id = selectedProductID.value
  form.sku_id = 0
  await loadSkuOptions(form.product_id)
}

function resetForm() {
  form.product_id = 0
  form.sku_id = 0
  form.activity_price = 0
  form.start_price = 0
  form.floor_price = 0
  form.limit_per_order = 0
  form.total_stock_limit = 0
  selectedProductID.value = 0
  selectedProductTitle.value = ''
  selectedProductCover.value = ''
  selectedSkuID.value = 0
  skuOptions.value = []
  error.value = ''
}

function openCreate() {
  editingID.value = 0
  resetForm()
  showForm.value = true
}

async function openEdit(row: any) {
  editingID.value = Number(row?.id || 0)
  form.product_id = Number(row?.product_id || 0)
  form.sku_id = Number(row?.sku_id || 0)
  form.activity_price = Number(row?.activity_price || 0)
  form.start_price = Number(row?.start_price || 0)
  form.floor_price = Number(row?.floor_price || 0)
  form.limit_per_order = Number(row?.limit_per_order || 0)
  form.total_stock_limit = Number(row?.total_stock_limit || 0)
  selectedProductID.value = form.product_id
  selectedProductTitle.value = String(row?.product_title || '')
  selectedProductCover.value = String(row?.product_cover || '')
  selectedSkuID.value = form.sku_id
  skuOptions.value = [{
    id: form.sku_id,
    label: formatSku(row?.sku_attrs, row?.sku_id),
    price: Number(row?.sku_price || 0),
    stock: Number(row?.sku_stock || 0),
    attrs: parseSkuAttrs(row?.sku_attrs),
  }]
  error.value = ''
  showForm.value = true
}

function validateForm() {
  form.product_id = Number(selectedProductID.value || 0)
  form.sku_id = Number(selectedSkuID.value || 0)
  if (!selectedActivityID.value) return t('marketingProduct.selectActivity')
  if (form.product_id <= 0 || form.sku_id <= 0) return t('marketingProduct.invalidProductSku')
  if (props.kind !== 'bargain') {
    if (form.activity_price <= 0) return t('marketingProduct.invalidActivityPrice')
  } else {
    if (form.start_price <= 0 || form.floor_price <= 0) return t('marketingProduct.invalidBargainPrice')
    if (form.floor_price > form.start_price) return t('marketingProduct.invalidBargainPrice')
  }
  if (form.limit_per_order < 0 || form.total_stock_limit < 0) return t('marketingProduct.invalidLimits')
  return ''
}

function toPayloadRows(rows: any[]) {
  return rows.map((item: any) => ({
    product_id: Number(item.product_id || 0),
    sku_id: Number(item.sku_id || 0),
    activity_price: Number(item.activity_price || 0),
    start_price: Number(item.start_price || 0),
    floor_price: Number(item.floor_price || 0),
    limit_per_order: Number(item.limit_per_order || 0),
    total_stock_limit: Number(item.total_stock_limit || 0),
  }))
}

async function persistRows(nextRows: any[]) {
  await request.put(`${activityEndpoint.value}/${selectedActivityID.value}/products`, toPayloadRows(nextRows))
}

async function save() {
  const message = validateForm()
  if (message) {
    error.value = message
    return
  }
  const duplicate = list.value.find((item: any) =>
    Number(item.id || 0) !== editingID.value &&
    Number(item.product_id || 0) === Number(form.product_id || 0) &&
    Number(item.sku_id || 0) === Number(form.sku_id || 0),
  )
  if (duplicate) {
    error.value = t('marketingProduct.duplicateSku')
    return
  }
  const nextRows = list.value.filter((item: any) => Number(item.id || 0) !== editingID.value)
  nextRows.push({
    ...form,
    product_id: Number(form.product_id || 0),
    sku_id: Number(form.sku_id || 0),
  })
  try {
    await persistRows(nextRows)
    showForm.value = false
    await loadProducts()
  } catch (err: any) {
    error.value = err?.message || t('common.saveFailed')
  }
}

async function removeRow(id: number) {
  const nextRows = list.value.filter((item: any) => Number(item.id || 0) !== Number(id || 0))
  await persistRows(nextRows)
  await loadProducts()
}

onMounted(async () => {
  await loadActivities()
  await loadProducts()
})
</script>
