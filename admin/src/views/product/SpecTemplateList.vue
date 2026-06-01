<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">规格模板</h2>
      <button
        class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition"
        @click="openCreate"
      >
        新建模板
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm p-4 mb-4 flex gap-3 border border-slate-100">
      <input
        v-model="query.keyword"
        placeholder="搜索模板名称"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm flex-1 focus:outline-none focus:border-blue-400"
      />
      <button
        class="px-4 py-2 bg-slate-100 rounded-lg text-sm hover:bg-slate-200"
        @click="loadTemplates"
      >
        搜索
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">ID</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">模板名称</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">分类ID</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">属性组</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">状态</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="item in list" :key="item.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 text-slate-400">{{ item.id }}</td>
            <td class="px-4 py-3 font-medium text-slate-700">{{ item.name }}</td>
            <td class="px-4 py-3 text-slate-600">{{ formatCategoryIDs(item.category_ids) }}</td>
            <td class="px-4 py-3 text-slate-600">{{ formatAttrs(item.attrs) }}</td>
            <td class="px-4 py-3">
              <span
                :class="Number(item.status || 0) === 1 ? 'text-green-600 bg-green-50' : 'text-slate-400 bg-slate-100'"
                class="px-2 py-1 rounded-full text-xs"
              >
                {{ Number(item.status || 0) === 1 ? '启用' : '禁用' }}
              </span>
            </td>
            <td class="px-4 py-3">
              <button class="text-blue-600 hover:underline text-xs mr-3" @click="openEdit(item)">编辑</button>
              <button class="text-red-500 hover:underline text-xs" @click="remove(item.id)">删除</button>
            </td>
          </tr>
          <tr v-if="!list.length">
            <td colspan="6" class="px-4 py-12 text-center text-slate-400">暂无模板</td>
          </tr>
        </tbody>
      </table>
      <div class="px-4 py-3 flex items-center justify-between border-t border-slate-100 text-sm text-slate-500">
        <span>共 {{ total }} 条</span>
        <div class="flex gap-2">
          <button
            :disabled="query.page <= 1"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40"
            @click="query.page--; loadTemplates()"
          >
            上一页
          </button>
          <button
            :disabled="query.page * query.size >= total"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40"
            @click="query.page++; loadTemplates()"
          >
            下一页
          </button>
        </div>
      </div>
    </div>

    <div v-if="showDialog" class="fixed inset-0 bg-black/35 flex items-center justify-center z-50" @click.self="showDialog = false">
      <div class="bg-white rounded-xl w-[720px] max-w-[95vw] p-6">
        <h3 class="text-lg font-semibold text-slate-800 mb-4">{{ form.id ? '编辑模板' : '新建模板' }}</h3>
        <div class="grid grid-cols-1 gap-3">
          <input v-model="form.name" placeholder="模板名称" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm" />
          <input
            v-model="form.categoryIDsText"
            placeholder="分类ID，逗号分隔（例：11,12）"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm"
          />
          <textarea
            v-model="form.attrsText"
            rows="8"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm font-mono"
            placeholder='属性组 JSON，例如：[{"name":"颜色","values":["黑色","白色"]},{"name":"尺码","values":["M","L"]}]'
          />
          <label class="flex items-center gap-2 text-sm text-slate-700">
            <input v-model="form.status" type="checkbox" />
            启用模板
          </label>
        </div>
        <div class="mt-5 flex justify-end gap-2">
          <button class="px-4 py-2 rounded-lg border border-slate-200 text-sm" @click="showDialog = false">取消</button>
          <button class="px-4 py-2 rounded-lg bg-blue-700 text-white text-sm" @click="save">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { createSpecTemplate, deleteSpecTemplate, getSpecTemplates, updateSpecTemplate } from '@/api/plugins'
import { confirmAction } from '@/utils/dialog'

type SpecTemplateAttr = { name: string; values: string[] }

const list = ref<any[]>([])
const total = ref(0)
const query = ref({ keyword: '', page: 1, size: 20 })
const showDialog = ref(false)
const form = ref({
  id: 0,
  name: '',
  categoryIDsText: '',
  attrsText: '[]',
  status: true,
})

function formatCategoryIDs(raw: any) {
  const ids = Array.isArray(raw) ? raw : []
  return ids.length ? ids.join(', ') : '-'
}

function formatAttrs(raw: any) {
  const attrs = Array.isArray(raw) ? raw : []
  if (!attrs.length) return '-'
  return attrs.map((item: any) => `${item.name}(${Array.isArray(item.values) ? item.values.length : 0})`).join(' / ')
}

function parseCategoryIDs(text: string) {
  return text
    .split(',')
    .map((item) => Number(String(item || '').trim()))
    .filter((id) => Number.isInteger(id) && id > 0)
}

function normalizeAttrs(raw: SpecTemplateAttr[]) {
  return raw
    .map((item) => ({
      name: String(item?.name || '').trim(),
      values: Array.isArray(item?.values)
        ? item.values.map((value: string) => String(value || '').trim()).filter((value: string) => !!value)
        : [],
    }))
    .filter((item) => item.name && item.values.length)
}

async function loadTemplates() {
  const data: any = await getSpecTemplates(query.value)
  list.value = Array.isArray(data?.list) ? data.list : []
  total.value = Number(data?.total || 0)
}

function openCreate() {
  form.value = {
    id: 0,
    name: '',
    categoryIDsText: '',
    attrsText: '[]',
    status: true,
  }
  showDialog.value = true
}

function openEdit(item: any) {
  form.value = {
    id: Number(item.id || 0),
    name: String(item.name || ''),
    categoryIDsText: Array.isArray(item.category_ids) ? item.category_ids.join(',') : '',
    attrsText: JSON.stringify(Array.isArray(item.attrs) ? item.attrs : [], null, 2),
    status: Number(item.status || 0) === 1,
  }
  showDialog.value = true
}

async function save() {
  const name = String(form.value.name || '').trim()
  if (!name) {
    window.alert('请填写模板名称')
    return
  }
  let attrs: SpecTemplateAttr[] = []
  try {
    const parsed = JSON.parse(String(form.value.attrsText || '[]'))
    if (!Array.isArray(parsed)) throw new Error('属性组格式错误')
    attrs = normalizeAttrs(parsed as SpecTemplateAttr[])
  } catch {
    window.alert('属性组 JSON 格式错误')
    return
  }
  const payload = {
    name,
    category_ids: parseCategoryIDs(form.value.categoryIDsText),
    attrs,
    status: form.value.status ? 1 : 0,
  }
  if (form.value.id > 0) {
    await updateSpecTemplate(form.value.id, payload)
  } else {
    await createSpecTemplate(payload)
  }
  showDialog.value = false
  await loadTemplates()
}

async function remove(id: number) {
  if (!confirmAction('确认删除该模板？')) return
  await deleteSpecTemplate(id)
  await loadTemplates()
}

onMounted(loadTemplates)
</script>

