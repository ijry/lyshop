<template>
  <div>
    <div v-if="showTitle" class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ title || $t('review.title') }}</h2>
    </div>

    <div class="bg-white rounded-xl shadow-sm p-4 mb-4 border border-slate-100 flex gap-3 flex-wrap">
      <input
        v-model="query.keyword"
        :placeholder="$t('review.searchPlaceholder')"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm flex-1 min-w-[220px]"
      />
      <input
        v-if="showProductFilter"
        v-model="query.product_id"
        type="number"
        :placeholder="$t('review.productId')"
        class="border border-slate-200 rounded-lg px-3 py-2 text-sm w-32"
      />
      <div
        v-else-if="effectiveProductID > 0"
        class="px-3 py-2 text-sm rounded-lg bg-slate-100 text-slate-600"
      >
        {{ $t('review.productId') }}：{{ effectiveProductID }}
      </div>
      <button @click="search" class="px-4 py-2 bg-slate-100 rounded-lg text-sm hover:bg-slate-200">{{ $t('common.search') }}</button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('review.reviewId') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('review.orderProduct') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('review.rating') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('review.content') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('review.appendReply') }}</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="rv in list" :key="rv.id" class="align-top hover:bg-slate-50">
            <td class="px-4 py-3 text-xs text-slate-500">{{ rv.id }}</td>
            <td class="px-4 py-3">
              <p class="text-xs text-slate-500">{{ rv.order_no || '-' }}</p>
              <p class="text-sm text-slate-700 mt-1">{{ rv.product?.title || rv.order_item?.title || '-' }}</p>
            </td>
            <td class="px-4 py-3 text-xs text-slate-600">
              <p>{{ $t('review.productScore') }}{{ rv.product_score }}</p>
              <p class="mt-1">{{ $t('review.logistics') }}{{ rv.logistics_score }}</p>
            </td>
            <td class="px-4 py-3 text-xs text-slate-600 max-w-[340px] whitespace-pre-wrap break-words">
              {{ rv.content || $t('review.noTextReview') }}
              <div v-if="rv.images?.length" class="flex flex-wrap gap-2 mt-2">
                <img
                  v-for="(img, idx) in rv.images"
                  :key="img + idx"
                  :src="img"
                  class="w-10 h-10 rounded object-cover border border-slate-100"
                />
              </div>
            </td>
            <td class="px-4 py-3 text-xs text-slate-500">
              <p>{{ $t('review.appendCount', { count: rv.appends?.length || 0 }) }}</p>
              <p class="mt-1" :class="rv.admin_reply ? 'text-emerald-600' : 'text-slate-400'">
                {{ rv.admin_reply ? $t('review.replied') : $t('review.noReply') }}
              </p>
            </td>
            <td class="px-4 py-3">
              <button class="text-blue-600 hover:underline text-xs mr-3" @click="openDetail(rv.id)">{{ $t('review.detail') }}</button>
              <button
                v-if="canReplyReview"
                class="text-emerald-600 hover:underline text-xs"
                @click="openReply(rv.id)"
              >
                {{ $t('review.reply') }}
              </button>
            </td>
          </tr>
          <tr v-if="!list.length">
            <td colspan="6" class="px-4 py-12 text-center text-slate-400">{{ $t('review.noReview') }}</td>
          </tr>
        </tbody>
      </table>
      <div class="px-4 py-3 flex items-center justify-between border-t border-slate-100 text-sm text-slate-500">
        <span>{{ $t('common.totalCount', { total }) }}</span>
        <div class="flex gap-2">
          <button
            :disabled="query.page <= 1"
            @click="prevPage"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40"
          >
            {{ $t('common.prevPage') }}
          </button>
          <button
            :disabled="query.page * query.size >= total"
            @click="nextPage"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40"
          >
            {{ $t('common.nextPage') }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="showDetail" class="fixed inset-0 bg-black/35 flex items-center justify-center z-50" @click.self="closeDetail">
      <div class="bg-white rounded-xl w-[720px] max-w-[92vw] max-h-[88vh] overflow-auto p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-base font-semibold text-slate-800">{{ $t('review.detailTitle', { id: detail?.id }) }}</h3>
          <button class="text-slate-400 hover:text-slate-600" @click="closeDetail">{{ $t('common.close') }}</button>
        </div>
        <div v-if="detail" class="space-y-4 text-sm">
          <div class="grid grid-cols-2 gap-3 text-slate-600">
            <p>{{ $t('review.orderNo') }}{{ detail.order_no || '-' }}</p>
            <p>{{ $t('review.user') }}{{ detail.user_nickname || '-' }}</p>
            <p>{{ $t('review.productRating') }}{{ detail.product_score }}</p>
            <p>{{ $t('review.logisticsRating') }}{{ detail.logistics_score }}</p>
          </div>
          <div class="p-3 rounded-lg bg-slate-50 text-slate-700 whitespace-pre-wrap break-words">
            {{ detail.content || $t('review.noTextReview') }}
          </div>
          <div v-if="detail.images?.length" class="flex flex-wrap gap-2">
            <img
              v-for="(img, idx) in detail.images"
              :key="img + idx"
              :src="img"
              class="w-16 h-16 rounded object-cover border border-slate-100"
            />
          </div>
          <div v-if="detail.appends?.length" class="p-3 rounded-lg bg-slate-50">
            <p class="font-medium text-slate-700 mb-2">{{ $t('review.appendRecords') }}</p>
            <div v-for="ap in detail.appends" :key="ap.id" class="mb-2">
              <p class="text-xs text-slate-600">- {{ ap.content || $t('review.imageOnlyAppend') }}</p>
              <div v-if="ap.images?.length" class="flex flex-wrap gap-2 mt-1">
                <img
                  v-for="(img, idx) in ap.images"
                  :key="img + idx"
                  :src="img"
                  class="w-12 h-12 rounded object-cover border border-slate-100"
                />
              </div>
            </div>
          </div>
          <div class="p-3 rounded-lg" :class="detail.admin_reply ? 'bg-emerald-50 text-emerald-700' : 'bg-amber-50 text-amber-700'">
            {{ detail.admin_reply ? `${t('review.currentReply')}${detail.admin_reply.content}` : $t('review.noReply') }}
          </div>
          <div v-if="!canReplyReview" class="p-3 rounded-lg bg-amber-50 text-amber-700">
            <p class="text-xs font-medium">{{ $t('review.permissionHint') }}</p>
            <p class="text-xs mt-1">{{ $t('review.noPermissionMsg') }}</p>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showReply" class="fixed inset-0 bg-black/35 flex items-center justify-center z-50" @click.self="closeReply">
      <div class="bg-white rounded-xl w-[560px] max-w-[92vw] p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-base font-semibold text-slate-800">{{ $t('review.replyTitle', { id: replyID }) }}</h3>
          <button class="text-slate-400 hover:text-slate-600" @click="closeReply">{{ $t('common.close') }}</button>
        </div>
        <textarea
          v-model="replyContent"
          class="w-full border border-slate-200 rounded-xl p-3 text-sm min-h-[120px] outline-none focus:border-blue-300"
          :placeholder="$t('review.replyPlaceholder')"
        />
        <div class="flex justify-end gap-2 mt-4">
          <button class="px-4 py-2 border rounded-lg text-sm" @click="closeReply">{{ $t('common.cancel') }}</button>
          <button
            class="px-4 py-2 bg-blue-700 text-white rounded-lg text-sm disabled:opacity-50"
            :disabled="replying"
            @click="submitReply"
          >
            {{ replying ? $t('common.submitting') : $t('review.submitReply') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { getReviewDetail, getReviews, replyReview } from '@/api/plugins'
import { useAuthStore } from '@/stores/auth'
import { notify } from '@/utils/notify'

const { t } = useI18n()

const props = withDefaults(defineProps<{
  title?: string
  showTitle?: boolean
  showProductFilter?: boolean
  fixedProductId?: number
}>(), {
  title: '',
  showTitle: true,
  showProductFilter: true,
  fixedProductId: 0,
})

const query = ref({ keyword: '', product_id: '', page: 1, size: 20 })
const list = ref<any[]>([])
const total = ref(0)
const showDetail = ref(false)
const detail = ref<any>(null)
const showReply = ref(false)
const replyID = ref(0)
const replyContent = ref('')
const replying = ref(false)
const auth = useAuthStore()
const canReplyReview = computed(() => auth.hasPermission('order:review-reply'))

function notifyReplyPermissionDenied() {
  notify(t('review.noPermission'))
}

const effectiveProductID = computed(() => {
  if (props.fixedProductId > 0) return Number(props.fixedProductId)
  return Number(query.value.product_id || 0)
})

async function load() {
  const data: any = await getReviews({
    keyword: query.value.keyword || undefined,
    product_id: effectiveProductID.value > 0 ? effectiveProductID.value : undefined,
    page: query.value.page,
    size: query.value.size,
  })
  list.value = data?.list || []
  total.value = Number(data?.total || 0)
}

function search() {
  query.value.page = 1
  load()
}

function prevPage() {
  if (query.value.page <= 1) return
  query.value.page -= 1
  load()
}

function nextPage() {
  if (query.value.page * query.value.size >= total.value) return
  query.value.page += 1
  load()
}

async function openDetail(id: number) {
  detail.value = await getReviewDetail(id)
  showDetail.value = true
}

function closeDetail() {
  showDetail.value = false
  detail.value = null
}

function openReply(id: number) {
  if (!canReplyReview.value) {
    notifyReplyPermissionDenied()
    return
  }
  replyID.value = id
  const target = list.value.find((item: any) => Number(item.id) === Number(id))
  replyContent.value = String(target?.admin_reply?.content || '')
  showReply.value = true
}

function closeReply() {
  showReply.value = false
  replyID.value = 0
  replyContent.value = ''
}

async function submitReply() {
  if (!canReplyReview.value) {
    notifyReplyPermissionDenied()
    return
  }
  const content = replyContent.value.trim()
  if (!content) {
    notify(t('review.replyPlaceholder'))
    return
  }
  replying.value = true
  try {
    await replyReview(replyID.value, content)
    closeReply()
    await load()
  } finally {
    replying.value = false
  }
}

watch(
  () => props.fixedProductId,
  (id) => {
    if (id > 0) {
      query.value.product_id = String(id)
    }
    query.value.page = 1
    load()
  },
  { immediate: true }
)
</script>
