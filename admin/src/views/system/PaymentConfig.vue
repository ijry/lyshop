<template>
  <div>
    <h2 class="text-xl font-semibold text-slate-800 mb-6">配置中心</h2>

    <div v-if="loading" class="text-center py-12 text-slate-400">加载中...</div>

    <div v-else-if="!schemas.length" class="text-center py-12 text-slate-400">
      暂无可配置的插件
    </div>

    <div v-else class="flex gap-6">
      <!-- Plugin tabs -->
      <div class="w-48 shrink-0">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
          <button v-for="s in schemas" :key="s.plugin"
            @click="selectPlugin(s.plugin)"
            :class="activePlugin === s.plugin ? 'bg-red-50 text-red-600 font-medium border-l-3 border-l-red-600' : 'text-slate-600 hover:bg-slate-50'"
            class="w-full text-left px-4 py-3 text-sm transition-colors border-b border-slate-50 last:border-0">
            {{ s.title }}
          </button>
        </div>
      </div>

      <!-- Config form -->
      <div class="flex-1 max-w-xl" v-if="activeSchema">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <div class="flex items-center justify-between mb-5">
            <h3 class="text-base font-semibold text-slate-800">{{ activeSchema.title }} 配置</h3>
            <span class="text-xs text-slate-400">插件: {{ activeSchema.plugin }}</span>
          </div>

          <div class="space-y-4">
            <div v-for="field in activeSchema.fields" :key="field.key">
              <label class="block text-sm font-medium text-slate-700 mb-1.5">
                {{ field.label }}
                <span v-if="field.required" class="text-red-500 ml-0.5">*</span>
              </label>

              <!-- select -->
              <select v-if="field.type === 'select'" v-model="configValues[field.key]"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500">
                <option v-for="opt in field.options" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
              </select>

              <!-- switch -->
              <div v-else-if="field.type === 'switch'" class="flex items-center gap-2">
                <input type="checkbox" :id="field.key" v-model="configValues[field.key]"
                  true-value="true" false-value="false"
                  class="w-4 h-4 accent-red-600" />
                <label :for="field.key" class="text-sm text-slate-600">{{ configValues[field.key] === 'true' ? '已开启' : '已关闭' }}</label>
              </div>

              <!-- textarea -->
              <textarea v-else-if="field.type === 'textarea'" v-model="configValues[field.key]"
                rows="5" :placeholder="field.placeholder || ''"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm font-mono resize-y focus:outline-none focus:border-red-500" />

              <!-- password -->
              <div v-else-if="field.type === 'password'" class="relative">
                <input v-model="configValues[field.key]" :type="showPw[field.key] ? 'text' : 'password'"
                  :placeholder="field.placeholder || ''"
                  class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm pr-10 focus:outline-none focus:border-red-500" />
                <button @click="showPw[field.key] = !showPw[field.key]"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 text-xs">
                  {{ showPw[field.key] ? '隐藏' : '显示' }}
                </button>
              </div>

              <!-- text / number -->
              <input v-else v-model="configValues[field.key]"
                :type="field.type === 'number' ? 'number' : 'text'"
                :placeholder="field.placeholder || ''"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
            </div>
          </div>

          <div class="flex items-center gap-3 mt-6">
            <button @click="save" :disabled="saving"
              class="px-6 py-2.5 bg-red-600 text-white rounded-xl text-sm font-medium hover:bg-red-700 transition disabled:opacity-40">
              {{ saving ? '保存中...' : '保存配置' }}
            </button>
            <span v-if="savedMsg" class="text-sm text-green-600">{{ savedMsg }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import request from '@/api/request'

interface ConfigField {
  key: string; label: string; type: string;
  placeholder?: string; required?: boolean;
  options?: { label: string; value: string }[]
}
interface PluginSchema { plugin: string; title: string; fields: ConfigField[] }

const schemas = ref<PluginSchema[]>([])
const activePlugin = ref('')
const configValues = ref<Record<string, string>>({})
const showPw = ref<Record<string, boolean>>({})
const loading = ref(true)
const saving = ref(false)
const savedMsg = ref('')

const activeSchema = computed(() => schemas.value.find(s => s.plugin === activePlugin.value))

function applyFieldDefaults(schema?: PluginSchema) {
  if (!schema) return
  for (const field of schema.fields) {
    if (field.type !== 'select') continue
    if (configValues.value[field.key]) continue
    const first = field.options?.[0]
    if (first?.value) configValues.value[field.key] = first.value
  }
}

async function selectPlugin(pluginName: string) {
  activePlugin.value = pluginName
  const data: any = await request.get(`/config/${pluginName}`)
  configValues.value = data || {}
  applyFieldDefaults(schemas.value.find(s => s.plugin === pluginName))
  showPw.value = {}
}

async function save() {
  saving.value = true
  try {
    await request.put(`/config/${activePlugin.value}`, configValues.value)
    savedMsg.value = '保存成功'
    setTimeout(() => savedMsg.value = '', 2000)
  } catch (e: any) {
    savedMsg.value = '保存失败: ' + (e.message || '')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    const data: any = await request.get('/config/schemas')
    const loaded = (data || []) as PluginSchema[]
    loaded.sort((a, b) => {
      if (a.plugin === 'storage_router') return -1
      if (b.plugin === 'storage_router') return 1
      return 0
    })
    schemas.value = loaded
    if (schemas.value.length) selectPlugin(schemas.value[0].plugin)
  } finally {
    loading.value = false
  }
})
</script>
