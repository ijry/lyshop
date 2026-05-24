<template>
  <div class="max-w-5xl mx-auto px-6 py-8" v-if="detail">
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">← 返回</button>
      <h1 class="text-xl font-bold text-gray-900">订单详情</h1>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-[2fr_1fr] gap-6">
      <div class="card p-5">
        <h3 class="font-semibold text-gray-800 mb-4">商品明细</h3>
        <div class="space-y-4">
          <div v-for="it in detail.items || []" :key="it.id" class="flex items-center gap-3">
            <img :src="it.cover" class="w-14 h-14 rounded-lg object-cover" />
            <div class="flex-1 min-w-0">
              <p class="text-sm text-gray-700 truncate">{{ it.title }}</p>
              <p class="text-xs text-gray-400">数量 x{{ it.qty }}</p>
            </div>
            <p class="text-sm font-medium text-gray-800">¥{{ money(it.price) }}</p>
          </div>
        </div>
      </div>

      <div class="space-y-4">
        <div class="card p-5">
          <h3 class="font-semibold text-gray-800 mb-3">订单信息</h3>
          <div class="space-y-2 text-sm text-gray-600">
            <p>订单号：{{ detail.order_no }}</p>
            <p>支付方式：{{ payLabel(detail.payment_method) }}</p>
            <p>状态：{{ statusLabel(detail.status) }}</p>
            <p>下单时间：{{ formatDate(detail.created_at) }}</p>
            <p v-if="detail.paid_at">支付时间：{{ formatDate(detail.paid_at) }}</p>
            <p v-if="detail.tracking_no">物流单号：{{ detail.tracking_no }}</p>
          </div>
        </div>
        <div class="card p-5">
          <h3 class="font-semibold text-gray-800 mb-3">价格体系</h3>
          <div class="space-y-2 text-sm text-gray-600">
            <div class="flex justify-between"><span>商品金额</span><span>¥{{ money(detail.amount_breakdown?.goods_amount ?? detail.goods_amount) }}</span></div>
            <div class="flex justify-between"><span>优惠金额</span><span>-¥{{ money(detail.amount_breakdown?.discount_amount ?? detail.discount_amount) }}</span></div>
            <div class="flex justify-between"><span>运费</span><span>¥{{ money(detail.amount_breakdown?.freight_amount ?? detail.freight_amount) }}</span></div>
            <div class="flex justify-between font-semibold text-gray-800 pt-2 border-t border-gray-100"><span>实付金额</span><span>¥{{ money(detail.amount_breakdown?.payable_amount ?? detail.total_amount) }}</span></div>
          </div>
        </div>

        <div class="card p-5">
          <h3 class="font-semibold text-gray-800 mb-3">评价信息</h3>
          <div v-if="reviewItems.length" class="space-y-3">
            <div v-for="rv in reviewItems" :key="rv.id" class="border border-gray-100 rounded-lg p-3">
              <p class="text-sm text-gray-700">{{ rv.product_title }}</p>
              <p class="text-xs text-gray-400 mt-1">商品 {{ rv.product_score }} / 物流 {{ rv.logistics_score }}</p>
              <p class="text-xs text-gray-500 mt-2">{{ rv.content || '未填写评价' }}</p>
              <div v-if="rv.images?.length" class="flex flex-wrap gap-2 mt-2">
                <img v-for="(img, idx) in rv.images" :key="img + idx" :src="img" class="w-14 h-14 rounded-md object-cover border border-gray-100" />
              </div>
              <div v-if="rv.appends?.length" class="mt-2 p-2 rounded bg-gray-50 space-y-2">
                <div v-for="ap in rv.appends" :key="ap.id">
                  <p class="text-xs text-gray-500">追加：{{ ap.content || '仅图片追评' }}</p>
                  <div v-if="ap.images?.length" class="flex flex-wrap gap-2 mt-1">
                    <img v-for="(img, idx) in ap.images" :key="img + idx" :src="img" class="w-12 h-12 rounded object-cover border border-gray-100" />
                  </div>
                </div>
              </div>
              <div class="flex justify-end mt-3">
                <button class="btn-outline !px-4 !py-2 text-xs" @click="goReview('append', rv.order_item_id)">追加评价</button>
              </div>
            </div>
          </div>
          <p v-else class="text-sm text-gray-400">暂无评价</p>
          <div class="flex gap-2 mt-4" v-if="hasUnreviewed">
            <button class="btn-primary !px-4 !py-2 text-xs" @click="goReview('root')">评价</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get } from '@/api/request'

const route = useRoute()
const router = useRouter()
const detail = ref<any>(null)
const reviewItems = ref<any[]>([])
const hasReviewed = ref(false)
const hasUnreviewed = ref(false)

const statusLabels: Record<number, string> = { 1: '待付款', 2: '待发货', 3: '待收货', 4: '已完成', 5: '售后' }
const statusLabel = (s: number) => statusLabels[s] || '未知'
const payLabel = (m: string) => m === 'wechat' ? '微信支付' : m === 'alipay' ? '支付宝支付' : '未支付'
const money = (v: number) => Number(v || 0).toFixed(2)
const formatDate = (v: string) => v ? v.slice(0, 19).replace('T', ' ') : '-'

function refreshReviewFlags() {
  const items = Array.isArray(detail.value?.items) ? detail.value.items : []
  hasReviewed.value = items.some((item: any) => !!item?.review?.id)
  hasUnreviewed.value = items.length > 0 && items.some((item: any) => !item?.review?.id)
}

function goReview(mode: 'root' | 'append', orderItemID?: number) {
  if (!detail.value?.id) return
  const itemQuery = mode === 'append' && Number(orderItemID || 0) > 0 ? `&item_id=${Number(orderItemID)}` : ''
  router.push(`/orders/${detail.value.id}/review?mode=${mode}${itemQuery}`)
}

onMounted(async () => {
  detail.value = await get<any>(`/api/v1/orders/${route.params.id}`)
  reviewItems.value = []
  for (const item of detail.value?.items || []) {
    if (item.review) {
      reviewItems.value.push({
        ...item.review,
        product_title: item.title,
      })
    }
  }
  refreshReviewFlags()
})
</script>
