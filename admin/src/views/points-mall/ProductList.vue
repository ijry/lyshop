<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">积分商品</h2>
      <button @click="openCreate" class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition">新增商品</button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 mb-4 p-4">
      <div class="flex gap-3">
        <select v-model="filters.type" @change="load" class="border border-slate-200 rounded-xl px-3 py-2 text-sm">
          <option value="">全部类型</option>
          <option value="coupon">优惠券</option>
          <option value="physical">实物</option>
          <option value="virtual">虚拟</option>
        </select>
        <select v-model="filters.status" @change="load" class="border border-slate-200 rounded-xl px-3 py-2 text-sm">
          <option value="-1">全部状态</option>
          <option value="1">上架</option>
          <option value="0">下架</option>
        </select>
      </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">商品信息</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">类型</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">积分价格</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">库存</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">已兑换</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">状态</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="item in list" :key="item.id" class="hover:bg-slate-50">
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <img v-if="item.cover" :src="item.cover" class="w-12 h-12 rounded-lg object-cover" />
                <div class="w-12 h-12 bg-slate-100 rounded-lg flex items-center justify-center" v-else>
                  <span class="text-slate-400 text-xs">无图</span>
                </div>
                <div>
                  <div class="font-medium text-slate-800">{{ item.title }}</div>
                  <div class="text-xs text-slate-500">ID: {{ item.id }}</div>
                </div>
              </div>
            </td>
            <td class="px-4 py-3 text-slate-600">
              <span v-if="item.type === 'coupon'" class="px-2 py-1 bg-orange-50 text-orange-600 rounded text-xs">优惠券</span>
              <span v-else-if="item.type === 'physical'" class="px-2 py-1 bg-blue-50 text-blue-600 rounded text-xs">实物</span>
              <span v-else class="px-2 py-1 bg-purple-50 text-purple-600 rounded text-xs">虚拟</span>
            </td>
            <td class="px-4 py-3 text-slate-600 font-medium">{{ item.points_price }} 积分</td>
            <td class="px-4 py-3 text-slate-600">{{ item.stock === 0 ? '无限' : item.stock }}</td>
            <td class="px-4 py-3 text-slate-600">{{ item.sold_count }}</td>
            <td class="px-4 py-3">
              <span v-if="item.status === 1" class="px-2 py-1 bg-green-50 text-green-600 rounded text-xs">上架</span>
              <span v-else class="px-2 py-1 bg-slate-100 text-slate-600 rounded text-xs">下架</span>
            </td>
            <td class="px-4 py-3">
              <button @click="openEdit(item)" class="text-blue-600 hover:text-blue-700 mr-3">编辑</button>
              <button @click="toggleStatus(item)" class="text-amber-600 hover:text-amber-700 mr-3">
                {{ item.status === 1 ? '下架' : '上架' }}
              </button>
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
      <div class="bg-white rounded-2xl p-6 w-[600px] max-h-[90vh] overflow-y-auto shadow-xl">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">{{ form.id ? '编辑商品' : '新增商品' }}</h3>
        <div class="space-y-3">
          <div>
            <label class="text-sm text-slate-600 mb-1 block">商品标题</label>
            <input v-model="form.title" placeholder="请输入商品标题" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">商品类型</label>
            <select v-model="form.type" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm">
              <option value="coupon">优惠券</option>
              <option value="physical">实物</option>
              <option value="virtual">虚拟</option>
            </select>
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">积分价格</label>
            <input v-model.number="form.points_price" type="number" placeholder="请输入积分价格" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">库存（0=无限）</label>
            <input v-model.number="form.stock" type="number" placeholder="0" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">封面图片</label>
            <input v-model="form.cover" placeholder="请输入图片URL" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">商品描述</label>
            <textarea v-model="form.description" placeholder="请输入商品描述" rows="3" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm"></textarea>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="text-sm text-slate-600 mb-1 block">每人限兑（0=不限）</label>
              <input v-model.number="form.limit_per_user" type="number" placeholder="0" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
            </div>
            <div>
              <label class="text-sm text-slate-600 mb-1 block">每日限兑（0=不限）</label>
              <input v-model.number="form.limit_per_day" type="number" placeholder="0" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
            </div>
          </div>
          <div v-if="form.type === 'coupon'">
            <label class="text-sm text-slate-600 mb-1 block">关联优惠券ID</label>
            <input v-model.number="form.coupon_id" type="number" placeholder="请输入优惠券ID" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div v-if="form.type === 'physical'">
            <label class="flex items-center gap-2 text-sm text-slate-600">
              <input v-model="form.need_address" type="checkbox" class="rounded" />
              需要收货地址
            </label>
          </div>
          <div v-if="form.type === 'virtual'">
            <label class="text-sm text-slate-600 mb-1 block">虚拟商品内容</label>
            <textarea v-model="form.virtual_content" placeholder="兑换后显示的内容（如兑换码）" rows="3" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm"></textarea>
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
import request from '@/api/request'
import { confirmAction } from '@/utils/dialog'

const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const filters = ref({ type: '', status: '-1' })
const visible = ref(false)
const form = ref<any>({
  id: 0,
  title: '',
  type: 'physical',
  points_price: 0,
  stock: 0,
  cover: '',
  description: '',
  limit_per_user: 0,
  limit_per_day: 0,
  coupon_id: 0,
  need_address: true,
  virtual_content: '',
  status: 1,
})

async function load() {
  const params: any = { page: page.value, size: size.value }
  if (filters.value.type) params.type = filters.value.type
  if (filters.value.status !== '-1') params.status = filters.value.status

  const data: any = await request.get('/points/products', params)
  list.value = data.list || []
  total.value = data.total || 0
}

function openCreate() {
  form.value = {
    id: 0,
    title: '',
    type: 'physical',
    points_price: 0,
    stock: 0,
    cover: '',
    description: '',
    limit_per_user: 0,
    limit_per_day: 0,
    coupon_id: 0,
    need_address: true,
    virtual_content: '',
    status: 1,
  }
  visible.value = true
}

function openEdit(row: any) {
  form.value = { ...row }
  visible.value = true
}

async function save() {
  if (form.value.id) {
    await request.put(`/points/products/${form.value.id}`, form.value)
  } else {
    await request.post('/points/products', form.value)
  }
  visible.value = false
  load()
}

async function toggleStatus(item: any) {
  const newStatus = item.status === 1 ? 0 : 1
  await request.put(`/points/products/${item.id}/status`, { status: newStatus })
  load()
}

async function remove(id: number) {
  if (!confirmAction('确定要删除该商品吗？')) return
  await request.delete(`/points/products/${id}`)
  load()
}

function changePage(newPage: number) {
  if (newPage < 1 || newPage > Math.ceil(total.value / size.value)) return
  page.value = newPage
  load()
}

onMounted(load)
</script>
