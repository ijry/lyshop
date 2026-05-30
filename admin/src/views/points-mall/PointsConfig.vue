<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">积分配置</h2>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">
      <div class="space-y-5">
        <div>
          <label class="text-sm text-slate-700 font-medium mb-2 block">积分兑换比例</label>
          <div class="flex items-center gap-3">
            <input v-model.number="config.points_to_yuan" type="number" class="border border-slate-200 rounded-xl px-4 py-2.5 text-sm w-32" />
            <span class="text-slate-600 text-sm">积分 = 1 元</span>
          </div>
          <div class="text-xs text-slate-500 mt-1">用于订单结算时积分抵扣金额的计算</div>
        </div>

        <div class="border-t pt-5">
          <label class="flex items-center gap-2 mb-3">
            <input v-model="config.enable_order_points" type="checkbox" class="rounded" />
            <span class="text-sm text-slate-700 font-medium">开启订单完成赠送积分</span>
          </label>
          <div v-if="config.enable_order_points" class="ml-6">
            <div class="flex items-center gap-3">
              <span class="text-slate-600 text-sm">消费</span>
              <input v-model.number="config.order_points_rate" type="number" step="0.01" class="border border-slate-200 rounded-xl px-4 py-2.5 text-sm w-24" />
              <span class="text-slate-600 text-sm">% 赠送积分</span>
            </div>
            <div class="text-xs text-slate-500 mt-1">例如：设置为 1，则消费 100 元赠送 100 积分</div>
          </div>
        </div>

        <div class="border-t pt-5">
          <label class="text-sm text-slate-700 font-medium mb-2 block">积分过期天数</label>
          <div class="flex items-center gap-3">
            <input v-model.number="config.points_expire_days" type="number" class="border border-slate-200 rounded-xl px-4 py-2.5 text-sm w-32" />
            <span class="text-slate-600 text-sm">天（0 = 永不过期）</span>
          </div>
          <div class="text-xs text-slate-500 mt-1">积分获得后超过此天数将自动过期扣除</div>
        </div>

        <div class="border-t pt-5">
          <button @click="save" class="px-6 py-2.5 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600">保存配置</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import request from '@/api/request'

const config = ref({
  points_to_yuan: 100,
  order_points_rate: 1,
  points_expire_days: 0,
  enable_order_points: false,
})

async function load() {
  const data: any = await request.get('/points/config')
  config.value = { ...config.value, ...data }
}

async function save() {
  await request.put('/points/config', config.value)
  alert('保存成功')
}

onMounted(load)
</script>
