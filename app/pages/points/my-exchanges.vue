<template>
  <view class="page">
    <view class="tabs">
      <view v-for="tab in tabs" :key="tab.value" :class="['tab-item', { active: activeTab === tab.value }]" @click="changeTab(tab.value)">
        {{ tab.label }}
      </view>
    </view>

    <view class="list">
      <view v-for="item in list" :key="item.id" class="exchange-card" @click="goToDetail(item.id)">
        <view class="card-header">
          <view class="exchange-no">{{ item.exchange_no }}</view>
          <view :class="['status', `status-${item.status}`]">{{ getStatusText(item.status) }}</view>
        </view>
        <view class="card-body">
          <image v-if="item.product_cover" :src="item.product_cover" class="product-thumb" mode="aspectFill" />
          <view class="product-info">
            <view class="product-title">{{ item.product_title }}</view>
            <view class="product-meta">
              <text>数量: {{ item.qty }}</text>
              <text class="points">{{ item.points_cost }} 积分</text>
            </view>
          </view>
        </view>
        <view class="card-footer">
          <view class="time">{{ formatDate(item.created_at) }}</view>
          <view v-if="item.status === 'shipped'" class="action-btn" @click.stop="confirmReceive(item.id)">确认收货</view>
        </view>
      </view>
    </view>

    <view v-if="!list.length" class="empty">暂无兑换记录</view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get, post } from '@/utils/request'

const activeTab = ref('')
const list = ref<any[]>([])

const tabs = [
  { label: '全部', value: '' },
  { label: '待发货', value: 'pending_ship' },
  { label: '已完成', value: 'completed' },
]

async function loadList() {
  try {
    const params: any = { page: 1, size: 100 }
    if (activeTab.value) params.status = activeTab.value
    const data: any = await get('/api/v1/points/exchanges', params)
    list.value = data.list || []
  } catch (e) {
    console.error(e)
  }
}

function changeTab(value: string) {
  activeTab.value = value
  loadList()
}

function goToDetail(id: number) {
  uni.navigateTo({ url: `/pages/points/exchange-detail?id=${id}` })
}

async function confirmReceive(id: number) {
  uni.showModal({
    title: '确认收货',
    content: '确认已收到商品吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await post(`/api/v1/points/exchanges/${id}/confirm`)
          uni.showToast({ title: '确认成功', icon: 'success' })
          loadList()
        } catch (e: any) {
          uni.showToast({ title: e.message || '操作失败', icon: 'none' })
        }
      }
    }
  })
}

function getStatusText(status: string) {
  const map: any = {
    pending_ship: '待发货',
    shipped: '已发货',
    completed: '已完成',
    canceled: '已取消',
  }
  return map[status] || status
}

function formatDate(date: string) {
  return date ? new Date(date).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }) : '-'
}

onMounted(loadList)
</script>

<style scoped>
.page { min-height: 100vh; background: #f5f5f5; }
.tabs { display: flex; background: #fff; padding: 0 30rpx; border-bottom: 1px solid #eee; }
.tab-item { flex: 1; text-align: center; padding: 30rpx 0; font-size: 28rpx; color: #666; position: relative; }
.tab-item.active { color: #667eea; font-weight: 600; }
.tab-item.active::after { content: ''; position: absolute; bottom: 0; left: 50%; transform: translateX(-50%); width: 40rpx; height: 4rpx; background: #667eea; border-radius: 2rpx; }
.list { padding: 20rpx; }
.exchange-card { background: #fff; border-radius: 16rpx; padding: 24rpx; margin-bottom: 20rpx; }
.card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20rpx; }
.exchange-no { font-size: 24rpx; color: #999; font-family: monospace; }
.status { font-size: 24rpx; padding: 4rpx 12rpx; border-radius: 20rpx; }
.status-pending_ship { background: #fff3e0; color: #f57c00; }
.status-shipped { background: #e3f2fd; color: #1976d2; }
.status-completed { background: #e8f5e9; color: #388e3c; }
.status-canceled { background: #f5f5f5; color: #999; }
.card-body { display: flex; gap: 20rpx; margin-bottom: 20rpx; }
.product-thumb { width: 120rpx; height: 120rpx; border-radius: 12rpx; }
.product-info { flex: 1; }
.product-title { font-size: 28rpx; color: #333; margin-bottom: 10rpx; }
.product-meta { display: flex; justify-content: space-between; font-size: 24rpx; color: #999; }
.points { color: #667eea; font-weight: bold; }
.card-footer { display: flex; justify-content: space-between; align-items: center; padding-top: 20rpx; border-top: 1px solid #f0f0f0; }
.time { font-size: 24rpx; color: #999; }
.action-btn { font-size: 24rpx; color: #667eea; padding: 8rpx 20rpx; border: 1px solid #667eea; border-radius: 20rpx; }
.empty { text-align: center; padding: 100rpx 0; color: #999; }
</style>
