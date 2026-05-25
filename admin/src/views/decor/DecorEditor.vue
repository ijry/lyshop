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

    <div class="flex gap-4 h-full" style="min-height: 600px">
      <!-- Component library -->
      <div class="w-52 bg-white rounded-xl border border-slate-100 p-4 shrink-0">
        <p class="text-xs font-medium text-slate-500 mb-3">组件库（拖拽到画布）</p>
        <div class="space-y-2">
          <div v-for="comp in componentLib" :key="comp.type"
            draggable="true" @dragstart="dragStart(comp)"
            class="flex items-center gap-2 px-3 py-2 border border-slate-200 rounded-xl cursor-grab hover:border-blue-300 hover:bg-blue-50 transition text-sm text-slate-700">
            <span>{{ comp.icon }}</span>
            <span>{{ comp.title }}</span>
          </div>
        </div>
      </div>

      <!-- Canvas -->
      <div class="flex-1 bg-white rounded-xl border border-slate-100 p-4 overflow-y-auto"
        @dragover.prevent @drop="onDrop">
        <p class="text-xs text-slate-400 text-center mb-4">拖拽组件到此处</p>
        <div class="max-w-sm mx-auto space-y-2">
          <div v-for="(c, i) in components" :key="c.id"
            @click="selectComp(i)"
            :class="selectedIndex === i ? 'ring-2 ring-blue-500' : 'hover:ring-1 hover:ring-slate-300'"
            class="relative bg-slate-50 rounded-xl p-3 cursor-pointer transition">
            <div class="flex justify-between items-center">
              <span class="text-sm text-slate-700 font-medium">{{ compTitle(c.type) }}</span>
              <div class="flex gap-1">
                <button @click.stop="moveUp(i)" :disabled="i===0" class="text-slate-400 hover:text-slate-600 text-xs px-1">↑</button>
                <button @click.stop="moveDown(i)" :disabled="i===components.length-1" class="text-slate-400 hover:text-slate-600 text-xs px-1">↓</button>
                <button @click.stop="remove(i)" class="text-red-400 hover:text-red-600 text-xs px-1">×</button>
              </div>
            </div>
            <p class="text-xs text-slate-400 mt-1">{{ compPreview(c) }}</p>
          </div>
          <div v-if="!components.length" class="text-center py-16 text-slate-300 text-sm border-2 border-dashed border-slate-200 rounded-xl">
            拖拽组件到这里
          </div>
        </div>
      </div>

      <!-- Props panel -->
      <div class="w-64 bg-white rounded-xl border border-slate-100 p-4 shrink-0">
        <p class="text-xs font-medium text-slate-500 mb-3">属性配置</p>
        <div v-if="selectedComp" class="space-y-3">
          <div>
            <label class="block text-xs text-slate-500 mb-1">组件类型</label>
            <span class="text-sm font-medium text-slate-700">{{ compTitle(selectedComp.type) }}</span>
          </div>
          <div>
            <label class="block text-xs text-slate-500 mb-1">JSON 属性</label>
            <textarea v-model="propsJson" rows="10"
              class="w-full border border-slate-200 rounded-xl px-3 py-2 text-xs font-mono resize-none focus:outline-none focus:border-blue-400"
              @change="updateProps" />
          </div>
        </div>
        <div v-else class="text-center py-8 text-slate-300 text-sm">选择组件编辑属性</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import request from '@/api/request'
import { notify } from '@/utils/notify'
import { confirmAction, promptText } from '@/utils/dialog'

const components = ref<any[]>([])
const variants = ref<any[]>([])
const currentVariantKey = ref('default')
const selectedIndex = ref<number | null>(null)
const saving = ref(false)
const propsJson = ref('')

let draggedComp: any = null

const componentLib = [
  { type: 'banner',          title: '轮播图',   icon: '🖼' },
  { type: 'category_nav',    title: '分类导航', icon: '📂' },
  { type: 'product_grid',    title: '商品列表', icon: '🛍' },
  { type: 'notice',          title: '公告栏',   icon: '📢' },
  { type: 'image_ad',        title: '广告图',   icon: '🎯' },
  { type: 'rich_text',       title: '富文本',   icon: '📝' },
  { type: 'marketing_zone',  title: '营销区块', icon: '🏷' },
  { type: 'spacer',          title: '间距',     icon: '↕' },
]

const compTitles: Record<string, string> = Object.fromEntries(componentLib.map(c => [c.type, c.title]))
const compTitle = (type: string) => compTitles[type] || type

const compPreview = (c: any) => {
  if (c.type === 'banner') return `${c.props?.images?.length || 0} 张图片`
  if (c.type === 'product_grid') return `来源: ${c.props?.source || 'hot'}, 限 ${c.props?.limit || 10} 条`
  if (c.type === 'notice') return `${c.props?.items?.length || (c.props?.text ? 1 : 0)} 条公告`
  return JSON.stringify(c.props || {}).slice(0, 40)
}

const selectedComp = computed(() =>
  selectedIndex.value !== null ? components.value[selectedIndex.value] : null
)

function dragStart(comp: any) { draggedComp = comp }

function onDrop() {
  if (!draggedComp) return
  components.value.push({
    type: draggedComp.type,
    id: `c_${Date.now()}`,
    props: defaultProps(draggedComp.type),
  })
  draggedComp = null
}

function defaultProps(type: string): any {
  const defaults: Record<string, any> = {
    banner:       { images: [], height: 350 },
    category_nav: { items: [] },
    product_grid: { source: 'hot', limit: 10, columns: 2 },
    notice:       {
      items: [
        { text: '欢迎来到 LYShop', link: '/pages/index/index' },
        { text: '新人下单立减，优惠券限时领取', link: '/pages/marketing/coupon' },
        { text: '精选好物每日上新，支持多端下单', link: '/pages/product/list' },
      ],
      color: '#f97316',
      bgColor: '#fff7ed',
      duration: 2500,
      mode: 'link',
    },
    image_ad:     { url: '', link: '' },
    rich_text:    { content: '' },
    marketing_zone: { type: 'seckill' },
    spacer:       { height: 16, background: '#ffffff' },
  }
  return defaults[type] || {}
}

function selectComp(i: number) {
  selectedIndex.value = i
  propsJson.value = JSON.stringify(components.value[i].props, null, 2)
}

function updateProps() {
  if (selectedIndex.value === null) return
  try {
    components.value[selectedIndex.value].props = JSON.parse(propsJson.value)
  } catch {}
}

function moveUp(i: number) {
  if (i === 0) return
  const arr = [...components.value]
  ;[arr[i-1], arr[i]] = [arr[i], arr[i-1]]
  components.value = arr
}

function moveDown(i: number) {
  if (i >= components.value.length - 1) return
  const arr = [...components.value]
  ;[arr[i], arr[i+1]] = [arr[i+1], arr[i]]
  components.value = arr
}

function remove(i: number) {
  components.value.splice(i, 1)
  if (selectedIndex.value === i) selectedIndex.value = null
}

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
  propsJson.value = ''
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
  await loadVariants()
  await loadCurrentVariant()
})
</script>
