<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">兑换记录</h2>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 mb-4 p-4">
      <div class="flex gap-3">
        <input v-model="filters.user_id" type="number" placeholder="用户ID" class="border border-slate-200 rounded-xl px-3 py-2 text-sm w-32" />
        <select v-model="filters.status" @change="load" class="border border-slate-200 rounded-xl px-3 py-2 text-sm">
          <option value="">全部状态</option>
          <option value="pending_ship">待发货</option>
          <option value="shipped">已发货</option>
          <option value="completed">已完成</option>
          <option value="canceled">已取消</option>
        </select>
        <button @click="load" class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600">查询</button>
      </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">兑换单号</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">用户ID</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">商品信息</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">积分</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">状态</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">时间</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="item in list" :key="item.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 font-mono text-xs text-slate-600">{{ item.exchange_no }}</td>
            <td class="px-4 py-3 text-slate-600">{{ item.user_id }}</td>
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <img v-if="item.product_cover" :src="item.product_cover" class="w-10 h-10 rounded object-cover" />
                <div>
                  <div class="font-medium text-slate-800">{{ item.product_title }}</div>
                  <div class="text-xs text-slate-500">数量: {{ item.qty }}</div>
                </div>
              </div>
            </td>
            <td class="px-4 py-3 text-slate-600 font-medium">{{ item.points_cost }}</td>
            <td class="px-4 py-3">
              <span v-if="item.status === 'pending_ship'" class="px-2 py-1 bg-amber-50 text-amber-600 rounded text-xs">待发货</span>
              <span v-else-if="item.status === 'shipped'" class="px-2 py-1 bg-blue-50 text-blue-600 rounded text-xs">已发货</span>
              <span v-else-if="item.status === 'completed'" class="px-2 py-1 bg-green-50 text-green-600 rounded text-xs">已完成</span>
              <span v-else class="px-2 py-1 bg-slate-100 text-slate-600 rounded text-xs">已取消</span>
            </td>
            <td class="px-4 py-3 text-slate-600 text-xs">{{ formatDate(item.created_at) }}</td>
            <td class="px-4 py-3">
              <button v-if="item.status === 'pending_ship'" @click="openShip(item)" class="text-blue-600 hover:text-blue-700 mr-3">发货</button>
              <button v-if="item.status === 'shipped'" @click="complete(item.id)" class="text-green-600 hover:text-green-700 mr-3">完成</button>
              <button v-if="item.status !== 'completed' && item.status !== 'canceled'" @click="openCancel(item)" class="text-red-500 hover:text-red-700">取消</button>
              <button @click="viewDetail(item)" class="text-slate-600 hover:text-slate-700">详情</button>
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

    <!-- 发货对话框 -->
    <div v-if="shipVisible" class="fixed inset-0 bg-black/30 flex items-center justify-center z-50" @click.self="shipVisible=false">
      <div class="bg-white rounded-2xl p-6 w-96 shadow-xl">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">发货</h3>
        <div class="space-y-3">
          <div>
            <label class="text-sm text-slate-600 mb-1 block">物流单号</label>
            <input v-model="shipForm.tracking_no" placeholder="请输入物流单号" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="shipVisible=false" class="flex-1 border border-slate-200 rounded-xl py-2.5 text-sm text-slate-600">取消</button>
          <button @click="ship" class="flex-1 bg-blue-700 text-white rounded-xl py-2.5 text-sm">确认发货</button>
        </div>
      </div>
    </div>

    <!-- 取消对话框 -->
    <div v-if="cancelVisible" class="fixed inset-0 bg-black/30 flex items-center justify-center z-50" @click.self="cancelVisible=false">
      <div class="bg-white rounded-2xl p-6 w-96 shadow-xl">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">取消兑换</h3>
        <div class="space-y-3">
          <div>
            <label class="text-sm text-slate-600 mb-1 block">取消原因</label>
            <textarea v-model="cancelForm.reason" placeholder="请输入取消原因" rows="3" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm"></textarea>
          </div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="cancelVisible=false" class="flex-1 border border-slate-200 rounded-xl py-2.5 text-sm text-slate-600">取消</button>
          <button @click="cancel" class="flex-1 bg-red-600 text-white rounded-xl py-2.5 text-sm">确认取消</button>
        </div>
      </div>
    </div>

    <!-- 详情对话框 -->
    <div v-if="detailVisible" class="fixed inset-0 bg-black/30 flex items-center justify-center z-50" @click.self="detailVisible=false">
      <div class="bg-white rounded-2xl p-6 w-[500px] max-h-[90vh] overflow-y-auto shadow-xl">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">兑换详情</h3>
        <div v-if="detail" class="space-y-3 text-sm">
          <div class="flex justify-between"><span class="text-slate-600">兑换单号:</span><span class="font-mono">{{ detail.exchange_no }}</span></div>
          <div class="flex justify-between"><span class="text-slate-600">用户ID:</span><span>{{ detail.user_id }}</span></div>
          <div class="flex justify-between"><span class="text-slate-600">商品:</span><span>{{ detail.product_title }}</span></div>
          <div class="flex justify-between"><span class="text-slate-600">数量:</span><span>{{ detail.qty }}</span></div>
          <div class="flex justify-between"><span class="text-slate-600">积分:</span><span class="font-medium">{{ detail.points_cost }}</span></div>
          <div class="flex justify-between"><span class="text-slate-600">状态:</span><span>{{ getStatusText(detail.status) }}</span></div>
          <div v-if="detail.tracking_no" class="flex justify-between"><span class="text-slate-600">物流单号:</span><span>{{ detail.tracking_no }}</span></div>
          <div v-if="detail.address_snapshot" class="border-t pt-3">
            <div class="text-slate-600 mb-2">收货地址:</div>
            <div class="bg-slate-50 p-3 rounded-lg text-xs">{{ formatAddress(detail.address_snapshot) }}</div>
          </div>
          <div v-if="detail.virtual_content" class="border-t pt-3">
            <div class="text-slate-600 mb-2">虚拟商品内容:</div>
            <div class="bg-slate-50 p-3 rounded-lg text-xs whitespace-pre-wrap">{{ detail.virtual_content }}</div>
          </div>
          <div class="flex justify-between text-xs text-slate-500"><span>创建时间:</span><span>{{ detail.created_at }}</span></div>
        </div>
        <div class="mt-5">
          <button @click="detailVisible=false" class="w-full border border-slate-200 rounded-xl py-2.5 text-sm text-slate-600">关闭</button>
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
const filters = ref({ user_id: '', status: '' })

const shipVisible = ref(false)
const shipForm = ref({ id: 0, tracking_no: '' })

const cancelVisible = ref(false)
const cancelForm = ref({ id: 0, reason: '' })

const detailVisible = ref(false)
const detail = ref<any>(null)

async function load() {
  const params: any = { page: page.value, size: size.value }
  if (filters.value.user_id) params.user_id = filters.value.user_id
  if (filters.value.status) params.status = filters.value.status

  const data: any = await request.get('/points/exchanges', params)
  list.value = data.list || []
  total.value = data.total || 0
}

function openShip(item: any) {
  shipForm.value = { id: item.id, tracking_no: '' }
  shipVisible.value = true
}

async function ship() {
  await request.put(`/points/exchanges/${shipForm.value.id}/ship`, { tracking_no: shipForm.value.tracking_no })
  shipVisible.value = false
  load()
}

async function complete(id: number) {
  if (!confirmAction('确认完成该兑换订单吗？')) return
  await request.put(`/points/exchanges/${id}/complete`)
  load()
}

function openCancel(item: any) {
  cancelForm.value = { id: item.id, reason: '' }
  cancelVisible.value = true
}

async function cancel() {
  await request.put(`/points/exchanges/${cancelForm.value.id}/cancel`, { reason: cancelForm.value.reason })
  cancelVisible.value = false
  load()
}

async function viewDetail(item: any) {
  const data: any = await request.get(`/points/exchanges/${item.id}`)
  detail.value = data
  detailVisible.value = true
}

function formatDate(date: string) {
  return date ? new Date(date).toLocaleString('zh-CN') : '-'
}

function getStatusText(status: string) {
  const map: any = {
    pending_ship: '待发货',
    shipped: '已发货',
    completed: '已完成',
    canceled: '已取消',
  }
  return map[status] || status
}

function formatAddress(snapshot: any) {
  if (typeof snapshot === 'string') {
    try {
      snapshot = JSON.parse(snapshot)
    } catch {
      return snapshot
    }
  }
  return `${snapshot.name} ${snapshot.phone}\n${snapshot.province} ${snapshot.city} ${snapshot.district} ${snapshot.address}`
}

function changePage(newPage: number) {
  if (newPage < 1 || newPage > Math.ceil(total.value / size.value)) return
  page.value = newPage
  load()
}

onMounted(load)
</script>
