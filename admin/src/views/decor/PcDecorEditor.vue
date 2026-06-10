<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-3">
        <h2 class="text-xl font-semibold text-slate-800">PC 首页装修</h2>
        <select v-model="currentVariantKey" @change="changeVariant"
          class="border border-slate-200 rounded-lg px-3 py-1.5 text-sm text-slate-700 bg-white">
          <option v-for="v in variants" :key="v.variant_key" :value="v.variant_key">
            {{ v.variant_name }}（{{ v.variant_key }}）{{ v.is_published ? ' · ' + t('decor.published') : '' }}
          </option>
        </select>
        <button @click="copyVariant"
          class="px-3 py-1.5 bg-slate-100 text-slate-700 rounded-lg text-xs hover:bg-slate-200 transition">
          {{ t('decor.copyVariant') }}
        </button>
        <button @click="renameVariant"
          class="px-3 py-1.5 bg-slate-100 text-slate-700 rounded-lg text-xs hover:bg-slate-200 transition">
          {{ t('decor.rename') }}
        </button>
        <button @click="deleteVariant"
          class="px-3 py-1.5 bg-red-50 text-red-600 rounded-lg text-xs hover:bg-red-100 transition">
          {{ t('decor.deleteVariant') }}
        </button>
      </div>
      <div class="flex gap-2">
        <button @click="save" :disabled="saving"
          class="px-4 py-2 bg-slate-100 text-slate-700 rounded-xl text-sm hover:bg-slate-200 transition">
          {{ saving ? '保存中...' : t('decor.saveDraft') }}
        </button>
        <button @click="publish"
          class="px-4 py-2 bg-blue-700 text-white rounded-xl text-sm hover:bg-blue-600 transition">
          {{ t('decor.publish') }}
        </button>
      </div>
    </div>

    <div class="flex gap-4" style="height: calc(100vh - 160px); min-height: 600px">
      <!-- Component library -->
      <div class="w-48 bg-white rounded-xl border border-slate-100 p-4 shrink-0 overflow-y-auto">
        <p class="text-xs font-medium text-slate-500 mb-3">组件库</p>
        <div class="space-y-2">
          <div v-for="comp in pcComponentLib" :key="comp.type"
            draggable="true" @dragstart="dragStart(comp)" @click="appendComp(comp)"
            class="flex items-center gap-2 px-3 py-2 border border-slate-200 rounded-xl cursor-grab hover:border-blue-300 hover:bg-blue-50 transition text-sm text-slate-700">
            <span>{{ comp.icon }}</span>
            <span>{{ comp.title }}</span>
          </div>
        </div>

        <p class="text-xs font-medium text-slate-500 mt-6 mb-3">已添加组件</p>
        <div class="space-y-1" @dragover.prevent @drop="onDrop">
          <div v-for="(c, i) in components" :key="c.id"
            @click="selectComp(i)"
            :class="selectedIndex === i ? 'bg-blue-50 border-blue-400 text-blue-700' : 'border-slate-200 text-slate-600 hover:border-slate-300'"
            class="flex items-center justify-between px-2.5 py-1.5 border rounded-lg cursor-pointer transition text-xs">
            <span class="truncate">{{ pcCompTitleMap[c.type] || c.type }}</span>
            <div class="flex gap-0.5 shrink-0">
              <button @click.stop="moveUp(i)" :disabled="i===0" class="text-slate-400 hover:text-slate-600 px-0.5 disabled:opacity-30">↑</button>
              <button @click.stop="moveDown(i)" :disabled="i===components.length-1" class="text-slate-400 hover:text-slate-600 px-0.5 disabled:opacity-30">↓</button>
              <button @click.stop="remove(i)" class="text-red-400 hover:text-red-600 px-0.5">×</button>
            </div>
          </div>
          <div v-if="!components.length" class="text-center py-6 text-slate-300 text-xs border-2 border-dashed border-slate-200 rounded-lg">
            从组件库拖入组件
          </div>
        </div>
      </div>

      <!-- Center: inline PC preview -->
      <div class="flex-1 flex flex-col bg-white rounded-xl border border-slate-100 overflow-hidden min-w-0">
        <div class="flex items-center justify-between px-4 py-2 border-b border-slate-100 shrink-0">
          <span class="text-xs text-slate-500">PC 预览</span>
          <div class="flex items-center gap-2">
            <button @click="previewScale = Math.max(0.3, previewScale - 0.1)"
              class="text-xs text-slate-400 hover:text-slate-600 px-1">-</button>
            <span class="text-xs text-slate-400 w-10 text-center">{{ Math.round(previewScale * 100) }}%</span>
            <button @click="previewScale = Math.min(1, previewScale + 0.1)"
              class="text-xs text-slate-400 hover:text-slate-600 px-1">+</button>
          </div>
        </div>
        <div class="flex-1 overflow-auto bg-slate-50 p-4">
          <div class="mx-auto flex justify-center" :style="{ width: `${1280 * previewScale}px`, minWidth: `${1280 * previewScale}px` }">
            <div class="border border-slate-200 rounded-xl overflow-hidden shadow-lg bg-white"
              :style="{ width: '1280px', transformOrigin: 'top center', transform: `scale(${previewScale})` }">
              <PcDecorPreview :components="components" :pageStyle="pagePayload.pageStyle" @select="selectComp" />
            </div>
          </div>
        </div>
      </div>

      <!-- Right: property editors -->
      <div class="w-80 bg-white rounded-xl border border-slate-100 p-4 shrink-0 overflow-y-auto">
        <p class="text-xs font-medium text-slate-500 mb-3">属性配置</p>
        <div class="mb-5">
          <p class="text-xs text-slate-500 mb-2">页面样式</p>
          <PcPageStyleEditor v-model="pagePayload.pageStyle" />
        </div>
        <div v-if="selectedComp">
          <div class="mb-4">
            <span class="text-sm font-medium text-slate-700">{{ pcCompTitleMap[selectedComp.type] || selectedComp.type }}</span>
          </div>
          <HeroEditor            v-if="selectedComp.type === 'hero'"             v-model="selectedComp.props" />
          <BannerEditor          v-else-if="selectedComp.type === 'banner'"       v-model="selectedComp.props" />
          <CategoryNavEditor     v-else-if="selectedComp.type === 'category_nav'" v-model="selectedComp.props" />
          <GridEditor            v-else-if="selectedComp.type === 'grid'"         v-model="selectedComp.props" />
          <ProductGridEditor     v-else-if="selectedComp.type === 'product_grid'" v-model="selectedComp.props" />
          <NoticeEditor          v-else-if="selectedComp.type === 'notice'"       v-model="selectedComp.props" />
          <ImageAdEditor         v-else-if="selectedComp.type === 'image_ad'"     v-model="selectedComp.props" />
          <RichTextEditor        v-else-if="selectedComp.type === 'rich_text'"    v-model="selectedComp.props" />
          <MarketingZoneEditor   v-else-if="selectedComp.type === 'marketing_zone'" v-model="selectedComp.props" />
          <FeaturesEditor        v-else-if="selectedComp.type === 'features'"     v-model="selectedComp.props" />
          <SpacerEditor          v-else-if="selectedComp.type === 'spacer'"       v-model="selectedComp.props" />
          <PcComponentStyleEditor v-model="selectedComp.style" />
        </div>
        <div v-else class="text-center py-4 text-slate-300 text-sm">点击左侧组件进行编辑</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'
import { notify } from '@/utils/notify'
import { confirmAction, promptText } from '@/utils/dialog'
import {
  pcComponentLib,
  pcCompTitleMap,
  createPcDefaultProps,
  createDefaultPcDecorPayload,
  type PcDecorPagePayload,
} from '@/types/decor'

import PcDecorPreview from './PcDecorPreview.vue'
import HeroEditor from './editors/HeroEditor.vue'
import BannerEditor from './editors/BannerEditor.vue'
import CategoryNavEditor from './editors/CategoryNavEditor.vue'
import GridEditor from './editors/GridEditor.vue'
import ProductGridEditor from './editors/ProductGridEditor.vue'
import NoticeEditor from './editors/NoticeEditor.vue'
import ImageAdEditor from './editors/ImageAdEditor.vue'
import RichTextEditor from './editors/RichTextEditor.vue'
import MarketingZoneEditor from './editors/MarketingZoneEditor.vue'
import FeaturesEditor from './editors/FeaturesEditor.vue'
import SpacerEditor from './editors/SpacerEditor.vue'
import PcPageStyleEditor from './widgets/PcPageStyleEditor.vue'
import PcComponentStyleEditor from './widgets/PcComponentStyleEditor.vue'

const { t } = useI18n()
const pagePayload = ref<PcDecorPagePayload>(createDefaultPcDecorPayload())
const variants = ref<any[]>([])
const currentVariantKey = ref('default')
const selectedIndex = ref<number | null>(null)
const saving = ref(false)
const previewScale = ref(0.8)

let draggedComp: any = null

const components = computed<any[]>({
  get: () => pagePayload.value.components || [],
  set: (list) => {
    pagePayload.value.components = Array.isArray(list) ? list : []
  },
})

const selectedComp = computed(() =>
  selectedIndex.value !== null ? components.value[selectedIndex.value] : null
)

function dragStart(comp: any) { draggedComp = comp }

function appendComp(comp: any) {
  components.value.push({
    type: comp.type,
    id: `pc_${Date.now()}_${Math.random().toString(16).slice(2, 6)}`,
    props: createPcDefaultProps(comp.type),
    style: {},
  })
  selectedIndex.value = components.value.length - 1
}

function onDrop() {
  if (!draggedComp) return
  appendComp(draggedComp)
  draggedComp = null
}

function selectComp(i: number) { selectedIndex.value = i }

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

async function save() {
  saving.value = true
  try {
    await request.put(`/decor/pc?variant=${encodeURIComponent(currentVariantKey.value)}`, { components: pagePayload.value })
    await loadVariants()
    notify('保存成功')
  } finally { saving.value = false }
}

async function publish() {
  await save()
  await request.post(`/decor/pc/publish?variant=${encodeURIComponent(currentVariantKey.value)}`)
  await loadVariants()
  notify(t('decor.publishedNote'))
}

async function loadVariants() {
  const data: any = await request.get('/decor/pc/variants')
  variants.value = Array.isArray(data) ? data : []
  if (!variants.value.length) {
    variants.value = [{
      variant_key: 'default',
      variant_name: t('decor.variantName', { key: 'default' }),
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
  const data: any = await request.get(`/decor/pc?variant=${encodeURIComponent(currentVariantKey.value)}`)
  pagePayload.value = normalizePayload(data?.components)
  selectedIndex.value = null
}

async function changeVariant() {
  await loadCurrentVariant()
}

function toVariantKey(raw: string) {
  return raw.trim().toLowerCase().replace(/\s+/g, '_').replace(/[^a-z0-9_-]/g, '')
}

async function copyVariant() {
  const keyRaw = promptText(t('decor.promptKey'))
  if (!keyRaw) return
  const newVariantKey = toVariantKey(keyRaw)
  if (!newVariantKey) {
    notify(t('decor.invalidKey'))
    return
  }
  const defaultName = t('decor.variantName', { key: newVariantKey })
  const newVariantName = promptText(t('decor.promptName'), defaultName) || defaultName
  await request.post('/decor/pc/copies', {
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
  const next = promptText(t('decor.promptName'), current?.variant_name || '')
  if (!next) return
  await request.put(`/decor/pc/variants/${encodeURIComponent(currentVariantKey.value)}`, {
    variant_name: next,
  })
  await loadVariants()
}

async function deleteVariant() {
  if (currentVariantKey.value === 'default') {
    notify(t('decor.defaultNoDelete'))
    return
  }
  if (!confirmAction(t('decor.confirmDeleteVariant', { key: currentVariantKey.value }))) return
  await request.delete(`/decor/pc/variants/${encodeURIComponent(currentVariantKey.value)}`)
  await loadVariants()
  await loadCurrentVariant()
}

onMounted(async () => {
  await loadVariants()
  await loadCurrentVariant()
})

function normalizePayload(raw: any): PcDecorPagePayload {
  const fallback = createDefaultPcDecorPayload()
  let parsed = raw
  if (typeof raw === 'string') {
    try {
      parsed = JSON.parse(raw)
    } catch {
      return fallback
    }
  }
  if (!parsed || typeof parsed !== 'object') return fallback
  if (!Array.isArray(parsed.components)) return fallback
  return {
    pageStyle: {
      ...fallback.pageStyle,
      ...(parsed.pageStyle || {}),
      background: {
        ...fallback.pageStyle.background,
        ...(parsed.pageStyle?.background || {}),
        gradient: {
          ...fallback.pageStyle.background.gradient,
          ...(parsed.pageStyle?.background?.gradient || {}),
        },
        image: {
          ...fallback.pageStyle.background.image,
          ...(parsed.pageStyle?.background?.image || {}),
        },
        overlay: {
          ...fallback.pageStyle.background.overlay,
          ...(parsed.pageStyle?.background?.overlay || {}),
        },
      },
      content: {
        ...fallback.pageStyle.content,
        ...(parsed.pageStyle?.content || {}),
      },
      surface: {
        ...fallback.pageStyle.surface,
        ...(parsed.pageStyle?.surface || {}),
      },
    },
    components: parsed.components,
  }
}
</script>
