<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { getImSessions } from '@/api/message'

const list = ref<any[]>([])

async function loadData() {
  const data: any = await getImSessions({ page: 1, size: 50 })
  list.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
}

onShow(loadData)
</script>

<template>
  <view class="page">
    <view v-if="!list.length" class="empty">暂无会话</view>
    <view v-for="item in list" :key="item.id" class="card">
      <view class="title">{{ item.user_nickname || item.user_id || '-' }}</view>
      <view class="desc">{{ item.last_message || '-' }}</view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 12rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.title { font-size: 28rpx; font-weight: 600; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
</style>
