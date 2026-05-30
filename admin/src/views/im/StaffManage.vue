<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-xl font-bold text-slate-800">客服坐席管理</h1>
      <button @click="openCreateModal" class="btn-primary">+ 添加客服</button>
    </div>

    <div class="card">
      <table class="w-full">
        <thead>
          <tr class="border-b border-slate-100">
            <th class="text-left py-3 px-4 text-sm font-semibold text-slate-600">ID</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-slate-600">管理员ID</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-slate-600">在线状态</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-slate-600">当前负载</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-slate-600">最大负载</th>
            <th class="text-right py-3 px-4 text-sm font-semibold text-slate-600">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="staff in staffList" :key="staff.id" class="border-b border-slate-50 hover:bg-slate-50">
            <td class="py-3 px-4 text-sm text-slate-700">{{ staff.id }}</td>
            <td class="py-3 px-4 text-sm text-slate-700">{{ staff.admin_id }}</td>
            <td class="py-3 px-4">
              <span :class="staff.is_online ? 'bg-green-100 text-green-700' : 'bg-slate-100 text-slate-500'"
                class="text-xs px-2 py-1 rounded-full">
                {{ staff.is_online ? '在线' : '离线' }}
              </span>
            </td>
            <td class="py-3 px-4 text-sm text-slate-700">{{ staff.current_load }}</td>
            <td class="py-3 px-4 text-sm text-slate-700">{{ staff.max_load }}</td>
            <td class="py-3 px-4 text-right">
              <button @click="openEditModal(staff)" class="text-sm text-blue-600 hover:text-blue-700 mr-3">编辑</button>
              <button @click="deleteStaff(staff.id)" class="text-sm text-red-600 hover:text-red-700">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="!staffList.length" class="text-center py-12 text-slate-400 text-sm">暂无客服人员</div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showModal = false">
      <div class="bg-white rounded-2xl p-6 w-96 shadow-xl">
        <h3 class="text-lg font-semibold mb-4">{{ editingStaff ? '编辑客服' : '添加客服' }}</h3>
        <div class="space-y-4">
          <div v-if="!editingStaff">
            <label class="block text-sm font-medium text-slate-700 mb-2">管理员ID *</label>
            <input v-model.number="form.admin_id" type="number" placeholder="输入管理员ID"
              class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">最大负载 *</label>
            <input v-model.number="form.max_load" type="number" placeholder="默认5"
              class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20" />
          </div>
        </div>
        <div class="flex gap-3 mt-6">
          <button @click="showModal = false"
            class="flex-1 px-4 py-2 border border-slate-200 text-slate-600 rounded-lg text-sm font-medium hover:bg-slate-50 transition">
            取消
          </button>
          <button @click="submitForm"
            class="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg text-sm font-medium hover:bg-blue-600 transition">
            确认
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'

const staffList = ref<any[]>([])
const showModal = ref(false)
const editingStaff = ref<any>(null)
const form = ref({ admin_id: null as number | null, max_load: 5 })

async function loadStaff() {
  const data: any = await request.get('/im/staff')
  staffList.value = data || []
}

function openCreateModal() {
  editingStaff.value = null
  form.value = { admin_id: null, max_load: 5 }
  showModal.value = true
}

function openEditModal(staff: any) {
  editingStaff.value = staff
  form.value = { admin_id: staff.admin_id, max_load: staff.max_load }
  showModal.value = true
}

async function submitForm() {
  if (!editingStaff.value && !form.value.admin_id) {
    alert('请输入管理员ID')
    return
  }
  if (!form.value.max_load || form.value.max_load < 1) {
    alert('最大负载必须大于0')
    return
  }

  try {
    if (editingStaff.value) {
      await request.put(`/im/staff/${editingStaff.value.id}`, { max_load: form.value.max_load })
    } else {
      await request.post('/im/staff', form.value)
    }
    showModal.value = false
    await loadStaff()
  } catch (err: any) {
    alert(err.message || '操作失败')
  }
}

async function deleteStaff(id: number) {
  if (!confirm('确认删除该客服？')) return
  try {
    await request.delete(`/im/staff/${id}`)
    await loadStaff()
  } catch (err: any) {
    alert(err.message || '删除失败')
  }
}

onMounted(loadStaff)
</script>
