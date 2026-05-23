<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <router-link to="/product/list" class="text-slate-400 hover:text-slate-600 text-sm">← 返回列表</router-link>
      <h2 class="text-xl font-semibold text-slate-800">{{ isEdit ? '编辑商品' : '新增商品' }}</h2>
    </div>

    <div class="grid grid-cols-1 xl:grid-cols-[2fr_1fr] gap-6">
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
        <div class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">商品名称 *</label>
            <input v-model="form.title" class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" placeholder="请输入商品名称" />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">售价 *</label>
              <input v-model.number="form.price" type="number" step="0.01"
                class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" placeholder="0.00" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">库存</label>
              <input v-model.number="form.stock" type="number"
                class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" placeholder="0" />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">分类</label>
            <select v-model="form.category_id"
              class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none">
              <option value="">请选择分类</option>
              <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">状态</label>
            <select v-model.number="form.status"
              class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none">
              <option :value="1">上架</option>
              <option :value="0">下架</option>
            </select>
          </div>

          <div>
            <div class="flex items-center justify-between mb-2">
              <label class="block text-sm font-medium text-slate-700">封面图URL</label>
              <span class="text-xs text-slate-400">可手输或由 AI 生成覆盖</span>
            </div>
            <input v-model="form.cover" class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" placeholder="https://..." />
          </div>

          <div>
            <div class="flex items-center justify-between mb-2">
              <label class="block text-sm font-medium text-slate-700">轮播图</label>
              <button class="text-xs text-blue-600 hover:underline" @click="addGalleryImage('')">+ 新增空白项</button>
            </div>
            <div class="space-y-2">
              <div v-for="(img, idx) in galleryImages" :key="idx" class="grid grid-cols-[1fr_auto] gap-2">
                <input v-model="img.url" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400" placeholder="https://..." />
                <button class="px-3 py-2 text-xs bg-slate-100 rounded-lg hover:bg-slate-200" @click="removeGalleryImage(idx)">删除</button>
              </div>
            </div>
          </div>

          <div>
            <div class="flex items-center justify-between mb-2">
              <label class="block text-sm font-medium text-slate-700">商品详情（JSON Blocks）</label>
              <span class="text-xs text-slate-400">支持 text/image，按顺序渲染</span>
            </div>
            <div class="space-y-3">
              <div v-for="(block, idx) in detailBlocks" :key="block.id" class="border border-slate-200 rounded-xl p-3">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-xs font-medium text-slate-500">Block {{ idx + 1 }} · {{ block.type }}</span>
                  <div class="flex gap-2">
                    <button class="px-2 py-1 text-xs rounded bg-slate-100 hover:bg-slate-200" @click="moveBlock(idx, -1)" :disabled="idx === 0">上移</button>
                    <button class="px-2 py-1 text-xs rounded bg-slate-100 hover:bg-slate-200" @click="moveBlock(idx, 1)" :disabled="idx === detailBlocks.length - 1">下移</button>
                    <button class="px-2 py-1 text-xs rounded bg-red-50 text-red-600 hover:bg-red-100" @click="removeBlock(idx)">删除</button>
                  </div>
                </div>
                <div v-if="block.type === 'text'">
                  <textarea
                    v-model="block.props.text"
                    rows="3"
                    class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 resize-none"
                    placeholder="输入文本内容"
                    @focus="currentBlockIndex = idx"
                  />
                </div>
                <div v-else-if="block.type === 'image'" class="space-y-2">
                  <input
                    v-model="block.props.url"
                    class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400"
                    placeholder="图片URL"
                    @focus="currentBlockIndex = idx"
                  />
                  <input
                    v-model="block.props.alt"
                    class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400"
                    placeholder="图片说明（可选）"
                    @focus="currentBlockIndex = idx"
                  />
                </div>
              </div>
              <div class="flex gap-2">
                <button class="px-3 py-2 text-xs rounded-lg bg-slate-100 hover:bg-slate-200" @click="addTextBlock()">+ 文本块</button>
                <button class="px-3 py-2 text-xs rounded-lg bg-slate-100 hover:bg-slate-200" @click="addImageBlock()">+ 图片块</button>
              </div>
            </div>
          </div>

          <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
          <div class="flex gap-3 pt-2">
            <button @click="save" :disabled="saving"
              class="px-6 py-3 bg-blue-700 text-white rounded-xl text-sm font-medium hover:bg-blue-600 transition disabled:opacity-60">
              {{ saving ? '保存中...' : '保 存' }}
            </button>
            <router-link to="/product/list"
              class="px-6 py-3 bg-slate-100 text-slate-600 rounded-xl text-sm font-medium hover:bg-slate-200 transition">
              取 消
            </router-link>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6 h-fit">
        <h3 class="font-semibold text-slate-700 mb-4">AI 图片助手</h3>
        <div class="space-y-3">
          <div>
            <label class="block text-xs text-slate-500 mb-1">生成目标</label>
            <select v-model="aiForm.biz_type" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm">
              <option value="cover">商品封面</option>
              <option value="gallery">轮播图</option>
              <option value="detail">详情图</option>
              <option value="intro">商品介绍图（预留）</option>
            </select>
          </div>
          <div>
            <label class="block text-xs text-slate-500 mb-1">模型</label>
            <select v-model.number="aiForm.model_id" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm">
              <option v-for="m in aiModels" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-xs text-slate-500 mb-1">描述 Prompt</label>
            <textarea v-model="aiForm.prompt" rows="3" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm resize-none" />
          </div>
          <div>
            <label class="block text-xs text-slate-500 mb-1">参考图（实物图）</label>
            <input
              type="file"
              accept="image/*"
              :disabled="!selectedModelSupportsRef"
              @change="onRefImageChange"
              class="w-full border border-slate-200 rounded-lg px-3 py-2 text-xs"
            />
            <p v-if="!selectedModelSupportsRef" class="text-xs text-orange-600 mt-1">该模型不支持参考图</p>
            <p v-if="aiForm.ref_image_url" class="text-xs text-slate-500 mt-1 truncate">已上传：{{ aiForm.ref_image_url }}</p>
          </div>
          <button
            class="w-full bg-blue-700 text-white py-2.5 rounded-lg text-sm hover:bg-blue-600 disabled:opacity-60"
            :disabled="aiGenerating || !aiForm.prompt.trim()"
            @click="generateWithAI"
          >
            {{ aiGenerating ? '生成中...' : '开始生成' }}
          </button>
          <div v-if="aiImages.length" class="grid grid-cols-2 gap-2 pt-1">
            <div v-for="(url, idx) in aiImages" :key="idx" class="border border-slate-200 rounded-lg p-1.5">
              <img :src="url" class="w-full h-24 object-cover rounded" />
              <button class="w-full mt-1 text-xs bg-slate-100 rounded py-1 hover:bg-slate-200" @click="applyAiImage(url)">
                应用此图
              </button>
            </div>
          </div>
          <p class="text-xs text-slate-400" v-if="aiForm.biz_type === 'detail'">提示：详情图将插入到当前编辑位置（Block {{ currentBlockIndex + 1 }} 之后）</p>
          <p v-if="aiNotice" class="text-xs text-emerald-600">{{ aiNotice }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createProduct, generateAiImage, getAiModels, getAiTask, getCategories, getProduct, updateProduct, uploadFile } from '@/api/plugins'

type DetailBlock = {
  id: string
  type: 'text' | 'image'
  props: Record<string, any>
}

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const categories = ref<any[]>([])
const saving = ref(false)
const error = ref('')
const currentBlockIndex = ref(0)
const galleryImages = ref<Array<{ url: string; sort: number }>>([])

const form = ref({
  title: '', price: 0, origin_price: 0, stock: 0,
  category_id: '', cover: '', status: 1,
})

const detailBlocks = ref<DetailBlock[]>([
  { id: `b-${Date.now()}`, type: 'text', props: { text: '' } }
])

const aiModels = ref<any[]>([])
const aiImages = ref<string[]>([])
const aiGenerating = ref(false)
const aiNotice = ref('')
const aiForm = ref({
  model_id: 0,
  scene: 'detail',
  biz_type: 'detail',
  prompt: '',
  ref_image_url: '',
  target_product_id: 0,
  params: { width: 750, height: 1000, count: 2, style: 'ecommerce' },
})

const selectedModel = computed(() => aiModels.value.find((m) => Number(m.id) === Number(aiForm.value.model_id)))
const selectedModelSupportsRef = computed(() => Number(selectedModel.value?.supports_ref_image || 0) === 1)

function makeDetailPayload() {
  return {
    version: 1,
    blocks: detailBlocks.value.map((block) => ({
      id: block.id,
      type: block.type,
      props: block.props,
    })),
  }
}

function parseDetail(detail: any) {
  const raw = typeof detail === 'string' ? (() => {
    try { return JSON.parse(detail) } catch { return null }
  })() : detail
  if (!raw || !Array.isArray(raw.blocks)) return
  detailBlocks.value = raw.blocks
    .filter((item: any) => item && (item.type === 'text' || item.type === 'image'))
    .map((item: any, idx: number) => ({
      id: item.id || `b-${idx}-${Date.now()}`,
      type: item.type,
      props: item.props || {},
    }))
  if (!detailBlocks.value.length) {
    detailBlocks.value = [{ id: `b-${Date.now()}`, type: 'text', props: { text: '' } }]
  }
}

function addTextBlock(position = detailBlocks.value.length) {
  detailBlocks.value.splice(position, 0, {
    id: `b-${Date.now()}-${Math.random().toString(16).slice(2, 6)}`,
    type: 'text',
    props: { text: '' },
  })
  currentBlockIndex.value = position
}

function addImageBlock(url = '', position = detailBlocks.value.length) {
  detailBlocks.value.splice(position, 0, {
    id: `b-${Date.now()}-${Math.random().toString(16).slice(2, 6)}`,
    type: 'image',
    props: { url, alt: '' },
  })
  currentBlockIndex.value = position
}

function removeBlock(index: number) {
  detailBlocks.value.splice(index, 1)
  if (!detailBlocks.value.length) addTextBlock(0)
  currentBlockIndex.value = Math.max(0, Math.min(currentBlockIndex.value, detailBlocks.value.length - 1))
}

function moveBlock(index: number, delta: number) {
  const target = index + delta
  if (target < 0 || target >= detailBlocks.value.length) return
  const [item] = detailBlocks.value.splice(index, 1)
  detailBlocks.value.splice(target, 0, item)
  currentBlockIndex.value = target
}

function addGalleryImage(url: string) {
  galleryImages.value.push({ url, sort: galleryImages.value.length })
}

function removeGalleryImage(index: number) {
  galleryImages.value.splice(index, 1)
  galleryImages.value.forEach((item, idx) => { item.sort = idx })
}

async function onRefImageChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  error.value = ''
  const result: any = await uploadFile(file)
  aiForm.value.ref_image_url = result?.url || ''
}

async function waitTaskDone(taskID: number, maxRetry = 20, intervalMs = 1200) {
  for (let i = 0; i < maxRetry; i += 1) {
    const detail: any = await getAiTask(taskID)
    if (Number(detail?.status) === 2 || Number(detail?.status) === 3) return detail
    await new Promise((resolve) => setTimeout(resolve, intervalMs))
  }
  throw new Error('AI 生成超时，请稍后到 AI 生图页面查看结果')
}

async function generateWithAI() {
  aiGenerating.value = true
  aiImages.value = []
  aiNotice.value = ''
  error.value = ''
  try {
    const task: any = await generateAiImage({
      model_id: aiForm.value.model_id,
      scene: aiForm.value.biz_type === 'detail' ? 'detail' : 'carousel',
      biz_type: aiForm.value.biz_type,
      prompt: aiForm.value.prompt,
      ref_image_url: aiForm.value.ref_image_url || undefined,
      target_product_id: aiForm.value.target_product_id || undefined,
      params: aiForm.value.params,
    })
    const taskDetail: any = await waitTaskDone(Number(task.id))
    if (Number(taskDetail?.status) === 3) {
      throw new Error(taskDetail?.error_msg || 'AI 生成失败')
    }
    if (taskDetail?.result_urls) {
      try {
        aiImages.value = JSON.parse(taskDetail.result_urls)
      } catch {
        aiImages.value = []
      }
    }
  } catch (e: any) {
    error.value = e.message || 'AI 生成失败'
  } finally {
    aiGenerating.value = false
  }
}

function applyAiImage(url: string) {
  if (!url) return
  aiNotice.value = ''
  if (aiForm.value.biz_type === 'cover') {
    form.value.cover = url
    aiNotice.value = '已应用到商品封面'
    return
  }
  if (aiForm.value.biz_type === 'gallery') {
    addGalleryImage(url)
    aiNotice.value = '已追加到轮播图'
    return
  }
  if (aiForm.value.biz_type === 'detail') {
    addImageBlock(url, Math.min(currentBlockIndex.value + 1, detailBlocks.value.length))
    aiNotice.value = `已插入到当前编辑位置（Block ${currentBlockIndex.value + 1}）`
    return
  }
  if (aiForm.value.biz_type === 'intro') {
    aiNotice.value = '商品介绍图已生成，当前版本已预留该能力'
  }
}

async function save() {
  if (!form.value.title) { error.value = '商品名称不能为空'; return }
  saving.value = true
  error.value = ''
  const payload = {
    ...form.value,
    detail: makeDetailPayload(),
  }
  const imagesPayload = galleryImages.value
    .filter((item) => item.url.trim())
    .map((item, idx) => ({ url: item.url.trim(), sort: idx }))

  try {
    if (isEdit.value) {
      await updateProduct(Number(route.params.id), { product: payload, images: imagesPayload })
    } else {
      await createProduct({ product: payload, skus: [], images: imagesPayload })
    }
    router.push('/product/list')
  } catch (e: any) {
    error.value = e.message || '保存失败'
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  categories.value = ((await getCategories()) || []) as any[]
  aiModels.value = ((await getAiModels()) || []) as any[]
  if (aiModels.value.length) aiForm.value.model_id = Number(aiModels.value[0].id)

  if (isEdit.value) {
    const data: any = await getProduct(Number(route.params.id))
    Object.assign(form.value, {
      title: data.title || '',
      price: data.price || 0,
      origin_price: data.origin_price || 0,
      stock: data.stock || 0,
      category_id: data.category_id || '',
      cover: data.cover || '',
      status: data.status ?? 1,
    })
    parseDetail(data.detail)
    galleryImages.value = Array.isArray(data.images) ? data.images.map((item: any, idx: number) => ({
      url: item.url || '',
      sort: item.sort ?? idx,
    })) : []
    aiForm.value.target_product_id = Number(route.params.id)
  }
})
</script>
