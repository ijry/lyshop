<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('wms.warehouseTitle') }}</h2>
      <button
        class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition"
        @click="openCreate"
      >
        {{ $t('wms.createWarehouse') }}
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm p-4 mb-4 flex gap-3 border border-slate-100">
      <input
        v-model.trim="keyword"
        :placeholder="$t('wms.searchWarehousePlaceholder')"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm flex-1 focus:outline-none focus:border-blue-400"
      />
      <button class="px-4 py-2 bg-slate-100 rounded-lg text-sm hover:bg-slate-200" @click="loadWarehouses">
        {{ $t('common.search') }}
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.id') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.name') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.warehouseCode') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.warehouseAddress') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('wms.contact') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.status') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="item in warehouses" :key="item.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-slate-500">{{ item.id }}</td>
            <td class="px-4 py-3 font-medium text-slate-800">{{ item.name }}</td>
            <td class="px-4 py-3 text-slate-500">{{ item.code || '-' }}</td>
            <td class="px-4 py-3 text-slate-500">{{ item.address || '-' }}</td>
            <td class="px-4 py-3 text-slate-500">{{ item.contact || '-' }} {{ item.phone || '' }}</td>
            <td class="px-4 py-3">
              <span :class="item.status === 1 ? 'bg-green-50 text-green-600' : 'bg-slate-100 text-slate-500'" class="px-2 py-1 rounded-full text-xs">
                {{ item.status === 1 ? $t('common.enabled') : $t('common.disabled') }}
              </span>
            </td>
            <td class="px-4 py-3">
              <button class="text-blue-600 hover:underline text-xs" @click="openEdit(item)">{{ $t('common.edit') }}</button>
            </td>
          </tr>
          <tr v-if="!warehouses.length">
            <td colspan="7" class="px-4 py-12 text-center text-slate-400">{{ $t('common.noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showForm" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50" @click.self="closeForm">
      <div class="bg-white rounded-2xl p-6 w-[540px] shadow-xl">
        <h3 class="font-semibold text-slate-800 mb-4">
          {{ form.id ? $t('wms.editWarehouse') : $t('wms.createWarehouse') }}
        </h3>
        <div class="grid grid-cols-2 gap-3">
          <input
            v-model.trim="form.name"
            :placeholder="$t('wms.warehouseNamePlaceholder')"
            class="border border-slate-200 rounded-xl px-4 py-2.5 text-sm"
          />
          <input
            v-model.trim="form.code"
            :placeholder="$t('wms.warehouseCodePlaceholder')"
            class="border border-slate-200 rounded-xl px-4 py-2.5 text-sm"
          />
          <input
            v-model.trim="form.contact"
            :placeholder="$t('wms.contactPlaceholder')"
            class="border border-slate-200 rounded-xl px-4 py-2.5 text-sm"
          />
          <input
            v-model.trim="form.phone"
            :placeholder="$t('wms.phonePlaceholder')"
            class="border border-slate-200 rounded-xl px-4 py-2.5 text-sm"
          />
          <input
            v-model.trim="form.address"
            :placeholder="$t('wms.warehouseAddressPlaceholder')"
            class="border border-slate-200 rounded-xl px-4 py-2.5 text-sm col-span-2"
          />
          <select v-model.number="form.status" class="border border-slate-200 rounded-xl px-4 py-2.5 text-sm col-span-2">
            <option :value="1">{{ $t('common.enabled') }}</option>
            <option :value="0">{{ $t('common.disabled') }}</option>
          </select>
        </div>
        <div class="flex gap-3 mt-5">
          <button class="flex-1 px-4 py-2 bg-slate-100 rounded-lg text-sm" @click="closeForm">{{ $t('common.cancel') }}</button>
          <button class="flex-1 px-4 py-2 bg-blue-700 text-white rounded-lg text-sm" @click="submitForm">{{ $t('common.save') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { createWarehouse, listWarehouses, updateWarehouse, type WmsWarehouse } from '@/api/wms'
import { notify } from '@/utils/notify'

const { t } = useI18n()
const warehouses = ref<WmsWarehouse[]>([])
const keyword = ref('')
const showForm = ref(false)
const form = ref({
  id: 0,
  name: '',
  code: '',
  address: '',
  contact: '',
  phone: '',
  status: 1,
})

function openCreate() {
  form.value = { id: 0, name: '', code: '', address: '', contact: '', phone: '', status: 1 }
  showForm.value = true
}

function openEdit(row: WmsWarehouse) {
  form.value = {
    id: row.id,
    name: row.name || '',
    code: row.code || '',
    address: row.address || '',
    contact: row.contact || '',
    phone: row.phone || '',
    status: Number(row.status || 0) === 1 ? 1 : 0,
  }
  showForm.value = true
}

function closeForm() {
  showForm.value = false
}

async function loadWarehouses() {
  const rows = await listWarehouses({ keyword: keyword.value || undefined })
  warehouses.value = Array.isArray(rows) ? rows : []
}

async function submitForm() {
  if (!form.value.name.trim()) {
    notify(t('wms.warehouseNameRequired'))
    return
  }
  const payload = {
    name: form.value.name.trim(),
    code: form.value.code.trim(),
    address: form.value.address.trim(),
    contact: form.value.contact.trim(),
    phone: form.value.phone.trim(),
    status: Number(form.value.status || 0) === 1 ? 1 : 0,
  }
  if (form.value.id > 0) {
    await updateWarehouse(form.value.id, payload)
  } else {
    await createWarehouse(payload)
  }
  closeForm()
  await loadWarehouses()
}

onMounted(loadWarehouses)
</script>
