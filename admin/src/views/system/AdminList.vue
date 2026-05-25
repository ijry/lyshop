<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('system.admin.title') }}</h2>
      <button @click="showCreate = true"
        class="bg-red-600 text-white px-4 py-2 rounded-xl text-sm hover:bg-red-700 transition">
        {{ $t('system.admin.add') }}
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 text-slate-500">
          <tr>
            <th class="px-4 py-3 text-left">{{ $t('common.id') }}</th>
            <th class="px-4 py-3 text-left">{{ $t('system.admin.username') }}</th>
            <th class="px-4 py-3 text-left">{{ $t('system.admin.role') }}</th>
            <th class="px-4 py-3 text-left">{{ $t('common.status') }}</th>
            <th class="px-4 py-3 text-left">{{ $t('system.admin.createTime') }}</th>
            <th class="px-4 py-3 text-left">{{ $t('common.action') }}</th>
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
                class="px-2 py-0.5 rounded text-xs">{{ a.status === 1 ? $t('system.admin.active') : $t('system.admin.banned') }}</span>
            </td>
            <td class="px-4 py-3 text-slate-400">{{ a.created_at?.slice(0, 10) }}</td>
            <td class="px-4 py-3">
              <button v-if="a.id !== 1" @click="deleteAdmin(a.id)"
                class="text-red-500 hover:text-red-700 text-xs">{{ $t('common.delete') }}</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create dialog -->
    <div v-if="showCreate" class="fixed inset-0 bg-black/30 flex items-center justify-center z-50" @click.self="showCreate=false">
      <div class="bg-white rounded-2xl p-6 w-96 shadow-xl">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">{{ $t('system.admin.addTitle') }}</h3>
        <div class="space-y-3">
          <input v-model="form.username" :placeholder="$t('system.admin.username')"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
          <input v-model="form.password" type="password" :placeholder="$t('system.admin.password')"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
          <select v-model="form.role_id"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500">
            <option value="">{{ $t('system.admin.selectRole') }}</option>
            <option v-for="r in roles" :key="r.id" :value="r.id">{{ r.name }}</option>
          </select>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="showCreate=false" class="flex-1 border border-slate-200 rounded-xl py-2.5 text-sm text-slate-600 hover:bg-slate-50">{{ $t('common.cancel') }}</button>
          <button @click="createAdmin" class="flex-1 bg-red-600 text-white rounded-xl py-2.5 text-sm hover:bg-red-700">{{ $t('common.create') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'
import { confirmAction } from '@/utils/dialog'

const { t } = useI18n()

const admins = ref<any[]>([])
const roles = ref<any[]>([])
const showCreate = ref(false)
const form = ref({ username: '', password: '', role_id: '' })

function roleName(roleId: number) {
  return roles.value.find(r => r.id === roleId)?.name || `${t('system.admin.role')}#${roleId}`
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
  if (!confirmAction(t('common.confirmDelete'))) return
  await request.delete(`/admins/${id}`)
  loadData()
}

onMounted(loadData)
</script>
