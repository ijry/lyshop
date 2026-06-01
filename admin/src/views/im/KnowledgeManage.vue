<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-2">
      <h1 class="text-xl font-bold text-slate-800">{{ $t('imkb.title') }}</h1>
      <div class="flex items-center gap-3">
        <button @click="testAI" class="text-sm text-slate-600 hover:text-slate-800 border border-slate-200 rounded-lg px-3 py-1.5">
          {{ $t('imkb.testAI') }}
        </button>
        <button @click="testRerank" class="text-sm text-slate-600 hover:text-slate-800 border border-slate-200 rounded-lg px-3 py-1.5">
          {{ $t('imkb.testRerank') }}
        </button>
        <button @click="reindex" class="text-sm text-indigo-600 hover:text-indigo-700 border border-indigo-200 rounded-lg px-3 py-1.5">
          {{ $t('imkb.reindex') }}
        </button>
        <button @click="openImport" class="text-sm text-emerald-600 hover:text-emerald-700 border border-emerald-200 rounded-lg px-3 py-1.5">
          {{ $t('imkb.import') }}
        </button>
        <button @click="openCreate" class="btn-primary">+ {{ $t('imkb.add') }}</button>
      </div>
    </div>
    <p class="text-sm text-slate-400 mb-6">{{ $t('imkb.subtitle') }}</p>

    <div class="mb-4 flex items-center gap-2">
      <input v-model="keyword" @keyup.enter="reload" :placeholder="$t('imkb.searchPlaceholder')"
        class="w-72 border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20" />
    </div>

    <div class="card">
      <table class="w-full">
        <thead>
          <tr class="border-b border-slate-100">
            <th class="text-left py-3 px-4 text-sm font-semibold text-slate-600">{{ $t('imkb.colTitle') }}</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-slate-600">{{ $t('imkb.colContent') }}</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-slate-600">{{ $t('imkb.colTags') }}</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-slate-600">{{ $t('imkb.colIndexed') }}</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-slate-600">{{ $t('imkb.colStatus') }}</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-slate-600">{{ $t('imkb.colSort') }}</th>
            <th class="text-right py-3 px-4 text-sm font-semibold text-slate-600">{{ $t('imkb.colActions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="k in list" :key="k.id" class="border-b border-slate-50 hover:bg-slate-50 align-top">
            <td class="py-3 px-4 text-sm text-slate-800 font-medium max-w-48 truncate">{{ k.title }}</td>
            <td class="py-3 px-4 text-sm text-slate-500 max-w-80">
              <div class="line-clamp-2">{{ k.content }}</div>
            </td>
            <td class="py-3 px-4 text-sm text-slate-500 max-w-40 truncate">{{ k.tags }}</td>
            <td class="py-3 px-4">
              <span :class="k.indexed ? 'bg-indigo-100 text-indigo-700' : 'bg-slate-100 text-slate-500'"
                class="text-xs px-2 py-1 rounded-full whitespace-nowrap">
                {{ k.indexed ? $t('imkb.indexed') : $t('imkb.notIndexed') }}
              </span>
            </td>
            <td class="py-3 px-4">
              <span :class="k.status ? 'bg-green-100 text-green-700' : 'bg-slate-100 text-slate-500'"
                class="text-xs px-2 py-1 rounded-full whitespace-nowrap">
                {{ k.status ? $t('imkb.statusOn') : $t('imkb.statusOff') }}
              </span>
            </td>
            <td class="py-3 px-4 text-sm text-slate-700">{{ k.sort }}</td>
            <td class="py-3 px-4 text-right whitespace-nowrap">
              <button @click="openEdit(k)" class="text-sm text-blue-600 hover:text-blue-700 mr-3">{{ $t('imkb.edit') }}</button>
              <button @click="remove(k.id)" class="text-sm text-red-600 hover:text-red-700">{{ $t('imkb.delete') }}</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="!list.length" class="text-center py-12 text-slate-400 text-sm">{{ $t('imkb.empty') }}</div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showModal = false">
      <div class="bg-white rounded-2xl p-6 w-[32rem] shadow-xl">
        <h3 class="text-lg font-semibold mb-4">{{ editing ? $t('imkb.edit') : $t('imkb.add') }}</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">{{ $t('imkb.formTitle') }} *</label>
            <input v-model="form.title"
              class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">{{ $t('imkb.formContent') }} *</label>
            <textarea v-model="form.content" rows="5"
              class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">{{ $t('imkb.formTags') }}</label>
            <input v-model="form.tags"
              class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20" />
          </div>
          <div class="flex items-center gap-6">
            <div class="flex items-center gap-2">
              <label class="text-sm font-medium text-slate-700">{{ $t('imkb.formSort') }}</label>
              <input v-model.number="form.sort" type="number"
                class="w-24 border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400" />
            </div>
            <label class="flex items-center gap-2 text-sm font-medium text-slate-700">
              <input v-model="form.status" type="checkbox" :true-value="1" :false-value="0" />
              {{ $t('imkb.formStatus') }}
            </label>
          </div>
        </div>
        <div class="flex gap-3 mt-6">
          <button @click="showModal = false"
            class="flex-1 px-4 py-2 border border-slate-200 text-slate-600 rounded-lg text-sm font-medium hover:bg-slate-50 transition">
            {{ $t('imkb.cancel') }}
          </button>
          <button @click="submit"
            class="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg text-sm font-medium hover:bg-blue-600 transition">
            {{ $t('imkb.save') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Import Document Modal -->
    <div v-if="showImport" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showImport = false">
      <div class="bg-white rounded-2xl p-6 w-[34rem] shadow-xl">
        <h3 class="text-lg font-semibold mb-1">{{ $t('imkb.import') }}</h3>
        <p class="text-sm text-slate-400 mb-4">{{ $t('imkb.importHint') }}</p>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">{{ $t('imkb.importFile') }} *</label>
            <input ref="fileInput" type="file" :accept="acceptExts"
              class="w-full text-sm text-slate-600 file:mr-3 file:py-2 file:px-4 file:rounded-lg file:border-0 file:bg-emerald-50 file:text-emerald-700 file:text-sm file:font-medium hover:file:bg-emerald-100" />
            <p class="text-xs text-slate-400 mt-1">{{ $t('imkb.importFormats') }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">{{ $t('imkb.formTitle') }}</label>
            <input v-model="importForm.title" :placeholder="$t('imkb.importTitlePlaceholder')"
              class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">{{ $t('imkb.formTags') }}</label>
            <input v-model="importForm.tags"
              class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400/20" />
          </div>
          <div class="flex items-center gap-6">
            <div class="flex items-center gap-2">
              <label class="text-sm font-medium text-slate-700">{{ $t('imkb.chunkSize') }}</label>
              <input v-model.number="importForm.chunk_size" type="number"
                class="w-24 border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400" />
            </div>
            <div class="flex items-center gap-2">
              <label class="text-sm font-medium text-slate-700">{{ $t('imkb.overlap') }}</label>
              <input v-model.number="importForm.overlap" type="number"
                class="w-24 border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400" />
            </div>
          </div>
        </div>
        <div class="flex gap-3 mt-6">
          <button @click="showImport = false" :disabled="importing"
            class="flex-1 px-4 py-2 border border-slate-200 text-slate-600 rounded-lg text-sm font-medium hover:bg-slate-50 transition disabled:opacity-50">
            {{ $t('imkb.cancel') }}
          </button>
          <button @click="submitImport" :disabled="importing"
            class="flex-1 px-4 py-2 bg-emerald-500 text-white rounded-lg text-sm font-medium hover:bg-emerald-600 transition disabled:opacity-50">
            {{ importing ? $t('imkb.importing') : $t('imkb.import') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'

const { t } = useI18n()

const list = ref<any[]>([])
const keyword = ref('')
const showModal = ref(false)
const editing = ref<any>(null)
const form = ref({ title: '', content: '', tags: '', sort: 0, status: 1 })

// Document import
const showImport = ref(false)
const importing = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const importForm = ref({ title: '', tags: '', chunk_size: 500, overlap: 50 })
const acceptExts = '.txt,.md,.markdown,.text,.log,.csv,.tsv,.json,.xml,.html,.htm,.docx,.pdf,.xlsx'

async function reload() {
  const data: any = await request.get('/im/knowledge', { params: { keyword: keyword.value, size: 100 } })
  list.value = data?.list || []
}

function openCreate() {
  editing.value = null
  form.value = { title: '', content: '', tags: '', sort: 0, status: 1 }
  showModal.value = true
}

function openEdit(k: any) {
  editing.value = k
  form.value = { title: k.title, content: k.content, tags: k.tags || '', sort: k.sort || 0, status: k.status }
  showModal.value = true
}

async function submit() {
  if (!form.value.title.trim() || !form.value.content.trim()) {
    alert(t('imkb.formTitle') + ' / ' + t('imkb.formContent'))
    return
  }
  try {
    if (editing.value) {
      await request.put(`/im/knowledge/${editing.value.id}`, form.value)
    } else {
      await request.post('/im/knowledge', form.value)
    }
    showModal.value = false
    await reload()
    alert(t('imkb.saved'))
  } catch (err: any) {
    alert(err.message || 'error')
  }
}

async function remove(id: number) {
  if (!confirm(t('imkb.confirmDelete'))) return
  try {
    await request.delete(`/im/knowledge/${id}`)
    await reload()
    alert(t('imkb.deleted'))
  } catch (err: any) {
    alert(err.message || 'error')
  }
}

async function reindex() {
  try {
    const data: any = await request.post('/im/knowledge/reindex')
    await reload()
    alert(t('imkb.reindexDone', { count: data?.indexed ?? 0 }))
  } catch (err: any) {
    alert(err.message || 'error')
  }
}

async function testAI() {
  try {
    const data: any = await request.post('/im/ai/test')
    alert(t('imkb.aiReplyPrefix') + (data?.reply || ''))
  } catch (err: any) {
    alert(err.message || 'error')
  }
}

async function testRerank() {
  try {
    const data: any = await request.post('/im/ai/rerank-test')
    alert(t('imkb.rerankReplyPrefix') + (data?.reply || ''))
  } catch (err: any) {
    alert(err.message || 'error')
  }
}

function openImport() {
  importForm.value = { title: '', tags: '', chunk_size: 500, overlap: 50 }
  showImport.value = true
}

async function submitImport() {
  const file = fileInput.value?.files?.[0]
  if (!file) {
    alert(t('imkb.importFile'))
    return
  }
  const fd = new FormData()
  fd.append('file', file)
  fd.append('title', importForm.value.title)
  fd.append('tags', importForm.value.tags)
  fd.append('chunk_size', String(importForm.value.chunk_size || 0))
  fd.append('overlap', String(importForm.value.overlap || 0))
  importing.value = true
  try {
    const data: any = await request.post('/im/knowledge/import', fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    showImport.value = false
    await reload()
    alert(t('imkb.importDone', { count: data?.chunks ?? 0 }))
  } catch (err: any) {
    alert(err.message || 'error')
  } finally {
    importing.value = false
  }
}

onMounted(reload)
</script>
