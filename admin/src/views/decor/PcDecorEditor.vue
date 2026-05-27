<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-xl font-semibold text-slate-800">PC 首页装修</h2>
      <div class="flex gap-2">
        <button @click="save" :disabled="saving"
          class="px-4 py-2 bg-slate-100 text-slate-700 rounded-xl text-sm hover:bg-slate-200 transition">
          {{ saving ? '保存中...' : '保存' }}
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
        <p class="text-xs font-medium text-slate-500 mb-3">组件库</p>
        <div class="space-y-2">
          <div v-for="comp in pcComponentLib" :key="comp.type"
            draggable="true" @dragstart="dragStart(comp)"
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

      <!-- Center: PC preview -->
      <div class="flex-1 flex flex-col bg-white rounded-xl border border-slate-100 overflow-hidden min-w-0">
        <div class="flex items-center justify-between px-4 py-2 border-b border-slate-100 shrink-0">
          <span class="text-xs text-slate-500">PC 预览（1280px）</span>
        </div>
        <div class="flex-1 overflow-auto bg-slate-50 p-4">
          <div class="mx-auto border border-slate-200 rounded-xl overflow-hidden shadow-lg bg-white" style="width: 1280px; transform-origin: top left; transform: scale(0.6);">
            <iframe ref="previewIframe" :src="previewUrl" class="w-full border-none" style="height: 1600px; pointer-events: none;" />
          </div>
        </div>
      </div>

      <!-- Right: property editors -->
      <div class="w-80 bg-white rounded-xl border border-slate-100 p-4 shrink-0 overflow-y-auto">
        <p class="text-xs font-medium text-slate-500 mb-3">属性配置</p>
        <div v-if="selectedComp">
          <div class="mb-4">
            <span class="text-sm font-medium text-slate-700">{{ pcCompTitleMap[selectedComp.type] || selectedComp.type }}</span>
          </div>
          <HeroEditor          v-if="selectedComp.type === 'hero'"           v-model="selectedComp.props" />
          <BannerEditor        v-else-if="selectedComp.type === 'banner'"       v-model="selectedComp.props" />
          <CategoryNavEditor   v-else-if="selectedComp.type === 'category_nav'" v-model="selectedComp.props" />
          <GridEditor          v-else-if="selectedComp.type === 'grid'"         v-model="selectedComp.props" />
          <ProductGridEditor   v-else-if="selectedComp.type === 'product_grid'" v-model="selectedComp.props" />
          <NoticeEditor        v-else-if="selectedComp.type === 'notice'"       v-model="selectedComp.props" />
          <ImageAdEditor       v-else-if="selectedComp.type === 'image_ad'"     v-model="selectedComp.props" />
          <RichTextEditor      v-else-if="selectedComp.type === 'rich_text'"    v-model="selectedComp.props" />
          <MarketingZoneEditor v-else-if="selectedComp.type === 'marketing_zone'" v-model="selectedComp.props" />
          <FeaturesEditor      v-else-if="selectedComp.type === 'features'"     v-model="selectedComp.props" />
          <SpacerEditor        v-else-if="selectedComp.type === 'spacer'"       v-model="selectedComp.props" />
        </div>
        <div v-else class="text-center py-8 text-slate-300 text-sm">点击左侧组件进行编辑</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import request from '@/api/request'
import { notify } from '@/utils/notify'
import { pcComponentLib, pcCompTitleMap, createPcDefaultProps } from '@/types/decor'

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

const components = ref<any[]>([])
const selectedIndex = ref<number | null>(null)
const saving = ref(false)

const previewUrl = (import.meta.env.VITE_PC_PREVIEW_URL || 'http://localhost:5174') + '/?preview=1'

let draggedComp: any = null

const selectedComp = computed(() =>
  selectedIndex.value !== null ? components.value[selectedIndex.value] : null
)

function dragStart(comp: any) { draggedComp = comp }

function onDrop() {
  if (!draggedComp) return
  components.value.push({
    type: draggedComp.type,
    id: `pc_${Date.now()}`,
    props: createPcDefaultProps(draggedComp.type),
  })
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
    await request.put('/decor/pc', { components: components.value })
    notify('保存成功')
  } finally { saving.value = false }
}

async function publish() {
  await save()
  await request.post('/decor/pc/publish')
  notify('已发布上线')
}

onMounted(async () => {
  const data: any = await request.get('/decor/pc')
  if (data?.components) {
    try {
      components.value = typeof data.components === 'string'
        ? JSON.parse(data.components)
        : data.components
    } catch {}
  }
})
</script>
