<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">优惠券管理</h2>
      <button @click="showForm = true"
        class="px-4 py-2 bg-blue-700 text-white text-sm rounded-xl hover:bg-blue-600 transition">
        + 创建优惠券
      </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">名称</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">类型</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">面值</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">门槛</th>
            <th class="px-4 py-3 text-left text-slate-500 font-medium">状态</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
          <tr v-for="c in coupons" :key="c.id" class="hover:bg-slate-50">
            <td class="px-4 py-3 font-medium text-slate-800">{{ c.name }}</td>
            <td class="px-4 py-3 text-slate-600">{{ typeLabel(c.type) }}</td>
            <td class="px-4 py-3 text-blue-700 font-medium">
              {{ c.type === 2 ? c.discount * 10 + '折' : '¥' + c.discount }}
            </td>
            <td class="px-4 py-3 text-slate-500">{{ c.min_amount > 0 ? '满¥' + c.min_amount : '无门槛' }}</td>
            <td class="px-4 py-3">
              <span :class="c.status===1 ? 'bg-green-50 text-green-600' : 'bg-slate-100 text-slate-400'"
                class="px-2 py-1 rounded-full text-xs">{{ c.status === 1 ? '启用' : '停用' }}</span>
            </td>
          </tr>
          <tr v-if="!coupons.length">
            <td colspan="5" class="px-4 py-12 text-center text-slate-400">暂无优惠券</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create coupon drawer -->
    <div v-if="showForm" class="fixed inset-0 bg-black/40 flex items-start justify-end z-50">
      <div class="bg-white w-96 h-full shadow-2xl p-6 overflow-y-auto">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-lg font-semibold text-slate-800">创建优惠券</h3>
          <button @click="showForm = false" class="text-slate-400 hover:text-slate-600">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">名称</label>
            <input v-model="form.name" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">类型</label>
            <select v-model.number="form.type" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm">
              <option :value="1">满减券</option>
              <option :value="2">折扣券</option>
              <option :value="3">无门槛</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">{{ form.type===2 ? '折扣率 (0-1)' : '减免金额 (¥)' }}</label>
            <input v-model.number="form.discount" type="number" step="0.01"
              class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">使用门槛 (¥，0=无门槛)</label>
            <input v-model.number="form.min_amount" type="number"
              class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">每人限领次数</label>
            <input v-model.number="form.per_limit" type="number"
              class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" />
          </div>
          <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
          <button @click="create" :disabled="saving"
            class="w-full bg-blue-700 text-white py-3 rounded-xl text-sm font-medium hover:bg-blue-600 transition">
            {{ saving ? '创建中...' : '创建' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'

const coupons = ref<any[]>([])
const showForm = ref(false)
const saving = ref(false)
const error = ref('')
const form = ref({ name: '', type: 1, discount: 0, min_amount: 0, per_limit: 1, status: 1 })

const typeLabels: Record<number, string> = { 1: '满减', 2: '折扣', 3: '无门槛' }
const typeLabel = (t: number) => typeLabels[t] || ''

async function loadCoupons() {
  const data: any = await request.get('/marketing/coupons')
  coupons.value = data.list || []
}

async function create() {
  if (!form.value.name) { error.value = '请输入名称'; return }
  saving.value = true; error.value = ''
  try {
    await request.post('/marketing/coupons', form.value)
    showForm.value = false
    form.value = { name: '', type: 1, discount: 0, min_amount: 0, per_limit: 1, status: 1 }
    loadCoupons()
  } catch (e: any) {
    error.value = e.message
  } finally { saving.value = false }
}

onMounted(loadCoupons)
</script>
