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
              <label class="block text-sm font-medium text-slate-700">{{ $t('product.form.skuSection') }}</label>
              <button class="text-xs text-blue-600 hover:underline" @click="addSku">{{ $t('product.form.addSku') }}</button>
            </div>
            <div class="space-y-3">
              <div v-for="(sku, skuIdx) in skus" :key="sku.local_key" class="border border-slate-200 rounded-xl p-3">
                <div class="flex items-center justify-between mb-3">
                  <span class="text-xs text-slate-500">{{ $t('product.form.skuLabel', { idx: skuIdx + 1, id: sku.id || '-' }) }}</span>
                  <button class="px-2 py-1 text-xs rounded bg-red-50 text-red-600 hover:bg-red-100" @click="removeSku(skuIdx)">{{ $t('common.delete') }}</button>
                </div>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                  <input
                    v-model="sku.sku_code"
                    class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400"
                    :placeholder="$t('product.form.skuCodePlaceholder')"
                  />
                  <input
                    v-model.number="sku.price"
                    type="number"
                    step="0.01"
                    class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400"
                    :placeholder="$t('product.form.price')"
                  />
                  <input
                    v-model.number="sku.stock"
                    type="number"
                    class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400"
                    :placeholder="$t('product.form.stock')"
                  />
                </div>

                <div class="mt-3">
                  <div class="flex items-center justify-between mb-2">
                    <span class="text-xs text-slate-500">{{ $t('product.form.skuAttrs') }}</span>
                    <button class="text-xs text-blue-600 hover:underline" @click="addSkuAttr(skuIdx)">{{ $t('product.form.addSkuAttr') }}</button>
                  </div>
                  <div class="space-y-2">
                    <div v-for="(attr, attrIdx) in sku.attrs" :key="attr.local_key" class="grid grid-cols-[1fr_1fr_auto] gap-2">
                      <input
                        v-model="attr.name"
                        class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400"
                        :placeholder="$t('product.form.skuAttrName')"
                      />
                      <input
                        v-model="attr.value"
                        class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400"
                        :placeholder="$t('product.form.skuAttrValue')"
                      />
                      <button class="px-3 py-2 text-xs bg-slate-100 rounded-lg hover:bg-slate-200" @click="removeSkuAttr(skuIdx, attrIdx)">{{ $t('common.delete') }}</button>
                    </div>
                  </div>
                </div>

                <div v-if="vipEnabled && vipLevels.length" class="mt-3 rounded-lg border border-amber-200 bg-amber-50/40 p-3">
                  <div class="text-xs font-medium text-amber-700 mb-2">{{ $t('product.form.vipPriceHook') }}</div>
                  <p v-if="!canEditVip" class="text-[11px] text-slate-500 mb-2">{{ $t('product.form.vipPriceReadonly') }}</p>
                  <p v-if="sku.id <= 0" class="text-[11px] text-slate-500 mb-2">{{ $t('product.form.vipPriceSaveHint') }}</p>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <div v-for="level in vipLevels" :key="`lvl-${level.id}`">
                      <label class="block text-xs text-slate-500 mb-1">{{ level.name }}</label>
                      <input
                        :value="getVipPriceValue(sku.id, Number(level.id))"
                        type="number"
                        step="0.01"
                        class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400"
                        :placeholder="$t('product.form.vipPricePlaceholder')"
                        :disabled="sku.id <= 0 || !canEditVip"
                        @input="setVipPriceValue(sku.id, Number(level.id), ($event.target as HTMLInputElement).value)"
                      />
                    </div>
                  </div>
                </div>
              </div>
              <p v-if="vipEnabled && !vipLevels.length" class="text-xs text-slate-400">{{ $t('product.form.vipPriceNoLevel') }}</p>
            </div>
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
                <div v-else-if="block.type === 'rich_text'">
                  <RichTextEditor
                    :model-value="{ content: block.props.content || '' }"
                    @update:modelValue="updateRichTextBlock(idx, $event.content || '')"
                  />
                </div>
              </div>
              <div class="flex gap-2">
                <button class="px-3 py-2 text-xs rounded-lg bg-slate-100 hover:bg-slate-200" @click="addTextBlock()">{{ $t('product.form.addText') }}</button>
                <button class="px-3 py-2 text-xs rounded-lg bg-slate-100 hover:bg-slate-200" @click="addImageBlock()">{{ $t('product.form.addImage') }}</button>
                <button class="px-3 py-2 text-xs rounded-lg bg-slate-100 hover:bg-slate-200" @click="addRichTextBlock()">{{ $t('product.form.addRichText') }}</button>
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
import { createProduct, createVipSkuPrice, deleteVipSkuPrice, generateAiImage, getAiModels, getAiTask, getCategories, getProduct, getVipLevels, getVipSkuPrices, updateProduct, updateVipSkuPrice, uploadFile } from '@/api/plugins'
import { getMenus, type AdminMenuGroupedResponse, type AdminMenuItem, type AdminMenuResponse } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import RichTextEditor from '@/views/decor/editors/RichTextEditor.vue'

type DetailBlock = {
  id: string
  type: 'text' | 'image' | 'rich_text'
  props: Record<string, any>
}

type SkuAttr = {
  local_key: string
  name: string
  value: string
}

type EditableSku = {
  local_key: string
  id: number
  sku_code: string
  price: number
  stock: number
  attrs: SkuAttr[]
}

type VipLevel = {
  id: number
  name: string
}

type VipSkuPrice = {
  id: number
  product_id: number
  sku_id: number
  level_id: number
  vip_price: number
  status: number
}

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const isEdit = computed(() => !!route.params.id)
const canEditVip = computed(() => {
  if (!Array.isArray(auth.perms) || auth.perms.length === 0) return true
  return auth.hasPermission('vip:edit')
})
const categories = ref<any[]>([])
const saving = ref(false)
const error = ref('')
const currentBlockIndex = ref(0)
const galleryImages = ref<Array<{ url: string; sort: number }>>([])
const vipEnabled = ref(false)
const vipLevels = ref<VipLevel[]>([])
const vipPriceRows = ref<VipSkuPrice[]>([])
const vipPriceMap = ref<Record<number, Record<number, string>>>({})

const form = ref({
  title: '', price: 0, origin_price: 0, stock: 0,
  category_id: '', cover: '', status: 1,
})

const detailBlocks = ref<DetailBlock[]>([
  { id: `b-${Date.now()}`, type: 'text', props: { text: '' } }
])

const skus = ref<EditableSku[]>([createEmptySku()])

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

function makeLocalKey(prefix: string) {
  return `${prefix}-${Date.now()}-${Math.random().toString(16).slice(2, 8)}`
}

function createEmptySkuAttr(): SkuAttr {
  return {
    local_key: makeLocalKey('sku-attr'),
    name: '',
    value: '',
  }
}

function createEmptySku(): EditableSku {
  return {
    local_key: makeLocalKey('sku'),
    id: 0,
    sku_code: '',
    price: Number(form.value.price || 0),
    stock: Number(form.value.stock || 0),
    attrs: [createEmptySkuAttr()],
  }
}

function parseSkuAttrs(raw: any): SkuAttr[] {
  let rows: any[] = []
  if (Array.isArray(raw)) {
    rows = raw
  } else if (typeof raw === 'string' && raw.trim()) {
    try {
      const parsed = JSON.parse(raw)
      if (Array.isArray(parsed)) rows = parsed
    } catch {
      rows = []
    }
  }
  const mapped = rows.map((item: any) => ({
    local_key: makeLocalKey('sku-attr'),
    name: String(item?.name || ''),
    value: String(item?.value || ''),
  }))
  return mapped.length ? mapped : [createEmptySkuAttr()]
}

function parseSkus(list: any[]) {
  const rows = Array.isArray(list) ? list : []
  skus.value = rows.map((item: any) => ({
    local_key: makeLocalKey('sku'),
    id: Number(item?.id || 0),
    sku_code: String(item?.sku_code || ''),
    price: Number(item?.price || 0),
    stock: Number(item?.stock || 0),
    attrs: parseSkuAttrs(item?.attrs),
  }))
  if (!skus.value.length) {
    skus.value = [createEmptySku()]
  }
}

function normalizeSkuAttrs(attrs: SkuAttr[]) {
  return attrs
    .map((item) => ({ name: String(item.name || '').trim(), value: String(item.value || '').trim() }))
    .filter((item) => item.name || item.value)
}

function buildSkusPayload() {
  return skus.value
    .map((sku) => {
      const row: any = {
        sku_code: String(sku.sku_code || '').trim(),
        price: Number(sku.price || 0),
        stock: Number(sku.stock || 0),
        attrs: normalizeSkuAttrs(sku.attrs),
      }
      if (sku.id > 0) row.id = sku.id
      return row
    })
    .filter((row: any) => row.sku_code || row.price > 0 || row.stock > 0 || row.attrs.length > 0)
}

function addSku() {
  skus.value.push(createEmptySku())
}

function removeSku(index: number) {
  const target = skus.value[index]
  if (!target) return
  if (target.id > 0) {
    delete vipPriceMap.value[target.id]
  }
  skus.value.splice(index, 1)
  if (!skus.value.length) {
    skus.value = [createEmptySku()]
  }
}

function addSkuAttr(skuIndex: number) {
  const target = skus.value[skuIndex]
  if (!target) return
  target.attrs.push(createEmptySkuAttr())
}

function removeSkuAttr(skuIndex: number, attrIndex: number) {
  const target = skus.value[skuIndex]
  if (!target) return
  target.attrs.splice(attrIndex, 1)
  if (!target.attrs.length) {
    target.attrs = [createEmptySkuAttr()]
  }
}

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
    .filter((item: any) => item && (item.type === 'text' || item.type === 'image' || item.type === 'rich_text'))
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

function addRichTextBlock(position = detailBlocks.value.length) {
  detailBlocks.value.splice(position, 0, {
    id: `b-${Date.now()}-${Math.random().toString(16).slice(2, 6)}`,
    type: 'rich_text',
    props: { content: '' },
  })
  currentBlockIndex.value = position
}

function updateRichTextBlock(index: number, content: string) {
  const block = detailBlocks.value[index]
  if (!block || block.type !== 'rich_text') return
  block.props.content = content
  currentBlockIndex.value = index
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

function isGroupedResponse(data: AdminMenuResponse): data is AdminMenuGroupedResponse {
  return !!data && !Array.isArray(data) && Array.isArray((data as AdminMenuGroupedResponse).groups)
}

function flattenMenus(rows: AdminMenuItem[]): AdminMenuItem[] {
  const list = Array.isArray(rows) ? rows : []
  return list.flatMap((row) => [row, ...flattenMenus(row.children || [])])
}

async function detectVipFeatureEnabled() {
  try {
    const data = await getMenus()
    const menuRows = isGroupedResponse(data)
      ? flattenMenus(data.groups.flatMap((group) => group.menus || []))
      : (Array.isArray(data) ? data : [])
    vipEnabled.value = menuRows.some((item) => String(item.path || '').startsWith('/vip'))
  } catch {
    vipEnabled.value = auth.hasPermission('vip:view') || auth.hasPermission('vip:edit')
  }
}

function setVipPriceMap(rows: VipSkuPrice[]) {
  const next: Record<number, Record<number, string>> = {}
  for (const row of rows) {
    const skuID = Number(row.sku_id || 0)
    const levelID = Number(row.level_id || 0)
    if (skuID <= 0 || levelID <= 0) continue
    if (!next[skuID]) next[skuID] = {}
    next[skuID][levelID] = Number(row.vip_price || 0) > 0 ? String(row.vip_price) : ''
  }
  vipPriceMap.value = next
}

async function loadVipMetaAndPrices(productID: number) {
  if (!vipEnabled.value || productID <= 0) return
  try {
    const levelData: any = await getVipLevels({ page: 1, size: 200 })
    vipLevels.value = Array.isArray(levelData?.list)
      ? levelData.list
      : (Array.isArray(levelData) ? levelData : [])

    const priceData: any = await getVipSkuPrices({ page: 1, size: 1000, product_id: productID, status: 1 })
    vipPriceRows.value = Array.isArray(priceData?.list)
      ? priceData.list
      : (Array.isArray(priceData) ? priceData : [])

    setVipPriceMap(vipPriceRows.value)
  } catch {
    vipLevels.value = []
    vipPriceRows.value = []
    vipPriceMap.value = {}
  }
}

function getVipPriceValue(skuID: number, levelID: number): string {
  if (skuID <= 0 || levelID <= 0) return ''
  return vipPriceMap.value[skuID]?.[levelID] || ''
}

function setVipPriceValue(skuID: number, levelID: number, value: string) {
  if (skuID <= 0 || levelID <= 0) return
  if (!vipPriceMap.value[skuID]) vipPriceMap.value[skuID] = {}
  vipPriceMap.value[skuID][levelID] = String(value || '').trim()
}

function vipRowKey(skuID: number, levelID: number) {
  return `${skuID}_${levelID}`
}

function buildDesiredVipPriceRows(productID: number) {
  const desired: Array<{ product_id: number; sku_id: number; level_id: number; vip_price: number; status: number }> = []
  for (const sku of skus.value) {
    if (Number(sku.id) <= 0) continue
    for (const level of vipLevels.value) {
      const levelID = Number(level.id || 0)
      if (levelID <= 0) continue
      const raw = getVipPriceValue(Number(sku.id), levelID)
      if (!raw) continue
      const price = Number(raw)
      if (!Number.isFinite(price) || price <= 0) {
        throw new Error(t('product.form.vipPriceInvalid', { skuId: sku.id, level: level.name || levelID }))
      }
      desired.push({
        product_id: productID,
        sku_id: Number(sku.id),
        level_id: levelID,
        vip_price: Number(price.toFixed(2)),
        status: 1,
      })
    }
  }
  return desired
}

async function syncVipPrices(productID: number) {
  if (!vipEnabled.value || !isEdit.value || !canEditVip.value || !vipLevels.value.length) return
  const desiredRows = buildDesiredVipPriceRows(productID)
  const desiredMap = new Map<string, { product_id: number; sku_id: number; level_id: number; vip_price: number; status: number }>()
  for (const row of desiredRows) {
    desiredMap.set(vipRowKey(row.sku_id, row.level_id), row)
  }

  const existingMap = new Map<string, VipSkuPrice>()
  for (const row of vipPriceRows.value) {
    existingMap.set(vipRowKey(Number(row.sku_id || 0), Number(row.level_id || 0)), row)
  }

  for (const [key, row] of desiredMap) {
    const existing = existingMap.get(key)
    if (!existing) {
      await createVipSkuPrice(row)
      continue
    }
    if (Number(existing.vip_price || 0) !== row.vip_price || Number(existing.status || 0) !== 1) {
      await updateVipSkuPrice(Number(existing.id), { vip_price: row.vip_price, status: 1 })
    }
  }

  for (const [key, row] of existingMap) {
    if (!desiredMap.has(key)) {
      await deleteVipSkuPrice(Number(row.id))
    }
  }

  await loadVipMetaAndPrices(productID)
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
  const skusPayload = buildSkusPayload()
  const imagesPayload = galleryImages.value
    .filter((item) => item.url.trim())
    .map((item, idx) => ({ url: item.url.trim(), sort: idx }))

  try {
    if (isEdit.value) {
      const productID = Number(route.params.id)
      await updateProduct(productID, { product: payload, skus: skusPayload, images: imagesPayload })
      if (vipEnabled.value && canEditVip.value) {
        await loadVipMetaAndPrices(productID)
        await syncVipPrices(productID)
      }
    } else {
      await createProduct({ product: payload, skus: skusPayload, images: imagesPayload })
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
  await detectVipFeatureEnabled()

  if (isEdit.value) {
    const productID = Number(route.params.id)
    const data: any = await getProduct(productID)
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
    parseSkus(data.skus)
    galleryImages.value = Array.isArray(data.images) ? data.images.map((item: any, idx: number) => ({
      url: item.url || '',
      sort: item.sort ?? idx,
    })) : []
    aiForm.value.target_product_id = productID
    if (vipEnabled.value) {
      await loadVipMetaAndPrices(productID)
    }
  }
})
</script>
