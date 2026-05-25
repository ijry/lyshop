<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">商品列表</h2>
      <router-link to="/product/form"
        class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition">
        + 新增商品
      </router-link>
    </div>

    <!-- Search bar -->
    <div class="bg-white rounded-xl shadow-sm p-4 mb-4 flex gap-3 border border-slate-100">
      <input v-model="query.keyword" placeholder="搜索商品名称"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm flex-1 focus:outline-none focus:border-blue-400" />
      <select v-model="query.category_id"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none">
        <option value="">全部分类</option>
        <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
      </select>
      <button @click="loadProducts" class="px-4 py-2 bg-slate-100 rounded-lg text-sm hover:bg-slate-200">搜索</button>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">ID</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">商品</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">价格</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">库存</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">收藏数</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">状态</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="p in products" :key="p.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-slate-400">{{ p.id }}</td>
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <img v-if="p.cover" :src="p.cover" class="w-10 h-10 rounded-lg object-cover" />
                <div class="w-10 h-10 rounded-lg bg-slate-100 flex items-center justify-center text-slate-400 text-xs" v-else>图</div>
                <span class="font-medium text-slate-700 truncate max-w-xs">{{ p.title }}</span>
              </div>
            </td>
            <td class="px-4 py-3 text-slate-700">¥{{ p.price }}</td>
            <td class="px-4 py-3 text-slate-700">{{ p.stock }}</td>
            <td class="px-4 py-3 text-slate-700">{{ p.favorite_count || 0 }}</td>
            <td class="px-4 py-3">
              <span :class="p.status===1 ? 'text-green-600 bg-green-50' : 'text-slate-400 bg-slate-100'"
                class="px-2 py-1 rounded-full text-xs">
                {{ p.status === 1 ? '上架' : '下架' }}
              </span>
            </td>
            <td class="px-4 py-3">
              <router-link :to="`/product/form/${p.id}`"
                class="text-blue-600 hover:underline text-xs mr-3">编辑</router-link>
              <button
                v-if="canManageReview"
                @click="openReview(p)"
                class="text-emerald-600 hover:underline text-xs mr-3"
              >
                管理评价
              </button>
              <button @click="remove(p.id)" class="text-red-500 hover:underline text-xs">删除</button>
            </td>
          </tr>
          <tr v-if="!products.length">
            <td colspan="7" class="px-4 py-12 text-center text-slate-400">暂无商品</td>
          </tr>
        </tbody>
      </table>
      <!-- Pagination -->
      <div class="px-4 py-3 flex items-center justify-between border-t border-slate-100 text-sm text-slate-500">
        <span>共 {{ total }} 条</span>
        <div class="flex gap-2">
          <button :disabled="query.page <= 1" @click="query.page--; loadProducts()"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40">上一页</button>
          <button :disabled="query.page * query.size >= total" @click="query.page++; loadProducts()"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40">下一页</button>
        </div>
      </div>
    </div>

    <div v-if="showReviewModal && reviewProduct" class="fixed inset-0 bg-black/35 flex items-center justify-center z-50" @click.self="closeReview">
      <div class="bg-white rounded-xl w-[1120px] max-w-[96vw] max-h-[92vh] overflow-auto p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-base font-semibold text-slate-800">管理评价 - {{ reviewProduct.title }}（ID: {{ reviewProduct.id }}）</h3>
          <button class="text-slate-400 hover:text-slate-600" @click="closeReview">关闭</button>
        </div>
        <ReviewManager
          :show-title="false"
          :show-product-filter="false"
          :fixed-product-id="Number(reviewProduct.id || 0)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getProducts, getCategories, deleteProduct } from '@/api/plugins'
import {
  getMenus,
  type AdminMenuGroupedResponse,
  type AdminMenuItem,
  type AdminMenuResponse,
} from '@/api/auth'
import ReviewManager from '@/components/review/ReviewManager.vue'
import { useAuthStore } from '@/stores/auth'
import { confirmAction } from '@/utils/dialog'

const products = ref<any[]>([])
const categories = ref<any[]>([])
const total = ref(0)
const query = ref({ keyword: '', category_id: '', page: 1, size: 20 })
const showReviewModal = ref(false)
const reviewProduct = ref<any>(null)
const canManageReview = ref(false)
const auth = useAuthStore()

function hasReviewMenu(rows: any[]): boolean {
  const list = Array.isArray(rows) ? rows : []
  for (const row of list) {
    if (String(row?.path || '') === '/review/list') return true
    if (hasReviewMenu(row?.children || [])) return true
  }
  return false
}

function isGroupedResponse(data: AdminMenuResponse): data is AdminMenuGroupedResponse {
  return !!data && !Array.isArray(data) && Array.isArray((data as AdminMenuGroupedResponse).groups)
}

function flattenMenus(rows: AdminMenuItem[]): AdminMenuItem[] {
  const list = Array.isArray(rows) ? rows : []
  return list.flatMap((row) => [row, ...flattenMenus(row.children || [])])
}

async function loadProducts() {
  const data: any = await getProducts(query.value)
  products.value = data.list || []
  total.value = data.total || 0
}

async function remove(id: number) {
  if (!confirmAction('确认删除该商品？')) return
  await deleteProduct(id)
  loadProducts()
}

function openReview(product: any) {
  if (!canManageReview.value) return
  reviewProduct.value = product
  showReviewModal.value = true
}

function closeReview() {
  showReviewModal.value = false
  reviewProduct.value = null
}

onMounted(async () => {
  getCategories().then((d: any) => categories.value = d || [])
  const canByPermission = auth.hasPermission('order:view')
  try {
    const data = await getMenus()
    const menuRows = isGroupedResponse(data)
      ? flattenMenus(data.groups.flatMap((group) => group.menus || []))
      : (Array.isArray(data) ? data : [])
    canManageReview.value = canByPermission || hasReviewMenu(menuRows as any[])
  } catch {
    canManageReview.value = canByPermission
  }
  loadProducts()
})
</script>
