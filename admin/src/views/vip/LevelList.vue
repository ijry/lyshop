<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('vip.level.title') }}</h2>
      <button @click="openCreate" class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition">{{ $t('vip.level.add') }}</button>
    </div>
    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('vip.level.name') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('vip.level.growthThreshold') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('vip.level.discount') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="item in list" :key="item.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 font-medium text-slate-800">{{ item.name }}</td>
            <td class="px-4 py-3 text-slate-600">{{ item.growth_min }}</td>
            <td class="px-4 py-3 text-slate-600">{{ item.discount_rate }}</td>
            <td class="px-4 py-3 text-slate-600">
              <button @click="openEdit(item)" class="text-blue-600 hover:text-blue-700 mr-3">{{ $t('common.edit') }}</button>
              <button @click="remove(item.id)" class="text-red-500 hover:text-red-700">{{ $t('common.delete') }}</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-if="visible" class="fixed inset-0 bg-black/30 flex items-center justify-center z-50" @click.self="visible=false">
      <div class="bg-white rounded-2xl p-6 w-96 shadow-xl">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">{{ form.id ? $t('vip.level.editTitle') : $t('vip.level.addTitle') }}</h3>
        <div class="space-y-3">
          <input v-model="form.name" :placeholder="$t('vip.level.name')" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          <input v-model.number="form.growth_min" type="number" :placeholder="$t('vip.level.growthThreshold')" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          <input v-model.number="form.discount_rate" type="number" step="0.01" :placeholder="$t('vip.level.discountHint')" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="visible=false" class="flex-1 border border-slate-200 rounded-xl py-2.5 text-sm text-slate-600">{{ $t('common.cancel') }}</button>
          <button @click="save" class="flex-1 bg-blue-700 text-white rounded-xl py-2.5 text-sm">{{ $t('common.save') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'
import { confirmAction } from '@/utils/dialog'

const { t } = useI18n()

const list = ref<any[]>([])
const visible = ref(false)
const form = ref<any>({ id: 0, name: '', growth_min: 0, discount_rate: 1 })

async function load() {
  const data: any = await request.get('/vip/levels')
  list.value = data.list || []
}
function openCreate() {
  form.value = { id: 0, name: '', growth_min: 0, discount_rate: 1 }
  visible.value = true
}
function openEdit(row: any) {
  form.value = { ...row }
  visible.value = true
}
async function save() {
  if (form.value.id) await request.put(`/vip/levels/${form.value.id}`, form.value)
  else await request.post('/vip/levels', form.value)
  visible.value = false
  load()
}
async function remove(id: number) {
  if (!confirmAction(t('vip.level.confirmDelete'))) return
  await request.delete(`/vip/levels/${id}`)
  load()
}

onMounted(load)
</script>
