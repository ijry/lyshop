<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">管理员管理</h2>
      <button @click="showCreate = true"
        class="bg-red-600 text-white px-4 py-2 rounded-xl text-sm hover:bg-red-700 transition">
        + 新增管理员
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 text-slate-500">
          <tr>
            <th class="px-4 py-3 text-left">ID</th>
            <th class="px-4 py-3 text-left">用户名</th>
            <th class="px-4 py-3 text-left">角色</th>
            <th class="px-4 py-3 text-left">状态</th>
            <th class="px-4 py-3 text-left">创建时间</th>
            <th class="px-4 py-3 text-left">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="a in admins" :key="a.id" class="border-t border-slate-50 hover:bg-slate-50/50">
            <td class="px-4 py-3 text-slate-500">{{ a.id }}</td>
            <td class="px-4 py-3 font-medium text-slate-800">{{ a.username }}</td>
            <td class="px-4 py-3">
              <span class="bg-blue-50 text-blue-600 px-2 py-0.5 rounded text-xs">
                {{ roleName(a.role_id) }}
              </span>
            </td>
            <td class="px-4 py-3">
              <span :class="a.status === 1 ? 'bg-green-50 text-green-600' : 'bg-red-50 text-red-500'"
                class="px-2 py-0.5 rounded text-xs">{{ a.status === 1 ? '正常' : '禁用' }}</span>
            </td>
            <td class="px-4 py-3 text-slate-400">{{ a.created_at?.slice(0, 10) }}</td>
            <td class="px-4 py-3">
              <button v-if="a.id !== 1" @click="deleteAdmin(a.id)"
                class="text-red-500 hover:text-red-700 text-xs">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create dialog -->
    <div v-if="showCreate" class="fixed inset-0 bg-black/30 flex items-center justify-center z-50" @click.self="showCreate=false">
      <div class="bg-white rounded-2xl p-6 w-96 shadow-xl">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">新增管理员</h3>
        <div class="space-y-3">
          <input v-model="form.username" placeholder="用户名"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
          <input v-model="form.password" type="password" placeholder="密码"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
          <select v-model="form.role_id"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500">
            <option value="">选择角色</option>
            <option v-for="r in roles" :key="r.id" :value="r.id">{{ r.name }}</option>
          </select>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="showCreate=false" class="flex-1 border border-slate-200 rounded-xl py-2.5 text-sm text-slate-600 hover:bg-slate-50">取消</button>
          <button @click="createAdmin" class="flex-1 bg-red-600 text-white rounded-xl py-2.5 text-sm hover:bg-red-700">创建</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'
import { confirmAction } from '@/utils/dialog'

const admins = ref<any[]>([])
const roles = ref<any[]>([])
const showCreate = ref(false)
const form = ref({ username: '', password: '', role_id: '' })

function roleName(roleId: number) {
  return roles.value.find(r => r.id === roleId)?.name || `角色#${roleId}`
}

async function loadData() {
  admins.value = (await request.get('/admins')) || []
  roles.value = (await request.get('/roles')) || []
}

async function createAdmin() {
  await request.post('/admins', { ...form.value, role_id: Number(form.value.role_id) })
  showCreate.value = false
  form.value = { username: '', password: '', role_id: '' }
  loadData()
}

async function deleteAdmin(id: number) {
  if (!confirmAction('确认删除？')) return
  await request.delete(`/admins/${id}`)
  loadData()
}

onMounted(loadData)
</script>
