<template>
  <div class="max-w-4xl mx-auto px-6 py-8">
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">← 返回</button>
      <h1 class="text-xl font-bold text-gray-900">订单评价</h1>
    </div>

    <div class="card p-5 mb-4" v-if="meta.order_no">
      <div class="flex items-center justify-between">
        <p class="text-sm text-gray-500">订单号：{{ meta.order_no }}</p>
        <div class="flex items-center gap-2" v-if="mode === 'root'">
          <span class="text-sm text-gray-500">物流评分</span>
          <div class="flex items-center gap-1">
            <button
              v-for="n in 5"
              :key="n"
              class="text-2xl leading-none transition"
              :class="n <= logisticsScore ? 'text-red-500' : 'text-gray-300 hover:text-red-300'"
              @click="logisticsScore = n"
            >
              ★
            </button>
          </div>
        </div>
      </div>
      <p class="text-xs text-gray-400 mt-2" v-if="mode === 'root'">根评价提交后，才能继续追加评论。</p>
      <p class="text-xs text-gray-400 mt-2" v-else>你正在追加已发布评价。</p>
    </div>

    <div v-if="mode === 'root'" class="space-y-4">
      <div class="card p-5" v-for="item in items" :key="item.order_item_id">
        <div class="flex items-start gap-4">
          <img :src="item.product_cover" class="w-16 h-16 rounded-lg object-cover border border-gray-100" />
          <div class="flex-1 min-w-0">
            <p class="text-sm text-gray-800 font-medium line-clamp-1">{{ item.product_title }}</p>
            <p class="text-xs text-gray-400 mt-1">订单商品ID：{{ item.order_item_id }}</p>
            <div class="flex items-center gap-2 mt-3">
              <span class="text-sm text-gray-500">商品评分</span>
              <div class="flex items-center gap-1">
                <button
                  v-for="n in 5"
                  :key="n"
                  class="text-2xl leading-none transition"
                  :class="n <= item.product_score ? 'text-red-500' : 'text-gray-300 hover:text-red-300'"
                  @click="item.product_score = n"
                >
                  ★
                </button>
              </div>
            </div>
            <textarea v-model="item.content" class="w-full mt-3 border border-gray-200 rounded-xl p-3 text-sm min-h-[88px] outline-none focus:border-red-300" placeholder="写点使用体验..." />
            <div class="mt-3">
              <div class="flex items-center justify-between text-xs text-gray-400 mb-2">
                <span>图片</span>
                <span>{{ item.images.length }}/9</span>
              </div>
              <div class="flex flex-wrap gap-2">
                <div v-for="(img, idx) in item.images" :key="img + idx" class="relative w-20 h-20">
                  <img :src="img" class="w-20 h-20 rounded-lg object-cover border border-gray-100 cursor-pointer" @click="previewImage(img)" />
                  <button class="absolute -top-1 -right-1 w-5 h-5 bg-black/65 text-white rounded-full text-xs leading-none" @click="removeImage(item.images, idx)">×</button>
                </div>
                <button
                  v-if="item.images.length < 9"
                  class="w-20 h-20 rounded-lg border border-dashed border-gray-300 text-xs text-gray-500 hover:border-gray-400"
                  @click="pickImages(item.images)"
                >
                  + 添加
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="card p-5">
      <div v-if="itemForAppend" class="space-y-4">
        <div class="flex items-start gap-4">
          <img :src="itemForAppend.product_cover" class="w-16 h-16 rounded-lg object-cover border border-gray-100" />
          <div class="flex-1 min-w-0">
            <p class="text-sm text-gray-800 font-medium line-clamp-1">{{ itemForAppend.product_title }}</p>
            <p class="text-xs text-gray-400 mt-1">订单商品ID：{{ itemForAppend.order_item_id }}</p>
            <p class="text-xs text-gray-500 mt-2">{{ itemForAppend.content || '未填写评价' }}</p>
            <div v-if="itemForAppend.images?.length" class="flex flex-wrap gap-2 mt-2">
              <img v-for="(img, idx) in itemForAppend.images" :key="img + idx" :src="img" class="w-14 h-14 rounded-md object-cover border border-gray-100" />
            </div>
          </div>
        </div>

        <textarea v-model="appendContent" class="w-full border border-gray-200 rounded-xl p-3 text-sm min-h-[90px] outline-none focus:border-red-300" placeholder="可选：补充追评内容" />
        <div>
          <div class="flex items-center justify-between text-xs text-gray-400 mb-2">
            <span>追评图片</span>
            <span>{{ appendImages.length }}/9</span>
          </div>
          <div class="flex flex-wrap gap-2">
            <div v-for="(img, idx) in appendImages" :key="img + idx" class="relative w-20 h-20">
              <img :src="img" class="w-20 h-20 rounded-lg object-cover border border-gray-100 cursor-pointer" @click="previewImage(img)" />
              <button class="absolute -top-1 -right-1 w-5 h-5 bg-black/65 text-white rounded-full text-xs leading-none" @click="removeImage(appendImages, idx)">×</button>
            </div>
            <button
              v-if="appendImages.length < 9"
              class="w-20 h-20 rounded-lg border border-dashed border-gray-300 text-xs text-gray-500 hover:border-gray-400"
              @click="pickImages(appendImages)"
            >
              + 添加
            </button>
          </div>
        </div>
      </div>
      <p v-else class="text-sm text-gray-400">未找到可追加的评价</p>
    </div>

    <div class="mt-6 space-y-3">
      <button v-if="mode === 'root'" class="btn-primary w-full !py-3" :disabled="savingRoot" @click="submitRootReview">
        {{ savingRoot ? '提交中...' : rootButtonText }}
      </button>
      <button v-else class="btn-primary w-full !py-3" :disabled="savingAppend" @click="submitAppendReview">
        {{ savingAppend ? '提交中...' : '提交追加评价' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get, post, upload } from '@/api/request'

const route = useRoute()
const router = useRouter()

const meta = ref<any>({})
const items = ref<any[]>([])
const itemForAppend = ref<any>(null)
const logisticsScore = ref(5)
const appendContent = ref('')
const appendImages = ref<string[]>([])
const savingRoot = ref(false)
const savingAppend = ref(false)
const mode = ref<'root' | 'append'>('root')

const rootButtonText = computed(() => {
  return items.value.some((item: any) => !item.has_review) ? '提交根评价' : '更新根评价'
})

function notify(msg: string) {
  alert(msg)
}

function normalizeItem(item: any) {
  return {
    ...item,
    product_score: Number(item.product_score || 5),
    content: String(item.content || ''),
    images: Array.isArray(item.images) ? item.images.map((img: any) => String(img || '')) : [],
  }
}

async function loadMeta(id: number) {
  const data = await get<any>(`/api/v1/orders/${id}/review`)
  meta.value = data || {}
  const options = Array.isArray(data?.options) ? data.options.map(normalizeItem) : []
  logisticsScore.value = Number(data?.logistics_score || 5)
  if (mode.value === 'append') {
    itemForAppend.value = options.find((item: any) => Number(item.order_item_id || 0) === Number(route.query.item_id || 0) && !!item.has_review) || null
    return
  }
  items.value = options
}

function pickFiles(maxCount: number): Promise<File[]> {
  return new Promise((resolve) => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = 'image/*'
    input.multiple = maxCount > 1
    input.onchange = () => {
      resolve(Array.from(input.files || []).slice(0, maxCount) as File[])
    }
    input.click()
  })
}

async function pickImages(target: string[]) {
  const remain = Math.max(0, 9 - target.length)
  if (!remain) return
  const files = await pickFiles(remain)
  for (const file of files) {
    const result: any = await upload('/api/v1/upload', file)
    if (result?.url) target.push(String(result.url))
  }
}

function removeImage(target: string[], index: number) {
  target.splice(index, 1)
}

function previewImage(url: string) {
  window.open(url, '_blank')
}

async function submitRootReview() {
  if (savingRoot.value) return
  const createItems = items.value.filter((item: any) => !item.has_review)
  const editItems = items.value.filter((item: any) => item.has_review)
  if (!createItems.length && !editItems.length) {
    notify('暂无可提交的根评价')
    return
  }
  savingRoot.value = true
  try {
    if (createItems.length) {
      await post(`/api/v1/orders/${meta.value.order_id}/review`, {
        mode: 'create',
        logistics_score: logisticsScore.value,
        items: createItems.map((item: any) => ({
          order_item_id: item.order_item_id,
          product_score: item.product_score,
          content: item.content,
          images: item.images || [],
        })),
      })
    }
    if (editItems.length) {
      await post(`/api/v1/orders/${meta.value.order_id}/review`, {
        mode: 'edit',
        logistics_score: logisticsScore.value,
        items: editItems.map((item: any) => ({
          order_item_id: item.order_item_id,
          product_score: item.product_score,
          content: item.content,
          images: item.images || [],
        })),
      })
    }
    notify('根评价已提交')
    router.back()
  } catch (error: any) {
    notify(error?.message || '提交失败')
  } finally {
    savingRoot.value = false
  }
}

async function submitAppendReview() {
  if (savingAppend.value) return
  if (!itemForAppend.value) {
    notify('请先完成根评价')
    return
  }
  const content = appendContent.value.trim()
  if (!content && appendImages.value.length === 0) {
    notify('请填写追评内容或上传图片')
    return
  }
  savingAppend.value = true
  try {
    await post(`/api/v1/orders/${meta.value.order_id}/review`, {
      mode: 'append',
      items: [{ order_item_id: itemForAppend.value.order_item_id }],
      append_content: content,
      append_images: appendImages.value.slice(),
    })
    notify('追评已提交')
    router.back()
  } catch (error: any) {
    notify(error?.message || '提交失败')
  } finally {
    savingAppend.value = false
  }
}

onMounted(async () => {
  mode.value = String(route.query.mode || 'root') === 'append' ? 'append' : 'root'
  await loadMeta(Number(route.params.id))
})
</script>
