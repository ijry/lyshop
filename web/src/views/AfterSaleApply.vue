<template>
  <div class="max-w-4xl mx-auto px-6 py-8" v-if="detail">
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">← {{ $t('orderDetail.back') }}</button>
      <h1 class="text-xl font-bold text-gray-900">{{ $t('afterSaleApply.title') }}</h1>
    </div>

    <div class="card p-5 mb-4">
      <p class="text-sm text-gray-600">{{ $t('afterSaleApply.orderNo') }}{{ detail.order_no }}</p>
      <p class="text-xs text-gray-400 mt-1">{{ $t('afterSaleApply.selectItemsTip') }}</p>
    </div>

    <div class="space-y-4 mb-4">
      <div v-for="item in formItems" :key="item.order_item_id" class="card p-5">
        <div class="flex items-start gap-3">
          <img :src="item.cover" class="w-14 h-14 rounded-lg object-cover" />
          <div class="flex-1 min-w-0">
            <p class="text-sm text-gray-800 truncate">{{ item.title }}</p>
            <p class="text-xs text-gray-400 mt-1">{{ $t('afterSaleApply.refundableQty') }}{{ item.max_qty }}</p>
            <div class="flex items-center gap-3 mt-3">
              <label class="text-xs text-gray-500">{{ $t('afterSaleApply.applyQty') }}</label>
              <input
                type="number"
                min="0"
                :max="item.max_qty"
                v-model.number="item.qty"
                class="w-24 border border-gray-200 rounded-lg px-2 py-1 text-sm"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="card p-5 mb-4">
      <div class="mb-3">
        <p class="text-sm text-gray-700 mb-2">{{ $t('afterSaleApply.type') }}</p>
        <div class="flex gap-3">
          <button
            class="px-4 py-2 rounded-lg text-sm border"
            :class="form.case_type === 'return' ? 'border-red-300 bg-red-50 text-red-600' : 'border-gray-200 text-gray-600'"
            @click="form.case_type = 'return'"
          >
            {{ $t('afterSaleApply.typeRefund') }}
          </button>
          <button
            class="px-4 py-2 rounded-lg text-sm border"
            :class="form.case_type === 'exchange' ? 'border-red-300 bg-red-50 text-red-600' : 'border-gray-200 text-gray-600'"
            @click="form.case_type = 'exchange'"
          >
            {{ $t('afterSaleApply.typeExchange') }}
          </button>
        </div>
      </div>

      <div class="mb-3">
        <p class="text-sm text-gray-700 mb-2">{{ $t('afterSaleApply.reason') }}</p>
        <input v-model="form.reason" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm" :placeholder="$t('afterSaleApply.reasonPlaceholder')" />
      </div>

      <div class="mb-3">
        <p class="text-sm text-gray-700 mb-2">{{ $t('afterSaleApply.description') }}</p>
        <textarea v-model="form.apply_content" class="w-full min-h-[88px] border border-gray-200 rounded-lg px-3 py-2 text-sm" :placeholder="$t('afterSaleApply.descriptionPlaceholder')" />
      </div>

      <div>
        <p class="text-sm text-gray-700 mb-2">{{ $t('afterSaleApply.images') }}</p>
        <div class="flex flex-wrap gap-2">
          <div v-for="(img, idx) in form.apply_images" :key="img + idx" class="relative w-20 h-20">
            <img :src="img" class="w-20 h-20 rounded-lg object-cover border border-gray-100 cursor-pointer" @click="previewImage(img)" />
            <button class="absolute -top-1 -right-1 w-5 h-5 bg-black/65 text-white rounded-full text-xs leading-none" @click="removeImage(idx)">×</button>
          </div>
          <button
            v-if="form.apply_images.length < 9"
            class="w-20 h-20 rounded-lg border border-dashed border-gray-300 text-xs text-gray-500 hover:border-gray-400"
            @click="pickImages"
          >
            {{ $t('afterSaleApply.addImage') }}
          </button>
        </div>
      </div>
    </div>

    <button class="btn-primary w-full !py-3" :disabled="submitting" @click="submit">
      {{ submitting ? $t('afterSaleApply.submitting') : $t('afterSaleApply.submit') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { get, post, upload } from '@/api/request'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const detail = ref<any>(null)
const submitting = ref(false)
const formItems = ref<any[]>([])
const form = ref<any>({
  case_type: 'return',
  reason: '',
  apply_content: '',
  apply_images: [],
})

function alertMsg(msg: string) {
  alert(msg)
}

function pickFiles(maxCount: number): Promise<File[]> {
  return new Promise((resolve) => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = 'image/*'
    input.multiple = maxCount > 1
    input.onchange = () => resolve(Array.from(input.files || []).slice(0, maxCount) as File[])
    input.click()
  })
}

async function pickImages() {
  const remain = Math.max(0, 9 - form.value.apply_images.length)
  if (!remain) return
  const files = await pickFiles(remain)
  for (const file of files) {
    const result: any = await upload('/api/v1/upload', file)
    if (result?.url) form.value.apply_images.push(String(result.url))
  }
}

function removeImage(index: number) {
  form.value.apply_images.splice(index, 1)
}

function previewImage(url: string) {
  if (!url) return
  window.open(url, '_blank')
}

function normalizeQty(raw: any, maxQty: number) {
  const num = Number(raw || 0)
  if (Number.isNaN(num) || num < 0) return 0
  return Math.min(Math.floor(num), maxQty)
}

async function submit() {
  if (submitting.value) return
  const reason = String(form.value.reason || '').trim()
  if (!reason) {
    alertMsg(t('afterSaleApply.reasonRequired'))
    return
  }
  const items = formItems.value
    .map((item: any) => ({
      order_item_id: Number(item.order_item_id),
      qty: normalizeQty(item.qty, Number(item.max_qty || 0)),
    }))
    .filter((item: any) => item.qty > 0)
  if (!items.length) {
    alertMsg(t('afterSaleApply.itemRequired'))
    return
  }
  submitting.value = true
  try {
    const result: any = await post(`/api/v1/orders/${route.params.id}/after-sales`, {
      case_type: form.value.case_type,
      reason,
      apply_content: String(form.value.apply_content || ''),
      apply_images: form.value.apply_images.slice(),
      items,
    })
    const caseID = Number(result?.id || 0)
    alertMsg(t('afterSaleApply.success'))
    if (caseID > 0) {
      router.replace(`/after-sales/${caseID}`)
    } else {
      router.replace(`/orders/${route.params.id}`)
    }
  } catch (error: any) {
    alertMsg(error?.message || t('afterSaleApply.failed'))
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  detail.value = await get<any>(`/api/v1/orders/${route.params.id}`)
  formItems.value = (detail.value?.items || []).map((item: any) => ({
    order_item_id: Number(item.id),
    title: String(item.title || ''),
    cover: String(item.cover || ''),
    max_qty: Math.max(1, Number(item.qty || 1)),
    qty: 0,
  }))
})
</script>
