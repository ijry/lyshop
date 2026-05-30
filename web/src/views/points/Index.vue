<template>
  <div class="points-mall-page">
    <div class="header">
      <div class="container">
        <div class="balance-card">
          <div class="balance-label">我的积分</div>
          <div class="balance-value">{{ userPoints }}</div>
        </div>
      </div>
    </div>

    <div class="container">
      <div class="tabs">
        <div v-for="tab in tabs" :key="tab.value" :class="['tab-item', { active: activeTab === tab.value }]" @click="activeTab = tab.value">
          {{ tab.label }}
        </div>
      </div>

      <div class="product-grid">
        <div v-for="item in filteredProducts" :key="item.id" class="product-card" @click="goToDetail(item.id)">
          <img v-if="item.cover" :src="item.cover" class="product-cover" />
          <div v-else class="product-cover-placeholder">无图</div>
          <div class="product-info">
            <div class="product-title">{{ item.title }}</div>
            <div class="product-footer">
              <div class="product-price">{{ item.points_price }} 积分</div>
              <div class="product-sold">已兑{{ item.sold_count }}</div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="!filteredProducts.length" class="empty">暂无商品</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import request from '@/api/request'

const router = useRouter()
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
    const data: any = await request.get('/api/v1/points/balance')
    userPoints.value = data.points || 0
  } catch (e) {
    console.error(e)
  }
}

async function loadProducts() {
  try {
    const data: any = await request.get('/api/v1/points/products', { page: 1, size: 100 })
    products.value = data.list || []
  } catch (e) {
    console.error(e)
  }
}

function goToDetail(id: number) {
  router.push(`/points/detail/${id}`)
}

onMounted(() => {
  loadUserPoints()
  loadProducts()
})
</script>

<style scoped>
.points-mall-page { min-height: 100vh; background: #f5f5f5; }
.header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 60px 0; }
.container { max-width: 1200px; margin: 0 auto; padding: 0 20px; }
.balance-card { text-align: center; color: #fff; }
.balance-label { font-size: 16px; opacity: 0.9; margin-bottom: 10px; }
.balance-value { font-size: 48px; font-weight: bold; }
.tabs { display: flex; background: #fff; border-radius: 12px; margin: -30px 0 30px; padding: 0 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.tab-item { flex: 1; text-align: center; padding: 20px 0; font-size: 16px; color: #666; cursor: pointer; position: relative; transition: color 0.3s; }
.tab-item:hover { color: #667eea; }
.tab-item.active { color: #667eea; font-weight: 600; }
.tab-item.active::after { content: ''; position: absolute; bottom: 0; left: 50%; transform: translateX(-50%); width: 30px; height: 3px; background: #667eea; border-radius: 2px; }
.product-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(250px, 1fr)); gap: 20px; margin-bottom: 40px; }
.product-card { background: #fff; border-radius: 12px; overflow: hidden; cursor: pointer; transition: transform 0.3s, box-shadow 0.3s; }
.product-card:hover { transform: translateY(-4px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.product-cover { width: 100%; height: 250px; object-fit: cover; }
.product-cover-placeholder { width: 100%; height: 250px; background: #f0f0f0; display: flex; align-items: center; justify-content: center; color: #999; font-size: 14px; }
.product-info { padding: 16px; }
.product-title { font-size: 16px; color: #333; margin-bottom: 12px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.product-footer { display: flex; justify-content: space-between; align-items: center; }
.product-price { font-size: 20px; color: #667eea; font-weight: bold; }
.product-sold { font-size: 14px; color: #999; }
.empty { text-align: center; padding: 80px 0; color: #999; font-size: 16px; }
</style>
