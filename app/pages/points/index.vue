<template>
  <view class="page">
    <view class="header">
      <view class="balance-card">
        <view class="balance-label">我的积分</view>
        <view class="balance-value">{{ userPoints }}</view>
      </view>
    </view>

    <view class="tabs">
      <view v-for="tab in tabs" :key="tab.value" :class="['tab-item', { active: activeTab === tab.value }]" @click="activeTab = tab.value">
        {{ tab.label }}
      </view>
    </view>

    <view class="product-list">
      <view v-for="item in filteredProducts" :key="item.id" class="product-card" @click="goToDetail(item.id)">
        <image v-if="item.cover" :src="item.cover" class="product-cover" mode="aspectFill" />
        <view v-else class="product-cover-placeholder">无图</view>
        <view class="product-info">
          <view class="product-title">{{ item.title }}</view>
          <view class="product-footer">
            <view class="product-price">{{ item.points_price }} 积分</view>
            <view class="product-sold">已兑{{ item.sold_count }}</view>
          </view>
        </view>
      </view>
    </view>

    <view v-if="!filteredProducts.length" class="empty">暂无商品</view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { get } from '@/utils/request'

const userPoints = ref(0)
const activeTab = ref('')
const products = ref<any[]>([])

const tabs = [
  { label: '全部', value: '' },
  { label: '优惠券', value: 'coupon' },
  { label: '实物', value: 'physical' },
  { label: '虚拟', value: 'virtual' },
]

const filteredProducts = computed(() => {
  if (!activeTab.value) return products.value
  return products.value.filter(p => p.type === activeTab.value)
})

async function loadUserPoints() {
  try {
    const data: any = await get('/api/v1/points/balance')
    userPoints.value = data.points || 0
  } catch (e) {
    console.error(e)
  }
}

async function loadProducts() {
  try {
    const data: any = await get('/api/v1/points/products', { page: 1, size: 100 })
    products.value = data.list || []
  } catch (e) {
    console.error(e)
  }
}

function goToDetail(id: number) {
  uni.navigateTo({ url: `/pages/points/detail?id=${id}` })
}

onMounted(() => {
  loadUserPoints()
  loadProducts()
})
</script>

<style scoped>
.page { min-height: 100vh; background: #f5f5f5; }
.header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 40rpx 30rpx; }
.balance-card { text-align: center; color: #fff; }
.balance-label { font-size: 28rpx; opacity: 0.9; margin-bottom: 10rpx; }
.balance-value { font-size: 60rpx; font-weight: bold; }
.tabs { display: flex; background: #fff; padding: 0 30rpx; border-bottom: 1px solid #eee; }
.tab-item { flex: 1; text-align: center; padding: 30rpx 0; font-size: 28rpx; color: #666; position: relative; }
.tab-item.active { color: #667eea; font-weight: 600; }
.tab-item.active::after { content: ''; position: absolute; bottom: 0; left: 50%; transform: translateX(-50%); width: 40rpx; height: 4rpx; background: #667eea; border-radius: 2rpx; }
.product-list { padding: 20rpx; display: grid; grid-template-columns: repeat(2, 1fr); gap: 20rpx; }
.product-card { background: #fff; border-radius: 16rpx; overflow: hidden; }
.product-cover { width: 100%; height: 300rpx; }
.product-cover-placeholder { width: 100%; height: 300rpx; background: #f0f0f0; display: flex; align-items: center; justify-content: center; color: #999; font-size: 24rpx; }
.product-info { padding: 20rpx; }
.product-title { font-size: 28rpx; color: #333; margin-bottom: 20rpx; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.product-footer { display: flex; justify-content: space-between; align-items: center; }
.product-price { font-size: 32rpx; color: #667eea; font-weight: bold; }
.product-sold { font-size: 24rpx; color: #999; }
.empty { text-align: center; padding: 100rpx 0; color: #999; }
</style>
