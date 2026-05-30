<template>
  <div v-if="visible" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50" @click.self="close">
    <div class="bg-white rounded-xl shadow-2xl w-[760px] max-w-[96vw] max-h-[88vh] flex flex-col">
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-100">
        <h3 class="text-base font-semibold text-slate-800">{{ $t('marketingProduct.selectProduct') }}</h3>
        <button @click="close" class="text-slate-400 hover:text-slate-600 text-lg leading-none">✕</button>
      </div>

      <!-- Search bar -->
      <div class="px-6 py-3 border-b border-slate-100 flex gap-3">
        <input
          v-model="keyword"
          :placeholder="$t('common.search')"
          class="border border-slate-200 rounded-lg px-3 py-2 text-sm flex-1 focus:outline-none focus:border-blue-400"
          @keyup.enter="search"
        />
        <select
          v-model="categoryId"
          class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none"
          @change="search"
        >
          <option value="">{{ $t('product.list.allCategory') }}</option>
          <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>
        <button @click="search" class="px-4 py-2 bg-slate-100 rounded-lg text-sm hover:bg-slate-200">
          {{ $t('common.search') }}
        </button>
      </div>

      <!-- Product list -->
      <div class="flex-1 overflow-y-auto">
        <table class="w-full text-sm">
          <thead class="bg-slate-50 border-b border-slate-100 sticky top-0">
            <tr>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.id') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('product.list.product') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('product.list.price') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('product.list.stock') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('product.list.favorites') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.status') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-50">
            <tr
              v-for="p in products"
              :key="p.id"
              class="hover:bg-blue-50 cursor-pointer"
              @click="select(p)"
            >
              <td class="px-4 py-3 text-slate-400">{{ p.id }}</td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <img v-if="p.cover" :src="p.cover" class="w-10 h-10 rounded-lg object-cover flex-shrink-0" />
                  <div v-else class="w-10 h-10 rounded-lg bg-slate-100 flex items-center justify-center text-slate-400 text-xs flex-shrink-0">
                    {{ $t('product.list.image') }}
                  </div>
                  <span class="font-medium text-slate-700 truncate max-w-[200px]">{{ p.title }}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-slate-700">¥{{ p.price }}</td>
              <td class="px-4 py-3 text-slate-700">{{ p.stock }}</td>
              <td class="px-4 py-3 text-slate-700">{{ p.favorite_count || 0 }}</td>
              <td class="px-4 py-3">
                <span
                  :class="p.status === 1 ? 'text-green-600 bg-green-50' : 'text-slate-400 bg-slate-100'"
                  class="px-2 py-1 rounded-full text-xs"
                >
                  {{ p.status === 1 ? $t('product.list.onSale') : $t('product.list.offSale') }}
                </span>
              </td>
            </tr>
            <tr v-if="!loading && !products.length">
              <td colspan="6" class="px-4 py-12 text-center text-slate-400">{{ $t('common.noData') }}</td>
            </tr>
            <tr v-if="loading">
              <td colspan="6" class="px-4 py-12 text-center text-slate-400">{{ $t('common.loading') }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="px-6 py-3 border-t border-slate-100 flex items-center justify-between text-sm text-slate-500">
        <span>{{ $t('common.totalCount', { total }) }}</span>
        <div class="flex gap-2">
          <button
            :disabled="page <= 1"
            @click="page--; load()"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40"
          >{{ $t('common.prevPage') }}</button>
          <button
            :disabled="page * pageSize >= total"
            @click="page++; load()"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40"
          >{{ $t('common.nextPage') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { getProducts, getCategories } from '@/api/plugins'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'select', product: any): void
}>()

const keyword = ref('')
const categoryId = ref('')
const page = ref(1)
const pageSize = 10
const total = ref(0)
const products = ref<any[]>([])
const categories = ref<any[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const data: any = await getProducts({
      keyword: keyword.value || undefined,
      category_id: categoryId.value || undefined,
      page: page.value,
      size: pageSize,
    })
    products.value = data?.list || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  load()
}

function select(product: any) {
  emit('select', product)
  emit('close')
}

function close() {
  emit('close')
}

watch(() => props.visible, (val) => {
  if (val) {
    keyword.value = ''
    categoryId.value = ''
    page.value = 1
    if (!categories.value.length) {
      getCategories().then((d: any) => { categories.value = d || [] })
    }
    load()
  }
})
</script>
