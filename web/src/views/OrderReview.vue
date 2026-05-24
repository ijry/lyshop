<template>
  <div class="max-w-4xl mx-auto px-6 py-8">
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">← 返回</button>
      <h1 class="text-xl font-bold text-gray-900">订单评价</h1>
    </div>

    <div class="card p-5 mb-4" v-if="meta.order_no">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-500">订单号：{{ meta.order_no }}</p>
          <div class="flex items-center gap-2">
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
    </div>

    <div class="space-y-4">
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
          </div>
        </div>
      </div>
    </div>

    <div class="card p-5 mt-4" v-if="canAppend">
      <p class="text-sm font-medium text-gray-800 mb-3">追加评论</p>
      <textarea v-model="appendContent" class="w-full border border-gray-200 rounded-xl p-3 text-sm min-h-[90px] outline-none focus:border-red-300" placeholder="可选：补充追评内容" />
    </div>

    <div class="mt-6">
      <button class="btn-primary w-full !py-3" :disabled="saving" @click="submitReview">
        {{ saving ? '提交中...' : '提交评价' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get, post } from '@/api/request'

const route = useRoute()
const router = useRouter()

const meta = ref<any>({})
const items = ref<any[]>([])
const logisticsScore = ref(5)
const appendContent = ref('')
const canAppend = ref(false)
const saving = ref(false)

async function loadMeta(id: number) {
  const data = await get<any>(`/api/v1/orders/${id}/review`)
  meta.value = data || {}
  items.value = Array.isArray(data?.options) ? data.options.map((item: any) => ({
    ...item,
    product_score: Number(item.product_score || 5),
    content: String(item.content || ''),
  })) : []
  logisticsScore.value = Number(data?.logistics_score || 5)
  canAppend.value = !!data?.can_append
}

async function submitReview() {
  if (saving.value) return
  saving.value = true
  try {
    const createItems = items.value.filter((item: any) => !item.has_review)
    const editItems = items.value.filter((item: any) => item.has_review)

    if (createItems.length) {
      await post(`/api/v1/orders/${meta.value.order_id}/review`, {
        mode: 'create',
        logistics_score: logisticsScore.value,
        items: createItems.map((item: any) => ({
          order_item_id: item.order_item_id,
          product_score: item.product_score,
          content: item.content,
        })),
        content: createItems[0]?.content || '',
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
        })),
        content: editItems[0]?.content || '',
      })
    }

    if (canAppend.value && appendContent.value.trim()) {
      await post(`/api/v1/orders/${meta.value.order_id}/review`, {
        mode: 'append',
        items: editItems.map((item: any) => ({ order_item_id: item.order_item_id })),
        append_content: appendContent.value.trim(),
        content: appendContent.value.trim(),
      })
    }

    alert('评价提交成功')
    router.replace('/orders')
  } catch (error: any) {
    alert(error?.message || '评价提交失败')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await loadMeta(Number(route.params.id))
})
</script>
