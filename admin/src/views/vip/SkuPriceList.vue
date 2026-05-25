<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('vip.skuPrice.title') }}</h2>
      <button @click="openCreate" class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition">{{ $t('vip.skuPrice.add') }}</button>
    </div>
    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('vip.skuPrice.skuId') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('vip.skuPrice.level') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('vip.skuPrice.price') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="item in list" :key="item.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 font-medium text-slate-800">{{ item.sku_id }}</td>
            <td class="px-4 py-3 text-slate-600">{{ item.level_name }}</td>
            <td class="px-4 py-3 text-slate-600">{{ item.vip_price }}</td>
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
        <h3 class="text-lg font-semibold text-slate-800 mb-4">{{ form.id ? $t('vip.skuPrice.editTitle') : $t('vip.skuPrice.addTitle') }}</h3>
        <div class="space-y-3">
          <input v-model.number="form.sku_id" type="number" :placeholder="$t('vip.skuPrice.skuId')" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          <input v-model="form.level_name" :placeholder="$t('vip.skuPrice.levelName')" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          <input v-model.number="form.vip_price" type="number" :placeholder="$t('vip.skuPrice.price')" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
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
const form = ref<any>({ id: 0, sku_id: 0, level_name: '', vip_price: 0 })

async function load() {
  const data: any = await request.get('/vip/sku-prices')
  list.value = data.list || []
}
function openCreate() {
  form.value = { id: 0, sku_id: 0, level_name: '', vip_price: 0 }
  visible.value = true
}
function openEdit(row: any) {
  form.value = { ...row }
  visible.value = true
}
async function save() {
  if (form.value.id) await request.put(`/vip/sku-prices/${form.value.id}`, form.value)
  else await request.post('/vip/sku-prices', form.value)
  visible.value = false
  load()
}
async function remove(id: number) {
  if (!confirmAction(t('vip.skuPrice.confirmDelete'))) return
  await request.delete(`/vip/sku-prices/${id}`)
  load()
}

onMounted(load)
</script>
