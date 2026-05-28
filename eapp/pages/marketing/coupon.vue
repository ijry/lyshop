<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getCoupons } from '@/api/marketing'

const list = ref<any[]>([])

async function loadData() {
  const data: any = await getCoupons({ page: 1, size: 30 })
  list.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
}

onShow(loadData)
</script>

<template>
  <view class="page">
    <view v-if="!list.length" class="empty">暂无数据</view>
    <view v-for="item in list" :key="item.id" class="card">
      <view class="row">
        <text>{{ item.name || '-' }}</text>
        <StatusTag :text="item.status_label || item.status || '-'" :type="item.status" />
      </view>
      <view class="desc">类型：{{ item.type || '-' }} · 面额：{{ item.amount || item.value || 0 }}</view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.row { display: flex; justify-content: space-between; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 23rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
</style>
