<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">秒杀活动</h2>
      <button @click="openCreate" class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition">新增活动</button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 mb-4 p-4">
      <div class="flex gap-3">
        <select v-model="filters.status" @change="load" class="border border-slate-200 rounded-xl px-3 py-2 text-sm">
          <option value="-1">全部状态</option>
          <option value="1">启用</option>
          <option value="0">禁用</option>
        </select>
      </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">活动名称</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">开始时间</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">结束时间</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">状态</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="item in list" :key="item.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 font-medium text-slate-800">{{ item.name }}</td>
            <td class="px-4 py-3 text-slate-600">{{ formatDate(item.start_at) }}</td>
            <td class="px-4 py-3 text-slate-600">{{ formatDate(item.end_at) }}</td>
            <td class="px-4 py-3">
              <span v-if="item.status === 1" class="px-2 py-1 bg-green-50 text-green-600 rounded text-xs">启用</span>
              <span v-else class="px-2 py-1 bg-slate-100 text-slate-600 rounded text-xs">禁用</span>
            </td>
            <td class="px-4 py-3">
              <button @click="goToProducts(item.id)" class="text-blue-600 hover:text-blue-700 mr-3">商品管理</button>
              <button @click="openEdit(item)" class="text-blue-600 hover:text-blue-700 mr-3">编辑</button>
              <button @click="remove(item.id)" class="text-red-500 hover:text-red-700">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="total > size" class="mt-4 flex justify-center">
      <div class="flex gap-2">
        <button @click="changePage(page - 1)" :disabled="page === 1" class="px-3 py-1 border rounded">上一页</button>
        <span class="px-3 py-1">{{ page }} / {{ Math.ceil(total / size) }}</span>
        <button @click="changePage(page + 1)" :disabled="page >= Math.ceil(total / size)" class="px-3 py-1 border rounded">下一页</button>
      </div>
    </div>

    <!-- 编辑对话框 -->
    <div v-if="visible" class="fixed inset-0 bg-black/30 flex items-center justify-center z-50" @click.self="visible=false">
      <div class="bg-white rounded-2xl p-6 w-[500px] shadow-xl">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">{{ form.id ? '编辑活动' : '新增活动' }}</h3>
        <div class="space-y-3">
          <div>
            <label class="text-sm text-slate-600 mb-1 block">活动名称</label>
            <input v-model="form.name" placeholder="请输入活动名称" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">开始时间</label>
            <input v-model="form.start_at" type="datetime-local" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">结束时间</label>
            <input v-model="form.end_at" type="datetime-local" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">排序</label>
            <input v-model.number="form.sort" type="number" placeholder="数字越大越靠前" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="flex items-center gap-2 text-sm text-slate-600">
              <input v-model="form.status" type="checkbox" :true-value="1" :false-value="0" class="rounded" />
              启用活动
            </label>
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
import { useRouter } from 'vue-router'
import request from '@/api/request'
import { confirmAction } from '@/utils/dialog'

const router = useRouter()
const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const filters = ref({ status: '-1' })
const visible = ref(false)
const form = ref<any>({
  id: 0,
  name: '',
  start_at: '',
  end_at: '',
  sort: 0,
  status: 1,
})

async function load() {
  const params: any = { page: page.value, size: size.value }
  if (filters.value.status !== '-1') params.status = filters.value.status

  const data: any = await request.get('/seckill/activities', params)
  list.value = data.list || []
  total.value = data.total || 0
}

function openCreate() {
  form.value = {
    id: 0,
    name: '',
    start_at: '',
    end_at: '',
    sort: 0,
    status: 1,
  }
  visible.value = true
}

function openEdit(row: any) {
  form.value = { ...row }
  // 转换时间格式
  if (form.value.start_at) {
    form.value.start_at = new Date(form.value.start_at).toISOString().slice(0, 16)
  }
  if (form.value.end_at) {
    form.value.end_at = new Date(form.value.end_at).toISOString().slice(0, 16)
  }
  visible.value = true
}

async function save() {
  const payload = { ...form.value }
  // 转换时间格式
  if (payload.start_at) {
    payload.start_at = new Date(payload.start_at).toISOString()
  }
  if (payload.end_at) {
    payload.end_at = new Date(payload.end_at).toISOString()
  }

  if (form.value.id) {
    await request.put(`/seckill/activities/${form.value.id}`, payload)
  } else {
    await request.post('/seckill/activities', payload)
  }
  visible.value = false
  load()
}

async function remove(id: number) {
  if (!confirmAction('确定要删除该活动吗？删除后商品数据也会被清除。')) return
  await request.delete(`/seckill/activities/${id}`)
  load()
}

function goToProducts(activityId: number) {
  router.push(`/seckill/products?activity_id=${activityId}`)
}

function formatDate(date: string) {
  return date ? new Date(date).toLocaleString('zh-CN') : '-'
}

function changePage(newPage: number) {
  if (newPage < 1 || newPage > Math.ceil(total.value / size.value)) return
  page.value = newPage
  load()
}

onMounted(load)
</script>
