<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <button class="text-slate-400 hover:text-slate-600 text-sm" @click="router.back()">← 返回</button>
      <h2 class="text-xl font-semibold text-slate-800">售后详情</h2>
    </div>

    <div v-if="detail" class="grid grid-cols-1 xl:grid-cols-[2fr_1fr] gap-6">
      <div class="space-y-4">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-4">售后信息</h3>
          <div class="space-y-2 text-sm text-slate-600">
            <p>售后单号：<span class="font-mono">{{ detail.case_no }}</span></p>
            <p>订单ID：{{ detail.order_id }}</p>
            <p>类型：{{ typeLabel(detail.case_type) }}</p>
            <p>状态：{{ detail.status }}</p>
            <p>原因：{{ detail.reason }}</p>
            <p>申请说明：{{ detail.apply_content || '-' }}</p>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-4">状态日志</h3>
          <div v-if="detail.logs?.length" class="space-y-3">
            <div v-for="log in detail.logs" :key="log.id" class="border border-slate-100 rounded-lg p-3">
              <p class="text-sm text-slate-700">{{ log.action }}：{{ log.from_status || '-' }} → {{ log.to_status }}</p>
              <p class="text-xs text-slate-400 mt-1">{{ log.content || '-' }}</p>
            </div>
          </div>
          <p v-else class="text-slate-400 text-sm">暂无日志</p>
        </div>
      </div>

      <div class="space-y-4">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-3">操作</h3>
          <div class="space-y-3">
            <button class="w-full px-4 py-2 rounded-lg bg-blue-700 text-white text-sm" @click="audit(true)">审核通过</button>
            <button class="w-full px-4 py-2 rounded-lg bg-slate-100 text-slate-700 text-sm" @click="audit(false)">审核拒绝</button>
            <button class="w-full px-4 py-2 rounded-lg bg-emerald-600 text-white text-sm" @click="receive">确认收货</button>
            <button class="w-full px-4 py-2 rounded-lg bg-orange-600 text-white text-sm" @click="refund">登记退款</button>
            <button class="w-full px-4 py-2 rounded-lg bg-purple-600 text-white text-sm" @click="complete">完结</button>
            <button class="w-full px-4 py-2 rounded-lg bg-red-600 text-white text-sm" @click="close">关闭</button>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
          <h3 class="font-semibold text-slate-700 mb-3">物流</h3>
          <div v-if="detail.shipments?.length" class="space-y-2">
            <div v-for="ship in detail.shipments" :key="ship.id" class="border border-slate-100 rounded-lg p-3 text-sm">
              <p>{{ ship.direction }} / {{ ship.biz_type }}</p>
              <p class="text-xs text-slate-400 mt-1">{{ ship.company }} · {{ ship.tracking_no }}</p>
            </div>
          </div>
          <p v-else class="text-slate-400 text-sm">暂无物流</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { auditAfterSale, closeAfterSale, completeAfterSale, getAfterSaleDetail, receiveAfterSale, refundAfterSale } from '@/api/plugins'

const route = useRoute()
const router = useRouter()
const detail = ref<any>(null)

const typeLabel = (v: string) => v === 'exchange' ? '换货' : '退货'

async function load() {
  detail.value = await getAfterSaleDetail(Number(route.params.id))
}

async function audit(approve: boolean) {
  await auditAfterSale(Number(route.params.id), { approve, audit_remark: '' })
  await load()
}
async function receive() { await receiveAfterSale(Number(route.params.id)); await load() }
async function refund() { await refundAfterSale(Number(route.params.id), { amount: detail.value?.refund_amount || 0, reason: detail.value?.reason || '' }); await load() }
async function complete() { await completeAfterSale(Number(route.params.id)); await load() }
async function close() { await closeAfterSale(Number(route.params.id), { reason: '人工关闭' }); await load() }

onMounted(load)
</script>
