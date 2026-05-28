<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getAfterSales } from '@/api/order'

const list = ref<any[]>([])
const loading = ref(false)

async function loadData() {
  loading.value = true
  try {
    const data: any = await getAfterSales({ page: 1, size: 20 })
    list.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/order/after-sale-detail?id=${id}` })
}

onShow(loadData)
</script>

<template>
  <view class="page">
    <view v-if="loading" class="empty">加载中...</view>
    <view v-else-if="!list.length" class="empty">暂无数据</view>
    <view v-for="item in list" :key="item.id" class="card" @click="goDetail(item.id)">
      <view class="row">
        <text>售后 #{{ item.id }}</text>
        <StatusTag :text="item.status_label || item.status || '-'" :type="item.status" />
      </view>
      <view class="desc">订单 #{{ item.order_id }}</view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 16rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; }
.row { display: flex; justify-content: space-between; align-items: center; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
</style>
