<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">← 返回</button>
      <h2 class="text-xl font-semibold text-slate-800">订单详情</h2>
    </div>

    <div v-if="detail" class="grid grid-cols-1 xl:grid-cols-[2fr_1fr] gap-6">
      <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
        <h3 class="font-semibold text-slate-700 mb-4">商品明细</h3>
        <div v-if="detail.items?.length" class="space-y-4">
          <div v-for="it in detail.items" :key="it.id" class="flex items-center gap-3 p-3 rounded-lg border border-slate-100">
            <img :src="it.cover" class="w-16 h-16 rounded object-cover border border-slate-200" />
            <div class="flex-1 min-w-0">
              <p class="text-sm text-slate-700 truncate">{{ it.title }}</p>
              <p class="text-xs text-slate-400 mt-1">数量 x{{ it.qty }}</p>
            </div>
            <p class="text-sm font-medium text-slate-800">¥{{ money(it.price) }}</p>
          </div>
        </div>
        <p v-else class="text-slate-400 text-sm">暂无商品数据</p>
      </div>

      <div class="space-y-4">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-3">订单信息</h3>
          <div class="space-y-2 text-sm text-slate-600">
            <p>订单号：<span class="font-mono">{{ detail.order_no }}</span></p>
            <p>用户ID：{{ detail.user_id }}</p>
            <p>支付方式：{{ payLabel(detail.payment_method) }}</p>
            <p>状态：{{ statusLabel(detail.status) }}</p>
            <p>下单时间：{{ formatDate(detail.created_at) }}</p>
            <p v-if="detail.paid_at">支付时间：{{ formatDate(detail.paid_at) }}</p>
            <p v-if="detail.tracking_no">物流单号：{{ detail.tracking_no }}</p>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-3">价格体系</h3>
          <div class="space-y-2 text-sm text-slate-600">
            <div class="flex items-center justify-between"><span>商品金额</span><span>¥{{ money(detail.amount_breakdown?.goods_amount ?? detail.goods_amount) }}</span></div>
            <div class="flex items-center justify-between"><span>优惠金额</span><span>-¥{{ money(detail.amount_breakdown?.discount_amount ?? detail.discount_amount) }}</span></div>
            <div class="flex items-center justify-between"><span>运费</span><span>¥{{ money(detail.amount_breakdown?.freight_amount ?? detail.freight_amount) }}</span></div>
            <div class="flex items-center justify-between text-slate-800 font-semibold pt-2 border-t border-slate-100"><span>实付金额</span><span>¥{{ money(detail.amount_breakdown?.payable_amount ?? detail.total_amount) }}</span></div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="bg-white rounded-xl shadow-sm border border-slate-100 p-10 text-center text-slate-400">
      加载中...
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getOrderDetail } from '@/api/plugins'

const route = useRoute()
const router = useRouter()
const detail = ref<any>(null)

const statusLabels: Record<number, string> = {
  1: '待付款', 2: '待发货', 3: '待收货', 4: '已完成', 5: '售后'
}
const statusLabel = (s: number) => statusLabels[s] || '未知'
const payLabel = (m: string) => m === 'wechat' ? '微信支付' : m === 'alipay' ? '支付宝支付' : '未支付'
const money = (v: number) => Number(v || 0).toFixed(2)
const formatDate = (v: string) => v ? v.slice(0, 19).replace('T', ' ') : '-'

onMounted(async () => {
  detail.value = await getOrderDetail(Number(route.params.id))
})
</script>
