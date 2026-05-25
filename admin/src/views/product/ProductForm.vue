<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <router-link to="/product/list" class="text-slate-400 hover:text-slate-600 text-sm">{{ $t('common.backToList') }}</router-link>
      <h2 class="text-xl font-semibold text-slate-800">{{ isEdit ? $t('product.form.editTitle') : $t('product.form.addTitle') }}</h2>
    </div>

    <div class="grid grid-cols-1 xl:grid-cols-[2fr_1fr] gap-6">
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
        <div class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('product.form.name') }}</label>
            <input v-model="form.title" class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" :placeholder="$t('product.form.namePlaceholder')" />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('product.form.price') }}</label>
              <input v-model.number="form.price" type="number" step="0.01"
                class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" placeholder="0.00" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('product.form.stock') }}</label>
              <input v-model.number="form.stock" type="number"
                class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" placeholder="0" />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('product.form.category') }}</label>
            <select v-model="form.category_id"
              class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none">
              <option value="">{{ $t('product.form.categoryPlaceholder') }}</option>
              <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('common.status') }}</label>
            <select v-model.number="form.status"
              class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none">
              <option :value="1">{{ $t('product.form.onSale') }}</option>
              <option :value="0">{{ $t('product.form.offSale') }}</option>
            </select>
          </div>

          <div>
            <div class="flex items-center justify-between mb-2">
              <label class="block text-sm font-medium text-slate-700">{{ $t('product.form.coverUrl') }}</label>
              <span class="text-xs text-slate-400">{{ $t('product.form.coverHint') }}</span>
            </div>
            <input v-model="form.cover" class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-400" placeholder="https://..." />
          </div>

          <div>
            <div class="flex items-center justify-between mb-2">
              <label class="block text-sm font-medium text-slate-700">{{ $t('product.form.carousel') }}</label>
              <button class="text-xs text-blue-600 hover:underline" @click="addGalleryImage('')">{{ $t('product.form.addBlank') }}</button>
            </div>
            <div class="space-y-2">
              <div v-for="(img, idx) in galleryImages" :key="idx" class="grid grid-cols-[1fr_auto] gap-2">
                <input v-model="img.url" class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400" placeholder="https://..." />
                <button class="px-3 py-2 text-xs bg-slate-100 rounded-lg hover:bg-slate-200" @click="removeGalleryImage(idx)">{{ $t('common.delete') }}</button>
              </div>
            </div>
          </div>

          <div>
            <div class="flex items-center justify-between mb-2">
              <label class="block text-sm font-medium text-slate-700">{{ $t('product.form.detail') }}</label>
              <span class="text-xs text-slate-400">{{ $t('product.form.detailHint') }}</span>
            </div>
            <div class="space-y-3">
              <div v-for="(block, idx) in detailBlocks" :key="block.id" class="border border-slate-200 rounded-xl p-3">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-xs font-medium text-slate-500">{{ $t('product.form.blockLabel', { idx: idx + 1, type: block.type }) }}</span>
                  <div class="flex gap-2">
                    <button class="px-2 py-1 text-xs rounded bg-slate-100 hover:bg-slate-200" @click="moveBlock(idx, -1)" :disabled="idx === 0">{{ $t('product.form.moveUp') }}</button>
                    <button class="px-2 py-1 text-xs rounded bg-slate-100 hover:bg-slate-200" @click="moveBlock(idx, 1)" :disabled="idx === detailBlocks.length - 1">{{ $t('product.form.moveDown') }}</button>
                    <button class="px-2 py-1 text-xs rounded bg-red-50 text-red-600 hover:bg-red-100" @click="removeBlock(idx)">{{ $t('common.delete') }}</button>
                  </div>
                </div>
                <div v-if="block.type === 'text'">
                  <textarea
                    v-model="block.props.text"
                    rows="3"
                    class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400 resize-none"
                    :placeholder="$t('product.form.textContent')"
                    @focus="currentBlockIndex = idx"
                  />
                </div>
                <div v-else-if="block.type === 'image'" class="space-y-2">
                  <input
                    v-model="block.props.url"
                    class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400"
                    :placeholder="$t('product.form.imageUrl')"
                    @focus="currentBlockIndex = idx"
                  />
                  <input
                    v-model="block.props.alt"
                    class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400"
                    :placeholder="$t('product.form.imageAlt')"
                    @focus="currentBlockIndex = idx"
                  />
                </div>
              </div>
              <div class="flex gap-2">
                <button class="px-3 py-2 text-xs rounded-lg bg-slate-100 hover:bg-slate-200" @click="addTextBlock()">{{ $t('product.form.addText') }}</button>
                <button class="px-3 py-2 text-xs rounded-lg bg-slate-100 hover:bg-slate-200" @click="addImageBlock()">{{ $t('product.form.addImage') }}</button>
              </div>
            </div>
          </div>

          <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
          <div class="flex gap-3 pt-2">
            <button @click="save" :disabled="saving"
              class="px-6 py-3 bg-blue-700 text-white rounded-xl text-sm font-medium hover:bg-blue-600 transition disabled:opacity-60">
              {{ saving ? $t('common.saving') : $t('common.save') }}
            </button>
            <router-link to="/product/list"
              class="px-6 py-3 bg-slate-100 text-slate-600 rounded-xl text-sm font-medium hover:bg-slate-200 transition">
              {{ $t('common.cancel') }}
            </router-link>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6 h-fit">
        <h3 class="font-semibold text-slate-700 mb-4">{{ $t('product.ai.title') }}</h3>
        <div class="space-y-3">
          <div>
            <label class="block text-xs text-slate-500 mb-1">{{ $t('product.ai.target') }}</label>
            <select v-model="aiForm.biz_type" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm">
              <option value="cover">{{ $t('product.ai.cover') }}</option>
              <option value="gallery">{{ $t('product.ai.carouselImage') }}</option>
              <option value="detail">{{ $t('product.ai.detailImage') }}</option>
              <option value="intro">{{ $t('product.ai.introImage') }}</option>
            </select>
          </div>
          <div>
            <label class="block text-xs text-slate-500 mb-1">{{ $t('product.ai.model') }}</label>
            <select v-model.number="aiForm.model_id" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm">
              <option v-for="m in aiModels" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-xs text-slate-500 mb-1">{{ $t('product.ai.prompt') }}</label>
            <textarea v-model="aiForm.prompt" rows="3" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm resize-none" />
          </div>
          <div>
            <label class="block text-xs text-slate-500 mb-1">{{ $t('product.ai.refImage') }}</label>
            <input
              type="file"
              accept="image/*"
              :disabled="!selectedModelSupportsRef"
              @change="onRefImageChange"
              class="w-full border border-slate-200 rounded-lg px-3 py-2 text-xs"
            />
            <p v-if="!selectedModelSupportsRef" class="text-xs text-orange-600 mt-1">{{ $t('product.ai.refNotSupported') }}</p>
            <p v-if="aiForm.ref_image_url" class="text-xs text-slate-500 mt-1 truncate">{{ $t('product.ai.refUploaded', { url: aiForm.ref_image_url }) }}</p>
          </div>
          <button
            class="w-full bg-blue-700 text-white py-2.5 rounded-lg text-sm hover:bg-blue-600 disabled:opacity-60"
            :disabled="aiGenerating || !aiForm.prompt.trim()"
            @click="generateWithAI"
          >
            {{ aiGenerating ? $t('product.ai.generating') : $t('product.ai.generate') }}
          </button>
          <div v-if="aiImages.length" class="grid grid-cols-2 gap-2 pt-1">
            <div v-for="(url, idx) in aiImages" :key="idx" class="border border-slate-200 rounded-lg p-1.5">
              <img :src="url" class="w-full h-24 object-cover rounded" />
              <button class="w-full mt-1 text-xs bg-slate-100 rounded py-1 hover:bg-slate-200" @click="applyAiImage(url)">
                {{ $t('product.ai.apply') }}
              </button>
            </div>
          </div>
          <p class="text-xs text-slate-400" v-if="aiForm.biz_type === 'detail'">{{ $t('product.ai.detailHint', { index: currentBlockIndex + 1 }) }}</p>
          <p v-if="aiNotice" class="text-xs text-emerald-600">{{ aiNotice }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { createProduct, generateAiImage, getAiModels, getAiTask, getCategories, getProduct, updateProduct, uploadFile } from '@/api/plugins'

type DetailBlock = {
  id: string
  type: 'text' | 'image'
  props: Record<string, any>
}

const { t } = useI18n()
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
  throw new Error(t('product.ai.timeout'))
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
      throw new Error(taskDetail?.error_msg || t('product.ai.failed'))
    }
    if (taskDetail?.result_urls) {
      try {
        aiImages.value = JSON.parse(taskDetail.result_urls)
      } catch {
        aiImages.value = []
      }
    }
  } catch (e: any) {
    error.value = e.message || t('product.ai.failed')
  } finally {
    aiGenerating.value = false
  }
}

function applyAiImage(url: string) {
  if (!url) return
  aiNotice.value = ''
  if (aiForm.value.biz_type === 'cover') {
    form.value.cover = url
    aiNotice.value = t('product.ai.appliedCover')
    return
  }
  if (aiForm.value.biz_type === 'gallery') {
    addGalleryImage(url)
    aiNotice.value = t('product.ai.appliedCarousel')
    return
  }
  if (aiForm.value.biz_type === 'detail') {
    addImageBlock(url, Math.min(currentBlockIndex.value + 1, detailBlocks.value.length))
    aiNotice.value = t('product.ai.appliedDetail', { index: currentBlockIndex.value + 1 })
    return
  }
  if (aiForm.value.biz_type === 'intro') {
    aiNotice.value = t('product.ai.appliedIntro')
  }
}

async function save() {
  if (!form.value.title) { error.value = t('product.form.nameRequired'); return }
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
    error.value = e.message || t('common.saveFailed')
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
