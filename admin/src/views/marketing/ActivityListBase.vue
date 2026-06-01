<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ pageTitle }}</h2>
      <button
        @click="openCreate"
        class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition"
      >
        {{ $t('common.create') }}
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.id') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.name') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('marketingActivity.startAt') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('marketingActivity.endAt') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.status') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="row in list" :key="row.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-slate-400">{{ row.id }}</td>
            <td class="px-4 py-3 text-slate-800">{{ row.name }}</td>
            <td class="px-4 py-3 text-slate-600">{{ formatDate(row.start_at) }}</td>
            <td class="px-4 py-3 text-slate-600">{{ formatDate(row.end_at) }}</td>
            <td class="px-4 py-3">
              <span :class="row.status === 1 ? 'text-green-600 bg-green-50' : 'text-slate-400 bg-slate-100'" class="px-2 py-1 rounded-full text-xs">
                {{ row.status === 1 ? $t('common.enabled') : $t('common.disabled') }}
              </span>
            </td>
            <td class="px-4 py-3">
              <button class="text-blue-600 hover:underline text-xs" @click="openEdit(row)">{{ $t('common.edit') }}</button>
            </td>
          </tr>
          <tr v-if="!list.length">
            <td colspan="6" class="px-4 py-12 text-center text-slate-400">{{ $t('common.noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showForm" class="fixed inset-0 bg-black/40 flex items-start justify-end z-50">
      <div class="bg-white w-[420px] h-full shadow-2xl p-6 overflow-y-auto">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-lg font-semibold text-slate-800">{{ editingID ? $t('common.edit') : $t('common.create') }}</h3>
          <button @click="showForm = false" class="text-slate-400 hover:text-slate-600">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('common.name') }}</label>
            <input v-model="form.name" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('marketingActivity.startAt') }}</label>
            <input v-model="form.start_at" type="datetime-local" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('marketingActivity.endAt') }}</label>
            <input v-model="form.end_at" type="datetime-local" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('common.status') }}</label>
            <select v-model.number="form.status" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm">
              <option :value="1">{{ $t('common.enabled') }}</option>
              <option :value="0">{{ $t('common.disabled') }}</option>
            </select>
          </div>
          <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
          <button
            @click="save"
            :disabled="saving"
            class="w-full bg-blue-700 text-white py-3 rounded-xl text-sm font-medium hover:bg-blue-600 transition disabled:opacity-60"
          >
            {{ saving ? $t('common.saving') : $t('common.save') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'

const props = defineProps<{ kind: 'seckill' | 'group-buy' | 'bargain' }>()
const { t } = useI18n()

const list = ref<any[]>([])
const showForm = ref(false)
const saving = ref(false)
const editingID = ref(0)
const error = ref('')
const form = reactive({
  name: '',
  start_at: '',
  end_at: '',
  status: 1,
})

const pageTitle = computed(() => {
  if (props.kind === 'seckill') return t('nav.seckillActivityManage')
  if (props.kind === 'group-buy') return t('nav.groupBuyActivityManage')
  return t('nav.bargainActivityManage')
})

const endpoint = computed(() => {
  if (props.kind === 'seckill') return '/seckill/activities'
  if (props.kind === 'group-buy') return '/group-buy/activities'
  return '/bargain/activities'
})

function toDateTimeInput(value?: string) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 16)
}

function toISO(value: string) {
  if (!value) return ''
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '' : date.toISOString()
}

function formatDate(value?: string) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

async function loadList() {
  const data: any = await request.get(endpoint.value)
  list.value = data?.list || []
}

function resetForm() {
  form.name = ''
  form.start_at = ''
  form.end_at = ''
  form.status = 1
  error.value = ''
}

function openCreate() {
  editingID.value = 0
  resetForm()
  showForm.value = true
}

function openEdit(row: any) {
  editingID.value = Number(row?.id || 0)
  form.name = String(row?.name || '')
  form.start_at = toDateTimeInput(row?.start_at)
  form.end_at = toDateTimeInput(row?.end_at)
  form.status = Number(row?.status || 0) === 1 ? 1 : 0
  error.value = ''
  showForm.value = true
}

async function save() {
  if (!form.name.trim()) {
    error.value = t('coupon.nameRequired')
    return
  }
  const startAt = toISO(form.start_at)
  const endAt = toISO(form.end_at)
  if (!startAt || !endAt || new Date(startAt).getTime() >= new Date(endAt).getTime()) {
    error.value = t('marketingActivity.invalidTimeRange')
    return
  }
  saving.value = true
  error.value = ''
  try {
    const payload = {
      name: form.name.trim(),
      start_at: startAt,
      end_at: endAt,
      status: form.status,
    }
    if (editingID.value > 0) {
      await request.put(`${endpoint.value}/${editingID.value}`, payload)
    } else {
      await request.post(endpoint.value, payload)
    }
    showForm.value = false
    await loadList()
  } catch (err: any) {
    error.value = err?.message || t('common.saveFailed')
  } finally {
    saving.value = false
  }
}

onMounted(loadList)
</script>
