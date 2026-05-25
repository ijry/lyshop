<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-3">
        <h2 class="text-xl font-semibold text-slate-800">首页装修</h2>
        <select v-model="currentVariantKey" @change="changeVariant"
          class="border border-slate-200 rounded-lg px-3 py-1.5 text-sm text-slate-700 bg-white">
          <option v-for="v in variants" :key="v.variant_key" :value="v.variant_key">
            {{ v.variant_name }}（{{ v.variant_key }}）{{ v.is_published ? ' · 已发布' : '' }}
          </option>
        </select>
        <button @click="copyVariant"
          class="px-3 py-1.5 bg-slate-100 text-slate-700 rounded-lg text-xs hover:bg-slate-200 transition">
          复制副本
        </button>
        <button @click="renameVariant"
          class="px-3 py-1.5 bg-slate-100 text-slate-700 rounded-lg text-xs hover:bg-slate-200 transition">
          重命名
        </button>
        <button @click="deleteVariant"
          class="px-3 py-1.5 bg-red-50 text-red-600 rounded-lg text-xs hover:bg-red-100 transition">
          删除副本
        </button>
      </div>
      <div class="flex gap-2">
        <button @click="save" :disabled="saving"
          class="px-4 py-2 bg-slate-100 text-slate-700 rounded-xl text-sm hover:bg-slate-200 transition">
          {{ saving ? '保存中...' : '保存草稿' }}
        </button>
        <button @click="publish"
          class="px-4 py-2 bg-blue-700 text-white rounded-xl text-sm hover:bg-blue-600 transition">
          发布上线
        </button>
      </div>
    </div>

    <div class="flex gap-4" style="height: calc(100vh - 160px); min-height: 600px">
      <!-- Component library -->
      <div class="w-48 bg-white rounded-xl border border-slate-100 p-4 shrink-0 overflow-y-auto">
        <p class="text-xs font-medium text-slate-500 mb-3">组件库（拖拽到画布）</p>
        <div class="space-y-2">
          <div v-for="comp in componentLib" :key="comp.type"
            draggable="true" @dragstart="dragStart(comp)"
            class="flex items-center gap-2 px-3 py-2 border border-slate-200 rounded-xl cursor-grab hover:border-blue-300 hover:bg-blue-50 transition text-sm text-slate-700">
            <span>{{ comp.icon }}</span>
            <span>{{ comp.title }}</span>
          </div>
        </div>

        <!-- Component list in canvas -->
        <p class="text-xs font-medium text-slate-500 mt-6 mb-3">画布组件</p>
        <div class="space-y-1" @dragover.prevent @drop="onDrop">
          <div v-for="(c, i) in components" :key="c.id"
            @click="selectComp(i)"
            :class="selectedIndex === i ? 'bg-blue-50 border-blue-400 text-blue-700' : 'border-slate-200 text-slate-600 hover:border-slate-300'"
            class="flex items-center justify-between px-2.5 py-1.5 border rounded-lg cursor-pointer transition text-xs">
            <span class="truncate">{{ compTitle(c.type) }}</span>
            <div class="flex gap-0.5 shrink-0">
              <button @click.stop="moveUp(i)" :disabled="i===0" class="text-slate-400 hover:text-slate-600 px-0.5 disabled:opacity-30">↑</button>
              <button @click.stop="moveDown(i)" :disabled="i===components.length-1" class="text-slate-400 hover:text-slate-600 px-0.5 disabled:opacity-30">↓</button>
              <button @click.stop="remove(i)" class="text-red-400 hover:text-red-600 px-0.5">×</button>
            </div>
          </div>
          <div v-if="!components.length" class="text-center py-6 text-slate-300 text-xs border-2 border-dashed border-slate-200 rounded-lg">
            拖拽组件到这里
          </div>
        </div>
      </div>

      <!-- Center: iframe preview -->
      <div class="flex-1 flex flex-col bg-white rounded-xl border border-slate-100 overflow-hidden min-w-0">
        <div class="flex items-center justify-between px-4 py-2 border-b border-slate-100 shrink-0">
          <span class="text-xs text-slate-500">实时预览</span>
          <div class="flex items-center gap-2">
            <span v-if="!previewReady" class="text-xs text-orange-500">等待 H5 就绪...</span>
            <span v-else class="text-xs text-emerald-500">已连接</span>
            <button @click="refreshPreview" class="text-xs text-blue-600 hover:underline">刷新</button>
          </div>
        </div>
        <div class="flex-1 flex items-start justify-center p-4 bg-slate-50 overflow-auto">
          <div class="w-[375px] h-[667px] border border-slate-200 rounded-2xl overflow-hidden shadow-lg bg-white shrink-0">
            <iframe
              ref="previewIframe"
              :src="previewUrl"
              class="w-full h-full border-none"
              @load="onIframeLoad"
            />
          </div>
        </div>
      </div>

      <!-- Right: property editors -->
      <div class="w-80 bg-white rounded-xl border border-slate-100 p-4 shrink-0 overflow-y-auto">
        <p class="text-xs font-medium text-slate-500 mb-3">属性配置</p>
        <div v-if="selectedComp">
          <div class="mb-4">
            <span class="text-sm font-medium text-slate-700">{{ compTitle(selectedComp.type) }}</span>
          </div>
          <BannerEditor       v-if="selectedComp.type === 'banner'"          v-model="selectedComp.props" />
          <CategoryNavEditor  v-else-if="selectedComp.type === 'category_nav'"  v-model="selectedComp.props" />
          <ProductGridEditor  v-else-if="selectedComp.type === 'product_grid'"  v-model="selectedComp.props" />
          <NoticeEditor       v-else-if="selectedComp.type === 'notice'"        v-model="selectedComp.props" />
          <ImageAdEditor      v-else-if="selectedComp.type === 'image_ad'"      v-model="selectedComp.props" />
          <RichTextEditor     v-else-if="selectedComp.type === 'rich_text'"     v-model="selectedComp.props" />
          <MarketingZoneEditor v-else-if="selectedComp.type === 'marketing_zone'" v-model="selectedComp.props" />
          <SpacerEditor       v-else-if="selectedComp.type === 'spacer'"        v-model="selectedComp.props" />
        </div>
        <div v-else class="text-center py-8 text-slate-300 text-sm">选择组件编辑属性</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import request from '@/api/request'
import { notify } from '@/utils/notify'
import { confirmAction, promptText } from '@/utils/dialog'
import { componentLib, compTitleMap, createDefaultProps } from '@/types/decor'

import BannerEditor from './editors/BannerEditor.vue'
import CategoryNavEditor from './editors/CategoryNavEditor.vue'
import ProductGridEditor from './editors/ProductGridEditor.vue'
import NoticeEditor from './editors/NoticeEditor.vue'
import ImageAdEditor from './editors/ImageAdEditor.vue'
import RichTextEditor from './editors/RichTextEditor.vue'
import MarketingZoneEditor from './editors/MarketingZoneEditor.vue'
import SpacerEditor from './editors/SpacerEditor.vue'

const components = ref<any[]>([])
const variants = ref<any[]>([])
const currentVariantKey = ref('default')
const selectedIndex = ref<number | null>(null)
const saving = ref(false)

// iframe preview state
const previewIframe = ref<HTMLIFrameElement>()
const previewReady = ref(false)
const previewUrl = import.meta.env.VITE_H5_PREVIEW_URL || 'http://localhost:5173/?preview=1'

let draggedComp: any = null

const compTitle = (type: string) => compTitleMap[type] || type

const selectedComp = computed(() =>
  selectedIndex.value !== null ? components.value[selectedIndex.value] : null
)

// ---- Drag & Drop ----
function dragStart(comp: any) { draggedComp = comp }

function onDrop() {
  if (!draggedComp) return
  components.value.push({
    type: draggedComp.type,
    id: `c_${Date.now()}`,
    props: createDefaultProps(draggedComp.type),
  })
  draggedComp = null
}

// ---- Component operations ----
function selectComp(i: number) {
  selectedIndex.value = i
}

function moveUp(i: number) {
  if (i === 0) return
  const arr = [...components.value]
  ;[arr[i-1], arr[i]] = [arr[i], arr[i-1]]
  components.value = arr
  if (selectedIndex.value === i) selectedIndex.value = i - 1
}

function moveDown(i: number) {
  if (i >= components.value.length - 1) return
  const arr = [...components.value]
  ;[arr[i], arr[i+1]] = [arr[i+1], arr[i]]
  components.value = arr
  if (selectedIndex.value === i) selectedIndex.value = i + 1
}

function remove(i: number) {
  components.value.splice(i, 1)
  if (selectedIndex.value === i) selectedIndex.value = null
}

// ---- iframe preview ----
function sendPreviewUpdate() {
  if (!previewReady.value || !previewIframe.value?.contentWindow) return
  previewIframe.value.contentWindow.postMessage({
    type: 'DECOR_PREVIEW_UPDATE',
    source: 'lyshop-admin',
    components: JSON.parse(JSON.stringify(components.value)),
  }, '*')
}

function onPreviewMessage(e: MessageEvent) {
  if (e.data?.type === 'DECOR_PREVIEW_READY' && e.data?.source === 'lyshop-app') {
    previewReady.value = true
    sendPreviewUpdate()
  }
}

function onIframeLoad() {
  // iframe might reload, reset readiness
  previewReady.value = false
}

function refreshPreview() {
  previewReady.value = false
  if (previewIframe.value) {
    previewIframe.value.src = previewUrl
  }
}

let debounceTimer: ReturnType<typeof setTimeout> | null = null
watch(components, () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(sendPreviewUpdate, 300)
}, { deep: true })

// ---- Variant management ----
async function save() {
  saving.value = true
  try {
    await request.put(`/decor/index?variant=${encodeURIComponent(currentVariantKey.value)}`, { components: components.value })
    await loadVariants()
  } finally { saving.value = false }
}

async function publish() {
  await save()
  await request.post(`/decor/index/publish?variant=${encodeURIComponent(currentVariantKey.value)}`)
  await loadVariants()
  notify('已发布上线（单活发布）')
}

async function loadVariants() {
  const data: any = await request.get('/decor/index/variants')
  variants.value = Array.isArray(data) ? data : []
  if (!variants.value.length) {
    variants.value = [{
      variant_key: 'default',
      variant_name: '默认副本',
      is_published: false,
    }]
  }
  const currentExists = variants.value.some(v => v.variant_key === currentVariantKey.value)
  if (!currentExists) {
    const published = variants.value.find(v => v.is_published)
    currentVariantKey.value = published?.variant_key || variants.value[0].variant_key || 'default'
  }
}

async function loadCurrentVariant() {
  const data: any = await request.get(`/decor/index?variant=${encodeURIComponent(currentVariantKey.value)}`)
  if (data?.components) {
    try { components.value = JSON.parse(data.components) } catch {}
  } else {
    components.value = []
  }
  selectedIndex.value = null
}

async function changeVariant() {
  await loadCurrentVariant()
}

function toVariantKey(raw: string) {
  return raw.trim().toLowerCase().replace(/\s+/g, '_').replace(/[^a-z0-9_-]/g, '')
}

async function copyVariant() {
  const keyRaw = promptText('请输入新副本标识（如 spring_festival_2027）')
  if (!keyRaw) return
  const newVariantKey = toVariantKey(keyRaw)
  if (!newVariantKey) {
    notify('副本标识不合法')
    return
  }
  const newVariantName = promptText('请输入新副本名称', `副本 ${newVariantKey}`) || `副本 ${newVariantKey}`
  await request.post('/decor/index/copies', {
    from_variant_key: currentVariantKey.value,
    new_variant_key: newVariantKey,
    new_variant_name: newVariantName,
  })
  await loadVariants()
  currentVariantKey.value = newVariantKey
  await loadCurrentVariant()
}

async function renameVariant() {
  const current = variants.value.find(v => v.variant_key === currentVariantKey.value)
  const next = promptText('请输入副本名称', current?.variant_name || '')
  if (!next) return
  await request.put(`/decor/index/variants/${encodeURIComponent(currentVariantKey.value)}`, {
    variant_name: next,
  })
  await loadVariants()
}

async function deleteVariant() {
  if (currentVariantKey.value === 'default') {
    notify('默认副本不支持删除')
    return
  }
  if (!confirmAction(`确认删除副本 ${currentVariantKey.value}？`)) return
  await request.delete(`/decor/index/variants/${encodeURIComponent(currentVariantKey.value)}`)
  await loadVariants()
  await loadCurrentVariant()
}

onMounted(async () => {
  window.addEventListener('message', onPreviewMessage)
  await loadVariants()
  await loadCurrentVariant()
})

onUnmounted(() => {
  window.removeEventListener('message', onPreviewMessage)
  if (debounceTimer) clearTimeout(debounceTimer)
})
</script>
