<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">角色管理</h2>
      <button @click="openCreate"
        class="bg-red-600 text-white px-4 py-2 rounded-xl text-sm hover:bg-red-700 transition">
        + 新增角色
      </button>
    </div>

    <!-- Role cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <div v-for="r in roles" :key="r.id"
        class="bg-white rounded-xl shadow-sm p-5 border border-slate-100 hover:shadow-md transition">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-base font-semibold text-slate-800">{{ r.name }}</h3>
          <span class="text-xs text-slate-400">#{{ r.id }}</span>
        </div>
        <div class="flex flex-wrap gap-1.5 mb-4">
          <span v-for="p in parsePerms(r.permissions)" :key="p"
            class="bg-slate-100 text-slate-600 px-2 py-0.5 rounded text-xs">{{ p }}</span>
          <span v-if="parsePerms(r.permissions).length === 0"
            class="text-slate-400 text-xs">无权限</span>
        </div>
        <div class="flex gap-2">
          <button @click="openEdit(r)"
            class="text-xs text-blue-600 hover:text-blue-800">编辑</button>
          <button v-if="r.id !== 1" @click="deleteRole(r.id)"
            class="text-xs text-red-500 hover:text-red-700">删除</button>
        </div>
      </div>
    </div>

    <!-- Edit/Create dialog -->
    <div v-if="showDialog" class="fixed inset-0 bg-black/30 flex items-center justify-center z-50" @click.self="showDialog=false">
      <div class="bg-white rounded-2xl p-6 w-[520px] max-h-[80vh] overflow-y-auto shadow-xl">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">{{ editId ? '编辑角色' : '新增角色' }}</h3>
        <div class="mb-4">
          <label class="text-sm text-slate-600 mb-1 block">角色名称</label>
          <input v-model="form.name" placeholder="如：客服、运营"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
        </div>
        <div class="mb-4">
          <label class="text-sm text-slate-600 mb-2 block">权限配置</label>
          <div class="border border-slate-200 rounded-xl p-4 max-h-64 overflow-y-auto">
            <div v-for="group in permGroups" :key="group.prefix" class="mb-3 last:mb-0">
              <p class="text-xs font-semibold text-slate-500 mb-1.5 uppercase">{{ group.prefix }}</p>
              <div class="flex flex-wrap gap-2">
                <label v-for="p in group.perms" :key="p"
                  class="flex items-center gap-1.5 text-sm cursor-pointer select-none"
                  :class="form.perms.includes(p) ? 'text-red-600' : 'text-slate-600'">
                  <input type="checkbox" :value="p" v-model="form.perms"
                    class="accent-red-600 w-3.5 h-3.5" />
                  {{ p }}
                </label>
              </div>
            </div>
          </div>
        </div>
        <div class="flex gap-3">
          <button @click="showDialog=false" class="flex-1 border border-slate-200 rounded-xl py-2.5 text-sm text-slate-600 hover:bg-slate-50">取消</button>
          <button @click="save" class="flex-1 bg-red-600 text-white rounded-xl py-2.5 text-sm hover:bg-red-700">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import request from '@/api/request'

const roles = ref<any[]>([])
const allPerms = ref<string[]>([])
const showDialog = ref(false)
const editId = ref<number | null>(null)
const form = ref({ name: '', perms: [] as string[] })

const permGroups = computed(() => {
  const groups: Record<string, string[]> = {}
  for (const p of allPerms.value) {
    const prefix = p.split(':')[0]
    if (!groups[prefix]) groups[prefix] = []
    groups[prefix].push(p)
  }
  return Object.entries(groups).map(([prefix, perms]) => ({ prefix, perms }))
})

function parsePerms(raw: any): string[] {
  if (Array.isArray(raw)) return raw
  try { return JSON.parse(raw) } catch { return [] }
}

function openCreate() {
  editId.value = null
  form.value = { name: '', perms: [] }
  showDialog.value = true
}

function openEdit(r: any) {
  editId.value = r.id
  form.value = { name: r.name, perms: [...parsePerms(r.permissions)] }
  showDialog.value = true
}

async function save() {
  if (editId.value) {
    await request.put(`/roles/${editId.value}`, { name: form.value.name, permissions: form.value.perms })
  } else {
    await request.post('/roles', { name: form.value.name, permissions: form.value.perms })
  }
  showDialog.value = false
  loadData()
}

async function deleteRole(id: number) {
  if (!confirm('确认删除？')) return
  await request.delete(`/roles/${id}`)
  loadData()
}

async function loadData() {
  roles.value = (await request.get('/roles')) || []
  allPerms.value = (await request.get('/permissions')) || []
}

onMounted(loadData)
</script>
