<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">秒杀商品管理</h2>
      <button @click="openAddProduct" class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition">添加商品</button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">商品信息</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">秒杀价</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">限购</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">库存/已售</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="item in list" :key="item.id" class="hover:bg-slate-50">
            <td class="px-4 py-3">
              <div class="font-medium text-slate-800">商品ID: {{ item.product_id }}</div>
              <div class="text-xs text-slate-500">SKU ID: {{ item.sku_id || '全部' }}</div>
            </td>
            <td class="px-4 py-3 text-slate-600 font-medium">¥{{ item.seckill_price }}</td>
            <td class="px-4 py-3 text-slate-600">{{ item.limit_per_order || '不限' }}</td>
            <td class="px-4 py-3 text-slate-600">{{ item.total_stock_limit || '不限' }} / {{ item.sold_qty }}</td>
            <td class="px-4 py-3">
              <button @click="openEdit(item)" class="text-blue-600 hover:text-blue-700 mr-3">编辑</button>
              <button @click="remove(item.id)" class="text-red-500 hover:text-red-700">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 商品编辑对话框 -->
    <div v-if="visible" class="fixed inset-0 bg-black/30 flex items-center justify-center z-50" @click.self="visible=false">
      <div class="bg-white rounded-2xl p-6 w-[500px] shadow-xl">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">{{ form.id ? '编辑商品' : '添加商品' }}</h3>
        <div class="space-y-3">
          <div>
            <label class="text-sm text-slate-600 mb-1 block">商品ID</label>
            <input v-model.number="form.product_id" type="number" placeholder="请输入商品ID" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">SKU ID（0=全部SKU）</label>
            <input v-model.number="form.sku_id" type="number" placeholder="0" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">秒杀价格</label>
            <input v-model.number="form.seckill_price" type="number" step="0.01" placeholder="请输入秒杀价格" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">单笔限购（0=不限）</label>
            <input v-model.number="form.limit_per_order" type="number" placeholder="0" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">活动库存（0=不限）</label>
            <input v-model.number="form.total_stock_limit" type="number" placeholder="0" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="visible=false" class="flex-1 border border-slate-200 rounded-xl py-2.5 text-sm text-slate-600">取消</button>
          <button @click="save" class="flex-1 bg-blue-700 text-white rounded-xl py-2.5 text-sm">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import request from '@/api/request'
import { confirmAction } from '@/utils/dialog'

const route = useRoute()
const activityId = ref(Number(route.query.activity_id))
const list = ref<any[]>([])
const visible = ref(false)
const form = ref<any>({
  id: 0,
  product_id: 0,
  sku_id: 0,
  seckill_price: 0,
  limit_per_order: 0,
  total_stock_limit: 0,
})

async function load() {
  const data: any = await request.get('/seckill/products', { params: { activity_id: activityId.value } })
  list.value = data.list || []
}

function openAddProduct() {
  form.value = {
    id: 0,
    product_id: 0,
    sku_id: 0,
    seckill_price: 0,
    limit_per_order: 0,
    total_stock_limit: 0,
  }
  visible.value = true
}

function openEdit(row: any) {
  form.value = { ...row }
  visible.value = true
}

async function save() {
  // 批量保存所有商品
  const products = form.value.id
    ? list.value.map(p => p.id === form.value.id ? form.value : p)
    : [...list.value, form.value]

  await request.put(`/seckill/activities/${activityId.value}/products`, products)
  visible.value = false
  load()
}

async function remove(id: number) {
  if (!confirmAction('确定要删除该商品吗？')) return
  const products = list.value.filter(p => p.id !== id)
  await request.put(`/seckill/activities/${activityId.value}/products`, products)
  load()
}

onMounted(load)
</script>
