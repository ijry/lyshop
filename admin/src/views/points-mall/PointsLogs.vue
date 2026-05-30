<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">积分日志</h2>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 mb-4 p-4">
      <div class="flex gap-3">
        <input v-model="filters.user_id" type="number" placeholder="用户ID" class="border border-slate-200 rounded-xl px-3 py-2 text-sm w-32" />
        <select v-model="filters.type" @change="load" class="border border-slate-200 rounded-xl px-3 py-2 text-sm">
          <option value="0">全部类型</option>
          <option value="1">签到</option>
          <option value="2">订单抵扣</option>
          <option value="3">兑换消耗</option>
          <option value="4">订单完成</option>
          <option value="5">管理员调整</option>
          <option value="6">过期扣除</option>
          <option value="7">活动奖励</option>
        </select>
        <button @click="load" class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600">查询</button>
        <button @click="openAdjust" class="px-4 py-2 bg-green-600 text-white text-sm rounded-xl hover:bg-green-500 ml-auto">调整积分</button>
      </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">用户ID</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">类型</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">积分变动</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">备注</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">时间</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="item in list" :key="item.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-slate-600">{{ item.user_id }}</td>
            <td class="px-4 py-3">
              <span class="px-2 py-1 rounded text-xs" :class="getTypeClass(item.type)">{{ getTypeText(item.type) }}</span>
            </td>
            <td class="px-4 py-3 font-medium" :class="item.points > 0 ? 'text-green-600' : 'text-red-600'">
              {{ item.points > 0 ? '+' : '' }}{{ item.points }}
            </td>
            <td class="px-4 py-3 text-slate-600 text-xs">{{ item.remark }}</td>
            <td class="px-4 py-3 text-slate-600 text-xs">{{ formatDate(item.created_at) }}</td>
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

    <!-- 调整积分对话框 -->
    <div v-if="adjustVisible" class="fixed inset-0 bg-black/30 flex items-center justify-center z-50" @click.self="adjustVisible=false">
      <div class="bg-white rounded-2xl p-6 w-96 shadow-xl">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">调整用户积分</h3>
        <div class="space-y-3">
          <div>
            <label class="text-sm text-slate-600 mb-1 block">用户ID</label>
            <input v-model.number="adjustForm.user_id" type="number" placeholder="请输入用户ID" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">积分变动（正数增加，负数减少）</label>
            <input v-model.number="adjustForm.points" type="number" placeholder="如：100 或 -50" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          </div>
          <div>
            <label class="text-sm text-slate-600 mb-1 block">备注</label>
            <textarea v-model="adjustForm.remark" placeholder="请输入调整原因" rows="3" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm"></textarea>
          </div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="adjustVisible=false" class="flex-1 border border-slate-200 rounded-xl py-2.5 text-sm text-slate-600">取消</button>
          <button @click="adjust" class="flex-1 bg-green-600 text-white rounded-xl py-2.5 text-sm">确认调整</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import request from '@/api/request'

const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const filters = ref({ user_id: '', type: '0' })

const adjustVisible = ref(false)
const adjustForm = ref({ user_id: 0, points: 0, remark: '' })

async function load() {
  const params: any = { page: page.value, size: size.value }
  if (filters.value.user_id) params.user_id = filters.value.user_id
  if (filters.value.type !== '0') params.type = filters.value.type

  const data: any = await request.get('/points/logs', params)
  list.value = data.list || []
  total.value = data.total || 0
}

function openAdjust() {
  adjustForm.value = { user_id: 0, points: 0, remark: '' }
  adjustVisible.value = true
}

async function adjust() {
  if (!adjustForm.value.user_id || !adjustForm.value.points) {
    alert('请填写完整信息')
    return
  }
  await request.post('/points/adjust', adjustForm.value)
  adjustVisible.value = false
  load()
}

function getTypeText(type: number) {
  const map: any = {
    1: '签到',
    2: '订单抵扣',
    3: '兑换消耗',
    4: '订单完成',
    5: '管理员调整',
    6: '过期扣除',
    7: '活动奖励',
  }
  return map[type] || '未知'
}

function getTypeClass(type: number) {
  const map: any = {
    1: 'bg-blue-50 text-blue-600',
    2: 'bg-orange-50 text-orange-600',
    3: 'bg-red-50 text-red-600',
    4: 'bg-green-50 text-green-600',
    5: 'bg-purple-50 text-purple-600',
    6: 'bg-slate-100 text-slate-600',
    7: 'bg-pink-50 text-pink-600',
  }
  return map[type] || 'bg-slate-100 text-slate-600'
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
