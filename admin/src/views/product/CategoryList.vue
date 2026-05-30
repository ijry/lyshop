<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('category.title') }}</h2>
      <button
        class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition"
        @click="openCreate"
      >
        {{ $t('category.add') }}
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.id') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('category.name') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('category.sort') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.status') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('category.updateTime') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="row in categories" :key="row.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-slate-400">{{ row.id }}</td>
            <td class="px-4 py-3 font-medium text-slate-700">{{ row.name }}</td>
            <td class="px-4 py-3 text-slate-600">{{ row.sort }}</td>
            <td class="px-4 py-3">
              <span
                :class="row.status === 1 ? 'bg-green-50 text-green-600' : 'bg-slate-100 text-slate-400'"
                class="px-2 py-1 rounded-full text-xs"
              >
                {{ row.status === 1 ? $t('common.enabled') : $t('common.disabled') }}
              </span>
            </td>
            <td class="px-4 py-3 text-slate-500">{{ formatDate(row.updated_at) }}</td>
            <td class="px-4 py-3">
              <button class="text-blue-600 hover:underline text-xs mr-3" @click="openEdit(row)">{{ $t('common.edit') }}</button>
              <button class="text-emerald-600 hover:underline text-xs mr-3" @click="toggleStatus(row)">
                {{ row.status === 1 ? $t('common.disabled') : $t('common.enabled') }}
              </button>
              <button class="text-red-500 hover:underline text-xs" @click="remove(row.id)">{{ $t('common.delete') }}</button>
            </td>
          </tr>
          <tr v-if="!categories.length">
            <td colspan="6" class="px-4 py-12 text-center text-slate-400">{{ $t('category.noCategory') }}</td>
          </tr>
        </tbody>
      </table>
      <div class="px-4 py-3 flex items-center justify-between border-t border-slate-100 text-sm text-slate-500">
        <span>{{ $t('common.totalCount', { total }) }}</span>
        <div class="flex gap-2">
          <button :disabled="query.page <= 1" @click="prevPage"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40">{{ $t('common.prevPage') }}</button>
          <button :disabled="query.page * query.size >= total" @click="nextPage"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40">{{ $t('common.nextPage') }}</button>
        </div>
      </div>
    </div>

    <div v-if="showDialog" class="fixed inset-0 z-50 bg-black/30 flex items-center justify-center p-4">
      <div class="w-full max-w-lg rounded-2xl bg-white shadow-xl p-6">
        <h3 class="text-lg font-semibold text-slate-900 mb-4">{{ editingID ? $t('category.editTitle') : $t('category.addTitle') }}</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm text-slate-600 mb-1">{{ $t('category.name') }}</label>
            <input
              v-model.trim="form.name"
              class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-blue-500"
              :placeholder="$t('category.namePlaceholder')"
            />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-sm text-slate-600 mb-1">{{ $t('category.sort') }}</label>
              <input
                v-model.number="form.sort"
                type="number"
                min="0"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-blue-500"
              />
            </div>
            <div>
              <label class="block text-sm text-slate-600 mb-1">{{ $t('common.status') }}</label>
              <select
                v-model.number="form.status"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-blue-500"
              >
                <option :value="1">{{ $t('common.enabled') }}</option>
                <option :value="0">{{ $t('common.disabled') }}</option>
              </select>
            </div>
          </div>
          <div>
            <label class="block text-sm text-slate-600 mb-1">{{ $t('category.icon') }}</label>
            <input
              v-model.trim="form.icon"
              class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-blue-500"
              :placeholder="$t('category.iconUrl')"
            />
          </div>
        </div>
        <div class="flex gap-3 mt-6">
          <button
            class="flex-1 border border-slate-200 rounded-xl py-2.5 text-sm text-slate-600 hover:bg-slate-50"
            @click="closeDialog"
          >
            {{ $t('common.cancel') }}
          </button>
          <button
            class="flex-1 bg-blue-700 text-white rounded-xl py-2.5 text-sm hover:bg-blue-600 disabled:opacity-50"
            :disabled="saving"
            @click="submit"
          >
            {{ saving ? $t('common.submitting') : $t('common.save') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { createCategory, deleteCategory, getCategories, updateCategory } from '@/api/plugins'
import { notify } from '@/utils/notify'
import { confirmAction } from '@/utils/dialog'

const { t } = useI18n()

type CategoryRow = {
  id: number
  parent_id: number
  name: string
  icon: string
  sort: number
  status: number
  updated_at?: string
}

const categories = ref<CategoryRow[]>([])
const total = ref(0)
const query = ref({ page: 1, size: 20 })
const showDialog = ref(false)
const saving = ref(false)
const editingID = ref(0)
const form = ref({
  name: '',
  icon: '',
  sort: 0,
  status: 1,
})

const formatDate = (v?: string) => v ? String(v).slice(0, 19).replace('T', ' ') : '-'

function resetForm() {
  form.value = {
    name: '',
    icon: '',
    sort: 0,
    status: 1,
  }
}

async function loadCategories() {
  const data: any = await getCategories({ page: query.value.page, size: query.value.size })
  categories.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
  total.value = Number(data?.total || categories.value.length)
}

function prevPage() {
  if (query.value.page <= 1) return
  query.value.page -= 1
  loadCategories()
}

function nextPage() {
  if (query.value.page * query.value.size >= total.value) return
  query.value.page += 1
  loadCategories()
}

function openCreate() {
  editingID.value = 0
  resetForm()
  showDialog.value = true
}

function openEdit(row: CategoryRow) {
  editingID.value = Number(row.id || 0)
  form.value = {
    name: String(row.name || ''),
    icon: String(row.icon || ''),
    sort: Number(row.sort || 0),
    status: Number(row.status || 0),
  }
  showDialog.value = true
}

function closeDialog() {
  showDialog.value = false
}

async function submit() {
  const name = form.value.name.trim()
  if (!name) {
    notify(t('category.nameRequired'))
    return
  }
  saving.value = true
  try {
    const payload = {
      parent_id: 0,
      name,
      icon: form.value.icon.trim(),
      sort: Number(form.value.sort || 0),
      status: Number(form.value.status || 0) === 1 ? 1 : 0,
    }
    if (editingID.value > 0) {
      await updateCategory(editingID.value, payload)
    } else {
      await createCategory(payload)
    }
    showDialog.value = false
    await loadCategories()
  } finally {
    saving.value = false
  }
}

async function toggleStatus(row: CategoryRow) {
  const nextStatus = Number(row.status || 0) === 1 ? 0 : 1
  await updateCategory(Number(row.id), { status: nextStatus })
  await loadCategories()
}

async function remove(id: number) {
  if (!confirmAction(t('category.confirmDelete'))) return
  await deleteCategory(id)
  await loadCategories()
}

onMounted(loadCategories)
</script>
