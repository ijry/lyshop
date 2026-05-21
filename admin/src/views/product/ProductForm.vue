<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <router-link to="/product/list" class="text-slate-400 hover:text-slate-600 text-sm">← 返回列表</router-link>
      <h2 class="text-xl font-semibold text-slate-800">{{ isEdit ? '编辑商品' : '新增商品' }}</h2>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6 max-w-2xl">
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">商品名称 *</label>
          <input v-model="form.title" class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" placeholder="请输入商品名称" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">售价 *</label>
            <input v-model.number="form.price" type="number" step="0.01"
              class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" placeholder="0.00" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">库存</label>
            <input v-model.number="form.stock" type="number"
              class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" placeholder="0" />
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">分类</label>
          <select v-model="form.category_id"
            class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none">
            <option value="">请选择分类</option>
            <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">状态</label>
          <select v-model.number="form.status"
            class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none">
            <option :value="1">上架</option>
            <option :value="0">下架</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">封面图URL</label>
          <input v-model="form.cover" class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" placeholder="https://..." />
        </div>
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">商品详情</label>
          <textarea v-model="form.detail" rows="5"
            class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400 resize-none" placeholder="商品详情描述..." />
        </div>
        <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
        <div class="flex gap-3 pt-2">
          <button @click="save" :disabled="saving"
            class="px-6 py-3 bg-blue-700 text-white rounded-xl text-sm font-medium hover:bg-blue-600 transition disabled:opacity-60">
            {{ saving ? '保存中...' : '保 存' }}
          </button>
          <router-link to="/product/list"
            class="px-6 py-3 bg-slate-100 text-slate-600 rounded-xl text-sm font-medium hover:bg-slate-200 transition">
            取 消
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getProduct, createProduct, updateProduct, getCategories } from '@/api/plugins'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const categories = ref<any[]>([])
const saving = ref(false)
const error = ref('')

const form = ref({
  title: '', price: 0, origin_price: 0, stock: 0,
  category_id: '', cover: '', detail: '', status: 1
})

onMounted(async () => {
  getCategories().then((d: any) => categories.value = d || [])
  if (isEdit.value) {
    const data: any = await getProduct(Number(route.params.id))
    Object.assign(form.value, data)
  }
})

async function save() {
  if (!form.value.title) { error.value = '商品名称不能为空'; return }
  saving.value = true; error.value = ''
  try {
    if (isEdit.value) {
      await updateProduct(Number(route.params.id), form.value)
    } else {
      await createProduct({ product: form.value, skus: [], images: [] })
    }
    router.push('/product/list')
  } catch (e: any) {
    error.value = e.message
  } finally {
    saving.value = false
  }
}
</script>
