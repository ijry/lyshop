<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { getCheckinLogs } from '@/api/checkin'

const { t } = useI18n()
const loading = ref(false)
const list = ref<any[]>([])

async function loadLogs() {
  loading.value = true
  try {
    const res = await getCheckinLogs({ page: 1, size: 50 })
    list.value = Array.isArray(res?.list) ? res.list : []
  } finally {
    loading.value = false
  }
}

onShow(() => loadLogs())
</script>

<template>
  <view class="page">
    <view class="title-bar">
      <text class="title">{{ t('checkin.logs') }}</text>
    </view>
    <view v-if="!loading && !list.length" class="empty">暂无签到记录</view>
    <view v-for="item in list" :key="item.id" class="card">
      <view class="row">
        <text class="name">用户 {{ item.user_id }}</text>
        <text class="date">{{ item.date }}</text>
      </view>
      <view class="desc">连续 {{ item.consecutive_days }} 天 · 获得 {{ item.points }} 积分</view>
    </view>
    <view v-if="loading" class="loading">加载中…</view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; align-content: start; }
.title-bar { display: flex; align-items: center; }
.title { font-size: 32rpx; font-weight: 700; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.row { display: flex; align-items: center; justify-content: space-between; }
.name { font-size: 28rpx; font-weight: 600; }
.date { font-size: 24rpx; color: var(--eapp-text-muted); }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>
