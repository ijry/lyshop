<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('ai.imageGen.title') }}</h2>
    </div>

    <!-- Generate form -->
    <div class="grid grid-cols-2 gap-6 mb-6">
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
        <h3 class="font-semibold text-slate-700 mb-4">{{ $t('ai.imageGen.subtitle') }}</h3>
        <div class="space-y-3">
          <div>
            <label class="block text-sm font-medium text-slate-600 mb-1">{{ $t('ai.imageGen.usage') }}</label>
            <select v-model="form.scene" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm">
              <option value="carousel">{{ $t('ai.imageGen.carousel') }}</option>
              <option value="detail">{{ $t('ai.imageGen.detail') }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-600 mb-1">{{ $t('ai.imageGen.model') }}</label>
            <select v-model="form.model_id" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm">
              <option value="">{{ $t('ai.imageGen.defaultModel') }}</option>
              <option v-for="m in models" :key="m.id" :value="m.id">{{ m.name }} ({{ m.driver }})</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-600 mb-1">{{ $t('ai.imageGen.prompt') }}</label>
            <textarea v-model="form.prompt" rows="3" :placeholder="$t('ai.imageGen.promptPlaceholder')"
              class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm resize-none" />
          </div>
          <div class="grid grid-cols-3 gap-2">
            <div>
              <label class="block text-xs text-slate-500 mb-1">{{ $t('ai.imageGen.style') }}</label>
              <select v-model="form.style" class="w-full border border-slate-200 rounded-lg px-2 py-1.5 text-xs">
                <option value="ecommerce">{{ $t('ai.imageGen.styleEcommerce') }}</option>
                <option value="realistic">{{ $t('ai.imageGen.styleRealistic') }}</option>
                <option value="illustration">{{ $t('ai.imageGen.styleIllustration') }}</option>
              </select>
            </div>
            <div>
              <label class="block text-xs text-slate-500 mb-1">{{ $t('ai.imageGen.count') }}</label>
              <input v-model.number="form.count" type="number" min="1" max="5"
                class="w-full border border-slate-200 rounded-lg px-2 py-1.5 text-xs" />
            </div>
            <div>
              <label class="block text-xs text-slate-500 mb-1">{{ $t('ai.imageGen.size') }}</label>
              <select v-model="form.size" class="w-full border border-slate-200 rounded-lg px-2 py-1.5 text-xs">
                <option value="750x750">750x750</option>
                <option value="750x1000">750x1000</option>
                <option value="1000x750">1000x750</option>
              </select>
            </div>
          </div>
          <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
          <button @click="generate" :disabled="generating"
            class="w-full bg-blue-700 text-white py-3 rounded-xl text-sm font-medium hover:bg-blue-600 transition disabled:opacity-60">
            {{ generating ? $t('ai.imageGen.generating') : $t('ai.imageGen.generate') }}
          </button>
        </div>
      </div>

      <!-- Result panel -->
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
        <h3 class="font-semibold text-slate-700 mb-4">{{ $t('ai.imageGen.result') }}</h3>
        <div v-if="!currentTask" class="flex items-center justify-center h-48 text-slate-300">
          <span>{{ $t('ai.imageGen.resultPlaceholder') }}</span>
        </div>
        <div v-else>
          <div class="flex items-center gap-2 mb-4">
            <span :class="taskStatusClass(currentTask.status)" class="px-2 py-1 rounded-full text-xs">
              {{ taskStatusLabel(currentTask.status) }}
            </span>
            <button v-if="currentTask.status === 1" @click="pollTask"
              class="text-xs text-blue-600 hover:underline">{{ $t('ai.imageGen.refreshStatus') }}</button>
          </div>
          <div v-if="resultURLs.length" class="grid grid-cols-2 gap-2">
            <div v-for="(url, i) in resultURLs" :key="i" class="relative group">
              <img :src="url" class="w-full rounded-xl object-cover" style="aspect-ratio:1" />
              <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 rounded-xl transition flex items-center justify-center gap-2">
                <a :href="url" download class="text-white text-xs bg-white/20 px-3 py-1 rounded-lg hover:bg-white/30">{{ $t('ai.imageGen.download') }}</a>
              </div>
            </div>
          </div>
          <p v-if="currentTask.error_msg" class="text-red-500 text-sm">{{ currentTask.error_msg }}</p>
        </div>
      </div>
    </div>

    <!-- History -->
    <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
      <h3 class="font-semibold text-slate-700 mb-4">{{ $t('ai.imageGen.history') }}</h3>
      <div class="space-y-2">
        <div v-for="t in tasks" :key="t.id"
          @click="viewTask(t)"
          class="flex items-center gap-4 p-3 rounded-xl hover:bg-slate-50 cursor-pointer transition">
          <span class="text-xs text-slate-400 w-6">#{{ t.id }}</span>
          <span :class="taskStatusClass(t.status)" class="text-xs px-2 py-0.5 rounded-full shrink-0">{{ taskStatusLabel(t.status) }}</span>
          <span class="text-sm text-slate-700 flex-1 truncate">{{ t.prompt }}</span>
          <span class="text-xs text-slate-400">{{ t.scene }}</span>
        </div>
        <div v-if="!tasks.length" class="text-center py-6 text-slate-400 text-sm">{{ $t('ai.imageGen.noHistory') }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'

const { t } = useI18n()

const models = ref<any[]>([])
const tasks = ref<any[]>([])
const currentTask = ref<any>(null)
const generating = ref(false)
const error = ref('')

const form = ref({
  prompt: '', scene: 'carousel', model_id: '', style: 'ecommerce',
  count: 3, size: '750x750'
})

const resultURLs = computed(() => {
  if (!currentTask.value?.result_urls) return []
  try { return JSON.parse(currentTask.value.result_urls) } catch { return [] }
})

const taskStatusLabels: Record<number, string> = { 1: t('ai.imageGen.generating'), 2: t('ai.imageGen.statusCompleted'), 3: t('ai.imageGen.statusFailed') }
const taskStatusColors: Record<number, string> = {
  1: 'bg-yellow-50 text-yellow-600', 2: 'bg-green-50 text-green-600', 3: 'bg-red-50 text-red-500'
}
const taskStatusLabel = (s: number) => taskStatusLabels[s] || ''
const taskStatusClass = (s: number) => taskStatusColors[s] || ''

async function generate() {
  if (!form.value.prompt) { error.value = t('ai.imageGen.promptRequired'); return }
  generating.value = true; error.value = ''
  try {
    const [w, h] = form.value.size.split('x').map(Number)
    const task: any = await request.post('/ai/generate', {
      model_id: form.value.model_id || undefined,
      scene: form.value.scene,
      prompt: form.value.prompt,
      params: { width: w, height: h, count: form.value.count, style: form.value.style }
    })
    currentTask.value = task
    // Poll until done
    const timer = setInterval(async () => {
      const t: any = await request.get(`/ai/tasks/${task.id}`)
      currentTask.value = t
      if (t.status !== 1) {
        clearInterval(timer)
        generating.value = false
        loadTasks()
      }
    }, 2000)
  } catch (e: any) {
    error.value = e.message
    generating.value = false
  }
}

async function pollTask() {
  if (!currentTask.value) return
  const t: any = await request.get(`/ai/tasks/${currentTask.value.id}`)
  currentTask.value = t
}

function viewTask(t: any) { currentTask.value = t }

async function loadTasks() {
  const data: any = await request.get('/ai/tasks', { params: { size: 20 } })
  tasks.value = data.list || []
}

onMounted(async () => {
  const data: any = await request.get('/ai/models')
  models.value = data || []
  loadTasks()
})
</script>
