<template>
  <div class="max-w-5xl mx-auto px-6 py-8">
    <div class="flex flex-col lg:flex-row gap-6">
      <!-- Sidebar -->
      <div class="w-72 shrink-0">
        <!-- Profile card -->
        <div class="card p-6 mb-4">
          <div class="flex items-center gap-4 mb-4">
            <img :src="user.avatar || 'https://api.dicebear.com/7.x/adventurer/svg?seed=default'"
              class="w-16 h-16 rounded-full border-2 border-red-100" />
            <div>
              <h2 class="text-lg font-bold text-gray-900">{{ user.nickname || '未登录' }}</h2>
              <p class="text-sm text-gray-500">{{ user.phone || '' }}</p>
            </div>
          </div>
          <div class="bg-gradient-to-r from-red-50 to-orange-50 rounded-xl p-4">
            <div class="flex-between">
              <span class="text-sm text-gray-600">我的积分</span>
              <span class="text-2xl font-bold text-red-500">{{ user.points || 0 }}</span>
            </div>
            <p class="text-xs text-gray-400 mt-1">100积分 = ¥1.00</p>
          </div>
        </div>

        <!-- Menu -->
        <div class="card overflow-hidden">
          <button v-for="item in menuItems" :key="item.label"
            @click="item.action"
            :class="activeMenu === item.key ? 'bg-red-50 text-red-600 font-medium border-l-3 border-red-600' : 'text-gray-600 hover:bg-gray-50'"
            class="w-full text-left px-5 py-3.5 text-sm transition-colors flex-between">
            <div class="flex items-center gap-3">
              <div :class="item.icon" class="text-lg" />
              <span>{{ item.label }}</span>
            </div>
            <div class="i-carbon-chevron-right text-xs text-gray-300" />
          </button>
        </div>
      </div>

      <!-- Main content -->
      <div class="flex-1 min-w-0">
        <!-- Order status shortcuts -->
        <div v-if="activeMenu === 'overview'" class="card p-6 mb-4">
          <div class="flex-between mb-4">
            <h3 class="text-base font-bold text-gray-900">我的订单</h3>
            <router-link to="/orders" class="text-sm text-red-600 hover:text-red-700">全部订单 →</router-link>
          </div>
          <div class="grid grid-cols-4 gap-4">
            <router-link v-for="s in orderStatuses" :key="s.label" to="/orders"
              class="flex flex-col items-center gap-2 p-4 rounded-xl hover:bg-gray-50 transition-colors">
              <div :class="s.icon" class="text-2xl text-gray-600" />
              <span class="text-sm text-gray-600">{{ s.label }}</span>
            </router-link>
          </div>
        </div>

        <!-- Coupons -->
        <div v-if="activeMenu === 'overview' || activeMenu === 'coupons'" class="card p-6 mb-4">
          <h3 class="text-base font-bold text-gray-900 mb-4">我的优惠券</h3>
          <div v-if="coupons.length" class="space-y-3">
            <div v-for="c in coupons" :key="c.id"
              class="flex items-center border rounded-xl overflow-hidden"
              :class="c.status === 1 ? 'border-red-200' : 'border-gray-100 opacity-60'">
              <div class="w-24 py-4 text-center shrink-0"
                :class="c.status === 1 ? 'bg-red-500 text-white' : 'bg-gray-200 text-gray-500'">
                <div class="text-xl font-bold">¥{{ c.discount || (c.type === 2 ? (c.discount_rate * 10) + '折' : '?') }}</div>
                <div class="text-xs mt-0.5">{{ c.min_amount > 0 ? `满${c.min_amount}可用` : '无门槛' }}</div>
              </div>
              <div class="flex-1 px-4 py-3">
                <p class="text-sm font-medium text-gray-800">{{ c.name || '优惠券' }}</p>
                <p class="text-xs text-gray-400 mt-1">{{ c.status === 1 ? '未使用' : c.status === 2 ? '已使用' : '已过期' }}</p>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-400 text-sm">暂无优惠券</div>
        </div>

        <!-- Points log -->
        <div v-if="activeMenu === 'overview' || activeMenu === 'points'" class="card p-6 mb-4">
          <h3 class="text-base font-bold text-gray-900 mb-4">积分明细</h3>
          <div v-if="pointsLogs.length" class="divide-y divide-gray-50">
            <div v-for="log in pointsLogs" :key="log.id" class="flex-between py-3">
              <div>
                <p class="text-sm text-gray-800">{{ log.remark }}</p>
                <p class="text-xs text-gray-400 mt-0.5">{{ log.created_at?.slice(0, 10) }}</p>
              </div>
              <span class="text-sm font-bold" :class="log.points > 0 ? 'text-green-600' : 'text-red-500'">
                {{ log.points > 0 ? '+' : '' }}{{ log.points }}
              </span>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-400 text-sm">暂无积分记录</div>
        </div>

        <!-- Addresses -->
        <div v-if="activeMenu === 'overview' || activeMenu === 'addresses'" class="card p-6">
          <div class="flex-between mb-4">
            <h3 class="text-base font-bold text-gray-900">收货地址</h3>
            <button class="btn-primary text-xs !px-4 !py-1.5">新增地址</button>
          </div>
          <div v-if="addresses.length" class="space-y-3">
            <div v-for="addr in addresses" :key="addr.id"
              class="border border-gray-100 rounded-xl p-4 hover:border-gray-200 transition-colors">
              <div class="flex-between mb-2">
                <div class="flex items-center gap-3">
                  <span class="text-sm font-semibold text-gray-800">{{ addr.name }}</span>
                  <span class="text-sm text-gray-500">{{ addr.phone }}</span>
                  <span v-if="addr.is_default"
                    class="text-xs bg-red-50 text-red-600 px-2 py-0.5 rounded">默认</span>
                </div>
              </div>
              <p class="text-sm text-gray-500">{{ addr.province }}{{ addr.city }}{{ addr.district }} {{ addr.detail }}</p>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-400 text-sm">暂无收货地址</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/api/request'

const user = ref<any>({})
const coupons = ref<any[]>([])
const pointsLogs = ref<any[]>([])
const addresses = ref<any[]>([])
const activeMenu = ref('overview')

const menuItems = [
  { key: 'overview', label: '账户总览', icon: 'i-carbon-dashboard', action: () => activeMenu.value = 'overview' },
  { key: 'coupons', label: '我的优惠券', icon: 'i-carbon-ticket', action: () => activeMenu.value = 'coupons' },
  { key: 'points', label: '我的积分', icon: 'i-carbon-star', action: () => activeMenu.value = 'points' },
  { key: 'addresses', label: '收货地址', icon: 'i-carbon-location', action: () => activeMenu.value = 'addresses' },
]

const orderStatuses = [
  { label: '待付款', icon: 'i-carbon-wallet' },
  { label: '待发货', icon: 'i-carbon-package' },
  { label: '待收货', icon: 'i-carbon-delivery-truck' },
  { label: '售后', icon: 'i-carbon-renew' },
]

onMounted(async () => {
  const [profile, couponData, pointsData, addrData] = await Promise.all([
    get<any>('/api/v1/user/profile'),
    get<any>('/api/v1/user/coupons'),
    get<any>('/api/v1/user/points/logs'),
    get<any[]>('/api/v1/addresses'),
  ])
  user.value = profile || {}
  coupons.value = couponData || []
  pointsLogs.value = pointsData?.list || []
  addresses.value = addrData || []
})
</script>
