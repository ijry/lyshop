<template>
  <view style="min-height: 100vh; background: #f5f5f5;">
    <u-navbar title="我的积分" :placeholder="true" />

    <!-- Points card -->
    <view style="margin: 16px; background: linear-gradient(135deg, #dc2626, #ef4444); border-radius: 16px; padding: 24px; color: #fff;">
      <text style="font-size: 13px; opacity: 0.8;">当前积分</text>
      <text style="font-size: 40px; font-weight: 800; display: block; margin-top: 4px;">{{ points }}</text>
      <text style="font-size: 12px; opacity: 0.6; margin-top: 8px; display: block;">100积分 = ¥1.00，下单时可抵扣</text>
    </view>

    <!-- Points log -->
    <view style="margin: 0 16px;">
      <text style="font-size: 15px; font-weight: 700; color: #111; display: block; margin-bottom: 12px;">积分明细</text>

      <view v-for="log in logs" :key="log.id"
        style="background: #fff; border-radius: 12px; padding: 14px 16px; margin-bottom: 8px; display: flex; align-items: center; justify-content: space-between; box-shadow: 0 1px 4px rgba(0,0,0,0.03);">
        <view>
          <text style="font-size: 14px; color: #333; display: block;">{{ log.remark }}</text>
          <text style="font-size: 12px; color: #999; margin-top: 4px; display: block;">{{ log.created_at?.slice(0,10) }}</text>
        </view>
        <text :style="{ fontSize: '16px', fontWeight: '700', color: log.points > 0 ? '#22c55e' : '#dc2626' }">
          {{ log.points > 0 ? '+' : '' }}{{ log.points }}
        </text>
      </view>

      <view v-if="!logs.length" style="text-align: center; padding: 60px 0; color: #999; font-size: 14px;">
        暂无积分记录
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'

const points = ref(0)
const logs = ref<any[]>([])

onMounted(async () => {
  const profile = await get<any>('/api/v1/user/profile')
  points.value = profile?.points || 0

  const data = await get<any>('/api/v1/user/points/logs')
  logs.value = data?.list || []
})
</script>
